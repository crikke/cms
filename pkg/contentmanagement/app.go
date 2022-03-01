package contentmanagement

import (
	contentcmd "github.com/crikke/cms/pkg/contentmanagement/content/command"
	"github.com/crikke/cms/pkg/contentmanagement/content/query"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Queries struct {
	GetContent  query.GetContentHandler
	ListContent query.ListChildContentHandler
}
type Commands struct {
	CreateContent contentcmd.CreateContentHandler
}
