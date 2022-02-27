/*
	Test TODOS:

	Tests involving version & status.
	Publish new version previous version should be set to  unpublished
	Updating a published version, the new version should be status draft
*/
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
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	cdRepo := contentdefinition.NewContentDefinitionRepository(c)
	contentRepo := content.NewContentRepository(c)

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
	actual, err := contentRepo.GetContent(context.Background(), contentId)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, contentId)
	_, ok := actual.Version[0]
	assert.True(t, ok)
}

func Test_CreateContent_Empty_ContentDefinition(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	contentRepo := content.NewContentRepository(c)
	cmd := CreateContent{}
	handler := CreateContentHandler{
		ContentDefinitionRepository: contentdefinition.NewContentDefinitionRepository(c),
		ContentRepository:           contentRepo,
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
				Version: map[int]content.ContentVersion{
					0: {
						Status: content.Draft,
					},
				},
			})
			assert.NoError(t, err)

			cmd := UpdateContent{
				Id:      contentId,
				Fields:  test.fields,
				Version: 0,
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

					assert.Equal(t, value, cont.Version[0].Properties[lang][field], value)
				}
			}
		})

	}
}

func Test_PublishContent(t *testing.T) {

	cfg := &siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("en-US"),
		},
	}

	tests := []struct {
		name        string
		contentdef  *contentdefinition.ContentDefinition
		content     content.Content
		expectedErr string
		expected    content.Content
	}{
		{
			name: "required field not set should return error",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:   uuid.New(),
						Name: "required_field",
						Type: "text",
						Validators: map[string]interface{}{
							"required": true,
						},
					},
				},
			},
			content: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Status: content.Draft,
						Properties: map[string]map[string]interface{}{
							"sv-SE": {
								content.NameField: "name sv",
							},
						},
					},
				},
			},

			expectedErr: "required",
		},
		{
			name: "required field set should return ok",
			contentdef: &contentdefinition.ContentDefinition{
				Name: "test",
				ID:   uuid.New(),
				Propertydefinitions: []contentdefinition.PropertyDefinition{
					{
						ID:   uuid.New(),
						Name: "required_field",
						Type: "text",
						Validators: map[string]interface{}{
							"required": true,
						},
					},
				},
			},
			content: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Properties: map[string]map[string]interface{}{
							"sv-SE": {
								content.NameField: "name sv",
								"required_field":  "ok",
							},
						},
					},
				},
			},
			expected: content.Content{
				PublishedVersion: 0,
				Version: map[int]content.ContentVersion{
					0: {
						Status: content.Published,
					},
				},
			},
		},
		// {
		// 	name: "new version is published",
		// 	contentdef: &contentdefinition.ContentDefinition{
		// 		Name: "test",
		// 		ID:   uuid.New(),
		// 		Propertydefinitions: []contentdefinition.PropertyDefinition{
		// 			{
		// 				ID:   uuid.New(),
		// 				Name: "required_field",
		// 				Type: "text",
		// 				Validators: map[string]interface{}{
		// 					"required": true,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	content: content.Content{
		// 		PublishedVersion: 0,
		// 		Version: map[int]content.ContentVersion{
		// 			0: {
		// 				Properties: map[string]map[string]interface{}{
		// 					"sv-SE": {
		// 						content.NameField: "name sv",
		// 						"required_field":  "ok",
		// 					},
		// 				},
		// 				Status: content.Published,
		// 			},
		// 			1: {
		// 				Properties: map[string]map[string]interface{}{
		// 					"sv-SE": {
		// 						content.NameField: "name sv",
		// 						"required_field":  "updated",
		// 					},
		// 				},
		// 				Status: content.Draft,
		// 			},
		// 		},
		// 	},
		// 	expected: content.Content{
		// 		PublishedVersion: 1,
		// 		Version: map[int]content.ContentVersion{
		// 			0: {
		// 				Status: content.PreviouslyPublished,
		// 			},
		// 			1: {
		// 				Status: content.Published,
		// 			},
		// 		},
		// 	},
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
			test.content.ContentDefinitionID = contentdefinitionId
			id, err := contentRepo.CreateContent(context.Background(), test.content)
			assert.NoError(t, err)

			cmd := PublishContent{
				ContentID: id,
			}

			handler := PublishContentHandler{
				ContentDefinitionRepository: cdRepo,
				ContentRepository:           contentRepo,
				SiteConfiguration:           cfg,
			}

			err = handler.Handle(context.Background(), cmd)
			if test.expectedErr != "" {
				assert.Equal(t, test.expectedErr, err.Error())
			}
			actual, err := contentRepo.GetContent(context.Background(), id)
			assert.NoError(t, err)

			assert.Equal(t, test.expected.PublishedVersion, actual.PublishedVersion)

			for v, contentver := range test.expected.Version {
				assert.Equal(t, actual.Version[v].Status, contentver.Status, "status")

				for lang, fields := range contentver.Properties {
					for field, value := range fields {
						assert.Equal(t, value, actual.Version[v].Properties[lang][field], v, lang, field)
					}
				}
			}
		})
	}
}
