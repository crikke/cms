//go:build integration

package content

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func Test_CreateAndFetchContent(t *testing.T) {

	// create content definition
	client, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	cd, _ := contentdefinition.NewContentDefinition("test contentdefinition", "test desc")

	cdRepo := contentdefinition.NewContentDefinitionRepository(client)
	contentRepo := domain.NewContentRepository(client)

	cdRepo.CreateContentDefinition(context.Background(), &cd)

	// initialize api endpoint
	ep := NewContentEndpoint(app.App{
		Commands: app.Commands{
			CreateContent: command.CreateContentHandler{
				ContentDefinitionRepository: cdRepo,
				ContentRepository:           contentRepo,
				Factory: domain.Factory{
					Cfg: &siteconfiguration.SiteConfiguration{
						Languages: []language.Tag{
							language.MustParse("sv-SE"),
						},
					}},
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

	type request struct {
		ContentDefinitionId uuid.UUID `json:"contentdefinitionid"`
		ParentId            uuid.UUID `json:"parentid"`
	}

	body := request{
		ContentDefinitionId: cd.ID,
	}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/content", &buf)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	// assert

	// get the created content
	// assert that its contentdefinition id matches the test contentdefinition id
	assert.Equal(t, res.Result().StatusCode, http.StatusCreated)

	location, err := res.Result().Location()
	assert.NoError(t, err)

	req2, err := http.NewRequest(http.MethodGet, location.String(), nil)
	res2 := httptest.NewRecorder()

	r.ServeHTTP(res2, req2)

	actual := &query.ContentReadModel{}

	err = json.NewDecoder(res2.Body).Decode(actual)
	assert.NoError(t, err)

	assert.Equal(t, cd.ID, actual.ContentDefinitionID)
	assert.Equal(t, content.Draft, actual.Status)
}
