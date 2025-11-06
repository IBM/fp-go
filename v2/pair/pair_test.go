// Copyright (c) 2024 - 2025 IBM Corp.
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

package pair

import (
	"fmt"
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	SG "github.com/IBM/fp-go/v2/semigroup"
	"github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	p := Of(42)
	assert.Equal(t, 42, Head(p))
	assert.Equal(t, 42, Tail(p))
}

func TestMakePair(t *testing.T) {
	p := MakePair("hello", 42)
	assert.Equal(t, "hello", Head(p))
	assert.Equal(t, 42, Tail(p))
}

func TestFromTuple(t *testing.T) {
	tup := tuple.MakeTuple2("world", 100)
	p := FromTuple(tup)
	assert.Equal(t, "world", Head(p))
	assert.Equal(t, 100, Tail(p))
}

func TestToTuple(t *testing.T) {
	p := MakePair("hello", 42)
	tup := ToTuple(p)
	assert.Equal(t, "hello", tup.F1)
	assert.Equal(t, 42, tup.F2)
}

func TestHeadAndTail(t *testing.T) {
	p := MakePair("test", 123)
	assert.Equal(t, "test", Head(p))
	assert.Equal(t, 123, Tail(p))
}

func TestFirstAndSecond(t *testing.T) {
	p := MakePair("first", "second")
	assert.Equal(t, "first", First(p))
	assert.Equal(t, "second", Second(p))
}

func TestMonadMapHead(t *testing.T) {
	p := MakePair(5, "hello")
	p2 := MonadMapHead(p, func(n int) string {
		return fmt.Sprintf("%d", n)
	})
	assert.Equal(t, "5", Head(p2))
	assert.Equal(t, "hello", Tail(p2))
}

func TestMonadMapTail(t *testing.T) {
	p := MakePair(5, "hello")
	p2 := MonadMapTail(p, func(s string) int {
		return len(s)
	})
	assert.Equal(t, 5, Head(p2))
	assert.Equal(t, 5, Tail(p2))
}

func TestMonadMap(t *testing.T) {
	p := MakePair(10, "test")
	p2 := MonadMap(p, func(n int) string {
		return fmt.Sprintf("value: %d", n)
	})
	assert.Equal(t, "value: 10", Head(p2))
	assert.Equal(t, "test", Tail(p2))
}

func TestMonadBiMap(t *testing.T) {
	p := MakePair(5, "hello")
	p2 := MonadBiMap(p,
		func(n int) string { return fmt.Sprintf("%d", n) },
		func(s string) int { return len(s) },
	)
	assert.Equal(t, "5", Head(p2))
	assert.Equal(t, 5, Tail(p2))
}

func TestMapHead(t *testing.T) {
	mapper := MapHead[string](func(n int) string {
		return fmt.Sprintf("%d", n)
	})
	p := MakePair(42, "world")
	p2 := mapper(p)
	assert.Equal(t, "42", Head(p2))
	assert.Equal(t, "world", Tail(p2))
}

func TestMapTail(t *testing.T) {
	mapper := MapTail[int](func(s string) int {
		return len(s)
	})
	p := MakePair(10, "hello")
	p2 := mapper(p)
	assert.Equal(t, 10, Head(p2))
	assert.Equal(t, 5, Tail(p2))
}

func TestMap(t *testing.T) {
	mapper := Map[int](func(s string) int {
		return len(s)
	})
	p := MakePair(10, "test")
	p2 := mapper(p)
	assert.Equal(t, 10, Head(p2))
	assert.Equal(t, 4, Tail(p2))
}

func TestBiMap(t *testing.T) {
	mapper := BiMap(
		func(n int) string { return fmt.Sprintf("n=%d", n) },
		func(s string) int { return len(s) },
	)
	p := MakePair(7, "hello")
	p2 := mapper(p)
	assert.Equal(t, "n=7", Head(p2))
	assert.Equal(t, 5, Tail(p2))
}

func TestSwap(t *testing.T) {
	p := MakePair("hello", 42)
	swapped := Swap(p)
	assert.Equal(t, 42, Head(swapped))
	assert.Equal(t, "hello", Tail(swapped))
}

func TestMonadChainHead(t *testing.T) {
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
	p := MakePair(5, "hello")
	p2 := MonadChainHead(strConcat, p, func(n int) Pair[string, string] {
		return MakePair(fmt.Sprintf("%d", n), "!")
	})
	assert.Equal(t, "5", Head(p2))
	assert.Equal(t, "hello!", Tail(p2))
}

