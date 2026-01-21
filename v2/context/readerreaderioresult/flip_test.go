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

package readerreaderioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type Config1 struct {
	value1 int
}

type Config2 struct {
	value2 string
}

func TestSequence(t *testing.T) {
	t.Run("swaps parameter order for simple types", func(t *testing.T) {
		ctx := t.Context()

		// Original: takes Config2, returns ReaderIOResult that may produce ReaderReaderIOResult[Config1, int]
		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx1 context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Of(func(cfg1 Config1) RIORES.ReaderIOResult[int] {
						return func(ctx2 context.Context) IOResult[int] {
							return func() Result[int] {
								return result.Of(cfg1.value1 + len(cfg2.value2))
							}
						}
					})
				}
			}
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := Sequence(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// Test original: Config2 -> Context -> Config1 -> Context
		result1 := original(cfg2)(ctx)()
		assert.True(t, result.IsRight(result1))
		innerFunc1, _ := result.Unwrap(result1)
		innerResult1 := innerFunc1(cfg1)(ctx)()
		assert.Equal(t, result.Of(15), innerResult1)

		// Test sequenced: Config1 -> Config2 -> Context
		innerFunc2 := sequenced(cfg1)
		innerResult2 := innerFunc2(cfg2)(ctx)()
		assert.Equal(t, result.Of(15), innerResult2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")

		// Original that returns an error
		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Left[ReaderReaderIOResult[Config1, int]](testErr)
				}
			}
		}

		sequenced := Sequence(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// Test sequenced preserves error
		innerFunc := sequenced(cfg1)
		outcome := innerFunc(cfg2)(ctx)()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("works with nested computations", func(t *testing.T) {
		ctx := t.Context()

		// Original with nested logic
		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, string]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, string]] {
				return func() Result[ReaderReaderIOResult[Config1, string]] {
					if len(cfg2.value2) == 0 {
						return result.Left[ReaderReaderIOResult[Config1, string]](errors.New("empty string"))
					}
					return result.Of(func(cfg1 Config1) RIORES.ReaderIOResult[string] {
						return func(ctx context.Context) IOResult[string] {
							return func() Result[string] {
								if cfg1.value1 < 0 {
									return result.Left[string](errors.New("negative value"))
								}
								return result.Of(fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1))
							}
						}
					})
				}
			}
		}

		sequenced := Sequence(original)

		// Test with valid inputs
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of("test:42"), result1)

		// Test with empty string
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, result.IsLeft(result2))

		// Test with negative value
		result3 := sequenced(Config1{value1: -1})(Config2{value2: "test"})(ctx)()
		assert.True(t, result.IsLeft(result3))
	})

	t.Run("works with zero values", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Of(func(cfg1 Config1) RIORES.ReaderIOResult[int] {
						return func(ctx context.Context) IOResult[int] {
							return func() Result[int] {
								return result.Of(cfg1.value1 + len(cfg2.value2))
							}
						}
					})
				}
			}
		}

		sequenced := Sequence(original)

		outcome := sequenced(Config1{value1: 0})(Config2{value2: ""})(ctx)()
		assert.Equal(t, result.Of(0), outcome)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Of(func(cfg1 Config1) RIORES.ReaderIOResult[int] {
						return func(ctx context.Context) IOResult[int] {
							return func() Result[int] {
								return result.Of(cfg1.value1 * len(cfg2.value2))
							}
						}
					})
				}
			}
		}

		sequenced := Sequence(original)

		cfg1 := Config1{value1: 3}
		cfg2 := Config2{value2: "test"}

		// Call multiple times with same inputs
		for range 5 {
			outcome := sequenced(cfg1)(cfg2)(ctx)()
			assert.Equal(t, result.Of(12), outcome)
		}
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("swaps parameter order for Reader types", func(t *testing.T) {
		ctx := t.Context()

		// Original: takes Config2, returns ReaderIOResult that may produce Reader[Config1, int]
		original := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, int]] {
			return func(ctx context.Context) IOResult[Reader[Config1, int]] {
				return func() Result[Reader[Config1, int]] {
					return result.Of(func(cfg1 Config1) int {
						return cfg1.value1 + len(cfg2.value2)
					})
				}
			}
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := SequenceReader(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// Test original
		result1 := original(cfg2)(ctx)()
		assert.True(t, result.IsRight(result1))
		innerFunc1, _ := result.Unwrap(result1)
		value1 := innerFunc1(cfg1)
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced(cfg1)
		result2 := innerFunc2(cfg2)(ctx)()
		assert.True(t, result.IsRight(result2))
		value2, _ := result.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")

		original := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, int]] {
			return func(ctx context.Context) IOResult[Reader[Config1, int]] {
				return func() Result[Reader[Config1, int]] {
					return result.Left[Reader[Config1, int]](testErr)
				}
			}
		}

		sequenced := SequenceReader(original)

		outcome := sequenced(Config1{value1: 10})(Config2{value2: "hello"})(ctx)()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("works with pure Reader computations", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, string]] {
			return func(ctx context.Context) IOResult[Reader[Config1, string]] {
				return func() Result[Reader[Config1, string]] {
					if len(cfg2.value2) == 0 {
						return result.Left[Reader[Config1, string]](errors.New("empty string"))
					}
					return result.Of(func(cfg1 Config1) string {
						return fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1)
					})
				}
			}
		}

		sequenced := SequenceReader(original)

		// Test with valid inputs
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of("test:42"), result1)

		// Test with empty string
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, result.IsLeft(result2))
	})

	t.Run("works with zero values", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, int]] {
			return func(ctx context.Context) IOResult[Reader[Config1, int]] {
				return func() Result[Reader[Config1, int]] {
					return result.Of(func(cfg1 Config1) int {
						return cfg1.value1 + len(cfg2.value2)
					})
				}
			}
		}

		sequenced := SequenceReader(original)

		outcome := sequenced(Config1{value1: 0})(Config2{value2: ""})(ctx)()
		assert.Equal(t, result.Of(0), outcome)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, int]] {
			return func(ctx context.Context) IOResult[Reader[Config1, int]] {
				return func() Result[Reader[Config1, int]] {
					return result.Of(func(cfg1 Config1) int {
						return cfg1.value1 * len(cfg2.value2)
					})
				}
			}
		}

		sequenced := SequenceReader(original)

		cfg1 := Config1{value1: 3}
		cfg2 := Config2{value2: "test"}

		// Call multiple times with same inputs
		for range 5 {
			outcome := sequenced(cfg1)(cfg2)(ctx)()
			assert.Equal(t, result.Of(12), outcome)
		}
	})
}

