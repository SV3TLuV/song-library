package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"song-library-api/src/pkg/music_info_client/model"
)

func main() {
	e := echo.New()

	e.GET("/info", func(c echo.Context) error {
		group := c.QueryParam("group")
		song := c.QueryParam("song")

		if group == "" || song == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Both 'group' and 'song' query parameters are required",
			})
		}

		if group == "500" {
			return echo.ErrInternalServerError
		}

		if group == "400" {
			return echo.ErrBadRequest
		}

		return c.JSON(http.StatusOK, model.SongDetail{
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		})
	})

	if err := e.Start("0.0.0.0:8081"); err != nil {
		log.Fatal(err)
	}
}
