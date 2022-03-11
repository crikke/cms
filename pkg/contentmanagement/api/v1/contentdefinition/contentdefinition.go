package contentdefinition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api/models"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

var contentKey = key("cid")
var propertyKey = key("pid")

type endpoint struct {
	app app.App
}

func NewContentDefinitionEndpoint(app app.App) endpoint {
	return endpoint{app}
}

func (c endpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/contentdefinitions", func(r chi.Router) {

		// ! TODO ContentDefinitionID should use Context instead to remove duplicate code

		r.Get("/", c.ListContentDefinitions())
		r.Post("/", c.CreateContentDefinition())

		r.Route("/{id}", func(r chi.Router) {
			r.Use(func(h http.Handler) http.Handler {
				return contentDefinitionIdContext(h, "id", contentKey)
			})
			r.Get("/", c.GetContentDefinition())
			r.Delete("/", c.DeleteContentDefinition())
			r.Put("/", c.UpdateContentDefinition())

			r.Route("/propertydefinitions", func(r chi.Router) {

				r.Post("/", c.CreatePropertyDefinition())

				r.Route("/{pid}", func(r chi.Router) {
					r.Use(func(h http.Handler) http.Handler {
						return contentDefinitionIdContext(h, "pid", propertyKey)
					})
					r.Get("/", c.GetPropertyDefinition())
					r.Put("/", c.UpdatePropertyDefinition())
					r.Delete("/", c.DeletePropertyDefinition())

				})
			})
		})

	})
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
// @Param						body		body	ContentDefinitionBody	true 	"request body"
// @Success						201			{object}	models.OKResult
// @Header						201			{string}	Location
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/contentdefinitions [post]
func (c endpoint) CreateContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := &ContentDefinitionBody{}

		err := json.NewDecoder(r.Body).Decode(req)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		id, err := c.app.Commands.CreateContentDefinition.Handle(r.Context(), command.CreateContentDefinition{
			Name:        req.Name,
			Description: req.Description,
		})

		if err != nil {
			models.WithError(r.Context(), err)
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
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Success						200			{object}	contentdefinition.ContentDefinition
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/contentdefinitions/{id} [get]
func (c endpoint) GetContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())

		cd, err := c.app.Queries.GetContentDefinition.Handle(r.Context(), query.GetContentDefinition{ID: id})

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		bytes, err := json.Marshal(&cd)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
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
// @Success						200			{object}	query.ListContentDefinitionModel
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/contentdefinitions [get]
func (c endpoint) ListContentDefinitions() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = c.app.Queries.ListContentDefinitions.Handle(r.Context(), query.ListContentDefinition{})
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
// @Success						200			{object}	models.OKResult
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/contentdefinitions/{id} [delete]
func (c endpoint) DeleteContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())

		err := c.app.Commands.DeleteContentDefinition.Handle(r.Context(), command.DeleteContentDefinition{ID: id})
		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
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
// @Param						body		body	ContentDefinitionBody	true 	"request body"
// @Success						200			{object}	models.OKResult
// @Failure						default		{object}	models.GenericError
// @Router						/contentmanagement/contentdefinitions/{id} [put]
func (c endpoint) UpdateContentDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())

		body := &ContentDefinitionBody{}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		c.app.Commands.UpdateContentDefinition.Handle(r.Context(), command.UpdateContentDefinition{
			ContentDefinitionID: id,
			Name:                body.Name,
			Description:         body.Description,
		})
	}
}

// CreatePropertyDefinition 	godoc
// @Summary 					Creates a new propertydefinition
// @Description 				Creates a new propertydefinition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						body		body	CreatePropertyDefinitionBody	true 	"request body"
// @Failure						default		{object}	models.GenericError
// @Success						201			{object}	models.OKResult
// @Header						201			{string}	Location
// @Router						/contentmanagement/contentdefinitions/{id}/propertydefinitions [post]
func (c endpoint) CreatePropertyDefinition() http.HandlerFunc {
	// ! TODO Type should not be string, probably enum

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		body := &CreatePropertyDefinitionBody{}
		json.NewDecoder(r.Body).Decode(body)

		pid, err := c.app.Commands.CreatePropertyDefinition.Handle(r.Context(), command.CreatePropertyDefinition{
			ContentDefinitionID: id,
			Name:                body.Name,
			Description:         body.Description,
			Type:                body.Type,
		})

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		url := r.URL.String()
		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, pid.String()))
		w.WriteHeader(http.StatusCreated)

	}
}

// UpdatePropertyDefinition 	godoc
// @Summary 					Updates an property definition
// @Description 				Updates an property definition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						pid			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						body		body	UpdatePropertyDefinitionBody	true 	"request body"
// @Failure						default		{object}	models.GenericError
// @Success						200			{object}	models.OKResult
// @Router						/contentmanagement/contentdefinitions/{id}/propertydefinitions/{pid} [put]
func (c endpoint) UpdatePropertyDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		body := &UpdatePropertyDefinitionBody{}

		id := withID(r.Context())
		pid := withPID(r.Context())

		err := json.NewDecoder(r.Body).Decode(body)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		cmd := command.UpdatePropertyDefinition{
			ContentDefinitionID:  id,
			PropertyDefinitionID: pid,
			Name:                 &body.Name,
			Description:          &body.Description,
			Localized:            &body.Localized,
			Rules:                body.Validation,
		}
		err = c.app.Commands.UpdatePropertyDefinition.Handle(r.Context(), cmd)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}
	}
}

// DeletePropertyDefinition 	godoc
// @Summary 					Deletes a propertydefinition
// @Description 				Deletes a propertydefinition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						pid			path	string	true 	"uuid formatted ID." format(uuid)
// @Failure						default		{object}	models.GenericError
// @Success						200			{object}	models.OKResult
// @Router						/contentmanagement/contentdefinitions/{id}/propertydefinitions/{pid} [delete]
func (c endpoint) DeletePropertyDefinition() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		pid := withID(r.Context())

		err := c.app.Commands.DeletePropertyDefinition.Handle(r.Context(), command.DeletePropertyDefinition{
			ContentDefinitionID:  id,
			PropertyDefinitionID: pid,
		})

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

	}
}

// GetPropertyDefinition 	godoc
// @Summary 					Gets a propertydefinition
// @Description 				Gets a propertydefinition
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param						pid			path	string	true 	"uuid formatted ID." format(uuid)
// @Failure						default		{object}	models.GenericError
// @Success						200			{object}	contentdefinition.PropertyDefinition
// @Router						/contentmanagement/contentdefinitions/{id}/propertydefinitions/{pid} [get]
func (c endpoint) GetPropertyDefinition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := withID(r.Context())
		pid := withID(r.Context())

		pd, err := c.app.Queries.GetPropertyDefinition.Handle(r.Context(), query.GetPropertyDefinition{
			ContentDefinitionID:  id,
			PropertyDefinitionID: pid,
		})

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		bytes, err := json.Marshal(&pd)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				Body:       models.ErrorBody{Message: err.Error()},
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		w.Write(bytes)
	}
}
