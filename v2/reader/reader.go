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

package reader

import (
	"github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/functor"
	T "github.com/IBM/fp-go/v2/tuple"
)

// Ask reads the current context and returns it as the result.
// This is the fundamental operation for accessing the environment.
//
// Example:
//
//	type Config struct { Host string }
//	r := reader.Ask[Config]()
//	config := r(Config{Host: "localhost"}) // Returns the config itself
//
//go:inline
func Ask[R any]() Reader[R, R] {
	return function.Identity[R]
}

// Asks projects a value from the global context in a Reader.
// It's essentially an identity function that makes the intent clearer.
//
// Example:
//
//	type Config struct { Port int }
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	port := getPort(Config{Port: 8080}) // Returns 8080
//
//go:inline
func Asks[R, A any](f Reader[R, A]) Reader[R, A] {
	return f
}

// AsksReader creates a Reader that depends on the environment to produce another Reader,
// then immediately executes that Reader with the same environment.
//
// This is useful when you need to dynamically choose a Reader based on the environment.
//
// Example:
//
//	type Config struct { UseCache bool }
//	r := reader.AsksReader(func(c Config) reader.Reader[Config, string] {
//	    if c.UseCache {
//	        return reader.Of[Config]("cached")
//	    }
//	    return reader.Of[Config]("fresh")
//	})
//
//go:inline
func AsksReader[R, A any](f Kleisli[R, R, A]) Reader[R, A] {
	//go:inline
	return func(r R) A {
		return f(r)(r)
	}
}

// MonadMap transforms the result value of a Reader using the provided function.
// This is the monadic version that takes the Reader as the first parameter.
//
// Example:
//
//	type Config struct { Port int }
//	getPort := func(c Config) int { return c.Port }
//	getPortStr := reader.MonadMap(getPort, strconv.Itoa)
//	result := getPortStr(Config{Port: 8080}) // "8080"
//
//go:inline
func MonadMap[E, A, B any](fa Reader[E, A], f func(A) B) Reader[E, B] {
	return function.Flow2(fa, f)
}

// MonadMapTo creates a new Reader that completely ignores the first Reader and returns a constant value.
// This is the monadic version that takes both the Reader and the constant value as parameters.
//
// IMPORTANT: Readers are pure functions with no side effects. This function does NOT compose or evaluate
// the first Reader - it completely ignores it and returns a new Reader that always returns the constant value.
// The first Reader is neither executed during composition nor when the resulting Reader runs.
//
// Type Parameters:
//   - E: The environment type
//   - A: The result type of the first Reader (completely ignored)
//   - B: The type of the constant value to return
//
// Parameters:
//   - _: The first Reader (completely ignored, never evaluated)
//   - b: The constant value to return
//
// Returns:
//   - A new Reader that ignores the environment and always returns b
//
// Example:
//
//	type Config struct { Counter int }
//	increment := func(c Config) int { return c.Counter + 1 }
//	// Create a Reader that ignores increment and returns "done"
//	r := reader.MonadMapTo(increment, "done")
//	result := r(Config{Counter: 5}) // "done" (increment was never evaluated)
//
//go:inline
func MonadMapTo[E, A, B any](_ Reader[E, A], b B) Reader[E, B] {
	return Of[E](b)
}

// Map transforms the result value of a Reader using the provided function.
// This is the Functor operation that allows you to transform values inside the Reader context.
//
// Map can be used to turn functions `func(A)B` into functions `(fa F[A])F[B]` whose argument and return types
// use the type constructor `F` to represent some computational context.
//
// Example:
//
//	type Config struct { Port int }
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	getPortStr := reader.Map(strconv.Itoa)(getPort)
//	result := getPortStr(Config{Port: 8080}) // "8080"
//
//go:inline
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return function.Bind2nd(MonadMap[E, A, B], f)
}

// MapTo creates an operator that completely ignores any Reader and returns a constant value.
// This is the curried version where the constant value is provided first,
// returning a function that can be applied to any Reader.
//
// IMPORTANT: Readers are pure functions with no side effects. This operator does NOT compose or evaluate
// the input Reader - it completely ignores it and returns a new Reader that always returns the constant value.
// The input Reader is neither executed during composition nor when the resulting Reader runs.
//
// Type Parameters:
//   - E: The environment type
//   - A: The result type of the input Reader (completely ignored)
//   - B: The type of the constant value to return
//
// Parameters:
//   - b: The constant value to return
//
// Returns:
//   - An Operator that takes a Reader[E, A] and returns Reader[E, B]
//
// Example:
//
//	type Config struct { Counter int }
//	increment := reader.Asks(func(c Config) int { return c.Counter + 1 })
//	// Create an operator that ignores any Reader and returns "done"
//	toDone := reader.MapTo[Config, int, string]("done")
//	pipeline := toDone(increment)
//	result := pipeline(Config{Counter: 5}) // "done" (increment was never evaluated)
//
// Example - In a functional pipeline:
//
//	type Env struct { Step int }
//	step1 := reader.Asks(func(e Env) int { return e.Step })
//	pipeline := F.Pipe1(
//	    step1,
//	    reader.MapTo[Env, int, string]("complete"),
//	)
//	output := pipeline(Env{Step: 1}) // "complete" (step1 was never evaluated)
//
//go:inline
func MapTo[E, A, B any](b B) Operator[E, A, B] {
	return Of[Reader[E, A]](Of[E](b))
}

