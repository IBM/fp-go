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

package option

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

func ExampleOption_extraction() {

	noneValue := None[int]()
	someValue := Of(42)

	// Convert Option[T] to T
	fromNone, okFromNone := Unwrap(noneValue) // 0, false
	fromSome, okFromSome := Unwrap(someValue) // 42, true

	// Convert Option[T] with a default value
	noneWithDefault := GetOrElse(F.Constant(0))(noneValue) // 0
	someWithDefault := GetOrElse(F.Constant(0))(someValue) // 42

	// Apply a different function on None/Some(...)
	doubleOrZero := Fold(
		F.Constant(0), // none case
		N.Mul(2),      // some case
	) // func(ma Option[int]) int

	doubleFromNone := doubleOrZero(noneValue) // 0
	doubleFromSome := doubleOrZero(someValue) // 84

	// Pro-tip: Fold is short for the following:
	doubleOfZeroBis := F.Flow2(
		Map(N.Mul(2)),            // some case
		GetOrElse(F.Constant(0)), // none case
	)
	doubleFromNoneBis := doubleOfZeroBis(noneValue) // 0
	doubleFromSomeBis := doubleOfZeroBis(someValue) // 84

	fmt.Printf("%d, %t\n", fromNone, okFromNone)
	fmt.Printf("%d, %t\n", fromSome, okFromSome)
	fmt.Println(noneWithDefault)
	fmt.Println(someWithDefault)
	fmt.Println(doubleFromNone)
	fmt.Println(doubleFromSome)
	fmt.Println(doubleFromNoneBis)
	fmt.Println(doubleFromSomeBis)

	// Output:
	// 0, false
	// 42, true
	// 0
	// 42
	// 0
	// 84
	// 0
	// 84
}
