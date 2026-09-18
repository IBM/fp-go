package context

import (
	"context"

	"github.com/IBM/fp-go/v2/internal/common"
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

// AskValue reads the value stored under key and asserts it to V.
// It returns None if the key is absent or the stored value is not a V.
func AskValue[V, K any](key K) Reader[common.Option[V]] {
	return func(ctx context.Context) common.Option[V] {
		return common.OptionInstanceOf[V](ctx.Value(key))
	}
}

// WithValueNopCancel derives a child context carrying key → value, paired
// with a no-op cancel function, suitable for Local-style combinators.
func WithValueNopCancel[K, V any](key K, value V) Reader[ContextCancel] {
	return func(ctx context.Context) ContextCancel {
		return NopCancel(WithValue[V](key)(value)(ctx))
	}
}
