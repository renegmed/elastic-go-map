package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	elastic "github.com/olivere/elastic"
)

const PRODUCER_INDEX = "producer"
const PRODUCER_TYPE = "producer"

type ProducerClient struct {
	client *elastic.Client
}

func NewProducerClient() (ProducerClient, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return ProducerClient{}, err
	}

	return ProducerClient{
		client: client,
	}, nil
}

// this works
// func (p *ProducerClient) MatchedProducers(phrase, startDate, endDate string) ([]Producer, error) {
// 	url := "http://localhost:9200/producer/producer/_search"
// 	query := []byte(`{
//             "query":{
// 				"range" : {
// 					"startdate" : { "gte": "2019-01-01", "format": "yyyy-MM-dd" }
// 				}
//             }
//             }`)
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("\n%s", string(body))

// 	return nil, nil
// }

// MatchedProducersEndDateRange() selects producers with end date within the consumers date range
func (p *ProducerClient) MatchedProducersEndDateRange(phrase, startDate, endDate string) ([]Producer, error) {

	fmt.Printf("==== MatchedProducersEndDateRange producers\n   phrase: %s  startdate: %s  enddate: %s\n", phrase, startDate, endDate)

	phraseQuery := elastic.NewMatchQuery("category", phrase)

	dateQuery := elastic.NewRangeQuery("enddate").
		Gte(startDate).
		Lte(endDate).
		Boost(2).
		Format("yyyy-MM-dd")

	termQuery := elastic.NewBoolQuery().Must(phraseQuery).Must(dateQuery)

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("matching error: %v\n", err)
		return nil, err
	}

	fmt.Printf("==== length: %d\n", len(searchResult.Hits.Hits))

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}

		list = append(list, producer)
	}

	return list, nil
}

// locateProducers() selects producers near the consumer
func (p *ProducerClient) LocateProducers(consumer Consumer, distance int64, unit string) ([]Producer, error) {

	phraseQuery := elastic.NewMatchQuery("category", consumer.Category)

	fmt.Println("--- locate producers \n", fmt.Sprintf("%f, %f", consumer.Location.Lat, consumer.Location.Lon))

	locationQuery := elastic.NewGeoDistanceQuery("location")
	locationQuery = locationQuery.Lat(consumer.Location.Lat)
	locationQuery = locationQuery.Lon(consumer.Location.Lon)
	locationQuery = locationQuery.Distance(fmt.Sprintf("%d%s", distance, unit))
	locationQuery = locationQuery.DistanceType("plane")

	termQuery := elastic.NewBoolQuery().Must(phraseQuery).Must(locationQuery)

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("matching error: %v\n", err)
		return nil, err
	}

	fmt.Printf("==== length: %d\n", len(searchResult.Hits.Hits))

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}
		list = append(list, producer)
	}

	return list, nil
}

// LocateProducersWithinDates() selects producers within consumer date range given the maximum distance
func (p *ProducerClient) LocateProducersWithinDates(consumer Consumer, distance int64, unit string) ([]Producer, error) {

	phraseQuery := elastic.NewMatchQuery("category", consumer.Category)

	//fmt.Println("--- locate producers \n", fmt.Sprintf("%f, %f", consumer.Location.Lat, consumer.Location.Lon))

	endDateQuery := elastic.NewRangeQuery("enddate").
		Gte(consumer.StartDate).
		Lte(consumer.EndDate).
		Boost(2).
		Format("yyyy-MM-dd")

	locationQuery := elastic.NewGeoDistanceQuery("location")
	locationQuery = locationQuery.Lat(consumer.Location.Lat)
	locationQuery = locationQuery.Lon(consumer.Location.Lon)
	locationQuery = locationQuery.Distance(fmt.Sprintf("%d%s", distance, unit))
	locationQuery = locationQuery.DistanceType("plane")

	termQuery := elastic.NewBoolQuery().Must(phraseQuery).Must(endDateQuery).Must(locationQuery)

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("matching error: %v\n", err)
		return nil, err
	}

	fmt.Printf("==== length: %d\n", len(searchResult.Hits.Hits))

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}
		list = append(list, producer)
	}

	return list, nil
}

