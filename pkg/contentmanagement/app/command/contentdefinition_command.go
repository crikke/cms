package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/google/uuid"
)

type CreateContentDefinition struct {
	Name        string
	Description string
}

func (c CreateContentDefinition) Valid() error {
	if c.Name == "" {
		return errors.New("missing field: Name")
	}
	return nil
}

type CreateContentDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
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

	id, err = c.Repo.CreateContentDefinition(ctx, &cd)

	return
}

type UpdateContentDefinition struct {
	ContentDefinitionID uuid.UUID `bson:"_id,omitempty"`
	Name                string    `bson:"omitempty"`
	Description         string    `bson:"omitempty"`
}

type UpdateContentDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (c UpdateContentDefinitionHandler) Handle(ctx context.Context, cmd UpdateContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("UpdateContentDefinitionHandler", cmd, err)
	}()

	err = c.Repo.UpdateContentDefinition(ctx, cmd.ContentDefinitionID, func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

		if cmd.Name != "" {
			cd.Name = cmd.Name
		}

		if cmd.Description != "" {
			cd.Description = cmd.Description
		}
		return cd, nil
	})
	return
}

type DeleteContentDefinition struct {
	ID uuid.UUID
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
