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

package readerreaderioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/context/reader"
	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type SimpleConfig struct {
	Port int
}

type DetailedConfig struct {
	Host string
	Port int
}

// TestLocalIOK tests LocalIOK functionality
func TestLocalIOK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic IO transformation", func(t *testing.T) {
		// IO effect that loads config from a path
		loadConfig := func(path string) io.IO[SimpleConfig] {
			return func() SimpleConfig {
				// Simulate loading config
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalIOK
		adapted := LocalIOK[string](loadConfig)(useConfig)
		res := adapted("config.json")(ctx)()

		assert.Equal(t, result.Of("Port: 8080"), res)
	})

	t.Run("IO transformation with side effects", func(t *testing.T) {
		var loadLog []string

		loadData := func(key string) io.IO[int] {
			return func() int {
				loadLog = append(loadLog, "Loading: "+key)
				return len(key) * 10
			}
		}

		processData := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Processed: %d", n))
				}
			}
		}

		adapted := LocalIOK[string](loadData)(processData)
		res := adapted("test")(ctx)()

		assert.Equal(t, result.Of("Processed: 40"), res)
		assert.Equal(t, []string{"Loading: test"}, loadLog)
	})

	t.Run("error propagation in ReaderReaderIOResult", func(t *testing.T) {
		loadConfig := func(path string) io.IO[SimpleConfig] {
			return func() SimpleConfig {
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderReaderIOResult that returns an error
		failingOperation := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Left[string](errors.New("operation failed"))
				}
			}
		}

		adapted := LocalIOK[string](loadConfig)(failingOperation)
		res := adapted("config.json")(ctx)()

		assert.True(t, result.IsLeft(res))
	})
}

// TestLocalIOEitherK tests LocalIOEitherK functionality
func TestLocalIOEitherK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic IOResult transformation", func(t *testing.T) {
		// IOResult effect that loads config from a path (can fail)
		loadConfig := func(path string) ioresult.IOResult[SimpleConfig] {
			return func() result.Result[SimpleConfig] {
				if path == "" {
					return result.Left[SimpleConfig](errors.New("empty path"))
				}
				return result.Of(SimpleConfig{Port: 8080})
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalIOEitherK
		adapted := LocalIOEitherK[string](loadConfig)(useConfig)

		// Success case
		res := adapted("config.json")(ctx)()
		assert.Equal(t, result.Of("Port: 8080"), res)

		// Failure case
		resErr := adapted("")(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})

	t.Run("error propagation from environment transformation", func(t *testing.T) {
		loadConfig := func(path string) ioresult.IOResult[SimpleConfig] {
			return func() result.Result[SimpleConfig] {
				return result.Left[SimpleConfig](errors.New("file not found"))
			}
		}

		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		adapted := LocalIOEitherK[string](loadConfig)(useConfig)
		res := adapted("missing.json")(ctx)()

		// Error from loadConfig should propagate
		assert.True(t, result.IsLeft(res))
	})
}

// TestLocalIOResultK tests LocalIOResultK functionality
func TestLocalIOResultK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic IOResult transformation", func(t *testing.T) {
		// IOResult effect that loads config from a path (can fail)
		loadConfig := func(path string) ioresult.IOResult[SimpleConfig] {
			return func() result.Result[SimpleConfig] {
				if path == "" {
					return result.Left[SimpleConfig](errors.New("empty path"))
				}
				return result.Of(SimpleConfig{Port: 8080})
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalIOResultK
		adapted := LocalIOResultK[string](loadConfig)(useConfig)

		// Success case
		res := adapted("config.json")(ctx)()
		assert.Equal(t, result.Of("Port: 8080"), res)

		// Failure case
		resErr := adapted("")(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})

	t.Run("compose multiple LocalIOResultK", func(t *testing.T) {
		// First transformation: string -> int (can fail)
		parseID := func(s string) ioresult.IOResult[int] {
			return func() result.Result[int] {
				if s == "" {
					return result.Left[int](errors.New("empty string"))
				}
				return result.Of(len(s) * 10)
			}
		}

		// Second transformation: int -> SimpleConfig (can fail)
		loadConfig := func(id int) ioresult.IOResult[SimpleConfig] {
			return func() result.Result[SimpleConfig] {
				if id < 0 {
					return result.Left[SimpleConfig](errors.New("invalid ID"))
				}
				return result.Of(SimpleConfig{Port: 8000 + id})
			}
		}

		// Use the config
		formatConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose transformations
		step1 := LocalIOResultK[string](loadConfig)(formatConfig)
		step2 := LocalIOResultK[string](parseID)(step1)

		// Success case
		res := step2("test")(ctx)()
		assert.Equal(t, result.Of("Port: 8040"), res)

		// Failure in first transformation
		resErr1 := step2("")(ctx)()
		assert.True(t, result.IsLeft(resErr1))
	})
}

