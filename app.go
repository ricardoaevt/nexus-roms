package main

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"romsRename/internal/db"
	"romsRename/internal/orchestrator"
	"romsRename/internal/renamer"
	"romsRename/internal/scraper"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	database     *db.DB
	orchestrator *orchestrator.Orchestrator
	mu           sync.Mutex
	isRunning    bool
}

// NewApp creates a new App struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Inicializar DB en el directorio del usuario
	home, _ := filepath.Abs(".")
	dbPath := filepath.Join(home, "roms_renamer.db")
	
	database, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Printf("Error al inicializar DB: %v\n", err)
		return
	}
	a.database = database
}

// SessionInfo contiene información sobre una sesión previa encontrada
type SessionInfo struct {
	Found       bool   `json:"found"`
	Status      string `json:"status"`
	RootPath    string `json:"root_path"`
	TotalFiles  int    `json:"total_files"`
	DoneFiles   int    `json:"done_files"`
}

// CheckPreviousSession verifica si existe una sesión previa para la ruta dada
func (a *App) CheckPreviousSession(rootPath string) SessionInfo {
	if a.database == nil {
		return SessionInfo{Found: false}
	}
	latest, err := a.database.GetLatestSession()
	if err != nil || latest == nil {
		return SessionInfo{Found: false}
	}
	// Informar si es la misma ruta, sin importar el estado
	if latest.RootPath != rootPath {
		return SessionInfo{Found: false}
	}
	total, done := a.database.GetSessionProgress(latest.ID)
	return SessionInfo{
		Found:      true,
		Status:     latest.Status,
		RootPath:   latest.RootPath,
		TotalFiles: total,
		DoneFiles:  done,
	}
}

// StartScraping inicia el proceso de scraping.
// forceRestart=true descarta la sesión anterior y empieza desde cero.
func (a *App) StartScraping(rootPath string, forceRestart bool) error {
	a.mu.Lock()
	if a.isRunning {
		// Ya hay un scraping en curso: detenerlo antes de continuar
		a.mu.Unlock()
		if a.orchestrator != nil {
			a.orchestrator.Stop()
		}
		a.mu.Lock()
	}
	a.isRunning = true
	a.mu.Unlock()

	// Si se pide reinicio, marcar sesión previa como completada
	if forceRestart {
		if latest, err := a.database.GetLatestSession(); err == nil && latest != nil && latest.RootPath == rootPath {
			a.database.UpdateSessionStatus(latest.ID, "completed")
		}
	}

	var scrapers []scraper.Scraper

	// 1. ScreenScraper
	ssCreds, err := a.database.GetCredentials("screenscraper")
	if err == nil && ssCreds.IsActive {
		scrapers = append(scrapers, scraper.NewRetryScraper(scraper.NewScreenScraperClient(ssCreds, a.database), 3, 2*time.Second))
	}

	// 2. TheGamesDB
	tgdbCreds, err := a.database.GetCredentials("thegamesdb")
	if err == nil && tgdbCreds.IsActive {
		scrapers = append(scrapers, scraper.NewRetryScraper(scraper.NewTheGamesDBClient(tgdbCreds, a.database), 3, 2*time.Second))
	}

	if len(scrapers) == 0 {
		a.mu.Lock()
		a.isRunning = false
		a.mu.Unlock()
		return fmt.Errorf("no hay ningun proveedor de scraping activo. Configura al menos uno en Settings")
	}

	tmpl := a.database.GetConfig("naming_template", "{Name} ({Region})")
	workerStr := a.database.GetConfig("worker_count", "4")
	workers := 4
	if w, err := strconv.Atoi(workerStr); err == nil && w > 0 {
		workers = w
	}

	a.orchestrator = orchestrator.NewOrchestrator(a.database, scrapers, func(p orchestrator.Progress) {
		runtime.EventsEmit(a.ctx, "progress", p)
		// Cuando termina (completed o stopped), liberar el flag
		if p.Status == "completed" || p.Status == "stopped" {
			a.mu.Lock()
			a.isRunning = false
			a.mu.Unlock()
		}
	})
	a.orchestrator.Template = tmpl
	a.orchestrator.MaxWorkers = workers

	return a.orchestrator.Start(rootPath)
}

func (a *App) PauseScraping() {
	if a.orchestrator != nil {
		a.orchestrator.Pause()
	}
}

func (a *App) ResumeScraping() {
	if a.orchestrator != nil {
		a.orchestrator.Resume()
	}
}

func (a *App) StopScraping() {
	if a.orchestrator != nil {
		a.orchestrator.Stop()
	}
}

func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Carpeta de ROMs",
	})
}

func (a *App) SaveAPICredentials(provider, username, password, apiKey, baseURL string, searchHash, searchName, isActive bool) error {
	creds := db.APICredentials{
		Provider:     provider,
		Username:     username,
		Password:     password,
		APIKey:       apiKey,
		BaseURL:      baseURL,
		IsActive:     isActive,
		SearchByHash: searchHash,
		SearchByName: searchName,
	}
	return a.database.SaveCredentials(creds)
}

func (a *App) GetAPICredentials(provider string) (*db.APICredentials, error) {
	return a.database.GetCredentials(provider)
}

func (a *App) SaveConfigValue(key, value string) error {
	return a.database.SaveConfig(key, value)
}

func (a *App) GetConfigValue(key, defaultValue string) string {
	return a.database.GetConfig(key, defaultValue)
}

type ErrorResult struct {
	Filename string `json:"filename"`
	Reason   string `json:"reason"`
}

// ApplyRenaming ejecuta el renombrado de los archivos seleccionados
func (a *App) ApplyRenaming(ids []int64) []ErrorResult {
	var errorsList []ErrorResult

	if a.database == nil {
		errorsList = append(errorsList, ErrorResult{Filename: "System", Reason: "base de datos no inicializada"})
		return errorsList
	}

	// Obtener la sesión actual para saber el rootPath
	session, err := a.database.GetLatestSession()
	if err != nil {
		errorsList = append(errorsList, ErrorResult{Filename: "System", Reason: fmt.Sprintf("no se encontró una sesión activa: %v", err)})
		return errorsList
	}

	successCount := 0

	for _, id := range ids {
		file, err := a.database.GetFileByID(id)
		if err != nil {
			errorsList = append(errorsList, ErrorResult{Filename: fmt.Sprintf("ID: %d", id), Reason: err.Error()})
			continue
		}

		if file.Status != "found" && file.Status != "error" {
			continue
		}

		err = renamer.RenameFile(session.RootPath, *file)
		if err != nil {
			file.Status = "error"
			file.ErrorMessage = sql.NullString{String: err.Error(), Valid: true}
			a.database.UpdateFileMetadata(*file)
			errorsList = append(errorsList, ErrorResult{Filename: file.Filename, Reason: err.Error()})
		} else {
			file.Status = "renamed"
			a.database.UpdateFileMetadata(*file)
			successCount++
		}
	}

	if len(errorsList) == 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Renombrado Completado",
			Message: fmt.Sprintf("Se han renombrado %d archivos correctamente.", successCount),
		})
	}

	return errorsList
}
