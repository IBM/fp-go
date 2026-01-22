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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Config represents a test configuration environment
type Config struct {
	Host    string
	Port    int
	Timeout int
	Debug   bool
}

var (
	defaultConfig = Config{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
		Debug:   false,
	}

	intAddMonoid = N.MonoidSum[int]()
	strMonoid    = S.Monoid
)

// TestApplicativeMonoid tests the ApplicativeMonoid function
func TestApplicativeMonoid(t *testing.T) {
	reMonoid := ApplicativeMonoid[Config, string](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := reMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, E.Right[string](0), result)
	})

	t.Run("concat two Right values", func(t *testing.T) {
		re1 := Right[Config, string](5)
		re2 := Right[Config, string](3)
		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](8), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(reMonoid.Empty(), re)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(re, reMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat with left error", func(t *testing.T) {
		reSuccess := Right[Config, string](5)
		reFailure := Left[Config, int]("error occurred")

		combined := reMonoid.Concat(reFailure, reSuccess)
		result := combined(defaultConfig)
		assert.Equal(t, E.Left[int]("error occurred"), result)
	})

	t.Run("concat with right error", func(t *testing.T) {
		reSuccess := Right[Config, string](5)
		reFailure := Left[Config, int]("error occurred")

		combined := reMonoid.Concat(reSuccess, reFailure)
		result := combined(defaultConfig)
		assert.Equal(t, E.Left[int]("error occurred"), result)
	})

	t.Run("concat both errors", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// First error is returned
		assert.Equal(t, E.Left[int]("error1"), result)
	})

	t.Run("concat multiple values", func(t *testing.T) {
		re1 := Right[Config, string](1)
		re2 := Right[Config, string](2)
		re3 := Right[Config, string](3)
		re4 := Right[Config, string](4)

		// Chain concat calls: ((1 + 2) + 3) + 4
		combined := reMonoid.Concat(
			reMonoid.Concat(
				reMonoid.Concat(re1, re2),
				re3,
			),
			re4,
		)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](10), result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		strREMonoid := ApplicativeMonoid[Config, string](strMonoid)

		re1 := Right[Config, string]("Hello")
		re2 := Right[Config, string](" ")
		re3 := Right[Config, string]("World")

		combined := strREMonoid.Concat(
			strREMonoid.Concat(re1, re2),
			re3,
		)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string]("Hello World"), result)
	})

	t.Run("environment dependent computation", func(t *testing.T) {
		// Create computations that use the environment
		re1 := Asks[string](func(cfg Config) int {
			return cfg.Port
		})
		re2 := Right[Config, string](100)

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// defaultConfig.Port is 8080, so 8080 + 100 = 8180
		assert.Equal(t, E.Right[string](8180), result)
	})

	t.Run("environment dependent with error", func(t *testing.T) {
		re1 := MonadChain(
			Ask[Config, string](),
			func(cfg Config) ReaderEither[Config, string, int] {
				if cfg.Debug {
					return Right[Config, string](cfg.Timeout)
				}
				return Left[Config, int]("debug mode disabled")
			},
		)
		re2 := Right[Config, string](50)

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// defaultConfig.Debug is false, so should fail
		assert.Equal(t, E.Left[int]("debug mode disabled"), result)
	})
}

