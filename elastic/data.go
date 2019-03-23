package elastic

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Producer struct {
	Account   string   `json:"account"`
	Name      string   `json:"name"`
	Category  string   `json:"category"`
	StartDate string   `json:"startdate"`
	EndDate   string   `json:"enddate"`
	Location  Location `json:"location"`
	Quantity  float32  `json:"quantity"`
	Icon      string   `json:"icon"`
	Client    string   `json:"client"`
}

type Consumer struct {
	Account   string   `json:"account"`
	Name      string   `json:"name"`
	Category  string   `json:"category"`
	StartDate string   `json:"startdate"`
	EndDate   string   `json:"enddate"`
	Location  Location `json:"location"`
	Icon      string   `json:"icon"`
	Quantity  float32  `json:"quantity"`
}
