package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server ServerConfig
	Logger Logger
}

type ServerConfig struct {
	Default DefaultServerConfig
	Grpc    GrpcServerConfig
}

type DefaultServerConfig struct {
	AppVersion   string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GrpcServerConfig struct {
	Port string
}

type Logger struct {
	Level string
}

func LoadConfig(filename string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
		return nil, err
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
		return nil, err
	}
	return &c, nil
}
