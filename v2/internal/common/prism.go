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

package common

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

type (
	// prismTag holds the display name for a Prism[S, A].
	//
	// It is embedded in Prism[S, A] as a non-generic carrier so that the
	// formatting and logging methods (String, Format, LogValue) are compiled
	// once and shared across all type-parameter instantiations, rather than
	// being duplicated for every distinct Prism[S, A].
	prismTag struct {
		// n is the end-user-facing name of the prism (e.g. "MyType.Variant").
		// It is returned by String() and LogValue().
		n string
	}

	// Prism is an optic used to select part of a sum type (tagged union).
	// It provides two operations:
	//   - GetOption: Try to extract a value of type A from S (may fail)
	//   - ReverseGet: Construct an S from an A (always succeeds)
	//
	// Prisms are useful for working with variant types like Either, Option,
	// or custom sum types where you want to focus on a specific variant.
	//
	// Type Parameters:
	//   - S: The source type (sum type)
	//   - A: The focus type (variant within the sum type)
	//
	// Example:
	//   type Result interface{ isResult() }
	//   type Success struct{ Value int }
	//   type Failure struct{ Error string }
	//
	//   successPrism := MakePrism(
	//       func(r Result) Option[int] {
	//           if s, ok := r.(Success); ok {
	//               return Some(s.Value)
	//           }
	//           return None[int]()
	//       },
	//       func(v int) Result { return Success{Value: v} },
	//   )
	Prism[S, A any] struct {
		prismTag

		// GetOption attempts to extract a value of type A from S.
		// Returns Some(a) if the extraction succeeds, None otherwise.
		GetOption OptionKleisli[S, A]

		// ReverseGet constructs an S from an A.
		// This operation always succeeds.
		ReverseGet func(A) S
	}

	// Kleisli represents a function that takes a value of type A and returns a Prism[S, B].
	// This is commonly used for composing prisms in a monadic style.
	//
	// Type Parameters:
	//   - S: The source type of the resulting prism
	//   - A: The input type to the function
	//   - B: The focus type of the resulting prism
	PrismKleisli[S, A, B any] = func(A) Prism[S, B]

	// Operator represents a function that transforms one prism into another.
	// It takes a Prism[S, A] and returns a Prism[S, B], allowing for prism transformations.
	//
	// Type Parameters:
	//   - S: The source type (remains constant)
	//   - A: The original focus type
	//   - B: The new focus type
	PrismOperator[S, A, B any] = func(Prism[S, A]) Prism[S, B]
)

// MakePrism constructs a Prism from GetOption and ReverseGet functions.
//
// Parameters:
//   - get: Function to extract A from S (returns Option[A])
//   - rev: Function to construct S from A
//
// Returns:
//   - A Prism[S, A] that uses the provided functions
//
// Example:
//
//	prism := MakePrism(
//	    func(opt Option[int]) Option[int] { return opt },
//	    func(n int) Option[int] { return Some(n) },
//	)
//
//go:inline
func MakePrism[S, A any](get OptionKleisli[S, A], rev func(A) S) Prism[S, A] {
	return MakePrismWithName(get, rev, "GenericPrism")
}

//go:inline
func MakePrismWithName[S, A any](get OptionKleisli[S, A], rev func(A) S, name string) Prism[S, A] {
	return Prism[S, A]{GetOption: get, ReverseGet: rev, prismTag: prismTag{n: name}}
}

// PrismId returns an identity prism that focuses on the entire value.
// GetOption always returns Some(s), and ReverseGet is the identity function.
//
// This is useful as a starting point for prism composition or when you need
// a prism that doesn't actually transform the value.
//
// Example:
//
//	idPrism := PrismId[int]()
//	value := idPrism.GetOption(42)    // Some(42)
//	result := idPrism.ReverseGet(42)  // 42
func PrismId[S any]() Prism[S, S] {
	return MakePrismWithName(OptionSome[S], F.Identity[S], "PrismIdentity")
}

