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

package io

import (
	"context"
	"sync"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

// Test FromIO
func TestFromIO(t *testing.T) {
	io := Of(42)
	result := FromIO(io)
	assert.Equal(t, 42, result())
}

// Test MonadOf
func TestMonadOf(t *testing.T) {
	result := MonadOf(42)
	assert.Equal(t, 42, result())
}

// Test MonadMapTo
func TestMonadMapTo(t *testing.T) {
	result := MonadMapTo(Of(1), "hello")
	assert.Equal(t, "hello", result())
}

// Test MapTo
func TestMapTo(t *testing.T) {
	result := F.Pipe1(Of(1), MapTo[int]("hello"))
	assert.Equal(t, "hello", result())
}

// Test MonadApSeq
func TestMonadApSeq(t *testing.T) {
	f := Of(N.Mul(2))
	result := MonadApSeq(f, Of(21))
	assert.Equal(t, 42, result())
}

// Test ApPar
func TestApPar(t *testing.T) {
	f := Of(N.Mul(2))
	result := F.Pipe1(f, ApPar[int](Of(21)))
	assert.Equal(t, 42, result())
}

// Test MonadChainFirst
func TestMonadChainFirst(t *testing.T) {
	sideEffect := 0
	result := MonadChainFirst(Of(42), func(x int) IO[string] {
		sideEffect = x
		return Of("ignored")
	})
	assert.Equal(t, 42, result())
	assert.Equal(t, 42, sideEffect)
}

// Test ChainFirst
func TestChainFirst(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(Of(42), ChainFirst(func(x int) IO[string] {
		sideEffect = x
		return Of("ignored")
	}))
	assert.Equal(t, 42, result())
	assert.Equal(t, 42, sideEffect)
}

// Test MonadApFirst
func TestMonadApFirst(t *testing.T) {
	result := MonadApFirst(Of(42), Of("ignored"))
	assert.Equal(t, 42, result())
}

// Test MonadApSecond
func TestMonadApSecond(t *testing.T) {
	result := MonadApSecond(Of("ignored"), Of(42))
	assert.Equal(t, 42, result())
}

// Test MonadChainTo
func TestMonadChainTo(t *testing.T) {
	result := MonadChainTo(Of(1), Of(42))
	assert.Equal(t, 42, result())
}

// Test ChainTo
func TestChainTo(t *testing.T) {
	result := F.Pipe1(Of(1), ChainTo[int](Of(42)))
	assert.Equal(t, 42, result())
}

// Test Defer
func TestDefer(t *testing.T) {
	counter := 0
	deferred := Defer(func() IO[int] {
		counter++
		return Of(counter)
	})

	assert.Equal(t, 1, deferred())
	assert.Equal(t, 2, deferred())
	assert.Equal(t, 3, deferred())
}

// Test MonadFlap
func TestMonadFlap(t *testing.T) {
	f := Of(N.Mul(2))
	result := MonadFlap(f, 21)
	assert.Equal(t, 42, result())
}

// Test Flap
func TestFlap(t *testing.T) {
	f := Of(N.Mul(2))
	result := F.Pipe1(f, Flap[int](21))
	assert.Equal(t, 42, result())
}

// Test After
func TestAfter(t *testing.T) {
	future := time.Now().Add(50 * time.Millisecond)
	start := time.Now()
	result := F.Pipe1(Of(42), After[int](future))
	value := result()
	elapsed := time.Since(start)

	assert.Equal(t, 42, value)
	assert.True(t, elapsed >= 50*time.Millisecond)
}

// Test WithTime
func TestWithTime(t *testing.T) {
	result := WithTime(Of(42))
	tuple := result()

	assert.Equal(t, 42, pair.Tail(tuple))
	rg := pair.Head(tuple)
	assert.True(t, pair.Head(rg).Before(pair.Tail(rg)) || pair.Head(rg).Equal(pair.Tail(rg)))
}

// Test WithDuration
func TestWithDuration(t *testing.T) {
	result := WithDuration(func() int {
		time.Sleep(10 * time.Millisecond)
		return 42
	})
	tuple := result()

	assert.Equal(t, 42, pair.Tail(tuple))
	assert.True(t, pair.Head(tuple) >= 10*time.Millisecond)
}

// Test Let
func TestLet(t *testing.T) {
	type State struct {
		value int
	}

	result := F.Pipe2(
		Of(State{value: 10}),
		Let(func(doubled int) func(s State) State {
			return func(s State) State {
				s.value = doubled
				return s
			}
		}, func(s State) int {
			return s.value * 2
		}),
		Map(func(s State) int { return s.value }),
	)

	assert.Equal(t, 20, result())
}

// Test LetTo
func TestLetTo(t *testing.T) {
	type State struct {
		value int
	}

	result := F.Pipe2(
		Of(State{value: 10}),
		LetTo(func(newVal int) func(s State) State {
			return func(s State) State {
				s.value = newVal
				return s
			}
		}, 42),
		Map(func(s State) int { return s.value }),
	)

	assert.Equal(t, 42, result())
}

// Test BindTo
func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	result := F.Pipe2(
		Of(42),
		BindTo(func(v int) State {
			return State{value: v}
		}),
		Map(func(s State) int { return s.value }),
	)

	assert.Equal(t, 42, result())
}

