// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readerioresult

import (
	"context"
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// idiomaticCtxKey is the key used to smuggle a value through the context, so the
// tests can assert the context is actually forwarded to the wrapped function.
type idiomaticCtxKey struct{}

var errIdiomatic = errors.New("idiomatic test error")

// idiomaticCtx returns a context carrying the given multiplier.
func idiomaticCtx(multiplier int) context.Context {
	return context.WithValue(context.Background(), idiomaticCtxKey{}, multiplier)
}

// idiomaticMultiplier extracts the multiplier from the context, defaulting to 1.
func idiomaticMultiplier(ctx context.Context) int {
	if m, ok := ctx.Value(idiomaticCtxKey{}).(int); ok {
		return m
	}
	return 1
}

// runIdiomatic executes a ReaderIOResult against a context carrying multiplier 3.
func runIdiomatic[A any](rio ReaderIOResult[A]) Result[A] {
	return rio(idiomaticCtx(3))()
}

// ---------------------------------------------------------------------------
// Idiomatic fixtures
// ---------------------------------------------------------------------------

// idiomaticDouble is an idiomatic Result Kleisli arrow.
func idiomaticDouble(n int) (int, error) { return n * 2, nil }

// idiomaticFail is an idiomatic Result Kleisli arrow that always fails.
func idiomaticFail(int) (int, error) { return 0, errIdiomatic }

// idiomaticDeferredDouble is an idiomatic IOResult Kleisli arrow.
func idiomaticDeferredDouble(n int) func() (int, error) {
	return func() (int, error) { return n * 2, nil }
}

// idiomaticDeferredFail is an idiomatic IOResult Kleisli arrow that always fails.
func idiomaticDeferredFail(int) func() (int, error) {
	return func() (int, error) { return 0, errIdiomatic }
}

// idiomaticScaleByCtx is an idiomatic ReaderResult Kleisli arrow.
func idiomaticScaleByCtx(n int) func(context.Context) (int, error) {
	return func(ctx context.Context) (int, error) { return n * idiomaticMultiplier(ctx), nil }
}

// idiomaticRejectByCtx is an idiomatic ReaderResult Kleisli arrow that always fails.
func idiomaticRejectByCtx(int) func(context.Context) (int, error) {
	return func(context.Context) (int, error) { return 0, errIdiomatic }
}

// idiomaticScale is an idiomatic ReaderIOResult Kleisli arrow.
func idiomaticScale(n int) func(context.Context) func() (int, error) {
	return func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return n * idiomaticMultiplier(ctx), nil }
	}
}

// idiomaticReject is an idiomatic ReaderIOResult Kleisli arrow that always fails.
func idiomaticReject(int) func(context.Context) func() (int, error) {
	return func(context.Context) func() (int, error) {
		return func() (int, error) { return 0, errIdiomatic }
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func TestFromIdiomatic(t *testing.T) {
	t.Run("success case, the context is forwarded", func(t *testing.T) {
		f := func(ctx context.Context) func() (int, error) {
			return func() (int, error) { return idiomaticMultiplier(ctx) * 2, nil }
		}
		assert.Equal(t, result.Of(6), runIdiomatic(FromIdiomatic(f)))
	})

	t.Run("error case", func(t *testing.T) {
		f := func(context.Context) func() (int, error) {
			return func() (int, error) { return 0, errIdiomatic }
		}
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromIdiomatic(f)))
	})
}

func TestFromResultI(t *testing.T) {
	assert.Equal(t, result.Of(42), runIdiomatic(FromResultI(42, nil)))
	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromResultI(0, errIdiomatic)))
}

func TestFromIOResultI(t *testing.T) {
	assert.Equal(t, result.Of("ok"), runIdiomatic(FromIOResultI(func() (string, error) { return "ok", nil })))
	assert.Equal(t, result.Left[string](errIdiomatic), runIdiomatic(FromIOResultI(func() (string, error) { return "", errIdiomatic })))
}

func TestFromReaderResultI(t *testing.T) {
	t.Run("success case, the context is forwarded", func(t *testing.T) {
		rr := func(ctx context.Context) (int, error) { return idiomaticMultiplier(ctx), nil }
		assert.Equal(t, result.Of(3), runIdiomatic(FromReaderResultI(rr)))
	})

	t.Run("error case", func(t *testing.T) {
		rr := func(context.Context) (int, error) { return 0, errIdiomatic }
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromReaderResultI(rr)))
	})
}

func TestFromReaderIOResultI(t *testing.T) {
	rr := func(ctx context.Context) func() (int, error) {
		return func() (int, error) { return idiomaticMultiplier(ctx) * 10, nil }
	}
	assert.Equal(t, result.Of(30), runIdiomatic(FromReaderIOResultI(rr)))
}

// ---------------------------------------------------------------------------
// ChainI / ChainFirstI / TapI
// ---------------------------------------------------------------------------

func TestMonadChainI(t *testing.T) {
	assert.Equal(t, result.Of(15), runIdiomatic(MonadChainI(Of(5), idiomaticScale)))
	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainI(Of(5), idiomaticReject)))
}

func TestChainI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(15), runIdiomatic(F.Pipe1(Of(5), ChainI(idiomaticScale))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainI(idiomaticReject))))
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) func(context.Context) func() (int, error) {
			called = true
			return idiomaticScale(n)
		}
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Left[int](errIdiomatic), ChainI(f))))
		assert.False(t, called)
	})
}

