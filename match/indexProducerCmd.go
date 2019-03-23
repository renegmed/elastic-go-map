package main

import (
	"fmt"
	"project-match/elastic"

	"github.com/spf13/cobra"
)

var indexProducerCmd = &cobra.Command{
	Use:   "index",
	Short: "index producer",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("index producer")
		f, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		client, err := elastic.NewProducerClient()
		err = client.IndexProducerJsonFile(f)
		if err != nil {
			fmt.Errorf("%v\n", err)
		}
	},
}

func init() {
	producerCmd.AddCommand(indexProducerCmd)
	indexProducerCmd.Flags().StringP("file", "f", "", "producer json file for indexing")
}
