//go:build unit

package locale

import (
	"net/http/httptest"
	"testing"

	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_GetLocale(t *testing.T) {

	config := &siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("nb-NO")},
	}

	tests := []struct {
		acceptlanguage string
		expected       language.Tag
		ok             bool
	}{
		{
			acceptlanguage: "sv-SE",
			expected:       language.MustParse("sv-SE"),
			ok:             true,
		},
		{
			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.5",
			expected:       language.MustParse("nb-NO"),
			ok:             true,
		},
		{
			acceptlanguage: "sv-SE;q=0.1, nb-NO;q=0.1",
			expected:       language.MustParse("sv-SE"),
			ok:             true,
		},
		{
			expected: language.MustParse("sv-SE"),
			ok:       true,
		},
		{
			acceptlanguage: "malformed",
			expected:       language.Tag{},
			ok:             false,
		},
	}

	for _, test := range tests {
		t.Run(test.acceptlanguage, func(t *testing.T) {

			router := gin.Default()

			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Add("Accept-Language", test.acceptlanguage)

			w := httptest.NewRecorder()

			router.GET("/", func(c *gin.Context) {

				tag, err := GetLocale(c.Request, config)

				ok := err == nil

				assert.Equal(t, test.ok, ok)
				assert.Equal(t, test.expected, tag)
			})

			router.ServeHTTP(w, r)
		})
	}
}
