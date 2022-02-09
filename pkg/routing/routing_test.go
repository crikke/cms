package routing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/mocks"
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

			router.GET("/*node", RoutingHandler(domain.SiteConfiguration{}, mocks.MockLoader{
				Nodes: nodes,
			}), func(c *gin.Context) {
				assert.Equal(t, test.expectedNode.ID, RoutedNode(*c).ID, test.description)
			})

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", test.url, nil)

			router.ServeHTTP(w, r)
		})
	}
}

func Test_GenerateUrl(t *testing.T) {
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
		description string
		node        domain.Content
		expected    string
	}{
		{
			description: "generate url node c",
			expected:    "/a/b/c/",
			node:        nodes[2],
		},
		{
			description: "generate url node b ",
			expected:    "/a/b/",
			node:        nodes[1],
		},
	}

	for _, test := range tests {

		t.Run(test.description, func(t *testing.T) {

			mock := mocks.MockLoader{
				Nodes: nodes,
			}

			result, err := GenerateUrl(context.TODO(), mock, test.node.ID)

			assert.NoError(t, err)

			assert.Equal(t, test.expected, result)
		})
	}
}
