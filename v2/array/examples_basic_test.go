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

package array

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// Example_basic adapts examples from [https://github.com/inato/fp-ts-cheatsheet#basic-manipulation]
func Example_basic() {

	someArray := From(0, 1, 2, 3, 4, 5, 6, 7, 8, 9) // []int

	isEven := func(num int) bool {
		return num%2 == 0
	}

	square := func(num int) int {
		return num * num
	}

	// filter and map
	result := F.Pipe2(
		someArray,
		Filter(isEven),
		Map(square),
	) // [0 4 16 36 64]

	// or in one go with filterMap
	resultFilterMap := F.Pipe1(
		someArray,
		FilterMap(
			F.Flow2(O.FromPredicate(isEven), O.Map(square)),
		),
	)

	fmt.Println(result)
	fmt.Println(resultFilterMap)

	// Output:
	// [0 4 16 36 64]
	// [0 4 16 36 64]
}
