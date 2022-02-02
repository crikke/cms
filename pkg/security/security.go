package security

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/routing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthorizationHandler Handler must be run after Routing Handler & Authentication Handler
func AuthorizationHandler(act string, cfg *config.ServerConfiguration) gin.HandlerFunc {
	return func(c *gin.Context) {

		a, err := mongodbadapter.NewAdapter(cfg.ConnectionString.Mongodb)

		if err != nil {
			panic(err)
		}

		e, err := casbin.NewEnforcer("model.conf", a)

		if err != nil {
			// TODO: Handle error gracefully
			panic(err)
		}

		e.LoadPolicy()
		node := routing.RoutedNode(*c)
		user, _ := CurrentUser(c)
		allowed, err := e.Enforce(user.GetID(), node.ID.ID, act)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if !allowed {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (Identity, bool) {
	u, exists := c.Get("user")
	return u.(Identity), exists
}

type Identity interface {
	GetID() uuid.UUID
	GetName() string
}

func MockAuthenticationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user", domain.User{ID: uuid.MustParse("b39ca6e6-c08d-4351-9007-d3e232259b5a"), Name: "test"})
		c.Next()
	}
}
