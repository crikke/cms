package query

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

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
	Version *int
}

type ListChildContentHandler struct {
	Repo content.ContentRepository
	Cfg  *siteconfiguration.SiteConfiguration
}

func (h ListChildContentHandler) Handle(ctx context.Context, query ListChildContent) ([]ContentListReadModel, error) {

	children, err := h.Repo.ListContent(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	result := []ContentListReadModel{}

	for _, ch := range children {

		name := ch.Version[ch.PublishedVersion].Properties[h.Cfg.Languages[0].String()][content.NameField].Value
		result = append(result, ContentListReadModel{
			ID:   ch.ID,
			Name: name.(string),
		})
	}

	return result, nil
}

type GetContentHandler struct {
	Repo content.ContentRepository
}

func (q GetContentHandler) Handle(ctx context.Context, query GetContent) (ContentReadModel, error) {

	c, err := q.Repo.GetContent(ctx, query.Id)

	if err != nil {
		return ContentReadModel{}, err
	}
	v := c.PublishedVersion

	if query.Version != nil {
		v = *query.Version
	}

	contentVer, ok := c.Version[v]

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

type ListChildContent struct {
	ID uuid.UUID
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
