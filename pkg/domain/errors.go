package domain

import (
	"net/http"
)

type ErrorResponse struct {
	ID         ContentReference
	Message    string
	StatusCode int
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func ErrContentNotFound(ref ContentReference) ErrorResponse {
	return ErrorResponse{ID: ref, StatusCode: http.StatusNotFound, Message: ref.String()}
}
