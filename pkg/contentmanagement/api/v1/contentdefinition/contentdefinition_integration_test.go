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
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
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
			UpdatePropertyDefinition: command.UpdatePropertyDefinitionHandler{
				Repo: cdRepo,
			},
		},
		Queries: app.Queries{
			GetContentDefinition: query.GetContentDefinitionHandler{
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

	updatePropertyDefinition := func(url url.URL) bool {
		ok := true
		type request struct {
			Name        string
			Description string
			Localized   bool
			Validation  map[string]interface{}
		}

		body := request{
			Name:        "new name",
			Description: "new description",
			Localized:   true,
			Validation: map[string]interface{}{
				validator.RuleRegex: "^[a-z]",
			},
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		ok = ok && assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, url.String(), &buf)
		ok = ok && assert.NoError(t, err)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		return ok
	}

	getContentDefinition := func(url url.URL) bool {
		ok := true
		req, err := http.NewRequest(http.MethodGet, url.String(), nil)
		ok = ok && assert.NoError(t, err)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		actual := contentdefinition.ContentDefinition{}
		json.NewDecoder(res.Body).Decode(&actual)

		expect := contentdefinition.ContentDefinition{
			Name:        "test contentdefinition ",
			Description: "test description",
			Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
				"new name": {
					Description: "new description",
					Localized:   true,
					Type:        "text",
					Validators: map[string]interface{}{
						validator.RuleRequired: false,
						validator.RuleRegex:    "^[a-z]",
					},
				},
			},
		}

		assert.Equal(t, expect.Name, actual.Name)
		assert.Equal(t, expect.Description, actual.Description)
		assert.Equal(t, expect.Propertydefinitions["new name"].Description, actual.Propertydefinitions["new name"].Description)
		assert.Equal(t, expect.Propertydefinitions["new name"].Localized, actual.Propertydefinitions["new name"].Localized)
		assert.Equal(t, expect.Propertydefinitions["new name"].Validators[validator.RuleRegex], actual.Propertydefinitions["new name"].Validators[validator.RuleRegex])
		assert.Equal(t, expect.Propertydefinitions["new name"].Validators[validator.RuleRequired], actual.Propertydefinitions["new name"].Validators[validator.RuleRequired])

		return ok
	}

	t.Run("test create and update content def", func(t *testing.T) {

		var propertyLocation url.URL

		contentLocation, ok := createContentDefinition()
		if ok {
			propertyLocation, ok = createPropertyDefinition(contentLocation)
		}

		ok = ok && updatePropertyDefinition(propertyLocation)
		ok = ok && getContentDefinition(contentLocation)
	})
}
