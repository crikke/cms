package locale

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLocaleHandler(t *testing.T) {

	config := config.Configuration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("nb-NO")},
	}

	tests := []struct {
		acceptlanguage string
		expected       language.Tag
		statuscode     int
	}{
		{
			acceptlanguage: "sv-SE",
			expected:       language.MustParse("sv-SE"),
			statuscode:     200,
		},
		{
			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.5",
			expected:       language.MustParse("nb-NO"),
			statuscode:     200,
		},
		{
			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.1",
			expected:       language.MustParse("sv-SE"),
			statuscode:     200,
		},
		{
			expected:   language.MustParse("sv-SE"),
			statuscode: 200,
		},
		{
			acceptlanguage: "malformed",
			statuscode:     400,
		},
	}

	for _, test := range tests {
		t.Run(test.acceptlanguage, func(t *testing.T) {

			assertHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expected, FromContext(r.Context()))
			})
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Add("Accept-Language", test.acceptlanguage)

			w := httptest.NewRecorder()
			h := Handler(assertHandler, config)
			h.ServeHTTP(w, r)

			assert.Equal(t, test.statuscode, w.Result().StatusCode)
		})
	}
}
