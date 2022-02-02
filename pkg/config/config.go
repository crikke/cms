package config

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type ServerConfiguration struct {
	ConnectionString struct {
		Mongodb string
	}
	LogLevel int
}

// Configuration bound to site, such as root page & configured languages.
// Since this configuration is configured by users. It should not be stored as a ConfigMap.
// TODO: this can wait and have hardcoded defaults for now.
type SiteConfiguration struct {
	// Languages are configured by contentdelivery api. The elements are prioritized.
	Languages []language.Tag
	RootPage  uuid.UUID
}

func LoadServerConfiguration() ServerConfiguration {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ffcms/")
	viper.AddConfigPath("$HOME/.ffcms/")
	viper.AddConfigPath(".")

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
