package contentdefinition

import (
	"time"

	"github.com/google/uuid"
)

type contentdefinitionQuerier struct {
}

type ContentDefinition struct {
	ID      uuid.UUID
	Name    string
	Created time.Time
	Updated time.Time

	Properties []PropertyDefinition
}

type PropertyDefinition struct {
	ID           uuid.UUID
	Name         string
	PropertyType string
}
