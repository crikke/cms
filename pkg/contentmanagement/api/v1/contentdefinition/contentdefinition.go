package contentdefinition

import "github.com/gin-gonic/gin"

type resource struct {
}

func RegisterEndpoints(r gin.IRouter) {
	res := resource{}

	group := r.Group("/contentdefinition")

	group.GET("/", res.GetAll)
	group.GET("/:id", res.Get)
	group.POST("/", res.Post)
	group.PUT("/", res.Put)
	group.DELETE("/", res.Delete)
}

func (r resource) GetAll(c *gin.Context) {

}

func (r resource) Get(c *gin.Context) {

}

func (r resource) Post(c *gin.Context) {

}

func (r resource) Put(c *gin.Context) {

}

func (r resource) Delete(c *gin.Context) {

}
