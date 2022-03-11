package command

import (
	"context"

	"github.com/crikke/cms/pkg/workspace"
	"github.com/google/uuid"
)

type CreateWorkspace struct {
	Name          string
	Description   string
	DefaultLocale string
}

type CreateWorkspaceHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h CreateWorkspaceHandler) Handle(ctx context.Context, cmd CreateWorkspace) (uuid.UUID, error) {

	ws, err := workspace.NewWorkspace(cmd.Name, cmd.Description, cmd.DefaultLocale)

	if err != nil {
		return uuid.UUID{}, err
	}

	return h.Repo.Create(ctx, ws)
}

type UpdateTag struct {
	WorkspaceId uuid.UUID
	Id          string
	Name        string
}

type UpdateTagHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h UpdateTagHandler) Handle(ctx context.Context, cmd UpdateTag) error {

	return h.Repo.Update(ctx, cmd.WorkspaceId, func(ctx context.Context, ws *workspace.Workspace) (*workspace.Workspace, error) {

		ws.Tags[cmd.Id] = cmd.Name

		return ws, nil
	})
}

type DeleteTag struct {
	WorkspaceId uuid.UUID
	Id          string
}

type DeleteTagHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h DeleteTagHandler) Handle(ctx context.Context, cmd DeleteTag) error {

	return h.Repo.Update(ctx, cmd.WorkspaceId, func(ctx context.Context, ws *workspace.Workspace) (*workspace.Workspace, error) {
		delete(ws.Tags, cmd.Id)

		return ws, nil
	})
}
