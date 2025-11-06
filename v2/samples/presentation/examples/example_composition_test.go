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

package examples

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
)

func Example_composition_pipe() {

	filter := func(i int) bool {
		return i%2 == 0
	}

	double := func(i int) int {
		return i * 2
	}

	input := []int{1, 2, 3, 4}

	res := F.Pipe2(
		input,
		A.Filter(filter),
		A.Map(double),
	)

	fmt.Println(res)

	// Output:
	// [4 8]
}

func Example_composition_flow() {

	filter := func(i int) bool {
		return i%2 == 0
	}

	double := func(i int) int {
		return i * 2
	}

	input := []int{1, 2, 3, 4}

	filterAndDouble := F.Flow2(
		A.Filter(filter),
		A.Map(double),
	) // func([]int) []int

	fmt.Println(filterAndDouble(input))

	// Output:
	// [4 8]
}
