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
			ChainReaderK(addMultiplier),
			ChainReaderIOK(logValue),
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

func TestChainFirstThunkK_Success(t *testing.T) {
	t.Run("executes thunk but preserves original value", func(t *testing.T) {
		sideEffectExecuted := false

		sideEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					sideEffectExecuted = true
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](sideEffect),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
		assert.True(t, sideEffectExecuted)
	})

	t.Run("chains multiple side effects", func(t *testing.T) {
		log := []string{}

		logValue := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("log: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe2(
			Of[TestConfig](10),
			ChainFirstThunkK[TestConfig](logValue),
			ChainFirstThunkK[TestConfig](logValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(10), outcome)
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "log: 10", log[0])
		assert.Equal(t, "log: 10", log[1])
	})

	t.Run("side effect can access runtime context", func(t *testing.T) {
		var capturedCtx context.Context

		captureContext := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					capturedCtx = ctx
					return result.Of[any](nil)
				}
			}
		}

		ctx := context.Background()
		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](captureContext),
		)
		outcome := computation(testConfig)(ctx)()

		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, ctx, capturedCtx)
	})

	t.Run("side effect result is discarded", func(t *testing.T) {
		returnDifferentValue := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					return result.Of("different value")
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](returnDifferentValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
	})
}

func TestChainFirstThunkK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := fmt.Errorf("previous error")
		sideEffectExecuted := false

		sideEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					sideEffectExecuted = true
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainFirstThunkK[TestConfig](sideEffect),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, sideEffectExecuted)
	})

	t.Run("propagates error from thunk side effect", func(t *testing.T) {
		testErr := fmt.Errorf("side effect error")

		failingSideEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					return result.Left[any](testErr)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](failingSideEffect),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("stops execution on first error", func(t *testing.T) {
		testErr := fmt.Errorf("first error")
		secondEffectExecuted := false

		failingEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					return result.Left[any](testErr)
				}
			}
		}

		secondEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					secondEffectExecuted = true
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe2(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](failingEffect),
			ChainFirstThunkK[TestConfig](secondEffect),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, secondEffectExecuted)
	})
}

func TestChainFirstThunkK_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		callCount := 0

		countCalls := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					callCount++
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](0),
			ChainFirstThunkK[TestConfig](countCalls),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(0), outcome)
		assert.Equal(t, 1, callCount)
	})

	t.Run("handles empty string", func(t *testing.T) {
		var capturedValue string

		captureValue := func(s string) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					capturedValue = s
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](""),
			ChainFirstThunkK[TestConfig](captureValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(""), outcome)
		assert.Equal(t, "", capturedValue)
	})

	t.Run("handles nil pointer", func(t *testing.T) {
		var capturedPtr *int

		capturePtr := func(ptr *int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					capturedPtr = ptr
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig]((*int)(nil)),
			ChainFirstThunkK[TestConfig](capturePtr),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of((*int)(nil)), outcome)
		assert.Nil(t, capturedPtr)
	})
}

func TestChainFirstThunkK_Integration(t *testing.T) {
	t.Run("composes with Map and Chain", func(t *testing.T) {
		log := []string{}

		logValue := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("value: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe3(
			Of[TestConfig](5),
			Map[TestConfig](func(x int) int { return x * 2 }),
			ChainFirstThunkK[TestConfig](logValue),
			Map[TestConfig](func(x int) int { return x + 3 }),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(13), outcome) // (5 * 2) + 3
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "value: 10", log[0])
	})

	t.Run("composes with ChainThunkK", func(t *testing.T) {
		log := []string{}

		logSideEffect := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("side-effect: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		transformValue := func(n int) readerioresult.ReaderIOResult[string] {
			return func(ctx context.Context) io.IO[result.Result[string]] {
				return func() result.Result[string] {
					log = append(log, fmt.Sprintf("transform: %d", n))
					return result.Of(fmt.Sprintf("Result: %d", n))
				}
			}
		}

		computation := F.Pipe2(
			Of[TestConfig](42),
			ChainFirstThunkK[TestConfig](logSideEffect),
			ChainThunkK[TestConfig](transformValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Result: 42"), outcome)
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "side-effect: 42", log[0])
		assert.Equal(t, "transform: 42", log[1])
	})

	t.Run("composes with ChainReaderK and ChainReaderIOK", func(t *testing.T) {
		log := []string{}

		addMultiplier := func(n int) reader.Reader[TestConfig, int] {
			return func(cfg TestConfig) int {
				return n + cfg.Multiplier
			}
		}

		logReaderIO := func(n int) readerio.ReaderIO[TestConfig, int] {
			return func(cfg TestConfig) io.IO[int] {
				return func() int {
					log = append(log, fmt.Sprintf("reader-io: %d", n))
					return n * 2
				}
			}
		}

		logThunk := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("thunk: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe3(
			Of[TestConfig](5),
			ChainReaderK(addMultiplier),
			ChainReaderIOK(logReaderIO),
			ChainFirstThunkK[TestConfig](logThunk),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(16), outcome) // (5 + 3) * 2
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "reader-io: 8", log[0])
		assert.Equal(t, "thunk: 16", log[1])
	})
}

