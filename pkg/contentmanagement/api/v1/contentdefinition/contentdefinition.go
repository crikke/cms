package contentdefinition

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contentEndpoint struct {
	app app.App
}

func NewContentDefinitionEndpoint(app app.App) contentEndpoint {
	return contentEndpoint{app}
}

func (c contentEndpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/contentdefinition", func(r chi.Router) {

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", c.GetContentDefinition())
		})
		r.Post("/", c.CreateContentDefinition())
	})
}

// swagger:route POST /contentdefinition contentdefinition CreateContentDefinition
//
// Creates a new content definition
//
// Creates a new contentdefinition. The contentdefinition
// acts as a template for creating new content,
// containing what properties to create & their validation.
//
//     Consumes:
//	   - application/json
//
//     Responses:
//       201: CreateContentDefinitionResponse
//		 400: genericError
//		 500: genericError
func (c contentEndpoint) CreateContentDefinition() http.HandlerFunc {

	// swagger:parameters request CreateContentDefinition
	type request struct {
		// Content definition Name
		// in:body
		Name string
		// Content definition description
		// in:body
		Description string
	}

	// swagger:response CreateContentDefinitionResponse
	type _ struct {
		Location string
	}
	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		err := json.NewDecoder(r.Body).Decode(req)

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		id, err := c.app.Commands.CreateContentDefinition.Handle(r.Context(), command.CreateContentDefinition{
			Name:        req.Name,
			Description: req.Description,
		})

		if err != nil {
			api.WithError(r.Context(), err)
			return
		}
		url := r.URL.String()
		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, id.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

// swagger:route GET /contentdefinition/{id} contentdefinition GetContentDefinition
//
// Gets a content definition
//
//     Consumes:
//	   - application/json
//
//     Responses:
//       200: GetContentDefinitionResponse
//		 400: genericError
//		 500: genericError
func (c contentEndpoint) GetContentDefinition() http.HandlerFunc {

	// swagger:parameters request GetContentDefinition
	type _ struct {
		ID uuid.UUID
	}

	// swagger:response GetContentDefinitionResponse
	type response struct {
		contentdefinition.ContentDefinition
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var id uuid.UUID

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			id = uid
		}

		cd, err := c.app.Queries.GetContentDefinition.Handle(r.Context(), query.GetContentDefinition{ID: id})

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		res := &response{cd}

		bytes, err := json.Marshal(res)

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		w.Write(bytes)
	}
}
