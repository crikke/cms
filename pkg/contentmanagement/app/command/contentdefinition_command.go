package command

import (
	"context"
	"fmt"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/workspace"
	"github.com/google/uuid"
)

type CreateContentDefinition struct {
	Name        string
	Description string
	WorkspaceId uuid.UUID
}

type CreateContentDefinitionHandler struct {
	WorkspaceRepo workspace.WorkspaceRepository
	Repo          contentdefinition.ContentDefinitionRepository
}

func (c CreateContentDefinitionHandler) Handle(ctx context.Context, cmd CreateContentDefinition) (id uuid.UUID, err error) {

	defer func() {
		// todo better logging
		fmt.Println("CreateContentDefinitionHandler", cmd, err)
	}()

	cd, err := contentdefinition.NewContentDefinition(cmd.Name, cmd.Description)
	if err != nil {
		return
	}

	id, err = c.Repo.CreateContentDefinition(ctx, &cd, cmd.WorkspaceId)

	return
}

type UpdateContentDefinition struct {
	ContentDefinitionID uuid.UUID `bson:"_id,omitempty"`
	Name                string    `bson:"omitempty"`
	Description         string    `bson:"omitempty"`
	WorkspaceId         uuid.UUID
	PropertyDefinitions map[string]struct {
		ID          uuid.UUID              `bson:_id, omitempty"`
		Description string                 `bson:"description, omitempty"`
		Localized   bool                   `bson:"localized, omitempty"`
		Validators  map[string]interface{} `bson:"validators, omitempty"`
	}
}

type UpdateContentDefinitionHandler struct {
	WorkspaceRepo            workspace.WorkspaceRepository
	Repo                     contentdefinition.ContentDefinitionRepository
	ContentDefinitionFactory contentdefinition.ContentDefinitionFactory

	// Properties are stored by their name, since name is unique
	// To allow chaning names map them by ID
}

func (c UpdateContentDefinitionHandler) Handle(ctx context.Context, cmd UpdateContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("UpdateContentDefinitionHandler", cmd, err)
	}()

	err = c.Repo.UpdateContentDefinition(ctx, cmd.ContentDefinitionID, cmd.WorkspaceId, func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

		if cmd.Name != "" {
			cd.Name = cmd.Name
		}

		if cmd.Description != "" {
			cd.Description = cmd.Description
		}

		// All propertydefinitions that havent been updated will be deleted.
		updatedProps := make(map[uuid.UUID]contentdefinition.PropertyDefinition, 0)

		for _, prop := range cmd.PropertyDefinitions {
			c.ContentDefinitionFactory.UpdatePropertyDefinition(
				cd,
				prop.ID,
				prop.Description,
				prop.Localized,
				prop.Validators)
			updatedProps[prop.ID] = contentdefinition.PropertyDefinition{
				ID:          prop.ID,
				Description: prop.Description,
				Localized:   prop.Localized,
				Validators:  prop.Validators,
			}

		}

		// check for props that should be deleted
		for name, prop := range cd.Propertydefinitions {

			if _, ok := updatedProps[prop.ID]; ok {
				continue
			}

			delete(cd.Propertydefinitions, name)

		}
		// todo: handle swapping propertynames
		for name, prop := range cmd.PropertyDefinitions {
			_, ok := cd.Propertydefinitions[name]

			if ok {
				continue
			}

			if err := c.ContentDefinitionFactory.UpdatePropertyDefinitionName(cd, prop.ID, name); err != nil {
				return nil, err
			}
		}
		return cd, nil
	})
	return
}

type DeleteContentDefinition struct {
	ID          uuid.UUID
	WorkspaceId uuid.UUID
}

type DeleteContentDefinitionHandler struct {
}

func (c DeleteContentDefinitionHandler) Handle(ctx context.Context, cmd DeleteContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("DeleteContentDefinitionHandler", cmd, err)
	}()

	if err != nil {
		return
	}

	return nil
}
