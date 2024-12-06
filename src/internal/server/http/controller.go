package http

import "github.com/labstack/echo/v4"

type SongController interface {
	GetList(echo.Context) error
	GetText(echo.Context) error
	Create(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}
