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
	"github.com/IBM/fp-go/v2/predicate"
)

func ExampleNot() {
	isPositive := func(n int) bool { return n > 0 }
	isNotPositive := predicate.Not(isPositive)
	fmt.Println(isNotPositive(5))
	fmt.Println(isNotPositive(-3))
	// Output:
	// false
	// true
}

func ExampleAnd() {
	isPositive := func(n int) bool { return n > 0 }
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
	isPositive := func(n int) bool { return n > 0 }
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
