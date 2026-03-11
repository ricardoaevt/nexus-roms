package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"romsRename/internal/db"
	"time"
)

type TheGamesDBClient struct {
	creds    *db.APICredentials
	http     *http.Client
	database *db.DB
}

func NewTheGamesDBClient(creds *db.APICredentials, database *db.DB) *TheGamesDBClient {
	return &TheGamesDBClient{
		creds:    creds,
		database: database,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *TheGamesDBClient) Name() string {
	return "TheGamesDB"
}

func (s *TheGamesDBClient) CanSearchByHash() bool {
	return s.creds.SearchByHash
}

func (s *TheGamesDBClient) CanSearchByName() bool {
	return s.creds.SearchByName
}

func (s *TheGamesDBClient) doRequest(ctx context.Context, apiPath string, params url.Values) ([]byte, error) {
	baseURL := "https://api.thegamesdb.net/v1"
	if s.creds.BaseURL != "" {
		baseURL = s.creds.BaseURL
	}

	params.Set("apikey", s.creds.APIKey)
	fullURL := fmt.Sprintf("%s/%s?%s", baseURL, apiPath, params.Encode())

	// Tracking
	period := time.Now().Format("2006-01")
	cfgKey := "api_tracker_thegamesdb"
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
	s.database.SaveConfig(cfgKey, string(newData))

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TGDB error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (s *TheGamesDBClient) SearchByHash(ctx context.Context, query SearchQuery) (*Metadata, error) {
	if s.creds.APIKey == "" || query.HashMD5 == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Add("hash", query.HashMD5)

	body, err := s.doRequest(ctx, "Games/ByHash", params)
	if err != nil || body == nil {
		return nil, err
	}

	var res struct {
		Data struct {
			Games []struct {
				GameTitle   string `json:"game_title"`
				ReleaseDate string `json:"release_date"`
			} `json:"games"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if len(res.Data.Games) == 0 {
		return nil, nil
	}

	game := res.Data.Games[0]
	year := ""
	if len(game.ReleaseDate) >= 4 {
		year = game.ReleaseDate[:4]
	}

	return &Metadata{
		Name: game.GameTitle,
		Year: year,
	}, nil
}

func (s *TheGamesDBClient) SearchByName(ctx context.Context, query SearchQuery) (*Metadata, error) {
	if s.creds.APIKey == "" || query.Filename == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Add("name", query.Filename)

	body, err := s.doRequest(ctx, "Games/ByGameName", params)
	if err != nil || body == nil {
		return nil, err
	}

	var res struct {
		Data struct {
			Games []struct {
				GameTitle   string `json:"game_title"`
				ReleaseDate string `json:"release_date"`
			} `json:"games"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if len(res.Data.Games) == 0 {
		return nil, nil
	}

	game := res.Data.Games[0]
	year := ""
	if len(game.ReleaseDate) >= 4 {
		year = game.ReleaseDate[:4]
	}

	return &Metadata{
		Name: game.GameTitle,
		Year: year,
	}, nil
}
