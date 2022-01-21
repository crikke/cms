package authorization

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2/persist"
)

type Authorizer interface {
	Authorize() bool
}

func Handler(ctx context.Context, policyAdapter persist.Adapter) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//TODO: TODO

		// e, err := casbin.NewEnforcer("model.conf", policyAdapter)

		// if err != nil {
		// 	// TODO: Handle error gracefully
		// 	panic(err)
		// }

		// node := node.RoutedNode(r.Context())

		// group := ""

		// e.Enforce(sub,node.,act)
	})
}
