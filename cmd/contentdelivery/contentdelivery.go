package main

import (
	"context"

	"github.com/crikke/cms/pkg/api"
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/config/siteconfiguration"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/locale"
	"github.com/crikke/cms/pkg/repository"
	"github.com/crikke/cms/pkg/services/loader"
	"github.com/gin-gonic/gin"
)

type Server struct {
	// Configuration config.SiteConfiguration
	Loader     loader.Loader
	SiteConfig *domain.SiteConfiguration
}

func main() {

	serverConfig := config.LoadServerConfiguration()

	db, err := repository.NewRepository(context.Background(), serverConfig)

	if err != nil {
		panic(err)
	}

	siteConfig, err := db.LoadSiteConfiguration(context.Background())

	closer, err := siteconfiguration.NewConfigurationWatcher(serverConfig.ConnectionString.Mongodb, siteConfig)

	if err != nil {
		panic(err)
	}

	defer func() {
		err = closer.Close()
		if err != nil {
			panic(err)
		}
	}()

	l := loader.NewLoader(db, siteConfig)
	if err != nil {
		panic(err)
	}

	server := Server{
		Loader:     l,
		SiteConfig: siteConfig,
	}

	panic(server.Start())
}

func (s Server) Start() error {

	r := gin.Default()
	r.Use(locale.Handler(s.SiteConfig))

	v1 := r.Group("/v1")
	{
		api.ContentHandler(v1, s.SiteConfig, s.Loader)
	}

	return r.Run()
}
