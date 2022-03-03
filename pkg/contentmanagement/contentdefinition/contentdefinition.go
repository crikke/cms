package contentdefinition

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ContentDefinition struct {
	ID                  uuid.UUID `bson:"_id"`
	Name                string    `bson:"name,omitempty"`
	Description         string    `bson:"description,omitempty"`
	Created             time.Time
	Propertydefinitions map[string]PropertyDefinition
}

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

var propertydefinitionTypes = map[string]struct{}{
	"text":        {},
	"shortstring": {},
	"number":      {},
	"bool":        {},
}

const ErrPropertyAlreadyExists = "propertydefinition already exists on contentdefinition"

func NewContentDefinition(name, desc string) (ContentDefinition, error) {

	if name == "" {
		return ContentDefinition{}, errors.New("missing field: Name")
	}

	return ContentDefinition{Name: name, Description: desc}, nil
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
