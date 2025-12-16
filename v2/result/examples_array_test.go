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
	S "github.com/IBM/fp-go/v2/string"
)

func ExampleTraverseArray() {
	parse := R.Eitherize1(strconv.Atoi)
	result := R.TraverseArray(parse)([]string{"1", "2", "3"})
	fmt.Println(R.GetOrElse(F.Constant1[error]([]int{}))(result))
	// Output:
	// [1 2 3]
}

func ExampleTraverseArrayWithIndex() {
	validate := func(i int, s string) R.Result[string] {
		if S.IsNonEmpty(s) {
			return R.Of(fmt.Sprintf("%d:%s", i, s))
		}
		return R.Left[string](fmt.Errorf("empty at index %d", i))
	}
	result := R.TraverseArrayWithIndex(validate)([]string{"a", "b"})
	fmt.Println(R.GetOrElse(F.Constant1[error]([]string{}))(result))
	// Output:
	// [0:a 1:b]
}

func ExampleSequenceArray() {
	eithers := []R.Result[int]{
		R.Of(1),
		R.Of(2),
		R.Of(3),
	}
	result := R.SequenceArray(eithers)
	fmt.Println(R.GetOrElse(F.Constant1[error]([]int{}))(result))
	// Output:
	// [1 2 3]
}

func ExampleCompactArray() {
	eithers := []R.Result[int]{
		R.Of(1),
		R.Left[int](errors.New("error")),
		R.Of(3),
	}
	result := R.CompactArray(eithers)
	fmt.Println(result)
	// Output:
	// [1 3]
}
