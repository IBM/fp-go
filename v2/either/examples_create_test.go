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

package either

import (
	"fmt"

	"github.com/IBM/fp-go/v2/errors"
)

func ExampleEither_creation() {
	// Build an Either
	leftValue := Left[string](fmt.Errorf("some error"))
	rightValue := Right[error]("value")

	// Build from a value
	fromNillable := FromNillable[string](fmt.Errorf("value was nil"))
	leftFromNil := fromNillable(nil)
	value := "value"
	rightFromPointer := fromNillable(&value)

	// some predicate
	isEven := func(num int) bool {
		return num%2 == 0
	}
	fromEven := FromPredicate(isEven, errors.OnSome[int]("%d is an odd number"))
	leftFromPred := fromEven(3)
	rightFromPred := fromEven(4)

	fmt.Println(leftValue)
	fmt.Println(rightValue)
	fmt.Println(leftFromNil)
	fmt.Println(IsRight(rightFromPointer))
	fmt.Println(leftFromPred)
	fmt.Println(rightFromPred)

	// Output:
	// Left[*errors.errorString](some error)
	// Right[string](value)
	// Left[*errors.errorString](value was nil)
	// true
	// Left[*errors.errorString](3 is an odd number)
	// Right[int](4)

}