// MonadAp applies a Reader containing a function to a Reader containing a value.
// Both Readers share the same environment and are evaluated with it.
// This is the monadic version that takes both parameters.
//
// Example:
//
//	type Config struct { X, Y int }
//	add := func(x int) func(int) int { return func(y int) int { return x + y } }
//	getX := func(c Config) func(int) int { return add(c.X) }
//	getY := func(c Config) int { return c.Y }
//	result := reader.MonadAp(getX, getY)
//	sum := result(Config{X: 3, Y: 4}) // 7
func MonadAp[B, R, A any](fab Reader[R, func(A) B], fa Reader[R, A]) Reader[R, B] {
	return func(r R) B {
		return fab(r)(fa(r))
	}
}

// Ap applies a Reader containing a function to a Reader containing a value.
// This is the Applicative operation for combining independent computations.
//
// Example:
//
//	type Config struct { X, Y int }
//	add := func(x int) func(int) int { return func(y int) int { return x + y } }
//	getX := reader.Map(add)(reader.Asks(func(c Config) int { return c.X }))
//	getY := reader.Asks(func(c Config) int { return c.Y })
//	getSum := reader.Ap(getY)(getX)
//	sum := getSum(Config{X: 3, Y: 4}) // 7
func Ap[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, R, A], fa)
}

// Of lifts a pure value into the Reader context.
// The resulting Reader ignores its environment and always returns the given value.
// This is the Pointed/Applicative pure operation.
//
// Example:
//
//	type Config struct { Host string }
//	r := reader.Of[Config]("constant value")
//	result := r(Config{Host: "any"}) // "constant value"
func Of[R, A any](a A) Reader[R, A] {
	return function.Constant1[R](a)
}

// MonadChain sequences two Reader computations where the second depends on the result of the first.
// Both computations share the same environment.
// This is the monadic bind operation (flatMap).
//
// Example:
//
//	type Config struct { UserId int }
//	getUser := func(c Config) int { return c.UserId }
//	getUserName := func(id int) reader.Reader[Config, string] {
//	    return func(c Config) string { return fmt.Sprintf("User%d", id) }
//	}
//	r := reader.MonadChain(getUser, getUserName)
//	name := r(Config{UserId: 42}) // "User42"
func MonadChain[R, A, B any](ma Reader[R, A], f Kleisli[R, A, B]) Reader[R, B] {
	return func(r R) B {
		return f(ma(r))(r)
	}
}

// Chain sequences two Reader computations where the second depends on the result of the first.
// This is the Monad operation that enables dependent computations.
//
// Relationship with Compose:
//
// Chain and Compose serve different purposes in Reader composition:
//
//   - Chain: Monadic composition - sequences Readers that share the SAME environment type.
//     The second Reader depends on the VALUE produced by the first Reader, but both
//     Readers receive the same environment R. This is the monadic bind (>>=) operation.
//     Signature: Chain[R, A, B](f: A -> Reader[R, B]) -> Reader[R, A] -> Reader[R, B]
//
//   - Compose: Function composition - chains Readers where the OUTPUT of the first
//     becomes the INPUT environment of the second. The environment types can differ.
//     This is standard function composition (.) for Readers as functions.
//     Signature: Compose[C, R, B](ab: Reader[R, B]) -> Reader[B, C] -> Reader[R, C]
//
// Key Differences:
//
//  1. Environment handling:
//     - Chain: Both Readers use the same environment R
//     - Compose: First Reader's output B becomes second Reader's input environment
//
//  2. Data flow:
//     - Chain: R -> A, then A -> Reader[R, B], both using same R
//     - Compose: R -> B, then B -> C (B is both output and environment)
//
//  3. Use cases:
//     - Chain: Dependent computations in the same context (e.g., fetch user, then fetch user's posts)
//     - Compose: Transforming nested environments (e.g., extract config from app state, then read from config)
//
// Example:
//
//	type Config struct { UserId int }
//	getUser := reader.Asks(func(c Config) int { return c.UserId })
//	getUserName := func(id int) reader.Reader[Config, string] {
//	    return reader.Of[Config](fmt.Sprintf("User%d", id))
//	}
//	r := reader.Chain(getUserName)(getUser)
//	name := r(Config{UserId: 42}) // "User42"
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return function.Bind2nd(MonadChain[R, A, B], f)
}

