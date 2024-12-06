package route

import (
	"github.com/labstack/echo/v4"
	"song-library-api/src/internal/server/http"
)

func InitSongRoutes(group *echo.Group, controller http.SongController) {
	g := group.Group("/songs")

	g.GET("", controller.GetList)
	g.GET("/:id/text", controller.GetText)
	g.POST("/new", controller.Create)
	g.PATCH("/:id/edit", controller.Update)
	g.DELETE("/:id/", controller.Delete)
}
