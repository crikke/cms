package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
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
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "baz",
								},
							},
						},
					},
				},
			},
			query: GetContent{},
			expect: ContentReadModel{
				Status: content.Published,
				Properties: content.ContentLanguage{
					"foo": {
						"bar": content.ContentField{
							Value: "baz",
						},
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
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "1",
								},
							},
						},
					},
					2: {
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "2",
								},
							},
						},
					},
					3: {
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "3",
								},
							},
						},
					},
				},
			},
			query: GetContent{Version: makeInt(1)},
			expect: ContentReadModel{
				Status: content.Published,
				Properties: content.ContentLanguage{
					"foo": {
						"bar": content.ContentField{
							Value: "1",
						},
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
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "2",
								},
							},
						},
					},
				},
			},
			query:     GetContent{Version: makeInt(2)},
			expect:    ContentReadModel{},
			expectErr: content.ErrMissingVersion,
		},
		{
			name: "negative version",
			content: content.Content{
				PublishedVersion: 1,
				Status:           content.Published,
				Version: map[int]content.ContentVersion{
					1: {
						Properties: content.ContentLanguage{
							"foo": {
								"bar": content.ContentField{
									Value: "3",
								},
							},
						},
					},
				},
			},
			query:     GetContent{Version: makeInt(-2)},
			expect:    ContentReadModel{},
			expectErr: content.ErrMissingVersion,
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

func Test_ListChildContent(t *testing.T) {

	uuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		name   string
		items  []content.Content
		id     uuid.UUID
		expect []ContentListReadModel
	}{
		{
			name: "get root children",
			items: []content.Content{
				{
					ID:               uuids[0],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "root",
									},
								},
							},
						},
					},
				},
				{
					ID:               uuids[1],
					ParentID:         uuids[0],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "page 1",
									},
								},
							},
						},
					},
				},
				{
					ID:               uuids[2],
					ParentID:         uuids[0],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "page 2",
									},
								},
							},
						},
					},
				},
			},
			id: uuids[0],
			expect: []ContentListReadModel{
				{
					ID:   uuids[1],
					Name: "page 1",
				},
				{
					ID:   uuids[2],
					Name: "page 2",
				},
			},
		},
		{
			name: "get children",
			items: []content.Content{
				{
					ID:               uuids[0],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "root",
									},
								},
							},
						},
					},
				},
				{
					ID:               uuids[1],
					ParentID:         uuids[0],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "page 1",
									},
								},
							},
						},
					},
				},
				{
					ID:               uuids[2],
					ParentID:         uuids[1],
					PublishedVersion: 0,
					Status:           content.Published,
					Version: map[int]content.ContentVersion{
						0: {
							Properties: content.ContentLanguage{
								"sv-SE": {
									content.NameField: content.ContentField{
										Value: "page 2",
									},
								},
							},
						},
					},
				},
			},
			id: uuids[1],
			expect: []ContentListReadModel{
				{
					ID:   uuids[2],
					Name: "page 2",
				},
			},
		},
	}

	cfg := &siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
		},
	}
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
			c.Database("cms").Collection("content").Drop(context.Background())

			repo := content.NewContentRepository(c)
			for _, cnt := range test.items {
				repo.CreateContent(context.Background(), cnt)
			}

			query := ListChildContent{
				ID: test.id,
			}
			handler := ListChildContentHandler{
				Repo: repo,
				Cfg:  cfg,
			}

			children, err := handler.Handle(context.Background(), query)
			assert.NoError(t, err)

			assert.Equal(t, len(test.expect), len(children))
			for _, ch := range children {

				ok := false

				for _, expect := range test.expect {
					if ch.ID == expect.ID {
						ok = true
						assert.Equal(t, expect.Name, ch.Name)
					}
				}

				assert.True(t, ok)
			}
		})
	}
}

func makeInt(n int) *int {
	i := int(n)
	return &i
}
