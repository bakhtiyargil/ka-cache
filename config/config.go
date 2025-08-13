package config

import (
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
