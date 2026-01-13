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

package state

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

// TestGet verifies that Get returns the current state as both state and value
func TestGet(t *testing.T) {
	initial := TestState{Counter: 42, Message: "test"}
	computation := Get[TestState]()

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, initial, pair.Tail(result), "value should equal state")
}

// TestGets verifies that Gets applies a function to state and returns the result
func TestGets(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	// Extract and double the counter
	computation := Gets(func(s TestState) int {
		return s.Counter * 2
	})

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 10, pair.Tail(result), "value should be doubled counter")
}

// TestPut verifies that Put replaces the state
func TestPut(t *testing.T) {
	newState := TestState{Counter: 10, Message: "new"}

	computation := Put[TestState]()

	result := computation(newState)

	assert.Equal(t, newState, pair.Head(result), "state should be replaced")
	assert.Equal(t, F.VOID, pair.Tail(result), "value should be Void")
}

// TestModify verifies that Modify transforms the state
func TestModify(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	increment := Modify(func(s TestState) TestState {
		return TestState{Counter: s.Counter + 1, Message: s.Message}
	})

	result := increment(initial)

	assert.Equal(t, 6, pair.Head(result).Counter, "counter should be incremented")
	assert.Equal(t, "test", pair.Head(result).Message, "message should be unchanged")
	assert.Equal(t, F.VOID, pair.Tail(result), "value should be Void")
}

// TestOf verifies that Of creates a computation with a value and unchanged state
func TestOf(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}
	computation := Of[TestState](42)

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 42, pair.Tail(result), "value should be 42")
}

// TestMonadMap verifies that MonadMap transforms the value
func TestMonadMap(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}
	computation := Of[TestState](10)

	mapped := MonadMap(computation, N.Mul(2))
	result := mapped(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 20, pair.Tail(result), "value should be doubled")
}

// TestMap verifies the curried version of MonadMap
func TestMap(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	double := Map[TestState](N.Mul(2))
	computation := F.Pipe1(
		Of[TestState](21),
		double,
	)

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 42, pair.Tail(result), "value should be doubled")
}

// TestMonadChain verifies that MonadChain sequences computations
func TestMonadChain(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	computation := Of[TestState](5)
	chained := MonadChain(computation, func(x int) State[TestState, string] {
		return func(s TestState) Pair[TestState, string] {
			newState := TestState{Counter: s.Counter + x, Message: fmt.Sprintf("value: %d", x)}
			return pair.MakePair(newState, fmt.Sprintf("result: %d", x*2))
		}
	})

	result := chained(initial)

	assert.Equal(t, "result: 10", pair.Tail(result), "value should be transformed")
	assert.Equal(t, 5, pair.Head(result).Counter, "counter should be updated")
	assert.Equal(t, "value: 5", pair.Head(result).Message, "message should be set")
}

// TestChain verifies the curried version of MonadChain
func TestChain(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	addToCounter := func(x int) State[TestState, int] {
		return func(s TestState) Pair[TestState, int] {
			newState := TestState{Counter: s.Counter + x, Message: s.Message}
			return pair.MakePair(newState, s.Counter+x)
		}
	}

	computation := F.Pipe1(
		Of[TestState](5),
		Chain(addToCounter),
	)

	result := computation(initial)

	assert.Equal(t, 5, pair.Tail(result), "value should be 5")
	assert.Equal(t, 5, pair.Head(result).Counter, "counter should be 5")
}

