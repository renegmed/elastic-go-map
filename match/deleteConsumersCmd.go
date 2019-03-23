package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"project-match/elastic"
)

var deleteConsumersCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete consumers",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete consumers")
		index, err := cmd.Flags().GetString("index")
		if err != nil {
			fmt.Errorf("%v\n", err)
			return
		}
		client, err := elastic.NewConsumerClient()
		err = client.DeleteIndex(index)
		if err != nil {
			fmt.Errorf("%v\n", err)
		}
		return
	},
}

func init() {
	consumerCmd.AddCommand(deleteConsumersCmd)
	deleteConsumersCmd.Flags().StringP("index", "i", "", "delete consumer index")
}
