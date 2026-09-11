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

package monoid_test

import (
	"fmt"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
)

func ExampleMakeMonoid() {
	addMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)
	fmt.Println(addMonoid.Concat(5, 3))
	fmt.Println(addMonoid.Empty())
	// Output:
	// 8
	// 0
}

func ExampleReverse() {
	subMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a - b },
		0,
	)
	reversed := monoid.Reverse(subMonoid)
	fmt.Println(subMonoid.Concat(10, 3))  // 10 - 3
	fmt.Println(reversed.Concat(10, 3))   // 3 - 10
	fmt.Println(reversed.Empty())         // identity unchanged
	// Output:
	// 7
	// -7
	// 0
}

func ExampleToSemigroup() {
	addMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)
	sg := monoid.ToSemigroup(addMonoid)
	fmt.Println(sg.Concat(5, 3))
	// Output:
	// 8
}

func ExampleConcatAll() {
	addMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)
	concatAll := monoid.ConcatAll(addMonoid)
	fmt.Println(concatAll([]int{1, 2, 3, 4, 5}))
	fmt.Println(concatAll([]int{}))
	// Output:
	// 15
	// 0
}

func ExampleFold() {
	mulMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a * b },
		1,
	)
	fold := monoid.Fold(mulMonoid)
	fmt.Println(fold([]int{2, 3, 4}))
	fmt.Println(fold([]int{}))
	// Output:
	// 24
	// 1
}

func ExampleGenericConcatAll() {
	type IntSlice []int
	addMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)
	concatAll := monoid.GenericConcatAll[IntSlice](addMonoid)
	fmt.Println(concatAll(IntSlice{10, 20, 30}))
	fmt.Println(concatAll(IntSlice{}))
	// Output:
	// 60
	// 0
}

func ExampleFunctionMonoid() {
	intAddMonoid := monoid.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)
	funcMonoid := monoid.FunctionMonoid[string](intAddMonoid)

	f1 := func(s string) int { return len(s) }
	f2 := func(s string) int { return len(s) * 2 }

	combined := funcMonoid.Concat(f1, f2)
	fmt.Println(combined("hello")) // 5 + 10

	emptyFunc := funcMonoid.Empty()
	fmt.Println(emptyFunc("anything")) // identity: 0
	// Output:
	// 15
	// 0
}

func ExampleVoidMonoid() {
	m := monoid.VoidMonoid()
	result := m.Concat(function.VOID, function.VOID)
	empty := m.Empty()
	fmt.Println(result == function.VOID)
	fmt.Println(empty == function.VOID)
	// Output:
	// true
	// true
}

func ExampleApplicativeMonoid() {
	// ApplicativeMonoid lifts a monoid into an applicative functor context.
	// Here the string monoid is lifted into Option: both values must be Some.
	m := monoid.ApplicativeMonoid(O.Of[string], O.Map[string, func(string) string], O.Ap[string, string], S.Monoid)

	getOrEmpty := O.GetOrElse(lazy.Of(""))

	result1 := m.Concat(O.Some("Hello, "), O.Some("World"))
	result2 := m.Concat(O.None[string](), O.Some("World"))

	fmt.Println(getOrEmpty(result1))
	fmt.Println(O.IsSome(result2))
	fmt.Println(getOrEmpty(m.Empty()))
	// Output:
	// Hello, World
	// false
	//
}

func ExampleAltMonoid() {
	// AltMonoid uses the first Some, falling back to the second.
	m := monoid.AltMonoid(O.None[int], O.Alt[int])

	getOrNeg1 := O.GetOrElse(lazy.Of(-1))

	result1 := m.Concat(O.Some(1), O.Some(2))     // first wins
	result2 := m.Concat(O.None[int](), O.Some(2)) // fall back to second
	result3 := m.Concat(O.None[int](), m.Empty())  // both empty

	fmt.Println(getOrNeg1(result1))
	fmt.Println(getOrNeg1(result2))
	fmt.Println(O.IsSome(result3))
	// Output:
	// 1
	// 2
	// false
}

func ExampleAlternativeMonoid() {
	// AlternativeMonoid combines applicative and alternative semantics.
	// When both values are Some, their int values are summed.
	// When either is None, the Alt fallback takes over.
	intAddMonoid := N.MonoidSum[int]()
	m := monoid.AlternativeMonoid(
		O.Of[int],
		O.Map[int, func(int) int],
		O.Ap[int, int],
		O.Alt[int],
		intAddMonoid,
	)

	getOrNeg1 := O.GetOrElse(lazy.Of(-1))

	result1 := m.Concat(O.Some(5), O.Some(3))      // applicative: 5+3
	result2 := m.Concat(O.None[int](), O.Some(3))  // alt fallback: Some(3)
	result3 := m.Concat(O.Some(5), O.None[int]())  // alt fallback: Some(5)

	fmt.Println(getOrNeg1(result1))
	fmt.Println(getOrNeg1(result2))
	fmt.Println(getOrNeg1(result3))
	// Output:
	// 8
	// 3
	// 5
}
