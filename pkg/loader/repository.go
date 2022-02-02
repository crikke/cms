package loader

import (
	"context"
	"errors"
	"time"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	GetContent(ctx context.Context, contentReference domain.ContentReference) (contentData, error)
}

type contentData struct {
	ID       uuid.UUID
	ParentID uuid.UUID
	Version  int
	Created  time.Time
	Updated  time.Time
	Data     map[int]contentVersion
}
type contentVersion struct {
	Properties []contentProperty
	Name       map[language.Tag]string
	URLSegment map[language.Tag]string
}

type contentProperty struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	Value     map[language.Tag]interface{}
}

type repository struct {
	connectionString string
	client           *mongo.Client
}

func NewRepository(ctx context.Context, cfg config.ServerConfiguration) (Repository, error) {
	r := repository{connectionString: cfg.ConnectionString.Mongodb}
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionString.Mongodb))

	defer func() {
		if err = c.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	r.client = c

	return r, nil
}

func (r repository) GetContent(ctx context.Context, contentReference domain.ContentReference) (contentData, error) {

	var doc contentData
	err := r.client.Database("").Collection("content").FindOne(ctx, bson.D{{"_id", contentReference.ID}}).Decode(doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return contentData{}, errors.New("not found")
		}
		panic(err)
	}
	return doc, nil
}
