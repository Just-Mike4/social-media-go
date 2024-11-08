package main

import (
	"social-media-api/config"
	"social-media-api/routes"
)

func main() {
	config.ConnectDatabase()
	r := routes.SetupRouter()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Run(":8080")
}
