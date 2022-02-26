package content

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "content"

type ContentRepository interface {
	CreateContent(ctx context.Context, content Content) (uuid.UUID, error)
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, id uuid.UUID, updateFn func(context.Context, *Content) (*Content, error)) error
}

type contentrepository struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewContentRepository(c *mongo.Client) ContentRepository {
	return contentrepository{
		client:   c,
		database: c.Database("cms"),
	}
}

func (c contentrepository) CreateContent(ctx context.Context, content Content) (uuid.UUID, error) {

	content.ID = uuid.New()
	_, err := c.database.
		Collection(collection).
		InsertOne(ctx, content)

	if err != nil {
		return uuid.UUID{}, err
	}
	return content.ID, nil
}

func (c contentrepository) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {

	content := &Content{}
	err := c.database.Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(content)

	if err != nil {
		return Content{}, err
	}

	return *content, nil
}

func (c contentrepository) UpdateContent(ctx context.Context, id uuid.UUID, updateFn func(context.Context, *Content) (*Content, error)) error {

	content := &Content{}
	err := c.database.Collection(collection).FindOne(ctx, bson.M{"_id": id}).Decode(content)

	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, content)

	if err != nil {
		return err
	}

	_, err = c.database.
		Collection(collection).
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "_id", Value: id}},
			bson.M{"$set": updated})

	if err != nil {
		return err
	}

	return nil
}