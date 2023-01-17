package main

import (
	"file_processing_service/config"
	"file_processing_service/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
