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

// Of constructs a single element NonEmptyArray.
// This is the simplest way to create a NonEmptyArray with exactly one element.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - first: The single element to include in the array
//
// Returns:
//   - NonEmptyArray[A]: A NonEmptyArray containing only the provided element
//
// Example:
//
//	arr := Of(42)           // NonEmptyArray[int]{42}
//	str := Of("hello")      // NonEmptyArray[string]{"hello"}
func Of[A any](first A) NonEmptyArray[A] {
	return G.Of[NonEmptyArray[A]](first)
}

// From constructs a NonEmptyArray from a set of variadic arguments.
// The first argument is required to ensure the array is non-empty, and additional
// elements can be provided as variadic arguments.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - first: The first element (required to ensure non-emptiness)
//   - data: Additional elements (optional)
//
// Returns:
//   - NonEmptyArray[A]: A NonEmptyArray containing all provided elements
//
// Example:
//
//	arr1 := From(1)              // NonEmptyArray[int]{1}
//	arr2 := From(1, 2, 3)        // NonEmptyArray[int]{1, 2, 3}
//	arr3 := From("a", "b", "c")  // NonEmptyArray[string]{"a", "b", "c"}
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

// IsEmpty always returns false for NonEmptyArray since it's guaranteed to have at least one element.
// This function exists for API consistency with regular arrays.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - _: The NonEmptyArray (unused, as the result is always false)
//
// Returns:
//   - bool: Always false
//
//go:inline
func IsEmpty[A any](_ NonEmptyArray[A]) bool {
	return false
}

// IsNonEmpty always returns true for NonEmptyArray since it's guaranteed to have at least one element.
// This function exists for API consistency with regular arrays.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - _: The NonEmptyArray (unused, as the result is always true)
//
// Returns:
//   - bool: Always true
//
//go:inline
func IsNonEmpty[A any](_ NonEmptyArray[A]) bool {
	return true
}

// MonadMap applies a function to each element of a NonEmptyArray, returning a new NonEmptyArray with the results.
// This is the monadic version of Map that takes the array as the first parameter.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - as: The input NonEmptyArray
//   - f: The function to apply to each element
//
// Returns:
//   - NonEmptyArray[B]: A new NonEmptyArray with the transformed elements
//
// Example:
//
//	arr := From(1, 2, 3)
//	doubled := MonadMap(arr, func(x int) int { return x * 2 })  // NonEmptyArray[int]{2, 4, 6}
//
//go:inline
func MonadMap[A, B any](as NonEmptyArray[A], f func(a A) B) NonEmptyArray[B] {
	return G.MonadMap[NonEmptyArray[A], NonEmptyArray[B]](as, f)
}

// Map applies a function to each element of a NonEmptyArray, returning a new NonEmptyArray with the results.
// This is the curried version that returns a function.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: The function to apply to each element
//
// Returns:
//   - Operator[A, B]: A function that transforms NonEmptyArray[A] to NonEmptyArray[B]
//
// Example:
//
//	double := Map(func(x int) int { return x * 2 })
//	result := double(From(1, 2, 3))  // NonEmptyArray[int]{2, 4, 6}
//
//go:inline
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return G.Map[NonEmptyArray[A], NonEmptyArray[B]](f)
}

// Reduce applies a function to each element of a NonEmptyArray from left to right,
// accumulating a result starting from an initial value.
//
// Type Parameters:
//   - A: The element type of the array
//   - B: The accumulator type
//
// Parameters:
//   - f: The reducer function that takes (accumulator, element) and returns a new accumulator
//   - initial: The initial value for the accumulator
//
// Returns:
//   - func(NonEmptyArray[A]) B: A function that reduces the array to a single value
//
// Example:
//
//	sum := Reduce(func(acc int, x int) int { return acc + x }, 0)
//	result := sum(From(1, 2, 3, 4))  // 10
//
//	concat := Reduce(func(acc string, x string) string { return acc + x }, "")
//	result := concat(From("a", "b", "c"))  // "abc"
func Reduce[A, B any](f func(B, A) B, initial B) func(NonEmptyArray[A]) B {
	return func(as NonEmptyArray[A]) B {
		return array.Reduce(as, f, initial)
	}
}

