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

package readerreaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/stretchr/testify/assert"
)

type Config1 struct {
	value1 int
}

type Config2 struct {
	value2 string
}

type Context struct {
	contextID string
}

func TestSequence(t *testing.T) {
	t.Run("swaps parameter order for simple types", func(t *testing.T) {
		// Original: takes Config2, returns ReaderIOEither that may produce ReaderReaderIOEither[Config1, Context, error, int]
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, int] {
				return RIOE.Of[Context, error](cfg1.value1 + len(cfg2.value2))
			})
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := Sequence(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// Test original: Config2 -> Context -> Config1 -> Context
		result1 := original(cfg2)(ctx)()
		assert.True(t, E.IsRight(result1))
		innerFunc1, _ := E.Unwrap(result1)
		innerResult1 := innerFunc1(cfg1)(ctx)()
		assert.Equal(t, E.Right[error](15), innerResult1)

		// Test sequenced: Config1 -> Config2 -> Context
		innerFunc2 := sequenced(cfg1)
		innerResult2 := innerFunc2(cfg2)(ctx)()
		assert.Equal(t, E.Right[error](15), innerResult2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		testErr := errors.New("test error")

		// Original that returns an error
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Left[Context, ReaderReaderIOEither[Config1, Context, error, int]](testErr)
		}

		sequenced := Sequence(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// Test sequenced preserves error
		innerFunc := sequenced(cfg1)
		result := innerFunc(cfg2)(ctx)()
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("works with nested computations", func(t *testing.T) {
		// Original with nested logic
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, string]] {
			if len(cfg2.value2) == 0 {
				return RIOE.Left[Context, ReaderReaderIOEither[Config1, Context, error, string]](errors.New("empty string"))
			}
			return RIOE.Of[Context, error](func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, string] {
				if cfg1.value1 < 0 {
					return RIOE.Left[Context, string](errors.New("negative value"))
				}
				return RIOE.Of[Context, error](fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1))
			})
		}

		sequenced := Sequence(original)

		ctx := Context{contextID: "test"}

		// Test with valid inputs
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, E.Right[error]("test:42"), result1)

		// Test with empty string
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, E.IsLeft(result2))

		// Test with negative value
		result3 := sequenced(Config1{value1: -1})(Config2{value2: "test"})(ctx)()
		assert.True(t, E.IsLeft(result3))
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, int] {
				return RIOE.Of[Context, error](cfg1.value1 + len(cfg2.value2))
			})
		}

		sequenced := Sequence(original)

		result := sequenced(Config1{value1: 0})(Config2{value2: ""})(Context{contextID: ""})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, int] {
				return RIOE.Of[Context, error](cfg1.value1 * len(cfg2.value2))
			})
		}

		sequenced := Sequence(original)

		cfg1 := Config1{value1: 3}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		// Call multiple times with same inputs
		for range 5 {
			result := sequenced(cfg1)(cfg2)(ctx)()
			assert.Equal(t, E.Right[error](12), result)
		}
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("swaps parameter order for Reader types", func(t *testing.T) {
		// Original: takes Config2, returns ReaderIOEither that may produce Reader[Config1, int]
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) int {
				return cfg1.value1 + len(cfg2.value2)
			})
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := SequenceReader(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// Test original
		result1 := original(cfg2)(ctx)()
		assert.True(t, E.IsRight(result1))
		innerFunc1, _ := E.Unwrap(result1)
		value1 := innerFunc1(cfg1)
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced(cfg1)
		result2 := innerFunc2(cfg2)(ctx)()
		assert.True(t, E.IsRight(result2))
		value2, _ := E.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		testErr := errors.New("test error")

		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, int]] {
			return RIOE.Left[Context, R.Reader[Config1, int]](testErr)
		}

		sequenced := SequenceReader(original)

		result := sequenced(Config1{value1: 10})(Config2{value2: "hello"})(Context{contextID: "test"})()
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("works with pure Reader computations", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, string]] {
			if len(cfg2.value2) == 0 {
				return RIOE.Left[Context, R.Reader[Config1, string]](errors.New("empty string"))
			}
			return RIOE.Of[Context, error](func(cfg1 Config1) string {
				return fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1)
			})
		}

		sequenced := SequenceReader(original)

		ctx := Context{contextID: "test"}

		// Test with valid inputs
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, E.Right[error]("test:42"), result1)

		// Test with empty string
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, E.IsLeft(result2))
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) int {
				return cfg1.value1 + len(cfg2.value2)
			})
		}

		sequenced := SequenceReader(original)

		result := sequenced(Config1{value1: 0})(Config2{value2: ""})(Context{contextID: ""})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) int {
				return cfg1.value1 * len(cfg2.value2)
			})
		}

		sequenced := SequenceReader(original)

		cfg1 := Config1{value1: 3}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		// Call multiple times with same inputs
		for range 5 {
			result := sequenced(cfg1)(cfg2)(ctx)()
			assert.Equal(t, E.Right[error](12), result)
		}
	})
}