func TestTapThunkK_Success(t *testing.T) {
	t.Run("is alias for ChainFirstThunkK", func(t *testing.T) {
		log := []string{}

		logValue := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("tapped: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			TapThunkK[TestConfig](logValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "tapped: 42", log[0])
	})

	t.Run("useful for logging without changing value", func(t *testing.T) {
		log := []string{}

		logStep := func(step string) func(int) readerioresult.ReaderIOResult[any] {
			return func(n int) readerioresult.ReaderIOResult[any] {
				return func(ctx context.Context) io.IO[result.Result[any]] {
					return func() result.Result[any] {
						log = append(log, fmt.Sprintf("%s: %d", step, n))
						return result.Of[any](nil)
					}
				}
			}
		}

		computation := F.Pipe4(
			Of[TestConfig](10),
			TapThunkK[TestConfig](logStep("start")),
			Map[TestConfig](func(x int) int { return x * 2 }),
			TapThunkK[TestConfig](logStep("after-map")),
			Map[TestConfig](func(x int) int { return x + 5 }),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(25), outcome) // (10 * 2) + 5
		assert.Equal(t, 2, len(log))
		assert.Equal(t, "start: 10", log[0])
		assert.Equal(t, "after-map: 20", log[1])
	})

	t.Run("can perform IO operations", func(t *testing.T) {
		var ioExecuted bool

		performIO := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					// Simulate IO operation
					ioExecuted = true
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			TapThunkK[TestConfig](performIO),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
		assert.True(t, ioExecuted)
	})
}

func TestTapThunkK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := fmt.Errorf("previous error")
		tapExecuted := false

		tapValue := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					tapExecuted = true
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			Fail[TestConfig, int](testErr),
			TapThunkK[TestConfig](tapValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, tapExecuted)
	})

	t.Run("propagates error from tap operation", func(t *testing.T) {
		testErr := fmt.Errorf("tap error")

		failingTap := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					return result.Left[any](testErr)
				}
			}
		}

		computation := F.Pipe1(
			Of[TestConfig](42),
			TapThunkK[TestConfig](failingTap),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}

func TestTapThunkK_EdgeCases(t *testing.T) {
	t.Run("handles multiple taps in sequence", func(t *testing.T) {
		log := []string{}

		tap1 := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, "tap1")
					return result.Of[any](nil)
				}
			}
		}

		tap2 := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, "tap2")
					return result.Of[any](nil)
				}
			}
		}

		tap3 := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, "tap3")
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe3(
			Of[TestConfig](42),
			TapThunkK[TestConfig](tap1),
			TapThunkK[TestConfig](tap2),
			TapThunkK[TestConfig](tap3),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, []string{"tap1", "tap2", "tap3"}, log)
	})
}

func TestTapThunkK_Integration(t *testing.T) {
	t.Run("real-world logging scenario", func(t *testing.T) {
		log := []string{}

		logStart := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("Starting computation with: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		logIntermediate := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("Intermediate result: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		logFinal := func(s string) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("Final result: %s", s))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe5(
			Of[TestConfig](10),
			TapThunkK[TestConfig](logStart),
			Map[TestConfig](func(x int) int { return x * 3 }),
			TapThunkK[TestConfig](logIntermediate),
			Map[TestConfig](func(x int) string { return fmt.Sprintf("Value: %d", x) }),
			TapThunkK[TestConfig](logFinal),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of("Value: 30"), outcome)
		assert.Equal(t, 3, len(log))
		assert.Equal(t, "Starting computation with: 10", log[0])
		assert.Equal(t, "Intermediate result: 30", log[1])
		assert.Equal(t, "Final result: Value: 30", log[2])
	})

	t.Run("composes with FromThunk", func(t *testing.T) {
		log := []string{}

		thunk := func(ctx context.Context) io.IO[result.Result[int]] {
			return func() result.Result[int] {
				return result.Of(100)
			}
		}

		logValue := func(n int) readerioresult.ReaderIOResult[any] {
			return func(ctx context.Context) io.IO[result.Result[any]] {
				return func() result.Result[any] {
					log = append(log, fmt.Sprintf("value: %d", n))
					return result.Of[any](nil)
				}
			}
		}

		computation := F.Pipe1(
			FromThunk[TestConfig](thunk),
			TapThunkK[TestConfig](logValue),
		)
		outcome := computation(testConfig)(context.Background())()

		assert.Equal(t, result.Of(100), outcome)
		assert.Equal(t, 1, len(log))
		assert.Equal(t, "value: 100", log[0])
	})
}

