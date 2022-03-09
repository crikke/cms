package query

import (
	"context"
	"errors"
	"time"

	"github.com/crikke/cms/pkg/content"
	"github.com/google/uuid"
)

type GetContentByID struct {
	ID       uuid.UUID
	Language string
}

type ContentResponse struct {
	ID                 uuid.UUID
	AvailableLanguages []string
	Fields             content.ContentFields
	Created            time.Time `bson:"created"`
}
type GetContentByIDHandler struct {
	Repo content.ContentRepository
}

func (h GetContentByIDHandler) Handle(ctx context.Context, query GetContentByID) (ContentResponse, error) {

	if query.ID == (uuid.UUID{}) {
		return ContentResponse{}, errors.New("missing id")
	}

	contentresult, err := h.Repo.GetContent(ctx, query.ID)

	if err != nil {
		return ContentResponse{}, err
	}

	cv, err := contentresult.GetPublishedVersion()

	if err != nil {
		return ContentResponse{}, err
	}

	fields, ok := cv.Properties[query.Language]

	if !ok {
		return ContentResponse{}, errors.New(content.ErrMissingLanguage)
	}

	response := ContentResponse{
		ID:                 contentresult.ID,
		AvailableLanguages: cv.AvailableLanguages(),
		Fields:             fields,
		Created:            cv.Created,
	}

	return response, nil
}

type ContentListResponse struct {
	Items []ContentResponse

	// how many items was returned
	Count int
	// how many items exists
	Total int
}

type GetContentByTags struct {
	Tags []string
	// What fields to return
	Fields []string
}

type GetContentByTagsHandler struct {
	Repo content.ContentRepository
}

//! TODO: Query builder, Find a way how it can be done loosly coupled.
//! Since in the future, there will probably exist filtering on more than just tags

//! TODO: Find a way to use projection to only return specified fields from repository instead
//! of fetching the whole content object and projecting fields on the server
//! this could otherwise become an performance issue.

//! This should work by getting all fields from default locale where field.localized: false
//! and fields from specified locale where field.localized:true
func (h GetContentByTagsHandler) Handle(ctx context.Context, query GetContentByTags) ([]ContentListResponse, error) {

	items, err := h.Repo.ListContentByTags(ctx, query.Tags)

	if err != nil {
		return nil, err
	}

	result := make([]ContentListResponse, 0)
	for _, item := range items {
		if item.Status != content.Published {
			continue
		}
		// result = append(result, ContentListResponse{
		// 	ID: item.ID,
		// 	Name: item,
		// })
	}

	return result, nil
}
