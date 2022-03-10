package content

import (
	"errors"
	"strings"
	"time"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

type ContentVersion struct {
	ContentID uuid.UUID     `bson:"contentId"`
	Version   int           `bson:"version"`
	Status    PublishStatus `bson:"status"`
}

// swagger:enum PublishStatus
type PublishStatus string

//! TODO: instead of storing the versionnumber in PublishedVersion, store the actual version there.
//! This will always be the most frequent requested data and by doing this unnecessary logic to query the published version wont be needed.
//! Older versions should be moved to another collection since read speeds isn´t as important.
//! Obviously ContentVersion will need a version field.
//! This will also solve problem with status, since status will only be needed on ContentVersion and not Content
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
	Tags    []string      `bson:"tags,omitempty"`
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

func (f Factory) NewContent(spec contentdefinition.ContentDefinition) (Content, error) {

	c := Content{
		ContentDefinitionID: spec.ID,
	}

	cv := ContentData{
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

	c.Data = cv
	cv.Properties[f.getDefaultLanguage()] = cf

	return c, nil
}

func (f Factory) AddLanguage(c *ContentData, language string) (ContentFields, error) {

	if c.Status != Draft {
		return nil, errors.New(ErrNotDraft)
	}

	cl := make(ContentFields)
	if _, ok := c.Properties[language]; !ok {
		c.Properties[language] = cl
	}

	return cl, nil
}

func (f Factory) NewContentVersion(c *Content, contentDefinition contentdefinition.ContentDefinition, version int) (*ContentData, error) {

	existing := c.Data

	cv := &ContentData{
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

	return cv, nil
}

// func (c Content) GetPublishedVersion() (ContentData, error) {
// 	cv, ok := c.Version[c.Data]

// 	if !ok {
// 		return ContentData{}, errors.New(ErrMissingVersion)
// 	}

// 	if cv.Status != Published {
// 		return ContentData{}, errors.New("not published")
// 	}

// 	return cv, nil
// }

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

func (f Factory) SetField(cv *ContentData, lang, fieldname string, value interface{}) error {

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
