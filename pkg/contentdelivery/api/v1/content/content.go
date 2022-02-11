package content

import (
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/routing"
	"github.com/crikke/cms/pkg/services/loader"
	"github.com/gin-gonic/gin"
)

type endpoint struct {
	cfg    *domain.SiteConfiguration
	loader loader.Loader
}

func RegisterEndpoints(r gin.IRouter, cfg *domain.SiteConfiguration, l loader.Loader) {

	e := endpoint{cfg, l}
	group := r.Group("/content", routing.RoutingHandler(cfg, l))

	group.GET("/*node", e.get)
}

func (e endpoint) get(c *gin.Context) {
	c.JSON(200, routing.RoutedNode(c))
}
