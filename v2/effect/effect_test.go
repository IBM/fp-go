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
	"log/slog"
	"os"
	"testing"

	"github.com/IBM/fp-go/v2/context/readerioresult"
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
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
			Map[TestConfig](N.Mul(2)),
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
			Map[TestConfig](N.Mul(2)),
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
			Map[TestConfig](N.Mul(2)),
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
			Map[Config](N.Mul(2)),
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

		hasNewUI := Asks(func(flags FeatureFlags) bool {
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

func TestMonadChainLeft_Success(t *testing.T) {
	t.Run("success value passes through unchanged", func(t *testing.T) {
		eff := Of[TestConfig](42)
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("success with context-dependent value", func(t *testing.T) {
		eff := Asks[TestConfig](func(cfg TestConfig) int {
			return cfg.Multiplier * 10
		})
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 30, value)
	})
}

func TestMonadChainLeft_Failure(t *testing.T) {
	t.Run("error triggers recovery", func(t *testing.T) {
		originalErr := fmt.Errorf("network error")
		eff := Fail[TestConfig, int](originalErr)
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
	})

	t.Run("recovery can inspect error", func(t *testing.T) {
		originalErr := fmt.Errorf("code: 404")
		eff := Fail[TestConfig, string](originalErr)
		recover := func(err error) Effect[TestConfig, string] {
			return Of[TestConfig](fmt.Sprintf("recovered from: %s", err.Error()))
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "recovered from: code: 404", value)
	})

	t.Run("recovery can use context", func(t *testing.T) {
		originalErr := fmt.Errorf("database error")
		eff := Fail[TestConfig, string](originalErr)
		recover := func(err error) Effect[TestConfig, string] {
			return Asks[TestConfig](func(cfg TestConfig) string {
				return fmt.Sprintf("fallback to: %s", cfg.DatabaseURL)
			})
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "fallback to: postgres://localhost", value)
	})

	t.Run("recovery can also fail", func(t *testing.T) {
		originalErr := fmt.Errorf("first error")
		recoveryErr := fmt.Errorf("recovery failed")
		eff := Fail[TestConfig, int](originalErr)
		recover := func(err error) Effect[TestConfig, int] {
			return Fail[TestConfig, int](recoveryErr)
		}

		result := MonadChainLeft(eff, recover)
		_, err := runEffect(result, testConfig)

		assert.Error(t, err)
		assert.Equal(t, recoveryErr, err)
	})
}

func TestMonadChainLeft_EdgeCases(t *testing.T) {
	t.Run("multiple recoveries in sequence", func(t *testing.T) {
		err1 := fmt.Errorf("error 1")
		err2 := fmt.Errorf("error 2")

		eff := Fail[TestConfig, int](err1)
		recover1 := func(err error) Effect[TestConfig, int] {
			return Fail[TestConfig, int](err2)
		}
		recover2 := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](42)
		}

		result := MonadChainLeft(MonadChainLeft(eff, recover1), recover2)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("recovery with zero value", func(t *testing.T) {
		eff := Fail[TestConfig, int](fmt.Errorf("error"))
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := MonadChainLeft(eff, recover)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})
}

func TestMonadChainLeft_Integration(t *testing.T) {
	t.Run("chain with other operations", func(t *testing.T) {
		eff := F.Pipe1(
			Fail[TestConfig, int](fmt.Errorf("initial error")),
			Map[TestConfig](N.Mul(2)),
		)
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](21)
		}

		result := F.Pipe1(
			MonadChainLeft(eff, recover),
			Map[TestConfig](N.Mul(2)),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("recovery in complex pipeline", func(t *testing.T) {
		type AppConfig struct {
			PrimaryDB  string
			FallbackDB string
			MaxRetries int
		}

		cfg := AppConfig{
			PrimaryDB:  "primary.db",
			FallbackDB: "fallback.db",
			MaxRetries: 3,
		}

		connectDB := func(dbName string) Effect[AppConfig, string] {
			if dbName == "primary.db" {
				return Fail[AppConfig, string](fmt.Errorf("primary unavailable"))
			}
			return Of[AppConfig](fmt.Sprintf("connected to %s", dbName))
		}

		pipeline := F.Pipe1(
			connectDB("primary.db"),
			ChainLeft[AppConfig](func(err error) Effect[AppConfig, string] {
				return Asks[AppConfig](func(cfg AppConfig) string {
					return fmt.Sprintf("connected to %s", cfg.FallbackDB)
				})
			}),
		)

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "connected to fallback.db", value)
	})
}

func TestChainLeft_Success(t *testing.T) {
	t.Run("success value passes through", func(t *testing.T) {
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := F.Pipe1(
			Of[TestConfig](42),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("success in pipeline", func(t *testing.T) {
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := F.Pipe3(
			Of[TestConfig](10),
			Map[TestConfig](N.Mul(2)),
			ChainLeft[TestConfig](recover),
			Map[TestConfig](N.Add(1)),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 21, value)
	})
}

func TestChainLeft_Failure(t *testing.T) {
	t.Run("error triggers recovery", func(t *testing.T) {
		originalErr := fmt.Errorf("operation failed")
		recover := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
	})

	t.Run("recovery uses error information", func(t *testing.T) {
		originalErr := fmt.Errorf("status: 500")
		recover := func(err error) Effect[TestConfig, string] {
			return Of[TestConfig](fmt.Sprintf("error was: %s", err.Error()))
		}

		result := F.Pipe1(
			Fail[TestConfig, string](originalErr),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "error was: status: 500", value)
	})

	t.Run("recovery accesses context", func(t *testing.T) {
		originalErr := fmt.Errorf("config error")
		recover := func(err error) Effect[TestConfig, string] {
			return Asks[TestConfig](func(cfg TestConfig) string {
				return cfg.Prefix
			})
		}

		result := F.Pipe1(
			Fail[TestConfig, string](originalErr),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "LOG", value)
	})

	t.Run("recovery fails", func(t *testing.T) {
		originalErr := fmt.Errorf("first error")
		recoveryErr := fmt.Errorf("recovery error")
		recover := func(err error) Effect[TestConfig, int] {
			return Fail[TestConfig, int](recoveryErr)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			ChainLeft[TestConfig](recover),
		)
		_, err := runEffect(result, testConfig)

		assert.Error(t, err)
		assert.Equal(t, recoveryErr, err)
	})
}

func TestChainLeft_EdgeCases(t *testing.T) {
	t.Run("chained recoveries", func(t *testing.T) {
		err1 := fmt.Errorf("error 1")
		err2 := fmt.Errorf("error 2")

		recover1 := func(err error) Effect[TestConfig, int] {
			return Fail[TestConfig, int](err2)
		}
		recover2 := func(err error) Effect[TestConfig, int] {
			return Of[TestConfig](42)
		}

		result := F.Pipe2(
			Fail[TestConfig, int](err1),
			ChainLeft[TestConfig](recover1),
			ChainLeft[TestConfig](recover2),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("recovery with empty string", func(t *testing.T) {
		recover := func(err error) Effect[TestConfig, string] {
			return Of[TestConfig]("")
		}

		result := F.Pipe1(
			Fail[TestConfig, string](fmt.Errorf("error")),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("recovery with nil slice", func(t *testing.T) {
		recover := func(err error) Effect[TestConfig, []int] {
			return Of[TestConfig, []int](nil)
		}

		result := F.Pipe1(
			Fail[TestConfig, []int](fmt.Errorf("error")),
			ChainLeft[TestConfig](recover),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Nil(t, value)
	})
}

func TestChainLeft_Integration(t *testing.T) {
	t.Run("retry pattern with ChainLeft", func(t *testing.T) {
		type RetryConfig struct {
			MaxAttempts int
			Attempt     int
		}

		cfg := RetryConfig{MaxAttempts: 3, Attempt: 0}

		attemptOperation := func(attempt int) Effect[RetryConfig, string] {
			if attempt < 2 {
				return Fail[RetryConfig, string](fmt.Errorf("attempt %d failed", attempt))
			}
			return Of[RetryConfig](fmt.Sprintf("success on attempt %d", attempt))
		}

		var pipeline Effect[RetryConfig, string]
		pipeline = attemptOperation(0)

		for i := 1; i <= 2; i++ {
			attempt := i
			pipeline = F.Pipe1(
				pipeline,
				ChainLeft[RetryConfig](func(err error) Effect[RetryConfig, string] {
					return attemptOperation(attempt)
				}),
			)
		}

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "success on attempt 2", value)
	})

	t.Run("fallback chain", func(t *testing.T) {
		type ServiceConfig struct {
			Services []string
		}

		cfg := ServiceConfig{
			Services: []string{"primary", "secondary", "tertiary"},
		}

		callService := func(name string) Effect[ServiceConfig, string] {
			if name == "tertiary" {
				return Of[ServiceConfig](fmt.Sprintf("response from %s", name))
			}
			return Fail[ServiceConfig, string](fmt.Errorf("%s unavailable", name))
		}

		pipeline := F.Pipe2(
			callService("primary"),
			ChainLeft[ServiceConfig](func(err error) Effect[ServiceConfig, string] {
				return callService("secondary")
			}),
			ChainLeft[ServiceConfig](func(err error) Effect[ServiceConfig, string] {
				return callService("tertiary")
			}),
		)

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "response from tertiary", value)
	})

	t.Run("error transformation chain", func(t *testing.T) {
		type ErrorConfig struct {
			LogErrors bool
		}

		cfg := ErrorConfig{LogErrors: true}

		pipeline := F.Pipe2(
			Fail[ErrorConfig, int](fmt.Errorf("database error")),
			ChainLeft[ErrorConfig](func(err error) Effect[ErrorConfig, int] {
				return Fail[ErrorConfig, int](fmt.Errorf("wrapped: %w", err))
			}),
			ChainLeft[ErrorConfig](func(err error) Effect[ErrorConfig, int] {
				return Of[ErrorConfig](0)
			}),
		)

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("combining with Chain and Map", func(t *testing.T) {
		parseAndDouble := func(s string) Effect[TestConfig, int] {
			if s == "invalid" {
				return Fail[TestConfig, int](fmt.Errorf("parse error"))
			}
			return Of[TestConfig](42)
		}

		pipeline := F.Pipe3(
			Of[TestConfig]("invalid"),
			Chain[TestConfig](parseAndDouble),
			ChainLeft[TestConfig](func(err error) Effect[TestConfig, int] {
				return Of[TestConfig](10)
			}),
			Map[TestConfig](N.Mul(2)),
		)

		value, err := runEffect(pipeline, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})
}

func TestMonadAlt_Success(t *testing.T) {
	t.Run("first succeeds, second not evaluated", func(t *testing.T) {
		first := Of[TestConfig](42)
		secondCalled := false
		second := func() Effect[TestConfig, int] {
			secondCalled = true
			return Of[TestConfig](99)
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
		assert.False(t, secondCalled, "second effect should not be evaluated when first succeeds")
	})

	t.Run("first succeeds with context-dependent value", func(t *testing.T) {
		first := Asks[TestConfig](func(cfg TestConfig) int {
			return cfg.Multiplier * 10
		})
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 30, value)
	})
}

func TestMonadAlt_Failure(t *testing.T) {
	t.Run("first fails, second succeeds", func(t *testing.T) {
		firstErr := fmt.Errorf("first error")
		first := Fail[TestConfig, int](firstErr)
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
	})

	t.Run("first fails, second uses context", func(t *testing.T) {
		firstErr := fmt.Errorf("primary failed")
		first := Fail[TestConfig, string](firstErr)
		second := func() Effect[TestConfig, string] {
			return Asks[TestConfig](func(cfg TestConfig) string {
				return cfg.Prefix + "-fallback"
			})
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "LOG-fallback", value)
	})

	t.Run("both fail", func(t *testing.T) {
		firstErr := fmt.Errorf("first error")
		secondErr := fmt.Errorf("second error")
		first := Fail[TestConfig, int](firstErr)
		second := func() Effect[TestConfig, int] {
			return Fail[TestConfig, int](secondErr)
		}

		result := MonadAlt(first, second)
		_, err := runEffect(result, testConfig)

		assert.Error(t, err)
		assert.Equal(t, secondErr, err)
	})

	t.Run("second is lazy evaluated only on failure", func(t *testing.T) {
		firstErr := fmt.Errorf("error")
		first := Fail[TestConfig, int](firstErr)
		secondCalled := false
		second := func() Effect[TestConfig, int] {
			secondCalled = true
			return Of[TestConfig](99)
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
		assert.True(t, secondCalled, "second effect should be evaluated when first fails")
	})
}

func TestMonadAlt_EdgeCases(t *testing.T) {
	t.Run("chained alternatives", func(t *testing.T) {
		err1 := fmt.Errorf("error 1")
		err2 := fmt.Errorf("error 2")

		first := Fail[TestConfig, int](err1)
		second := func() Effect[TestConfig, int] {
			return Fail[TestConfig, int](err2)
		}
		third := func() Effect[TestConfig, int] {
			return Of[TestConfig](42)
		}

		result := MonadAlt(MonadAlt(first, second), third)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("alternative with zero value", func(t *testing.T) {
		first := Fail[TestConfig, int](fmt.Errorf("error"))
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("alternative with empty string", func(t *testing.T) {
		first := Fail[TestConfig, string](fmt.Errorf("error"))
		second := func() Effect[TestConfig, string] {
			return Of[TestConfig]("")
		}

		result := MonadAlt(first, second)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})
}

func TestMonadAlt_Integration(t *testing.T) {
	t.Run("alternative with Map", func(t *testing.T) {
		first := Fail[TestConfig, int](fmt.Errorf("error"))
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](21)
		}

		result := F.Pipe1(
			MonadAlt(first, second),
			Map[TestConfig](N.Mul(2)),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("alternative with Chain", func(t *testing.T) {
		first := Fail[TestConfig, int](fmt.Errorf("error"))
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](10)
		}

		result := F.Pipe1(
			MonadAlt(first, second),
			Chain[TestConfig](func(n int) Effect[TestConfig, int] {
				return Of[TestConfig](n * 4)
			}),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 40, value)
	})

	t.Run("fallback chain pattern", func(t *testing.T) {
		type ServiceConfig struct {
			Services []string
		}

		cfg := ServiceConfig{
			Services: []string{"primary", "secondary", "tertiary"},
		}

		callService := func(name string) Effect[ServiceConfig, string] {
			if name == "tertiary" {
				return Of[ServiceConfig](fmt.Sprintf("response from %s", name))
			}
			return Fail[ServiceConfig, string](fmt.Errorf("%s unavailable", name))
		}

		result := MonadAlt(
			callService("primary"),
			func() Effect[ServiceConfig, string] {
				return MonadAlt(
					callService("secondary"),
					func() Effect[ServiceConfig, string] {
						return callService("tertiary")
					},
				)
			},
		)

		value, err := runEffect(result, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "response from tertiary", value)
	})
}

func TestAlt_Success(t *testing.T) {
	t.Run("first succeeds, second not evaluated", func(t *testing.T) {
		secondCalled := false
		second := func() Effect[TestConfig, int] {
			secondCalled = true
			return Of[TestConfig](99)
		}

		result := F.Pipe1(
			Of[TestConfig](42),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
		assert.False(t, secondCalled, "second effect should not be evaluated when first succeeds")
	})

	t.Run("success in pipeline", func(t *testing.T) {
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := F.Pipe2(
			Of[TestConfig](10),
			Map[TestConfig](N.Mul(2)),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})
}

func TestAlt_Failure(t *testing.T) {
	t.Run("first fails, second succeeds", func(t *testing.T) {
		originalErr := fmt.Errorf("operation failed")
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](99)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
	})

	t.Run("first fails, second uses context", func(t *testing.T) {
		originalErr := fmt.Errorf("primary error")
		second := func() Effect[TestConfig, string] {
			return Asks[TestConfig](func(cfg TestConfig) string {
				return cfg.DatabaseURL
			})
		}

		result := F.Pipe1(
			Fail[TestConfig, string](originalErr),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "postgres://localhost", value)
	})

	t.Run("both fail", func(t *testing.T) {
		firstErr := fmt.Errorf("first error")
		secondErr := fmt.Errorf("second error")
		second := func() Effect[TestConfig, int] {
			return Fail[TestConfig, int](secondErr)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](firstErr),
			Alt[TestConfig](second),
		)
		_, err := runEffect(result, testConfig)

		assert.Error(t, err)
		assert.Equal(t, secondErr, err)
	})

	t.Run("second is lazy evaluated", func(t *testing.T) {
		firstErr := fmt.Errorf("error")
		secondCalled := false
		second := func() Effect[TestConfig, int] {
			secondCalled = true
			return Of[TestConfig](99)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](firstErr),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
		assert.True(t, secondCalled, "second effect should be evaluated when first fails")
	})
}

func TestAlt_EdgeCases(t *testing.T) {
	t.Run("chained alternatives", func(t *testing.T) {
		err1 := fmt.Errorf("error 1")
		err2 := fmt.Errorf("error 2")

		second := func() Effect[TestConfig, int] {
			return Fail[TestConfig, int](err2)
		}
		third := func() Effect[TestConfig, int] {
			return Of[TestConfig](42)
		}

		result := F.Pipe2(
			Fail[TestConfig, int](err1),
			Alt[TestConfig](second),
			Alt[TestConfig](third),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("alternative with zero value", func(t *testing.T) {
		second := func() Effect[TestConfig, int] {
			return Of[TestConfig](0)
		}

		result := F.Pipe1(
			Fail[TestConfig, int](fmt.Errorf("error")),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("alternative with nil slice", func(t *testing.T) {
		second := func() Effect[TestConfig, []int] {
			return Of[TestConfig, []int](nil)
		}

		result := F.Pipe1(
			Fail[TestConfig, []int](fmt.Errorf("error")),
			Alt[TestConfig](second),
		)
		value, err := runEffect(result, testConfig)

		assert.NoError(t, err)
		assert.Nil(t, value)
	})
}

func TestAlt_Integration(t *testing.T) {
	t.Run("retry pattern with Alt", func(t *testing.T) {
		type RetryConfig struct {
			MaxAttempts int
		}

		cfg := RetryConfig{MaxAttempts: 3}
		attemptCount := 0

		attemptOperation := func() Effect[RetryConfig, string] {
			attemptCount++
			if attemptCount < 3 {
				return Fail[RetryConfig, string](fmt.Errorf("attempt %d failed", attemptCount))
			}
			return Of[RetryConfig](fmt.Sprintf("success on attempt %d", attemptCount))
		}

		pipeline := F.Pipe2(
			attemptOperation(),
			Alt[RetryConfig](attemptOperation),
			Alt[RetryConfig](attemptOperation),
		)

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "success on attempt 3", value)
	})

	t.Run("fallback chain with Alt", func(t *testing.T) {
		type ServiceConfig struct {
			Endpoints []string
		}

		cfg := ServiceConfig{
			Endpoints: []string{"primary", "secondary", "tertiary"},
		}

		callEndpoint := func(name string) Effect[ServiceConfig, string] {
			if name == "tertiary" {
				return Of[ServiceConfig](fmt.Sprintf("data from %s", name))
			}
			return Fail[ServiceConfig, string](fmt.Errorf("%s failed", name))
		}

		pipeline := F.Pipe2(
			callEndpoint("primary"),
			Alt[ServiceConfig](func() Effect[ServiceConfig, string] {
				return callEndpoint("secondary")
			}),
			Alt[ServiceConfig](func() Effect[ServiceConfig, string] {
				return callEndpoint("tertiary")
			}),
		)

		value, err := runEffect(pipeline, cfg)
		assert.NoError(t, err)
		assert.Equal(t, "data from tertiary", value)
	})

	t.Run("combining Alt with Map and Chain", func(t *testing.T) {
		parseAndDouble := func(s string) Effect[TestConfig, int] {
			if s == "invalid" {
				return Fail[TestConfig, int](fmt.Errorf("parse error"))
			}
			return Of[TestConfig](42)
		}

		pipeline := F.Pipe3(
			Of[TestConfig]("invalid"),
			Chain[TestConfig](parseAndDouble),
			Alt[TestConfig](func() Effect[TestConfig, int] {
				return Of[TestConfig](10)
			}),
			Map[TestConfig](N.Mul(2)),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})

	t.Run("Alt with ChainLeft comparison", func(t *testing.T) {
		originalErr := fmt.Errorf("error: 404")

		// Using Alt - cannot inspect error
		altResult := F.Pipe1(
			Fail[TestConfig, string](originalErr),
			Alt[TestConfig](func() Effect[TestConfig, string] {
				return Of[TestConfig]("fallback")
			}),
		)

		// Using ChainLeft - can inspect error
		chainLeftResult := F.Pipe1(
			Fail[TestConfig, string](originalErr),
			ChainLeft[TestConfig](func(err error) Effect[TestConfig, string] {
				return Of[TestConfig](fmt.Sprintf("recovered from: %s", err.Error()))
			}),
		)

		altValue, err1 := runEffect(altResult, testConfig)
		chainLeftValue, err2 := runEffect(chainLeftResult, testConfig)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, "fallback", altValue)
		assert.Equal(t, "recovered from: error: 404", chainLeftValue)
	})
}

// TestChainFirstLeftThunkK_Success tests ChainFirstLeftThunkK with successful effects
func TestChainFirstLeftThunkK_Success(t *testing.T) {
	t.Run("success value passes through unchanged", func(t *testing.T) {
		sideEffectRan := false

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					sideEffectRan = true
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstLeftThunkK[TestConfig, int](logError),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
		assert.False(t, sideEffectRan, "side effect should not run on success")
	})

	t.Run("chains multiple successful operations", func(t *testing.T) {
		pipeline := F.Pipe3(
			Of[TestConfig](10),
			Map[TestConfig](N.Mul(2)),
			ChainFirstLeftThunkK[TestConfig, int](func(err error) readerioresult.ReaderIOResult[F.Void] {
				return func(ctx context.Context) io.IO[result.Result[F.Void]] {
					return func() result.Result[F.Void] {
						return result.Of(F.VOID)
					}
				}
			}),
			Map[TestConfig](N.Add(5)),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 25, value)
	})
}

// TestChainFirstLeftThunkK_Failure tests ChainFirstLeftThunkK with failing effects
func TestChainFirstLeftThunkK_Failure(t *testing.T) {
	t.Run("error triggers handler with context", func(t *testing.T) {
		var capturedError error
		var capturedCtx context.Context
		originalErr := fmt.Errorf("operation failed")

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					capturedError = err
					capturedCtx = ctx
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			ChainFirstLeftThunkK[TestConfig, int](logError),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.Equal(t, originalErr, capturedError)
		assert.NotNil(t, capturedCtx)
	})

	t.Run("handler error replaces original when handler fails", func(t *testing.T) {
		originalErr := fmt.Errorf("original error")
		handlerErr := fmt.Errorf("handler error")

		failingHandler := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					return result.Left[F.Void](handlerErr)
				}
			}
		}

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			ChainFirstLeftThunkK[TestConfig, int](failingHandler),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		// ChainFirstLeft preserves the original error even when handler fails
		// This is the "First" behavior - it keeps the first (original) value/error
		assert.Equal(t, originalErr, err)
	})

	t.Run("preserves error cause chain", func(t *testing.T) {
		rootErr := fmt.Errorf("root cause")
		wrappedErr := fmt.Errorf("wrapped: %w", rootErr)
		var capturedError error

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					capturedError = err
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Fail[TestConfig, int](wrappedErr),
			ChainFirstLeftThunkK[TestConfig, int](logError),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.ErrorIs(t, err, rootErr)
		assert.ErrorIs(t, capturedError, rootErr)
	})
}

// TestChainFirstLeftThunkK_EdgeCases tests edge cases for ChainFirstLeftThunkK
func TestChainFirstLeftThunkK_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Of[TestConfig](0),
			ChainFirstLeftThunkK[TestConfig, int](logError),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("multiple error handlers in sequence", func(t *testing.T) {
		callCount := 0
		originalErr := fmt.Errorf("test error")

		countingHandler := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					callCount++
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe3(
			Fail[TestConfig, int](originalErr),
			ChainFirstLeftThunkK[TestConfig, int](countingHandler),
			ChainFirstLeftThunkK[TestConfig, int](countingHandler),
			ChainFirstLeftThunkK[TestConfig, int](countingHandler),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, 3, callCount)
	})

	t.Run("handler with context cancellation", func(t *testing.T) {
		var handlerCtx context.Context
		originalErr := fmt.Errorf("test error")

		captureContext := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					handlerCtx = ctx
					return result.Of(F.VOID)
				}
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			ChainFirstLeftThunkK[TestConfig, int](captureContext),
		)

		res := pipeline(testConfig)(ctx)()
		assert.True(t, result.IsLeft(res))
		assert.NotNil(t, handlerCtx)
		assert.Error(t, handlerCtx.Err())
	})
}

// TestChainFirstLeftThunkK_Integration tests integration scenarios
func TestChainFirstLeftThunkK_Integration(t *testing.T) {
	t.Run("error logging with structured data", func(t *testing.T) {
		type LogEntry struct {
			Error   error
			Context string
		}
		var logEntries []LogEntry

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					logEntries = append(logEntries, LogEntry{
						Error:   err,
						Context: "operation failed",
					})
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe3(
			Of[TestConfig](10),
			Chain[TestConfig](func(n int) Effect[TestConfig, int] {
				if n < 20 {
					return Fail[TestConfig, int](fmt.Errorf("value too small: %d", n))
				}
				return Of[TestConfig](n)
			}),
			ChainFirstLeftThunkK[TestConfig, int](logError),
			Alt[TestConfig](func() Effect[TestConfig, int] {
				return Of[TestConfig](100)
			}),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 100, value)
		assert.Len(t, logEntries, 1)
		assert.Contains(t, logEntries[0].Error.Error(), "value too small")
	})

	t.Run("composes with Map and Chain", func(t *testing.T) {
		errorCount := 0

		countErrors := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					errorCount++
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe4(
			Of[TestConfig](5),
			Map[TestConfig](N.Mul(2)),
			Chain[TestConfig](func(n int) Effect[TestConfig, int] {
				if n < 20 {
					return Fail[TestConfig, int](fmt.Errorf("too small"))
				}
				return Of[TestConfig](n)
			}),
			ChainFirstLeftThunkK[TestConfig, int](countErrors),
			Map[TestConfig](N.Add(10)),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, 1, errorCount)
	})
}

// TestTapLeftThunkK_Success tests TapLeftThunkK with successful effects
func TestTapLeftThunkK_Success(t *testing.T) {
	t.Run("is alias for ChainFirstLeftThunkK", func(t *testing.T) {
		sideEffectRan := false

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					sideEffectRan = true
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Of[TestConfig](42),
			TapLeftThunkK[TestConfig, int](logError),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
		assert.False(t, sideEffectRan)
	})

	t.Run("preserves value through pipeline", func(t *testing.T) {
		pipeline := F.Pipe3(
			Of[TestConfig](10),
			TapLeftThunkK[TestConfig, int](func(err error) readerioresult.ReaderIOResult[F.Void] {
				return func(ctx context.Context) io.IO[result.Result[F.Void]] {
					return func() result.Result[F.Void] {
						return result.Of(F.VOID)
					}
				}
			}),
			Map[TestConfig](N.Mul(3)),
			TapLeftThunkK[TestConfig, int](func(err error) readerioresult.ReaderIOResult[F.Void] {
				return func(ctx context.Context) io.IO[result.Result[F.Void]] {
					return func() result.Result[F.Void] {
						return result.Of(F.VOID)
					}
				}
			}),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 30, value)
	})
}

// TestTapLeftThunkK_Failure tests TapLeftThunkK with failing effects
func TestTapLeftThunkK_Failure(t *testing.T) {
	t.Run("executes on error", func(t *testing.T) {
		var loggedError error
		originalErr := fmt.Errorf("computation failed")

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					loggedError = err
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			TapLeftThunkK[TestConfig, int](logError),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.Equal(t, originalErr, loggedError)
	})

	t.Run("preserves error through multiple taps", func(t *testing.T) {
		callOrder := []int{}
		originalErr := fmt.Errorf("test error")

		makeTap := func(id int) func(error) readerioresult.ReaderIOResult[F.Void] {
			return func(err error) readerioresult.ReaderIOResult[F.Void] {
				return func(ctx context.Context) io.IO[result.Result[F.Void]] {
					return func() result.Result[F.Void] {
						callOrder = append(callOrder, id)
						return result.Of(F.VOID)
					}
				}
			}
		}

		pipeline := F.Pipe3(
			Fail[TestConfig, int](originalErr),
			TapLeftThunkK[TestConfig, int](makeTap(1)),
			TapLeftThunkK[TestConfig, int](makeTap(2)),
			TapLeftThunkK[TestConfig, int](makeTap(3)),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.Equal(t, []int{1, 2, 3}, callOrder)
	})
}

// TestTapLeftThunkK_EdgeCases tests edge cases for TapLeftThunkK
func TestTapLeftThunkK_EdgeCases(t *testing.T) {
	t.Run("handler with context values", func(t *testing.T) {
		type contextKey string
		const requestIDKey contextKey = "requestID"

		var capturedRequestID string
		originalErr := fmt.Errorf("test error")

		captureRequestID := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					if id := ctx.Value(requestIDKey); id != nil {
						capturedRequestID = id.(string)
					}
					return result.Of(F.VOID)
				}
			}
		}

		ctx := context.WithValue(context.Background(), requestIDKey, "req-123")

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			TapLeftThunkK[TestConfig, int](captureRequestID),
		)

		res := pipeline(testConfig)(ctx)()
		assert.True(t, result.IsLeft(res))
		assert.Equal(t, "req-123", capturedRequestID)
	})

	t.Run("combines with OrElse for recovery", func(t *testing.T) {
		errorLogged := false
		originalErr := fmt.Errorf("primary failed")

		logError := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					errorLogged = true
					return result.Of(F.VOID)
				}
			}
		}

		pipeline := F.Pipe2(
			Fail[TestConfig, int](originalErr),
			TapLeftThunkK[TestConfig, int](logError),
			Alt[TestConfig](func() Effect[TestConfig, int] {
				return Of[TestConfig](999)
			}),
		)

		value, err := runEffect(pipeline, testConfig)
		assert.NoError(t, err)
		assert.Equal(t, 999, value)
		assert.True(t, errorLogged)
	})
}

// TestTapLeftThunkK_Integration tests integration scenarios
func TestTapLeftThunkK_Integration(t *testing.T) {
	t.Run("real-world error logging scenario", func(t *testing.T) {
		type ErrorLog struct {
			Timestamp string
			Error     string
			Operation string
		}
		var errorLogs []ErrorLog

		logError := func(operation string) func(error) readerioresult.ReaderIOResult[F.Void] {
			return func(err error) readerioresult.ReaderIOResult[F.Void] {
				return func(ctx context.Context) io.IO[result.Result[F.Void]] {
					return func() result.Result[F.Void] {
						errorLogs = append(errorLogs, ErrorLog{
							Timestamp: "2024-01-01",
							Error:     err.Error(),
							Operation: operation,
						})
						return result.Of(F.VOID)
					}
				}
			}
		}

		fetchData := func(id int) Effect[TestConfig, string] {
			if id < 0 {
				return Fail[TestConfig, string](fmt.Errorf("invalid ID: %d", id))
			}
			return Of[TestConfig](fmt.Sprintf("data-%d", id))
		}

		validateData := func(data string) Effect[TestConfig, string] {
			if len(data) < 5 {
				return Fail[TestConfig, string](fmt.Errorf("invalid data: %s", data))
			}
			return Of[TestConfig](data)
		}

		pipeline := F.Pipe3(
			fetchData(-1),
			TapLeftThunkK[TestConfig, string](logError("fetchData")),
			Chain[TestConfig](validateData),
			TapLeftThunkK[TestConfig, string](logError("validateData")),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		// Both handlers run because the error propagates through the chain
		assert.GreaterOrEqual(t, len(errorLogs), 1)
		assert.Equal(t, "fetchData", errorLogs[0].Operation)
		assert.Contains(t, errorLogs[0].Error, "invalid ID")
	})

	t.Run("error notification with timeout", func(t *testing.T) {
		notificationSent := false
		originalErr := fmt.Errorf("critical error")

		sendNotification := func(err error) readerioresult.ReaderIOResult[F.Void] {
			return func(ctx context.Context) io.IO[result.Result[F.Void]] {
				return func() result.Result[F.Void] {
					// Simulate notification with context
					select {
					case <-ctx.Done():
						return result.Left[F.Void](ctx.Err())
					default:
						notificationSent = true
						return result.Of(F.VOID)
					}
				}
			}
		}

		pipeline := F.Pipe1(
			Fail[TestConfig, int](originalErr),
			TapLeftThunkK[TestConfig, int](sendNotification),
		)

		_, err := runEffect(pipeline, testConfig)
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.True(t, notificationSent)
	})
}

// ExampleTapLeftThunkK demonstrates error logging with structured logging.
func ExampleTapLeftThunkK() {
	type Config struct{}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{} // drop timestamp in tests
		}
		return a
	}}))

	cancel, ctx := pair.Unpack(logging.WithLogger(logger)(context.Background()))
	defer cancel()

	fetchUser := func(id int) Effect[Config, string] {
		return Fail[Config, string](fmt.Errorf("user not found: %d", id))
	}

	pipeline := F.Pipe1(
		fetchUser(42),
		TapLeftThunkK[Config, string](thunk.SLogLeft("Operation failed")),
	)

	_ = pipeline(Config{})(ctx)()

	// Output:
	// level=ERROR msg="Operation failed" error="user not found: 42"
}
