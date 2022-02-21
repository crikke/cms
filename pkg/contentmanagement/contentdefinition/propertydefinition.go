package contentdefinition

import (
	"errors"

	"github.com/google/uuid"
)

type PropertyDefinition struct {
	ID          uuid.UUID `bson:"id"`
	Name        string    `bson:"name,omitempty"`
	Description string    `bson:"description,omitempty"`
	Type        string    `bson:"type,omitempty"`
	Localized   bool      `bson:"localized,omitempty"`
}

var propertydefinitionTypes = map[string]struct{}{
	"text":        {},
	"shortstring": {},
	"number":      {},
	"bool":        {},
}

func NewPropertyDefinition(cid uuid.UUID, name, description, propertytype string) (PropertyDefinition, error) {

	pd := PropertyDefinition{
		Name:        name,
		Description: description,
		Type:        propertytype,
	}

	if err := pd.Valid(); err != nil {
		return PropertyDefinition{}, err
	}

	return pd, nil
}

func (p PropertyDefinition) Valid() error {

	if p.Name == "" {
		return errors.New("name is empty")
	}

	if _, ok := propertydefinitionTypes[p.Type]; !ok {
		return errors.New("invalid property definition type")
	}

	return nil
}