// TestAlternativeMonoid tests the AlternativeMonoid function
func TestAlternativeMonoid(t *testing.T) {
	reMonoid := AlternativeMonoid[Config, string](intAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := reMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, E.Right[string](0), result)
	})

	t.Run("concat two Right values", func(t *testing.T) {
		re1 := Right[Config, string](5)
		re2 := Right[Config, string](3)
		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Alternative combines successful values
		assert.Equal(t, E.Right[string](8), result)
	})

	t.Run("concat Left then Right - fallback behavior", func(t *testing.T) {
		reFailure := Left[Config, int]("error")
		reSuccess := Right[Config, string](42)

		combined := reMonoid.Concat(reFailure, reSuccess)
		result := combined(defaultConfig)
		// Should fall back to second when first fails
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat Right then Left - uses first", func(t *testing.T) {
		reSuccess := Right[Config, string](42)
		reFailure := Left[Config, int]("error")

		combined := reMonoid.Concat(reSuccess, reFailure)
		result := combined(defaultConfig)
		// Should use first successful value
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat both errors", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Second error is returned when both fail
		assert.Equal(t, E.Left[int]("error2"), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(reMonoid.Empty(), re)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(re, reMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("multiple values with some failures", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Right[Config, string](5)
		re3 := Left[Config, int]("error2")
		re4 := Right[Config, string](10)

		// Alternative should skip failures and accumulate successes
		combined := reMonoid.Concat(
			reMonoid.Concat(
				reMonoid.Concat(re1, re2),
				re3,
			),
			re4,
		)
		result := combined(defaultConfig)
		// Should accumulate successful values: 5 + 10 = 15
		assert.Equal(t, E.Right[string](15), result)
	})

	t.Run("fallback chain", func(t *testing.T) {
		// Simulate trying multiple sources until one succeeds
		primary := Left[Config, string]("primary failed")
		secondary := Left[Config, string]("secondary failed")
		tertiary := Right[Config, string]("tertiary success")

		strREMonoid := AlternativeMonoid[Config, string](strMonoid)

		// Chain concat: try primary, then secondary, then tertiary
		combined := strREMonoid.Concat(
			strREMonoid.Concat(primary, secondary),
			tertiary,
		)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string]("tertiary success"), result)
	})

	t.Run("environment dependent with fallback", func(t *testing.T) {
		// First computation fails
		re1 := Left[Config, int]("error")
		// Second computation uses environment
		re2 := Asks[string](func(cfg Config) int {
			return cfg.Timeout
		})

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Should fall back to second computation
		assert.Equal(t, E.Right[string](30), result)
	})

	t.Run("all failures in chain", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")
		re3 := Left[Config, int]("error3")

		combined := reMonoid.Concat(
			reMonoid.Concat(re1, re2),
			re3,
		)
		result := combined(defaultConfig)
		// Last error is returned
		assert.Equal(t, E.Left[int]("error3"), result)
	})
}

// TestAltMonoid tests the AltMonoid function
func TestAltMonoid(t *testing.T) {
	zero := func() ReaderEither[Config, string, int] {
		return Left[Config, int]("no value")
	}
	reMonoid := AltMonoid(zero)

	t.Run("empty element", func(t *testing.T) {
		empty := reMonoid.Empty()
		result := empty(defaultConfig)
		assert.Equal(t, E.Left[int]("no value"), result)
	})

	t.Run("concat Left then Right - uses second", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Right[Config, string](42)

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Should use first success
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat Right then Right - uses first", func(t *testing.T) {
		re1 := Right[Config, string](100)
		re2 := Right[Config, string](200)

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Uses first success, doesn't combine
		assert.Equal(t, E.Right[string](100), result)
	})

	t.Run("concat Right then Left - uses first", func(t *testing.T) {
		re1 := Right[Config, string](42)
		re2 := Left[Config, int]("error")

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat both errors", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// Second error is returned
		assert.Equal(t, E.Left[int]("error2"), result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(reMonoid.Empty(), re)
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		re := Right[Config, string](42)
		combined := reMonoid.Concat(re, reMonoid.Empty())
		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string](42), result)
	})

	t.Run("fallback chain with first success", func(t *testing.T) {
		re1 := Right[Config, string](10)
		re2 := Right[Config, string](20)
		re3 := Right[Config, string](30)

		combined := reMonoid.Concat(
			reMonoid.Concat(re1, re2),
			re3,
		)
		result := combined(defaultConfig)
		// First success is used
		assert.Equal(t, E.Right[string](10), result)
	})

	t.Run("fallback chain with middle success", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Right[Config, string](20)
		re3 := Right[Config, string](30)

		combined := reMonoid.Concat(
			reMonoid.Concat(re1, re2),
			re3,
		)
		result := combined(defaultConfig)
		// First success (re2) is used
		assert.Equal(t, E.Right[string](20), result)
	})

	t.Run("fallback chain with last success", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")
		re3 := Right[Config, string](30)

		combined := reMonoid.Concat(
			reMonoid.Concat(re1, re2),
			re3,
		)
		result := combined(defaultConfig)
		// Last success is used
		assert.Equal(t, E.Right[string](30), result)
	})

	t.Run("environment dependent fallback", func(t *testing.T) {
		re1 := MonadChain(
			Ask[Config, string](),
			func(cfg Config) ReaderEither[Config, string, int] {
				if cfg.Debug {
					return Right[Config, string](cfg.Port)
				}
				return Left[Config, int]("debug disabled")
			},
		)
		re2 := Right[Config, string](9999)

		combined := reMonoid.Concat(re1, re2)
		result := combined(defaultConfig)
		// First fails (debug is false), falls back to second
		assert.Equal(t, E.Right[string](9999), result)
	})
}

