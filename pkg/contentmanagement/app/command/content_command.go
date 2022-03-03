package command

import (
	"context"
	"time"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionId uuid.UUID
}

type CreateContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
}

func (h CreateContentHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cmd.ContentDefinitionId)

	if err != nil {
		return uuid.UUID{}, err
	}

	c := content.Content{
		ContentDefinitionID: cd.ID,
		Version: map[int]content.ContentVersion{
			0: {
				Created: time.Now().UTC(),
			},
		},
	}

	id, err := h.ContentRepository.CreateContent(ctx, c)

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

type UpdateField struct {
	Id       uuid.UUID
	Version  int
	Language string
	Field    string
	Value    interface{}
}

type UpdateFieldHandler struct {
	ContentRepository           content.ContentRepository
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	Factory                     content.Factory
}

func (h UpdateFieldHandler) Handle(ctx context.Context, cmd UpdateField) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.Id, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		cv := c.Version[cmd.Version]

		err := h.Factory.SetField(&cv, cmd.Language, cmd.Field, cmd.Value)

		if err != nil {
			return nil, err
		}

		return c, nil
	})
}

// type UpdateContent struct {
// 	Id      uuid.UUID
// 	Version int
// 	Fields  []struct {
// 		Language string
// 		Field    string
// 		Value    interface{}
// 	}
// }

// type UpdateContentHandler struct {
// 	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
// 	ContentRepository           content.ContentRepository
// 	SiteConfiguration           *siteconfiguration.SiteConfiguration
// }

// func (h UpdateContentHandler) Handle(ctx context.Context, cmd UpdateContent) error {

// 	return h.ContentRepository.UpdateContent(ctx, cmd.Id, func(ctx context.Context, c *content.Content) (*content.Content, error) {
// 		cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, c.ContentDefinitionID)

// 		if err != nil {
// 			return nil, err
// 		}

// 		properties := make(map[string]contentdefinition.PropertyDefinition)
// 		for _, pd := range cd.Propertydefinitions {
// 			properties[strings.ToLower(pd.Name)] = pd
// 		}

// 		contentVer := content.ContentVersion{}
// 		// get the version of which the update should be based on
// 		if cv, ok := c.Version[cmd.Version]; ok {

// 			// this is necessary to dereference contentver
// 			contentVer.Created = cv.Created
// 			contentVer.Properties = make(map[string]map[string]interface{})
// 			for lng, fields := range cv.Properties {
// 				contentVer.Properties[lng] = make(map[string]interface{})
// 				for field, val := range fields {
// 					contentVer.Properties[lng][field] = val
// 				}
// 			}

// 		} else {
// 			return nil, errors.New(content.ErrVersionNotExists)
// 		}

// 		if contentVer.Properties == nil {
// 			contentVer.Properties = make(map[string]map[string]interface{})
// 		}

// 		for _, configuredLanguage := range h.SiteConfiguration.Languages {

// 			contentVer.Properties[configuredLanguage.String()] = make(map[string]interface{})
// 		}

// 		// check if field exists in contentdefinition
// 		for _, f := range cmd.Fields {

// 			field := strings.ToLower(f.Field)
// 			// name & urlsegment is not handled here, they are handled separately
// 			if field == content.NameField || field == content.UrlSegmentField {

// 				contentVer.Properties[f.Language][field] = f.Value
// 				continue
// 			}

// 			if f.Field == "" {
// 				return nil, errors.New("property with empty name")
// 			}

// 			pd, ok := properties[field]
// 			if !ok {
// 				return nil, errors.New("property does not exist on propertydefinition")
// 			}

// 			lang := h.SiteConfiguration.Languages[0].String()

// 			// if prop is not localized and
// 			// field.language exist and is not default
// 			if !pd.Localized {
// 				if f.Language != "" && f.Language != h.SiteConfiguration.Languages[0].String() {
// 					return nil, errors.New(content.ErrUnlocalizedPropLocalizedValue)
// 				}
// 			}

// 			if f.Language != "" {
// 				exists := false
// 				for _, l := range h.SiteConfiguration.Languages {

// 					if l.String() == f.Language {
// 						exists = true
// 						break
// 					}
// 				}

