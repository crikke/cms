package repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"golang.org/x/text/language"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	// GetContent(ctx context.Context, contentReference domain.ContentReference) (ContentData, error)
	// GetChildren(ctx context.Context, contentReference domain.ContentReference) ([]ContentData, error)
	LoadSiteConfiguration(ctx context.Context) (*domain.SiteConfiguration, error)
}

// ContentData is not part of domain to make it clear that the repository is responsible this struct,

type repository struct {
	connectionString string
	client           *mongo.Client
	database         *mongo.Database
}

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	tTag        = reflect.TypeOf(language.Tag{})
	uuidSubtype = byte(0x04)
	registry    = bson.NewRegistryBuilder().
			RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(encodeUUID)).
			RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(decodeUUID)).
			RegisterTypeEncoder(tTag, bsoncodec.ValueEncoderFunc(encodeTag)).
			RegisterTypeDecoder(tTag, bsoncodec.ValueDecoderFunc(decodeTag)).
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
	r.database = c.Database("cms")

	return r, nil
}

func (r repository) LoadSiteConfiguration(ctx context.Context) (*domain.SiteConfiguration, error) {

	cfg := &domain.SiteConfiguration{}
	err := r.database.
		Collection("configuration").
		FindOne(ctx, bson.D{}).
		Decode(cfg)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return cfg, nil
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

func encodeTag(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	if !val.IsValid() || val.Type() != tTag {
		return bsoncodec.ValueEncoderError{Name: "encodeTag", Types: []reflect.Type{tTag}, Received: val}
	}

	b := val.Interface().(language.Tag)

	return vw.WriteString(b.String())
}

func decodeTag(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tTag {
		return bsoncodec.ValueDecoderError{Name: "decodeTag", Types: []reflect.Type{tTag}, Received: val}
	}

	var str string
	var err error

	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		str, err = vr.ReadString()
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

	tag, err := language.Parse(str)

	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(tag))
	return nil
}

func encodeUUID(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "encodeUUID", Types: []reflect.Type{tUUID}, Received: val}
	}

	b := val.Interface().(uuid.UUID)

	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}
