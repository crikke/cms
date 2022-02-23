package command

import (
	"context"

	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionType uuid.UUID
}

type CreateCommandHandler struct {
}

func (h CreateCommandHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	return uuid.UUID{}, nil
}
