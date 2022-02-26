package command

import (
	"context"
	"errors"
	"strings"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
	"github.com/crikke/cms/pkg/siteconfiguration"
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
		Status:              content.Draft,
	}

	id, err := h.ContentRepository.CreateContent(ctx, c)

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

type UpdateContent struct {
	Id     uuid.UUID
	Fields []struct {
		Language string
		Field    string
		Value    interface{}
	}
}

type UpdateContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
	SiteConfiguration           *siteconfiguration.SiteConfiguration
}

func (h UpdateContentHandler) Handle(ctx context.Context, cmd UpdateContent) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.Id, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, c.ContentDefinitionID)

		if err != nil {
			return nil, err
		}

		properties := make(map[string]contentdefinition.PropertyDefinition)
		for _, pd := range cd.Propertydefinitions {
			properties[strings.ToLower(pd.Name)] = pd
		}

		if c.Properties == nil {
			c.Properties = make(map[string]map[string]interface{})
		}

		for _, configuredLanguage := range h.SiteConfiguration.Languages {

			c.Properties[configuredLanguage.String()] = make(map[string]interface{})
		}

		// check if field exists in contentdefinition
		for _, f := range cmd.Fields {

			field := strings.ToLower(f.Field)
			// name & urlsegment is not handled here, they are handled separately
			if field == content.NameField || field == content.UrlSegmentField {

				c.Properties[f.Language][field] = f.Value
				continue
			}

			if f.Field == "" {
				return nil, errors.New("property with empty name")
			}

			pd, ok := properties[field]
			if !ok {
				return nil, errors.New("property does not exist on propertydefinition")
			}

			lang := h.SiteConfiguration.Languages[0].String()

			// if prop is not localized and
			// field.language exist and is not default
			if !pd.Localized {
				if f.Language != "" && f.Language != h.SiteConfiguration.Languages[0].String() {
					return nil, errors.New(content.ErrUnlocalizedPropLocalizedValue)
				}
			}

			if f.Language != "" {
				exists := false
				for _, l := range h.SiteConfiguration.Languages {

					if l.String() == f.Language {
						exists = true
						break
					}
				}

				if !exists {
					return nil, errors.New(content.ErrNotConfiguredLocale)
				}

				lang = f.Language
			}

			c.Properties[lang][field] = f.Value
		}

		// ensure that content name is set for at least default language
		// if not set for other languages, it is set to default language
		// todo validate name
		defaultName, ok := c.Properties[h.SiteConfiguration.Languages[0].String()][content.NameField]
		if !ok {
			return nil, errors.New("content name cannot be empty for configured default language")
		}

		for _, l := range h.SiteConfiguration.Languages[1:] {
			_, ok := c.Properties[l.String()][content.NameField]

			if !ok {
				c.Properties[l.String()][content.NameField] = defaultName
			}
		}

		for _, l := range h.SiteConfiguration.Languages {
			url, ok := c.Properties[l.String()][content.UrlSegmentField]

			if !ok {
				url = c.Properties[l.String()][content.NameField]
			}

			str, ok := url.(string)
			if !ok {
				return nil, errors.New("urlsegment is not of type string")
			}
			str = strings.Replace(str, " ", "-", -1)
			c.Properties[l.String()][content.UrlSegmentField] = str
		}

		c.Version++
		return c, nil
	})
}

type PublishContent struct {
	ContentID uuid.UUID
}

type PublishContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
	SiteConfiguration           *siteconfiguration.SiteConfiguration
}

func (h PublishContentHandler) Handle(ctx context.Context, cmd PublishContent) error {

	cont, err := h.ContentRepository.GetContent(ctx, cmd.ContentID)

	if err != nil {
		return err
	}

	contentDefinition, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cont.ContentDefinitionID)

	if err != nil {
		return err
	}

	for _, pd := range contentDefinition.Propertydefinitions {

		var propvalues []interface{}
		var validators []validator.Validator

		for typ, v := range pd.Validators {
			val, err := validator.Parse(typ, v)

			if err != nil {
				return err
			}

			validators = append(validators, val)
		}

		if pd.Localized {

			for _, l := range h.SiteConfiguration.Languages {

				p := getPropertyValue(cont, pd.Name, l.String())
				propvalues = append(propvalues, p)
			}
		} else {
			p := getPropertyValue(cont, pd.Name, h.SiteConfiguration.Languages[0].String())
			propvalues = append(propvalues, p)

		}

		for _, value := range propvalues {
			for _, v := range validators {
				err := v.Validate(ctx, value)

				if err != nil {
					return err
				}
			}
		}
	}

	return h.ContentRepository.UpdateContent(ctx, cmd.ContentID, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		c.Status = content.Published
		return c, nil
	})
}

func getPropertyValue(c content.Content, name, locale string) interface{} {

	properties, ok := c.Properties[locale]

	if !ok {
		return nil
	}

	return properties[name]
}
