package handlers

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/captjt/saddle/models"
)

var info *debug.BuildInfo

func init() {
	info, _ = debug.ReadBuildInfo()
}

func (h *handlers) getStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(models.StatusResponse{
			Version:    h.config.Version,
			CompiledAt: h.config.CompiledAt,
			ExecutedAt: h.config.ExecutedAt.Format(time.RFC3339),
			Uptime:     time.Now().UTC().Sub(h.config.ExecutedAt).String(),
			BuildInfo:  info,
		})
	}
}
