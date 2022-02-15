package contentdefinition

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ContentDefinition struct {
	ID          uuid.UUID
	Name        string
	Description string
	Created     time.Time
}

func NewContentDefinition(name, desc string) (ContentDefinition, error) {

	if name == "" {
		return ContentDefinition{}, errors.New("missing field: Name")
	}

	return ContentDefinition{Name: name, Description: desc}, nil
}