// Test Bracket
func TestBracket(t *testing.T) {
	acquired := false
	released := false

	acquire := func() int {
		acquired = true
		return 42
	}

	use := func(x int) IO[int] {
		return Of(x * 2)
	}

	release := func(x int, result int) IO[Void] {
		return FromImpure(func() {
			released = true
		})
	}

	result := Bracket(Of(acquire()), use, release)
	value := result()

	assert.Equal(t, 84, value)
	assert.True(t, acquired)
	assert.True(t, released)
}

// Test WithResource
func TestWithResource(t *testing.T) {
	acquired := false
	released := false

	onCreate := func() int {
		acquired = true
		return 42
	}

	onRelease := func(x int) IO[Void] {
		return FromImpure(func() {
			released = true
		})
	}

	withRes := WithResource[int, int](Of(onCreate()), onRelease)

	result := withRes(func(x int) IO[int] {
		return Of(x * 2)
	})

	value := result()

	assert.Equal(t, 84, value)
	assert.True(t, acquired)
	assert.True(t, released)
}

// Test WithLock - simplified test
func TestWithLock(t *testing.T) {
	var mu sync.Mutex
	counter := 0

	// Create a lock IO that acquires the lock and returns a release function
	lock := func() func() {
		mu.Lock()
		return func() { mu.Unlock() }
	}

	// Create operation that increments counter
	operation := func() int {
		counter++
		return counter
	}

	// Apply lock to operation - cast to context.CancelFunc
	var lockIO IO[context.CancelFunc] = func() context.CancelFunc {
		return lock()
	}
	safeOp := WithLock[int](lockIO)(operation)

	// Run the operation
	result := safeOp()

	assert.Equal(t, 1, result)
	assert.Equal(t, 1, counter)
}

// Test Printf
func TestPrintf(t *testing.T) {
	result := F.Pipe1(Of(42), ChainFirst(Printf[int]("Value: %d\n")))
	assert.Equal(t, 42, result())
}

// Test ApplySemigroup
func TestApplySemigroup(t *testing.T) {
	intAdd := N.MonoidSum[int]()
	ioSemigroup := ApplySemigroup(intAdd)

	result := ioSemigroup.Concat(Of(10), Of(32))
	assert.Equal(t, 42, result())
}

// Test ApplicativeMonoid
func TestApplicativeMonoid(t *testing.T) {
	intAdd := N.MonoidSum[int]()
	ioMonoid := ApplicativeMonoid(intAdd)

	result := ioMonoid.Concat(Of(10), Of(32))
	assert.Equal(t, 42, result())

	empty := ioMonoid.Empty()
	assert.Equal(t, 0, empty())
}

