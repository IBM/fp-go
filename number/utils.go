// Copyright (c) 2023 IBM Corp.
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

package number

import (
	C "github.com/IBM/fp-go/constraints"
)

type Number interface {
	C.Integer | C.Float | C.Complex
}

// Add is a curried function used to add two numbers
func Add[T Number](left T) func(T) T {
	return func(right T) T {
		return left + right
	}
}

// Inc is a function that increments a number
func Inc[T Number](value T) T {
	return value + 1
}

// Min takes the minimum of two values. If they are considered equal, the first argument is chosen
func Min[A C.Ordered](a, b A) A {
	if a < b {
		return a
	}
	return b
}

// Max takes the maximum of two values. If they are considered equal, the first argument is chosen
func Max[A C.Ordered](a, b A) A {
	if a > b {
		return a
	}
	return b
}
