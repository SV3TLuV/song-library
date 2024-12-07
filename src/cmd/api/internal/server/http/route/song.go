package route

import (
	"github.com/labstack/echo/v4"
	v1 "song-library-api/src/cmd/api/internal/server/http/v1"
)

func InitSongRoutes(group *echo.Group, controller *v1.SongController) {
	g := group.Group("/songs")

	g.GET("", controller.GetList)
	g.GET("/:id/text", controller.GetText)
	g.POST("", controller.Create)
	g.PATCH("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}
