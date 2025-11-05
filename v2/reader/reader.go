// Copyright (c) 2023 IBM Corp.
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
func AsksReader[R, A any](f func(R) Reader[R, A]) Reader[R, A] {
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
func MonadMap[E, A, B any](fa Reader[E, A], f func(A) B) Reader[E, B] {
	return function.Flow2(fa, f)
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
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return function.Bind2nd(MonadMap[E, A, B], f)
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
func MonadChain[R, A, B any](ma Reader[R, A], f func(A) Reader[R, B]) Reader[R, B] {
	return func(r R) B {
		return f(ma(r))(r)
	}
}

// Chain sequences two Reader computations where the second depends on the result of the first.
// This is the Monad operation that enables dependent computations.
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
func Chain[R, A, B any](f func(A) Reader[R, B]) Operator[R, A, B] {
	return function.Bind2nd(MonadChain[R, A, B], f)
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
func Flatten[R, A any](mma func(R) Reader[R, A]) Reader[R, A] {
	return MonadChain(mma, function.Identity[Reader[R, A]])
}

// Compose composes two Readers sequentially, where the output environment of the first
// becomes the input environment of the second.
//
// Example:
//
//	type Config struct { Port int }
//	type Env struct { Config Config }
//	getConfig := func(e Env) Config { return e.Config }
//	getPort := func(c Config) int { return c.Port }
//	getPortFromEnv := reader.Compose(getConfig)(getPort)
func Compose[R, B, C any](ab Reader[R, B]) func(Reader[B, C]) Reader[R, C] {
	return func(bc Reader[B, C]) Reader[R, C] {
		return function.Flow2(ab, bc)
	}
}

// First applies a Reader to the first element of a tuple, leaving the second element unchanged.
// This is useful for working with paired data where only one element needs transformation.
//
// Example:
//
//	double := func(x int) int { return x * 2 }
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
//	double := func(x int) int { return x * 2 }
//	r := reader.Second[string, int, int](double)
//	result := r(tuple.MakeTuple2("hello", 5)) // ("hello", 10)
func Second[A, B, C any](pbc Reader[B, C]) Reader[T.Tuple2[A, B], T.Tuple2[A, C]] {
	return func(tab T.Tuple2[A, B]) T.Tuple2[A, C] {
		return T.MakeTuple2(tab.F1, pbc(tab.F2))
	}
}

// Promap is the profunctor map operation that transforms both the input and output of a Reader.
// It applies f to the input (contravariantly) and g to the output (covariantly).
//
// Example:
//
//	type Config struct { Port int }
//	type Env struct { Config Config }
//	getPort := func(c Config) int { return c.Port }
//	extractConfig := func(e Env) Config { return e.Config }
//	toString := func(i int) string { return strconv.Itoa(i) }
//	r := reader.Promap(extractConfig, toString)(getPort)
//	result := r(Env{Config: Config{Port: 8080}}) // "8080"
func Promap[E, A, D, B any](f func(D) E, g func(A) B) func(Reader[E, A]) Reader[D, B] {
	return func(fea Reader[E, A]) Reader[D, B] {
		return function.Flow3(f, fea, g)
	}
}

// Local changes the value of the local context during the execution of the action `ma`.
// This is similar to Contravariant's contramap and allows you to modify the environment
// before passing it to a Reader.
//
// Example:
//
//	type DetailedConfig struct { Host string; Port int }
//	type SimpleConfig struct { Host string }
//	getHost := func(c SimpleConfig) string { return c.Host }
//	simplify := func(d DetailedConfig) SimpleConfig { return SimpleConfig{Host: d.Host} }
//	r := reader.Local(simplify)(getHost)
//	result := r(DetailedConfig{Host: "localhost", Port: 8080}) // "localhost"
func Local[R2, R1, A any](f func(R2) R1) func(Reader[R1, A]) Reader[R2, A] {
	return Compose[R2, R1, A](f)
}

// Read applies a context to a Reader to obtain its value.
// This is the "run" operation that executes a Reader with a specific environment.
//
// Example:
//
//	type Config struct { Port int }
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	run := reader.Read(Config{Port: 8080})
//	port := run(getPort) // 8080
func Read[E, A any](e E) func(Reader[E, A]) A {
	return I.Ap[A](e)
}

// MonadFlap is the monadic version of Flap.
// It takes a Reader containing a function and a value, and returns a Reader that applies the function to the value.
//
// Example:
//
//	type Config struct { Multiplier int }
//	getMultiplier := func(c Config) func(int) int {
//	    return func(x int) int { return x * c.Multiplier }
//	}
//	r := reader.MonadFlap(getMultiplier, 5)
//	result := r(Config{Multiplier: 3}) // 15
func MonadFlap[R, A, B any](fab Reader[R, func(A) B], a A) Reader[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap takes a value and returns a function that applies a Reader containing a function to that value.
// This is useful for partial application in the Reader context.
//
// Example:
//
//	type Config struct { Multiplier int }
//	getMultiplier := reader.Asks(func(c Config) func(int) int {
//	    return func(x int) int { return x * c.Multiplier }
//	})
//	applyTo5 := reader.Flap[Config](5)
//	r := applyTo5(getMultiplier)
//	result := r(Config{Multiplier: 3}) // 15
func Flap[R, A, B any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}
