package content

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crikke/cms/pkg/contentmanagement/app"
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
		})
	})
}

// swagger:route GET /content pets users listPets
//
// Get content by id and optionally version
//
// Gets content by Id.
//
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: contentResponse
//       404
func (c contentEndpoint) GetContent() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var uid uuid.UUID
		var err error
		var ver *int
		if id := chi.URLParam(r, "id"); id != "" {
			uid, err = uuid.Parse(id)

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("bad formatted id"))
				return
			}
		}

		if v := r.URL.Query().Get("v"); v != "" {

			i, err := strconv.Atoi(v)

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("bad formatted version"))
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
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		response := &ContentResponse{Body: res}
		data, err := json.Marshal(response)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write(data)
	}
}

// ContentResponse is the representation of the content for the Content management API
// It contains all information of given content for every configured language.
//
// swagger:response contentResponse
type ContentResponse struct {
	// in: body
	Body query.ContentReadModel
}
