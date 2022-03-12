package workspace

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
)

type Workspace struct {
	ID          uuid.UUID         `bson:"_id"`
	Name        string            `bson:"name"`
	Description string            `bson:"description"`
	Languages   []string          `bson:"languages"`
	Tags        map[string]string `bson:"tags"`
}

func NewWorkspace(name, description, defaultLocale string) (Workspace, error) {

	if name == "" {
		return Workspace{}, errors.New("missing: name")
	}

	_, err := language.Parse(defaultLocale)

	if err != nil {
		return Workspace{}, err
	}

	ws := Workspace{
		Name:        name,
		Description: description,
		Languages:   []string{defaultLocale},
	}

	return ws, nil
}

const workspaceCollection = "workspace"

type WorkspaceRepository struct {
	client *mongo.Client
}

func NewWorkspaceRepository(client *mongo.Client) WorkspaceRepository {

	wr := WorkspaceRepository{
		client: client,
	}

	return wr
}

func (r WorkspaceRepository) Create(ctx context.Context, ws Workspace) (uuid.UUID, error) {

	if ws.ID == (uuid.UUID{}) {
		ws.ID = uuid.New()
	}

	_, err := r.client.Database("cms").
		Collection(workspaceCollection).
		InsertOne(ctx, ws)

	if err != nil {
		return uuid.UUID{}, err
	}

	return ws.ID, nil
}

func (r WorkspaceRepository) Get(ctx context.Context, id uuid.UUID) (Workspace, error) {

	ws := &Workspace{}

	err := r.client.Database("cms").
		Collection(workspaceCollection).
		FindOne(ctx, bson.M{"_id": id}).
		Decode(ws)

	if err != nil {
		return Workspace{}, err
	}

	return *ws, nil
}

func (r WorkspaceRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(ctx context.Context, ws *Workspace) (*Workspace, error)) error {
	ws := &Workspace{}

	err := r.client.Database("cms").
		Collection(workspaceCollection).
		FindOne(ctx, bson.M{"_id": id}).
		Decode(ws)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, ws)
	if err != nil {
		return err
	}

	_, err = r.client.Database("cms").
		Collection(workspaceCollection).
		UpdateOne(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": updated})

	return err
}

func (r WorkspaceRepository) ListAll(ctx context.Context) ([]Workspace, error) {
	cursor, err := r.client.Database("cms").
		Collection(workspaceCollection).
		Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	items := make([]Workspace, 0)
	for cursor.Next(ctx) {
		ws := &Workspace{}

		err = cursor.Decode(&ws)

		if err != nil {
			return nil, err
		}
		items = append(items, *ws)
	}

	return items, nil
}

func (r WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.client.Database("cms").
		Collection(workspaceCollection).
		DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}
	return r.client.Database(id.String()).Drop(ctx)
}
