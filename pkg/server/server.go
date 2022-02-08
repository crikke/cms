package server

import (
	"github.com/crikke/cms/pkg/api"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/locale"
	"github.com/crikke/cms/pkg/services/loader"
	"github.com/gin-gonic/gin"
)

type Server struct {
	// Configuration config.SiteConfiguration
	Loader loader.Loader
}

func NewServer(loader loader.Loader) (Server, error) {
	return Server{loader}, nil
}

func (s Server) Start() error {

	cfg := domain.SiteConfiguration{}
	r := gin.Default()
	r.Use(locale.Handler(cfg))

	v1 := r.Group("/v1")
	{
		api.ContentHandler(v1, cfg, s.Loader)
	}

	return r.Run()
}