// PrismFromPredicate creates a prism that matches values satisfying a predicate.
// GetOption returns Some(s) if the predicate is true, None otherwise.
// ReverseGet is the identity function (doesn't validate the predicate).
//
// Parameters:
//   - pred: Predicate function to test values
//
// Returns:
//   - A Prism[S, S] that filters based on the predicate
//
// Example:
//
//	positivePrism := PrismFromPredicate(N.MoreThan(0))
//	value := positivePrism.GetOption(42)  // Some(42)
//	value = positivePrism.GetOption(-5)   // None[int]
func PrismFromPredicate[S any](pred func(S) bool) Prism[S, S] {
	return MakePrismWithName(OptionFromPredicate(pred), F.Identity[S], "PrismWithPredicate")
}

// Compose composes two prisms to create a prism that focuses deeper into a structure.
// The resulting prism first applies the outer prism (S → A), then the inner prism (A → B).
//
// Type Parameters:
//   - S: The outermost source type
//   - A: The intermediate type
//   - B: The innermost focus type
//
// Parameters:
//   - ab: The inner prism (A → B)
//
// Returns:
//   - A function that takes the outer prism (S → A) and returns the composed prism (S → B)
//
// Example:
//
//	outerPrism := MakePrism(...)  // Prism[Outer, Inner]
//	innerPrism := MakePrism(...)  // Prism[Inner, Value]
//	composed := Compose[Outer](innerPrism)(outerPrism)  // Prism[Outer, Value]
func PrismComposePrism[S, A, B any](ab Prism[A, B]) PrismOperator[S, A, B] {
	return func(sa Prism[S, A]) Prism[S, B] {
		return MakePrismWithName(F.Flow2(
			sa.GetOption,
			OptionChain(ab.GetOption),
		), F.Flow2(
			ab.ReverseGet,
			sa.ReverseGet,
		),
			fmt.Sprintf("PrismCompose[%s x %s]", ab, sa),
		)
	}
}

// prismModifyOption applies a transformation function through a prism,
// returning Some(modified S) if the prism matches, None otherwise.
// This is an internal helper function.
func prismModifyOption[S, A any](f Endomorphism[A], sa Prism[S, A], s S) Option[S] {
	return F.Pipe2(
		s,
		sa.GetOption,
		OptionMap(F.Flow2(
			f,
			sa.ReverseGet,
		)),
	)
}

// prismModify applies a transformation function through a prism.
// If the prism matches, it extracts the value, applies the function,
// and reconstructs the result. If the prism doesn't match, returns the original value.
// This is an internal helper function.
func prismModify[S, A any](f Endomorphism[A], sa Prism[S, A], s S) S {
	return F.Pipe1(
		prismModifyOption(f, sa, s),
		OptionGetOrElse(F.Constant(s)),
	)
}

// prismSet is an internal helper that creates a setter function.
//
// Deprecated: Use Set instead.
func prismSet[S, A any](a A) func(Prism[S, A]) Endomorphism[S] {
	return F.Curry3(prismModify[S, A])(F.Constant1[A](a))
}

// PrismSet creates a function that sets a value through a prism.
// If the prism matches, it replaces the focused value with the new value.
// If the prism doesn't match, it returns the original value unchanged.
//
// Parameters:
//   - a: The new value to set
//
// Returns:
//   - A function that takes a prism and returns an endomorphism (S → S)
//
// Example:
//
//	somePrism := MakePrism(...)
//	setter := PrismSet[Option[int], int](100)
//	result := setter(somePrism)(Some(42))  // Some(100)
//	result = setter(somePrism)(None[int]()) // None[int]() (unchanged)
func PrismSet[S, A any](a A) func(Prism[S, A]) Endomorphism[S] {
	return F.Curry3(prismModify[S, A])(F.Constant1[A](a))
}

// PrismSome creates a prism that focuses on the PrismSome variant of an Option within a structure.
// It composes the provided prism (which focuses on an Option[A]) with a prism that
// extracts the value from PrismSome.
//
// Type Parameters:
//   - S: The source type
//   - A: The value type within the Option
//
// Parameters:
//   - soa: A prism that focuses on an Option[A] within S
//
// Returns:
//   - A prism that focuses on the A value within PrismSome
//
// Example:
//
//	type Config struct { Timeout Option[int] }
//	configPrism := MakePrism(...)  // Prism[Config, Option[int]]
//	timeoutPrism := PrismSome(configPrism)  // Prism[Config, int]
//	value := timeoutPrism.GetOption(Config{Timeout: PrismSome(30)})  // PrismSome(30)
func PrismSome[S, A any](soa Prism[S, Option[A]]) Prism[S, A] {
	return PrismComposePrism[S](PrismFromOption[A]())(soa)
}

