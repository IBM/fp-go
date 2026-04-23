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
	"testing"

	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var (
	intAddMonoid = N.MonoidSum[int]()
	strMonoid    = S.Monoid
)

// TestApplicativeMonoid tests the ApplicativeMonoid function
func TestApplicativeMonoid(t *testing.T) {
	roMonoid := ApplicativeMonoid[Config](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := roMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, O.Some(0), result)
	})

	t.Run("concat two Some values", func(t *testing.T) {
		ro1 := Of[Config](5)
		ro2 := Of[Config](3)
		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(8), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(roMonoid.Empty(), ro)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(ro, roMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat with left None", func(t *testing.T) {
		roSuccess := Of[Config](5)
		roFailure := None[Config, int]()

		combined := roMonoid.Concat(roFailure, roSuccess)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat with right None", func(t *testing.T) {
		roSuccess := Of[Config](5)
		roFailure := None[Config, int]()

		combined := roMonoid.Concat(roSuccess, roFailure)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat both None", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat multiple values", func(t *testing.T) {
		ro1 := Of[Config](1)
		ro2 := Of[Config](2)
		ro3 := Of[Config](3)
		ro4 := Of[Config](4)

		// Chain concat calls: ((1 + 2) + 3) + 4
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(ro1, ro2),
				ro3,
			),
			ro4,
		)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		strROMonoid := ApplicativeMonoid[Config](strMonoid)

		ro1 := Of[Config]("Hello")
		ro2 := Of[Config](" ")
		ro3 := Of[Config]("World")

		combined := strROMonoid.Concat(
			strROMonoid.Concat(ro1, ro2),
			ro3,
		)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some("Hello World"), result)
	})

	t.Run("environment dependent computation", func(t *testing.T) {
		// Create computations that use the environment
		ro1 := Asks(func(cfg Config) int {
			return cfg.Port
		})
		ro2 := Of[Config](100)

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		// defaultConfig.Port is 8080, so 8080 + 100 = 8180
		assert.Equal(t, O.Some(8180), result)
	})
}

