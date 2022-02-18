package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
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
	repo contentdefinition.ContentDefinitionRepository
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

	id, err = c.repo.CreateContentDefinition(ctx, &cd)

	return
}
