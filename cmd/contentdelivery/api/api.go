package api

import (
	"net/http"

	"github.com/crikke/cms/cmd/contentdelivery/app"
	"github.com/go-chi/chi/v5"
)

func NewContentDeliveryAPI(app app.App) http.Handler {

	r := chi.NewRouter()

	// r.Mount("/v1/content", content.NewContentRoute(app))
	return r
}
