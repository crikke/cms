package app

import (
	contentcmd "github.com/crikke/cms/cmd/contentmanagement/app/command"
	"github.com/crikke/cms/cmd/contentmanagement/app/query"
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

	WorkspaceQueries WorkspaceQueries
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

	WorkspaceCommands WorkspaceCommands
}

type WorkspaceCommands struct {
	CreateWorkspace contentcmd.CreateWorkspaceHandler
	UpdateWorkspace contentcmd.UpdateWorkspaceHandler
	UpdateTag       contentcmd.UpdateTagHandler
	DeleteTag       contentcmd.DeleteTagHandler
}

type WorkspaceQueries struct {
	GetTag         query.GetTagHandler
	GetWorkspace   query.GetWorkspaceHandler
	ListTags       query.ListTagsHandler
	ListWorkspaces query.ListWorkspaceHandler
}
