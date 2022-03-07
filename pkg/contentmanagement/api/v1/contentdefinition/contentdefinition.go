package contentdefinition

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
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

	router.Route("/contentdefinitions", func(r chi.Router) {

		// ! TODO ContentDefinitionID should use Context instead to remove duplicate code

		r.Get("/", c.ListContentDefinitions())
		r.Post("/", c.CreateContentDefinition())

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", c.GetContentDefinition())
			r.Delete("/", c.DeleteContentDefinition())
			r.Put("/", c.UpdateContentDefinition())

			r.Route("/propertydefinitions", func(r chi.Router) {

				r.Post("/", c.CreatePropertyDefinition())

				r.Route("/{pid}", func(r chi.Router) {
					r.Get("/", c.GetPropertyDefinition())
					r.Put("/", c.UpdatePropertyDefinition())
					r.Delete("/", c.DeletePropertyDefinition())

					// r.Put("/validator", c.UpdatePropertyDefinitionValidator())
				})
			})
		})

	})
}

// swagger:route POST /contentdefinitions contentdefinition CreateContentDefinition
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
//       201: Location
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

// swagger:route GET /contentdefinitions/{id} contentdefinition GetContentDefinition
//
// Gets a content definition
//
// Gets a content definition by ID
//
//     Consumes:
//	   - application/json
//
//     Responses:
//       200: ContentDefinition
//		 400: genericError
//		 500: genericError
func (c contentEndpoint) GetContentDefinition() http.HandlerFunc {

	// swagger:parameters request GetContentDefinition
	type _ struct {
		ID uuid.UUID
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

		bytes, err := json.Marshal(&cd)

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

// swagger:route GET /contentdefinitions contentdefinition ListContentDefinitions
//
// Get all content definitions
//
// Gets all existing contentdefinitions
//
//     Produces:
//	   - application/json
//
//     Responses:
//       200: ListContentDefinitionsResponse
//		 400: genericError
//		 500: genericError
func (c contentEndpoint) ListContentDefinitions() http.HandlerFunc {

	// swagger:response ListContentDefinitionsResponse
	type _ struct {

		// in: body
		ContentDefinitions []struct {
			Name string
			ID   uuid.UUID
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// swagger:route DELETE /contentdefinitions/{id} contentdefinition DeleteContentDefinition
//
// Delete a content definition
//
//     Responses:
//       200: OK
//		 400: genericError
//		 404: genericError
//		 500: genericError
func (c contentEndpoint) DeleteContentDefinition() http.HandlerFunc {

	// swagger:parameters DeleteContentDefinition
	type _ struct {
		ID uuid.UUID
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

		err := c.app.Commands.DeleteContentDefinition.Handle(r.Context(), command.DeleteContentDefinition{ID: id})
		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}
}

// swagger:route PUT /contentdefinitions/{id} contentdefinition UpdateContentDefinition
//
// Updates a contentdefinition
//
//     Responses:
//       200: OK
//		 400: genericError
//		 404: genericError
//		 500: genericError
func (c contentEndpoint) UpdateContentDefinition() http.HandlerFunc {

	// swagger:parameters UpdateContentDefinition
	type request struct {
		// ID
		Id uuid.UUID
		// in:body
		Body struct {
			Name        string
			Description string
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.Id = uid
		}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		c.app.Commands.UpdateContentDefinition.Handle(r.Context(), command.UpdateContentDefinition{
			ContentDefinitionID: req.Id,
		})
	}
}

// swagger:route POST /contentdefinitions/{id}/propertydefinitions contentdefinition propertydefinition CreatePropertyDefinition
//
// Creates a new propertydefinition
//
//     Responses:
//       201: Location
//		 400: genericError
//		 404: genericError
//		 500: genericError
func (c contentEndpoint) CreatePropertyDefinition() http.HandlerFunc {
	// ! TODO Type should not be string, probably enum

	// swagger:parameters CreatePropertyDefinition
	type request struct {
		Id uuid.UUID
		// in: body
		Body *struct {
			Name        string
			Description string
			Type        string
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{Body: &struct {
			Name        string
			Description string
			Type        string
		}{}}

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.Id = uid
		}

		json.NewDecoder(r.Body).Decode(req.Body)

		pid, err := c.app.Commands.CreatePropertyDefinition.Handle(r.Context(), command.CreatePropertyDefinition{
			ContentDefinitionID: req.Id,
			Name:                req.Body.Name,
			Description:         req.Body.Description,
			Type:                req.Body.Type,
		})

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		url := r.URL.String()
		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, pid.String()))
		w.WriteHeader(http.StatusCreated)

	}
}

// swagger:route PUT /contentdefinitions/{id}/propertydefinitions/{pid} contentdefinition propertydefinition UpdatePropertyDefinition
//
// Updates an property definition
//
//     Responses:
//		 200: OK
//		 500: genericError
func (c contentEndpoint) UpdatePropertyDefinition() http.HandlerFunc {

	// swagger:parameters UpdatePropertyDefinition
	type request struct {
		// required:true
		ContentDefinitionID uuid.UUID
		// required:true
		PropertyDefinitionID uuid.UUID
		// in:body
		Body *struct {
			// required:true
			Name string
			// required:true
			Description string
			// required:true
			Localized bool
			// Validators
			Validation map[string]interface{}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{Body: &struct {
			Name        string
			Description string
			Localized   bool
			Validation  map[string]interface{}
		}{}}

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.ContentDefinitionID = uid
		}

		if param := chi.URLParam(r, "pid"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.PropertyDefinitionID = uid
		}

		err := json.NewDecoder(r.Body).Decode(req.Body)

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		cmd := command.UpdatePropertyDefinition{
			ContentDefinitionID:  req.ContentDefinitionID,
			PropertyDefinitionID: req.PropertyDefinitionID,
			Name:                 &req.Body.Name,
			Description:          &req.Body.Description,
			Localized:            &req.Body.Localized,
			Rules:                req.Body.Validation,
		}
		err = c.app.Commands.UpdatePropertyDefinition.Handle(r.Context(), cmd)

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}
}

// swagger:route DELETE /contentdefinitions/{id}/propertydefinitions/{pid} contentdefinition propertydefinition DeletePropertyDefinition
//
// Deletes a propertydefinition
//
// Deletes an property definition
//
//     Responses:
//		 200: OK
//		 500: genericError
func (c contentEndpoint) DeletePropertyDefinition() http.HandlerFunc {

	// swagger:parameters DeletePropertyDefinition
	type request struct {
		// required:true
		ContentDefinitionID uuid.UUID
		// required:true
		PropertyDefinitionID uuid.UUID
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.ContentDefinitionID = uid
		}

		if param := chi.URLParam(r, "pid"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.PropertyDefinitionID = uid
		}

		err := c.app.Commands.DeletePropertyDefinition.Handle(r.Context(), command.DeletePropertyDefinition{
			ContentDefinitionID:  req.ContentDefinitionID,
			PropertyDefinitionID: req.PropertyDefinitionID,
		})

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

	}
}

// swagger:route GET /contentdefinitions/{id}/propertydefinitions/{pid} contentdefinition propertydefinition GetPropertyDefinition
//
// Gets a propertydefinition
//
// Gets a propertydefinition
//
//     Responses:
//		 200: PropertyDefinition
//		 500: genericError
func (c contentEndpoint) GetPropertyDefinition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// swagger:parameters GetPropertyDefinition
		type request struct {
			// required:true
			ContentDefinitionID uuid.UUID
			// required:true
			PropertyDefinitionID uuid.UUID
		}

		req := &request{}

		if param := chi.URLParam(r, "id"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.ContentDefinitionID = uid
		}

		if param := chi.URLParam(r, "pid"); param != "" {

			uid, err := uuid.Parse(param)

			if err != nil {
				api.WithError(r.Context(), api.GenericError{
					Body:       api.ErrorBody{Message: err.Error()},
					StatusCode: http.StatusBadRequest,
				})
				return
			}
			req.PropertyDefinitionID = uid
		}

		pd, err := c.app.Queries.GetPropertyDefinition.Handle(r.Context(), query.GetPropertyDefinition{
			ContentDefinitionID:  req.ContentDefinitionID,
			PropertyDefinitionID: req.PropertyDefinitionID,
		})

		if err != nil {
			api.WithError(r.Context(), api.GenericError{
				Body:       api.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		bytes, err := json.Marshal(&pd)

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