// Test Applicative type class
func TestApplicativeTypeClass(t *testing.T) {
	app := Applicative[int, int]()

	// Test Of
	io1 := app.Of(21)
	assert.Equal(t, 21, io1())

	// Test Map
	io2 := app.Map(N.Mul(2))(io1)
	assert.Equal(t, 42, io2())
}

// Test Monad type class
func TestMonadTypeClass(t *testing.T) {
	m := Monad[int, int]()

	result := F.Pipe2(
		m.Of(21),
		m.Chain(func(x int) IO[int] {
			return m.Of(x * 2)
		}),
		m.Map(N.Add(1)),
	)

	assert.Equal(t, 43, result())
}

// Test TraverseArrayWithIndex
func TestTraverseArrayWithIndex(t *testing.T) {
	src := []string{"a", "b", "c"}

	result := F.Pipe1(
		src,
		TraverseArrayWithIndex(func(i int, s string) IO[string] {
			return Of(F.Pipe1(i, func(idx int) string {
				return s + string(rune('0'+idx))
			}))
		}),
	)

	assert.Equal(t, []string{"a0", "b1", "c2"}, result())
}

// Test TraverseRecord
func TestTraverseRecord(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := F.Pipe1(
		src,
		TraverseRecord[string](func(x int) IO[int] {
			return Of(x * 2)
		}),
	)

	values := result()
	assert.Equal(t, 2, values["a"])
	assert.Equal(t, 4, values["b"])
}

// Test TraverseRecordWithIndex
func TestTraverseRecordWithIndex(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := F.Pipe1(
		src,
		TraverseRecordWithIndex(func(k string, x int) IO[string] {
			return Of(k)
		}),
	)

	values := result()
	assert.Equal(t, "a", values["a"])
	assert.Equal(t, "b", values["b"])
}

// Test SequenceRecord
func TestSequenceRecord(t *testing.T) {
	src := map[string]IO[int]{
		"a": Of(1),
		"b": Of(2),
	}

	result := SequenceRecord(src)
	values := result()

	assert.Equal(t, 1, values["a"])
	assert.Equal(t, 2, values["b"])
}

// Test MonadTraverseRecord
func TestMonadTraverseRecord(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := MonadTraverseRecord(src, func(x int) IO[int] {
		return Of(x * 2)
	})

	values := result()
	assert.Equal(t, 2, values["a"])
	assert.Equal(t, 4, values["b"])
}

// Test TraverseArrayWithIndexSeq
func TestTraverseArrayWithIndexSeq(t *testing.T) {
	var order []int
	src := []int{1, 2, 3}

	result := F.Pipe1(
		src,
		TraverseArrayWithIndexSeq(func(i int, x int) IO[int] {
			return func() int {
				order = append(order, i)
				return x * 2
			}
		}),
	)

	values := result()
	assert.Equal(t, []int{2, 4, 6}, values)
	assert.Equal(t, []int{0, 1, 2}, order) // Sequential order
}

// Test SequenceArraySeq
func TestSequenceArraySeq(t *testing.T) {
	var order []int
	src := []IO[int]{
		func() int { order = append(order, 0); return 1 },
		func() int { order = append(order, 1); return 2 },
		func() int { order = append(order, 2); return 3 },
	}

	result := SequenceArraySeq(src)
	values := result()

	assert.Equal(t, []int{1, 2, 3}, values)
	assert.Equal(t, []int{0, 1, 2}, order) // Sequential order
}

// Test MonadTraverseArraySeq
func TestMonadTraverseArraySeq(t *testing.T) {
	var order []int
	src := []int{1, 2, 3}

	result := MonadTraverseArraySeq(src, func(x int) IO[int] {
		return func() int {
			order = append(order, x)
			return x * 2
		}
	})

	values := result()
	assert.Equal(t, []int{2, 4, 6}, values)
	assert.Equal(t, []int{1, 2, 3}, order) // Sequential order
}

