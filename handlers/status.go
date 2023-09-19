package handlers

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/captjt/saddle/models"
)

var info *debug.BuildInfo

func init() {
	info, _ = debug.ReadBuildInfo()
}

func (h *handlers) getStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.StatusResponse{
			Version:    h.config.Version,
			CompiledAt: h.config.CompiledAt,
			ExecutedAt: h.config.ExecutedAt.Format(time.RFC3339),
			Uptime:     time.Now().UTC().Sub(h.config.ExecutedAt).String(),
			BuildInfo:  info,
		})
	}
}