func TestSequenceReaderIO(t *testing.T) {
	t.Run("swaps parameter order for ReaderIO types", func(t *testing.T) {
		ctx := t.Context()

		// Original: takes Config2, returns ReaderIOResult that may produce ReaderIO[Config1, int]
		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, int]] {
				return func() Result[ReaderIO[Config1, int]] {
					return result.Of(func(cfg1 Config1) io.IO[int] {
						return io.Of(cfg1.value1 + len(cfg2.value2))
					})
				}
			}
		}

		// Sequence swaps Config1 and Config2 order
		sequenced := SequenceReaderIO(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// Test original
		result1 := original(cfg2)(ctx)()
		assert.True(t, result.IsRight(result1))
		innerFunc1, _ := result.Unwrap(result1)
		value1 := innerFunc1(cfg1)()
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced(cfg1)
		result2 := innerFunc2(cfg2)(ctx)()
		assert.True(t, result.IsRight(result2))
		value2, _ := result.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, int]] {
				return func() Result[ReaderIO[Config1, int]] {
					return result.Left[ReaderIO[Config1, int]](testErr)
				}
			}
		}

		sequenced := SequenceReaderIO(original)

		outcome := sequenced(Config1{value1: 10})(Config2{value2: "hello"})(ctx)()
		assert.Equal(t, result.Left[int](testErr), outcome)
	})

	t.Run("works with IO effects", func(t *testing.T) {
		ctx := t.Context()
		sideEffect := 0

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, string]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, string]] {
				return func() Result[ReaderIO[Config1, string]] {
					if len(cfg2.value2) == 0 {
						return result.Left[ReaderIO[Config1, string]](errors.New("empty string"))
					}
					return result.Of(func(cfg1 Config1) io.IO[string] {
						return func() string {
							sideEffect = cfg1.value1
							return fmt.Sprintf("%s:%d", cfg2.value2, cfg1.value1)
						}
					})
				}
			}
		}

		sequenced := SequenceReaderIO(original)

		// Test with valid inputs
		sideEffect = 0
		result1 := sequenced(Config1{value1: 42})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of("test:42"), result1)
		assert.Equal(t, 42, sideEffect)

		// Test with empty string
		sideEffect = 0
		result2 := sequenced(Config1{value1: 42})(Config2{value2: ""})(ctx)()
		assert.True(t, result.IsLeft(result2))
		assert.Equal(t, 0, sideEffect) // Side effect should not occur
	})

	t.Run("works with zero values", func(t *testing.T) {
		ctx := t.Context()

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, int]] {
				return func() Result[ReaderIO[Config1, int]] {
					return result.Of(func(cfg1 Config1) io.IO[int] {
						return io.Of(cfg1.value1 + len(cfg2.value2))
					})
				}
			}
		}

		sequenced := SequenceReaderIO(original)

		outcome := sequenced(Config1{value1: 0})(Config2{value2: ""})(ctx)()
		assert.Equal(t, result.Of(0), outcome)
	})

	t.Run("executes IO effects correctly", func(t *testing.T) {
		ctx := t.Context()
		counter := 0

		original := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, int]] {
				return func() Result[ReaderIO[Config1, int]] {
					return result.Of(func(cfg1 Config1) io.IO[int] {
						return func() int {
							counter++
							return cfg1.value1 + len(cfg2.value2)
						}
					})
				}
			}
		}

		sequenced := SequenceReaderIO(original)

		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// Each execution should increment counter
		counter = 0
		result1 := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of(15), result1)
		assert.Equal(t, 1, counter)

		result2 := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of(15), result2)
		assert.Equal(t, 2, counter)
	})
}

