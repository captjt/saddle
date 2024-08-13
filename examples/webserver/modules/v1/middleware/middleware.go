package middleware

import (
	"saddle/examples/webserver/modules/v1/models"

	"github.com/gofiber/fiber/v2"
)

func APIData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch c.Path() {
		case "/v1/hello-world":
			c.Locals("model", new(models.HelloWorldRequest))
		}
		return c.Next()
	}
}
