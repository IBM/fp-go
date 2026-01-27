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

package readerioeither

import (
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
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
		// ReaderIOEither that reads port from SimpleConfig
		getPort := func(c SimpleConfig) IOEither[string, int] {
			return IOE.Of[string](c.Port)
		}

		// Transform DetailedConfig to SimpleConfig and int to string
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap[SimpleConfig, string](simplify, toString)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Equal(t, E.Of[string]("8080"), result)
	})

	t.Run("handles error case", func(t *testing.T) {
		// ReaderIOEither that returns an error
		getError := func(c SimpleConfig) IOEither[string, int] {
			return IOE.Left[int]("error occurred")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap[SimpleConfig, string](simplify, toString)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Equal(t, E.Left[string]("error occurred"), result)
	})
}

// TestContramapBasic tests basic Contramap functionality
func TestContramapBasic(t *testing.T) {
	t.Run("environment adaptation", func(t *testing.T) {
		// ReaderIOEither that reads from SimpleConfig
		getPort := func(c SimpleConfig) IOEither[string, int] {
			return IOE.Of[string](c.Port)
		}

		// Adapt to work with DetailedConfig
		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[string, int](simplify)(getPort)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})()

		assert.Equal(t, E.Of[string](9000), result)
	})

	t.Run("preserves error", func(t *testing.T) {
		getError := func(c SimpleConfig) IOEither[string, int] {
			return IOE.Left[int]("config error")
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}

		adapted := Contramap[string, int](simplify)(getError)
		result := adapted(DetailedConfig{Host: "localhost", Port: 9000})()

		assert.Equal(t, E.Left[int]("config error"), result)
	})
}

// TestPromapWithIO tests Promap with actual IO effects
func TestPromapWithIO(t *testing.T) {
	t.Run("transform IO result", func(t *testing.T) {
		counter := 0

		// ReaderIOEither with side effect
		getPortWithEffect := func(c SimpleConfig) IOEither[string, int] {
			return func() E.Either[string, int] {
				counter++
				return E.Of[string](c.Port)
			}
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Port: d.Port}
		}
		toString := strconv.Itoa

		adapted := Promap[SimpleConfig, string](simplify, toString)(getPortWithEffect)
		result := adapted(DetailedConfig{Host: "localhost", Port: 8080})()

		assert.Equal(t, E.Of[string]("8080"), result)
		assert.Equal(t, 1, counter) // Side effect occurred
	})
}

