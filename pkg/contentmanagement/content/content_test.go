package content

import (
	"testing"
	"time"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_CreateContent(t *testing.T) {

	tests := []struct {
		name       string
		contentdef contentdefinition.ContentDefinition
		expect     Content
	}{
		{
			name: "create content ok",
			contentdef: contentdefinition.ContentDefinition{
				ID:          uuid.New(),
				Name:        "test",
				Description: "test desc",
				Created:     time.Now(),
				Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
					"field1": {
						ID:          uuid.New(),
						Type:        "text",
						Localized:   false,
						Description: "text field",
					},
					"field2": {
						ID:          uuid.New(),
						Type:        "text",
						Localized:   true,
						Description: "text field",
					},
				},
			},
			expect: Content{
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"field1": ContentField{
									Type:      "text",
									Localized: false,
								},
								"field2": ContentField{
									Type:      "text",
									Localized: true,
								},
							},
						},
					},
				},
				Status: Draft,
			},
		},
	}

	cfg := &siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Factory{
				Cfg: cfg,
			}

			actual, err := f.NewContent(test.contentdef)
			assert.NoError(t, err)

			assert.Equal(t, test.contentdef.ID, actual.ContentDefinitionID)
			assert.Equal(t, test.expect.Status, actual.Status)

			for version, cv := range test.expect.Version {

				acv, ok := actual.Version[version]
				assert.True(t, ok)

				assert.Equal(t, cv.Status, acv.Status)

				for lang, cl := range cv.Properties {
					acl, ok := acv.Properties[lang]
					assert.True(t, ok)

					for name, value := range cl {

						aval, ok := acl[name]
						assert.True(t, ok)

						assert.Equal(t, test.contentdef.Propertydefinitions[name].ID, aval.ID)
						assert.Equal(t, value.Localized, aval.Localized)
						assert.Equal(t, value.Type, aval.Type)
					}
				}
			}

		})
	}
}

func Test_UpdateField(t *testing.T) {

	tests := []struct {
		name      string
		content   Content
		expect    Content
		expectErr string
		update    struct {
			version   int
			fieldname string
			lang      string
			value     interface{}
		}
	}{
		{
			name: "update valid field",
			content: Content{
				ID:               uuid.New(),
				PublishedVersion: 0,
				Status:           Draft,
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: false,
								},
							},
						},
					},
				},
			},
			expect: Content{
				PublishedVersion: 0,
				Status:           Draft,
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"foo": ContentField{
									Type:      "text",
									Localized: false,
									Value:     "bar",
								},
							},
						},
					},
				},
			},
			update: struct {
				version   int
				fieldname string
				lang      string
				value     interface{}
			}{
				version:   0,
				fieldname: "foo",
				lang:      "sv-SE",
				value:     "bar",
			},
		},
		{
			name: "update localized field no language set",
			content: Content{
				ID:               uuid.New(),
				PublishedVersion: 0,
				Status:           Draft,
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
								},
							},
						},
					},
				},
			},
			expect: Content{
				PublishedVersion: 0,
				Status:           Draft,
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"foo": ContentField{
									Type:      "text",
									Localized: false,
									Value:     "bar",
								},
							},
						},
					},
				},
			},
			update: struct {
				version   int
				fieldname string
				lang      string
				value     interface{}
			}{
				version:   0,
				fieldname: "foo",
				lang:      "",
				value:     "bar",
			},
		},
		{
			name: "update not draft should return error",
			content: Content{
				ID:               uuid.New(),
				PublishedVersion: 0,
				Status:           Published,
				Version: map[int]ContentVersion{
					0: {
						Status: Published,
						Properties: ContentLanguage{
							"nb-NO": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
								},
							},
							"sv-SE": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
								},
							},
						},
					},
				},
			},
			expectErr: ErrNotDraft,
			update: struct {
				version   int
				fieldname string
				lang      string
				value     interface{}
			}{
				version:   0,
				fieldname: "foo",
				lang:      "nb-NO",
				value:     "bar",
			},
		},

		{
			name: "update unlocalized field not default language return error",
			content: Content{
				ID:               uuid.New(),
				PublishedVersion: 0,
				Status:           Draft,
				Version: map[int]ContentVersion{
					0: {
						Status: Draft,
						Properties: ContentLanguage{
							"nb-NO": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
								},
							},
							"sv-SE": ContentFields{
								"foo": ContentField{
									ID:        uuid.New(),
									Type:      "text",
									Localized: true,
								},
							},
						},
					},
				},
			},
			expectErr: ErrUnlocalizedPropLocalizedValue,
			update: struct {
				version   int
				fieldname string
				lang      string
				value     interface{}
			}{
				version:   0,
				fieldname: "foo",
				lang:      "nb-NO",
				value:     "bar",
			},
		},
	}

	factory := Factory{Cfg: &siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			cv := test.content.Version[test.update.version]
			err := factory.SetField(&cv, test.update.lang, test.update.fieldname, test.update.value)

			if test.expectErr != "" {
				assert.Equal(t, test.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			for ver, ecv := range test.expect.Version {
				for lang, ecl := range ecv.Properties {
					for field, val := range ecl {

						assert.Equal(t, val.Type, test.content.Version[ver].Properties[lang][field].Type)
						assert.Equal(t, val.Value, test.content.Version[ver].Properties[lang][field].Value)
					}
				}
			}
		})
	}
}

// todo test property ID
func Test_CreateNewVersion(t *testing.T) {
	tests := []struct {
		name       string
		contentDef contentdefinition.ContentDefinition
		content    Content
		expect     Content
		updateFn   func(cd contentdefinition.ContentDefinition) contentdefinition.ContentDefinition
	}{
		{
			name: "rename field",
			contentDef: contentdefinition.ContentDefinition{
				ID: uuid.New(),
				Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
					"field1": {
						ID:   uuid.New(),
						Type: "bool",
					},
				},
			},
			content: Content{
				Version: map[int]ContentVersion{
					0: {
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"field1": ContentField{
									Type: "bool",
								},
							},
						},
					},
				},
			},
			updateFn: func(cd contentdefinition.ContentDefinition) contentdefinition.ContentDefinition {

				cd.Propertydefinitions["field_new"] = cd.Propertydefinitions["field1"]
				return cd
			},
			expect: Content{
				Version: map[int]ContentVersion{
					1: {
						Properties: ContentLanguage{
							"sv-SE": ContentFields{
								"field_new": ContentField{
									Type: "bool",
								},
							},
						},
					},
				},
			},
		},
	}

	f := Factory{
		Cfg: &siteconfiguration.SiteConfiguration{Languages: []language.Tag{
			language.MustParse("sv-SE"),
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			cd := test.updateFn(test.contentDef)
			_, err := f.NewContentVersion(&test.content, cd, 0)
			assert.NoError(t, err)

			// assert
			for ver, ecv := range test.expect.Version {
				for lan, ecl := range ecv.Properties {
					for name, ecf := range ecl {

						acv, ok := test.content.Version[ver]
						assert.True(t, ok)
						acl, ok := acv.Properties[lan]
						assert.Equal(t, Draft, acv.Status)
						assert.True(t, ok)
						acf, ok := acl[name]
						assert.True(t, ok)

						assert.Equal(t, ecf.Localized, acf.Localized)
						assert.Equal(t, ecf.Type, acf.Type)
						// assert.Equal(t, ecf.ID, acf.ID)
					}
				}
			}
		})
	}
}
