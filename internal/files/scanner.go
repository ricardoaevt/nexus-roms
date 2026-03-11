package files

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode"
	"github.com/saracen/go7z"
)

var RomExtensions = map[string]bool{
	".nes": true, ".sfc": true, ".smc": true, ".gb": true, ".gbc": true, ".gba": true,
	".n64": true, ".z64": true, ".v64": true, ".md": true, ".gen": true, ".sms": true,
	".gg": true, ".pce": true, ".cue": true, ".bin": true, ".iso": true, ".chd": true,
}

var ArchiveExtensions = map[string]bool{
	".zip": true, ".rar": true, ".7z": true,
}

type ScannedFile struct {
	RelativePath  string
	Filename      string
	ContainerPath string // Si está dentro de un ZIP/RAR/7Z
	IsArchive     bool
}

func ScanDirectory(root string) ([]ScannedFile, error) {
	var files []ScannedFile

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "duplicados" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		relPath, _ := filepath.Rel(root, path)

		if RomExtensions[ext] {
			files = append(files, ScannedFile{
				RelativePath: relPath,
				Filename:     info.Name(),
			})
		} else if ArchiveExtensions[ext] {
			// Explorar dentro del archivo comprimido
			archived, err := listArchiveContents(path, relPath)
			if err == nil {
				if isValidArchiveCollection(archived) {
					files = append(files, archived...)
				}
			}
		}

		return nil
	})

	return files, err
}

// isValidArchiveCollection aplica heurística para determinar si los archivos dentro
// de un contenedor comprimido están relacionados (ej. discos múltiples) o si
// es un paquete desconectado (romset/recopilación).
func isValidArchiveCollection(files []ScannedFile) bool {
	if len(files) == 0 {
		return false
	}
	if len(files) == 1 {
		return true
	}
	if len(files) > 10 {
		return false // Umbral lógico excedido, probablemente es un romset
	}

	prefix := files[0].Filename
	for i := 1; i < len(files); i++ {
		prefix = longestCommonPrefix(prefix, files[i].Filename)
	}

	prefix = strings.TrimSpace(prefix)

	if len(prefix) < 4 {
		return false
	}

	minLen := len(files[0].Filename)
	for _, f := range files {
		if len(f.Filename) < minLen {
			minLen = len(f.Filename)
		}
	}

	if float64(len(prefix)) < float64(minLen)*0.3 {
		return false
	}

	return true
}

func longestCommonPrefix(a, b string) string {
	minLength := len(a)
	if len(b) < minLength {
		minLength = len(b)
	}
	for i := 0; i < minLength; i++ {
		if strings.ToLower(string(a[i])) != strings.ToLower(string(b[i])) {
			return a[:i]
		}
	}
	return a[:minLength]
}

func listArchiveContents(path, relPath string) ([]ScannedFile, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".zip":
		return listZipContents(path, relPath)
	case ".rar":
		return listRarContents(path, relPath)
	case ".7z":
		return list7zContents(path, relPath)
	}
	return nil, fmt.Errorf("formato no soportado: %s", ext)
}

func listZipContents(path, relPath string) ([]ScannedFile, error) {
	var result []ScannedFile
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && RomExtensions[strings.ToLower(filepath.Ext(f.Name))] {
			result = append(result, ScannedFile{
				RelativePath:  f.Name,
				Filename:      filepath.Base(f.Name),
				ContainerPath: relPath,
			})
		}
	}
	return result, nil
}

func listRarContents(path, relPath string) ([]ScannedFile, error) {
	var result []ScannedFile
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rr, err := rardecode.NewReader(f, "")
	if err != nil {
		return nil, err
	}
	for {
		header, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !header.IsDir && RomExtensions[strings.ToLower(filepath.Ext(header.Name))] {
			result = append(result, ScannedFile{
				RelativePath:  header.Name,
				Filename:      filepath.Base(header.Name),
				ContainerPath: relPath,
			})
		}
	}
	return result, nil
}

func list7zContents(path, relPath string) ([]ScannedFile, error) {
	var result []ScannedFile
	sz, err := go7z.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer sz.Close()
	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !hdr.IsEmptyStream && !hdr.IsEmptyFile && RomExtensions[strings.ToLower(filepath.Ext(hdr.Name))] {
			result = append(result, ScannedFile{
				RelativePath:  hdr.Name,
				Filename:      filepath.Base(hdr.Name),
				ContainerPath: relPath,
			})
		}
	}
	return result, nil
}