// TestLocalReaderIOEitherK tests LocalReaderIOEitherK functionality
func TestLocalReaderIOEitherK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic ReaderIOResult transformation", func(t *testing.T) {
		// ReaderIOResult effect that loads config from a path (can fail, uses context)
		loadConfig := func(path string) readerioresult.ReaderIOResult[SimpleConfig] {
			return func(ctx context.Context) ioresult.IOResult[SimpleConfig] {
				return func() result.Result[SimpleConfig] {
					if path == "" {
						return result.Left[SimpleConfig](errors.New("empty path"))
					}
					// Could use context here for cancellation, logging, etc.
					return result.Of(SimpleConfig{Port: 8080})
				}
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalReaderIOEitherK
		adapted := LocalReaderIOEitherK[string](loadConfig)(useConfig)

		// Success case
		res := adapted("config.json")(ctx)()
		assert.Equal(t, result.Of("Port: 8080"), res)

		// Failure case
		resErr := adapted("")(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})

	t.Run("context propagation", func(t *testing.T) {
		type ctxKey string
		const key ctxKey = "test-key"

		// ReaderIOResult that reads from context
		loadFromContext := func(path string) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					if val := ctx.Value(key); val != nil {
						return result.Of(val.(string))
					}
					return result.Left[string](errors.New("key not found in context"))
				}
			}
		}

		// ReaderReaderIOResult that uses the loaded value
		useValue := func(val string) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of("Loaded: " + val)
				}
			}
		}

		adapted := LocalReaderIOEitherK[string](loadFromContext)(useValue)

		// With context value
		ctxWithValue := context.WithValue(ctx, key, "test-value")
		res := adapted("ignored")(ctxWithValue)()
		assert.Equal(t, result.Of("Loaded: test-value"), res)

		// Without context value
		resErr := adapted("ignored")(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})
}

// TestLocalReaderIOResultK tests LocalReaderIOResultK functionality
func TestLocalReaderIOResultK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic ReaderIOResult transformation", func(t *testing.T) {
		// ReaderIOResult effect that loads config from a path (can fail, uses context)
		loadConfig := func(path string) readerioresult.ReaderIOResult[SimpleConfig] {
			return func(ctx context.Context) ioresult.IOResult[SimpleConfig] {
				return func() result.Result[SimpleConfig] {
					if path == "" {
						return result.Left[SimpleConfig](errors.New("empty path"))
					}
					return result.Of(SimpleConfig{Port: 8080})
				}
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalReaderIOResultK
		adapted := LocalReaderIOResultK[string](loadConfig)(useConfig)

		// Success case
		res := adapted("config.json")(ctx)()
		assert.Equal(t, result.Of("Port: 8080"), res)

		// Failure case
		resErr := adapted("")(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})

	t.Run("real-world: load and validate config with context", func(t *testing.T) {
		type ConfigFile struct {
			Path string
		}

		// Read file with context (can fail, uses context for cancellation)
		readFile := func(cf ConfigFile) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					// Check context cancellation
					select {
					case <-ctx.Done():
						return result.Left[string](ctx.Err())
					default:
					}

					if cf.Path == "" {
						return result.Left[string](errors.New("empty path"))
					}
					return result.Of(`{"port":9000}`)
				}
			}
		}

		// Parse config with context (can fail)
		parseConfig := func(content string) readerioresult.ReaderIOResult[SimpleConfig] {
			return func(ctx context.Context) ioresult.IOResult[SimpleConfig] {
				return func() result.Result[SimpleConfig] {
					if content == "" {
						return result.Left[SimpleConfig](errors.New("empty content"))
					}
					return result.Of(SimpleConfig{Port: 9000})
				}
			}
		}

		// Use the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Using port: %d", cfg.Port))
				}
			}
		}

		// Compose the pipeline
		step1 := LocalReaderIOResultK[string](parseConfig)(useConfig)
		step2 := LocalReaderIOResultK[string](readFile)(step1)

		// Success case
		res := step2(ConfigFile{Path: "app.json"})(ctx)()
		assert.Equal(t, result.Of("Using port: 9000"), res)

		// Failure case
		resErr := step2(ConfigFile{Path: ""})(ctx)()
		assert.True(t, result.IsLeft(resErr))
	})
}

