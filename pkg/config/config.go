package config

import (
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	ConnectionString struct {
		Mongodb  string
		RabbitMQ string
	}
	LogLevel int
}

func LoadServerConfiguration() ServerConfiguration {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ffcms/")
	viper.AddConfigPath("$HOME/.ffcms/")
	viper.AddConfigPath(".")

	viper.SetDefault("ConnectionString.Mongodb", "mongodb://0.0.0.0")
	viper.SetDefault("ConnectionString.RabbitMQ", "amqp://0.0.0.0")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	c := &ServerConfiguration{}

	err := viper.Unmarshal(c)

	if err != nil {
		panic(err)
	}

	return *c
}
