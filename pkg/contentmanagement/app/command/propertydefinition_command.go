package command

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
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
	Repo contentdefinition.ContentDefinitionRepository
}

func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

	if cmd.ContentDefinitionID == (uuid.UUID{}) {
		return uuid.UUID{}, errors.New("empty contentdefinition id")
	}

	var id uuid.UUID
	err := h.Repo.UpdateContentDefinition(
		ctx,
		cmd.ContentDefinitionID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {
			pd, err := contentdefinition.NewPropertyDefinition(cd, cmd.Name, cmd.Description, cmd.Type)
			if err != nil {
				return nil, err
			}

			id = pd.ID

			return cd, nil
		})

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
	repo contentdefinition.ContentDefinitionRepository
}

func (h UpdatePropertyDefinitionHandler) Handle(ctx context.Context, cmd UpdatePropertyDefinition) error {

	if cmd.ContentDefinitionID == (uuid.UUID{}) {
		return errors.New("empty contentdefinition id")
	}

	if cmd.PropertyDefinitionID == (uuid.UUID{}) {
		return errors.New("empty propertydefinition id")
	}

	return h.repo.UpdateContentDefinition(
		ctx, cmd.ContentDefinitionID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

			pd := contentdefinition.PropertyDefinition{}
			pdName := ""
			for n, p := range cd.Propertydefinitions {
				if p.ID == cmd.PropertyDefinitionID {
					pd = p
					pdName = n
					break
				}
			}

			if pd.ID == (uuid.UUID{}) {
				return nil, errors.New("propertydefinition not found")
			}

			if cmd.Description != nil {
				pd.Description = *cmd.Description
			}

			if cmd.Name != nil && *cmd.Name != "" && *cmd.Name != pdName {

				delete(cd.Propertydefinitions, pdName)
				pdName = *cmd.Name
			}

			if cmd.Localized != nil {
				pd.Localized = *cmd.Localized
			}
			cd.Propertydefinitions[pdName] = pd
			return cd, nil
		})
}

type DeletePropertyDefinition struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
}

type DeletePropertyDefinitionHandler struct {
	repo contentdefinition.ContentDefinitionRepository
}

func (h DeletePropertyDefinitionHandler) Handle(ctx context.Context, cmd DeletePropertyDefinition) error {

	if cmd.ContentDefinitionID == (uuid.UUID{}) {
		return errors.New("empty contentdefinition id")
	}

	if cmd.PropertyDefinitionID == (uuid.UUID{}) {
		return errors.New("empty propertydefinition id")
	}

	return h.repo.UpdateContentDefinition(
		ctx, cmd.ContentDefinitionID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

			name := ""
			for n, p := range cd.Propertydefinitions {
				if p.ID == cmd.PropertyDefinitionID {
					name = n
					break
				}
			}
			if name == "" {
				return nil, errors.New("propertydefinition not found")
			}

			delete(cd.Propertydefinitions, name)
			return cd, nil
		})
}

type UpdateValidator struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
	ValidatorName        string
	Value                interface{}
}

type UpdateValidatorHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h UpdateValidatorHandler) Handle(ctx context.Context, cmd UpdateValidator) error {

	v, err := validator.Parse(cmd.ValidatorName, cmd.Value)

	if err != nil {
		return err
	}

	return h.Repo.UpdateContentDefinition(ctx,
		cmd.ContentDefinitionID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

			// idx := 0
			pd := contentdefinition.PropertyDefinition{}
			name := ""
			for n, p := range cd.Propertydefinitions {

				if p.ID == cmd.PropertyDefinitionID {
					pd = p
					name = n
					break
				}
			}

			if pd.ID == (uuid.UUID{}) {
				return nil, errors.New("propertydefinition not found")
			}

			if pd.Validators == nil {
				pd.Validators = make(map[string]interface{})
			}

			pd.Validators[cmd.ValidatorName] = v
			cd.Propertydefinitions[name] = pd
			return cd, nil
		})
}