// TestAlternativeMonoid tests the AlternativeMonoid function
func TestAlternativeMonoid(t *testing.T) {
	roMonoid := AlternativeMonoid[Config](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := roMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, O.Some(0), result)
	})

	t.Run("concat two Some values", func(t *testing.T) {
		ro1 := Of[Config](5)
		ro2 := Of[Config](3)
		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		// Alternative combines successful values
		assert.Equal(t, O.Some(8), result)
	})

	t.Run("concat None then Some - fallback behavior", func(t *testing.T) {
		roFailure := None[Config, int]()
		roSuccess := Of[Config](42)

		combined := roMonoid.Concat(roFailure, roSuccess)
		result := combined(defaultConfig)
		// Should fall back to second when first fails
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat Some then None - uses first", func(t *testing.T) {
		roSuccess := Of[Config](42)
		roFailure := None[Config, int]()

		combined := roMonoid.Concat(roSuccess, roFailure)
		result := combined(defaultConfig)
		// Should use first successful value
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat both None", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(roMonoid.Empty(), ro)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(ro, roMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("multiple values with some failures", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := Of[Config](5)
		ro3 := None[Config, int]()
		ro4 := Of[Config](10)

		// Alternative should skip failures and accumulate successes
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(ro1, ro2),
				ro3,
			),
			ro4,
		)
		result := combined(defaultConfig)
		// Should accumulate successful values: 5 + 10 = 15
		assert.Equal(t, O.Some(15), result)
	})

	t.Run("fallback chain", func(t *testing.T) {
		// Simulate trying multiple sources until one succeeds
		primary := None[Config, string]()
		secondary := None[Config, string]()
		tertiary := Of[Config]("tertiary success")

		strROMonoid := AlternativeMonoid[Config](strMonoid)

		// Chain concat: try primary, then secondary, then tertiary
		combined := strROMonoid.Concat(
			strROMonoid.Concat(primary, secondary),
			tertiary,
		)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some("tertiary success"), result)
	})

	t.Run("environment dependent with fallback", func(t *testing.T) {
		// First computation fails
		ro1 := None[Config, int]()
		// Second computation uses environment
		ro2 := Asks(func(cfg Config) int {
			return cfg.Timeout
		})

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		// Should fall back to second computation
		assert.Equal(t, O.Some(30), result)
	})

	t.Run("all failures in chain", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()
		ro3 := None[Config, int]()

		combined := roMonoid.Concat(
			roMonoid.Concat(ro1, ro2),
			ro3,
		)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestMonoidLaws verifies that the monoid laws hold for ApplicativeMonoid
func TestMonoidLaws(t *testing.T) {
	roMonoid := ApplicativeMonoid[Config](intAddMonoid)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(roMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(x, roMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Of[Config](1)
		y := Of[Config](2)
		z := Of[Config](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with None values", func(t *testing.T) {
		// Verify associativity even with None values
		x := Of[Config](5)
		y := None[Config, int]()
		z := Of[Config](10)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestAlternativeMonoidLaws verifies that the monoid laws hold for AlternativeMonoid
func TestAlternativeMonoidLaws(t *testing.T) {
	roMonoid := AlternativeMonoid[Config](intAddMonoid)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(roMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(x, roMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Of[Config](1)
		y := Of[Config](2)
		z := Of[Config](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with None values", func(t *testing.T) {
		// Verify associativity even with None values
		x := None[Config, int]()
		y := Of[Config](5)
		z := Of[Config](10)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestApplicativeVsAlternative compares the behavior of both monoids
func TestApplicativeVsAlternative(t *testing.T) {
	applicativeMonoid := ApplicativeMonoid[Config](intAddMonoid)
	alternativeMonoid := AlternativeMonoid[Config](intAddMonoid)

	t.Run("both succeed - same result", func(t *testing.T) {
		ro1 := Of[Config](5)
		ro2 := Of[Config](3)

		appResult := applicativeMonoid.Concat(ro1, ro2)(defaultConfig)
		altResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		assert.Equal(t, O.Some(8), appResult)
		assert.Equal(t, O.Some(8), altResult)
		assert.Equal(t, appResult, altResult)
	})

	t.Run("first fails - different behavior", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := Of[Config](42)

		appResult := applicativeMonoid.Concat(ro1, ro2)(defaultConfig)
		altResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		// Applicative fails if any fails
		assert.Equal(t, O.None[int](), appResult)
		// Alternative falls back to second
		assert.Equal(t, O.Some(42), altResult)
	})

	t.Run("second fails - different behavior", func(t *testing.T) {
		ro1 := Of[Config](42)
		ro2 := None[Config, int]()

		appResult := applicativeMonoid.Concat(ro1, ro2)(defaultConfig)
		altResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		// Applicative fails if any fails
		assert.Equal(t, O.None[int](), appResult)
		// Alternative uses first success
		assert.Equal(t, O.Some(42), altResult)
	})

	t.Run("both fail - same result", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()

		appResult := applicativeMonoid.Concat(ro1, ro2)(defaultConfig)
		altResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		assert.Equal(t, O.None[int](), appResult)
		assert.Equal(t, O.None[int](), altResult)
		assert.Equal(t, appResult, altResult)
	})
}

// TestComplexScenarios tests more complex real-world scenarios
func TestComplexScenarios(t *testing.T) {
	t.Run("accumulate configuration values", func(t *testing.T) {
		roMonoid := ApplicativeMonoid[Config](intAddMonoid)

		// Accumulate multiple configuration values
		getPort := Asks(func(cfg Config) int { return cfg.Port })
		getTimeout := Asks(func(cfg Config) int { return cfg.Timeout })
		getConstant := Of[Config](100)

		combined := roMonoid.Concat(
			roMonoid.Concat(getPort, getTimeout),
			getConstant,
		)

		result := combined(defaultConfig)
		// 8080 + 30 + 100 = 8210
		assert.Equal(t, O.Some(8210), result)
	})

	t.Run("fallback configuration loading", func(t *testing.T) {
		roMonoid := AlternativeMonoid[Config](strMonoid)

		// Simulate trying to load config from multiple sources
		fromEnv := None[Config, string]()
		fromFile := None[Config, string]()
		fromDefault := Of[Config]("default-config")

		combined := roMonoid.Concat(
			roMonoid.Concat(fromEnv, fromFile),
			fromDefault,
		)

		result := combined(defaultConfig)
		assert.Equal(t, O.Some("default-config"), result)
	})

	t.Run("partial success accumulation", func(t *testing.T) {
		roMonoid := AlternativeMonoid[Config](intAddMonoid)

		// Simulate collecting metrics where some may fail
		metric1 := Of[Config](100)
		metric2 := None[Config, int]() // Failed to collect
		metric3 := Of[Config](200)
		metric4 := None[Config, int]() // Failed to collect
		metric5 := Of[Config](300)

		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(
					roMonoid.Concat(metric1, metric2),
					metric3,
				),
				metric4,
			),
			metric5,
		)

		result := combined(defaultConfig)
		// Should accumulate only successful metrics: 100 + 200 + 300 = 600
		assert.Equal(t, O.Some(600), result)
	})
}

// TestAltMonoid tests the AltMonoid function
func TestAltMonoid(t *testing.T) {
	roMonoid := AltMonoid(lazy.Of(None[Config, int]()))

	t.Run("empty element", func(t *testing.T) {
		empty := roMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat two Some values - returns first", func(t *testing.T) {
		ro1 := Of[Config](5)
		ro2 := Of[Config](10)
		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		// Alt returns first success, does NOT combine values
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("concat None then Some - fallback behavior", func(t *testing.T) {
		roFailure := None[Config, int]()
		roSuccess := Of[Config](42)

		combined := roMonoid.Concat(roFailure, roSuccess)
		result := combined(defaultConfig)
		// Should fall back to second when first fails
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat Some then None - uses first", func(t *testing.T) {
		roSuccess := Of[Config](42)
		roFailure := None[Config, int]()

		combined := roMonoid.Concat(roSuccess, roFailure)
		result := combined(defaultConfig)
		// Should use first successful value
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat both None", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(roMonoid.Empty(), ro)
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		ro := Of[Config](42)
		combined := roMonoid.Concat(ro, roMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("multiple values - returns first success", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := Of[Config](5)
		ro3 := Of[Config](10)
		ro4 := Of[Config](15)

		// Alt should return first success, not accumulate
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(ro1, ro2),
				ro3,
			),
			ro4,
		)
		result := combined(defaultConfig)
		// Should return first successful value: 5
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("fallback chain - tries until success", func(t *testing.T) {
		// Simulate trying multiple sources until one succeeds
		primary := None[Config, string]()
		secondary := None[Config, string]()
		tertiary := Of[Config]("tertiary success")
		quaternary := Of[Config]("quaternary success")

		strROMonoid := AltMonoid(lazy.Of(None[Config, string]()))

		// Chain concat: try primary, then secondary, then tertiary
		combined := strROMonoid.Concat(
			strROMonoid.Concat(
				strROMonoid.Concat(primary, secondary),
				tertiary,
			),
			quaternary,
		)
		result := combined(defaultConfig)
		// Should return first success (tertiary), not quaternary
		assert.Equal(t, O.Some("tertiary success"), result)
	})

	t.Run("environment dependent with fallback", func(t *testing.T) {
		// First computation fails
		ro1 := None[Config, int]()
		// Second computation uses environment
		ro2 := Asks(func(cfg Config) int {
			return cfg.Timeout
		})

		combined := roMonoid.Concat(ro1, ro2)
		result := combined(defaultConfig)
		// Should fall back to second computation
		assert.Equal(t, O.Some(30), result)
	})

	t.Run("all failures in chain", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()
		ro3 := None[Config, int]()

		combined := roMonoid.Concat(
			roMonoid.Concat(ro1, ro2),
			ro3,
		)
		result := combined(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("custom empty value", func(t *testing.T) {
		// Create monoid with custom empty value
		customEmpty := Of[Config](100)
		customMonoid := AltMonoid(lazy.Of(customEmpty))

		empty := customMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, O.Some(100), result)

		// Verify it acts as identity
		ro := Of[Config](42)
		combined := customMonoid.Concat(ro, customMonoid.Empty())
		result2 := combined(defaultConfig)
		assert.Equal(t, O.Some(42), result2)
	})
}

// TestAltMonoidLaws verifies that the monoid laws hold for AltMonoid
func TestAltMonoidLaws(t *testing.T) {
	roMonoid := AltMonoid(lazy.Of(None[Config, int]()))

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(roMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Of[Config](42)
		result1 := roMonoid.Concat(x, roMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Of[Config](1)
		y := Of[Config](2)
		z := Of[Config](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with None values", func(t *testing.T) {
		// Verify associativity even with None values
		x := None[Config, int]()
		y := Of[Config](5)
		z := Of[Config](10)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with mixed None and Some", func(t *testing.T) {
		x := Of[Config](1)
		y := None[Config, int]()
		z := Of[Config](3)

		left := roMonoid.Concat(roMonoid.Concat(x, y), z)(defaultConfig)
		right := roMonoid.Concat(x, roMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestAltVsAlternativeMonoid compares AltMonoid with AlternativeMonoid
func TestAltVsAlternativeMonoid(t *testing.T) {
	altMonoid := AltMonoid(lazy.Of(None[Config, int]()))
	alternativeMonoid := AlternativeMonoid[Config](intAddMonoid)

	t.Run("both succeed - different behavior", func(t *testing.T) {
		ro1 := Of[Config](5)
		ro2 := Of[Config](3)

		altResult := altMonoid.Concat(ro1, ro2)(defaultConfig)
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		// Alt returns first success
		assert.Equal(t, O.Some(5), altResult)
		// Alternative combines values
		assert.Equal(t, O.Some(8), alternativeResult)
	})

	t.Run("first fails - same behavior", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := Of[Config](42)

		altResult := altMonoid.Concat(ro1, ro2)(defaultConfig)
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		// Both fall back to second
		assert.Equal(t, O.Some(42), altResult)
		assert.Equal(t, O.Some(42), alternativeResult)
	})

	t.Run("second fails - same behavior", func(t *testing.T) {
		ro1 := Of[Config](42)
		ro2 := None[Config, int]()

		altResult := altMonoid.Concat(ro1, ro2)(defaultConfig)
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		// Both use first success
		assert.Equal(t, O.Some(42), altResult)
		assert.Equal(t, O.Some(42), alternativeResult)
	})

	t.Run("both fail - same behavior", func(t *testing.T) {
		ro1 := None[Config, int]()
		ro2 := None[Config, int]()

		altResult := altMonoid.Concat(ro1, ro2)(defaultConfig)
		alternativeResult := alternativeMonoid.Concat(ro1, ro2)(defaultConfig)

		assert.Equal(t, O.None[int](), altResult)
		assert.Equal(t, O.None[int](), alternativeResult)
	})

	t.Run("multiple successes - key difference", func(t *testing.T) {
		ro1 := Of[Config](10)
		ro2 := Of[Config](20)
		ro3 := Of[Config](30)

		altResult := altMonoid.Concat(
			altMonoid.Concat(ro1, ro2),
			ro3,
		)(defaultConfig)

		alternativeResult := alternativeMonoid.Concat(
			alternativeMonoid.Concat(ro1, ro2),
			ro3,
		)(defaultConfig)

		// Alt returns first success
		assert.Equal(t, O.Some(10), altResult)
		// Alternative accumulates all successes
		assert.Equal(t, O.Some(60), alternativeResult)
	})
}

// TestAltMonoidRealWorldScenarios tests practical use cases
func TestAltMonoidRealWorldScenarios(t *testing.T) {
	t.Run("configuration source fallback", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[Config, string]()))

		// Simulate trying to load config from multiple sources
		fromEnv := None[Config, string]()
		fromFile := None[Config, string]()
		fromDefault := Of[Config]("default-config")
		fromHardcoded := Of[Config]("hardcoded-config")

		// Try sources in priority order
		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(fromEnv, fromFile),
				fromDefault,
			),
			fromHardcoded,
		)

		result := combined(defaultConfig)
		// Should use first available (default-config)
		assert.Equal(t, O.Some("default-config"), result)
	})

	t.Run("service discovery with fallback", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[Config, string]()))

		// Simulate service discovery from multiple registries
		fromConsul := Of[Config]("consul-service")
		fromEtcd := Of[Config]("etcd-service")
		fromStatic := Of[Config]("static-service")

		combined := roMonoid.Concat(
			roMonoid.Concat(fromConsul, fromEtcd),
			fromStatic,
		)

		result := combined(defaultConfig)
		// Should use first available service
		assert.Equal(t, O.Some("consul-service"), result)
	})

	t.Run("cache lookup with fallback to database", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[Config, int]()))

		// Simulate cache miss, then database lookup
		fromCache := None[Config, int]()
		fromDatabase := Asks(func(cfg Config) int {
			return cfg.Port // Simulate database query
		})

		combined := roMonoid.Concat(fromCache, fromDatabase)
		result := combined(defaultConfig)
		// Should fall back to database
		assert.Equal(t, O.Some(8080), result)
	})

	t.Run("retry with first success", func(t *testing.T) {
		roMonoid := AltMonoid(lazy.Of(None[Config, string]()))

		// Simulate retrying an operation
		attempt1 := None[Config, string]()
		attempt2 := None[Config, string]()
		attempt3 := Of[Config]("success on third try")
		attempt4 := Of[Config]("would succeed but not reached")

		combined := roMonoid.Concat(
			roMonoid.Concat(
				roMonoid.Concat(attempt1, attempt2),
				attempt3,
			),
			attempt4,
		)

		result := combined(defaultConfig)
		// Should return first success
		assert.Equal(t, O.Some("success on third try"), result)
	})
}
