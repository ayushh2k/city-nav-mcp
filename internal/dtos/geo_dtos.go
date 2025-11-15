package dtos

type GeocodeResponse struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	DisplayName string  `json:"display_name"`
}

type NearbyPlace struct {
	Name string                 `json:"name"`
	Lat  float64                `json:"lat"`
	Lon  float64                `json:"lon"`
	Tags map[string]interface{} `json:"tags"`
}

type NearbyResponse struct {
	Places []NearbyPlace `json:"places"`
}

type GeocodeRequestParams struct {
	City        string `form:"city" binding:"required"`
	CountryHint string `form:"country_hint"`
}

type NearbyRequestParams struct {
	Lat     float64 `form:"lat" binding:"required"`
	Lon     float64 `form:"lon" binding:"required"`
	Query   string  `form:"query" binding:"required"`
	RadiusM int     `form:"radius_m" binding:"required,lte=5000"`
	Limit   int     `form:"limit" binding:"required,lte=15"`
}
