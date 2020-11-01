package mapmodels

type Results struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

type Result struct {
	AddressComponents []Address `json:"address_components"`
	FormattedAddress  string    `json:"formatted_address"`
	Geometry          Geometry  `json:"geometry"`
	PlaceId           string    `json:"place_id"`
	Types             []string  `json:"types"`
}
