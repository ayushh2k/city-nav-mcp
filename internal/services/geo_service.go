package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
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

type nominatimGeocodeResult struct {
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

	var results []nominatimGeocodeResult
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

type nominatimNearbyResult struct {
	Lat         string            `json:"lat"`
	Lon         string            `json:"lon"`
	DisplayName string            `json:"display_name"`
	Extratags   map[string]string `json:"extratags"`
}

func (s *GeoService) GetNearby(params dtos.NearbyRequestParams) (*dtos.NearbyResponse, error) {
	log.Printf("Service: Calling external API for nearby: q=%s\n", params.Query)

	lon1, lat1, lon2, lat2 := calculateBoundingBox(params.Lat, params.Lon, params.RadiusM)
	viewbox := fmt.Sprintf("%f,%f,%f,%f", lon1, lat1, lon2, lat2)

	baseURL, _ := url.Parse("https://nominatim.openstreetmap.org/search")
	query := baseURL.Query()
	query.Set("q", params.Query)
	query.Set("viewbox", viewbox)
	query.Set("bounded", "1")
	query.Set("format", "json")
	query.Set("extratags", "1")
	query.Set("limit", fmt.Sprintf("%d", params.Limit))
	baseURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, errors.New("failed to create nearby request")
	}
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server (Go-http-client)")

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error calling Nominatim (nearby): %v", err)
		return nil, errors.New("external API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Nominatim returned non-200 status: %s", resp.Status)
		return nil, fmt.Errorf("external API returned status: %s", resp.Status)
	}

	var results []nominatimNearbyResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return nil, errors.New("failed to decode API response")
	}

	places := make([]dtos.NearbyPlace, 0, len(results))
	for _, res := range results {
		var lat, lon float64
		fmt.Sscanf(res.Lat, "%f", &lat)
		fmt.Sscanf(res.Lon, "%f", &lon)

		tags := make(map[string]interface{}, len(res.Extratags))
		for k, v := range res.Extratags {
			tags[k] = v
		}

		places = append(places, dtos.NearbyPlace{
			Name: res.DisplayName,
			Lat:  lat,
			Lon:  lon,
			Tags: tags,
		})
	}

	return &dtos.NearbyResponse{Places: places}, nil
}

func calculateBoundingBox(lat, lon float64, radiusM int) (float64, float64, float64, float64) {
	const earthRadius = 6371000.0

	latDelta := float64(radiusM) / earthRadius * (180.0 / math.Pi)
	lonDelta := float64(radiusM) / (earthRadius * math.Cos(lat*(math.Pi/180.0))) * (180.0 / math.Pi)

	lat1 := lat - latDelta
	lon1 := lon - lonDelta

	lat2 := lat + latDelta
	lon2 := lon + lonDelta

	return lon1, lat1, lon2, lat2
}
