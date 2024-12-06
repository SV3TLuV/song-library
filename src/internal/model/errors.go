package model

import "github.com/pkg/errors"

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)
