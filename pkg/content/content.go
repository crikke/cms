package content

import (
	"errors"
	"strings"
	"time"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/google/uuid"
)

type ContentVersion struct {
	ContentID uuid.UUID     `bson:"contentId"`
	Version   int           `bson:"version"`
	Status    PublishStatus `bson:"status"`
}

// swagger:enum PublishStatus
type PublishStatus string

type Content struct {
	ID                  uuid.UUID   `bson:"_id"`
	ContentDefinitionID uuid.UUID   `bson:"contentdefinition_id"`
	Data                ContentData `bson:"data"`
	Created             time.Time   `bson:"created"`
	Updated             time.Time   `bson:"updated"`
}

// swagger: model ContentData
type ContentData struct {
	ContentID  uuid.UUID       `bson:"contentId"`
	Version    int             `bson:"version"`
	Properties ContentLanguage `bson:"properties"`
	// TODO: does ContentData need a Created Field?
	Created time.Time     `bson:"created"`
	Status  PublishStatus `bson:"status"`
	// Tag IDs
	Tags []string `bson:"tags,omitempty"`
}

//! TODO: Is it better to handle localized values in field directly?
//! This could make it easier to query content, since only one ContentFields would be needed to be fetched.
//! Currently to get content for given locale:
//! - Get ContentVersion
//! - Get ContentLanguage with default locale
//! - Get ContentLanguage for current locale
//! - For each field in default locale
//! - 	if field is localize
//! -	return field from current locale
//!
//! Instead this could be done by moving the map from ContentVersion.Properties to ContentField.Value
//! This would help with when filtering which fields to return, since the field only exist once.
//!
//! IMPORTANT!! However, this would make it difficult when updating content since now each language can be treated as its own object,
//! So when doing a http PUT to update some localized contentdata, replace the whole object.
//! By doing above, the contentdata cannot be just replaced. Instead every property in every language needs to be sent and updated at once.

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

type ContentFactory struct {
}

func (f ContentFactory) NewContent(contentDefinition contentdefinition.ContentDefinition, defaultLanguage string) Content {

	c := Content{
		ContentDefinitionID: contentDefinition.ID,
		Data: ContentData{
			Status:     Draft,
			Created:    time.Now(),
			Properties: make(ContentLanguage),
		},
	}

	f.AddLanguage(&c.Data, defaultLanguage, true, contentDefinition)

	return c
}

// Language is added to ContentData
func (f ContentFactory) AddLanguage(c *ContentData, language string, addLocalized bool, contentDefinition contentdefinition.ContentDefinition) error {

	if c.Status != Draft {
		return errors.New(ErrNotDraft)
	}

	if c.Properties == nil {
		c.Properties = make(ContentLanguage)
	}

	// check if language already exists on contentdata
	if _, ok := c.Properties[language]; ok {
		return errors.New("language already exist")
	}

	// if propertydefinition is localized add if localized
	// if propertydefinition is not localized add
	cf := make(ContentFields)
	for name, val := range contentDefinition.Propertydefinitions {
		if val.Localized && !addLocalized {
			continue
		}

		cf[name] = ContentField{
			ID:        val.ID,
			Type:      val.Type,
			Localized: val.Localized,
		}
	}

	c.Properties[language] = cf
	return nil
}

// Creates a new content version from an existing version.
func (f ContentFactory) NewContentVersion(c *Content, contentDefinition contentdefinition.ContentDefinition, version int, defaultLanguage string) (*ContentData, error) {

	old := c.Data

	contentData := &ContentData{
		Status:     Draft,
		Created:    time.Now(),
		Properties: make(ContentLanguage),
	}

	for lang, oldFields := range old.Properties {

		addLocalized := lang == defaultLanguage
		f.AddLanguage(contentData, lang, addLocalized, contentDefinition)

		lookupfields := map[uuid.UUID]ContentField{}
		for fieldname, field := range oldFields {

			// happends if the fields name has changed
			// checking for ID is edge case when the old fields name has changed and a new field has the old fields name.00
			if newfield, ok := contentData.Properties[lang][fieldname]; ok && field.ID == newfield.ID {
				newfield.Value = field.Value
			} else {
				lookupfields[field.ID] = field
			}
		}

		// find all the fields with changed names by id
		for newName, field := range contentDefinition.Propertydefinitions {

			// if there is no match, the field is deleted from the contentdefinition
			if match, ok := lookupfields[field.ID]; ok {

				existingField := oldFields[newName]

				existingField.Value = match.Value
				contentData.Properties[lang][newName] = existingField
			}
		}
	}

	return contentData, nil
}

func (c ContentData) AvailableLanguages() []string {

	res := make([]string, 0)
	for lang := range c.Properties {
		res = append(res, lang)
	}

	return res
}

func (c ContentData) CanEdit() bool {
	return c.Status == Draft
}

func (f ContentFactory) SetField(cv *ContentData, lang, fieldname string, value interface{}) error {

	normalizedFieldname := strings.ToLower(fieldname)
	if !cv.CanEdit() {
		return errors.New(ErrNotDraft)
	}
	if lang == "" {
		return errors.New(ErrMissingLanguage)
	}

	cf, ok := cv.Properties[lang]

	if !ok {
		return errors.New(ErrMissingLanguage)
	}

	// instead of having to check if language is default language, the property should only exist on the default contentlanguage
	field, ok := cf[normalizedFieldname]

	if !ok {
		return errors.New(ErrMissingField)
	}

	field.Value = value
	cf[normalizedFieldname] = field

	return nil
}
