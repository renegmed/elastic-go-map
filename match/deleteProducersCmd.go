package main

import (
	"fmt"
	"project-match/elastic"
	"github.com/spf13/cobra"
)
var index string

var deleteProducersCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete producer",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete producers")
		index, err := cmd.Flags().GetString("index")
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}
		client, err := elastic.NewProducerClient()
		err = client.DeleteIndex(index)
		if err != nil {
			fmt.Errorf("%v\n", err)
		}
		return
	},
}

func init() {
	producerCmd.AddCommand(deleteProducersCmd)
	//deleteProducersCmd.Flags().StringVarP(&index, "index", "i", "", "delete index")
	deleteProducersCmd.Flags().StringP("index", "i", "", "delete producer index")
}
