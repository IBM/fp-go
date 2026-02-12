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

// Package iter provides functional programming utilities for Go 1.23+ iterators.
//
// This package offers a comprehensive set of operations for working with lazy sequences
// using Go's native iter.Seq and iter.Seq2 types. It follows functional programming
// principles and provides monadic operations, transformations, and reductions.
//
// The package supports:
//   - Functor operations (Map, MapWithIndex, MapWithKey)
//   - Monad operations (Chain, Flatten, Ap)
//   - Filtering (Filter, FilterMap, FilterWithIndex, FilterWithKey)
//   - Folding and reduction (Reduce, Fold, FoldMap)
//   - Sequence construction (Of, From, MakeBy, Replicate)
//   - Sequence combination (Zip, Prepend, Append)
//
// All operations are lazy and only execute when the sequence is consumed via iteration.
//
// Example usage:
//
//	// Create a sequence and transform it
//	seq := From(1, 2, 3, 4, 5)
//	doubled := Map(N.Mul(2))(seq)
//
//	// Filter and reduce
//	evens := Filter(func(x int) bool { return x%2 == 0 })(doubled)
//	sum := MonadReduce(evens, func(acc, x int) int { return acc + x }, 0)
//	// sum = 20 (2+4+6+8+10 from doubled evens)
package iter

import (
	"slices"

	I "iter"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	G "github.com/IBM/fp-go/v2/internal/iter"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

// Of creates a sequence containing a single element.
//
// Example:
//
//	seq := Of(42)
//	// yields: 42
//
//go:inline
func Of[A any](a A) Seq[A] {
	return G.Of[Seq[A]](a)
}

// Of2 creates a key-value sequence containing a single key-value pair.
//
// Example:
//
//	seq := Of2("key", 100)
//	// yields: ("key", 100)
func Of2[K, A any](k K, a A) Seq2[K, A] {
	return func(yield func(K, A) bool) {
		yield(k, a)
	}
}

// MonadMap transforms each element in a sequence using the provided function.
// This is the monadic version that takes the sequence as the first parameter.
//
// RxJS Equivalent: [map] - https://rxjs.dev/api/operators/map
//
// Example:
//
//	seq := From(1, 2, 3)
//	result := MonadMap(seq, N.Mul(2))
//	// yields: 2, 4, 6
func MonadMap[A, B any](as Seq[A], f func(A) B) Seq[B] {
	return func(yield Predicate[B]) {
		for a := range as {
			if !yield(f(a)) {
				return
			}
		}
	}
}

// Map returns a function that transforms each element in a sequence.
// This is the curried version of MonadMap.
//
// Example:
//
//	double := Map(N.Mul(2))
//	seq := From(1, 2, 3)
//	result := double(seq)
//	// yields: 2, 4, 6
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

// MonadMapWithIndex transforms each element in a sequence using a function that also receives the element's index.
//
// Example:
//
//	seq := From("a", "b", "c")
//	result := MonadMapWithIndex(seq, func(i int, s string) string {
//	    return fmt.Sprintf("%d:%s", i, s)
//	})
//	// yields: "0:a", "1:b", "2:c"
func MonadMapWithIndex[A, B any](as Seq[A], f func(int, A) B) Seq[B] {
	return func(yield Predicate[B]) {
		var i int
		for a := range as {
			if !yield(f(i, a)) {
				return
			}
			i += 1
		}
	}
}

// MapWithIndex returns a function that transforms elements with their indices.
// This is the curried version of MonadMapWithIndex.
//
// Example:
//
//	addIndex := MapWithIndex(func(i int, s string) string {
//	    return fmt.Sprintf("%d:%s", i, s)
//	})
//	seq := From("a", "b", "c")
//	result := addIndex(seq)
//	// yields: "0:a", "1:b", "2:c"
//
//go:inline
func MapWithIndex[A, B any](f func(int, A) B) Operator[A, B] {
	return F.Bind2nd(MonadMapWithIndex[A, B], f)
}

// MonadMapWithKey transforms values in a key-value sequence using a function that receives both key and value.
//
// Example:
//
//	seq := Of2("x", 10)
//	result := MonadMapWithKey(seq, func(k string, v int) int { return v * 2 })
//	// yields: ("x", 20)
func MonadMapWithKey[K, A, B any](as Seq2[K, A], f func(K, A) B) Seq2[K, B] {
	return func(yield func(K, B) bool) {
		for k, a := range as {
			if !yield(k, f(k, a)) {
				return
			}
		}
	}
}

// MapWithKey returns a function that transforms values using their keys.
// This is the curried version of MonadMapWithKey.
//
// Example:
//
//	doubleValue := MapWithKey(func(k string, v int) int { return v * 2 })
//	seq := Of2("x", 10)
//	result := doubleValue(seq)
//	// yields: ("x", 20)
//
//go:inline
func MapWithKey[K, A, B any](f func(K, A) B) Operator2[K, A, B] {
	return F.Bind2nd(MonadMapWithKey[K, A, B], f)
}

// MonadFilter returns a sequence containing only elements that satisfy the predicate.
//
// RxJS Equivalent: [filter] - https://rxjs.dev/api/operators/filter
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
//	// yields: 2, 4
func MonadFilter[A any](as Seq[A], pred func(A) bool) Seq[A] {
	return func(yield Predicate[A]) {
		for a := range as {
			if pred(a) {
				if !yield(a) {
					return
				}
			}
		}
	}
}

// Filter returns a function that filters elements based on a predicate.
// This is the curried version of MonadFilter.
//
// Example:
//
//	evens := Filter(func(x int) bool { return x%2 == 0 })
//	seq := From(1, 2, 3, 4, 5)
//	result := evens(seq)
//	// yields: 2, 4
//
//go:inline
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return F.Bind2nd(MonadFilter[A], pred)
}

