package main

import (
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/server"
)

func main() {
	// cfg := config.LoadConfiguration()
	s, err := server.NewServer(config.SiteConfiguration{}, nil)

	if err != nil {
		panic(err)
	}

	panic(s.Start())
}
