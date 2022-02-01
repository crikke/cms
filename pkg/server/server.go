package server

import (
	"net/http"

	"github.com/crikke/cms/pkg/api"
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/loader"
	"github.com/crikke/cms/pkg/locale"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Configuration config.SiteConfiguration
	Loader        loader.Loader
}

func NewServer(cfg config.SiteConfiguration, loader loader.Loader) (Server, error) {
	return Server{cfg, loader}, nil
}

func (s Server) Start() http.Handler {

	r := gin.Default()
	r.Use(locale.Handler(s.Configuration))

	v1 := r.Group("/v1")
	{
		api.ContentHandler(v1, s.Configuration, s.Loader)
	}
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

	})
}