// Test TraverseRecordSeq
func TestTraverseRecordSeq(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := F.Pipe1(
		src,
		TraverseRecordSeq[string](func(x int) IO[int] {
			return Of(x * 2)
		}),
	)

	values := result()
	assert.Equal(t, 2, values["a"])
	assert.Equal(t, 4, values["b"])
}

// Test TraverseRecordWithIndeSeq (note the typo in the function name)
func TestTraverseRecordWithIndeSeq(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := F.Pipe1(
		src,
		TraverseRecordWithIndeSeq(func(k string, x int) IO[string] {
			return Of(k)
		}),
	)

	values := result()
	assert.Equal(t, "a", values["a"])
	assert.Equal(t, "b", values["b"])
}

// Test SequenceRecordSeq
func TestSequenceRecordSeq(t *testing.T) {
	src := map[string]IO[int]{
		"a": Of(1),
		"b": Of(2),
	}

	result := SequenceRecordSeq(src)
	values := result()

	assert.Equal(t, 1, values["a"])
	assert.Equal(t, 2, values["b"])
}

// Test MonadTraverseRecordSeq
func TestMonadTraverseRecordSeq(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}

	result := MonadTraverseRecordSeq(src, func(x int) IO[int] {
		return Of(x * 2)
	})

	values := result()
	assert.Equal(t, 2, values["a"])
	assert.Equal(t, 4, values["b"])
}

// Test SequenceTuple2
func TestSequenceTuple2(t *testing.T) {
	io1 := Of(10)
	io2 := Of(32)

	tup := T.MakeTuple2(io1, io2)
	result := SequenceTuple2(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 32, tuple.F2)
}

// Test SequenceT2
func TestSequenceT2(t *testing.T) {
	result := SequenceT2(Of(10), Of(32))
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 32, tuple.F2)
}

// Test SequenceTuple3
func TestSequenceTuple3(t *testing.T) {
	io1 := Of(10)
	io2 := Of(20)
	io3 := Of(30)

	tup := T.MakeTuple3(io1, io2, io3)
	result := SequenceTuple3(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
}

// Test SequenceSeqTuple2
func TestSequenceSeqTuple2(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 10 }
	io2 := func() int { order = append(order, 2); return 20 }

	tup := T.MakeTuple2(io1, io2)
	result := SequenceSeqTuple2(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, []int{1, 2}, order) // Sequential execution
}

// Test SequenceParTuple2
func TestSequenceParTuple2(t *testing.T) {
	io1 := Of(10)
	io2 := Of(20)

	tup := T.MakeTuple2(io1, io2)
	result := SequenceParTuple2(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
}

// Test SequenceT3
func TestSequenceT3(t *testing.T) {
	result := SequenceT3(Of(10), Of(20), Of(30))
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
}

// Test SequenceSeqT2
func TestSequenceSeqT2(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 10 }
	io2 := func() int { order = append(order, 2); return 20 }

	result := SequenceSeqT2(io1, io2)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, []int{1, 2}, order)
}

// Test SequenceParT2
func TestSequenceParT2(t *testing.T) {
	result := SequenceParT2(Of(10), Of(20))
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
}

// Test SequenceT4
func TestSequenceT4(t *testing.T) {
	result := SequenceT4(Of(1), Of(2), Of(3), Of(4))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
}

