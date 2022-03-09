package main

import (
	"context"
	"net/http"

	_ "github.com/crikke/cms/cmd/contentmanagement/docs"
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	// Configuration config.SiteConfiguration
	Database     *mongo.Client
	Eventhandler siteconfiguration.ConfigurationEventHandler
	SiteConfig   *siteconfiguration.SiteConfiguration
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	serverConfig := config.LoadServerConfiguration()

	c, err := db.Connect(context.Background(), serverConfig.ConnectionString.Mongodb)

	if err != nil {
		panic(err)
	}

	configRepo := siteconfiguration.NewConfigurationRepository(c)
	siteConfig, err := configRepo.LoadConfiguration(context.Background())

	server := Server{
		Database:   c,
		SiteConfig: siteConfig,
	}

	if err != nil {
		panic(err)
	}
	if serverConfig.ConnectionString.RabbitMQ != "" {
		eventhandler, err := siteconfiguration.NewConfigurationEventHandler(serverConfig.ConnectionString.RabbitMQ)

		if err != nil {
			panic(err)
		}

		defer func() {
			err = eventhandler.Close()
			if err != nil {
				panic(err)
			}
		}()
		eventhandler.Watch(siteConfig)
		server.Eventhandler = eventhandler
	}

	panic(server.Start())
}

func (s Server) Start() error {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/swagger", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Mount("/contentmanagement", api.NewContentManagementAPI(s.Database, s.Eventhandler, s.SiteConfig))

	return http.ListenAndServe(":8080", r)
}
