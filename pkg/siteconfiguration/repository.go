package siteconfiguration

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "configuration"

type ConfigurationRepository interface {
	LoadConfiguration(ctx context.Context) (*SiteConfiguration, error)
	UpdateConfiguration(ctx context.Context, cfg SiteConfiguration) error
}

type repository struct {
	db *mongo.Database
}

func NewConfigurationRepository(c *mongo.Client) ConfigurationRepository {

	db := c.Database("cms")
	repo := repository{db}
	return repo
}

func (r repository) LoadConfiguration(ctx context.Context) (*SiteConfiguration, error) {
	cfg := &SiteConfiguration{}
	err := r.db.
		Collection(collectionName).
		FindOne(ctx, bson.D{}).
		Decode(cfg)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return cfg, nil
}

func (r repository) UpdateConfiguration(ctx context.Context, cfg SiteConfiguration) error {

	_, err := r.db.Collection(collectionName).
		UpdateByID(
			ctx,
			cfg.ID,
			bson.M{"$set": cfg})

	return err
}
