package contentdefinition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crikke/cms/cmd/contentmanagement/api/handlers"
	"github.com/crikke/cms/cmd/contentmanagement/api/models"
	"github.com/crikke/cms/cmd/contentmanagement/app"
	"github.com/crikke/cms/cmd/contentmanagement/app/command"
	"github.com/crikke/cms/cmd/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

var contentKey = key("cid")
var propertyKey = key("pid")

type endpoint struct {
	app app.App
}

func NewContentDefinitionRoute(app app.App) http.Handler {

	c := endpoint{app: app}
	r := chi.NewRouter()

	r.Get("/", c.ListContentDefinitions())
	r.Post("/", c.CreateContentDefinition())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return contentDefinitionIdContext(h, "id", contentKey)
		})
		r.Get("/", c.GetContentDefinition())
		r.Delete("/", c.DeleteContentDefinition())
		r.Put("/", c.UpdateContentDefinition())
	})
	return r
}

func contentDefinitionIdContext(next http.Handler, param string, key key) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, param)
		if id == "" {

			models.WithError(r.Context(), models.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: models.ErrorBody{
					FieldName: "contentid",
					Message:   "parameter contentid is required",
				},
			})
			return
		}

		uid, err := uuid.Parse(id)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: models.ErrorBody{
					FieldName: "contentid",
					Message:   "bad format",
				},
			})
			return
		}

		ctx := context.WithValue(r.Context(), key, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withID(ctx context.Context) uuid.UUID {

	var id uuid.UUID

	if r := ctx.Value(contentKey); r != nil {
		id = r.(uuid.UUID)
	}

	return id
}

func withPID(ctx context.Context) uuid.UUID {

	var id uuid.UUID

	if r := ctx.Value(propertyKey); r != nil {
		id = r.(uuid.UUID)
	}

	return id
}

// CreateContentDefinition 		godoc
// @Summary 					Creates a new content definition
// @Description 				Creates a new contentdefinition. The contentdefinition
// @Description 				acts as a template for creating new content,
// @Description 				containing what properties to create & their validation.
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param						body		body	ContentDefinitionBody	true 	"request body"
// @Success						201			{object}	models.OKResult
// @Header						201			{string}	Location
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/workspaces/{workspace}/contentdefinitions [post]
func (c endpoint) CreateContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := &ContentDefinitionBody{}
		ws := handlers.WithWorkspace(r.Context())
		err := json.NewDecoder(r.Body).Decode(req)

		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := c.app.Commands.CreateContentDefinition.Handle(r.Context(), command.CreateContentDefinition{
			Name:        req.Name,
			Description: req.Description,
			WorkspaceId: ws.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := r.URL.String()
		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, id.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

// GetContentDefinition 		godoc
// @Summary 					Gets a content definition
// @Description 				Gets a content definition by ID
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Success						200			{object}	contentdefinition.ContentDefinition
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/workspaces/{workspace}/contentdefinitions/{id} [get]
func (c endpoint) GetContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		cd, err := c.app.Queries.GetContentDefinition.Handle(r.Context(), query.GetContentDefinition{ID: id, WorkspaceID: ws.ID})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		bytes, err := json.Marshal(&cd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(bytes)
	}
}

// ListContentDefinitions 		godoc
// @Summary 					Get all content definitions
// @Description 				Gets all existing contentdefinitions
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Success						200			{object}	[]query.ListContentDefinitionModel
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/workspaces/{workspace}/contentdefinitions [get]
func (c endpoint) ListContentDefinitions() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())
		cd, err := c.app.Queries.ListContentDefinitions.Handle(r.Context(), query.ListContentDefinition{WorkspaceID: ws.ID})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(cd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(data)
	}
}

// DeleteContentDefinition 		godoc
// @Summary 					Delete a content definition
// @Description 				Delete a content definition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Success						200			{object}	models.OKResult
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/workspaces/{workspace}/contentdefinitions/{id} [delete]
func (c endpoint) DeleteContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		err := c.app.Commands.DeleteContentDefinition.Handle(r.Context(), command.DeleteContentDefinition{ID: id, WorkspaceId: ws.ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// UpdateContentDefinition 		godoc
// @Summary 					Updates a contentdefinition
// @Description 				Updates a contentdefinition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param						body		body	ContentDefinitionBody	true 	"request body"
// @Success						200			{object}	models.OKResult
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/workspaces/{workspace}/contentdefinitions/{id} [put]
func (c endpoint) UpdateContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		body := &ContentDefinitionBody{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c.app.Commands.UpdateContentDefinition.Handle(r.Context(), command.UpdateContentDefinition{
			ContentDefinitionID: id,
			Name:                body.Name,
			Description:         body.Description,
			WorkspaceId:         ws.ID,
			PropertyDefinitions: body.PropertyDefinitions,
		})
	}
}
