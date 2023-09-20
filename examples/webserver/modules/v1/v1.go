package v1

import (
	"github.com/captjt/saddle/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"

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
	exporter trace.Tracer,
	test string, // Purely to show passing configurations through to handlers.
) *v1 {
	return &v1{
		handlers: handlers.New(
			logger,
			exporter,
			test,
		),
	}
}

func (v *v1) Add(e *echo.Echo) error {
	return v.handlers.Route(e, basePath)
}
