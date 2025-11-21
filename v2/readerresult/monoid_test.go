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

package readerresult

import (
	"errors"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var (
	intAddMonoid = N.MonoidSum[int]()
	strMonoid    = S.Monoid
)

func TestApplicativeMonoid(t *testing.T) {
	rrMonoid := ApplicativeMonoid[MyContext](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := rrMonoid.Empty()
		assert.Equal(t, result.Of(0), empty(defaultContext))
	})

	t.Run("concat two success values", func(t *testing.T) {
		rr1 := Of[MyContext](5)
		rr2 := Of[MyContext](3)
		combined := rrMonoid.Concat(rr1, rr2)
		assert.Equal(t, result.Of(8), combined(defaultContext))
	})

	t.Run("concat with empty", func(t *testing.T) {
		rr := Of[MyContext](42)
		combined1 := rrMonoid.Concat(rr, rrMonoid.Empty())
		combined2 := rrMonoid.Concat(rrMonoid.Empty(), rr)

		assert.Equal(t, result.Of(42), combined1(defaultContext))
		assert.Equal(t, result.Of(42), combined2(defaultContext))
	})

	t.Run("concat with left failure", func(t *testing.T) {
		rrSuccess := Of[MyContext](5)
		rrFailure := Left[MyContext, int](testError)

		combined := rrMonoid.Concat(rrFailure, rrSuccess)
		assert.True(t, result.IsLeft(combined(defaultContext)))
	})

	t.Run("concat with right failure", func(t *testing.T) {
		rrSuccess := Of[MyContext](5)
		rrFailure := Left[MyContext, int](testError)

		combined := rrMonoid.Concat(rrSuccess, rrFailure)
		assert.True(t, result.IsLeft(combined(defaultContext)))
	})

	t.Run("concat multiple values", func(t *testing.T) {
		rr1 := Of[MyContext](1)
		rr2 := Of[MyContext](2)
		rr3 := Of[MyContext](3)
		rr4 := Of[MyContext](4)

		// Chain concat calls: ((1 + 2) + 3) + 4
		combined := rrMonoid.Concat(
			rrMonoid.Concat(
				rrMonoid.Concat(rr1, rr2),
				rr3,
			),
			rr4,
		)
		assert.Equal(t, result.Of(10), combined(defaultContext))
	})

	t.Run("string concatenation", func(t *testing.T) {
		strRRMonoid := ApplicativeMonoid[MyContext](strMonoid)

		rr1 := Of[MyContext]("Hello")
		rr2 := Of[MyContext](" ")
		rr3 := Of[MyContext]("World")

		combined := strRRMonoid.Concat(
			strRRMonoid.Concat(rr1, rr2),
			rr3,
		)
		assert.Equal(t, result.Of("Hello World"), combined(defaultContext))
	})
}

func TestAltMonoid(t *testing.T) {
	zero := func() ReaderResult[MyContext, int] {
		return Left[MyContext, int](errors.New("empty"))
	}

	rrMonoid := AltMonoid(zero)

	t.Run("empty element", func(t *testing.T) {
		empty := rrMonoid.Empty()
		assert.True(t, result.IsLeft(empty(defaultContext)))
	})

	t.Run("concat two success values - uses first", func(t *testing.T) {
		rr1 := Of[MyContext](5)
		rr2 := Of[MyContext](3)
		combined := rrMonoid.Concat(rr1, rr2)
		// AltMonoid takes the first successful value
		assert.Equal(t, result.Of(5), combined(defaultContext))
	})

	t.Run("concat failure then success", func(t *testing.T) {
		rrFailure := Left[MyContext, int](testError)
		rrSuccess := Of[MyContext](42)

		combined := rrMonoid.Concat(rrFailure, rrSuccess)
		// Should fall back to second when first fails
		assert.Equal(t, result.Of(42), combined(defaultContext))
	})

	t.Run("concat success then failure", func(t *testing.T) {
		rrSuccess := Of[MyContext](42)
		rrFailure := Left[MyContext, int](testError)

		combined := rrMonoid.Concat(rrSuccess, rrFailure)
		// Should use first successful value
		assert.Equal(t, result.Of(42), combined(defaultContext))
	})

	t.Run("concat two failures", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		rr1 := Left[MyContext, int](err1)
		rr2 := Left[MyContext, int](err2)

		combined := rrMonoid.Concat(rr1, rr2)
		// Should use second error when both fail
		assert.True(t, result.IsLeft(combined(defaultContext)))
	})

	t.Run("concat with empty", func(t *testing.T) {
		rr := Of[MyContext](42)
		combined1 := rrMonoid.Concat(rr, rrMonoid.Empty())
		combined2 := rrMonoid.Concat(rrMonoid.Empty(), rr)

		assert.Equal(t, result.Of(42), combined1(defaultContext))
		assert.Equal(t, result.Of(42), combined2(defaultContext))
	})

	t.Run("fallback chain", func(t *testing.T) {
		// Simulate trying multiple sources until one succeeds
		primary := Left[MyContext, string](errors.New("primary failed"))
		secondary := Left[MyContext, string](errors.New("secondary failed"))
		tertiary := Of[MyContext]("tertiary success")

		strZero := func() ReaderResult[MyContext, string] {
			return Left[MyContext, string](errors.New("all failed"))
		}
		strMonoid := AltMonoid(strZero)

		// Chain concat: try primary, then secondary, then tertiary
		combined := strMonoid.Concat(
			strMonoid.Concat(primary, secondary),
			tertiary,
		)
		assert.Equal(t, result.Of("tertiary success"), combined(defaultContext))
	})
}

