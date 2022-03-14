//go:build unit

package contentdefinition

import (
	"testing"

	"github.com/crikke/cms/pkg/contentdefinition/validator"
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
