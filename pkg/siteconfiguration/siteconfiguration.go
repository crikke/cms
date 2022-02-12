package siteconfiguration

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type SiteConfiguration struct {
	Languages []language.Tag
	RootPage  uuid.UUID
}
