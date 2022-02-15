package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
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
}

func (c CreateContentDefinitionHandler) Handle(ctx context.Context, cmd CreateContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("CreateContentDefinitionHandler", cmd, err)
	}()

	_, err = contentdefinition.NewContentDefinition(cmd.Name, cmd.Description)

	if err != nil {
		return
	}

	return nil
}
