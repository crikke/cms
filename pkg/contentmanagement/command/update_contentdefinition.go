package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UpdateContentDefinition struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type UpdateContentDefinitionHandler struct {
}

func (c UpdateContentDefinitionHandler) Handle(ctx context.Context, cmd UpdateContentDefinition) (err error) {

	defer func() {
		// todo better logging
		fmt.Println("UpdateContentDefinitionHandler", cmd, err)
	}()

	return nil
}
