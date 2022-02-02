package main

import (
	"context"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/loader"
	"github.com/crikke/cms/pkg/server"
	"golang.org/x/text/language"
)

func main() {
	cfg := config.LoadServerConfiguration()
	siteCfg := config.SiteConfiguration{
		Languages: []language.Tag{
			language.Swedish,
		},
	}

	db, err := loader.NewRepository(context.Background(), cfg)
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
