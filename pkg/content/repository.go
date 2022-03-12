package content

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contentCollection = "content"
const contentVersionCollection = "contentversion"

type ContentManagementRepository struct {
	client *mongo.Client
}

func NewContentRepository(c *mongo.Client) ContentManagementRepository {
	return ContentManagementRepository{
		client: c,
	}
}

//! TODO: This should be done in a transaction
func (c ContentManagementRepository) CreateContent(ctx context.Context, content Content, workspace uuid.UUID) (uuid.UUID, error) {

	if content.ID == (uuid.UUID{}) {
		content.ID = uuid.New()
	}

	content.Data.ContentID = content.ID

	_, err := c.client.Database(workspace.String()).
		Collection(contentCollection).
		InsertOne(ctx, content)

	if err != nil {
		return uuid.UUID{}, err
	}

	_, err = c.client.Database(workspace.String()).
		Collection(contentVersionCollection).
		InsertOne(ctx, content.Data)

	if err != nil {
		return uuid.UUID{}, err
	}

	return content.ID, nil
}

func (c ContentManagementRepository) GetContent(ctx context.Context, id uuid.UUID, version int, workspace uuid.UUID) (Content, error) {

	content := &Content{}
	contentData := &ContentData{}
	filter := bson.M{
		"contentId": id,
		"version":   version,
	}

	err := c.client.Database(workspace.String()).Collection(contentVersionCollection).FindOne(ctx, filter).Decode(contentData)
	if err != nil {
		return Content{}, err
	}

	err = c.client.Database(workspace.String()).
		Collection(contentCollection).
		FindOne(
			ctx,
			bson.M{"_id": id},
			options.FindOne().SetProjection(bson.M{"data": 0})).
		Decode(content)

	if err != nil {
		return Content{}, err
	}

	content.Data = *contentData

	return *content, nil
}

func (c ContentManagementRepository) UpdateContentData(
	ctx context.Context,
	id uuid.UUID,
	version int,
	workspace uuid.UUID,
	updateFn func(context.Context, *ContentData) (*ContentData, error)) error {

	contentData := &ContentData{}
	err := c.client.Database(workspace.String()).Collection(contentVersionCollection).FindOne(ctx, bson.M{"contentId": id, "version": version}).Decode(contentData)

	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, contentData)

	if err != nil {
		return err
	}

	_, err = c.client.Database(workspace.String()).
		Collection(contentVersionCollection).
		UpdateOne(
			ctx,
			bson.M{"contentId": id, "version": updated.Version},
			bson.M{"$set": updated}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

func (c ContentManagementRepository) UpdateContent(
	ctx context.Context,
	id uuid.UUID,
	workspace uuid.UUID,
	updateFn func(context.Context, *Content) (*Content, error)) error {

	content := &Content{}
	err := c.client.Database(workspace.String()).Collection(contentCollection).FindOne(ctx, bson.M{"_id": id}).Decode(content)

	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, content)

	if err != nil {
		return err
	}

	_, err = c.client.Database(workspace.String()).
		Collection(contentCollection).
		UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": updated})

	if err != nil {
		return err
	}

	return nil
}

func (c ContentManagementRepository) ListContent(ctx context.Context, contentDefinitionTypes []uuid.UUID, tags []string, workspace uuid.UUID) ([]Content, error) {

	query := bson.M{}

	query["status"] = bson.M{"$ne": Archived}

	if len(contentDefinitionTypes) > 0 {
		query["contentdefinition_id"] = bson.M{
			"$in": bson.A{contentDefinitionTypes},
		}
	}

	if len(tags) > 0 {
		query["tags"] = bson.M{
			"$in": bson.A{tags},
		}
	}

	cur, err := c.client.Database(workspace.String()).
		Collection(contentCollection).
		Find(
			ctx,
			query)

	if err != nil {
		return nil, err
	}

	result := []Content{}
	for cur.Next(ctx) {
		data := &Content{}
		err = cur.Decode(data)

		if err != nil {
			return nil, err
		}

		result = append(result, *data)
	}

	return result, nil
}

func (c ContentManagementRepository) ListContentByTags(ctx context.Context, tags []string, workspace uuid.UUID) ([]Content, error) {

	c.client.Database(workspace.String()).Collection(contentCollection).Find(ctx, bson.M{})
	// for _, field := range tags {

	// }
	return nil, nil
}

func (c ContentManagementRepository) ListContentVersions(ctx context.Context, id uuid.UUID, workspace uuid.UUID) ([]ContentVersion, error) {
	filter := bson.M{"contentId": id}
	projection := bson.M{
		"_id":       0,
		"contentId": 1,
		"version":   1,
		"status":    1,
	}
	cursor, err := c.client.Database(workspace.String()).
		Collection(contentVersionCollection).
		Find(ctx, filter, options.Find().SetProjection(projection))

	if err != nil {
		return nil, err
	}

	items := make([]ContentVersion, 0)

	for cursor.Next(ctx) {
		item := &ContentVersion{}
		err := cursor.Decode(item)

		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}
	return items, nil
}
