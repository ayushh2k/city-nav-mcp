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

type WikidataService struct {
	Client *http.Client
}

func NewWikidataService(client *http.Client) *WikidataService {
	return &WikidataService{Client: client}
}

func (s *WikidataService) Query(sparqlQuery string) (*dtos.WikidataResponse, error) {
	log.Println("Service: Calling Wikidata SPARQL API")

	endpoint := "https://query.wikidata.org/sparql"

	data := url.Values{}
	data.Set("query", sparqlQuery)
	body := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, errors.New("failed to create SPARQL request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "CityNavigator-MCP-Server (Go-http-client)")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.New("external SPARQL API call failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wikidata API returned status: %s", resp.Status)
	}

	var result dtos.WikidataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New("failed to decode SPARQL response")
	}

	return &result, nil
}
