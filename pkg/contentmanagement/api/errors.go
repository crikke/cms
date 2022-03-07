package api

import (
	"context"
	"net/http"
)

func init() {

}

type key string

var errorKey = key("error")

func HandleHttpError(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		ctx := r.Context()
		if err := ctx.Value(errorKey); err != nil {

			if ge, ok := err.(GenericError); ok {
				http.Error(w, ge.Error(), ge.StatusCode)
				return
			}

			http.Error(w, (err.(error)).Error(), http.StatusInternalServerError)
		}
	})

}

func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

type ErrorBody struct {
	// required: true
	Message   string
	FieldName string
}

// GenericError
// swagger:model genericError
type GenericError struct {
	// in: body
	Body ErrorBody
	// swagger:ignore
	StatusCode int
}

func (g GenericError) Error() string {
	return g.Body.Message
}
