package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

var contentKey = key("content")
var versionKey = key("version")

type contentEndpoint struct {
	app app.App
}

func NewContentEndpoint(app app.App) contentEndpoint {
	return contentEndpoint{app}
}
func (c contentEndpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/content", func(r chi.Router) {

		r.Get("/", c.ListContent())
		r.Post("/", c.CreateContent())
		r.Route("/{id}", func(r chi.Router) {
			r.Use(contentIdContext)
			r.Use(contentVersionContext)
			r.Use(api.HandleHttpError)
			r.Get("/", c.GetContent())
			r.Put("/", c.UpdateContent())
			r.Delete("/", c.ArchiveContent())
		})
	})
}

func contentIdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := ContentID{}
		contentID := chi.URLParam(r, "id")
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

		req.ID = cid

		ctx := context.WithValue(r.Context(), contentKey, cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contentVersionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		version := r.URL.Query().Get("version")

		if version == "" {
			// api.WithError(r.Context(), api.GenericError{
			// 	Body: api.ErrorBody{
			// 		Message:   "parameter version is required",
			// 		FieldName: "version",
			// 	},
			// 	StatusCode: http.StatusBadRequest,
			// })
			next.ServeHTTP(w, r)
		}

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

		ctx := context.WithValue(r.Context(), versionKey, v)
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
	// required:true
	Version int
}

// swagger:route GET /content content ListContent
//
// Gets all content
//
// Gets all content
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: body:[]ContentListReadModel
//       404: genericError
//       400: genericError
func (c contentEndpoint) ListContent() http.HandlerFunc {

	// swagger:parameters ListContent
	type _ struct {
		// cid
		//
		// in: query
		cid []uuid.UUID
	}

	return func(rw http.ResponseWriter, r *http.Request) {

		val, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {

			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusNotFound,
			})

			return
		}

		q := query.ListContent{
			ContentDefinitionIDs: make([]uuid.UUID, 0),
		}

		if ids, ok := val["cid"]; ok {
			for _, id := range ids {
				uid, _ := uuid.Parse(id)

				q.ContentDefinitionIDs = append(q.ContentDefinitionIDs, uid)
			}
		}

		res, err := c.app.Queries.ListContent.Handle(r.Context(), q)

		if err != nil {

			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusNotFound,
			})

			return
		}

		data, err := json.Marshal(&res)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write(data)
	}
}

// swagger:route GET /content/{id} content GetContent
//
// Get content by id and optionally version
//
// Gets content by Id.
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: body:Contentresponse
//       404: genericError
//       400: genericError
func (c contentEndpoint) GetContent() http.HandlerFunc {

	// swagger:parameters GetContent
	type request struct {
		ContentID
	}

	return func(rw http.ResponseWriter, r *http.Request) {

		req := request{}

		if r := r.Context().Value(contentKey); r != nil {
			req.ContentID.ID = r.(uuid.UUID)
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

		data, err := json.Marshal(&res)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write(data)
	}
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
//       201: Location
//		 400: genericError
func (c contentEndpoint) CreateContent() http.HandlerFunc {

	// swagger:parameters CreateContent
	type request struct {
		// Contentdefinition ID
		// in: body
		// required: true
		ContentDefinitionId uuid.UUID `json:"contentdefinitionid"`
		// ParentId
		// in: body
		ParentId uuid.UUID `json:"parentid"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {

		req := &request{}

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

// swagger:route PUT /content/{id} content UpdateContent
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

	// swagger:parameters UpdateContent
	type request struct {
		// ID
		//
		// in: path
		// required:true
		ID uuid.UUID

		// in:body
		Body struct {
			// Version
			// required:true
			Version int
			// Language
			// required:true
			Language string
			// Properties
			// required:true
			Fields map[string]interface{}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := request{}

		if r := r.Context().Value(contentKey); r != nil {
			req.ID = r.(uuid.UUID)
		}

		bod := &struct {
			Version  int
			Language string
			Fields   map[string]interface{}
		}{}

		err := json.NewDecoder(r.Body).Decode(bod)
		if err != nil {
			api.WithError(r.Context(), err)
			return
		}

		req.Body = *bod

		err = c.app.Commands.UpdateContentFields.Handle(r.Context(), command.UpdateContentFields{
			ContentID: req.ID,
			Version:   req.Body.Version,
			Language:  req.Body.Language,
			Fields:    req.Body.Fields,
		})

		if err != nil {
			api.WithError(r.Context(), err)
			return
		}
	}
}

// swagger:route DELETE /content/{id} content ArchiveContent
//
// Archives content
//
// Archives content with ID
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
func (c contentEndpoint) ArchiveContent() http.HandlerFunc {

	// swagger:parameters request DeleteContent
	type _ struct {
		// ID
		//
		// in: path
		// required:true
		ID uuid.UUID
	}

	return func(w http.ResponseWriter, r *http.Request) {

		id := uuid.UUID{}

		if r := r.Context().Value(contentKey); r != nil {
			uid := r.(uuid.UUID)
			id = uid
		}
		err := c.app.Commands.ArchiveContent.Handle(
			r.Context(),
			command.ArchiveContent{
				ID: id,
			})

		if err != nil {
			api.WithError(r.Context(), err)
		}
	}
}
