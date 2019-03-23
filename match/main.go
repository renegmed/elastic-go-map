/*

	curl -XGET 'localhost:9200/_cat/indices?v&pretty'

   $ match producer index file -f producer.json
   $ match producer delete index -i producer
   $ match producer search phrase -p soy

   $ match consumer index file -f consumer.json
   $ match consumer delete index -i consumer
   $ match consumer search phrase -p soy

   $ match producer consumer id -i 33324
   	phrase: soy bean  startdate: 2019-01-15  enddate: 2019-04-15

		name: producer-101, description: soy bean, start date: 2019-01-15, end date: 2019-03-15
		name: producer-102, description: soy bean, start date: 2019-02-01, end date: 2019-03-30

   $ match producer location consumer -i 33324 -d 2000 -u km

*/
package main

import "github.com/spf13/cobra"

var rootCmd *cobra.Command

func init() {

	rootCmd = &cobra.Command{
		Use:   "match",
		Short: "match producer or consumer",
	}
	rootCmd.AddCommand(producerCmd)
	rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(combinedGeoDataCmd)
}

func main() {
	rootCmd.Execute()
}
