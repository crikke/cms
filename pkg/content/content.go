package content

import (
	"errors"
	"strings"
	"time"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

// swagger:enum PublishStatus
type PublishStatus string

type Content struct {
	ID                  uuid.UUID `bson:"_id"`
	ContentDefinitionID uuid.UUID `bson:"contentdefinition_id"`
	ParentID            uuid.UUID `bson:"parentid"`
	PublishedVersion    int
	Version             map[int]ContentVersion `bson:"version"`
	// Which languages this content has been translated to
	AvailableLanguages []string      `bson:"version"`
	Status             PublishStatus `bson:"status"`
}

// swagger: model ContentVersion
type ContentVersion struct {
	Properties ContentLanguage `bson:"properties"`
	Created    time.Time       `bson:"created"`
	Status     PublishStatus   `bson:"status"`
}

// ContentField describes the property aswell as its value
// swagger: model ContentField
type ContentField struct {
	ID        uuid.UUID `bson:"id"`
	Type      string    `bson:"type,omitempty"`
	Localized bool      `bson:"localized,omitempty"`
	Value     interface{}
}

// ContentField is a map where key is field name and value is the ContentField
// swagger: model ContentFields
type ContentFields map[string]ContentField

// ContentLanguage is a map containing a given languages ContentFields
// swagger: model ContentLanguage
type ContentLanguage map[string]ContentFields

const (
	Draft               PublishStatus = "draft"
	Published           PublishStatus = "published"
	PreviouslyPublished PublishStatus = "previouslyPublished"
	Archived            PublishStatus = "archived"
)

type Factory struct {
	Cfg *siteconfiguration.SiteConfiguration
}

func (f Factory) getDefaultLanguage() string {
	return f.Cfg.Languages[0].String()
}

func (f Factory) NewContent(spec contentdefinition.ContentDefinition, parentID uuid.UUID) (Content, error) {

	c := Content{
		ContentDefinitionID: spec.ID,
		Version:             map[int]ContentVersion{},
		ParentID:            parentID,
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

	normalizedFieldname := strings.ToLower(fieldname)
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

	field, ok := cf[normalizedFieldname]

	if !ok {
		return errors.New(ErrMissingField)
	}

	if !field.Localized {
		if lang != f.getDefaultLanguage() {
			return errors.New(ErrUnlocalizedPropLocalizedValue)
		}
	}

	field.Value = value
	cf[normalizedFieldname] = field

	return nil
}
