// Copyright (c) 2025 IBM Corp.
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

// Package either provides implementations of the Either type and related operations.
//
// This package implements several Fantasy Land algebraic structures:
//   - Filterable: https://github.com/fantasyland/fantasy-land#filterable
//
// The Filterable specification defines operations for filtering and partitioning
// data structures based on predicates and mapping functions.
package either

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// Partition separates an [Either] value into a [Pair] based on a predicate function.
// It returns a function that takes an Either and produces a Pair of Either values,
// where the first element contains values that fail the predicate and the second
// contains values that pass the predicate.
//
// This function implements the Filterable specification's partition operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is Left, both elements of the resulting Pair will be the same Left value
//   - If the input is Right and the predicate returns true, the result is (Left(empty), Right(value))
//   - If the input is Right and the predicate returns false, the result is (Right(value), Left(empty))
//
// This function is useful for separating Either values into two categories based on
// a condition, commonly used in filtering operations where you want to keep track of
// both the values that pass and fail a test.
//
// Parameters:
//   - p: A predicate function that tests values of type A
//   - empty: The default Left value to use when creating Left instances for partitioning
//
// Returns:
//
//	A function that takes an Either[E, A] and returns a Pair where:
//	  - First element: Either values that fail the predicate (or original Left)
//	  - Second element: Either values that pass the predicate (or original Left)
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    N "github.com/IBM/fp-go/v2/number"
//	    P "github.com/IBM/fp-go/v2/pair"
//	)
//
//	// Partition positive and non-positive numbers
//	isPositive := N.MoreThan(0)
//	partition := E.Partition(isPositive, "not positive")
//
//	// Right value that passes predicate
//	result1 := partition(E.Right[string](5))
//	// result1 = Pair(Left("not positive"), Right(5))
//	left1, right1 := P.Unpack(result1)
//	// left1 = Left("not positive"), right1 = Right(5)
//
//	// Right value that fails predicate
//	result2 := partition(E.Right[string](-3))
//	// result2 = Pair(Right(-3), Left("not positive"))
//	left2, right2 := P.Unpack(result2)
//	// left2 = Right(-3), right2 = Left("not positive")
//
//	// Left value passes through unchanged in both positions
//	result3 := partition(E.Left[int]("error"))
//	// result3 = Pair(Left("error"), Left("error"))
//	left3, right3 := P.Unpack(result3)
//	// left3 = Left("error"), right3 = Left("error")
func Partition[E, A any](p Predicate[A], empty E) func(Either[E, A]) Pair[Either[E, A], Either[E, A]] {
	l := Left[A](empty)
	return func(e Either[E, A]) Pair[Either[E, A], Either[E, A]] {
		if e.isLeft {
			return pair.Of(e)
		}
		if p(e.r) {
			return pair.MakePair(l, e)
		}
		return pair.MakePair(e, l)
	}
}

// Filter creates a filtering operation for [Either] values based on a predicate function.
// It returns a function that takes an Either and produces an Either, where Right values
// that fail the predicate are converted to Left values with the provided empty value.
//
// This function implements the Filterable specification's filter operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is Left, it passes through unchanged
//   - If the input is Right and the predicate returns true, the Right value passes through unchanged
//   - If the input is Right and the predicate returns false, it's converted to Left(empty)
//
// This function is useful for conditional validation or filtering of Either values,
// where you want to reject Right values that don't meet certain criteria by converting
// them to Left values with a default error.
//
// Parameters:
//   - p: A predicate function that tests values of type A
//   - empty: The default Left value to use when filtering out Right values that fail the predicate
//
// Returns:
//
//	An Operator function that takes an Either[E, A] and returns an Either[E, A] where:
//	  - Left values pass through unchanged
//	  - Right values that pass the predicate remain as Right
//	  - Right values that fail the predicate become Left(empty)
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	// Filter to keep only positive numbers
//	isPositive := N.MoreThan(0)
//	filterPositive := E.Filter(isPositive, "not positive")
//
//	// Right value that passes predicate - remains Right
//	result1 := filterPositive(E.Right[string](5))
//	// result1 = Right(5)
//
//	// Right value that fails predicate - becomes Left
//	result2 := filterPositive(E.Right[string](-3))
//	// result2 = Left("not positive")
//
//	// Left value passes through unchanged
//	result3 := filterPositive(E.Left[int]("original error"))
//	// result3 = Left("original error")
//
//	// Chaining filters
//	isEven := func(n int) bool { return n%2 == 0 }
//	filterEven := E.Filter(isEven, "not even")
//
//	// Apply multiple filters in sequence
//	result4 := filterEven(filterPositive(E.Right[string](4)))
//	// result4 = Right(4) - passes both filters
//
//	result5 := filterEven(filterPositive(E.Right[string](3)))
//	// result5 = Left("not even") - passes first, fails second
func Filter[E, A any](p Predicate[A], empty E) Operator[E, A, A] {
	l := Left[A](empty)
	return func(e Either[E, A]) Either[E, A] {
		if e.isLeft || p(e.r) {
			return e
		}
		return l
	}
}

