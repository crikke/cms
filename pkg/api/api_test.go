package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApiSlug(t *testing.T) {

	router := gin.Default()

	router.GET("/content/*node", func(c *gin.Context) {
		assert.Equal(t, "/a/b/c", c.Param("node"))
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/content/a/b/c?foo=11", nil)

	router.ServeHTTP(w, r)

}
