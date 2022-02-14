package locale

import (
	"errors"
	"net/http"

	"github.com/crikke/cms/pkg/siteconfiguration"
	"golang.org/x/text/language"
)

func GetLocale(r *http.Request, cfg *siteconfiguration.SiteConfiguration) (language.Tag, error) {

	accept := r.Header.Get("Accept-Language")

	t, _, err := language.ParseAcceptLanguage(accept)

	if err != nil {

		return language.Tag{}, errors.New("Accept-Language")
	}

	matcher := language.NewMatcher(cfg.Languages)

	tag, _, _ := matcher.Match(t...)

	base, _ := tag.Base()
	region, _ := tag.Region()
	tag, err = language.Compose(base, region)

	if err != nil {
		return language.Tag{}, errors.New("Accept-Language")
	}

	return tag, nil
}
