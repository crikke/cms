package contentdefinition

import (
	"errors"
	"time"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
	"github.com/google/uuid"
)

const (
	NameField       = "name"
	UrlSegmentField = "url"
)

// swagger:model ContentDefinition
type ContentDefinition struct {
	ID                  uuid.UUID `bson:"_id"`
	Name                string    `bson:"name,omitempty"`
	Description         string    `bson:"description,omitempty"`
	Created             time.Time
	Propertydefinitions map[string]PropertyDefinition
}

// swagger:model PropertyDefinition
type PropertyDefinition struct {
	ID uuid.UUID `bson:"id"`
	// Name        string    `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	Type        string `bson:"type,omitempty"`
	Localized   bool   `bson:"localized,omitempty"`
	// instead of using map[strin]validator.Validator, interface{} is used
	// this wont be a problem becuase they will be translated to validator.Validator in GetValidatorQueury
	Validators map[string]interface{} `bson:"validators,omitempty"`
}

type PropertyDefinitionFactory struct {
}

func (f PropertyDefinitionFactory) NewPropertyDefinition(typ, description string, localized bool) (PropertyDefinition, error) {

	pd := PropertyDefinition{
		ID:          uuid.New(),
		Description: description,
		Type:        typ,
		Localized:   localized,
		Validators: map[string]interface{}{
			validator.RuleRequired: validator.Required(false),
		},
	}

	switch typ {
	case PropertyTypeText:
		pd.Validators[validator.RuleRegex] = validator.Regex("")
		pd.Validators[validator.RuleRange] = validator.Range{}
	case PropertyTypeBool:
		break
	case PropertyTypeNumber:
		pd.Validators[validator.RuleRange] = validator.Range{}
	default:
		return PropertyDefinition{}, errors.New("propertydefinition type does not exist")
	}

	return pd, nil
}

const (
	PropertyTypeText   = "text"
	PropertyTypeNumber = "number"
	PropertyTypeBool   = "bool"
)

// todo this can be done better
var propertydefinitionTypes = map[string]struct{}{
	"text":   {},
	"number": {},
	"bool":   {},
}

const ErrPropertyAlreadyExists = "propertydefinition already exists on contentdefinition"

func NewContentDefinition(name, desc string) (ContentDefinition, error) {

	if name == "" {
		return ContentDefinition{}, errors.New("name required")
	}

	return ContentDefinition{
		Name:        name,
		Description: desc,
		Propertydefinitions: map[string]PropertyDefinition{
			NameField: {
				ID:        uuid.New(),
				Type:      "text",
				Localized: true,
				Validators: map[string]interface{}{
					"required": validator.Required(true),
				},
			},
		}}, nil
}

func (cd ContentDefinition) PropertyValid(field, lang string, value interface{}) error {

	pd, ok := cd.Propertydefinitions[field]

	if !ok {
		return errors.New("property does not exist")
	}

	if !pd.Localized && lang != "" {
		return errors.New("content.ErrUnlocalizedPropLocalizedValue")
	}
	return nil
}

func NewPropertyDefinition(contentDefinition *ContentDefinition, name, description, propertytype string) (PropertyDefinition, error) {

	pd := PropertyDefinition{
		ID:          uuid.New(),
		Description: description,
		Type:        propertytype,
		Validators:  make(map[string]interface{}),
	}

	if err := pd.Valid(); err != nil {
		return PropertyDefinition{}, err
	}

	if _, exist := contentDefinition.Propertydefinitions[name]; exist {
		return PropertyDefinition{}, errors.New(ErrPropertyAlreadyExists)
	}

	contentDefinition.Propertydefinitions[name] = pd
	return pd, nil
}

// Checks if PropertyDefinition is valid.
func (p PropertyDefinition) Valid() error {

	if _, ok := propertydefinitionTypes[p.Type]; !ok {
		return errors.New("invalid property definition type")
	}

	return nil
}
