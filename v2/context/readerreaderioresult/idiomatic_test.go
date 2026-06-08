// Copyright (c) 2024 IBM Corp.
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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

var idiomaticTestErr = errors.New("idiomatic test error")

// double is a simple idiomatic Kleisli that doubles an int.
func doubleIdiomatic(n int) func(context.Context, AppConfig) (int, error) {
	return func(ctx context.Context, cfg AppConfig) (int, error) {
		return n * 2, nil
	}
}

// failIdiomatic always returns an error.
func failIdiomatic(n int) func(context.Context, AppConfig) (int, error) {
	return func(ctx context.Context, cfg AppConfig) (int, error) {
		return 0, idiomaticTestErr
	}
}

func TestFromIdiomatic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		kleisli := FromIdiomatic(doubleIdiomatic)
		outcome := kleisli(5)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error", func(t *testing.T) {
		kleisli := FromIdiomatic(failIdiomatic)
		outcome := kleisli(5)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Left[int](idiomaticTestErr), outcome)
	})

	t.Run("accesses config", func(t *testing.T) {
		f := func(n int) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				return n + len(cfg.LogLevel), nil
			}
		}
		kleisli := FromIdiomatic(f)
		outcome := kleisli(10)(defaultConfig)(t.Context())()
		// 10 + len("info") = 14
		assert.Equal(t, result.Of(14), outcome)
	})
}

func TestMonadChainI(t *testing.T) {
	t.Run("success chains value", func(t *testing.T) {
		fa := Of[AppConfig](5)
		outcome := MonadChainI(fa, doubleIdiomatic)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error in fa short-circuits", func(t *testing.T) {
		fa := Left[AppConfig, int](idiomaticTestErr)
		outcome := MonadChainI(fa, doubleIdiomatic)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})

	t.Run("error in f propagates", func(t *testing.T) {
		fa := Of[AppConfig](5)
		outcome := MonadChainI(fa, failIdiomatic)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Left[int](idiomaticTestErr), outcome)
	})
}

func TestMonadChainFirstI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			sideEffectRan = true
			return n * 100, nil
		}
	}

	t.Run("returns original value", func(t *testing.T) {
		sideEffectRan = false
		fa := Of[AppConfig](7)
		outcome := MonadChainFirstI(fa, sideEffect)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(7), outcome)
		assert.True(t, sideEffectRan)
	})

	t.Run("error in fa short-circuits", func(t *testing.T) {
		sideEffectRan = false
		fa := Left[AppConfig, int](idiomaticTestErr)
		outcome := MonadChainFirstI(fa, sideEffect)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
		assert.False(t, sideEffectRan)
	})

	t.Run("error in f propagates", func(t *testing.T) {
		fa := Of[AppConfig](7)
		outcome := MonadChainFirstI(fa, failIdiomatic)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestMonadTapI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, AppConfig) (string, error) {
		return func(ctx context.Context, cfg AppConfig) (string, error) {
			sideEffectRan = true
			return "logged", nil
		}
	}

	t.Run("returns original value after side effect", func(t *testing.T) {
		sideEffectRan = false
		fa := Of[AppConfig](42)
		outcome := MonadTapI(fa, sideEffect)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
		assert.True(t, sideEffectRan)
	})

	t.Run("error in fa skips side effect", func(t *testing.T) {
		sideEffectRan = false
		fa := Left[AppConfig, int](idiomaticTestErr)
		outcome := MonadTapI(fa, sideEffect)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
		assert.False(t, sideEffectRan)
	})
}

func TestChainI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		outcome := F.Pipe1(
			Of[AppConfig](5),
			ChainI(doubleIdiomatic),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error propagates", func(t *testing.T) {
		outcome := F.Pipe1(
			Of[AppConfig](5),
			ChainI(failIdiomatic),
		)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestChainFirstI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			sideEffectRan = true
			return n * 100, nil
		}
	}

	t.Run("returns original value", func(t *testing.T) {
		sideEffectRan = false
		outcome := F.Pipe1(
			Of[AppConfig](3),
			ChainFirstI(sideEffect),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(3), outcome)
		assert.True(t, sideEffectRan)
	})

	t.Run("error propagates from f", func(t *testing.T) {
		outcome := F.Pipe1(
			Of[AppConfig](3),
			ChainFirstI(failIdiomatic),
		)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestTapI(t *testing.T) {
	logged := ""
	logF := func(n int) func(context.Context, AppConfig) (string, error) {
		return func(ctx context.Context, cfg AppConfig) (string, error) {
			logged = cfg.LogLevel
			return logged, nil
		}
	}

	t.Run("passes through original value", func(t *testing.T) {
		logged = ""
		outcome := F.Pipe1(
			Of[AppConfig](99),
			TapI(logF),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(99), outcome)
		assert.Equal(t, "info", logged)
	})
}

func TestMonadChainLeftI(t *testing.T) {
	recover := func(err error) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			return -1, nil
		}
	}

	t.Run("success value passes through", func(t *testing.T) {
		fa := Of[AppConfig](42)
		outcome := MonadChainLeftI(fa, recover)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("error triggers recovery", func(t *testing.T) {
		fa := Left[AppConfig, int](idiomaticTestErr)
		outcome := MonadChainLeftI(fa, recover)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(-1), outcome)
	})

	t.Run("recovery can also fail", func(t *testing.T) {
		recoveryErr := errors.New("recovery failed")
		failRecover := func(err error) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				return 0, recoveryErr
			}
		}
		fa := Left[AppConfig, int](idiomaticTestErr)
		outcome := MonadChainLeftI(fa, failRecover)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Left[int](recoveryErr), outcome)
	})
}

