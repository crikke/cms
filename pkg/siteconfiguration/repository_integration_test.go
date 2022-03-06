//go:build integration

package siteconfiguration

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
)

func Test_GetEmptyConfig(t *testing.T) {

	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	database := c.Database("cms")
	assert.NoError(t, err)
	r := NewConfigurationRepository(c)

	seed := func(input SiteConfiguration) {
		_, err := database.
			Collection("configuration").
			InsertOne(context.TODO(), input)

		assert.NoError(t, err)
	}

	tests := []struct {
		name   string
		input  *SiteConfiguration
		expect SiteConfiguration
	}{
		{
			name:   "Empty config",
			input:  nil,
			expect: SiteConfiguration{},
		},
		{
			name: "Existing config",
			input: &SiteConfiguration{
				Languages: []language.Tag{language.Swahili},
			},
			expect: SiteConfiguration{
				Languages: []language.Tag{language.Swahili},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err = truncateCollection("configuration", database)
			assert.NoError(t, err)

			if test.input != nil {
				seed(*test.input)
			}

			res, err := r.LoadConfiguration(context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, test.expect, *res)
		})
	}
}

func truncateCollection(col string, db *mongo.Database) error {

	return db.Collection(col).Drop(context.TODO())
}
