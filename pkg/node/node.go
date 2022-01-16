package node

import (
	"context"
	"strings"

	"github.com/crikke/cms/pkg/config"
	"github.com/google/uuid"
)

type nodeKey string

var NodeKey nodeKey = "context-nodes"

type Node struct {
	URLSegment map[string]string
	ID         uuid.UUID
	Parent     *Node
}

// Checks if a node matches a urlsegment
func (n Node) Match(ctx context.Context, remaining []string) (match bool, segments []string) {
	segments = remaining
	lang := ctx.Value(config.LanguageKey).(string)

	if len(remaining) == 0 {
		match = false
		return
	}

	segment := remaining[0]
	nodeSegment, exist := n.URLSegment[lang]
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
