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

package readerreaderioresult

import (
	"context"
	"errors"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

type recConfig struct {
	step  int
	limit int
}

type recState struct {
	n   int
	acc int
}

// recSum adds the numbers from n down to 1, reading the decrement and an
// iteration limit from the outer environment.
func recSum(calls *int) Kleisli[recConfig, recState, Trampoline[recState, int]] {
	return func(s recState) ReaderReaderIOResult[recConfig, Trampoline[recState, int]] {
		return func(cfg recConfig) ReaderIOResult[context.Context, Trampoline[recState, int]] {
			return func(_ context.Context) IOResult[Trampoline[recState, int]] {
				return func() Result[Trampoline[recState, int]] {
					*calls++
					if *calls > cfg.limit {
						return result.Left[Trampoline[recState, int]](errors.New("limit exceeded"))
					}
					if s.n <= 0 {
						return result.Of(tailrec.Land[recState](s.acc))
					}
					return result.Of(tailrec.Bounce[int](recState{s.n - cfg.step, s.acc + s.n}))
				}
			}
		}
	}
}

func TestTailRec_Success(t *testing.T) {
	t.Run("sums numbers using the environment", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		assert.Equal(t, result.Of(15), sum(recState{5, 0})(recConfig{1, 100})(t.Context())())
		assert.Equal(t, 6, calls)
	})

	t.Run("lands immediately", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		assert.Equal(t, result.Of(7), sum(recState{0, 7})(recConfig{1, 100})(t.Context())())
		assert.Equal(t, 1, calls)
	})

	t.Run("step function sees the context", func(t *testing.T) {
		type key struct{}
		ctx := context.WithValue(t.Context(), key{}, 3)

		countdown := TailRec(func(n int) ReaderReaderIOResult[recConfig, Trampoline[int, int]] {
			return func(_ recConfig) ReaderIOResult[context.Context, Trampoline[int, int]] {
				return func(ctx context.Context) IOResult[Trampoline[int, int]] {
					return func() Result[Trampoline[int, int]] {
						if n <= 0 {
							return result.Of(tailrec.Land[int](ctx.Value(key{}).(int)))
						}
						return result.Of(tailrec.Bounce[int](n - 1))
					}
				}
			}
		})

		assert.Equal(t, result.Of(3), countdown(10)(recConfig{})(ctx)())
	})
}

func TestTailRec_Failure(t *testing.T) {
	t.Run("error stops the recursion", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		res := sum(recState{100, 0})(recConfig{1, 3})(t.Context())()

		assert.Equal(t, result.Left[int](errors.New("limit exceeded")), res)
		assert.Equal(t, 4, calls)
	})

	t.Run("cancelled context stops before the first step", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		res := sum(recState{5, 0})(recConfig{1, 100})(ctx)()

		assert.True(t, result.IsLeft(res))
		assert.Equal(t, 0, calls)
	})

	t.Run("cancellation during the recursion stops it", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		calls := 0

		endless := TailRec(func(n int) ReaderReaderIOResult[recConfig, Trampoline[int, int]] {
			return FromIO[recConfig](func() Trampoline[int, int] {
				calls++
				if calls == 5 {
					cancel()
				}
				return tailrec.Bounce[int](n + 1)
			})
		})

		res := endless(0)(recConfig{})(ctx)()

		assert.Equal(t, result.Left[int](context.Canceled), res)
		assert.Equal(t, 5, calls)
	})

	t.Run("timeout stops the recursion with its cause", func(t *testing.T) {
		cause := errors.New("too slow")
		ctx, cancel := context.WithTimeoutCause(t.Context(), 20*time.Millisecond, cause)
		defer cancel()

		endless := TailRec(func(n int) ReaderReaderIOResult[recConfig, Trampoline[int, int]] {
			return F.Pipe1(
				Of[recConfig](tailrec.Bounce[int](n+1)),
				Delay[recConfig, Trampoline[int, int]](time.Millisecond),
			)
		})

		assert.Equal(t, result.Left[int](cause), endless(0)(recConfig{})(ctx)())
	})
}

func TestTailRec_EdgeCases(t *testing.T) {
	t.Run("is stack safe", func(t *testing.T) {
		countdown := TailRec(func(n int) ReaderReaderIOResult[recConfig, Trampoline[int, int]] {
			if n <= 0 {
				return Of[recConfig](tailrec.Land[int](0))
			}
			return Of[recConfig](tailrec.Bounce[int](n - 1))
		})

		assert.Equal(t, result.Of(0), countdown(1_000_000)(recConfig{})(t.Context())())
	})

	t.Run("does not call the step function before the IO runs", func(t *testing.T) {
		calls := 0
		io := TailRec(recSum(&calls))(recState{3, 0})(recConfig{1, 100})(t.Context())

		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(6), io())
		assert.Equal(t, 4, calls)
	})
}

func TestTailRec_Integration(t *testing.T) {
	t.Run("composes with Chain and Map", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		program := F.Pipe2(
			Asks(func(cfg recConfig) int { return cfg.step * 4 }),
			Chain(func(n int) ReaderReaderIOResult[recConfig, int] {
				return sum(recState{n, 0})
			}),
			Map[recConfig](func(n int) int { return n * 10 }),
		)

		assert.Equal(t, result.Of(100), program(recConfig{1, 100})(t.Context())())
	})
}
