package locale

import (
	"net/http"

	"github.com/crikke/cms/pkg/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

var languageKey = "languagekey"

/*
	Prefered language is set by contentmanagement API,
	If Accept-Header isnt set, configured default language is used as fallback
*/
func Handler(cfg *domain.SiteConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept-Language")

		t, _, err := language.ParseAcceptLanguage(accept)

		if err != nil {
			c.String(http.StatusBadRequest, "Accept-Language")
			return
		}

		matcher := language.NewMatcher(cfg.Languages)

		tag, _, _ := matcher.Match(t...)
		c.Set(string(languageKey), tag)
		c.Next()
	}
}

func FromContext(c *gin.Context) language.Tag {
	t, exists := c.Get(languageKey)

	if !exists {
		t = language.Tag{}
	}

	return t.(language.Tag)
}
