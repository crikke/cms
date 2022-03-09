package app

import (
	"github.com/crikke/cms/pkg/contentdelivery/app/query"
	"github.com/crikke/cms/pkg/siteconfiguration"
)

type Queries struct {
	GetContentByID query.GetContentByIDHandler
}
type App struct {
	Queries Queries

	SiteConfiguration *siteconfiguration.SiteConfiguration
}
