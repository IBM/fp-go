package context

import (
	"context"

	"github.com/IBM/fp-go/v2/pair"
)

func noop() {}

func NopCancel(ctx context.Context) ContextCancel {
	return pair.MakePair[context.CancelFunc](noop, ctx)
}

func WithValue[A, K any](key K) Kleisli[A, context.Context] {
	return func(val A) Reader[context.Context] {
		return func(ctx context.Context) context.Context {
			return context.WithValue(ctx, key, val)
		}
	}
}
