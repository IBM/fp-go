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

package iter

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/option"
)

// FindFirstMap finds the first element in a sequence that successfully maps to a value.
// It applies a selector function that returns an Option to each element, and returns
// the first Some result wrapped in an IOOption. If no element produces Some, returns None.
//
// This function combines FilterMap and First, providing an efficient way to search for
// and transform an element in a single pass. The sequence is consumed lazily, stopping
// as soon as a matching element is found.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of the transformed result
//
// Parameters:
//   - sel: A selector function that returns Some(B) for matching elements, None otherwise
//
// Returns:
//   - iooption.Kleisli[Seq[A], B]: A function that takes a sequence and returns an
//     IOOption that produces the first successfully mapped value when executed
//
// Example:
//
//	parsePositive := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil || n <= 0 {
//	        return None[int]()
//	    }
//	    return Some(n)
//	}
//	seq := From("invalid", "0", "42", "100")
//	result := FindFirstMap(parsePositive)(seq)()
//	// Returns: Some(42)
//
// Example with no match:
//
//	seq := From("invalid", "0", "-5")
//	result := FindFirstMap(parsePositive)(seq)()
//	// Returns: None
//
// See Also:
//   - FindFirst: Find first element matching a predicate
//   - FindLastMap: Find last element that maps successfully
//   - FilterMap: Filter and map all elements
//   - First: Get first element of a sequence
func FindFirstMap[A, B any](sel option.Kleisli[A, B]) iooption.Kleisli[Seq[A], B] {
	return F.Flow2(
		FilterMap(sel),
		First,
	)
}

// FindFirst finds the first element in a sequence that satisfies a predicate.
// It returns the first matching element wrapped in an IOOption. If no element
// matches, returns None.
//
// This function combines Filter and First, providing an efficient way to search
// for an element. The sequence is consumed lazily, stopping as soon as a matching
// element is found.
//
// Type Parameters:
//   - A: The type of elements in the sequence
//
// Parameters:
//   - p: A predicate function that returns true for matching elements
//
// Returns:
//   - iooption.Kleisli[Seq[A], A]: A function that takes a sequence and returns an
//     IOOption that produces the first matching element when executed
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := FindFirst(func(x int) bool { return x > 3 })(seq)()
//	// Returns: Some(4)
//
// Example with no match:
//
//	seq := From(1, 2, 3)
//	result := FindFirst(func(x int) bool { return x > 10 })(seq)()
//	// Returns: None
//
// Example with strings:
//
//	seq := From("apple", "banana", "cherry")
//	result := FindFirst(func(s string) bool { return len(s) > 6 })(seq)()
//	// Returns: Some("banana")
//
// See Also:
//   - FindFirstMap: Find and transform first matching element
//   - FindLast: Find last element matching a predicate
//   - Filter: Filter all elements matching a predicate
//   - First: Get first element of a sequence
func FindFirst[A any](p Predicate[A]) iooption.Kleisli[Seq[A], A] {
	return F.Flow2(
		Filter(p),
		First,
	)
}

// FindLastMap finds the last element in a sequence that successfully maps to a value.
// It applies a selector function that returns an Option to each element, and returns
// the last Some result wrapped in an IOOption. If no element produces Some, returns None.
//
// This function combines FilterMap and Last. Unlike FindFirstMap, it must consume
// the entire sequence to determine the last matching element.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of the transformed result
//
// Parameters:
//   - sel: A selector function that returns Some(B) for matching elements, None otherwise
//
// Returns:
//   - iooption.Kleisli[Seq[A], B]: A function that takes a sequence and returns an
//     IOOption that produces the last successfully mapped value when executed
//
// Example:
//
//	parsePositive := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil || n <= 0 {
//	        return None[int]()
//	    }
//	    return Some(n)
//	}
//	seq := From("invalid", "42", "100", "0")
//	result := FindLastMap(parsePositive)(seq)()
//	// Returns: Some(100)
//
// Example with no match:
//
//	seq := From("invalid", "0", "-5")
//	result := FindLastMap(parsePositive)(seq)()
//	// Returns: None
//
// See Also:
//   - FindLast: Find last element matching a predicate
//   - FindFirstMap: Find first element that maps successfully
//   - FilterMap: Filter and map all elements
//   - Last: Get last element of a sequence
func FindLastMap[A, B any](sel option.Kleisli[A, B]) iooption.Kleisli[Seq[A], B] {
	return F.Flow2(
		FilterMap(sel),
		Last,
	)
}

// FindLast finds the last element in a sequence that satisfies a predicate.
// It returns the last matching element wrapped in an IOOption. If no element
// matches, returns None.
//
// This function combines Filter and Last. Unlike FindFirst, it must consume
// the entire sequence to determine the last matching element.
//
// Type Parameters:
//   - A: The type of elements in the sequence
//
// Parameters:
//   - p: A predicate function that returns true for matching elements
//
// Returns:
//   - iooption.Kleisli[Seq[A], A]: A function that takes a sequence and returns an
//     IOOption that produces the last matching element when executed
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := FindLast(func(x int) bool { return x > 3 })(seq)()
//	// Returns: Some(5)
//
// Example with no match:
//
//	seq := From(1, 2, 3)
//	result := FindLast(func(x int) bool { return x > 10 })(seq)()
//	// Returns: None
//
// Example with strings:
//
//	seq := From("apple", "banana", "cherry", "date")
//	result := FindLast(func(s string) bool { return len(s) > 5 })(seq)()
//	// Returns: Some("cherry")
//
// See Also:
//   - FindLastMap: Find and transform last matching element
//   - FindFirst: Find first element matching a predicate
//   - Filter: Filter all elements matching a predicate
//   - Last: Get last element of a sequence
func FindLast[A any](p Predicate[A]) iooption.Kleisli[Seq[A], A] {
	return F.Flow2(
		Filter(p),
		Last,
	)
}
