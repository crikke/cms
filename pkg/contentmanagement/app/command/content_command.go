package command

import (
	"context"

	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/contentdefinition/validator"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
)

type CreateContent struct {
	ContentDefinitionId uuid.UUID
	WorkspaceId         uuid.UUID
}

type CreateContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentManagementRepository
	Factory                     content.Factory
}

func (h CreateContentHandler) Handle(ctx context.Context, cmd CreateContent) (uuid.UUID, error) {

	cd, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, cmd.ContentDefinitionId, cmd.WorkspaceId)
	if err != nil {
		return uuid.UUID{}, err
	}

	c, err := h.Factory.NewContent(cd)
	if err != nil {
		return uuid.UUID{}, err
	}

	return h.ContentRepository.CreateContent(ctx, c)
}

type UpdateContentFields struct {
	ContentID uuid.UUID
	Version   int
	Language  string

	Fields      map[string]interface{}
	WorkspaceId uuid.UUID
}

type UpdateContentFieldsHandler struct {
	ContentRepository           content.ContentManagementRepository
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	Factory                     content.Factory
}

func (h UpdateContentFieldsHandler) Handle(ctx context.Context, cmd UpdateContentFields) error {

	return h.ContentRepository.UpdateContentData(ctx, cmd.ContentID, cmd.Version, func(ctx context.Context, c *content.ContentData) (*content.ContentData, error) {

		// if this version is a draft, update it directly.
		// Otherwise create a new version based on this version.

		contentData := *c

		if c.Status != content.Draft {
			versions, err := h.ContentRepository.ListContentVersions(ctx, cmd.ContentID)
			if err != nil {
				return nil, err
			}

			contentData.Version = len(versions)
			contentData.Status = content.Draft
		}

		for f, v := range cmd.Fields {
			err := h.Factory.SetField(c, cmd.Language, f, v)
			if err != nil {
				return nil, err
			}
		}

		return &contentData, nil
	})
}

type PublishContent struct {
	ContentID   uuid.UUID
	Version     int
	WorkspaceId uuid.UUID
}

type PublishContentHandler struct {
	ContentDefinitionRepository contentdefinition.ContentDefinitionRepository
	ContentRepository           content.ContentManagementRepository
	SiteConfiguration           *siteconfiguration.SiteConfiguration
}

func (h PublishContentHandler) Handle(ctx context.Context, cmd PublishContent) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.ContentID, func(ctx context.Context, c *content.Content) (*content.Content, error) {
		previousVersion := c.Data.Version

		contentDefinition, err := h.ContentDefinitionRepository.GetContentDefinition(ctx, c.ContentDefinitionID, cmd.WorkspaceId)
		if err != nil {
			return nil, err
		}

		// set new version to status published
		err = h.ContentRepository.UpdateContentData(ctx, cmd.ContentID, cmd.Version, func(ctx context.Context, cd *content.ContentData) (*content.ContentData, error) {

			for propName, pd := range contentDefinition.Propertydefinitions {

				var propvalues []interface{}
				var validators []validator.Validator

				for typ, v := range pd.Validators {
					val, err := validator.Parse(typ, v)

					if err != nil {
						return nil, err
					}

					validators = append(validators, val)
				}

				if pd.Localized {

					for _, l := range h.SiteConfiguration.Languages {

						p := getPropertyValue(*cd, propName, l.String())
						propvalues = append(propvalues, p)
					}
				} else {
					p := getPropertyValue(*cd, propName, h.SiteConfiguration.Languages[0].String())
					propvalues = append(propvalues, p)

				}

				for _, value := range propvalues {
					for _, v := range validators {
						err := v.Validate(ctx, value)

						if err != nil {
							return nil, err
						}
					}
				}
			}

			cd.Status = content.Published
			c.Data = *cd
			return cd, nil
		})
		if err != nil {
			return nil, err
		}

		// set previous version to previouslypublished
		if previousVersion != cmd.Version {
			err = h.ContentRepository.UpdateContentData(ctx, cmd.ContentID, previousVersion, func(ctx context.Context, cd *content.ContentData) (*content.ContentData, error) {
				cd.Status = content.PreviouslyPublished
				return cd, nil
			})
			if err != nil {
				return nil, err
			}
		}
		return c, nil
	})
}

func getPropertyValue(c content.ContentData, name, locale string) interface{} {

	properties, ok := c.Properties[locale]

	if !ok {
		return nil
	}

	return properties[name].Value
}

type ArchiveContent struct {
	ID          uuid.UUID
	WorkspaceId uuid.UUID
}
type ArchiveContentHandler struct {
	ContentRepository content.ContentManagementRepository
}

func (h ArchiveContentHandler) Handle(ctx context.Context, cmd ArchiveContent) error {

	return h.ContentRepository.UpdateContent(ctx, cmd.ID, func(ctx context.Context, c *content.Content) (*content.Content, error) {

		err := h.ContentRepository.UpdateContentData(ctx, cmd.ID, c.Data.Version, func(ctx context.Context, cd *content.ContentData) (*content.ContentData, error) {
			cd.Status = content.PreviouslyPublished
			return cd, nil
		})

		if err != nil {
			return nil, err
		}
		return c, nil
	})
}
