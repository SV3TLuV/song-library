package service

import (
	"context"
	"github.com/google/uuid"
	"song-library-api/src/internal/model"
)

type SongService interface {
	GetSongs(ctx context.Context, filters *model.SongFilter, page, pageSize uint) (*model.PaginatedList[model.Song], error)
	GetSongText(ctx context.Context, id uuid.UUID, page, pageSize uint) (*model.PaginatedList[string], error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Song, error)
	Add(ctx context.Context, song, group string) (*model.Song, error)
	Edit(ctx context.Context, song model.Song) (*model.Song, error)
	Delete(ctx context.Context, id uuid.UUID) (*model.Song, error)
}

type GroupService interface {
	GetByName(ctx context.Context, group string) (*model.Group, error)
}
