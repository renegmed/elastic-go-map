package main

import "project-match/web"

// http://localhost:8000/consumer/33324/producers/10/km
// http://localhost:8000/consumer/33324/producersmap/10/km
// http://localhost:8000/consumer/33324/allproducersmap
func main() {
	r := web.RegisterRoutes()
	r.Static("/public", "./public")
	r.Run(":8000")
}
