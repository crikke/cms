package content

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/crikke/cms/pkg/contentmanagement"
	"github.com/crikke/cms/pkg/contentmanagement/content/query"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type contentEndpoint struct {
	app contentmanagement.App
}

func NewContentEndpoint(app contentmanagement.App) contentEndpoint {
	return contentEndpoint{app}
}
func (c contentEndpoint) RegisterEndpoints(router chi.Router) {
	router.Route("/content", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", c.GetContent())
		})
	})
}

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

		data, err := json.Marshal(&res)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write(data)
	}
}
