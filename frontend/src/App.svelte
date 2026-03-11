<script>
  import { onMount } from 'svelte';
  import appIcon from './assets/images/app-icon.png';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import {
    StartScraping,
    PauseScraping,
    ResumeScraping,
    StopScraping,
    SelectDirectory,
    SaveAPICredentials,
    GetAPICredentials,
    SaveConfigValue,
    GetConfigValue,
    ApplyRenaming,
    CheckPreviousSession
  } from '../wailsjs/go/main/App';
  let activeTab = 'dashboard';
  let progress = { total: 0, processed: 0, current_file: '', proposed_name: '', status: 'idle' };
  let logs = [];
  let rootPath = '';
  let matchedFiles = [];
  let selectedAll = true;

  // Modal de sesión previa
  let showSessionModal = false;
  let sessionInfo = null; // { found, status, total_files, done_files }
  let isScraping = false;

  // Config fields
  let screenScraperUser = '';
  let screenScraperPass = '';
  let tgdbApiKey = '';
  let namingTemplate = '{Name} ({Region})';
  let workerCount = 4;
  let ssSearchByHash = true;
  let ssSearchByName = true;
  let ssActive = true;
  let tgdbSearchByHash = true;
  let tgdbSearchByName = true;
  let tgdbActive = false;
  
  let apiCounts = {};

  onMount(async () => {
    EventsOn('progress', (data) => {
      progress = data;
      if (data.api_counts) {
        for (let provider in data.api_counts) {
          apiCounts[provider] = { Count: data.api_counts[provider] };
        }
        apiCounts = { ...apiCounts }; // Trigger Svelte reactivity
      }
      if (data.message) {
        addLog(`${data.current_file}: ${data.message}`, data.status === 'error' ? 'error' : 'info');
      }
      if (data.proposed_name) {
        matchedFiles = [{ id: data.id, original: data.current_file, proposed: data.proposed_name, selected: true }, ...matchedFiles.slice(0, 49)];
      }
      if (data.status === 'error' ) {
          addLog(`Error: ${data.current_file}`, 'error');
      }
      if (data.status === 'completed') {
          addLog('Proceso finalizado con éxito', 'success');
          isScraping = false;
      }
      if (data.status === 'stopped') {
          addLog('Proceso detenido.', 'info');
          isScraping = false;
      }
    });

    try {
      const ssCreds = await GetAPICredentials('screenscraper');
      if (ssCreds) {
        screenScraperUser = ssCreds.Username;
        screenScraperPass = ssCreds.Password;
        ssActive = ssCreds.IsActive;
        ssSearchByHash = ssCreds.SearchByHash;
        ssSearchByName = ssCreds.SearchByName;
      }
      
      const tgdb = await GetAPICredentials('thegamesdb');
      if (tgdb) {
        tgdbApiKey = tgdb.APIKey;
        tgdbActive = tgdb.IsActive;
        tgdbSearchByHash = tgdb.SearchByHash;
        tgdbSearchByName = tgdb.SearchByName;
      }

      namingTemplate = await GetConfigValue('naming_template', '{Name} ({Region})');
      workerCount = parseInt(await GetConfigValue('worker_count', '4'));
      searchByHash = (await GetConfigValue('search_by_hash', 'true')) === 'true';
      searchByName = (await GetConfigValue('search_by_name', 'true')) === 'true';
      
      await refreshAPICounts();
      console.log("API Counts Initialized:", apiCounts);

    } catch (e) {
      console.log("Initialization: Some defaults used");
    }
  });

  function addLog(msg, type = 'info') {
    logs = [{ msg, type, time: new Date().toLocaleTimeString() }, ...logs.slice(0, 99)];
  }

  async function handleSelectDir() {
    const result = await SelectDirectory();
    if (result) {
      rootPath = result;
      addLog(`Directorio seleccionado: ${rootPath}`, 'success');
    }
  }

  async function start() {
    if (!rootPath) {
      addLog('Error: Selecciona un directorio primero', 'error');
      return;
    }
    // Verificar si hay sesión previa para esta ruta
    try {
      const prev = await CheckPreviousSession(rootPath);
      if (prev && prev.found) {
        sessionInfo = prev;
        showSessionModal = true;
        return; // Esperar decision del usuario en el modal
      }
    } catch (_) {}
    await launchScraping(false);
  }

  async function launchScraping(forceRestart) {
    if (isScraping) return;
    isScraping = true;
    showSessionModal = false;
    matchedFiles = []; // SIEMPRE limpiar la tabla antes de un nuevo escaneo
    try {
      const msg = forceRestart
        ? 'Iniciando sesión nueva desde cero...'
        : (sessionInfo?.status === 'completed' ? 'Iniciando nuevo escaneo...' : 'Reanudando sesión anterior...');
      addLog(msg, 'info');
      await StartScraping(rootPath, forceRestart);
    } catch (e) {
      addLog(`Error fatal: ${e}`, 'error');
      isScraping = false;
    }
  }

  async function saveConfig() {
    try {
      await SaveAPICredentials('screenscraper', screenScraperUser, screenScraperPass, '', '', ssSearchByHash, ssSearchByName, ssActive);
      await SaveAPICredentials('thegamesdb', '', '', tgdbApiKey, '', tgdbSearchByHash, tgdbSearchByName, tgdbActive);
      
      await SaveConfigValue('naming_template', namingTemplate);
      await SaveConfigValue('worker_count', workerCount.toString());

      addLog('Configuración sincronizada correctamente', 'success');
      await refreshAPICounts();
    } catch (e) {
      addLog(`Error al sincronizar: ${e}`, 'error');
    }
  }

  async function refreshAPICounts() {
    const ssTracker = await GetConfigValue('api_tracker_screenscraper', '{}');
    try {
      if (ssTracker) apiCounts['screenscraper'] = JSON.parse(ssTracker);
    } catch(e){}
    
    const tgdbTracker = await GetConfigValue('api_tracker_thegamesdb', '{}');
    try {
      if (tgdbTracker) apiCounts['thegamesdb'] = JSON.parse(tgdbTracker);
    } catch(e){}
    
    apiCounts = { ...apiCounts };
  }

  function pause() { PauseScraping(); addLog('Proceso pausado'); }
  function resume() { ResumeScraping(); addLog('Proceso reanudado'); }
  function stop() { StopScraping(); addLog('Proceso detenido'); }

  function toggleAll() {
    matchedFiles = matchedFiles.map(f => ({ ...f, selected: selectedAll }));
  }

  let renameErrors = [];
  let showRenameErrorsModal = false;

  async function handleRenameSelected() {
    const selectedIds = matchedFiles.filter(f => f.selected).map(f => f.id);
    if (selectedIds.length === 0) {
      addLog('No hay archivos seleccionados para renombrar', 'warning');
      return;
    }
    
    try {
      const errors = await ApplyRenaming(selectedIds);
      
      matchedFiles = matchedFiles.map(f => {
        if (selectedIds.includes(f.id)) {
          // If in errors, leave it as is or mark error
          const isError = errors && errors.some(e => e.filename === f.original || e.filename === f.proposed);
          if (!isError) {
            return { ...f, status: 'renamed', selected: false };
          }
        }
        return f;
      });

      if (errors && errors.length > 0) {
        renameErrors = errors;
        showRenameErrorsModal = true;
        addLog(`Hubo ${errors.length} errores al renombrar`, 'error');
      } else {
        addLog(`Renombrado de ${selectedIds.length} archivos exitoso`, 'success');
      }

    } catch (e) {
      addLog(`Error al renombrar: ${e}`, 'error');
    }
  }
