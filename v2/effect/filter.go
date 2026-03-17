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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/option"
)

// Filter lifts a filtering operation on a higher-kinded type into an Effect operator.
// This is a generic function that works with any filterable data structure by taking
// a filter function and returning an operator that can be used in effect chains.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - HKTA: The higher-kinded type being filtered (e.g., []A, Seq[A])
//   - A: The element type being filtered
//
// # Parameters
//
//   - filter: A function that takes a predicate and returns an endomorphism on HKTA
//
// # Returns
//
//   - func(Predicate[A]) Operator[C, HKTA, HKTA]: A function that takes a predicate
//     and returns an operator that filters effects containing HKTA values
//
// # Example Usage
//
//	import A "github.com/IBM/fp-go/v2/array"
//
//	// Create a custom filter operator for arrays
//	filterOp := Filter[MyContext](A.Filter[int])
//	isEven := func(n int) bool { return n%2 == 0 }
//
//	pipeline := F.Pipe2(
//	    Succeed[MyContext]([]int{1, 2, 3, 4, 5}),
//	    filterOp(isEven),
//	    Map[MyContext](func(arr []int) int { return len(arr) }),
//	)
//	// Result: Effect that produces 2 (count of even numbers)
//
// # See Also
//
//   - FilterArray: Specialized version for array filtering
//   - FilterIter: Specialized version for iterator filtering
//   - FilterMap: For filtering and mapping simultaneously
//
//go:inline
func Filter[C, HKTA, A any](
	filter func(Predicate[A]) Endomorphism[HKTA],
) func(Predicate[A]) Operator[C, HKTA, HKTA] {
	return readerreaderioresult.Filter[C](filter)
}

// FilterArray creates an operator that filters array elements within an Effect based on a predicate.
// Elements that satisfy the predicate are kept, while others are removed.
// This is a specialized version of Filter for arrays.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The element type in the array
//
// # Parameters
//
//   - p: A predicate function that tests each element
//
// # Returns
//
//   - Operator[C, []A, []A]: An operator that filters array elements in an effect
//
// # Example Usage
//
//	isPositive := func(n int) bool { return n > 0 }
//	filterPositive := FilterArray[MyContext](isPositive)
//
//	pipeline := F.Pipe1(
//	    Succeed[MyContext]([]int{-2, -1, 0, 1, 2, 3}),
//	    filterPositive,
//	)
//	// Result: Effect that produces []int{1, 2, 3}
//
// # See Also
//
//   - Filter: Generic version for any filterable type
//   - FilterIter: For filtering iterators
//   - FilterMapArray: For filtering and mapping arrays simultaneously
//
//go:inline
func FilterArray[C, A any](p Predicate[A]) Operator[C, []A, []A] {
	return readerreaderioresult.FilterArray[C](p)
}

// FilterIter creates an operator that filters iterator elements within an Effect based on a predicate.
// Elements that satisfy the predicate are kept in the resulting iterator, while others are removed.
// This is a specialized version of Filter for iterators (Seq).
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The element type in the iterator
//
// # Parameters
//
//   - p: A predicate function that tests each element
//
// # Returns
//
//   - Operator[C, Seq[A], Seq[A]]: An operator that filters iterator elements in an effect
//
// # Example Usage
//
//	isEven := func(n int) bool { return n%2 == 0 }
//	filterEven := FilterIter[MyContext](isEven)
//
//	pipeline := F.Pipe1(
//	    Succeed[MyContext](slices.Values([]int{1, 2, 3, 4, 5, 6})),
//	    filterEven,
//	)
//	// Result: Effect that produces an iterator over [2, 4, 6]
//
// # See Also
//
//   - Filter: Generic version for any filterable type
//   - FilterArray: For filtering arrays
//   - FilterMapIter: For filtering and mapping iterators simultaneously
//
//go:inline
func FilterIter[C, A any](p Predicate[A]) Operator[C, Seq[A], Seq[A]] {
	return readerreaderioresult.FilterIter[C](p)
}

