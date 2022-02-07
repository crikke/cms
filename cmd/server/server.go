package main

import (
	"context"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/repository"
	"github.com/crikke/cms/pkg/server"
	"github.com/crikke/cms/pkg/services/loader"
	"golang.org/x/text/language"
)

func main() {
	cfg := config.LoadServerConfiguration()
	siteCfg := config.SiteConfiguration{
		Languages: []language.Tag{
			language.Swedish,
		},
	}

	db, err := repository.NewRepository(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	loader := loader.NewLoader(db, siteCfg)
	s, err := server.NewServer(siteCfg, loader)

	if err != nil {
		panic(err)
	}

	panic(s.Start())
}
