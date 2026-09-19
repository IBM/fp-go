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

package pair_test

import (
	"strconv"
	"strings"
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	SG "github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

var (
	intSum    = N.SemigroupSum[int]()
	strConcat = SG.MakeSemigroup(S.Monoid.Concat)
)

func TestOf(t *testing.T) {
	p := pair.Of("x")
	assert.Equal(t, "x", pair.Head(p))
	assert.Equal(t, "x", pair.Tail(p))
}

func TestTupleRoundTrip(t *testing.T) {
	tp := tuple.MakeTuple2("hello", 42)
	p := pair.FromTuple(tp)

	assert.Equal(t, "hello", pair.Head(p))
	assert.Equal(t, 42, pair.Tail(p))
	assert.Equal(t, tp, pair.ToTuple(p))
	assert.Equal(t, p, pair.FromTuple(pair.ToTuple(p)))
}

func TestFromHeadAndFromTail(t *testing.T) {
	assert.Equal(t, pair.MakePair("k", 1), pair.FromHead[int]("k")(1))
	assert.Equal(t, pair.MakePair("k", 1), pair.FromTail[string](1)("k"))
}

func TestMapVariants(t *testing.T) {
	p := pair.MakePair(5, "hello")

	assert.Equal(t, pair.MakePair(5, 5), pair.MonadMap(p, S.Size))
	assert.Equal(t, pair.MakePair(5, 5), pair.MonadMapTail(p, S.Size))
	assert.Equal(t, pair.MakePair(5, 5), pair.Map[int](S.Size)(p))
	assert.Equal(t, pair.MakePair(5, 5), pair.MapTail[int](S.Size)(p))
	assert.Equal(t, pair.MakePair("5", "hello"), pair.MonadMapHead(p, strconv.Itoa))
	assert.Equal(t, pair.MakePair("5", "hello"), pair.MapHead[string](strconv.Itoa)(p))
	assert.Equal(t, pair.MakePair("5", 5), pair.MonadBiMap(p, strconv.Itoa, S.Size))
	assert.Equal(t, pair.MakePair("5", 5), pair.BiMap(strconv.Itoa, S.Size)(p))
}

func TestChainTail(t *testing.T) {
	p := pair.MakePair(5, "hello")
	f := func(s string) pair.Pair[int, int] {
		return pair.MakePair(S.Size(s), S.Size(s)*2)
	}
	expected := pair.MakePair(10, 10)

	assert.Equal(t, expected, pair.MonadChainTail(intSum, p, f))
	assert.Equal(t, expected, pair.MonadChain(intSum, p, f))
	assert.Equal(t, expected, pair.ChainTail(intSum, f)(p))
	assert.Equal(t, expected, pair.Chain(intSum, f)(p))
}

func TestChainTailCombinesHeadsLeftToRight(t *testing.T) {
	p := pair.MakePair("a", 1)
	f := func(n int) pair.Pair[string, int] {
		return pair.MakePair("b", n+1)
	}
	// the head of the input is the left operand, the head of the result the right one
	assert.Equal(t, pair.MakePair("ab", 2), pair.ChainTail(strConcat, f)(p))
}

func TestChainHead(t *testing.T) {
	p := pair.MakePair(5, "hello")
	f := func(n int) pair.Pair[string, string] {
		return pair.MakePair(strconv.Itoa(n), "!")
	}
	expected := pair.MakePair("5", "hello!")

	assert.Equal(t, expected, pair.MonadChainHead(strConcat, p, f))
	assert.Equal(t, expected, pair.ChainHead(strConcat, f)(p))
}

func TestApTail(t *testing.T) {
	pf := pair.MakePair("f", S.Size)
	pv := pair.MakePair("v", "hello")
	// the head of the value is the left operand, the head of the function the right one
	expected := pair.MakePair("vf", 5)

	assert.Equal(t, expected, pair.MonadApTail(strConcat, pf, pv))
	assert.Equal(t, expected, pair.MonadAp(strConcat, pf, pv))
	assert.Equal(t, expected, pair.ApTail[string, string, int](strConcat, pv)(pf))
	assert.Equal(t, expected, pair.Ap[string, string, int](strConcat, pv)(pf))
}

func TestApHead(t *testing.T) {
	pf := pair.MakePair(strconv.Itoa, "!")
	pv := pair.MakePair(42, "hello")
	// the tail of the value is the left operand, the tail of the function the right one
	expected := pair.MakePair("42", "hello!")

	assert.Equal(t, expected, pair.MonadApHead(strConcat, pf, pv))
	assert.Equal(t, expected, pair.ApHead[string, int, string](strConcat, pv)(pf))
}

func TestSwapIsInvolution(t *testing.T) {
	p := pair.MakePair("hello", 42)
	assert.Equal(t, pair.MakePair(42, "hello"), pair.Swap(p))
	assert.Equal(t, p, pair.Swap(pair.Swap(p)))
}

func TestPairedAndUnpaired(t *testing.T) {
	sub := func(a, b int) int { return a - b }

	paired := pair.Paired(sub)
	assert.Equal(t, 7, paired(pair.MakePair(10, 3)))

	unpaired := pair.Unpaired(paired)
	assert.Equal(t, 7, unpaired(10, 3))
	assert.Equal(t, sub(3, 10), unpaired(3, 10))
}

func TestMerge(t *testing.T) {
	// data-last: the tail is the argument of the outer function, the head the data
	merge := pair.Merge(N.Sub[int])
	assert.Equal(t, 7, merge(pair.MakePair(10, 3)))

	format := pair.Merge(func(n int) func(string) string {
		return func(s string) string { return s + strconv.Itoa(n) }
	})
	assert.Equal(t, "a1", format(pair.MakePair("a", 1)))
}

func TestEq(t *testing.T) {
	pairEq := pair.Eq(EQ.FromStrictEquals[string](), EQ.FromStrictEquals[int]())

	assert.True(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 1)))
	assert.False(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("b", 1)))
	assert.False(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 2)))
	assert.False(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("b", 2)))
}

