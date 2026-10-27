package vault

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// medvaultMagic is the binary file signature to identify .medvault containers.
var medvaultMagic = [8]byte{'M', 'V', 'L', 'T', '2', '5', 'Y', 0x01}

// payloadContainer is the unencrypted structure before compression.
type payloadContainer struct {
	Records []MedicalRecord `json:"records"`
}

// Archive is the in-memory structure representing an opened vault.
type Archive struct {
	Manifest VaultManifest   `json:"manifest"`
	Records  []MedicalRecord `json:"records"`
}

// CreateVault compresses and encrypts a slice of medical records into a .medvault file.
func CreateVault(records []MedicalRecord, orgName, passphrase, destPath string) (*VaultManifest, error) {
	manifest := VaultManifest{
		VaultID:            uuid.NewString(),
		OrganizationName:   orgName,
		ArchivedAt:         time.Now().UTC(),
		RetentionExpiresAt: time.Now().UTC().AddDate(25, 0, 0),
		TotalRecords:       len(records),
	}

	// 1. Serialize records data first
	recordsJSON, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("marshal records: %w", err)
	}
	manifest.UncompressedBytes = int64(len(recordsJSON))
	manifest.IntegritySHA256 = SHA256Hex(recordsJSON)

	// 2. Wrap manifest + records into full payload
	container := Archive{
		Manifest: manifest,
		Records:  records,
	}
	payloadJSON, err := json.Marshal(container)
	if err != nil {
		return nil, fmt.Errorf("marshal container: %w", err)
	}

	// 3. Compress
	var cbuf bytes.Buffer
	w := zlib.NewWriter(&cbuf)
	if _, err := w.Write(payloadJSON); err != nil {
		return nil, fmt.Errorf("compress: %w", err)
	}
	w.Close()
	compressed := cbuf.Bytes()
	manifest.CompressedBytes = int64(len(compressed))
	manifest.CompressionRatio = (1.0 - float64(len(compressed))/float64(len(payloadJSON))) * 100.0

	// 4. Derive key + encrypt
	salt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}
	key := DeriveKey(passphrase, salt)
	ciphertext, nonce, err := Encrypt(compressed, key)
	if err != nil {
		return nil, err
	}
	manifest.EncryptedChecksum = SHA256Hex(ciphertext)

	// 5. Write binary container
	f, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	write32 := func(n uint32) { binary.Write(f, binary.LittleEndian, n) }
	write64 := func(n uint64) { binary.Write(f, binary.LittleEndian, n) }

	f.Write(medvaultMagic[:])
	write32(uint32(len(salt)))
	f.Write(salt)
	write32(uint32(len(nonce)))
	f.Write(nonce)
	write64(uint64(len(ciphertext)))
	f.Write(ciphertext)

	return &manifest, nil
}

// OpenVault decrypts and loads a .medvault archive into memory.
func OpenVault(path, passphrase string) (*Archive, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baca file: %w", err)
	}

	r := bytes.NewReader(data)

	// Validate magic
	var magic [8]byte
	if _, err := io.ReadFull(r, magic[:]); err != nil || magic != medvaultMagic {
		return nil, errors.New("bukan berkas arsip MedVault yang valid (.medvault)")
	}

	read32 := func() (uint32, error) {
		var n uint32
		return n, binary.Read(r, binary.LittleEndian, &n)
	}
	read64 := func() (uint64, error) {
		var n uint64
		return n, binary.Read(r, binary.LittleEndian, &n)
	}

	saltLen, _ := read32()
	salt := make([]byte, saltLen)
	io.ReadFull(r, salt)

	nonceLen, _ := read32()
	nonce := make([]byte, nonceLen)
	io.ReadFull(r, nonce)

	payloadLen, _ := read64()
	ciphertext := make([]byte, payloadLen)
	io.ReadFull(r, ciphertext)

	key := DeriveKey(passphrase, salt)
	compressed, err := Decrypt(ciphertext, key, nonce)
	if err != nil {
		return nil, err
	}

	// Decompress
	zr, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, fmt.Errorf("dekompresi gagal: %w", err)
	}
	payloadJSON, err := io.ReadAll(zr)
	zr.Close()
	if err != nil {
		return nil, fmt.Errorf("dekompresi gagal: %w", err)
	}

	var arc Archive
	if err := json.Unmarshal(payloadJSON, &arc); err != nil {
		return nil, fmt.Errorf("parse arsip: %w", err)
	}

	// Verify data integrity against re-marshaled records
	recJSON, _ := json.Marshal(arc.Records)
	actualChecksum := SHA256Hex(recJSON)
	if actualChecksum != arc.Manifest.IntegritySHA256 {
		arc.Manifest.IsTampered = true
	}

	return &arc, nil
}

// SearchArchive performs an in-memory search against a decrypted archive.
func SearchArchive(arc *Archive, filter SearchFilter) *SearchResult {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	start := time.Now()
	q := strings.ToLower(strings.TrimSpace(filter.Query))
	var matched []MedicalRecord

	for _, rec := range arc.Records {
		if len(matched) >= filter.Limit {
			break
		}
		if q == "" ||
			strings.Contains(strings.ToLower(rec.PatientNoRM), q) ||
			strings.Contains(strings.ToLower(rec.NIK), q) ||
			strings.Contains(strings.ToLower(rec.FullName), q) {

			if filter.DateFrom != "" && rec.VisitDate < filter.DateFrom {
				continue
			}
			if filter.DateTo != "" && rec.VisitDate > filter.DateTo {
				continue
			}
			matched = append(matched, rec)
		}
	}

	return &SearchResult{
		TotalMatched: len(matched),
		DurationMs:   time.Since(start).Milliseconds(),
		Records:      matched,
	}
}
