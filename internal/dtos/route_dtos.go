package dtos

type Point struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

type EtaRequest struct {
	Profile string  `json:"profile" binding:"required,oneof=foot bike car"`
	Points  []Point `json:"points" binding:"required,min=2"`
}

type EtaResponse struct {
	DistanceKm  float64 `json:"distance_km"`
	DurationMin float64 `json:"duration_min"`
	Polyline    string  `json:"polyline"`
}
