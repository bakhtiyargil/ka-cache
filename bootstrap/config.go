package bootstrap

import (
	"github.com/spf13/viper"
	"ka-cache/config"
	"log"
)

func LoadConfig(filename string) *config.Config {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	var c config.Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}
	return &c
}
