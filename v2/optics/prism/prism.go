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

package prism

import (
	"github.com/IBM/fp-go/v2/internal/common"
	O "github.com/IBM/fp-go/v2/option"
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
func MakePrism[S, A any](get O.Kleisli[S, A], rev func(A) S) Prism[S, A] {
	return common.MakePrism(get, rev)
}

//go:inline
func MakePrismWithName[S, A any](get O.Kleisli[S, A], rev func(A) S, name string) Prism[S, A] {
	return common.MakePrismWithName(get, rev, name)
}

// Id returns an identity prism that focuses on the entire value.
// GetOption always returns Some(s), and ReverseGet is the identity function.
//
// This is useful as a starting point for prism composition or when you need
// a prism that doesn't actually transform the value.
//
// Example:
//
//	idPrism := Id[int]()
//	value := idPrism.GetOption(42)    // Some(42)
//	result := idPrism.ReverseGet(42)  // 42
func Id[S any]() Prism[S, S] {
	return common.PrismId[S]()
}

// FromPredicate creates a prism that matches values satisfying a predicate.
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
//	positivePrism := FromPredicate(N.MoreThan(0))
//	value := positivePrism.GetOption(42)  // Some(42)
//	value = positivePrism.GetOption(-5)   // None[int]
func FromPredicate[S any](pred func(S) bool) Prism[S, S] {
	return common.PrismFromPredicate(pred)
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
func Compose[S, A, B any](ab Prism[A, B]) Operator[S, A, B] {
	return common.PrismComposePrism[S](ab)
}

// Set creates a function that sets a value through a prism.
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
//	setter := Set[Option[int], int](100)
//	result := setter(somePrism)(Some(42))  // Some(100)
//	result = setter(somePrism)(None[int]()) // None[int]() (unchanged)
func Set[S, A any](a A) func(Prism[S, A]) Endomorphism[S] {
	return common.PrismSet[S](a)
}

// Some creates a prism that focuses on the Some variant of an Option within a structure.
// It composes the provided prism (which focuses on an Option[A]) with a prism that
// extracts the value from Some.
//
// Type Parameters:
//   - S: The source type
//   - A: The value type within the Option
//
// Parameters:
//   - soa: A prism that focuses on an Option[A] within S
//
// Returns:
//   - A prism that focuses on the A value within Some
//
// Example:
//
//	type Config struct { Timeout Option[int] }
//	configPrism := MakePrism(...)  // Prism[Config, Option[int]]
//	timeoutPrism := Some(configPrism)  // Prism[Config, int]
//	value := timeoutPrism.GetOption(Config{Timeout: Some(30)})  // Some(30)
func Some[S, A any](soa Prism[S, Option[A]]) Prism[S, A] {
	return common.PrismSome(soa)
}

// IMap bidirectionally maps the focus type of a prism.
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
//	stringPrism := IMap[Result](
//	    strconv.Itoa,
//	    func(s string) int { n, _ := strconv.Atoi(s); return n },
//	)(intPrism)  // Prism[Result, string]
func IMap[S any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) Operator[S, A, B] {
	return common.PrismIMap[S](ab, ba)
}