func TestAlternativeMonoid(t *testing.T) {
	rrMonoid := AlternativeMonoid[MyContext](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := rrMonoid.Empty()
		assert.Equal(t, result.Of(0), empty(defaultContext))
	})

	t.Run("concat two success values", func(t *testing.T) {
		rr1 := Of[MyContext](5)
		rr2 := Of[MyContext](3)
		combined := rrMonoid.Concat(rr1, rr2)
		assert.Equal(t, result.Of(8), combined(defaultContext))
	})

	t.Run("concat failure then success", func(t *testing.T) {
		rrFailure := Left[MyContext, int](testError)
		rrSuccess := Of[MyContext](42)

		combined := rrMonoid.Concat(rrFailure, rrSuccess)
		// Alternative falls back to second when first fails
		assert.Equal(t, result.Of(42), combined(defaultContext))
	})

	t.Run("concat success then failure", func(t *testing.T) {
		rrSuccess := Of[MyContext](42)
		rrFailure := Left[MyContext, int](testError)

		combined := rrMonoid.Concat(rrSuccess, rrFailure)
		// Should use first successful value
		assert.Equal(t, result.Of(42), combined(defaultContext))
	})

	t.Run("concat with empty", func(t *testing.T) {
		rr := Of[MyContext](42)
		combined1 := rrMonoid.Concat(rr, rrMonoid.Empty())
		combined2 := rrMonoid.Concat(rrMonoid.Empty(), rr)

		assert.Equal(t, result.Of(42), combined1(defaultContext))
		assert.Equal(t, result.Of(42), combined2(defaultContext))
	})

	t.Run("multiple values with some failures", func(t *testing.T) {
		rr1 := Left[MyContext, int](errors.New("fail 1"))
		rr2 := Of[MyContext](5)
		rr3 := Left[MyContext, int](errors.New("fail 2"))
		rr4 := Of[MyContext](10)

		// Alternative should skip failures and accumulate successes
		combined := rrMonoid.Concat(
			rrMonoid.Concat(
				rrMonoid.Concat(rr1, rr2),
				rr3,
			),
			rr4,
		)
		// Should accumulate successful values: 5 + 10 = 15
		assert.Equal(t, result.Of(15), combined(defaultContext))
	})
}

// Test monoid laws
func TestMonoidLaws(t *testing.T) {
	rrMonoid := ApplicativeMonoid[MyContext](intAddMonoid)

	// Left identity: empty <> x == x
	t.Run("left identity", func(t *testing.T) {
		x := Of[MyContext](42)
		result1 := rrMonoid.Concat(rrMonoid.Empty(), x)(defaultContext)
		result2 := x(defaultContext)
		assert.Equal(t, result2, result1)
	})

	// Right identity: x <> empty == x
	t.Run("right identity", func(t *testing.T) {
		x := Of[MyContext](42)
		result1 := rrMonoid.Concat(x, rrMonoid.Empty())(defaultContext)
		result2 := x(defaultContext)
		assert.Equal(t, result2, result1)
	})

	// Associativity: (x <> y) <> z == x <> (y <> z)
	t.Run("associativity", func(t *testing.T) {
		x := Of[MyContext](1)
		y := Of[MyContext](2)
		z := Of[MyContext](3)

		left := rrMonoid.Concat(rrMonoid.Concat(x, y), z)(defaultContext)
		right := rrMonoid.Concat(x, rrMonoid.Concat(y, z))(defaultContext)

		assert.Equal(t, right, left)
	})
}
