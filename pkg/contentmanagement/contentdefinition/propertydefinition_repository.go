package contentdefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyDefinitionRepository interface {
	CreatePropertyDefinition(ctx context.Context, cid uuid.UUID, pd *PropertyDefinition) (uuid.UUID, error)
	UpdatePropertyDefinition(ctx context.Context, cid, pid uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error
	DeletePropertyDefinition(ctx context.Context, cid, pid uuid.UUID) error
	GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID) (PropertyDefinition, error)
}

type propertydefinition_repository struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewPropertyDefinitionRepository(client *mongo.Client) PropertyDefinitionRepository {

	db := client.Database("cms")
	r := propertydefinition_repository{client: client, database: db}

	return r
}

func (r propertydefinition_repository) CreatePropertyDefinition(ctx context.Context, cid uuid.UUID, pd *PropertyDefinition) (uuid.UUID, error) {
	pd.ID = uuid.New()

	_, err := r.database.
		Collection(collection).
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "_id", Value: cid}},
			bson.M{"$push": bson.M{"propertydefinitions": pd}})

	if err != nil {
		return uuid.UUID{}, err
	}

	return pd.ID, err
}

func (r propertydefinition_repository) UpdatePropertyDefinition(ctx context.Context, cid, pid uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error {

	entry, err := r.GetPropertyDefinition(ctx, cid, pid)
	if err != nil {
		return err
	}
	e, err := updateFn(ctx, &entry)
	if err != nil {
		return err
	}

	_, err = r.database.
		Collection(collection).
		UpdateOne(
			ctx,
			bson.D{
				bson.E{Key: "_id", Value: cid},
				bson.E{Key: "propertydefinitions.id", Value: pid}},
			bson.M{"$set": bson.M{"propertydefinitions.$": e}})

	if err != nil {
		return err
	}

	return nil
}

func (r propertydefinition_repository) DeletePropertyDefinition(ctx context.Context, cid, pid uuid.UUID) error {
	_, err := r.database.
		Collection(collection).
		DeleteOne(
			ctx,
			bson.D{
				bson.E{Key: "_id", Value: cid},
				bson.E{Key: "propertydefinitions.id", Value: pid}})

	if err != nil {
		return err
	}

	return nil
}

func (r propertydefinition_repository) GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID) (PropertyDefinition, error) {
	var res struct {
		PropertyDefinitions []PropertyDefinition `bson:"propertydefinitions,omitempty"`
	}

	err := r.database.
		Collection(collection).
		FindOne(
			ctx,
			bson.D{
				bson.E{Key: "_id", Value: cid},
				bson.E{Key: "propertydefinitions.id", Value: pid}}).
		Decode(&res)

	if err != nil {
		return PropertyDefinition{}, err
	}

	return res.PropertyDefinitions[0], nil
}
