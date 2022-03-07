package query

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

type ListContent struct {
	ContentDefinitionIDs []uuid.UUID
}

type ListContentHandler struct {
	Repo content.ContentRepository
	Cfg  *siteconfiguration.SiteConfiguration
}

func (h ListContentHandler) Handle(ctx context.Context, query ListContent) ([]ContentListReadModel, error) {

	children, err := h.Repo.ListContent(ctx, query.ContentDefinitionIDs)

	if err != nil {
		return nil, err
	}

	result := []ContentListReadModel{}

	for _, ch := range children {

		name := ch.Version[ch.PublishedVersion].Properties[h.Cfg.Languages[0].String()][contentdefinition.NameField].Value
		result = append(result, ContentListReadModel{
			ID:   ch.ID,
			Name: name.(string),
		})
	}

	return result, nil
}

// ContentReadModel is the representation of the content for the Content management API
// It contains all information of given content for every configured language.
//
// swagger:model Contentresponse
type ContentReadModel struct {
	ID                  uuid.UUID
	ContentDefinitionID uuid.UUID             `bson:"contentdefinition_id"`
	Status              content.PublishStatus `bson:"publishstatus"`
	// properties for the content
	Properties content.ContentLanguage `bson:"properties"`
}

// In contentmanagement, all languages should be retrived for content of given version
// If Version is nil, return publishedversion
type GetContent struct {
	Id      uuid.UUID
	Version int
}

type GetContentHandler struct {
	Repo content.ContentRepository
}

func (q GetContentHandler) Handle(ctx context.Context, query GetContent) (ContentReadModel, error) {

	c, err := q.Repo.GetContent(ctx, query.Id)

	if err != nil {
		return ContentReadModel{}, err
	}

	contentVer, ok := c.Version[query.Version]

	if !ok {
		return ContentReadModel{}, errors.New(content.ErrMissingVersion)
	}

	rm := ContentReadModel{
		ID:                  c.ID,
		ContentDefinitionID: c.ContentDefinitionID,
		Status:              c.Status,
		Properties:          contentVer.Properties,
	}

	return rm, nil
}

type ContentListReadModel struct {
	ID   uuid.UUID
	Name string
}

// type ListChildContentHandler struct {
// 	Repo content.ContentRepository
// 	Cfg  *siteconfiguration.SiteConfiguration
// }

// func (h ListChildContentHandler) Handle(ctx context.Context, query ListChildContent) ([]ContentListReadModel, error) {

// 	children, err := h.Repo.ListContent(ctx, query.ID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	result := []ContentListReadModel{}

// 	for _, ch := range children {

// 		name := ch.Version[ch.PublishedVersion].Properties[h.Cfg.Languages[0].String()][content.NameField]
// 		result = append(result, ContentListReadModel{
// 			ID:   ch.ID,
// 			Name: name.(string),
// 		})
// 	}

// 	return result, nil
// }
