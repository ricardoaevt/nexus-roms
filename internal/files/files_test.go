package files_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"romsRename/internal/files"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesScanner(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Scan directory with ROMs", func(t *testing.T) {
		p1 := filepath.Join(tmpDir, "game1.sfc")
		_ = os.WriteFile(p1, []byte("content"), 0644)
		_ = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not a rom"), 0644)

		scanned, err := files.ScanDirectory(tmpDir)
		assert.NoError(t, err)
		assert.Len(t, scanned, 1)
		assert.Equal(t, "game1.sfc", scanned[0].Filename)
		
		_ = os.Remove(p1)
	})

	t.Run("Scan directory with ZIP", func(t *testing.T) {
		zipPath := filepath.Join(tmpDir, "collection.zip")
		f, _ := os.Create(zipPath)
		w := zip.NewWriter(f)
		f1, _ := w.Create("mario.sfc")
		_, _ = f1.Write([]byte("mario rom"))
		_ = w.Close()
		_ = f.Close()

		scanned, err := files.ScanDirectory(tmpDir)
		assert.NoError(t, err)
		assert.Len(t, scanned, 1)
		assert.Equal(t, "mario.sfc", scanned[0].Filename)
		assert.Equal(t, "collection.zip", scanned[0].ContainerPath)
	})
}

func TestFilesHashing(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.bin")
	content := []byte("hello world")
	_ = os.WriteFile(path, content, 0644)

	t.Run("HashFile", func(t *testing.T) {
		hashes, err := files.HashFile(path)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashes.MD5)
		// md5 for "hello world" is 5eb63bbbe01eeed093cb22bb8f5acdc3
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hashes.MD5)
	})

	t.Run("HashCompressedFile (ZIP)", func(t *testing.T) {
		zipPath := filepath.Join(tmpDir, "hashing.zip")
		f, _ := os.Create(zipPath)
		w := zip.NewWriter(f)
		f1, _ := w.Create("inner.rom")
		_, _ = f1.Write([]byte("hello world"))
		_ = w.Close()
		_ = f.Close()

		hashes, err := files.HashCompressedFile(zipPath, "inner.rom")
		assert.NoError(t, err)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", hashes.MD5)
	})
}
