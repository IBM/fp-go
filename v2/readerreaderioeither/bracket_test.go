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

package readerreaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// bracketEnv is the outer environment used by the Bracket tests
type bracketEnv struct {
	pool string
}

// bracketCtx is the inner environment used by the Bracket tests
type bracketCtx struct {
	timeout int
}

// bracketRes is the resource managed by the Bracket tests
type bracketRes struct {
	id string
}

// bracketEvent records the execution of a single step together with the environments it saw
type bracketEvent struct {
	name string
	env  bracketEnv
	ctx  bracketCtx
}

// bracketRecorder collects the events of all steps in execution order
type bracketRecorder struct {
	events []bracketEvent
}

func (rec *bracketRecorder) names() []string {
	names := make([]string, len(rec.events))
	for i, ev := range rec.events {
		names[i] = ev.name
	}
	return names
}

// recordStep returns a computation that records its execution under the given name and yields res
func recordStep[A any](rec *bracketRecorder, name string, res Either[error, A]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, A] {
	return func(env bracketEnv) ReaderIOEither[bracketCtx, error, A] {
		return func(ctx bracketCtx) IOEither[error, A] {
			return func() Either[error, A] {
				rec.events = append(rec.events, bracketEvent{name: name, env: env, ctx: ctx})
				return res
			}
		}
	}
}

var (
	testBracketEnv = bracketEnv{pool: "pool1"}
	testBracketCtx = bracketCtx{timeout: 30}
)

// TestBracketOutcomes covers every combination of use and release succeeding or failing
func TestBracketOutcomes(t *testing.T) {
	useErr := errors.New("use failed")
	releaseErr := errors.New("release failed")

	tests := []struct {
		name          string
		useResult     Either[error, string]
		releaseResult Either[error, F.Void]
		expected      Either[error, string]
	}{
		{
			name:          "use and release succeed",
			useResult:     E.Right[error]("used"),
			releaseResult: E.Right[error](F.VOID),
			expected:      E.Right[error]("used"),
		},
		{
			name:          "use fails, release succeeds",
			useResult:     E.Left[string](useErr),
			releaseResult: E.Right[error](F.VOID),
			expected:      E.Left[string](useErr),
		},
		{
			name:          "use succeeds, release fails",
			useResult:     E.Right[error]("used"),
			releaseResult: E.Left[F.Void](releaseErr),
			expected:      E.Left[string](releaseErr),
		},
		{
			name:          "use and release fail, release error wins",
			useResult:     E.Left[string](useErr),
			releaseResult: E.Left[F.Void](releaseErr),
			expected:      E.Left[string](releaseErr),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &bracketRecorder{}
			resource := bracketRes{id: "res1"}

			var (
				usedWith     bracketRes
				releasedWith bracketRes
				outcome      Either[error, string]
			)

			use := func(r bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
				usedWith = r
				return recordStep(rec, "use", tt.useResult)
			}

			release := func(r bracketRes, result Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
				releasedWith = r
				outcome = result
				return recordStep(rec, "release", tt.releaseResult)
			}

			program := Bracket(recordStep(rec, "acquire", E.Right[error](resource)), use, release)

			assert.Equal(t, tt.expected, program(testBracketEnv)(testBracketCtx)())
			assert.Equal(t, []string{"acquire", "use", "release"}, rec.names())
			assert.Equal(t, resource, usedWith, "use must receive the acquired resource")
			assert.Equal(t, resource, releasedWith, "release must receive the acquired resource")
			assert.Equal(t, tt.useResult, outcome, "release must receive the outcome of use")
		})
	}
}

// TestBracketAcquireFailure verifies that neither use nor release run if acquire fails
func TestBracketAcquireFailure(t *testing.T) {
	rec := &bracketRecorder{}
	acquireErr := errors.New("acquire failed")

	use := func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		return recordStep(rec, "use", E.Right[error]("used"))
	}

	release := func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
		return recordStep(rec, "release", E.Right[error](F.VOID))
	}

	program := Bracket(recordStep(rec, "acquire", E.Left[bracketRes](acquireErr)), use, release)

	assert.Equal(t, E.Left[string](acquireErr), program(testBracketEnv)(testBracketCtx)())
	assert.Equal(t, []string{"acquire"}, rec.names())
}

// TestBracketEnvironments verifies that acquire, use and release all see both environments
func TestBracketEnvironments(t *testing.T) {
	rec := &bracketRecorder{}
	env := bracketEnv{pool: "production"}
	ctx := bracketCtx{timeout: 60}

	// acquire derives the resource from the outer environment
	acquire := F.Pipe1(
		Asks[bracketCtx, error](func(e bracketEnv) bracketRes {
			return bracketRes{id: e.pool + "-resource"}
		}),
		ChainFirst(func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
			return recordStep(rec, "acquire", E.Right[error](F.VOID))
		}),
	)

	use := func(r bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		return recordStep(rec, "use", E.Right[error](r.id))
	}

	release := func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
		return recordStep(rec, "release", E.Right[error](F.VOID))
	}

	result := Bracket(acquire, use, release)(env)(ctx)()

	assert.Equal(t, E.Right[error]("production-resource"), result)
	assert.Equal(t, []bracketEvent{
		{name: "acquire", env: env, ctx: ctx},
		{name: "use", env: env, ctx: ctx},
		{name: "release", env: env, ctx: ctx},
	}, rec.events)
}