// TestMonadAp verifies applicative application
func TestMonadAp(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	fab := Of[TestState](N.Mul(3))
	fa := Of[TestState](7)

	result := MonadAp(fab, fa)(initial)

	assert.Equal(t, 21, pair.Tail(result), "value should be 21")
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestAp verifies the curried version of MonadAp
func TestAp(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	computation := F.Pipe1(
		Of[TestState](N.Mul(4)),
		Ap[int](Of[TestState](10)),
	)

	result := computation(initial)

	assert.Equal(t, 40, pair.Tail(result), "value should be 40")
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestMonadChainFirst verifies that ChainFirst keeps the first value
func TestMonadChainFirst(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	computation := Of[TestState](42)
	increment := func(x int) State[TestState, Void] {
		return Modify(func(s TestState) TestState {
			return TestState{Counter: s.Counter + 1, Message: s.Message}
		})
	}

	result := MonadChainFirst(computation, increment)(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be preserved")
	assert.Equal(t, 1, pair.Head(result).Counter, "counter should be incremented")
}

// TestChainFirst verifies the curried version of MonadChainFirst
func TestChainFirst(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	increment := func(x int) State[TestState, Void] {
		return Modify(func(s TestState) TestState {
			return TestState{Counter: s.Counter + 1, Message: s.Message}
		})
	}

	computation := F.Pipe1(Of[TestState](42), ChainFirst(increment))
	result := computation(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be preserved")
	assert.Equal(t, 6, pair.Head(result).Counter, "counter should be incremented")
}

// TestFlatten verifies that Flatten removes one level of nesting
func TestFlatten(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	nested := Of[TestState](Of[TestState](42))
	flattened := Flatten(nested)

	result := flattened(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be 42")
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestExecute verifies that Execute returns only the final state
func TestExecute(t *testing.T) {
	initial := TestState{Counter: 5, Message: "old"}

	computation := Modify(func(s TestState) TestState {
		return TestState{Counter: s.Counter + 1, Message: "new"}
	})

	finalState := Execute[Void](initial)(computation)

	assert.Equal(t, 6, finalState.Counter, "counter should be incremented")
	assert.Equal(t, "new", finalState.Message, "message should be updated")
}

// TestEvaluate verifies that Evaluate returns only the value
func TestEvaluate(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	computation := Of[TestState](42)

	value := Evaluate[int](initial)(computation)

	assert.Equal(t, 42, value, "value should be 42")
}

// TestMonadFlap verifies that MonadFlap applies a value to a function in State
func TestMonadFlap(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	fab := Of[TestState](func(x int) int { return x * 2 })
	result := MonadFlap(fab, 21)(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be 42")
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestFlap verifies the curried version of MonadFlap
func TestFlap(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	applyTwentyOne := Flap[TestState, int, int](21)
	computation := F.Pipe1(
		Of[TestState](func(x int) int { return x * 2 }),
		applyTwentyOne,
	)

	result := computation(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be 42")
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestChainedOperations verifies complex chained operations
func TestChainedOperations(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	incrementCounter := func(x int) State[TestState, int] {
		return func(s TestState) Pair[TestState, int] {
			newState := TestState{Counter: s.Counter + x, Message: s.Message}
			return pair.MakePair(newState, newState.Counter)
		}
	}

	setMessage := func(count int) State[TestState, string] {
		return func(s TestState) Pair[TestState, string] {
			msg := fmt.Sprintf("Count is %d", count)
			newState := TestState{Counter: s.Counter, Message: msg}
			return pair.MakePair(newState, msg)
		}
	}

	computation := F.Pipe2(
		Of[TestState](5),
		Chain(incrementCounter),
		Chain(setMessage),
	)

	result := computation(initial)

	assert.Equal(t, "Count is 5", pair.Tail(result), "value should be message")
	assert.Equal(t, 5, pair.Head(result).Counter, "counter should be 5")
	assert.Equal(t, "Count is 5", pair.Head(result).Message, "message should be set")
}

// TestMapPreservesState verifies that Map operations don't modify state
func TestMapPreservesState(t *testing.T) {
	initial := TestState{Counter: 42, Message: "important"}

	computation := F.Pipe2(
		Of[TestState](10),
		Map[TestState](N.Mul(2)),
		Map[TestState](N.Add(5)),
	)

	result := computation(initial)

	// Value should be transformed: 10 * 2 + 5 = 25
	assert.Equal(t, 25, pair.Tail(result), "value should be 25")
	// State should be unchanged
	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
}

// TestChainModifiesState verifies that Chain operations can modify state
func TestChainModifiesState(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	addOne := func(x int) State[TestState, int] {
		return func(s TestState) Pair[TestState, int] {
			newState := TestState{Counter: s.Counter + 1, Message: s.Message}
			return pair.MakePair(newState, x+1)
		}
	}

	computation := F.Pipe2(
		Of[TestState](0),
		Chain(addOne),
		Chain(addOne),
	)

	result := computation(initial)

	assert.Equal(t, 2, pair.Tail(result), "value should be 2")
	assert.Equal(t, 2, pair.Head(result).Counter, "counter should be 2")
}

// TestApplicativeComposition verifies applicative composition
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

	result := computation(initial)

	assert.Equal(t, 42, pair.Tail(result), "value should be 42")
}

// TestStatefulComputation verifies a computation that reads and modifies state
func TestStatefulComputation(t *testing.T) {
	initial := TestState{Counter: 10, Message: "start"}

	// A computation that reads and modifies state
	computation := func(s TestState) Pair[TestState, int] {
		newState := TestState{
			Counter: s.Counter * 2,
			Message: fmt.Sprintf("%s -> doubled", s.Message),
		}
		return pair.MakePair(newState, newState.Counter)
	}

	result := computation(initial)

	assert.Equal(t, 20, pair.Tail(result), "value should be 20")
	assert.Equal(t, 20, pair.Head(result).Counter, "counter should be 20")
	assert.Equal(t, "start -> doubled", pair.Head(result).Message, "message should be updated")
}

// TestGetAndModify verifies combining Get and Modify
func TestGetAndModify(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	step1 := Chain(func(s TestState) State[TestState, Void] {
		return Modify(func(_ TestState) TestState {
			return TestState{Counter: s.Counter * 2, Message: s.Message + "!"}
		})
	})

	step2 := Chain(func(_ Void) State[TestState, TestState] {
		return Get[TestState]()
	})

	computation := step2(step1(Get[TestState]()))

	result := computation(initial)

	assert.Equal(t, 10, pair.Tail(result).Counter, "counter should be doubled")
	assert.Equal(t, "test!", pair.Tail(result).Message, "message should have exclamation")
}

// TestGetsWithComplexExtraction verifies Gets with complex state extraction
func TestGetsWithComplexExtraction(t *testing.T) {
	initial := TestState{Counter: 5, Message: "hello"}

	computation := Gets(func(s TestState) string {
		return fmt.Sprintf("%s: %d", s.Message, s.Counter)
	})

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, "hello: 5", pair.Tail(result), "value should be formatted string")
}

// TestMultipleModifications verifies multiple state modifications
func TestMultipleModifications(t *testing.T) {
	initial := TestState{Counter: 0, Message: ""}

	increment := Modify(func(s TestState) TestState {
		return TestState{Counter: s.Counter + 1, Message: s.Message}
	})

	setMessage := Modify(func(s TestState) TestState {
		return TestState{Counter: s.Counter, Message: "done"}
	})

	step1 := ChainFirst(func(_ Void) State[TestState, Void] { return increment })
	step2 := ChainFirst(func(_ Void) State[TestState, Void] { return setMessage })

	computation := step2(step1(increment))

	result := computation(initial)

	assert.Equal(t, 2, pair.Head(result).Counter, "counter should be 2")
	assert.Equal(t, "done", pair.Head(result).Message, "message should be 'done'")
}

// TestExecuteWithComplexState verifies Execute with complex state transformations
func TestExecuteWithComplexState(t *testing.T) {
	initial := TestState{Counter: 1, Message: "start"}

	step1 := Modify(func(s TestState) TestState {
		return TestState{Counter: s.Counter * 2, Message: s.Message}
	})

	step2 := ChainFirst(func(_ Void) State[TestState, Void] {
		return Modify(func(s TestState) TestState {
			return TestState{Counter: s.Counter + 10, Message: "end"}
		})
	})

	computation := step2(step1)

	finalState := Execute[Void](initial)(computation)

	assert.Equal(t, 12, finalState.Counter, "counter should be (1*2)+10 = 12")
	assert.Equal(t, "end", finalState.Message, "message should be 'end'")
}

// TestEvaluateWithChain verifies Evaluate with chained computations
func TestEvaluateWithChain(t *testing.T) {
	initial := TestState{Counter: 5, Message: "test"}

	computation := F.Pipe2(
		Of[TestState](10),
		Map[TestState](N.Mul(2)),
		Chain(func(x int) State[TestState, string] {
			return Of[TestState](fmt.Sprintf("result: %d", x))
		}),
	)

	value := Evaluate[string](initial)(computation)

	assert.Equal(t, "result: 20", value, "value should be 'result: 20'")
}
