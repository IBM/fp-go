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

package readerioresult

import (
	"context"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
)

// ExampleTraverseArrayParN processes a list with at most two computations in flight.
func ExampleTraverseArrayParN() {
	double := F.Flow2(N.Mul(2), Of[int])

	result, err := exRun(context.Background(), F.Pipe1(
		[]int{1, 2, 3, 4, 5},
		TraverseArrayParN(2, double),
	))

	fmt.Println(result, err)
	// Output: [2 4 6 8 10] <nil>
}

// ExampleTraverseArrayParN_failure shows that the first failure is returned.
func ExampleTraverseArrayParN_failure() {
	parse := F.Flow2(R.Eitherize1(strconv.Atoi), FromResult[int])

	_, err := exRun(context.Background(), F.Pipe1(
		[]string{"1", "x", "3"},
		TraverseArrayParN(2, parse),
	))

	fmt.Println(err)
	// Output: strconv.Atoi: parsing "x": invalid syntax
}

// ExampleWithSemaphore limits concurrent access to a shared resource to one slot.
func ExampleWithSemaphore() {
	sem := make(chan Void, 1)
	guarded := F.Flow2(
		F.Flow2(N.Add(10), Of[int]),
		WithSemaphore[int](sem),
	)

	result, err := exRun(context.Background(), F.Pipe1(
		[]int{1, 2, 3},
		TraverseArrayPar(guarded),
	))

	fmt.Println(result, err)
	// Output: [11 12 13] <nil>
}

// ExampleAnnotateLogs adds a request ID to every log record written inside a computation.
func ExampleAnnotateLogs() {
	fmt.Println(exRunLogged(F.Pipe3(
		exFetchUser(42),
		Map(exGetName),
		TapSLog[string]("fetched user"),
		AnnotateLogs[string]("requestID", "r-1"),
	)))

	// Output:
	// level=INFO msg="fetched user" requestID=r-1 value=Alice
	// Alice <nil>
}
