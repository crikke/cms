package loader

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
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
	content []contentData
}

func newMockRepo() mockRepo {

	repo := mockRepo{}

	cd := contentData{
		ID:      uuid.UUID{},
		Version: 0,
		Data: map[int]contentVersion{
			0: {
				Name: map[language.Tag]string{
					language.Swedish: "foo",
				},
				URLSegment: map[language.Tag]string{
					language.Swedish: "foo",
				},
				Properties: []contentProperty{
					{
						ID:        uuid.UUID{},
						Name:      "prop",
						Type:      "text",
						Localized: false,
						Value: map[language.Tag]interface{}{
							language.Swedish: "bar",
						},
					},
				},
			},
		},
	}
	repo.content = []contentData{cd}

	return repo
}
func (m mockRepo) GetContent(ctx context.Context, contentReference domain.ContentReference) (contentData, error) {
	return m.content[0], nil
}
