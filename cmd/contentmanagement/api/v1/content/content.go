package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/crikke/cms/cmd/contentmanagement/api/handlers"
	"github.com/crikke/cms/cmd/contentmanagement/api/models"
	"github.com/crikke/cms/cmd/contentmanagement/app"
	"github.com/crikke/cms/cmd/contentmanagement/app/command"
	"github.com/crikke/cms/cmd/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

var contentKey = key("content")
var versionKey = key("version")

type contentEndpoint struct {
	app app.App
}

func NewContentRoute(app app.App) http.Handler {
	r := chi.NewRouter()
	c := contentEndpoint{app}

	r.Get("/", c.ListContent())
	r.Post("/", c.CreateContent())
	r.Route("/{id}", func(r chi.Router) {
		r.Use(contentIdContext)
		r.Put("/", c.UpdateContent())
		r.Delete("/", c.ArchiveContent())

		r.Route("/", func(r chi.Router) {
			r.Use(contentVersionContext)
			r.Get("/", c.GetContent())
		})

		r.Post("/publish", c.PublishContent())
	})
	return r
}

func contentIdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentID := chi.URLParam(r, "id")
		if contentID == "" {

			models.WithError(r.Context(), models.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: models.ErrorBody{
					FieldName: "contentid",
					Message:   "parameter contentid is required",
				},
			})
			return
		}

		cid, err := uuid.Parse(contentID)

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

		ctx := context.WithValue(r.Context(), contentKey, cid)
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

func withVersion(ctx context.Context) int {
	var version int

	if r := ctx.Value(versionKey); r != nil {
		i := r.(int)
		version = i
	}

	return version
}

func contentVersionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		version := r.URL.Query().Get("version")

		if version == "" {
			models.WithError(r.Context(), models.GenericError{
				Body: models.ErrorBody{
					Message:   "version is required",
					FieldName: "version",
				},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		v, err := strconv.Atoi(version)
		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body: models.ErrorBody{
					Message:   "bad formatted version",
					FieldName: "version",
				},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		ctx := context.WithValue(r.Context(), versionKey, v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ListContent 		godoc
// @Summary 		List all content
// @Description 	list all content
// @Tags 			content
// @Accept 			json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Produces 		json
// @Param			cid			query	[]string	true 	"uuid formatted ID." format(uuid)
// @Param			tag			query	[]string	true 	"tag id"
// @Success			200			{object}	[]query.ContentListReadModel
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content [get]
func (c contentEndpoint) ListContent() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())

		val, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q := query.ListContent{
			ContentDefinitionIDs: make([]uuid.UUID, 0),
			WorkspaceId:          ws.ID,
		}

		if ids, ok := val["cid"]; ok {
			for _, id := range ids {
				uid, _ := uuid.Parse(id)

				q.ContentDefinitionIDs = append(q.ContentDefinitionIDs, uid)
			}
		}

		if tags, ok := val["tag"]; ok {
			q.Tags = append(q.Tags, tags...)
		}

		res, err := c.app.Queries.ListContent.Handle(r.Context(), q)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data, err := json.Marshal(&res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(data)
	}
}

// GetContent 		godoc
// @Summary 		Get content by id
// @Description 	Get content by id and optionally version
// @Tags 			content
// @Accept 			json
// @Produces 		json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			version		query	int		false 	"content version"
// @Success			200			{object}	query.ContentReadModel
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content/{id} [get]
func (c contentEndpoint) GetContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		version := withVersion(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		q := query.GetContent{
			Id:          id,
			Version:     version,
			WorkspaceId: ws.ID,
		}

		res, err := c.app.Queries.GetContent.Handle(r.Context(), q)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data, err := json.Marshal(&res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(data)
	}
}

// CreateContent 	godoc
// @Summary 		Create new content
// @Description 	Creates new content basen on a contentdefinition
// @Tags 			content
// @Accept 			json
// @Produces 		json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param			contentdefinitionid	body CreateContentRequest true "contentdefinitionid"
// @Success						201			{object}	content.Content
// @Header						201			{string}	Location
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content [post]
func (c contentEndpoint) CreateContent() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		body := &CreateContentRequest{}
		ws := handlers.WithWorkspace(r.Context())
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cid, err := uuid.Parse(body.ContentDefinitionId.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := c.app.Commands.CreateContent.Handle(r.Context(),
			command.CreateContent{
				ContentDefinitionId: cid,
				WorkspaceId:         ws.ID,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		created, err := c.app.Queries.GetContent.Handle(r.Context(), query.GetContent{
			Id:          id,
			Version:     0,
			WorkspaceId: ws.ID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		data, err := json.Marshal(&created)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		url := r.URL.String()
		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, id.String()))
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

// UpdateContent 	godoc
// @Summary 		Update content
// @Description 	Update content
// @Tags 			content
// @Accept 			json
// @Produces 		json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			requestbody	body UpdateContentRequestBody true "body"
// @Success			200		{object}		models.OKResult
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content/{id} [put]
func (c contentEndpoint) UpdateContent() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		body := &UpdateContentRequestBody{}

		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.app.Commands.UpdateContentFields.Handle(r.Context(), command.UpdateContentFields{
			ContentID: id,
			Version:   body.Version,
			Language:  body.Language,
			Fields:    body.Fields,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

// ArchivesContent 	godoc
// @Summary 		Archives content
// @Description 	Archives content with ID
// @Tags 			content
// @Accept 			json
// @Produces 		json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Success			200		{object}		OKResult
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content/{id} [delete]
func (c contentEndpoint) ArchiveContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		err := c.app.Commands.ArchiveContent.Handle(
			r.Context(),
			command.ArchiveContent{
				ID:          id,
				WorkspaceId: ws.ID,
			})

		if err != nil {
			models.WithError(r.Context(), err)
		}
	}
}

// PublishContent 	godoc
// @Summary 		Publishes content
// @Description 	Publishes content with ID
// @Tags 			content
// @Accept 			json
// @Produces 		json
// @Param			workspace	path	string	true 	"uuid formatted ID." format(uuid)
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			version		query	int		true 	"content version"
// @Success			200			{object}		OKResult
// @Failure			default		{object}		models.GenericError
// @Router			/contentmanagement/workspaces/{workspace}/content/{id}/publish [post]
func (c contentEndpoint) PublishContent() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		version := withVersion(r.Context())
		ws := handlers.WithWorkspace(r.Context())

		err := c.app.Commands.PublishContent.Handle(
			r.Context(),
			command.PublishContent{
				ContentID:   id,
				Version:     version,
				WorkspaceId: ws.ID,
			})

		if err != nil {
			models.WithError(r.Context(), err)
		}
	}
}
