package content

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/crikke/cms/pkg/siteconfiguration"
// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/text/language"
// )

// type ContentRepository interface {
// 	GetContent(ctx context.Context, contentReference ContentReference) (Content, error)
// 	GetChildren(ctx context.Context, contentReference ContentReference) ([]Content, error)
// }

// type repository struct {
// 	database *mongo.Database
// 	cfg      *siteconfiguration.SiteConfiguration
// }

// type contentData struct {
// 	ID               uuid.UUID `bson:"_id"`
// 	ParentID         uuid.UUID
// 	PublishedVersion int
// 	Created          time.Time
// 	Updated          time.Time
// 	Data             map[int]contentVersion
// }
// type contentVersion struct {
// 	Properties []contentProperty
// 	Name       map[string]string
// 	URLSegment map[string]string
// }

// type contentProperty struct {
// 	ID        uuid.UUID
// 	Name      string
// 	Type      string
// 	Localized bool
// 	Value     map[string]interface{}
// }

// type ErrorResponse struct {
// 	ID         ContentReference
// 	Message    string
// 	StatusCode int
// }

// func (e ErrorResponse) Error() string {
// 	return e.Message
// }

// func ErrContentNotFound(ref ContentReference) ErrorResponse {
// 	return ErrorResponse{ID: ref, StatusCode: http.StatusNotFound, Message: ref.String()}
// }

// func NewContentRepository(c *mongo.Client, cfg *siteconfiguration.SiteConfiguration) ContentRepository {
// 	db := c.Database("cms")
// 	return repository{database: db, cfg: cfg}
// }

// func (r repository) GetContent(ctx context.Context, contentReference ContentReference) (Content, error) {

// 	doc := &contentData{}
// 	err := r.database.Collection("content").FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: contentReference.ID}}).Decode(doc)

// 	if err != nil {
// 		return Content{}, err
// 	}

// 	lang := []language.Tag{}

// 	if contentReference.Locale != nil {
// 		lang = append(lang, *contentReference.Locale)
// 	}
// 	lang = append(lang, r.cfg.Languages...)

// 	return r.transform(ctx, doc, lang)
// }

// func (r repository) GetChildren(ctx context.Context, contentReference ContentReference) ([]Content, error) {

// 	cur, err := r.database.Collection("content").Find(ctx, bson.D{primitive.E{Key: "parentID", Value: contentReference.ID}})

// 	if err != nil {
// 		return nil, err
// 	}

// 	lang := []language.Tag{}

// 	if contentReference.Locale != nil {
// 		lang = append(lang, *contentReference.Locale)
// 	}
// 	lang = append(lang, r.cfg.Languages...)

// 	result := []Content{}

// 	for cur.Next(ctx) {
// 		data := &contentData{}
// 		err = cur.Decode(data)

// 		if err != nil {
// 			return nil, err
// 		}

// 		doc, err := r.transform(ctx, data, lang)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, doc)
// 	}

// 	return result, nil
// }

// // if property is localized. Find localized value by priorited language.Tags
// // else get property by default language
// func (r repository) transform(ctx context.Context, in *contentData, lang []language.Tag) (Content, error) {

// 	defaultLang := r.cfg.Languages[0]

// 	result := Content{
// 		ID:       ContentReference{ID: in.ID, Locale: &lang[0], Version: nil},
// 		ParentID: in.ParentID,
// 		Created:  in.Created,
// 		Updated:  in.Updated,
// 	}

// 	// todo
// 	data, exist := in.Data[0]

// 	if !exist {
// 		return Content{}, ErrContentNotFound(result.ID)
// 	}

// 	result.URLSegment, exist = data.URLSegment[defaultLang.String()]

// 	// Localized content must have a URL segment for given locale.
// 	if !exist {
// 		return Content{}, ErrContentNotFound(result.ID)
// 	}

// 	result.Name = lookupLocalizedValueString(data.Name, lang...)

// 	for _, prop := range data.Properties {

// 		val := prop.Value[defaultLang.String()]
// 		if prop.Localized {

// 			val = lookupLocalizedValueInterface(prop.Value, lang...)
// 		}

// 		cp := Property{
// 			ID:        prop.ID,
// 			Name:      prop.Name,
// 			Type:      prop.Type,
// 			Localized: prop.Localized,
// 			Value:     val,
// 		}
// 		result.Properties = append(result.Properties, cp)
// 		continue
// 	}
// 	return result, nil
// }

// // todo use generics when available
// func lookupLocalizedValueString(in map[string]string, tags ...language.Tag) string {

// 	for _, tag := range tags {

// 		if str, exist := in[tag.String()]; exist {
// 			return str
// 		}
// 	}

// 	return ""
// }

// // todo use generics when available
// func lookupLocalizedValueInterface(in map[string]interface{}, tags ...language.Tag) interface{} {

// 	for _, tag := range tags {

// 		if val, exist := in[tag.String()]; exist {
// 			return val
// 		}
// 	}

// 	return nil
// }
