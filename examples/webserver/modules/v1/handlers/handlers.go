package handlers

import (
	"net/http"

	"github.com/captjt/saddle/middleware"
	"github.com/captjt/saddle/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	v1mw "saddle/examples/webserver/modules/v1/middleware"
)

type (
	Handlers struct {
		logger    *logger.Logger
		validator *validator.Validate
		test      string // Purely to show passing configurations through to handlers.
	}
)

func New(
	logger *logger.Logger,
	validator *validator.Validate,
	test string,
) *Handlers {
	return &Handlers{
		validator: validator,
		logger:    logger,
		test:      test,
	}
}

func (h *Handlers) Route(e *fiber.App, basePath string) error {
	g := e.Group(basePath)

	// Chain middleware for all routes.
	g.Use(v1mw.APIData())
	g.Use(middleware.Validate(h.validator))

	// Define routes.
	g.Add(http.MethodPost, "/hello-world",
		h.helloWorld(),
	)
	return nil
}
