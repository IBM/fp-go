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
