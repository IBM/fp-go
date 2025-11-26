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
//	double := N.Mul(2)
//	increment := N.Add(1)
//
//	// Compose them (RIGHT-TO-LEFT execution)
//	composed := endomorphism.Compose(double, increment)
//	result := composed(5) // increment(5) then double: (5 + 1) * 2 = 12
//
//	// Chain them (LEFT-TO-RIGHT execution)
//	chained := endomorphism.MonadChain(double, increment)
//	result2 := chained(5) // double(5) then increment: (5 * 2) + 1 = 11
//
// # Monoid Operations
//
// Endomorphisms form a monoid, which means you can combine multiple endomorphisms.
// The monoid uses Compose, which executes RIGHT-TO-LEFT:
//
//	import (
//		"github.com/IBM/fp-go/v2/endomorphism"
//		M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	// Get the monoid for int endomorphisms
//	monoid := endomorphism.Monoid[int]()
//
//	// Combine multiple endomorphisms (RIGHT-TO-LEFT execution)
//	combined := M.ConcatAll(monoid)(
//		N.Mul(2),  // applied third
//		N.Add(1),  // applied second
//		N.Mul(3),  // applied first
//	)
//	result := combined(5) // (5 * 3) = 15, (15 + 1) = 16, (16 * 2) = 32
//
// # Monad Operations
//
// The package also provides monadic operations for endomorphisms.
// MonadChain executes LEFT-TO-RIGHT, unlike Compose:
//
//	// Chain allows sequencing of endomorphisms (LEFT-TO-RIGHT)
//	f := N.Mul(2)
//	g := N.Add(1)
//	chained := endomorphism.MonadChain(f, g)  // f first, then g
//	result := chained(5) // (5 * 2) + 1 = 11
//
// # Compose vs Chain
//
// The key difference between Compose and Chain/MonadChain is execution order:
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//
//	// Compose: RIGHT-TO-LEFT (mathematical composition)
//	composed := endomorphism.Compose(double, increment)
//	result1 := composed(5) // increment(5) * 2 = (5 + 1) * 2 = 12
//
//	// MonadChain: LEFT-TO-RIGHT (sequential application)
//	chained := endomorphism.MonadChain(double, increment)
//	result2 := chained(5) // double(5) + 1 = (5 * 2) + 1 = 11
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
