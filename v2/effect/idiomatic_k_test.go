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

package effect

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Idiomatic fixtures, one per layer of the stack
// ---------------------------------------------------------------------------

// ikDouble is an idiomatic Result Kleisli arrow: it needs nothing but the value.
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

// ikByEnv is an idiomatic ReaderResult Kleisli arrow: it reads the environment.
func ikByEnv(n int) func(TestContext) (int, error) {
	return func(c TestContext) (int, error) { return n + len(c.Value), nil }
}

// ikByEnvFail is an idiomatic ReaderResult Kleisli arrow that always fails.
func ikByEnvFail(int) func(TestContext) (int, error) {
	return func(TestContext) (int, error) { return 0, idiomaticTestErr }
}

// ikThunk is an idiomatic Thunk Kleisli arrow: it sees the context and performs
// IO, but does not read the environment.
func ikThunk(n int) func(context.Context) func() (int, error) {
	return func(ctx context.Context) func() (int, error) {
		return func() (int, error) {
			if err := ctx.Err(); err != nil {
				return 0, err
			}
			return n * 10, nil
		}
	}
}

// ikThunkFail is an idiomatic Thunk Kleisli arrow that always fails.
func ikThunkFail(int) func(context.Context) func() (int, error) {
	return func(context.Context) func() (int, error) {
		return func() (int, error) { return 0, idiomaticTestErr }
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func TestFromResultI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v, err := runEffect(FromResultI[TestContext](42, nil), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(FromResultI[TestContext](0, idiomaticTestErr), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

func TestFromIOResultI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v, err := runEffect(FromIOResultI[TestContext](func() (string, error) { return "ok", nil }), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, "ok", v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(FromIOResultI[TestContext](func() (string, error) { return "", idiomaticTestErr }), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

func TestFromReaderResultI(t *testing.T) {
	t.Run("success, the environment is forwarded", func(t *testing.T) {
		rr := func(c TestContext) (string, error) { return c.Value, nil }
		v, err := runEffect(FromReaderResultI(rr), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, "test", v)
	})

	t.Run("error", func(t *testing.T) {
		rr := func(TestContext) (string, error) { return "", idiomaticTestErr }
		_, err := runEffect(FromReaderResultI(rr), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

func TestFromThunkI(t *testing.T) {
	t.Run("success, the context is forwarded", func(t *testing.T) {
		f := func(ctx context.Context) func() (string, error) {
			return func() (string, error) { return "ok", ctx.Err() }
		}
		v, err := runEffect(FromThunkI[TestContext](f), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, "ok", v)
	})

	t.Run("error", func(t *testing.T) {
		f := func(context.Context) func() (string, error) {
			return func() (string, error) { return "", idiomaticTestErr }
		}
		_, err := runEffect(FromThunkI[TestContext](f), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

func TestChainResultIK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainResultIK[TestContext](ikDouble)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainResultIK[TestContext](ikFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) (int, error) {
			called = true
			return n, nil
		}
		_, err := runEffect(F.Pipe1(Fail[TestContext, int](idiomaticTestErr), ChainResultIK[TestContext](f)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
		assert.False(t, called)
	})

	t.Run("uncurried form", func(t *testing.T) {
		v, err := runEffect(MonadChainResultIK(Of[TestContext](5), ikDouble), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})
}

func TestChainFirstResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		for _, eff := range []Effect[TestContext, int]{
			F.Pipe1(Of[TestContext](5), ChainFirstResultIK[TestContext](ikDouble)),
			F.Pipe1(Of[TestContext](5), TapResultIK[TestContext](ikDouble)),
			MonadChainFirstResultIK(Of[TestContext](5), ikDouble),
			MonadTapResultIK(Of[TestContext](5), ikDouble),
		} {
			v, err := runEffect(eff, testCtx)
			assert.NoError(t, err)
			assert.Equal(t, 5, v)
		}
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), TapResultIK[TestContext](ikFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

func TestChainReaderResultIK(t *testing.T) {
	t.Run("success, the environment is forwarded", func(t *testing.T) {
		// 5 + len("test") = 9
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainReaderResultIK(ikByEnv)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 9, v)
	})

	t.Run("a different environment yields a different result", func(t *testing.T) {
		// 5 + len("longer") = 11
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainReaderResultIK(ikByEnv)), TestContext{Value: "longer"})
		assert.NoError(t, err)
		assert.Equal(t, 11, v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainReaderResultIK(ikByEnvFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("uncurried form", func(t *testing.T) {
		v, err := runEffect(MonadChainReaderResultIK(Of[TestContext](5), ikByEnv), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 9, v)
	})
}

func TestChainFirstReaderResultIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		for _, eff := range []Effect[TestContext, int]{
			F.Pipe1(Of[TestContext](5), ChainFirstReaderResultIK(ikByEnv)),
			F.Pipe1(Of[TestContext](5), TapReaderResultIK(ikByEnv)),
			MonadChainFirstReaderResultIK(Of[TestContext](5), ikByEnv),
			MonadTapReaderResultIK(Of[TestContext](5), ikByEnv),
		} {
			v, err := runEffect(eff, testCtx)
			assert.NoError(t, err)
			assert.Equal(t, 5, v)
		}
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), TapReaderResultIK(ikByEnvFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

func TestChainIOResultIK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainIOResultIK[TestContext](ikDeferredDouble)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainIOResultIK[TestContext](ikDeferredFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("uncurried form", func(t *testing.T) {
		v, err := runEffect(MonadChainIOResultIK(Of[TestContext](5), ikDeferredDouble), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("the effect runs once per execution", func(t *testing.T) {
		calls := 0
		f := func(n int) func() (int, error) {
			return func() (int, error) {
				calls++
				return n, nil
			}
		}
		eff := F.Pipe1(Of[TestContext](5), ChainIOResultIK[TestContext](f))
		assert.Equal(t, 0, calls)
		_, _ = runEffect(eff, testCtx)
		assert.Equal(t, 1, calls)
		_, _ = runEffect(eff, testCtx)
		assert.Equal(t, 2, calls)
	})
}

// ---------------------------------------------------------------------------
// Thunk level
// ---------------------------------------------------------------------------

func TestChainThunkIK(t *testing.T) {
	t.Run("success, the context is available", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainThunkIK[TestContext](ikThunk)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 50, v)
	})

	t.Run("error", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainThunkIK[TestContext](ikThunkFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("uncurried form", func(t *testing.T) {
		v, err := runEffect(MonadChainThunkIK(Of[TestContext](5), ikThunk), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 50, v)
	})

	t.Run("a cancelled context surfaces as an error", func(t *testing.T) {
		cancelled, cancel := context.WithCancel(context.Background())
		cancel()

		eff := F.Pipe1(Of[TestContext](5), ChainThunkIK[TestContext](ikThunk))
		_, err := RunSync(Provide[int](testCtx)(eff))(cancelled)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestChainFirstThunkIK(t *testing.T) {
	t.Run("the original value is preserved", func(t *testing.T) {
		for _, eff := range []Effect[TestContext, int]{
			F.Pipe1(Of[TestContext](5), ChainFirstThunkIK[TestContext, int](ikThunk)),
			F.Pipe1(Of[TestContext](5), TapThunkIK[TestContext, int](ikThunk)),
		} {
			v, err := runEffect(eff, testCtx)
			assert.NoError(t, err)
			assert.Equal(t, 5, v)
		}
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), TapThunkIK[TestContext, int](ikThunkFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

func TestChainFirstLeftThunkIK(t *testing.T) {
	t.Run("the handler runs on the error path and the error is preserved", func(t *testing.T) {
		var seen error
		report := func(err error) func(context.Context) func() (string, error) {
			return func(context.Context) func() (string, error) {
				return func() (string, error) {
					seen = err
					return "reported", nil
				}
			}
		}

		_, err := runEffect(F.Pipe1(Fail[TestContext, int](idiomaticTestErr), ChainFirstLeftThunkIK[TestContext, int](report)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
		assert.Equal(t, idiomaticTestErr, seen)
	})

	t.Run("a successful value is untouched and the handler is not called", func(t *testing.T) {
		called := false
		report := func(error) func(context.Context) func() (string, error) {
			return func(context.Context) func() (string, error) {
				return func() (string, error) {
					called = true
					return "", nil
				}
			}
		}

		v, err := runEffect(F.Pipe1(Of[TestContext](5), TapLeftThunkIK[TestContext, int](report)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 5, v)
		assert.False(t, called)
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
	chain := ChainOptionIK[TestContext, string, string](func() error { return idiomaticTestErr })

	t.Run("the comma-ok pair is lifted into a success", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext]("host"), chain(lookup)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", v)
	})

	t.Run("a false flag becomes the error from onNone", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext]("port"), chain(lookup)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

// ---------------------------------------------------------------------------
// Integration
// ---------------------------------------------------------------------------

// TestIdiomaticKPipeline combines the layer-specific adapters with the
// full-stack ChainI in a single pipeline.
func TestIdiomaticKPipeline(t *testing.T) {
	audited := 0

	eff := F.Pipe4(
		FromResultI[TestContext](2, nil),
		ChainReaderResultIK(ikByEnv),         // 2 + len("test") = 6
		ChainResultIK[TestContext](ikDouble), // 12
		TapIOResultIK[TestContext](func(n int) func() (int, error) {
			return func() (int, error) {
				audited = n
				return n, nil
			}
		}),
		ChainI(doubleIdiomatic), // 24
	)

	v, err := runEffect(eff, testCtx)
	assert.NoError(t, err)
	assert.Equal(t, 24, v)
	assert.Equal(t, 12, audited)
}

// TestIdiomaticKPipelineError verifies the pipeline short-circuits on the first error.
func TestIdiomaticKPipelineError(t *testing.T) {
	reached := false

	eff := F.Pipe2(
		FromResultI[TestContext](2, nil),
		ChainResultIK[TestContext](ikFail),
		ChainResultIK[TestContext](func(n int) (int, error) {
			reached = true
			return n, nil
		}),
	)

	_, err := runEffect(eff, testCtx)
	assert.Equal(t, idiomaticTestErr, err)
	assert.False(t, reached)
}

func TestChainFirstIOResultIK(t *testing.T) {
	t.Run("the original value is preserved and the effect runs", func(t *testing.T) {
		for name, op := range map[string]func(func(int) func() (int, error)) Operator[TestContext, int, int]{
			"ChainFirstIOResultIK": ChainFirstIOResultIK[TestContext, int, int],
			"TapIOResultIK":        TapIOResultIK[TestContext, int, int],
		} {
			t.Run(name, func(t *testing.T) {
				seen := 0
				f := func(n int) func() (int, error) {
					return func() (int, error) {
						seen = n
						return n * 2, nil
					}
				}
				v, err := runEffect(F.Pipe1(Of[TestContext](5), op(f)), testCtx)
				assert.NoError(t, err)
				assert.Equal(t, 5, v)
				assert.Equal(t, 5, seen)
			})
		}
	})

	t.Run("a failure of the side effect is propagated", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainFirstIOResultIK[TestContext](ikDeferredFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
		_, err = runEffect(F.Pipe1(Of[TestContext](5), TapIOResultIK[TestContext](ikDeferredFail)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("the source error is propagated, f is not called", func(t *testing.T) {
		called := false
		f := func(n int) func() (int, error) {
			return func() (int, error) {
				called = true
				return n, nil
			}
		}
		_, err := runEffect(F.Pipe1(Fail[TestContext, int](idiomaticTestErr), ChainFirstIOResultIK[TestContext](f)), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
		assert.False(t, called)
	})
}
