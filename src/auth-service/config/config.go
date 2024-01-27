package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Port                 string        `mapstructure:"PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPass               string        `mapstructure:"DB_PASS"`
	DBName               string        `mapstructure:"DB_NAME"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBPort               string        `mapstructure:"DB_PORT"`
	RedisHost            string        `mapstructure:"REDIS_HOST"`
	RedisPass            string        `mapstructure:"REDIS_PASS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SecretKeyJWT         string        `mapstructure:"SECRET_KEY"`
	Logger               Logger
	Server               ServerConfig
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

func LoadConfig(path string, name string) (config *Config, err error) {
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	return config, nil
}
