//go:build unit

package contentdefinition

import (
	"testing"

	"github.com/crikke/cms/pkg/contentdefinition/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewPropertyDefinition(t *testing.T) {
	tests := []struct {
		name        string
		contentDef  ContentDefinition
		typ         string
		desc        string
		localized   bool
		expect      PropertyDefinition
		expectedErr string
	}{
		{
			name: "property ok",
			contentDef: ContentDefinition{
				Name: "test",
			},
			typ:       PropertyTypeBool,
			localized: false,
			desc:      "desc",
			expect: PropertyDefinition{
				Description: "desc",
				Type:        PropertyTypeBool,
				Localized:   false,
				Validators: map[string]interface{}{
					validator.RuleRequired: validator.Required(false),
				},
			},
		},
		{
			name: "prop already exist",
			contentDef: ContentDefinition{
				Name: "test",
				Propertydefinitions: map[string]PropertyDefinition{
					"test": {},
				},
			},
			expectedErr: ErrPropertyAlreadyExists,
			expect:      PropertyDefinition{},
		},
		{
			name: "proptype not exists",
			contentDef: ContentDefinition{
				Name: "test",
			},
			typ:         "some random string",
			expectedErr: ErrPropertyTypeNotExists,
			expect:      PropertyDefinition{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := ContentDefinitionFactory{}

			err := f.NewPropertyDefinition(&test.contentDef, "test", test.typ, test.desc, test.localized)

			if test.expectedErr != "" && assert.Error(t, err) {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expect.Description, test.contentDef.Propertydefinitions["test"].Description)
			assert.Equal(t, test.expect.Type, test.contentDef.Propertydefinitions["test"].Type)
			assert.Equal(t, test.expect.Localized, test.contentDef.Propertydefinitions["test"].Localized)
			assert.Equal(t, test.expect.Validators, test.contentDef.Propertydefinitions["test"].Validators)
		})
	}
}

func Test_UpdatePropertyDefinitionName(t *testing.T) {
	tests := []struct {
		name        string
		contentDef  ContentDefinition
		propName    string
		expect      ContentDefinition
		expectedErr string
	}{
		{
			name:     "name ok",
			propName: "new",
			contentDef: ContentDefinition{
				Name: "test contentdef",
				Propertydefinitions: map[string]PropertyDefinition{
					"old": {
						ID:         uuid.MustParse("5cb686d2-bc1c-4261-80c9-e539222bda73"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
				},
			},
			expect: ContentDefinition{
				Name: "test contentdef",
				Propertydefinitions: map[string]PropertyDefinition{
					"new": {
						ID:         uuid.MustParse("5cb686d2-bc1c-4261-80c9-e539222bda73"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
				},
			},
		},
		{
			name:     "name already exist",
			propName: "existing",
			contentDef: ContentDefinition{
				Name: "test contentdef",
				Propertydefinitions: map[string]PropertyDefinition{
					"old": {
						ID:         uuid.MustParse("5cb686d2-bc1c-4261-80c9-e539222bda73"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
					"existing": {
						ID:         uuid.MustParse("ae4b4e24-d3e5-4efa-83a5-f9d6eeadabe9"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
				},
			},
			expect: ContentDefinition{
				Name: "test contentdef",
				Propertydefinitions: map[string]PropertyDefinition{
					"old": {
						ID:         uuid.MustParse("5cb686d2-bc1c-4261-80c9-e539222bda73"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
					"existing": {
						ID:         uuid.MustParse("ae4b4e24-d3e5-4efa-83a5-f9d6eeadabe9"),
						Type:       "text",
						Validators: make(map[string]interface{}),
					},
				},
			},
			expectedErr: "property with name already exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := ContentDefinitionFactory{}

			err := f.UpdatePropertyDefinitionName(&test.contentDef, uuid.MustParse("5cb686d2-bc1c-4261-80c9-e539222bda73"), test.propName)

			if test.expectedErr != "" && assert.Error(t, err) {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expect, test.contentDef)
		})
	}
}

func Test_UpdatePropertyDefinition(t *testing.T) {
	tests := []struct {
		name        string
		contentDef  ContentDefinition
		expect      PropertyDefinition
		expectedErr string

		id         uuid.UUID
		desc       string
		localized  bool
		validators map[string]interface{}
	}{
		{
			name:      "update property ok",
			id:        uuid.MustParse("ae4b4e24-d3e5-4efa-83a5-f9d6eeadabe9"),
			desc:      "new",
			localized: true,
			validators: map[string]interface{}{
				validator.RuleRequired: true,
			},
			contentDef: ContentDefinition{
				Propertydefinitions: map[string]PropertyDefinition{
					"prop": {
						ID:          uuid.MustParse("ae4b4e24-d3e5-4efa-83a5-f9d6eeadabe9"),
						Description: "",
						Type:        "text",
						Localized:   false,
						Validators: map[string]interface{}{
							validator.RuleRequired: validator.Required(false),
						},
					},
				},
			},
			expect: PropertyDefinition{
				ID:          uuid.MustParse("ae4b4e24-d3e5-4efa-83a5-f9d6eeadabe9"),
				Description: "new",
				Type:        "text",
				Localized:   true,
				Validators: map[string]interface{}{
					validator.RuleRequired: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := ContentDefinitionFactory{}

			err := f.UpdatePropertyDefinition(&test.contentDef, test.id, test.desc, test.localized, test.validators)

			if test.expectedErr != "" && assert.Error(t, err) {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expect, test.contentDef.Propertydefinitions["prop"])
		})
	}
}
