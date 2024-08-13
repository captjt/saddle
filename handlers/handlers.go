package handlers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"

	log "github.com/captjt/saddle/pkg/logger"
)

type (
	handlers struct {
		logger    *log.Logger
		validator *validator.Validate

		config *Config
	}

	Config struct {
		CompiledAt string
		ExecutedAt time.Time
		GitBranch  string
		GitCommit  string
		Version    string
	}
)

const (
	healthEndpointURI = "/healthz"
	statusEndpointURI = "/status"
)

// healthCheckRegex contains parts of a request URL used to bypass metrics during calls to a health check endpoint.
var healthCheckRegex *regexp.Regexp

func init() {
	healthCheckRegex = regexp.MustCompile("^(kube-probe|GoogleHC/)")
}

func New(config *Config, logger *log.Logger, validator *validator.Validate) *handlers {
	return &handlers{
		config:    config,
		logger:    logger,
		validator: validator,
	}
}

func (h *handlers) Route(e *fiber.App, basePath string) {
	g := e.Group(basePath)

	g.Add(http.MethodGet, healthEndpointURI, h.getHealth())
	g.Add(http.MethodGet, statusEndpointURI, h.getStatus())
}

// Skipper is used for specifying which route(s) should be opted out by the open-telemetry collector.
func Skipper(c echo.Context) bool {
	if c.Path() == healthEndpointURI || c.Path() == statusEndpointURI {
		return true
	}
	return healthCheckRegex.MatchString(c.Request().UserAgent())
}
