package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	AppVersion   string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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
	}
	return &c, nil
}
