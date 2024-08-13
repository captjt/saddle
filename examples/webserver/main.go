package main

import (
	"fmt"
	"os"

	"github.com/captjt/saddle"
	"github.com/captjt/saddle/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		app         *fiber.App
		logger      *logger.Logger
		name        string
		validator   *validator.Validate
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
	service := New(description, name, validator.New())
	webserver := saddle.Command(saddle.Instantiate(service))
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

func New(description, name string, validator *validator.Validate) *Service {
	return &Service{
		config:      &models.Config{},
		description: description,
		name:        name,
		validator:   validator,
	}
}

func (s *Service) Attach(e *fiber.App, logger *logger.Logger, validator *validator.Validate) (
	func(), error,
) {
	s.app = e
	s.logger = logger
	s.validator = validator

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
		s.validator,
		s.extended.test, // Purely to show passing configurations through to handlers.
	).Add(s.app); err != nil {
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

func (s *Service) Validator() *validator.Validate {
	return s.validator
}

func (s *Service) shutdown() {
	for _, f := range _cleanup {
		f()
	}
	s.logger.Info("gracefully shutdown service")
}