// FilterMap combines filtering and mapping operations for [Either] values using an [Option]-returning function.
// It returns a function that takes an Either[E, A] and produces an Either[E, B], where Right values
// are transformed by applying the function f. If f returns Some(B), the result is Right(B). If f returns
// None, the result is Left(empty).
//
// This function implements the Filterable specification's filterMap operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is Left, it passes through with its error value preserved as Left[B]
//   - If the input is Right and f returns Some(B), the result is Right(B)
//   - If the input is Right and f returns None, the result is Left(empty)
//
// This function is useful for operations that combine validation/filtering with transformation,
// such as parsing strings to numbers (where invalid strings result in None), or extracting
// optional fields from structures.
//
// Parameters:
//   - f: An Option Kleisli function that transforms values of type A to Option[B]
//   - empty: The default Left value to use when f returns None
//
// Returns:
//
//	An Operator function that takes an Either[E, A] and returns an Either[E, B] where:
//	  - Left values pass through with error preserved
//	  - Right values are transformed by f: Some(B) becomes Right(B), None becomes Left(empty)
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    O "github.com/IBM/fp-go/v2/option"
//	    "strconv"
//	)
//
//	// Parse string to int, filtering out invalid values
//	parseInt := func(s string) O.Option[int] {
//	    if n, err := strconv.Atoi(s); err == nil {
//	        return O.Some(n)
//	    }
//	    return O.None[int]()
//	}
//	filterMapInt := E.FilterMap(parseInt, "invalid number")
//
//	// Valid number string - transforms to Right(int)
//	result1 := filterMapInt(E.Right[string]("42"))
//	// result1 = Right(42)
//
//	// Invalid number string - becomes Left
//	result2 := filterMapInt(E.Right[string]("abc"))
//	// result2 = Left("invalid number")
//
//	// Left value passes through with error preserved
//	result3 := filterMapInt(E.Left[string]("original error"))
//	// result3 = Left("original error")
//
//	// Extract optional field from struct
//	type Person struct {
//	    Name  string
//	    Email O.Option[string]
//	}
//	extractEmail := func(p Person) O.Option[string] { return p.Email }
//	filterMapEmail := E.FilterMap(extractEmail, "no email")
//
//	result4 := filterMapEmail(E.Right[string](Person{Name: "Alice", Email: O.Some("alice@example.com")}))
//	// result4 = Right("alice@example.com")
//
//	result5 := filterMapEmail(E.Right[string](Person{Name: "Bob", Email: O.None[string]()}))
//	// result5 = Left("no email")
func FilterMap[E, A, B any](f option.Kleisli[A, B], empty E) Operator[E, A, B] {
	l := Left[B](empty)
	return func(e Either[E, A]) Either[E, B] {
		if e.isLeft {
			return Left[B](e.l)
		}
		if b, ok := option.Unwrap(f(e.r)); ok {
			return Right[E](b)
		}
		return l
	}
}