// AllProducersWithinDateRange() selects producers within consumer date range regardless of distance
func (p *ProducerClient) AllProducersWithinDateRange(consumer Consumer) ([]Producer, error) {

	phraseQuery := elastic.NewMatchQuery("category", consumer.Category)

	endDateQuery := elastic.NewRangeQuery("enddate").
		Gte(consumer.StartDate).
		Lte(consumer.EndDate).
		Boost(2).
		Format("yyyy-MM-dd")

	termQuery := elastic.NewBoolQuery().Must(phraseQuery).Must(endDateQuery)

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("matching error: %v\n", err)
		return nil, err
	}

	//fmt.Printf("==== length: %d\n", len(searchResult.Hits.Hits))

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}
		list = append(list, producer)
	}

	return list, nil
}

// Matched
func (p *ProducerClient) MatchedProducersWithinDateRange(phrase, startDate, endDate string) ([]Producer, error) {

	//fmt.Printf("==== MatchedProducersWithinDateRange producers\n    phrase: %s  startdate: %s  enddate: %s\n", phrase, startDate, endDate)

	phraseQuery := elastic.NewMatchQuery("category", phrase)

	// dateQuery := elastic.NewRangeQuery("startdate").
	// Gte(startDate).
	// Lte(endDate).
	// Boost(2).
	// Format("yyyy-MM-dd")

	startDateQuery := elastic.NewRangeQuery("startdate").
		Gte(startDate).
		Boost(2).
		Format("yyyy-MM-dd")

	endDateQuery := elastic.NewRangeQuery("enddate").
		Lte(endDate).
		Boost(2).
		Format("yyyy-MM-dd")
	termQuery := elastic.NewBoolQuery().Must(phraseQuery).Must(startDateQuery).Must(endDateQuery)

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("matching error: %v\n", err)
		return nil, err
	}

	fmt.Printf("==== length: %d\n", len(searchResult.Hits.Hits))

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}

		list = append(list, producer)
	}

	return list, nil
}

func (p *ProducerClient) DeleteIndex(index string) error {
	ctx := context.Background()

	exists, err := p.client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		// Delete an index.
		deleteIndex, err := p.client.DeleteIndex(index).Do(ctx)
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return nil
}

func (p *ProducerClient) Search(phrase string) ([]Producer, error) {

	// var termQuery elastic.Query

	// words := strings.Split(strings.Trim(phrase, " "), " ")
	// //fmt.Printf("++++ words: %v\n", words)

	// tQuery := elastic.NewBoolQuery()

	// for _, word := range words {
	// 	termQuery = tQuery.Must(elastic.NewTermQuery("content", word))
	// }

	// termQuery := elastic.NewMultiMatchQuery(
	// 	phrase, "name", "category", "startdate", "enddate").Fuzziness("AUTO:2,5")

	// searchResult, err := p.client.Search().
	// 	Index(PRODUCER_INDEX). // search in index "producer"
	// 	Query(termQuery).      // specify the query
	// 	//Sort("topic.keyword", true). // sort by "topic" field, ascending
	// 	From(0).Size(2000).      // take documents 0-2000
	// 	Pretty(true).            // pretty print request and response JSON
	// 	Do(context.Background()) // execute

	termQuery := elastic.NewMultiMatchQuery(
		phrase, "category").Type("phrase_prefix")

	searchResult, err := p.client.Search().
		Index(PRODUCER_INDEX). // name of Index (dev / prod)
		Type(PRODUCER_TYPE).   // type of Index
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("search error: %v\n", err)
		return nil, err
	}

	list := []Producer{}

	for _, hit := range searchResult.Hits.Hits {
		var producer Producer
		err := json.Unmarshal(*hit.Source, &producer)
		if err != nil {
			return nil, err
		}

		list = append(list, producer)
	}

	return list, nil
}

func (p *ProducerClient) IndexProducerJsonFile(file string) error {

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	producers := []Producer{}
	err = json.Unmarshal(byteValue, &producers)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// for _, producer := range producers {
	// 	fmt.Printf("name: %s, description: %s, start date: %s, end date: %s\n",
	// 		producer.Name, producer.Description, producer.StartDate, producer.EndDate)
	// }

	return indexProducer(p, producers)
}

func indexProducer(p *ProducerClient, producers []Producer) error {
	ctx := context.Background()
	exists, err := p.client.IndexExists(PRODUCER_INDEX).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := p.client.CreateIndex(PRODUCER_INDEX).Body(ProducerMapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	for _, producer := range producers {
		addProducerToIndex(p.client, producer)
	}

	return nil
}

func addProducerToIndex(client *elastic.Client, p Producer) error {
	_, err := client.Index().
		Index(PRODUCER_INDEX).
		Type(PRODUCER_TYPE).
		//Id(id).
		BodyJson(p).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
