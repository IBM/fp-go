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
)

func Example_find() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	fmt.Println(FindFirst(pred)(data1))

	// Output:
	// Some[int](1)
}

func Example_find_filter() {

	pred := func(val int) bool {
		return val&2 == 0
	}

	data1 := From(1, 2, 3)

	Find := F.Flow2(
		Filter(pred),
		Head[int],
	)

	fmt.Println(Find(data1))

	// Output:
	// Some[int](1)
}
