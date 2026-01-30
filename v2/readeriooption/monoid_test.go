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

package readeriooption

import (
	"context"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid_BothSuccess(t *testing.T) {
	// Test combining two successful computations
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](5)
	ro2 := Of[context.Context](3)

	result := m.Concat(ro1, ro2)
	expected := O.Of(8)

	assert.Equal(t, expected, result(context.Background())())
}

func TestApplicativeMonoid_FirstFailure(t *testing.T) {
	// Test when first computation fails
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro1 := None[context.Context, int]()
	ro2 := Of[context.Context](3)

	result := m.Concat(ro1, ro2)
	expected := O.None[int]()

	assert.Equal(t, expected, result(context.Background())())
}

func TestApplicativeMonoid_SecondFailure(t *testing.T) {
	// Test when second computation fails
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](5)
	ro2 := None[context.Context, int]()

	result := m.Concat(ro1, ro2)
	expected := O.None[int]()

	assert.Equal(t, expected, result(context.Background())())
}

func TestApplicativeMonoid_BothFailure(t *testing.T) {
	// Test when both computations fail
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro1 := None[context.Context, int]()
	ro2 := None[context.Context, int]()

	result := m.Concat(ro1, ro2)
	expected := O.None[int]()

	assert.Equal(t, expected, result(context.Background())())
}

func TestApplicativeMonoid_LeftIdentity(t *testing.T) {
	// Test left identity: Concat(Empty(), x) = x
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro := Of[context.Context](5)
	result := m.Concat(m.Empty(), ro)

	assert.Equal(t, O.Of(5), result(context.Background())())
}

func TestApplicativeMonoid_RightIdentity(t *testing.T) {
	// Test right identity: Concat(x, Empty()) = x
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro := Of[context.Context](5)
	result := m.Concat(ro, m.Empty())

	assert.Equal(t, O.Of(5), result(context.Background())())
}

func TestApplicativeMonoid_Associativity(t *testing.T) {
	// Test associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](2)
	ro2 := Of[context.Context](3)
	ro3 := Of[context.Context](5)

	left := m.Concat(m.Concat(ro1, ro2), ro3)
	right := m.Concat(ro1, m.Concat(ro2, ro3))

	assert.Equal(t, O.Of(10), left(context.Background())())
	assert.Equal(t, O.Of(10), right(context.Background())())
}

func TestApplicativeMonoid_StringConcat(t *testing.T) {
	// Test with string concatenation monoid
	strConcat := S.Monoid
	m := ApplicativeMonoid[context.Context](strConcat)

	ro1 := Of[context.Context]("Hello")
	ro2 := Of[context.Context](" ")
	ro3 := Of[context.Context]("World")

	result := m.Concat(m.Concat(ro1, ro2), ro3)
	expected := O.Of("Hello World")

	assert.Equal(t, expected, result(context.Background())())
}

func TestApplicativeMonoid_WithEnvironment(t *testing.T) {
	// Test that environment is properly passed through
	type Config struct {
		Factor int
	}

	intAdd := N.MonoidSum[int]()
	m := ApplicativeMonoid[Config](intAdd)

	ro1 := func(cfg Config) IOOption[int] {
		return func() Option[int] {
			return O.Of(10 * cfg.Factor)
		}
	}
	ro2 := func(cfg Config) IOOption[int] {
		return func() Option[int] {
			return O.Of(5 * cfg.Factor)
		}
	}

	result := m.Concat(ro1, ro2)
	cfg := Config{Factor: 2}
	expected := O.Of(30) // (10*2) + (5*2) = 30

	assert.Equal(t, expected, result(cfg)())
}

