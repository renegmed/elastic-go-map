package main

import (
	"encoding/json"
	"fmt"
	"os"

	"project-match/elastic"

	"github.com/spf13/cobra"
)

// $ match consumer geodata

var createConsumerGeoDataCmd = &cobra.Command{
	Use:   "geodata",
	Short: "generate consumer json geodata",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate consumer json geodata")

		consClient, err := elastic.NewConsumerClient()
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}
		consumers, err := consClient.GetAllConsumers()
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}

		//fmt.Printf("consumers:\n%v\n", consumers)

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

		geoData := elastic.NewConsumerGeoData(consumers)

		fmt.Printf("%v\n", geoData)
		out, err := os.Create("geoData-Consumers.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer out.Close()

		encodeJSON := json.NewEncoder(out)
		err = encodeJSON.Encode(geoData)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	consumerCmd.AddCommand(createConsumerGeoDataCmd)
}
