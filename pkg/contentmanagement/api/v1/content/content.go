package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type contentEndpoint struct {
	app app.App
}

func NewContentEndpoint(app app.App) contentEndpoint {
	return contentEndpoint{app}
}
func (c contentEndpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/content", func(r chi.Router) {

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", c.GetContent())
			r.Post("/", c.CreateContent())
		})

	})
}

type ErrorBody struct {
	// required: true
	Message   string
	FieldName string
}

// GenericError
// swagger:response genericError
type GenericError struct {
	// in: body
	Body ErrorBody
	// swagger:ignore
	StatusCode int
}

func (g *GenericError) WriteResponse(rw http.ResponseWriter) {
	b, err := json.Marshal(g)

	if err != nil {
		panic(err)
	}

	rw.Write(b)
	rw.WriteHeader(401)
}

// GetContentResponse is the representation of the content for the Content management API
// It contains all information of given content for every configured language.
//
// swagger:response contentResponse
type GetContentResponse struct {
	// in: body
	Body query.ContentReadModel
}

// swagger:parameters GetContentRequest GetContent
type GetContentRequest struct {
	// ID
	//
	// in: path
	// required:true
	ID string
	// Version
	//
	// in: query
	// required:false
	Version string
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

		req := GetContentRequest{
			ID:      chi.URLParam(r, "id"),
			Version: r.URL.Query().Get("version"),
		}
		var uid uuid.UUID
		var err error
		var ver *int
		if req.ID == "" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("missing id"))
			return
		}
		uid, err = uuid.Parse(req.ID)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			e := &GenericError{
				Body: ErrorBody{
					Message:   "bad formatted id",
					FieldName: "id",
				},
				StatusCode: http.StatusBadRequest,
			}

			e.WriteResponse(rw)
			return
		}

		if req.Version != "" {

			i, err := strconv.Atoi(req.Version)

			if err != nil {
				e := &GenericError{
					Body: ErrorBody{
						Message:   "bad formatted version",
						FieldName: "version",
					},
					StatusCode: http.StatusBadRequest,
				}
				e.WriteResponse(rw)
				return
			}

			ver = &i
		}

		q := query.GetContent{
			Id:      uid,
			Version: ver,
		}

		res, err := c.app.Queries.GetContent.Handle(r.Context(), q)

		if err != nil {

			e := &GenericError{
				Body: ErrorBody{
					Message: err.Error(),
				},
				StatusCode: 404,
			}

			e.WriteResponse(rw)
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
	ContentDefinitionId string `json:"contentdefinitionid"`
	// ParentId
	// in: body
	ParentId string `json:"parentid"`
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
			e := &GenericError{
				Body: ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			}
			e.WriteResponse(rw)
			return
		}

		cid, err := uuid.Parse(req.ContentDefinitionId)
		if err != nil {
			e := &GenericError{
				Body: ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			}
			e.WriteResponse(rw)
			return
		}

		pid, err := uuid.Parse(req.ParentId)
		if err != nil {
			e := &GenericError{
				Body: ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			}
			e.WriteResponse(rw)
			return

		}

		id, err := c.app.Commands.CreateContent.Handle(r.Context(),
			command.CreateContent{
				ContentDefinitionId: cid,
				ParentID:            pid,
			},
		)
		if err != nil {
			e := &GenericError{
				Body: ErrorBody{
					Message: err.Error(),
				},
				StatusCode: http.StatusBadRequest,
			}
			e.WriteResponse(rw)
			return

		}

		url := r.URL.String()
		rw.Header().Add("Location", fmt.Sprintf("%s/%s", url, id.String()))
		rw.WriteHeader(http.StatusCreated)
	}
}