func TestChainFirstI(t *testing.T) {
	assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstI(idiomaticScale))))
	assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstI(Of(5), idiomaticScale)))
	assert.Equal(t, result.Of(5), runIdiomatic(MonadTapI(Of(5), idiomaticScale)))
	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainFirstI(idiomaticReject))))
}

func TestTapI(t *testing.T) {
	seen := 0
	tap := func(n int) func(context.Context) func() (int, error) {
		return func(context.Context) func() (int, error) {
			return func() (int, error) {
				seen = n
				return n, nil
			}
		}
	}
	assert.Equal(t, result.Of(7), runIdiomatic(F.Pipe1(Of(7), TapI(tap))))
	assert.Equal(t, 7, seen)
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

func TestChainResultIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of(5), ChainResultIK(idiomaticDouble))))
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of(5), ChainEitherIK(idiomaticDouble))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainResultIK(idiomaticFail))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainEitherIK(idiomaticFail))))
	})
}

func TestMonadChainResultIK(t *testing.T) {
	assert.Equal(t, result.Of(10), runIdiomatic(MonadChainResultIK(Of(5), idiomaticDouble)))
	assert.Equal(t, result.Of(10), runIdiomatic(MonadChainEitherIK(Of(5), idiomaticDouble)))
	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainResultIK(Of(5), idiomaticFail)))
}

func TestChainFirstResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstResultIK(idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstEitherIK(idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapResultIK(idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapEitherIK(idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstResultIK(Of(5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstEitherIK(Of(5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapResultIK(Of(5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapEitherIK(Of(5), idiomaticDouble)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainFirstResultIK(idiomaticFail))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), TapEitherIK(idiomaticFail))))
	})
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

func TestChainIOResultIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of(5), ChainIOResultIK(idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of(5), ChainIOEitherIK(idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainIOResultIK(Of(5), idiomaticDeferredDouble)))
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainIOEitherIK(Of(5), idiomaticDeferredDouble)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainIOResultIK(idiomaticDeferredFail))))
	})
}

func TestChainFirstIOResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstIOResultIK(idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstIOEitherIK(idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapIOResultIK(idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapIOEitherIK(idiomaticDeferredDouble))))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), TapIOResultIK(idiomaticDeferredFail))))
	})
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

func TestChainReaderResultIK(t *testing.T) {
	t.Run("success case, the context is forwarded", func(t *testing.T) {
		assert.Equal(t, result.Of(15), runIdiomatic(F.Pipe1(Of(5), ChainReaderResultIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(15), runIdiomatic(F.Pipe1(Of(5), ChainReaderEitherIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(15), runIdiomatic(MonadChainReaderResultIK(Of(5), idiomaticScaleByCtx)))
		assert.Equal(t, result.Of(15), runIdiomatic(MonadChainReaderEitherIK(Of(5), idiomaticScaleByCtx)))
	})

	t.Run("a different context yields a different result", func(t *testing.T) {
		rio := F.Pipe1(Of(5), ChainReaderResultIK(idiomaticScaleByCtx))
		assert.Equal(t, result.Of(50), rio(idiomaticCtx(10))())
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), ChainReaderResultIK(idiomaticRejectByCtx))))
	})
}

func TestChainFirstReaderResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstReaderResultIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), ChainFirstReaderEitherIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapReaderResultIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of(5), TapReaderEitherIK(idiomaticScaleByCtx))))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstReaderResultIK(Of(5), idiomaticScaleByCtx)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstReaderEitherIK(Of(5), idiomaticScaleByCtx)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapReaderResultIK(Of(5), idiomaticScaleByCtx)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapReaderEitherIK(Of(5), idiomaticScaleByCtx)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of(5), TapReaderResultIK(idiomaticRejectByCtx))))
	})
}

// ---------------------------------------------------------------------------
// Option level
// ---------------------------------------------------------------------------

func TestChainOptionIK(t *testing.T) {
	lookup := map[string]string{"host": "localhost"}
	get := func(k string) (string, bool) {
		v, ok := lookup[k]
		return v, ok
	}
	chain := ChainOptionIK[string, string](func() error { return errIdiomatic })

	assert.Equal(t, result.Of("localhost"), runIdiomatic(F.Pipe1(Of("host"), chain(get))))
	assert.Equal(t, result.Left[string](errIdiomatic), runIdiomatic(F.Pipe1(Of("port"), chain(get))))
}

// ---------------------------------------------------------------------------
// Integration
// ---------------------------------------------------------------------------

// TestIdiomaticPipeline combines the different idiomatic levels in a single
// pipeline to make sure they compose.
func TestIdiomaticPipeline(t *testing.T) {
	audited := 0

	rio := F.Pipe4(
		FromResultI(2, nil),
		ChainReaderResultIK(idiomaticScaleByCtx), // 2 * 3 = 6
		ChainResultIK(idiomaticDouble),           // 12
		TapIOResultIK(func(n int) func() (int, error) {
			return func() (int, error) {
				audited = n
				return n, nil
			}
		}),
		ChainI(idiomaticScale), // 12 * 3 = 36
	)

	assert.Equal(t, result.Of(36), runIdiomatic(rio))
	assert.Equal(t, 12, audited)
}

// TestIdiomaticPipelineError verifies the pipeline short-circuits on the first error.
func TestIdiomaticPipelineError(t *testing.T) {
	reached := false

	rio := F.Pipe2(
		FromResultI(2, nil),
		ChainResultIK(idiomaticFail),
		ChainResultIK(func(n int) (int, error) {
			reached = true
			return n, nil
		}),
	)

	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(rio))
	assert.False(t, reached)
}
