package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"saddle/examples/webserver/modules/v1/models"
)

func (h *Handlers) helloWorld() fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := c.Locals("request").(*models.HelloWorldRequest)

		h.logger.Info("retrieving application embedding",
			zap.String("message", in.Message),
		)

		out := models.Response{
			Message: "Hello there, saddle up friend!",
		}

		return c.JSON(out)
	}
}
