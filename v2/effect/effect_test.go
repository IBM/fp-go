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

	"github.com/IBM/fp-go/v2/context/readerioresult"
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

func TestFromThunk_Success(t *testing.T) {
	t.Run("lifts successful thunk into effect", func(t *testing.T) {
		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				return result.Of(42)
			}
		}

		computation := FromThunk[TestConfig](thunk)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("thunk ignores effect context", func(t *testing.T) {
		thunk := func(ctx context.Context) io.IO[result.Result[string]] {
			return func() result.Result[string] {
				// Thunk doesn't use TestConfig, only runtime context
				return result.Of("context-independent")
			}
		}

		computation := FromThunk[TestConfig](thunk)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("context-independent"), outcome)
	})

	t.Run("thunk can perform IO operations", func(t *testing.T) {
		counter := 0
		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				counter++
				return result.Of(counter)
			}
		}

		computation := FromThunk[TestConfig](thunk)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(1), outcome)
		assert.Equal(t, 1, counter)
	})
}

func TestFromThunk_Failure(t *testing.T) {
	t.Run("propagates error from thunk", func(t *testing.T) {
		testErr := fmt.Errorf("thunk error")
		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				return result.Left[int](testErr)
			}
		}

		computation := FromThunk[TestConfig](thunk)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}

func TestFromThunk_Integration(t *testing.T) {
	t.Run("composes with other effect operations", func(t *testing.T) {
		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				return result.Of(10)
			}
		}

		computation := F.Pipe2(
			FromThunk[TestConfig](thunk),
			Map[TestConfig](func(x int) int { return x * 2 }),
			Map[TestConfig](func(x int) string { return fmt.Sprintf("Result: %d", x) }),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Result: 20"), outcome)
	})
}

func TestChainThunkK_Success(t *testing.T) {
	t.Run("chains with thunk that performs IO", func(t *testing.T) {
		processValue := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Processed: %d", n))
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainThunkK[TestConfig](processValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Processed: 42"), outcome)
	})

	t.Run("chains multiple thunk operations", func(t *testing.T) {
		log := []string{}

		logAndAdd := func(n int) readerioresult.ReaderIOResult[int] {
			return func(ctx context.Context) io.IO[result.Result[int]] {
				return func() result.Result[int] {
					log = append(log, fmt.Sprintf("add: %d", n))
					return result.Of(n + 10)
				}
			}
		}

		logAndMultiply := func(n int) readerioresult.ReaderIOResult[int] {
			return func(ctx context.Context) io.IO[result.Result[int]] {
				return func() result.Result[int] {
					log = append(log, fmt.Sprintf("multiply: %d", n))
					return result.Of(n * 2)
				}
			}
		}

		computation := F.Pipe2(
			Of[TestConfig](5),
			ChainThunkK[TestConfig](logAndAdd),
			ChainThunkK[TestConfig](logAndMultiply),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(30), outcome) // (5 + 10) * 2
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "add: 5", log[0])
		assert.Equal(t, "multiply: 15", log[1])
	})

	t.Run("thunk can access runtime context", func(t *testing.T) {
		var capturedCtx context.Context

		captureContext := func(n int) readerioresult.ReaderIOResult[int] {
			return func(ctx context.Context) io.IO[result.Result[int]] {
				return func() result.Result[int] {
					capturedCtx = ctx
					return result.Of(n * 2)
				}
			}
		}

		ctx := context.Background()
		computation := F.Pipe1(
			Of[TestConfig](21),
			ChainThunkK[TestConfig](captureContext),
		)
		outcome := computation(testConfig)(ctx)()

		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, ctx, capturedCtx)
	})
}

func TestChainThunkK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := fmt.Errorf("previous error")
		sideEffectExecuted := false

		processValue := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					sideEffectExecuted = true
					return result.Of(fmt.Sprintf("%d", n))
				}
			}
		}

		computation := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainThunkK[TestConfig](processValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[string](testErr), outcome)
		assert.False(t, sideEffectExecuted)
	})

	t.Run("propagates error from thunk", func(t *testing.T) {
		testErr := fmt.Errorf("thunk error")

		failingThunk := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					return result.Left[string](testErr)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainThunkK[TestConfig](failingThunk),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[string](testErr), outcome)
	})
}

func TestChainThunkK_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		callCount := 0

		processZero := func(n int) readerioresult.ReaderIOResult[int] {
			return func(ctx context.Context) io.IO[result.Result[int]] {
				return func() result.Result[int] {
					callCount++
					return result.Of(n + 1)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](0),
			ChainThunkK[TestConfig](processZero),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(1), outcome)
		assert.Equal(t, 1, callCount)
	})

	t.Run("handles empty string", func(t *testing.T) {
		addPrefix := func(s string) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					return result.Of("prefix:" + s)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](""),
			ChainThunkK[TestConfig](addPrefix),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("prefix:"), outcome)
	})
}

func TestChainThunkK_Integration(t *testing.T) {
	t.Run("composes with Map and Chain", func(t *testing.T) {
		log := []string{}

		logAndProcess := func(n int) readerioresult.ReaderIOResult[int] {
			return func(ctx context.Context) io.IO[result.Result[int]] {
				return func() result.Result[int] {
					log = append(log, fmt.Sprintf("processing: %d", n))
					return result.Of(n * 3)
				}
			}
		}

		computation := F.Pipe3(
			Of[TestConfig](7),
			Map[TestConfig](func(x int) int { return x + 3 }),
			ChainThunkK[TestConfig](logAndProcess),
			Map[TestConfig](func(x int) string { return fmt.Sprintf("Final: %d", x) }),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Final: 30"), outcome) // (7 + 3) * 3
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "processing: 10", log[0])
	})

	t.Run("composes with ChainReaderK and ChainReaderIOK", func(t *testing.T) {
		log := []string{}

		addMultiplier := func(n int) reader.Reader[TestConfig, int] {
			return func(cfg TestConfig) int {
				return n + cfg.Multiplier
			}
		}

		logValue := func(n int) readerio.ReaderIO[TestConfig, int] {
			return func(cfg TestConfig) io.IO[int] {
				return func() int {
					log = append(log, fmt.Sprintf("reader-io: %d", n))
					return n * 2
				}
			}
		}

		processThunk := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					log = append(log, fmt.Sprintf("thunk: %d", n))
					return result.Of(fmt.Sprintf("Result: %d", n))
				}
			}
		}

		computation := F.Pipe3(
			Of[TestConfig](5),
			ChainReaderK[TestConfig](addMultiplier),
			ChainReaderIOK[TestConfig](logValue),
			ChainThunkK[TestConfig](processThunk),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Result: 16"), outcome) // (5 + 3) * 2
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "reader-io: 8", log[0])
		assert.Equal(t, "thunk: 16", log[1])
	})

	t.Run("composes FromThunk with ChainThunkK", func(t *testing.T) {
		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				return result.Of(100)
			}
		}

		processValue := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					return result.Of(fmt.Sprintf("Value: %d", n))
				}
			}
		}

		computation := F.Pipe1(
			FromThunk[TestConfig](thunk),
			ChainThunkK[TestConfig](processValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Value: 100"), outcome)
	})
}
