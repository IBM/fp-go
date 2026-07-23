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

package identity_test

import (
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	N "github.com/IBM/fp-go/v2/number"
	T "github.com/IBM/fp-go/v2/tuple"
)

// ExampleOf demonstrates that Of is the identity function — it returns its argument unchanged.
func ExampleOf() {
	value := I.Of(42)
	fmt.Println(value)
	// Output: 42
}

// ExampleMonadAp demonstrates uncurried function application in the Identity monad.
func ExampleMonadAp() {
	result := I.MonadAp(func(n int) int { return n * 2 }, 21)
	fmt.Println(result)
	// Output: 42
}

// ExampleAp demonstrates applying a wrapped function to a value using Pipe.
func ExampleAp() {
	double := func(n int) int { return n * 2 }
	result := F.Pipe1(double, I.Ap[int](21))
	fmt.Println(result)
	// Output: 42
}

// ExampleMonadMap demonstrates uncurried value transformation in the Identity monad.
func ExampleMonadMap() {
	result := I.MonadMap(21, func(n int) int { return n * 2 })
	fmt.Println(result)
	// Output: 42
}

// ExampleMap demonstrates transforming a value using Pipe.
func ExampleMap() {
	result := F.Pipe1(21, I.Map(func(n int) int { return n * 2 }))
	fmt.Println(result)
	// Output: 42
}

// ExampleMonadMapTo demonstrates replacing a value with a constant (uncurried).
func ExampleMonadMapTo() {
	result := I.MonadMapTo("ignored", 42)
	fmt.Println(result)
	// Output: 42
}

// ExampleMapTo demonstrates replacing any value with a constant using Pipe.
func ExampleMapTo() {
	result := F.Pipe1("ignored", I.MapTo[string](42))
	fmt.Println(result)
	// Output: 42
}

// ExampleMonadChain demonstrates uncurried monadic bind in the Identity monad.
func ExampleMonadChain() {
	result := I.MonadChain(21, func(n int) int { return n * 2 })
	fmt.Println(result)
	// Output: 42
}

// ExampleChain demonstrates sequential composition using Pipe.
func ExampleChain() {
	result := F.Pipe2(
		10,
		I.Chain(N.Mul(2)),
		I.Chain(N.Add(5)),
	)
	fmt.Println(result)
	// Output: 25
}

// ExampleMonadChainFirst demonstrates that ChainFirst executes a side-effecting computation
// but preserves the original value (uncurried).
func ExampleMonadChainFirst() {
	result := I.MonadChainFirst(42, func(n int) string {
		return strconv.Itoa(n) // side effect: produces a string, but is discarded
	})
	fmt.Println(result)
	// Output: 42
}

// ExampleChainFirst demonstrates executing a computation for its effect while keeping
// the original value.
func ExampleChainFirst() {
	result := F.Pipe1(
		42,
		I.ChainFirst(func(n int) string {
			return strconv.Itoa(n) // side effect discarded
		}),
	)
	fmt.Println(result)
	// Output: 42
}

// ExampleMonadFlap demonstrates uncurried flapped application.
func ExampleMonadFlap() {
	double := func(n int) int { return n * 2 }
	result := I.MonadFlap(double, 21)
	fmt.Println(result)
	// Output: 42
}

// ExampleFlap demonstrates applying a fixed value to a function using Pipe.
func ExampleFlap() {
	double := func(n int) int { return n * 2 }
	result := F.Pipe1(double, I.Flap[int](21))
	fmt.Println(result)
	// Output: 42
}

// ExampleExtract demonstrates that Extract is the identity function for the Comonad interface.
func ExampleExtract() {
	value := I.Extract(42)
	fmt.Println(value)
	// Output: 42
}

// ExampleExtend demonstrates the Comonad extend operation, which is just function application.
func ExampleExtend() {
	result := F.Pipe1(21, I.Extend(func(n int) int { return n * 2 }))
	fmt.Println(result)
	// Output: 42
}

// ExampleDo demonstrates initialising a do-notation context.
func ExampleDo() {
	type State struct{ X int }
	s := I.Do(State{})
	fmt.Println(s)
	// Output: {0}
}

// ExampleBind demonstrates sequential do-notation composition where each step
// can read the accumulated state.
func ExampleBind() {
	type State struct{ X, Y int }

	result := F.Pipe2(
		I.Do(State{}),
		I.Bind(
			func(x int) func(State) State {
				return func(s State) State { s.X = x; return s }
			},
			func(State) int { return 42 },
		),
		I.Bind(
			func(y int) func(State) State {
				return func(s State) State { s.Y = y; return s }
			},
			func(s State) int { return s.X * 2 },
		),
	)
	fmt.Println(result)
	// Output: {42 84}
}

// ExampleLet demonstrates computing a derived value from the accumulated state.
func ExampleLet() {
	type State struct{ X, Y, Sum int }

	result := F.Pipe1(
		I.Do(State{X: 10, Y: 20}),
		I.Let(
			func(sum int) func(State) State {
				return func(s State) State { s.Sum = sum; return s }
			},
			func(s State) int { return s.X + s.Y },
		),
	)
	fmt.Println(result)
	// Output: {10 20 30}
}

// ExampleLetTo demonstrates attaching a constant value to the context.
func ExampleLetTo() {
	type State struct {
		X        int
		Constant string
	}

	result := F.Pipe1(
		I.Do(State{X: 10}),
		I.LetTo(
			func(c string) func(State) State {
				return func(s State) State { s.Constant = c; return s }
			},
			"fixed",
		),
	)
	fmt.Println(result)
	// Output: {10 fixed}
}

// ExampleBindTo demonstrates lifting an initial value into a do-notation context.
func ExampleBindTo() {
	type State struct{ X, Y int }

	result := F.Pipe2(
		42,
		I.BindTo(func(x int) State { return State{X: x} }),
		I.Bind(
			func(y int) func(State) State {
				return func(s State) State { s.Y = y; return s }
			},
			func(s State) int { return s.X * 2 },
		),
	)
	fmt.Println(result)
	// Output: {42 84}
}

// ExampleApS demonstrates combining independent values into a context using
// the Applicative (rather than Monadic) interface.
func ExampleApS() {
	type State struct{ X, Y int }

	result := F.Pipe2(
		I.Do(State{}),
		I.ApS(
			func(x int) func(State) State {
				return func(s State) State { s.X = x; return s }
			},
			42,
		),
		I.ApS(
			func(y int) func(State) State {
				return func(s State) State { s.Y = y; return s }
			},
			100,
		),
	)
	fmt.Println(result)
	// Output: {42 100}
}

// ExampleMakeTraversable demonstrates that traversing an Identity value is
// equivalent to mapping: the transformation function is applied directly.
func ExampleMakeTraversable() {
	traverseWithItoa := I.MakeTraversable[int, string, string]()
	result := traverseWithItoa(strconv.Itoa)(42)
	fmt.Println(result)
	// Output: 42
}

// ExampleSequenceTuple2 demonstrates sequencing a 2-tuple of Identity values.
func ExampleSequenceTuple2() {
	tuple := T.MakeTuple2(1, 2)
	result := I.SequenceTuple2(tuple)
	fmt.Println(result)
	// Output: Tuple2[int, int](1, 2)
}

// ExampleTraverseTuple2 demonstrates traversing a 2-tuple with per-element transformations.
func ExampleTraverseTuple2() {
	result := I.TraverseTuple2(N.Mul(2), N.Mul(3))(T.MakeTuple2(1, 2))
	fmt.Println(result)
	// Output: Tuple2[int, int](2, 6)
}