// MonadFilterWithIndex filters elements using a predicate that also receives the element's index.
//
// Example:
//
//	seq := From("a", "b", "c", "d")
//	result := MonadFilterWithIndex(seq, func(i int, s string) bool { return i%2 == 0 })
//	// yields: "a", "c" (elements at even indices)
func MonadFilterWithIndex[A any](as Seq[A], pred func(int, A) bool) Seq[A] {
	return func(yield Predicate[A]) {
		var i int
		for a := range as {
			if pred(i, a) {
				if !yield(a) {
					return
				}
			}
			i++
		}
	}
}

// FilterWithIndex returns a function that filters elements based on their index and value.
// This is the curried version of MonadFilterWithIndex.
//
// Example:
//
//	evenIndices := FilterWithIndex(func(i int, s string) bool { return i%2 == 0 })
//	seq := From("a", "b", "c", "d")
//	result := evenIndices(seq)
//	// yields: "a", "c"
//
//go:inline
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A] {
	return F.Bind2nd(MonadFilterWithIndex[A], pred)
}

// MonadFilterWithKey filters key-value pairs using a predicate that receives both key and value.
//
// Example:
//
//	seq := Of2("x", 10)
//	result := MonadFilterWithKey(seq, func(k string, v int) bool { return v > 5 })
//	// yields: ("x", 10)
func MonadFilterWithKey[K, A any](as Seq2[K, A], pred func(K, A) bool) Seq2[K, A] {
	return func(yield func(K, A) bool) {
		for k, a := range as {
			if pred(k, a) {
				if !yield(k, a) {
					return
				}
			}
		}
	}
}

// FilterWithKey returns a function that filters key-value pairs based on a predicate.
// This is the curried version of MonadFilterWithKey.
//
// Example:
//
//	largeValues := FilterWithKey(func(k string, v int) bool { return v > 5 })
//	seq := Of2("x", 10)
//	result := largeValues(seq)
//	// yields: ("x", 10)
//
//go:inline
func FilterWithKey[K, A any](pred func(K, A) bool) Operator2[K, A, A] {
	return F.Bind2nd(MonadFilterWithKey[K, A], pred)
}

