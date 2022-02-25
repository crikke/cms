package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_CreateContent(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	cdRepo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := cdRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test2",
	})
	assert.NoError(t, err)

	cmd := CreateContent{
		ContentDefinitionId: cid,
	}
	handler := CreateContentHandler{
		ContentDefinitionRepository: cdRepo,
		ContentRepository:           content.NewContentRepository(c),
	}

	contentId, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, contentId)
}

func Test_CreateContent_Empty_ContentDefinition(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	cmd := CreateContent{}
	handler := CreateContentHandler{
		ContentDefinitionRepository: contentdefinition.NewContentDefinitionRepository(c),
		ContentRepository:           content.NewContentRepository(c),
	}

	contentId, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Equal(t, uuid.UUID{}, contentId)
}

func Test_UpdateContent(t *testing.T) {

	cfg := siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("en-US"),
		},
	}

	tests := []struct {
		name       string
		contentdef *contentdefinition.ContentDefinition
		fields     []struct {
			Language string
			Field    string
			Value    interface{}
		}
		expectedErr    string
		expectedValues map[string]map[string]interface{}
	}{
		{
			name: "common fields only all configured languages",
			contentdef: &contentdefinition.ContentDefinition{
				Name:                "test",
				ID:                  uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
				{
					Language: "en-US",
					Field:    content.NameField,
					Value:    "name en",
				},
				{
					Language: "en-US",
					Field:    content.UrlSegmentField,
					Value:    "url en",
				},
			},
			expectedValues: map[string]map[string]interface{}{
				"sv-SE": {
					content.NameField:       "name sv",
					content.UrlSegmentField: "url-sv",
				},
				"en-US": {
					content.NameField:       "name en",
					content.UrlSegmentField: "url-en",
				},
			},
		},
		{
			name: "common fields only default language",
			contentdef: &contentdefinition.ContentDefinition{
				Name:                "test",
				ID:                  uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
			},
			expectedValues: map[string]map[string]interface{}{
				"sv-SE": {
					content.NameField:       "name sv",
					content.UrlSegmentField: "url-sv",
				},
				"en-US": {
					content.NameField:       "name sv",
					content.UrlSegmentField: "name-sv",
				},
			},
		},
		{
			name: "not localized field with empty language ok",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:        uuid.New(),
						Name:      "not_localized",
						Type:      "text",
						Localized: false,
					},
				},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
				{
					Language: "sv-SE",
					Field:    "not_localized",
					Value:    "ok",
				},
			},
			expectedValues: map[string]map[string]interface{}{
				"sv-SE": {
					"not_localized": "ok",
				},
				"en-US": {
					"not_localized": nil,
				},
			},
		},
		{
			name: "localized field with not configured language should return error",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:        uuid.New(),
						Name:      "localized_field",
						Type:      "text",
						Localized: true,
					},
				},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
				{
					Language: "nb-NO",
					Field:    "localized_field",
					Value:    "ok",
				},
			},
			expectedErr: content.ErrNotConfiguredLocale,
		},
		{
			name: "localized field ok",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:        uuid.New(),
						Name:      "localized_field",
						Type:      "text",
						Localized: true,
					},
				},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
				{
					Language: "en-US",
					Field:    "localized_field",
					Value:    "ok",
				},
			},
			expectedValues: map[string]map[string]interface{}{
				"en-US": {
					"localized_field": "ok",
				},
			},
		},

		{
			name: "not localized field with not default language should return error",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:        uuid.New(),
						Name:      "not_localized",
						Type:      "text",
						Localized: false,
					},
				},
			},
			fields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url sv",
				},
				{
					Language: "en-US",
					Field:    "not_localized",
					Value:    "not ok",
				},
			},
			expectedErr: content.ErrUnlocalizedPropLocalizedValue,
		},
		// {
		// 	name: "field with wrong type should return error",
		// },
	}

	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
			c.Database("cms").Collection("content").Drop(context.Background())

			cdRepo := contentdefinition.NewContentDefinitionRepository(c)
			contentdefinitionId, err := cdRepo.CreateContentDefinition(context.Background(), test.contentdef)
			assert.NoError(t, err)

			contentRepo := content.NewContentRepository(c)

			contentId, err := contentRepo.CreateContent(context.Background(), content.Content{
				ContentDefinitionID: contentdefinitionId,
			})
			assert.NoError(t, err)

			cmd := UpdateContent{
				Id:     contentId,
				Fields: test.fields,
			}

			handler := UpdateContentHandler{
				ContentDefinitionRepository: cdRepo,
				ContentRepository:           contentRepo,
				SiteConfiguration:           &cfg,
			}
			err = handler.Handle(context.Background(), cmd)

			if test.expectedErr != "" {
				assert.Equal(t, test.expectedErr, err.Error())
			}

			cont, err := contentRepo.GetContent(context.Background(), contentId)
			assert.NoError(t, err)

			for lang, fields := range test.expectedValues {

				for field, value := range fields {

					assert.Equal(t, value, cont.Properties[lang][field], value)
				}
			}
		})

	}

	// assert.Equal(t, "name test", cont.Properties[cfg.Languages[0].String()][content.NameField])
	// assert.Equal(t, "url-test", cont.Properties[cfg.Languages[0].String()][content.UrlSegmentField])
	// assert.Equal(t, "name test", cont.Properties[cfg.Languages[1].String()][content.NameField])
	// assert.Equal(t, "name-test", cont.Properties[cfg.Languages[1].String()][content.UrlSegmentField])
	// assert.Equal(t, "testfield 123", cont.Properties[cfg.Languages[0].String()]["field_not_localized"])
	// assert.Equal(t, nil, cont.Properties[cfg.Languages[1].String()]["field_not_localized"])
}
