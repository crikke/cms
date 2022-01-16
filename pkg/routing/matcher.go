package routing

import "context"

type Matcher interface {
	Match(ctx context.Context, remainingPath string) bool
}
