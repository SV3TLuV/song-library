package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"song-library-api/src/cmd/api/internal/model"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		status := http.StatusInternalServerError
		if err != nil {
			switch {
			case errors.Is(err, model.ErrBadRequest), errors.Is(err, echo.ErrBadRequest):
				status = http.StatusBadRequest
			case errors.Is(err, model.ErrNotFound):
				status = http.StatusNotFound
			}

			return c.JSON(status, model.APIError{Message: err.Error()})
		}

		return nil
	}
}
