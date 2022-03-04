package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/crikke/cms/pkg/contentdelivery/config"
	contentapi "github.com/crikke/cms/pkg/contentmanagement/api/v1/content"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	if serverConfig.ConnectionString.RabbitMQ != "" {
		closer, err := siteconfiguration.NewConfigurationWatcher(serverConfig.ConnectionString.RabbitMQ, siteConfig)

		if err != nil {
			panic(err)
		}

		defer func() {
			err = closer.Close()
			if err != nil {
				panic(err)
			}
		}()

	}

	server := Server{
		database:   c,
		SiteConfig: siteConfig,
	}

	panic(server.Start())
}

func (s Server) Start() error {

	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")

	if err != nil {
		panic(err)
	}
	contentRepo := content.NewContentRepository(c)
	app := app.App{
		Queries: app.Queries{
			GetContent: query.GetContentHandler{
				Repo: contentRepo,
			},
		},
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ep := contentapi.NewContentEndpoint(app)

	ep.RegisterEndpoints(r)

	r.Get("/swagger", func(rw http.ResponseWriter, r *http.Request) {
		dat, err := ioutil.ReadFile("./swagger.json")
		if err != nil {
			rw.Write([]byte(err.Error()))
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Write(dat)
	})
	return http.ListenAndServe(":8080", r)
}
