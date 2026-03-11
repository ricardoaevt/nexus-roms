package orchestrator

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"romsRename/internal/db"
	"romsRename/internal/files"
	"romsRename/internal/renamer"
	"romsRename/internal/scraper"
	"sync"
	"sync/atomic"
	"time"
)

type Progress struct {
	ID           int64  `json:"id,omitempty"`
	Total        int32  `json:"total"`
	Processed    int32  `json:"processed"`
	CurrentFile  string `json:"current_file"`
	ProposedName string `json:"proposed_name"`
	Status       string         `json:"status"`
	Message      string         `json:"message,omitempty"`
	APICounts    map[string]int `json:"api_counts,omitempty"`
}

type Orchestrator struct {
	db       *db.DB
	scrapers []scraper.Scraper
	
	// Control de flujo
	cancel     context.CancelFunc
	pauseChan  chan struct{}
	resumeChan chan struct{}
	isPaused   int32 // 0 = running, 1 = paused
	
	progCallback func(Progress)
	
	// Configuración
	Template    string
	MaxWorkers  int
	
	wg sync.WaitGroup
}

func NewOrchestrator(database *db.DB, scrapers []scraper.Scraper, progCb func(Progress)) *Orchestrator {
	return &Orchestrator{
		db:           database,
		scrapers:     scrapers,
		pauseChan:    make(chan struct{}),
		resumeChan:   make(chan struct{}),
		progCallback: progCb,
		Template:     database.GetConfig("naming_template", "{Name} ({Region})"),
		MaxWorkers:   1, // será sobreescrito por app.go leyendo la BD
		wg:           sync.WaitGroup{},
	}
}

func (o *Orchestrator) getAPICounts() map[string]int {
	counts := make(map[string]int)
	providers := []string{"screenscraper", "thegamesdb"}
	for _, p := range providers {
		val := o.db.GetConfig("api_tracker_"+p, "{}")
		var data struct{ Count int }
		_ = json.Unmarshal([]byte(val), &data)
		counts[p] = data.Count
	}
	return counts
}

func (o *Orchestrator) Start(rootPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	o.cancel = cancel

	// 1. Verificar si existe sesión previa incompleta para esta ruta
	latest, _ := o.db.GetLatestSession()
	var sessionID int64
	var err error

	if latest != nil && latest.RootPath == rootPath && (latest.Status == "running" || latest.Status == "paused" || latest.Status == "stopped") {
		sessionID = latest.ID
		o.db.UpdateSessionStatus(sessionID, "running")
	} else {
		// Crear nueva sesión
		sessionID, err = o.db.CreateSession(rootPath)
		if err != nil {
			return err
		}

		// Escanear archivos (solo si es nueva sesión)
		scanned, err := files.ScanDirectory(rootPath)
		if err != nil {
			return err
		}

		for _, f := range scanned {
			_, err := o.db.AddFile(db.SessionFile{
				SessionID:     sessionID,
				RelativePath:  f.RelativePath,
				Filename:      f.Filename,
				ContainerPath: sql.NullString{String: f.ContainerPath, Valid: f.ContainerPath != ""},
			})
			if err != nil {
				return err
			}
		}
	}

	// 2. Obtener total real
	pending, _ := o.db.GetPendingFiles(sessionID)

	// 3. Iniciar trabajadores
	go o.workerLoop(ctx, sessionID, int32(len(pending)), rootPath)

	return nil
}

func (o *Orchestrator) workerLoop(ctx context.Context, sessionID int64, total int32, rootPath string) {
	pending, err := o.db.GetPendingFiles(sessionID)
	if err != nil {
		return
	}

	var processed int32
	fileChan := make(chan db.SessionFile, len(pending))
	for _, f := range pending {
		fileChan <- f
	}
	close(fileChan)

	numWorkers := o.MaxWorkers
	if numWorkers <= 0 {
		numWorkers = 1 // mínimo seguro si la BD devuelve un valor inválido
	}
	o.wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go o.runWorker(ctx, fileChan, &processed, total, rootPath)
	}

	o.wg.Wait()
	status := "completed"
	if ctx.Err() != nil {
		status = "stopped"
	}
	o.db.UpdateSessionStatus(sessionID, status)
	
	if o.progCallback != nil {
		o.progCallback(Progress{Status: status, APICounts: o.getAPICounts()})
	}
}

func (o *Orchestrator) runWorker(ctx context.Context, fileChan <-chan db.SessionFile, processed *int32, total int32, rootPath string) {
	defer o.wg.Done()
	for f := range fileChan {
		if atomic.LoadInt32(&o.isPaused) == 1 {
			<-o.resumeChan
		}

		select {
		case <-ctx.Done():
			return
		default:
			o.processFile(ctx, f, rootPath, processed, total)
		}
	}
}

