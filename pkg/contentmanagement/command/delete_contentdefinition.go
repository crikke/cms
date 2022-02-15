package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

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
