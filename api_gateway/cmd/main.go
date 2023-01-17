package main

import (
	"api_gateway/config"
	"api_gateway/internal/app"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
