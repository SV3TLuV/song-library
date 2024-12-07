package v1

import (
	"github.com/google/uuid"
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
	songService  service.SongService
	groupService service.GroupService
}

func NewSongController(songService service.SongService) *songController {
	return &songController{songService: songService}
}

func (c *songController) GetList(ctx echo.Context) error {
	var request song.GetListRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	request.SetDefaults()

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	context := ctx.Request().Context()
	filters := &model.SongFilter{
		GroupID:     uuid.Nil,
		Song:        request.Song,
		Text:        request.Text,
		Link:        request.Link,
		ReleaseDate: request.ReleaseDate,
	}

	if request.Group != nil {
		group, err := c.groupService.GetByName(context, *request.Group)
		if err != nil {
			return err
		}
		filters.GroupID = group.ID
	}

	songs, err := c.songService.GetSongs(context, filters, request.Page, request.PageSize)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, model.PaginatedList[model.SongView]{
		Page:       songs.Page,
		PageSize:   songs.PageSize,
		TotalPages: songs.TotalPages,
		Items:      converter.ToViewsFromSong(songs.Items),
	})
}

func (c *songController) GetText(ctx echo.Context) error {
	var request song.GetTextRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	request.SetDefaults()

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return model.ErrBadRequest
	}
	request.ID = id

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	context := ctx.Request().Context()
	verses, err := c.songService.GetSongText(context, request.ID, request.Page, request.PageSize)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, verses)
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

	return ctx.JSON(http2.StatusOK, converter.ToViewFromSong(*createdSong))
}

func (c *songController) Update(ctx echo.Context) error {
	var request song.UpdateRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return model.ErrBadRequest
	}
	request.ID = id

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	song := model.Song{ID: request.ID}

	if request.Group != nil {
		song.Group = *request.Group
	}
	if request.Song != nil {
		song.Song = *request.Song
	}
	if request.Text != nil {
		song.Text = *request.Text
	}
	if request.Link != nil {
		song.Link = *request.Link
	}
	if request.ReleaseDate != nil {
		song.ReleaseDate = *request.ReleaseDate
	}

	context := ctx.Request().Context()
	entity, err := c.songService.Edit(context, song)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, converter.ToViewFromSong(*entity))
}

func (c *songController) Delete(ctx echo.Context) error {
	var request song.DeleteRequest
	if err := ctx.Bind(&request); err != nil {
		return model.ErrBadRequest
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return model.ErrBadRequest
	}
	request.ID = id

	if err := ctx.Validate(&request); err != nil {
		return model.ErrBadRequest
	}

	context := ctx.Request().Context()
	entity, err := c.songService.Delete(context, request.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, converter.ToViewFromSong(*entity))
}
