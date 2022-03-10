//go:build integration

package command

import (
	"context"
	"testing"
	"time"

	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/contentdefinition/validator"
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
		Data: content.ContentData{
			Status:  content.Draft,
			Created: time.Now().UTC(),
		},
	}
)

func Test_CreateContent(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("contentversion").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	cdRepo := contentdefinition.NewContentDefinitionRepository(c)
	contentRepo := content.NewContentRepository(c)

	cid, err := cdRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test2",
		Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
			contentdefinition.NameField: {
				ID:          uuid.New(),
				Description: "Content name",
				Type:        "text",
				Localized:   true,
				Validators: map[string]interface{}{
					"required": validator.Required(true),
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
	actual, err := contentRepo.GetContent(context.Background(), contentId, 0)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, contentId)
	assert.Equal(t, content.Draft, actual.Data.Status)
}

// func Test_CreateContent_Empty_ContentDefinition(t *testing.T) {
// 	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
// 	assert.NoError(t, err)
// 	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
// 	c.Database("cms").Collection("content").Drop(context.Background())
// 	c.Database("cms").Collection("contentversion").Drop(context.Background())

// 	contentRepo := content.NewContentRepository(c)
// 	cmd := CreateContent{}
// 	handler := CreateContentHandler{
// 		ContentDefinitionRepository: contentdefinition.NewContentDefinitionRepository(c),
// 		ContentRepository:           contentRepo,
// 	}

// 	contentId, err := handler.Handle(context.Background(), cmd)

// 	assert.Error(t, err)
// 	assert.Equal(t, uuid.UUID{}, contentId)
// }

func Test_UpdateContent(t *testing.T) {

	cfg := siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("en-US"),
		},
	}

	tests := []struct {
		name        string
		contentdef  *contentdefinition.ContentDefinition
		existing    content.ContentData
		cmd         UpdateContentFields
		expectedErr string
		expected    content.Content
	}{
		{
			name:       "update fields should return ok",
			contentdef: &emptyContentDef,
			existing: content.ContentData{
				Version: 0,
				Status:  content.Draft,
				Properties: content.ContentLanguage{
					"en-US": content.ContentFields{
						contentdefinition.NameField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: true,
							Value:     "",
						},
						contentdefinition.UrlSegmentField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: true,
							Value:     "",
						},
					},
					"sv-SE": content.ContentFields{
						contentdefinition.NameField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: true,
							Value:     "",
						},
						contentdefinition.UrlSegmentField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: true,
							Value:     "",
						},
					},
				},
			},
			cmd: UpdateContentFields{
				Language: "sv-SE",
				Version:  0,
				Fields: map[string]interface{}{
					contentdefinition.NameField:       "name-sv",
					contentdefinition.UrlSegmentField: "url-sv",
				},
			},
			expected: content.Content{
				Data: content.ContentData{
					Status: content.Draft,
					Properties: content.ContentLanguage{
						"sv-SE": content.ContentFields{
							contentdefinition.NameField: content.ContentField{
								ID:        uuid.New(),
								Type:      "text",
								Localized: true,
								Value:     "name-sv",
							},
							contentdefinition.UrlSegmentField: content.ContentField{
								ID:        uuid.New(),
								Type:      "text",
								Localized: true,
								Value:     "url-sv",
							},
						},
					},
				},
			},
		},
		{
			name:       "localized field with not configured language should return error",
			contentdef: &emptyContentDef,
			existing: content.ContentData{
				Version: 0,
				Status:  content.Draft,
				Properties: content.ContentLanguage{
					"sv-SE": content.ContentFields{
						contentdefinition.NameField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: true,
							Value:     "",
						},
					},
				},
			},
			cmd: UpdateContentFields{
				Language: "nb-NO",
				Version:  0,
				Fields: map[string]interface{}{
					contentdefinition.NameField: "test error",
				},
			},
			expectedErr: content.ErrMissingLanguage,
		},

		{
			name:       "not localized field with not default language should return error",
			contentdef: &emptyContentDef,
			existing: content.ContentData{
				Status: content.Draft,
				Properties: content.ContentLanguage{
					"sv-SE": content.ContentFields{
						contentdefinition.NameField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: false,
							Value:     "",
						},
					},
					"en-US": content.ContentFields{
						contentdefinition.NameField: content.ContentField{
							ID:        uuid.New(),
							Type:      "text",
							Localized: false,
							Value:     "",
						},
					},
				},
			},
			cmd: UpdateContentFields{
				Language: "en-US",
				Version:  0,
				Fields: map[string]interface{}{
					contentdefinition.NameField: "url-sv",
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
			c.Database("cms").Collection("contentversion").Drop(context.Background())

			contentDefinitionRepo := contentdefinition.NewContentDefinitionRepository(c)
			contentRepo := content.NewContentRepository(c)

			contentdefinitionId, err := contentDefinitionRepo.CreateContentDefinition(context.Background(), test.contentdef)
			assert.NoError(t, err)

			factory := content.Factory{
				Cfg: &cfg,
			}

			newContent, err := factory.NewContent(*test.contentdef)
			assert.NoError(t, err)
			newContent.ContentDefinitionID = contentdefinitionId
			contentId, err := contentRepo.CreateContent(context.Background(), newContent)
			assert.NoError(t, err)
			test.existing.ContentID = contentId

			contentRepo.UpdateContentData(context.Background(), contentId, test.existing.Version, func(ctx context.Context, cd *content.ContentData) (*content.ContentData, error) {
				return &test.existing, nil
			})

			handler := UpdateContentFieldsHandler{
				ContentRepository:           contentRepo,
				ContentDefinitionRepository: contentDefinitionRepo,
				Factory:                     content.Factory{Cfg: &cfg},
			}

			test.cmd.ContentID = contentId
			err = handler.Handle(context.Background(), test.cmd)
			if test.expectedErr != "" {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			actual, err := contentRepo.GetContent(context.Background(), contentId, test.cmd.Version)
			assert.NoError(t, err)

			for lang, fields := range test.expected.Data.Properties {

				for field, value := range fields {
					assert.Equal(t, value.Value, actual.Data.Properties[lang][field].Value, value)
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
		name            string
		contentdef      *contentdefinition.ContentDefinition
		contentVersions []content.ContentData
		publishVer      int
		expectedErr     string
		expected        content.Content
	}{
		{
			name:       "required field not set should return error",
			contentdef: &reqfieldContentDef,
			contentVersions: []content.ContentData{
				{
					Version: 0,
					Status:  content.Draft,
					Properties: content.ContentLanguage{
						"sv-SE": {},
					},
				},
			},
			expectedErr: "required",
			expected: content.Content{
				Data: content.ContentData{
					Status: content.Draft,
				},
			},
		},
		{
			name:       "required field set should return ok",
			contentdef: &reqfieldContentDef,
			contentVersions: []content.ContentData{
				{
					Version: 0,
					Status:  content.Draft,
					Properties: content.ContentLanguage{
						"sv-SE": {
							contentdefinition.NameField: content.ContentField{
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
			expected: content.Content{
				Data: content.ContentData{
					Status: content.Published,
				},
			},
		},
		{
			name:       "new version is published",
			contentdef: &emptyContentDef,
			contentVersions: []content.ContentData{
				{
					Version: 0,
					Status:  content.Published,
					Properties: content.ContentLanguage{
						"sv-SE": {
							contentdefinition.NameField: content.ContentField{
								Type:      "text",
								Localized: false,
								Value:     "version 1",
							},
						},
					},
				},
				{
					Status:  content.Draft,
					Version: 1,
					Properties: content.ContentLanguage{
						"sv-SE": {
							contentdefinition.NameField: content.ContentField{
								Type:      "text",
								Localized: false,
								Value:     "version 2",
							},
						},
					},
				},
			},
			publishVer: 1,
			expected: content.Content{
				Data: content.ContentData{
					Version: 1,
					Status:  content.Published,
				},
			},
		},
	}

	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
			c.Database("cms").Collection("content").Drop(context.Background())
			c.Database("cms").Collection("contentversion").Drop(context.Background())

			cdRepo := contentdefinition.NewContentDefinitionRepository(c)
			_, err := cdRepo.CreateContentDefinition(context.Background(), test.contentdef)
			assert.NoError(t, err)

			contentRepo := content.NewContentRepository(c)

			factory := content.Factory{Cfg: cfg}
			newContent, err := factory.NewContent(*test.contentdef)

			id, err := contentRepo.CreateContent(context.Background(), newContent)
			assert.NoError(t, err)

			for _, v := range test.contentVersions {
				v.ContentID = id

				err := contentRepo.UpdateContentData(context.Background(), id, 0, func(ctx context.Context, cd *content.ContentData) (*content.ContentData, error) {
					return &v, nil
				})
				assert.NoError(t, err)
			}

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

			actual, err := contentRepo.GetContent(context.Background(), id, test.publishVer)
			assert.NoError(t, err)
			assert.Equal(t, test.expected.Data.Status, actual.Data.Status)

			for lang, fields := range test.expected.Data.Properties {
				for field, value := range fields {
					assert.Equal(t, value, actual.Data.Properties[lang][field], lang, field)
				}
			}

			versions, err := contentRepo.ListContentVersions(context.Background(), id)
			assert.NoError(t, err)

			for _, cv := range versions {
				if cv.Version == actual.Data.Version {
					assert.Equal(t, actual.Data.Status, cv.Status)
				} else {
					assert.NotEqual(t, content.Published, cv.Status)
				}
			}
		})
	}
}
