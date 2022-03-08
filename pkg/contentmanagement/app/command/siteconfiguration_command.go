package command

import (
	"context"

	"github.com/crikke/cms/pkg/siteconfiguration"
	"golang.org/x/text/language"
)

type UpdateSiteConfiguration struct {
	Languages []string
}

type UpdateSiteConfigurationHandler struct {
	Cfg          *siteconfiguration.SiteConfiguration
	Repo         siteconfiguration.ConfigurationRepository
	Eventhandler siteconfiguration.ConfigurationEventHandler
}

func (h UpdateSiteConfigurationHandler) Handle(ctx context.Context, cmd UpdateSiteConfiguration) error {

	languages := []language.Tag{}

	for _, lang := range cmd.Languages {

		l, err := language.Parse(lang)

		if err != nil {
			return err
		}

		languages = append(languages, l)

	}
	cfg := siteconfiguration.SiteConfiguration{
		ID:        h.Cfg.ID,
		Languages: languages,
	}
	if err := h.Repo.UpdateConfiguration(ctx, cfg); err != nil {
		return err
	}

	// if there is no messagebroker configured
	// update siteconfig directly, otherwise let eventhandler be responsible for watching for changes and updating it
	if h.Eventhandler != nil {
		return h.Eventhandler.Publish(cfg)
	}

	*h.Cfg = cfg
	return nil
}
