package app

import (
	contentcmd "github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Queries struct {
	GetContent  query.GetContentHandler
	ListContent query.ListChildContentHandler

	GetContentDefinition query.GetContentDefinitionHandler
}
type Commands struct {
	CreateContent      contentcmd.CreateContentHandler
	UpdateContentField contentcmd.UpdateContentFieldHandler
	DeleteContent      contentcmd.DeleteContentHandler

	CreateContentDefinition contentcmd.CreateContentDefinitionHandler
	UpdateContentDefinition contentcmd.UpdateContentDefinitionHandler
	DeleteContentDefinition contentcmd.DeleteContentDefinitionHandler

	CreatePropertyDefinition contentcmd.CreatePropertyDefinitionHandler
}
