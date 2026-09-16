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
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// idiomaticEnv is the environment used by the idiomatic tests. It carries a value
// so that tests can assert the environment is actually forwarded to the wrapped
// idiomatic function.
type idiomaticEnv struct {
	Multiplier int
}

var (
	idiomaticTestEnv = idiomaticEnv{Multiplier: 3}
	errIdiomatic     = errors.New("idiomatic test error")
)

// runIdiomatic executes a ReaderIOResult against the shared test environment.
func runIdiomatic[A any](rio ReaderIOResult[idiomaticEnv, A]) Result[A] {
	return rio(idiomaticTestEnv)()
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func TestFromIdiomatic(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		f := func(e idiomaticEnv) func() (int, error) {
			return func() (int, error) { return e.Multiplier * 2, nil }
		}
		assert.Equal(t, result.Of(6), runIdiomatic(FromIdiomatic(f)))
	})

	t.Run("error case", func(t *testing.T) {
		f := func(idiomaticEnv) func() (int, error) {
			return func() (int, error) { return 0, errIdiomatic }
		}
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromIdiomatic(f)))
	})

	t.Run("the IO action is deferred until the thunk is invoked", func(t *testing.T) {
		calls := 0
		f := func(idiomaticEnv) func() (int, error) {
			return func() (int, error) {
				calls++
				return calls, nil
			}
		}
		rio := FromIdiomatic(f)
		assert.Equal(t, 0, calls)

		io := rio(idiomaticTestEnv)
		assert.Equal(t, 0, calls)

		assert.Equal(t, result.Of(1), io())
		assert.Equal(t, result.Of(2), io())
	})
}

func TestFromResultI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(42), runIdiomatic(FromResultI[idiomaticEnv](42, nil)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromResultI[idiomaticEnv](0, errIdiomatic)))
	})
}

func TestFromIOResultI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mr := func() (string, error) { return "ok", nil }
		assert.Equal(t, result.Of("ok"), runIdiomatic(FromIOResultI[idiomaticEnv](mr)))
	})

	t.Run("error case", func(t *testing.T) {
		mr := func() (string, error) { return "", errIdiomatic }
		assert.Equal(t, result.Left[string](errIdiomatic), runIdiomatic(FromIOResultI[idiomaticEnv](mr)))
	})
}

func TestFromReaderResultI(t *testing.T) {
	t.Run("success case, environment is forwarded", func(t *testing.T) {
		rr := func(e idiomaticEnv) (int, error) { return e.Multiplier, nil }
		assert.Equal(t, result.Of(3), runIdiomatic(FromReaderResultI(rr)))
	})

	t.Run("error case", func(t *testing.T) {
		rr := func(idiomaticEnv) (int, error) { return 0, errIdiomatic }
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromReaderResultI(rr)))
	})
}

func TestFromReaderIOResultI(t *testing.T) {
	t.Run("success case, environment is forwarded", func(t *testing.T) {
		rr := func(e idiomaticEnv) func() (int, error) {
			return func() (int, error) { return e.Multiplier * 10, nil }
		}
		assert.Equal(t, result.Of(30), runIdiomatic(FromReaderIOResultI(rr)))
	})

	t.Run("error case", func(t *testing.T) {
		rr := func(idiomaticEnv) func() (int, error) {
			return func() (int, error) { return 0, errIdiomatic }
		}
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(FromReaderIOResultI(rr)))
	})
}

// ---------------------------------------------------------------------------
// ChainI / ChainFirstI / TapI
// ---------------------------------------------------------------------------

// idiomaticScale is an idiomatic ReaderIOResult Kleisli arrow.
func idiomaticScale(n int) func(idiomaticEnv) func() (int, error) {
	return func(e idiomaticEnv) func() (int, error) {
		return func() (int, error) { return n * e.Multiplier, nil }
	}
}

