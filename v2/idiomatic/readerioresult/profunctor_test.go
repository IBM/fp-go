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

package readerioresult

import (
	"fmt"
	"strconv"
	"testing"

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
		// ReaderIOResult that reads port from SimpleConfig
		getPort := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				return c.Port, nil
			}
		}

		// Transform DetailedConfig to SimpleConfig and int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getPort)
		result, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.NoError(t, err)
		assert.Equal(t, "8080", result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderIOResult that returns an error
		getError := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				return 0, fmt.Errorf("error occurred")
			}
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getError)
		_, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Error(t, err)
		assert.Equal(t, "error occurred", err.Error())
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		// ReaderIOResult that reads from SimpleConfig
		getPort := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				return c.Port, nil
			}
		}

		// Adapt to work with DetailedConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getPort)
		result, err := adapted(DetailedConfig{Host: "localhost", Port: 9000})()

		assert.NoError(t, err)
		assert.Equal(t, 9000, result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				return 0, fmt.Errorf("config error")
			}
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getError)
		_, err := adapted(DetailedConfig{Host: "localhost", Port: 9000})()

		assert.Error(t, err)
		assert.Equal(t, "config error", err.Error())
	})
}

// TestPromapWithIO tests Promap with actual IO effects
func TestPromapWithIO(t *testing.T) {
	t.Run("transform IO result", func(t *testing.T) {
		counter := 0

		// ReaderIOResult with side effect
		getPortWithEffect := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				counter++
				return c.Port, nil
			}
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getPortWithEffect)
		result, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.NoError(t, err)
		assert.Equal(t, "8080", result)
		assert.Equal(t, 1, counter) // Side effect occurred
	})

	t.Run("side effect occurs even on error", func(t *testing.T) {
		counter := 0

		getErrorWithEffect := func(c SimpleConfig) func() (int, error) {
			return func() (int, error) {
				counter++
				return 0, fmt.Errorf("io error")
			}
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getErrorWithEffect)
		_, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Error(t, err)
		assert.Equal(t, 1, counter) // Side effect occurred before error
	})
}

// TestPromapComposition tests that Promap can be composed
func TestPromapComposition(t *testing.T) {
	t.Run("compose two Promap transformations", func(t *testing.T) {
		type Config1 struct{ Value int }
		type Config2 struct{ Value int }
		type Config3 struct{ Value int }

		reader := func(c Config1) func() (int, error) {
			return func() (int, error) {
				return c.Value, nil
			}
		}

		f1 := func(c2 Config2) Config1 { return Config1{Value: c2.Value} }
		g1 := N.Mul(2)

		f2 := func(c3 Config3) Config2 { return Config2{Value: c3.Value} }
		g2 := N.Add(10)

		// Apply two Promap transformations
		step1 := Promap(f1, g1)(reader)
		step2 := Promap(f2, g2)(step1)

		result, err := step2(Config3{Value: 5})()

		// (5 * 2) + 10 = 20
		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})
}