func TestSequenceReaderIO(t *testing.T) {
	t.Run("swaps parameter order for ReaderIO types", func(t *testing.T) {
		// Original: takes Config2, returns ReaderIOEither that may produce ReaderIO[Config1, int]
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) io.IO[int] {
				return io.Of(cfg1.value1 + len(cfg2.value2))
			})
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := SequenceReaderIO(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// Test original
		result1 := original(cfg2)(ctx)()
		assert.True(t, E.IsRight(result1))
		innerFunc1, _ := E.Unwrap(result1)
		value1 := innerFunc1(cfg1)()
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced(cfg1)
		result2 := innerFunc2(cfg2)(ctx)()
		assert.True(t, E.IsRight(result2))
		value2, _ := E.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		testErr := errors.New("test error")

		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, int]] {
			return RIOE.Left[Context, readerio.ReaderIO[Config1, int]](testErr)
		}

		sequenced := SequenceReaderIO(original)

		result := sequenced(Config1{value1: 10})(Config2{value2: "hello"})(Context{contextID: "test"})()
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("works with IO effects", func(t *testing.T) {
		sideEffect := 0

		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, string]] {
			if len(cfg2.value2) == 0 {
				return RIOE.Left[Context, readerio.ReaderIO[Config1, string]](errors.New("empty string"))
			}
			return RIOE.Of[Context, error](func(cfg1 Config1) io.IO[string] {
				return func() string {
					sideEffect = cfg1.value1
					return fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1)
				}
			})
		}

		sequenced := SequenceReaderIO(original)

		ctx := Context{contextID: "test"}

		// Test with valid inputs
		sideEffect = 0
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, E.Right[error]("test:42"), result1)
		assert.Equal(t, 42, sideEffect)

		// Test with empty string
		sideEffect = 0
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, E.IsLeft(result2))
		assert.Equal(t, 0, sideEffect) // Side effect should not occur
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) io.IO[int] {
				return io.Of(cfg1.value1 + len(cfg2.value2))
			})
		}

		sequenced := SequenceReaderIO(original)

		result := sequenced(Config1{value1: 0})(Config2{value2: ""})(Context{contextID: ""})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("executes IO effects correctly", func(t *testing.T) {
		counter := 0

		original := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, int]] {
			return RIOE.Of[Context, error](func(cfg1 Config1) io.IO[int] {
				return func() int {
					counter++
					return cfg1.value1 + len(cfg2.value2)
				}
			})
		}

		sequenced := SequenceReaderIO(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// Each execution should increment counter
		counter = 0
		result1 := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error](15), result1)
		assert.Equal(t, 1, counter)

		result2 := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error](15), result2)
		assert.Equal(t, 2, counter)
	})
}

