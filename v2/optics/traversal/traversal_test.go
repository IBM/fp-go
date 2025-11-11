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

package traversal

import (
	"testing"

	AR "github.com/IBM/fp-go/v2/array"
	C "github.com/IBM/fp-go/v2/constant"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	AT "github.com/IBM/fp-go/v2/optics/traversal/array/const"
	AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {

	as := AR.From(1, 2, 3)

	tr := AT.FromArray[[]int, int](AR.Monoid[int]())

	sa := F.Pipe1(
		Id[[]int, C.Const[[]int, []int]](),
		Compose[[]int, []int, int, C.Const[[]int, []int]](tr),
	)

	getall := GetAll[[]int, int](as)(sa)

	assert.Equal(t, AR.From(1, 2, 3), getall)
}

func TestFold(t *testing.T) {

	monoidSum := N.MonoidSum[int]()

	as := AR.From(1, 2, 3)

	tr := AT.FromArray[int, int](monoidSum)

	sa := F.Pipe1(
		Id[[]int, C.Const[int, []int]](),
		Compose[[]int, []int, int, C.Const[int, []int]](tr),
	)

	folded := Fold(sa)(as)

	assert.Equal(t, 6, folded)
}

func TestTraverse(t *testing.T) {

	as := AR.From(1, 2, 3)

	tr := AI.FromArray[int]()

	sa := F.Pipe1(
		Id[[]int, []int](),
		Compose[[]int, []int, int, []int](tr),
	)

	res := sa(utils.Double)(as)

	assert.Equal(t, AR.From(2, 4, 6), res)
}
