package contentdefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyDefinitionRepository interface {
	CreatePropertyDefinition(ctx context.Context, cid uuid.UUID, pd *PropertyDefinition) (uuid.UUID, error)
	UpdatePropertyDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error
	DeletePropertyDefinition(ctx context.Context, id uuid.UUID) error
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

func (r propertydefinition_repository) UpdatePropertyDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error {

	entry := &PropertyDefinition{}
	err := r.database.Collection(collection).FindOne(ctx, primitive.M{"_id": id}).Decode(entry)
	if err != nil {
		return err
	}
	e, err := updateFn(ctx, entry)
	if err != nil {
		return err
	}

	_, err = r.database.
		Collection(collection).
		UpdateOne(
			ctx,
			primitive.D{primitive.E{Key: "_id", Value: id}},
			bson.M{"$set": e})

	if err != nil {
		return err
	}

	// return session.CommitTransaction(sc)

	// })

	// if err != nil {
	// 	if abortErr := session.AbortTransaction(ctx); abortErr != nil {
	// 		return err
	// 	}
	// 	return err
	// }

	return nil
}

func (r propertydefinition_repository) DeletePropertyDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r propertydefinition_repository) GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID) (PropertyDefinition, error) {

	res := &ContentDefinition{}
	err := r.database.
		Collection(collection).
		FindOne(
			ctx,
			bson.D{
				bson.E{Key: "_id", Value: cid},
				bson.E{Key: "propertydefinitions.id", Value: pid}}).
		Decode(res)

	if err != nil {

		return PropertyDefinition{}, err
	}
	return res.PropertyDefinitions[0], nil
}
