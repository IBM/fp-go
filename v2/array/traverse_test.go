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

package array

import (
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type ArrayType = []int

func TestTraverse(t *testing.T) {

	traverse := Traverse(
		O.Of[ArrayType],
		O.Map[ArrayType, func(int) ArrayType],
		O.Ap[ArrayType, int],

		func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.None[int]()
			}
			return O.Of(n)
		})

	assert.Equal(t, O.None[[]int](), traverse(ArrayType{1, 2}))
	assert.Equal(t, O.Of(ArrayType{1, 3}), traverse(ArrayType{1, 3}))
}
