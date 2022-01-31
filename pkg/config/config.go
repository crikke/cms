package config

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type Configuration struct {
	// Languages are configured by contentdelivery api. The elements are prioritized.
	Languages []language.Tag
	RootPage  uuid.UUID
}
