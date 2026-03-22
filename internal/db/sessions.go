package db

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        int64
	RootPath  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionFile struct {
	ID              int64
	SessionID       int64
	RelativePath    string
	Filename        string
	ContainerPath   sql.NullString
	HashMD5         sql.NullString
	HashSHA1        sql.NullString
	HashCRC32       sql.NullString
	NameMetadata    sql.NullString
	RegionMetadata  sql.NullString
	YearMetadata    sql.NullString
	CompanyMetadata sql.NullString
	NewName         sql.NullString
	Status          string
	ErrorMessage    sql.NullString
}

func (db *DB) CreateSession(rootPath string) (int64, error) {
	query := `INSERT INTO sessions (root_path, status) VALUES (?, ?)`
	res, err := db.Conn.Exec(query, rootPath, "running")
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) GetLatestSession() (*Session, error) {
	query := `SELECT id, root_path, status, created_at, updated_at FROM sessions ORDER BY id DESC LIMIT 1`
	row := db.Conn.QueryRow(query)

	var s Session
	err := row.Scan(&s.ID, &s.RootPath, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (db *DB) UpdateSessionStatus(id int64, status string) error {
	query := `UPDATE sessions SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Conn.Exec(query, status, id)
	return err
}

func (db *DB) AddFile(file SessionFile) (int64, error) {
	query := `
		INSERT INTO session_files (
			session_id, relative_path, filename, container_path, 
			status
		) VALUES (?, ?, ?, ?, ?)
	`
	res, err := db.Conn.Exec(query, 
		file.SessionID, file.RelativePath, file.Filename, file.ContainerPath, 
		"pending",
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) UpdateFileMetadata(file SessionFile) error {
	query := `
		UPDATE session_files SET
			hash_md5 = ?, hash_sha1 = ?, hash_crc32 = ?,
			name_metadata = ?, region_metadata = ?, year_metadata = ?, company_metadata = ?,
			new_name = ?, status = ?, error_message = ?
		WHERE id = ?
	`
	_, err := db.Conn.Exec(query,
		file.HashMD5, file.HashSHA1, file.HashCRC32,
		file.NameMetadata, file.RegionMetadata, file.YearMetadata, file.CompanyMetadata,
		file.NewName, file.Status, file.ErrorMessage,
		file.ID,
	)
	return err
}

func (db *DB) GetPendingFiles(sessionID int64) ([]SessionFile, error) {
	query := `SELECT id, session_id, relative_path, filename, container_path, status FROM session_files WHERE session_id = ? AND status NOT IN ('found', 'renamed')`
	rows, err := db.Conn.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []SessionFile
	for rows.Next() {
		var f SessionFile
		err := rows.Scan(&f.ID, &f.SessionID, &f.RelativePath, &f.Filename, &f.ContainerPath, &f.Status)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (db *DB) GetFileByID(id int64) (*SessionFile, error) {
	query := `SELECT id, session_id, relative_path, filename, container_path, new_name, status FROM session_files WHERE id = ?`
	row := db.Conn.QueryRow(query, id)

	var f SessionFile
	err := row.Scan(&f.ID, &f.SessionID, &f.RelativePath, &f.Filename, &f.ContainerPath, &f.NewName, &f.Status)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// GetSessionProgress devuelve (total, done) para una sesión
func (db *DB) GetSessionProgress(sessionID int64) (int, int) {
	var total, done int
	db.Conn.QueryRow(`SELECT COUNT(*) FROM session_files WHERE session_id = ?`, sessionID).Scan(&total)
	db.Conn.QueryRow(`SELECT COUNT(*) FROM session_files WHERE session_id = ? AND status IN ('found', 'renamed', 'not_found', 'error')`, sessionID).Scan(&done)
	return total, done
}
