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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestLeft(t *testing.T) {
	err := errors.New("test error")
	res := Left[int](err)()
	assert.Equal(t, result.Left[int](err), res)
}

func TestRight(t *testing.T) {
	res := Right(42)()
	assert.Equal(t, result.Of(42), res)
}

func TestOf(t *testing.T) {
	res := Of(42)()
	assert.Equal(t, result.Of(42), res)
}

func TestMonadOf(t *testing.T) {
	res := MonadOf(42)()
	assert.Equal(t, result.Of(42), res)
}

func TestLeftIO(t *testing.T) {
	err := errors.New("test error")
	res := LeftIO[int](io.Of(err))()
	assert.Equal(t, result.Left[int](err), res)
}

func TestRightIO(t *testing.T) {
	res := RightIO(io.Of(42))()
	assert.Equal(t, result.Of(42), res)
}

func TestFromEither(t *testing.T) {
	t.Run("from Right", func(t *testing.T) {
		either := result.Of(42)
		res := FromEither(either)()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("from Left", func(t *testing.T) {
		err := errors.New("test error")
		either := result.Left[int](err)
		res := FromEither(either)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestFromResult(t *testing.T) {
	t.Run("from success", func(t *testing.T) {
		res := FromResult(result.Of(42))()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("from error", func(t *testing.T) {
		err := errors.New("test error")
		res := FromResult(result.Left[int](err))()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestFromEitherI(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		res := FromEitherI(42, nil)()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("with error", func(t *testing.T) {
		err := errors.New("test error")
		res := FromEitherI(0, err)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestFromResultI(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		res := FromResultI(42, nil)()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("with error", func(t *testing.T) {
		err := errors.New("test error")
		res := FromResultI(0, err)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestFromOption_Success(t *testing.T) {
	onNone := func() error {
		return errors.New("none")
	}

	t.Run("from Some", func(t *testing.T) {
		res := FromOption[int](onNone)(O.Some(42))()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("from None", func(t *testing.T) {
		res := FromOption[int](onNone)(O.None[int]())()
		assert.Equal(t, result.Left[int](errors.New("none")), res)
	})
}

func TestFromIO(t *testing.T) {
	ioValue := io.Of(42)
	res := FromIO(ioValue)()
	assert.Equal(t, result.Of(42), res)
}

func TestFromLazy(t *testing.T) {
	lazy := func() int { return 42 }
	res := FromLazy(lazy)()
	assert.Equal(t, result.Of(42), res)
}

func TestMonadMap(t *testing.T) {
	t.Run("map over Right", func(t *testing.T) {
		res := MonadMap(Of(5), utils.Double)()
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("map over Left", func(t *testing.T) {
		err := errors.New("test error")
		res := MonadMap(Left[int](err), utils.Double)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestMap_Comprehensive(t *testing.T) {
	double := func(n int) int { return n * 2 }

	t.Run("map Right", func(t *testing.T) {
		res := F.Pipe1(Of(5), Map(double))()
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("map Left", func(t *testing.T) {
		err := errors.New("test error")
		res := F.Pipe1(Left[int](err), Map(double))()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestMonadMapTo(t *testing.T) {
	t.Run("mapTo Right", func(t *testing.T) {
		res := MonadMapTo(Of(5), "constant")()
		assert.Equal(t, result.Of("constant"), res)
	})

	t.Run("mapTo Left", func(t *testing.T) {
		err := errors.New("test error")
		res := MonadMapTo(Left[int](err), "constant")()
		assert.Equal(t, result.Left[string](err), res)
	})
}

func TestMapTo(t *testing.T) {
	res := F.Pipe1(Of(5), MapTo[int]("constant"))()
	assert.Equal(t, result.Of("constant"), res)
}

func TestMonadChain(t *testing.T) {
	f := func(n int) IOResult[int] {
		return Of(n * 2)
	}

	t.Run("chain Right to Right", func(t *testing.T) {
		res := MonadChain(Of(5), f)()
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("chain Right to Left", func(t *testing.T) {
		err := errors.New("test error")
		f := func(n int) IOResult[int] {
			return Left[int](err)
		}
		res := MonadChain(Of(5), f)()
		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("chain Left", func(t *testing.T) {
		err := errors.New("test error")
		res := MonadChain(Left[int](err), f)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestChain_Comprehensive(t *testing.T) {
	f := func(n int) IOResult[string] {
		if n > 0 {
			return Of(fmt.Sprintf("%d", n))
		}
		return Left[string](errors.New("negative"))
	}

	t.Run("chain positive", func(t *testing.T) {
		res := F.Pipe1(Of(5), Chain(f))()
		assert.Equal(t, result.Of("5"), res)
	})

	t.Run("chain negative", func(t *testing.T) {
		res := F.Pipe1(Of(-5), Chain(f))()
		assert.Equal(t, result.Left[string](errors.New("negative")), res)
	})
}

func TestMonadChainEitherK(t *testing.T) {
	f := func(n int) result.Result[int] {
		if n > 0 {
			return result.Of(n * 2)
		}
		return result.Left[int](errors.New("non-positive"))
	}

	t.Run("chain to success", func(t *testing.T) {
		res := MonadChainEitherK(Of(5), f)()
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("chain to error", func(t *testing.T) {
		res := MonadChainEitherK(Of(-5), f)()
		assert.Equal(t, result.Left[int](errors.New("non-positive")), res)
	})
}

func TestMonadChainResultK(t *testing.T) {
	f := func(n int) result.Result[int] {
		return result.Of(n * 2)
	}

	res := MonadChainResultK(Of(5), f)()
	assert.Equal(t, result.Of(10), res)
}

func TestChainResultK(t *testing.T) {
	f := func(n int) result.Result[int] {
		return result.Of(n * 2)
	}

	res := F.Pipe1(Of(5), ChainResultK(f))()
	assert.Equal(t, result.Of(10), res)
}

func TestMonadAp_Comprehensive(t *testing.T) {
	t.Run("apply Right function to Right value", func(t *testing.T) {
		mab := Of(utils.Double)
		ma := Of(5)
		res := MonadAp(mab, ma)()
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("apply Left function", func(t *testing.T) {
		err := errors.New("function error")
		mab := Left[func(int) int](err)
		ma := Of(5)
		res := MonadAp(mab, ma)()
		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("apply to Left value", func(t *testing.T) {
		err := errors.New("value error")
		mab := Of(utils.Double)
		ma := Left[int](err)
		res := MonadAp(mab, ma)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestAp_Comprehensive(t *testing.T) {
	ma := Of(5)
	res := F.Pipe1(Of(utils.Double), Ap[int, int](ma))()
	assert.Equal(t, result.Of(10), res)
}

func TestApPar(t *testing.T) {
	ma := Of(5)
	res := F.Pipe1(Of(utils.Double), ApPar[int, int](ma))()
	assert.Equal(t, result.Of(10), res)
}

func TestApSeq(t *testing.T) {
	ma := Of(5)
	res := F.Pipe1(Of(utils.Double), ApSeq[int, int](ma))()
	assert.Equal(t, result.Of(10), res)
}

func TestFlatten_Comprehensive(t *testing.T) {
	t.Run("flatten Right(Right)", func(t *testing.T) {
		nested := Of(Of(42))
		res := Flatten(nested)()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("flatten Right(Left)", func(t *testing.T) {
		err := errors.New("inner error")
		nested := Of(Left[int](err))
		res := Flatten(nested)()
		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("flatten Left", func(t *testing.T) {
		err := errors.New("outer error")
		nested := Left[IOResult[int]](err)
		res := Flatten(nested)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestTryCatch(t *testing.T) {
	t.Run("successful function", func(t *testing.T) {
		f := func() (int, error) {
			return 42, nil
		}
		res := TryCatch(f, F.Identity[error])()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("failing function", func(t *testing.T) {
		err := errors.New("test error")
		f := func() (int, error) {
			return 0, err
		}
		res := TryCatch(f, F.Identity[error])()
		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("with error transformation", func(t *testing.T) {
		err := errors.New("original")
		f := func() (int, error) {
			return 0, err
		}
		onThrow := func(e error) error {
			return fmt.Errorf("wrapped: %w", e)
		}
		res := TryCatch(f, onThrow)()
		assert.True(t, result.IsLeft(res))
	})
}

func TestTryCatchError_Comprehensive(t *testing.T) {
	t.Run("successful function", func(t *testing.T) {
		f := func() (int, error) {
			return 42, nil
		}
		res := TryCatchError(f)()
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("failing function", func(t *testing.T) {
		err := errors.New("test error")
		f := func() (int, error) {
			return 0, err
		}
		res := TryCatchError(f)()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestMemoize_Comprehensive(t *testing.T) {
	callCount := 0
	ioRes := func() Result[int] {
		callCount++
		return result.Of(42)
	}

	memoized := Memoize(ioRes)

	// First call
	res1 := memoized()
	assert.Equal(t, result.Of(42), res1)
	assert.Equal(t, 1, callCount)

	// Second call should use cached value
	res2 := memoized()
	assert.Equal(t, result.Of(42), res2)
	assert.Equal(t, 1, callCount)
}

func TestMonadMapLeft(t *testing.T) {
	t.Run("map Left error", func(t *testing.T) {
		err := errors.New("original")
		f := func(e error) string {
			return e.Error()
		}
		res := MonadMapLeft(Left[int](err), f)()
		// Result is IOEither[string, int], check it's a left
		assert.True(t, ET.IsLeft(res))
	})

	t.Run("map Right unchanged", func(t *testing.T) {
		f := func(e error) string {
			return e.Error()
		}
		res := MonadMapLeft(Of(42), f)()
		// MapLeft changes the error type, so result is IOEither[string, int]
		assert.True(t, ET.IsRight(res))
		assert.Equal(t, 42, ET.MonadFold(res, func(string) int { return 0 }, F.Identity[int]))
	})
}

func TestMapLeft_Comprehensive(t *testing.T) {
	f := func(e error) string {
		return fmt.Sprintf("wrapped: %s", e.Error())
	}

	t.Run("map Left", func(t *testing.T) {
		err := errors.New("original")
		res := F.Pipe1(Left[int](err), MapLeft[int](f))()
		// Result is IOEither[string, int], check it's a left
		assert.True(t, ET.IsLeft(res))
	})

	t.Run("map Right unchanged", func(t *testing.T) {
		res := F.Pipe1(Of(42), MapLeft[int](f))()
		// MapLeft changes the error type, so result is IOEither[string, int]
		assert.True(t, ET.IsRight(res))
		assert.Equal(t, 42, ET.MonadFold(res, func(string) int { return 0 }, F.Identity[int]))
	})
}

func TestMonadBiMap(t *testing.T) {
	leftF := func(e error) string {
		return e.Error()
	}
	rightF := func(n int) string {
		return fmt.Sprintf("%d", n)
	}

	t.Run("bimap Right", func(t *testing.T) {
		res := MonadBiMap(Of(42), leftF, rightF)()
		// BiMap changes both types, so result is IOEither[string, string]
		assert.True(t, ET.IsRight(res))
		assert.Equal(t, "42", ET.MonadFold(res, F.Identity[string], F.Identity[string]))
	})

	t.Run("bimap Left", func(t *testing.T) {
		err := errors.New("test")
		res := MonadBiMap(Left[int](err), leftF, rightF)()
		// Result is IOEither[string, string], check it's a left
		assert.True(t, ET.IsLeft(res))
	})
}

func TestBiMap_Comprehensive(t *testing.T) {
	leftF := func(e error) string {
		return e.Error()
	}
	rightF := func(n int) string {
		return fmt.Sprintf("%d", n)
	}

	t.Run("bimap Right", func(t *testing.T) {
		res := F.Pipe1(Of(42), BiMap(leftF, rightF))()
		// BiMap changes both types, so result is IOEither[string, string]
		assert.True(t, ET.IsRight(res))
		assert.Equal(t, "42", ET.MonadFold(res, F.Identity[string], F.Identity[string]))
	})

	t.Run("bimap Left", func(t *testing.T) {
		err := errors.New("test")
		res := F.Pipe1(Left[int](err), BiMap(leftF, rightF))()
		// Result is IOEither[string, string], check it's a left
		assert.True(t, ET.IsLeft(res))
	})
}

func TestFold_Comprehensive(t *testing.T) {
	onLeft := func(e error) io.IO[string] {
		return io.Of(fmt.Sprintf("error: %s", e.Error()))
	}
	onRight := func(n int) io.IO[string] {
		return io.Of(fmt.Sprintf("value: %d", n))
	}

	t.Run("fold Right", func(t *testing.T) {
		res := Fold(onLeft, onRight)(Of(42))()
		assert.Equal(t, "value: 42", res)
	})

	t.Run("fold Left", func(t *testing.T) {
		err := errors.New("test")
		res := Fold(onLeft, onRight)(Left[int](err))()
		assert.Equal(t, "error: test", res)
	})
}

func TestGetOrElse_Comprehensive(t *testing.T) {
	onLeft := func(e error) io.IO[int] {
		return io.Of(0)
	}

	t.Run("get Right value", func(t *testing.T) {
		res := GetOrElse(onLeft)(Of(42))()
		assert.Equal(t, 42, res)
	})

	t.Run("get default on Left", func(t *testing.T) {
		err := errors.New("test")
		res := GetOrElse(onLeft)(Left[int](err))()
		assert.Equal(t, 0, res)
	})
}

func TestGetOrElseOf(t *testing.T) {
	onLeft := func(e error) int {
		return 0
	}

	t.Run("get Right value", func(t *testing.T) {
		res := GetOrElseOf(onLeft)(Of(42))()
		assert.Equal(t, 42, res)
	})

	t.Run("get default on Left", func(t *testing.T) {
		err := errors.New("test")
		res := GetOrElseOf(onLeft)(Left[int](err))()
		assert.Equal(t, 0, res)
	})
}

func TestMonadChainTo(t *testing.T) {
	t.Run("chain Right to Right", func(t *testing.T) {
		res := MonadChainTo(Of(1), Of(2))()
		assert.Equal(t, result.Of(2), res)
	})

	t.Run("chain Right to Left", func(t *testing.T) {
		err := errors.New("test")
		res := MonadChainTo(Of(1), Left[int](err))()
		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("chain Left", func(t *testing.T) {
		err := errors.New("test")
		res := MonadChainTo(Left[int](err), Of(2))()
		assert.Equal(t, result.Left[int](err), res)
	})
}

func TestChainLazyK(t *testing.T) {
	f := func(n int) Lazy[string] {
		return func() string {
			return fmt.Sprintf("%d", n)
		}
	}

	res := F.Pipe1(Of(42), ChainLazyK(f))()
	assert.Equal(t, result.Of("42"), res)
}

// Made with Bob
