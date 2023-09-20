package main

import (
	"fmt"
	"os"

	"github.com/captjt/saddle"
	"github.com/captjt/saddle/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"saddle/examples/webserver/models"
	v1 "saddle/examples/webserver/modules/v1"
)

type (
	extended struct {
		test string
	}

	Service struct {
		config      *models.Config
		description string
		echo        *echo.Echo
		logger      *logger.Logger
		name        string
		tracer      trace.Tracer
		extended
	}

	cleanup []func()
)

const (
	// description contains the description of the service.
	description = "Example webserver to leverage the scaffolding power of saddle!"
	// name contains the name of the service.
	name = "webserver"
)

var (
	_cleanup cleanup

	command *cobra.Command
	version string
)

func init() {

	command = saddle.New(version)
	webserver := saddle.Command(saddle.Instantiate(New(description, name)))
	command.AddCommand(webserver)

	// - define command-line parameters ↴
	webserver.PersistentFlags().StringP("environment", "e", "", "environment of service deployment")
	webserver.PersistentFlags().StringP("address", "a", "", "address | interface to listen for incoming requests")
	webserver.PersistentFlags().Set("environment", os.Getenv("ENVIRONMENT"))
	webserver.PersistentFlags().Set("address", os.Getenv("ADDRESS"))

	// - mark properties as required ↴
	webserver.MarkPersistentFlagRequired("environment")
	webserver.MarkPersistentFlagRequired("address")

	// - bind persistent command-line flags to viper so they are accessible within the service ↴
	_ = viper.BindPFlag(fmt.Sprintf("%s.%s", name, "environment"), webserver.PersistentFlags().Lookup("environment"))
	_ = viper.BindPFlag(fmt.Sprintf("%s.%s", name, "address"), webserver.PersistentFlags().Lookup("address"))
}

func main() {
	command.Execute()
}

func New(description, name string) *Service {
	return &Service{
		config:      &models.Config{},
		description: description,
		name:        name,
	}
}

func (s *Service) Attach(e *echo.Echo, logger *logger.Logger, tracer trace.Tracer) (
	func(), error,
) {
	s.echo = e
	s.tracer = tracer
	s.logger = logger

	if err := s.construct(); err != nil {
		return s.shutdown, err
	}
	return s.shutdown, nil
}

func (s *Service) construct() error {

	// Purely to show passing configurations through to handlers.
	s.extended.test = s.config.V1.Test

	// Instantiate and load v1 module: handlers, middleware, etc.
	if err := v1.New(
		s.logger,
		s.tracer,
		s.extended.test, // Purely to show passing configurations through to handlers.
	).Add(s.echo); err != nil {
		return err
	}
	s.logger.Info("webserver service module loaded",
		zap.String("module", "v1"),
	)

	return nil
}

func (s *Service) Config() any {
	return s.config
}

func (s *Service) Description() string {
	return s.description
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) shutdown() {
	for _, f := range _cleanup {
		f()
	}
	s.logger.Info("gracefully shutdown service")
}
