package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	App AppConfig
	Db  DbConfig
}

func NewConfig() *Config {
	viper.AddConfigPath("./")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Config{
		App: LoadAppConfig(),
		Db:  LoadDbConfig(),
	}
}
