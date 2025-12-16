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

package nonempty

import (
	G "github.com/IBM/fp-go/v2/array/generic"
	EM "github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Of constructs a single element array
func Of[A any](first A) NonEmptyArray[A] {
	return G.Of[NonEmptyArray[A]](first)
}

// From constructs a [NonEmptyArray] from a set of variadic arguments
func From[A any](first A, data ...A) NonEmptyArray[A] {
	count := len(data)
	if count == 0 {
		return Of(first)
	}
	// allocate the requested buffer
	buffer := make(NonEmptyArray[A], count+1)
	buffer[0] = first
	copy(buffer[1:], data)
	return buffer
}

//go:inline
func IsEmpty[A any](_ NonEmptyArray[A]) bool {
	return false
}

//go:inline
func IsNonEmpty[A any](_ NonEmptyArray[A]) bool {
	return true
}

//go:inline
func MonadMap[A, B any](as NonEmptyArray[A], f func(a A) B) NonEmptyArray[B] {
	return G.MonadMap[NonEmptyArray[A], NonEmptyArray[B]](as, f)
}

//go:inline
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return G.Map[NonEmptyArray[A], NonEmptyArray[B]](f)
}

func Reduce[A, B any](f func(B, A) B, initial B) func(NonEmptyArray[A]) B {
	return func(as NonEmptyArray[A]) B {
		return array.Reduce(as, f, initial)
	}
}

func ReduceRight[A, B any](f func(A, B) B, initial B) func(NonEmptyArray[A]) B {
	return func(as NonEmptyArray[A]) B {
		return array.ReduceRight(as, f, initial)
	}
}

//go:inline
func Tail[A any](as NonEmptyArray[A]) []A {
	return as[1:]
}

//go:inline
func Head[A any](as NonEmptyArray[A]) A {
	return as[0]
}

//go:inline
func First[A any](as NonEmptyArray[A]) A {
	return as[0]
}

//go:inline
func Last[A any](as NonEmptyArray[A]) A {
	return as[len(as)-1]
}

//go:inline
func Size[A any](as NonEmptyArray[A]) int {
	return G.Size(as)
}

func Flatten[A any](mma NonEmptyArray[NonEmptyArray[A]]) NonEmptyArray[A] {
	return G.Flatten(mma)
}

func MonadChain[A, B any](fa NonEmptyArray[A], f Kleisli[A, B]) NonEmptyArray[B] {
	return G.MonadChain(fa, f)
}

func Chain[A, B any](f func(A) NonEmptyArray[B]) Operator[A, B] {
	return G.Chain[NonEmptyArray[A]](f)
}

func MonadAp[B, A any](fab NonEmptyArray[func(A) B], fa NonEmptyArray[A]) NonEmptyArray[B] {
	return G.MonadAp[NonEmptyArray[B]](fab, fa)
}

func Ap[B, A any](fa NonEmptyArray[A]) func(NonEmptyArray[func(A) B]) NonEmptyArray[B] {
	return G.Ap[NonEmptyArray[B], NonEmptyArray[func(A) B]](fa)
}

// FoldMap maps and folds a [NonEmptyArray]. Map the [NonEmptyArray] passing each value to the iterating function. Then fold the results using the provided [Semigroup].
func FoldMap[A, B any](s S.Semigroup[B]) func(func(A) B) func(NonEmptyArray[A]) B {
	return func(f func(A) B) func(NonEmptyArray[A]) B {
		return func(as NonEmptyArray[A]) B {
			return array.Reduce(Tail(as), func(cur B, a A) B {
				return s.Concat(cur, f(a))
			}, f(Head(as)))
		}
	}
}

// Fold folds the [NonEmptyArray] using the provided [Semigroup].
func Fold[A any](s S.Semigroup[A]) func(NonEmptyArray[A]) A {
	return func(as NonEmptyArray[A]) A {
		return array.Reduce(Tail(as), s.Concat, Head(as))
	}
}

// Prepend prepends a single value to an array
func Prepend[A any](head A) EM.Endomorphism[NonEmptyArray[A]] {
	return array.Prepend[EM.Endomorphism[NonEmptyArray[A]]](head)
}

// ToNonEmptyArray attempts to convert a regular slice into a NonEmptyArray.
// This function provides a safe way to create a NonEmptyArray from a slice that might be empty,
// returning an Option type to handle the case where the input slice is empty.
//
// Type Parameters:
//   - A: The element type of the array
//
// Parameters:
//   - as: A regular slice that may or may not be empty
//
// Returns:
//   - Option[NonEmptyArray[A]]: Some(NonEmptyArray) if the input slice is non-empty, None if empty
//
// Behavior:
//   - If the input slice is empty, returns None
//   - If the input slice has at least one element, wraps it in Some and returns it as a NonEmptyArray
//   - The conversion is a type cast, so no data is copied
//
// Example:
//
//	// Convert non-empty slice
//	numbers := []int{1, 2, 3}
//	result := ToNonEmptyArray(numbers)  // Some(NonEmptyArray[1, 2, 3])
//
//	// Convert empty slice
//	empty := []int{}
//	result := ToNonEmptyArray(empty)  // None
//
//	// Use with Option methods
//	numbers := []int{1, 2, 3}
//	result := ToNonEmptyArray(numbers)
//	if O.IsSome(result) {
//	    nea := O.GetOrElse(F.Constant(From(0)))(result)
//	    head := Head(nea)  // 1
//	}
//
// Use cases:
//   - Safely converting user input or external data to NonEmptyArray
//   - Validating that a collection has at least one element before processing
//   - Converting results from functions that return regular slices
//   - Ensuring type safety when working with collections that must not be empty
//
// Example with validation:
//
//	func processItems(items []string) Option[string] {
//	    return F.Pipe2(
//	        items,
//	        ToNonEmptyArray[string],
//	        O.Map(func(nea NonEmptyArray[string]) string {
//	            return Head(nea)  // Safe to get head since we know it's non-empty
//	        }),
//	    )
//	}
//
// Example with error handling:
//
//	items := []int{1, 2, 3}
//	result := ToNonEmptyArray(items)
//	switch {
//	case O.IsSome(result):
//	    nea := O.GetOrElse(F.Constant(From(0)))(result)
//	    fmt.Println("First item:", Head(nea))
//	case O.IsNone(result):
//	    fmt.Println("Array is empty")
//	}
//
// Example with chaining:
//
//	// Process only if non-empty
//	result := F.Pipe3(
//	    []int{1, 2, 3},
//	    ToNonEmptyArray[int],
//	    O.Map(Map(func(x int) int { return x * 2 })),
//	    O.Map(Head[int]),
//	)  // Some(2)
//
// Note: This function is particularly useful when working with APIs or functions
// that return regular slices but you need the type-level guarantee that the
// collection is non-empty for subsequent operations.
func ToNonEmptyArray[A any](as []A) Option[NonEmptyArray[A]] {
	if G.IsEmpty(as) {
		return option.None[NonEmptyArray[A]]()
	}
	return option.Some(NonEmptyArray[A](as))
}
