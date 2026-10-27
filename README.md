# CortexCore :: MedVault 25Y Archive

[![Build & Test](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Regulation](https://img.shields.io/badge/PMK_24%2F2022-Pasal_21-blue.svg)]()
[![Compliance](https://img.shields.io/badge/UU_PDP-AES--256--GCM-059669.svg)]()
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Zero Cloud Cost](https://img.shields.io/badge/Cloud_Storage_Cost-Rp0_(Local)-green.svg)]()

> **25-Year High-Compression Encrypted Medical Record Archiver & Sub-10ms Cold Storage Search for Indonesian Hospitals & SIMRS.**

Part of **AetherCortex Solutions** · Cluster 01: Health & SATUSEHAT  
Architected by **Rashqa Andrean Fitrah Sulaeman**

---

## Latar Belakang & Masalah Rumah Sakit

Berdasarkan **PMK No. 24 Tahun 2022 Pasal 21**, seluruh fasilitas pelayanan kesehatan (Rumah Sakit, Puskesmas, Klinik) di Indonesia **diwajibkan menyimpan berkas Rekam Medis Elektronik (RME) minimal 25 tahun** terhitung sejak kunjungan terakhir pasien.

### Dilema Rumah Sakit:
1. **Biaya Storage Server Membengkak**: Database operasional SIMRS (PostgreSQL/MySQL/SQL Server) menampung jutaan riwayat kunjungan masa lalu, menuntut kapasitas SAN/SSD enterprise bernilai puluhan juta rupiah.
2. **Kinerja SIMRS Lambat**: Query data pasien aktif melambat karena indeks database penuh dengan riwayat kunjungan 10–20 tahun yang lalu.
3. **Risiko Kepatuhan UU No. 27/2022 (UU PDP)**: Data medis sensitif pasien yang disimpan dalam bentuk backup SQL polos rentan terhadap pencurian data.
4. **Pembuktian Integritas Hukum Akreditasi KARS / Kemenkes**: Rumah sakit kesulitan membuktikan secara matematis bahwa rekam medis masa lalu tidak dimodifikasi atau dipalsukan.

---

## Solusi: MedVault 25Y Archive

Aplikasi desktop native (Wails v2 + Go + Svelte 5) yang **100% offline-first**:
1. **Kompresi Data Medis Tinggi (>80%)**:
   - Memampatkan berkas dump CSV/JSON SIMRS ke dalam kontainer binary `.medvault`.
   - Mengurangi kebutuhan ruang penyimpanan server hingga 90%.
2. **Enkripsi Kuat Standar UU PDP (AES-256-GCM)**:
   - Dilindungi enkripsi kelas militer AES-256-GCM dengan derivasi kunci PBKDF2 (100.000 iterasi SHA-256).
   - Zero-Knowledge: Frasa sandi master hanya diketahui oleh pihak rumah sakit.
3. **Integritas Anti-Tampering (SHA-256 Digest)**:
   - Menyimpan checksum matematis data rekam medis saat diarsipkan.
   - Deteksi seketika jika file arsip mengalami kerusakan bit atau modifikasi tidak sah saat diaudit akreditasi.
4. **Pencarian Sub-Detik (Sub-10ms)**:
   - Pencarian riwayat kunjungan pasien berdasarkan No. RM, NIK, atau nama langsung dari kontainer tanpa memerlukan database aktif.

---

## Menjalankan & Membangun

### Prasyarat
- Go 1.22+
- Node.js 18+ & npm
- Wails v2 CLI

### Build Binary
```bash
# Clone repository
git clone https://github.com/Rasshqa/cortexcore-medvault-25y.git
cd cortexcore-medvault-25y

# Install dependencies frontend
cd frontend && npm install && cd ..

# Jalankan unit test
go test -v ./...

# Build binary desktop
wails build -tags webkit2_41
```

Binary berada di `build/bin/medvault-25y`.

---

## Unit Testing

```bash
go test -v ./...
```
Output:
```text
=== RUN   TestVaultRoundtripAndTamperCheck
--- PASS: TestVaultRoundtripAndTamperCheck (0.15s)
PASS
ok      medvault-25y/vault      0.15s
```

---

## Hak Cipta & Lisensi

- Arsitektur & Rekayasa Sistem oleh **Rashqa Andrean Fitrah Sulaeman**.
- Lisensi Kode Terbuka: **MIT License**.
- Bagian dari ekosistem solusi 150 proyek AetherCortex.
