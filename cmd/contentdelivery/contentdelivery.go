package main

import (
	"context"

	"github.com/crikke/cms/pkg/config"
	contentapi "github.com/crikke/cms/pkg/contentdelivery/api/v1/content"
	"github.com/crikke/cms/pkg/contentdelivery/content"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/crikke/cms/pkg/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	// Configuration config.SiteConfiguration
	database   *mongo.Client
	SiteConfig *siteconfiguration.SiteConfiguration
}

func main() {

	serverConfig := config.LoadServerConfiguration()

	c, err := db.Connect(context.Background(), serverConfig.ConnectionString.Mongodb)

	if err != nil {
		panic(err)
	}

	configRepo := siteconfiguration.NewConfigurationRepository(c)
	siteConfig, err := configRepo.LoadConfiguration(context.Background())
	if err != nil {
		panic(err)
	}

	if serverConfig.ConnectionString.Mongodb != "" {
		cfgQueue, err := siteconfiguration.NewConfigurationEventHandler(serverConfig.ConnectionString.RabbitMQ)

		if err != nil {
			panic(err)
		}

		defer func() {
			err = cfgQueue.Close()
			if err != nil {
				panic(err)
			}
		}()
		cfgQueue.Watch(siteConfig)
	}

	server := Server{
		database:   c,
		SiteConfig: siteConfig,
	}

	panic(server.Start())
}

func (s Server) Start() error {

	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Use(
		telemetry.Handle(),
		gin.Recovery(),
		// locale.Handler(s.SiteConfig),
	)

	v1 := r.Group("/v1")
	{
		contentRepo := content.NewContentRepository(s.database, s.SiteConfig)
		contentapi.RegisterEndpoints(v1, s.SiteConfig, contentRepo)
	}

	return r.Run()
}
