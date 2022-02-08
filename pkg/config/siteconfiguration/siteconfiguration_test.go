package siteconfiguration

import (
	"encoding/json"
	"testing"

	"github.com/crikke/cms/pkg/domain"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_MessageHandlerUpdatesConfiguration(t *testing.T) {

	ch := make(chan amqp.Delivery)
	done := make(chan error)
	cfg := &domain.SiteConfiguration{
		Languages: []language.Tag{
			language.Swahili,
		},
	}

	newCfg := &domain.SiteConfiguration{
		Languages: []language.Tag{
			language.Norwegian,
		},
	}

	data, err := json.Marshal(newCfg)

	assert.NoError(t, err)

	go messageHandler(cfg, ch, done)
	ch <- amqp.Delivery{
		Body: data,
	}
	close(ch)

	err = <-done
	assert.NoError(t, err)
	assert.Equal(t, language.Norwegian, cfg.Languages[0])
}
