package config

import (
	"time"
)

type Config struct {
	Server ServerConfig
	Logger Logger
	Cache  Cache
}

type ServerConfig struct {
	Default DefaultServerConfig
	Grpc    GrpcServerConfig
}

type DefaultServerConfig struct {
	AppVersion     string        `yaml:"appVersion"`
	Port           string        `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
	AllowOrigins   []string      `yaml:"allowOrigins"`
}

type GrpcServerConfig struct {
	Port string `yaml:"port"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Cache struct {
	Capacity        int           `mapstructure:"cap"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval"`
}
