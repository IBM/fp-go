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

func Example_any() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	fmt.Println(Any(pred)(data1))

	// Output:
	// true
}

func Example_any_filter() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	// Any tests if any of the entries in the array matches the condition
	Any := F.Flow2(
		Filter(pred),
		IsNonEmpty[int],
	)

	fmt.Println(Any(data1))

	// Output:
	// true
}

func Example_any_find() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	// Any tests if any of the entries in the array matches the condition
	Any := F.Flow2(
		FindFirst(pred),
		O.IsSome[int],
	)

	fmt.Println(Any(data1))

	// Output:
	// true
}
