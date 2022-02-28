package query

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/google/uuid"
)

type ContentReadModel struct {
	ID                  uuid.UUID
	ContentDefinitionID uuid.UUID                         `bson:"contentdefinition_id"`
	Status              content.SaveStatus                `bson:"status"`
	Properties          map[string]map[string]interface{} `bson:"properties"`
}

// In contentmanagement, all languages should be retrived for content of given version
// If Version is nil, return publishedversion
type GetContent struct {
	Id      uuid.UUID
	Version *int
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
		return ContentReadModel{}, errors.New(content.ErrVersionNotExists)
	}

	rm := ContentReadModel{
		ID:                  c.ID,
		ContentDefinitionID: c.ContentDefinitionID,
		Status:              c.Status,
		Properties:          contentVer.Properties,
	}

	return rm, nil
}
