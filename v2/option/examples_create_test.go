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
