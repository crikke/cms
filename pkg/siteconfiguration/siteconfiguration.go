package siteconfiguration

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type SiteConfiguration struct {
	// Site ID
	ID        uuid.UUID      `bson:"_id"`
	Languages []language.Tag `bson:"languages"`
	RootPage  uuid.UUID      `bson:"rootpage"`
}
