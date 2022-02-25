package content

import (
	"github.com/google/uuid"
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

const (
	NameField       = "name"
	UrlSegmentField = "url"
)

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
	ID                  uuid.UUID                         `bson:"_id"`
	ContentDefinitionID uuid.UUID                         `bson:"contentdefinition_id"`
	Version             int                               `bson:"version"`
	Properties          map[string]map[string]interface{} `bson:"properties"`
	Status              SaveStatus                        `bson:"status"`
}
