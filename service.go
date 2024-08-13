package saddle

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/captjt/saddle/handlers"
	"github.com/captjt/saddle/middleware"
	log "github.com/captjt/saddle/pkg/logger"
)

type (
	// Service contains functions references attached to a service.
	Service interface {
		// Attach attaches the service to execute | expose.
		Attach(*fiber.App, *log.Logger, *validator.Validate) (func(), error)
		// Config returns the configuration of a service.
		Config() any
		// Description returns the description of a service.
		Description() string
		// Name returns the name of a service.
		Name() string
		// Validator returns the validator of a service.
		Validator() *validator.Validate
	}

	// Project contains elements, functions and references attached to a project.
	Project[T Service] struct {
		// App contains the referenced Fiber framework app instance attached to the project.
		App       *fiber.App
		validator *validator.Validate

		service  T
		shutdown func(func())
	}

	Validate struct {
		Validator *validator.Validate
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

// new instantiates a new project instance.
func new[T Service](service T, logger *log.Logger, validator *validator.Validate) (
	*Project[T], error,
) {

	s := &Project[T]{
		App:       fiber.New(),
		service:   service,
		validator: validator,
	}

	s.App.Use(middleware.RequestID())
	s.App.Use(middleware.RequestLog(logger))

	// route saddle-specific handlers ↴
	h := handlers.New(&handlers.Config{
		CompiledAt: compiledAt,
		ExecutedAt: executedAt,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		Version:    version,
	},
		logger,
		s.validator,
	)
	h.Route(s.App, "")

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
	sd, err := s.service.Attach(s.App, logger, s.validator)
	s.shutdown(sd)
	return s, err
}
