package db



func (db *DB) GetConfig(key string, defaultValue string) string {
	query := `SELECT value FROM config WHERE key = ?`
	row := db.Conn.QueryRow(query, key)
	var value string
	err := row.Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

func (db *DB) SaveConfig(key, value string) error {
	query := `INSERT INTO config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
	          ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP`
	_, err := db.Conn.Exec(query, key, value)
	return err
}
