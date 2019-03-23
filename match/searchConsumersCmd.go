package main

import (
	"fmt"

	"project-match/elastic"

	"github.com/spf13/cobra"
)

var searchConsumersCmd = &cobra.Command{
	Use:   "search",
	Short: "search consumers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		phrase, err := cmd.Flags().GetString("phrase")

		fmt.Printf("search consumers phrase: %s\n", phrase)

		client, err := elastic.NewConsumerClient()
		consumers, err := client.Search(phrase)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("consumers:\n%v\n", consumers)

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
	},
}

func init() {
	consumerCmd.AddCommand(searchConsumersCmd)
	searchConsumersCmd.Flags().StringP("phrase", "p", "", "phrase to search")
}