// prismImap is an internal helper that bidirectionally maps a prism's focus type.
func prismImap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](sa Prism[S, A], ab AB, ba BA) Prism[S, B] {
	return MakePrismWithName(
		F.Flow2(sa.GetOption, OptionMap(ab)),
		F.Flow2(ba, sa.ReverseGet),
		fmt.Sprintf("PrismIMap[%s]", sa),
	)
}

// PrismIMap bidirectionally maps the focus type of a prism.
// It transforms a Prism[S, A] into a Prism[S, B] using two functions:
// one to map A → B and another to map B → A.
//
// Type Parameters:
//   - S: The source type
//   - A: The original focus type
//   - B: The new focus type
//   - AB: Function type A → B
//   - BA: Function type B → A
//
// Parameters:
//   - ab: Function to map from A to B
//   - ba: Function to map from B to A
//
// Returns:
//   - A function that transforms Prism[S, A] to Prism[S, B]
//
// Example:
//
//	intPrism := MakePrism(...)  // Prism[Result, int]
//	stringPrism := PrismIMap[Result](
//	    strconv.Itoa,
//	    func(s string) int { n, _ := strconv.Atoi(s); return n },
//	)(intPrism)  // Prism[Result, string]
func PrismIMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) PrismOperator[S, A, B] {
	return func(sa Prism[S, A]) Prism[S, B] {
		return prismImap(sa, ab, ba)
	}
}

// PrismFromOption creates a prism for extracting values from Option types.
// It provides a safe way to work with Option values, focusing on the Some case
// and handling the None case gracefully through the prism's GetOption behavior.
//
// The prism's GetOption is the identity function - it returns the Option as-is.
// If the Option is Some(value), GetOption returns Some(value); if it's None, it returns None.
// This allows the prism to naturally handle the presence or absence of a value.
//
// The prism's ReverseGet wraps a value into Some, always succeeding.
//
// Type Parameters:
//   - T: The value type contained in the Option
//
// Returns:
//   - A Prism[Option[T], T] that safely extracts values from Options
//
// Example:
//
//	// Create a prism for extracting int values from Option[int]
//	optPrism := FromOption[int]()
//
//	// Extract from Some
//	someValue := option.Some(42)
//	result := optPrism.GetOption(someValue)  // Some(42)
//
//	// Extract from None
//	noneValue := option.None[int]()
//	result = optPrism.GetOption(noneValue)  // None[int]()
//
//	// Wrap value into Some
//	wrapped := optPrism.ReverseGet(100)  // Some(100)
//
//	// Use with Set to update Some values
//	setter := Set[Option[int], int](200)
//	result := setter(optPrism)(someValue)  // Some(200)
//	result = setter(optPrism)(noneValue)   // None[int]() (unchanged)
//
//	// Compose with other prisms for nested extraction
//	// Extract int from Option[Option[int]]
//	nestedPrism := Compose[Option[Option[int]], Option[int], int](
//	    FromOption[Option[int]](),
//	    FromOption[int](),
//	)
//	nested := option.Some(option.Some(42))
//	value := nestedPrism.GetOption(nested)  // Some(42)
//
// Common use cases:
//   - Extracting values from optional fields
//   - Working with nullable data in a type-safe way
//   - Composing with other prisms to handle nested Options
//   - Filtering and transforming optional values in pipelines
//   - Converting between Option and other optional representations
//
// Key insight: This prism treats Option[T] as a "container" that may or may not
// hold a value of type T. The prism focuses on the value inside, allowing you to
// work with it when present and gracefully handle its absence when not.
//
//go:inline
func PrismFromOption[T any]() Prism[Option[T], T] {
	return MakePrismWithName(
		F.Identity[Option[T]],
		OptionSome[T],
		"PrismFromOption",
	)
}
