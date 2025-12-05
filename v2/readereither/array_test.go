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

package readereither

import (
	"context"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArrayWithIndex(t *testing.T) {
	ctx := context.Background()

	t.Run("empty array", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a string) ReaderEither[context.Context, error, string] {
			return Of[context.Context, error](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{})(ctx)
		assert.Equal(t, ET.Right[error]([]string{}), result)
	})

	t.Run("transformation with index", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a string) ReaderEither[context.Context, error, string] {
			return Of[context.Context, error](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{"a", "b", "c"})(ctx)
		assert.Equal(t, ET.Right[error]([]string{"0:a", "1:b", "2:c"}), result)
	})

	t.Run("fails at specific index", func(t *testing.T) {
		expectedErr := fmt.Errorf("error at index 1")
		f := TraverseArrayWithIndex(func(i int, a string) ReaderEither[context.Context, error, string] {
			if i == 1 {
				return Left[context.Context, string](expectedErr)
			}
			return Of[context.Context, error](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{"a", "b", "c"})(ctx)
		assert.Equal(t, ET.Left[[]string](expectedErr), result)
	})
}

func TestSequenceArray(t *testing.T) {
	ctx := context.Background()

	t.Run("basic functionality", func(t *testing.T) {
		n := 10
		readers := A.MakeBy(n, Of[context.Context, error, int])
		exp := ET.Of[error](A.MakeBy(n, F.Identity[int]))

		g := F.Pipe1(
			readers,
			SequenceArray[context.Context, error, int],
		)

		assert.Equal(t, exp, g(ctx))
	})

	t.Run("empty array", func(t *testing.T) {
		computations := []ReaderEither[context.Context, error, int]{}
		result := SequenceArray(computations)(ctx)
		assert.Equal(t, ET.Right[error]([]int{}), result)
	})

	t.Run("all successful", func(t *testing.T) {
		computations := []ReaderEither[context.Context, error, int]{
			Of[context.Context, error](1),
			Of[context.Context, error](2),
			Of[context.Context, error](3),
		}
		result := SequenceArray(computations)(ctx)
		assert.Equal(t, ET.Right[error]([]int{1, 2, 3}), result)
	})

	t.Run("first computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("first error")
		computations := []ReaderEither[context.Context, error, int]{
			Left[context.Context, int](expectedErr),
			Of[context.Context, error](2),
			Of[context.Context, error](3),
		}
		result := SequenceArray(computations)(ctx)
		assert.Equal(t, ET.Left[[]int](expectedErr), result)
	})

	t.Run("middle computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("middle error")
		computations := []ReaderEither[context.Context, error, int]{
			Of[context.Context, error](1),
			Left[context.Context, int](expectedErr),
			Of[context.Context, error](3),
		}
		result := SequenceArray(computations)(ctx)
		assert.Equal(t, ET.Left[[]int](expectedErr), result)
	})
}
