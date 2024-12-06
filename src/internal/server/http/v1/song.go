package v1

import (
	"github.com/labstack/echo/v4"
	http2 "net/http"
	"song-library-api/src/internal/converter"
	"song-library-api/src/internal/model"
	"song-library-api/src/internal/server/http"
	"song-library-api/src/internal/server/http/v1/requests/song"
	"song-library-api/src/internal/service"
)

var _ http.SongController = (*songController)(nil)

type songController struct {
	songService service.SongService
}

func NewSongController(songService service.SongService) *songController {
	return &songController{songService: songService}
}

func (c *songController) GetList(ctx echo.Context) error {

}

func (c *songController) GetText(ctx echo.Context) error {

}

func (c *songController) Create(ctx echo.Context) error {
	var request song.CreateRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	context := ctx.Request().Context()
	createdSong, err := c.songService.Add(context, request.Song, request.Group)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK)
}

func (c *songController) Update(ctx echo.Context) error {

}

func (c *songController) Delete(ctx echo.Context) error {
	var request song.DeleteRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	context := ctx.Request().Context()
	entity, err := c.songService.GetByID(context, request.ID)
	if err != nil {
		return err
	}

	if err = c.songService.Delete(context, *entity); err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, converter.ToViewFromSong(*entity))
}