// ReduceRight applies a function to each element of a NonEmptyArray from right to left,
// accumulating a result starting from an initial value.
//
// Type Parameters:
//   - A: The element type of the array
//   - B: The accumulator type
//
// Parameters:
//   - f: The reducer function that takes (element, accumulator) and returns a new accumulator
//   - initial: The initial value for the accumulator
//
// Returns:
//   - func(NonEmptyArray[A]) B: A function that reduces the array to a single value
//
// Example:
//
//	concat := ReduceRight(func(x string, acc string) string { return acc + x }, "")
//	result := concat(From("a", "b", "c"))  // "cba"
func ReduceRight[A, B any](f func(A, B) B, initial B) func(NonEmptyArray[A]) B {
	return func(as NonEmptyArray[A]) B {
		return array.ReduceRight(as, f, initial)
	}
}

// Tail returns all elements of a NonEmptyArray except the first one.
// Returns an empty slice if the array has only one element.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - []A: A slice containing all elements except the first (may be empty)
//
// Example:
//
//	arr := From(1, 2, 3, 4)
//	tail := Tail(arr)  // []int{2, 3, 4}
//
//	single := From(1)
//	tail := Tail(single)  // []int{}
//
//go:inline
func Tail[A any](as NonEmptyArray[A]) []A {
	return as[1:]
}

// Head returns the first element of a NonEmptyArray.
// This operation is always safe since NonEmptyArray is guaranteed to have at least one element.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - A: The first element
//
// Example:
//
//	arr := From(1, 2, 3)
//	first := Head(arr)  // 1
//
//go:inline
func Head[A any](as NonEmptyArray[A]) A {
	return as[0]
}

// First returns the first element of a NonEmptyArray.
// This is an alias for Head.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - A: The first element
//
// Example:
//
//	arr := From(1, 2, 3)
//	first := First(arr)  // 1
//
//go:inline
func First[A any](as NonEmptyArray[A]) A {
	return as[0]
}

// Last returns the last element of a NonEmptyArray.
// This operation is always safe since NonEmptyArray is guaranteed to have at least one element.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - A: The last element
//
// Example:
//
//	arr := From(1, 2, 3)
//	last := Last(arr)  // 3
//
//go:inline
func Last[A any](as NonEmptyArray[A]) A {
	return as[len(as)-1]
}

// Size returns the number of elements in a NonEmptyArray.
// The result is always at least 1.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - int: The number of elements (always >= 1)
//
// Example:
//
//	arr := From(1, 2, 3)
//	size := Size(arr)  // 3
//
//go:inline
func Size[A any](as NonEmptyArray[A]) int {
	return G.Size(as)
}

// Flatten flattens a NonEmptyArray of NonEmptyArrays into a single NonEmptyArray.
// This operation concatenates all inner arrays into one.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - mma: A NonEmptyArray of NonEmptyArrays
//
// Returns:
//   - NonEmptyArray[A]: A flattened NonEmptyArray containing all elements
//
// Example:
//
//	nested := From(From(1, 2), From(3, 4), From(5))
//	flat := Flatten(nested)  // NonEmptyArray[int]{1, 2, 3, 4, 5}
func Flatten[A any](mma NonEmptyArray[NonEmptyArray[A]]) NonEmptyArray[A] {
	return G.Flatten(mma)
}

// MonadChain applies a function that returns a NonEmptyArray to each element and flattens the results.
// This is the monadic bind operation (flatMap) that takes the array as the first parameter.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - fa: The input NonEmptyArray
//   - f: A function that takes an element and returns a NonEmptyArray
//
// Returns:
//   - NonEmptyArray[B]: The flattened result
//
// Example:
//
//	arr := From(1, 2, 3)
//	result := MonadChain(arr, func(x int) NonEmptyArray[int] {
//	    return From(x, x*10)
//	})  // NonEmptyArray[int]{1, 10, 2, 20, 3, 30}
func MonadChain[A, B any](fa NonEmptyArray[A], f Kleisli[A, B]) NonEmptyArray[B] {
	return G.MonadChain(fa, f)
}