// TestLocalIOK tests LocalIOK functionality
func TestLocalIOK(t *testing.T) {
	t.Run("basic IO transformation", func(t *testing.T) {
		// IO effect that loads config from a path
		loadConfig := func(path string) IOE.IO[SimpleConfig] {
			return func() SimpleConfig {
				// Simulate loading config
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderIOEither that uses the config
		useConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Port: " + strconv.Itoa(cfg.Port))
		}

		// Compose using LocalIOK
		adapted := LocalIOK[string, string, SimpleConfig, string](loadConfig)(useConfig)
		result := adapted("config.json")()

		assert.Equal(t, E.Of[string]("Port: 8080"), result)
	})

	t.Run("IO transformation with side effects", func(t *testing.T) {
		var loadLog []string

		loadData := func(key string) IOE.IO[int] {
			return func() int {
				loadLog = append(loadLog, "Loading: "+key)
				return len(key) * 10
			}
		}

		processData := func(n int) IOEither[string, string] {
			return IOE.Of[string]("Processed: " + strconv.Itoa(n))
		}

		adapted := LocalIOK[string, string, int, string](loadData)(processData)
		result := adapted("test")()

		assert.Equal(t, E.Of[string]("Processed: 40"), result)
		assert.Equal(t, []string{"Loading: test"}, loadLog)
	})

	t.Run("error propagation in ReaderIOEither", func(t *testing.T) {
		loadConfig := func(path string) IOE.IO[SimpleConfig] {
			return func() SimpleConfig {
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderIOEither that returns an error
		failingOperation := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Left[string]("operation failed")
		}

		adapted := LocalIOK[string, string, SimpleConfig, string](loadConfig)(failingOperation)
		result := adapted("config.json")()

		assert.Equal(t, E.Left[string]("operation failed"), result)
	})

	t.Run("compose multiple LocalIOK", func(t *testing.T) {
		// First transformation: string -> int
		parseID := func(s string) IOE.IO[int] {
			return func() int {
				id, _ := strconv.Atoi(s)
				return id
			}
		}

		// Second transformation: int -> SimpleConfig
		loadConfig := func(id int) IOE.IO[SimpleConfig] {
			return func() SimpleConfig {
				return SimpleConfig{Port: 8000 + id}
			}
		}

		// Use the config
		formatConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Port: " + strconv.Itoa(cfg.Port))
		}

		// Compose transformations
		step1 := LocalIOK[string, string, SimpleConfig, int](loadConfig)(formatConfig)
		step2 := LocalIOK[string, string, int, string](parseID)(step1)

		result := step2("42")()
		assert.Equal(t, E.Of[string]("Port: 8042"), result)
	})
}

// TestLocalIOEitherK tests LocalIOEitherK functionality
func TestLocalIOEitherK(t *testing.T) {
	t.Run("basic IOEither transformation", func(t *testing.T) {
		// IOEither effect that loads config from a path (can fail)
		loadConfig := func(path string) IOEither[string, SimpleConfig] {
			return func() E.Either[string, SimpleConfig] {
				if path == "" {
					return E.Left[SimpleConfig]("empty path")
				}
				return E.Of[string](SimpleConfig{Port: 8080})
			}
		}

		// ReaderIOEither that uses the config
		useConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Port: " + strconv.Itoa(cfg.Port))
		}

		// Compose using LocalIOEitherK
		adapted := LocalIOEitherK[string, SimpleConfig, string, string](loadConfig)(useConfig)

		// Success case
		result := adapted("config.json")()
		assert.Equal(t, E.Of[string]("Port: 8080"), result)

		// Failure case
		resultErr := adapted("")()
		assert.Equal(t, E.Left[string]("empty path"), resultErr)
	})

	t.Run("error propagation from environment transformation", func(t *testing.T) {
		loadConfig := func(path string) IOEither[string, SimpleConfig] {
			return func() E.Either[string, SimpleConfig] {
				return E.Left[SimpleConfig]("file not found")
			}
		}

		useConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Port: " + strconv.Itoa(cfg.Port))
		}

		adapted := LocalIOEitherK[string, SimpleConfig, string, string](loadConfig)(useConfig)
		result := adapted("missing.json")()

		// Error from loadConfig should propagate
		assert.Equal(t, E.Left[string]("file not found"), result)
	})

	t.Run("error propagation from ReaderIOEither", func(t *testing.T) {
		loadConfig := func(path string) IOEither[string, SimpleConfig] {
			return IOE.Of[string](SimpleConfig{Port: 8080})
		}

		// ReaderIOEither that returns an error
		failingOperation := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Left[string]("operation failed")
		}

		adapted := LocalIOEitherK[string, SimpleConfig, string, string](loadConfig)(failingOperation)
		result := adapted("config.json")()

		// Error from ReaderIOEither should propagate
		assert.Equal(t, E.Left[string]("operation failed"), result)
	})

	t.Run("compose multiple LocalIOEitherK", func(t *testing.T) {
		// First transformation: string -> int (can fail)
		parseID := func(s string) IOEither[string, int] {
			return func() E.Either[string, int] {
				id, err := strconv.Atoi(s)
				if err != nil {
					return E.Left[int]("invalid ID")
				}
				return E.Of[string](id)
			}
		}

		// Second transformation: int -> SimpleConfig (can fail)
		loadConfig := func(id int) IOEither[string, SimpleConfig] {
			return func() E.Either[string, SimpleConfig] {
				if id < 0 {
					return E.Left[SimpleConfig]("invalid ID")
				}
				return E.Of[string](SimpleConfig{Port: 8000 + id})
			}
		}

		// Use the config
		formatConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Port: " + strconv.Itoa(cfg.Port))
		}

		// Compose transformations
		step1 := LocalIOEitherK[string, SimpleConfig, int, string](loadConfig)(formatConfig)
		step2 := LocalIOEitherK[string, int, string, string](parseID)(step1)

		// Success case
		result := step2("42")()
		assert.Equal(t, E.Of[string]("Port: 8042"), result)

		// Failure in first transformation
		resultErr1 := step2("invalid")()
		assert.Equal(t, E.Left[string]("invalid ID"), resultErr1)

		// Failure in second transformation
		resultErr2 := step2("-5")()
		assert.Equal(t, E.Left[string]("invalid ID"), resultErr2)
	})

	t.Run("real-world: load and validate config", func(t *testing.T) {
		type ConfigFile struct {
			Path string
		}

		// Read file (can fail)
		readFile := func(cf ConfigFile) IOEither[string, string] {
			return func() E.Either[string, string] {
				if cf.Path == "" {
					return E.Left[string]("empty path")
				}
				return E.Of[string](`{"port":9000}`)
			}
		}

		// Parse config (can fail)
		parseConfig := func(content string) IOEither[string, SimpleConfig] {
			return func() E.Either[string, SimpleConfig] {
				if content == "" {
					return E.Left[SimpleConfig]("empty content")
				}
				return E.Of[string](SimpleConfig{Port: 9000})
			}
		}

		// Use the config
		useConfig := func(cfg SimpleConfig) IOEither[string, string] {
			return IOE.Of[string]("Using port: " + strconv.Itoa(cfg.Port))
		}

		// Compose the pipeline
		step1 := LocalIOEitherK[string, SimpleConfig, string, string](parseConfig)(useConfig)
		step2 := LocalIOEitherK[string, string, ConfigFile, string](readFile)(step1)

		// Success case
		result := step2(ConfigFile{Path: "app.json"})()
		assert.Equal(t, E.Of[string]("Using port: 9000"), result)

		// Failure case
		resultErr := step2(ConfigFile{Path: ""})()
		assert.Equal(t, E.Left[string]("empty path"), resultErr)
	})
}
