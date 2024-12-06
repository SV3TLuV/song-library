package api

import (
	"log"
	"song-library-api/src/cmd/api/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
