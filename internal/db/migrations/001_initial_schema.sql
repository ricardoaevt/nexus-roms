-- Configuración de la aplicación
CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sesiones de scraping
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    root_path TEXT NOT NULL,
    status TEXT NOT NULL, -- 'running', 'paused', 'completed', 'stopped'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Archivos dentro de una sesión
CREATE TABLE IF NOT EXISTS session_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    relative_path TEXT NOT NULL,
    filename TEXT NOT NULL,
    container_path TEXT, -- Ruta del archivo ZIP/RAR/7Z si está dentro de uno
    hash_md5 TEXT,
    hash_sha1 TEXT,
    hash_crc32 TEXT,
    name_metadata TEXT, -- Nombre obtenido de la API
    region_metadata TEXT,
    year_metadata TEXT,
    company_metadata TEXT,
    new_name TEXT,
    status TEXT NOT NULL, -- 'pending', 'hashing', 'scraping', 'found', 'not_found', 'renamed', 'error'
    error_message TEXT,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Credenciales de API (cifradas en la aplicación)
CREATE TABLE IF NOT EXISTS api_credentials (
    provider TEXT PRIMARY KEY, -- 'screenscraper', 'igdb', etc.
    username TEXT,
    password TEXT, -- Cifrado
    api_key TEXT,   -- Cifrado
    base_url TEXT,
    is_active INTEGER DEFAULT 1,
    search_by_hash INTEGER DEFAULT 1,
    search_by_name INTEGER DEFAULT 1
);
