// Copyright (c) 2024 - 2025 IBM Corp.
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

package stateio

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

type TestState struct {
	Counter int
	Message string
}

func TestOf(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}
	computation := Of[TestState](42)

	result := computation(initial)()

	assert.Equal(t, 42, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestMonadMap(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}
	computation := Of[TestState](10)

	mapped := MonadMap(computation, N.Mul(2))
	result := mapped(initial)()

	assert.Equal(t, 20, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestMap(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	double := Map[TestState](N.Mul(2))
	computation := F.Pipe1(
		Of[TestState](21),
		double,
	)

	result := computation(initial)()

	assert.Equal(t, 42, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestMonadChain(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	computation := Of[TestState](5)
	chained := MonadChain(computation, func(x int) StateIO[TestState, string] {
		return func(s TestState) IO[Pair[TestState, string]] {
			return func() Pair[TestState, string] {
				newState := TestState{Counter: s.Counter + x, Message: fmt.Sprintf("value: %d", x)}
				return pair.MakePair(newState, fmt.Sprintf("result: %d", x*2))
			}
		}
	})

	result := chained(initial)()

	assert.Equal(t, "result: 10", pair.Tail(result))
	assert.Equal(t, 5, pair.Head(result).Counter)
	assert.Equal(t, "value: 5", pair.Head(result).Message)
}

func TestChain(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	addToCounter := func(x int) StateIO[TestState, int] {
		return func(s TestState) IO[Pair[TestState, int]] {
			return func() Pair[TestState, int] {
				newState := TestState{Counter: s.Counter + x, Message: s.Message}
				return pair.MakePair(newState, s.Counter+x)
			}
		}
	}

	computation := F.Pipe1(
		Of[TestState](5),
		Chain(addToCounter),
	)

	result := computation(initial)()

	assert.Equal(t, 5, pair.Tail(result))
	assert.Equal(t, 5, pair.Head(result).Counter)
}

func TestMonadAp(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	fab := Of[TestState](N.Mul(3))
	fa := Of[TestState](7)

	result := MonadAp(fab, fa)(initial)()

	assert.Equal(t, 21, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestAp(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	computation := F.Pipe1(
		Of[TestState](N.Mul(4)),
		Ap[int](Of[TestState](10)),
	)

	result := computation(initial)()

	assert.Equal(t, 40, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestFromIO(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	ioComputation := func() int { return 42 }
	computation := FromIO[TestState](ioComputation)

	result := computation(initial)()

	assert.Equal(t, 42, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestFromIOK(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	ioFunc := func(x int) IO[string] {
		return func() string {
			return fmt.Sprintf("value: %d", x*2)
		}
	}

	kleisli := FromIOK[TestState](ioFunc)
	computation := kleisli(21)

	result := computation(initial)()

	assert.Equal(t, "value: 42", pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestChainedOperations(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	incrementCounter := func(x int) StateIO[TestState, int] {
		return func(s TestState) IO[Pair[TestState, int]] {
			return func() Pair[TestState, int] {
				newState := TestState{Counter: s.Counter + x, Message: s.Message}
				return pair.MakePair(newState, newState.Counter)
			}
		}
	}

	setMessage := func(count int) StateIO[TestState, string] {
		return func(s TestState) IO[Pair[TestState, string]] {
			return func() Pair[TestState, string] {
				msg := fmt.Sprintf("Count is %d", count)
				newState := TestState{Counter: s.Counter, Message: msg}
				return pair.MakePair(newState, msg)
			}
		}
	}

	computation := F.Pipe2(
		Of[TestState](5),
		Chain(incrementCounter),
		Chain(setMessage),
	)

	result := computation(initial)()

	assert.Equal(t, "Count is 5", pair.Tail(result))
	assert.Equal(t, 5, pair.Head(result).Counter)
	assert.Equal(t, "Count is 5", pair.Head(result).Message)
}

func TestStatefulComputation(t *testing.T) {
	initial := TestState{Counter: 10, Message: "start"}

	// A computation that reads and modifies state
	computation := func(s TestState) IO[Pair[TestState, int]] {
		return func() Pair[TestState, int] {
			newState := TestState{
				Counter: s.Counter * 2,
				Message: fmt.Sprintf("%s -> doubled", s.Message),
			}
			return pair.MakePair(newState, newState.Counter)
		}
	}

	result := computation(initial)()

	assert.Equal(t, 20, pair.Tail(result))
	assert.Equal(t, 20, pair.Head(result).Counter)
	assert.Equal(t, "start -> doubled", pair.Head(result).Message)
}

func TestMapPreservesState(t *testing.T) {
	initial := TestState{Counter: 42, Message: "important"}

	computation := F.Pipe2(
		Of[TestState](10),
		Map[TestState](N.Mul(2)),
		Map[TestState](N.Add(5)),
	)

	result := computation(initial)()

	// Value should be transformed: 10 * 2 + 5 = 25
	assert.Equal(t, 25, pair.Tail(result))
	// State should be unchanged
	assert.Equal(t, initial, pair.Head(result))
}

func TestChainModifiesState(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	addOne := func(x int) StateIO[TestState, int] {
		return func(s TestState) IO[Pair[TestState, int]] {
			return func() Pair[TestState, int] {
				newState := TestState{Counter: s.Counter + 1, Message: s.Message}
				return pair.MakePair(newState, x+1)
			}
		}
	}

	computation := F.Pipe2(
		Of[TestState](0),
		Chain(addOne),
		Chain(addOne),
	)

	result := computation(initial)()

	assert.Equal(t, 2, pair.Tail(result))
	assert.Equal(t, 2, pair.Head(result).Counter)
}

func TestApplicativeComposition(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	add := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}

	computation := F.Pipe1(
		Of[TestState](add(10)),
		Ap[int](Of[TestState](32)),
	)

	result := computation(initial)()

	assert.Equal(t, 42, pair.Tail(result))
}
