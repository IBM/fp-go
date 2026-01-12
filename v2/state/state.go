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
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

// Get returns a State computation that retrieves the current state and returns it as the value.
// The state is unchanged by this operation.
//
// Example:
//
//	type Counter struct { count int }
//	getState := Get[Counter]()
//	result := getState(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: Counter{count: 5}}
//
//go:inline
func Get[S any]() State[S, S] {
	return pair.Of[S]
}

// Gets applies a function to the current state and returns the result as the value.
// The state itself remains unchanged. This is useful for extracting or computing
// values from the state without modifying it.
//
// Example:
//
//	type Counter struct { count int }
//	getDouble := Gets(func(c Counter) int { return c.count * 2 })
//	result := getDouble(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 10}
//
//go:line
func Gets[FCT ~func(S) A, A, S any](f FCT) State[S, A] {
	return func(s S) Pair[S, A] {
		return pair.MakePair(s, f(s))
	}
}

// Put returns a State computation that replaces the current state with a new state.
// The returned value is Void, indicating this operation is performed for its side effect.
//
// Example:
//
//	type Counter struct { count int }
//	setState := Put[Counter]()
//	result := setState(Counter{count: 10})
//	// result = Pair{head: Counter{count: 10}, tail: Void}
//
//go:inline
func Put[S any]() State[S, Void] {
	return Of[S](function.VOID)
}

// Modify applies a transformation function to the current state, producing a new state.
// The returned value is Void, indicating this operation is performed for its side effect.
//
// Example:
//
//	type Counter struct { count int }
//	increment := Modify(func(c Counter) Counter { return Counter{count: c.count + 1} })
//	result := increment(Counter{count: 5})
//	// result = Pair{head: Counter{count: 6}, tail: Void}
func Modify[FCT ~func(S) S, S any](f FCT) State[S, Void] {
	return function.Flow2(
		f,
		Put[S](),
	)
}

// Of creates a State computation that returns the given value without modifying the state.
// This is the Pointed interface implementation for State, lifting a pure value into
// the State context.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Of[Counter](42)
//	result := computation(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
//
//go:inline
func Of[S, A any](a A) State[S, A] {
	return pair.FromTail[S](a)
}

// MonadMap transforms the value produced by a State computation using the given function,
// while preserving the state. This is the Functor interface implementation for State.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Of[Counter](10)
//	doubled := MonadMap(computation, func(x int) int { return x * 2 })
//	result := doubled(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 20}
//
//go:inline
func MonadMap[S any, FCT ~func(A) B, A, B any](fa State[S, A], f FCT) State[S, B] {
	return reader.MonadMap(fa, pair.Map[S](f))
}

// Map returns a function that transforms the value of a State computation.
// This is the curried version of MonadMap, useful for composition and pipelines.
//
// Example:
//
//	type Counter struct { count int }
//	double := Map[Counter](func(x int) int { return x * 2 })
//	computation := function.Pipe1(Of[Counter](10), double)
//	result := computation(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 20}
//
//go:inline
func Map[S any, FCT ~func(A) B, A, B any](f FCT) Operator[S, A, B] {
	return reader.Map[S](pair.Map[S](f))
}

// MonadChain sequences two State computations, where the second computation depends
// on the value produced by the first. The state is threaded through both computations.
// This is the Monad interface implementation for State.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Of[Counter](5)
//	chained := MonadChain(computation, func(x int) State[Counter, int] {
//	    return func(s Counter) Pair[Counter, int] {
//	        newState := Counter{count: s.count + x}
//	        return pair.MakePair(newState, x * 2)
//	    }
//	})
//	result := chained(Counter{count: 10})
//	// result = Pair{head: Counter{count: 15}, tail: 10}
func MonadChain[S any, FCT ~func(A) State[S, B], A, B any](fa State[S, A], f FCT) State[S, B] {
	return func(s S) Pair[S, B] {
		a := fa(s)
		return f(pair.Tail(a))(pair.Head(a))
	}
}

// Chain returns a function that sequences State computations.
// This is the curried version of MonadChain, useful for composition and pipelines.
//
// Example:
//
//	type Counter struct { count int }
//	addToCounter := func(x int) State[Counter, int] {
//	    return func(s Counter) Pair[Counter, int] {
//	        newState := Counter{count: s.count + x}
//	        return pair.MakePair(newState, newState.count)
//	    }
//	}
//	computation := function.Pipe1(Of[Counter](5), Chain(addToCounter))
//	result := computation(Counter{count: 10})
//	// result = Pair{head: Counter{count: 15}, tail: 15}
//
//go:inline
func Chain[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, B] {
	return function.Bind2nd(MonadChain[S, FCT, A, B], f)
}

