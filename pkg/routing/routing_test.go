package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchRoute(t *testing.T) {

	// tests := []struct {
	// 	Description string
	// 	URL	string
	// 	Language string
	// }

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	req := httptest.NewRequest("GET", "https://www.foo.com", nil)
	handler := RouteHandler(testHandler, nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

}
