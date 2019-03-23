package web

import (
	"fmt"
	"net/http"
	"project-match/elastic"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GanttBar struct {
	Id         string
	Name       string
	StartYear  int
	StartMonth int
	StartDay   int
	EndYear    int
	EndMonth   int
	EndDay     int
	IsHeader   bool
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

func RegisterRoutes() *gin.Engine {

	r := gin.Default()

	r.LoadHTMLGlob("templates/**/*.html")

	r.GET("/store", func(c *gin.Context) {
		c.HTML(http.StatusOK, "store-locator.html", nil)

	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)

	})

	r.GET("/geodata/consumer/:id", func(c *gin.Context) {
		id := c.Param("id")
		//fmt.Printf("    id: %s\n", id)

		consClient, err := elastic.NewConsumerClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - consumer client creation failed.")
			return
		}

		consumer, err := consClient.GetConsumer(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - consumer creation failed.")
			return
		}

		consumers := []elastic.Consumer{}
		consumers = append(consumers, consumer)

		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - producer client creation failed.")
			return
		}

		producers, err := prodClient.LocateProducers(consumer, 100000, "km")
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - producers creation failed.")
			return
		}

		geoData := elastic.NewConsumerProducersGeoData(consumers, producers)

		//fmt.Printf("    geoData: %v\n", geoData)

		c.JSON(http.StatusOK, geoData)

	})

	r.POST("/", func(c *gin.Context) {
		account := c.PostForm("account")
		display := c.PostForm("btn-display")

		fmt.Printf("  account: %s, display: %s\n", account, display)

		consumer, err := getConsumer(account)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - Consumer not Found")
			return
		}
		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating new Producer Client")
			return
		}

		if strings.Contains(display, "chart") {
			producers, err := prodClient.LocateProducersWithinDates(consumer, 1000000, "km")
			if err != nil {
				fmt.Println(err)
				c.String(http.StatusNotFound, "404 - Problem creating chart")
				return
			}
			// for _, producer := range producers {
			// 	fmt.Printf("----- account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f\n",
			// 		producer.Account,
			// 		producer.Name,
			// 		producer.Category,
			// 		producer.StartDate,
			// 		producer.EndDate,
			// 		producer.Location.Lat,
			// 		producer.Location.Lon,
			// 		producer.Quantity,
			// 	)
			// }
			ganttBars, err := ganttData(producers, consumer)
			if err != nil {
				fmt.Println(err)
				c.String(http.StatusNotFound, "404 - Problem creating gantt chart bars")
				return
			}
			consumer.Category = strings.Title(consumer.Category)

			params := map[string]interface{}{
				"consumer":  consumer,
				"producers": producers,
				"ganttBars": ganttBars,
				"isChart":   true,
			}
			c.HTML(http.StatusOK, "index.html", params)
		} else { // display map
			producers, err := prodClient.LocateProducersWithinDates(consumer, 1000000, "km")

			if err != nil {
				fmt.Println(err)
				c.String(http.StatusNotFound, "404 - Problem producers map")
				return
			}
			// for _, producer := range producers {
			// 	fmt.Printf("----- LocateProducersWithinDates() account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f\n",
			// 		producer.Account,
			// 		producer.Name,
			// 		producer.Category,
			// 		producer.StartDate,
			// 		producer.EndDate,
			// 		producer.Location.Lat,
			// 		producer.Location.Lon,
			// 		producer.Quantity,
			// 	)
			// }
			producers = append(producers, elastic.Producer{
				Account:   consumer.Account,
				Name:      consumer.Name,
				Category:  consumer.Category,
				StartDate: consumer.StartDate,
				EndDate:   consumer.EndDate,
				Location:  elastic.Location{Lat: consumer.Location.Lat, Lon: consumer.Location.Lon},
				Quantity:  consumer.Quantity,
				Icon:      consumer.Icon,
			})

			params := map[string]interface{}{
				"consumer":  consumer,
				"producers": producers,
				"isMap":     true,
				"apiKey":    GOOGLE_MAP_API_KEY,
			}
			c.HTML(http.StatusOK, "index.html", params)
		}

	})

	r.GET("/consumer/:id/producers/:distance/:unit", func(c *gin.Context) {
		id := c.Param("id")
		distance := c.Param("distance")
		unit := c.Param("unit")

		numDistance, err := strconv.ParseInt(distance, 10, 64)
		if err != nil {
			c.String(http.StatusNotFound, "404 - Invalid distance number - "+distance)
			return
		}

		//fmt.Printf("++++ Customer ID: %s, Distance: %d, Unit: %s\n", id, numDistance, unit)

		consumer, err := getConsumer(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - Consumer not Found")
			return
		}
		consumer.Category = strings.Title(consumer.Category)

		// fmt.Printf("account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 	consumer.Account,
		// 	consumer.Name,
		// 	consumer.Category,
		// 	consumer.StartDate,
		// 	consumer.EndDate,
		// 	consumer.Location.Lat,
		// 	consumer.Location.Lon,
		// 	consumer.Quantity,
		// 	consumer.Icon,
		// )

		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating product client")
			return
		}

		//producers, err := prodClient.LocateProducers(consumer, numDistance, unit)
		producers, err := prodClient.LocateProducersWithinDates(consumer, numDistance, unit)

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating producers.")
			return
		}
		// for _, producer := range producers {
		// 	fmt.Printf("----- account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 		producer.Account,
		// 		producer.Name,
		// 		producer.Category,
		// 		producer.StartDate,
		// 		producer.EndDate,
		// 		producer.Location.Lat,
		// 		producer.Location.Lon,
		// 		producer.Quantity,
		// 		producer.Icon,
		// 	)
		// }
		ganttBars, err := ganttData(producers, consumer)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating gantt chart bars")
			return
		}

		params := map[string]interface{}{
			"consumer":  consumer,
			"producers": producers,
			"ganttBars": ganttBars,
		}
		c.HTML(http.StatusOK, "consumer.html", params)
	})

	r.GET("/consumer/:id/producersmap/:distance/:unit", func(c *gin.Context) {
		id := c.Param("id")
		distance := c.Param("distance")
		unit := c.Param("unit")

		numDistance, err := strconv.ParseInt(distance, 10, 64)
		if err != nil {
			c.String(http.StatusNotFound, "404 - Invalid distance number - "+distance)
			return
		}

		//fmt.Printf("++++ Customer ID: %s, Distance: %d, Unit: %s\n", id, numDistance, unit)

		consumer, err := getConsumer(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - Consumer not Found")
			return
		}
		consumer.Category = strings.Title(consumer.Category)

		// fmt.Printf("account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 	consumer.Account,
		// 	consumer.Name,
		// 	consumer.Category,
		// 	consumer.StartDate,
		// 	consumer.EndDate,
		// 	consumer.Location.Lat,
		// 	consumer.Location.Lon,
		// 	consumer.Quantity,
		// 	consumer.Icon,
		// )

		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating producer client")
			return
		}

		//producers, err := prodClient.LocateProducers(consumer, numDistance, unit)
		producers, err := prodClient.LocateProducersWithinDates(consumer, numDistance, unit)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - Problem creating producers.")
			return
		}
		// for _, producer := range producers {
		// 	fmt.Printf(
		// 		"----- LocateProducersWithinDates() account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 		producer.Account,
		// 		producer.Name,
		// 		producer.Category,
		// 		producer.StartDate,
		// 		producer.EndDate,
		// 		producer.Location.Lat,
		// 		producer.Location.Lon,
		// 		producer.Quantity,
		// 		producer.Icon,
		// 	)
		// }
		producers = append(producers, elastic.Producer{
			Account:   consumer.Account,
			Name:      consumer.Name,
			Category:  consumer.Category,
			StartDate: consumer.StartDate,
			EndDate:   consumer.EndDate,
			Location:  elastic.Location{Lat: consumer.Location.Lat, Lon: consumer.Location.Lon},
			Quantity:  consumer.Quantity,
			Icon:      consumer.Icon,
		})

		params := map[string]interface{}{
			"consumer":  consumer,
			"producers": producers,
			"apiKey":    GOOGLE_MAP_API_KEY,
		}
		c.HTML(http.StatusOK, "producersmap.html", params)
	})

	r.GET("/consumer/:id/allproducersmap", func(c *gin.Context) {
		id := c.Param("id")

		consumer, err := getConsumer(id)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusNotFound, "404 - Consumer not Found")
			return
		}
		consumer.Category = strings.Title(consumer.Category) // capitalize the first character

		// fmt.Printf("account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 	consumer.Account,
		// 	consumer.Name,
		// 	consumer.Category,
		// 	consumer.StartDate,
		// 	consumer.EndDate,
		// 	consumer.Location.Lat,
		// 	consumer.Location.Lon,
		// 	consumer.Quantity,
		// 	consumer.Icon,
		// )

		prodClient, err := elastic.NewProducerClient()
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - problem creating producer client.")
			return
		}

		producers, err := prodClient.AllProducersWithinDateRange(consumer)

		if err != nil {
			fmt.Println(err)
			c.String(http.StatusNotFound, "404 - problem creating producer.")
			return
		}

		producers = append(producers, elastic.Producer{
			Account:   consumer.Account,
			Name:      consumer.Name,
			Category:  consumer.Category,
			StartDate: consumer.StartDate,
			EndDate:   consumer.EndDate,
			Location:  elastic.Location{Lat: consumer.Location.Lat, Lon: consumer.Location.Lon},
			Quantity:  consumer.Quantity,
			Icon:      consumer.Icon,
		})

		// for _, producer := range producers {
		// 	fmt.Printf(
		// 		"----- AllProducersWithinDateRange account: %s, name: %s, category: %s, start date: %s, end date: %s, location: lat %f, lon %f, quantity: %f, icon: %s\n",
		// 		producer.Account,
		// 		producer.Name,
		// 		producer.Category,
		// 		producer.StartDate,
		// 		producer.EndDate,
		// 		producer.Location.Lat,
		// 		producer.Location.Lon,
		// 		producer.Quantity,
		// 		producer.Icon,
		// 	)
		// }

		params := map[string]interface{}{
			"consumer":  consumer,
			"producers": producers,
			"apiKey":    GOOGLE_MAP_API_KEY,
		}
		c.HTML(http.StatusOK, "producersmap.html", params)
	})

	return r
}

