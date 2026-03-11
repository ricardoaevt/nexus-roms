package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001_initial_schema.sql
var initialSchema string

type DB struct {
	Conn *sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	// Asegurar que el directorio existe
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("error al crear directorio de base de datos: %w", err)
		}
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir base de datos: %w", err)
	}

	// Habilitar claves foráneas
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("error al habilitar foreign_keys: %w", err)
	}

	database := &DB{Conn: conn}

	if err := database.migrate(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *DB) migrate() error {
	_, err := db.Conn.Exec(initialSchema)
	if err != nil {
		return fmt.Errorf("error al ejecutar migración inicial: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
