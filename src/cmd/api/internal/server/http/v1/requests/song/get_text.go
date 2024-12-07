package song

import "github.com/google/uuid"

type GetTextRequest struct {
	ID       uuid.UUID `json:"id" validate:"required,uuid"`
	Page     uint      `query:"page" validate:"gte=0"`
	PageSize uint      `query:"pageSize" validate:"gte=0"`
}

func (r *GetTextRequest) SetDefaults() {
	if r.PageSize == 0 {
		r.PageSize = 1
	}
	if r.Page == 0 {
		r.Page = 0
	}
}
