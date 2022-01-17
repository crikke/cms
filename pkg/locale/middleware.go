package locale

import (
	"context"
	"net/http"

	"github.com/crikke/cms/pkg/config"
	"golang.org/x/text/language"
)

type key int

var languageKey key

/*
Prefered language is set by contentmanagement API,
*/
func Handler(next http.Handler, cfg config.Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept-Language")

		t, _, err := language.ParseAcceptLanguage(accept)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Accept-Language"))
			return
		}

		matcher := language.NewMatcher(cfg.Languages)

		tag, _, _ := matcher.Match(t...)
		ctx := WithLocale(r.Context(), tag)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) language.Tag {
	return ctx.Value(languageKey).(language.Tag)
}

func WithLocale(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, languageKey, tag)
}
