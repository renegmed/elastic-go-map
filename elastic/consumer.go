package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	elastic "github.com/olivere/elastic"
)

const CONSUMER_INDEX = "consumer"
const CONSUMER_TYPE = "consumer"

type ConsumerClient struct {
	client *elastic.Client
}

func NewConsumerClient() (ConsumerClient, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return ConsumerClient{}, err
	}

	return ConsumerClient{
		client: client,
	}, nil
}

func (c *ConsumerClient) GetAllConsumers() ([]Consumer, error) {

	termQuery := elastic.NewMatchAllQuery()

	searchResult, err := c.client.Search().
		Index(CONSUMER_INDEX).
		Type(CONSUMER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("search error: %v\n", err)
		return []Consumer{}, err
	}

	list := []Consumer{}

	//fmt.Printf("---- consumer length: %d\n", len(searchResult.Hits.Hits))

	if len(searchResult.Hits.Hits) == 0 {
		return []Consumer{}, nil
	}

	for _, hit := range searchResult.Hits.Hits {
		var consumer Consumer
		err := json.Unmarshal(*hit.Source, &consumer)
		if err != nil {
			return []Consumer{}, err
		}
		//fmt.Printf("---- consumer: %v\n", consumer)
		list = append(list, consumer)
	}

	return list, nil
}

func (c *ConsumerClient) GetConsumer(account string) (Consumer, error) {

	termQuery := elastic.NewMultiMatchQuery(
		account, "account").Type("phrase_prefix")

	searchResult, err := c.client.Search().
		Index(CONSUMER_INDEX).
		Type(CONSUMER_TYPE).
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("search error: %v\n", err)
		return Consumer{}, err
	}

	list := []Consumer{}

	//fmt.Printf("---- consumer length: %d\n", len(searchResult.Hits.Hits))

	if len(searchResult.Hits.Hits) == 0 {
		return Consumer{}, nil
	}

	for _, hit := range searchResult.Hits.Hits {
		var consumer Consumer
		err := json.Unmarshal(*hit.Source, &consumer)
		if err != nil {
			return Consumer{}, err
		}
		//fmt.Printf("---- consumer: %v\n", consumer)
		list = append(list, consumer)
		break
	}

	return list[0], nil
}

func (c *ConsumerClient) DeleteIndex(index string) error {
	ctx := context.Background()

	exists, err := c.client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		deleteIndex, err := c.client.DeleteIndex(index).Do(ctx)
		if err != nil {
			return err
		}
		if !deleteIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return nil
}

func (c *ConsumerClient) Search(phrase string) ([]Consumer, error) {

	termQuery := elastic.NewMultiMatchQuery(
		phrase, "category").Type("phrase_prefix")

	searchResult, err := c.client.Search().
		Index(CONSUMER_INDEX). // name of Index (dev / prod)
		Type(CONSUMER_TYPE).   // type of Index
		Query(termQuery).
		Do(context.Background())

	if err != nil {
		fmt.Printf("search error: %v\n", err)
		return nil, err
	}

	list := []Consumer{}

	for _, hit := range searchResult.Hits.Hits {
		var consumer Consumer
		err := json.Unmarshal(*hit.Source, &consumer)
		if err != nil {
			return nil, err
		}

		list = append(list, consumer)
	}

	return list, nil
}

func (c *ConsumerClient) IndexConsumerJsonFile(file string) error {

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	consumers := []Consumer{}
	err = json.Unmarshal(byteValue, &consumers)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, consumer := range consumers {
		fmt.Printf("account: %s name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
			consumer.Account,
			consumer.Name,
			consumer.Category,
			consumer.StartDate,
			consumer.EndDate,
			consumer.Location.Lat,
			consumer.Location.Lon,
			consumer.Quantity,
			consumer.Icon,
		)
	}

	return indexConsumer(c, consumers)
}

func indexConsumer(c *ConsumerClient, consumers []Consumer) error {
	ctx := context.Background()
	exists, err := c.client.IndexExists(CONSUMER_INDEX).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Printf("---- index '%s' doesn't exists. will create new index\n ", CONSUMER_INDEX)
		createIndex, err := c.client.CreateIndex(CONSUMER_INDEX).Body(ConsumerMapping).Do(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	fmt.Printf("---- consumers length: %d\n", len(consumers))

	for _, consumer := range consumers {
		addConsumerToIndex(c.client, consumer)
	}

	return nil
}

func addConsumerToIndex(client *elastic.Client, c Consumer) error {
	fmt.Println("---- adding new consumer to index -----")

	_, err := client.Index().
		Index(CONSUMER_INDEX).
		Type(CONSUMER_TYPE).
		//Id(id).
		BodyJson(c).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
