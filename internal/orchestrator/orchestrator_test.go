package orchestrator_test

import (
	"os"
	"path/filepath"
	"romsRename/internal/db"
	"romsRename/internal/orchestrator"
	"romsRename/internal/scraper"
	"romsRename/internal/scraper/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrchestratorStart(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	database, err := db.InitDB(dbPath)
	assert.NoError(t, err)
	defer database.Close()

	// Create a test file
	romRel := "mario.sfc"
	romPath := filepath.Join(tmpDir, romRel)
	content := []byte("mario rom content")
	err = os.WriteFile(romPath, content, 0644)
	assert.NoError(t, err)

	// Mock scraper
	mockScraper := mocks.NewScraper(t)
	mockScraper.On("Name").Return("MockScraper").Maybe()
	mockScraper.On("CanSearchByHash").Return(true).Maybe()
	mockScraper.On("CanSearchByName").Return(true).Maybe()
	mockScraper.On("SearchByHash", mock.Anything, mock.Anything).Return(&scraper.Metadata{
		Name:   "Super Mario World",
		Region: "USA",
	}, nil).Maybe()

	// Progress callback
	done := make(chan bool, 1)
	var maxTotal, maxProcessed int32
	progCb := func(p orchestrator.Progress) {
		if p.Total > 0 {
			maxTotal = p.Total
		}
		if p.Processed > 0 {
			maxProcessed = p.Processed
		}
		if p.Status == "completed" {
			select {
			case done <- true:
			default:
			}
		}
	}

	o := orchestrator.NewOrchestrator(database, []scraper.Scraper{mockScraper}, progCb)
	o.MaxWorkers = 2

	err = o.Start(tmpDir)
	assert.NoError(t, err)

	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for orchestrator to finish")
	}

	assert.Equal(t, int32(1), maxTotal)
	assert.Equal(t, int32(1), maxProcessed)

	// Verify DB update
	session, err := database.GetLatestSession()
	assert.NoError(t, err)
	assert.Equal(t, "completed", session.Status)
}

func TestOrchestratorNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	database, _ := db.InitDB(":memory:")
	defer database.Close()

	_ = os.WriteFile(filepath.Join(tmpDir, "unknown.sfc"), []byte("unknown"), 0644)

	mockScraper := mocks.NewScraper(t)
	mockScraper.On("CanSearchByHash").Return(true).Maybe()
	mockScraper.On("CanSearchByName").Return(true).Maybe()
	mockScraper.On("SearchByHash", mock.Anything, mock.Anything).Return(nil, nil)
	mockScraper.On("SearchByName", mock.Anything, mock.Anything).Return(nil, nil)

	done := make(chan bool, 1)
	progCb := func(p orchestrator.Progress) {
		if p.Status == "completed" {
			done <- true
		}
	}

	o := orchestrator.NewOrchestrator(database, []scraper.Scraper{mockScraper}, progCb)
	_ = o.Start(tmpDir)

	<-done
	
	f, _ := database.GetFileByID(1)
	assert.Equal(t, "not_found", f.Status)
}

func TestOrchestratorPauseResume(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_pause.db")
	database, _ := db.InitDB(dbPath)
	defer database.Close()

	mockScraper := mocks.NewScraper(t)
	mockScraper.On("Name").Return("MockScraper").Maybe()
	mockScraper.On("CanSearchByHash").Return(true).Maybe()

	o := orchestrator.NewOrchestrator(database, []scraper.Scraper{mockScraper}, nil)
	o.MaxWorkers = 1

	o.Pause()
	o.Resume()
}
