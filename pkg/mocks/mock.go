package mocks

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/domain"
	"github.com/google/uuid"
)

type MockLoader struct {
	Nodes []domain.Content
}

func (m MockLoader) GetContent(ctx context.Context, id domain.ContentReference) (domain.Content, error) {

	for _, node := range m.Nodes {
		if node.ID.ID == id.ID {
			return node, nil
		}
	}
	return domain.Content{}, nil
}

func (m MockLoader) GetChildNodes(ctx context.Context, id domain.ContentReference) ([]domain.Content, error) {

	if id.ID == uuid.Nil {
		return nil, errors.New("empty uuid")
	}

	result := []domain.Content{}

	for _, node := range m.Nodes {
		if node.ParentID == id.ID {
			result = append(result, node)
		}
	}
	return result, nil
}
