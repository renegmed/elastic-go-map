package main

import (
	"fmt"
	"project-match/elastic"

	"github.com/spf13/cobra"
)

var searchProducersCmd = &cobra.Command{
	Use:   "search",
	Short: "search producers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("search producers")
		phrase, err := cmd.Flags().GetString("phrase")

		fmt.Printf("search producers phrase: %s\n", phrase)

		client, err := elastic.NewProducerClient()
		producers, err := client.Search(phrase)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("producers:\n%v\n", producers)

		for _, producer := range producers {
			fmt.Printf("----- account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s, client: %s\n",
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
	producerCmd.AddCommand(searchProducersCmd)
	searchProducersCmd.Flags().StringP("phrase", "p", "", "phrase to search")
}