// MonadChainTo completely ignores the first Reader and returns the second Reader.
// This is the monadic version that takes both Readers as parameters.
//
// IMPORTANT: Readers are pure functions with no side effects. This function does NOT compose or evaluate
// the first Reader - it completely ignores it and returns the second Reader directly.
// The first Reader is neither executed during composition nor when the resulting Reader runs.
//
// Type Parameters:
//   - A: The result type of the first Reader (completely ignored)
//   - R: The environment type
//   - B: The result type of the second Reader
//
// Parameters:
//   - _: The first Reader (completely ignored, never evaluated)
//   - b: The second Reader to return
//
// Returns:
//   - The second Reader unchanged
//
// Example:
//
//	type Config struct { Counter int; Message string }
//	increment := func(c Config) int { return c.Counter + 1 }
//	getMessage := func(c Config) string { return c.Message }
//	// Ignore increment and return getMessage
//	r := reader.MonadChainTo(increment, getMessage)
//	result := r(Config{Counter: 5, Message: "done"}) // "done" (increment was never evaluated)
//
//go:inline
func MonadChainTo[A, R, B any](_ Reader[R, A], b Reader[R, B]) Reader[R, B] {
	return b
}

// ChainTo creates an operator that completely ignores any Reader and returns a specific Reader.
// This is the curried version where the second Reader is provided first,
// returning a function that can be applied to any first Reader (which will be ignored).
//
// IMPORTANT: Readers are pure functions with no side effects. This operator does NOT compose or evaluate
// the input Reader - it completely ignores it and returns the specified Reader directly.
// The input Reader is neither executed during composition nor when the resulting Reader runs.
//
// Type Parameters:
//   - A: The result type of the first Reader (completely ignored)
//   - R: The environment type
//   - B: The result type of the second Reader
//
// Parameters:
//   - b: The Reader to return (ignoring any input Reader)
//
// Returns:
//   - An Operator that takes a Reader[R, A] and returns Reader[R, B]
//
// Example:
//
//	type Config struct { Counter int; Message string }
//	getMessage := func(c Config) string { return c.Message }
//	// Create an operator that ignores any Reader and returns getMessage
//	thenGetMessage := reader.ChainTo[int, Config, string](getMessage)
//
//	increment := func(c Config) int { return c.Counter + 1 }
//	pipeline := thenGetMessage(increment)
//	result := pipeline(Config{Counter: 5, Message: "done"}) // "done" (increment was never evaluated)
//
// Example - In a functional pipeline:
//
//	type Env struct { Step int; Result string }
//	step1 := reader.Asks(func(e Env) int { return e.Step })
//	getResult := reader.Asks(func(e Env) string { return e.Result })
//
//	pipeline := F.Pipe1(
//	    step1,
//	    reader.ChainTo[int, Env, string](getResult),
//	)
//	output := pipeline(Env{Step: 1, Result: "success"}) // "success" (step1 was never evaluated)
//
//go:inline
func ChainTo[A, R, B any](b Reader[R, B]) Operator[R, A, B] {
	return Of[Reader[R, A]](b)
}

// Flatten removes one level of Reader nesting.
// Converts Reader[R, Reader[R, A]] to Reader[R, A].
//
// Example:
//
//	type Config struct { Value int }
//	nested := func(c Config) reader.Reader[Config, int] {
//	    return func(c2 Config) int { return c.Value + c2.Value }
//	}
//	flat := reader.Flatten(nested)
//	result := flat(Config{Value: 5}) // 10 (5 + 5)
func Flatten[R, A any](mma Reader[R, Reader[R, A]]) Reader[R, A] {
	return MonadChain(mma, function.Identity[Reader[R, A]])
}

