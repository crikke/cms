package content

// type key string

// const langKey = key("language")
// const contentKey = key("language")

// type endpoint struct {
// 	app app.App
// }

// func NewContentRoute(app app.App) http.Handler {

// 	r := chi.NewRouter()
// 	ep := endpoint{app: app}

// 	r.Use(ep.localeContext)

// 	r.Get("/", ep.ListContentByTags())
// 	r.Route("/{id}", func(r chi.Router) {
// 		r.Use(idContext)
// 		r.Get("/", ep.GetContentById())
// 	})

// 	return r
// }

// func idContext(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		contentID := chi.URLParam(r, "id")
// 		if contentID == "" {

// 			models.WithError(r.Context(), models.GenericError{
// 				StatusCode: http.StatusBadRequest,
// 				Body: models.ErrorBody{
// 					FieldName: "contentid",
// 					Message:   "parameter contentid is required",
// 				},
// 			})
// 			return
// 		}

// 		cid, err := uuid.Parse(contentID)

// 		if err != nil {
// 			models.WithError(r.Context(), models.GenericError{
// 				StatusCode: http.StatusBadRequest,
// 				Body: models.ErrorBody{
// 					FieldName: "contentid",
// 					Message:   "bad format",
// 				},
// 			})
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), contentKey, cid)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func withID(ctx context.Context) uuid.UUID {

// 	var id uuid.UUID

// 	if r := ctx.Value(contentKey); r != nil {
// 		id = r.(uuid.UUID)
// 	}

// 	return id
// }

// func withLocale(ctx context.Context) string {

// 	if tag := ctx.Value(langKey); tag != nil {

// 		t := tag.(language.Tag)

// 		return t.String()
// 	}

// 	return ""
// }

// func (ep endpoint) localeContext(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		accept := r.Header.Get("Accept-Language")

// 		if accept == "" {

// 			l := ep.app.SiteConfiguration.Languages[0]
// 			ctx = context.WithValue(ctx, langKey, l)

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 			return
// 		}

// 		t, _, err := language.ParseAcceptLanguage(accept)

// 		if err != nil {

// 			models.WithError(r.Context(), models.GenericError{
// 				StatusCode: http.StatusBadRequest,
// 				Body: models.ErrorBody{
// 					FieldName: "",
// 					Message:   "Accept-language",
// 				},
// 			})
// 			return
// 		}
// 		matcher := language.NewMatcher(ep.app.SiteConfiguration.Languages)

// 		tag, _, _ := matcher.Match(t...)

// 		base, _ := tag.Base()
// 		region, _ := tag.Region()
// 		tag, err = language.Compose(base, region)

// 		if err != nil {
// 			models.WithError(r.Context(), models.GenericError{
// 				StatusCode: http.StatusBadRequest,
// 				Body: models.ErrorBody{
// 					FieldName: "",
// 					Message:   "Accept-language",
// 				},
// 			})
// 			return
// 		}
// 		ctx = context.WithValue(ctx, langKey, tag)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// // GetContentById 		godoc
// // @Summary 					Get content by ID
// // @Description 				Gets content by ID and language. If Accept-Language header is not set,
// // @Description					the default language will be used.
// //
// // @Tags 						content
// // @Accept 						json
// // @Produces 					json
// // @Param						id					path	string	true 	"uuid formatted ID." format(uuid)
// // @Param 						Accept-Language 	header 	string 	false 	"content language"
// // @Success						200			{object}	query.ContentResponse
// // @Failure						default		{object}	models.GenericError
// // @Router						/contentdelivery/content/{id} [get]
// func (ep endpoint) GetContentById() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		lang := withLocale(r.Context())
// 		id := withID(r.Context())

// 		content, err := ep.app.Queries.GetContentByID.Handle(r.Context(), query.GetContentByID{
// 			ID:       id,
// 			Language: lang,
// 		})

// 		if err != nil {
// 			models.WithError(r.Context(), err)
// 			return
// 		}

// 		data, err := json.Marshal(&content)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		w.Write(data)
// 	}
// }

// // ListContentByTags 			godoc
// // @Summary 					List content by tags
// // @Description 				Returns a list of content which has specified tags
// //
// // @Tags 						content
// // @Accept 						json
// // @Produces 					json
// // @Param						id					path	string	true 	"uuid formatted ID." format(uuid)
// // @Param 						Accept-Language 	header 	string 	false 	"content language"
// // @Success						200			{object}	query.ContentResponse
// // @Failure						default		{object}	models.GenericError
// // @Router						/contentdelivery/content/ [get]
// func (ep endpoint) ListContentByTags() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