func TestAlternativeMonoid_BothSuccess(t *testing.T) {
	// Test combining two successful computations
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](5)
	ro2 := Of[context.Context](3)

	result := m.Concat(ro1, ro2)
	expected := O.Of(8)

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_FirstFailure(t *testing.T) {
	// Test fallback when first computation fails
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := None[context.Context, int]()
	ro2 := Of[context.Context](10)

	result := m.Concat(ro1, ro2)
	expected := O.Of(10)

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_SecondFailure(t *testing.T) {
	// Test using first success when second fails
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](5)
	ro2 := None[context.Context, int]()

	result := m.Concat(ro1, ro2)
	expected := O.Of(5)

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_BothFailure(t *testing.T) {
	// Test when both computations fail
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := None[context.Context, int]()
	ro2 := None[context.Context, int]()

	result := m.Concat(ro1, ro2)
	expected := O.None[int]()

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_LeftIdentity(t *testing.T) {
	// Test left identity: Concat(Empty(), x) = x
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro := Of[context.Context](5)
	result := m.Concat(m.Empty(), ro)

	assert.Equal(t, O.Of(5), result(context.Background())())
}

func TestAlternativeMonoid_RightIdentity(t *testing.T) {
	// Test right identity: Concat(x, Empty()) = x
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro := Of[context.Context](5)
	result := m.Concat(ro, m.Empty())

	assert.Equal(t, O.Of(5), result(context.Background())())
}

func TestAlternativeMonoid_Associativity(t *testing.T) {
	// Test associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](2)
	ro2 := Of[context.Context](3)
	ro3 := Of[context.Context](5)

	left := m.Concat(m.Concat(ro1, ro2), ro3)
	right := m.Concat(ro1, m.Concat(ro2, ro3))

	assert.Equal(t, O.Of(10), left(context.Background())())
	assert.Equal(t, O.Of(10), right(context.Background())())
}

func TestAlternativeMonoid_FallbackChain(t *testing.T) {
	// Test chaining multiple fallbacks
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := None[context.Context, int]()
	ro2 := None[context.Context, int]()
	ro3 := Of[context.Context](7)
	ro4 := Of[context.Context](3)

	// None + None = None, then None + 7 = 7, then 7 + 3 = 10
	result := m.Concat(m.Concat(m.Concat(ro1, ro2), ro3), ro4)
	expected := O.Of(10)

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_WithEnvironment(t *testing.T) {
	// Test that environment is properly passed through with fallback
	type Config struct {
		UseCache bool
		Factor   int
	}

	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[Config](intAdd)

	cacheValue := func(cfg Config) IOOption[int] {
		return func() Option[int] {
			if cfg.UseCache {
				return O.Of(100)
			}
			return O.None[int]()
		}
	}

	dbValue := func(cfg Config) IOOption[int] {
		return func() Option[int] {
			return O.Of(50 * cfg.Factor)
		}
	}

	result := m.Concat(cacheValue, dbValue)

	// With cache enabled, both succeed so values are combined: 100 + (50*2) = 200
	cfg1 := Config{UseCache: true, Factor: 2}
	assert.Equal(t, O.Of(200), result(cfg1)())

	// With cache disabled, should fall back to DB value: 0 + (50*2) = 100
	cfg2 := Config{UseCache: false, Factor: 2}
	assert.Equal(t, O.Of(100), result(cfg2)())
}

func TestAlternativeMonoid_StringConcat(t *testing.T) {
	// Test with string concatenation and fallback
	strConcat := S.Monoid
	m := AlternativeMonoid[context.Context](strConcat)

	ro1 := None[context.Context, string]()
	ro2 := Of[context.Context]("Hello")
	ro3 := Of[context.Context](" World")

	result := m.Concat(m.Concat(ro1, ro2), ro3)
	expected := O.Of("Hello World")

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_MultipleSuccesses(t *testing.T) {
	// Test accumulating multiple successful values
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](1)
	ro2 := Of[context.Context](2)
	ro3 := Of[context.Context](3)
	ro4 := Of[context.Context](4)

	result := m.Concat(m.Concat(m.Concat(ro1, ro2), ro3), ro4)
	expected := O.Of(10)

	assert.Equal(t, expected, result(context.Background())())
}

func TestAlternativeMonoid_InterspersedFailures(t *testing.T) {
	// Test with failures interspersed between successes
	intAdd := N.MonoidSum[int]()
	m := AlternativeMonoid[context.Context](intAdd)

	ro1 := Of[context.Context](5)
	ro2 := None[context.Context, int]()
	ro3 := Of[context.Context](3)
	ro4 := None[context.Context, int]()
	ro5 := Of[context.Context](2)

	result := m.Concat(m.Concat(m.Concat(m.Concat(ro1, ro2), ro3), ro4), ro5)
	expected := O.Of(10) // 5 + 3 + 2

	assert.Equal(t, expected, result(context.Background())())
}
