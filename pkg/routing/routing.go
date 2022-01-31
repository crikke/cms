package routing

import (
	"context"
	"net/http"
	"strings"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/loader"
)

type key int

var nodeKey key

/*
Routing logic works as following:
	1. Split url path into url segments
	2. While remaining segments is not empty
	3.   Get child nodes from previous matched node through context
	4.   Loop through child nodes and check if node contain this segment then pop matched segment from remainingsegments and
	     set matchedNode
	5. When done looping through segments, set matchedNode to context
*/
func RoutingHandler(next http.Handler, loader loader.Loader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var segments []string
		segments = strings.Split(r.URL.Path, "/")

		// first item is always rootnode
		currentNode, err := loader.GetContent(r.Context(), config.SiteConfiguration.RootPage)
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

			nodes, err := loader.GetChildNodes(r.Context(), currentNode.ID)

			if err != nil {
				// TODO: Handle error
				panic(err)
			}

			match := false
			for _, child := range nodes {

				match, segments = child.Match(r.Context(), segments)

				if match {
					currentNode = child
					break
				}
			}
		}
		ctx := WithNode(r.Context(), currentNode)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func WithNode(ctx context.Context, node content.Content) context.Context {
	return context.WithValue(ctx, nodeKey, node)
}

func RoutedNode(ctx context.Context) content.Content {
	return ctx.Value(nodeKey).(content.Content)
}
