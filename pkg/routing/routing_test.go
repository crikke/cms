package routing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatchRoute(t *testing.T) {

	a := uuid.New()
	b := uuid.New()
	c := uuid.New()

	nodes := []domain.Content{
		{
			ID:         domain.ContentReference{ID: a},
			URLSegment: "a",
		},
		{
			ID:         domain.ContentReference{ID: b},
			ParentID:   a,
			URLSegment: "b",
		},
		{
			ParentID:   b,
			ID:         domain.ContentReference{ID: c},
			URLSegment: "c",
		},
	}

	tests := []struct {
		description  string
		url          string
		language     string
		expectedNode domain.Content
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

		t.Run(test.description, func(t *testing.T) {
			router := gin.Default()

			router.GET("/*nodes", RoutingHandler(config.Configuration{}, mockLoader{
				nodes: nodes,
			}), func(c *gin.Context) {
				assert.Equal(t, test.expectedNode.ID, RoutedNode(*c).ID, test.description)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", test.url, nil)

			router.ServeHTTP(w, r)
		})
	}
}

type mockLoader struct {
	nodes []domain.Content
}

func (m mockLoader) GetContent(ctx context.Context, id domain.ContentReference) (domain.Content, error) {

	for _, node := range m.nodes {
		if node.ID == id {
			return node, nil
		}
	}
	return domain.Content{}, nil
}

func (m mockLoader) GetChildNodes(ctx context.Context, id domain.ContentReference) ([]domain.Content, error) {

	result := []domain.Content{}

	for _, node := range m.nodes {
		if node.ParentID == id.ID {
			result = append(result, node)
		}
	}
	return result, nil
}