func TestEqUsesComponentEquality(t *testing.T) {
	caseInsensitive := EQ.FromEquals(strings.EqualFold)
	pairEq := pair.Eq(caseInsensitive, EQ.FromStrictEquals[int]())

	assert.True(t, pairEq.Equals(pair.MakePair("Hello", 1), pair.MakePair("hELLO", 1)))
	assert.False(t, pairEq.Equals(pair.MakePair("Hello", 1), pair.MakePair("hELLO", 2)))
}

func TestFromStrictEquals(t *testing.T) {
	pairEq := pair.FromStrictEquals[string, int]()

	assert.True(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 1)))
	assert.False(t, pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 2)))
}

func TestTailInstances(t *testing.T) {
	m := N.MonoidSum[int]()
	p := pair.MakePair(5, "hello")

	assert.Equal(t, pair.MakePair(0, "x"), pair.Pointed[string](m).Of("x"))
	assert.Equal(t, pair.MakePair(0, "x"), pair.PointedTail[string](m).Of("x"))

	assert.Equal(t, pair.MakePair(5, 5), pair.Functor[string, int, int]().Map(S.Size)(p))
	assert.Equal(t, pair.MakePair(5, 5), pair.FunctorTail[string, int, int]().Map(S.Size)(p))

	app := pair.Applicative[string, int, int](m)
	assert.Equal(t, pair.MakePair(0, "x"), app.Of("x"))
	assert.Equal(t, pair.MakePair(5, 5), app.Map(S.Size)(p))
	assert.Equal(t, pair.MakePair(15, 5), app.Ap(p)(pair.MakePair(10, S.Size)))

	appT := pair.ApplicativeTail[string, int, int](m)
	assert.Equal(t, pair.MakePair(0, "x"), appT.Of("x"))
	assert.Equal(t, pair.MakePair(5, 5), appT.Map(S.Size)(p))
	assert.Equal(t, pair.MakePair(15, 5), appT.Ap(p)(pair.MakePair(10, S.Size)))

	f := func(s string) pair.Pair[int, int] {
		return pair.MakePair(1, S.Size(s))
	}
	mon := pair.Monad[string, int, int](m)
	assert.Equal(t, pair.MakePair(0, "x"), mon.Of("x"))
	assert.Equal(t, pair.MakePair(5, 5), mon.Map(S.Size)(p))
	assert.Equal(t, pair.MakePair(6, 5), mon.Chain(f)(p))
	assert.Equal(t, pair.MakePair(15, 5), mon.Ap(p)(pair.MakePair(10, S.Size)))

	monT := pair.MonadTail[string, int, int](m)
	assert.Equal(t, pair.MakePair(0, "x"), monT.Of("x"))
	assert.Equal(t, pair.MakePair(5, 5), monT.Map(S.Size)(p))
	assert.Equal(t, pair.MakePair(6, 5), monT.Chain(f)(p))
	assert.Equal(t, pair.MakePair(15, 5), monT.Ap(p)(pair.MakePair(10, S.Size)))
}

