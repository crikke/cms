package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestMarshalContent(t *testing.T) {
	c := Content{
		ID: ContentReference{
			Version: 0,
			Locale:  &language.Swedish,
		},
		URLSegment: "",
		Name:       "Root",
		Created:    time.Now(),
		Updated:    time.Now(),
		Properties: []Property{
			{
				ID:        uuid.UUID{},
				Name:      "header",
				Type:      "text",
				Localized: false,
				Value:     "Hello world",
			},
		},
	}

	b, err := json.Marshal(c)
	assert.NoError(t, err)
	assert.Contains(t, string(b), "Hello world")
}