// Test SequenceTuple4
func TestSequenceTuple4(t *testing.T) {
	tup := T.MakeTuple4(Of(1), Of(2), Of(3), Of(4))
	result := SequenceTuple4(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
}

// Test SequenceT5
func TestSequenceT5(t *testing.T) {
	result := SequenceT5(Of(1), Of(2), Of(3), Of(4), Of(5))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
}

// Test TraverseTuple2
func TestTraverseTuple2(t *testing.T) {
	inputTuple := T.MakeTuple2(10, 20)
	result := TraverseTuple2(
		func(x int) IO[int] { return Of(x * 2) },
		func(x int) IO[int] { return Of(x + 1) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 20, tuple.F1)
	assert.Equal(t, 21, tuple.F2)
}

// Test TraverseTuple3
func TestTraverseTuple3(t *testing.T) {
	inputTuple := T.MakeTuple3(10, 20, 30)
	result := TraverseTuple3(
		func(x int) IO[int] { return Of(x * 2) },
		func(x int) IO[int] { return Of(x + 1) },
		func(x int) IO[int] { return Of(x - 1) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 20, tuple.F1)
	assert.Equal(t, 21, tuple.F2)
	assert.Equal(t, 29, tuple.F3)
}

// Test SequenceTuple1
func TestSequenceTuple1(t *testing.T) {
	tup := T.MakeTuple1(Of(42))
	result := SequenceTuple1(tup)
	tuple := result()

	assert.Equal(t, 42, tuple.F1)
}

// Test SequenceT1
func TestSequenceT1(t *testing.T) {
	result := SequenceT1(Of(42))
	tuple := result()

	assert.Equal(t, 42, tuple.F1)
}

// Test SequenceSeqT3
func TestSequenceSeqT3(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 10 }
	io2 := func() int { order = append(order, 2); return 20 }
	io3 := func() int { order = append(order, 3); return 30 }

	result := SequenceSeqT3(io1, io2, io3)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
	assert.Equal(t, []int{1, 2, 3}, order)
}

// Test SequenceParT3
func TestSequenceParT3(t *testing.T) {
	result := SequenceParT3(Of(10), Of(20), Of(30))
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
}

// Test TraverseSeqTuple3
func TestTraverseSeqTuple3(t *testing.T) {
	var order []int
	inputTuple := T.MakeTuple3(10, 20, 30)
	result := TraverseSeqTuple3(
		func(x int) IO[int] { return func() int { order = append(order, 1); return x * 2 } },
		func(x int) IO[int] { return func() int { order = append(order, 2); return x + 1 } },
		func(x int) IO[int] { return func() int { order = append(order, 3); return x - 1 } },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 20, tuple.F1)
	assert.Equal(t, 21, tuple.F2)
	assert.Equal(t, 29, tuple.F3)
	assert.Equal(t, []int{1, 2, 3}, order)
}

// Test TraverseParTuple3
func TestTraverseParTuple3(t *testing.T) {
	inputTuple := T.MakeTuple3(10, 20, 30)
	result := TraverseParTuple3(
		func(x int) IO[int] { return Of(x * 2) },
		func(x int) IO[int] { return Of(x + 1) },
		func(x int) IO[int] { return Of(x - 1) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 20, tuple.F1)
	assert.Equal(t, 21, tuple.F2)
	assert.Equal(t, 29, tuple.F3)
}

// Test SequenceSeqTuple3
func TestSequenceSeqTuple3(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 10 }
	io2 := func() int { order = append(order, 2); return 20 }
	io3 := func() int { order = append(order, 3); return 30 }

	tup := T.MakeTuple3(io1, io2, io3)
	result := SequenceSeqTuple3(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
	assert.Equal(t, []int{1, 2, 3}, order)
}

// Test SequenceParTuple3
func TestSequenceParTuple3(t *testing.T) {
	tup := T.MakeTuple3(Of(10), Of(20), Of(30))
	result := SequenceParTuple3(tup)
	tuple := result()

	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 20, tuple.F2)
	assert.Equal(t, 30, tuple.F3)
}

// Test SequenceSeqT4
func TestSequenceSeqT4(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }

	result := SequenceSeqT4(io1, io2, io3, io4)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, []int{1, 2, 3, 4}, order)
}

// Test SequenceParT4
func TestSequenceParT4(t *testing.T) {
	result := SequenceParT4(Of(1), Of(2), Of(3), Of(4))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
}

