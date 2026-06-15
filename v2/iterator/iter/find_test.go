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

package iter

import (
	"strconv"
	"strings"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestFindFirst_Success(t *testing.T) {
	t.Run("finds first element matching predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := FindFirst(func(x int) bool { return x > 3 })(seq)()
		assert.Equal(t, O.Some(4), result)
	})

	t.Run("finds first element at start of sequence", func(t *testing.T) {
		seq := From(10, 2, 3, 4, 5)
		result := FindFirst(func(x int) bool { return x > 5 })(seq)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("finds first element at end of sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 10)
		result := FindFirst(func(x int) bool { return x > 5 })(seq)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("finds first with string predicate", func(t *testing.T) {
		seq := From("apple", "banana", "cherry", "date")
		result := FindFirst(func(s string) bool { return len(s) > 5 })(seq)()
		assert.Equal(t, O.Some("banana"), result)
	})

	t.Run("finds first even number", func(t *testing.T) {
		seq := From(1, 3, 5, 6, 7, 8)
		result := FindFirst(func(x int) bool { return x%2 == 0 })(seq)()
		assert.Equal(t, O.Some(6), result)
	})
}

func TestFindFirst_Failure(t *testing.T) {
	t.Run("returns None when no element matches", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := FindFirst(func(x int) bool { return x > 10 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := FindFirst(func(x int) bool { return x > 0 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None when all elements fail predicate", func(t *testing.T) {
		seq := From("a", "b", "c")
		result := FindFirst(func(s string) bool { return len(s) > 5 })(seq)()
		assert.Equal(t, O.None[string](), result)
	})
}

func TestFindFirst_EdgeCases(t *testing.T) {
	t.Run("handles single element sequence matching", func(t *testing.T) {
		seq := From(42)
		result := FindFirst(func(x int) bool { return x > 0 })(seq)()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("handles single element sequence not matching", func(t *testing.T) {
		seq := From(42)
		result := FindFirst(func(x int) bool { return x < 0 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("stops at first match without consuming rest", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(int) bool) {
			for i := range 10 {
				consumed++
				if !yield(i) {
					return
				}
			}
		}
		result := FindFirst(func(x int) bool { return x == 3 })(seq)()
		assert.Equal(t, O.Some(3), result)
		// First will consume up to and including the match
		assert.Equal(t, 4, consumed)
	})

	t.Run("works with zero values", func(t *testing.T) {
		seq := From(1, 0, 2, 3)
		result := FindFirst(func(x int) bool { return x == 0 })(seq)()
		assert.Equal(t, O.Some(0), result)
	})
}

func TestFindFirstMap_Success(t *testing.T) {
	t.Run("finds and transforms first matching element", func(t *testing.T) {
		parsePositive := func(s string) O.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil || n <= 0 {
				return O.None[int]()
			}
			return O.Some(n)
		}
		seq := From("invalid", "0", "42", "100")
		result := FindFirstMap(parsePositive)(seq)()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("finds first valid transformation", func(t *testing.T) {
		toUpper := func(s string) O.Option[string] {
			if len(s) > 3 {
				return O.Some(strings.ToUpper(s))
			}
			return O.None[string]()
		}
		seq := From("a", "ab", "abc", "abcd", "abcde")
		result := FindFirstMap(toUpper)(seq)()
		assert.Equal(t, O.Some("ABCD"), result)
	})

	t.Run("transforms with different output type", func(t *testing.T) {
		strlen := func(s string) O.Option[int] {
			if s != "" {
				return O.Some(len(s))
			}
			return O.None[int]()
		}
		seq := From("", "", "hello", "world")
		result := FindFirstMap(strlen)(seq)()
		assert.Equal(t, O.Some(5), result)
	})
}

func TestFindFirstMap_Failure(t *testing.T) {
	t.Run("returns None when no element maps successfully", func(t *testing.T) {
		parsePositive := func(s string) O.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil || n <= 0 {
				return O.None[int]()
			}
			return O.Some(n)
		}
		seq := From("invalid", "0", "-5")
		result := FindFirstMap(parsePositive)(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for empty sequence", func(t *testing.T) {
		mapper := func(x int) O.Option[string] {
			return O.Some(strconv.Itoa(x))
		}
		seq := Empty[int]()
		result := FindFirstMap(mapper)(seq)()
		assert.Equal(t, O.None[string](), result)
	})
}

func TestFindLast_Success(t *testing.T) {
	t.Run("finds last element matching predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := FindLast(func(x int) bool { return x > 3 })(seq)()
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("finds last element when multiple match", func(t *testing.T) {
		seq := From(2, 4, 6, 8, 10)
		result := FindLast(func(x int) bool { return x%2 == 0 })(seq)()
		assert.Equal(t, O.Some(10), result)
	})

	t.Run("finds last with string predicate", func(t *testing.T) {
		seq := From("apple", "banana", "cherry", "date")
		result := FindLast(func(s string) bool { return len(s) > 5 })(seq)()
		assert.Equal(t, O.Some("cherry"), result)
	})

	t.Run("finds last when only one matches", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 100)
		result := FindLast(func(x int) bool { return x > 50 })(seq)()
		assert.Equal(t, O.Some(100), result)
	})
}

func TestFindLast_Failure(t *testing.T) {
	t.Run("returns None when no element matches", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := FindLast(func(x int) bool { return x > 10 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := FindLast(func(x int) bool { return x > 0 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None when all elements fail predicate", func(t *testing.T) {
		seq := From("a", "b", "c")
		result := FindLast(func(s string) bool { return len(s) > 5 })(seq)()
		assert.Equal(t, O.None[string](), result)
	})
}

func TestFindLast_EdgeCases(t *testing.T) {
	t.Run("handles single element sequence matching", func(t *testing.T) {
		seq := From(42)
		result := FindLast(func(x int) bool { return x > 0 })(seq)()
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("handles single element sequence not matching", func(t *testing.T) {
		seq := From(42)
		result := FindLast(func(x int) bool { return x < 0 })(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("consumes entire sequence", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(int) bool) {
			for i := range 10 {
				consumed++
				if !yield(i) {
					return
				}
			}
		}
		result := FindLast(func(x int) bool { return x < 5 })(seq)()
		assert.Equal(t, O.Some(4), result)
		assert.Equal(t, 10, consumed, "should consume entire sequence")
	})

	t.Run("works with zero values", func(t *testing.T) {
		seq := From(1, 2, 0, 3)
		result := FindLast(func(x int) bool { return x == 0 })(seq)()
		assert.Equal(t, O.Some(0), result)
	})
}

func TestFindLastMap_Success(t *testing.T) {
	t.Run("finds and transforms last matching element", func(t *testing.T) {
		parsePositive := func(s string) O.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil || n <= 0 {
				return O.None[int]()
			}
			return O.Some(n)
		}
		seq := From("invalid", "42", "100", "0")
		result := FindLastMap(parsePositive)(seq)()
		assert.Equal(t, O.Some(100), result)
	})

	t.Run("finds last valid transformation", func(t *testing.T) {
		toUpper := func(s string) O.Option[string] {
			if len(s) > 3 {
				return O.Some(strings.ToUpper(s))
			}
			return O.None[string]()
		}
		seq := From("abcd", "abc", "abcde", "ab")
		result := FindLastMap(toUpper)(seq)()
		assert.Equal(t, O.Some("ABCDE"), result)
	})

	t.Run("transforms with different output type", func(t *testing.T) {
		strlen := func(s string) O.Option[int] {
			if s != "" {
				return O.Some(len(s))
			}
			return O.None[int]()
		}
		seq := From("hello", "", "world", "")
		result := FindLastMap(strlen)(seq)()
		assert.Equal(t, O.Some(5), result)
	})
}

func TestFindLastMap_Failure(t *testing.T) {
	t.Run("returns None when no element maps successfully", func(t *testing.T) {
		parsePositive := func(s string) O.Option[int] {
			n, err := strconv.Atoi(s)
			if err != nil || n <= 0 {
				return O.None[int]()
			}
			return O.Some(n)
		}
		seq := From("invalid", "0", "-5")
		result := FindLastMap(parsePositive)(seq)()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for empty sequence", func(t *testing.T) {
		mapper := func(x int) O.Option[string] {
			return O.Some(strconv.Itoa(x))
		}
		seq := Empty[int]()
		result := FindLastMap(mapper)(seq)()
		assert.Equal(t, O.None[string](), result)
	})
}

func TestFind_Integration(t *testing.T) {
	t.Run("FindFirst vs FindLast with same predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		pred := func(x int) bool { return x > 5 }

		first := FindFirst(pred)(seq)()
		last := FindLast(pred)(seq)()

		assert.Equal(t, O.Some(6), first)
		assert.Equal(t, O.Some(10), last)
	})

	t.Run("FindFirstMap vs FindLastMap with same mapper", func(t *testing.T) {
		mapper := func(x int) O.Option[int] {
			if x%2 == 0 {
				return O.Some(x * 2)
			}
			return O.None[int]()
		}
		seq := From(1, 2, 3, 4, 5, 6)

		first := FindFirstMap(mapper)(seq)()
		last := FindLastMap(mapper)(seq)()

		assert.Equal(t, O.Some(4), first)
		assert.Equal(t, O.Some(12), last)
	})

	t.Run("chaining with other operations", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		doubled := Map(func(x int) int { return x * 2 })(seq)
		result := FindFirst(func(x int) bool { return x > 10 })(doubled)()
		assert.Equal(t, O.Some(12), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		seq := From(
			User{"Alice", 25},
			User{"Bob", 30},
			User{"Charlie", 35},
		)
		result := FindFirst(func(u User) bool { return u.Age > 28 })(seq)()
		assert.Equal(t, O.Some(User{"Bob", 30}), result)
	})
}

func TestFind_WithComplexTypes(t *testing.T) {
	t.Run("FindFirst with struct type", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		seq := From(
			Point{1, 2},
			Point{3, 4},
			Point{5, 6},
		)
		result := FindFirst(func(p Point) bool { return p.X > 2 })(seq)()
		assert.Equal(t, O.Some(Point{3, 4}), result)
	})

	t.Run("FindFirstMap with struct transformation", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		extractAdult := func(p Person) O.Option[string] {
			if p.Age >= 18 {
				return O.Some(p.Name)
			}
			return O.None[string]()
		}
		seq := From(
			Person{"Alice", 15},
			Person{"Bob", 20},
			Person{"Charlie", 17},
		)
		result := FindFirstMap(extractAdult)(seq)()
		assert.Equal(t, O.Some("Bob"), result)
	})
}

func BenchmarkFindFirst(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	finder := FindFirst(func(x int) bool { return x > 5 })

	b.ResetTimer()
	for range b.N {
		_ = finder(seq)()
	}
}

func BenchmarkFindLast(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	finder := FindLast(func(x int) bool { return x > 5 })

	b.ResetTimer()
	for range b.N {
		_ = finder(seq)()
	}
}

func BenchmarkFindFirstMap(b *testing.B) {
	seq := From("1", "2", "3", "4", "5", "6", "7", "8", "9", "10")
	mapper := func(s string) O.Option[int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return O.None[int]()
		}
		return O.Some(n)
	}
	finder := FindFirstMap(mapper)

	b.ResetTimer()
	for range b.N {
		_ = finder(seq)()
	}
}

func BenchmarkFindLastMap(b *testing.B) {
	seq := From("1", "2", "3", "4", "5", "6", "7", "8", "9", "10")
	mapper := func(s string) O.Option[int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return O.None[int]()
		}
		return O.Some(n)
	}
	finder := FindLastMap(mapper)

	b.ResetTimer()
	for range b.N {
		_ = finder(seq)()
	}
}
