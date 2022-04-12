package app

import (
	"github.com/crikke/cms/cmd/contentdelivery/app/query"
)

type Queries struct {
	GetContentByID query.GetContentByIDHandler
}
type App struct {
	Queries Queries
}
