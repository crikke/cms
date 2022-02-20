package command

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/google/uuid"
)

type CreatePropertyDefinition struct {
	ContentDefinitionID uuid.UUID
	Name                string
	Description         string
	Type                string
}

type CreatePropertyDefinitionHandler struct {
	repo propertydefinition.PropertyDefinitionRepository
}

func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

	pd, err := propertydefinition.NewPropertyDefinition(cmd.Name, cmd.Description, cmd.Type)

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := h.repo.CreatePropertyDefinition(ctx, &pd)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