// MonadAp applies a State computation containing a function to a State computation
// containing a value. Both computations are executed sequentially, threading the state
// through both. This is the Applicative interface implementation for State.
//
// Example:
//
//	type Counter struct { count int }
//	fab := Of[Counter](func(x int) int { return x * 2 })
//	fa := Of[Counter](21)
//	result := MonadAp(fab, fa)(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
func MonadAp[B, S, A any](fab State[S, func(A) B], fa State[S, A]) State[S, B] {
	return func(s S) Pair[S, B] {
		f := fab(s)
		a := fa(pair.Head(f))

		return pair.MakePair(pair.Head(a), pair.Tail(f)(pair.Tail(a)))
	}
}

// Ap returns a function that applies a State computation containing a function
// to a State computation containing a value. This is the curried version of MonadAp.
//
// Example:
//
//	type Counter struct { count int }
//	computation := function.Pipe1(
//	    Of[Counter](func(x int) int { return x * 2 }),
//	    Ap[int](Of[Counter](21)),
//	)
//	result := computation(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
//
//go:inline
func Ap[B, S, A any](ga State[S, A]) Operator[S, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, S, A], ga)
}

// MonadChainFirst sequences two State computations but returns the value from the first
// computation while still threading the state through both. This is useful when you want
// to perform a stateful side effect but keep the original value.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Of[Counter](42)
//	increment := func(x int) State[Counter, Void] {
//	    return Modify(func(c Counter) Counter { return Counter{count: c.count + 1} })
//	}
//	result := MonadChainFirst(computation, increment)(Counter{count: 5})
//	// result = Pair{head: Counter{count: 6}, tail: 42}
func MonadChainFirst[S any, FCT ~func(A) State[S, B], A, B any](ma State[S, A], f FCT) State[S, A] {
	return chain.MonadChainFirst(
		MonadChain[S, func(A) State[S, A], A, A],
		MonadMap[S, func(B) A],
		ma,
		f,
	)
}

// ChainFirst returns a function that sequences State computations but keeps the first value.
// This is the curried version of MonadChainFirst, useful for composition and pipelines.
//
// Example:
//
//	type Counter struct { count int }
//	increment := func(x int) State[Counter, Void] {
//	    return Modify(func(c Counter) Counter { return Counter{count: c.count + 1} })
//	}
//	computation := function.Pipe1(Of[Counter](42), ChainFirst(increment))
//	result := computation(Counter{count: 5})
//	// result = Pair{head: Counter{count: 6}, tail: 42}
func ChainFirst[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, A] {
	return chain.ChainFirst(
		Chain[S, func(A) State[S, A], A, A],
		Map[S, func(B) A],
		f,
	)
}

// Flatten removes one level of nesting from a State computation that produces another
// State computation. This is equivalent to MonadChain with the identity function.
//
// Example:
//
//	type Counter struct { count int }
//	nested := Of[Counter](Of[Counter](42))
//	flattened := Flatten(nested)
//	result := flattened(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
//
//go:inline
func Flatten[S, A any](mma State[S, State[S, A]]) State[S, A] {
	return MonadChain(mma, function.Identity[State[S, A]])
}

// Execute runs a State computation with the given initial state and returns only
// the final state, discarding the computed value. This is useful when you only
// care about the state transformations.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Modify(func(c Counter) Counter { return Counter{count: c.count + 1} })
//	finalState := Execute[Void, Counter](Counter{count: 5})(computation)
//	// finalState = Counter{count: 6}
func Execute[A, S any](s S) func(State[S, A]) S {
	return func(fa State[S, A]) S {
		return pair.Head(fa(s))
	}
}

// Evaluate runs a State computation with the given initial state and returns only
// the computed value, discarding the final state. This is useful when you only
// care about the result of the computation.
//
// Example:
//
//	type Counter struct { count int }
//	computation := Of[Counter](42)
//	value := Evaluate[int, Counter](Counter{count: 5})(computation)
//	// value = 42
func Evaluate[A, S any](s S) func(State[S, A]) A {
	return func(fa State[S, A]) A {
		return pair.Tail(fa(s))
	}
}

// MonadFlap applies a fixed value to a State computation containing a function.
// This is the reverse of MonadAp, where the value is known but the function is
// in the State context.
//
// Example:
//
//	type Counter struct { count int }
//	fab := Of[Counter](func(x int) int { return x * 2 })
//	result := MonadFlap(fab, 21)(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
func MonadFlap[FAB ~func(A) B, S, A, B any](fab State[S, FAB], a A) State[S, B] {
	return functor.MonadFlap(
		MonadMap[S, func(FAB) B],
		fab,
		a)
}

// Flap returns a function that applies a fixed value to a State computation containing
// a function. This is the curried version of MonadFlap.
//
// Example:
//
//	type Counter struct { count int }
//	applyTwentyOne := Flap[Counter, int, int](21)
//	computation := function.Pipe1(
//	    Of[Counter](func(x int) int { return x * 2 }),
//	    applyTwentyOne,
//	)
//	result := computation(Counter{count: 5})
//	// result = Pair{head: Counter{count: 5}, tail: 42}
func Flap[S, A, B any](a A) Operator[S, func(A) B, B] {
	return functor.Flap(
		Map[S, func(func(A) B) B],
		a)
}
