package routing

import (
	"context"
	"fmt"
	"strings"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/locale"
	"github.com/crikke/cms/pkg/services/loader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const nodeKey = "nodeKey"

/*
Routing logic works as following:
	1. Split url path into url segments
	2. While remaining segments is not empty
	3.   Get child nodes from previous matched node through context
	4.   Loop through child nodes and check if node contain this segment then pop matched segment from remainingsegments and
	     set matchedNode
	5. When done looping through segments, set matchedNode to context
*/
func RoutingHandler(cfg config.SiteConfiguration, loader loader.Loader) gin.HandlerFunc {

	return func(c *gin.Context) {
		var segments []string
		segments = strings.Split(c.Param("node"), "/")

		// first item is always rootnode
		locale := locale.FromContext(*c)
		contentReference := domain.ContentReference{
			ID:     cfg.RootPage,
			Locale: &locale,
		}

		currentNode, err := loader.GetContent(c.Request.Context(), contentReference)
		if err != nil {
			// TODO: Handle error
			panic(err)
		}

		for len(segments) > 0 {

			// if empty segment remove it
			if segments[0] == "" {
				segments = segments[1:]
				continue
			}

			nodes, err := loader.GetChildNodes(c.Request.Context(), contentReference)

			if err != nil {
				// TODO: Handle error
				panic(err)
			}

			m := false
			for _, child := range nodes {

				m, segments = match(child, segments)

				if m {
					currentNode = child
					contentReference = child.ID
					break
				}
			}
		}
		c.Set(nodeKey, currentNode)
		c.Next()
	}
}

func GenerateUrl(ctx context.Context, loader loader.Loader, contentReference domain.ContentReference) (string, error) {

	c, err := loader.GetContent(ctx, contentReference)

	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("/%s/", c.URLSegment)

	next := c

	for (next.ParentID != uuid.UUID{}) {

		pRef := domain.ContentReference{
			ID:     next.ParentID,
			Locale: contentReference.Locale,
		}

		c, err = loader.GetContent(ctx, pRef)
		if err != nil {
			return "", err

		}

		path = fmt.Sprintf("/%s%s", c.URLSegment, path)

		next = c
	}

	return path, nil
}

func match(c domain.Content, remaining []string) (match bool, segments []string) {
	segments = remaining

	segment := remaining[0]

	nodeSegment := c.URLSegment

	match = strings.EqualFold(segment, nodeSegment)

	if match {
		// pop matched segment
		segments = segments[1:]
	}
	return
}

func RoutedNode(c gin.Context) domain.Content {

	node, exist := c.Get(nodeKey)

	if !exist {
		node = domain.Content{}
	}

	return node.(domain.Content)
}