</script>

<div class="layout">
  <nav>
    <div class="logo-container">
      <div class="logo-box">
        <img src={appIcon} alt="Logo" style="width: 100%; height: 100%; object-fit: contain; border-radius: 8px;" />
      </div>
      <div>
        <h1 style="font-size: 1.1rem; font-weight: 700; letter-spacing: -0.02em;">NEXUS ROMS</h1>
        <div style="font-size: 0.65rem; color: var(--primary); font-weight: 800; letter-spacing: 0.1em; margin-top: -2px;">PREMIUM ENGINE</div>
      </div>
    </div>
    <div class="tab-container">
      <button class="tab-btn {activeTab === 'dashboard' ? 'active' : ''}" on:click={() => activeTab = 'dashboard'}>DASHBOARD</button>
      <button class="tab-btn {activeTab === 'settings' ? 'active' : ''}" on:click={() => activeTab = 'settings'}>SETTINGS</button>
    </div>
  </nav>

  <main class="content">
    {#if activeTab === 'dashboard'}
      <div class="dashboard-grid">
        <!-- Sidebar Controls -->
        <div style="display: flex; flex-direction: column; gap: 2rem;">
          <div class="card" style="{isScraping ? 'position: relative; z-index: 1001;' : ''}">
            <h2>Session control</h2>
            
            <div class="field-group">
              <span class="label">Target Directory</span>
              <div class="input-wrapper">
                <input type="text" readonly value={rootPath} placeholder="Choose a folder...">
                <button class="btn btn-primary" style="position: absolute; right: 4px; padding: 0.5rem; height: calc(100% - 8px);" on:click={handleSelectDir} disabled={isScraping}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"/></svg>
                </button>
              </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-top: 0.5rem;">
              {#if progress.status === 'idle' || progress.status === 'stopped' || progress.status === 'completed'}
                <button class="btn btn-primary" style="grid-column: span 2; padding: 1.25rem;" on:click={start}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.5-1 1-4c1.5 0 3 .5 3 .5L9 12z"/><path d="M15 6h5s.5 1 1 4c-1.5 0-3-.5-3-.5L15 6z"/></svg>
                  START ENGINE
                </button>
              {:else}
                {#if progress.status === 'paused'}
                   <button class="btn btn-ghost" style="color: var(--success);" on:click={resume}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    Resume
                   </button>
                {:else}
                   <button class="btn btn-ghost" style="color: var(--warning);" on:click={pause}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
                    Pause
                   </button>
                {/if}
                <button class="btn btn-ghost" style="color: var(--error);" on:click={stop}>
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/></svg>
                  Stop
                </button>
              {/if}
            </div>
          </div>

          <div class="card">
            <h2>Scraping analytics</h2>
            <div class="progress-section">
              <div class="progress-header">
                <span class="label">Overall Completion</span>
                <span class="percentage">{progress.total > 0 ? Math.round((progress.processed / progress.total) * 100) : 0}%</span>
              </div>
              <div class="progress-track">
                <div class="progress-fill" style="width: {progress.total > 0 ? (progress.processed / progress.total) * 100 : 0}%"></div>
              </div>
              
              <div class="stat-grid">
                <div class="stat-card">
                  <div class="label">Processed</div>
                  <div class="stat-value">{progress.processed}</div>
                </div>
                <div class="stat-card">
                  <div class="label">Queue</div>
                  <div class="stat-value">{progress.total}</div>
                </div>
              </div>
              
              <div style="background: rgba(0,0,0,0.2); padding: 0.8rem; border-radius: 12px; border: 1px solid var(--glass-border);">
                <div class="label" style="font-size: 0.6rem; margin-bottom: 0.2rem;">Live Stream</div>
                <div style="font-size: 0.75rem; color: var(--text-primary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; font-family: 'JetBrains Mono'; font-weight: 500;">
                  {progress.current_file || 'Awaiting transmission...'}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- System Console & Preview -->
        <div style="display: flex; flex-direction: column; gap: 2rem; min-height: 0;">
          <div class="card" style="flex: 2; display: flex; flex-direction: column; min-height: 0;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
              <h2>Live identification</h2>
              <div style="display: flex; gap: 1rem; align-items: center;">
                {#if matchedFiles.length > 0}
                  <button class="btn btn-primary" style="padding: 0.5rem 1rem; font-size: 0.75rem;" on:click={handleRenameSelected}>
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                    RENAME SELECTED
                  </button>
                  <span class="status-badge" style="font-size: 0.6rem; padding: 2px 8px;">{matchedFiles.length} MATCHES</span>
                {/if}
              </div>
            </div>
            <div class="preview-table-container scrollbar" style="flex: 1; overflow-y: auto;">
              <table class="preview-table">
                <thead>
                  <tr>
                    <th style="width: 40px; text-align: center;">
                      <input type="checkbox" bind:checked={selectedAll} on:change={toggleAll}>
                    </th>
                    <th>Original Filename</th>
                    <th>Proposed Name</th>
                  </tr>
                </thead>
                <tbody>
                  {#each matchedFiles as match}
                    <tr class={match.status === 'renamed' ? 'row-renamed' : ''}>
                      <td style="text-align: center;">
                        <input type="checkbox" bind:checked={match.selected} disabled={match.status === 'renamed'}>
                      </td>
                      <td class="original">{match.original}</td>
                      <td class="proposed">
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" style="margin-right: 8px; color: var(--success);"><polyline points="20 6 9 17 4 12"/></svg>
                        {match.proposed}
                      </td>
                    </tr>
                  {/each}
                  {#if matchedFiles.length === 0}
                    <tr>
                      <td colspan="3" style="text-align: center; padding: 3rem; opacity: 0.3; font-size: 0.8rem; letter-spacing: 0.05em;">
                        NO MATCHES FOUND YET
                      </td>
                    </tr>
                  {/if}
                </tbody>
              </table>
            </div>
          </div>

          <div class="card" style="flex: 1; min-height: 200px; display: flex; flex-direction: column;">
            <h2>System output</h2>
            <div class="console scrollbar" style="flex: 1;">
              {#each logs as log}
                <div class="log-item">
                  <span class="log-time">[{log.time}]</span>
                  <span class="log-msg-{log.type === 'error' ? 'error' : (log.type === 'success' ? 'success' : 'info')}">{log.msg}</span>
                </div>
              {/each}
              {#if logs.length === 0}
                <div style="height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; opacity: 0.3; gap: 1rem;">
                  <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
                  <div style="font-size: 0.7rem; letter-spacing: 0.1em;">NO ACTIVITY</div>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>

    {:else}
      <div class="card" style="flex: 1; overflow-y: auto;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 3rem;">
          <div>
            <h2 style="font-size: 1.75rem; color: var(--text-primary); margin-bottom: 0.25rem;">Engine Configuration</h2>
            <p style="color: var(--text-secondary); font-size: 0.9rem;">Fine-tune your scraping and renaming strategies.</p>
          </div>
          <button class="btn btn-primary" style="padding: 1rem 2.5rem;" on:click={saveConfig}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><path d="M7 3v5h8"/></svg>
            SYNC CHANGES
          </button>
        </div>

        <div style="max-width: 1000px; display: grid; grid-template-columns: 1fr 1fr; gap: 3rem;">
          <section style="display: flex; flex-direction: column; gap: 1.5rem;">
            <h3 style="font-size: 0.8rem; letter-spacing: 0.1em; color: var(--primary); font-weight: 800; text-transform: uppercase;">Providers</h3>
            
            <div class="card" style="padding: 1rem; border: 1px solid var(--glass-border);">
               <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
                  <span style="font-weight: 700;">ScreenScraper.fr</span>
                  <div style="display: flex; align-items: center; gap: 0.5rem;">
                    {#if ssActive}
                      <span class="status-badge" style="font-size: 0.6rem; padding: 2px 6px;">{apiCounts.screenscraper?.Count || 0} CALLS THIS MONTH</span>
                    {/if}
                    <input type="checkbox" bind:checked={ssActive}>
                  </div>
               </div>
               <div class="field-group">
                 <span class="label">Developer ID</span>
                 <input type="text" bind:value={screenScraperUser} placeholder="Enter ID">
               </div>
               <div class="field-group">
                 <span class="label">API Token</span>
                 <input type="password" bind:value={screenScraperPass} placeholder="••••••••••••">
               </div>
               <div style="display: flex; gap: 1rem; margin-top: 0.5rem;">
                  <label style="font-size: 0.7rem; display: flex; align-items: center; gap: 4px;">
                    <input type="checkbox" bind:checked={ssSearchByHash}> Hash
                  </label>
                  <label style="font-size: 0.7rem; display: flex; align-items: center; gap: 4px;">
                    <input type="checkbox" bind:checked={ssSearchByName}> Name
                  </label>
               </div>
            </div>

            <div class="card" style="padding: 1rem; border: 1px solid var(--glass-border);">
               <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
                  <span style="font-weight: 700;">TheGamesDB.net</span>
                  <input type="checkbox" bind:checked={tgdbActive}>
               </div>
               <div class="field-group">
                 <span class="label">API Key</span>
                 <input type="password" bind:value={tgdbApiKey} placeholder="Enter API Key">
               </div>
               <div style="display: flex; gap: 1rem; margin-top: 0.5rem;">
                  <label style="font-size: 0.7rem; display: flex; align-items: center; gap: 4px;">
                    <input type="checkbox" bind:checked={tgdbSearchByHash}> Hash
                  </label>
                  <label style="font-size: 0.7rem; display: flex; align-items: center; gap: 4px;">
                    <input type="checkbox" bind:checked={tgdbSearchByName}> Name
                  </label>
               </div>
            </div>
          </section>

          <section style="display: flex; flex-direction: column; gap: 1.5rem;">
            <h3 style="font-size: 0.8rem; letter-spacing: 0.1em; color: var(--secondary); font-weight: 800; text-transform: uppercase;">Engine Settings</h3>
            
            <div class="field-group">
              <span class="label">Naming Template</span>
              <input type="text" bind:value={namingTemplate} placeholder="{'{Name} ({Region}) {Year}'}">
              <div style="font-size: 0.6rem; color: var(--text-secondary); margin-top: 0.5rem; line-height: 1.8;">
                <span style="color: var(--primary); font-weight: 700;">Base:</span> {'{Name}'} &nbsp;{'{Region}'} &nbsp;{'{Languages}'} &nbsp;{'{Year}'}<br>
                <span style="color: var(--secondary); font-weight: 700;">Extra:</span> {'{Company}'} &nbsp;{'{Developer}'} &nbsp;{'{Genre}'} &nbsp;{'{Players}'} &nbsp;{'{Rating}'} &nbsp;{'{RomType}'} &nbsp;{'{Hash}'}
              </div>
            </div>

            <div class="field-group">
              <span class="label">Parallel Workers (CPU/IO)</span>
              <input type="number" bind:value={workerCount} min="1" max="16">
            </div>
          </section>
        </div>
      </div>
    {/if}
  </main>

  <footer>
    <div class="status-badge">
      <div class="status-dot {progress.status !== 'idle' ? 'active' : ''}"></div>
      SYSTEM STATUS: <span style="color: {progress.status === 'error' ? 'var(--error)' : 'var(--text-primary)'}">{progress.status.toUpperCase()}</span>
    </div>
    <div style="display: flex; gap: 2rem; align-items: center;">
      <span>CORE: <span style="color: var(--text-primary);">WAILS 2.0</span></span>
      <span>RUNTIME: <span style="color: var(--text-primary);">GO + SVELTE</span></span>
      <span style="background: rgba(255,255,255,0.05); padding: 2px 8px; border-radius: 4px; border: 1px solid var(--glass-border);">v1.2.0-ELITE</span>
    </div>
  </footer>
</div>

<!-- Modal: Sesión previa detectada -->
{#if showSessionModal && sessionInfo}
<div class="modal-backdrop" on:click|self={() => showSessionModal = false} on:keydown|self={(e) => e.key === 'Escape' && (showSessionModal = false)} role="dialog" aria-modal="true">
  <div class="modal-box">
    <div style="display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem;">
      {#if sessionInfo.status === 'completed'}
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--success)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        <h2 style="font-size: 1.1rem; font-weight: 700; color: var(--text-primary); margin: 0;">Escaneo anterior completado</h2>
      {:else}
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--warning)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        <h2 style="font-size: 1.1rem; font-weight: 700; color: var(--text-primary); margin: 0;">Sesión anterior incompleta</h2>
      {/if}
    </div>

    <p style="color: var(--text-secondary); font-size: 0.85rem; margin-bottom: 0.75rem;">
      {#if sessionInfo.status === 'completed'}
        Este directorio ya fue escaneado por completo. ¿Deseas iniciar un nuevo escaneo?
      {:else}
        Este directorio tiene un escaneo en proceso. ¿Deseas continuar donde lo dejaste o empezar desde cero?
      {/if}
    </p>

    <div style="background: rgba(0,0,0,0.25); border: 1px solid var(--glass-border); border-radius: 12px; padding: 1rem; margin-bottom: 1.5rem; font-size: 0.8rem;">
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.5rem;">
        <span style="color: var(--text-secondary);">Estado</span>
        <span style="color: {sessionInfo.status === 'completed' ? 'var(--success)' : 'var(--warning)'}; font-weight: 700; text-transform: uppercase;">{sessionInfo.status}</span>
      </div>
      <div style="display: flex; justify-content: space-between; margin-bottom: 0.5rem;">
        <span style="color: var(--text-secondary);">Archivos totales</span>
        <span style="color: var(--text-primary); font-weight: 700;">{sessionInfo.total_files}</span>
      </div>
      <div style="display: flex; justify-content: space-between;">
        <span style="color: var(--text-secondary);">Encontrados</span>
        <span style="color: var(--success); font-weight: 700;">{sessionInfo.done_files} / {sessionInfo.total_files}</span>
      </div>
      <div style="margin-top: 0.75rem;">
        <div style="height: 4px; background: rgba(255,255,255,0.08); border-radius: 2px; overflow: hidden;">
          <div style="height: 100%; width: {sessionInfo.total_files > 0 ? Math.round(sessionInfo.done_files/sessionInfo.total_files*100) : 0}%; background: {sessionInfo.status === 'completed' ? 'var(--success)' : 'var(--primary)'}; border-radius: 2px;"></div>
        </div>
      </div>
    </div>

    {#if sessionInfo.status === 'completed'}
      <!-- Sesión completa: nuevo escaneo o cancelar -->
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem;">
        <button class="btn btn-ghost" style="border: 1px solid var(--glass-border);" on:click={() => showSessionModal = false} disabled={isScraping}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          Cancelar
        </button>
        <button class="btn btn-primary" on:click={() => launchScraping(true)} disabled={isScraping}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-3.48"/></svg>
          {isScraping ? 'Iniciando...' : 'Nuevo escaneo'}
        </button>
      </div>
    {:else}
      <!-- Sesión incompleta: cancelar, reiniciar o reanudar -->
      <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 0.5rem;">
        <button class="btn btn-ghost" style="border: 1px solid var(--glass-border); font-size: 0.78rem;" on:click={() => showSessionModal = false} disabled={isScraping}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          Cancelar
        </button>
        <button class="btn btn-ghost" style="border: 1px solid var(--error); color: var(--error); font-size: 0.78rem;" on:click={() => launchScraping(true)} disabled={isScraping}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-3.48"/></svg>
          Reiniciar
        </button>
        <button class="btn btn-primary" style="font-size: 0.78rem;" on:click={() => launchScraping(false)} disabled={isScraping}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>
          {isScraping ? '...' : 'Reanudar'}
        </button>
      </div>
    {/if}
  </div>
</div>
{/if}

{#if isScraping}
<div class="active-process-overlay">
  <div class="spinner"></div>
  <h2 style="font-size: 1.5rem; font-weight: 700; color: var(--text-primary);">ENGINE ACTIVE</h2>
  <p style="color: var(--primary); font-family: 'JetBrains Mono', monospace; font-size: 0.9rem; margin-top: 0.5rem;">
    {progress.processed} / {progress.total} FILES PROCESSED
  </p>
</div>
{/if}

{#if showRenameErrorsModal}
<div class="modal-backdrop" on:click|self={() => showRenameErrorsModal = false} on:keydown|self={(e) => e.key === 'Escape' && (showRenameErrorsModal = false)} role="dialog" aria-modal="true">
  <div class="modal-box" style="width: 600px; max-width: 95vw;">
    <div style="display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem;">
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="var(--error)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
      <h2 style="font-size: 1.1rem; font-weight: 700; color: var(--text-primary); margin: 0;">Errores de Renombramiento</h2>
    </div>

    <p style="color: var(--text-secondary); font-size: 0.85rem; margin-bottom: 1rem;">
      Se encontraron los siguientes problemas al procesar los archivos:
    </p>

    <div class="scrollbar" style="background: rgba(0,0,0,0.25); border: 1px solid var(--glass-border); border-radius: 12px; padding: 1rem; max-height: 300px; overflow-y: auto; margin-bottom: 1.5rem;">
      {#each renameErrors as err}
        <div style="margin-bottom: 0.75rem; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 0.5rem;">
          <div style="font-weight: bold; color: var(--primary); font-family: monospace; font-size: 0.8rem;">{err.filename}</div>
          <div style="color: var(--error); font-size: 0.75rem;">{err.reason}</div>
        </div>
      {/each}
    </div>

    <div style="display: flex; justify-content: flex-end;">
      <button class="btn btn-primary" on:click={() => showRenameErrorsModal = false}>
        Cerrar
      </button>
    </div>
  </div>
</div>
{/if}

<style>
  :global(body) {
    user-select: none;
    -webkit-user-select: none;
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.65);
    backdrop-filter: blur(6px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 0.15s ease;
  }

  .modal-box {
    background: #1a1d2e;
    border: 1px solid rgba(255,255,255,0.12);
    border-radius: 20px;
    padding: 2rem;
    width: 440px;
    max-width: 92vw;
    box-shadow: 0 32px 100px rgba(0,0,0,0.85), 0 0 0 1px rgba(255,255,255,0.04) inset;
    animation: slideUp 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to   { opacity: 1; }
  }
  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to   { transform: translateY(0);    opacity: 1; }
  }
</style>

