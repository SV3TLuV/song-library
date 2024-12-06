package song

type CreateRequest struct {
	Group string `body:"group" validate:"required,max=255"`
	Song  string `body:"song" validate:"required,max=255"`
}
