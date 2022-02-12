package content

import (
	"github.com/crikke/cms/pkg/contentdelivery/content"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/gin-gonic/gin"
)

type endpoint struct {
	cfg  *siteconfiguration.SiteConfiguration
	repo content.ContentRepository
}

func RegisterEndpoints(r gin.IRouter, cfg *siteconfiguration.SiteConfiguration, repo content.ContentRepository) {

	e := endpoint{cfg, repo}
	group := r.Group("/content")

	group.Use(content.RoutingHandler(cfg, repo))
	group.GET("/*node", e.get)
}

func (e endpoint) get(c *gin.Context) {
	c.JSON(200, content.RoutedNode(c))
}