func (o *Orchestrator) processFile(ctx context.Context, f db.SessionFile, rootPath string, processed *int32, total int32) {
	f.Status = "hashing"
	o.db.UpdateFileMetadata(f)

	var hashes *files.FileHashes
	var err error

	// Find displaying name (container if exists, otherwise filename)
	displayName := f.Filename
	if f.ContainerPath.Valid {
		displayName = filepath.Base(f.ContainerPath.String)
	}

	fullPath := filepath.Join(rootPath, f.RelativePath)
	if f.ContainerPath.Valid {
		fullPath = filepath.Join(rootPath, f.ContainerPath.String)
		if o.progCallback != nil {
			o.progCallback(Progress{Total: total, Processed: *processed, CurrentFile: displayName, Status: "running", Message: "Hashing internal: " + f.RelativePath, APICounts: o.getAPICounts()})
		}
		hashes, err = files.HashCompressedFile(fullPath, f.RelativePath)
	} else {
		if o.progCallback != nil {
			o.progCallback(Progress{Total: total, Processed: *processed, CurrentFile: displayName, Status: "running", Message: "Hashing file...", APICounts: o.getAPICounts()})
		}
		hashes, err = files.HashFile(fullPath)
	}

	if err != nil {
		f.Status = "error"
		f.ErrorMessage = sql.NullString{String: err.Error(), Valid: true}
		o.db.UpdateFileMetadata(f)
		return
	}

	f.HashMD5 = sql.NullString{String: hashes.MD5, Valid: true}
	f.HashSHA1 = sql.NullString{String: hashes.SHA1, Valid: true}
	f.HashCRC32 = sql.NullString{String: hashes.CRC32, Valid: true}
	f.Status = "scraping"
	o.db.UpdateFileMetadata(f)

	// Scraping
	if o.progCallback != nil {
		o.progCallback(Progress{Total: total, Processed: *processed, CurrentFile: displayName, Status: "running", Message: "Searching in APIs...", APICounts: o.getAPICounts()})
	}

	query := scraper.SearchQuery{
		Filename: f.Filename,
		HashMD5:  hashes.MD5,
		HashSHA1: hashes.SHA1,
		HashCRC32: hashes.CRC32,
	}

	for _, s := range o.scrapers {
		if !s.CanSearchByHash() {
			continue
		}

		// Rate limit simple
		select {
		case <-ctx.Done():
			return
		case <-time.After(500 * time.Millisecond):
		}

		meta, err := s.SearchByHash(ctx, query)
		if err == nil && meta != nil {
			o.applyMetadata(&f, meta, hashes.MD5, processed, total, s.Name())
			return
		}
	}

	// Fallback a nombre si no se encontró por hash
	for _, s := range o.scrapers {
		if !s.CanSearchByName() {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(500 * time.Millisecond):
		}

		meta, err := s.SearchByName(ctx, query)
		if err == nil && meta != nil {
			o.applyMetadata(&f, meta, hashes.MD5, processed, total, s.Name())
			return
		}
	}

	f.Status = "not_found"
	o.db.UpdateFileMetadata(f)
	
	p := atomic.AddInt32(processed, 1)
	if o.progCallback != nil {
		o.progCallback(Progress{
			ID:          f.ID,
			Total:       total,
			Processed:   p,
			CurrentFile: displayName,
			Status:      "running",
			APICounts:   o.getAPICounts(),
		})
	}
}

func (o *Orchestrator) applyMetadata(f *db.SessionFile, meta *scraper.Metadata, md5 string, processed *int32, total int32, provider string) {
	f.NameMetadata = sql.NullString{String: meta.Name, Valid: true}
	f.RegionMetadata = sql.NullString{String: meta.Region, Valid: true}
	f.YearMetadata = sql.NullString{String: meta.Year, Valid: true}
	f.CompanyMetadata = sql.NullString{String: meta.Company, Valid: true}
	f.Status = "found"

	// Generar nuevo nombre basado en plantilla dinámica
	newName := renamer.FormatName(o.Template, renamer.TemplateData{
		Name:      meta.Name,
		Region:    meta.Region,
		Languages: meta.Languages,
		Year:      meta.Year,
		Company:   meta.Company,
		Developer: meta.Developer,
		Genre:     meta.Genre,
		Players:   meta.Players,
		Rating:    meta.Rating,
		RomType:   meta.RomType,
		Hash:      md5,
	})
	f.NewName = sql.NullString{String: newName, Valid: true}
	
	o.db.UpdateFileMetadata(*f)

	// Find displaying name and extension
	displayName := f.Filename
	var ext string
	if f.ContainerPath.Valid {
		displayName = filepath.Base(f.ContainerPath.String)
		ext = filepath.Ext(f.ContainerPath.String)
	} else {
		ext = filepath.Ext(f.RelativePath)
	}

	p := atomic.AddInt32(processed, 1)
	if o.progCallback != nil {
		o.progCallback(Progress{
			ID:           f.ID,
			Total:        total,
			Processed:    p,
			CurrentFile:  displayName,
			ProposedName: newName + ext,
			Status:       "running",
			Message:      "Identified by " + provider,
			APICounts:    o.getAPICounts(),
		})
	}
}

func (o *Orchestrator) Pause() {
	if atomic.CompareAndSwapInt32(&o.isPaused, 0, 1) {
		if o.progCallback != nil {
			o.progCallback(Progress{Status: "paused"})
		}
	}
}

func (o *Orchestrator) Resume() {
	if atomic.CompareAndSwapInt32(&o.isPaused, 1, 0) {
		close(o.resumeChan)
		o.resumeChan = make(chan struct{})
		if o.progCallback != nil {
			o.progCallback(Progress{Status: "running"})
		}
	}
}

func (o *Orchestrator) Stop() {
	if o.cancel != nil {
		o.cancel()
		o.Resume() // Asegurar que los trabajadores salgan de la pausa
	}
}
