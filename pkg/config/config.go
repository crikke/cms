package config

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type ServerConfiguration struct {
	ConnectionStrings struct {
		Mongodb string
	}
}

// Configuration bound to site, such as root page & configured languages
type SiteConfiguration struct {
	// Languages are configured by contentdelivery api. The elements are prioritized.
	Languages []language.Tag
	RootPage  uuid.UUID
}

func LoadConfiguration() {

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
}
