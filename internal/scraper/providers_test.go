package scraper_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"romsRename/internal/db"
	"romsRename/internal/scraper"
	"testing"

	"github.com/stretchr/testify/assert"
)

const inMemoryDB = ":memory:"

func TestScreenScraperClient(t *testing.T) {
	database, _ := db.InitDB(inMemoryDB)
	defer database.Close()

	t.Run("successful search", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "nexus", r.URL.Query().Get("devid"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"response":{"jeu":{"noms":[{"text":"Super Mario World","region":"us"}],"rom":{"romregions":"us"},"editeur":{"text":"Nintendo"}}}}`)
		}))
		defer server.Close()

		creds := &db.APICredentials{
			Username:     "nexus",
			BaseURL:      server.URL,
			SearchByHash: true,
		}
		client := scraper.NewScreenScraperClient(creds, database)
		
		meta, err := client.SearchByHash(context.Background(), scraper.SearchQuery{HashMD5: "123"})
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, "Super Mario World", meta.Name)
	})

	t.Run("successful SearchByName", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Mario", r.URL.Query().Get("romnom"))
			assert.Equal(t, "Nexus Roms v1.2.0", r.URL.Query().Get("softname"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"response":{"jeu":{"noms":[{"text":"Super Mario World","region":"us"}],"rom":{"romregions":"us"},"editeur":{"text":"Nintendo"}}}}`)
		}))
		defer server.Close()

		creds := &db.APICredentials{
			Username:     "nexus",
			BaseURL:      server.URL,
			SearchByName: true,
		}
		client := scraper.NewScreenScraperClient(creds, database)
		
		meta, err := client.SearchByName(context.Background(), scraper.SearchQuery{Filename: "Mario"})
		assert.NoError(t, err)
		assert.Equal(t, "Super Mario World", meta.Name)
	})
}

func TestTheGamesDBClient(t *testing.T) {
	database, _ := db.InitDB(inMemoryDB)
	defer database.Close()

	t.Run("successful search", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "key123", r.URL.Query().Get("apikey"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data":{"games":[{"game_title":"Sonic the Hedgehog","release_date":"1991-06-23"}]}}`)
		}))
		defer server.Close()

		creds := &db.APICredentials{
			APIKey:       "key123",
			BaseURL:      server.URL,
			SearchByHash: true,
		}
		client := scraper.NewTheGamesDBClient(creds, database)
		
		meta, err := client.SearchByHash(context.Background(), scraper.SearchQuery{HashMD5: "123"})
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, "Sonic the Hedgehog", meta.Name)
		assert.Equal(t, "1991", meta.Year)
	})

	t.Run("successful SearchByName", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Sonic", r.URL.Query().Get("name"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data":{"games":[{"game_title":"Sonic the Hedgehog","release_date":"1991-06-23"}]}}`)
		}))
		defer server.Close()

		creds := &db.APICredentials{
			APIKey:       "key123",
			BaseURL:      server.URL,
			SearchByName: true,
		}
		client := scraper.NewTheGamesDBClient(creds, database)
		
		meta, err := client.SearchByName(context.Background(), scraper.SearchQuery{Filename: "Sonic"})
		assert.NoError(t, err)
		assert.Equal(t, "Sonic the Hedgehog", meta.Name)
	})
}