func TestHeadInstances(t *testing.T) {
	m := S.Monoid
	p := pair.MakePair(42, "hello")

	assert.Equal(t, pair.MakePair(1, ""), pair.PointedHead[int](m).Of(1))
	assert.Equal(t, pair.MakePair("42", "hello"), pair.FunctorHead[int, string, string]().Map(strconv.Itoa)(p))

	app := pair.ApplicativeHead[int, string, string](m)
	assert.Equal(t, pair.MakePair(1, ""), app.Of(1))
	assert.Equal(t, pair.MakePair("42", "hello"), app.Map(strconv.Itoa)(p))
	assert.Equal(t, pair.MakePair("42", "hello!"), app.Ap(p)(pair.MakePair(strconv.Itoa, "!")))

	f := func(n int) pair.Pair[string, string] {
		return pair.MakePair(strconv.Itoa(n), "!")
	}
	mon := pair.MonadHead[int, string, string](m)
	assert.Equal(t, pair.MakePair(1, ""), mon.Of(1))
	assert.Equal(t, pair.MakePair("42", "hello"), mon.Map(strconv.Itoa)(p))
	assert.Equal(t, pair.MakePair("42", "hello!"), mon.Chain(f)(p))
	assert.Equal(t, pair.MakePair("42", "hello!"), mon.Ap(p)(pair.MakePair(strconv.Itoa, "!")))
}

func TestMonadSequence(t *testing.T) {
	mmap := O.MonadMap[int, pair.Pair[string, int]]

	assert.Equal(t, O.Some(pair.MakePair("k", 42)), pair.MonadSequence(mmap, pair.MakePair("k", O.Some(42))))
	assert.Equal(t, O.None[pair.Pair[string, int]](), pair.MonadSequence(mmap, pair.MakePair("k", O.None[int]())))
}

func TestSequence(t *testing.T) {
	seq := pair.Sequence[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
		O.Map[int, pair.Pair[string, int]],
	)

	assert.Equal(t, O.Some(pair.MakePair("k", 42)), seq(pair.MakePair("k", O.Some(42))))
	assert.Equal(t, O.None[pair.Pair[string, int]](), seq(pair.MakePair("k", O.None[int]())))
}

func TestMonadTraverse(t *testing.T) {
	mmap := O.MonadMap[int, pair.Pair[string, int]]
	positive := O.FromPredicate(N.MoreThan(0))

	assert.Equal(t, O.Some(pair.MakePair("k", 42)), pair.MonadTraverse(mmap, positive, pair.MakePair("k", 42)))
	assert.Equal(t, O.None[pair.Pair[string, int]](), pair.MonadTraverse(mmap, positive, pair.MakePair("k", -1)))
}

func TestTraverse(t *testing.T) {
	positive := O.FromPredicate(N.MoreThan(0))
	trav := pair.Traverse[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
		O.Map[int, pair.Pair[string, int]],
	)(positive)

	assert.Equal(t, O.Some(pair.MakePair("k", 42)), trav(pair.MakePair("k", 42)))
	assert.Equal(t, O.None[pair.Pair[string, int]](), trav(pair.MakePair("k", -1)))
}

func TestTraverseEqualsMapThenSequence(t *testing.T) {
	positive := O.FromPredicate(N.MoreThan(0))
	mmap := O.Map[int, pair.Pair[string, int]]

	viaTraverse := pair.Traverse[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](mmap)(positive)
	viaSequence := F.Flow2(
		pair.MapTail[string](positive),
		pair.Sequence[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](mmap),
	)

	for _, n := range []int{-1, 0, 1, 42} {
		p := pair.MakePair("k", n)
		assert.Equal(t, viaSequence(p), viaTraverse(p))
	}
}
