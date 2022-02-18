package propertydefinition

import (
	"errors"

	"github.com/google/uuid"
)

type PropertyDefinition struct {
	ID          uuid.UUID
	Type        string
	Name        string
	Description string
	Localized   bool
}

var propertydefinitionTypes = map[string]struct{}{
	"text":        {},
	"shortstring": {},
	"number":      {},
	"bool":        {},
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
