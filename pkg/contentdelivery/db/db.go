package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/language"
)

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

func Connect(ctx context.Context, connectionstring string) (*mongo.Client, error) {
	c, err := mongo.Connect(ctx, options.Client().SetRegistry(registry), options.Client().ApplyURI(connectionstring))

	// defer func() {
	// 	if err = c.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	if err != nil {
		return nil, err
	}

	return c, nil
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
