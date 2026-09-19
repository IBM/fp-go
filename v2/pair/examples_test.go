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
	"fmt"
	"strconv"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	SG "github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/IBM/fp-go/v2/tuple"
)

func ExampleOf() {
	p := pair.Of(42)
	fmt.Println(p)
	// Output:
	// Pair[int, int](42, 42)
}

func ExampleFromTuple() {
	t := tuple.MakeTuple2("hello", 42)
	p := pair.FromTuple(t)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleFromHead() {
	makePair := pair.FromHead[int]("hello")
	p := makePair(42)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleFromTail() {
	makePair := pair.FromTail[string](42)
	p := makePair("hello")
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleToTuple() {
	p := pair.MakePair("hello", 42)
	t := pair.ToTuple(p)
	fmt.Println(t)
	// Output:
	// Tuple2[string, int](hello, 42)
}

func ExampleMakePair() {
	p := pair.MakePair("hello", 42)
	fmt.Println(p)
	// Output:
	// Pair[string, int](hello, 42)
}

func ExampleHead() {
	p := pair.MakePair("hello", 42)
	h := pair.Head(p)
	fmt.Println(h)
	// Output:
	// hello
}

func ExampleTail() {
	p := pair.MakePair("hello", 42)
	t := pair.Tail(p)
	fmt.Println(t)
	// Output:
	// 42
}

func ExampleFirst() {
	p := pair.MakePair("hello", 42)
	f := pair.First(p)
	fmt.Println(f)
	// Output:
	// hello
}

func ExampleSecond() {
	p := pair.MakePair("hello", 42)
	s := pair.Second(p)
	fmt.Println(s)
	// Output:
	// 42
}

func ExampleSwap() {
	p := pair.MakePair("hello", 42)
	swapped := pair.Swap(p)
	fmt.Println(swapped)
	// Output:
	// Pair[int, string](42, hello)
}

func ExampleMap() {
	p := pair.MakePair("hello", 42)
	doubled := pair.Map[string](N.Mul(2))(p)
	fmt.Println(doubled)
	// Output:
	// Pair[string, int](hello, 84)
}

func ExampleMapHead() {
	p := pair.MakePair("hello", 42)
	upper := pair.MapHead[int](func(s string) string { return s + "!" })(p)
	fmt.Println(upper)
	// Output:
	// Pair[string, int](hello!, 42)
}

func ExampleBiMap() {
	p := pair.MakePair("hello", 42)
	result := pair.BiMap(
		func(s string) string { return s + "!" },
		N.Mul(2),
	)(p)
	fmt.Println(result)
	// Output:
	// Pair[string, int](hello!, 84)
}

func ExampleMapTail() {
	p := pair.MakePair("hello", "world")
	result := pair.MapTail[string](S.Size)(p)
	fmt.Println(result)
	// Output:
	// Pair[string, int](hello, 5)
}

func ExampleChain() {
	// a Writer-like computation: the head accumulates a count, the tail carries the value
	countLen := func(s string) pair.Pair[int, int] {
		return pair.MakePair(1, S.Size(s))
	}
	result := F.Pipe1(
		pair.MakePair(1, "hello"),
		pair.Chain(N.SemigroupSum[int](), countLen),
	)
	fmt.Println(result)
	// Output:
	// Pair[int, int](2, 5)
}

func ExampleChainHead() {
	strConcat := SG.MakeSemigroup(S.Monoid.Concat)
	toString := func(n int) pair.Pair[string, string] {
		return pair.MakePair(strconv.Itoa(n), "!")
	}
	result := F.Pipe1(
		pair.MakePair(5, "hello"),
		pair.ChainHead(strConcat, toString),
	)
	fmt.Println(result)
	// Output:
	// Pair[string, string](5, hello!)
}

func ExampleAp() {
	pv := pair.MakePair(5, "hello")
	pf := pair.MakePair(10, S.Size)
	result := pair.Ap[int, string, int](N.SemigroupSum[int](), pv)(pf)
	fmt.Println(result)
	// Output:
	// Pair[int, int](15, 5)
}

func ExampleApHead() {
	strConcat := SG.MakeSemigroup(S.Monoid.Concat)
	pv := pair.MakePair(42, "hello")
	pf := pair.MakePair(strconv.Itoa, "!")
	result := pair.ApHead[string, int, string](strConcat, pv)(pf)
	fmt.Println(result)
	// Output:
	// Pair[string, string](42, hello!)
}

func ExamplePaired() {
	add := pair.Paired(N.MonoidSum[int]().Concat)
	fmt.Println(add(pair.MakePair(3, 4)))
	// Output:
	// 7
}

func ExampleUnpaired() {
	sum := func(p pair.Pair[int, int]) int {
		return pair.Head(p) + pair.Tail(p)
	}
	fmt.Println(pair.Unpaired(sum)(3, 4))
	// Output:
	// 7
}

func ExampleMerge() {
	// N.Sub is data-last: N.Sub(b)(a) == a - b
	sub := pair.Merge(N.Sub[int])
	fmt.Println(sub(pair.MakePair(10, 3)))
	// Output:
	// 7
}

func ExampleUnpack() {
	name, age := pair.Unpack(pair.MakePair("Alice", 30))
	fmt.Printf("%s is %d years old\n", name, age)
	// Output:
	// Alice is 30 years old
}

func ExampleZero() {
	fmt.Println(pair.Zero[string, int]())
	// Output:
	// Pair[string, int](, 0)
}

func ExampleEq() {
	pairEq := pair.Eq(EQ.FromStrictEquals[string](), EQ.FromStrictEquals[int]())
	fmt.Println(pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 1)))
	fmt.Println(pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("a", 2)))
	// Output:
	// true
	// false
}

func ExampleFromStrictEquals() {
	pairEq := pair.FromStrictEquals[string, int]()
	fmt.Println(pairEq.Equals(pair.MakePair("a", 1), pair.MakePair("b", 1)))
	// Output:
	// false
}

func ExampleSequence() {
	seq := pair.Sequence[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
		O.Map[int, pair.Pair[string, int]],
	)
	fmt.Println(seq(pair.MakePair("key", O.Some(42))))
	fmt.Println(seq(pair.MakePair("key", O.None[int]())))
	// Output:
	// Some[pair.Pair[string,int]](Pair[string, int](key, 42))
	// None[pair.Pair[string,int]]
}

func ExampleTraverse() {
	positive := O.FromPredicate(N.MoreThan(0))
	trav := pair.Traverse[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
		O.Map[int, pair.Pair[string, int]],
	)(positive)
	fmt.Println(trav(pair.MakePair("key", 42)))
	fmt.Println(trav(pair.MakePair("key", -1)))
	// Output:
	// Some[pair.Pair[string,int]](Pair[string, int](key, 42))
	// None[pair.Pair[string,int]]
}

func ExampleMonad() {
	// Pair as a Writer monad: the head accumulates a log using a monoid
	m := pair.Monad[int, string, int](S.Monoid)

	double := func(n int) pair.Pair[string, int] {
		return pair.MakePair("doubled;", n*2)
	}
	addTen := func(n int) pair.Pair[string, int] {
		return pair.MakePair("added 10;", n+10)
	}

	result := F.Pipe2(
		m.Of(5),
		m.Chain(double),
		m.Chain(addTen),
	)
	fmt.Println(result)
	// Output:
	// Pair[string, int](doubled;added 10;, 20)
}
