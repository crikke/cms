package api

import (
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/loader"
	"github.com/crikke/cms/pkg/routing"
	"github.com/crikke/cms/pkg/security"
	"github.com/gin-gonic/gin"
)

func ContentHandler(group gin.IRouter, c config.Configuration, l loader.Loader) {

	r := group.Group("/content", routing.RoutingHandler(c, l))

	r.GET("/*node", security.Authorization("read", nil), func(c *gin.Context) {
		c.JSON(200, routing.RoutedNode(*c))
	})
}