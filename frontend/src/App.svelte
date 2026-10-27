<script>
  import { onMount } from 'svelte';
  import {
    ImportCSV,
    OpenCSVFile,
    CreateVaultFromImported,
    OpenMedvaultFile,
    SearchVault,
    GetVaultManifest
  } from '../wailsjs/go/main/App.js';

  let activeTab = $state("import"); // import | archive | search
  let csvContent = $state("");
  let importResult = $state(null);
  let orgName = $state("RSUD Dr. Moewardi Surakarta");
  let passphrase = $state("");
  let manifest = $state(null);
  let isTampered = $state(false);
  let searchQuery = $state("");
  let searchResults = $state(null);
  let createStatus = $state("");
  let isProcessing = $state(false);

  async function handleOpenCSV() {
    const res = await OpenCSVFile();
    if (res.content) {
      csvContent = res.content;
      const parsed = await ImportCSV(csvContent);
      importResult = parsed;
    }
  }

  async function handleCreateVault() {
    if (!passphrase || passphrase.length < 8) {
      createStatus = "Frasa sandi minimal 8 karakter!";
      return;
    }
    isProcessing = true;
    createStatus = "Mengompresi & mengenkripsi arsip...";
    const res = await CreateVaultFromImported(orgName, passphrase);
    isProcessing = false;
    if (res.error) {
      createStatus = "Gagal: " + res.error;
    } else {
      manifest = res.manifest;
      createStatus = "Arsip tersimpan: " + res.saved;
      activeTab = "archive";
    }
  }

  async function handleOpenVault() {
    if (!passphrase) {
      createStatus = "Masukkan frasa sandi dekripsi!";
      return;
    }
    isProcessing = true;
    const res = await OpenMedvaultFile(passphrase);
    isProcessing = false;
    if (res.error) {
      createStatus = "Gagal: " + res.error;
    } else {
      manifest = res.manifest;
      isTampered = res.is_tampered;
      activeTab = "search";
    }
  }

  async function handleSearch() {
    const res = await SearchVault(searchQuery, "", "", 50);
    if (res.error) {
      searchResults = { error: res.error };
    } else {
      searchResults = res;
    }
  }
</script>

