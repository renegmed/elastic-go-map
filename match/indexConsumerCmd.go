package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"project-match/elastic"

)

var indexConsumerCmd = &cobra.Command{
	Use:   "index",
	Short: "index consumer",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("index consumer")
		f, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		client, err := elastic.NewConsumerClient()
		err = client.IndexConsumerJsonFile(f)
		if err != nil {
			fmt.Errorf("%v\n", err)
		}
	},
}

func init() {
	consumerCmd.AddCommand(indexConsumerCmd)
	indexConsumerCmd.Flags().StringP("file", "f", "", "consumer json file for indexing")
}
