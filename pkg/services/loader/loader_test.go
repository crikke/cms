package loader

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLoadContentWithDefaultLanguage(t *testing.T) {

	cfg := config.SiteConfiguration{
		Languages: []language.Tag{
			language.Swedish,
		},
	}
	loader := NewLoader(newMockRepo(), cfg)

	content, err := loader.GetContent(context.TODO(), domain.ContentReference{})

	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	assert.Equal(t, "foo", content.Name)
	assert.Equal(t, "foo", content.URLSegment)
	assert.Equal(t, "prop", content.Properties[0].Name)
	assert.Equal(t, "text", content.Properties[0].Type)
	assert.Equal(t, "bar", content.Properties[0].Value)

}

type mockRepo struct {
	content []repository.ContentData
}

func newMockRepo() mockRepo {

	repo := mockRepo{}

	cd := repository.ContentData{
		ID:               uuid.UUID{},
		PublishedVersion: 0,
		Data: map[int]repository.ContentVersion{
			0: {
				Name: map[string]string{
					"sv": "foo",
				},
				URLSegment: map[string]string{
					"sv": "foo",
				},
				Properties: []repository.ContentProperty{
					{
						ID:        uuid.UUID{},
						Name:      "prop",
						Type:      "text",
						Localized: false,
						Value: map[string]interface{}{
							"sv": "bar",
						},
					},
				},
			},
		},
	}
	repo.content = []repository.ContentData{cd}

	return repo
}
func (m mockRepo) GetContent(ctx context.Context, contentReference domain.ContentReference) (repository.ContentData, error) {
	return m.content[0], nil
}

func (m mockRepo) GetChildren(ctx context.Context, contentReference domain.ContentReference) ([]repository.ContentData, error) {
	return nil, nil
}