func TestTraverse(t *testing.T) {
	t.Run("transforms and swaps parameter order", func(t *testing.T) {
		// Original computation depending on Config2
		original := Of[Config2, Context, error](42)

		// Transformation that introduces Config1 dependency
		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, string] {
			return func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, string] {
				return RIOE.Of[Context, error](fmt.Sprintf("value=%d, cfg1=%d", n, cfg1.value1))
			}
		}

		// Apply traverse to swap order and transform
		traversed := Traverse[Config2](transform)(original)

		cfg1 := Config1{value1: 100}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		result := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error]("value=42, cfg1=100"), result)
	})

	t.Run("preserves error handling in original", func(t *testing.T) {
		testErr := errors.New("test error")
		original := Left[Config2, Context, int](testErr)

		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, string] {
			return Of[Config1, Context, error](fmt.Sprintf("%d", n))
		}

		traversed := Traverse[Config2](transform)(original)

		result := traversed(Config1{value1: 100})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Left[string](testErr), result)
	})

	t.Run("preserves error handling in transformation", func(t *testing.T) {
		original := Of[Config2, Context, error](42)
		testErr := errors.New("transform error")

		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, string] {
			if n < 0 {
				return Left[Config1, Context, string](testErr)
			}
			return Of[Config1, Context, error](fmt.Sprintf("%d", n))
		}

		// Test with negative value
		originalNeg := Of[Config2, Context, error](-1)
		traversedNeg := Traverse[Config2](transform)(originalNeg)
		resultNeg := traversedNeg(Config1{value1: 100})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Left[string](testErr), resultNeg)

		// Test with positive value
		traversedPos := Traverse[Config2](transform)(original)
		resultPos := traversedPos(Config1{value1: 100})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Right[error]("42"), resultPos)
	})

	t.Run("works with complex transformations", func(t *testing.T) {
		original := Of[Config2, Context, error](10)

		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, int] {
			return func(cfg1 Config1) RIOE.ReaderIOEither[Context, error, int] {
				return func(ctx Context) IOE.IOEither[error, int] {
					return IOE.Of[error](n * cfg1.value1)
				}
			}
		}

		traversed := Traverse[Config2](transform)(original)

		result := traversed(Config1{value1: 5})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Right[error](50), result)
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		original := Of[Config2, Context, error](10)

		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, int] {
			return Of[Config1, Context, error](n * 2)
		}

		result := F.Pipe2(
			original,
			Traverse[Config2](transform),
			func(k Kleisli[Config2, Context, error, Config1, int]) ReaderReaderIOEither[Config2, Context, error, int] {
				return k(Config1{value1: 5})
			},
		)

		outcome := result(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Right[error](20), outcome)
	})
}

func TestTraverseReader(t *testing.T) {
	t.Run("transforms with pure Reader and swaps parameter order", func(t *testing.T) {
		// Original computation depending on Config2
		original := Of[Config2, Context, error](100)

		// Pure Reader transformation that introduces Config1 dependency
		formatWithConfig := func(value int) R.Reader[Config1, string] {
			return func(cfg1 Config1) string {
				return fmt.Sprintf("value=%d, multiplier=%d, result=%d", value, cfg1.value1, value*cfg1.value1)
			}
		}

		// Apply traverse to introduce Config1 and swap order
		traversed := TraverseReader[Config2, Config1, Context, error](formatWithConfig)(original)

		cfg1 := Config1{value1: 5}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		result := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error]("value=100, multiplier=5, result=500"), result)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		testErr := errors.New("test error")
		original := Left[Config2, Context, int](testErr)

		transform := func(n int) R.Reader[Config1, string] {
			return R.Of[Config1](fmt.Sprintf("%d", n))
		}

		traversed := TraverseReader[Config2, Config1, Context, error](transform)(original)

		result := traversed(Config1{value1: 5})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Left[string](testErr), result)
	})

	t.Run("works with pure computations", func(t *testing.T) {
		original := Of[Config2, Context, error](42)

		// Pure transformation using Reader
		double := func(n int) R.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		traversed := TraverseReader[Config2, Config1, Context, error](double)(original)

		result := traversed(Config1{value1: 3})(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Right[error](126), result)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := Of[Config2, Context, error](0)

		transform := func(n int) R.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n + cfg1.value1
			}
		}

		traversed := TraverseReader[Config2, Config1, Context, error](transform)(original)

		result := traversed(Config1{value1: 0})(Config2{value2: ""})(Context{contextID: ""})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		original := Of[Config2, Context, error](10)

		transform := func(n int) R.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		traversed := TraverseReader[Config2, Config1, Context, error](transform)(original)

		cfg1 := Config1{value1: 5}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		// Call multiple times with same inputs
		for range 5 {
			result := traversed(cfg1)(cfg2)(ctx)()
			assert.Equal(t, E.Right[error](50), result)
		}
	})

	t.Run("can be used in composition", func(t *testing.T) {
		original := Of[Config2, Context, error](10)

		multiply := func(n int) R.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		result := F.Pipe2(
			original,
			TraverseReader[Config2, Config1, Context, error](multiply),
			func(k Kleisli[Config2, Context, error, Config1, int]) ReaderReaderIOEither[Config2, Context, error, int] {
				return k(Config1{value1: 3})
			},
		)

		outcome := result(Config2{value2: "test"})(Context{contextID: "test"})()
		assert.Equal(t, E.Right[error](30), outcome)
	})
}

