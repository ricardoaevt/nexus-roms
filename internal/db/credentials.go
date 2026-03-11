package db

import (
	"romsRename/internal/crypto"
)

type APICredentials struct {
	Provider     string
	Username     string
	Password     string
	APIKey       string
	BaseURL      string
	IsActive     bool
	SearchByHash bool
	SearchByName bool
}

func (db *DB) SaveCredentials(creds APICredentials) error {
	key, _ := crypto.GetKey()
	
	encPass := ""
	if creds.Password != "" {
		var err error
		encPass, err = crypto.Encrypt(creds.Password, key)
		if err != nil {
			return err
		}
	}

	encKey := ""
	if creds.APIKey != "" {
		var err error
		encKey, err = crypto.Encrypt(creds.APIKey, key)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO api_credentials (provider, username, password, api_key, base_url, is_active, search_by_hash, search_by_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(provider) DO UPDATE SET
			username=excluded.username,
			password=excluded.password,
			api_key=excluded.api_key,
			base_url=excluded.base_url,
			is_active=excluded.is_active,
			search_by_hash=excluded.search_by_hash,
			search_by_name=excluded.search_by_name
	`
	_, err := db.Conn.Exec(query, 
		creds.Provider, 
		creds.Username, 
		encPass, 
		encKey, 
		creds.BaseURL, 
		creds.IsActive, 
		creds.SearchByHash, 
		creds.SearchByName,
	)
	return err
}

func (db *DB) GetCredentials(provider string) (*APICredentials, error) {
	query := `SELECT provider, username, password, api_key, base_url, is_active, search_by_hash, search_by_name FROM api_credentials WHERE provider = ?`
	row := db.Conn.QueryRow(query, provider)

	var creds APICredentials
	var encPass, encKey string
	err := row.Scan(&creds.Provider, &creds.Username, &encPass, &encKey, &creds.BaseURL, &creds.IsActive, &creds.SearchByHash, &creds.SearchByName)
	if err != nil {
		return nil, err
	}

	key, _ := crypto.GetKey()
	if encPass != "" {
		creds.Password, _ = crypto.Decrypt(encPass, key)
	}
	if encKey != "" {
		creds.APIKey, _ = crypto.Decrypt(encKey, key)
	}

	return &creds, nil
}
