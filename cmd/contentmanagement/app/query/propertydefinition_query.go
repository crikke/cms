package query

import (
	"context"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/google/uuid"
)

type GetContentDefinition struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
}

type GetContentDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h GetContentDefinitionHandler) Handle(ctx context.Context, query GetContentDefinition) (contentdefinition.ContentDefinition, error) {

	return h.Repo.GetContentDefinition(ctx, query.ID, query.WorkspaceID)
}

type GetPropertyDefinition struct {
	ContentDefinitionID  uuid.UUID
	WorkspaceID          uuid.UUID
	PropertyDefinitionID uuid.UUID
}

type GetPropertyDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h GetPropertyDefinitionHandler) Handle(ctx context.Context, query GetPropertyDefinition) (contentdefinition.PropertyDefinition, error) {

	return h.Repo.GetPropertyDefinition(ctx, query.ContentDefinitionID, query.PropertyDefinitionID, query.WorkspaceID)
}

type ListContentDefinition struct {
	WorkspaceID uuid.UUID
}

type ListContentDefinitionModel struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type ListContentDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h ListContentDefinitionHandler) Handle(ctx context.Context, query ListContentDefinition) ([]ListContentDefinitionModel, error) {

	items, err := h.Repo.ListContentDefinitions(ctx, query.WorkspaceID)

	if err != nil {
		return nil, err
	}

	result := make([]ListContentDefinitionModel, 0)

	for _, item := range items {
		result = append(result, ListContentDefinitionModel{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return result, nil
}
