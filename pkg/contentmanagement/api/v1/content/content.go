package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

var contentkey = key("content")

type contentEndpoint struct {
	app app.App
}

func NewContentEndpoint(app app.App) contentEndpoint {
	return contentEndpoint{app}
}
func (c contentEndpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/content", func(r chi.Router) {

		r.Post("/", c.CreateContent())
		r.Route("/{id}", func(r chi.Router) {
			r.Use(contentCtx)
			r.Use(api.HandleHttpError)
			r.Get("/", c.GetContent())
			r.Put("/", c.UpdateContent())
		})
	})
}

func contentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := ContentRequest{}
		contentID := chi.URLParam(r, "contentid")
		if contentID == "" {

			api.WithError(r.Context(), api.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: api.ErrorBody{
					FieldName: "contentid",
					Message:   "parameter contentid is required",
				},
			})
			return
		}

		cid, err := uuid.Parse(contentID)

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: api.ErrorBody{
					FieldName: "contentid",
					Message:   "bad format",
				},
			})
			return
		}
		version := r.URL.Query().Get("version")

		if version != "" {
			v, err := strconv.Atoi(version)
			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body: api.ErrorBody{
						Message:   "bad formatted version",
						FieldName: "version",
					},
					StatusCode: http.StatusBadRequest,
				})
				return
			}

			req.Version = &v
		}

		req.ID = cid

		ctx := context.WithValue(r.Context(), contentkey, cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ContentID struct {
	// ID
	//
	// in: path
	// required:true
	ID uuid.UUID
	// Version
	//
	// in: query
	// required:false
	Version *int
}

// GetContentResponse is the representation of the content for the Content management API
// It contains all information of given content for every configured language.
//
// swagger:response contentResponse
type GetContentResponse struct {
	// in: body
	Body query.ContentReadModel
}

// swagger:parameters ContentRequest GetContent
type ContentRequest struct {
	ContentID
}

// swagger:route GET /content/{id} content GetContent
//
// Get content by id and optionally version
//
// Gets content by Id. If version is not specified, the published version will be returned.
// If there is no version published, the version with highest version number will be returned
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: contentResponse
//       404: genericError
//       400: genericError
func (c contentEndpoint) GetContent() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var req ContentRequest

		if r := r.Context().Value(contentkey); r != nil {
			req = r.(ContentRequest)
		}

		q := query.GetContent{
			Id:      req.ID,
			Version: req.Version,
		}

		res, err := c.app.Queries.GetContent.Handle(r.Context(), q)

		if err != nil {

			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusNotFound,
			})

			return
		}

		response := &GetContentResponse{Body: res}
		data, err := json.Marshal(response)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write(data)
	}
}

// swagger:parameters CreateContentRequest CreateContent
type CreateContentRequest struct {
	// Contentdefinition ID
	// in: body
	// required: true
	ContentDefinitionId uuid.UUID `json:"contentdefinitionid"`
	// ParentId
	// in: body
	ParentId uuid.UUID `json:"parentid"`
}

// swagger:response CreateContentResponse
type CreateContentResponse struct {
	Location string
}

// swagger:route POST /content content CreateContent
//
// Create new content
//
// Creates new content node under the parent. The content is created from the specified contentdefinition
// which acts as a template, containing what properties to create & their validation.
//
//     Consumes:
//	   - application/json
//     Produces:
//     - application/json
//
//     Responses:
//       201: CreateContentResponse
//		 400: genericError
func (c contentEndpoint) CreateContent() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		req := &CreateContentRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		id, err := c.app.Commands.CreateContent.Handle(r.Context(),
			command.CreateContent{
				ContentDefinitionId: req.ContentDefinitionId,
				ParentID:            req.ParentId,
			},
		)
		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			})
			return

		}

		url := r.URL.String()
		rw.Header().Add("Location", fmt.Sprintf("%s/%s", url, id.String()))
		rw.WriteHeader(http.StatusCreated)
	}
}

// swager:route PUT /content/{id} content UpdateContent
//
// Update content
//
// Updates content node
//
//		Consumes:
//		- application/json
//		Produces:
//		- application/json
//
//		Responses:
//		  200: OK
//		  404: genericError
//        400: genericError
func (c contentEndpoint) UpdateContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