// Test SequenceSeqTuple4
func TestSequenceSeqTuple4(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }

	tup := T.MakeTuple4(io1, io2, io3, io4)
	result := SequenceSeqTuple4(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, []int{1, 2, 3, 4}, order)
}

// Test SequenceParTuple4
func TestSequenceParTuple4(t *testing.T) {
	tup := T.MakeTuple4(Of(1), Of(2), Of(3), Of(4))
	result := SequenceParTuple4(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
}

// Test TraverseTuple4
func TestTraverseTuple4(t *testing.T) {
	inputTuple := T.MakeTuple4(1, 2, 3, 4)
	result := TraverseTuple4(
		func(x int) IO[int] { return Of(x * 10) },
		func(x int) IO[int] { return Of(x * 20) },
		func(x int) IO[int] { return Of(x * 30) },
		func(x int) IO[int] { return Of(x * 40) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 40, tuple.F2)
	assert.Equal(t, 90, tuple.F3)
	assert.Equal(t, 160, tuple.F4)
}

// Test SequenceSeqT5
func TestSequenceSeqT5(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }

	result := SequenceSeqT5(io1, io2, io3, io4, io5)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, order)
}

// Test SequenceParT5
func TestSequenceParT5(t *testing.T) {
	result := SequenceParT5(Of(1), Of(2), Of(3), Of(4), Of(5))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
}

// Test SequenceSeqTuple1
func TestSequenceSeqTuple1(t *testing.T) {
	var executed bool
	io1 := func() int { executed = true; return 42 }

	tup := T.MakeTuple1(io1)
	result := SequenceSeqTuple1(tup)
	tuple := result()

	assert.Equal(t, 42, tuple.F1)
	assert.True(t, executed)
}

// Test SequenceParTuple1
func TestSequenceParTuple1(t *testing.T) {
	tup := T.MakeTuple1(Of(42))
	result := SequenceParTuple1(tup)
	tuple := result()

	assert.Equal(t, 42, tuple.F1)
}