func TestTraverse(t *testing.T) {
	t.Run("transforms and swaps parameter order", func(t *testing.T) {
		ctx := t.Context()

		// Original computation depending on Config2
		original := Of[Config2](42)

		// Transformation that introduces Config1 dependency
		transform := func(n int) ReaderReaderIOResult[Config1, string] {
			return func(cfg1 Config1) RIORES.ReaderIOResult[string] {
				return func(ctx context.Context) IOResult[string] {
					return func() Result[string] {
						return result.Of(fmt.Sprintf("value=%d, cfg1=%d", n, cfg1.value1))
					}
				}
			}
		}

		// Apply traverse to swap order and transform
		traversed := Traverse[Config2](transform)(original)

		cfg1 := Config1{value1: 100}
		cfg2 := Config2{value2: "test"}

		outcome := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of("value=42, cfg1=100"), outcome)
	})

	t.Run("preserves error handling in original", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")
		original := Left[Config2, int](testErr)

		transform := func(n int) ReaderReaderIOResult[Config1, string] {
			return Of[Config1](fmt.Sprintf("%d", n))
		}

		traversed := Traverse[Config2](transform)(original)

		outcome := traversed(Config1{value1: 100})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Left[string](testErr), outcome)
	})

	t.Run("preserves error handling in transformation", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](42)
		testErr := errors.New("transform error")

		transform := func(n int) ReaderReaderIOResult[Config1, string] {
			if n < 0 {
				return Left[Config1, string](testErr)
			}
			return Of[Config1](fmt.Sprintf("%d", n))
		}

		// Test with negative value
		originalNeg := Of[Config2](-1)
		traversedNeg := Traverse[Config2](transform)(originalNeg)
		resultNeg := traversedNeg(Config1{value1: 100})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Left[string](testErr), resultNeg)

		// Test with positive value
		traversedPos := Traverse[Config2](transform)(original)
		resultPos := traversedPos(Config1{value1: 100})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of("42"), resultPos)
	})

	t.Run("works with complex transformations", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](10)

		transform := func(n int) ReaderReaderIOResult[Config1, int] {
			return func(cfg1 Config1) RIORES.ReaderIOResult[int] {
				return func(ctx context.Context) IOResult[int] {
					return func() Result[int] {
						return result.Of(n * cfg1.value1)
					}
				}
			}
		}

		traversed := Traverse[Config2](transform)(original)

		outcome := traversed(Config1{value1: 5})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of(50), outcome)
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](10)

		transform := func(n int) ReaderReaderIOResult[Config1, int] {
			return Of[Config1](n * 2)
		}

		outcome := F.Pipe2(
			original,
			Traverse[Config2](transform),
			func(k Kleisli[Config2, Config1, int]) ReaderReaderIOResult[Config2, int] {
				return k(Config1{value1: 5})
			},
		)

		res := outcome(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of(20), res)
	})
}

