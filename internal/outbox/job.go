package outbox

import "context"

type Job interface {
	Name() string
	Handle(ctx context.Context, payload string) error
}
