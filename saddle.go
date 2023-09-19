package saddle

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/captjt/saddle/models"
	log "github.com/captjt/saddle/pkg/logger"
	"github.com/captjt/saddle/pkg/telemetry"
)

const (
	configFolder = ".config"

	logo = `
	________     __       ________   ________   ___       _______
	/"       )   /""\     |"      "\ |"      "\ |"  |     /"     "|
   (:   \___/   /    \    (.  ___  :)(.  ___  :)||  |    (: ______)
	\___  \    /' /\  \   |: \   ) |||: \   ) |||:  |     \/    |
	 __/  \\  //  __'  \  (| (___\ ||(| (___\ || \  |___  // ___)_
	/" \   :)/   /  \\  \ |:       :)|:       :)( \_|:  \(:      "|
   (_______/(___/    \___)(________/ (________/  \_______)\_______)
`
)

var logger *log.Logger

func init() {
	logger = log.New(log.Unknown, "saddle up! service initialization")
}

func New(version string) *cobra.Command {
	return &cobra.Command{
		Use:     "up",
		Long:    logo,
		Version: version,
	}
}

func config[T Service](service T, environment string) *models.Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	// viper.SetEnvPrefix(service.Name())

	// - handle import of any configuration file; set by referenced environment ↴
	viper.SetConfigName(environment)
	viper.AddConfigPath(configFolder)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatal("config file not found",
				zap.String("config folder", configFolder),
				zap.String("environment", environment),
			)
		} else {
			logger.Fatal("config file error",
				zap.Error(err),
			)
		}
	}

	v := validator.New()

	// - deserialize | validate saddle configuration(s) ↴
	hc := &models.Config{}
	if err := viper.Unmarshal(hc); err != nil {
		logger.Fatal("harnesser config deserialization error",
			zap.Error(err),
		)
	}
	if err := v.Struct(hc); err != nil {
		logger.Fatal("harnesser config validation error",
			zap.Error(err),
		)
	}

	// - deserialize | validate service configuration(s) ↴
	sc := service.Config()
	if err := viper.Unmarshal(sc); err != nil {
		logger.Fatal("service config deserialization error",
			zap.Error(err),
		)
	}
	if err := v.Struct(sc); err != nil {
		logger.Fatal("service config validation error",
			zap.Error(err),
		)
	}

	return hc
}

func Command[T Service](service T, entry func(*cobra.Command, []string) error) *cobra.Command {
	sn := service.Name()
	return &cobra.Command{
		Use:   sn,
		Short: fmt.Sprintf("%s service", sn),
		RunE:  entry,
	}
}

func Instantiate[T Service](service T) (T, func(cmd *cobra.Command, args []string) error) {
	return service, func(cmd *cobra.Command, args []string) error {
		env, address := viper.GetString(fmt.Sprintf("%s.%s", service.Name(), "environment")),
			viper.GetString(fmt.Sprintf("%s.%s", service.Name(), "address"))
		hc := config(service, env)

		// display project logo w/ service name, environment and description
		fmt.Printf("%s\n%s [%s]\n   ⤷ %s\n\n", logo, service.Name(), env,
			service.Description())

		// update logger for proper env and service
		logger.SetEnvironment(log.Environment(env), service.Name())

		// - instantiate open-telemetry tracing exporter ↴
		var (
			err error
			tp  *trace.TracerProvider
		)

		switch {
		case hc.Saddle.CloudTrace != nil:
			// TODO: Add project ID somehow in here for cloud trace settings.
			tp, err = telemetry.NewExporter(telemetry.CloudTrace, service.Name(), env,
				hc.Saddle.CloudTrace.SampleRate)
			if err != nil {
				logger.Fatal("unable to attach cloud trace telemetry exporter",
					zap.String("service", service.Name()),
					zap.String("environment", env),
					zap.Float64("sample rate %", float64(hc.Saddle.CloudTrace.SampleRate)),
					zap.Error(err),
				)
			}
		case hc.Saddle.Jaeger != nil:
			tp, err = telemetry.NewExporter(telemetry.Jaeger, service.Name(), env,
				hc.Saddle.Jaeger.SampleRate, hc.Saddle.Jaeger.URI)
			if err != nil {
				logger.Fatal("unable to attach jaeger telemetry exporter",
					zap.String("service", service.Name()),
					zap.String("environment", env),
					zap.String("jaeger uri", hc.Saddle.Jaeger.URI),
					zap.Float64("sample rate %", float64(hc.Saddle.Jaeger.SampleRate)),
					zap.Error(err),
				)
			}
		case hc.Saddle.StdOut != nil:
			tp, err = telemetry.NewExporter(telemetry.StdOut, service.Name(), env,
				hc.Saddle.StdOut.SampleRate)
			if err != nil {
				logger.Fatal("unable to attach stdout telemetry exporter",
					zap.String("service", service.Name()),
					zap.String("environment", env),
					zap.Float64("sample rate %", float64(hc.Saddle.StdOut.SampleRate)),
					zap.Error(err),
				)
			}
		default: // no exporter was specified in the service configuration; use stdout as default
			const defaultSampleRate = 10 // defaulted to 10% of overall requests are sampled
			tp, err = telemetry.NewExporter(telemetry.StdOut, service.Name(), env, defaultSampleRate)
			if err != nil {
				logger.Fatal("unable to attach stdout telemetry exporter (defaulted)",
					zap.String("service", service.Name()),
					zap.String("environment", env),
					zap.Float64("sample rate %", defaultSampleRate),
					zap.Error(err),
				)
			}
		}
		defer func() {
			ctx := context.Background()

			tp.ForceFlush(ctx) // flush any pending spans
			if err := tp.Shutdown(ctx); err != nil {
				logger.Error("unable to shutdown tracer provider",
					zap.Error(err),
				)
			}
		}()

		// - instantiate new service ↴
		s, err := new(service, logger, otel.Tracer(service.Name()))
		if err != nil {
			logger.Fatal("unable to attach service",
				zap.String("service", service.Name()),
				zap.Error(err),
			)
		}

		// - execute | expose service ↴
		logger.Info("listening for requests",
			zap.String("address", address),
		)
		return s.Echo.Start(address)
	}
}
