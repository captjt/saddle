package saddle

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"

	"github.com/captjt/saddle/handlers"
	"github.com/captjt/saddle/middleware"
	log "github.com/captjt/saddle/pkg/logger"
)

type (
	// Service contains functions references attached to a service.
	Service interface {
		// Attach attaches the service to execute | expose.
		Attach(*echo.Echo, *log.Logger, trace.Tracer) (func(), error)
		// Config returns the configuration of a service.
		Config() any
		// Description returns the description of a service.
		Description() string
		// Name returns the name of a service.
		Name() string
	}

	// Project contains elements, functions and references attached to a project.
	Project[T Service] struct {
		// Echo contains the referenced Echo framework instance attached to the project.
		Echo *echo.Echo

		logger   *log.Logger
		service  T
		shutdown func(func())
		tracer   trace.Tracer
	}

	validate struct {
		validator *validator.Validate
	}
)

var (
	// compiledAt contains the datetime stamp representing when the service was built.
	compiledAt string
	// executedAt contains the datetime stamp representing when the service was executed.
	executedAt time.Time
	// gitBranch contains the GIT branch of the service.
	gitBranch string
	// gitCommit contains the GIT commit of the service.
	gitCommit string
	// version contains the version of the service.
	version string
)

func init() {
	executedAt = time.Now().UTC()
}

// Validate processes validation on a struct | model.
func (v *validate) Validate(i any) error {
	return v.validator.Struct(i)
}

// new instantiates a new project instance.
func new[T Service](service T, logger *log.Logger, tracer trace.Tracer) (
	*Project[T], error,
) {
	s := &Project[T]{
		Echo:    echo.New(),
		service: service,
	}

	s.Echo.HideBanner = true
	s.Echo.HidePort = true
	s.Echo.Validator = &validate{validator.New()}

	s.Echo.Use(otelecho.Middleware(service.Name(), otelecho.WithSkipper(handlers.Skipper)))
	s.Echo.Use(middleware.RequestID())

	// route harnesser-specific handlers ↴
	h := handlers.New(&handlers.Config{
		CompiledAt: compiledAt,
		ExecutedAt: executedAt,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		Version:    version,
	},
		logger,
		tracer,
	)
	h.Route(s.Echo, "")

	// construct safe shutdown monitor ↴
	s.shutdown = func(shutdown func()) {
		if shutdown != nil {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				defer logger.Sync() // flush any pending log(s) before exiting
				<-c
				logger.Info("initiating shutdown")
				shutdown()
				os.Exit(0)
			}()
		}
	}

	// attach service with service-specific safe shutdown ↴
	sd, err := s.service.Attach(s.Echo, logger, tracer)
	s.shutdown(sd)
	return s, err
}
