package node

import (
	"context"
	"strings"

	"github.com/crikke/cms/pkg/locale"
	"github.com/google/uuid"
)

type key int

var nodeKey key

type Node struct {
	// A node are required to have a localized URLSegment for each configured locale.
	URLSegment map[string]string
	ID         uuid.UUID
	ParentID   uuid.UUID
}

// Checks if a node matches a urlsegment
func (n Node) Match(ctx context.Context, remaining []string) (match bool, segments []string) {
	segments = remaining
	lang := locale.FromContext(ctx)

	segment := remaining[0]
	nodeSegment, exist := n.URLSegment[lang.String()]
	// this should never fail because if URLSegment is not set explicitly, the localized name will be used as URLSegment
	if !exist {
		match = false
		return
	}

	match = strings.EqualFold(segment, nodeSegment)

	if match {
		// pop matched segment
		segments = segments[1:]
	}
	return
}

func WithNode(ctx context.Context, node Node) context.Context {
	return context.WithValue(ctx, nodeKey, node)
}

func RoutedNode(ctx context.Context) Node {
	return ctx.Value(nodeKey).(Node)
}
