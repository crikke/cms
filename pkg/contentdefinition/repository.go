package contentdefinition

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const contentdefinitionCollection = "contentdefinition"

type ContentDefinitionRepository struct {
	client *mongo.Client
}

func NewContentDefinitionRepository(client *mongo.Client) ContentDefinitionRepository {

	r := &ContentDefinitionRepository{client: client}
	return *r
}

func (r ContentDefinitionRepository) CreateContentDefinition(ctx context.Context, cd *ContentDefinition, workspaceId uuid.UUID) (uuid.UUID, error) {
	cd.ID = uuid.New()

	_, err := r.client.Database(workspaceId.String()).Collection(contentdefinitionCollection).InsertOne(ctx, cd)

	if err != nil {
		return uuid.UUID{}, err
	}

	return cd.ID, err
}

func (r ContentDefinitionRepository) UpdateContentDefinition(ctx context.Context, id uuid.UUID, workspaceId uuid.UUID, updateFn func(ctx context.Context, cd *ContentDefinition) (*ContentDefinition, error)) error {

	entry := &ContentDefinition{}
	err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
		FindOne(ctx, bson.M{"_id": id}).Decode(entry)

	if err != nil {
		return err
	}
	e, err := updateFn(ctx, entry)
	if err != nil {
		return err
	}

	_, err = r.client.Database(workspaceId.String()).
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

func (r ContentDefinitionRepository) DeleteContentDefinition(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r ContentDefinitionRepository) GetContentDefinition(ctx context.Context, id uuid.UUID, workspaceId uuid.UUID) (ContentDefinition, error) {

	res := &ContentDefinition{}
	err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
		FindOne(ctx, bson.D{bson.E{Key: "_id", Value: id}}).
		Decode(res)

	if err != nil {
		return ContentDefinition{}, err
	}
	return *res, nil
}

func (r ContentDefinitionRepository) ListContentDefinitions(ctx context.Context, workspaceId uuid.UUID) ([]ContentDefinition, error) {
	cursor, err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
		Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	items := make([]ContentDefinition, 0)

	for cursor.Next(ctx) {

		res := &ContentDefinition{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		items = append(items, *res)
	}

	return items, nil
}

func (r ContentDefinitionRepository) CreatePropertyDefinition(ctx context.Context, cid uuid.UUID, workspaceId uuid.UUID, pd *PropertyDefinition) (uuid.UUID, error) {
	pd.ID = uuid.New()

	_, err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "_id", Value: cid}},
			bson.M{"$push": bson.M{"propertydefinitions": pd}})

	if err != nil {

		return uuid.UUID{}, err
	}

	return pd.ID, err
}

func (r ContentDefinitionRepository) DeletePropertyDefinition(ctx context.Context, cid, pid uuid.UUID, workspaceId uuid.UUID) error {
	_, err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
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

func (r ContentDefinitionRepository) GetPropertyDefinition(ctx context.Context, cid, pid uuid.UUID, workspaceId uuid.UUID) (PropertyDefinition, error) {
	var res struct {
		PropertyDefinitions []PropertyDefinition `bson:"propertydefinitions,omitempty"`
	}

	err := r.client.Database(workspaceId.String()).
		Collection(contentdefinitionCollection).
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