// FilterMap lifts a filter-map operation on a higher-kinded type into an Effect operator.
// This combines filtering and mapping in a single operation: elements are transformed
// using a function that returns Option, and only Some values are kept in the result.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - HKTA: The input higher-kinded type (e.g., []A, Seq[A])
//   - HKTB: The output higher-kinded type (e.g., []B, Seq[B])
//   - A: The input element type
//   - B: The output element type
//
// # Parameters
//
//   - filter: A function that takes an option.Kleisli and returns a transformation from HKTA to HKTB
//
// # Returns
//
//   - func(option.Kleisli[A, B]) Operator[C, HKTA, HKTB]: A function that takes a Kleisli arrow
//     and returns an operator that filter-maps effects
//
// # Example Usage
//
//	import A "github.com/IBM/fp-go/v2/array"
//	import O "github.com/IBM/fp-go/v2/option"
//
//	// Parse and filter positive integers
//	parsePositive := func(s string) O.Option[int] {
//	    var n int
//	    if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n > 0 {
//	        return O.Some(n)
//	    }
//	    return O.None[int]()
//	}
//
//	filterMapOp := FilterMap[MyContext](A.FilterMap[string, int])
//	pipeline := F.Pipe1(
//	    Succeed[MyContext]([]string{"1", "-2", "3", "invalid", "5"}),
//	    filterMapOp(parsePositive),
//	)
//	// Result: Effect that produces []int{1, 3, 5}
//
// # See Also
//
//   - FilterMapArray: Specialized version for arrays
//   - FilterMapIter: Specialized version for iterators
//   - Filter: For filtering without transformation
//
//go:inline
func FilterMap[C, HKTA, HKTB, A, B any](
	filter func(option.Kleisli[A, B]) Reader[HKTA, HKTB],
) func(option.Kleisli[A, B]) Operator[C, HKTA, HKTB] {
	return readerreaderioresult.FilterMap[C](filter)
}

// FilterMapArray creates an operator that filters and maps array elements within an Effect.
// Each element is transformed using a function that returns Option[B]. Elements that
// produce Some(b) are kept in the result array, while None values are filtered out.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input element type
//   - B: The output element type
//
// # Parameters
//
//   - p: A Kleisli arrow from A to Option[B] that transforms and filters elements
//
// # Returns
//
//   - Operator[C, []A, []B]: An operator that filter-maps array elements in an effect
//
// # Example Usage
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	// Double even numbers, filter out odd numbers
//	doubleEven := func(n int) O.Option[int] {
//	    if n%2 == 0 {
//	        return O.Some(n * 2)
//	    }
//	    return O.None[int]()
//	}
//
//	pipeline := F.Pipe1(
//	    Succeed[MyContext]([]int{1, 2, 3, 4, 5}),
//	    FilterMapArray[MyContext](doubleEven),
//	)
//	// Result: Effect that produces []int{4, 8}
//
// # See Also
//
//   - FilterMap: Generic version for any filterable type
//   - FilterMapIter: For filter-mapping iterators
//   - FilterArray: For filtering without transformation
//
//go:inline
func FilterMapArray[C, A, B any](p option.Kleisli[A, B]) Operator[C, []A, []B] {
	return readerreaderioresult.FilterMapArray[C](p)
}

// FilterMapIter creates an operator that filters and maps iterator elements within an Effect.
// Each element is transformed using a function that returns Option[B]. Elements that
// produce Some(b) are kept in the resulting iterator, while None values are filtered out.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The input element type
//   - B: The output element type
//
// # Parameters
//
//   - p: A Kleisli arrow from A to Option[B] that transforms and filters elements
//
// # Returns
//
//   - Operator[C, Seq[A], Seq[B]]: An operator that filter-maps iterator elements in an effect
//
// # Example Usage
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	// Parse strings to integers, keeping only valid ones
//	parseInt := func(s string) O.Option[int] {
//	    var n int
//	    if _, err := fmt.Sscanf(s, "%d", &n); err == nil {
//	        return O.Some(n)
//	    }
//	    return O.None[int]()
//	}
//
//	pipeline := F.Pipe1(
//	    Succeed[MyContext](slices.Values([]string{"1", "2", "invalid", "3"})),
//	    FilterMapIter[MyContext](parseInt),
//	)
//	// Result: Effect that produces an iterator over [1, 2, 3]
//
// # See Also
//
//   - FilterMap: Generic version for any filterable type
//   - FilterMapArray: For filter-mapping arrays
//   - FilterIter: For filtering without transformation
//
//go:inline
func FilterMapIter[C, A, B any](p option.Kleisli[A, B]) Operator[C, Seq[A], Seq[B]] {
	return readerreaderioresult.FilterMapIter[C](p)
}
