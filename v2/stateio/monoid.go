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
	M "github.com/IBM/fp-go/v2/monoid"
)

// ApplicativeMonoid lifts a monoid into the StateIO applicative functor context.
//
// This function creates a monoid for StateIO[S, A] values given a monoid for the base type A.
// It uses the StateIO monad's applicative operations (Of, MonadMap, MonadAp) to lift the
// monoid operations into the StateIO context, allowing you to combine stateful computations
// with side effects that produce monoidal values.
//
// The resulting monoid combines StateIO computations by:
//  1. Threading the state through both computations sequentially
//  2. Executing the IO effects of both computations in sequence
//  3. Combining the produced values using the underlying monoid's Concat operation
//  4. Returning a new StateIO computation with the combined value and final state
//
// The empty element is a StateIO computation that returns the underlying monoid's empty value
// without modifying the state or performing any side effects.
//
// This is particularly useful for:
//   - Accumulating results across multiple stateful computations with side effects
//   - Building complex state transformations that aggregate values while performing IO
//   - Combining independent stateful operations that produce monoidal results
//
// Type Parameters:
//   - S: The state type that is threaded through the computations
//   - A: The value type that has a monoid structure
//
// Parameters:
//   - m: A monoid for the base type A that defines how to combine values
//
// Returns:
//   - A Monoid[StateIO[S, A]] that combines stateful IO computations using the base monoid
//
// The resulting monoid satisfies the standard monoid laws:
//   - Associativity: Concat(Concat(s1, s2), s3) = Concat(s1, Concat(s2, s3))
//   - Left identity: Concat(Empty(), s) = s
//   - Right identity: Concat(s, Empty()) = s
//
// Example with integer addition:
//
//	import (
//	    N "github.com/IBM/fp-go/v2/number"
//	    "github.com/IBM/fp-go/v2/io"
//	    "github.com/IBM/fp-go/v2/pair"
//	)
//
//	type Counter struct {
//	    count int
//	}
//
//	// Create a monoid for StateIO[Counter, int] using integer addition
//	intAddMonoid := N.MonoidSum[int]()
//	stateIOMonoid := stateio.ApplicativeMonoid[Counter](intAddMonoid)
//
//	// Create two stateful IO computations
//	s1 := stateio.Of[Counter](5)  // Returns 5, state unchanged
//	s2 := stateio.Of[Counter](3)  // Returns 3, state unchanged
//
//	// Combine them using the monoid
//	combined := stateIOMonoid.Concat(s1, s2)
//	result := combined(Counter{count: 10})()
//	// result = Pair{head: Counter{count: 10}, tail: 8}  // 5 + 3
//
//	// Empty element
//	empty := stateIOMonoid.Empty()
//	emptyResult := empty(Counter{count: 10})()
//	// emptyResult = Pair{head: Counter{count: 10}, tail: 0}
//
// Example with string concatenation and state modification:
//
//	import (
//	    S "github.com/IBM/fp-go/v2/string"
//	    "github.com/IBM/fp-go/v2/io"
//	)
//
//	type Logger struct {
//	    logs []string
//	}
//
//	strMonoid := S.Monoid
//	stateIOMonoid := stateio.ApplicativeMonoid[Logger](strMonoid)
//
//	// Stateful IO computation that logs and returns a message
//	logMessage := func(msg string) stateio.StateIO[Logger, string] {
//	    return func(s Logger) io.IO[pair.Pair[Logger, string]] {
//	        return func() pair.Pair[Logger, string] {
//	            newState := Logger{logs: append(s.logs, msg)}
//	            return pair.MakePair(newState, msg)
//	        }
//	    }
//	}
//
//	s1 := logMessage("Hello")
//	s2 := logMessage(" World")
//
//	// Combine the computations - both log entries are added, messages concatenated
//	combined := stateIOMonoid.Concat(s1, s2)
//	result := combined(Logger{logs: []string{}})()
//	// result.head.logs = ["Hello", " World"]
//	// result.tail = "Hello World"
//
// Example demonstrating monoid laws:
//
//	intAddMonoid := N.MonoidSum[int]()
//	m := stateio.ApplicativeMonoid[Counter](intAddMonoid)
//
//	s1 := stateio.Of[Counter](1)
//	s2 := stateio.Of[Counter](2)
//	s3 := stateio.Of[Counter](3)
//
//	initialState := Counter{count: 0}
//
//	// Associativity
//	left := m.Concat(m.Concat(s1, s2), s3)
//	right := m.Concat(s1, m.Concat(s2, s3))
//	// Both produce: Pair{head: Counter{count: 0}, tail: 6}
//
//	// Left identity
//	leftId := m.Concat(m.Empty(), s1)
//	// Produces: Pair{head: Counter{count: 0}, tail: 1}
//
//	// Right identity
//	rightId := m.Concat(s1, m.Empty())
//	// Produces: Pair{head: Counter{count: 0}, tail: 1}
//
//go:inline
func ApplicativeMonoid[S, A any](m M.Monoid[A]) M.Monoid[StateIO[S, A]] {
	return M.ApplicativeMonoid(
		Of[S, A],
		MonadMap[S, A, func(A) A],
		MonadAp[A, S, A],
		m)
}
