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

func Example_extraction() {
	leftValue, leftErr := Left[int](fmt.Errorf("Division by Zero!"))
	rightValue, rightErr := Right(10)

	// Convert Either[A] to A with a default value
	leftWithDefault := GetOrElse(F.Constant1[error](0))(leftValue, leftErr)    // 0
	rightWithDefault := GetOrElse(F.Constant1[error](0))(rightValue, rightErr) // 10

	// Apply a different function on Left(...)/Right(...)
	doubleOrZero := Fold(F.Constant1[error](0), N.Mul(2)) // func(int, error) int
	doubleFromLeft := doubleOrZero(leftValue, leftErr)    // 0
	doubleFromRight := doubleOrZero(rightValue, rightErr) // 20

	// You can also chain operations using Map
	tripled, tripledErr := Pipe2(
		rightValue,
		Right[int],
		Map(N.Mul(3)),
	)
	tripledResult := GetOrElse(F.Constant1[error](0))(tripled, tripledErr) // 30

	fmt.Println(ToString(leftValue, leftErr))
	fmt.Println(ToString(rightValue, rightErr))
	fmt.Println(leftWithDefault)
	fmt.Println(rightWithDefault)
	fmt.Println(doubleFromLeft)
	fmt.Println(doubleFromRight)
	fmt.Println(tripledResult)

	// Output:
	// Left(Division by Zero!)
	// Right[int](10)
	// 0
	// 10
	// 0
	// 20
	// 30
}
