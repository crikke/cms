package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/db"
	"github.com/stretchr/testify/assert"
)

func Test_GetContent(t *testing.T) {
	tests := []struct {
		name      string
		content   content.Content
		query     GetContent
		expect    ContentReadModel
		expectErr string
	}{
		{
			name: "get latest version",
			content: content.Content{
				PublishedVersion: 1,
				Status:           content.Published,
				Version: map[int]content.ContentVersion{
					1: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "baz",
							},
						},
					},
				},
			},
			query: GetContent{},
			expect: ContentReadModel{
				Status: content.Published,
				Properties: map[string]map[string]interface{}{
					"foo": {
						"bar": "baz",
					},
				},
			},
		},
		{
			name: "get previous version",
			content: content.Content{
				PublishedVersion: 3,
				Status:           content.Published,
				Version: map[int]content.ContentVersion{
					1: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "1",
							},
						},
					},
					2: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "2",
							},
						},
					},
					3: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "3",
							},
						},
					},
				},
			},
			query: GetContent{Version: makeInt(1)},
			expect: ContentReadModel{
				Status: content.Published,
				Properties: map[string]map[string]interface{}{
					"foo": {
						"bar": "1",
					},
				},
			},
		},
		{
			name: "version not exist",
			content: content.Content{
				PublishedVersion: 1,
				Status:           content.Published,
				Version: map[int]content.ContentVersion{
					1: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "1",
							},
						},
					},
				},
			},
			query:     GetContent{Version: makeInt(2)},
			expect:    ContentReadModel{},
			expectErr: content.ErrVersionNotExists,
		},
		{
			name: "negative version",
			content: content.Content{
				PublishedVersion: 1,
				Status:           content.Published,
				Version: map[int]content.ContentVersion{
					1: {
						Properties: map[string]map[string]interface{}{
							"foo": {
								"bar": "1",
							},
						},
					},
				},
			},
			query:     GetContent{Version: makeInt(-2)},
			expect:    ContentReadModel{},
			expectErr: content.ErrVersionNotExists,
		},
	}

	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
			c.Database("cms").Collection("content").Drop(context.Background())

			contentRepo := content.NewContentRepository(c)

			id, err := contentRepo.CreateContent(context.Background(), test.content)
			assert.NoError(t, err)

			test.query.Id = id

			handler := GetContentHandler{Repo: contentRepo}
			actual, err := handler.Handle(context.Background(), test.query)

			if test.expectErr != "" {
				assert.Equal(t, test.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expect.Status, actual.Status)
			eq := reflect.DeepEqual(test.expect.Properties, actual.Properties)
			assert.True(t, eq)
		})
	}
}

func makeInt(n int) *int {
	i := int(n)
	return &i
}
