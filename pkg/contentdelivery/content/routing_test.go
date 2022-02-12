package content

// func TestMatchRoute(t *testing.T) {

// 	root := uuid.New()
// 	a := uuid.New()
// 	b := uuid.New()

// 	nodes := []Content{
// 		{
// 			ID:         ContentReference{ID: root},
// 			URLSegment: "",
// 		},
// 		{
// 			ID:         ContentReference{ID: a},
// 			ParentID:   root,
// 			URLSegment: "a",
// 		},
// 		{
// 			ParentID:   a,
// 			ID:         ContentReference{ID: b},
// 			URLSegment: "b",
// 		},
// 	}

// 	tests := []struct {
// 		description        string
// 		url                string
// 		language           string
// 		expectedNode       Content
// 		expectedStatusCode int
// 	}{
// 		{
// 			description:        "node matched with sv language",
// 			url:                "/a/b",
// 			expectedNode:       nodes[2],
// 			expectedStatusCode: http.StatusOK,
// 		},
// 		{
// 			description:        "multiple slashes in path ",
// 			url:                `/a///b`,
// 			expectedNode:       nodes[2],
// 			expectedStatusCode: http.StatusOK,
// 		},
// 		{
// 			description:        "path ends with '/' ",
// 			url:                "/a/b/",
// 			expectedNode:       nodes[2],
// 			expectedStatusCode: http.StatusOK,
// 		},
// 		{
// 			description:        "route not found ",
// 			url:                "/a/b/c/d",
// 			expectedNode:       Content{},
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 		{
// 			description:        "root",
// 			url:                "/",
// 			expectedNode:       nodes[0],
// 			expectedStatusCode: http.StatusOK,
// 		},
// 	}

// 	for _, test := range tests {

// 		t.Run(test.description, func(t *testing.T) {
// 			router := gin.Default()

// 			router.GET("/*node", RoutingHandler(&domain.SiteConfiguration{
// 				RootPage: root,
// 			}, mocks.MockLoader{
// 				Nodes: nodes,
// 			}), func(c *gin.Context) {
// 				assert.Equal(t, test.expectedNode.ID, RoutedNode(c).ID, test.description)
// 			})

// 			w := httptest.NewRecorder()
// 			r, _ := http.NewRequest("GET", test.url, nil)

// 			router.ServeHTTP(w, r)
// 			assert.Equal(t, test.expectedStatusCode, w.Code)
// 		})
// 	}
// }

// func Test_GenerateUrl(t *testing.T) {
// 	a := uuid.New()
// 	b := uuid.New()
// 	c := uuid.New()

// 	nodes := []Content{
// 		{
// 			ID:         ContentReference{ID: a},
// 			URLSegment: "a",
// 		},
// 		{
// 			ID:         ContentReference{ID: b},
// 			ParentID:   a,
// 			URLSegment: "b",
// 		},
// 		{
// 			ParentID:   b,
// 			ID:         ContentReference{ID: c},
// 			URLSegment: "c",
// 		},
// 	}

// 	tests := []struct {
// 		description string
// 		node        Content
// 		expected    string
// 	}{
// 		{
// 			description: "generate url node c",
// 			expected:    "/a/b/c/",
// 			node:        nodes[2],
// 		},
// 		{
// 			description: "generate url node b ",
// 			expected:    "/a/b/",
// 			node:        nodes[1],
// 		},
// 	}

// 	for _, test := range tests {

// 		t.Run(test.description, func(t *testing.T) {

// 			mock := mocks.MockLoader{
// 				Nodes: nodes,
// 			}

// 			result, err := GenerateUrl(context.TODO(), mock, test.node.ID)

// 			assert.NoError(t, err)

// 			assert.Equal(t, test.expected, result)
// 		})
// 	}
// }
