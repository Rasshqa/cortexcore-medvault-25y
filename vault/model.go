package vault

import (
	"time"
)

// MedicalRecord represents a standardized clinical encounter for 25-year archiving.
type MedicalRecord struct {
	RecordID        string `json:"record_id"`         // UUID or internal SIMRS ID
	PatientNoRM     string `json:"patient_no_rm"`     // Nomor Rekam Medis (e.g. "RM-2024-00128")
	NIK             string `json:"nik"`               // 16-digit Indonesian NIK
	FullName        string `json:"full_name"`         // Patient full name
	BirthDate       string `json:"birth_date"`        // YYYY-MM-DD
	Gender          string `json:"gender"`            // male | female
	VisitDate       string `json:"visit_date"`        // YYYY-MM-DD
	DischargeDate   string `json:"discharge_date"`    // YYYY-MM-DD
	Department      string `json:"department"`        // Poli Penyakit Dalam, IGD, etc.
	DoctorName      string `json:"doctor_name"`       // Attending physician
	DiagnosisICD10  string `json:"diagnosis_icd10"`   // e.g. "I10", "E11.9"
	DiagnosisText   string `json:"diagnosis_text"`    // e.g. "Hipertensi Esensial"
	ProcedureICD9   string `json:"procedure_icd9"`    // e.g. "89.52"
	TreatmentNotes  string `json:"treatment_notes"`   // Summary SOAP
}

// VaultManifest contains cryptographic and regulatory metadata of the archive.
type VaultManifest struct {
	VaultID            string    `json:"vault_id"`
	OrganizationName   string    `json:"organization_name"`
	ArchivedAt         time.Time `json:"archived_at"`
	RetentionExpiresAt time.Time `json:"retention_expires_at"` // ArchivedAt + 25 Years (PMK 24/2022)
	TotalRecords       int       `json:"total_records"`
	UncompressedBytes  int64     `json:"uncompressed_bytes"`
	CompressedBytes    int64     `json:"compressed_bytes"`
	CompressionRatio   float64   `json:"compression_ratio"`    // Percentage saved
	IntegritySHA256    string    `json:"integrity_sha256"`     // Checksum of plaintext data
	EncryptedChecksum  string    `json:"encrypted_checksum"`   // Checksum of final file
	IsTampered         bool      `json:"is_tampered"`
}

// SearchFilter specifies search criteria across the archive.
type SearchFilter struct {
	Query     string `json:"query"`      // Matches NoRM, NIK, or Patient Name
	DateFrom  string `json:"date_from"`  // Filter visit dates
	DateTo    string `json:"date_to"`
	Limit     int    `json:"limit"`
}

// SearchResult returns matched records and performance metrics.
type SearchResult struct {
	TotalMatched int             `json:"total_matched"`
	DurationMs   int64           `json:"duration_ms"`
	Records      []MedicalRecord `json:"records"`
}
