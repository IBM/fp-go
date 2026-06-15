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

package number

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// SemigroupSum creates a Semigroup for addition operations.
// It provides an associative binary operation that combines two values by adding them.
//
// The type parameter A can be any numeric type (integers, floats, complex numbers)
// or string type. For strings, concatenation is performed.
//
// Type Parameters:
//   - A: Any numeric type or string type
//
// Returns:
//   - S.Semigroup[A]: A semigroup instance with addition as the binary operation
//
// Example:
//
//	addSemigroup := SemigroupSum[int]()
//	result := addSemigroup.Concat(10, 20) // returns 30
//
//	floatSemigroup := SemigroupSum[float64]()
//	result := floatSemigroup.Concat(3.14, 2.86) // returns 6.0
//
//	stringSemigroup := SemigroupSum[string]()
//	result := stringSemigroup.Concat("Hello, ", "World!") // returns "Hello, World!"
//
// See Also:
//   - MonoidSum: Addition semigroup with identity element (0)
//   - SemigroupProduct: Multiplication semigroup
func SemigroupSum[A Number | ~string]() S.Semigroup[A] {
	return S.MakeSemigroup(Sum[A])
}

// SemigroupProduct creates a Semigroup for multiplication operations.
// It provides an associative binary operation that combines two values by multiplying them.
//
// The type parameter A can be any numeric type (integers, floats, complex numbers).
// Unlike SemigroupSum, this does not support string types.
//
// Type Parameters:
//   - A: Any numeric type
//
// Returns:
//   - S.Semigroup[A]: A semigroup instance with multiplication as the binary operation
//
// Example:
//
//	mulSemigroup := SemigroupProduct[int]()
//	result := mulSemigroup.Concat(5, 3) // returns 15
//
//	floatSemigroup := SemigroupProduct[float64]()
//	result := floatSemigroup.Concat(2.5, 4.0) // returns 10.0
//
//	complexSemigroup := SemigroupProduct[complex128]()
//	c1 := complex(2, 1)
//	c2 := complex(3, 4)
//	result := complexSemigroup.Concat(c1, c2) // returns (2+11i)
//
// See Also:
//   - MonoidProduct: Multiplication semigroup with identity element (1)
//   - SemigroupSum: Addition semigroup
func SemigroupProduct[A Number]() S.Semigroup[A] {
	return S.MakeSemigroup(Prod[A])
}
