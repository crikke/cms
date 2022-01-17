package routing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/locale"
	"github.com/crikke/cms/pkg/node"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestMatchRoute(t *testing.T) {

	a := uuid.New()
	b := uuid.New()
	c := uuid.New()

	nodes := []node.Node{
		{
			ID: a,
			URLSegment: map[string]string{
				"sv": "a",
				"en": "aa",
			},
		},
		{
			ID:       b,
			ParentID: a,
			URLSegment: map[string]string{
				"sv": "b",
				"en": "bb",
			},
		},
		{
			ParentID: b,
			ID:       c,
			URLSegment: map[string]string{
				"sv": "c",
				"en": "cc",
			},
		},
	}

	tests := []struct {
		description  string
		url          string
		language     string
		expectedNode node.Node
	}{
		{
			description:  "node matched with sv language",
			url:          "/a/b/c",
			language:     "sv",
			expectedNode: nodes[2],
		},
		{
			description:  "node matched with en language",
			url:          "/aa/bb/cc",
			language:     "en",
			expectedNode: nodes[2],
		},
		{
			description:  "multiple slashes in path ",
			url:          "///a//b////c",
			language:     "sv",
			expectedNode: nodes[2],
		},
		{
			description:  "path ends with '/' ",
			url:          "/a/b/c/",
			language:     "sv",
			expectedNode: nodes[2],
		},
	}

	for _, test := range tests {

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routedNode := node.RoutedNode(r.Context())
			assert.Equal(t, test.expectedNode.ID, routedNode.ID)
		})

		t.Run(test.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			ctx := locale.WithLocale(req.Context(), language.MustParse(test.language))
			req = req.WithContext(ctx)
			handler := Handler(testHandler, mockLoader{
				nodes: nodes,
			},
			)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

type mockLoader struct {
	nodes []node.Node
}

func (m mockLoader) GetContent(ctx context.Context, id uuid.UUID) (node.Node, error) {

	for _, node := range m.nodes {
		if node.ID == id {
			return node, nil
		}
	}
	return node.Node{}, nil
}

func (m mockLoader) GetChildNodes(ctx context.Context, id uuid.UUID) ([]node.Node, error) {

	result := []node.Node{}

	for _, node := range m.nodes {
		if node.ParentID == id {
			result = append(result, node)
		}
	}
	return result, nil
}
