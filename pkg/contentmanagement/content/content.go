package content

import (
	"errors"
	"time"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

/*

Saving content example:

HTTP POST contentmanagement/content
{
	"contentdefinition": uuid,
}

HTTP PUT contentmanagement/content/{contentID}
{
	"name": {
		"sv_SE": "exempel",
		"en_US": "example"
	},
	"properties": [
		{
			"id": uuid
			"value": {
				"sv_SE": "example value"
			},
			"localized": false|true // todo should PUT contain this? isnt contentdefinition responsible for this?
		}
	]
}

*/

const (
	NameField       = "name"
	UrlSegmentField = "url"
)

// swagger:enum PublishStatus
type PublishStatus string

type Content struct {
	ID                  uuid.UUID `bson:"_id"`
	ContentDefinitionID uuid.UUID `bson:"contentdefinition_id"`
	ParentID            uuid.UUID `bson:"parentid"`
	PublishedVersion    int
	Version             map[int]ContentVersion `bson:"version"`
	Status              PublishStatus          `bson:"status"`
}

type ContentVersion struct {
	Properties ContentLanguage `bson:"properties"`
	Created    time.Time       `bson:"created"`
	Status     PublishStatus   `bson:"status"`
}
type ContentField struct {
	ID        uuid.UUID `bson:"id"`
	Type      string    `bson:"type,omitempty"`
	Localized bool      `bson:"localized,omitempty"`
	Value     interface{}
}

type ContentFields map[string]ContentField
type ContentLanguage map[string]ContentFields

const (
	Draft     PublishStatus = "draft"
	Published PublishStatus = "published"
	Archived  PublishStatus = "archived"
)

type Factory struct {
	Cfg *siteconfiguration.SiteConfiguration
}

func (f Factory) getDefaultLanguage() string {
	return f.Cfg.Languages[0].String()
}

func (f Factory) NewContent(spec contentdefinition.ContentDefinition) (Content, error) {

	c := Content{
		ContentDefinitionID: spec.ID,
		Version:             map[int]ContentVersion{},
		Status:              Draft,
	}

	cv := ContentVersion{
		Status:     Draft,
		Created:    time.Now(),
		Properties: make(ContentLanguage),
	}
	cf := make(ContentFields)
	for name, val := range spec.Propertydefinitions {
		cf[name] = ContentField{
			ID:        val.ID,
			Type:      val.Type,
			Localized: val.Localized,
		}
	}

	c.Version[0] = cv
	cv.Properties[f.getDefaultLanguage()] = cf

	return c, nil
}

func (f Factory) AddLanguage(c *ContentVersion, language string) (ContentFields, error) {

	if c.Status != Draft {
		return nil, errors.New(ErrNotDraft)
	}

	cl := make(ContentFields)
	if _, ok := c.Properties[language]; !ok {
		c.Properties[language] = cl
	}
	return cl, nil
}

func (f Factory) NewContentVersion(c *Content, contentDefinition contentdefinition.ContentDefinition, version int) (*ContentVersion, error) {

	existing, ok := c.Version[version]

	if !ok {
		return nil, errors.New(ErrMissingVersion)
	}

	cv := &ContentVersion{
		Status:     Draft,
		Created:    time.Now(),
		Properties: make(ContentLanguage),
	}

	for lang, cl := range existing.Properties {
		cf := make(ContentFields)

		// create new properties from propertydefinitions
		for name, val := range contentDefinition.Propertydefinitions {
			cf[name] = ContentField{
				ID:        val.ID,
				Type:      val.Type,
				Localized: val.Localized,
			}
		}

		lookupfields := map[uuid.UUID]ContentField{}
		for fieldname, field := range cl {

			// happends if the fields name has changed
			// checking for ID is edge case when the old fields name has changed and a new field has the old fields name.00
			if newfield, ok := cf[fieldname]; ok && field.ID == newfield.ID {
				newfield.Value = field.Value
			} else {
				lookupfields[field.ID] = field
			}
		}

		// find all the fields with changed names by id
		for name, field := range contentDefinition.Propertydefinitions {

			// if there is no match, the field is deleted from the contentdefinition
			if match, ok := lookupfields[field.ID]; ok {

				newfield := cl[name]

				newfield.Value = match.Value
				cl[name] = newfield
			}
		}

		cv.Properties[lang] = cf
	}

	c.Version[len(c.Version)] = *cv
	return cv, nil
}

func (c ContentVersion) CanEdit() bool {
	return c.Status == Draft
}

func (f Factory) SetField(cv *ContentVersion, lang, fieldname string, value interface{}) error {

	if !cv.CanEdit() {
		return errors.New(ErrNotDraft)
	}
	if lang == "" {
		lang = f.getDefaultLanguage()
	}

	cf, ok := cv.Properties[lang]

	if !ok {
		return errors.New(ErrMissingLanguage)
	}

	field, ok := cf[fieldname]

	if !ok {
		return errors.New(ErrMissingField)
	}

	if !field.Localized {
		if lang != f.getDefaultLanguage() {
			return errors.New(ErrUnlocalizedPropLocalizedValue)
		}
	}

	field.Value = value
	cf[fieldname] = field

	return nil
}
