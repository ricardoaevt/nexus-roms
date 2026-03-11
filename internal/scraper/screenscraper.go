package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"romsRename/internal/db"
	"strings"
	"time"
)

const ssSoftName = "Nexus Roms v1.2.0"

type ScreenScraperClient struct {
	creds    *db.APICredentials
	http     *http.Client
	database *db.DB
}

func NewScreenScraperClient(creds *db.APICredentials, database *db.DB) *ScreenScraperClient {
	return &ScreenScraperClient{
		creds:    creds,
		database: database,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *ScreenScraperClient) Name() string {
	return "ScreenScraper"
}

func (s *ScreenScraperClient) CanSearchByHash() bool {
	return s.creds.SearchByHash
}

func (s *ScreenScraperClient) CanSearchByName() bool {
	return s.creds.SearchByName
}

// ssResponse es la estructura común de la respuesta de ScreenScraper
type ssResponse struct {
	Response struct {
		Jeu struct {
			Noms []struct {
				Text   string `json:"text"`
				Region string `json:"region"`
			} `json:"noms"`
			Dates []struct {
				Region string `json:"region"`
				Text   string `json:"text"`
			} `json:"dates"`
			Editeur struct {
				Text string `json:"text"`
			} `json:"editeur"`
			Developpeur struct {
				Text string `json:"text"`
			} `json:"developpeur"`
			Note struct {
				Text string `json:"text"`
			} `json:"note"`
			Genres []struct {
				Noms []struct {
					Langue string `json:"langue"`
					Text   string `json:"text"`
				} `json:"noms"`
			} `json:"genres"`
			NbJoueurs string `json:"nbjoueurs"`
			Rom       struct {
				RomRegions string `json:"romregions"`
				RomLangues string `json:"romlangues"`
				Beta       string `json:"beta"`
				Demo       string `json:"demo"`
				Proto      string `json:"proto"`
				Hack       string `json:"hack"`
			} `json:"rom"`
		} `json:"jeu"`
	} `json:"response"`
}

// parseMetadata extrae todos los campos disponibles de la respuesta de ScreenScraper
func parseSSMetadata(data ssResponse) *Metadata {
	jeu := data.Response.Jeu
	meta := &Metadata{}

	// --- Nombre: elegir el primero que tenga región válida ---
	for _, n := range jeu.Noms {
		if n.Region != "ss" && n.Region != "" {
			meta.Name = n.Text
			break
		}
	}
	if meta.Name == "" && len(jeu.Noms) > 0 {
		meta.Name = jeu.Noms[0].Text
	}

	// --- Región: preferir datos del rom concreto, luego los noms ---
	if jeu.Rom.RomRegions != "" {
		// La región del rom es la más precisa (ej: "us", "eu,us")
		parts := strings.Split(jeu.Rom.RomRegions, ",")
		var regions []string
		for _, p := range parts {
			r := strings.ToUpper(strings.TrimSpace(p))
			if r != "" && r != "SS" {
				regions = append(regions, r)
			}
		}
		meta.Region = strings.Join(regions, ",")
	} else {
		// Fallback: acumular regiones desde los noms
		seen := make(map[string]bool)
		var regions []string
		for _, n := range jeu.Noms {
			r := strings.ToUpper(n.Region)
			if r != "" && r != "SS" && !seen[r] {
				seen[r] = true
				regions = append(regions, r)
			}
		}
		meta.Region = strings.Join(regions, ",")
	}

	// --- Idiomas del rom concreto ---
	if jeu.Rom.RomLangues != "" {
		parts := strings.Split(jeu.Rom.RomLangues, ",")
		var langs []string
		for _, p := range parts {
			l := strings.ToUpper(strings.TrimSpace(p))
			if l != "" {
				langs = append(langs, l)
			}
		}
		meta.Languages = strings.Join(langs, ",")
	}

	// --- Año: tomar el primero disponible ---
	if len(jeu.Dates) > 0 {
		year := jeu.Dates[0].Text
		if len(year) >= 4 {
			meta.Year = year[:4] // Solo el año (YYYY)
		} else {
			meta.Year = year
		}
	}

	// --- Compañía (publisher) ---
	meta.Company = jeu.Editeur.Text

	// --- Desarrollador ---
	meta.Developer = jeu.Developpeur.Text

	// --- Género (en inglés preferido) ---
	if len(jeu.Genres) > 0 {
		for _, nom := range jeu.Genres[0].Noms {
			if nom.Langue == "en" {
				meta.Genre = nom.Text
				break
			}
		}
		// Fallback al primero si no hay en inglés
		if meta.Genre == "" && len(jeu.Genres[0].Noms) > 0 {
			meta.Genre = jeu.Genres[0].Noms[0].Text
		}
	}

	// --- Jugadores ---
	meta.Players = jeu.NbJoueurs

	// --- Nota media ---
	meta.Rating = jeu.Note.Text

	// --- Tipo de ROM ---
	rom := jeu.Rom
	switch {
	case rom.Beta == "1":
		meta.RomType = "Beta"
	case rom.Demo == "1":
		meta.RomType = "Demo"
	case rom.Proto == "1":
		meta.RomType = "Proto"
	case rom.Hack == "1":
		meta.RomType = "Hack"
	default:
		meta.RomType = ""
	}

	return meta
}

func (s *ScreenScraperClient) doRequest(ctx context.Context, params url.Values) ([]byte, error) {
	baseURL := "https://www.screenscraper.fr/api2/jeuInfos.php"
	if s.creds.BaseURL != "" {
		baseURL = s.creds.BaseURL
	}

	// devid y devpassword son las credenciales de desarrollador de ScreenScraper
	// Se leen desde la DB (Username = devid, Password = devpassword)
	if s.creds.Username == "" {
		return nil, fmt.Errorf("screenscraper: devid (Username) no configurado en Settings")
	}
	params.Set("devid", s.creds.Username)
	params.Set("devpassword", s.creds.Password)
	params.Set("softname", ssSoftName)
	params.Set("output", "json")

	// Si además tiene cuenta personal de usuario, incluirla para mayor cuota
	if s.creds.APIKey != "" {
		params.Set("ssid", s.creds.APIKey)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	var err error

	// API Tracking
	period := time.Now().Format("2006-01")
	cfgKey := "api_tracker_screenscraper"
	cfgValue := s.database.GetConfig(cfgKey, "{}")
	var trackerData struct {
		Period string `json:"Period"`
		Count  int    `json:"Count"`
	}
	_ = json.Unmarshal([]byte(cfgValue), &trackerData)

	if trackerData.Period != period {
		trackerData.Period = period
		trackerData.Count = 1
	} else {
		trackerData.Count++
	}
	newData, _ := json.Marshal(trackerData)
	err = s.database.SaveConfig(cfgKey, string(newData))
	if err != nil {
		fmt.Printf("[DEBUG] Error saving API tracker: %v\n", err)
	} else {
		fmt.Printf("[DEBUG] API Tracker updated for %s: count is now %d\n", cfgKey, trackerData.Count)
	}

	fmt.Printf("[DEBUG] Making API request to: %s\n", fullURL)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("acceso denegado (403): verifica tus credenciales o límites de cuota")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error de API: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	if len(body) > 0 && body[0] == '<' {
		return nil, fmt.Errorf("la API devolvió HTML en lugar de JSON. Servicio caído o bloqueado")
	}
	return body, nil
}

func (s *ScreenScraperClient) SearchByHash(ctx context.Context, query SearchQuery) (*Metadata, error) {
	if !s.creds.SearchByHash || query.HashMD5 == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Set("md5", query.HashMD5)

	body, err := s.doRequest(ctx, params)
	if err != nil || body == nil {
		return nil, err
	}

	var data ssResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error al procesar respuesta JSON: %v", err)
	}

	meta := parseSSMetadata(data)
	if meta.Name == "" {
		return nil, nil
	}
	return meta, nil
}

func (s *ScreenScraperClient) SearchByName(ctx context.Context, query SearchQuery) (*Metadata, error) {
	if !s.creds.SearchByName || query.Filename == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Set("romnom", query.Filename)

	body, err := s.doRequest(ctx, params)
	if err != nil || body == nil {
		return nil, err
	}

	var data ssResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error al procesar respuesta JSON: %v", err)
	}

	meta := parseSSMetadata(data)
	if meta.Name == "" {
		return nil, nil
	}
	return meta, nil
}
