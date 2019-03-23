package elastic

type Properties struct {
	Account   string  `json:"account"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	StartDate string  `json:"startdate"`
	EndDate   string  `json:"enddate"`
	Quantity  float32 `json:"quantity"`
	Icon      string  `json:"icon"`
	Client    string  `json:"client"`
}
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type Feature struct {
	Geometry   Geometry   `json:"geometry"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type GeoData struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func NewProducerGeoData(producers []Producer) GeoData {
	geoData := GeoData{}
	geoData.Type = "FeatureCollection"

	for _, producer := range producers {
		feature := Feature{}
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, producer.Location.Lon)
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, producer.Location.Lat)
		feature.Type = "Feature"
		feature.Properties = Properties{
			Account:   producer.Account,
			Name:      producer.Name,
			Category:  producer.Category,
			StartDate: producer.StartDate,
			EndDate:   producer.EndDate,
			Quantity:  producer.Quantity,
			Icon:      producer.Icon,
			Client:    producer.Client,
		}
		geoData.Features = append(geoData.Features, feature)
	}

	return geoData
}

func NewConsumerGeoData(consumers []Consumer) GeoData {
	geoData := GeoData{}
	geoData.Type = "FeatureCollection"

	for _, consumer := range consumers {
		feature := Feature{}
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, consumer.Location.Lon)
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, consumer.Location.Lat)
		feature.Type = "Feature"
		feature.Properties = Properties{
			Account:   consumer.Account,
			Name:      consumer.Name,
			Category:  consumer.Category,
			StartDate: consumer.StartDate,
			EndDate:   consumer.EndDate,
			Quantity:  consumer.Quantity,
			Icon:      consumer.Icon,
		}
		geoData.Features = append(geoData.Features, feature)
	}

	return geoData
}

func NewConsumerProducersGeoData(consumers []Consumer, producers []Producer) GeoData {
	geoData := GeoData{}
	geoData.Type = "FeatureCollection"

	geoData = consumersGeoData(geoData, consumers)
	geoData = producersGeoData(geoData, producers)

	return geoData
}

func consumersGeoData(geoData GeoData, consumers []Consumer) GeoData {
	for _, consumer := range consumers {
		feature := Feature{}
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, consumer.Location.Lon)
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, consumer.Location.Lat)
		feature.Type = "Feature"
		feature.Properties = Properties{
			Account:   consumer.Account,
			Name:      consumer.Name,
			Category:  consumer.Category,
			StartDate: consumer.StartDate,
			EndDate:   consumer.EndDate,
			Quantity:  consumer.Quantity,
			Icon:      consumer.Icon,
			Client:    "",
		}
		geoData.Features = append(geoData.Features, feature)
	}

	return geoData
}

func producersGeoData(geoData GeoData, producers []Producer) GeoData {
	for _, producer := range producers {
		feature := Feature{}
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, producer.Location.Lon)
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, producer.Location.Lat)
		feature.Type = "Feature"
		feature.Properties = Properties{
			Account:   producer.Account,
			Name:      producer.Name,
			Category:  producer.Category,
			StartDate: producer.StartDate,
			EndDate:   producer.EndDate,
			Quantity:  producer.Quantity,
			Icon:      producer.Icon,
			Client:    producer.Client,
		}
		geoData.Features = append(geoData.Features, feature)
	}

	return geoData
}
