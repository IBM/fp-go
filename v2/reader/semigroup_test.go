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

package reader

import (
	"strconv"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// SemigroupConfig represents a test configuration environment for semigroup tests
type SemigroupConfig struct {
	Host       string
	Port       int
	Timeout    int
	MaxRetries int
	Debug      bool
}

var (
	defaultSemigroupConfig = SemigroupConfig{
		Host:       "localhost",
		Port:       8080,
		Timeout:    30,
		MaxRetries: 3,
		Debug:      false,
	}

	semigroupIntAddMonoid = N.MonoidSum[int]()
	semigroupIntMulMonoid = N.MonoidProduct[int]()
	semigroupStrMonoid    = S.Monoid
)

// TestApplicativeMonoidSemigroup tests the ApplicativeMonoid function
func TestApplicativeMonoidSemigroup(t *testing.T) {
	readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

	t.Run("empty element", func(t *testing.T) {
		empty := readerMonoid.Empty()
		result := empty(defaultSemigroupConfig)
		assert.Equal(t, 0, result)
	})

	t.Run("concat two readers", func(t *testing.T) {
		r1 := func(c SemigroupConfig) int { return c.Port }
		r2 := func(c SemigroupConfig) int { return c.Timeout }
		combined := readerMonoid.Concat(r1, r2)
		result := combined(defaultSemigroupConfig)
		// 8080 + 30 = 8110
		assert.Equal(t, 8110, result)
	})

	t.Run("concat with empty - left identity", func(t *testing.T) {
		r := func(c SemigroupConfig) int { return c.Port }
		combined := readerMonoid.Concat(readerMonoid.Empty(), r)
		result := combined(defaultSemigroupConfig)
		assert.Equal(t, 8080, result)
	})

	t.Run("concat with empty - right identity", func(t *testing.T) {
		r := func(c SemigroupConfig) int { return c.Port }
		combined := readerMonoid.Concat(r, readerMonoid.Empty())
		result := combined(defaultSemigroupConfig)
		assert.Equal(t, 8080, result)
	})

	t.Run("concat multiple readers", func(t *testing.T) {
		r1 := func(c SemigroupConfig) int { return c.Port }
		r2 := func(c SemigroupConfig) int { return c.Timeout }
		r3 := func(c SemigroupConfig) int { return c.MaxRetries }
		r4 := Of[SemigroupConfig](100)

		// Chain concat calls: ((r1 + r2) + r3) + r4
		combined := readerMonoid.Concat(
			readerMonoid.Concat(
				readerMonoid.Concat(r1, r2),
				r3,
			),
			r4,
		)
		result := combined(defaultSemigroupConfig)
		// 8080 + 30 + 3 + 100 = 8213
		assert.Equal(t, 8213, result)
	})

	t.Run("concat constant readers", func(t *testing.T) {
		r1 := Of[SemigroupConfig](10)
		r2 := Of[SemigroupConfig](20)
		r3 := Of[SemigroupConfig](30)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(r1, r2),
			r3,
		)
		result := combined(defaultSemigroupConfig)
		assert.Equal(t, 60, result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		strReaderMonoid := ApplicativeMonoid[SemigroupConfig](semigroupStrMonoid)

		r1 := func(c SemigroupConfig) string { return c.Host }
		r2 := Of[SemigroupConfig](":")
		r3 := Asks(func(c SemigroupConfig) string {
			return strconv.Itoa(c.Port)
		})

		combined := strReaderMonoid.Concat(
			strReaderMonoid.Concat(r1, r2),
			r3,
		)
		result := combined(defaultSemigroupConfig)
		assert.Equal(t, "localhost:8080", result)
	})

	t.Run("multiplication monoid", func(t *testing.T) {
		mulReaderMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntMulMonoid)

		r1 := func(c SemigroupConfig) int { return c.MaxRetries }
		r2 := Of[SemigroupConfig](10)
		r3 := Of[SemigroupConfig](2)

		combined := mulReaderMonoid.Concat(
			mulReaderMonoid.Concat(r1, r2),
			r3,
		)
		result := combined(defaultSemigroupConfig)
		// 3 * 10 * 2 = 60
		assert.Equal(t, 60, result)
	})

	t.Run("environment dependent computation", func(t *testing.T) {
		// Create readers that use different parts of the environment
		getPort := Asks(func(c SemigroupConfig) int { return c.Port })
		getTimeout := Asks(func(c SemigroupConfig) int { return c.Timeout })
		getRetries := Asks(func(c SemigroupConfig) int { return c.MaxRetries })

		combined := readerMonoid.Concat(
			readerMonoid.Concat(getPort, getTimeout),
			getRetries,
		)

		result := combined(defaultSemigroupConfig)
		// 8080 + 30 + 3 = 8113
		assert.Equal(t, 8113, result)
	})

	t.Run("mixed constant and environment readers", func(t *testing.T) {
		r1 := Of[SemigroupConfig](1000)
		r2 := func(c SemigroupConfig) int { return c.Port }
		r3 := Of[SemigroupConfig](5)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(r1, r2),
			r3,
		)
		result := combined(defaultSemigroupConfig)
		// 1000 + 8080 + 5 = 9085
		assert.Equal(t, 9085, result)
	})

	t.Run("different environment values", func(t *testing.T) {
		r1 := func(c SemigroupConfig) int { return c.Port }
		r2 := func(c SemigroupConfig) int { return c.Timeout }

		combined := readerMonoid.Concat(r1, r2)

		// Test with different configs
		config1 := SemigroupConfig{Port: 3000, Timeout: 60}
		config2 := SemigroupConfig{Port: 9000, Timeout: 120}

		result1 := combined(config1)
		result2 := combined(config2)

		assert.Equal(t, 3060, result1)
		assert.Equal(t, 9120, result2)
	})

	t.Run("conditional reader based on environment", func(t *testing.T) {
		r1 := func(c SemigroupConfig) int {
			if c.Debug {
				return c.Port * 2
			}
			return c.Port
		}
		r2 := func(c SemigroupConfig) int { return c.Timeout }

		combined := readerMonoid.Concat(r1, r2)

		// Test with debug off
		result1 := combined(defaultSemigroupConfig)
		assert.Equal(t, 8110, result1) // 8080 + 30

		// Test with debug on
		debugConfig := defaultSemigroupConfig
		debugConfig.Debug = true
		result2 := combined(debugConfig)
		assert.Equal(t, 16190, result2) // (8080 * 2) + 30
	})
}

