package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "producer services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Producer services")
	},
}
