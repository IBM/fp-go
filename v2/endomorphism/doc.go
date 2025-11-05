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

// Package endomorphism provides functional programming utilities for working with endomorphisms.
//
// An endomorphism is a function from a type to itself: func(A) A. This package provides
// various operations and algebraic structures for composing and manipulating endomorphisms.
//
// # Core Concepts
//
// An Endomorphism[A] is simply a function that takes a value of type A and returns a value
// of the same type A. This simple concept has powerful algebraic properties:
//
//   - Identity: The identity function is an endomorphism
//   - Composition: Endomorphisms can be composed to form new endomorphisms
//   - Monoid: Endomorphisms form a monoid under composition with identity as the empty element
//
// # Basic Usage
//
// Creating and composing endomorphisms:
//
//	import (
//		"github.com/IBM/fp-go/v2/endomorphism"
//	)
//
//	// Define some endomorphisms
//	double := func(x int) int { return x * 2 }
//	increment := func(x int) int { return x + 1 }
//
//	// Compose them
//	doubleAndIncrement := endomorphism.Compose(double, increment)
//	result := doubleAndIncrement(5) // (5 * 2) + 1 = 11
//
// # Monoid Operations
//
// Endomorphisms form a monoid, which means you can combine multiple endomorphisms:
//
//	import (
//		"github.com/IBM/fp-go/v2/endomorphism"
//		M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	// Get the monoid for int endomorphisms
//	monoid := endomorphism.Monoid[int]()
//
//	// Combine multiple endomorphisms
//	combined := M.ConcatAll(monoid)(
//		func(x int) int { return x * 2 },
//		func(x int) int { return x + 1 },
//		func(x int) int { return x * 3 },
//	)
//	result := combined(5) // ((5 * 2) + 1) * 3 = 33
//
// # Monad Operations
//
// The package also provides monadic operations for endomorphisms:
//
//	// Chain allows sequencing of endomorphisms
//	f := func(x int) int { return x * 2 }
//	g := func(x int) int { return x + 1 }
//	chained := endomorphism.MonadChain(f, g)
//
// # Type Safety
//
// The package uses Go generics to ensure type safety. All operations preserve
// the type of the endomorphism, preventing type mismatches at compile time.
//
// # Related Packages
//
//   - function: Provides general function composition utilities
//   - identity: Provides identity functor operations
//   - monoid: Provides monoid algebraic structure
//   - semigroup: Provides semigroup algebraic structure
package endomorphism
