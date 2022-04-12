//go:build integration

package content

// func Test_CreateAndUpdateNewContent(t *testing.T) {

// 	// create content definition
// 	client, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
// 	assert.NoError(t, err)

// 	client.Database("cms").Collection("content").Drop(context.Background())
// 	client.Database("cms").Collection("contentdefinition").Drop(context.Background())

// 	r := api.NewContentManagementAPI(client, nil, nil)

// 	cd, _ := contentdefinition.NewContentDefinition("test contentdefinition", "test desc")

// 	createContent := func(contentdefid uuid.UUID) (url.URL, uuid.UUID, bool) {
// 		t.Helper()

// 		ok := true
// 		type request struct {
// 			ContentDefinitionId uuid.UUID `json:"contentdefinitionid"`
// 		}

// 		body := request{ContentDefinitionId: contentdefid}
// 		var buf bytes.Buffer
// 		err = json.NewEncoder(&buf).Encode(body)
// 		ok = ok && assert.NoError(t, err)
// 		req, err := http.NewRequest(http.MethodPost, "/content", &buf)
// 		ok = ok && assert.NoError(t, err)

// 		res := httptest.NewRecorder()
// 		r.ServeHTTP(res, req)
// 		ok = ok && assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

// 		location, err := res.Result().Location()
// 		ok = ok && assert.NoError(t, err)

// 		actual := &domain.Content{}
// 		json.NewDecoder(res.Body).Decode(actual)

// 		return *location, actual.ID, ok
// 	}

// 	getContent := func(url url.URL, expect domain.ContentData) (uuid.UUID, bool) {
// 		t.Helper()

// 		ok := true
// 		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?version=0", url.String()), nil)
// 		ok = ok && assert.NoError(t, err)

// 		res := httptest.NewRecorder()

// 		r.ServeHTTP(res, req)

// 		actual := &domain.Content{}
// 		json.NewDecoder(res.Body).Decode(actual)

// 		ok = ok && assert.Equal(t, cd.ID, actual.ContentDefinitionID)
// 		ok = ok && assert.Equal(t, expect.Status, actual.Data.Status)

// 		for lang, fields := range expect.Properties {
// 			for fieldname, field := range fields {
// 				ok = ok && assert.Equal(t, field.Value, actual.Data.Properties[lang][fieldname].Value)
// 			}
// 		}
// 		return actual.ID, ok
// 	}

// 	updateContent := func(contentID uuid.UUID) bool {

// 		t.Helper()
// 		ok := true

// 		type request struct {
// 			Version  int
// 			Language string
// 			Fields   map[string]interface{}
// 		}

// 		body := request{
// 			Version:  0,
// 			Language: "sv-SE",
// 			Fields: map[string]interface{}{
// 				contentdefinition.NameField: "updated content",
// 			},
// 		}
// 		var buf bytes.Buffer
// 		err = json.NewEncoder(&buf).Encode(body)
// 		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/content/%s", contentID.String()), &buf)
// 		ok = ok && assert.NoError(t, err)

// 		res := httptest.NewRecorder()
// 		r.ServeHTTP(res, req)

// 		ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)
// 		ok = ok && assert.Equal(t, len(res.Body.Bytes()), 0)

// 		return ok
// 	}

// 	archiveContent := func(contentID uuid.UUID) bool {
// 		t.Helper()
// 		ok := true

// 		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/content/%s", contentID.String()), nil)
// 		ok = ok && assert.NoError(t, err)

// 		res := httptest.NewRecorder()

// 		r.ServeHTTP(res, req)

// 		ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)
// 		ok = ok && assert.Equal(t, len(res.Body.Bytes()), 0)

// 		return ok
// 	}
// 	//TODO: list content
// 	listContent := func(expect []query.ContentListReadModel) bool {
// 		t.Helper()
// 		ok := true

// 		// req, err := http.NewRequest(http.MethodGet, "/content", nil)
// 		// ok = ok && assert.NoError(t, err)

// 		// res := httptest.NewRecorder()

// 		// r.ServeHTTP(res, req)

// 		// result := []query.ContentListReadModel{}

// 		// err = json.NewDecoder(res.Body).Decode(&result)
// 		// ok = ok && assert.NoError(t, err)
// 		// ok = ok && assert.Equal(t, http.StatusOK, res.Result().StatusCode)
// 		// ok = ok && assert.Equal(t, len(expect), len(result))
// 		return ok
// 	}

// 	t.Run("create new content and update it", func(t *testing.T) {

// 		location, contentID, ok := createContent(cd.ID)

// 		ok = ok && updateContent(contentID)

// 		if ok {
// 			_, ok = getContent(location, domain.ContentData{
// 				Status: domain.Draft,

// 				Version: 0,
// 				Properties: domain.ContentLanguage{
// 					"sv-SE": domain.ContentFields{
// 						"name": domain.ContentField{
// 							Value: "updated content",
// 						},
// 					},
// 				},
// 			})
// 		}

// 		ok = ok && archiveContent(contentID)
// 		ok = ok && listContent([]query.ContentListReadModel{})
// 	})
// }
