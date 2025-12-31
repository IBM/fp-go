// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

// Test FromNillable
func TestFromNillable(t *testing.T) {
	var nilPtr *int = nil
	assert.Equal(t, None[*int](), FromNillable(nilPtr))

	val := 42
	ptr := &val
	result := FromNillable(ptr)
	assert.True(t, IsSome(result))
	unwrapped, ok := Unwrap(result)
	assert.True(t, ok)
	assert.Equal(t, &val, unwrapped)
}

// Test FromValidation
func TestFromValidation(t *testing.T) {
	validate := func(x int) (int, bool) {
		if x > 0 {
			return x * 2, true
		}
		return 0, false
	}

	f := FromValidation(validate)
	assert.Equal(t, Some(10), f(5))
	assert.Equal(t, None[int](), f(-1))
}

// Test MonadAp
func TestMonadAp(t *testing.T) {
	double := N.Mul(2)

	assert.Equal(t, Some(10), MonadAp(Some(double), Some(5)))
	assert.Equal(t, None[int](), MonadAp(Some(double), None[int]()))
	assert.Equal(t, None[int](), MonadAp(None[func(int) int](), Some(5)))
	assert.Equal(t, None[int](), MonadAp(None[func(int) int](), None[int]()))
}

// Test MonadMap
func TestMonadMap(t *testing.T) {
	double := N.Mul(2)

	assert.Equal(t, Some(10), MonadMap(Some(5), double))
	assert.Equal(t, None[int](), MonadMap(None[int](), double))
}

// Test MonadMapTo
func TestMonadMapTo(t *testing.T) {
	assert.Equal(t, Some("hello"), MonadMapTo(Some(42), "hello"))
	assert.Equal(t, None[string](), MonadMapTo(None[int](), "hello"))
}

// Test MapTo
func TestMapTo(t *testing.T) {
	replaceWith42 := MapTo[string](42)
	assert.Equal(t, Some(42), replaceWith42(Some("hello")))
	assert.Equal(t, None[int](), replaceWith42(None[string]()))
}

// Test MonadGetOrElse
func TestMonadGetOrElse(t *testing.T) {
	defaultVal := func() int { return 0 }

	assert.Equal(t, 42, MonadGetOrElse(Some(42), defaultVal))
	assert.Equal(t, 0, MonadGetOrElse(None[int](), defaultVal))
}

// Test GetOrElse
func TestGetOrElse(t *testing.T) {
	getOrZero := GetOrElse(func() int { return 0 })

	assert.Equal(t, 42, getOrZero(Some(42)))
	assert.Equal(t, 0, getOrZero(None[int]()))
}

// Test MonadChain
func TestMonadChain(t *testing.T) {
	validate := func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	assert.Equal(t, Some(10), MonadChain(Some(5), validate))
	assert.Equal(t, None[int](), MonadChain(Some(-1), validate))
	assert.Equal(t, None[int](), MonadChain(None[int](), validate))
}

// Test MonadChainTo
func TestMonadChainTo(t *testing.T) {
	assert.Equal(t, Some("hello"), MonadChainTo(Some(42), Some("hello")))
	assert.Equal(t, None[string](), MonadChainTo(Some(42), None[string]()))
	assert.Equal(t, None[string](), MonadChainTo(None[int](), Some("hello")))
}

// Test ChainTo
func TestChainTo(t *testing.T) {
	replaceWith := ChainTo[int](Some("hello"))
	assert.Equal(t, Some("hello"), replaceWith(Some(42)))
	assert.Equal(t, None[string](), replaceWith(None[int]()))
}

// Test MonadChainFirst
func TestMonadChainFirst(t *testing.T) {
	sideEffect := func(x int) Option[string] {
		return Some(fmt.Sprintf("%d", x))
	}

	assert.Equal(t, Some(5), MonadChainFirst(Some(5), sideEffect))
	assert.Equal(t, None[int](), MonadChainFirst(None[int](), sideEffect))
}

