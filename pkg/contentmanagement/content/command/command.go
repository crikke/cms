package command

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionType uuid.UUID
}

type CreateCommandHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
}

func (h CreateCommandHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	_, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cmd.ContentDefinitionType)

	if err != nil {
		return uuid.UUID{}, err
	}

	// _, err := h.ContentRepository.
	return uuid.UUID{}, nil
}
