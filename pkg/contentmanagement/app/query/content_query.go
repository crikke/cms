package query

import (
	"context"
	"time"

	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/crikke/cms/pkg/workspace"
	"github.com/google/uuid"
)

// swagger:model ContentListReadModel
type ContentListReadModel struct {
	ID   uuid.UUID
	Name string
}

type ListContent struct {
	ContentDefinitionIDs []uuid.UUID
	Tags                 []string
	WorkspaceId          uuid.UUID
}

type ListContentHandler struct {
	Repo content.ContentManagementRepository
	Cfg  *siteconfiguration.SiteConfiguration
}

//! TODO Should name be returned for current locale?
func (h ListContentHandler) Handle(ctx context.Context, query ListContent) ([]ContentListReadModel, error) {

	items, err := h.Repo.ListContent(ctx, query.ContentDefinitionIDs, query.Tags, query.WorkspaceId)

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

// ContentReadModel is the representation of the content for the Content management API
// It contains all information of given content for every configured language.
type ContentReadModel struct {
	ID                  uuid.UUID
	ContentDefinitionID uuid.UUID `bson:"contentdefinition_id"`
	Version             int       `bson:"version"`

	Status content.PublishStatus `bson:"status"`
	// properties for the content
	Properties content.ContentLanguage `bson:"properties"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`

	Tags map[string]string `bson:"tags,omitempty"`
}

// In contentmanagement, all languages should be retrived for content of given version
// If Version is nil, return publishedversion
type GetContent struct {
	Id          uuid.UUID
	WorkspaceId uuid.UUID
	Version     int
}

type GetContentHandler struct {
	Repo          content.ContentManagementRepository
	WorkspaceRepo workspace.WorkspaceRepository
}

func (q GetContentHandler) Handle(ctx context.Context, query GetContent) (ContentReadModel, error) {

	c, err := q.Repo.GetContent(ctx, query.Id, query.Version, query.WorkspaceId)
	if err != nil {
		return ContentReadModel{}, err
	}

	// ContentData only stores tagId, instead of only returning tagId, tagId & tagName is returned.
	// This is done by getting all tags on workspace, performance is not prioritized in ContentManagement API so
	// doing it this was is OK.

	ws, err := q.WorkspaceRepo.Get(ctx, query.WorkspaceId)
	if err != nil {
		return ContentReadModel{}, err
	}

	crm := ContentReadModel{
		ID:                  c.ID,
		ContentDefinitionID: c.ContentDefinitionID,
		Version:             c.Data.Version,
		Status:              c.Data.Status,
		Created:             c.Created,
		Updated:             c.Updated,
		Properties:          c.Data.Properties,
		Tags:                make(map[string]string),
	}

	for _, tagId := range c.Data.Tags {

		tagName, ok := ws.Tags[tagId]

		// tag exist on content but does not exist in workspace
		// TODO: handle this
		if !ok {
			continue
		}

		crm.Tags[tagId] = tagName
	}

	return crm, nil
}
