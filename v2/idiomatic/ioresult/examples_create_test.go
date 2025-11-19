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

package ioresult

import (
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
)

func ExampleIOResult_creation() {
	// Build an IOResult
	leftValue := Left[string](fmt.Errorf("some error"))
	rightValue := Right("value")

	// Convert from Either
	eitherValue := E.Of[error](42)
	ioFromEither := FromEither(eitherValue)

	// some predicate
	isEven := func(num int) (int, error) {
		if num%2 == 0 {
			return num, nil
		}
		return 0, fmt.Errorf("%d is an odd number", num)
	}
	fromEven := Eitherize1(isEven)
	leftFromPred := fromEven(3)
	rightFromPred := fromEven(4)

	// Convert results to Either for display
	val1, err1 := leftValue()
	if err1 != nil {
		fmt.Printf("Left[*errors.errorString](%s)\n", err1.Error())
	} else {
		fmt.Printf("Right[string](%s)\n", val1)
	}

	val2, err2 := rightValue()
	if err2 != nil {
		fmt.Printf("Left[*errors.errorString](%s)\n", err2.Error())
	} else {
		fmt.Printf("Right[string](%s)\n", val2)
	}

	val3, err3 := ioFromEither()
	if err3 != nil {
		fmt.Printf("Left[*errors.errorString](%s)\n", err3.Error())
	} else {
		fmt.Printf("Right[int](%d)\n", val3)
	}

	val4, err4 := leftFromPred()
	if err4 != nil {
		fmt.Printf("Left[*errors.errorString](%s)\n", err4.Error())
	} else {
		fmt.Printf("Right[int](%d)\n", val4)
	}

	val5, err5 := rightFromPred()
	if err5 != nil {
		fmt.Printf("Left[*errors.errorString](%s)\n", err5.Error())
	} else {
		fmt.Printf("Right[int](%d)\n", val5)
	}

	// Output:
	// Left[*errors.errorString](some error)
	// Right[string](value)
	// Right[int](42)
	// Left[*errors.errorString](3 is an odd number)
	// Right[int](4)

}
