package routing

import (
	"context"
	"net/http"
	"strings"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/contentloader"
	"github.com/crikke/cms/pkg/node"
)

/*
Routing logic works as following:
	1. Split url path into url segments
	2. While remaining segments is not empty
	3.   Get child nodes from previous matched node through context
	4.   Loop through child nodes and check if node contain this segment then pop matched segment from remainingsegments and
	     set push matched node to context.
*/
func RouteHandler(next http.Handler, contentLoader contentloader.Loader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// TODO: Maybe move localization to its own middleware
		language := r.Header.Get("Accept-Language")
		if language == "" {
			language = config.SiteConfiguration.DefaultLanguage
		}

		ctx = context.WithValue(ctx, config.LanguageKey, language)
		var segments []string
		segments = strings.Split(r.URL.Path, "/")

		// first item is always rootnode
		currentNode, err := contentLoader.GetContent(ctx, config.SiteConfiguration.RootPage)
		if err != nil {
			// TODO: Handle error
			panic(err)
		}

		for i := 1; i < len(segments); i++ {

			nodes, err := contentLoader.GetChildNodes(ctx, currentNode.ID)

			if err != nil {
				// TODO: Handle error
				panic(err)
			}

			match := false
			for _, child := range nodes {

				match, segments = child.Match(ctx, segments)

				if match {
					currentNode = child
					break
				}
			}
		}
		ctx = context.WithValue(ctx, node.NodeKey, currentNode)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
