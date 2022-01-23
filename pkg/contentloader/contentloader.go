package contentloader

import (
	"context"

	"github.com/crikke/cms/pkg/content"
	"github.com/google/uuid"
)

type Loader interface {
	GetContent(ctx context.Context, id uuid.UUID) (content.Content, error)
	GetChildNodes(ctx context.Context, id uuid.UUID) ([]content.Content, error)
	// GetNode()
}