// PartitionMap separates and transforms an [Either] value into a [Pair] of Either values using a mapping function.
// It returns a function that takes an Either[E, A] and produces a Pair of Either values, where the mapping
// function f transforms the Right value into Either[B, C]. The result is partitioned based on whether f
// produces a Left or Right value.
//
// This function implements the Filterable specification's partitionMap operation:
// https://github.com/fantasyland/fantasy-land#filterable
//
// The behavior is as follows:
//   - If the input is Left, both elements of the resulting Pair will be Left with the original error
//   - If the input is Right and f returns Left(B), the result is (Right(B), Left(empty))
//   - If the input is Right and f returns Right(C), the result is (Left(empty), Right(C))
//
// This function is useful for operations that need to categorize and transform values simultaneously,
// such as separating valid and invalid data while applying different transformations to each category.
//
// Parameters:
//   - f: A Kleisli function that transforms values of type A to Either[B, C]
//   - empty: The default error value to use when creating Left instances for partitioning
//
// Returns:
//
//	A function that takes an Either[E, A] and returns a Pair[Either[E, B], Either[E, C]] where:
//	  - If input is Left: (Left(original_error), Left(original_error))
//	  - If f returns Left(B): (Right(B), Left(empty))
//	  - If f returns Right(C): (Left(empty), Right(C))
//
// Example:
//
//	import (
//	    E "github.com/IBM/fp-go/v2/either"
//	    P "github.com/IBM/fp-go/v2/pair"
//	)
//
//	// Classify and transform numbers: negative -> error message, positive -> squared value
//	classifyNumber := func(n int) E.Either[string, int] {
//	    if n < 0 {
//	        return E.Left[int]("negative: " + strconv.Itoa(n))
//	    }
//	    return E.Right[string](n * n)
//	}
//	partitionMap := E.PartitionMap(classifyNumber, "not classified")
//
//	// Positive number - goes to right side as squared value
//	result1 := partitionMap(E.Right[string](5))
//	// result1 = Pair(Left("not classified"), Right(25))
//	left1, right1 := P.Unpack(result1)
//	// left1 = Left("not classified"), right1 = Right(25)
//
//	// Negative number - goes to left side with error message
//	result2 := partitionMap(E.Right[string](-3))
//	// result2 = Pair(Right("negative: -3"), Left("not classified"))
//	left2, right2 := P.Unpack(result2)
//	// left2 = Right("negative: -3"), right2 = Left("not classified")
//
//	// Original Left value - appears in both positions
//	result3 := partitionMap(E.Left[int]("original error"))
//	// result3 = Pair(Left("original error"), Left("original error"))
//	left3, right3 := P.Unpack(result3)
//	// left3 = Left("original error"), right3 = Left("original error")
//
//	// Validate and transform user input
//	type ValidationError struct{ Field, Message string }
//	type User struct{ Name string; Age int }
//
//	validateUser := func(input map[string]string) E.Either[ValidationError, User] {
//	    name, hasName := input["name"]
//	    ageStr, hasAge := input["age"]
//	    if !hasName {
//	        return E.Left[User](ValidationError{"name", "missing"})
//	    }
//	    if !hasAge {
//	        return E.Left[User](ValidationError{"age", "missing"})
//	    }
//	    age, err := strconv.Atoi(ageStr)
//	    if err != nil {
//	        return E.Left[User](ValidationError{"age", "invalid"})
//	    }
//	    return E.Right[ValidationError](User{name, age})
//	}
//	partitionUsers := E.PartitionMap(validateUser, ValidationError{"", "not processed"})
//
//	validInput := map[string]string{"name": "Alice", "age": "30"}
//	result4 := partitionUsers(E.Right[string](validInput))
//	// result4 = Pair(Left(ValidationError{"", "not processed"}), Right(User{"Alice", 30}))
//
//	invalidInput := map[string]string{"name": "Bob"}
//	result5 := partitionUsers(E.Right[string](invalidInput))
//	// result5 = Pair(Right(ValidationError{"age", "missing"}), Left(ValidationError{"", "not processed"}))
func PartitionMap[E, A, B, C any](f Kleisli[B, A, C], empty E) func(Either[E, A]) Pair[Either[E, B], Either[E, C]] {
	return func(e Either[E, A]) Pair[Either[E, B], Either[E, C]] {
		if e.isLeft {
			return pair.MakePair(Left[B](e.l), Left[C](e.l))
		}
		res := f(e.r)
		if res.isLeft {
			return pair.MakePair(Right[E](res.l), Left[C](empty))
		}
		return pair.MakePair(Left[B](empty), Right[E](res.r))
	}
}