// MonadFilterMap applies a function that returns an Option to each element,
// keeping only the Some values and unwrapping them.
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	result := MonadFilterMap(seq, func(x int) Option[int] {
//	    if x%2 == 0 {
//	        return option.Some(x * 10)
//	    }
//	    return option.None[int]()
//	})
//	// yields: 20, 40
func MonadFilterMap[A, B any](as Seq[A], f option.Kleisli[A, B]) Seq[B] {
	return func(yield Predicate[B]) {
		for a := range as {
			if b, ok := option.Unwrap(f(a)); ok {
				if !yield(b) {
					return
				}
			}
		}
	}
}

// FilterMap returns a function that filters and maps in one operation.
// This is the curried version of MonadFilterMap.
//
// Example:
//
//	evenDoubled := FilterMap(func(x int) Option[int] {
//	    if x%2 == 0 {
//	        return option.Some(x * 2)
//	    }
//	    return option.None[int]()
//	})
//	seq := From(1, 2, 3, 4)
//	result := evenDoubled(seq)
//	// yields: 4, 8
//
//go:inline
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B] {
	return F.Bind2nd(MonadFilterMap[A, B], f)
}

// MonadFilterMapWithIndex applies a function with index that returns an Option,
// keeping only the Some values.
//
// Example:
//
//	seq := From("a", "b", "c")
//	result := MonadFilterMapWithIndex(seq, func(i int, s string) Option[string] {
//	    if i%2 == 0 {
//	        return option.Some(fmt.Sprintf("%d:%s", i, s))
//	    }
//	    return option.None[string]()
//	})
//	// yields: "0:a", "2:c"
func MonadFilterMapWithIndex[A, B any](as Seq[A], f func(int, A) Option[B]) Seq[B] {
	return func(yield Predicate[B]) {
		var i int
		for a := range as {
			if b, ok := option.Unwrap(f(i, a)); ok {
				if !yield(b) {
					return
				}
			}
			i++
		}
	}
}

// FilterMapWithIndex returns a function that filters and maps with index.
// This is the curried version of MonadFilterMapWithIndex.
//
// Example:
//
//	evenIndexed := FilterMapWithIndex(func(i int, s string) Option[string] {
//	    if i%2 == 0 {
//	        return option.Some(s)
//	    }
//	    return option.None[string]()
//	})
//	seq := From("a", "b", "c", "d")
//	result := evenIndexed(seq)
//	// yields: "a", "c"
//
//go:inline
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B] {
	return F.Bind2nd(MonadFilterMapWithIndex[A, B], f)
}

// MonadFilterMapWithKey applies a function with key that returns an Option to key-value pairs,
// keeping only the Some values.
//
// Example:
//
//	seq := Of2("x", 10)
//	result := MonadFilterMapWithKey(seq, func(k string, v int) Option[int] {
//	    if v > 5 {
//	        return option.Some(v * 2)
//	    }
//	    return option.None[int]()
//	})
//	// yields: ("x", 20)
func MonadFilterMapWithKey[K, A, B any](as Seq2[K, A], f func(K, A) Option[B]) Seq2[K, B] {
	return func(yield func(K, B) bool) {
		for k, a := range as {
			if b, ok := option.Unwrap(f(k, a)); ok {
				if !yield(k, b) {
					return
				}
			}
		}
	}
}

// FilterMapWithKey returns a function that filters and maps key-value pairs.
// This is the curried version of MonadFilterMapWithKey.
//
// Example:
//
//	largeDoubled := FilterMapWithKey(func(k string, v int) Option[int] {
//	    if v > 5 {
//	        return option.Some(v * 2)
//	    }
//	    return option.None[int]()
//	})
//	seq := Of2("x", 10)
//	result := largeDoubled(seq)
//	// yields: ("x", 20)
//
//go:inline
func FilterMapWithKey[K, A, B any](f func(K, A) Option[B]) Operator2[K, A, B] {
	return F.Bind2nd(MonadFilterMapWithKey[K, A, B], f)
}

// MonadChain applies a function that returns a sequence to each element and flattens the results.
// This is the monadic bind operation (flatMap).
//
// RxJS Equivalent: [mergeMap/flatMap] - https://rxjs.dev/api/operators/mergeMap
//
// Example:
//
//	seq := From(1, 2, 3)
//	result := MonadChain(seq, func(x int) Seq[int] {
//	    return From(x, x*10)
//	})
//	// yields: 1, 10, 2, 20, 3, 30
func MonadChain[A, B any](as Seq[A], f Kleisli[A, B]) Seq[B] {
	return func(yield Predicate[B]) {
		for a := range as {
			for b := range f(a) {
				if !yield(b) {
					return
				}
			}
		}
	}
}