func getConsumer(id string) (elastic.Consumer, error) {
	consClient, err := elastic.NewConsumerClient()
	if err != nil {
		return elastic.Consumer{}, err
	}
	consumer, err := consClient.GetConsumer(id)
	if err != nil {
		return elastic.Consumer{}, err
	}
	return consumer, nil
}

func ganttData(producers []elastic.Producer, consumer elastic.Consumer) ([]GanttBar, error) {
	bars := []GanttBar{}

	var cStartDate = strings.Split(consumer.StartDate, "-")
	var cEndDate = strings.Split(consumer.EndDate, "-")
	cStartYear, _ := strconv.Atoi(cStartDate[0])
	cStartMonth, _ := strconv.Atoi(cStartDate[1])
	cStartDay, _ := strconv.Atoi(cStartDate[2])
	cEndYear, _ := strconv.Atoi(cEndDate[0])
	cEndMonth, _ := strconv.Atoi(cEndDate[1])
	cEndDay, _ := strconv.Atoi(cEndDate[2])

	bars = append(bars, GanttBar{
		Id:       "888888",
		Name:     "CLIENT:",
		IsHeader: true,
	})

	bars = append(bars, GanttBar{
		Id:         consumer.Account,
		Name:       consumer.Name,
		StartYear:  cStartYear,
		StartMonth: cStartMonth - 1,
		StartDay:   cStartDay,
		EndYear:    cEndYear,
		EndMonth:   cEndMonth - 1,
		EndDay:     cEndDay,
		IsHeader:   false,
	})

	bars = append(bars, GanttBar{
		Id:       "999999",
		Name:     "SUPPLIERS:",
		IsHeader: true,
	})

	for _, producer := range producers {

		var startDate = strings.Split(producer.StartDate, "-")
		var endDate = strings.Split(producer.EndDate, "-")

		startYear, _ := strconv.Atoi(startDate[0])
		startMonth, _ := strconv.Atoi(startDate[1])
		startDay, _ := strconv.Atoi(startDate[2])
		endYear, _ := strconv.Atoi(endDate[0])
		endMonth, _ := strconv.Atoi(endDate[1])
		endDay, _ := strconv.Atoi(endDate[2])
		bars = append(bars, GanttBar{
			Id:         producer.Account,
			Name:       producer.Name,
			StartYear:  startYear,
			StartMonth: startMonth - 1,
			StartDay:   startDay,
			EndYear:    endYear,
			EndMonth:   endMonth - 1,
			EndDay:     endDay,
			IsHeader:   false,
		})
	}

	fmt.Println(bars)

	return bars, nil
}