// Test ChainFirst
func TestChainFirst(t *testing.T) {
	sideEffect := func(x int) Option[string] {
		return Some(fmt.Sprintf("%d", x))
	}
	chainFirst := ChainFirst(sideEffect)

	assert.Equal(t, Some(5), chainFirst(Some(5)))
	assert.Equal(t, None[int](), chainFirst(None[int]()))
}

// Test MonadAlt
func TestMonadAlt(t *testing.T) {
	alternative := func() Option[int] { return Some(10) }

	assert.Equal(t, Some(5), MonadAlt(Some(5), alternative))
	assert.Equal(t, Some(10), MonadAlt(None[int](), alternative))
}

// Test MonadSequence2
func TestMonadSequence2(t *testing.T) {
	combine := func(a, b int) Option[int] {
		return Some(a + b)
	}

	assert.Equal(t, Some(5), MonadSequence2(Some(2), Some(3), combine))
	assert.Equal(t, None[int](), MonadSequence2(None[int](), Some(3), combine))
	assert.Equal(t, None[int](), MonadSequence2(Some(2), None[int](), combine))
}

// Test Sequence2
func TestSequence2(t *testing.T) {
	add := Sequence2(func(a, b int) Option[int] { return Some(a + b) })

	assert.Equal(t, Some(5), add(Some(2), Some(3)))
	assert.Equal(t, None[int](), add(None[int](), Some(3)))
}

// Test Filter
func TestFilter(t *testing.T) {
	isPositive := Filter(N.MoreThan(0))

	assert.Equal(t, Some(5), isPositive(Some(5)))
	assert.Equal(t, None[int](), isPositive(Some(-1)))
	assert.Equal(t, None[int](), isPositive(None[int]()))
}

// Test MonadFlap
func TestMonadFlap(t *testing.T) {
	double := N.Mul(2)

	assert.Equal(t, Some(10), MonadFlap(Some(double), 5))
	assert.Equal(t, None[int](), MonadFlap(None[func(int) int](), 5))
}

// Test Flap
func TestFlap(t *testing.T) {
	applyFive := Flap[int](5)
	double := N.Mul(2)

	assert.Equal(t, Some(10), applyFive(Some(double)))
	assert.Equal(t, None[int](), applyFive(None[func(int) int]()))
}

