package security

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/crikke/cms/pkg/domain"
	"github.com/crikke/cms/pkg/routing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Authorization Handler must be run after Routing Handler & Authentication Handler
func Authorization(act string, policyAdapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {

		e, err := casbin.NewEnforcer("model.conf", policyAdapter)

		if err != nil {
			// TODO: Handle error gracefully
			panic(err)
		}

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
