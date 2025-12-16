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

package semigroup_test

import (
	"fmt"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/semigroup"
)

func ExampleSemigroup() {
	sum := N.SemigroupSum[int]()
	result := sum.Concat(sum.Concat(1, 2), 3)
	fmt.Println(result)
	// Output:
	// 6
}

func ExampleMakeSemigroup() {
	strConcat := semigroup.MakeSemigroup(func(a, b string) string {
		return a + b
	})
	result := strConcat.Concat("Hello, ", "World!")
	fmt.Println(result)
	// Output:
	// Hello, World!
}

func ExampleReverse() {
	sub := semigroup.MakeSemigroup(func(a, b int) int { return a - b })
	reversed := semigroup.Reverse(sub)
	result1 := sub.Concat(10, 3)
	result2 := reversed.Concat(10, 3)
	fmt.Println(result1, result2)
	// Output:
	// 7 -7
}

func ExampleFunctionSemigroup() {
	intSum := N.SemigroupSum[int]()
	funcSG := semigroup.FunctionSemigroup[string](intSum)

	f := func(s string) int { return len(s) }
	g := func(s string) int { return len(s) * 2 }
	combined := funcSG.Concat(f, g)
	result := combined("hello")
	fmt.Println(result)
	// Output:
	// 15
}

func ExampleFirst() {
	first := semigroup.First[int]()
	result := first.Concat(1, 2)
	fmt.Println(result)
	// Output:
	// 1
}

func ExampleLast() {
	last := semigroup.Last[int]()
	result := last.Concat(1, 2)
	fmt.Println(result)
	// Output:
	// 2
}

func ExampleToMagma() {
	sg := semigroup.First[int]()
	magma := semigroup.ToMagma(sg)
	result := magma.Concat(1, 2)
	fmt.Println(result)
	// Output:
	// 1
}
