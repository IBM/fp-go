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
	N "github.com/IBM/fp-go/v2/number"
)

func Example_map() {

	f := func(i int) int {
		return i * 2
	}

	input := []int{1, 2, 3, 4}

	// idiomatic go
	res1 := make([]int, 0, len(input))
	for _, i := range input {
		res1 = append(res1, f(i))
	}
	fmt.Println(res1)

	// map
	res2 := A.Map(f)(input)
	fmt.Println(res2)

	// Output:
	// [2 4 6 8]
	// [2 4 6 8]
}

func Example_reduce() {

	input := []int{1, 2, 3, 4}

	// reduce
	red := A.Reduce(N.MonoidSum[int]().Concat, 0)(input)
	fmt.Println(red)

	// fold
	fld := A.Fold(N.MonoidSum[int]())(input)
	fmt.Println(fld)

	// Output:
	// 10
	// 10
}
