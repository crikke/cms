package workspace

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tag struct {
	Id   uuid.UUID `bson:"_id"`
	Name string    `bson:"name"`
}
type Workspace struct {
	ID        uuid.UUID `bson:"_id"`
	Name      string    `bson:"name"`
	Languages []string  `bson:"languages"`
	Tags      []Tag     `bson:"tags"`
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
