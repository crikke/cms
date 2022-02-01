package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetContent(t *testing.T) {

	root := domain.ContentReference{
		ID: uuid.New(),
	}
	loader := mocks.MockLoader{
		Nodes: []domain.Content{
			{
				ID:         root,
				Name:       "root",
				URLSegment: "",
			},
			{
				ParentID:   root.ID,
				Name:       "foo",
				URLSegment: "foo",
				Properties: []domain.Property{
					{
						ID:    uuid.New(),
						Name:  "header",
						Type:  "text",
						Value: "hello world",
					},
				},
			},
		},
	}
	router := gin.Default()
	ContentHandler(router, config.SiteConfiguration{RootPage: loader.Nodes[0].ID.ID}, loader)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/content/foo", nil)

	router.ServeHTTP(w, r)

	assert.Contains(t, w.Body.String(), "hello world")
}
