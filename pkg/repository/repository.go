package repository

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language"
)

type Repository interface {
	GetContent(ctx context.Context, contentReference domain.ContentReference) (ContentData, error)
}

type ContentData struct {
	ID               uuid.UUID `bson:"_id"`
	ParentID         uuid.UUID
	PublishedVersion int
	Created          time.Time
	Updated          time.Time
	Data             map[int]ContentVersion
}
type ContentVersion struct {
	Properties []ContentProperty
	Name       map[string]string
	URLSegment map[string]string
}

type ContentProperty struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	Value     map[string]interface{}
}

type repository struct {
	connectionString string
	client           *mongo.Client
}

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	tTag        = reflect.TypeOf(language.Tag{})
	uuidSubtype = byte(0x04)
	registry    = bson.NewRegistryBuilder().
			RegisterTypeDecoder(tTag, bsoncodec.ValueDecoderFunc(decodeLanguageTag)).
			RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(encodeUUID)).
			RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(decodeUUID)).
			RegisterTypeEncoder(tTag, bsoncodec.ValueEncoderFunc(encodeLanguageTag)).
			Build()
)

func NewRepository(ctx context.Context, cfg config.ServerConfiguration) (Repository, error) {
	r := repository{connectionString: cfg.ConnectionString.Mongodb}
	c, err := mongo.Connect(ctx, options.Client().SetRegistry(registry), options.Client().ApplyURI(cfg.ConnectionString.Mongodb))

	// defer func() {
	// 	if err = c.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	if err != nil {
		return nil, err
	}

	r.client = c

	return r, nil
}

func (r repository) GetContent(ctx context.Context, contentReference domain.ContentReference) (ContentData, error) {

	doc := &ContentData{}
	err := r.client.Database("cms").Collection("content").FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: contentReference.ID}}).Decode(doc)

	if err != nil {
		return ContentData{}, err
	}
	return *doc, nil
}

func decodeUUID(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "decodeUUID", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error

	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}

	uuid, err := uuid.FromBytes(data)

	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(uuid))
	return nil
}

func encodeUUID(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "encodeUUID", Types: []reflect.Type{tUUID}, Received: val}
	}

	b := val.Interface().(uuid.UUID)

	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}
