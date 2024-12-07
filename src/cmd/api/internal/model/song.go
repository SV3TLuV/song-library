package model

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID          uuid.UUID
	GroupID     uuid.UUID
	Group       string
	Song        string
	Text        string
	Link        string
	ReleaseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SongView struct {
	ID          uuid.UUID
	Group       string
	Song        string
	Text        string
	Link        string
	ReleaseDate string
	CreatedAt   string
	UpdatedAt   string
}

type SongFilter struct {
	GroupID     uuid.UUID
	Song        *string
	Text        *string
	Link        *string
	ReleaseDate *time.Time
}
