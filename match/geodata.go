package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var geodataCmd = &cobra.Command{
	Use:   "geodata",
	Short: "geodata services",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Geo data services")
	},
}
