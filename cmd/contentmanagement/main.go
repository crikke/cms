package main

import (
	"context"
	"net/http"

	"github.com/crikke/cms/cmd/contentmanagement/api"
	_ "github.com/crikke/cms/cmd/contentmanagement/docs"
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	// Configuration config.SiteConfiguration
	Database *mongo.Client
	Logger   *zap.SugaredLogger
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:8080
// @BasePath  /
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	serverConfig := config.LoadServerConfiguration()

	c, err := db.Connect(context.Background(), serverConfig.ConnectionString.Mongodb)

	if err != nil {
		panic(err)
	}

	server := Server{
		Database: c,
		Logger:   sugar,
	}

	panic(server.Start())
}

func (s Server) Start() error {

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return origin == "http://localhost:3000"
		},
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/swagger", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Mount("/contentmanagement", api.NewContentManagementAPI(s.Database))

	return http.ListenAndServe(":8080", r)
}
