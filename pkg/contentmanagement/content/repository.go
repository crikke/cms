package content

import (
	"context"

	"github.com/google/uuid"
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

	return Content{}, nil
}

func (c contentrepository) UpdateContent(ctx context.Context, id uuid.UUID, updateFn func(context.Context, *Content) (*Content, error)) error {
	return nil
}
