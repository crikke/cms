//go:build integration

package query

// func Test_GetContent(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		content   content.Content
// 		query     GetContent
// 		expect    ContentReadModel
// 		expectErr string
// 	}{
// 		{
// 			name: "get latest version",
// 			content: content.Content{
// 				Data: content.ContentData{
// 					Version: 1,
// 					Status:  content.Published,
// 					Properties: content.ContentLanguage{
// 						"foo": {
// 							"bar": content.ContentField{
// 								Value: "baz",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			query: GetContent{},
// 			expect: ContentReadModel{
// 				Status: content.Published,
// 				Properties: content.ContentLanguage{
// 					"foo": {
// 						"bar": content.ContentField{
// 							Value: "baz",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "get previous version",
// 			content: content.Content{
// 				Data:   3,
// 				Status: content.Published,
// 				Version: map[int]content.ContentData{
// 					1: {
// 						Status: content.PreviouslyPublished,
// 						Properties: content.ContentLanguage{
// 							"foo": {
// 								"bar": content.ContentField{
// 									Value: "1",
// 								},
// 							},
// 						},
// 					},
// 					2: {
// 						Properties: content.ContentLanguage{
// 							"foo": {
// 								"bar": content.ContentField{
// 									Value: "2",
// 								},
// 							},
// 						},
// 					},
// 					3: {
// 						Status: content.Published,
// 						Properties: content.ContentLanguage{
// 							"foo": {
// 								"bar": content.ContentField{
// 									Value: "3",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			query: GetContent{Version: makeInt(1)},
// 			expect: ContentReadModel{
// 				Status: content.PreviouslyPublished,
// 				Properties: content.ContentLanguage{
// 					"foo": {
// 						"bar": content.ContentField{
// 							Value: "1",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "version not exist",
// 			content: content.Content{
// 				Data:   1,
// 				Status: content.Published,
// 				Version: map[int]content.ContentData{
// 					1: {
// 						Properties: content.ContentLanguage{
// 							"foo": {
// 								"bar": content.ContentField{
// 									Value: "2",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			query:     GetContent{Version: makeInt(44)},
// 			expect:    ContentReadModel{},
// 			expectErr: content.ErrMissingVersion,
// 		},
// 		{
// 			name: "negative version",
// 			content: content.Content{
// 				Data:   1,
// 				Status: content.Published,
// 				Version: map[int]content.ContentData{
// 					1: {
// 						Properties: content.ContentLanguage{
// 							"foo": {
// 								"bar": content.ContentField{
// 									Value: "3",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			query:     GetContent{Version: makeInt(-2)},
// 			expect:    ContentReadModel{},
// 			expectErr: content.ErrMissingVersion,
// 		},
// 	}

// 	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
// 	assert.NoError(t, err)

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
// 			c.Database("cms").Collection("content").Drop(context.Background())

// 			contentRepo := content.NewContentRepository(c)

// 			id, err := contentRepo.CreateContent(context.Background(), test.content)
// 			assert.NoError(t, err)

// 			test.query.Id = id

// 			handler := GetContentHandler{Repo: contentRepo}
// 			actual, err := handler.Handle(context.Background(), test.query)

// 			if test.expectErr != "" {
// 				assert.Equal(t, test.expectErr, err.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			assert.Equal(t, test.expect.Status, actual.Status)
// 			eq := reflect.DeepEqual(test.expect.Properties, actual.Properties)
// 			assert.True(t, eq)
// 		})
// 	}
// }

// func Test_ListChildContent(t *testing.T) {

// 	tests := []struct {
// 		name   string
// 		items  []content.Content
// 		expect []ContentListReadModel
// 	}{
// 		{
// 			name: "get root children",
// 			items: []content.Content{
// 				{
// 					ID:     uuid.New(),
// 					Data:   0,
// 					Status: content.Published,
// 					Version: map[int]content.ContentData{
// 						0: {
// 							Properties: content.ContentLanguage{
// 								"sv-SE": {
// 									contentdefinition.NameField: content.ContentField{
// 										Value: "root",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				{
// 					ID:     uuid.New(),
// 					Data:   0,
// 					Status: content.Published,
// 					Version: map[int]content.ContentData{
// 						0: {
// 							Properties: content.ContentLanguage{
// 								"sv-SE": {
// 									contentdefinition.NameField: content.ContentField{
// 										Value: "page 1",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				{
// 					ID:     uuid.New(),
// 					Data:   0,
// 					Status: content.Published,
// 					Version: map[int]content.ContentData{
// 						0: {
// 							Properties: content.ContentLanguage{
// 								"sv-SE": {
// 									contentdefinition.NameField: content.ContentField{
// 										Value: "page 2",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			expect: []ContentListReadModel{
// 				{
// 					Name: "page 1",
// 				},
// 				{
// 					Name: "page 2",
// 				},
// 				{
// 					Name: "root",
// 				},
// 			},
// 		},
// 	}

// 	cfg := &siteconfiguration.SiteConfiguration{
// 		Languages: []language.Tag{
// 			language.MustParse("sv-SE"),
// 		},
// 	}
// 	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
// 	assert.NoError(t, err)

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			c.Database("cms").Collection("contentdefinition").Drop(context.Background())
// 			c.Database("cms").Collection("content").Drop(context.Background())

// 			repo := content.NewContentRepository(c)
// 			for _, cnt := range test.items {
// 				repo.CreateContent(context.Background(), cnt)
// 			}

// 			query := ListContent{}
// 			handler := ListContentHandler{
// 				Repo: repo,
// 				Cfg:  cfg,
// 			}

// 			children, err := handler.Handle(context.Background(), query)
// 			assert.NoError(t, err)

// 			assert.Equal(t, len(test.expect), len(children))
// 			for _, ch := range children {

// 				ok := false

// 				for _, expect := range test.expect {
// 					if ch.Name == expect.Name {
// 						ok = true
// 						assert.Equal(t, expect.Name, ch.Name)
// 					}
// 				}

// 				assert.True(t, ok)
// 			}
// 		})
// 	}
// }

// func makeInt(n int) *int {
// 	return &n
// }