// TestMonoidLawsSemigroup verifies that the monoid laws hold for ApplicativeMonoid
func TestMonoidLawsSemigroup(t *testing.T) {
	readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

	t.Run("left identity law", func(t *testing.T) {
		// empty <> x == x
		x := func(c SemigroupConfig) int { return c.Port }
		result1 := readerMonoid.Concat(readerMonoid.Empty(), x)(defaultSemigroupConfig)
		result2 := x(defaultSemigroupConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("right identity law", func(t *testing.T) {
		// x <> empty == x
		x := func(c SemigroupConfig) int { return c.Port }
		result1 := readerMonoid.Concat(x, readerMonoid.Empty())(defaultSemigroupConfig)
		result2 := x(defaultSemigroupConfig)
		assert.Equal(t, result2, result1)
	})

	t.Run("associativity law", func(t *testing.T) {
		// (x <> y) <> z == x <> (y <> z)
		x := func(c SemigroupConfig) int { return c.Port }
		y := func(c SemigroupConfig) int { return c.Timeout }
		z := func(c SemigroupConfig) int { return c.MaxRetries }

		left := readerMonoid.Concat(readerMonoid.Concat(x, y), z)(defaultSemigroupConfig)
		right := readerMonoid.Concat(x, readerMonoid.Concat(y, z))(defaultSemigroupConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with constants", func(t *testing.T) {
		x := Of[SemigroupConfig](10)
		y := Of[SemigroupConfig](20)
		z := Of[SemigroupConfig](30)

		left := readerMonoid.Concat(readerMonoid.Concat(x, y), z)(defaultSemigroupConfig)
		right := readerMonoid.Concat(x, readerMonoid.Concat(y, z))(defaultSemigroupConfig)

		assert.Equal(t, right, left)
	})

	t.Run("associativity with mixed readers", func(t *testing.T) {
		x := func(c SemigroupConfig) int { return c.Port }
		y := Of[SemigroupConfig](100)
		z := func(c SemigroupConfig) int { return c.Timeout }

		left := readerMonoid.Concat(readerMonoid.Concat(x, y), z)(defaultSemigroupConfig)
		right := readerMonoid.Concat(x, readerMonoid.Concat(y, z))(defaultSemigroupConfig)

		assert.Equal(t, right, left)
	})
}

// TestMonoidWithDifferentTypesSemigroup tests monoid with various types
func TestMonoidWithDifferentTypesSemigroup(t *testing.T) {
	t.Run("string monoid", func(t *testing.T) {
		strReaderMonoid := ApplicativeMonoid[SemigroupConfig](semigroupStrMonoid)

		r1 := func(c SemigroupConfig) string { return "Host: " }
		r2 := func(c SemigroupConfig) string { return c.Host }
		r3 := Of[SemigroupConfig](" | Port: ")
		r4 := Asks(func(c SemigroupConfig) string { return strconv.Itoa(c.Port) })

		combined := strReaderMonoid.Concat(
			strReaderMonoid.Concat(
				strReaderMonoid.Concat(r1, r2),
				r3,
			),
			r4,
		)

		result := combined(defaultSemigroupConfig)
		assert.Equal(t, "Host: localhost | Port: 8080", result)
	})

	t.Run("product monoid", func(t *testing.T) {
		mulReaderMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntMulMonoid)

		r1 := func(c SemigroupConfig) int { return 2 }
		r2 := func(c SemigroupConfig) int { return c.MaxRetries }
		r3 := Of[SemigroupConfig](5)

		combined := mulReaderMonoid.Concat(
			mulReaderMonoid.Concat(r1, r2),
			r3,
		)

		result := combined(defaultSemigroupConfig)
		// 2 * 3 * 5 = 30
		assert.Equal(t, 30, result)
	})
}

// TestComplexScenariosSemigroup tests more complex real-world scenarios
func TestComplexScenariosSemigroup(t *testing.T) {
	t.Run("accumulate configuration values", func(t *testing.T) {
		readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

		// Accumulate multiple configuration values
		getPort := Asks(func(c SemigroupConfig) int { return c.Port })
		getTimeout := Asks(func(c SemigroupConfig) int { return c.Timeout })
		getRetries := Asks(func(c SemigroupConfig) int { return c.MaxRetries })
		getConstant := Of[SemigroupConfig](1000)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(
				readerMonoid.Concat(getPort, getTimeout),
				getRetries,
			),
			getConstant,
		)

		result := combined(defaultSemigroupConfig)
		// 8080 + 30 + 3 + 1000 = 9113
		assert.Equal(t, 9113, result)
	})

	t.Run("build connection string", func(t *testing.T) {
		strReaderMonoid := ApplicativeMonoid[SemigroupConfig](semigroupStrMonoid)

		protocol := Of[SemigroupConfig]("http://")
		host := func(c SemigroupConfig) string { return c.Host }
		colon := Of[SemigroupConfig](":")
		port := Asks(func(c SemigroupConfig) string { return strconv.Itoa(c.Port) })

		buildURL := strReaderMonoid.Concat(
			strReaderMonoid.Concat(
				strReaderMonoid.Concat(protocol, host),
				colon,
			),
			port,
		)

		result := buildURL(defaultSemigroupConfig)
		assert.Equal(t, "http://localhost:8080", result)
	})

	t.Run("calculate total score", func(t *testing.T) {
		type ScoreConfig struct {
			BaseScore        int
			BonusPoints      int
			Multiplier       int
			PenaltyDeduction int
		}

		scoreConfig := ScoreConfig{
			BaseScore:        100,
			BonusPoints:      50,
			Multiplier:       2,
			PenaltyDeduction: 10,
		}

		readerMonoid := ApplicativeMonoid[ScoreConfig](semigroupIntAddMonoid)

		getBase := func(c ScoreConfig) int { return c.BaseScore }
		getBonus := func(c ScoreConfig) int { return c.BonusPoints }
		getPenalty := func(c ScoreConfig) int { return -c.PenaltyDeduction }

		totalScore := readerMonoid.Concat(
			readerMonoid.Concat(getBase, getBonus),
			getPenalty,
		)

		result := totalScore(scoreConfig)
		// 100 + 50 - 10 = 140
		assert.Equal(t, 140, result)
	})

	t.Run("compose multiple readers with empty", func(t *testing.T) {
		readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

		r1 := func(c SemigroupConfig) int { return c.Port }
		r2 := readerMonoid.Empty()
		r3 := func(c SemigroupConfig) int { return c.Timeout }
		r4 := readerMonoid.Empty()
		r5 := Of[SemigroupConfig](100)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(
				readerMonoid.Concat(
					readerMonoid.Concat(r1, r2),
					r3,
				),
				r4,
			),
			r5,
		)

		result := combined(defaultSemigroupConfig)
		// 8080 + 0 + 30 + 0 + 100 = 8210
		assert.Equal(t, 8210, result)
	})
}

// TestEdgeCasesSemigroup tests edge cases and boundary conditions
func TestEdgeCasesSemigroup(t *testing.T) {
	t.Run("empty config struct", func(t *testing.T) {
		type EmptyConfig struct{}
		emptyConfig := EmptyConfig{}

		readerMonoid := ApplicativeMonoid[EmptyConfig](semigroupIntAddMonoid)

		r1 := Of[EmptyConfig](10)
		r2 := Of[EmptyConfig](20)

		combined := readerMonoid.Concat(r1, r2)
		result := combined(emptyConfig)
		assert.Equal(t, 30, result)
	})

	t.Run("zero values", func(t *testing.T) {
		readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

		r1 := Of[SemigroupConfig](0)
		r2 := Of[SemigroupConfig](0)
		r3 := Of[SemigroupConfig](0)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(r1, r2),
			r3,
		)

		result := combined(defaultSemigroupConfig)
		assert.Equal(t, 0, result)
	})

	t.Run("negative values", func(t *testing.T) {
		readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

		r1 := Of[SemigroupConfig](-100)
		r2 := Of[SemigroupConfig](50)
		r3 := Of[SemigroupConfig](-30)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(r1, r2),
			r3,
		)

		result := combined(defaultSemigroupConfig)
		// -100 + 50 - 30 = -80
		assert.Equal(t, -80, result)
	})

	t.Run("large values", func(t *testing.T) {
		readerMonoid := ApplicativeMonoid[SemigroupConfig](semigroupIntAddMonoid)

		r1 := Of[SemigroupConfig](1000000)
		r2 := Of[SemigroupConfig](2000000)
		r3 := Of[SemigroupConfig](3000000)

		combined := readerMonoid.Concat(
			readerMonoid.Concat(r1, r2),
			r3,
		)

		result := combined(defaultSemigroupConfig)
		assert.Equal(t, 6000000, result)
	})
}
