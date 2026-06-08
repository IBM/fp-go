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

package effect

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

// doubleIdiomatic is a simple idiomatic Kleisli that doubles an int.
func doubleIdiomatic(n int) func(context.Context, TestContext) (int, error) {
	return func(ctx context.Context, c TestContext) (int, error) {
		return n * 2, nil
	}
}

// failIdiomatic always returns an error.
func failIdiomatic(n int) func(context.Context, TestContext) (int, error) {
	return func(ctx context.Context, c TestContext) (int, error) {
		return 0, idiomaticTestErr
	}
}

var testCtx = TestContext{Value: "test"}

func TestFromIdiomatic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		kleisli := FromIdiomatic(doubleIdiomatic)
		v, err := runEffect(kleisli(5), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error", func(t *testing.T) {
		kleisli := FromIdiomatic(failIdiomatic)
		_, err := runEffect(kleisli(5), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})

	t.Run("accesses context value", func(t *testing.T) {
		f := func(n int) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				return n + len(c.Value), nil
			}
		}
		kleisli := FromIdiomatic(f)
		v, err := runEffect(kleisli(10), testCtx)
		assert.NoError(t, err)
		// 10 + len("test") = 14
		assert.Equal(t, 14, v)
	})
}

func TestMonadChainI(t *testing.T) {
	t.Run("success chains value", func(t *testing.T) {
		fa := Of[TestContext](5)
		v, err := runEffect(MonadChainI(fa, doubleIdiomatic), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error in fa short-circuits", func(t *testing.T) {
		fa := Fail[TestContext, int](idiomaticTestErr)
		_, err := runEffect(MonadChainI(fa, doubleIdiomatic), testCtx)
		assert.Error(t, err)
	})

	t.Run("error in f propagates", func(t *testing.T) {
		fa := Of[TestContext](5)
		_, err := runEffect(MonadChainI(fa, failIdiomatic), testCtx)
		assert.Equal(t, idiomaticTestErr, err)
	})
}

func TestMonadChainFirstI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			sideEffectRan = true
			return n * 100, nil
		}
	}

	t.Run("returns original value", func(t *testing.T) {
		sideEffectRan = false
		fa := Of[TestContext](7)
		v, err := runEffect(MonadChainFirstI(fa, sideEffect), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 7, v)
		assert.True(t, sideEffectRan)
	})

	t.Run("error in fa short-circuits", func(t *testing.T) {
		sideEffectRan = false
		fa := Fail[TestContext, int](idiomaticTestErr)
		_, err := runEffect(MonadChainFirstI(fa, sideEffect), testCtx)
		assert.Error(t, err)
		assert.False(t, sideEffectRan)
	})

	t.Run("error in f propagates", func(t *testing.T) {
		fa := Of[TestContext](7)
		_, err := runEffect(MonadChainFirstI(fa, failIdiomatic), testCtx)
		assert.Error(t, err)
	})
}

func TestMonadTapI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, TestContext) (string, error) {
		return func(ctx context.Context, c TestContext) (string, error) {
			sideEffectRan = true
			return "logged", nil
		}
	}

	t.Run("returns original value after side effect", func(t *testing.T) {
		sideEffectRan = false
		fa := Of[TestContext](42)
		v, err := runEffect(MonadTapI(fa, sideEffect), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
		assert.True(t, sideEffectRan)
	})

	t.Run("error in fa skips side effect", func(t *testing.T) {
		sideEffectRan = false
		fa := Fail[TestContext, int](idiomaticTestErr)
		_, err := runEffect(MonadTapI(fa, sideEffect), testCtx)
		assert.Error(t, err)
		assert.False(t, sideEffectRan)
	})
}

func TestChainI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext](5), ChainI(doubleIdiomatic)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error propagates", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](5), ChainI(failIdiomatic)), testCtx)
		assert.Error(t, err)
	})
}

func TestChainFirstI(t *testing.T) {
	sideEffectRan := false
	sideEffect := func(n int) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			sideEffectRan = true
			return n * 100, nil
		}
	}

	t.Run("returns original value", func(t *testing.T) {
		sideEffectRan = false
		v, err := runEffect(F.Pipe1(Of[TestContext](3), ChainFirstI(sideEffect)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 3, v)
		assert.True(t, sideEffectRan)
	})

	t.Run("error propagates from f", func(t *testing.T) {
		_, err := runEffect(F.Pipe1(Of[TestContext](3), ChainFirstI(failIdiomatic)), testCtx)
		assert.Error(t, err)
	})
}

