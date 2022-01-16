package config

import "github.com/google/uuid"

// constants

type languageKey string

var (
	LanguageKey       languageKey
	SiteConfiguration Configuration
)

type Configuration struct {
	DefaultLanguage string
	RootPage        uuid.UUID
}
