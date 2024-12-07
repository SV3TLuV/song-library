package main

import (
	"log"
	_ "song-library-api/src/cmd/api/docs"
	"song-library-api/src/cmd/api/internal/app"
)

// @title Song Library API
// @version 1.0
// @description Song Library API
// @BasePath /
func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
