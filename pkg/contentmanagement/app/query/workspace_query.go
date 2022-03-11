package query

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/workspace"
	"github.com/google/uuid"
)

type GetWorkspaceHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h GetWorkspaceHandler) Handle(ctx context.Context, id uuid.UUID) (workspace.Workspace, error) {

	if id == (uuid.UUID{}) {
		return workspace.Workspace{}, errors.New("empty id")
	}

	ws, err := h.Repo.Get(ctx, id)

	if err != nil {
		return workspace.Workspace{}, err
	}

	return ws, nil
}
