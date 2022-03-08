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
	ListContent query.ListContentHandler

	GetContentDefinition   query.GetContentDefinitionHandler
	GetPropertyDefinition  query.GetPropertyDefinitionHandler
	ListContentDefinitions query.ListContentDefinitionHandler
}
type Commands struct {
	CreateContent       contentcmd.CreateContentHandler
	UpdateContentFields contentcmd.UpdateContentFieldsHandler
	ArchiveContent      contentcmd.ArchiveContentHandler
	PublishContent      contentcmd.PublishContentHandler

	CreateContentDefinition contentcmd.CreateContentDefinitionHandler
	UpdateContentDefinition contentcmd.UpdateContentDefinitionHandler
	DeleteContentDefinition contentcmd.DeleteContentDefinitionHandler

	CreatePropertyDefinition contentcmd.CreatePropertyDefinitionHandler
	UpdatePropertyDefinition contentcmd.UpdatePropertyDefinitionHandler
	DeletePropertyDefinition contentcmd.DeletePropertyDefinitionHandler

	UpdateSiteConfiguration contentcmd.UpdateSiteConfigurationHandler
}