func TestMonadChainTail(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	p := MakePair(5, "hello")
	p2 := MonadChainTail(intSum, p, func(s string) Pair[int, int] {
		return MakePair(len(s), len(s)*2)
	})
	assert.Equal(t, 10, Head(p2)) // 5 + 5
	assert.Equal(t, 10, Tail(p2))
}

func TestMonadChain(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	p := MakePair(3, "test")
	p2 := MonadChain(intSum, p, func(s string) Pair[int, int] {
		return MakePair(len(s), len(s)*3)
	})
	assert.Equal(t, 7, Head(p2)) // 3 + 4
	assert.Equal(t, 12, Tail(p2))
}

func TestChainHead(t *testing.T) {
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
	chain := ChainHead(strConcat, func(n int) Pair[string, string] {
		return MakePair(fmt.Sprintf("%d", n), "!")
	})
	p := MakePair(42, "hello")
	p2 := chain(p)
	assert.Equal(t, "42", Head(p2))
	assert.Equal(t, "hello!", Tail(p2))
}

func TestChainTail(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	chain := ChainTail(intSum, func(s string) Pair[int, int] {
		return MakePair(len(s), len(s)*2)
	})
	p := MakePair(10, "world")
	p2 := chain(p)
	assert.Equal(t, 15, Head(p2)) // 10 + 5
	assert.Equal(t, 10, Tail(p2))
}

func TestChain(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	chain := Chain(intSum, func(s string) Pair[int, int] {
		return MakePair(len(s), len(s)*2)
	})
	p := MakePair(5, "hi")
	p2 := chain(p)
	assert.Equal(t, 7, Head(p2)) // 5 + 2
	assert.Equal(t, 4, Tail(p2))
}

func TestMonadApHead(t *testing.T) {
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
	pf := MakePair(func(n int) string { return fmt.Sprintf("%d", n) }, "!")
	pv := MakePair(42, "hello")
	result := MonadApHead(strConcat, pf, pv)
	assert.Equal(t, "42", Head(result))
	assert.Equal(t, "hello!", Tail(result))
}

func TestMonadApTail(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	pf := MakePair(10, func(s string) int { return len(s) })
	pv := MakePair(5, "hello")
	result := MonadApTail(intSum, pf, pv)
	assert.Equal(t, 15, Head(result)) // 5 + 10
	assert.Equal(t, 5, Tail(result))
}

func TestMonadAp(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	pf := MakePair(7, func(s string) int { return len(s) * 2 })
	pv := MakePair(3, "test")
	result := MonadAp(intSum, pf, pv)
	assert.Equal(t, 10, Head(result)) // 3 + 7
	assert.Equal(t, 8, Tail(result))  // len("test") * 2
}

func TestApHead(t *testing.T) {
	strConcat := SG.MakeSemigroup(func(a, b string) string { return a + b })
	pv := MakePair(100, "world")
	ap := ApHead[string, int, string](strConcat, pv)
	pf := MakePair(func(n int) string { return fmt.Sprintf("num=%d", n) }, "!")
	result := ap(pf)
	assert.Equal(t, "num=100", Head(result))
	assert.Equal(t, "world!", Tail(result))
}

func TestApTail(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	pv := MakePair(20, "hello")
	ap := ApTail[int, string, int](intSum, pv)
	pf := MakePair(5, func(s string) int { return len(s) })
	result := ap(pf)
	assert.Equal(t, 25, Head(result)) // 20 + 5
	assert.Equal(t, 5, Tail(result))
}

func TestAp(t *testing.T) {
	intSum := N.SemigroupSum[int]()
	pv := MakePair(15, "test")
	ap := Ap[int, string, int](intSum, pv)
	pf := MakePair(10, func(s string) int { return len(s) * 3 })
	result := ap(pf)
	assert.Equal(t, 25, Head(result)) // 15 + 10
	assert.Equal(t, 12, Tail(result)) // len("test") * 3
}

func TestPaired(t *testing.T) {
	add := func(a, b int) int { return a + b }
	pairedAdd := Paired(add)
	result := pairedAdd(MakePair(3, 4))
	assert.Equal(t, 7, result)
}

func TestUnpaired(t *testing.T) {
	pairedAdd := func(p Pair[int, int]) int {
		return Head(p) + Tail(p)
	}
	add := Unpaired(pairedAdd)
	result := add(5, 7)
	assert.Equal(t, 12, result)
}

