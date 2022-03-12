package query

import (
	"context"

	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

// swagger:model ContentListReadModel
type ContentListReadModel struct {
	ID   uuid.UUID
	Name string
}

type ListContent struct {
	ContentDefinitionIDs []uuid.UUID
	WorkspaceId          uuid.UUID
}

type ListContentHandler struct {
	Repo content.ContentManagementRepository
	Cfg  *siteconfiguration.SiteConfiguration
}

//! TODO Should name be returned for current locale?
func (h ListContentHandler) Handle(ctx context.Context, query ListContent) ([]ContentListReadModel, error) {

	items, err := h.Repo.ListContentByContentDefinition(ctx, query.ContentDefinitionIDs, query.WorkspaceId)

	if err != nil {
		return nil, err
	}

	result := []ContentListReadModel{}

	for _, ch := range items {
		name := ch.Data.Properties[h.Cfg.Languages[0].String()][contentdefinition.NameField].Value
		result = append(result, ContentListReadModel{
			ID:   ch.ID,
			Name: name.(string),
		})
	}

	return result, nil
}

// // ContentReadModel is the representation of the content for the Content management API
// // It contains all information of given content for every configured language.
// //
// // swagger:model Contentresponse
// type ContentReadModel struct {
// 	ID                  uuid.UUID
// 	ContentDefinitionID uuid.UUID             `bson:"contentdefinition_id"`
// 	Status              content.PublishStatus `bson:"publishstatus"`
// 	// properties for the content
// 	Properties content.ContentLanguage `bson:"properties"`
// }

// In contentmanagement, all languages should be retrived for content of given version
// If Version is nil, return publishedversion
type GetContent struct {
	Id          uuid.UUID
	WorkspaceId uuid.UUID
	Version     int
}

type GetContentHandler struct {
	Repo content.ContentManagementRepository
}

func (q GetContentHandler) Handle(ctx context.Context, query GetContent) (content.Content, error) {

	c, err := q.Repo.GetContent(ctx, query.Id, query.Version, query.WorkspaceId)

	if err != nil {
		return content.Content{}, err
	}

	return c, nil
}
