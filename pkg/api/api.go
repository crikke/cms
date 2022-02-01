package api

import (
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/loader"
	"github.com/crikke/cms/pkg/routing"
	"github.com/gin-gonic/gin"
)

func ContentHandler(group *gin.RouterGroup, c config.Configuration, l loader.Loader) {

	r := group.Group("/content", routing.RoutingHandler(config.Configuration{}, nil))

	r.GET("/*node", func(c *gin.Context) {
		c.JSON(200, routing.RoutedNode(*c))
	})
}
