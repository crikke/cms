package contentdefinition

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type ContentDefinitionRepository interface {
	CreateContentDefinition(ctx context.Context, cd *ContentDefinition) error
	UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error
	DeleteContentDefinition(ctx context.Context, id uuid.UUID) error
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

func (r repository) CreateContentDefinition(ctx context.Context, cd *ContentDefinition) error {
	_, err := r.database.Collection("contentdefinition").InsertOne(ctx, cd)
	return err
}

func (r repository) UpdateContentDefinition(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error {

	// TODO: Write/Read concern
	// wc := writeconcern.WMajority()
	// rc := readconcern.Majority()
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err != session.StartTransaction() {
			return err
		}

		entry := &ContentDefinition{}
		err = r.database.Collection("contentdefinition").FindOne(sc, primitive.E{Key: "_id", Value: id}).Decode(entry)
		if err != nil {
			return err
		}
		e, err := updateFn(sc, entry)

		if err != nil {
			return err
		}

		_, err = r.database.Collection("contentdefinition").UpdateOne(sc, primitive.E{Key: "_id", Value: id}, e)

		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)

	})

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return err
		}
		return err
	}

	return nil
}

func (r repository) DeleteContentDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}