// 				if !exists {
// 					return nil, errors.New(content.ErrNotConfiguredLocale)
// 				}

// 				lang = f.Language
// 			}

// 			contentVer.Properties[lang][field] = f.Value
// 		}

// 		// ensure that content name is set for at least default language
// 		// if not set for other languages, it is set to default language
// 		// todo validate name
// 		defaultName, ok := contentVer.Properties[h.SiteConfiguration.Languages[0].String()][content.NameField]
// 		if !ok {
// 			return nil, errors.New("content name cannot be empty for configured default language")
// 		}

// 		for _, l := range h.SiteConfiguration.Languages[1:] {
// 			_, ok := contentVer.Properties[l.String()][content.NameField]

// 			if !ok {
// 				contentVer.Properties[l.String()][content.NameField] = defaultName
// 			}
// 		}

// 		for _, l := range h.SiteConfiguration.Languages {
// 			url, ok := contentVer.Properties[l.String()][content.UrlSegmentField]

// 			if !ok {
// 				url = contentVer.Properties[l.String()][content.NameField]
// 			}

// 			str, ok := url.(string)
// 			if !ok {
// 				return nil, errors.New("urlsegment is not of type string")
// 			}
// 			str = strings.Replace(str, " ", "-", -1)
// 			contentVer.Properties[l.String()][content.UrlSegmentField] = str
// 		}

// 		if cmd.Version == c.PublishedVersion && c.Status == content.Published {
// 			nextver := len(c.Version)
// 			contentVer.Created = time.Now().UTC()
// 			c.Version[nextver] = contentVer
// 		} else {
// 			c.Version[cmd.Version] = contentVer
// 		}

// 		return c, nil
// 	})
// }

// type PublishContent struct {
// 	ContentID uuid.UUID
// 	Version   int
// }

// type PublishContentHandler struct {
// 	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
// 	ContentRepository           content.ContentRepository
// 	SiteConfiguration           *siteconfiguration.SiteConfiguration
// }

// func (h PublishContentHandler) Handle(ctx context.Context, cmd PublishContent) error {

// 	cont, err := h.ContentRepository.GetContent(ctx, cmd.ContentID)

// 	if err != nil {
// 		return err
// 	}

// 	contentDefinition, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cont.ContentDefinitionID)

// 	if err != nil {
// 		return err
// 	}

// 	contentver, ok := cont.Version[cmd.Version]
// 	if !ok {
// 		return errors.New(content.ErrVersionNotExists)
// 	}

// 	for _, pd := range contentDefinition.Propertydefinitions {

// 		var propvalues []interface{}
// 		var validators []validator.Validator

// 		for typ, v := range pd.Validators {
// 			val, err := validator.Parse(typ, v)

// 			if err != nil {
// 				return err
// 			}

// 			validators = append(validators, val)
// 		}

// 		if pd.Localized {

// 			for _, l := range h.SiteConfiguration.Languages {

// 				p := getPropertyValue(contentver, pd.Name, l.String())
// 				propvalues = append(propvalues, p)
// 			}
// 		} else {
// 			p := getPropertyValue(contentver, pd.Name, h.SiteConfiguration.Languages[0].String())
// 			propvalues = append(propvalues, p)

// 		}

// 		for _, value := range propvalues {
// 			for _, v := range validators {
// 				err := v.Validate(ctx, value)

// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	}

// 	return h.ContentRepository.UpdateContent(ctx, cmd.ContentID, func(ctx context.Context, c *content.Content) (*content.Content, error) {

// 		// todo use ptr
// 		current := c.Version[c.PublishedVersion]
// 		// current.Status = content.PreviouslyPublished
// 		c.Version[c.PublishedVersion] = current

// 		// contentver.Status = content.Published
// 		c.PublishedVersion = cmd.Version
// 		c.Version[cmd.Version] = contentver
// 		return c, nil
// 	})
// }

// func getPropertyValue(c content.ContentVersion, name, locale string) interface{} {

// 	properties, ok := c.Properties[locale]

// 	if !ok {
// 		return nil
// 	}

// 	return properties[name]
// }
