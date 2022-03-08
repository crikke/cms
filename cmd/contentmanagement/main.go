package main

import (
	"context"
	"net/http"

	_ "github.com/crikke/cms/cmd/contentmanagement/docs"
	"github.com/crikke/cms/pkg/contentdelivery/config"
	contentapi "github.com/crikke/cms/pkg/contentmanagement/api/v1/content"
	contentdefapi "github.com/crikke/cms/pkg/contentmanagement/api/v1/contentdefinition"
	cfgapi "github.com/crikke/cms/pkg/contentmanagement/api/v1/siteconfiguration"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	// Configuration config.SiteConfiguration
	database     *mongo.Client
	eventhandler siteconfiguration.ConfigurationEventHandler
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
		database:   c,
		SiteConfig: siteConfig,
	}

	if err != nil {
		panic(err)
	}
	if serverConfig.ConnectionString.Mongodb != "" {
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
		server.eventhandler = eventhandler
	}

	panic(server.Start())
}

func (s Server) Start() error {

	contentDefinitionRepo := contentdefinition.NewContentDefinitionRepository(s.database)
	contentRepo := content.NewContentRepository(s.database)
	app := app.App{
		Queries: app.Queries{
			GetContent: query.GetContentHandler{
				Repo: contentRepo,
			},
			ListContent: query.ListContentHandler{
				Repo: contentRepo,
			},
			GetContentDefinition: query.GetContentDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
			GetPropertyDefinition: query.GetPropertyDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
			ListContentDefinitions: query.ListContentDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
		},
		Commands: app.Commands{
			CreateContent: command.CreateContentHandler{
				ContentDefinitionRepository: contentDefinitionRepo,
				ContentRepository:           contentRepo,
				Factory:                     content.Factory{Cfg: s.SiteConfig},
			},
			UpdateContentFields: command.UpdateContentFieldsHandler{
				ContentRepository:           contentRepo,
				ContentDefinitionRepository: contentDefinitionRepo,
				Factory:                     content.Factory{Cfg: s.SiteConfig},
			},
			ArchiveContent: command.ArchiveContentHandler{
				ContentRepository: contentRepo,
			},
			PublishContent: command.PublishContentHandler{
				ContentDefinitionRepository: contentDefinitionRepo,
				ContentRepository:           contentRepo,
				SiteConfiguration:           s.SiteConfig,
			},
			CreateContentDefinition: command.CreateContentDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
			UpdateContentDefinition: command.UpdateContentDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
			DeleteContentDefinition: command.DeleteContentDefinitionHandler{},
			CreatePropertyDefinition: command.CreatePropertyDefinitionHandler{
				Repo:    contentDefinitionRepo,
				Factory: contentdefinition.PropertyDefinitionFactory{},
			},
			UpdatePropertyDefinition: command.UpdatePropertyDefinitionHandler{
				Repo: contentDefinitionRepo,
			},
			DeletePropertyDefinition: command.DeletePropertyDefinitionHandler{},
			UpdateSiteConfiguration: command.UpdateSiteConfigurationHandler{
				Cfg:          s.SiteConfig,
				Repo:         siteconfiguration.NewConfigurationRepository(s.database),
				Eventhandler: s.eventhandler,
			},
		},
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/swagger", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	contentendpoint := contentapi.NewContentEndpoint(app)
	contentendpoint.RegisterEndpoints(r)

	contentdefendpoint := contentdefapi.NewContentDefinitionEndpoint(app)
	contentdefendpoint.RegisterEndpoints(r)

	r.Mount("/siteconfiguration", cfgapi.NewSiteConfigurationRouter(s.SiteConfig))

	return http.ListenAndServe(":8080", r)
}
