// Copyright (c) 2025 IBM Corp.
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
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

type SimpleConfig struct {
	Port int
}

type DetailedConfig struct {
	Host string
	Port int
}

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both input and output", func(t *testing.T) {
		// ReaderEither that reads port from SimpleConfig
		getPort := func(c SimpleConfig) Either[string, int] {
			return E.Of[string](c.Port)
		}

		// Transform DetailedConfig to SimpleConfig and int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap[SimpleConfig, string](simplify, toString)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.Equal(t, E.Of[string]("8080"), result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderEither that returns an error
		getError := func(c SimpleConfig) Either[string, int] {
			return E.Left[int]("error occurred")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap[SimpleConfig, string](simplify, toString)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.Equal(t, E.Left[string]("error occurred"), result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		// ReaderEither that reads from SimpleConfig
		getPort := func(c SimpleConfig) Either[string, int] {
			return E.Of[string](c.Port)
		}

		// Adapt to work with DetailedConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[string, int](simplify)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.Equal(t, E.Of[string](9000), result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(c SimpleConfig) Either[string, int] {
			return E.Left[int]("config error")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[string, int](simplify)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.Equal(t, E.Left[int]("config error"), result)
	})
}

// TestPromapComposition tests that Promap can be composed
func TestPromapComposition(t *testing.T) {
	t.Run("compose two Promap transformations", func(t *testing.T) {
		type Config1 struct{ Value int }
		type Config2 struct{ Value int }
		type Config3 struct{ Value int }

		reader := func(c Config1) Either[string, int] {
			return E.Of[string](c.Value)
		}

		f1 := func(c2 Config2) Config1 { return Config1{Value: c2.Value} }
		g1 := N.Mul(2)

		f2 := func(c3 Config3) Config2 { return Config2{Value: c3.Value} }
		g2 := N.Add(10)

		// Apply two Promap transformations
		step1 := Promap[Config1, string](f1, g1)(reader)
		step2 := Promap[Config2, string](f2, g2)(step1)

		result := step2(Config3{Value: 5})

		// (5 * 2) + 10 = 20
		assert.Equal(t, E.Of[string](20), result)
	})
}
