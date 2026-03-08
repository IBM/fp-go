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

package iooption

import (
	"fmt"
	"testing"
	"time"

	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	I "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	result := Of(42)()
	assert.Equal(t, O.Some(42), result)
}

func TestSome(t *testing.T) {
	result := Some("test")()
	assert.Equal(t, O.Some("test"), result)
}

func TestNone(t *testing.T) {
	result := None[int]()()
	assert.Equal(t, O.None[int](), result)
}

func TestMonadOf(t *testing.T) {
	result := MonadOf(100)()
	assert.Equal(t, O.Some(100), result)
}

func TestFromOptionComprehensive(t *testing.T) {
	t.Run("from Some", func(t *testing.T) {
		result := FromOption(O.Some(42))()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("from None", func(t *testing.T) {
		result := FromOption(O.None[int]())()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestFromIO(t *testing.T) {
	ioValue := I.Of(42)
	result := FromIO(ioValue)()
	assert.Equal(t, O.Some(42), result)
}

func TestMonadMap(t *testing.T) {
	t.Run("map over Some", func(t *testing.T) {
		result := MonadMap(Of(5), utils.Double)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("map over None", func(t *testing.T) {
		result := MonadMap(None[int](), utils.Double)()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestMonadChain(t *testing.T) {
	t.Run("chain Some to Some", func(t *testing.T) {
		f := func(n int) IOOption[int] {
			return Of(n * 2)
		}
		result := MonadChain(Of(5), f)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("chain Some to None", func(t *testing.T) {
		f := func(n int) IOOption[int] {
			return None[int]()
		}
		result := MonadChain(Of(5), f)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("chain None", func(t *testing.T) {
		f := func(n int) IOOption[int] {
			return Of(n * 2)
		}
		result := MonadChain(None[int](), f)()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestChain(t *testing.T) {
	f := func(n int) IOOption[string] {
		if n > 0 {
			return Of("positive")
		}
		return None[string]()
	}

	t.Run("chain positive", func(t *testing.T) {
		result := F.Pipe1(Of(5), Chain(f))()
		assert.Equal(t, O.Some("positive"), result)
	})

	t.Run("chain negative", func(t *testing.T) {
		result := F.Pipe1(Of(-5), Chain(f))()
		assert.Equal(t, O.None[string](), result)
	})
}

func TestMonadAp(t *testing.T) {
	t.Run("apply Some function to Some value", func(t *testing.T) {
		mab := Of(utils.Double)
		ma := Of(5)
		result := MonadAp(mab, ma)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("apply None function", func(t *testing.T) {
		mab := None[func(int) int]()
		ma := Of(5)
		result := MonadAp(mab, ma)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("apply to None value", func(t *testing.T) {
		mab := Of(utils.Double)
		ma := None[int]()
		result := MonadAp(mab, ma)()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestAp(t *testing.T) {
	ma := Of(5)
	result := F.Pipe1(Of(utils.Double), Ap[int, int](ma))()
	assert.Equal(t, O.Some(10), result)
}

func TestApSeq(t *testing.T) {
	ma := Of(5)
	result := F.Pipe1(Of(utils.Double), ApSeq[int, int](ma))()
	assert.Equal(t, O.Some(10), result)
}

func TestApPar(t *testing.T) {
	ma := Of(5)
	result := F.Pipe1(Of(utils.Double), ApPar[int, int](ma))()
	assert.Equal(t, O.Some(10), result)
}

func TestFlatten(t *testing.T) {
	t.Run("flatten Some(Some)", func(t *testing.T) {
		nested := Of(Of(42))
		result := Flatten(nested)()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("flatten Some(None)", func(t *testing.T) {
		nested := Of(None[int]())
		result := Flatten(nested)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("flatten None", func(t *testing.T) {
		nested := None[IOOption[int]]()
		result := Flatten(nested)()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestOptionize0(t *testing.T) {
	f := func() (int, bool) {
		return 42, true
	}
	result := Optionize0(f)()()
	assert.Equal(t, O.Some(42), result)

	f2 := func() (int, bool) {
		return 0, false
	}
	result2 := Optionize0(f2)()()
	assert.Equal(t, O.None[int](), result2)
}

func TestOptionize2(t *testing.T) {
	f := func(a, b int) (int, bool) {
		if b != 0 {
			return a / b, true
		}
		return 0, false
	}

	result := Optionize2(f)(10, 2)()
	assert.Equal(t, O.Some(5), result)

	result2 := Optionize2(f)(10, 0)()
	assert.Equal(t, O.None[int](), result2)
}

func TestOptionize3(t *testing.T) {
	f := func(a, b, c int) (int, bool) {
		if c != 0 {
			return (a + b) / c, true
		}
		return 0, false
	}

	result := Optionize3(f)(10, 5, 3)()
	assert.Equal(t, O.Some(5), result)

	result2 := Optionize3(f)(10, 5, 0)()
	assert.Equal(t, O.None[int](), result2)
}

func TestOptionize4(t *testing.T) {
	f := func(a, b, c, d int) (int, bool) {
		if d != 0 {
			return (a + b + c) / d, true
		}
		return 0, false
	}

	result := Optionize4(f)(10, 5, 3, 2)()
	assert.Equal(t, O.Some(9), result)

	result2 := Optionize4(f)(10, 5, 3, 0)()
	assert.Equal(t, O.None[int](), result2)
}

func TestMemoize(t *testing.T) {
	callCount := 0
	ioOpt := func() Option[int] {
		callCount++
		return O.Some(42)
	}

	memoized := Memoize(ioOpt)

	// First call
	result1 := memoized()
	assert.Equal(t, O.Some(42), result1)
	assert.Equal(t, 1, callCount)

	// Second call should use cached value
	result2 := memoized()
	assert.Equal(t, O.Some(42), result2)
	assert.Equal(t, 1, callCount)
}

func TestFold(t *testing.T) {
	onNone := I.Of("none")
	onSome := func(n int) I.IO[string] {
		return I.Of(fmt.Sprintf("%d", n))
	}

	t.Run("fold Some", func(t *testing.T) {
		result := Fold(onNone, onSome)(Of(42))()
		assert.Equal(t, "42", result)
	})

	t.Run("fold None", func(t *testing.T) {
		result := Fold(onNone, onSome)(None[int]())()
		assert.Equal(t, "none", result)
	})
}

func TestDefer(t *testing.T) {
	callCount := 0
	gen := func() IOOption[int] {
		callCount++
		return Of(42)
	}

	deferred := Defer(gen)

	// Each call should invoke the generator
	result1 := deferred()
	assert.Equal(t, O.Some(42), result1)
	assert.Equal(t, 1, callCount)

	result2 := deferred()
	assert.Equal(t, O.Some(42), result2)
	assert.Equal(t, 2, callCount)
}

func TestFromEither(t *testing.T) {
	t.Run("from Right", func(t *testing.T) {
		either := ET.Right[string](42)
		result := FromEither(either)()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("from Left", func(t *testing.T) {
		either := ET.Left[int]("error")
		result := FromEither(either)()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestMonadAlt(t *testing.T) {
	t.Run("first is Some", func(t *testing.T) {
		result := MonadAlt(Of(1), Of(2))()
		assert.Equal(t, O.Some(1), result)
	})

	t.Run("first is None, second is Some", func(t *testing.T) {
		result := MonadAlt(None[int](), Of(2))()
		assert.Equal(t, O.Some(2), result)
	})

	t.Run("both are None", func(t *testing.T) {
		result := MonadAlt(None[int](), None[int]())()
		assert.Equal(t, O.None[int](), result)
	})
}

func TestAlt(t *testing.T) {
	t.Run("first is Some", func(t *testing.T) {
		result := F.Pipe1(Of(1), Alt(Of(2)))()
		assert.Equal(t, O.Some(1), result)
	})

	t.Run("first is None", func(t *testing.T) {
		result := F.Pipe1(None[int](), Alt(Of(2)))()
		assert.Equal(t, O.Some(2), result)
	})
}

func TestMonadChainFirst(t *testing.T) {
	sideEffect := 0
	f := func(n int) IOOption[string] {
		sideEffect = n * 2
		return Of("side effect")
	}

	result := MonadChainFirst(Of(5), f)()
	assert.Equal(t, O.Some(5), result)
	assert.Equal(t, 10, sideEffect)
}

func TestChainFirst(t *testing.T) {
	sideEffect := 0
	f := func(n int) IOOption[string] {
		sideEffect = n * 2
		return Of("side effect")
	}

	result := F.Pipe1(Of(5), ChainFirst(f))()
	assert.Equal(t, O.Some(5), result)
	assert.Equal(t, 10, sideEffect)
}

func TestMonadChainFirstIOK(t *testing.T) {
	sideEffect := 0
	f := func(n int) I.IO[string] {
		return func() string {
			sideEffect = n * 2
			return "side effect"
		}
	}

	result := MonadChainFirstIOK(Of(5), f)()
	assert.Equal(t, O.Some(5), result)
	assert.Equal(t, 10, sideEffect)
}

func TestChainFirstIOK(t *testing.T) {
	sideEffect := 0
	f := func(n int) I.IO[string] {
		return func() string {
			sideEffect = n * 2
			return "side effect"
		}
	}

	result := F.Pipe1(Of(5), ChainFirstIOK(f))()
	assert.Equal(t, O.Some(5), result)
	assert.Equal(t, 10, sideEffect)
}

func TestDelay(t *testing.T) {
	start := time.Now()
	delay := 50 * time.Millisecond

	result := F.Pipe1(Of(42), Delay[int](delay))()

	elapsed := time.Since(start)
	assert.Equal(t, O.Some(42), result)
	assert.True(t, elapsed >= delay, "Expected delay of at least %v, got %v", delay, elapsed)
}

func TestAfter(t *testing.T) {
	timestamp := time.Now().Add(50 * time.Millisecond)

	result := F.Pipe1(Of(42), After[int](timestamp))()

	assert.Equal(t, O.Some(42), result)
	assert.True(t, time.Now().After(timestamp) || time.Now().Equal(timestamp))
}

func TestMonadChainIOK(t *testing.T) {
	f := func(n int) I.IO[string] {
		return I.Of(fmt.Sprintf("%d", n))
	}

	t.Run("chain Some", func(t *testing.T) {
		result := MonadChainIOK(Of(42), f)()
		assert.Equal(t, O.Some("42"), result)
	})

	t.Run("chain None", func(t *testing.T) {
		result := MonadChainIOK(None[int](), f)()
		assert.Equal(t, O.None[string](), result)
	})
}

// Made with Bob
