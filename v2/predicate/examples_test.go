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

package predicate_test

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/predicate"
)

func ExampleNot() {
	isPositive := N.MoreThan(0)
	isNotPositive := predicate.Not(isPositive)
	fmt.Println(isNotPositive(5))
	fmt.Println(isNotPositive(-3))
	// Output:
	// false
	// true
}

func ExampleAnd() {
	isPositive := N.MoreThan(0)
	isEven := func(n int) bool { return n%2 == 0 }
	isPositiveAndEven := F.Pipe1(isPositive, predicate.And(isEven))
	fmt.Println(isPositiveAndEven(4))
	fmt.Println(isPositiveAndEven(-2))
	fmt.Println(isPositiveAndEven(3))
	// Output:
	// true
	// false
	// false
}

func ExampleOr() {
	isPositive := N.MoreThan(0)
	isEven := func(n int) bool { return n%2 == 0 }
	isPositiveOrEven := F.Pipe1(isPositive, predicate.Or(isEven))
	fmt.Println(isPositiveOrEven(4))
	fmt.Println(isPositiveOrEven(-2))
	fmt.Println(isPositiveOrEven(3))
	fmt.Println(isPositiveOrEven(-3))
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleContraMap() {
	type Person struct{ Age int }
	isAdult := func(age int) bool { return age >= 18 }
	getAge := func(p Person) int { return p.Age }
	isPersonAdult := F.Pipe1(isAdult, predicate.ContraMap(getAge))
	fmt.Println(isPersonAdult(Person{Age: 25}))
	fmt.Println(isPersonAdult(Person{Age: 15}))
	// Output:
	// true
	// false
}

func ExampleAlways() {
	alwaysTrue := predicate.Always[int]()
	fmt.Println(alwaysTrue(42))
	fmt.Println(alwaysTrue(-10))
	fmt.Println(alwaysTrue(0))
	// Output:
	// true
	// true
	// true
}

func ExampleNever() {
	neverTrue := predicate.Never[int]()
	fmt.Println(neverTrue(42))
	fmt.Println(neverTrue(-10))
	fmt.Println(neverTrue(0))
	// Output:
	// false
	// false
	// false
}

func ExampleAlways_withAnd() {
	isPositive := N.MoreThan(0)
	// Always AND isPositive == isPositive
	combined := F.Pipe1(predicate.Always[int](), predicate.And(isPositive))
	fmt.Println(combined(5))
	fmt.Println(combined(-5))
	// Output:
	// true
	// false
}

func ExampleNever_withOr() {
	isPositive := N.MoreThan(0)
	// Never OR isPositive == isPositive
	combined := F.Pipe1(predicate.Never[int](), predicate.Or(isPositive))
	fmt.Println(combined(5))
	fmt.Println(combined(-5))
	// Output:
	// true
	// false
}

// ExampleFold_basic demonstrates eliminating a predicate into a string label.
//
// Both branch functions receive the tested value because the bool tag alone carries
// no payload — compare with option.ExampleFold where the None branch is a thunk.
func ExampleFold_basic() {
	isPositive := N.MoreThan(0)
	classify := F.Pipe1(isPositive, predicate.Fold(
		func(n int) string { return "not positive" },
		func(n int) string { return "positive" },
	))
	fmt.Println(classify(5))
	fmt.Println(classify(0))
	fmt.Println(classify(-3))
	// Output:
	// positive
	// not positive
	// not positive
}

// ExampleFold_withValue shows that both branches always receive A, allowing the
// original value to be used in the result regardless of which branch is taken.
// This mirrors option.Fold's onSome branch, but predicate.Fold must also pass A
// to onFalse because the bool tag carries no payload of its own.
func ExampleFold_withValue() {
	isEven := func(n int) bool { return n%2 == 0 }
	describe := F.Pipe1(isEven, predicate.Fold(
		func(n int) string { return fmt.Sprintf("%d is odd", n) },
		func(n int) string { return fmt.Sprintf("%d is even", n) },
	))
	fmt.Println(describe(4))
	fmt.Println(describe(7))
	// Output:
	// 4 is even
	// 7 is odd
}
