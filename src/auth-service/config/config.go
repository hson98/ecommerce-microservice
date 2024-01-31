package config

import (
	"fmt"
	"github.com/hson98/ecommerce-microservice/src/auth-service/pkg/constants"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	Logger  Logger       `mapstructure:"logger"`
	Server  ServerConfig `mapstructure:"server"`
	Metrics Metrics      `mapstructure:"metrics"`
	Jaeger  Jaeger       `mapstructure:"jaeger"`
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disableCaller"`
	DisableStacktrace bool   `mapstructure:"disableStacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}
type Metrics struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"serviceName"`
}

// Jaeger
type Jaeger struct {
	Host        string `mapstructure:"host"`
	ServiceName string `mapstructure:"serviceName"`
	LogSpans    bool   `mapstructure:"logSpans"`
}

type ServerConfig struct {
	AppVersion           string        `mapstructure:"appVersion"`
	Mode                 string        `mapstructure:"mode"`
	SSL                  bool          `mapstructure:"ssl"`
	MaxConnectionIdle    time.Duration `mapstructure:"maxConnectionIdle"`
	Timeout              time.Duration `mapstructure:"timeout"`
	MaxConnectionAge     time.Duration `mapstructure:"maxConnectionAge"`
	Port                 string        `mapstructure:"port"`
	DBUser               string        `mapstructure:"dbUser"`
	DBPass               string        `mapstructure:"dbPass"`
	DBName               string        `mapstructure:"dbName"`
	DBHost               string        `mapstructure:"dbHost"`
	DBPort               string        `mapstructure:"dbBPort"`
	AccessTokenDuration  time.Duration `mapstructure:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `mapstructure:"refreshTokenDuration"`
	SecretKeyJWT         string        `mapstructure:"secretKeyJWT"`
}

func LoadConfig(configPath string) (config *Config, err error) {
	if configPath == "" {
		getwd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "os.Getwd")
		}
		configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
	}
	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}
	return cfg, nil
}
