package main

import (
	"fmt"

	"project-match/elastic"

	"github.com/spf13/cobra"
)

// $ match producer consumer id -i 33324

var matchToProducersCmd = &cobra.Command{
	Use:   "consumer",
	Short: "match consumer to producers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("match consumer to producers")
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}
		fmt.Printf("match consumer to producers: %s\n", id)

		consClient, err := elastic.NewConsumerClient()
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}
		consumer, err := consClient.GetConsumer(id)
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}

		fmt.Printf("account: %s, name: %s, category: %s, start date: %s, end date: %s\n",
			consumer.Account, consumer.Name, consumer.Category, consumer.StartDate, consumer.EndDate)

		prodClient, err := elastic.NewProducerClient()

		//producers, err := prodClient.MatchedProducersWithinDateRange(

		producers, err := prodClient.MatchedProducersEndDateRange(
			consumer.Category,
			consumer.StartDate,
			consumer.EndDate)
		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Printf("producers:\n%v\n", producers)

		for _, producer := range producers {
			fmt.Printf(
				"account: %s, name: %s, category: %s, start date: %s, end date: %s location: lat %f, lon %f, quantity: %f, icon: %s, client: %s\n",
				producer.Account,
				producer.Name,
				producer.Category,
				producer.StartDate,
				producer.EndDate,
				producer.Location.Lat,
				producer.Location.Lon,
				producer.Quantity,
				producer.Icon,
				producer.Client,
			)
		}
	},
}

func init() {
	producerCmd.AddCommand(matchToProducersCmd)
	matchToProducersCmd.Flags().StringP("id", "i", "d", "consumer ID")
}