func TestAsks_Success(t *testing.T) {
	t.Run("extracts a field from context", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		getHost := Asks(func(cfg Config) string {
			return cfg.Host
		})

		result, err := runEffect(getHost, Config{Host: "localhost", Port: 8080})

		assert.NoError(t, err)
		assert.Equal(t, "localhost", result)
	})

	t.Run("extracts multiple fields and computes derived value", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		getURL := Asks(func(cfg Config) string {
			return fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
		})

		result, err := runEffect(getURL, Config{Host: "example.com", Port: 443})

		assert.NoError(t, err)
		assert.Equal(t, "http://example.com:443", result)
	})

	t.Run("extracts numeric field", func(t *testing.T) {
		getPort := Asks(func(cfg TestConfig) int {
			return cfg.Multiplier
		})

		result, err := runEffect(getPort, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 3, result)
	})

	t.Run("computes value from context", func(t *testing.T) {
		type Config struct {
			Width  int
			Height int
		}

		getArea := Asks(func(cfg Config) int {
			return cfg.Width * cfg.Height
		})

		result, err := runEffect(getArea, Config{Width: 10, Height: 20})

		assert.NoError(t, err)
		assert.Equal(t, 200, result)
	})

	t.Run("transforms string field", func(t *testing.T) {
		getUpperPrefix := Asks(func(cfg TestConfig) string {
			return fmt.Sprintf("[%s]", cfg.Prefix)
		})

		result, err := runEffect(getUpperPrefix, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "[LOG]", result)
	})
}

