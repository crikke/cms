package command

import (
	"context"
	"fmt"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type UpdateContentDefinition struct {
	ID          uuid.UUID `bson:"_id,omitempty"`
	Name        string    `bson:"omitempty"`
	Description string    `bson:"omitempty"`
}

type UpdateContentDefinitionHandler struct {
	repo contentdefinition.ContentDefinitionRepository
}

func (c UpdateContentDefinitionHandler) Handle(ctx context.Context, cmd UpdateContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("UpdateContentDefinitionHandler", cmd, err)
	}()

	err = c.repo.UpdateContentDefinition(ctx, cmd.ID, func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

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
