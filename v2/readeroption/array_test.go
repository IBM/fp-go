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

package readeroption

import (
	"context"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	n := 10

	readers := A.MakeBy(n, Of[context.Context, int])
	exp := O.Of(A.MakeBy(n, F.Identity[int]))

	g := F.Pipe1(
		readers,
		SequenceArray[context.Context, int],
	)

	assert.Equal(t, exp, g(context.Background()))
}

func TestTraverseArray(t *testing.T) {
	// Function that doubles a number if it's positive
	doubleIfPositive := func(x int) ReaderOption[context.Context, int] {
		if x > 0 {
			return Of[context.Context](x * 2)
		}
		return None[context.Context, int]()
	}

	// Test with all positive numbers
	input1 := []int{1, 2, 3}
	g1 := F.Pipe1(
		Of[context.Context](input1),
		Chain(TraverseArray(doubleIfPositive)),
	)
	assert.Equal(t, O.Of([]int{2, 4, 6}), g1(context.Background()))

	// Test with a negative number (should return None)
	input2 := []int{1, -2, 3}
	g2 := F.Pipe1(
		Of[context.Context](input2),
		Chain(TraverseArray(doubleIfPositive)),
	)
	assert.Equal(t, O.None[[]int](), g2(context.Background()))

	// Test with empty array
	input3 := []int{}
	g3 := F.Pipe1(
		Of[context.Context](input3),
		Chain(TraverseArray(doubleIfPositive)),
	)
	assert.Equal(t, O.Of([]int{}), g3(context.Background()))
}

func TestTraverseArrayWithIndex(t *testing.T) {
	// Function that multiplies value by its index if index is even
	multiplyByIndexIfEven := func(idx int, x int) ReaderOption[context.Context, int] {
		if idx%2 == 0 {
			return Of[context.Context](x * idx)
		}
		return Of[context.Context](x)
	}

	input := []int{10, 20, 30, 40}
	g := TraverseArrayWithIndex(multiplyByIndexIfEven)(input)

	// Expected: [10*0, 20, 30*2, 40] = [0, 20, 60, 40]
	assert.Equal(t, O.Of([]int{0, 20, 60, 40}), g(context.Background()))
}

func TestTraverseArrayWithIndexNone(t *testing.T) {
	// Function that returns None for odd indices
	noneForOdd := func(idx int, x int) ReaderOption[context.Context, int] {
		if idx%2 == 0 {
			return Of[context.Context](x)
		}
		return None[context.Context, int]()
	}

	input := []int{10, 20, 30}
	g := TraverseArrayWithIndex(noneForOdd)(input)

	// Should return None because index 1 returns None
	assert.Equal(t, O.None[[]int](), g(context.Background()))
}
