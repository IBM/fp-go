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
Package effect provides a functional effect system for service code with typed dependencies.

# Overview

An Effect[C, A] describes a computation that needs dependencies of type C, runs with a
context.Context, may perform I/O, and either fails with an error or succeeds with a value
of type A. It is built on top of ReaderReaderIOResult:

	Effect[C, A] = func(C) func(context.Context) func() Result[A]

# Dependencies and Context

The two environments of an Effect have distinct roles:
  - C carries the dependencies: long-lived collaborators such as repositories, database
    and HTTP clients, configuration or environment access. They are supplied once, with
    Provide.
  - context.Context carries the request scope: cancellation, deadlines and request-scoped
    values such as request IDs. It is supplied on every run, with RunSync.
  - Per-call input is a function argument: a Kleisli[C, A, B] is a func(A) Effect[C, B].

Declare dependencies in C, where the compiler checks that they are provided. Do not store
them in the context.Context and do not pass them as parameters through every layer.

# Naming Conventions

The naming conventions in this package are modeled after effect-ts (https://effect.website/),
a popular TypeScript library for functional effect systems. This alignment helps developers
familiar with effect-ts to quickly understand and use this Go implementation.

# Basic Operations

The examples below use this dependency type:

	type Deps struct {
		Greeting string
	}

	func getGreeting(d Deps) string { return d.Greeting }

Creating effects:

	effect.Succeed[Deps]("hello")                   // Effect[Deps, string]
	effect.Fail[Deps, string](errors.New("failed")) // Effect[Deps, string]
	effect.Of[Deps](42)                             // Effect[Deps, int]
	effect.Asks(getGreeting)                        // pure projection of the dependencies
	effect.Ask[Deps]()                              // the whole dependency value

Lifting Go functions that receive the dependencies and the context:

	// func(Deps, context.Context) (string, error) -> Effect[Deps, string]
	effect.Eitherize(loadGreeting)

	// func(Deps, context.Context, int) (User, error) -> Kleisli[Deps, int, User]
	effect.Eitherize1(findUser)

Transforming effects:

	// Map over the success value; C cannot be inferred, so it is given explicitly
	effect.Map[Deps](strconv.Itoa)

	// Chain effects together (flatMap)
	effect.Chain(effect.Eitherize1(findUser))

	// Run a side effect without changing the value
	effect.TapIOK[Deps](io.Logf[int]("value: %d"))

# Dependency Injection

Functions declare the narrowest dependency type they need. Local adapts an effect that
needs C2 to an environment C1, given a projection from C1 to C2:

	type App struct {
		Deps Deps
		Name string
	}

	func getDeps(a App) Deps { return a.Deps }

	// Effect[Deps, string] -> Effect[App, string]
	effect.Local[string](getDeps)(effect.Asks(getGreeting))

LocalReaderK, LocalIOK, LocalResultK, LocalIOResultK, LocalThunkK and LocalEffectK derive
the inner dependency with an effect, for example when it has to be loaded or validated.

# Do Notation

Do notation accumulates the results of several effects in a struct:

	type State struct {
		Greeting string
		Length   int
	}

	func setGreeting(g string) func(State) State {
		return func(s State) State { s.Greeting = g; return s }
	}

	func setLength(n int) func(State) State {
		return func(s State) State { s.Length = n; return s }
	}

	func greetingLength(s State) int { return len(s.Greeting) }

	result := F.Pipe2(
		effect.Do[Deps](State{}),
		effect.ApS(setGreeting, effect.Asks(getGreeting)),
		effect.Let[Deps](setLength, greetingLength),
	)

Generated lenses (see the lens command of the code generator) provide these setters, and
the L variants (BindL, ApSL, LetL, LetToL) accept a lens directly.

# Bind Operations

The package provides various bind operations for integrating with other effect types:

  - BindIOK: Bind an IO operation
  - BindIOResultK: Bind an IOResult operation
  - BindReaderK: Bind a Reader operation
  - BindReaderIOK: Bind a ReaderIO operation
  - BindResultK: Bind a Result operation

Each bind operation has a corresponding "L" variant for working with lenses:
  - BindL, BindIOKL, BindReaderKL, etc.

# Applicative Operations

Combine independent effects:

	// Apply a function effect to a value effect: Ap[B, C, A]
	effect.Ap[string](valueEffect)(functionEffect)

	// Add an independent effect's result to the do-notation state
	effect.ApS(setter, effect1)

# Traversal

Traverse collections with effects:

	// Map an array with an effectful function
	effect.TraverseArray(F.Flow2(strconv.Itoa, effect.Succeed[Deps, string]))

# Retry Logic

Retry effects with configurable policies:

	effect.Retrying(
		retry.LimitRetries(3),
		func(retry.RetryStatus) effect.Effect[Deps, string] {
			return fetchData()
		},
		result.IsLeft[string], // retry on error
	)

# Monoids

Combine effects using monoid operations:

	// Combine effects using applicative semantics
	effect.ApplicativeMonoid[Deps](S.Monoid)

	// Combine effects using alternative semantics (first success)
	effect.AlternativeMonoid[Deps](S.Monoid)

# Running Effects

Provide the dependencies once, then run the resulting thunk with a context.Context:

	// Provide the dependencies; A cannot be inferred, so it is given explicitly
	thunk := effect.Provide[string](deps)(myEffect)

	// Run synchronously
	value, err := effect.RunSync(thunk)(ctx)

# Integration with Other Packages

The effect package integrates seamlessly with other fp-go packages:
  - either: For error handling
  - io: For I/O operations
  - reader: For dependency injection
  - result: For result types
  - retry: For retry logic
  - monoid: For combining effects

# Example

A repository is a dependency. The interface method is lifted with Eitherize1, which turns
the receiver into the dependency, and Local fetches it from the application's dependencies:

	type UserRepo interface {
		FindUser(ctx context.Context, id int) (User, error)
	}

	type Deps struct {
		Users UserRepo
	}

	func getUsers(d Deps) UserRepo { return d.Users }

	func fetchUser() effect.Kleisli[Deps, int, User] {
		return F.Flow2(
			effect.Eitherize1(UserRepo.FindUser), // Kleisli[UserRepo, int, User]
			effect.Local[User](getUsers),         // Kleisli[Deps, int, User]
		)
	}

	func main() {
		deps := Deps{Users: newUserRepo()}

		user, err := effect.RunSync(effect.Provide[User](deps)(fetchUser()(42)))(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %+v\n", user)
	}
*/
package effect

//go:generate go run ../main.go lens --dir . --filename gen_lens.go --include-test-files
