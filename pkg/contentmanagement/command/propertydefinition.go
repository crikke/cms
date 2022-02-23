package command

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition/validator"
	"github.com/google/uuid"
)

/*
example:

POST /contentdefinition/propertydefinition
{
	cid: uuid,
	name: name
	description: desc,
	type: text,
}
*/

/*
PUT /contentdefinition/{cid}/propertydefinition/{pid}
{
	required: true,
	regex: "^(foo)",
	unique: true,
	localized: true,
	boundary: { // coordinates - polygon in which a coord needs to be to be valid
		x,y
		x,y
		x,y
	}
}
*/

type CreatePropertyDefinition struct {
	ContentDefinitionID uuid.UUID
	Name                string
	Description         string
	Type                string
}

type CreatePropertyDefinitionHandler struct {
	Repo propertydefinition.PropertyDefinitionRepository
}

func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

	pd, err := propertydefinition.NewPropertyDefinition(cmd.ContentDefinitionID, cmd.Name, cmd.Description, cmd.Type)

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := h.Repo.CreatePropertyDefinition(ctx, cmd.ContentDefinitionID, &pd)
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
	repo propertydefinition.PropertyDefinitionRepository
}

func (h UpdatePropertyDefinitionHandler) Handle(ctx context.Context, cmd UpdatePropertyDefinition) error {

	h.repo.UpdatePropertyDefinition(
		ctx, cmd.ContentDefinitionID,
		cmd.PropertyDefinitionID,
		func(ctx context.Context, pd *propertydefinition.PropertyDefinition) (*propertydefinition.PropertyDefinition, error) {

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
	repo propertydefinition.PropertyDefinitionRepository
}

func (h DeletePropertyDefinitionHandler) Handle(ctx context.Context, cmd DeletePropertyDefinition) error {
	return h.repo.DeletePropertyDefinition(ctx, cmd.ContentDefinitionID, cmd.PropertyDefinitionID)
}

type UpdateValidator struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
	ValidatorName        string
	Value                interface{}
}

type UpdateValidatorHandler struct {
	Repo propertydefinition.PropertyDefinitionRepository
}

func (h UpdateValidatorHandler) Handle(ctx context.Context, cmd UpdateValidator) error {

	v, err := validator.Parse(cmd.ValidatorName, cmd.Value)

	if err != nil {
		return err
	}

	return h.Repo.UpdatePropertyDefinition(
		ctx,
		cmd.ContentDefinitionID,
		cmd.PropertyDefinitionID,
		func(ctx context.Context, pd *propertydefinition.PropertyDefinition) (*propertydefinition.PropertyDefinition, error) {
			if pd.Validators == nil {
				pd.Validators = make(map[string]interface{})
			}

			pd.Validators[cmd.ValidatorName] = v
			return pd, nil
		})
}
