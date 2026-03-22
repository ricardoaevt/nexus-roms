package db_test

import (
	"romsRename/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

const inMemoryDB = ":memory:"

func TestDBSessions(t *testing.T) {
	database, err := db.InitDB(inMemoryDB)
	assert.NoError(t, err)
	defer database.Close()

	t.Run("Create and Get Latest Session", func(t *testing.T) {
		id, err := database.CreateSession("/root/path")
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))

		latest, err := database.GetLatestSession()
		assert.NoError(t, err)
		assert.Equal(t, "/root/path", latest.RootPath)
		assert.Equal(t, "running", latest.Status)
	})

	t.Run("Add File and Update Metadata", func(t *testing.T) {
		sessionID, _ := database.CreateSession("/path")
		file := db.SessionFile{
			SessionID:    sessionID,
			RelativePath: "game.sfc",
			Filename:     "game.sfc",
		}
		fileID, err := database.AddFile(file)
		assert.NoError(t, err)

		// Fetch and update
		pending, err := database.GetPendingFiles(sessionID)
		assert.NoError(t, err)
		assert.Len(t, pending, 1)
		
		f := pending[0]
		f.Status = "found"
		err = database.UpdateFileMetadata(f)
		assert.NoError(t, err)

		updated, err := database.GetFileByID(fileID)
		assert.NoError(t, err)
		assert.Equal(t, "found", updated.Status)
	})

	t.Run("Update Session Status and Progress", func(t *testing.T) {
		sessionID, _ := database.CreateSession("/status_path")
		err := database.UpdateSessionStatus(sessionID, "paused")
		assert.NoError(t, err)

		latest, _ := database.GetLatestSession()
		assert.Equal(t, "paused", latest.Status)

		_, done := database.GetSessionProgress(sessionID)
		assert.Equal(t, 0, done)
	})
}

func TestDBConfig(t *testing.T) {
	database, _ := db.InitDB(inMemoryDB)
	defer database.Close()

	t.Run("Get and Save Config", func(t *testing.T) {
		val := database.GetConfig("key", "default")
		assert.Equal(t, "default", val)

		err := database.SaveConfig("key", "new_value")
		assert.NoError(t, err)

		val = database.GetConfig("key", "default")
		assert.Equal(t, "new_value", val)
	})
}

func TestDBCredentials(t *testing.T) {
	database, _ := db.InitDB(inMemoryDB)
	defer database.Close()

	t.Run("Save and Get Credentials", func(t *testing.T) {
		creds := db.APICredentials{
			Provider:     "ScreenScraper",
			Username:     "user",
			Password:     "secret",
			IsActive:     true,
			SearchByHash: true,
		}
		err := database.SaveCredentials(creds)
		assert.NoError(t, err)

		saved, err := database.GetCredentials("ScreenScraper")
		assert.NoError(t, err)
		assert.Equal(t, "user", saved.Username)
		assert.Equal(t, "secret", saved.Password)
		assert.True(t, saved.IsActive)
	})
}
