package renamer

import (
	"database/sql"
	"os"
	"path/filepath"
	"romsRename/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatName(t *testing.T) {
	data := TemplateData{
		Name:      "Super Mario World",
		Region:    "USA",
		Languages: "En",
		Year:      "1991",
		Company:   "Nintendo",
		Developer: "Nintendo",
		Genre:     "Platform",
		Players:   "2",
		Rating:    "5/5",
		RomType:   "rom",
		Hash:      "123abc456def",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "default template",
			template: "",
			expected: "Super Mario World (USA)",
		},
		{
			name:     "custom template all fields",
			template: "{Name} [{Region}] [{Year}] [{Genre}]",
			expected: "Super Mario World [USA] [1991] [Platform]",
		},
		{
			name:     "invalid characters in metadata",
			template: "{Name}",
			expected: "Super Mario World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatName(tt.template, data)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("cleaning invalid characters", func(t *testing.T) {
		dataWithInvalid := TemplateData{Name: "Mario: The <Game>?"}
		result := FormatName("{Name}", dataWithInvalid)
		assert.Equal(t, "Mario- The -Game--", result)
	})
}

func TestRenameFile(t *testing.T) {
	tmpRoot := t.TempDir()

	t.Run("success direct rename", func(t *testing.T) {
		relPath := "subdir/old_rom.sfc"
		oldFullPath := filepath.Join(tmpRoot, relPath)
		err := os.MkdirAll(filepath.Dir(oldFullPath), 0755)
		assert.NoError(t, err)
		err = os.WriteFile(oldFullPath, []byte("rom content"), 0644)
		assert.NoError(t, err)

		file := db.SessionFile{
			RelativePath: relPath,
			NewName:      sql.NullString{String: "New Super Mario", Valid: true},
		}

		err = RenameFile(tmpRoot, file)
		assert.NoError(t, err)

		newPath := filepath.Join(tmpRoot, "subdir", "New Super Mario.sfc")
		assert.FileExists(t, newPath)
		assert.NoFileExists(t, oldFullPath)
	})

	t.Run("success container rename", func(t *testing.T) {
		containerRel := "archive.zip"
		containerPath := filepath.Join(tmpRoot, containerRel)
		err := os.WriteFile(containerPath, []byte("zip content"), 0644)
		assert.NoError(t, err)

		file := db.SessionFile{
			RelativePath:  "internal/file.sfc",
			ContainerPath: sql.NullString{String: containerRel, Valid: true},
			NewName:       sql.NullString{String: "Zipped Mario", Valid: true},
		}

		err = RenameFile(tmpRoot, file)
		assert.NoError(t, err)

		newPath := filepath.Join(tmpRoot, "Zipped Mario.zip")
		assert.FileExists(t, newPath)
		assert.NoFileExists(t, containerPath)
	})

	t.Run("collision handling - move to duplicates", func(t *testing.T) {
		// Setup original file
		oldRel := "game.sfc"
		oldPath := filepath.Join(tmpRoot, oldRel)
		os.WriteFile(oldPath, []byte("content"), 0644)

		// Setup collision file
		collisionPath := filepath.Join(tmpRoot, "Collision.sfc")
		os.WriteFile(collisionPath, []byte("existing"), 0644)

		file := db.SessionFile{
			RelativePath: oldRel,
			NewName:      sql.NullString{String: "Collision", Valid: true},
		}

		err := RenameFile(tmpRoot, file)
		assert.NoError(t, err)

		// Should be in duplicados/.
		dupePath := filepath.Join(tmpRoot, "duplicados", "Collision.sfc")
		assert.FileExists(t, dupePath)
	})

	t.Run("missing new name error", func(t *testing.T) {
		file := db.SessionFile{
			RelativePath: "any.sfc",
			NewName:      sql.NullString{Valid: false},
		}
		err := RenameFile(tmpRoot, file)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "el archivo no tiene un nuevo nombre validado")
	})
}
