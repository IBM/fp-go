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

package result_test

import (
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/result"
)

func ExampleCurry0() {
	getConfig := func() (string, error) { return "config", nil }
	curried := R.Curry0(getConfig)
	result := curried()
	fmt.Println(R.GetOrElse(F.Constant1[error](""))(result))
	// Output:
	// config
}

func ExampleCurry1() {
	curried := R.Curry1(strconv.Atoi)
	result := curried("42")
	fmt.Println(R.GetOrElse(F.Constant1[error](0))(result))
	// Output:
	// 42
}

func ExampleCurry2() {
	divide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, errors.New("div by zero")
		}
		return a / b, nil
	}
	curried := R.Curry2(divide)
	result := curried(10)(2)
	fmt.Println(R.GetOrElse(F.Constant1[error](0))(result))
	// Output:
	// 5
}

func ExampleUncurry0() {
	curried := func() R.Result[string] { return R.Of("value") }
	uncurried := R.Uncurry0(curried)
	result, err := uncurried()
	fmt.Println(result, err)
	// Output:
	// value <nil>
}

func ExampleUncurry1() {
	curried := func(x int) R.Result[string] { return R.Of(strconv.Itoa(x)) }
	uncurried := R.Uncurry1(curried)
	result, err := uncurried(42)
	fmt.Println(result, err)
	// Output:
	// 42 <nil>
}
