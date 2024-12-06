package song

import "github.com/google/uuid"

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}
