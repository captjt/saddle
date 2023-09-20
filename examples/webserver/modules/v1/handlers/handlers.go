package handlers

import (
	"net/http"

	"github.com/captjt/saddle/middleware"
	"github.com/captjt/saddle/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"

	v1mw "saddle/examples/webserver/modules/v1/middleware"
)

type (
	Handlers struct {
		tracer trace.Tracer
		logger *logger.Logger
		test   string // Purely to show passing configurations through to handlers.
	}
)

func New(
	logger *logger.Logger,
	tracer trace.Tracer,
	test string,
) *Handlers {
	return &Handlers{
		tracer: tracer,
		logger: logger,
		test:   test,
	}
}

func (h *Handlers) Route(e *echo.Echo, basePath string) error {
	g := e.Group(basePath)

	g.Add(http.MethodPost, "/hello-world", h.helloWorld(),
		middleware.RequestID(),
		middleware.RequestLog(h.logger),
		v1mw.APIData(),
		middleware.Validate(h.logger),
	)
	return nil
}
