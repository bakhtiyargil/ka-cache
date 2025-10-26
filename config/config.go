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
	EnableSecure   bool          `yaml:"enableSecure"`
	SecurePort     string        `yaml:"securePort"`
	CertFile       string        `yaml:"certFile"`
	KeyFile        string        `yaml:"keyFile"`
}

type GrpcServerConfig struct {
	Port         string `yaml:"port"`
	EnableSecure bool   `yaml:"enableSecure"`
	SecurePort   string `yaml:"securePort"`
	CertFile     string `yaml:"certFile"`
	KeyFile      string `yaml:"keyFile"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Cache struct {
	InitCapacity    int           `mapstructure:"initCapacity"`
	MaxCapacity     int           `mapstructure:"maxCapacity"`
	CleanupInterval time.Duration `mapstructure:"cleanupInterval"`
	Shrink          Shrink
}

type Shrink struct {
	UtilizationThreshold float64 `yaml:"utilizationThreshold"`
}
