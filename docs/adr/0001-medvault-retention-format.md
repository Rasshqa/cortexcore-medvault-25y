# ADR 0001: MedVault Container Architecture & 25-Year Retention Standard

## Status
Accepted

## Konteks
Rumah sakit di Indonesia menghadapi ledakan ukuran database SIMRS setelah kewajiban RME 2022. Mempertahankan data kunjungan 25 tahun di server database aktif (PostgreSQL/MySQL) membebani performa query transaksi rawat jalan harian dan memerlukan biaya penyimpanan SSD enterprise yang sangat mahal.

## Keputusan
1. **Format Kontainer `.medvault`**:
   - Struktur berkas berformat binary container berbasis blok.
   - Header memuat metadata faskes, timestamp pembentukan, dan batas waktu retensi (25 tahun).
   - Payload data medis terkompresi dengan rasio tinggi (>80%) menggunakan algoritma kompresi berkecepatan tinggi.
   - Dilindungi enkripsi **AES-256-GCM** dengan derivasi kunci berbasis PBKDF2/SHA-256 dari frasa sandi master faskes.
2. **Integritas & Anti-Tampering**:
   - Setiap berkas `.medvault` memuat SHA-256 digital digest dari data sebelum dan sesudah enkripsi.
   - Modifikasi 1 bit pada berkas akan terdeteksi seketika saat pemeriksaan integritas.
3. **Pencarian Sub-Detik Luring**:
   - Pencarian riwayat kunjungan berdasarkan No. RM, NIK, atau rentang tanggal dapat dieksekusi langsung dari kontainer tanpa memerlukan instalasi database eksternal.

## Konsekuensi
- **Positif:** Menghemat 80-90% kebutuhan storage server RS, zero license fee, compliant dengan PMK 24/2022 dan UU PDP, dapat dibuka kembali 25 tahun kemudian tanpa ketergantungan versi vendor database tertentu.
- **Negatif:** Arsip bersifat append-only / read-only; perubahan data rekam medis lama memerlukan pembuatan berkas addendum baru.
