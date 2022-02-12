package siteconfiguration

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigurationRepository interface {
	LoadConfiguration(ctx context.Context) (*SiteConfiguration, error)
}

type repository struct {
	db *mongo.Database
}

func NewConfigurationRepository(db *mongo.Database) ConfigurationRepository {

	repo := repository{db}
	return repo
}

func (r repository) LoadConfiguration(ctx context.Context) (*SiteConfiguration, error) {
	cfg := &SiteConfiguration{}
	err := r.db.
		Collection("configuration").
		FindOne(ctx, bson.D{}).
		Decode(cfg)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return cfg, nil
}
