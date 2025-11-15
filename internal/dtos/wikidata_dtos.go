package dtos

type WikidataRequest struct {
	Query string `json:"query" binding:"required"`
}

type WikidataResponse struct {
	Results struct {
		Bindings []map[string]interface{} `json:"bindings"`
	} `json:"results"`
}
