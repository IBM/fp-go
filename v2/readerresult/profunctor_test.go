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

package readerresult

import (
	"fmt"
	"strconv"
	"testing"

	R "github.com/IBM/fp-go/v2/result"
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
		// ReaderResult that reads port from SimpleConfig
		getPort := func(c SimpleConfig) Result[int] {
			return R.Of(c.Port)
		}

		// Transform DetailedConfig to SimpleConfig and int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.Equal(t, R.Of("8080"), result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderResult that returns an error
		getError := func(c SimpleConfig) Result[int] {
			return R.Left[int](fmt.Errorf("error occurred"))
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.True(t, R.IsLeft(result))
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		// ReaderResult that reads from SimpleConfig
		getPort := func(c SimpleConfig) Result[int] {
			return R.Of(c.Port)
		}

		// Adapt to work with DetailedConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.Equal(t, R.Of(9000), result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(c SimpleConfig) Result[int] {
			return R.Left[int](fmt.Errorf("config error"))
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.True(t, R.IsLeft(result))
	})
}
