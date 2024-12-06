package service

import (
	"context"
	"github.com/google/uuid"
	"song-library-api/src/internal/model"
)

type SongService interface {
	GetSongs(ctx context.Context, filters *model.SongFilter, page, pageSize uint) (*model.PaginatedList[model.SongWithGroup], error)
	GetSongText(ctx context.Context, id uuid.UUID, page, pageSize uint) (*model.PaginatedList[string], error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.SongWithGroup, error)
	Add(ctx context.Context, song, group string) (*model.SongWithGroup, error)
	Edit(ctx context.Context, id uuid.UUID, song, group string) (*model.SongWithGroup, error)
	Delete(ctx context.Context, song model.SongWithGroup) error
}
