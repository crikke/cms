//go:build unit

package content

import (
	"testing"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewContentVersion(t *testing.T) {

	tests := []struct {
		name       string
		existing   Content
		expect     Content
		contentdef contentdefinition.ContentDefinition
	}{
		{
			name: "not default locale assert no unlocalized properties ",
			contentdef: contentdefinition.ContentDefinition{
				ID:   uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Name: "test contentdef",
				Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
					"localized": {
						ID:        uuid.MustParse("6973eba3-24b1-44f3-ade1-83a5e3de5d1b"),
						Type:      "string",
						Localized: true,
					},
					"unlocalized": {
						ID:        uuid.MustParse("dfddadc9-0aaa-48e4-8465-43a39559d94d"),
						Type:      "string",
						Localized: false,
					},
				},
			},
			existing: Content{
				ID:                  uuid.New(),
				ContentDefinitionID: uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Data: ContentData{
					Version: 0,
					Status:  Draft,
					Properties: ContentLanguage{
						"defaultlang": ContentFields{
							"localized": ContentField{
								ID:    uuid.MustParse("6973eba3-24b1-44f3-ade1-83a5e3de5d1b"),
								Value: "localized value default locale",
							},
							"unlocalized": ContentField{
								ID:    uuid.MustParse("dfddadc9-0aaa-48e4-8465-43a39559d94d"),
								Value: "unlocalized value default locale",
							},
						},
						"other": ContentFields{
							"localized": ContentField{
								ID:    uuid.MustParse("6973eba3-24b1-44f3-ade1-83a5e3de5d1b"),
								Value: "localized value other locale",
							},
						},
					},
				},
			},
			expect: Content{
				ID:                  uuid.New(),
				ContentDefinitionID: uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Data: ContentData{
					Version: 0,
					Status:  Draft,
					Properties: ContentLanguage{
						"defaultlang": ContentFields{
							"localized": ContentField{
								ID:        uuid.MustParse("6973eba3-24b1-44f3-ade1-83a5e3de5d1b"),
								Value:     "localized value default locale",
								Type:      "string",
								Localized: true,
							},
							"unlocalized": ContentField{
								ID:        uuid.MustParse("dfddadc9-0aaa-48e4-8465-43a39559d94d"),
								Value:     "unlocalized value default locale",
								Type:      "string",
								Localized: false},
						},
						"other": ContentFields{
							"localized": ContentField{
								ID:        uuid.MustParse("6973eba3-24b1-44f3-ade1-83a5e3de5d1b"),
								Value:     "localized value other locale",
								Type:      "string",
								Localized: true},
						},
					},
				},
			},
		},
		{
			name: "switch name on properties",
			contentdef: contentdefinition.ContentDefinition{
				ID:   uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Name: "test contentdef",
				// IDs are switched, so
				Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
					"a": {
						ID:        uuid.MustParse("aaaaaaaa-0aaa-48e4-8465-43a39559d94d"),
						Type:      "string",
						Localized: false,
					},
					"b": {
						ID:        uuid.MustParse("bbbbbbbb-24b1-44f3-ade1-83a5e3de5d1b"),
						Type:      "string",
						Localized: false,
					},
				},
			},
			existing: Content{
				ID:                  uuid.New(),
				ContentDefinitionID: uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Data: ContentData{
					Version: 0,
					Status:  Draft,
					Properties: ContentLanguage{
						"defaultlang": ContentFields{
							"a": ContentField{
								ID:    uuid.MustParse("bbbbbbbb-24b1-44f3-ade1-83a5e3de5d1b"),
								Type:  "string",
								Value: "foo",
							},
							"b": ContentField{
								ID:    uuid.MustParse("aaaaaaaa-0aaa-48e4-8465-43a39559d94d"),
								Type:  "string",
								Value: "bar",
							},
						},
					},
				},
			},
			expect: Content{
				ID:                  uuid.New(),
				ContentDefinitionID: uuid.MustParse("d5d2ba13-7ef2-4ed3-b196-1d96ecca3bcb"),
				Data: ContentData{
					Version: 0,
					Status:  Draft,
					Properties: ContentLanguage{
						"defaultlang": ContentFields{
							"a": ContentField{
								ID:        uuid.MustParse("aaaaaaaa-0aaa-48e4-8465-43a39559d94d"),
								Value:     "bar",
								Type:      "string",
								Localized: false,
							},
							"b": ContentField{
								ID:        uuid.MustParse("bbbbbbbb-24b1-44f3-ade1-83a5e3de5d1b"),
								Value:     "foo",
								Type:      "string",
								Localized: false},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			f := ContentFactory{}

			contentData, err := f.NewContentVersion(test.existing, test.contentdef, test.existing.Data.Version, "defaultlang")
			assert.NoError(t, err)

			for lang, fields := range test.expect.Data.Properties {

				assert.Equal(t, fields, contentData.Properties[lang])
			}
		})
	}

}

func Test_SetField(t *testing.T) {

	tests := []struct {
		name      string
		content   ContentData
		expect    ContentData
		lang      string
		fieldname string
		value     interface{}
		expectErr string
	}{
		{
			name: "set field ok",
			content: ContentData{
				Status: Draft,
				Properties: ContentLanguage{
					"default": ContentFields{
						"field": ContentField{
							Value: "foo",
						},
					},
				},
			},
			expect: ContentData{
				Status: Draft,
				Properties: ContentLanguage{
					"default": ContentFields{
						"field": ContentField{
							Value: "bar",
						},
					},
				},
			},
			lang:      "default",
			fieldname: "field",
			value:     "bar",
		},
		{
			name: "missing locale",
			content: ContentData{
				Status:     Draft,
				Properties: ContentLanguage{},
			},
			expect: ContentData{
				Status:     Draft,
				Properties: ContentLanguage{},
			},
			lang:      "default",
			fieldname: "field",
			value:     "bar",
			expectErr: ErrMissingLanguage,
		},
		{
			name: "not draft",
			content: ContentData{
				Status:     Published,
				Properties: ContentLanguage{},
			},
			expect: ContentData{
				Status:     Published,
				Properties: ContentLanguage{},
			},
			lang:      "default",
			fieldname: "field",
			value:     "bar",
			expectErr: ErrNotDraft,
		},
		{
			name: "field not exists",
			content: ContentData{
				Status: Draft,
				Properties: ContentLanguage{
					"default": ContentFields{},
				},
			},
			expect: ContentData{
				Status: Draft,
				Properties: ContentLanguage{
					"default": ContentFields{},
				},
			},
			lang:      "default",
			fieldname: "field",
			value:     "bar",
			expectErr: ErrMissingField,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			f := ContentFactory{}

			err := f.SetField(&test.content, test.lang, test.fieldname, test.value)

			if test.expectErr != "" {
				if assert.Error(t, err) {
					assert.Equal(t, test.expectErr, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expect, test.content)
		})
	}
}
