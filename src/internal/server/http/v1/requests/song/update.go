package song

import (
	"github.com/google/uuid"
	"time"
)

type UpdateRequest struct {
	ID          uuid.UUID  `json:"id" validate:"required,uuid"`
	Group       *string    `body:"group" validate:"max=255"`
	Song        *string    `body:"song" validate:"max=255"`
	Link        *string    `body:"link" validate:"max=2048"`
	Text        *string    `body:"text" validate:"max=2048"`
	ReleaseDate *time.Time `body:"releaseDate"`
}
