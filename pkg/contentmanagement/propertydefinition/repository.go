package propertydefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyDefinitionRepository interface {
	CreatePropertyDefinition(ctx context.Context, pd *PropertyDefinition) (uuid.UUID, error)
	UpdatePropertyDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error
	DeletePropertyDefinition(ctx context.Context, id uuid.UUID) error
	GetPropertyDefinition(ctx context.Context, id uuid.UUID) (PropertyDefinition, error)
}

type repository struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewPropertyDefinitionRepository(client *mongo.Client) PropertyDefinitionRepository {

	db := client.Database("cms")
	r := repository{client: client, database: db}

	return r
}

func (r repository) CreatePropertyDefinition(ctx context.Context, pd *PropertyDefinition) (uuid.UUID, error) {
	pd.ID = uuid.New()
	_, err := r.database.Collection("contentdefinition").InsertOne(ctx, pd)

	if err != nil {
		return uuid.UUID{}, err
	}

	return pd.ID, err
}

func (r repository) UpdatePropertyDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, pd *PropertyDefinition) (*PropertyDefinition, error)) error {

	// TODO: Write/Read concern
	// wc := writeconcern.WMajority()
	// rc := readconcern.Majority()
	// session, err := r.client.StartSession()
	// if err != nil {
	// 	return err
	// }
	// defer session.EndSession(ctx)

	// err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
	// if err != session.StartTransaction() {
	// 	return err
	// }

	entry := &PropertyDefinition{}
	err := r.database.Collection("contentdefinition").FindOne(ctx, primitive.M{"_id": id}).Decode(entry)
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

func (r repository) DeletePropertyDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r repository) GetPropertyDefinition(ctx context.Context, id uuid.UUID) (PropertyDefinition, error) {

	res := &PropertyDefinition{}
	err := r.database.Collection("contentdefinition").FindOne(ctx, bson.D{bson.E{Key: "_id", Value: id}}).Decode(res)
	if err != nil {
		return PropertyDefinition{}, err
	}
	return *res, nil
}
