package fixture_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const OSM_BASE_URL = "https://nominatim.openstreetmap.org"

type Result struct {
	PlaceID     int `json:"place_id"`
	Name        string
	DisplayName string `json:"display_name"`
}

type OpenStreetMap struct {
	BaseURL string
}

func (m OpenStreetMap) baseURL() string {
	if m.BaseURL != "" {
		return m.BaseURL
	}
	return OSM_BASE_URL
}

func (m OpenStreetMap) Search(q string) ([]Result, error) {
	url, err := url.Parse(m.BaseURL)
	if err != nil {
		return nil, err
	}
	url.Path = "/search"
	url.RawQuery = "q=" + q
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []Result
	json.Unmarshal(data, &results)
	return results, nil
}
