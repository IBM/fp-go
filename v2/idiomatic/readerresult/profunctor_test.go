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

	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/reader"
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
		getPort := func(c SimpleConfig) (int, error) {
			return c.Port, nil
		}

		// Transform DetailedConfig to SimpleConfig and int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getPort)
		result, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.NoError(t, err)
		assert.Equal(t, "8080", result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderResult that returns an error
		getError := func(c SimpleConfig) (int, error) {
			return 0, fmt.Errorf("error occurred")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap(simplify, toString)(getError)
		_, err := adapted(DetailedConfig{Host: "localhost", Port: 8080})

		assert.Error(t, err)
		assert.Equal(t, "error occurred", err.Error())
	})

	t.Run("environment transformation with complex types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type AppConfig struct {
			DB Database
		}

		getConnection := func(db Database) (string, error) {
			if db.ConnectionString == "" {
				return "", fmt.Errorf("empty connection string")
			}
			return db.ConnectionString, nil
		}

		extractDB := func(cfg AppConfig) Database {
			return cfg.DB
		}
		addPrefix := func(s string) string {
			return "postgres://" + s
		}

		adapted := Promap(extractDB, addPrefix)(getConnection)
		result, err := adapted(AppConfig{DB: Database{ConnectionString: "localhost:5432"}})

		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost:5432", result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		// ReaderResult that reads from SimpleConfig
		getPort := func(c SimpleConfig) (int, error) {
			return c.Port, nil
		}

		// Adapt to work with DetailedConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getPort)
		result, err := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.NoError(t, err)
		assert.Equal(t, 9000, result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(c SimpleConfig) (int, error) {
			return 0, fmt.Errorf("config error")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[int](simplify)(getError)
		_, err := adapted(DetailedConfig{Host: "localhost", Port: 9000})

		assert.Error(t, err)
		assert.Equal(t, "config error", err.Error())
	})

	t.Run("multiple field extraction", func(t *testing.T) {
		type FullConfig struct {
			Host     string
			Port     int
			Protocol string
		}

		getURL := func(c DetailedConfig) (string, error) {
			return fmt.Sprintf("%s:%d", c.Host, c.Port), nil
		}

		extractHostPort := func(fc FullConfig) DetailedConfig {
			return DetailedConfig{Host: fc.Host, Port: fc.Port}
		}

		adapted := Contramap[string](extractHostPort)(getURL)
		result, err := adapted(FullConfig{Host: "example.com", Port: 443, Protocol: "https"})

		assert.NoError(t, err)
		assert.Equal(t, "example.com:443", result)
	})
}

// TestPromapComposition tests that Promap can be composed
func TestPromapComposition(t *testing.T) {
	t.Run("compose two Promap transformations", func(t *testing.T) {
		type Config1 struct{ Value int }
		type Config2 struct{ Value int }
		type Config3 struct{ Value int }

		reader := func(c Config1) (int, error) {
			return c.Value, nil
		}

		f1 := func(c2 Config2) Config1 { return Config1{Value: c2.Value} }
		g1 := N.Mul(2)

		f2 := func(c3 Config3) Config2 { return Config2{Value: c3.Value} }
		g2 := N.Add(10)

		// Apply two Promap transformations
		step1 := Promap(f1, g1)(reader)
		step2 := Promap(f2, g2)(step1)

		result, err := step2(Config3{Value: 5})

		// (5 * 2) + 10 = 20
		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})

	t.Run("compose Promap and Contramap", func(t *testing.T) {
		type Config1 struct{ Value int }
		type Config2 struct{ Value int }

		reader := func(c Config1) (int, error) {
			return c.Value * 3, nil
		}

		// First apply Contramap
		f1 := func(c2 Config2) Config1 { return Config1{Value: c2.Value} }
		step1 := Contramap[int](f1)(reader)

		// Then apply Promap
		f2 := func(c2 Config2) Config2 { return c2 }
		g2 := func(n int) string { return fmt.Sprintf("result: %d", n) }
		step2 := Promap(f2, g2)(step1)

		result, err := step2(Config2{Value: 7})

		// 7 * 3 = 21
		assert.NoError(t, err)
		assert.Equal(t, "result: 21", result)
	})
}

// TestPromapIdentityLaws tests profunctor identity laws
func TestPromapIdentityLaws(t *testing.T) {
	t.Run("identity law", func(t *testing.T) {
		// Promap with identity functions should be identity
		reader := func(c SimpleConfig) (int, error) {
			return c.Port, nil
		}

		identity := R.Ask[SimpleConfig]()
		identityInt := R.Ask[int]()

		adapted := Promap(identity, identityInt)(reader)

		config := SimpleConfig{Port: 8080}
		result1, err1 := reader(config)
		result2, err2 := adapted(config)

		assert.Equal(t, err1, err2)
		assert.Equal(t, result1, result2)
	})
}