func TestFlipIntegration(t *testing.T) {
	t.Run("Sequence and Traverse work together", func(t *testing.T) {
		// Create a nested computation
		nested := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Of[Context, error](Of[Config1, Context, error](len(cfg2.value2)))
		}

		// Sequence it
		sequenced := Sequence(nested)

		// Then traverse with a transformation
		transform := func(n int) ReaderReaderIOEither[Config1, Context, error, string] {
			return Of[Config1, Context, error](fmt.Sprintf("length=%d", n))
		}

		// Apply both operations
		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}
		ctx := Context{contextID: "test"}

		// First sequence
		intermediate := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error](5), intermediate)

		// Then apply traverse on a new computation
		original := Of[Config2, Context, error](5)
		traversed := Traverse[Config2](transform)(original)
		result := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, E.Right[error]("length=5"), result)
	})

	t.Run("all flip functions preserve error semantics", func(t *testing.T) {
		testErr := errors.New("test error")
		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "test"}
		ctx := Context{contextID: "test"}

		// Test Sequence with error
		seqErr := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, ReaderReaderIOEither[Config1, Context, error, int]] {
			return RIOE.Left[Context, ReaderReaderIOEither[Config1, Context, error, int]](testErr)
		}
		seqResult := Sequence(seqErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, E.IsLeft(seqResult))

		// Test SequenceReader with error
		seqReaderErr := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, R.Reader[Config1, int]] {
			return RIOE.Left[Context, R.Reader[Config1, int]](testErr)
		}
		seqReaderResult := SequenceReader(seqReaderErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, E.IsLeft(seqReaderResult))

		// Test SequenceReaderIO with error
		seqReaderIOErr := func(cfg2 Config2) RIOE.ReaderIOEither[Context, error, readerio.ReaderIO[Config1, int]] {
			return RIOE.Left[Context, readerio.ReaderIO[Config1, int]](testErr)
		}
		seqReaderIOResult := SequenceReaderIO(seqReaderIOErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, E.IsLeft(seqReaderIOResult))

		// Test Traverse with error
		travErr := Left[Config2, Context, int](testErr)
		travTransform := func(n int) ReaderReaderIOEither[Config1, Context, error, string] {
			return Of[Config1, Context, error](fmt.Sprintf("%d", n))
		}
		travResult := Traverse[Config2](travTransform)(travErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, E.IsLeft(travResult))

		// Test TraverseReader with error
		travReaderErr := Left[Config2, Context, int](testErr)
		travReaderTransform := func(n int) R.Reader[Config1, string] {
			return R.Of[Config1](fmt.Sprintf("%d", n))
		}
		travReaderResult := TraverseReader[Config2, Config1, Context, error](travReaderTransform)(travReaderErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, E.IsLeft(travReaderResult))
	})
}
