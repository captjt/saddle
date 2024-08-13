package saddle

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/captjt/saddle/models"
	log "github.com/captjt/saddle/pkg/logger"
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
		logger.Fatal("saddle config deserialization error",
			zap.Error(err),
		)
	}
	if err := v.Struct(hc); err != nil {
		logger.Fatal("saddle config validation error",
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
		// Have to call config() here to ensure the environment is set before the logger is updated.
		_ = config(service, env)
		// display project logo w/ service name, environment and description
		fmt.Printf("%s\n%s [%s]\n   ⤷ %s\n\n", logo, service.Name(), env,
			service.Description())

		// update logger for proper env and service
		logger.SetEnvironment(log.Environment(env), service.Name())

		var (
			err error
		)

		// - instantiate new service ↴
		s, err := new(service, logger, service.Validator())
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
		return s.App.Listen(address)
	}
}