// Chain returns a function that chains (flatMaps) a sequence transformation.
// This is the curried version of MonadChain.
//
// Example:
//
//	duplicate := Chain(func(x int) Seq[int] { return From(x, x) })
//	seq := From(1, 2, 3)
//	result := duplicate(seq)
//	// yields: 1, 1, 2, 2, 3, 3
//
//go:inline
func Chain[A, B any](f func(A) Seq[B]) Operator[A, B] {
	return F.Bind2nd(MonadChain[A, B], f)
}

//go:inline
func FlatMap[A, B any](f func(A) Seq[B]) Operator[A, B] {
	return Chain(f)
}

// Flatten flattens a sequence of sequences into a single sequence.
//
// RxJS Equivalent: [mergeAll] - https://rxjs.dev/api/operators/mergeAll
//
// Example:
//
//	nested := From(From(1, 2), From(3, 4), From(5))
//	result := Flatten(nested)
//	// yields: 1, 2, 3, 4, 5
//
//go:inline
func Flatten[A any](mma Seq[Seq[A]]) Seq[A] {
	return MonadChain(mma, F.Identity[Seq[A]])
}

// MonadAp applies a sequence of functions to a sequence of values.
// This is the applicative apply operation.
//
// Example:
//
//	fns := From(N.Mul(2), N.Add(10))
//	vals := From(5, 3)
//	result := MonadAp(fns, vals)
//	// yields: 10, 6, 15, 13 (each function applied to each value)
//
//go:inline
func MonadAp[B, A any](fab Seq[func(A) B], fa Seq[A]) Seq[B] {
	return MonadChain(fab, F.Bind1st(MonadMap[A, B], fa))
}

// Ap returns a function that applies functions to values.
// This is the curried version of MonadAp.
//
// Example:
//
//	applyTo5 := Ap(From(5))
//	fns := From(N.Mul(2), N.Add(10))
//	result := applyTo5(fns)
//	// yields: 10, 15
//
//go:inline
func Ap[B, A any](fa Seq[A]) Operator[func(A) B, B] {
	return Chain(F.Bind1st(MonadMap[A, B], fa))
}

// From creates a sequence from a variadic list of elements.
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	// yields: 1, 2, 3, 4, 5
//
//go:inline
func From[A any](data ...A) Seq[A] {
	return slices.Values(data)
}

// Empty returns an empty sequence that yields no elements.
//
// Example:
//
//	seq := Empty[int]()
//	// yields nothing
//
//go:inline
func Empty[A any]() Seq[A] {
	return G.Empty[Seq[A]]()
}

// MakeBy creates a sequence of n elements by applying a function to each index.
// Returns an empty sequence if n <= 0.
//
// Example:
//
//	seq := MakeBy(5, func(i int) int { return i * i })
//	// yields: 0, 1, 4, 9, 16
func MakeBy[A any](n int, f func(int) A) Seq[A] {
	// sanity check
	if n <= 0 {
		return Empty[A]()
	}
	// run the generator function across the input
	return func(yield Predicate[A]) {
		for i := range n {
			if !yield(f(i)) {
				return
			}
		}
	}
}

// Replicate creates a sequence containing n copies of the same element.
//
// Example:
//
//	seq := Replicate(3, "hello")
//	// yields: "hello", "hello", "hello"
//
//go:inline
func Replicate[A any](n int, a A) Seq[A] {
	return MakeBy(n, F.Constant1[int](a))
}

// MonadReduce reduces a sequence to a single value by applying a function to each element
// and an accumulator, starting with an initial value.
//
// RxJS Equivalent: [reduce] - https://rxjs.dev/api/operators/reduce
//
// Example:
//
//	seq := From(1, 2, 3, 4, 5)
//	sum := MonadReduce(seq, func(acc, x int) int { return acc + x }, 0)
//	// returns: 15
//
//go:inline
func MonadReduce[A, B any](fa Seq[A], f func(B, A) B, initial B) B {
	return G.MonadReduce(fa, f, initial)
}

