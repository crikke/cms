package contentdefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContentDefinitionRepository interface {
	CreateContentDefinition(ctx context.Context, cd *ContentDefinition) (uuid.UUID, error)
	UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error
	DeleteContentDefinition(ctx context.Context, id uuid.UUID) error
	GetContentDefinition(ctx context.Context, id uuid.UUID) (ContentDefinition, error)
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
	_, err := r.database.Collection("contentdefinition").InsertOne(ctx, cd)

	if err != nil {
		return uuid.UUID{}, err
	}

	return cd.ID, err
}

func (r repository) UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error {

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

	entry := &ContentDefinition{}
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

func (r repository) DeleteContentDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r repository) GetContentDefinition(ctx context.Context, id uuid.UUID) (ContentDefinition, error) {

	res := &ContentDefinition{}
	err := r.database.Collection("contentdefinition").FindOne(ctx, bson.D{bson.E{Key: "_id", Value: id}}).Decode(res)
	if err != nil {
		return ContentDefinition{}, err
	}
	return *res, nil
}
