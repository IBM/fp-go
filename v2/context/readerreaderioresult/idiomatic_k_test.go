// Copyright (c) 2024 IBM Corp.
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

package readerreaderioresult

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// runIK executes an effect against the shared test config and the test context.
func runIK[A any](t *testing.T, eff ReaderReaderIOResult[AppConfig, A]) Result[A] {
	t.Helper()
	return eff(defaultConfig)(t.Context())()
}

// ---------------------------------------------------------------------------
// Idiomatic fixtures, one per layer of the stack
// ---------------------------------------------------------------------------

// ikDouble is an idiomatic Result Kleisli arrow: no config, no context, no IO.
func ikDouble(n int) (int, error) { return n * 2, nil }

// ikFail is an idiomatic Result Kleisli arrow that always fails.
func ikFail(int) (int, error) { return 0, idiomaticTestErr }

// ikDeferredDouble is an idiomatic IOResult Kleisli arrow.
func ikDeferredDouble(n int) func() (int, error) {
	return func() (int, error) { return n * 2, nil }
}

// ikDeferredFail is an idiomatic IOResult Kleisli arrow that always fails.
func ikDeferredFail(int) func() (int, error) {
	return func() (int, error) { return 0, idiomaticTestErr }
}

// ikByConfig is an idiomatic ReaderResult Kleisli arrow: it reads the config.
func ikByConfig(n int) func(AppConfig) (int, error) {
	return func(cfg AppConfig) (int, error) { return n + len(cfg.LogLevel), nil }
}

// ikByConfigFail is an idiomatic ReaderResult Kleisli arrow that always fails.
func ikByConfigFail(int) func(AppConfig) (int, error) {
	return func(AppConfig) (int, error) { return 0, idiomaticTestErr }
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func TestFromResultI(t *testing.T) {
	assert.Equal(t, result.Of(42), runIK(t, FromResultI[AppConfig](42, nil)))
	assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, FromResultI[AppConfig](0, idiomaticTestErr)))
}

func TestFromIOResultI(t *testing.T) {
	assert.Equal(t, result.Of("ok"), runIK(t, FromIOResultI[AppConfig](func() (string, error) { return "ok", nil })))
	assert.Equal(t, result.Left[string](idiomaticTestErr), runIK(t, FromIOResultI[AppConfig](func() (string, error) { return "", idiomaticTestErr })))
}

func TestFromReaderResultI(t *testing.T) {
	t.Run("success case, the config is forwarded", func(t *testing.T) {
		rr := func(cfg AppConfig) (string, error) { return cfg.LogLevel, nil }
		assert.Equal(t, result.Of("info"), runIK(t, FromReaderResultI(rr)))
	})

	t.Run("error case", func(t *testing.T) {
		rr := func(AppConfig) (string, error) { return "", idiomaticTestErr }
		assert.Equal(t, result.Left[string](idiomaticTestErr), runIK(t, FromReaderResultI(rr)))
	})
}

func TestFromReaderIOResultI(t *testing.T) {
	t.Run("success case, the config is forwarded", func(t *testing.T) {
		rr := func(cfg AppConfig) func() (string, error) {
			return func() (string, error) { return cfg.DatabaseURL, nil }
		}
		assert.Equal(t, result.Of("postgres://localhost"), runIK(t, FromReaderIOResultI(rr)))
	})

	t.Run("error case", func(t *testing.T) {
		rr := func(AppConfig) func() (string, error) {
			return func() (string, error) { return "", idiomaticTestErr }
		}
		assert.Equal(t, result.Left[string](idiomaticTestErr), runIK(t, FromReaderIOResultI(rr)))
	})

	t.Run("the IO action is deferred", func(t *testing.T) {
		calls := 0
		rr := func(AppConfig) func() (int, error) {
			return func() (int, error) {
				calls++
				return calls, nil
			}
		}
		eff := FromReaderIOResultI(rr)
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(1), runIK(t, eff))
		assert.Equal(t, result.Of(2), runIK(t, eff))
	})
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

func TestChainResultIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIK(t, F.Pipe1(Of[AppConfig](5), ChainResultIK[AppConfig](ikDouble))))
		assert.Equal(t, result.Of(10), runIK(t, F.Pipe1(Of[AppConfig](5), ChainEitherIK[AppConfig](ikDouble))))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainResultIK[AppConfig](ikFail))))
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainEitherIK[AppConfig](ikFail))))
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) (int, error) {
			called = true
			return n, nil
		}
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Left[AppConfig, int](idiomaticTestErr), ChainResultIK[AppConfig](f))))
		assert.False(t, called)
	})
}

func TestMonadChainResultIK(t *testing.T) {
	assert.Equal(t, result.Of(10), runIK(t, MonadChainResultIK(Of[AppConfig](5), ikDouble)))
	assert.Equal(t, result.Of(10), runIK(t, MonadChainEitherIK(Of[AppConfig](5), ikDouble)))
	assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, MonadChainResultIK(Of[AppConfig](5), ikFail)))
}

func TestChainFirstResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstResultIK[AppConfig](ikDouble))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstEitherIK[AppConfig](ikDouble))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), TapResultIK[AppConfig](ikDouble))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), TapEitherIK[AppConfig](ikDouble))))
		assert.Equal(t, result.Of(5), runIK(t, MonadChainFirstResultIK(Of[AppConfig](5), ikDouble)))
		assert.Equal(t, result.Of(5), runIK(t, MonadChainFirstEitherIK(Of[AppConfig](5), ikDouble)))
		assert.Equal(t, result.Of(5), runIK(t, MonadTapResultIK(Of[AppConfig](5), ikDouble)))
		assert.Equal(t, result.Of(5), runIK(t, MonadTapEitherIK(Of[AppConfig](5), ikDouble)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstResultIK[AppConfig](ikFail))))
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), TapEitherIK[AppConfig](ikFail))))
	})
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

func TestChainIOResultIK(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		assert.Equal(t, result.Of(10), runIK(t, F.Pipe1(Of[AppConfig](5), ChainIOResultIK[AppConfig](ikDeferredDouble))))
		assert.Equal(t, result.Of(10), runIK(t, F.Pipe1(Of[AppConfig](5), ChainIOEitherIK[AppConfig](ikDeferredDouble))))
		assert.Equal(t, result.Of(10), runIK(t, MonadChainIOResultIK(Of[AppConfig](5), ikDeferredDouble)))
		assert.Equal(t, result.Of(10), runIK(t, MonadChainIOEitherIK(Of[AppConfig](5), ikDeferredDouble)))
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainIOResultIK[AppConfig](ikDeferredFail))))
	})

	t.Run("the effect runs once per execution", func(t *testing.T) {
		calls := 0
		f := func(n int) func() (int, error) {
			return func() (int, error) {
				calls++
				return n, nil
			}
		}
		eff := F.Pipe1(Of[AppConfig](5), ChainIOResultIK[AppConfig](f))
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(5), runIK(t, eff))
		assert.Equal(t, 1, calls)
		assert.Equal(t, result.Of(5), runIK(t, eff))
		assert.Equal(t, 2, calls)
	})
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

func TestChainReaderResultIK(t *testing.T) {
	t.Run("success case, the config is forwarded", func(t *testing.T) {
		// 5 + len("info") = 9
		assert.Equal(t, result.Of(9), runIK(t, F.Pipe1(Of[AppConfig](5), ChainReaderResultIK(ikByConfig))))
		assert.Equal(t, result.Of(9), runIK(t, F.Pipe1(Of[AppConfig](5), ChainReaderEitherIK(ikByConfig))))
		assert.Equal(t, result.Of(9), runIK(t, MonadChainReaderResultIK(Of[AppConfig](5), ikByConfig)))
		assert.Equal(t, result.Of(9), runIK(t, MonadChainReaderEitherIK(Of[AppConfig](5), ikByConfig)))
	})

	t.Run("a different config yields a different result", func(t *testing.T) {
		eff := F.Pipe1(Of[AppConfig](5), ChainReaderResultIK(ikByConfig))
		// 5 + len("debug") = 10
		assert.Equal(t, result.Of(10), eff(AppConfig{LogLevel: "debug"})(t.Context())())
	})

	t.Run("error case", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainReaderResultIK(ikByConfigFail))))
	})
}

func TestChainFirstReaderResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstReaderResultIK(ikByConfig))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstReaderEitherIK(ikByConfig))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), TapReaderResultIK(ikByConfig))))
		assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), TapReaderEitherIK(ikByConfig))))
		assert.Equal(t, result.Of(5), runIK(t, MonadChainFirstReaderResultIK(Of[AppConfig](5), ikByConfig)))
		assert.Equal(t, result.Of(5), runIK(t, MonadChainFirstReaderEitherIK(Of[AppConfig](5), ikByConfig)))
		assert.Equal(t, result.Of(5), runIK(t, MonadTapReaderResultIK(Of[AppConfig](5), ikByConfig)))
		assert.Equal(t, result.Of(5), runIK(t, MonadTapReaderEitherIK(Of[AppConfig](5), ikByConfig)))
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), TapReaderResultIK(ikByConfigFail))))
	})
}

// ---------------------------------------------------------------------------
// Option level
// ---------------------------------------------------------------------------

func TestChainOptionIK(t *testing.T) {
	table := map[string]string{"host": "localhost"}
	lookup := func(k string) (string, bool) {
		v, ok := table[k]
		return v, ok
	}
	chain := ChainOptionIK[AppConfig, string, string](func() error { return idiomaticTestErr })

	assert.Equal(t, result.Of("localhost"), runIK(t, F.Pipe1(Of[AppConfig]("host"), chain(lookup))))
	assert.Equal(t, result.Left[string](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig]("port"), chain(lookup))))
}

// ---------------------------------------------------------------------------
// Integration
// ---------------------------------------------------------------------------

// TestIdiomaticKPipeline combines the layer-specific adapters with the
// full-stack [ChainI] in a single pipeline.
func TestIdiomaticKPipeline(t *testing.T) {
	audited := 0

	eff := F.Pipe4(
		FromResultI[AppConfig](2, nil),
		ChainReaderResultIK(ikByConfig),    // 2 + len("info") = 6
		ChainResultIK[AppConfig](ikDouble), // 12
		TapIOResultIK[AppConfig](func(n int) func() (int, error) {
			return func() (int, error) {
				audited = n
				return n, nil
			}
		}),
		ChainI(doubleIdiomatic), // 24
	)

	assert.Equal(t, result.Of(24), runIK(t, eff))
	assert.Equal(t, 12, audited)
}

// TestIdiomaticKPipelineError verifies the pipeline short-circuits on the first error.
func TestIdiomaticKPipelineError(t *testing.T) {
	reached := false

	eff := F.Pipe2(
		FromResultI[AppConfig](2, nil),
		ChainResultIK[AppConfig](ikFail),
		ChainResultIK[AppConfig](func(n int) (int, error) {
			reached = true
			return n, nil
		}),
	)

	assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, eff))
	assert.False(t, reached)
}

func TestChainFirstIOResultIK(t *testing.T) {
	t.Run("the original value is preserved and the effect runs", func(t *testing.T) {
		for name, op := range map[string]func(func(int) func() (int, error)) Operator[AppConfig, int, int]{
			"ChainFirstIOEitherIK": ChainFirstIOEitherIK[AppConfig, int, int],
			"ChainFirstIOResultIK": ChainFirstIOResultIK[AppConfig, int, int],
			"TapIOEitherIK":        TapIOEitherIK[AppConfig, int, int],
			"TapIOResultIK":        TapIOResultIK[AppConfig, int, int],
		} {
			t.Run(name, func(t *testing.T) {
				seen := 0
				f := func(n int) func() (int, error) {
					return func() (int, error) {
						seen = n
						return n * 2, nil
					}
				}
				assert.Equal(t, result.Of(5), runIK(t, F.Pipe1(Of[AppConfig](5), op(f))))
				assert.Equal(t, 5, seen)
			})
		}
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstIOEitherIK[AppConfig](ikDeferredFail))))
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), ChainFirstIOResultIK[AppConfig](ikDeferredFail))))
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), TapIOEitherIK[AppConfig](ikDeferredFail))))
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Of[AppConfig](5), TapIOResultIK[AppConfig](ikDeferredFail))))
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) func() (int, error) {
			return func() (int, error) {
				called = true
				return n, nil
			}
		}
		assert.Equal(t, result.Left[int](idiomaticTestErr), runIK(t, F.Pipe1(Left[AppConfig, int](idiomaticTestErr), ChainFirstIOResultIK[AppConfig](f))))
		assert.False(t, called)
	})
}
