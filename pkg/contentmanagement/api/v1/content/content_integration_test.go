//go:build integration

package content

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
	"github.com/crikke/cms/pkg/contentmanagement/content"
	domain "github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_CreateAndUpdateNewContent(t *testing.T) {

	// create content definition
	client, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	client.Database("cms").Collection("content").Drop(context.Background())
	client.Database("cms").Collection("contentdefinition").Drop(context.Background())
	cd, _ := contentdefinition.NewContentDefinition("test contentdefinition", "test desc")

	cdRepo := contentdefinition.NewContentDefinitionRepository(client)
	contentRepo := domain.NewContentRepository(client)
	cdRepo.CreateContentDefinition(context.Background(), &cd)

	factory := domain.Factory{
		Cfg: &siteconfiguration.SiteConfiguration{
			Languages: []language.Tag{
				language.MustParse("sv-SE"),
			},
		}}
	// initialize api endpoint
	ep := NewContentEndpoint(app.App{
		Commands: app.Commands{
			CreateContent: command.CreateContentHandler{
				ContentDefinitionRepository: cdRepo,
				ContentRepository:           contentRepo,
				Factory:                     factory,
			},
			UpdateContentFields: command.UpdateContentFieldsHandler{
				ContentRepository:           contentRepo,
				ContentDefinitionRepository: cdRepo,
				Factory:                     factory,
			},
		},
		Queries: app.Queries{
			GetContent: query.GetContentHandler{
				Repo: contentRepo,
			},
		},
	})

	r := chi.NewRouter()
	ep.RegisterEndpoints(r)

	createContent := func(contentdefid uuid.UUID) (url.URL, bool) {
		t.Helper()

		ok := true
		type request struct {
			ContentDefinitionId uuid.UUID `json:"contentdefinitionid"`
		}

		body := request{ContentDefinitionId: contentdefid}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		ok = ok && assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/content", &buf)
		ok = ok && assert.NoError(t, err)

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		ok = ok && assert.Equal(t, res.Result().StatusCode, http.StatusCreated)

		location, err := res.Result().Location()
		ok = ok && assert.NoError(t, err)
		return *location, ok
	}

	getCreatedContent := func(url url.URL) (uuid.UUID, bool) {
		t.Helper()

		ok := true
		req, err := http.NewRequest(http.MethodGet, url.String(), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		actual := &query.ContentReadModel{}

		err = json.NewDecoder(res.Body).Decode(actual)
		ok = ok && assert.NoError(t, err)
		ok = ok && assert.Equal(t, cd.ID, actual.ContentDefinitionID)
		ok = ok && assert.Equal(t, content.Draft, actual.Status)

		return actual.ID, ok
	}

	updateContent := func(contentID uuid.UUID) bool {

		t.Helper()
		ok := true

		type request struct {
			Version  int
			Language string
			Fields   map[string]interface{}
		}

		body := request{
			Version:  0,
			Language: "sv-SE",
			Fields: map[string]interface{}{
				contentdefinition.NameField: "updated content",
			},
		}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/content/%s?l=sv-SE", contentID.String()), &buf)
		ok = ok && assert.NoError(t, err)

		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		ok = ok && assert.Equal(t, len(res.Body.Bytes()), 0)

		return ok
	}

	t.Run("create new content and update it", func(t *testing.T) {
		location, ok := createContent(cd.ID)

		var contentID uuid.UUID
		if ok {
			contentID, ok = getCreatedContent(location)
		}

		ok = ok && updateContent(contentID)
	})
}
