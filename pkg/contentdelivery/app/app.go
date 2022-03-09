package app

import "github.com/crikke/cms/pkg/siteconfiguration"

type Queries struct {
}
type App struct {
	Queries Queries

	SiteConfiguration *siteconfiguration.SiteConfiguration
}
