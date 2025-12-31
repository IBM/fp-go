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

package readerioeither

import (
	"testing"

	"github.com/IBM/fp-go/v2/either"
	L "github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

type testCtx struct {
	value int
}

func TestApplicativeMonoidSeq(t *testing.T) {
	ctx := testCtx{value: 1}
	m := ApplicativeMonoidSeq[testCtx, error](N.MonoidSum[int]())

	t.Run("combines two Right values sequentially", func(t *testing.T) {
		a := Of[testCtx, error](10)
		b := Of[testCtx, error](20)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 30, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("returns Left if first value is Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Of[testCtx, error](20)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsLeft(outcome))
	})

	t.Run("returns Left if second value is Left", func(t *testing.T) {
		a := Of[testCtx, error](10)
		b := Left[testCtx, int](assert.AnError)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsLeft(outcome))
	})

	t.Run("empty returns identity", func(t *testing.T) {
		empty := m.Empty()
		outcome := empty(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 0, either.GetOrElse(func(error) int { return -1 })(outcome))
	})
}

func TestApplicativeMonoidPar(t *testing.T) {
	ctx := testCtx{value: 1}
	m := ApplicativeMonoidPar[testCtx, error](N.MonoidSum[int]())

	t.Run("combines two Right values in parallel", func(t *testing.T) {
		a := Of[testCtx, error](15)
		b := Of[testCtx, error](25)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 40, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("returns Left if first value is Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Of[testCtx, error](25)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsLeft(outcome))
	})

	t.Run("empty returns identity", func(t *testing.T) {
		empty := m.Empty()
		outcome := empty(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 0, either.GetOrElse(func(error) int { return -1 })(outcome))
	})
}

func TestAlternativeMonoid(t *testing.T) {
	ctx := testCtx{value: 1}
	m := AlternativeMonoid[testCtx, error](N.MonoidSum[int]())

	t.Run("combines two Right values", func(t *testing.T) {
		a := Of[testCtx, error](5)
		b := Of[testCtx, error](7)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 12, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("uses second value if first is Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Of[testCtx, error](7)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 7, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("returns Left if both are Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Left[testCtx, int](assert.AnError)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsLeft(outcome))
	})
}

func TestAltMonoid(t *testing.T) {
	ctx := testCtx{value: 1}
	zero := L.Of(Of[testCtx, error](0))
	m := AltMonoid(zero)

	t.Run("returns first Right value", func(t *testing.T) {
		a := Of[testCtx, error](42)
		b := Of[testCtx, error](100)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 42, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("returns second value if first is Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Of[testCtx, error](100)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 100, either.GetOrElse(func(error) int { return 0 })(outcome))
	})

	t.Run("returns Left if both are Left", func(t *testing.T) {
		a := Left[testCtx, int](assert.AnError)
		b := Left[testCtx, int](assert.AnError)

		result := m.Concat(a, b)
		outcome := result(ctx)()

		assert.True(t, either.IsLeft(outcome))
	})

	t.Run("empty returns zero value", func(t *testing.T) {
		empty := m.Empty()
		outcome := empty(ctx)()

		assert.True(t, either.IsRight(outcome))
		assert.Equal(t, 0, either.GetOrElse(func(error) int { return -1 })(outcome))
	})
}
