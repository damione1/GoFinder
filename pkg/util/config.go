package util

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	HTTPServerPort    string `mapstructure:"HTTP_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	GCSBucketName     string `mapstructure:"GCS_BUCKET_NAME"`
	Domain            string `mapstructure:"DOMAIN"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()

	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("HTTP_SERVER_PORT")
	viper.BindEnv("GRPC_SERVER_PORT")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("GCS_BUCKET_NAME")
	viper.BindEnv("DOMAIN")

	if err = viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
