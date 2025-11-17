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

package result

import (
	"fmt"

	"github.com/IBM/fp-go/v2/errors"
)

func Example_creation() {
	// Build an Either
	leftValue, leftErr := Left[string](fmt.Errorf("some error"))
	rightValue, rightErr := Right("value")

	// Build from a value
	fromNillable := FromNillable[string](fmt.Errorf("value was nil"))
	leftFromNil, nilErr := fromNillable(nil)
	value := "value"
	rightFromPointer, ptrErr := fromNillable(&value)

	// some predicate
	isEven := func(num int) bool {
		return num%2 == 0
	}
	fromEven := FromPredicate(isEven, errors.OnSome[int]("%d is an odd number"))
	leftFromPred, predErrOdd := fromEven(3)
	rightFromPred, predErrEven := fromEven(4)

	fmt.Println(ToString(leftValue, leftErr))
	fmt.Println(ToString(rightValue, rightErr))
	fmt.Println(ToString(leftFromNil, nilErr))
	fmt.Println(IsRight(rightFromPointer, ptrErr))
	fmt.Println(ToString(leftFromPred, predErrOdd))
	fmt.Println(ToString(rightFromPred, predErrEven))

	// Output:
	// Left(some error)
	// Right[string](value)
	// Left(value was nil)
	// true
	// Left(3 is an odd number)
	// Right[int](4)

}
