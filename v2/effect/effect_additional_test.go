// Copyright (c) 2023 - 2025 IBM Corp.
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

package effect

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestSucceed tests the Succeed function
func TestSucceed_Success(t *testing.T) {
	t.Run("creates successful effect with int", func(t *testing.T) {
		eff := Succeed[TestConfig](42)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("creates successful effect with string", func(t *testing.T) {
		eff := Succeed[TestConfig]("hello")
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("hello"), outcome)
	})

	t.Run("creates successful effect with zero value", func(t *testing.T) {
		eff := Succeed[TestConfig](0)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(0), outcome)
	})
}

// TestFail tests the Fail function
func TestFail_Failure(t *testing.T) {
	t.Run("creates failed effect with error", func(t *testing.T) {
		testErr := errors.New("test error")
		eff := Fail[TestConfig, int](testErr)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("preserves error message", func(t *testing.T) {
		testErr := errors.New("specific error message")
		eff := Fail[TestConfig, string](testErr)
		outcome := eff(testConfig)(context.Background())()
		assert.True(t, result.IsLeft(outcome))
		extractedErr := result.MonadFold(outcome,
			F.Identity[error],
			func(string) error { return nil },
		)
		assert.Equal(t, testErr, extractedErr)
	})
}

// TestOf tests the Of function
func TestOf_Success(t *testing.T) {
	t.Run("creates successful effect with value", func(t *testing.T) {
		eff := Of[TestConfig](100)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(100), outcome)
	})

	t.Run("is equivalent to Succeed", func(t *testing.T) {
		value := "test"
		eff1 := Of[TestConfig](value)
		eff2 := Succeed[TestConfig](value)
		outcome1 := eff1(testConfig)(context.Background())()
		outcome2 := eff2(testConfig)(context.Background())()
		assert.Equal(t, outcome1, outcome2)
	})
}

// TestMap tests the Map function
func TestMap_Success(t *testing.T) {
	t.Run("transforms success value", func(t *testing.T) {
		eff := F.Pipe1(
			Of[TestConfig](42),
			Map[TestConfig](func(x int) int { return x * 2 }),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(84), outcome)
	})

	t.Run("transforms type", func(t *testing.T) {
		eff := F.Pipe1(
			Of[TestConfig](42),
			Map[TestConfig](func(x int) string { return strconv.Itoa(x) }),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("42"), outcome)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		eff := F.Pipe2(
			Of[TestConfig](10),
			Map[TestConfig](func(x int) int { return x + 5 }),
			Map[TestConfig](func(x int) int { return x * 2 }),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
	})
}

func TestMap_Failure(t *testing.T) {
	t.Run("propagates error unchanged", func(t *testing.T) {
		testErr := errors.New("test error")
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			Map[TestConfig](func(x int) int { return x * 2 }),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}

// TestChain tests the Chain function
func TestChain_Success(t *testing.T) {
	t.Run("sequences two effects", func(t *testing.T) {
		eff := F.Pipe1(
			Of[TestConfig](42),
			Chain(func(x int) Effect[TestConfig, string] {
				return Of[TestConfig](strconv.Itoa(x))
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("42"), outcome)
	})

	t.Run("chains multiple effects", func(t *testing.T) {
		eff := F.Pipe2(
			Of[TestConfig](10),
			Chain(func(x int) Effect[TestConfig, int] {
				return Of[TestConfig](x + 5)
			}),
			Chain(func(x int) Effect[TestConfig, int] {
				return Of[TestConfig](x * 2)
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
	})
}

func TestChain_Failure(t *testing.T) {
	t.Run("propagates error from first effect", func(t *testing.T) {
		testErr := errors.New("first error")
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			Chain(func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("should not execute")
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[string](testErr), outcome)
	})

	t.Run("propagates error from second effect", func(t *testing.T) {
		testErr := errors.New("second error")
		eff := F.Pipe1(
			Of[TestConfig](42),
			Chain(func(x int) Effect[TestConfig, string] {
				return Fail[TestConfig, string](testErr)
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[string](testErr), outcome)
	})
}

// TestChainIOK tests the ChainIOK function
func TestChainIOK_Success(t *testing.T) {
	t.Run("chains with IO action", func(t *testing.T) {
		counter := 0
		eff := F.Pipe1(
			Of[TestConfig](42),
			ChainIOK[TestConfig](func(x int) io.IO[string] {
				return func() string {
					counter++
					return fmt.Sprintf("Value: %d", x)
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("Value: 42"), outcome)
		assert.Equal(t, 1, counter)
	})

	t.Run("chains multiple IO actions", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe2(
			Of[TestConfig](10),
			ChainIOK[TestConfig](func(x int) io.IO[int] {
				return func() int {
					log = append(log, "first")
					return x + 5
				}
			}),
			ChainIOK[TestConfig](func(x int) io.IO[int] {
				return func() int {
					log = append(log, "second")
					return x * 2
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
		assert.Equal(t, []string{"first", "second"}, log)
	})
}

func TestChainIOK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := errors.New("test error")
		executed := false
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainIOK[TestConfig](func(x int) io.IO[string] {
				return func() string {
					executed = true
					return "should not execute"
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[string](testErr), outcome)
		assert.False(t, executed)
	})
}

// TestChainFirstIOK tests the ChainFirstIOK function
func TestChainFirstIOK_Success(t *testing.T) {
	t.Run("executes IO but preserves value", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe1(
			Of[TestConfig](42),
			ChainFirstIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log = append(log, fmt.Sprintf("logged: %d", x))
					return nil
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, []string{"logged: 42"}, log)
	})

	t.Run("chains multiple side effects", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe2(
			Of[TestConfig](10),
			ChainFirstIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log = append(log, "first")
					return nil
				}
			}),
			ChainFirstIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log = append(log, "second")
					return nil
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(10), outcome)
		assert.Equal(t, []string{"first", "second"}, log)
	})
}

func TestChainFirstIOK_Failure(t *testing.T) {
	t.Run("propagates error without executing IO", func(t *testing.T) {
		testErr := errors.New("test error")
		executed := false
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			ChainFirstIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					executed = true
					return nil
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, executed)
	})
}

// TestTapIOK tests the TapIOK function
func TestTapIOK_Success(t *testing.T) {
	t.Run("executes IO but preserves value", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe1(
			Of[TestConfig](42),
			TapIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log = append(log, fmt.Sprintf("tapped: %d", x))
					return nil
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, []string{"tapped: 42"}, log)
	})

	t.Run("is equivalent to ChainFirstIOK", func(t *testing.T) {
		log1 := []string{}
		log2 := []string{}

		eff1 := F.Pipe1(
			Of[TestConfig](10),
			TapIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log1 = append(log1, "tap")
					return nil
				}
			}),
		)

		eff2 := F.Pipe1(
			Of[TestConfig](10),
			ChainFirstIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					log2 = append(log2, "tap")
					return nil
				}
			}),
		)

		outcome1 := eff1(testConfig)(context.Background())()
		outcome2 := eff2(testConfig)(context.Background())()

		assert.Equal(t, outcome1, outcome2)
		assert.Equal(t, log1, log2)
	})
}

func TestTapIOK_Failure(t *testing.T) {
	t.Run("propagates error without executing IO", func(t *testing.T) {
		testErr := errors.New("test error")
		executed := false
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			TapIOK[TestConfig](func(x int) io.IO[any] {
				return func() any {
					executed = true
					return nil
				}
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, executed)
	})
}

// TestChainResultK tests the ChainResultK function
func TestChainResultK_Success(t *testing.T) {
	t.Run("chains with Result-returning function", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := F.Pipe1(
			Of[TestConfig]("42"),
			ChainResultK[TestConfig](parseIntResult),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("chains multiple Result operations", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := F.Pipe2(
			Of[TestConfig]("10"),
			ChainResultK[TestConfig](parseIntResult),
			ChainResultK[TestConfig](func(x int) result.Result[string] {
				return result.Of(fmt.Sprintf("Value: %d", x*2))
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("Value: 20"), outcome)
	})
}

func TestChainResultK_Failure(t *testing.T) {
	t.Run("propagates error from previous effect", func(t *testing.T) {
		testErr := errors.New("test error")
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := F.Pipe1(
			Fail[TestConfig, string](testErr),
			ChainResultK[TestConfig](parseIntResult),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("propagates error from Result function", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := F.Pipe1(
			Of[TestConfig]("not a number"),
			ChainResultK[TestConfig](parseIntResult),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.True(t, result.IsLeft(outcome))
	})
}

// TestAp tests the Ap function
func TestAp_Success(t *testing.T) {
	t.Run("applies function effect to value effect", func(t *testing.T) {
		fnEff := Of[TestConfig](func(x int) int { return x * 2 })
		valEff := Of[TestConfig](21)
		eff := Ap[int](valEff)(fnEff)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("applies function with different types", func(t *testing.T) {
		fnEff := Of[TestConfig](func(x int) string { return strconv.Itoa(x) })
		valEff := Of[TestConfig](42)
		eff := Ap[string](valEff)(fnEff)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of("42"), outcome)
	})
}

func TestAp_Failure(t *testing.T) {
	t.Run("propagates error from function effect", func(t *testing.T) {
		testErr := errors.New("function error")
		fnEff := Fail[TestConfig, func(int) int](testErr)
		valEff := Of[TestConfig](42)
		eff := Ap[int](valEff)(fnEff)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("propagates error from value effect", func(t *testing.T) {
		testErr := errors.New("value error")
		fnEff := Of[TestConfig](func(x int) int { return x * 2 })
		valEff := Fail[TestConfig, int](testErr)
		eff := Ap[int](valEff)(fnEff)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}

// TestSuspend tests the Suspend function
func TestSuspend_Success(t *testing.T) {
	t.Run("delays evaluation of effect", func(t *testing.T) {
		counter := 0
		eff := Suspend(func() Effect[TestConfig, int] {
			counter++
			return Of[TestConfig](42)
		})
		assert.Equal(t, 0, counter, "should not evaluate immediately")
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, 1, counter, "should evaluate when run")
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("enables recursive effects", func(t *testing.T) {
		var factorial func(int) Effect[TestConfig, int]
		factorial = func(n int) Effect[TestConfig, int] {
			if n <= 1 {
				return Of[TestConfig](1)
			}
			return Suspend(func() Effect[TestConfig, int] {
				return F.Pipe1(
					factorial(n-1),
					Map[TestConfig](func(x int) int { return x * n }),
				)
			})
		}

		outcome := factorial(5)(testConfig)(context.Background())()
		assert.Equal(t, result.Of(120), outcome)
	})
}

// TestTap tests the Tap function
func TestTap_Success(t *testing.T) {
	t.Run("executes side effect but preserves value", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe1(
			Of[TestConfig](42),
			Tap(func(x int) Effect[TestConfig, any] {
				log = append(log, fmt.Sprintf("tapped: %d", x))
				return Of[TestConfig, any](nil)
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, []string{"tapped: 42"}, log)
	})

	t.Run("chains multiple taps", func(t *testing.T) {
		log := []string{}
		eff := F.Pipe2(
			Of[TestConfig](10),
			Tap(func(x int) Effect[TestConfig, any] {
				log = append(log, "first")
				return Of[TestConfig, any](nil)
			}),
			Tap(func(x int) Effect[TestConfig, any] {
				log = append(log, "second")
				return Of[TestConfig, any](nil)
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Of(10), outcome)
		assert.Equal(t, []string{"first", "second"}, log)
	})
}

func TestTap_Failure(t *testing.T) {
	t.Run("propagates error without executing tap", func(t *testing.T) {
		testErr := errors.New("test error")
		executed := false
		eff := F.Pipe1(
			Fail[TestConfig, int](testErr),
			Tap(func(x int) Effect[TestConfig, any] {
				executed = true
				return Of[TestConfig, any](nil)
			}),
		)
		outcome := eff(testConfig)(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
		assert.False(t, executed)
	})
}

// TestTernary tests the Ternary function
func TestTernary_Success(t *testing.T) {
	t.Run("executes onTrue when predicate is true", func(t *testing.T) {
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("large")
			},
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("small")
			},
		)
		outcome := kleisli(15)(testConfig)(context.Background())()
		assert.Equal(t, result.Of("large"), outcome)
	})

	t.Run("executes onFalse when predicate is false", func(t *testing.T) {
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("large")
			},
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("small")
			},
		)
		outcome := kleisli(5)(testConfig)(context.Background())()
		assert.Equal(t, result.Of("small"), outcome)
	})

	t.Run("works with boundary value", func(t *testing.T) {
		kleisli := Ternary(
			func(x int) bool { return x >= 10 },
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("gte")
			},
			func(x int) Effect[TestConfig, string] {
				return Of[TestConfig]("lt")
			},
		)
		outcome := kleisli(10)(testConfig)(context.Background())()
		assert.Equal(t, result.Of("gte"), outcome)
	})
}

// TestRead tests the Read function
func TestRead_Success(t *testing.T) {
	t.Run("provides context to effect", func(t *testing.T) {
		eff := Of[TestConfig](42)
		thunk := Read[int](testConfig)(eff)
		outcome := thunk(context.Background())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("converts effect to thunk", func(t *testing.T) {
		eff := F.Pipe1(
			Of[TestConfig](10),
			Map[TestConfig](func(x int) int { return x * testConfig.Multiplier }),
		)
		thunk := Read[int](testConfig)(eff)
		outcome := thunk(context.Background())()
		assert.Equal(t, result.Of(30), outcome)
	})

	t.Run("works with different contexts", func(t *testing.T) {
		cfg1 := TestConfig{Multiplier: 2, Prefix: "A", DatabaseURL: ""}
		cfg2 := TestConfig{Multiplier: 5, Prefix: "B", DatabaseURL: ""}

		// Create an effect that uses the context's Multiplier
		eff := F.Pipe1(
			Of[TestConfig](10),
			ChainReaderK(func(x int) reader.Reader[TestConfig, int] {
				return func(cfg TestConfig) int {
					return x * cfg.Multiplier
				}
			}),
		)

		thunk1 := Read[int](cfg1)(eff)
		thunk2 := Read[int](cfg2)(eff)

		outcome1 := thunk1(context.Background())()
		outcome2 := thunk2(context.Background())()

		assert.Equal(t, result.Of(20), outcome1)
		assert.Equal(t, result.Of(50), outcome2)
	})
}

func TestRead_Failure(t *testing.T) {
	t.Run("propagates error from effect", func(t *testing.T) {
		testErr := errors.New("test error")
		eff := Fail[TestConfig, int](testErr)
		thunk := Read[int](testConfig)(eff)
		outcome := thunk(context.Background())()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})
}
