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

type Tag struct {
	Id   string
	Name string
}
type GetTag struct {
	Id          string
	WorkspaceId uuid.UUID
}
type GetTagHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h GetTagHandler) Handle(ctx context.Context, query GetTag) (Tag, error) {
	ws, err := h.Repo.Get(ctx, query.WorkspaceId)

	if err != nil {
		return Tag{}, err
	}
	tag, ok := ws.Tags[query.Id]

	if !ok {
		return Tag{}, errors.New("not found")
	}

	t := Tag{
		Id:   query.Id,
		Name: tag,
	}

	return t, nil
}

type ListTagsHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h ListTagsHandler) Handle(ctx context.Context, workspaceID uuid.UUID) ([]Tag, error) {

	tags := make([]Tag, 0)

	ws, err := h.Repo.Get(ctx, workspaceID)

	if err != nil {
		return nil, err
	}

	for tagId, tagName := range ws.Tags {
		tags = append(tags, Tag{
			Id:   tagId,
			Name: tagName,
		})
	}

	return tags, nil
}

type ListWorkspaceResult struct {
	Name string
	Id   uuid.UUID
}

type ListWorkspaceHandler struct {
	Repo workspace.WorkspaceRepository
}

func (h ListWorkspaceHandler) Handle(ctx context.Context) ([]ListWorkspaceResult, error) {
	workspaces, err := h.Repo.ListAll(ctx)

	if err != nil {
		return nil, err
	}

	items := make([]ListWorkspaceResult, 0)

	for _, ws := range workspaces {
		items = append(items, ListWorkspaceResult{
			Name: ws.Name,
			Id:   ws.ID,
		})
	}

	return items, nil
}
