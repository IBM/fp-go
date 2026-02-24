// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
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

package effect

import (
	"context"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Multiplier  int
	Prefix      string
	DatabaseURL string
}

var testConfig = TestConfig{
	Multiplier:  3,
	Prefix:      "LOG",
	DatabaseURL: "postgres://localhost",
}

func TestChainReaderK_Success(t *testing.T) {
	t.Run("chains with Reader that uses context", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](5),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n * cfg.Multiplier
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(15), outcome)
	})

	t.Run("chains multiple Reader operations", func(t *testing.T) {
		computation := F.Pipe2(
			Of[TestConfig](10),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n + cfg.Multiplier
				}
			}),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n * cfg.Multiplier
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(39), outcome)
	})

	t.Run("chains with Reader that returns string", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainReaderK(func(n int) reader.Reader[TestConfig, string] {
				return func(cfg TestConfig) string {
					return fmt.Sprintf("%s: %d", cfg.Prefix, n)
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of("LOG: 42"), outcome)
	})

	t.Run("accesses context fields", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](10),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n + len(cfg.DatabaseURL)
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
	})
}

func TestChainReaderK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		computation := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n * cfg.Multiplier
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}

func TestChainReaderK_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](0),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n * cfg.Multiplier
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(0), outcome)
	})

	t.Run("handles empty string", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](""),
			ChainReaderK(func(s string) reader.Reader[TestConfig, string] {
				return func(cfg TestConfig) string {
					return cfg.Prefix + s
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of("LOG"), outcome)
	})
}

func TestChainReaderK_Integration(t *testing.T) {
	t.Run("composes with Map and Chain", func(t *testing.T) {
		computation := F.Pipe3(
			Of[TestConfig](10),
			Map[TestConfig](func(x int) int { return x + 5 }),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n * cfg.Multiplier
				}
			}),
			Map[TestConfig](func(x int) string { return fmt.Sprintf("Result: %d", x) }),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of("Result: 45"), outcome)
	})
}

func TestChainReaderIOK_Success(t *testing.T) {
	t.Run("chains with ReaderIO that performs IO", func(t *testing.T) {
		counter := 0
		computation := F.Pipe1(
			Of[TestConfig](7),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						counter++
						return n * cfg.Multiplier
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(21), outcome)
		assert.Equal(t, 1, counter)
	})

	t.Run("chains multiple ReaderIO operations", func(t *testing.T) {
		log := []string{}
		computation := F.Pipe2(
			Of[TestConfig](5),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						log = append(log, fmt.Sprintf("add: %d", n))
						return n + cfg.Multiplier
					}
				}
			}),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						log = append(log, fmt.Sprintf("multiply: %d", n))
						return n * cfg.Multiplier
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(24), outcome)
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "add: 5", log[0])
		assert.Equal(t, "multiply: 8", log[1])
	})

	t.Run("chains with ReaderIO that formats string", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](100),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, string] {
				return func(cfg TestConfig) io.IO[string] {
					return func() string {
						return fmt.Sprintf("%s: %d", cfg.Prefix, n)
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of("LOG: 100"), outcome)
	})

	t.Run("accesses context in IO operation", func(t *testing.T) {
		computation := F.Pipe1(
			Of[TestConfig](10),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						return n + len(cfg.DatabaseURL)
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
	})
}

func TestChainReaderIOK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		sideEffectExecuted := false
		computation := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						sideEffectExecuted = true
						return n * cfg.Multiplier
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, sideEffectExecuted)
	})
}

func TestChainReaderIOK_EdgeCases(t *testing.T) {
	t.Run("handles zero value with side effects", func(t *testing.T) {
		callCount := 0
		computation := F.Pipe1(
			Of[TestConfig](0),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						callCount++
						return n * cfg.Multiplier
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(0), outcome)
		assert.Equal(t, 1, callCount)
	})

	t.Run("handles empty context fields", func(t *testing.T) {
		emptyConfig := TestConfig{Multiplier: 0, Prefix: "", DatabaseURL: ""}
		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, string] {
				return func(cfg TestConfig) io.IO[string] {
					return func() string {
						return fmt.Sprintf("%s%d", cfg.Prefix, n)
					}
				}
			}),
		)
		outcome := computation(emptyConfig)(context.Background())()
		assert.Equal(t, result.Of("42"), outcome)
	})
}

func TestChainReaderIOK_Integration(t *testing.T) {
	t.Run("composes with Map and Chain", func(t *testing.T) {
		log := []string{}
		computation := F.Pipe3(
			Of[TestConfig](8),
			Map[TestConfig](func(x int) int { return x + 2 }),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						log = append(log, fmt.Sprintf("processing: %d", n))
						return n * cfg.Multiplier
					}
				}
			}),
			Map[TestConfig](func(x int) string { return fmt.Sprintf("Final: %d", x) }),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of("Final: 30"), outcome)
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "processing: 10", log[0])
	})

	t.Run("composes ChainReaderK and ChainReaderIOK together", func(t *testing.T) {
		log := []string{}
		computation := F.Pipe2(
			Of[TestConfig](5),
			ChainReaderK(func(n int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return n + cfg.Multiplier
				}
			}),
			ChainReaderIOK(func(n int) readerio.ReaderIO[TestConfig, int] {
				return func(cfg TestConfig) io.IO[int] {
					return func() int {
						log = append(log, fmt.Sprintf("value: %d", n))
						return n * cfg.Multiplier
					}
				}
			}),
		)
		outcome := computation(testConfig)(context.Background())()
		assert.Equal(t, result.Of(24), outcome)
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "value: 8", log[0])
	})
}