// TestMonoidLaws verifies that the monoid laws hold for ApplicativeMonoid
func TestApplicativeMonoidLaws(t *testing.T) {
	reMonoid := ApplicativeMonoid[Config, string](intAddMonoid)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(reMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(x, reMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Right[Config, string](1)
		y := Right[Config, string](2)
		z := Right[Config, string](3)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with Left values", func(t *testing.T) {
		// Verify associativity even with Left values
		x := Right[Config, string](5)
		y := Left[Config, int]("error")
		z := Right[Config, string](10)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestAlternativeMonoidLaws verifies that the monoid laws hold for AlternativeMonoid
func TestAlternativeMonoidLaws(t *testing.T) {
	reMonoid := AlternativeMonoid[Config, string](intAddMonoid)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(reMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(x, reMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Right[Config, string](1)
		y := Right[Config, string](2)
		z := Right[Config, string](3)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with Left values", func(t *testing.T) {
		// Verify associativity even with Left values
		x := Left[Config, int]("error1")
		y := Right[Config, string](5)
		z := Right[Config, string](10)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestAltMonoidLaws verifies that the monoid laws hold for AltMonoid
func TestAltMonoidLaws(t *testing.T) {
	zero := func() ReaderEither[Config, string, int] {
		return Left[Config, int]("no value")
	}
	reMonoid := AltMonoid(zero)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(reMonoid.Empty(), x)(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := Right[Config, string](42)
		result1 := reMonoid.Concat(x, reMonoid.Empty())(defaultConfig)
		result2 := x(defaultConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := Right[Config, string](1)
		y := Right[Config, string](2)
		z := Right[Config, string](3)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with Left values", func(t *testing.T) {
		// Verify associativity even with Left values
		x := Left[Config, int]("error1")
		y := Left[Config, int]("error2")
		z := Right[Config, string](10)

		left := reMonoid.Concat(reMonoid.Concat(x, y), z)(defaultConfig)
		right := reMonoid.Concat(x, reMonoid.Concat(y, z))(defaultConfig)

		assert.Equal(t, right, left)
	})
}

// TestApplicativeVsAlternative compares the behavior of both monoids
func TestApplicativeVsAlternative(t *testing.T) {
	applicativeMonoid := ApplicativeMonoid[Config, string](intAddMonoid)
	alternativeMonoid := AlternativeMonoid[Config, string](intAddMonoid)

	t.Run("both succeed - same result", func(t *testing.T) {
		re1 := Right[Config, string](5)
		re2 := Right[Config, string](3)

		appResult := applicativeMonoid.Concat(re1, re2)(defaultConfig)
		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)

		assert.Equal(t, E.Right[string](8), appResult)
		assert.Equal(t, E.Right[string](8), altResult)
		assert.Equal(t, appResult, altResult)
	})

	t.Run("first fails - different behavior", func(t *testing.T) {
		re1 := Left[Config, int]("error")
		re2 := Right[Config, string](42)

		appResult := applicativeMonoid.Concat(re1, re2)(defaultConfig)
		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)

		// Applicative fails if any fails
		assert.Equal(t, E.Left[int]("error"), appResult)
		// Alternative falls back to second
		assert.Equal(t, E.Right[string](42), altResult)
	})

	t.Run("second fails - different behavior", func(t *testing.T) {
		re1 := Right[Config, string](42)
		re2 := Left[Config, int]("error")

		appResult := applicativeMonoid.Concat(re1, re2)(defaultConfig)
		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)

		// Applicative fails if any fails
		assert.Equal(t, E.Left[int]("error"), appResult)
		// Alternative uses first success
		assert.Equal(t, E.Right[string](42), altResult)
	})

	t.Run("both fail - different behavior", func(t *testing.T) {
		re1 := Left[Config, int]("error1")
		re2 := Left[Config, int]("error2")

		appResult := applicativeMonoid.Concat(re1, re2)(defaultConfig)
		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)

		// Applicative returns first error
		assert.Equal(t, E.Left[int]("error1"), appResult)
		// Alternative returns second error
		assert.Equal(t, E.Left[int]("error2"), altResult)
	})
}

// TestAlternativeVsAlt compares AlternativeMonoid and AltMonoid
func TestAlternativeVsAlt(t *testing.T) {
	alternativeMonoid := AlternativeMonoid[Config, string](intAddMonoid)
	zero := func() ReaderEither[Config, string, int] {
		return Left[Config, int]("no value")
	}
	altMonoid := AltMonoid(zero)

	t.Run("both succeed - different behavior", func(t *testing.T) {
		re1 := Right[Config, string](5)
		re2 := Right[Config, string](3)

		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)
		altMonoidResult := altMonoid.Concat(re1, re2)(defaultConfig)

		// Alternative combines values
		assert.Equal(t, E.Right[string](8), altResult)
		// AltMonoid uses first success
		assert.Equal(t, E.Right[string](5), altMonoidResult)
	})

	t.Run("first fails - same behavior", func(t *testing.T) {
		re1 := Left[Config, int]("error")
		re2 := Right[Config, string](42)

		altResult := alternativeMonoid.Concat(re1, re2)(defaultConfig)
		altMonoidResult := altMonoid.Concat(re1, re2)(defaultConfig)

		// Both fall back to second
		assert.Equal(t, E.Right[string](42), altResult)
		assert.Equal(t, E.Right[string](42), altMonoidResult)
	})
}

// TestComplexScenarios tests more complex real-world scenarios
func TestComplexScenarios(t *testing.T) {
	t.Run("accumulate configuration values", func(t *testing.T) {
		reMonoid := ApplicativeMonoid[Config, string](intAddMonoid)

		// Accumulate multiple configuration values
		getPort := Asks[string](func(cfg Config) int {
			return cfg.Port
		})
		getTimeout := Asks[string](func(cfg Config) int {
			return cfg.Timeout
		})
		getConstant := Right[Config, string](100)

		combined := reMonoid.Concat(
			reMonoid.Concat(getPort, getTimeout),
			getConstant,
		)

		result := combined(defaultConfig)
		// 8080 + 30 + 100 = 8210
		assert.Equal(t, E.Right[string](8210), result)
	})

	t.Run("fallback configuration loading", func(t *testing.T) {
		reMonoid := AlternativeMonoid[Config, string](strMonoid)

		// Simulate trying to load config from multiple sources
		fromEnv := Left[Config, string]("env not set")
		fromFile := Left[Config, string]("file not found")
		fromDefault := Right[Config, string]("default-config")

		combined := reMonoid.Concat(
			reMonoid.Concat(fromEnv, fromFile),
			fromDefault,
		)

		result := combined(defaultConfig)
		assert.Equal(t, E.Right[string]("default-config"), result)
	})

	t.Run("partial success accumulation", func(t *testing.T) {
		reMonoid := AlternativeMonoid[Config, string](intAddMonoid)

		// Simulate collecting metrics where some may fail
		metric1 := Right[Config, string](100)
		metric2 := Left[Config, int]("metric2 failed") // Failed to collect
		metric3 := Right[Config, string](200)
		metric4 := Left[Config, int]("metric4 failed") // Failed to collect
		metric5 := Right[Config, string](300)

		combined := reMonoid.Concat(
			reMonoid.Concat(
				reMonoid.Concat(
					reMonoid.Concat(metric1, metric2),
					metric3,
				),
				metric4,
			),
			metric5,
		)

		result := combined(defaultConfig)
		// Should accumulate only successful metrics: 100 + 200 + 300 = 600
		assert.Equal(t, E.Right[string](600), result)
	})

	t.Run("cascading fallback with AltMonoid", func(t *testing.T) {
		zero := func() ReaderEither[Config, string, string] {
			return Left[Config, string]("all sources failed")
		}
		reMonoid := AltMonoid(zero)

		// Try multiple data sources in order
		primaryDB := Left[Config, string]("primary DB down")
		secondaryDB := Left[Config, string]("secondary DB down")
		cache := Right[Config, string]("cached-data")
		fallback := Right[Config, string]("fallback-data")

		combined := reMonoid.Concat(
			reMonoid.Concat(
				reMonoid.Concat(primaryDB, secondaryDB),
				cache,
			),
			fallback,
		)

		result := combined(defaultConfig)
		// Should use first successful source (cache)
		assert.Equal(t, E.Right[string]("cached-data"), result)
	})
}