func TestAsks_EdgeCases(t *testing.T) {
	t.Run("handles zero values", func(t *testing.T) {
		type Config struct {
			Value int
		}

		getValue := Asks(func(cfg Config) int {
			return cfg.Value
		})

		result, err := runEffect(getValue, Config{Value: 0})

		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("handles empty string", func(t *testing.T) {
		type Config struct {
			Name string
		}

		getName := Asks(func(cfg Config) string {
			return cfg.Name
		})

		result, err := runEffect(getName, Config{Name: ""})

		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("handles nil pointer fields", func(t *testing.T) {
		type Config struct {
			Data *string
		}

		hasData := Asks(func(cfg Config) bool {
			return cfg.Data != nil
		})

		result, err := runEffect(hasData, Config{Data: nil})

		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("handles complex nested structures", func(t *testing.T) {
		type Database struct {
			Host string
			Port int
		}
		type Config struct {
			DB Database
		}

		getDBHost := Asks(func(cfg Config) string {
			return cfg.DB.Host
		})

		result, err := runEffect(getDBHost, Config{
			DB: Database{Host: "db.example.com", Port: 5432},
		})

		assert.NoError(t, err)
		assert.Equal(t, "db.example.com", result)
	})
}

func TestAsks_Integration(t *testing.T) {
	t.Run("composes with Map", func(t *testing.T) {
		type Config struct {
			Value int
		}

		computation := F.Pipe1(
			Asks(func(cfg Config) int {
				return cfg.Value
			}),
			Map[Config](func(x int) int { return x * 2 }),
		)

		result, err := runEffect(computation, Config{Value: 21})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("composes with Chain", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		computation := F.Pipe1(
			Asks(func(cfg Config) int {
				return cfg.Multiplier
			}),
			Chain(func(mult int) Effect[Config, int] {
				return Of[Config](mult * 10)
			}),
		)

		result, err := runEffect(computation, Config{Multiplier: 5})

		assert.NoError(t, err)
		assert.Equal(t, 50, result)
	})

	t.Run("composes with ChainReaderK", func(t *testing.T) {
		computation := F.Pipe1(
			Asks(func(cfg TestConfig) int {
				return cfg.Multiplier
			}),
			ChainReaderK(func(mult int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return mult + len(cfg.Prefix)
				}
			}),
		)

		result, err := runEffect(computation, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 6, result) // 3 + len("LOG")
	})

	t.Run("composes with ChainReaderIOK", func(t *testing.T) {
		log := []string{}

		computation := F.Pipe1(
			Asks(func(cfg TestConfig) string {
				return cfg.Prefix
			}),
			ChainReaderIOK(func(prefix string) readerio.ReaderIO[TestConfig, string] {
				return func(cfg TestConfig) io.IO[string] {
					return func() string {
						log = append(log, "executed")
						return fmt.Sprintf("%s:%d", prefix, cfg.Multiplier)
					}
				}
			}),
		)

		result, err := runEffect(computation, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "LOG:3", result)
		assert.Equal(t, 1, len(log))
	})

	t.Run("multiple Asks in sequence", func(t *testing.T) {
		type Config struct {
			First  string
			Second string
		}

		computation := F.Pipe2(
			Asks(func(cfg Config) string {
				return cfg.First
			}),
			Chain(func(_ string) Effect[Config, string] {
				return Asks(func(cfg Config) string {
					return cfg.Second
				})
			}),
			Map[Config](func(s string) string {
				return "Result: " + s
			}),
		)

		result, err := runEffect(computation, Config{First: "A", Second: "B"})

		assert.NoError(t, err)
		assert.Equal(t, "Result: B", result)
	})

	t.Run("Asks combined with Ask", func(t *testing.T) {
		type Config struct {
			Value int
		}

		computation := F.Pipe1(
			Ask[Config](),
			Chain(func(cfg Config) Effect[Config, int] {
				return Asks(func(c Config) int {
					return c.Value * 2
				})
			}),
		)

		result, err := runEffect(computation, Config{Value: 15})

		assert.NoError(t, err)
		assert.Equal(t, 30, result)
	})
}

func TestAsks_Comparison(t *testing.T) {
	t.Run("Asks vs Ask with Map", func(t *testing.T) {
		type Config struct {
			Port int
		}

		// Using Asks
		asksVersion := Asks(func(cfg Config) int {
			return cfg.Port
		})

		// Using Ask + Map
		askMapVersion := F.Pipe1(
			Ask[Config](),
			Map[Config](func(cfg Config) int {
				return cfg.Port
			}),
		)

		cfg := Config{Port: 8080}

		result1, err1 := runEffect(asksVersion, cfg)
		result2, err2 := runEffect(askMapVersion, cfg)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, result1, result2)
		assert.Equal(t, 8080, result1)
	})

	t.Run("Asks is more concise than Ask + Map", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		// Asks is more direct for field extraction
		getHost := Asks(func(cfg Config) string {
			return cfg.Host
		})

		result, err := runEffect(getHost, Config{Host: "api.example.com", Port: 443})

		assert.NoError(t, err)
		assert.Equal(t, "api.example.com", result)
	})
}

func TestAsks_RealWorldScenarios(t *testing.T) {
	t.Run("extract database connection string", func(t *testing.T) {
		type DatabaseConfig struct {
			Host     string
			Port     int
			Database string
			User     string
		}

		getConnectionString := Asks(func(cfg DatabaseConfig) string {
			return fmt.Sprintf("postgres://%s@%s:%d/%s",
				cfg.User, cfg.Host, cfg.Port, cfg.Database)
		})

		result, err := runEffect(getConnectionString, DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "myapp",
			User:     "admin",
		})

		assert.NoError(t, err)
		assert.Equal(t, "postgres://admin@localhost:5432/myapp", result)
	})

	t.Run("compute API endpoint from config", func(t *testing.T) {
		type APIConfig struct {
			Protocol string
			Host     string
			Port     int
			BasePath string
		}

		getEndpoint := Asks(func(cfg APIConfig) string {
			return fmt.Sprintf("%s://%s:%d%s",
				cfg.Protocol, cfg.Host, cfg.Port, cfg.BasePath)
		})

		result, err := runEffect(getEndpoint, APIConfig{
			Protocol: "https",
			Host:     "api.example.com",
			Port:     443,
			BasePath: "/v1",
		})

		assert.NoError(t, err)
		assert.Equal(t, "https://api.example.com:443/v1", result)
	})

	t.Run("validate configuration", func(t *testing.T) {
		type Config struct {
			Timeout    int
			MaxRetries int
		}

		isValid := Asks(func(cfg Config) bool {
			return cfg.Timeout > 0 && cfg.MaxRetries >= 0
		})

		// Valid config
		result1, err1 := runEffect(isValid, Config{Timeout: 30, MaxRetries: 3})
		assert.NoError(t, err1)
		assert.True(t, result1)

		// Invalid config
		result2, err2 := runEffect(isValid, Config{Timeout: 0, MaxRetries: 3})
		assert.NoError(t, err2)
		assert.False(t, result2)
	})

	t.Run("extract feature flags", func(t *testing.T) {
		type FeatureFlags struct {
			EnableNewUI     bool
			EnableBetaAPI   bool
			EnableAnalytics bool
		}

		hasNewUI := Asks[FeatureFlags](func(flags FeatureFlags) bool {
			return flags.EnableNewUI
		})

		result, err := runEffect(hasNewUI, FeatureFlags{
			EnableNewUI:     true,
			EnableBetaAPI:   false,
			EnableAnalytics: true,
		})

		assert.NoError(t, err)
		assert.True(t, result)
	})
}
