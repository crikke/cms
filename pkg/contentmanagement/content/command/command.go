package command

import (
	"context"
	"errors"
	"strings"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
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
}

func (h UpdateContentHandler) Handle(ctx context.Context, cmd UpdateContent) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.Id, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, c.ContentDefinitionID)

		if err != nil {
			return nil, err
		}

		properties := make(map[string]contentdefinition.PropertyDefinition)

		for _, pd := range cd.Properties {
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

			// todo set to default language
			lang := f.Language
			// if propertydefinition is not localized, field language should be undefinied
			if !pd.Localized && lang != language.Und {
				return nil, errors.New("cannot set localized value on unlocalized property")
			}

			// todo: ensure field & property value is same type
			c.Properties[lang][field] = f.Value
		}

		return c, nil
	})
}