func TestChainLeftI(t *testing.T) {
	recover := func(err error) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			return 0, nil
		}
	}

	t.Run("success passes through", func(t *testing.T) {
		outcome := F.Pipe1(
			Of[AppConfig](42),
			ChainLeftI(recover),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
	})

	t.Run("error triggers recovery", func(t *testing.T) {
		outcome := F.Pipe1(
			Left[AppConfig, int](idiomaticTestErr),
			ChainLeftI(recover),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(0), outcome)
	})
}

func TestRetryingI(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				attempts++
				return 42, nil
			}
		}
		outcome := RetryingI(retry.LimitRetries(3), action, result.IsLeft[int])(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, 1, attempts)
	})

	t.Run("retries until success", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				attempts++
				if attempts < 3 {
					return 0, idiomaticTestErr
				}
				return 42, nil
			}
		}
		outcome := RetryingI(retry.LimitRetries(5), action, result.IsLeft[int])(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(42), outcome)
		assert.Equal(t, 3, attempts)
	})

	t.Run("exhausts retries", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				attempts++
				return 0, idiomaticTestErr
			}
		}
		outcome := RetryingI(retry.LimitRetries(2), action, result.IsLeft[int])(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
		assert.Equal(t, 3, attempts) // initial + 2 retries
	})
}

func TestTraverseArrayI(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		outcome := TraverseArrayI(doubleIdiomatic)([]int{1, 2, 3})(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of([]int{2, 4, 6}), outcome)
	})

	t.Run("one element fails", func(t *testing.T) {
		f := func(n int) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				if n == 2 {
					return 0, idiomaticTestErr
				}
				return n * 2, nil
			}
		}
		outcome := TraverseArrayI(f)([]int{1, 2, 3})(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})

	t.Run("empty slice", func(t *testing.T) {
		outcome := TraverseArrayI(doubleIdiomatic)([]int{})(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of([]int{}), outcome)
	})
}

func TestBindI(t *testing.T) {
	type State struct {
		Value int
	}

	setter := func(v int) func(State) State {
		return func(s State) State {
			s.Value = v
			return s
		}
	}

	f := func(s State) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			return s.Value * 2, nil
		}
	}

	t.Run("success", func(t *testing.T) {
		outcome := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindI(setter, f),
			Map[AppConfig](func(s State) int { return s.Value }),
		)(defaultConfig)(t.Context())()
		assert.Equal(t, result.Of(10), outcome)
	})

	t.Run("error propagates", func(t *testing.T) {
		fErr := func(s State) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				return 0, idiomaticTestErr
			}
		}
		outcome := F.Pipe2(
			Do[AppConfig](State{Value: 5}),
			BindI(setter, fErr),
			Map[AppConfig](func(s State) int { return s.Value }),
		)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}

func TestBindIL(t *testing.T) {
	type State struct {
		Count int
	}

	countLens := lens.MakeLens(
		func(s State) int { return s.Count },
		func(s State, v int) State {
			s.Count = v
			return s
		},
	)

	f := func(n int) func(context.Context, AppConfig) (int, error) {
		return func(ctx context.Context, cfg AppConfig) (int, error) {
			return n + len(cfg.LogLevel), nil
		}
	}

	t.Run("success updates via lens", func(t *testing.T) {
		outcome := F.Pipe2(
			Do[AppConfig](State{Count: 10}),
			BindIL(countLens, f),
			Map[AppConfig](func(s State) int { return s.Count }),
		)(defaultConfig)(t.Context())()
		// 10 + len("info") = 14
		assert.Equal(t, result.Of(14), outcome)
	})

	t.Run("error propagates", func(t *testing.T) {
		fErr := func(n int) func(context.Context, AppConfig) (int, error) {
			return func(ctx context.Context, cfg AppConfig) (int, error) {
				return 0, idiomaticTestErr
			}
		}
		outcome := F.Pipe2(
			Do[AppConfig](State{Count: 10}),
			BindIL(countLens, fErr),
			Map[AppConfig](func(s State) int { return s.Count }),
		)(defaultConfig)(t.Context())()
		assert.True(t, result.IsLeft(outcome))
	})
}
