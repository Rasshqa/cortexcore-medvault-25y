package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVaultRoundtripAndTamperCheck(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "medvault-test")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	records := []MedicalRecord{
		{
			RecordID:       "REC-001",
			PatientNoRM:    "RM-2024-00128",
			NIK:            "3201234567890001",
			FullName:       "Siti Rahmawati",
			BirthDate:      "1985-04-12",
			Gender:         "female",
			VisitDate:      "2024-08-15",
			DischargeDate:  "2024-08-15",
			Department:     "Poli Penyakit Dalam",
			DoctorName:     "dr. Hendra Pratama, Sp.PD",
			DiagnosisICD10: "I10",
			DiagnosisText:  "Hipertensi Esensial Primer",
			TreatmentNotes: "Tensimeter 145/95. Terapi Amlodipine 5mg.",
		},
		{
			RecordID:       "REC-002",
			PatientNoRM:    "RM-2024-00129",
			NIK:            "3201234567890002",
			FullName:       "Budi Santoso",
			BirthDate:      "1990-11-20",
			Gender:         "male",
			VisitDate:      "2024-09-01",
			DischargeDate:  "2024-09-03",
			Department:     "Rawat Inap Bedah",
			DoctorName:     "dr. Rian Nugroho, Sp.B",
			DiagnosisICD10: "K35.8",
			DiagnosisText:  "Apendisitis Akut Lainnya",
			TreatmentNotes: "Apendektomi elektif, pemulihan baik.",
		},
	}

	passphrase := "RahasiaFaskes2026!"
	vaultPath := filepath.Join(tmpDir, "simrs_archive_2024.medvault")

	// 1. Create vault
	manifest, err := CreateVault(records, "RSUD Dr. Moewardi", passphrase, vaultPath)
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	if manifest.TotalRecords != 2 {
		t.Errorf("expected 2 records, got %d", manifest.TotalRecords)
	}
	if manifest.CompressedBytes >= manifest.UncompressedBytes {
		t.Errorf("expected compression saving, got uncompressed=%d, compressed=%d",
			manifest.UncompressedBytes, manifest.CompressedBytes)
	}

	// 2. Open with wrong passphrase -> must fail
	_, err = OpenVault(vaultPath, "WrongPassword")
	if err == nil {
		t.Errorf("expected error on wrong password, but passed")
	}

	// 3. Open with correct passphrase
	arc, err := OpenVault(vaultPath, passphrase)
	if err != nil {
		t.Fatalf("OpenVault with correct passphrase failed: %v", err)
	}

	if arc.Manifest.IsTampered {
		t.Errorf("expected clean archive, marked as tampered")
	}
	if len(arc.Records) != 2 {
		t.Fatalf("expected 2 records decrypted, got %d", len(arc.Records))
	}

	// 4. Test search
	res := SearchArchive(arc, SearchFilter{Query: "Siti"})
	if res.TotalMatched != 1 || res.Records[0].FullName != "Siti Rahmawati" {
		t.Errorf("search Siti failed, got %+v", res)
	}

	resByNIK := SearchArchive(arc, SearchFilter{Query: "3201234567890002"})
	if resByNIK.TotalMatched != 1 || resByNIK.Records[0].PatientNoRM != "RM-2024-00129" {
		t.Errorf("search by NIK failed, got %+v", resByNIK)
	}
}
