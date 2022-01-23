package content

// var contentKey key

// Returns specified version of the content or published version if none is specified.
//   returns http.StatusNotFound if version does not exist
//   returns http.StatusBadRequest if invalid parameter
// func VersionHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		node := RoutedNode(r.Context())
// 		v := r.URL.Query().Get("v")

// 		if v == "" {
// 			ctx := WithContent(r.Context(), node.content[node.version])

// 			r = r.WithContext(ctx)
// 			next.ServeHTTP(rw, r)
// 			return
// 		}
// 		version, err := strconv.Atoi(v)

// 		if err != nil {
// 			r.Response.StatusCode = http.StatusBadRequest
// 			return
// 		}

// 		content, exist := node.content[version]

// 		if !exist {
// 			r.Response.StatusCode = http.StatusNotFound
// 			return
// 		}

// 		ctx := WithContent(r.Context(), content)

// 		r = r.WithContext(ctx)
// 		next.ServeHTTP(rw, r)
// 	})
// }

// func WithContent(ctx context.Context, content ContentData) context.Context {
// 	return context.WithValue(ctx, contentKey, content)
// }

// func Content(ctx context.Context) ContentData {
// 	return ctx.Value(contentKey).(ContentData)
// }
