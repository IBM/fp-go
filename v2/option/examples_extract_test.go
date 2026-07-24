// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

func ExampleOption_extraction() {

	noneValue := None[int]()
	someValue := Of(42)

	// Convert Option[T] to T
	fromNone, okFromNone := Unwrap(noneValue) // 0, false
	fromSome, okFromSome := Unwrap(someValue) // 42, true

	// Convert Option[T] with a default value
	noneWithDefault := F.Pipe1(noneValue, GetOrElse(F.Constant(0))) // 0
	someWithDefault := F.Pipe1(someValue, GetOrElse(F.Constant(0))) // 42

	// Apply a different function on None/Some(...)
	doubleOrZero := Fold(
		F.Constant(0), // none case
		N.Mul(2),      // some case
	) // func(ma Option[int]) int

	doubleFromNone := F.Pipe1(noneValue, doubleOrZero) // 0
	doubleFromSome := F.Pipe1(someValue, doubleOrZero) // 84

	// Pro-tip: Fold is short for the following:
	doubleOfZeroBis := F.Flow2(
		Map(N.Mul(2)),            // some case
		GetOrElse(F.Constant(0)), // none case
	)
	doubleFromNoneBis := F.Pipe1(noneValue, doubleOfZeroBis) // 0
	doubleFromSomeBis := F.Pipe1(someValue, doubleOfZeroBis) // 84

	fmt.Printf("%d, %t\n", fromNone, okFromNone)
	fmt.Printf("%d, %t\n", fromSome, okFromSome)
	fmt.Println(noneWithDefault)
	fmt.Println(someWithDefault)
	fmt.Println(doubleFromNone)
	fmt.Println(doubleFromSome)
	fmt.Println(doubleFromNoneBis)
	fmt.Println(doubleFromSomeBis)

	// Output:
	// 0, false
	// 42, true
	// 0
	// 42
	// 0
	// 84
	// 0
	// 84
}

// ExampleFold_basic demonstrates eliminating an Option into a string label.
//
// The None branch is a thunk (func() B) because None carries no payload. Compare
// with predicate.ExampleFold where both branches receive A — because a bool tag
// alone carries no payload, so A must be threaded through explicitly.
func ExampleFold_basic() {
	classify := Fold(
		func() string { return "absent" },
		func(n int) string { return fmt.Sprintf("present: %d", n) },
	)
	fmt.Println(classify(Some(42)))
	fmt.Println(classify(None[int]()))
	// Output:
	// present: 42
	// absent
}

// ExampleFold_fromPredicate shows that predicate.Fold and option.Fold are
// interderivable via FromPredicate.
//
// Any call to predicate.Fold can be expressed as option.Fold composed with
// FromPredicate, at the cost of closing over the input in the None branch:
//
//	predicate.Fold(onFalse, onTrue)(p)(a)
//	  == option.Fold(func() B { return onFalse(a) }, onTrue)(FromPredicate(p)(a))
func ExampleFold_fromPredicate() {
	isPositive := func(n int) bool { return n > 0 }

	// Direct: predicate.Fold — both branches receive n
	directClassify := func(n int) string {
		if isPositive(n) {
			return fmt.Sprintf("%d is positive", n)
		}
		return fmt.Sprintf("%d is not positive", n)
	}

	// Via option.Fold + FromPredicate — equivalent, but None branch closes over n
	viaOption := func(n int) string {
		return Fold(
			func() string { return fmt.Sprintf("%d is not positive", n) },
			func(v int) string { return fmt.Sprintf("%d is positive", v) },
		)(FromPredicate(isPositive)(n))
	}

	for _, n := range []int{5, 0, -3} {
		d := directClassify(n)
		v := viaOption(n)
		fmt.Println(d, "|", v, "|", d == v)
	}
	// Output:
	// 5 is positive | 5 is positive | true
	// 0 is not positive | 0 is not positive | true
	// -3 is not positive | -3 is not positive | true
}
