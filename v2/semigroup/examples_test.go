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

	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
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

func ExampleConcatWith() {
	// ConcatWith fixes the left operand and returns a function waiting for the right.
	sum := N.SemigroupSum[int]()
	add5 := semigroup.ConcatWith(sum)(5)
	fmt.Println(add5(3))  // 5 + 3
	fmt.Println(add5(10)) // 5 + 10
	// Output:
	// 8
	// 15
}

func ExampleConcatWith_prefix() {
	// ConcatWith with string concatenation: fixes a prefix.
	sg := semigroup.MakeSemigroup(func(a, b string) string { return a + b })
	greet := semigroup.ConcatWith(sg)("Hello, ")
	fmt.Println(greet("World"))
	fmt.Println(greet("Go"))
	// Output:
	// Hello, World
	// Hello, Go
}

func ExampleAppendTo() {
	// AppendTo fixes the right operand and returns a function waiting for the left.
	sg := semigroup.MakeSemigroup(func(a, b string) string { return a + b })
	addBang := semigroup.AppendTo(sg)("!")
	fmt.Println(addBang("Hello"))
	fmt.Println(addBang("World"))
	// Output:
	// Hello!
	// World!
}

func ExampleAppendTo_nonCommutative() {
	// AppendTo vs ConcatWith on a non-commutative operation.
	sub := semigroup.MakeSemigroup(func(a, b int) int { return a - b })
	// AppendTo fixes the right: x - 3
	subtract3 := semigroup.AppendTo(sub)(3)
	fmt.Println(subtract3(10)) // 10 - 3 = 7
	// ConcatWith fixes the left: 10 - x
	subtract10From := semigroup.ConcatWith(sub)(10)
	fmt.Println(subtract10From(3)) // 10 - 3 = 7
	// Output:
	// 7
	// 7
}

func ExampleConcatAll() {
	// ConcatAll reduces a slice to a single value using an initial seed.
	sum := N.SemigroupSum[int]()
	result := semigroup.ConcatAll(sum)(0)([]int{1, 2, 3, 4})
	fmt.Println(result)
	// Output:
	// 10
}

func ExampleMonadConcatAll() {
	// MonadConcatAll is the uncurried form: takes the slice and seed together.
	sum := N.SemigroupSum[int]()
	result := semigroup.MonadConcatAll(sum)([]int{1, 2, 3, 4}, 0)
	fmt.Println(result)
	// Output:
	// 10
}

func ExampleGenericConcatAll() {
	// GenericConcatAll works with any named slice type.
	type Scores []int
	sum := N.SemigroupSum[int]()
	result := semigroup.GenericConcatAll[Scores](sum)(0)(Scores{10, 20, 30})
	fmt.Println(result)
	// Output:
	// 60
}

func ExampleGenericMonadConcatAll() {
	// GenericMonadConcatAll is the uncurried form for named slice types.
	type Scores []int
	sum := N.SemigroupSum[int]()
	result := semigroup.GenericMonadConcatAll[Scores](sum)(Scores{10, 20, 30}, 0)
	fmt.Println(result)
	// Output:
	// 60
}

func ExampleAltSemigroup() {
	// AltSemigroup combines two Option values using the Alt operation:
	// the first Some wins; None falls through to the second value.
	sg := semigroup.AltSemigroup(O.MonadAlt[int])

	getOrNeg1 := O.GetOrElse(lazy.Of(-1))

	result1 := sg.Concat(O.Some(1), O.Some(2)) // first wins
	result2 := sg.Concat(O.None[int](), O.Some(2)) // fallback to second
	result3 := sg.Concat(O.None[int](), O.None[int]()) // both absent

	fmt.Println(O.IsSome(result1), getOrNeg1(result1))
	fmt.Println(O.IsSome(result2), getOrNeg1(result2))
	fmt.Println(O.IsSome(result3))
	// Output:
	// true 1
	// true 2
	// false
}

func ExampleApplySemigroup() {
	// ApplySemigroup lifts a Semigroup into an applicative context.
	// Here the int sum semigroup is lifted into Option: both values must be
	// Some for the result to be Some; if either is None the result is None.
	sum := N.SemigroupSum[int]()
	optSG := semigroup.ApplySemigroup(
		O.MonadMap[int, func(int) int],
		O.MonadAp[int, int],
		sum,
	)

	getOrNeg1 := O.GetOrElse(lazy.Of(-1))

	result1 := optSG.Concat(O.Some(5), O.Some(10))
	result2 := optSG.Concat(O.None[int](), O.Some(10))

	fmt.Println(O.IsSome(result1), getOrNeg1(result1))
	fmt.Println(O.IsSome(result2))
	// Output:
	// true 15
	// false
}
