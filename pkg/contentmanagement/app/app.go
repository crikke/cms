package app

import (
	contentcmd "github.com/crikke/cms/pkg/contentmanagement/app/command"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Queries struct {
	// GetContent  query.GetContentHandler
	// ListContent query.ListChildContentHandler
}
type Commands struct {
	CreateContent contentcmd.CreateContentHandler
}
