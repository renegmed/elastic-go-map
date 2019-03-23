package main

import (
	"fmt"
	"project-match/elastic"
	"strconv"

	"github.com/spf13/cobra"
)

// $ clea

var locateProducersCmd = &cobra.Command{
	Use:   "location",
	Short: "location of producers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("consumer")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		// fmt.Printf("    id: %s\n", id)

		distance, err := cmd.Flags().GetString("distance")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		numDistance, err := strconv.ParseInt(distance, 10, 64)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		// fmt.Printf("   distance: %s\n", distance)

		unit, err := cmd.Flags().GetString("unit")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		// fmt.Printf("   unit: %s\n", unit)

		consClient, err := elastic.NewConsumerClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		consumer, err := consClient.GetConsumer(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		// fmt.Printf("account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f\n",
		// 	consumer.Account,
		// 	consumer.Name,
		// 	consumer.Category,
		// 	consumer.StartDate,
		// 	consumer.EndDate,
		// 	consumer.Location.Lat,
		// 	consumer.Location.Lon,
		// 	consumer.Quantity,
		// )

		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		producers, err := prodClient.LocateProducers(consumer, numDistance, unit)
		if err != nil {
			fmt.Println(err)
			return
		}
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
	producerCmd.AddCommand(locateProducersCmd)
	locateProducersCmd.Flags().StringP("consumer", "i", "", "consumer ID")
	locateProducersCmd.Flags().StringP("distance", "d", "", "maximum distance from consumer")
	locateProducersCmd.Flags().StringP("unit", "u", "", "unit of distance, e.g. km, mi")
}
