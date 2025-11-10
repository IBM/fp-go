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

package result

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

func ExampleEither_extraction() {
	leftValue := Left[int](fmt.Errorf("Division by Zero!"))
	rightValue := Right(10)

	// Convert Result[A] to A with a default value
	leftWithDefault := GetOrElse(F.Constant1[error](0))(leftValue)   // 0
	rightWithDefault := GetOrElse(F.Constant1[error](0))(rightValue) // 10

	// Apply a different function on Left(...)/Right(...)
	doubleOrZero := Fold(F.Constant1[error](0), N.Mul(2)) // func(Result[int]) int
	doubleFromLeft := doubleOrZero(leftValue)             // 0
	doubleFromRight := doubleOrZero(rightValue)           // 20

	// Pro-tip: Fold is short for the following:
	doubleOrZeroBis := F.Flow2(
		Map(N.Mul(2)),
		GetOrElse(F.Constant1[error](0)),
	)
	doubleFromLeftBis := doubleOrZeroBis(leftValue)   // 0
	doubleFromRightBis := doubleOrZeroBis(rightValue) // 20

	fmt.Println(leftValue)
	fmt.Println(rightValue)
	fmt.Println(leftWithDefault)
	fmt.Println(rightWithDefault)
	fmt.Println(doubleFromLeft)
	fmt.Println(doubleFromRight)
	fmt.Println(doubleFromLeftBis)
	fmt.Println(doubleFromRightBis)

	// Output:
	// Left[*errors.errorString](Division by Zero!)
	// Right[int](10)
	// 0
	// 10
	// 0
	// 20
	// 0
	// 20
}