// Reduce returns a function that reduces a sequence to a single value.
// This is the curried version of MonadReduce.
//
// Example:
//
//	sum := Reduce(func(acc, x int) int { return acc + x }, 0)
//	seq := From(1, 2, 3, 4, 5)
//	result := sum(seq)
//	// returns: 15
func Reduce[A, B any](f func(B, A) B, initial B) func(Seq[A]) B {
	return func(fa Seq[A]) B {
		return MonadReduce(fa, f, initial)
	}
}

// MonadReduceWithIndex reduces a sequence using a function that also receives the element's index.
//
// Example:
//
//	seq := From(10, 20, 30)
//	result := MonadReduceWithIndex(seq, func(i, acc, x int) int {
//	    return acc + (i * x)
//	}, 0)
//	// returns: 0*10 + 1*20 + 2*30 = 80
//
//go:inline
func MonadReduceWithIndex[A, B any](fa Seq[A], f func(int, B, A) B, initial B) B {
	return G.MonadReduceWithIndex(fa, f, initial)
}

// ReduceWithIndex returns a function that reduces with index.
// This is the curried version of MonadReduceWithIndex.
//
// Example:
//
//	weightedSum := ReduceWithIndex(func(i, acc, x int) int {
//	    return acc + (i * x)
//	}, 0)
//	seq := From(10, 20, 30)
//	result := weightedSum(seq)
//	// returns: 80
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func(Seq[A]) B {
	return func(fa Seq[A]) B {
		return MonadReduceWithIndex(fa, f, initial)
	}
}

// MonadReduceWithKey reduces a key-value sequence using a function that receives the key.
//
// Example:
//
//	seq := Of2("x", 10)
//	result := MonadReduceWithKey(seq, func(k string, acc int, v int) int {
//	    return acc + v
//	}, 0)
//	// returns: 10
func MonadReduceWithKey[K, A, B any](fa Seq2[K, A], f func(K, B, A) B, initial B) B {
	current := initial
	for k, a := range fa {
		current = f(k, current, a)
	}
	return current
}

// ReduceWithKey returns a function that reduces key-value pairs.
// This is the curried version of MonadReduceWithKey.
//
// Example:
//
//	sumValues := ReduceWithKey(func(k string, acc int, v int) int {
//	    return acc + v
//	}, 0)
//	seq := Of2("x", 10)
//	result := sumValues(seq)
//	// returns: 10
func ReduceWithKey[K, A, B any](f func(K, B, A) B, initial B) func(Seq2[K, A]) B {
	return func(fa Seq2[K, A]) B {
		return MonadReduceWithKey(fa, f, initial)
	}
}

// MonadFold folds a sequence using a monoid's concat operation and empty value.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/number"
//	seq := From(1, 2, 3, 4, 5)
//	sum := MonadFold(seq, number.MonoidSum[int]())
//	// returns: 15
//
//go:inline
func MonadFold[A any](fa Seq[A], m M.Monoid[A]) A {
	return MonadReduce(fa, m.Concat, m.Empty())
}

// Fold returns a function that folds a sequence using a monoid.
// This is the curried version of MonadFold.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/number"
//	sumAll := Fold(number.MonoidSum[int]())
//	seq := From(1, 2, 3, 4, 5)
//	result := sumAll(seq)
//	// returns: 15
//
//go:inline
func Fold[A any](m M.Monoid[A]) func(Seq[A]) A {
	return Reduce(m.Concat, m.Empty())
}

// MonadFoldMap maps each element to a monoid value and combines them using the monoid.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//	seq := From(1, 2, 3)
//	result := MonadFoldMap(seq, func(x int) string {
//	    return fmt.Sprintf("%d ", x)
//	}, string.Monoid)
//	// returns: "1 2 3 "
//
//go:inline
func MonadFoldMap[A, B any](fa Seq[A], f func(A) B, m M.Monoid[B]) B {
	return MonadFold(MonadMap(fa, f), m)
}

// FoldMap returns a function that maps and folds using a monoid.
// This is the curried version of MonadFoldMap.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//	stringify := FoldMap(string.Monoid)(func(x int) string {
//	    return fmt.Sprintf("%d ", x)
//	})
//	seq := From(1, 2, 3)
//	result := stringify(seq)
//	// returns: "1 2 3 "
//
//go:inline
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func(Seq[A]) B {
	return F.Pipe1(
		Map[A, B],
		reader.Map[func(A) B](reader.Map[Seq[A]](Fold(m))),
	)
}