// Chain applies a function that returns a NonEmptyArray to each element and flattens the results.
// This is the curried version of MonadChain.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that takes an element and returns a NonEmptyArray
//
// Returns:
//   - Operator[A, B]: A function that transforms NonEmptyArray[A] to NonEmptyArray[B]
//
// Example:
//
//	duplicate := Chain(func(x int) NonEmptyArray[int] { return From(x, x) })
//	result := duplicate(From(1, 2, 3))  // NonEmptyArray[int]{1, 1, 2, 2, 3, 3}
func Chain[A, B any](f func(A) NonEmptyArray[B]) Operator[A, B] {
	return G.Chain[NonEmptyArray[A]](f)
}

// MonadAp applies a NonEmptyArray of functions to a NonEmptyArray of values.
// Each function is applied to each value, producing a cartesian product of results.
//
// Type Parameters:
//   - B: The output element type
//   - A: The input element type
//
// Parameters:
//   - fab: A NonEmptyArray of functions
//   - fa: A NonEmptyArray of values
//
// Returns:
//   - NonEmptyArray[B]: The result of applying all functions to all values
//
// Example:
//
//	fns := From(func(x int) int { return x * 2 }, func(x int) int { return x + 10 })
//	vals := From(1, 2)
//	result := MonadAp(fns, vals)  // NonEmptyArray[int]{2, 4, 11, 12}
func MonadAp[B, A any](fab NonEmptyArray[func(A) B], fa NonEmptyArray[A]) NonEmptyArray[B] {
	return G.MonadAp[NonEmptyArray[B]](fab, fa)
}

// Ap applies a NonEmptyArray of functions to a NonEmptyArray of values.
// This is the curried version of MonadAp.
//
// Type Parameters:
//   - B: The output element type
//   - A: The input element type
//
// Parameters:
//   - fa: A NonEmptyArray of values
//
// Returns:
//   - func(NonEmptyArray[func(A) B]) NonEmptyArray[B]: A function that applies functions to the values
//
// Example:
//
//	vals := From(1, 2)
//	applyTo := Ap[int](vals)
//	fns := From(func(x int) int { return x * 2 }, func(x int) int { return x + 10 })
//	result := applyTo(fns)  // NonEmptyArray[int]{2, 4, 11, 12}
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

// Prepend prepends a single value to the beginning of a NonEmptyArray.
// Returns a new NonEmptyArray with the value at the front.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - head: The value to prepend
//
// Returns:
//   - EM.Endomorphism[NonEmptyArray[A]]: A function that prepends the value to a NonEmptyArray
//
// Example:
//
//	arr := From(2, 3, 4)
//	prepend1 := Prepend(1)
//	result := prepend1(arr)  // NonEmptyArray[int]{1, 2, 3, 4}
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

// Extract returns the first element of a NonEmptyArray.
// This is an alias for Head and is part of the Comonad interface.
//
// Type Parameters:
//   - A: The element type
//
// Parameters:
//   - as: The input NonEmptyArray
//
// Returns:
//   - A: The first element
//
// Example:
//
//	arr := From(1, 2, 3)
//	first := Extract(arr)  // 1
//
//go:inline
func Extract[A any](as NonEmptyArray[A]) A {
	return Head(as)
}

// Extend applies a function to all suffixes of a NonEmptyArray.
// For each position i, it applies the function to the subarray starting at position i.
// This is part of the Comonad interface.
//
// Type Parameters:
//   - A: The input element type
//   - B: The output element type
//
// Parameters:
//   - f: A function that takes a NonEmptyArray and returns a value
//
// Returns:
//   - Operator[A, B]: A function that transforms NonEmptyArray[A] to NonEmptyArray[B]
//
// Example:
//
//	arr := From(1, 2, 3, 4)
//	sumSuffix := Extend(func(xs NonEmptyArray[int]) int {
//	    sum := 0
//	    for _, x := range xs {
//	        sum += x
//	    }
//	    return sum
//	})
//	result := sumSuffix(arr)  // NonEmptyArray[int]{10, 9, 7, 4}
//	// [1,2,3,4] -> 10, [2,3,4] -> 9, [3,4] -> 7, [4] -> 4
//
//go:inline
func Extend[A, B any](f func(NonEmptyArray[A]) B) Operator[A, B] {
	return func(as NonEmptyArray[A]) NonEmptyArray[B] {
		return G.MakeBy[NonEmptyArray[B]](len(as), func(i int) B { return f(as[i:]) })
	}
}
