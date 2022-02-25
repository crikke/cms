package command

import (
	"context"
	"errors"
	"strings"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"golang.org/x/text/language"
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
		Language language.Tag
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

		// check if field exists in contentdefinition
		for _, f := range cmd.Fields {

			field := strings.ToLower(f.Field)
			if f.Field == "" {
				return nil, errors.New("property with empty name")
			}

			pd, ok := properties[field]
			if !ok {
				return nil, errors.New("property does not exist on propertydefinition")
			}

			lang := h.SiteConfiguration.Languages[0]

			if f.Language != language.Und {

				exists := false
				for _, l := range h.SiteConfiguration.Languages {

					if l == f.Language {
						exists = true
						break
					}
				}

				if !exists {
					return nil, errors.New("language is not configured in siteconfiguration")
				}

				lang = f.Language
			}

			if !pd.Localized && lang != language.Und {
				return nil, errors.New("cannot set localized value on unlocalized property")
			}

			// todo: ensure field & property value is same type
			c.Properties[lang][field] = f.Value

			// also validate urlsegment

			// ensure that content name is set for at least default language
			// if not set for other languages, it is set to default language
			// todo validate name
			defaultName, ok := c.Properties[h.SiteConfiguration.Languages[0]][content.NameField]
			if !ok {
				return nil, errors.New("content name cannot be empty for configured default language")
			}

			for _, l := range h.SiteConfiguration.Languages[1:] {
				_, ok := c.Properties[l][content.NameField]

				if !ok {
					c.Properties[l][content.NameField] = defaultName
				}
			}

			for _, l := range h.SiteConfiguration.Languages {
				url, ok := c.Properties[l][content.UrlSegmentField]

				if !ok {
					url = c.Properties[l][content.NameField]
				}

				str, ok := url.(string)
				if !ok {
					return nil, errors.New("urlsegment is not of type string")
				}
				str = strings.Replace(str, " ", "-", -1)
				c.Properties[l][content.UrlSegmentField] = str
			}
		}

		return c, nil
	})
}
