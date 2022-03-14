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

// func Test_UpdateField(t *testing.T) {

// 	tests := []struct {
// 		name      string
// 		content   Content
// 		expect    Content
// 		expectErr string
// 		update    struct {
// 			version   int
// 			fieldname string
// 			lang      string
// 			value     interface{}
// 		}
// 	}{
// 		// {
// 		// 	name: "update valid field",
// 		// 	content: Content{
// 		// 		ID:     uuid.New(),
// 		// 		Data: ContentData{
// 		// 				Status: Draft,
// 		// 				Properties: ContentLanguage{
// 		// 					"sv-SE": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: false,
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	expect: Content{
// 		// 		Data: ontentData{
// 		// 			Status: Draft,
// 		// 			Properties: ContentLanguage{
// 		// 				"sv-SE": ContentFields{
// 		// 					"foo": ContentField{
// 		// 						Type:      "text",
// 		// 						Localized: false,
// 		// 						Value:     "bar",
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	update: struct {
// 		// 		version   int
// 		// 		fieldname string
// 		// 		lang      string
// 		// 		value     interface{}
// 		// 	}{
// 		// 		version:   0,
// 		// 		fieldname: "foo",
// 		// 		lang:      "sv-SE",
// 		// 		value:     "bar",
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "update localized field no language set",
// 		// 	content: Content{
// 		// 		ID:     uuid.New(),
// 		// 		Data:   0,
// 		// 		Status: Draft,
// 		// 		Version: map[int]ContentData{
// 		// 			0: {
// 		// 				Status: Draft,
// 		// 				Properties: ContentLanguage{
// 		// 					"sv-SE": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: true,
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	expect: Content{
// 		// 		Data: ContentData{
// 		// 				Status: Draft,
// 		// 				Properties: ContentLanguage{
// 		// 					"sv-SE": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							Type:      "text",
// 		// 							Localized: false,
// 		// 							Value:     "bar",
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 		},
// 		// 	},
// 		// 	update: struct {
// 		// 		version   int
// 		// 		fieldname string
// 		// 		lang      string
// 		// 		value     interface{}
// 		// 	}{
// 		// 		version:   0,
// 		// 		fieldname: "foo",
// 		// 		lang:      "",
// 		// 		value:     "bar",
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "update not draft should return error",
// 		// 	content: Content{
// 		// 		ID:     uuid.New(),
// 		// 		Data: ContentData{
// 		// 				Status: Published,
// 		// 				Properties: ContentLanguage{
// 		// 					"nb-NO": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: true,
// 		// 						},
// 		// 					},
// 		// 					"sv-SE": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: true,
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 		},
// 		// 	},
// 		// 	expectErr: ErrNotDraft,
// 		// 	update: struct {
// 		// 		version   int
// 		// 		fieldname string
// 		// 		lang      string
// 		// 		value     interface{}
// 		// 	}{
// 		// 		version:   0,
// 		// 		fieldname: "foo",
// 		// 		lang:      "nb-NO",
// 		// 		value:     "bar",
// 		// 	},
// 		// },

// 		// {
// 		// 	name: "update unlocalized field not default language return error",
// 		// 	content: Content{
// 		// 		ID:     uuid.New(),
// 		// 		Data: ContentData{
// 		// 				Status: Draft,
// 		// 				Properties: ContentLanguage{
// 		// 					"nb-NO": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: false,
// 		// 						},
// 		// 					},
// 		// 					"sv-SE": ContentFields{
// 		// 						"foo": ContentField{
// 		// 							ID:        uuid.New(),
// 		// 							Type:      "text",
// 		// 							Localized: false,
// 		// 						},
// 		// 					},
// 		// 				},
// 		// 		},
// 		// 	},
// 		// 	expectErr: ErrUnlocalizedPropLocalizedValue,
// 		// 	update: struct {
// 		// 		version   int
// 		// 		fieldname string
// 		// 		lang      string
// 		// 		value     interface{}
// 		// 	}{
// 		// 		version:   0,
// 		// 		fieldname: "foo",
// 		// 		lang:      "nb-NO",
// 		// 		value:     "bar",
// 		// 	},
// 		// },
// 	}

// 	factory := Factory{Cfg: &siteconfiguration.SiteConfiguration{
// 		Languages: []language.Tag{
// 			language.MustParse("sv-SE"),
// 		},
// 	}}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {

// 			cv := test.content.Version[test.update.version]
// 			err := factory.SetField(&cv, test.update.lang, test.update.fieldname, test.update.value)

// 			if test.expectErr != "" {
// 				assert.Equal(t, test.expectErr, err.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			for ver, ecv := range test.expect.Version {
// 				for lang, ecl := range ecv.Properties {
// 					for field, val := range ecl {

// 						assert.Equal(t, val.Type, test.content.Version[ver].Properties[lang][field].Type)
// 						assert.Equal(t, val.Value, test.content.Version[ver].Properties[lang][field].Value)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// // todo test property ID
// func Test_CreateNewVersion(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		contentDef contentdefinition.ContentDefinition
// 		content    Content
// 		expect     Content
// 		updateFn   func(cd contentdefinition.ContentDefinition) contentdefinition.ContentDefinition
// 	}{
// 		{
// 			name: "rename field",
// 			contentDef: contentdefinition.ContentDefinition{
// 				ID: uuid.New(),
// 				Propertydefinitions: map[string]contentdefinition.PropertyDefinition{
// 					"field1": {
// 						ID:   uuid.New(),
// 						Type: "bool",
// 					},
// 				},
// 			},
// 			content: Content{
// 				Version: map[int]ContentData{
// 					0: {
// 						Properties: ContentLanguage{
// 							"sv-SE": ContentFields{
// 								"field1": ContentField{
// 									Type: "bool",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			updateFn: func(cd contentdefinition.ContentDefinition) contentdefinition.ContentDefinition {

// 				cd.Propertydefinitions["field_new"] = cd.Propertydefinitions["field1"]
// 				return cd
// 			},
// 			expect: Content{
// 				Version: map[int]ContentData{
// 					1: {
// 						Properties: ContentLanguage{
// 							"sv-SE": ContentFields{
// 								"field_new": ContentField{
// 									Type: "bool",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	f := Factory{
// 		Cfg: &siteconfiguration.SiteConfiguration{Languages: []language.Tag{
// 			language.MustParse("sv-SE"),
// 		}},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {

// 			cd := test.updateFn(test.contentDef)
// 			_, err := f.NewContentVersion(&test.content, cd, 0)
// 			assert.NoError(t, err)

// 			// assert
// 			for ver, ecv := range test.expect.Version {
// 				for lan, ecl := range ecv.Properties {
// 					for name, ecf := range ecl {

// 						acv, ok := test.content.Version[ver]
// 						assert.True(t, ok)
// 						acl, ok := acv.Properties[lan]
// 						assert.Equal(t, Draft, acv.Status)
// 						assert.True(t, ok)
// 						acf, ok := acl[name]
// 						assert.True(t, ok)

// 						assert.Equal(t, ecf.Localized, acf.Localized)
// 						assert.Equal(t, ecf.Type, acf.Type)
// 						// assert.Equal(t, ecf.ID, acf.ID)
// 					}
// 				}
// 			}
// 		})
// 	}
// }