<div class="app-shell">
  <!-- Topbar -->
  <header class="topbar">
    <div class="topbar-brand">
      <span class="brand-dot"></span>
      <span class="brand-name">CortexCore :: MedVault 25Y Archive</span>
      <span class="brand-sub">PMK 24/2022 · UU PDP · Retensi 25 Tahun</span>
    </div>
    <div style="display: flex; gap: 8px;">
      <button class="btn btn-ghost" onclick={handleOpenCSV}>Impor CSV SIMRS</button>
      <button class="btn btn-ghost" onclick={handleOpenVault}>Buka .medvault</button>
    </div>
  </header>

  <!-- Metrics Strip -->
  {#if manifest}
    <div class="metrics-strip">
      <div class="metric-card">
        <div class="metric-label">Total Rekam Medis</div>
        <div class="metric-val">{manifest.total_records?.toLocaleString() ?? '—'}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Kompresi</div>
        <div class="metric-val emerald">{manifest.compression_ratio?.toFixed(1) ?? '—'}%</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Kadaluarsa Retensi</div>
        <div class="metric-val">{manifest.retention_expires_at?.slice(0,10) ?? '—'}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Integritas SHA-256</div>
        <div class="metric-val" style="font-size: 10px; word-break: break-all;">
          {#if isTampered}
            <span style="color: var(--mv-crimson);">DIPALSUKAN</span>
          {:else}
            <span style="color: var(--mv-emerald-lt);">TERVERIFIKASI</span>
          {/if}
        </div>
      </div>
    </div>
  {/if}

  <!-- Tabs -->
  <nav class="nav-tabs">
    <button class="tab-btn {activeTab === 'import' ? 'active' : ''}" onclick={() => activeTab = 'import'}>
      Impor CSV
    </button>
    <button class="tab-btn {activeTab === 'archive' ? 'active' : ''}" onclick={() => activeTab = 'archive'}>
      Arsip Vault
    </button>
    <button class="tab-btn {activeTab === 'search' ? 'active' : ''}" onclick={() => activeTab = 'search'}>
      Cari Rekam Medis
    </button>
  </nav>

  <!-- Workspace -->
  <div class="workspace">
    <div class="pane-container">
      {#if activeTab === 'import'}
        <!-- Import CSV Tab -->
        <div class="input-row">
          <input class="text-input" bind:value={orgName} placeholder="Nama Faskes / Rumah Sakit" />
        </div>
        <div class="input-row">
          <input class="text-input" type="password" bind:value={passphrase} placeholder="Frasa Sandi Enkripsi (min. 8 karakter)" />
          <button class="btn btn-cyan" onclick={handleCreateVault} disabled={isProcessing || !importResult}>
            {#if isProcessing}Mengompresi...{:else}Buat Arsip .medvault{/if}
          </button>
        </div>

        {#if createStatus}
          <div style="padding: 8px 0; font-size: 11px; color: var(--mv-amber);">{createStatus}</div>
        {/if}

        {#if importResult && importResult.total}
          <div style="font-size: 12px; margin-bottom: 8px; color: var(--mv-cyan-lt);">
            {importResult.total} rekam medis terdeteksi dari CSV ({importResult.columns?.length ?? 0} kolom)
          </div>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>No RM</th>
                  <th>NIK</th>
                  <th>Nama Pasien</th>
                  <th>Tgl Kunjungan</th>
                  <th>Poli / Unit</th>
                  <th>ICD-10</th>
                  <th>Diagnosis</th>
                </tr>
              </thead>
              <tbody>
                {#each importResult.preview as rec}
                  <tr>
                    <td>{rec.patient_no_rm}</td>
                    <td>{rec.nik}</td>
                    <td>{rec.full_name}</td>
                    <td>{rec.visit_date}</td>
                    <td>{rec.department}</td>
                    <td>{rec.diagnosis_icd10}</td>
                    <td>{rec.diagnosis_text}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if !importResult}
          <div style="text-align: center; color: var(--mv-muted); padding: 48px;">
            <div style="font-size: 28px; opacity: 0.3; margin-bottom: 8px;">📦</div>
            <div style="font-weight: 600; margin-bottom: 4px;">Belum ada data</div>
            <div style="font-size: 11px;">Klik "Impor CSV SIMRS" untuk memuat dump rekam medis dari SIMRS Anda.</div>
          </div>
        {/if}

      {:else if activeTab === 'archive'}
        <!-- Archive Manifest Tab -->
        {#if manifest}
          <div style="font-size: 12px; margin-bottom: 12px;">
            <strong>Arsip:</strong> {manifest.organization_name} ·
            <strong>ID:</strong> <span style="font-family: monospace; color: var(--mv-cyan-lt);">{manifest.vault_id}</span>
          </div>
          <div class="table-container">
            <table class="data-table">
              <thead><tr><th>Properti</th><th>Nilai</th></tr></thead>
              <tbody>
                <tr><td>Total Rekam Medis</td><td>{manifest.total_records}</td></tr>
                <tr><td>Ukuran Sebelum Kompresi</td><td>{(manifest.uncompressed_bytes / 1024).toFixed(1)} KB</td></tr>
                <tr><td>Ukuran Sesudah Kompresi</td><td>{(manifest.compressed_bytes / 1024).toFixed(1)} KB</td></tr>
                <tr><td>Rasio Kompresi</td><td style="color: var(--mv-emerald-lt);">{manifest.compression_ratio?.toFixed(1)}%</td></tr>
                <tr><td>SHA-256 Integritas Data</td><td style="font-family: monospace; font-size: 10px; word-break: break-all;">{manifest.integrity_sha256}</td></tr>
                <tr><td>Tanggal Arsip</td><td>{manifest.archived_at?.slice(0,10)}</td></tr>
                <tr><td>Kadaluarsa Retensi (25 th)</td><td style="color: var(--mv-amber);">{manifest.retention_expires_at?.slice(0,10)}</td></tr>
              </tbody>
            </table>
          </div>
        {:else}
          <div style="text-align: center; color: var(--mv-muted); padding: 48px;">
            Buat arsip dari tab Impor CSV atau buka file .medvault yang sudah ada.
          </div>
        {/if}

      {:else if activeTab === 'search'}
        <!-- Search Tab -->
        <div class="input-row">
          <input class="text-input" bind:value={searchQuery} placeholder="Cari No RM, NIK, atau Nama Pasien..." />
          <button class="btn btn-cyan" onclick={handleSearch}>Cari</button>
        </div>

        {#if searchResults}
          {#if searchResults.error}
            <div style="color: var(--mv-crimson); font-size: 12px;">{searchResults.error}</div>
          {:else}
            <div style="font-size: 11px; color: var(--mv-muted); margin-bottom: 8px;">
              {searchResults.total_matched} hasil ditemukan dalam {searchResults.duration_ms}ms
            </div>
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>No RM</th>
                    <th>NIK</th>
                    <th>Nama</th>
                    <th>Tgl Kunjungan</th>
                    <th>Poli</th>
                    <th>ICD-10</th>
                    <th>Diagnosis</th>
                    <th>Dokter</th>
                  </tr>
                </thead>
                <tbody>
                  {#each searchResults.records as rec}
                    <tr>
                      <td>{rec.patient_no_rm}</td>
                      <td>{rec.nik}</td>
                      <td>{rec.full_name}</td>
                      <td>{rec.visit_date}</td>
                      <td>{rec.department}</td>
                      <td>{rec.diagnosis_icd10}</td>
                      <td>{rec.diagnosis_text}</td>
                      <td>{rec.doctor_name}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        {:else}
          <div style="text-align: center; color: var(--mv-muted); padding: 48px;">
            <div style="font-size: 28px; opacity: 0.3; margin-bottom: 8px;">🔍</div>
            <div style="font-weight: 600; margin-bottom: 4px;">Pencarian Arsip 25 Tahun</div>
            <div style="font-size: 11px;">Buka file .medvault lalu cari berdasarkan No. RM, NIK, atau nama pasien.</div>
          </div>
        {/if}
      {/if}
    </div>
  </div>

  <!-- Attribution Footer -->
  <footer class="footer-attribution">
    <div>
      <span class="footer-brand">AetherCortex</span> Solution · Cluster 01: Health & SATUSEHAT
    </div>
    <div>
      Architected by <strong>Rashqa Andrean Fitrah Sulaeman</strong> · PMK 24/2022 Compliant · UU PDP Encrypted
    </div>
  </footer>
</div>