// MonadFoldMapWithIndex maps each element with its index to a monoid value and combines them.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//	seq := From("a", "b", "c")
//	result := MonadFoldMapWithIndex(seq, func(i int, s string) string {
//	    return fmt.Sprintf("%d:%s ", i, s)
//	}, string.Monoid)
//	// returns: "0:a 1:b 2:c "
//
//go:inline
func MonadFoldMapWithIndex[A, B any](fa Seq[A], f func(int, A) B, m M.Monoid[B]) B {
	return MonadReduceWithIndex(fa, func(i int, b B, a A) B {
		return m.Concat(b, f(i, a))
	}, m.Empty())
}

// FoldMapWithIndex returns a function that maps with index and folds.
// This is the curried version of MonadFoldMapWithIndex.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//	indexedStringify := FoldMapWithIndex(string.Monoid)(func(i int, s string) string {
//	    return fmt.Sprintf("%d:%s ", i, s)
//	})
//	seq := From("a", "b", "c")
//	result := indexedStringify(seq)
//	// returns: "0:a 1:b 2:c "
//
//go:inline
func FoldMapWithIndex[A, B any](m M.Monoid[B]) func(func(int, A) B) func(Seq[A]) B {
	return func(f func(int, A) B) func(Seq[A]) B {
		return func(as Seq[A]) B {
			return MonadFoldMapWithIndex(as, f, m)
		}
	}
}

// MonadFoldMapWithKey maps each key-value pair to a monoid value and combines them.
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/string"
//	seq := Of2("x", 10)
//	result := MonadFoldMapWithKey(seq, func(k string, v int) string {
//	    return fmt.Sprintf("%s:%d ", k, v)
//	}, string.Monoid)
//	// returns: "x:10 "
//
//go:inline
func MonadFoldMapWithKey[K, A, B any](fa Seq2[K, A], f func(K, A) B, m M.Monoid[B]) B {
	return MonadReduceWithKey(fa, func(k K, b B, a A) B {
		return m.Concat(b, f(k, a))
	}, m.Empty())
}

// FoldMapWithKey returns a function that maps with key and folds.
// This is the curried version of MonadFoldMapWithKey.
//
//go:inline
func FoldMapWithKey[K, A, B any](m M.Monoid[B]) func(func(K, A) B) func(Seq2[K, A]) B {
	return func(f func(K, A) B) func(Seq2[K, A]) B {
		return func(as Seq2[K, A]) B {
			return MonadFoldMapWithKey(as, f, m)
		}
	}
}

// MonadFlap applies a fixed value to a sequence of functions.
// This is the dual of MonadAp.
//
// Example:
//
//	fns := From(N.Mul(2), N.Add(10))
//	result := MonadFlap(fns, 5)
//	// yields: 10, 15
//
//go:inline
func MonadFlap[B, A any](fab Seq[func(A) B], a A) Seq[B] {
	return functor.MonadFlap(MonadMap[func(A) B, B], fab, a)
}

// Flap returns a function that applies a fixed value to functions.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return functor.Flap(Map[func(A) B, B], a)
}

// Prepend returns a function that adds an element to the beginning of a sequence.
//
// RxJS Equivalent: [startWith] - https://rxjs.dev/api/operators/startWith
//
// Example:
//
//	seq := From(2, 3, 4)
//	result := Prepend(1)(seq)
//	// yields: 1, 2, 3, 4
//
//go:inline
func Prepend[A any](head A) Operator[A, A] {
	return G.Prepend[Seq[A]](head)
}

// Append returns a function that adds an element to the end of a sequence.
//
// RxJS Equivalent: [endWith] - https://rxjs.dev/api/operators/endWith
//
// Example:
//
//	seq := From(1, 2, 3)
//	result := Append(4)(seq)
//	// yields: 1, 2, 3, 4
//
//go:inline
func Append[A any](tail A) Operator[A, A] {
	return G.Append[Seq[A]](tail)
}

