package siteconfiguration

import (
	"encoding/json"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api/models"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/go-chi/chi/v5"
)

func NewSiteConfigurationRouter(cfg *siteconfiguration.SiteConfiguration) http.Handler {

	r := chi.NewRouter()
	r.Get("/", getSiteConfiguration(cfg))
	r.Post("/", updateSiteConfiguration())
	return r
}

type SiteConfigurationRequest struct {
	Languages []string
}

// @Summary 					Get siteconfiguration
// @Description 				Gets siteconfiguration for this site.
//
// @Tags 						siteconfiguration
// @Produces 					json
// @Success						200			{object}	SiteConfigurationRequest
// @Failure						default		{object}	models.GenericError
// @Router						/siteconfiguration [get]
func getSiteConfiguration(cfg *siteconfiguration.SiteConfiguration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if cfg == nil {

			models.WithError(r.Context(), models.GenericError{
				Body: models.ErrorBody{
					Message: "siteconfiguration is missing",
				},
				StatusCode: http.StatusInternalServerError,
			})

			return
		}
		data, err := json.Marshal(cfg)

		if err != nil {
			models.WithError(r.Context(), err)
			return
		}

		w.Write(data)
	}
}

// @Summary 					Updates siteconfiguration
// @Description 				Updates siteconfiguration for this site.
//
// @Tags 						siteconfiguration
// @Produces 					json
// @Param						body		body	SiteConfigurationRequest	true 	"request body"
// @Success						200			{object}	models.OKResult
// @Failure						default		{object}	models.GenericError
// @Router						/siteconfiguration [put]
func updateSiteConfiguration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
