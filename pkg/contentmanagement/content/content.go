package content

import (
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

/*

Saving content example:

HTTP POST contentmanagement/content
{
	"contentdefinition": uuid,
}

HTTP PUT contentmanagement/content/{contentID}
{
	"name": {
		"sv_SE": "exempel",
		"en_US": "example"
	},
	"properties": [
		{
			"id": uuid
			"value": {
				"sv_SE": "example value"
			},
			"localized": false|true // todo should PUT contain this? isnt contentdefinition responsible for this?
		}
	]
}

*/

type SaveStatus int64

const (

	// Status is Draft when the content is saved but the version of given content has not previously been published
	Draft SaveStatus = iota
	// Indicates that there is a newer version available
	Unpublished
	Published
	// When content is archived, it wont be available for consumers.
	Archived
)

type Content struct {
	ID                  uuid.UUID
	ContentDefinitionID uuid.UUID
	Version             int
	Properties          map[language.Tag]map[string]interface{}
	Status              SaveStatus
}
