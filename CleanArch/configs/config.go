package configs

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "root")
	viper.SetDefault("DB_NAME", "orders")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("WEB_SERVER_PORT", "8000")
	viper.SetDefault("GRPC_SERVER_PORT", "50051")
	viper.SetDefault("GRAPHQL_SERVER_PORT", "8080")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) && !os.IsNotExist(err) {
			return nil, err
		}
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
