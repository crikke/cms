package command

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/contentdefinition/validator"
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
	WorkspaceID         uuid.UUID
	Name                string
	Description         string
	Type                string
}

type CreatePropertyDefinitionHandler struct {
	Repo    contentdefinition.ContentDefinitionRepository
	Factory contentdefinition.ContentDefinitionFactory
}

func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

	if cmd.ContentDefinitionID == (uuid.UUID{}) {
		return uuid.UUID{}, errors.New("empty contentdefinition id")
	}

	var id uuid.UUID
	err := h.Repo.UpdateContentDefinition(
		ctx,
		cmd.ContentDefinitionID,
		cmd.WorkspaceID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {
			err := h.Factory.NewPropertyDefinition(cd, cmd.Name, cmd.Type, cmd.Description, false)
			if err != nil {
				return nil, err
			}

			id = cd.Propertydefinitions[cmd.Name].ID

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
	WorkspaceID          uuid.UUID
	Name                 *string
	Description          *string
	Localized            *bool
	Rules                map[string]interface{}
}

type UpdatePropertyDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h UpdatePropertyDefinitionHandler) Handle(ctx context.Context, cmd UpdatePropertyDefinition) error {

	if cmd.ContentDefinitionID == (uuid.UUID{}) {
		return errors.New("empty contentdefinition id")
	}

	if cmd.PropertyDefinitionID == (uuid.UUID{}) {
		return errors.New("empty propertydefinition id")
	}

	return h.Repo.UpdateContentDefinition(
		ctx,
		cmd.ContentDefinitionID,
		cmd.WorkspaceID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

			f := contentdefinition.ContentDefinitionFactory{}

			if *cmd.Name != "" && cmd.Name != nil {
				err := f.UpdatePropertyDefinitionName(cd, cmd.PropertyDefinitionID, *cmd.Name)

				if err != nil {
					return nil, err
				}
			}

			err := f.UpdatePropertyDefinition(cd, cmd.PropertyDefinitionID, *cmd.Description, *cmd.Localized, cmd.Rules)
			if err != nil {
				return nil, err
			}

			return cd, nil
		})
}

type DeletePropertyDefinition struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
	WorkspaceID          uuid.UUID
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
		ctx,
		cmd.ContentDefinitionID,
		cmd.WorkspaceID,
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
	WorkspaceID          uuid.UUID
}

type UpdateValidatorHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h UpdateValidatorHandler) Handle(ctx context.Context, cmd UpdateValidator) error {

	v, err := validator.Parse(cmd.ValidatorName, cmd.Value)

	if err != nil {
		return err
	}

	return h.Repo.UpdateContentDefinition(
		ctx,
		cmd.ContentDefinitionID,
		cmd.WorkspaceID,
		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

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
