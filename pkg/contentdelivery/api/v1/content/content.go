package content

import (
	"context"
	"net/http"

	"github.com/crikke/cms/pkg/contentdelivery/app"
	"github.com/crikke/cms/pkg/contentmanagement/api/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

type key string

const langKey = key("language")

type endpoint struct {
	app app.App
}

func NewContentRoute(app app.App) http.Handler {

	r := chi.NewRouter()
	ep := endpoint{app: app}

	r.Use(ep.localeContext)
	r.Get("/{id}", ep.GetContentById())

	return r
}

// GetContentById 		godoc
// @Summary 					Get content by ID
// @Description 				Gets content by ID and language. If Accept-Language header is not set,
// @Description					the default language will be used.
//
// @Tags 						contentdefinition
// @Accept 						json
// @Produces 					json
// @Param						id					path	string	true 	"uuid formatted ID." format(uuid)
// @Param 						Accept-Language 	header 	string 	false 	"content language"
// @Success						200			{object}	query.ContentResponse
// @Failure						default		{object}	models.GenericError
// @Router						/content/{id} [get]
func (ep endpoint) GetContentById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func withLocale(ctx context.Context) string {

	if tag := ctx.Value(langKey); tag != nil {

		t := tag.(language.Tag)

		return t.String()
	}

	return ""
}

func (ep endpoint) localeContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accept := r.Header.Get("Accept-Language")

		if accept == "" {

			l := ep.app.SiteConfiguration.Languages[0]
			ctx = context.WithValue(ctx, langKey, l)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		t, _, err := language.ParseAcceptLanguage(accept)

		if err != nil {

			models.WithError(r.Context(), models.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: models.ErrorBody{
					FieldName: "",
					Message:   "Accept-language",
				},
			})
			return
		}
		matcher := language.NewMatcher(ep.app.SiteConfiguration.Languages)

		tag, _, _ := matcher.Match(t...)

		base, _ := tag.Base()
		region, _ := tag.Region()
		tag, err = language.Compose(base, region)

		if err != nil {
			models.WithError(r.Context(), models.GenericError{
				StatusCode: http.StatusBadRequest,
				Body: models.ErrorBody{
					FieldName: "",
					Message:   "Accept-language",
				},
			})
			return
		}
		ctx = context.WithValue(ctx, langKey, tag)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
