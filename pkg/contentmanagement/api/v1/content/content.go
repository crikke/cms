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
			r.Use(contentIdContext)
			r.Use(api.HandleHttpError)
			r.Get("/", c.GetContent())
			r.Put("/", c.UpdateContent())
		})
	})
}

func contentIdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := ContentID{}
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
			api.WithError(r.Context(), api.GenericError{
				Body: api.ErrorBody{
					Message:   "parameter version is required",
					FieldName: "version",
				},
				StatusCode: http.StatusBadRequest,
			})
			return
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

		req.Version = v

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
	// required:true
	Version int
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
//       200: contentResponse
//       404: genericError
//       400: genericError
func (c contentEndpoint) GetContent() http.HandlerFunc {

	// response is the representation of the content for the Content management API
	// It contains all information of given content for every configured language.
	//
	// swagger:response contentResponse
	type response struct {
		// in: body
		Body query.ContentReadModel
	}

	// swagger:parameters GetContent
	type request struct {
		ContentID
	}

	return func(rw http.ResponseWriter, r *http.Request) {

		var req request

		if r := r.Context().Value(contentkey); r != nil {
			req = r.(request)
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

		response := &response{Body: res}
		data, err := json.Marshal(response)
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
//       201: CreateContentResponse
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

	// swagger:response response
	type _ struct {
		Location string
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
		ContentID
		Language string

		// ! TODO remove swagger ignore
		// swagger:ignore
		Fields []struct {
			Name  string
			Value interface{}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req request

		if r := r.Context().Value(contentkey); r != nil {
			id := r.(ContentID)
			req = request{ContentID: id}
		}

		// ! TODO: this needs to be remade.
		// ! eighter by implementing transactions,
		// ! or not having an array of fields.
		// ! otherwise partial updates can happen.
		for _, field := range req.Fields {

			err := c.app.Commands.UpdateContentField.Handle(r.Context(), command.UpdateField{
				ContentID: req.ID,
				Version:   req.Version,
				Value:     field.Name,
				Language:  req.Language,
				Name:      field.Name,
			})

			if err != nil {
				api.WithError(r.Context(), err)
				return
			}
		}
	}
}

// swagger:route DELETE /content/{id} content DeleteContent
//
// Delete content
//
// Deletes a content node
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
func (c contentEndpoint) DeleteContent() http.HandlerFunc {

	// swagger:parameters request DeleteContent
	type request struct {
		// ! TODO: Split ContentID and ContentVersion
		// ! DeleteContent does not need ContentVersion
		ContentID
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var req request

		if r := r.Context().Value(contentkey); r != nil {
			id := r.(ContentID)
			req = request{ContentID: id}
		}
		err := c.app.Commands.DeleteContent.Handle(
			r.Context(),
			command.DeleteContent{
				ID: req.ID,
			})

		if err != nil {
			api.WithError(r.Context(), err)
		}
	}
}