// MonadZip combines two sequences into a sequence of pairs.
// The resulting sequence stops when either input sequence is exhausted.
//
// RxJS Equivalent: [zip] - https://rxjs.dev/api/operators/zip
//
// Example:
//
//	seqA := From(1, 2, 3)
//	seqB := From("a", "b")
//	result := MonadZip(seqA, seqB)
//	// yields: (1, "a"), (2, "b")
func MonadZip[A, B any](fa Seq[A], fb Seq[B]) Seq2[A, B] {

	return func(yield func(A, B) bool) {
		na, sa := I.Pull(fa)
		defer sa()

		for b := range fb {
			a, ok := na()
			if !ok {
				return
			}

			if !yield(a, b) {
				return
			}
		}
	}
}

// Zip returns a function that zips a sequence with another sequence.
// This is the curried version of MonadZip.
//
// Example:
//
//	seqA := From(1, 2, 3)
//	zipWithA := Zip(seqA)
//	seqB := From("a", "b", "c")
//	result := zipWithA(seqB)
//	// yields: (1, "a"), (2, "b"), (3, "c")
//
//go:inline
func Zip[A, B any](fb Seq[B]) func(Seq[A]) Seq2[A, B] {
	return F.Bind2nd(MonadZip[A, B], fb)
}

// MonadMapToArray maps each element in a sequence using a function and collects the results into an array.
// This is a convenience function that combines Map and collection into a single operation.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of elements in the output array
//
// Parameters:
//   - fa: The input sequence to map
//   - f: The mapping function to apply to each element
//
// Returns:
//   - A slice containing all mapped elements
//
// Example:
//
//	seq := From(1, 2, 3)
//	result := MonadMapToArray(seq, N.Mul(2))
//	// returns: []int{2, 4, 6}
//
//go:inline
func MonadMapToArray[A, B any](fa Seq[A], f func(A) B) []B {
	return G.MonadMapToArray[Seq[A], []B](fa, f)
}

// MapToArray returns a function that maps elements and collects them into an array.
// This is the curried version of MonadMapToArray.
//
// Type Parameters:
//   - A: The type of elements in the input sequence
//   - B: The type of elements in the output array
//
// Parameters:
//   - f: The mapping function to apply to each element
//
// Returns:
//   - A function that takes a sequence and returns a slice of mapped elements
//
// Example:
//
//	double := MapToArray(N.Mul(2))
//	seq := From(1, 2, 3)
//	result := double(seq)
//	// returns: []int{2, 4, 6}
//
//go:inline
func MapToArray[A, B any](f func(A) B) func(Seq[A]) []B {
	return G.MapToArray[Seq[A], []B](f)
}

// ToSeqPair converts a key-value sequence (Seq2) into a sequence of Pairs.
//
// This function transforms a Seq2[A, B] (which yields key-value pairs when iterated)
// into a Seq[Pair[A, B]] (which yields Pair objects). This is useful when you need
// to work with pairs as first-class values rather than as separate key-value arguments.
//
// Type Parameters:
//   - A: The type of the first element (key) in each pair
//   - B: The type of the second element (value) in each pair
//
// Parameters:
//   - as: A Seq2 that yields key-value pairs
//
// Returns:
//   - A Seq that yields Pair objects containing the key-value pairs
//
// Example - Basic conversion:
//
//	seq2 := iter.MonadZip(iter.From("a", "b", "c"), iter.From(1, 2, 3))
//	pairs := iter.ToSeqPair(seq2)
//	// yields: Pair("a", 1), Pair("b", 2), Pair("c", 3)
//
// Example - Using with Map:
//
//	seq2 := iter.MonadZip(iter.From(1, 2, 3), iter.From(10, 20, 30))
//	pairs := iter.ToSeqPair(seq2)
//	sums := iter.MonadMap(pairs, func(p Pair[int, int]) int {
//	    return p.Fst + p.Snd
//	})
//	// yields: 11, 22, 33
//
// Example - Empty sequence:
//
//	seq2 := iter.Empty[int]()
//	zipped := iter.MonadZip(seq2, iter.Empty[string]())
//	pairs := iter.ToSeqPair(zipped)
//	// yields: nothing (empty sequence)
func ToSeqPair[A, B any](as Seq2[A, B]) Seq[Pair[A, B]] {
	return func(yield Predicate[Pair[A, B]]) {
		for a, b := range as {
			if !yield(pair.MakePair(a, b)) {
				return
			}
		}
	}
}
