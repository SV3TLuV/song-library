package repository

import (
	"context"
	"github.com/google/uuid"
	"song-library-api/src/internal/model"
)

type SongRepository interface {
	GetSongs(ctx context.Context, filters *model.SongFilter, limit, offset uint) ([]model.SongWithGroup, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.SongWithGroup, error)
	GetByNameAndGroup(ctx context.Context, group, name string) (*model.SongWithGroup, error)
	Count(ctx context.Context, filters *model.SongFilter) (uint, error)
	Create(ctx context.Context, entity model.Song) (*model.Song, error)
	Update(ctx context.Context, entity model.Song) (*model.Song, error)
	Delete(ctx context.Context, song model.Song) error
}

type GroupRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	GetByName(ctx context.Context, name string) (*model.Group, error)
	Create(ctx context.Context, entity model.Group) (*model.Group, error)
	Update(ctx context.Context, entity model.Group) (*model.Group, error)
}
