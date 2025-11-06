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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
)

// MonadAp applies an endomorphism to a value in a monadic context.
//
// This function applies the endomorphism fab to the value fa, returning the result.
// It's the monadic application operation for endomorphisms.
//
// Parameters:
//   - fab: An endomorphism to apply
//   - fa: The value to apply the endomorphism to
//
// Returns:
//   - The result of applying fab to fa
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	result := endomorphism.MonadAp(double, 5) // Returns: 10
func MonadAp[A any](fab Endomorphism[A], fa A) A {
	return identity.MonadAp(fab, fa)
}

// Ap returns a function that applies a value to an endomorphism.
//
// This is the curried version of MonadAp. It takes a value and returns a function
// that applies that value to any endomorphism.
//
// Parameters:
//   - fa: The value to be applied
//
// Returns:
//   - A function that takes an endomorphism and applies fa to it
//
// Example:
//
//	applyFive := endomorphism.Ap(5)
//	double := func(x int) int { return x * 2 }
//	result := applyFive(double) // Returns: 10
func Ap[A any](fa A) func(Endomorphism[A]) A {
	return identity.Ap[A](fa)
}

// Compose composes two endomorphisms into a single endomorphism.
//
// Given two endomorphisms f1 and f2, Compose returns a new endomorphism that
// applies f1 first, then applies f2 to the result. This is function composition:
// Compose(f1, f2)(x) = f2(f1(x))
//
// Composition is associative: Compose(Compose(f, g), h) = Compose(f, Compose(g, h))
//
// Parameters:
//   - f1: The first endomorphism to apply
//   - f2: The second endomorphism to apply
//
// Returns:
//   - A new endomorphism that is the composition of f1 and f2
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	increment := func(x int) int { return x + 1 }
//	doubleAndIncrement := endomorphism.Compose(double, increment)
//	result := doubleAndIncrement(5) // (5 * 2) + 1 = 11
func Compose[A any](f1, f2 Endomorphism[A]) Endomorphism[A] {
	return function.Flow2(f1, f2)
}

// MonadChain chains two endomorphisms together.
//
// This is the monadic bind operation for endomorphisms. It composes two endomorphisms
// ma and f, returning a new endomorphism that applies ma first, then f.
// MonadChain is equivalent to Compose.
//
// Parameters:
//   - ma: The first endomorphism in the chain
//   - f: The second endomorphism in the chain
//
// Returns:
//   - A new endomorphism that chains ma and f
//
// Example:
//
//	double := func(x int) int { return x * 2 }
//	increment := func(x int) int { return x + 1 }
//	chained := endomorphism.MonadChain(double, increment)
//	result := chained(5) // (5 * 2) + 1 = 11
func MonadChain[A any](ma Endomorphism[A], f Endomorphism[A]) Endomorphism[A] {
	return Compose(ma, f)
}

// Chain returns a function that chains an endomorphism with another.
//
// This is the curried version of MonadChain. It takes an endomorphism f and returns
// a function that chains any endomorphism with f.
//
// Parameters:
//   - f: The endomorphism to chain with
//
// Returns:
//   - A function that takes an endomorphism and chains it with f
//
// Example:
//
//	increment := func(x int) int { return x + 1 }
//	chainWithIncrement := endomorphism.Chain(increment)
//	double := func(x int) int { return x * 2 }
//	chained := chainWithIncrement(double)
//	result := chained(5) // (5 * 2) + 1 = 11
func Chain[A any](f Endomorphism[A]) Endomorphism[Endomorphism[A]] {
	return function.Bind2nd(MonadChain, f)
}
