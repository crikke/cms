package contentdefinition

import (
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/go-chi/chi"
)

type contentEndpoint struct {
	app app.App
}

func NewContentEndpoint(app app.App) contentEndpoint {
	return contentEndpoint{app}
}

func (c contentEndpoint) RegisterEndpoints(router chi.Router) {

	router.Route("/contentdefinition", func(r chi.Router) {

		r.Route("/{id}", func(r chi.Router) {
			// r.Get("/", c.GetContentDefinition())
			// r.Post("/", c.CreateContentDefinition())
		})

	})
}