// Compose composes two Readers sequentially, where the output environment of the first
// becomes the input environment of the second.
//
// Relationship with Chain:
//
// Compose and Chain serve different purposes in Reader composition:
//
//   - Compose: Function composition - chains Readers where the OUTPUT of the first
//     becomes the INPUT environment of the second. The environment types can differ.
//     This is standard function composition (.) for Readers as functions.
//     Signature: Compose[C, R, B](ab: Reader[R, B]) -> Reader[B, C] -> Reader[R, C]
//
//   - Chain: Monadic composition - sequences Readers that share the SAME environment type.
//     The second Reader depends on the VALUE produced by the first Reader, but both
//     Readers receive the same environment R. This is the monadic bind (>>=) operation.
//     Signature: Chain[R, A, B](f: A -> Reader[R, B]) -> Reader[R, A] -> Reader[R, B]
//
// Key Differences:
//
//  1. Environment handling:
//     - Compose: First Reader's output B becomes second Reader's input environment
//     - Chain: Both Readers use the same environment R
//
//  2. Data flow:
//     - Compose: R -> B, then B -> C (B is both output and environment)
//     - Chain: R -> A, then A -> Reader[R, B], both using same R
//
//  3. Use cases:
//     - Compose: Transforming nested environments (e.g., extract config from app state, then read from config)
//     - Chain: Dependent computations in the same context (e.g., fetch user, then fetch user's posts)
//
// Visual Comparison:
//
//	// Compose: Environment transformation
//	type AppState struct { Config Config }
//	type Config struct { Port int }
//	getConfig := func(s AppState) Config { return s.Config }
//	getPort := func(c Config) int { return c.Port }
//	getPortFromState := reader.Compose(getConfig)(getPort)
//	// Flow: AppState -> Config -> int (Config is both output and next input)
//
//	// Chain: Same environment, dependent values
//	type Env struct { UserId int; Users map[int]string }
//	getUserId := func(e Env) int { return e.UserId }
//	getUser := func(id int) reader.Reader[Env, string] {
//	    return func(e Env) string { return e.Users[id] }
//	}
//	getUserName := reader.Chain(getUser)(getUserId)
//	// Flow: Env -> int, then int -> Reader[Env, string] (Env used twice)
//
// Example:
//
//	type Config struct { Port int }
//	type Env struct { Config Config }
//	getConfig := func(e Env) Config { return e.Config }
//	getPort := func(c Config) int { return c.Port }
//	getPortFromEnv := reader.Compose(getConfig)(getPort)
//
//go:inline
func Compose[C, R, B any](ab Reader[R, B]) Kleisli[R, Reader[B, C], C] {
	return function.Bind1st(function.Flow2[Reader[R, B], Reader[B, C]], ab)
}

// First applies a Reader to the first element of a tuple, leaving the second element unchanged.
// This is useful for working with paired data where only one element needs transformation.
//
// Example:
//
//	double := N.Mul(2)
//	r := reader.First[int, int, string](double)
//	result := r(tuple.MakeTuple2(5, "hello")) // (10, "hello")
func First[A, B, C any](pab Reader[A, B]) Reader[T.Tuple2[A, C], T.Tuple2[B, C]] {
	return func(tac T.Tuple2[A, C]) T.Tuple2[B, C] {
		return T.MakeTuple2(pab(tac.F1), tac.F2)
	}
}

// Second applies a Reader to the second element of a tuple, leaving the first element unchanged.
// This is useful for working with paired data where only one element needs transformation.
//
// Example:
//
//	double := N.Mul(2)
//	r := reader.Second[string, int, int](double)
//	result := r(tuple.MakeTuple2("hello", 5)) // ("hello", 10)
func Second[A, B, C any](pbc Reader[B, C]) Reader[T.Tuple2[A, B], T.Tuple2[A, C]] {
	return func(tab T.Tuple2[A, B]) T.Tuple2[A, C] {
		return T.MakeTuple2(tab.F1, pbc(tab.F2))
	}
}

// Read applies a context to a Reader to obtain its value.
// This is the "run" operation that executes a Reader with a specific environment.
//
// Note: Read is functionally identical to identity.Flap[A](e). Both take a value and
// return a function that applies that value to a function. The difference is semantic:
// - identity.Flap: Generic function application (applies value to any function)
// - reader.Read: Reader-specific execution (applies environment to a Reader)
//
// Recommendation: Use reader.Read when working in a Reader context, as it makes the
// intent clearer that you're executing a Reader computation with an environment.
// Use identity.Flap for general-purpose function application outside the Reader context.
//
// Example:
//
//	type Config struct { Port int }
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	run := reader.Read(Config{Port: 8080})
//	port := run(getPort) // 8080
//
//go:inline
func Read[A, E any](e E) func(Reader[E, A]) A {
	return I.Flap[A](e)
}

// MonadFlap is the monadic version of Flap.
// It takes a Reader containing a function and a value, and returns a Reader that applies the function to the value.
//
// Example:
//
//	type Config struct { Multiplier int }
//	getMultiplier := func(c Config) func(int) int {
//	    return N.Mul(c.Multiplier)
//	}
//	r := reader.MonadFlap(getMultiplier, 5)
//	result := r(Config{Multiplier: 3}) // 15
//
//go:inline
func MonadFlap[R, B, A any](fab Reader[R, func(A) B], a A) Reader[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap takes a value and returns a function that applies a Reader containing a function to that value.
// This is useful for partial application in the Reader context.
//
// Example:
//
//	type Config struct { Multiplier int }
//	getMultiplier := reader.Asks(func(c Config) func(int) int {
//	    return N.Mul(c.Multiplier)
//	})
//	applyTo5 := reader.Flap[Config](5)
//	r := applyTo5(getMultiplier)
//	result := r(Config{Multiplier: 3}) // 15
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}
