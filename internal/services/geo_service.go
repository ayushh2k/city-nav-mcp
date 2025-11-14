package services

import (
	"errors"
	"fmt"
	"mcp-server/internal/dtos"
)

type GeoService struct{}

func NewGeoService() *GeoService {
	return &GeoService{}
}

func (s *GeoService) GetGeocode(params dtos.GeocodeRequestParams) (*dtos.GeocodeResponse, error) {

	fmt.Printf("Service: Received geocode request: city=%s, country_hint=%s\n", params.City, params.CountryHint)

	if params.City == "kyoto" {
		return &dtos.GeocodeResponse{
			Lat:         35.0116,
			Lon:         135.7681,
			DisplayName: "Kyoto, Kyoto Prefecture, Japan",
		}, nil
	}

	return nil, errors.New("city not found (placeholder)")
}

func (s *GeoService) GetNearby(params dtos.NearbyRequestParams) (*dtos.NearbyResponse, error) {

	fmt.Printf("Service: Received nearby request: %+v\n", params)

	return &dtos.NearbyResponse{
		Places: []dtos.NearbyPlace{
			{
				Name: fmt.Sprintf("Placeholder %s 1", params.Query),
				Lat:  params.Lat + 0.001,
				Lon:  params.Lon + 0.001,
				Tags: map[string]interface{}{"type": params.Query},
			},
			{
				Name: fmt.Sprintf("Placeholder %s 2", params.Query),
				Lat:  params.Lat - 0.001,
				Lon:  params.Lon - 0.001,
				Tags: map[string]interface{}{"type": params.Query},
			},
		},
	}, nil
}