// Test Unwrap
func TestUnwrap(t *testing.T) {
	val, ok := Unwrap(Some(42))
	assert.True(t, ok)
	assert.Equal(t, 42, val)

	val, ok = Unwrap(None[int]())
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

// Test String and Format
func TestStringFormat(t *testing.T) {
	opt := Some(42)
	str := opt.String()
	assert.Contains(t, str, "Some")
	assert.Contains(t, str, "42")

	none := None[int]()
	str = none.String()
	assert.Contains(t, str, "None")
}

// Test Semigroup
func TestSemigroup(t *testing.T) {
	intSemigroup := N.MonoidSum[int]()
	optSemigroup := Semigroup[int]()(intSemigroup)

	assert.Equal(t, Some(5), optSemigroup.Concat(Some(2), Some(3)))
	assert.Equal(t, Some(2), optSemigroup.Concat(Some(2), None[int]()))
	assert.Equal(t, Some(3), optSemigroup.Concat(None[int](), Some(3)))
	assert.Equal(t, None[int](), optSemigroup.Concat(None[int](), None[int]()))
}

// Test Monoid
func TestMonoid(t *testing.T) {
	intSemigroup := N.MonoidSum[int]()
	optMonoid := Monoid[int]()(intSemigroup)

	assert.Equal(t, Some(5), optMonoid.Concat(Some(2), Some(3)))
	assert.Equal(t, None[int](), optMonoid.Empty())
}

// Test ApplySemigroup
func TestApplySemigroup(t *testing.T) {
	intSemigroup := N.MonoidSum[int]()
	optSemigroup := ApplySemigroup(intSemigroup)

	assert.Equal(t, Some(5), optSemigroup.Concat(Some(2), Some(3)))
	assert.Equal(t, None[int](), optSemigroup.Concat(Some(2), None[int]()))
}

// Test ApplicativeMonoid
func TestApplicativeMonoid(t *testing.T) {
	intMonoid := N.MonoidSum[int]()
	optMonoid := ApplicativeMonoid(intMonoid)

	assert.Equal(t, Some(5), optMonoid.Concat(Some(2), Some(3)))
	assert.Equal(t, Some(0), optMonoid.Empty())
}

// Test AlternativeMonoid
func TestAlternativeMonoid(t *testing.T) {
	intMonoid := N.MonoidSum[int]()
	optMonoid := AlternativeMonoid(intMonoid)

	// AlternativeMonoid uses applicative semantics, so it combines values
	assert.Equal(t, Some(5), optMonoid.Concat(Some(2), Some(3)))
	assert.Equal(t, Some(3), optMonoid.Concat(None[int](), Some(3)))
	assert.Equal(t, Some(0), optMonoid.Empty())
}

// Test AltMonoid
func TestAltMonoid(t *testing.T) {
	optMonoid := AltMonoid[int]()

	assert.Equal(t, Some(2), optMonoid.Concat(Some(2), Some(3)))
	assert.Equal(t, Some(3), optMonoid.Concat(None[int](), Some(3)))
	assert.Equal(t, None[int](), optMonoid.Empty())
}

// Test Do, Let, LetTo, BindTo
func TestDoLetLetToBindTo(t *testing.T) {
	type State struct {
		x        int
		y        int
		computed int
		name     string
	}

	result := F.Pipe4(
		Do(State{}),
		Let(func(c int) func(State) State {
			return func(s State) State { s.x = c; return s }
		}, func(s State) int { return 5 }),
		LetTo(func(n string) func(State) State {
			return func(s State) State { s.name = n; return s }
		}, "test"),
		Bind(func(y int) func(State) State {
			return func(s State) State { s.y = y; return s }
		}, func(s State) Option[int] { return Some(10) }),
		Map(func(s State) State {
			s.computed = s.x + s.y
			return s
		}),
	)

	expected := Some(State{x: 5, y: 10, computed: 15, name: "test"})
	assert.Equal(t, expected, result)
}

// Test BindTo
func TestBindToFunction(t *testing.T) {
	type State struct {
		value int
	}

	result := F.Pipe1(
		Some(42),
		BindTo(func(x int) State { return State{value: x} }),
	)

	assert.Equal(t, Some(State{value: 42}), result)
}

// Test Functor
func TestFunctor(t *testing.T) {
	f := Functor[int, string]()
	mapper := f.Map(strconv.Itoa)

	assert.Equal(t, Some("42"), mapper(Some(42)))
	assert.Equal(t, None[string](), mapper(None[int]()))
}

// Test Monad
func TestMonad(t *testing.T) {
	m := Monad[int, string]()

	// Test Of
	assert.Equal(t, Some(42), m.Of(42))

	// Test Map
	mapper := m.Map(strconv.Itoa)
	assert.Equal(t, Some("42"), mapper(Some(42)))

	// Test Chain
	chainer := m.Chain(func(x int) Option[string] {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	})
	assert.Equal(t, Some("42"), chainer(Some(42)))

	// Test Ap
	double := func(x int) string { return fmt.Sprintf("%d", x*2) }
	ap := m.Ap(Some(5))
	assert.Equal(t, Some("10"), ap(Some(double)))
}

// Test Pointed
func TestPointed(t *testing.T) {
	p := Pointed[int]()
	assert.Equal(t, Some(42), p.Of(42))
}

// Test ToAny
func TestToAny(t *testing.T) {
	result := ToAny(42)
	assert.True(t, IsSome(result))

	val, ok := Unwrap(result)
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

// Test TraverseArray
func TestTraverseArray(t *testing.T) {
	validate := func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	result := TraverseArray(validate)([]int{1, 2, 3})
	assert.Equal(t, Some([]int{2, 4, 6}), result)

	result = TraverseArray(validate)([]int{1, -1, 3})
	assert.Equal(t, None[[]int](), result)
}

// Test TraverseArrayWithIndex
func TestTraverseArrayWithIndex(t *testing.T) {
	f := func(i int, x int) Option[int] {
		if x > i {
			return Some(x + i)
		}
		return None[int]()
	}

	result := TraverseArrayWithIndex(f)([]int{1, 2, 3})
	assert.Equal(t, Some([]int{1, 3, 5}), result)
}

// Test TraverseRecord
func TestTraverseRecord(t *testing.T) {
	validate := func(x int) Option[string] {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	input := map[string]int{"a": 1, "b": 2}
	result := TraverseRecord[string](validate)(input)

	expected := Some(map[string]string{"a": "1", "b": "2"})
	assert.Equal(t, expected, result)
}

// Test TraverseRecordWithIndex
func TestTraverseRecordWithIndex(t *testing.T) {
	f := func(k string, v int) Option[string] {
		return Some(fmt.Sprintf("%s:%d", k, v))
	}

	input := map[string]int{"a": 1, "b": 2}
	result := TraverseRecordWithIndex(f)(input)

	assert.True(t, IsSome(result))
}

// Test SequencePair
func TestSequencePair(t *testing.T) {
	pair := P.MakePair(Some(1), Some("hello"))
	result := SequencePair(pair)

	assert.True(t, IsSome(result))

	pair2 := P.MakePair(Some(1), None[string]())
	result2 := SequencePair(pair2)
	assert.True(t, IsNone(result2))
}

// Test Optionize functions
func TestOptionize0(t *testing.T) {
	f := func() (int, bool) {
		return 42, true
	}

	optF := Optionize0(f)
	assert.Equal(t, Some(42), optF())
}

func TestOptionize1(t *testing.T) {
	f := func(x int) (int, bool) {
		if x > 0 {
			return x * 2, true
		}
		return 0, false
	}

	optF := Optionize1(f)
	assert.Equal(t, Some(10), optF(5))
	assert.Equal(t, None[int](), optF(-1))
}

func TestOptionize2(t *testing.T) {
	f := func(x, y int) (int, bool) {
		if x > 0 && y > 0 {
			return x + y, true
		}
		return 0, false
	}

	optF := Optionize2(f)
	assert.Equal(t, Some(5), optF(2, 3))
	assert.Equal(t, None[int](), optF(-1, 3))
}

// Test Unoptionize functions
func TestUnoptionize0(t *testing.T) {
	f := func() Option[int] {
		return Some(42)
	}

	unoptF := Unoptionize0(f)
	val, ok := unoptF()
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestUnoptionize1(t *testing.T) {
	f := func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	unoptF := Unoptionize1(f)
	val, ok := unoptF(5)
	assert.True(t, ok)
	assert.Equal(t, 10, val)

	_, ok = unoptF(-1)
	assert.False(t, ok)
}

// Test SequenceTuple functions
func TestSequenceTuple2(t *testing.T) {
	tuple := T.MakeTuple2(Some(1), Some("hello"))
	result := SequenceTuple2(tuple)

	expected := Some(T.MakeTuple2(1, "hello"))
	assert.Equal(t, expected, result)
}

func TestSequenceTuple3(t *testing.T) {
	tuple := T.MakeTuple3(Some(1), Some("hello"), Some(true))
	result := SequenceTuple3(tuple)

	expected := Some(T.MakeTuple3(1, "hello", true))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple functions
func TestTraverseTuple2(t *testing.T) {
	f1 := func(x int) Option[int] { return Some(x * 2) }
	f2 := func(s string) Option[string] { return Some(s + "!") }

	traverse := TraverseTuple2(f1, f2)
	tuple := T.MakeTuple2(5, "hello")
	result := traverse(tuple)

	expected := Some(T.MakeTuple2(10, "hello!"))
	assert.Equal(t, expected, result)
}

// Test FromStrictCompare
func TestFromStrictCompare(t *testing.T) {
	optOrd := FromStrictCompare[int]()

	assert.Equal(t, 0, optOrd.Compare(Some(5), Some(5)))
	assert.Equal(t, -1, optOrd.Compare(Some(3), Some(5)))
	assert.Equal(t, 1, optOrd.Compare(Some(5), Some(3)))
	assert.Equal(t, -1, optOrd.Compare(None[int](), Some(5)))
	assert.Equal(t, 1, optOrd.Compare(Some(5), None[int]()))
}
