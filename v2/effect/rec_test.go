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

package effect

import (
	"context"
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

type recConfig struct {
	Step  int
	Limit int
}

type recState struct {
	n   int
	acc int
}

// recSum adds the numbers from n down to 1, reading the decrement and an
// iteration limit from the context.
func recSum(calls *int) Kleisli[recConfig, recState, Trampoline[recState, int]] {
	return func(s recState) Effect[recConfig, Trampoline[recState, int]] {
		return F.Pipe1(
			Asks(F.Identity[recConfig]),
			Chain(func(cfg recConfig) Effect[recConfig, Trampoline[recState, int]] {
				*calls++
				if *calls > cfg.Limit {
					return Fail[recConfig, Trampoline[recState, int]](errors.New("limit exceeded"))
				}
				if s.n <= 0 {
					return Of[recConfig](tailrec.Land[recState](s.acc))
				}
				return Of[recConfig](tailrec.Bounce[int](recState{s.n - cfg.Step, s.acc + s.n}))
			}),
		)
	}
}

func TestTailRec_Success(t *testing.T) {
	t.Run("sums numbers using the context", func(t *testing.T) {
		calls := 0
		res, err := runEffect(TailRec(recSum(&calls))(recState{5, 0}), recConfig{1, 100})

		assert.NoError(t, err)
		assert.Equal(t, 15, res)
		assert.Equal(t, 6, calls)
	})

	t.Run("context changes the step", func(t *testing.T) {
		calls := 0
		res, err := runEffect(TailRec(recSum(&calls))(recState{10, 0}), recConfig{2, 100})

		assert.NoError(t, err)
		assert.Equal(t, 30, res)
	})

	t.Run("lands immediately", func(t *testing.T) {
		calls := 0
		res, err := runEffect(TailRec(recSum(&calls))(recState{0, 7}), recConfig{1, 100})

		assert.NoError(t, err)
		assert.Equal(t, 7, res)
		assert.Equal(t, 1, calls)
	})
}

func TestTailRec_Failure(t *testing.T) {
	t.Run("error stops the recursion", func(t *testing.T) {
		calls := 0
		_, err := runEffect(TailRec(recSum(&calls))(recState{100, 0}), recConfig{1, 3})

		assert.EqualError(t, err, "limit exceeded")
		assert.Equal(t, 4, calls)
	})

	t.Run("cancelled context stops the recursion", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		calls := 0

		endless := TailRec(func(n int) Effect[recConfig, Trampoline[int, int]] {
			return FromIO[recConfig](func() Trampoline[int, int] {
				calls++
				if calls == 5 {
					cancel()
				}
				return tailrec.Bounce[int](n + 1)
			})
		})

		_, err := RunSync(Provide[int](recConfig{})(endless(0)))(ctx)

		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 5, calls)
	})
}

func TestTailRec_EdgeCases(t *testing.T) {
	t.Run("is stack safe", func(t *testing.T) {
		countdown := TailRec(func(n int) Effect[recConfig, Trampoline[int, int]] {
			if n <= 0 {
				return Of[recConfig](tailrec.Land[int](0))
			}
			return Of[recConfig](tailrec.Bounce[int](n - 1))
		})

		res, err := runEffect(countdown(1_000_000), recConfig{})

		assert.NoError(t, err)
		assert.Equal(t, 0, res)
	})

	t.Run("does not run before the effect is run", func(t *testing.T) {
		calls := 0
		eff := TailRec(recSum(&calls))(recState{3, 0})

		assert.Equal(t, 0, calls)

		res, err := runEffect(eff, recConfig{1, 100})

		assert.NoError(t, err)
		assert.Equal(t, 6, res)
		assert.Equal(t, 4, calls)
	})
}

func TestTailRec_Integration(t *testing.T) {
	t.Run("composes with Chain and Map", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		program := F.Pipe2(
			Asks(func(cfg recConfig) recState { return recState{cfg.Step * 4, 0} }),
			Chain(sum),
			Map[recConfig](func(n int) int { return n * 10 }),
		)

		res, err := runEffect(program, recConfig{1, 100})

		assert.NoError(t, err)
		assert.Equal(t, 100, res)
	})
}