func TestTraverseReader(t *testing.T) {
	t.Run("transforms with pure Reader and swaps parameter order", func(t *testing.T) {
		ctx := t.Context()

		// Original computation depending on Config2
		original := Of[Config2](100)

		// Pure Reader transformation that introduces Config1 dependency
		formatWithConfig := func(value int) reader.Reader[Config1, string] {
			return func(cfg1 Config1) string {
				return fmt.Sprintf("value=%d, multiplier=%d, result=%d", value, cfg1.value1, value*cfg1.value1)
			}
		}

		// Apply traverse to introduce Config1 and swap order
		traversed := TraverseReader[Config2](formatWithConfig)(original)

		cfg1 := Config1{value1: 5}
		cfg2 := Config2{value2: "test"}

		outcome := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of("value=100, multiplier=5, result=500"), outcome)
	})

	t.Run("preserves error handling", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")
		original := Left[Config2, int](testErr)

		transform := func(n int) reader.Reader[Config1, string] {
			return reader.Of[Config1](fmt.Sprintf("%d", n))
		}

		traversed := TraverseReader[Config2](transform)(original)

		outcome := traversed(Config1{value1: 5})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Left[string](testErr), outcome)
	})

	t.Run("works with pure computations", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](42)

		// Pure transformation using Reader
		double := func(n int) reader.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		traversed := TraverseReader[Config2](double)(original)

		outcome := traversed(Config1{value1: 3})(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of(126), outcome)
	})

	t.Run("works with zero values", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](0)

		transform := func(n int) reader.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n + cfg1.value1
			}
		}

		traversed := TraverseReader[Config2](transform)(original)

		outcome := traversed(Config1{value1: 0})(Config2{value2: ""})(ctx)()
		assert.Equal(t, result.Of(0), outcome)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](10)

		transform := func(n int) reader.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		traversed := TraverseReader[Config2](transform)(original)

		cfg1 := Config1{value1: 5}
		cfg2 := Config2{value2: "test"}

		// Call multiple times with same inputs
		for range 5 {
			outcome := traversed(cfg1)(cfg2)(ctx)()
			assert.Equal(t, result.Of(50), outcome)
		}
	})

	t.Run("can be used in composition", func(t *testing.T) {
		ctx := t.Context()
		original := Of[Config2](10)

		multiply := func(n int) reader.Reader[Config1, int] {
			return func(cfg1 Config1) int {
				return n * cfg1.value1
			}
		}

		outcome := F.Pipe2(
			original,
			TraverseReader[Config2](multiply),
			func(k Kleisli[Config2, Config1, int]) ReaderReaderIOResult[Config2, int] {
				return k(Config1{value1: 3})
			},
		)

		res := outcome(Config2{value2: "test"})(ctx)()
		assert.Equal(t, result.Of(30), res)
	})
}

func TestFlipIntegration(t *testing.T) {
	t.Run("Sequence and Traverse work together", func(t *testing.T) {
		ctx := t.Context()

		// Create a nested computation
		nested := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Of(Of[Config1](len(cfg2.value2)))
				}
			}
		}

		// Sequence it
		sequenced := Sequence(nested)

		// Then traverse with a transformation
		transform := func(n int) ReaderReaderIOResult[Config1, string] {
			return Of[Config1](fmt.Sprintf("length=%d", n))
		}

		// Apply both operations
		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "hello"}

		// First sequence
		intermediate := sequenced(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of(5), intermediate)

		// Then apply traverse on a new computation
		original := Of[Config2](5)
		traversed := Traverse[Config2](transform)(original)
		outcome := traversed(cfg1)(cfg2)(ctx)()
		assert.Equal(t, result.Of("length=5"), outcome)
	})

	t.Run("all flip functions preserve error semantics", func(t *testing.T) {
		ctx := t.Context()
		testErr := errors.New("test error")
		cfg1 := Config1{value1: 10}
		cfg2 := Config2{value2: "test"}

		// Test Sequence with error
		seqErr := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderReaderIOResult[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderReaderIOResult[Config1, int]] {
				return func() Result[ReaderReaderIOResult[Config1, int]] {
					return result.Left[ReaderReaderIOResult[Config1, int]](testErr)
				}
			}
		}
		seqResult := Sequence(seqErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, result.IsLeft(seqResult))

		// Test SequenceReader with error
		seqReaderErr := func(cfg2 Config2) RIORES.ReaderIOResult[Reader[Config1, int]] {
			return func(ctx context.Context) IOResult[Reader[Config1, int]] {
				return func() Result[Reader[Config1, int]] {
					return result.Left[Reader[Config1, int]](testErr)
				}
			}
		}
		seqReaderResult := SequenceReader(seqReaderErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, result.IsLeft(seqReaderResult))

		// Test SequenceReaderIO with error
		seqReaderIOErr := func(cfg2 Config2) RIORES.ReaderIOResult[ReaderIO[Config1, int]] {
			return func(ctx context.Context) IOResult[ReaderIO[Config1, int]] {
				return func() Result[ReaderIO[Config1, int]] {
					return result.Left[ReaderIO[Config1, int]](testErr)
				}
			}
		}
		seqReaderIOResult := SequenceReaderIO(seqReaderIOErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, result.IsLeft(seqReaderIOResult))

		// Test Traverse with error
		travErr := Left[Config2, int](testErr)
		travTransform := func(n int) ReaderReaderIOResult[Config1, string] {
			return Of[Config1](fmt.Sprintf("%d", n))
		}
		travResult := Traverse[Config2](travTransform)(travErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, result.IsLeft(travResult))

		// Test TraverseReader with error
		travReaderErr := Left[Config2, int](testErr)
		travReaderTransform := func(n int) reader.Reader[Config1, string] {
			return reader.Of[Config1](fmt.Sprintf("%d", n))
		}
		travReaderResult := TraverseReader[Config2](travReaderTransform)(travReaderErr)(cfg1)(cfg2)(ctx)()
		assert.True(t, result.IsLeft(travReaderResult))
	})
}
