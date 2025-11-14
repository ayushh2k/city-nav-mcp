package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcp-server/internal/dtos"
	"net/http"
	"net/url"
)

type GeoService struct {
	Client *http.Client
}

func NewGeoService(client *http.Client) *GeoService {
	return &GeoService{Client: client}
}

type nominatimResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func (s *GeoService) GetGeocode(params dtos.GeocodeRequestParams) (*dtos.GeocodeResponse, error) {
	log.Printf("Service: Calling external API for geocode: city=%s\n", params.City)

	baseURL, _ := url.Parse("https://nominatim.openstreetmap.org/search")
	query := baseURL.Query()
	query.Set("city", params.City)
	if params.CountryHint != "" {
		query.Set("country", params.CountryHint)
	}
	query.Set("format", "json")
	query.Set("limit", "1")
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, errors.New("failed to create geocode request")
	}

	req.Header.Set("User-Agent", "CityNavigator-MCP-Server (Go-http-client)")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error calling Nominatim: %v", err)
		return nil, errors.New("external API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Nominatim returned non-200 status: %s", resp.Status)
		return nil, fmt.Errorf("external API returned status: %s", resp.Status)
	}

	var results []nominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return nil, errors.New("failed to decode API response")
	}

	if len(results) == 0 {
		log.Printf("No results found for city: %s", params.City)
		return nil, errors.New("city not found")
	}

	topResult := results[0]
	var lat, lon float64
	fmt.Sscanf(topResult.Lat, "%f", &lat)
	fmt.Sscanf(topResult.Lon, "%f", &lon)

	return &dtos.GeocodeResponse{
		Lat:         lat,
		Lon:         lon,
		DisplayName: topResult.DisplayName,
	}, nil
}

func (s *GeoService) GetNearby(params dtos.NearbyRequestParams) (*dtos.NearbyResponse, error) {

	fmt.Printf("Service: Received nearby request (placeholder): %+v\n", params)

	return &dtos.NearbyResponse{
		Places: []dtos.NearbyPlace{
			{
				Name: fmt.Sprintf("Placeholder %s 1", params.Query),
				Lat:  params.Lat + 0.001,
				Lon:  params.Lon + 0.001,
				Tags: map[string]interface{}{"type": params.Query},
			},
		},
	}, nil
}
