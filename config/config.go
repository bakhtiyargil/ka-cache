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
	AppVersion     string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type GrpcServerConfig struct {
	Port string
}

type Logger struct {
	Level string
}

func LoadConfig(filename string) *Config {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}
	return &c
}