func TestMerge(t *testing.T) {
	add := func(b int) func(a int) int {
		return func(a int) int { return a + b }
	}
	merge := Merge(add)
	result := merge(MakePair(3, 4))
	assert.Equal(t, 7, result)
}

func TestEq(t *testing.T) {
	pairEq := Eq(
		EQ.FromStrictEquals[string](),
		EQ.FromStrictEquals[int](),
	)
	p1 := MakePair("hello", 42)
	p2 := MakePair("hello", 42)
	p3 := MakePair("world", 42)
	p4 := MakePair("hello", 100)

	assert.True(t, pairEq.Equals(p1, p2))
	assert.False(t, pairEq.Equals(p1, p3))
	assert.False(t, pairEq.Equals(p1, p4))
}

func TestFromStrictEquals(t *testing.T) {
	pairEq := FromStrictEquals[string, int]()
	p1 := MakePair("test", 123)
	p2 := MakePair("test", 123)
	p3 := MakePair("test", 456)

	assert.True(t, pairEq.Equals(p1, p2))
	assert.False(t, pairEq.Equals(p1, p3))
}

func TestString(t *testing.T) {
	p := MakePair("hello", 42)
	str := p.String()
	assert.Contains(t, str, "Pair")
	assert.Contains(t, str, "hello")
	assert.Contains(t, str, "42")
}

func TestFormat(t *testing.T) {
	p := MakePair("test", 100)
	str := fmt.Sprintf("%s", p)
	assert.Contains(t, str, "Pair")
	assert.Contains(t, str, "test")
	assert.Contains(t, str, "100")
}

func TestMonadHead(t *testing.T) {
	stringMonoid := M.MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)
	monad := MonadHead[int, string, string](stringMonoid)

	// Test Of
	p := monad.Of(42)
	assert.Equal(t, 42, Head(p))
	assert.Equal(t, "", Tail(p))

	// Test Map
	mapper := monad.Map(func(n int) string { return fmt.Sprintf("%d", n) })
	p2 := mapper(MakePair(100, "!"))
	assert.Equal(t, "100", Head(p2))
	assert.Equal(t, "!", Tail(p2))

	// Test Chain
	chain := monad.Chain(func(n int) Pair[string, string] {
		return MakePair(fmt.Sprintf("n=%d", n), "!")
	})
	p3 := chain(MakePair(7, "hello"))
	assert.Equal(t, "n=7", Head(p3))
	assert.Equal(t, "hello!", Tail(p3))

	// Test Ap
	pv := MakePair(5, "world")
	ap := monad.Ap(pv)
	pf := MakePair(func(n int) string { return fmt.Sprintf("%d", n*2) }, "!")
	p4 := ap(pf)
	assert.Equal(t, "10", Head(p4))
	assert.Equal(t, "world!", Tail(p4))
}

func TestPointedHead(t *testing.T) {
	stringMonoid := M.MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)
	pointed := PointedHead[int, string](stringMonoid)
	p := pointed.Of(42)
	assert.Equal(t, 42, Head(p))
	assert.Equal(t, "", Tail(p))
}

func TestFunctorHead(t *testing.T) {
	functor := FunctorHead[int, string, string]()
	mapper := functor.Map(func(n int) string { return fmt.Sprintf("value=%d", n) })
	p := MakePair(42, "test")
	p2 := mapper(p)
	assert.Equal(t, "value=42", Head(p2))
	assert.Equal(t, "test", Tail(p2))
}

func TestApplicativeHead(t *testing.T) {
	stringMonoid := M.MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)
	applicative := ApplicativeHead[int, string, string](stringMonoid)

	// Test Of
	p := applicative.Of(100)
	assert.Equal(t, 100, Head(p))
	assert.Equal(t, "", Tail(p))

	// Test Map
	mapper := applicative.Map(func(n int) string { return fmt.Sprintf("%d", n) })
	p2 := mapper(MakePair(42, "!"))
	assert.Equal(t, "42", Head(p2))
	assert.Equal(t, "!", Tail(p2))

	// Test Ap
	pv := MakePair(7, "hello")
	ap := applicative.Ap(pv)
	pf := MakePair(func(n int) string { return fmt.Sprintf("n=%d", n) }, "!")
	p3 := ap(pf)
	assert.Equal(t, "n=7", Head(p3))
	assert.Equal(t, "hello!", Tail(p3))
}

