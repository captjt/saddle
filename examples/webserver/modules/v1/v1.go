package v1

import (
	"github.com/captjt/saddle/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"saddle/examples/webserver/modules/v1/handlers"
)

type (
	v1 struct {
		handlers *handlers.Handlers
	}
)

const basePath = "/v1"

func New(
	logger *logger.Logger,
	validator *validator.Validate,
	test string, // Purely to show passing configurations through to handlers.
) *v1 {
	return &v1{
		handlers: handlers.New(
			logger,
			validator,
			test,
		),
	}
}

func (v *v1) Add(e *fiber.App) error {
	return v.handlers.Route(e, basePath)
}
