package song

import (
	"time"
)

type GetListRequest struct {
	Group       *string    `query:"group"`
	Song        *string    `query:"song"`
	Text        *string    `query:"text"`
	Link        *string    `query:"link"`
	ReleaseDate *time.Time `query:"releaseDate"`
	Page        uint       `query:"page" validate:"gte=0"`
	PageSize    uint       `query:"pageSize" validate:"gte=0"`
}

func (r *GetListRequest) SetDefaults() {
	if r.PageSize == 0 {
		r.PageSize = 5
	}
	if r.Page == 0 {
		r.Page = 1
	}
}
