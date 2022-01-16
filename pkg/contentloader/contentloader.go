package contentloader

import (
	"context"

	"github.com/crikke/cms/pkg/node"
	"github.com/google/uuid"
)

type Loader interface {
	GetContent(ctx context.Context, id uuid.UUID) (node.Node, error)
	GetChildNodes(ctx context.Context, id uuid.UUID) ([]node.Node, error)
	// GetNode()
}
