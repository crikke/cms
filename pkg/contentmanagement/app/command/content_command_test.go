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
	"time"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

var (
	emptyContentDef contentdefinition.ContentDefinition = contentdefinition.ContentDefinition{
		Name:                "test",
		ID:                  uuid.New(),
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	}

	reqfieldContentDef contentdefinition.ContentDefinition = contentdefinition.ContentDefinition{
		Name: "test",
		ID:   uuid.New(),
		Propertydefinitions: map[string]contentdefinition.PropertyDefinition{

			"required_field": {
				ID:   uuid.New(),
				Type: "text",
				Validators: map[string]interface{}{
					"required": true,
				},
			},
		},
	}

	emptyContent content.Content = content.Content{
		Version: map[int]content.ContentVersion{
			0: {
				Status:  content.Draft,
				Created: time.Now().UTC(),
			},
		},
	}
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
		Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
			content.NameField: {
				ID:          uuid.New(),
				Description: "Content name",
				Type:        "text",
				Localized:   true,
				Validators: map[string]interface{}{
					"required": validator.RequiredRule(true),
				},
			},
		},
	})

	assert.NoError(t, err)

	cmd := CreateContent{
		ContentDefinitionId: cid,
	}
	handler := CreateContentHandler{
		ContentDefinitionRepository: cdRepo,
		ContentRepository:           content.NewContentRepository(c),
		Factory:                     content.Factory{Cfg: &siteconfiguration.SiteConfiguration{Languages: []language.Tag{language.MustParse("sv-SE")}}},
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
		name         string
		contentdef   *contentdefinition.ContentDefinition
		existing     content.Content
		updatefields []struct {
			Language string
			Field    string
			Value    interface{}
		}
		version     int
		expectedErr string
		expected    content.Content
	}{
		{
			name:       "common fields only all configured languages",
			contentdef: &emptyContentDef,
			existing: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Status:  content.Draft,
						Created: time.Now().UTC(),
						Properties: content.ContentLanguage{
							"sv-SE": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "",
								},
								content.UrlSegmentField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "",
								},
							},
							"en-US": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "",
								},
								content.UrlSegmentField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "",
								},
							},
						},
					},
				},
			},
			updatefields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "sv-SE",
					Field:    content.NameField,
					Value:    "name-sv",
				},
				{
					Language: "sv-SE",
					Field:    content.UrlSegmentField,
					Value:    "url-sv",
				},
				{
					Language: "en-US",
					Field:    content.NameField,
					Value:    "name-en",
				},
				{
					Language: "en-US",
					Field:    content.UrlSegmentField,
					Value:    "url-en",
				},
			},
			expected: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Status: content.Draft,
						Properties: content.ContentLanguage{
							"sv-SE": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "name-sv",
								},
								content.UrlSegmentField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "url-sv",
								},
							},
							"en-US": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "name-en",
								},
								content.UrlSegmentField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "url-en",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "localized field with not configured language should return error",
			contentdef: &emptyContentDef,
			existing: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Status:  content.Draft,
						Created: time.Now().UTC(),
						Properties: content.ContentLanguage{
							"sv-SE": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "",
								},
							},
						},
					},
				},
			},
			updatefields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "nb-NO",
					Field:    content.NameField,
					Value:    "url-sv",
				},
			},
			expectedErr: content.ErrMissingLanguage,
		},

		{
			name:       "not localized field with not default language should return error",
			contentdef: &emptyContentDef,
			existing: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Status:  content.Draft,
						Created: time.Now().UTC(),
						Properties: content.ContentLanguage{
							"sv-SE": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: false,
									Value:     "",
								},
							},
							"en-US": content.ContentFields{
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: false,
									Value:     "",
								},
							},
						},
					},
				},
			},
			updatefields: []struct {
				Language string
				Field    string
				Value    interface{}
			}{
				{
					Language: "en-US",
					Field:    content.NameField,
					Value:    "url-sv",
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

			test.existing.ContentDefinitionID = contentdefinitionId
			contentId, err := contentRepo.CreateContent(context.Background(), test.existing)
			assert.NoError(t, err)

			for _, f := range test.updatefields {

				cmd := UpdateField{
					Id:       contentId,
					Version:  test.version,
					Field:    f.Field,
					Language: f.Language,
					Value:    f.Value,
				}

				handler := UpdateFieldHandler{
					ContentRepository:           contentRepo,
					ContentDefinitionRepository: cdRepo,
					Factory:                     content.Factory{Cfg: &cfg},
				}

				err := handler.Handle(context.Background(), cmd)
				if test.expectedErr != "" {
					assert.Equal(t, test.expectedErr, err.Error())
				} else {
					assert.NoError(t, err)
				}

			}
			actual, err := contentRepo.GetContent(context.Background(), contentId)
			assert.NoError(t, err)
			for version, contentver := range test.expected.Version {

				for lang, fields := range contentver.Properties {

					for field, value := range fields {
						assert.Equal(t, value.Value, actual.Version[version].Properties[lang][field].Value, value)
					}
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
		publishVer  int
		expectedErr string
		expected    content.Content
	}{
		{
			name:       "required field not set should return error",
			contentdef: &reqfieldContentDef,
			content: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Properties: content.ContentLanguage{
							"sv-SE": {},
						},
					},
				},
			},

			expectedErr: "required",
		},
		{
			name:       "required field set should return ok",
			contentdef: &reqfieldContentDef,
			content: content.Content{
				Version: map[int]content.ContentVersion{
					0: {
						Properties: content.ContentLanguage{
							"sv-SE": {
								content.NameField: content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
									Value:     "name sv",
								},
								"required_field": content.ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: false,
									Value:     "ok",
								},
							},
						},
					},
				},
			},
			expected: content.Content{
				PublishedVersion: 0,
				Version: map[int]content.ContentVersion{
					0: {},
				},
			},
		},
		// {
		// 	name:       "new version is published",
		// 	contentdef: &emptyContentDef,
		// 	content: content.Content{
		// 		PublishedVersion: 0,
		// 		Version: map[int]content.ContentVersion{
		// 			0: {
		// 				Properties: map[string]map[string]interface{}{
		// 					"sv-SE": {
		// 						content.NameField: "name sv",
		// 					},
		// 				},
		// 			},
		// 			1: {
		// 				Properties: map[string]map[string]interface{}{
		// 					"sv-SE": {
		// 						content.NameField: "name sv",
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	publishVer: 1,
		// 	expected: content.Content{
		// 		PublishedVersion: 1,
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
				Version:   test.publishVer,
			}

			handler := PublishContentHandler{
				ContentDefinitionRepository: cdRepo,
				ContentRepository:           contentRepo,
				SiteConfiguration:           cfg,
			}

			err = handler.Handle(context.Background(), cmd)
			if test.expectedErr != "" {
				if assert.Error(t, err) {
					assert.Equal(t, test.expectedErr, err.Error())
				}
			} else {
				assert.NoError(t, err)

			}
			actual, err := contentRepo.GetContent(context.Background(), id)
			assert.NoError(t, err)

			assert.Equal(t, test.expected.PublishedVersion, actual.PublishedVersion)

			if test.expectedErr != "" {
				assert.Equal(t, test.publishVer, actual.PublishedVersion)
			}
			for v, contentver := range test.expected.Version {

				for lang, fields := range contentver.Properties {
					for field, value := range fields {
						assert.Equal(t, value, actual.Version[v].Properties[lang][field], v, lang, field)
					}
				}
			}
		})
	}
}
