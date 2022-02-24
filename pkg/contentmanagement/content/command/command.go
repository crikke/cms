package command

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionId uuid.UUID
}

type CreateContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
}

func (h CreateContentHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cmd.ContentDefinitionId)

	if err != nil {
		return uuid.UUID{}, err
	}

	c := content.Content{
		ContentDefinitionID: cd.ID,
		Status:              content.Draft,
	}

	id, err := h.ContentRepository.CreateContent(ctx, c)

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
