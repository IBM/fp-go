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

import "fmt"

func ExampleOption_creation() {

	// Build an Option
	none1 := None[int]()
	some1 := Some("value")

	// Build from a value
	fromNillable := FromNillable[string]
	nonFromNil := fromNillable(nil) // None[*string]
	value := "value"
	someFromPointer := fromNillable(&value) // Some[*string](xxx)

	// some predicate
	isEven := func(num int) bool {
		return num%2 == 0
	}

	fromEven := FromPredicate(isEven)
	noneFromPred := fromEven(3) // None[int]
	someFromPred := fromEven(4) // Some[int](4)

	fmt.Println(none1)
	fmt.Println(some1)
	fmt.Println(nonFromNil)
	fmt.Println(IsSome(someFromPointer))
	fmt.Println(noneFromPred)
	fmt.Println(someFromPred)

	// Output:
	// None[int]
	// Some[string](value)
	// None[string]
	// true
	// None[int]
	// Some[int](4)
}

// ExampleFromNillable2_nil shows that a nil pointer produces None.
func ExampleFromNillable2_nil() {
	result := FromNillable2[int](nil)
	fmt.Println(result)
	// Output:
	// None[int]
}

// ExampleFromNillable2_nonNil shows that a non-nil pointer produces Some with
// the dereferenced value — not the pointer — making the result suitable for
// direct use with Map, Chain, and other Option combinators.
func ExampleFromNillable2_nonNil() {
	val := 42
	result := FromNillable2(&val)
	fmt.Println(result)
	// Output:
	// Some[int](42)
}

// ExampleToNillable2_some shows that Some(v) converts back to a non-nil pointer.
func ExampleToNillable2_some() {
	ptr := ToNillable2(Some(42))
	fmt.Println(ptr != nil)
	fmt.Println(*ptr)
	// Output:
	// true
	// 42
}

// ExampleToNillable2_none shows that None converts to a nil pointer.
func ExampleToNillable2_none() {
	ptr := ToNillable2(None[int]())
	fmt.Println(ptr == nil)
	// Output:
	// true
}