// TestBracketIsLazy verifies that no step runs before the IO is executed and that every
// execution acquires and releases a fresh resource
func TestBracketIsLazy(t *testing.T) {
	rec := &bracketRecorder{}

	use := func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		return recordStep(rec, "use", E.Right[error]("used"))
	}

	release := func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
		return recordStep(rec, "release", E.Right[error](F.VOID))
	}

	io := Bracket(recordStep(rec, "acquire", E.Right[error](bracketRes{id: "res1"})), use, release)(testBracketEnv)(testBracketCtx)

	assert.Empty(t, rec.events, "no step must run before the IO is executed")

	assert.Equal(t, E.Right[error]("used"), io())
	assert.Equal(t, E.Right[error]("used"), io())
	assert.Equal(t, []string{"acquire", "use", "release", "acquire", "use", "release"}, rec.names())
}

// TestBracketNested verifies that nested brackets release their resources in reverse order
// of acquisition, also if the innermost step fails
func TestBracketNested(t *testing.T) {
	expected := []string{"acquire outer", "acquire inner", "use", "release inner", "release outer"}

	tests := []struct {
		name      string
		useResult Either[error, string]
	}{
		{name: "success", useResult: E.Right[error]("used")},
		{name: "failure", useResult: E.Left[string](errors.New("inner failed"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &bracketRecorder{}

			releaseAs := func(name string) func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
				return func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
					return recordStep(rec, name, E.Right[error](F.VOID))
				}
			}

			useOuter := func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
				return Bracket(
					recordStep(rec, "acquire inner", E.Right[error](bracketRes{id: "inner"})),
					func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
						return recordStep(rec, "use", tt.useResult)
					},
					releaseAs("release inner"),
				)
			}

			program := Bracket(
				recordStep(rec, "acquire outer", E.Right[error](bracketRes{id: "outer"})),
				useOuter,
				releaseAs("release outer"),
			)

			assert.Equal(t, tt.useResult, program(testBracketEnv)(testBracketCtx)())
			assert.Equal(t, expected, rec.names())
		})
	}
}

// TestBracketNestedInnerAcquireFailure verifies that the outer resource is released if the
// inner resource cannot be acquired, while the inner release does not run
func TestBracketNestedInnerAcquireFailure(t *testing.T) {
	rec := &bracketRecorder{}
	acquireErr := errors.New("inner acquire failed")

	useOuter := func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		return Bracket(
			recordStep(rec, "acquire inner", E.Left[bracketRes](acquireErr)),
			func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
				return recordStep(rec, "use", E.Right[error]("used"))
			},
			func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
				return recordStep(rec, "release inner", E.Right[error](F.VOID))
			},
		)
	}

	var outerOutcome Either[error, string]

	program := Bracket(
		recordStep(rec, "acquire outer", E.Right[error](bracketRes{id: "outer"})),
		useOuter,
		func(_ bracketRes, result Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
			outerOutcome = result
			return recordStep(rec, "release outer", E.Right[error](F.VOID))
		},
	)

	assert.Equal(t, E.Left[string](acquireErr), program(testBracketEnv)(testBracketCtx)())
	assert.Equal(t, []string{"acquire outer", "acquire inner", "release outer"}, rec.names())
	assert.Equal(t, E.Left[string](acquireErr), outerOutcome)
}

// TestBracketReleaseReceivesOutcome verifies that release receives the exact Either[E, B]
// that use produced — Right with the computed value on success, Left with the error on failure.
func TestBracketReleaseReceivesOutcome(t *testing.T) {
	t.Run("release receives Right outcome when use succeeds", func(t *testing.T) {
		var capturedOutcome Either[error, string]
		resource := bracketRes{id: "r1"}

		use := func(r bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			return Right[bracketEnv, bracketCtx, error]("ok")
		}

		release := func(r bracketRes, outcome Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
			capturedOutcome = outcome
			return Right[bracketEnv, bracketCtx, error](F.VOID)
		}

		program := Bracket(Right[bracketEnv, bracketCtx, error](resource), use, release)
		program(testBracketEnv)(testBracketCtx)()

		assert.Equal(t, E.Right[error]("ok"), capturedOutcome)
	})

	t.Run("release receives Left outcome when use fails", func(t *testing.T) {
		var capturedOutcome Either[error, string]
		useErr := errors.New("use error")
		resource := bracketRes{id: "r1"}

		use := func(r bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			return Left[bracketEnv, bracketCtx, string](useErr)
		}

		release := func(r bracketRes, outcome Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
			capturedOutcome = outcome
			return Right[bracketEnv, bracketCtx, error](F.VOID)
		}

		program := Bracket(Right[bracketEnv, bracketCtx, error](resource), use, release)
		program(testBracketEnv)(testBracketCtx)()

		assert.Equal(t, E.Left[string](useErr), capturedOutcome)
	})
}

