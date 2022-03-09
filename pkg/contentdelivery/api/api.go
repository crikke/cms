package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewContentDeliveryAPI() http.Handler {

	r := chi.NewRouter()

	return r
}
