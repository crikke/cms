//go:build integration

package contentdefinition

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_CreateContentDefinition(t *testing.T) {
	client, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	client.Database("cms").Collection("content").Drop(context.Background())
	client.Database("cms").Collection("contentdefinition").Drop(context.Background())

	cdRepo := contentdefinition.NewContentDefinitionRepository(client)

	// initialize api endpoint
	ep := NewContentDefinitionEndpoint(app.App{
		Commands: app.Commands{
			CreateContentDefinition: command.CreateContentDefinitionHandler{
				Repo: cdRepo,
			},
			UpdateContentDefinition: command.UpdateContentDefinitionHandler{
				Repo: cdRepo,
			},
			CreatePropertyDefinition: command.CreatePropertyDefinitionHandler{
				Repo: cdRepo,
			},
		},
	})
	r := chi.NewRouter()
	ep.RegisterEndpoints(r)

	createContentDefinition := func() (url.URL, bool) {
		t.Helper()

		ok := true

		type request struct {
			Name        string
			Description string
		}

		body := request{
			Name:        "test contentdefinition ",
			Description: "test description",
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		ok = ok && assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/contentdefinitions", &buf)
		ok = ok && assert.NoError(t, err)

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		ok = ok && assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

		location, err := res.Result().Location()
		ok = ok && assert.NoError(t, err)
		return *location, ok
	}

	createPropertyDefinition := func(url url.URL) (url.URL, bool) {
		ok := true
		type request struct {
			Name        string
			Description string
			Type        string
		}

		body := request{
			Name:        "test_property",
			Description: "test_description",
			Type:        "text",
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		ok = ok && assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/propertydefinitions", url.String()), &buf)
		ok = ok && assert.NoError(t, err)

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		ok = ok && assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		location, err := res.Result().Location()
		ok = ok && assert.NoError(t, err)
		return *location, ok
	}

	t.Run("test create and update content def", func(t *testing.T) {

		// var propertyLocation url.URL

		contentLocation, ok := createContentDefinition()
		if ok {
			_, ok = createPropertyDefinition(contentLocation)
		}

	})
}
