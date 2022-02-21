package command

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type CreatePropertyDefinition struct {
	ContentDefinitionID uuid.UUID
	Name                string
	Description         string
	Type                string
}

type CreatePropertyDefinitionHandler struct {
	repo contentdefinition.PropertyDefinitionRepository
}

func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

	pd, err := contentdefinition.NewPropertyDefinition(cmd.ContentDefinitionID, cmd.Name, cmd.Description, cmd.Type)

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := h.repo.CreatePropertyDefinition(ctx, cmd.ContentDefinitionID, &pd)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
