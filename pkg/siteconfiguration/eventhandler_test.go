//go:build unit

package siteconfiguration

// func Test_MessageHandlerUpdatesConfiguration(t *testing.T) {

// 	ch := make(chan amqp.Delivery)
// 	done := make(chan error)
// 	cfg := &SiteConfiguration{
// 		Languages: []language.Tag{
// 			language.Swahili,
// 		},
// 	}

// 	newCfg := &SiteConfiguration{
// 		Languages: []language.Tag{
// 			language.Norwegian,
// 		},
// 	}

// 	data, err := json.Marshal(newCfg)

// 	assert.NoError(t, err)

// 	go messageHandler(cfg, ch, done)
// 	ch <- amqp.Delivery{
// 		Body: data,
// 	}
// 	close(ch)

// 	err = <-done
// 	assert.NoError(t, err)
// 	assert.Equal(t, language.Norwegian, cfg.Languages[0])
// }