// Test TraverseTuple1
func TestTraverseTuple1(t *testing.T) {
	inputTuple := T.MakeTuple1(21)
	result := TraverseTuple1(
		func(x int) IO[int] { return Of(x * 2) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 42, tuple.F1)
}

// Test TraverseSeqTuple1
func TestTraverseSeqTuple1(t *testing.T) {
	var executed bool
	inputTuple := T.MakeTuple1(21)
	result := TraverseSeqTuple1(
		func(x int) IO[int] { return func() int { executed = true; return x * 2 } },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 42, tuple.F1)
	assert.True(t, executed)
}

// Test TraverseParTuple1
func TestTraverseParTuple1(t *testing.T) {
	inputTuple := T.MakeTuple1(21)
	result := TraverseParTuple1(
		func(x int) IO[int] { return Of(x * 2) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 42, tuple.F1)
}

// Test SequenceTuple5
func TestSequenceTuple5(t *testing.T) {
	tup := T.MakeTuple5(Of(1), Of(2), Of(3), Of(4), Of(5))
	result := SequenceTuple5(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
}

// Test SequenceSeqTuple5
func TestSequenceSeqTuple5(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }

	tup := T.MakeTuple5(io1, io2, io3, io4, io5)
	result := SequenceSeqTuple5(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, order)
}

// Test SequenceParTuple5
func TestSequenceParTuple5(t *testing.T) {
	tup := T.MakeTuple5(Of(1), Of(2), Of(3), Of(4), Of(5))
	result := SequenceParTuple5(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
}

// Test TraverseTuple5
func TestTraverseTuple5(t *testing.T) {
	inputTuple := T.MakeTuple5(1, 2, 3, 4, 5)
	result := TraverseTuple5(
		func(x int) IO[int] { return Of(x * 10) },
		func(x int) IO[int] { return Of(x * 20) },
		func(x int) IO[int] { return Of(x * 30) },
		func(x int) IO[int] { return Of(x * 40) },
		func(x int) IO[int] { return Of(x * 50) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 40, tuple.F2)
	assert.Equal(t, 90, tuple.F3)
	assert.Equal(t, 160, tuple.F4)
	assert.Equal(t, 250, tuple.F5)
}

// Test TraverseSeqTuple5
func TestTraverseSeqTuple5(t *testing.T) {
	var order []int
	inputTuple := T.MakeTuple5(1, 2, 3, 4, 5)
	result := TraverseSeqTuple5(
		func(x int) IO[int] { return func() int { order = append(order, 1); return x * 10 } },
		func(x int) IO[int] { return func() int { order = append(order, 2); return x * 20 } },
		func(x int) IO[int] { return func() int { order = append(order, 3); return x * 30 } },
		func(x int) IO[int] { return func() int { order = append(order, 4); return x * 40 } },
		func(x int) IO[int] { return func() int { order = append(order, 5); return x * 50 } },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 40, tuple.F2)
	assert.Equal(t, 90, tuple.F3)
	assert.Equal(t, 160, tuple.F4)
	assert.Equal(t, 250, tuple.F5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, order)
}

// Test TraverseParTuple5
func TestTraverseParTuple5(t *testing.T) {
	inputTuple := T.MakeTuple5(1, 2, 3, 4, 5)
	result := TraverseParTuple5(
		func(x int) IO[int] { return Of(x * 10) },
		func(x int) IO[int] { return Of(x * 20) },
		func(x int) IO[int] { return Of(x * 30) },
		func(x int) IO[int] { return Of(x * 40) },
		func(x int) IO[int] { return Of(x * 50) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 10, tuple.F1)
	assert.Equal(t, 40, tuple.F2)
	assert.Equal(t, 90, tuple.F3)
	assert.Equal(t, 160, tuple.F4)
	assert.Equal(t, 250, tuple.F5)
}

// Test SequenceT6
func TestSequenceT6(t *testing.T) {
	result := SequenceT6(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
}

// Test SequenceSeqT6
func TestSequenceSeqT6(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }

	result := SequenceSeqT6(io1, io2, io3, io4, io5, io6)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, order)
}

// Test SequenceParT6
func TestSequenceParT6(t *testing.T) {
	result := SequenceParT6(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
}

// Test SequenceT7
func TestSequenceT7(t *testing.T) {
	result := SequenceT7(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
}

// Test SequenceT8
func TestSequenceT8(t *testing.T) {
	result := SequenceT8(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
}

// Test SequenceT9
func TestSequenceT9(t *testing.T) {
	result := SequenceT9(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
}

// Test SequenceT10
func TestSequenceT10(t *testing.T) {
	result := SequenceT10(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9), Of(10))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
	assert.Equal(t, 10, tuple.F10)
}

// Test SequenceTuple6
func TestSequenceTuple6(t *testing.T) {
	tup := T.MakeTuple6(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6))
	result := SequenceTuple6(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
}

// Test SequenceSeqTuple6
func TestSequenceSeqTuple6(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }

	tup := T.MakeTuple6(io1, io2, io3, io4, io5, io6)
	result := SequenceSeqTuple6(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, order)
}

// Test SequenceParTuple6
func TestSequenceParTuple6(t *testing.T) {
	tup := T.MakeTuple6(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6))
	result := SequenceParTuple6(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
}

// Test TraverseTuple6
func TestTraverseTuple6(t *testing.T) {
	inputTuple := T.MakeTuple6(1, 2, 3, 4, 5, 6)
	result := TraverseTuple6(
		func(x int) IO[int] { return Of(x * 1) },
		func(x int) IO[int] { return Of(x * 2) },
		func(x int) IO[int] { return Of(x * 3) },
		func(x int) IO[int] { return Of(x * 4) },
		func(x int) IO[int] { return Of(x * 5) },
		func(x int) IO[int] { return Of(x * 6) },
	)(inputTuple)

	tuple := result()
	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 4, tuple.F2)
	assert.Equal(t, 9, tuple.F3)
	assert.Equal(t, 16, tuple.F4)
	assert.Equal(t, 25, tuple.F5)
	assert.Equal(t, 36, tuple.F6)
}

// Test SequenceTuple7
func TestSequenceTuple7(t *testing.T) {
	tup := T.MakeTuple7(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7))
	result := SequenceTuple7(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
}

// Test SequenceSeqT7
func TestSequenceSeqT7(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }
	io7 := func() int { order = append(order, 7); return 7 }

	result := SequenceSeqT7(io1, io2, io3, io4, io5, io6, io7)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, order)
}

// Test SequenceParT7
func TestSequenceParT7(t *testing.T) {
	result := SequenceParT7(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
}

// Test SequenceTuple8
func TestSequenceTuple8(t *testing.T) {
	tup := T.MakeTuple8(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8))
	result := SequenceTuple8(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
}

// Test SequenceSeqT8
func TestSequenceSeqT8(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }
	io7 := func() int { order = append(order, 7); return 7 }
	io8 := func() int { order = append(order, 8); return 8 }

	result := SequenceSeqT8(io1, io2, io3, io4, io5, io6, io7, io8)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, order)
}

// Test SequenceParT8
func TestSequenceParT8(t *testing.T) {
	result := SequenceParT8(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
}

// Test SequenceTuple9
func TestSequenceTuple9(t *testing.T) {
	tup := T.MakeTuple9(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9))
	result := SequenceTuple9(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
}

// Test SequenceSeqT9
func TestSequenceSeqT9(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }
	io7 := func() int { order = append(order, 7); return 7 }
	io8 := func() int { order = append(order, 8); return 8 }
	io9 := func() int { order = append(order, 9); return 9 }

	result := SequenceSeqT9(io1, io2, io3, io4, io5, io6, io7, io8, io9)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, order)
}

// Test SequenceParT9
func TestSequenceParT9(t *testing.T) {
	result := SequenceParT9(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
}

// Test SequenceTuple10
func TestSequenceTuple10(t *testing.T) {
	tup := T.MakeTuple10(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9), Of(10))
	result := SequenceTuple10(tup)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
	assert.Equal(t, 10, tuple.F10)
}

// Test SequenceSeqT10
func TestSequenceSeqT10(t *testing.T) {
	var order []int
	io1 := func() int { order = append(order, 1); return 1 }
	io2 := func() int { order = append(order, 2); return 2 }
	io3 := func() int { order = append(order, 3); return 3 }
	io4 := func() int { order = append(order, 4); return 4 }
	io5 := func() int { order = append(order, 5); return 5 }
	io6 := func() int { order = append(order, 6); return 6 }
	io7 := func() int { order = append(order, 7); return 7 }
	io8 := func() int { order = append(order, 8); return 8 }
	io9 := func() int { order = append(order, 9); return 9 }
	io10 := func() int { order = append(order, 10); return 10 }

	result := SequenceSeqT10(io1, io2, io3, io4, io5, io6, io7, io8, io9, io10)
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
	assert.Equal(t, 10, tuple.F10)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, order)
}

// Test SequenceParT10
func TestSequenceParT10(t *testing.T) {
	result := SequenceParT10(Of(1), Of(2), Of(3), Of(4), Of(5), Of(6), Of(7), Of(8), Of(9), Of(10))
	tuple := result()

	assert.Equal(t, 1, tuple.F1)
	assert.Equal(t, 2, tuple.F2)
	assert.Equal(t, 3, tuple.F3)
	assert.Equal(t, 4, tuple.F4)
	assert.Equal(t, 5, tuple.F5)
	assert.Equal(t, 6, tuple.F6)
	assert.Equal(t, 7, tuple.F7)
	assert.Equal(t, 8, tuple.F8)
	assert.Equal(t, 9, tuple.F9)
	assert.Equal(t, 10, tuple.F10)
}
