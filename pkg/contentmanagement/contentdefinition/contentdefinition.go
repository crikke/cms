package contentdefinition

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ContentDefinition struct {
	ID          uuid.UUID `bson:"_id"`
	Name        string    `bson:"name,omitempty"`
	Description string    `bson:"description,omitempty"`
	Created     time.Time
	// todo: ensure name is unique
	Propertydefinitions map[string]PropertyDefinition
}

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
