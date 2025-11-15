package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mcp-server/internal/dtos"
	"net/http"
	"net/url"
	"strings"
)

type RouteService struct {
	Client *http.Client
}

func NewRouteService(client *http.Client) *RouteService {
	return &RouteService{Client: client}
}

type osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry string  `json:"geometry"`
}

type osrmResponse struct {
	Routes []osrmRoute `json:"routes"`
}

func (s *RouteService) GetEta(params dtos.EtaRequest) (*dtos.EtaResponse, error) {
	log.Printf("Service: Calling OSRM API for profile: %s", params.Profile)

	var coordStrings []string
	for _, p := range params.Points {
		coordStrings = append(coordStrings, fmt.Sprintf("%f,%f", p.Lon, p.Lat))
	}
	coordinates := strings.Join(coordStrings, ";")

	path := fmt.Sprintf("/route/v1/%s/%s", params.Profile, url.PathEscape(coordinates))
	baseURL, _ := url.Parse("https://router.project-osrm.org" + path)

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, errors.New("failed to create OSRM request")
	}
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("external OSRM API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSRM API returned status: %s", resp.Status)
	}

	var data osrmResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.New("failed to decode OSRM response")
	}

	if len(data.Routes) == 0 {
		return nil, errors.New("no route found")
	}

	topRoute := data.Routes[0]
	return &dtos.EtaResponse{
		DistanceKm:  topRoute.Distance / 1000.0,
		DurationMin: topRoute.Duration / 60.0,
		Polyline:    topRoute.Geometry,
	}, nil
}
