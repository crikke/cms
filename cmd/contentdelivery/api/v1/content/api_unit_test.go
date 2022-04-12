//go:build unit

package content

// func Test_GetLocale(t *testing.T) {

// 	tests := []struct {
// 		acceptlanguage string
// 		expected       language.Tag
// 		ok             bool
// 	}{
// 		{
// 			acceptlanguage: "sv-SE",
// 			expected:       language.MustParse("sv-SE"),
// 			ok:             true,
// 		},
// 		{
// 			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.5",
// 			expected:       language.MustParse("nb-NO"),
// 			ok:             true,
// 		},
// 		{
// 			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.1",
// 			expected:       language.MustParse("sv-SE"),
// 			ok:             true,
// 		},
// 		{
// 			expected: language.MustParse("sv-SE"),
// 			ok:       true,
// 		},
// 		{
// 			acceptlanguage: "malformed",
// 			expected:       language.Tag{},
// 			ok:             false,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.acceptlanguage, func(t *testing.T) {

// 			req := httptest.NewRequest("GET", "/", nil)
// 			req.Header.Add("Accept-Language", test.acceptlanguage)

// 			w := httptest.NewRecorder()
// 			ep := endpoint{app: app.App{SiteConfiguration: config}}

// 			ep.localeContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				actual := withLocale(r.Context())

// 				assert.Equal(t, test.expected.String(), actual)
// 			})).ServeHTTP(w, req)

// 		})
// 	}
// }
