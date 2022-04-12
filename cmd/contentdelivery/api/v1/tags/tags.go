package tags

import (
	"net/http"

	"github.com/crikke/cms/cmd/contentdelivery/app"
	"github.com/go-chi/chi/v5"
)

func NewTagsRoute(app app.App) http.Handler {
	r := chi.NewRouter()

	return r
}
