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

package predicate

import (
	F "github.com/IBM/fp-go/v2/function"
)

// ContraMap creates a new predicate by transforming the input before applying an existing predicate.
//
// This is a contravariant functor operation that allows you to adapt a predicate for type A
// to work with type B by providing a function that converts B to A. The resulting predicate
// first applies the mapping function f to transform the input, then applies the original predicate.
//
// This is particularly useful when you have a predicate for one type and want to reuse it
// for a related type without rewriting the predicate logic.
//
// Parameters:
//   - f: A function that converts values of type B to type A
//
// Returns:
//   - An Operator that transforms a Predicate[A] into a Predicate[B]
//
// Example:
//
//	type Person struct { Age int }
//	isAdult := func(age int) bool { return age >= 18 }
//	getAge := func(p Person) int { return p.Age }
//	isPersonAdult := F.Pipe1(isAdult, ContraMap(getAge))
//	isPersonAdult(Person{Age: 25}) // true
//	isPersonAdult(Person{Age: 15}) // false
func ContraMap[A, B any](f func(B) A) Operator[A, B] {
	return func(pred Predicate[A]) Predicate[B] {
		return F.Flow2(
			f,
			pred,
		)
	}
}
