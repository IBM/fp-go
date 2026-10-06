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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

type recOuter struct {
	step int
}

type recInner struct {
	limit int
}

type recState struct {
	n   int
	acc int
}

// recSum adds the numbers from n down to 1, reading the decrement from the outer
// and the iteration limit from the inner environment.
func recSum(calls *int) Kleisli[recOuter, recInner, error, recState, Trampoline[recState, int]] {
	return func(s recState) ReaderReaderIOEither[recOuter, recInner, error, Trampoline[recState, int]] {
		return func(o recOuter) ReaderIOEither[recInner, error, Trampoline[recState, int]] {
			return func(i recInner) IOEither[error, Trampoline[recState, int]] {
				return func() Either[error, Trampoline[recState, int]] {
					*calls++
					if *calls > i.limit {
						return E.Left[Trampoline[recState, int]](errors.New("limit exceeded"))
					}
					if s.n <= 0 {
						return E.Right[error](tailrec.Land[recState](s.acc))
					}
					return E.Right[error](tailrec.Bounce[int](recState{s.n - o.step, s.acc + s.n}))
				}
			}
		}
	}
}

func TestTailRec_Success(t *testing.T) {
	t.Run("sums numbers using both environments", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		assert.Equal(t, E.Right[error](15), sum(recState{5, 0})(recOuter{1})(recInner{100})())
		assert.Equal(t, 6, calls)
	})

	t.Run("outer environment changes the step", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		// 10 + 8 + 6 + 4 + 2
		assert.Equal(t, E.Right[error](30), sum(recState{10, 0})(recOuter{2})(recInner{100})())
	})

	t.Run("lands immediately", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		assert.Equal(t, E.Right[error](7), sum(recState{0, 7})(recOuter{1})(recInner{100})())
		assert.Equal(t, 1, calls)
	})
}

func TestTailRec_Failure(t *testing.T) {
	t.Run("error stops the recursion", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		res := sum(recState{100, 0})(recOuter{1})(recInner{3})()

		assert.Equal(t, E.Left[int](errors.New("limit exceeded")), res)
		assert.Equal(t, 4, calls)
	})

	t.Run("error in the first step", func(t *testing.T) {
		failing := TailRec(func(n int) ReaderReaderIOEither[recOuter, recInner, string, Trampoline[int, int]] {
			return Left[recOuter, recInner, Trampoline[int, int]]("boom")
		})

		assert.Equal(t, E.Left[int]("boom"), failing(1)(recOuter{})(recInner{})())
	})
}

func TestTailRec_EdgeCases(t *testing.T) {
	t.Run("is stack safe", func(t *testing.T) {
		countdown := TailRec(func(n int) ReaderReaderIOEither[recOuter, recInner, error, Trampoline[int, int]] {
			if n <= 0 {
				return Of[recOuter, recInner, error](tailrec.Land[int](0))
			}
			return Of[recOuter, recInner, error](tailrec.Bounce[int](n - 1))
		})

		assert.Equal(t, E.Right[error](0), countdown(1_000_000)(recOuter{})(recInner{})())
	})

	t.Run("does not call the step function before the IO runs", func(t *testing.T) {
		calls := 0
		io := TailRec(recSum(&calls))(recState{3, 0})(recOuter{1})(recInner{100})

		assert.Equal(t, 0, calls)
		assert.Equal(t, E.Right[error](6), io())
		assert.Equal(t, 4, calls)
	})

	t.Run("running the IO twice repeats the recursion", func(t *testing.T) {
		calls := 0
		io := TailRec(recSum(&calls))(recState{3, 0})(recOuter{1})(recInner{100})

		assert.Equal(t, io(), io())
		assert.Equal(t, 8, calls)
	})
}

func TestTailRec_Integration(t *testing.T) {
	t.Run("composes with Chain and Map", func(t *testing.T) {
		calls := 0
		sum := TailRec(recSum(&calls))

		program := F.Pipe2(
			Asks[recInner, error](func(o recOuter) int { return o.step * 4 }),
			Chain(func(n int) ReaderReaderIOEither[recOuter, recInner, error, int] {
				return sum(recState{n, 0})
			}),
			Map[recOuter, recInner, error](func(n int) int { return n * 10 }),
		)

		assert.Equal(t, E.Right[error](100), program(recOuter{1})(recInner{100})())
	})
}
