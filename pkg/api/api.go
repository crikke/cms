package api

import (
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/routing"
	"github.com/crikke/cms/pkg/security"
	"github.com/crikke/cms/pkg/services/loader"
	"github.com/gin-gonic/gin"
)

func ContentHandler(group gin.IRouter, c domain.SiteConfiguration, l loader.Loader) {

	r := group.Group("/content", routing.RoutingHandler(c, l))

	r.GET("/*node", security.AuthorizationHandler("read", nil), func(c *gin.Context) {
		c.JSON(200, routing.RoutedNode(*c))
	})
}