// idiomaticReject is an idiomatic ReaderIOResult Kleisli arrow that always fails.
func idiomaticReject(n int) func(idiomaticEnv) func() (int, error) {
	return func(idiomaticEnv) func() (int, error) {
		return func() (int, error) { return 0, fmt.Errorf("cannot handle %d: %w", n, errIdiomatic) }
	}
}

func TestMonadChainI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(15), runIdiomatic(MonadChainI(Of[idiomaticEnv](5), idiomaticScale)))
	})

	t.Run("error in the Kleisli arrow short-circuits", func(t *testing.T) {
		res := runIdiomatic(MonadChainI(Of[idiomaticEnv](5), idiomaticReject))
		assert.True(t, result.IsLeft(res))
	})

	t.Run("error in the source is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) func(idiomaticEnv) func() (int, error) {
			called = true
			return idiomaticScale(n)
		}
		res := runIdiomatic(MonadChainI(Left[idiomaticEnv, int](errIdiomatic), f))
		assert.Equal(t, result.Left[int](errIdiomatic), res)
		assert.False(t, called)
	})
}

func TestChainI(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		rio := F.Pipe1(Of[idiomaticEnv](5), ChainI(idiomaticScale))
		assert.Equal(t, result.Of(15), runIdiomatic(rio))
	})

	t.Run("error case", func(t *testing.T) {
		rio := F.Pipe1(Of[idiomaticEnv](5), ChainI(idiomaticReject))
		assert.True(t, result.IsLeft(runIdiomatic(rio)))
	})
}

func TestMonadChainFirstI(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstI(Of[idiomaticEnv](5), idiomaticScale)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.True(t, result.IsLeft(runIdiomatic(MonadChainFirstI(Of[idiomaticEnv](5), idiomaticReject))))
	})
}

func TestChainFirstI(t *testing.T) {
	rio := F.Pipe1(Of[idiomaticEnv](5), ChainFirstI(idiomaticScale))
	assert.Equal(t, result.Of(5), runIdiomatic(rio))
}

func TestMonadTapI(t *testing.T) {
	assert.Equal(t, result.Of(5), runIdiomatic(MonadTapI(Of[idiomaticEnv](5), idiomaticScale)))
}

func TestTapI(t *testing.T) {
	seen := 0
	tap := func(n int) func(idiomaticEnv) func() (int, error) {
		return func(idiomaticEnv) func() (int, error) {
			return func() (int, error) {
				seen = n
				return n, nil
			}
		}
	}
	rio := F.Pipe1(Of[idiomaticEnv](7), TapI(tap))
	assert.Equal(t, result.Of(7), runIdiomatic(rio))
	assert.Equal(t, 7, seen)
}

// ---------------------------------------------------------------------------
// Result level: ChainEitherIK / ChainResultIK and friends
// ---------------------------------------------------------------------------

// idiomaticDouble is an idiomatic Result Kleisli arrow.
func idiomaticDouble(n int) (int, error) { return n * 2, nil }

// idiomaticFail is an idiomatic Result Kleisli arrow that always fails.
func idiomaticFail(int) (int, error) { return 0, errIdiomatic }

func TestMonadChainEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainEitherIK(Of[idiomaticEnv](5), idiomaticDouble)))
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainResultIK(Of[idiomaticEnv](5), idiomaticDouble)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainEitherIK(Of[idiomaticEnv](5), idiomaticFail)))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainResultIK(Of[idiomaticEnv](5), idiomaticFail)))
	})
}

func TestChainEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainEitherIK[idiomaticEnv](idiomaticDouble))))
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainResultIK[idiomaticEnv](idiomaticDouble))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainEitherIK[idiomaticEnv](idiomaticFail))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainResultIK[idiomaticEnv](idiomaticFail))))
	})

	t.Run("the source error is propagated", func(t *testing.T) {
		res := runIdiomatic(F.Pipe1(Left[idiomaticEnv, int](errIdiomatic), ChainResultIK[idiomaticEnv](idiomaticDouble)))
		assert.Equal(t, result.Left[int](errIdiomatic), res)
	})
}

func TestChainFirstEitherIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstEitherIK[idiomaticEnv](idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstResultIK[idiomaticEnv](idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapEitherIK[idiomaticEnv](idiomaticDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapResultIK[idiomaticEnv](idiomaticDouble))))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstEitherIK[idiomaticEnv](idiomaticFail))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapResultIK[idiomaticEnv](idiomaticFail))))
	})
}

func TestMonadChainFirstEitherIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstEitherIK(Of[idiomaticEnv](5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstResultIK(Of[idiomaticEnv](5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapEitherIK(Of[idiomaticEnv](5), idiomaticDouble)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapResultIK(Of[idiomaticEnv](5), idiomaticDouble)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainFirstResultIK(Of[idiomaticEnv](5), idiomaticFail)))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadTapEitherIK(Of[idiomaticEnv](5), idiomaticFail)))
	})
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

// idiomaticDeferredDouble is an idiomatic IOResult Kleisli arrow.
func idiomaticDeferredDouble(n int) func() (int, error) {
	return func() (int, error) { return n * 2, nil }
}

// idiomaticDeferredFail is an idiomatic IOResult Kleisli arrow that always fails.
func idiomaticDeferredFail(int) func() (int, error) {
	return func() (int, error) { return 0, errIdiomatic }
}

func TestMonadChainIOEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainIOEitherIK(Of[idiomaticEnv](5), idiomaticDeferredDouble)))
		assert.Equal(t, result.Of(10), runIdiomatic(MonadChainIOResultIK(Of[idiomaticEnv](5), idiomaticDeferredDouble)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainIOEitherIK(Of[idiomaticEnv](5), idiomaticDeferredFail)))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainIOResultIK(Of[idiomaticEnv](5), idiomaticDeferredFail)))
	})
}

func TestChainIOEitherIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainIOEitherIK[idiomaticEnv](idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(10), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainIOResultIK[idiomaticEnv](idiomaticDeferredDouble))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainIOResultIK[idiomaticEnv](idiomaticDeferredFail))))
	})

	t.Run("the effect runs once per execution", func(t *testing.T) {
		calls := 0
		f := func(n int) func() (int, error) {
			return func() (int, error) {
				calls++
				return n, nil
			}
		}
		rio := F.Pipe1(Of[idiomaticEnv](5), ChainIOResultIK[idiomaticEnv](f))
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(5), runIdiomatic(rio))
		assert.Equal(t, 1, calls)
		assert.Equal(t, result.Of(5), runIdiomatic(rio))
		assert.Equal(t, 2, calls)
	})
}

func TestChainFirstIOEitherIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstIOEitherIK[idiomaticEnv](idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstIOResultIK[idiomaticEnv](idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapIOEitherIK[idiomaticEnv](idiomaticDeferredDouble))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapIOResultIK[idiomaticEnv](idiomaticDeferredDouble))))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstIOResultIK[idiomaticEnv](idiomaticDeferredFail))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapIOEitherIK[idiomaticEnv](idiomaticDeferredFail))))
	})
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

// idiomaticScaleByEnv is an idiomatic ReaderResult Kleisli arrow.
func idiomaticScaleByEnv(n int) func(idiomaticEnv) (int, error) {
	return func(e idiomaticEnv) (int, error) { return n * e.Multiplier, nil }
}

// idiomaticRejectByEnv is an idiomatic ReaderResult Kleisli arrow that always fails.
func idiomaticRejectByEnv(int) func(idiomaticEnv) (int, error) {
	return func(idiomaticEnv) (int, error) { return 0, errIdiomatic }
}

func TestMonadChainReaderEitherIK(t *testing.T) {
	t.Run("success case, the environment is forwarded", func(t *testing.T) {
		assert.Equal(t, result.Of(15), runIdiomatic(MonadChainReaderEitherIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
		assert.Equal(t, result.Of(15), runIdiomatic(MonadChainReaderResultIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainReaderEitherIK(Of[idiomaticEnv](5), idiomaticRejectByEnv)))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainReaderResultIK(Of[idiomaticEnv](5), idiomaticRejectByEnv)))
	})
}