func TestTapI(t *testing.T) {
	logged := ""
	logF := func(n int) func(context.Context, TestContext) (string, error) {
		return func(ctx context.Context, c TestContext) (string, error) {
			logged = c.Value
			return logged, nil
		}
	}

	t.Run("passes through original value", func(t *testing.T) {
		logged = ""
		v, err := runEffect(F.Pipe1(Of[TestContext](99), TapI(logF)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 99, v)
		assert.Equal(t, "test", logged)
	})
}

func TestMonadChainLeftI(t *testing.T) {
	recover := func(err error) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			return -1, nil
		}
	}

	t.Run("success value passes through", func(t *testing.T) {
		fa := Of[TestContext](42)
		v, err := runEffect(MonadChainLeftI(fa, recover), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("error triggers recovery", func(t *testing.T) {
		fa := Fail[TestContext, int](idiomaticTestErr)
		v, err := runEffect(MonadChainLeftI(fa, recover), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, -1, v)
	})

	t.Run("recovery can also fail", func(t *testing.T) {
		recoveryErr := errors.New("recovery failed")
		failRecover := func(e error) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				return 0, recoveryErr
			}
		}
		fa := Fail[TestContext, int](idiomaticTestErr)
		_, err := runEffect(MonadChainLeftI(fa, failRecover), testCtx)
		assert.Equal(t, recoveryErr, err)
	})
}

func TestChainLeftI(t *testing.T) {
	recover := func(e error) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			return 0, nil
		}
	}

	t.Run("success passes through", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Of[TestContext](42), ChainLeftI(recover)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
	})

	t.Run("error triggers recovery", func(t *testing.T) {
		v, err := runEffect(F.Pipe1(Fail[TestContext, int](idiomaticTestErr), ChainLeftI(recover)), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 0, v)
	})
}

func TestRetryingI(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				attempts++
				return 42, nil
			}
		}
		v, err := runEffect(RetryingI(retry.LimitRetries(3), action, result.IsLeft[int]), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
		assert.Equal(t, 1, attempts)
	})

	t.Run("retries until success", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				attempts++
				if attempts < 3 {
					return 0, idiomaticTestErr
				}
				return 42, nil
			}
		}
		v, err := runEffect(RetryingI(retry.LimitRetries(5), action, result.IsLeft[int]), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 42, v)
		assert.Equal(t, 3, attempts)
	})

	t.Run("exhausts retries", func(t *testing.T) {
		attempts := 0
		action := func(_ retry.RetryStatus) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				attempts++
				return 0, idiomaticTestErr
			}
		}
		_, err := runEffect(RetryingI(retry.LimitRetries(2), action, result.IsLeft[int]), testCtx)
		assert.Error(t, err)
		assert.Equal(t, 3, attempts) // initial + 2 retries
	})
}

func TestTraverseArrayI(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		v, err := runEffect(TraverseArrayI(doubleIdiomatic)([]int{1, 2, 3}), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4, 6}, v)
	})

	t.Run("one element fails", func(t *testing.T) {
		f := func(n int) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				if n == 2 {
					return 0, idiomaticTestErr
				}
				return n * 2, nil
			}
		}
		_, err := runEffect(TraverseArrayI(f)([]int{1, 2, 3}), testCtx)
		assert.Error(t, err)
	})

	t.Run("empty slice", func(t *testing.T) {
		v, err := runEffect(TraverseArrayI(doubleIdiomatic)([]int{}), testCtx)
		assert.NoError(t, err)
		assert.Equal(t, []int{}, v)
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

	f := func(s State) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			return s.Value * 2, nil
		}
	}

	t.Run("success", func(t *testing.T) {
		eff := F.Pipe2(
			Do[TestContext](State{Value: 5}),
			BindI(setter, f),
			Map[TestContext](func(s State) int { return s.Value }),
		)
		v, err := runEffect(eff, testCtx)
		assert.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("error propagates", func(t *testing.T) {
		fErr := func(s State) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				return 0, idiomaticTestErr
			}
		}
		eff := F.Pipe2(
			Do[TestContext](State{Value: 5}),
			BindI(setter, fErr),
			Map[TestContext](func(s State) int { return s.Value }),
		)
		_, err := runEffect(eff, testCtx)
		assert.Error(t, err)
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

	f := func(n int) func(context.Context, TestContext) (int, error) {
		return func(ctx context.Context, c TestContext) (int, error) {
			return n + len(c.Value), nil
		}
	}

	t.Run("success updates via lens", func(t *testing.T) {
		eff := F.Pipe2(
			Do[TestContext](State{Count: 10}),
			BindIL(countLens, f),
			Map[TestContext](func(s State) int { return s.Count }),
		)
		v, err := runEffect(eff, testCtx)
		assert.NoError(t, err)
		// 10 + len("test") = 14
		assert.Equal(t, 14, v)
	})

	t.Run("error propagates", func(t *testing.T) {
		fErr := func(n int) func(context.Context, TestContext) (int, error) {
			return func(ctx context.Context, c TestContext) (int, error) {
				return 0, idiomaticTestErr
			}
		}
		eff := F.Pipe2(
			Do[TestContext](State{Count: 10}),
			BindIL(countLens, fErr),
			Map[TestContext](func(s State) int { return s.Count }),
		)
		_, err := runEffect(eff, testCtx)
		assert.Error(t, err)
	})
}
