package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handlers) getHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}