func TestChainReaderResultIK(t *testing.T) {
	t.Run("success case, the environment is forwarded", func(t *testing.T) {
		assert.Equal(t, result.Of(15), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainReaderResultIK(idiomaticScaleByEnv))))
		assert.Equal(t, result.Of(15), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainReaderEitherIK(idiomaticScaleByEnv))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainReaderResultIK(idiomaticRejectByEnv))))
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) func(idiomaticEnv) (int, error) {
			called = true
			return idiomaticScaleByEnv(n)
		}
		res := runIdiomatic(F.Pipe1(Left[idiomaticEnv, int](errIdiomatic), ChainReaderResultIK(f)))
		assert.Equal(t, result.Left[int](errIdiomatic), res)
		assert.False(t, called)
	})

	t.Run("a different environment yields a different result", func(t *testing.T) {
		rio := F.Pipe1(Of[idiomaticEnv](5), ChainReaderResultIK(idiomaticScaleByEnv))
		assert.Equal(t, result.Of(50), rio(idiomaticEnv{Multiplier: 10})())
	})
}

func TestMonadChainFirstReaderEitherIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstReaderEitherIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadChainFirstReaderResultIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapReaderEitherIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
		assert.Equal(t, result.Of(5), runIdiomatic(MonadTapReaderResultIK(Of[idiomaticEnv](5), idiomaticScaleByEnv)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadChainFirstReaderResultIK(Of[idiomaticEnv](5), idiomaticRejectByEnv)))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(MonadTapReaderEitherIK(Of[idiomaticEnv](5), idiomaticRejectByEnv)))
	})
}

func TestChainFirstReaderResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstReaderEitherIK(idiomaticScaleByEnv))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstReaderResultIK(idiomaticScaleByEnv))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapReaderEitherIK(idiomaticScaleByEnv))))
		assert.Equal(t, result.Of(5), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapReaderResultIK(idiomaticScaleByEnv))))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), ChainFirstReaderResultIK(idiomaticRejectByEnv))))
		assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv](5), TapReaderResultIK(idiomaticRejectByEnv))))
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
	onNone := func() error { return errIdiomatic }
	chain := ChainOptionIK[idiomaticEnv, string, string](onNone)

	t.Run("the comma-ok pair is lifted into a success", func(t *testing.T) {
		assert.Equal(t, result.Of("localhost"), runIdiomatic(F.Pipe1(Of[idiomaticEnv]("host"), chain(get))))
	})

	t.Run("a false flag is converted into the error from onNone", func(t *testing.T) {
		assert.Equal(t, result.Left[string](errIdiomatic), runIdiomatic(F.Pipe1(Of[idiomaticEnv]("port"), chain(get))))
	})
}

// ---------------------------------------------------------------------------
// Integration
// ---------------------------------------------------------------------------

// TestIdiomaticPipeline combines the different idiomatic levels in a single pipeline
// to make sure they compose.
func TestIdiomaticPipeline(t *testing.T) {
	audited := 0

	rio := F.Pipe4(
		FromResultI[idiomaticEnv](2, nil),
		ChainReaderResultIK(idiomaticScaleByEnv),     // 2 * 3 = 6
		ChainResultIK[idiomaticEnv](idiomaticDouble), // 12
		TapIOResultIK[idiomaticEnv](func(n int) func() (int, error) {
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
		FromResultI[idiomaticEnv](2, nil),
		ChainResultIK[idiomaticEnv](idiomaticFail),
		ChainResultIK[idiomaticEnv](func(n int) (int, error) {
			reached = true
			return n, nil
		}),
	)

	assert.Equal(t, result.Left[int](errIdiomatic), runIdiomatic(rio))
	assert.False(t, reached)
}