// TestBracketResourceIdentity verifies that use and release receive the identical resource
// value that acquire produced (pointer/value equality, not just structural equality).
func TestBracketResourceIdentity(t *testing.T) {
	type heap struct{ id int }

	acquired := &heap{id: 42}
	var usedPtr, releasedPtr *heap

	acquire := Right[bracketEnv, bracketCtx, error](acquired)

	use := func(r *heap) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		usedPtr = r
		return Right[bracketEnv, bracketCtx, error]("done")
	}

	release := func(r *heap, _ Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
		releasedPtr = r
		return Right[bracketEnv, bracketCtx, error](F.VOID)
	}

	Bracket(acquire, use, release)(testBracketEnv)(testBracketCtx)()

	assert.Same(t, acquired, usedPtr, "use must receive the exact resource pointer from acquire")
	assert.Same(t, acquired, releasedPtr, "release must receive the exact resource pointer from acquire")
}

// TestBracketReleaseRunsBeforeResultPropagates verifies that the release step is fully
// executed before the final result is returned to the caller.
func TestBracketReleaseRunsBeforeResultPropagates(t *testing.T) {
	var order []string

	use := func(bracketRes) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
		return func(env bracketEnv) ReaderIOEither[bracketCtx, error, string] {
			return func(ctx bracketCtx) IOEither[error, string] {
				return func() Either[error, string] {
					order = append(order, "use")
					return E.Right[error]("result")
				}
			}
		}
	}

	release := func(bracketRes, Either[error, string]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, F.Void] {
		return func(env bracketEnv) ReaderIOEither[bracketCtx, error, F.Void] {
			return func(ctx bracketCtx) IOEither[error, F.Void] {
				return func() Either[error, F.Void] {
					order = append(order, "release")
					return E.Right[error](F.VOID)
				}
			}
		}
	}

	program := Bracket(Right[bracketEnv, bracketCtx, error](bracketRes{id: "r"}), use, release)
	result := program(testBracketEnv)(testBracketCtx)()

	assert.Equal(t, []string{"use", "release"}, order, "release must run before result is returned")
	assert.Equal(t, E.Right[error]("result"), result)
}

// TestMonadChainReaderReaderIO verifies that monadChainReaderReaderIO threads the full Either[E, A] into the
// continuation — it must run the continuation for both Right and Left values.
func TestMonadChainReaderReaderIO(t *testing.T) {
	t.Run("continuation receives Right value", func(t *testing.T) {
		var received Either[error, int]

		fa := Right[bracketEnv, bracketCtx, error](42)
		result := monadChainReaderReaderIO(fa, func(e Either[error, int]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			received = e
			return Right[bracketEnv, bracketCtx, error]("ok")
		})

		assert.Equal(t, E.Right[error]("ok"), result(testBracketEnv)(testBracketCtx)())
		assert.Equal(t, E.Right[error](42), received)
	})

	t.Run("continuation receives Left value — runs even when fa is Left", func(t *testing.T) {
		err := errors.New("original")
		var received Either[error, int]

		fa := Left[bracketEnv, bracketCtx, int](err)
		result := monadChainReaderReaderIO(fa, func(e Either[error, int]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			received = e
			return Right[bracketEnv, bracketCtx, error]("recovered")
		})

		assert.Equal(t, E.Right[error]("recovered"), result(testBracketEnv)(testBracketCtx)())
		assert.Equal(t, E.Left[int](err), received)
	})

	t.Run("continuation can itself return Left", func(t *testing.T) {
		contErr := errors.New("cont error")
		fa := Right[bracketEnv, bracketCtx, error](1)

		result := monadChainReaderReaderIO(fa, func(_ Either[error, int]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			return Left[bracketEnv, bracketCtx, string](contErr)
		})

		assert.Equal(t, E.Left[string](contErr), result(testBracketEnv)(testBracketCtx)())
	})

	t.Run("continuation receives both environments", func(t *testing.T) {
		env := bracketEnv{pool: "prod"}
		ctx := bracketCtx{timeout: 99}

		fa := Right[bracketEnv, bracketCtx, error](0)
		result := monadChainReaderReaderIO(fa, func(_ Either[error, int]) ReaderReaderIOEither[bracketEnv, bracketCtx, error, string] {
			return func(e bracketEnv) ReaderIOEither[bracketCtx, error, string] {
				return func(c bracketCtx) IOEither[error, string] {
					return Right[struct{}, struct{}, error](
						e.pool + ":" + fmt.Sprintf("%d", c.timeout),
					)(struct{}{})(struct{}{})
				}
			}
		})

		assert.Equal(t, E.Right[error]("prod:99"), result(env)(ctx)())
	})
}
