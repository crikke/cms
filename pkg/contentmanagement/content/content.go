package content

import (
	"time"

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

// swagger:enum PublishStatus
type PublishStatus string

const (
	Draft     PublishStatus = "draft"
	Published PublishStatus = "published"
	Archived  PublishStatus = "archived"
)

type Content struct {
	ID                  uuid.UUID `bson:"_id"`
	ContentDefinitionID uuid.UUID `bson:"contentdefinition_id"`
	ParentID            uuid.UUID `bson:"parentid"`
	PublishedVersion    int
	Version             map[int]ContentVersion `bson:"version"`
	Status              PublishStatus          `bson:"status"`
}

type ContentVersion struct {
	Properties map[string]map[string]interface{} `bson:"properties"`
	Created    time.Time
}
