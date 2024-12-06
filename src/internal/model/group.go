package model

import (
	"github.com/google/uuid"
	"time"
)

type Group struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