// TestLocalReaderK tests LocalReaderK functionality
func TestLocalReaderK(t *testing.T) {
	ctx := context.Background()

	t.Run("basic Reader transformation", func(t *testing.T) {
		// Reader that transforms string path to SimpleConfig using context
		loadConfig := func(path string) reader.Reader[SimpleConfig] {
			return func(ctx context.Context) SimpleConfig {
				// Could extract values from context here
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderReaderIOResult that uses the config
		useConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose using LocalReaderK
		adapted := LocalReaderK[string](loadConfig)(useConfig)
		res := adapted("config.json")(ctx)()

		assert.Equal(t, result.Of("Port: 8080"), res)
	})

	t.Run("extract config from context", func(t *testing.T) {
		type ctxKey string
		const configKey ctxKey = "config"

		// Reader that extracts config from context
		extractConfig := func(path string) reader.Reader[DetailedConfig] {
			return func(ctx context.Context) DetailedConfig {
				if cfg, ok := ctx.Value(configKey).(DetailedConfig); ok {
					return cfg
				}
				// Default config if not in context
				return DetailedConfig{Host: "localhost", Port: 8080}
			}
		}

		// Use the config
		useConfig := func(cfg DetailedConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
				}
			}
		}

		adapted := LocalReaderK[string](extractConfig)(useConfig)

		// With context value
		ctxWithConfig := context.WithValue(ctx, configKey, DetailedConfig{Host: "api.example.com", Port: 443})
		res := adapted("ignored")(ctxWithConfig)()
		assert.Equal(t, result.Of("api.example.com:443"), res)

		// Without context value (uses default)
		resDefault := adapted("ignored")(ctx)()
		assert.Equal(t, result.Of("localhost:8080"), resDefault)
	})

	t.Run("context-aware transformation", func(t *testing.T) {
		type ctxKey string
		const multiplierKey ctxKey = "multiplier"

		// Reader that uses context to compute environment
		computeValue := func(base int) reader.Reader[int] {
			return func(ctx context.Context) int {
				if mult, ok := ctx.Value(multiplierKey).(int); ok {
					return base * mult
				}
				return base
			}
		}

		// Use the computed value
		formatValue := func(val int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Value: %d", val))
				}
			}
		}

		adapted := LocalReaderK[string](computeValue)(formatValue)

		// With multiplier in context
		ctxWithMult := context.WithValue(ctx, multiplierKey, 10)
		res := adapted(5)(ctxWithMult)()
		assert.Equal(t, result.Of("Value: 50"), res)

		// Without multiplier (uses base value)
		resBase := adapted(5)(ctx)()
		assert.Equal(t, result.Of("Value: 5"), resBase)
	})

	t.Run("compose multiple LocalReaderK", func(t *testing.T) {
		type ctxKey string
		const prefixKey ctxKey = "prefix"

		// First transformation: int -> string using context
		intToString := func(n int) reader.Reader[string] {
			return func(ctx context.Context) string {
				if prefix, ok := ctx.Value(prefixKey).(string); ok {
					return fmt.Sprintf("%s-%d", prefix, n)
				}
				return fmt.Sprintf("%d", n)
			}
		}

		// Second transformation: string -> SimpleConfig
		stringToConfig := func(s string) reader.Reader[SimpleConfig] {
			return func(ctx context.Context) SimpleConfig {
				return SimpleConfig{Port: len(s) * 100}
			}
		}

		// Use the config
		formatConfig := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Port: %d", cfg.Port))
				}
			}
		}

		// Compose transformations
		step1 := LocalReaderK[string](stringToConfig)(formatConfig)
		step2 := LocalReaderK[string](intToString)(step1)

		// With prefix in context
		ctxWithPrefix := context.WithValue(ctx, prefixKey, "test")
		res := step2(42)(ctxWithPrefix)()
		// "test-42" has length 7, so port = 700
		assert.Equal(t, result.Of("Port: 700"), res)

		// Without prefix
		resNoPrefix := step2(42)(ctx)()
		// "42" has length 2, so port = 200
		assert.Equal(t, result.Of("Port: 200"), resNoPrefix)
	})

	t.Run("error propagation in ReaderReaderIOResult", func(t *testing.T) {
		// Reader transformation (pure, cannot fail)
		loadConfig := func(path string) reader.Reader[SimpleConfig] {
			return func(ctx context.Context) SimpleConfig {
				return SimpleConfig{Port: 8080}
			}
		}

		// ReaderReaderIOResult that returns an error
		failingOperation := func(cfg SimpleConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Left[string](errors.New("operation failed"))
				}
			}
		}

		adapted := LocalReaderK[string](loadConfig)(failingOperation)
		res := adapted("config.json")(ctx)()

		// Error from the ReaderReaderIOResult should propagate
		assert.True(t, result.IsLeft(res))
	})

	t.Run("real-world: environment selection based on context", func(t *testing.T) {
		type Environment string
		const (
			Dev  Environment = "dev"
			Prod Environment = "prod"
		)

		type ctxKey string
		const envKey ctxKey = "environment"

		type EnvConfig struct {
			Name string
		}

		// Reader that selects config based on context environment
		selectConfig := func(envName EnvConfig) reader.Reader[DetailedConfig] {
			return func(ctx context.Context) DetailedConfig {
				env := Dev
				if e, ok := ctx.Value(envKey).(Environment); ok {
					env = e
				}

				switch env {
				case Prod:
					return DetailedConfig{Host: "api.production.com", Port: 443}
				default:
					return DetailedConfig{Host: "localhost", Port: 8080}
				}
			}
		}

		// Use the selected config
		useConfig := func(cfg DetailedConfig) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) ioresult.IOResult[string] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Connecting to %s:%d", cfg.Host, cfg.Port))
				}
			}
		}

		adapted := LocalReaderK[string](selectConfig)(useConfig)

		// Production environment
		ctxProd := context.WithValue(ctx, envKey, Prod)
		resProd := adapted(EnvConfig{Name: "app"})(ctxProd)()
		assert.Equal(t, result.Of("Connecting to api.production.com:443"), resProd)

		// Development environment (default)
		resDev := adapted(EnvConfig{Name: "app"})(ctx)()
		assert.Equal(t, result.Of("Connecting to localhost:8080"), resDev)
	})
}
