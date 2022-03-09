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
