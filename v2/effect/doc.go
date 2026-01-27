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

/*
Package effect provides a functional effect system for managing side effects in Go.

# Overview

The effect package is a high-level abstraction for composing effectful computations
that may fail, require dependencies (context), and perform I/O operations. It is built
on top of ReaderReaderIOResult, providing a clean API for dependency injection and
error handling.

# Naming Conventions

The naming conventions in this package are modeled after effect-ts (https://effect.website/),
a popular TypeScript library for functional effect systems. This alignment helps developers
familiar with effect-ts to quickly understand and use this Go implementation.

# Core Type

The central type is Effect[C, A], which represents:
  - C: The context/dependency type required by the effect
  - A: The success value type produced by the effect

An Effect can:
  - Succeed with a value of type A
  - Fail with an error
  - Require a context of type C
  - Perform I/O operations

# Basic Operations

Creating Effects:

	// Create a successful effect
	effect.Succeed[MyContext, string]("hello")

	// Create a failed effect
	effect.Fail[MyContext, string](errors.New("failed"))

	// Lift a pure value into an effect
	effect.Of[MyContext, int](42)

Transforming Effects:

	// Map over the success value
	effect.Map[MyContext](func(x int) string {
		return strconv.Itoa(x)
	})

	// Chain effects together (flatMap)
	effect.Chain[MyContext](func(x int) Effect[MyContext, string] {
		return effect.Succeed[MyContext, string](strconv.Itoa(x))
	})

	// Tap into an effect without changing its value
	effect.Tap[MyContext](func(x int) Effect[MyContext, any] {
		return effect.Succeed[MyContext, any](fmt.Println(x))
	})

# Dependency Injection

Effects can access their required context:

	// Transform the context before passing it to an effect
	effect.Local[OuterCtx, InnerCtx](func(outer OuterCtx) InnerCtx {
		return outer.Inner
	})

	// Provide a context to run an effect
	effect.Provide[MyContext, string](myContext)

# Do Notation

The package provides "do notation" for composing effects in a sequential, imperative style:

	type State struct {
		X int
		Y string
	}

	result := effect.Do[MyContext](State{}).
		Bind(func(y string) func(State) State {
			return func(s State) State {
				s.Y = y
				return s
			}
		}, fetchString).
		Let(func(x int) func(State) State {
			return func(s State) State {
				s.X = x
				return s
			}
		}, func(s State) int {
			return len(s.Y)
		})

# Bind Operations

The package provides various bind operations for integrating with other effect types:

  - BindIOK: Bind an IO operation
  - BindIOEitherK: Bind an IOEither operation
  - BindIOResultK: Bind an IOResult operation
  - BindReaderK: Bind a Reader operation
  - BindReaderIOK: Bind a ReaderIO operation
  - BindEitherK: Bind an Either operation

Each bind operation has a corresponding "L" variant for working with lenses:
  - BindL, BindIOKL, BindReaderKL, etc.

# Applicative Operations

Apply effects in parallel:

	// Apply a function effect to a value effect
	effect.Ap[string, MyContext](valueEffect)(functionEffect)

	// Apply effects to build up a structure
	effect.ApS[MyContext](setter, effect1)

# Traversal

Traverse collections with effects:

	// Map an array with an effectful function
	effect.TraverseArray[MyContext](func(x int) Effect[MyContext, string] {
		return effect.Succeed[MyContext, string](strconv.Itoa(x))
	})

# Retry Logic

Retry effects with configurable policies:

	effect.Retrying[MyContext, string](
		retryPolicy,
		func(status retry.RetryStatus) Effect[MyContext, string] {
			return fetchData()
		},
		func(result Result[string]) bool {
			return result.IsLeft() // retry on error
		},
	)

# Monoids

Combine effects using monoid operations:

	// Combine effects using applicative semantics
	effect.ApplicativeMonoid[MyContext](stringMonoid)

	// Combine effects using alternative semantics (first success)
	effect.AlternativeMonoid[MyContext](stringMonoid)

# Running Effects

To execute an effect:

	// Provide the context
	ioResult := effect.Provide[MyContext, string](myContext)(myEffect)

	// Run synchronously
	readerResult := effect.RunSync(ioResult)

	// Execute with a context.Context
	value, err := readerResult(ctx)

# Integration with Other Packages

The effect package integrates seamlessly with other fp-go packages:
  - either: For error handling
  - io: For I/O operations
  - reader: For dependency injection
  - result: For result types
  - retry: For retry logic
  - monoid: For combining effects

# Example

	type Config struct {
		APIKey string
		BaseURL string
	}

	func fetchUser(id int) Effect[Config, User] {
		return effect.Chain[Config](func(cfg Config) Effect[Config, User] {
			// Use cfg.APIKey and cfg.BaseURL
			return effect.Succeed[Config, User](User{ID: id})
		})(effect.Of[Config, Config](Config{}))
	}

	func main() {
		cfg := Config{APIKey: "key", BaseURL: "https://api.example.com"}
		userEffect := fetchUser(42)

		// Run the effect
		ioResult := effect.Provide(cfg)(userEffect)
		readerResult := effect.RunSync(ioResult)
		user, err := readerResult(context.Background())

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %+v\n", user)
	}
*/
package effect

//go:generate go run ../main.go lens --dir . --filename gen_lens.go --include-test-files
