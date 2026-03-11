package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"romsRename/internal/db"
	"strings"
	"time"
)

type TemplateData struct {
	Name      string
	Region    string
	Languages string
	Year      string
	Company   string
	Developer string
	Genre     string
	Players   string
	Rating    string
	RomType   string
	Hash      string
}

// FormatName genera un nombre basado en una plantilla y metadatos
func FormatName(template string, data TemplateData) string {
	if template == "" { // Corrected the line that was malformed in the instruction
		template = "{Name} ({Region})"
	}

	result := template
	result = strings.ReplaceAll(result, "{Name}", data.Name)
	result = strings.ReplaceAll(result, "{Region}", data.Region)
	result = strings.ReplaceAll(result, "{Languages}", data.Languages)
	result = strings.ReplaceAll(result, "{Year}", data.Year)
	result = strings.ReplaceAll(result, "{Company}", data.Company)
	result = strings.ReplaceAll(result, "{Developer}", data.Developer)
	result = strings.ReplaceAll(result, "{Genre}", data.Genre)
	result = strings.ReplaceAll(result, "{Players}", data.Players)
	result = strings.ReplaceAll(result, "{Rating}", data.Rating)
	result = strings.ReplaceAll(result, "{RomType}", data.RomType)
	result = strings.ReplaceAll(result, "{Hash}", data.Hash)

	// Limpiar caracteres no válidos para nombres de archivo
	invalid := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "-")
	}

	return strings.TrimSpace(result)
}

// RenameFile ejecuta el renombrado físico en el disco
func RenameFile(rootPath string, file db.SessionFile) error {
	if !file.NewName.Valid || file.NewName.String == "" {
		return fmt.Errorf("el archivo no tiene un nuevo nombre validado")
	}

	var oldPath string
	var relPath string
	var ext string

	if file.ContainerPath.Valid {
		// Renombramos el contenedor
		oldPath = filepath.Join(rootPath, file.ContainerPath.String)
		relPath = file.ContainerPath.String
		ext = filepath.Ext(file.ContainerPath.String)
	} else {
		// Renombrado directo
		oldPath = filepath.Join(rootPath, file.RelativePath)
		relPath = file.RelativePath
		ext = filepath.Ext(file.RelativePath)
	}

	newDir := filepath.Dir(oldPath)
	newPath := filepath.Join(newDir, file.NewName.String+ext)

	if oldPath == newPath {
		return nil
	}

	// Si hay colisión, mover a duplicados manteniendo la jerarquía relativa
	if _, err := os.Stat(newPath); err == nil {
		dupeDir := filepath.Join(rootPath, "duplicados", filepath.Dir(relPath))
		if mkErr := os.MkdirAll(dupeDir, 0755); mkErr != nil {
			return fmt.Errorf("error al crear directorio de duplicados: %w", mkErr)
		}
		
		// Intentar nombre en carpeta duplicados
		newPath = filepath.Join(dupeDir, file.NewName.String+ext)
		
		// Si también hay colisión en duplicados, añadir sufijo temporal/único
		if _, err := os.Stat(newPath); err == nil {
			newPath = filepath.Join(dupeDir, file.NewName.String+"_sep_"+fmt.Sprintf("%x", time.Now().UnixNano()%0xFFF)+ext)
		}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("error al renombrar (posible bloqueo u origen faltante): %w", err)
	}
	return nil
}
