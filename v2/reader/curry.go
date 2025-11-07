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
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// These functions curry a Go function with the context as the first parameter into a Reader
// with the context as the last parameter, which is equivalent to a function returning a Reader
// of that context.
//
// This follows the Go convention (https://pkg.go.dev/context) of putting the context as the
// first parameter, while Reader monad convention has the context as the last parameter.

// Curry0 converts a function that takes a context and returns a value into a Reader.
//
// Example:
//
//	type Config struct { Value int }
//	getValue := func(c Config) int { return c.Value }
//	r := reader.Curry0(getValue)
//	result := r(Config{Value: 42}) // 42
func Curry0[R, A any](f func(R) A) Reader[R, A] {
	return G.Curry0[Reader[R, A]](f)
}

// Curry1 converts a function with context as first parameter into a curried function
// returning a Reader. The context parameter is moved to the end (Reader position).
//
// Example:
//
//	type Config struct { Prefix string }
//	addPrefix := func(c Config, s string) string { return c.Prefix + s }
//	curried := reader.Curry1(addPrefix)
//	r := curried("hello")
//	result := r(Config{Prefix: ">> "}) // ">> hello"
func Curry1[R, T1, A any](f func(R, T1) A) Kleisli[R, T1, A] {
	return G.Curry1[Reader[R, A]](f)
}

// Curry2 converts a function with context as first parameter and 2 other parameters
// into a curried function returning a Reader.
//
// Example:
//
//	type Config struct { Sep string }
//	join := func(c Config, a, b string) string { return a + c.Sep + b }
//	curried := reader.Curry2(join)
//	r := curried("hello")("world")
//	result := r(Config{Sep: "-"}) // "hello-world"
func Curry2[R, T1, T2, A any](f func(R, T1, T2) A) func(T1) func(T2) Reader[R, A] {
	return G.Curry2[Reader[R, A]](f)
}

// Curry3 converts a function with context as first parameter and 3 other parameters
// into a curried function returning a Reader.
//
// Example:
//
//	type Config struct { Format string }
//	format := func(c Config, a, b, d string) string {
//	    return fmt.Sprintf(c.Format, a, b, d)
//	}
//	curried := reader.Curry3(format)
//	r := curried("a")("b")("c")
//	result := r(Config{Format: "%s-%s-%s"}) // "a-b-c"
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) A) func(T1) func(T2) func(T3) Reader[R, A] {
	return G.Curry3[Reader[R, A]](f)
}

// Curry4 converts a function with context as first parameter and 4 other parameters
// into a curried function returning a Reader.
//
// Example:
//
//	type Config struct { Multiplier int }
//	sum := func(c Config, a, b, d, e int) int {
//	    return (a + b + d + e) * c.Multiplier
//	}
//	curried := reader.Curry4(sum)
//	r := curried(1)(2)(3)(4)
//	result := r(Config{Multiplier: 10}) // 100
func Curry4[R, T1, T2, T3, T4, A any](f func(R, T1, T2, T3, T4) A) func(T1) func(T2) func(T3) func(T4) Reader[R, A] {
	return G.Curry4[Reader[R, A]](f)
}

// Uncurry0 converts a Reader back into a regular function with context as first parameter.
//
// Example:
//
//	type Config struct { Value int }
//	r := reader.Of[Config](42)
//	f := reader.Uncurry0(r)
//	result := f(Config{Value: 0}) // 42
func Uncurry0[R, A any](f Reader[R, A]) func(R) A {
	return G.Uncurry0(f)
}

// Uncurry1 converts a curried function returning a Reader back into a regular function
// with context as first parameter.
//
// Example:
//
//	type Config struct { Prefix string }
//	curried := func(s string) reader.Reader[Config, string] {
//	    return reader.Asks(func(c Config) string { return c.Prefix + s })
//	}
//	f := reader.Uncurry1(curried)
//	result := f(Config{Prefix: ">> "}, "hello") // ">> hello"
func Uncurry1[R, T1, A any](f Kleisli[R, T1, A]) func(R, T1) A {
	return G.Uncurry1(f)
}

// Uncurry2 converts a curried function with 2 parameters returning a Reader back into
// a regular function with context as first parameter.
//
// Example:
//
//	type Config struct { Sep string }
//	curried := func(a string) func(string) reader.Reader[Config, string] {
//	    return func(b string) reader.Reader[Config, string] {
//	        return reader.Asks(func(c Config) string { return a + c.Sep + b })
//	    }
//	}
//	f := reader.Uncurry2(curried)
//	result := f(Config{Sep: "-"}, "hello", "world") // "hello-world"
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) Reader[R, A]) func(R, T1, T2) A {
	return G.Uncurry2(f)
}

// Uncurry3 converts a curried function with 3 parameters returning a Reader back into
// a regular function with context as first parameter.
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) Reader[R, A]) func(R, T1, T2, T3) A {
	return G.Uncurry3(f)
}

// Uncurry4 converts a curried function with 4 parameters returning a Reader back into
// a regular function with context as first parameter.
func Uncurry4[R, T1, T2, T3, T4, A any](f func(T1) func(T2) func(T3) func(T4) Reader[R, A]) func(R, T1, T2, T3, T4) A {
	return G.Uncurry4(f)
}
