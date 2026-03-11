package scraper

import "context"

type Metadata struct {
	Name      string
	Region    string
	Languages string // ej: "En,Es"
	Year      string
	Company   string
	Developer string
	Genre     string
	Players   string
	Rating    string // nota media de la comunidad
	RomType   string // "rom", "beta", "demo", "hack", "proto"
	HashMD5   string
}

type SearchQuery struct {
	Filename string
	HashMD5  string
	HashSHA1 string
	HashCRC32 string
}

type Scraper interface {
	Name() string
	SearchByHash(ctx context.Context, query SearchQuery) (*Metadata, error)
	SearchByName(ctx context.Context, query SearchQuery) (*Metadata, error)
	CanSearchByHash() bool
	CanSearchByName() bool
}
