package command

import (
	"context"
	"errors"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api"
	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionId uuid.UUID
	ParentID            uuid.UUID
}

type CreateContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentRepository
	Factory                     content.Factory
}

func (h CreateContentHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cmd.ContentDefinitionId)
	if err != nil {
		return uuid.UUID{}, err
	}

	// if empty parentid this is a root object
	if cmd.ParentID != (uuid.UUID{}) {

		parent, err := h.ContentRepository.GetContent(ctx, cmd.ParentID)

		if err != nil {
			return uuid.UUID{}, err
		}

		if parent.Version[parent.PublishedVersion].Status != content.Published {
			return uuid.UUID{}, errors.New("cannot create content under unpublished content")
		}
	}

	c, err := h.Factory.NewContent(cd, cmd.ParentID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return h.ContentRepository.CreateContent(ctx, c)
}

type UpdateContentFields struct {
	ContentID uuid.UUID
	Version   int
	Language  string

	Fields map[string]interface{}
	// Name  string
	// Value interface{}
}

type UpdateContentFieldsHandler struct {
	ContentRepository           content.ContentRepository
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	Factory                     content.Factory
}

func (h UpdateContentFieldsHandler) Handle(ctx context.Context, cmd UpdateContentFields) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.ContentID, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		cv, ok := c.Version[cmd.Version]

		if !ok {
			return nil, api.GenericError{
				Body: api.ErrorBody{
					Message:   content.ErrMissingVersion,
					FieldName: "version",
				},
				StatusCode: http.StatusNotFound,
			}
		}
		for f, v := range cmd.Fields {
			err := h.Factory.SetField(&cv, cmd.Language, f, v)
			if err != nil {
				return nil, err
			}
		}

		return c, nil
	})
}

type PublishContent struct {
	ContentID uuid.UUID
	Version   int
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

	contentver, ok := cont.Version[cmd.Version]
	if !ok {
		return errors.New(content.ErrMissingVersion)
	}

	for propName, pd := range contentDefinition.Propertydefinitions {

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

				p := getPropertyValue(contentver, propName, l.String())
				propvalues = append(propvalues, p)
			}
		} else {
			p := getPropertyValue(contentver, propName, h.SiteConfiguration.Languages[0].String())
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

		current := c.Version[c.PublishedVersion]
		current.Status = content.Archived
		c.Version[c.PublishedVersion] = current

		contentver.Status = content.Published
		c.PublishedVersion = cmd.Version
		c.Version[cmd.Version] = contentver
		return c, nil
	})
}

func getPropertyValue(c content.ContentVersion, name, locale string) interface{} {

	properties, ok := c.Properties[locale]

	if !ok {
		return nil
	}

	return properties[name].Value
}

type DeleteContent struct {
	ID uuid.UUID
}
type DeleteContentHandler struct {
}

func (h DeleteContentHandler) Handle(ctx context.Context, cmd DeleteContent) error {
	// TODO: Implement
	return nil
}
