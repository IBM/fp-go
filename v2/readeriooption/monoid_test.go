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

	"github.com/IBM/fp-go/v2/lazy"
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

// TestAltMonoid tests the AltMonoid function
func TestAltMonoid(t *testing.T) {
	roMonoid := AltMonoid(lazy.Of(None[context.Context, int]()))

	t.Run("empty element", func(t *testing.T) {
		empty := roMonoid.Empty()
		result := empty(context.Background())()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat two Some values - returns first", func(t *testing.T) {
		ro1 := Of[context.Context](5)
		ro2 := Of[context.Context](10)
		combined := roMonoid.Concat(ro1, ro2)
		result := combined(context.Background())()
		// Alt returns first success, does NOT combine values
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("concat None then Some - fallback behavior", func(t *testing.T) {
		roFailure := None[context.Context, int]()
		roSuccess := Of[context.Context](42)

		combined := roMonoid.Concat(roFailure, roSuccess)
		result := combined(context.Background())()
		// Should fall back to second when first fails
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat Some then None - uses first", func(t *testing.T) {
		roSuccess := Of[context.Context](42)
		roFailure := None[context.Context, int]()

		combined := roMonoid.Concat(roSuccess, roFailure)
		result := combined(context.Background())()
		// Should use first successful value
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat both None", func(t *testing.T) {
		ro1 := None[context.Context, int]()
		ro2 := None[context.Context, int]()

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(context.Background())()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		ro := Of[context.Context](42)
		combined := roMonoid.Concat(roMonoid.Empty(), ro)
		result := combined(context.Background())()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		ro := Of[context.Context](42)
		combined := roMonoid.Concat(ro, roMonoid.Empty())
		result := combined(context.Background())()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("multiple values - returns first success", func(t *testing.T) {
		ro1 := None[context.Context, int]()
		ro2 := Of[context.Context](5)
		ro3 := Of[context.Context](10)
		ro4 := Of[context.Context](15)

		// Alt should return first success, not accumulate
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(ro1, ro2),
				ro3,
			),
			ro4,
		)
		result := combined(context.Background())()
		// Should return first successful value: 5
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("fallback chain - tries until success", func(t *testing.T) {
		// Simulate trying multiple sources until one succeeds
		primary := None[context.Context, string]()
		secondary := None[context.Context, string]()
		tertiary := Of[context.Context]("tertiary success")
		quaternary := Of[context.Context]("quaternary success")

		strROMonoid := AltMonoid(lazy.Of(None[context.Context, string]()))

		// Chain concat: try primary, then secondary, then tertiary
		combined := strROMonoid.Concat(
			strROMonoid.Concat(
				strROMonoid.Concat(primary, secondary),
				tertiary,
			),
			quaternary,
		)
		result := combined(context.Background())()
		// Should return first success (tertiary), not quaternary
		assert.Equal(t, O.Some("tertiary success"), result)
	})

	t.Run("all failures in chain", func(t *testing.T) {
		ro1 := None[context.Context, int]()
		ro2 := None[context.Context, int]()
		ro3 := None[context.Context, int]()

		combined := roMonoid.Concat(
			roMonoid.Concat(ro1, ro2),
			ro3,
		)
		result := combined(context.Background())()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("custom empty value", func(t *testing.T) {
		// Create monoid with custom empty value
		customEmpty := Of[context.Context](100)
		customMonoid := AltMonoid(lazy.Of(customEmpty))

		empty := customMonoid.Empty()
		result := empty(context.Background())()
		assert.Equal(t, O.Some(100), result)

		// Verify it acts as identity
		ro := Of[context.Context](42)
		combined := customMonoid.Concat(ro, customMonoid.Empty())
		result2 := combined(context.Background())()
		assert.Equal(t, O.Some(42), result2)
	})
}

// TestAltMonoidLaws verifies that the monoid laws hold for AltMonoid
func TestAltMonoidLaws(t *testing.T) {
	roMonoid := AltMonoid(lazy.Of(None[context.Context, int]()))

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Of[context.Context](42)
		result1 := roMonoid.Concat(roMonoid.Empty(), x)(context.Background())()
		result2 := x(context.Background())()
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Of[context.Context](42)
		result1 := roMonoid.Concat(x, roMonoid.Empty())(context.Background())()
		result2 := x(context.Background())()
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Of[context.Context](1)
		y := Of[context.Context](2)
		z := Of[context.Context](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(context.Background())()
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(context.Background())()

		assert.Equal(t, right, left)
	})

	t.Run("associativity with None values", func(t *testing.T) {
		// Verify associativity even with None values
		x := None[context.Context, int]()
		y := Of[context.Context](5)
		z := Of[context.Context](10)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(context.Background())()
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(context.Background())()

		assert.Equal(t, right, left)
	})

	t.Run("associativity with mixed None and Some", func(t *testing.T) {
		x := Of[context.Context](1)
		y := None[context.Context, int]()
		z := Of[context.Context](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(context.Background())()
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(context.Background())()

		assert.Equal(t, right, left)
	})
}

// TestAltVsAlternativeMonoid compares AltMonoid with AlternativeMonoid
func TestAltVsAlternativeMonoid(t *testing.T) {
	altMonoid := AltMonoid(lazy.Of(None[context.Context, int]()))
	alternativeMonoid := AlternativeMonoid[context.Context](N.MonoidSum[int]())

	t.Run("both succeed - different behavior", func(t *testing.T) {
		ro1 := Of[context.Context](5)
		ro2 := Of[context.Context](3)

		altResult := altMonoid.Concat(ro1, ro2)(context.Background())()
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(context.Background())()

		// Alt returns first success
		assert.Equal(t, O.Some(5), altResult)
		// Alternative combines values
		assert.Equal(t, O.Some(8), alternativeResult)
	})

	t.Run("first fails - same behavior", func(t *testing.T) {
		ro1 := None[context.Context, int]()
		ro2 := Of[context.Context](42)

		altResult := altMonoid.Concat(ro1, ro2)(context.Background())()
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(context.Background())()

		// Both fall back to second
		assert.Equal(t, O.Some(42), altResult)
		assert.Equal(t, O.Some(42), alternativeResult)
	})

	t.Run("second fails - same behavior", func(t *testing.T) {
		ro1 := Of[context.Context](42)
		ro2 := None[context.Context, int]()

		altResult := altMonoid.Concat(ro1, ro2)(context.Background())()
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(context.Background())()

		// Both use first success
		assert.Equal(t, O.Some(42), altResult)
		assert.Equal(t, O.Some(42), alternativeResult)
	})

	t.Run("both fail - same behavior", func(t *testing.T) {
		ro1 := None[context.Context, int]()
		ro2 := None[context.Context, int]()

		altResult := altMonoid.Concat(ro1, ro2)(context.Background())()
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(context.Background())()

		assert.Equal(t, O.None[int](), altResult)
		assert.Equal(t, O.None[int](), alternativeResult)
	})

	t.Run("multiple successes - key difference", func(t *testing.T) {
		ro1 := Of[context.Context](10)
		ro2 := Of[context.Context](20)
		ro3 := Of[context.Context](30)

		altResult := altMonoid.Concat(
			altMonoid.Concat(ro1, ro2),
			ro3,
		)(context.Background())()

		alternativeResult := alternativeMonoid.Concat(
			alternativeMonoid.Concat(ro1, ro2),
			ro3,
		)(context.Background())()

		// Alt returns first success
		assert.Equal(t, O.Some(10), altResult)
		// Alternative accumulates all successes
		assert.Equal(t, O.Some(60), alternativeResult)
	})
}

// TestAltMonoidWithEnvironment tests AltMonoid with environment-dependent computations
func TestAltMonoidWithEnvironment(t *testing.T) {
	type Config struct {
		UseCache bool
		Factor   int
	}

	roMonoid := AltMonoid(lazy.Of(None[Config, int]()))

	t.Run("environment dependent with fallback", func(t *testing.T) {
		// First computation depends on environment
		cacheValue := func(cfg Config) IOOption[int] {
			return func() Option[int] {
				if cfg.UseCache {
					return O.Of(100)
				}
				return O.None[int]()
			}
		}

		// Second computation also depends on environment
		dbValue := func(cfg Config) IOOption[int] {
			return func() Option[int] {
				return O.Of(50 * cfg.Factor)
			}
		}

		combined := roMonoid.Concat(cacheValue, dbValue)

		// With cache enabled, should use cache value (first success)
		cfg1 := Config{UseCache: true, Factor: 2}
		assert.Equal(t, O.Some(100), combined(cfg1)())

		// With cache disabled, should fall back to DB value
		cfg2 := Config{UseCache: false, Factor: 2}
		assert.Equal(t, O.Some(100), combined(cfg2)())
	})
}

// TestAltMonoidRealWorldScenarios tests practical use cases
func TestAltMonoidRealWorldScenarios(t *testing.T) {
	t.Run("configuration source fallback", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[context.Context, string]()))

		// Simulate trying to load config from multiple sources
		fromEnv := None[context.Context, string]()
		fromFile := None[context.Context, string]()
		fromDefault := Of[context.Context]("default-config")
		fromHardcoded := Of[context.Context]("hardcoded-config")

		// Try sources in priority order
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(fromEnv, fromFile),
				fromDefault,
			),
			fromHardcoded,
		)

		result := combined(context.Background())()
		// Should use first available (default-config)
		assert.Equal(t, O.Some("default-config"), result)
	})

	t.Run("service discovery with fallback", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[context.Context, string]()))

		// Simulate service discovery from multiple registries
		fromConsul := Of[context.Context]("consul-service")
		fromEtcd := Of[context.Context]("etcd-service")
		fromStatic := Of[context.Context]("static-service")

		combined := roMonoid.Concat(
			roMonoid.Concat(fromConsul, fromEtcd),
			fromStatic,
		)

		result := combined(context.Background())()
		// Should use first available service
		assert.Equal(t, O.Some("consul-service"), result)
	})

	t.Run("cache lookup with fallback to database", func(t *testing.T) {
		type Config struct {
			CacheEnabled bool
		}

		roMonoid := AltMonoid(lazy.Of(None[Config, int]()))

		// Simulate cache miss, then database lookup
		fromCache := func(cfg Config) IOOption[int] {
			return func() Option[int] {
				if cfg.CacheEnabled {
					return O.None[int]() // Cache miss
				}
				return O.None[int]()
			}
		}

		fromDatabase := func(cfg Config) IOOption[int] {
			return func() Option[int] {
				return O.Of(42) // Database query
			}
		}

		combined := roMonoid.Concat(fromCache, fromDatabase)
		cfg := Config{CacheEnabled: true}
		result := combined(cfg)()
		// Should fall back to database
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("retry with first success", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[context.Context, string]()))

		// Simulate retrying an operation
		attempt1 := None[context.Context, string]()
		attempt2 := None[context.Context, string]()
		attempt3 := Of[context.Context]("success on third try")
		attempt4 := Of[context.Context]("would succeed but not reached")

		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(attempt1, attempt2),
				attempt3,
			),
			attempt4,
		)

		result := combined(context.Background())()
		// Should return first success
		assert.Equal(t, O.Some("success on third try"), result)
	})
}
