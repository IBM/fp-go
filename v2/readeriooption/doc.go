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

// Package readeriooption provides a monad transformer that combines the Reader, IO, and Option monads.
//
// # Overview
//
// ReaderIOOption[R, A] represents a computation that:
//   - Depends on a shared environment of type R (Reader monad)
//   - Performs side effects (IO monad)
//   - May fail to produce a value of type A (Option monad)
//
// This is particularly useful for computations that need:
//   - Dependency injection or configuration access
//   - Side effects like I/O operations
//   - Optional results without using error types
//
// The ReaderIOOption monad is defined as: Reader[R, IOOption[A]]
//
// # Fantasy Land Specification
//
// This package implements the following Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//   - Profunctor: https://github.com/fantasyland/fantasy-land#profunctor
//
// # Core Operations
//
// Creating ReaderIOOption values:
//   - Of/Some: Wraps a value in a successful ReaderIOOption
//   - None: Creates a ReaderIOOption representing no value
//   - FromOption: Lifts an Option into ReaderIOOption
//   - FromReader: Lifts a Reader into ReaderIOOption
//   - Ask/Asks: Accesses the environment
//
// Transforming values:
//   - Map: Transforms the value inside a ReaderIOOption
//   - Chain: Sequences ReaderIOOption computations
//   - Ap: Applies a function wrapped in ReaderIOOption
//   - Alt: Provides alternative computation on failure
//
// Extracting values:
//   - Fold: Extracts value by providing handlers for both cases
//   - GetOrElse: Returns value or default
//   - Read: Executes the computation with an environment
//
// # Basic Example
//
//	type Config struct {
//	    DatabaseURL string
//	    Timeout     int
//	}
//
//	// A computation that may or may not find a user
//	func findUser(id int) readeriooption.ReaderIOOption[Config, User] {
//	    return readeriooption.Asks(func(cfg Config) iooption.IOOption[User] {
//	        return func() option.Option[User] {
//	            // Use cfg.DatabaseURL to query database
//	            // Return Some(user) if found, None() if not found
//	            user, found := queryDB(cfg.DatabaseURL, id)
//	            if found {
//	                return option.Some(user)
//	            }
//	            return option.None[User]()
//	        }
//	    })
//	}
//
//	// Chain multiple operations
//	result := F.Pipe2(
//	    findUser(123),
//	    readeriooption.Chain(func(user User) readeriooption.ReaderIOOption[Config, Profile] {
//	        return loadProfile(user.ProfileID)
//	    }),
//	    readeriooption.Map(func(profile Profile) string {
//	        return profile.DisplayName
//	    }),
//	)
//
//	// Execute with config
//	config := Config{DatabaseURL: "localhost:5432", Timeout: 30}
//	displayName := result(config)() // Returns Option[string]
//
// # Do-Notation Style
//
// The package supports do-notation style composition for building complex computations:
//
//	type State struct {
//	    User    User
//	    Profile Profile
//	    Posts   []Post
//	}
//
//	result := F.Pipe3(
//	    readeriooption.Do[Config](State{}),
//	    readeriooption.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readeriooption.ReaderIOOption[Config, User] {
//	            return findUser(123)
//	        },
//	    ),
//	    readeriooption.Bind(
//	        func(profile Profile) func(State) State {
//	            return func(s State) State { s.Profile = profile; return s }
//	        },
//	        func(s State) readeriooption.ReaderIOOption[Config, Profile] {
//	            return loadProfile(s.User.ProfileID)
//	        },
//	    ),
//	    readeriooption.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readeriooption.ReaderIOOption[Config, []Post] {
//	            return loadPosts(s.User.ID)
//	        },
//	    ),
//	)
//
// # Alternative Computations
//
// Use Alt to provide fallback behavior when computations fail:
//
//	// Try cache first, fall back to database
//	result := F.Pipe1(
//	    findUserInCache(123),
//	    readeriooption.Alt(func() readeriooption.ReaderIOOption[Config, User] {
//	        return findUserInDB(123)
//	    }),
//	)
//
// # Array Operations
//
// Transform arrays where each element may fail:
//
//	userIDs := []int{1, 2, 3, 4, 5}
//	users := F.Pipe1(
//	    readeriooption.Of[Config](userIDs),
//	    readeriooption.Chain(readeriooption.TraverseArray[Config](findUser)),
//	)
//	// Returns Some([]User) if all users found, None otherwise
//
// # Monoid Operations
//
// Combine multiple ReaderIOOption computations:
//
//	import N "github.com/IBM/fp-go/v2/number"
//
//	// Applicative monoid - all must succeed
//	intAdd := N.MonoidSum[int]()
//	roMonoid := readeriooption.ApplicativeMonoid[Config](intAdd)
//	combined := roMonoid.Concat(
//	    readeriooption.Of[Config](5),
//	    readeriooption.Of[Config](3),
//	)
//	// Returns Some(8)
//
//	// Alternative monoid - provides fallback
//	altMonoid := readeriooption.AlternativeMonoid[Config](intAdd)
//	withFallback := altMonoid.Concat(
//	    readeriooption.None[Config, int](),
//	    readeriooption.Of[Config](10),
//	)
//	// Returns Some(10)
//
// # Profunctor Operations
//
// Transform both input and output:
//
//	type GlobalConfig struct {
//	    DB DBConfig
//	}
//
//	type DBConfig struct {
//	    Host string
//	}
//
//	// Adapt environment and transform result
//	adapted := F.Pipe1(
//	    queryDB, // ReaderIOOption[DBConfig, User]
//	    readeriooption.Promap(
//	        func(g GlobalConfig) DBConfig { return g.DB },
//	        func(u User) string { return u.Name },
//	    ),
//	)
//	// Now: ReaderIOOption[GlobalConfig, string]
//
// # Tail Recursion
//
// For recursive computations, use TailRec to avoid stack overflow:
//
//	func factorial(n int) readeriooption.ReaderIOOption[Config, int] {
//	    return readeriooption.TailRec(func(acc int) readeriooption.ReaderIOOption[Config, tailrec.Trampoline[int, int]] {
//	        if n <= 1 {
//	            return readeriooption.Of[Config](tailrec.Done[int](acc))
//	        }
//	        return readeriooption.Of[Config](tailrec.Continue[int](acc * n))
//	    })(1)
//	}
//
// # Relationship to Other Monads
//
// ReaderIOOption is related to other monads in the fp-go library:
//   - reader: ReaderIOOption adds IO and Option capabilities
//   - readerio: ReaderIOOption adds Option capability
//   - readeroption: ReaderIOOption adds IO capability
//   - iooption: ReaderIOOption adds Reader capability
//   - option: ReaderIOOption adds Reader and IO capabilities
//
// # Type Safety
//
// The type system ensures:
//   - Environment dependencies are explicit in the type signature
//   - Side effects are tracked through the IO layer
//   - Optional results are handled explicitly
//   - Composition maintains type safety
//
// # Performance Considerations
//
// ReaderIOOption computations are lazy and only execute when:
//  1. An environment is provided (Reader layer)
//  2. The IO action is invoked (IO layer)
//
// This allows for efficient composition without premature execution.
package readeriooption