func TestMonadTail(t *testing.T) {
	intSum := N.MonoidSum[int]()
	monad := MonadTail[string, int, int](intSum)

	// Test Of
	p := monad.Of("hello")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "hello", Tail(p))

	// Test Map
	mapper := monad.Map(func(s string) int { return len(s) })
	p2 := mapper(MakePair(5, "world"))
	assert.Equal(t, 5, Head(p2))
	assert.Equal(t, 5, Tail(p2))

	// Test Chain
	chain := monad.Chain(func(s string) Pair[int, int] {
		return MakePair(len(s), len(s)*2)
	})
	p3 := chain(MakePair(10, "test"))
	assert.Equal(t, 14, Head(p3)) // 10 + 4
	assert.Equal(t, 8, Tail(p3))

	// Test Ap
	pv := MakePair(5, "hello")
	ap := monad.Ap(pv)
	pf := MakePair(10, func(s string) int { return len(s) })
	p4 := ap(pf)
	assert.Equal(t, 15, Head(p4)) // 5 + 10
	assert.Equal(t, 5, Tail(p4))
}

func TestPointedTail(t *testing.T) {
	intSum := N.MonoidSum[int]()
	pointed := PointedTail[string, int](intSum)
	p := pointed.Of("test")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "test", Tail(p))
}

func TestFunctorTail(t *testing.T) {
	functor := FunctorTail[string, int, int]()
	mapper := functor.Map(func(s string) int { return len(s) * 2 })
	p := MakePair(10, "hello")
	p2 := mapper(p)
	assert.Equal(t, 10, Head(p2))
	assert.Equal(t, 10, Tail(p2))
}

func TestApplicativeTail(t *testing.T) {
	intSum := N.MonoidSum[int]()
	applicative := ApplicativeTail[string, int, int](intSum)

	// Test Of
	p := applicative.Of("world")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "world", Tail(p))

	// Test Map
	mapper := applicative.Map(func(s string) int { return len(s) })
	p2 := mapper(MakePair(5, "test"))
	assert.Equal(t, 5, Head(p2))
	assert.Equal(t, 4, Tail(p2))

	// Test Ap
	pv := MakePair(10, "hello")
	ap := applicative.Ap(pv)
	pf := MakePair(5, func(s string) int { return len(s) * 2 })
	p3 := ap(pf)
	assert.Equal(t, 15, Head(p3)) // 10 + 5
	assert.Equal(t, 10, Tail(p3))
}

func TestMonad(t *testing.T) {
	intSum := N.MonoidSum[int]()
	monad := Monad[string, int, int](intSum)

	p := monad.Of("test")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "test", Tail(p))
}

func TestPointed(t *testing.T) {
	intSum := N.MonoidSum[int]()
	pointed := Pointed[string, int](intSum)

	p := pointed.Of("hello")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "hello", Tail(p))
}

func TestFunctor(t *testing.T) {
	functor := Functor[string, int, int]()
	mapper := functor.Map(func(s string) int { return len(s) })
	p := MakePair(7, "world")
	p2 := mapper(p)
	assert.Equal(t, 7, Head(p2))
	assert.Equal(t, 5, Tail(p2))
}

func TestApplicative(t *testing.T) {
	intSum := N.MonoidSum[int]()
	applicative := Applicative[string, int, int](intSum)

	p := applicative.Of("test")
	assert.Equal(t, 0, Head(p))
	assert.Equal(t, "test", Tail(p))
}

// Test edge cases and complex scenarios
func TestComplexChaining(t *testing.T) {
	intSum := N.SemigroupSum[int]()

	// Chain multiple operations
	p := MakePair(1, "a")
	p2 := MonadChainTail(intSum, p, func(s string) Pair[int, string] {
		return MakePair(len(s), s+"b")
	})
	p3 := MonadChainTail(intSum, p2, func(s string) Pair[int, string] {
		return MakePair(len(s), s+"c")
	})

	assert.Equal(t, 4, Head(p3)) // 1 + 1 + 2
	assert.Equal(t, "abc", Tail(p3))
}

func TestBiMapWithDifferentTypes(t *testing.T) {
	p := MakePair(3.14, true)
	p2 := MonadBiMap(p,
		func(f float64) int { return int(f * 10) },
		func(b bool) string {
			if b {
				return "yes"
			}
			return "no"
		},
	)
	assert.Equal(t, 31, Head(p2))
	assert.Equal(t, "yes", Tail(p2))
}

func TestSwapTwice(t *testing.T) {
	p := MakePair("original", 999)
	swapped := Swap(p)
	swappedBack := Swap(swapped)
	assert.Equal(t, "original", Head(swappedBack))
	assert.Equal(t, 999, Tail(swappedBack))
}
