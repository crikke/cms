package routing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/content"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatchRoute(t *testing.T) {

	a := uuid.New()
	b := uuid.New()
	c := uuid.New()

	nodes := []content.Content{
		{
			ID:         content.ContentReference{ID: a},
			URLSegment: "a",
		},
		{
			ID:         content.ContentReference{ID: b},
			ParentID:   a,
			URLSegment: "b",
		},
		{
			ParentID:   b,
			ID:         content.ContentReference{ID: c},
			URLSegment: "c",
		},
	}

	tests := []struct {
		description  string
		url          string
		language     string
		expectedNode content.Content
	}{
		{
			description:  "node matched with sv language",
			url:          "/a/b/c",
			expectedNode: nodes[2],
		},
		{
			description:  "multiple slashes in path ",
			url:          "///a//b////c",
			expectedNode: nodes[2],
		},
		{
			description:  "path ends with '/' ",
			url:          "/a/b/c/",
			expectedNode: nodes[2],
		},
	}

	for _, test := range tests {

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routedNode := RoutedNode(r.Context())
			assert.Equal(t, test.expectedNode.ID, routedNode.ID)
		})

		t.Run(test.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			handler := RoutingHandler(testHandler,
				config.Configuration{},
				mockLoader{
					nodes: nodes,
				},
			)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

type mockLoader struct {
	nodes []content.Content
}

func (m mockLoader) GetContent(ctx context.Context, id content.ContentReference) (content.Content, error) {

	for _, node := range m.nodes {
		if node.ID == id {
			return node, nil
		}
	}
	return content.Content{}, nil
}

func (m mockLoader) GetChildNodes(ctx context.Context, id content.ContentReference) ([]content.Content, error) {

	result := []content.Content{}

	for _, node := range m.nodes {
		if node.ParentID == id.ID {
			result = append(result, node)
		}
	}
	return result, nil
}
