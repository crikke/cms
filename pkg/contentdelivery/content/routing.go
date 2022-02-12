package content

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/crikke/cms/pkg/locale"
	"github.com/crikke/cms/pkg/siteconfiguration"
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
// TODO move to api and not as middleware
func RoutingHandler(cfg *siteconfiguration.SiteConfiguration, repo ContentRepository) gin.HandlerFunc {

	return func(c *gin.Context) {
		var segments []string
		segments = strings.Split(c.Param("node"), "/")

		// first item is always rootnode
		locale := locale.FromContext(c)
		contentReference := ContentReference{
			ID:     cfg.RootPage,
			Locale: &locale,
		}

		currentNode, err := repo.GetContent(c.Request.Context(), contentReference)
		if err != nil {
			c.Error(err)
			return
		}

		for len(segments) > 0 {

			// if empty segment remove it
			if segments[0] == "" {
				segments = segments[1:]
				continue
			}

			nodes, err := repo.GetChildren(c.Request.Context(), contentReference)

			if err != nil {
				c.Error(err)
				return
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

			if !m {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		c.Set(nodeKey, currentNode)
		c.Next()
	}
}

func GenerateUrl(ctx context.Context, repo ContentRepository, contentReference ContentReference) (string, error) {

	c, err := repo.GetContent(ctx, contentReference)

	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("/%s/", c.URLSegment)

	next := c

	for (next.ParentID != uuid.UUID{}) {

		pRef := ContentReference{
			ID:     next.ParentID,
			Locale: contentReference.Locale,
		}

		c, err = repo.GetContent(ctx, pRef)
		if err != nil {
			return "", err

		}

		path = fmt.Sprintf("/%s%s", c.URLSegment, path)

		next = c
	}

	return path, nil
}

func match(c Content, remaining []string) (match bool, segments []string) {
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

func RoutedNode(c *gin.Context) Content {

	node, exist := c.Get(nodeKey)

	if !exist {
		node = Content{}
	}

	return node.(Content)
}
