package repository

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_GetDocument(t *testing.T) {

	tests := []struct {
		name         string
		uid          uuid.UUID
		version      int
		locale       string
		expectedName string
	}{
		{
			name:         "Test GetDefault",
			uid:          uuid.MustParse("a1f6da93-80c9-4315-a012-1ea4249d7413"),
			expectedName: "Hejsan!",
			locale:       "sv",
		},
		{
			name:         "Test GetDefault - en",
			uid:          uuid.MustParse("a1f6da93-80c9-4315-a012-1ea4249d7413"),
			expectedName: "Hello!",
			locale:       "en",
		},
	}

	cfg := config.LoadServerConfiguration()
	r, err := NewRepository(context.TODO(), cfg)
	assert.NoError(t, err)

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			cref := domain.ContentReference{ID: test.uid,
				Version: test.version,
			}
			l := language.MustParse(test.locale)

			cref.Locale = &l
			c, err := r.GetContent(context.Background(), cref)
			assert.NoError(t, err)

			assert.Equal(t, test.uid, c.ID)
			assert.Equal(t, test.expectedName, c.Data[test.version].Name[test.locale])
		})
	}
}

func Test_GetChildDocuments(t *testing.T) {
	tests := []struct {
		name        string
		uid         uuid.UUID
		version     int
		locale      string
		returnedIds []uuid.UUID
	}{
		{
			name:   "Return 1 child node",
			uid:    uuid.MustParse("a1f6da93-80c9-4315-a012-1ea4249d7413"),
			locale: "sv",
			returnedIds: []uuid.UUID{
				uuid.MustParse("b2184714-4bae-4c50-9642-98fc5cadab86"),
			},
		},
	}
	cfg := config.LoadServerConfiguration()
	r, err := NewRepository(context.TODO(), cfg)
	assert.NoError(t, err)

	for _, test := range tests {

		cref := domain.ContentReference{ID: test.uid,
			Version: test.version,
		}
		l := language.MustParse(test.locale)

		cref.Locale = &l

		t.Run(test.name, func(t *testing.T) {
			returned, err := r.GetChildren(context.TODO(), cref)
			assert.NoError(t, err)

			for _, data := range returned {

				match := false
				for i := 0; i < len(test.returnedIds); i++ {

					match = data.ID == test.returnedIds[i]

					if match {
						break
					}
				}

				assert.True(t, match)
			}
		})
	}
}
