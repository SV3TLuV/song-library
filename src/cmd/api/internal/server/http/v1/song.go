package v1

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	http2 "net/http"
	"song-library-api/src/cmd/api/internal/converter"
	"song-library-api/src/cmd/api/internal/model"
	"song-library-api/src/cmd/api/internal/server/http/v1/requests/song"
	"song-library-api/src/cmd/api/internal/service"
	"time"
)

type SongController struct {
	songService  service.SongService
	groupService service.GroupService
}

func NewSongController(
	songService service.SongService,
	groupService service.GroupService) *SongController {
	return &SongController{
		songService:  songService,
		groupService: groupService,
	}
}

// GetList godoc
// @Summary      Get list of songs
// @Description  Retrieves a paginated list of songs with optional filters
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        group       query     string  false  "Filter by group name"
// @Param        song        query     string  false  "Filter by song name"
// @Param        text        query     string  false  "Filter by text"
// @Param        link        query     string  false  "Filter by link"
// @Param        releaseDate query     string  false  "Filter by release date (DD.MM.YYYY)"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        pageSize    query     int     false  "Page size (default: 5)"
// @Success      200         {object}  model.PaginatedList[model.SongView]
// @Failure      400         {object}  model.APIError  "Bad request"
// @Failure      500         {object}  model.APIError  "Internal server error"
// @Router       /songs [get]
func (c *SongController) GetList(ctx echo.Context) error {
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
		GroupID: uuid.Nil,
		Song:    request.Song,
		Text:    request.Text,
		Link:    request.Link,
	}

	if request.Group != nil {
		group, err := c.groupService.GetByName(context, *request.Group)
		if err != nil {
			return err
		}
		filters.GroupID = group.ID
	}

	if request.ReleaseDate != nil {
		releaseDate, err := time.Parse("02.01.2006", *request.ReleaseDate)
		if err != nil {
			return errors.Wrap(err, "failed to parse release date")
		}
		filters.ReleaseDate = &releaseDate
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

// GetText godoc
// @Summary      Get song text
// @Description  Retrieves song text by song ID with optional pagination for verses
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id        path      string  true   "Song ID"
// @Param        page      query     int     false  "Page number (default: 1)"
// @Param        pageSize  query     int     false  "Page size (default: 1)"
// @Success      200       {object}  []string
// @Failure      400       {object}  model.APIError  "Bad request"
// @Failure      404       {object}  model.APIError  "Song not found"
// @Failure      500       {object}  model.APIError  "Internal server error"
// @Router       /songs/{id}/text [get]
func (c *SongController) GetText(ctx echo.Context) error {
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

// Create godoc
// @Summary      Create a new song
// @Description  Adds a new song to the library
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song  body      song.CreateRequest  true  "Song data"
// @Success      200   {object}  model.SongView
// @Failure      400   {object}  model.APIError  "Bad request"
// @Failure      500   {object}  model.APIError  "Internal server error"
// @Router       /songs [post]
func (c *SongController) Create(ctx echo.Context) error {
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

// Update godoc
// @Summary      Update a song
// @Description  Updates the details of an existing song
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id    path      string          true   "Song ID"
// @Param        song  body      song.UpdateRequest true "Updated song data"
// @Success      200   {object}  model.SongView
// @Failure      400   {object}  model.APIError  "Bad request"
// @Failure      404   {object}  model.APIError  "Song not found"
// @Failure      500   {object}  model.APIError  "Internal server error"
// @Router       /songs/{id} [patch]
func (c *SongController) Update(ctx echo.Context) error {
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
		releaseDate, err := time.Parse("02.01.2006", *request.ReleaseDate)
		if err != nil {
			return errors.Wrap(err, "failed to parse release date")
		}
		song.ReleaseDate = releaseDate
	}

	context := ctx.Request().Context()
	entity, err := c.songService.Edit(context, song)
	if err != nil {
		return err
	}

	return ctx.JSON(http2.StatusOK, converter.ToViewFromSong(*entity))
}

// Delete godoc
// @Summary      Delete a song
// @Description  Deletes a song by its ID
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id    path      string  true   "Song ID"
// @Success      200   {object}  model.SongView
// @Failure      400   {object}  model.APIError  "Bad request"
// @Failure      404   {object}  model.APIError  "Song not found"
// @Failure      500   {object}  model.APIError  "Internal server error"
// @Router       /songs/{id} [delete]
func (c *SongController) Delete(ctx echo.Context) error {
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
