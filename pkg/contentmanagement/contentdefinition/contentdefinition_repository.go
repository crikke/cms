package contentdefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "contentdefinition"

type ContentDefinitionRepository interface {
	CreateContentDefinition(ctx context.Context, cd *ContentDefinition) (uuid.UUID, error)
	UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error
	DeleteContentDefinition(ctx context.Context, id uuid.UUID) error
	GetContentDefinition(ctx context.Context, id uuid.UUID) (ContentDefinition, error)

	GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID) (PropertyDefinition, error)
}

type repository struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewContentDefinitionRepository(client *mongo.Client) ContentDefinitionRepository {

	db := client.Database("cms")
	r := repository{client: client, database: db}

	return r
}

func (r repository) CreateContentDefinition(ctx context.Context, cd *ContentDefinition) (uuid.UUID, error) {
	cd.ID = uuid.New()
	_, err := r.database.Collection(collection).InsertOne(ctx, cd)

	if err != nil {
		return uuid.UUID{}, err
	}

	return cd.ID, err
}

func (r repository) UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error {

	// TODO: Write/Read concern
	entry := &ContentDefinition{}
	err := r.database.Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(entry)
	if err != nil {
		return err
	}
	e, err := updateFn(ctx, entry)
	if err != nil {
		return err
	}

	_, err = r.database.
		Collection("contentdefinition").
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "_id", Value: id}},
			bson.M{"$set": e})

	if err != nil {
		return err
	}

	return nil
}

func (r repository) DeleteContentDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r repository) GetContentDefinition(ctx context.Context, id uuid.UUID) (ContentDefinition, error) {

	res := &ContentDefinition{}
	err := r.database.Collection(collection).FindOne(ctx, bson.D{bson.E{Key: "_id", Value: id}}).Decode(res)
	if err != nil {
		return ContentDefinition{}, err
	}
	return *res, nil
}

func (r repository) CreatePropertyDefinition(ctx context.Context, cid uuid.UUID, pd *PropertyDefinition) (uuid.UUID, error) {
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

func (r repository) UpdatePropertyDefinition(ctx context.Context, cid, pid uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error {

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

func (r repository) DeletePropertyDefinition(ctx context.Context, cid, pid uuid.UUID) error {
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

func (r repository) GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID) (PropertyDefinition, error) {
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
