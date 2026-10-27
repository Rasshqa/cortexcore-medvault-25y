package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"medvault-25y/vault"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	openArc   *vault.Archive
	allImported []vault.MedicalRecord
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// ImportCSV parses a CSV dump from SIMRS and returns preview + total count.
func (a *App) ImportCSV(csvContent string) map[string]interface{} {
	records, headers, err := parseCSVRecords(csvContent)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	a.allImported = records

	preview := records
	if len(preview) > 10 {
		preview = records[:10]
	}

	return map[string]interface{}{
		"total":   len(records),
		"preview": preview,
		"columns": headers,
	}
}

// OpenCSVFile shows file picker and returns CSV content.
func (a *App) OpenCSVFile() map[string]interface{} {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Pilih Berkas CSV Ekspor SIMRS",
		Filters: []runtime.FileFilter{{DisplayName: "CSV Files", Pattern: "*.csv"}},
	})
	if err != nil || path == "" {
		return map[string]interface{}{"error": "file tidak dipilih"}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": "gagal membaca file"}
	}

	return map[string]interface{}{"path": path, "content": string(raw)}
}

// CreateVaultFromImported creates a .medvault from previously imported CSV records.
func (a *App) CreateVaultFromImported(orgName, passphrase string) map[string]interface{} {
	if len(a.allImported) == 0 {
		return map[string]interface{}{"error": "Belum ada data CSV yang diimpor"}
	}

	destPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Simpan Arsip MedVault (.medvault)",
		DefaultFilename: strings.ReplaceAll(orgName, " ", "_") + "_archive.medvault",
		Filters:         []runtime.FileFilter{{DisplayName: "MedVault Archive", Pattern: "*.medvault"}},
	})
	if err != nil || destPath == "" {
		return map[string]interface{}{"error": "simpan dibatalkan"}
	}

	manifest, err := vault.CreateVault(a.allImported, orgName, passphrase, destPath)
	if err != nil {
		return map[string]interface{}{"error": "gagal membuat vault: " + err.Error()}
	}

	return map[string]interface{}{
		"saved":    destPath,
		"manifest": manifest,
	}
}

// OpenMedvaultFile opens a .medvault archive and stores it in memory.
func (a *App) OpenMedvaultFile(passphrase string) map[string]interface{} {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Buka Arsip MedVault",
		Filters: []runtime.FileFilter{{DisplayName: "MedVault Archive", Pattern: "*.medvault"}},
	})
	if err != nil || path == "" {
		return map[string]interface{}{"error": "file tidak dipilih"}
	}

	arc, err := vault.OpenVault(path, passphrase)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	a.openArc = arc
	return map[string]interface{}{
		"manifest":    arc.Manifest,
		"is_tampered": arc.Manifest.IsTampered,
	}
}

// SearchVault searches the currently opened vault.
func (a *App) SearchVault(query, dateFrom, dateTo string, limit int) map[string]interface{} {
	if a.openArc == nil {
		return map[string]interface{}{"error": "Belum ada arsip yang dibuka"}
	}

	result := vault.SearchArchive(a.openArc, vault.SearchFilter{
		Query:    query,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
	})

	return map[string]interface{}{
		"total_matched": result.TotalMatched,
		"duration_ms":   result.DurationMs,
		"records":       result.Records,
	}
}

// GetVaultManifest returns manifest of currently opened vault.
func (a *App) GetVaultManifest() map[string]interface{} {
	if a.openArc == nil {
		return map[string]interface{}{"error": "Belum ada arsip yang dibuka"}
	}
	b, _ := json.Marshal(a.openArc.Manifest)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}

// ── CSV parser helper ──

func parseCSVRecords(csvContent string) ([]vault.MedicalRecord, []string, error) {
	r := csv.NewReader(strings.NewReader(csvContent))
	rows, err := r.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("CSV tidak valid: %w", err)
	}
	if len(rows) < 2 {
		return nil, nil, fmt.Errorf("CSV kosong atau hanya header")
	}

	headers := rows[0]
	idx := func(names ...string) int {
		for _, name := range names {
			for i, h := range headers {
				if strings.EqualFold(strings.TrimSpace(h), name) {
					return i
				}
			}
		}
		return -1
	}

	get := func(row []string, i int) string {
		if i < 0 || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}

	iRecordID   := idx("record_id", "id", "no_kunjungan", "kunjungan_id")
	iNoRM       := idx("patient_no_rm", "no_rm", "nomor_rm", "norm")
	iNIK        := idx("nik", "no_ktp", "ktp")
	iFullName   := idx("full_name", "nama", "nama_pasien", "patient_name")
	iBirthDate  := idx("birth_date", "tgl_lahir", "tanggal_lahir")
	iGender     := idx("gender", "jenis_kelamin", "sex")
	iVisitDate  := idx("visit_date", "tgl_kunjungan", "tanggal_kunjungan", "tgl_masuk")
	iDischarge  := idx("discharge_date", "tgl_keluar", "tanggal_keluar")
	iDepartment := idx("department", "poli", "ruang", "unit")
	iDoctor     := idx("doctor_name", "dokter", "nama_dokter")
	iDiagICD10  := idx("diagnosis_icd10", "icd10", "kode_icd", "kode_diagnosa")
	iDiagText   := idx("diagnosis_text", "diagnosa", "diagnosis", "nama_diagnosa")
	iProcICD9   := idx("procedure_icd9", "icd9", "kode_tindakan")
	iNotes      := idx("treatment_notes", "catatan", "soap", "anamnesis")

	var records []vault.MedicalRecord
	for _, row := range rows[1:] {
		rec := vault.MedicalRecord{
			RecordID:       fmt.Sprintf("%s-%s", get(row, iNoRM), get(row, iVisitDate)),
			PatientNoRM:    get(row, iNoRM),
			NIK:            get(row, iNIK),
			FullName:       get(row, iFullName),
			BirthDate:      get(row, iBirthDate),
			Gender:         get(row, iGender),
			VisitDate:      get(row, iVisitDate),
			DischargeDate:  get(row, iDischarge),
			Department:     get(row, iDepartment),
			DoctorName:     get(row, iDoctor),
			DiagnosisICD10: get(row, iDiagICD10),
			DiagnosisText:  get(row, iDiagText),
			ProcedureICD9:  get(row, iProcICD9),
			TreatmentNotes: get(row, iNotes),
		}
		if iRecordID >= 0 {
			rec.RecordID = get(row, iRecordID)
		}
		records = append(records, rec)
	}

	return records, headers, nil
}
