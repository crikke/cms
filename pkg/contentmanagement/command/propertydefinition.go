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

type UpdatePropertyDefinition struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
	Name                 *string
	Description          *string
	Localized            *bool
}

type UpdatePropertyDefinitionHandler struct {
	repo contentdefinition.PropertyDefinitionRepository
}

func (h UpdatePropertyDefinitionHandler) Handle(ctx context.Context, cmd UpdatePropertyDefinition) error {

	h.repo.UpdatePropertyDefinition(
		ctx, cmd.ContentDefinitionID,
		cmd.PropertyDefinitionID,
		func(ctx context.Context, pd *contentdefinition.PropertyDefinition) (*contentdefinition.PropertyDefinition, error) {

			if cmd.Description != nil {
				pd.Description = *cmd.Description
			}

			if cmd.Name != nil {
				pd.Name = *cmd.Name
			}

			if cmd.Localized != nil {
				pd.Localized = *cmd.Localized
			}

			return pd, nil
		})

	return nil
}

type DeletePropertyDefinition struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
}

type DeletePropertyDefinitionHandler struct {
	repo contentdefinition.PropertyDefinitionRepository
}

func (h DeletePropertyDefinitionHandler) Handle(ctx context.Context, cmd DeletePropertyDefinition) error {
	return h.repo.DeletePropertyDefinition(ctx, cmd.ContentDefinitionID, cmd.PropertyDefinitionID)
}
