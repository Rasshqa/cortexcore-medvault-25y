# GLOSSARY: MedVault 25Y Archive (Klaster 01 - Health)

## Terminologi Regulasi & Teknis
- **PMK No. 24 Tahun 2022 Pasal 21**: Regulasi Kemenkes RI yang mewajibkan seluruh fasilitas pelayanan kesehatan menyimpan Rekam Medis Elektronik (RME) minimal 25 tahun sejak kunjungan terakhir pasien.
- **UU No. 27 Tahun 2022 (UU PDP)**: Kewajiban enkripsi data pribadi kesehatan (kategori data sensitif/spesifik) saat disimpan dalam media penyimpanan arsip (data at rest).
- **Nomor Rekam Medis (No. RM)**: Nomor identitas rekam medis internal rumah sakit yang unik per pasien.
- **Data at Rest Encryption**: Perlindungan kriptografi AES-256-GCM terhadap berkas arsip agar tidak dapat dibaca jika media fisik (hard disk / server) diambil pihak tidak berwenang.
- **Integritas Kriptografis (SHA-256 Checksum)**: Bukti matematis bahwa arsip medis tidak mengalami modifikasi, kerusakan bit (bit-rot), atau pemalsuan sejak hari pertama pengarsipan untuk kebutuhan audit akreditasi KARS / Kemenkes.
- **Columnar Archive**: Struktur penyimpanan data berorientasi kolom yang memungkinkan kompresi ekstrem (hingga 90%) dan pemindaian cepat tanpa memuat seluruh baris ke RAM.
