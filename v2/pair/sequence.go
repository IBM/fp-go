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

import F "github.com/IBM/fp-go/v2/function"

// MonadSequence turns a [Pair] whose tail is a higher kinded type into a higher kinded type
// of a [Pair], keeping the head unchanged.
//
// The higher kinded type is described by its map function mmap, which receives the effect HKT[A]
// and the function that re-attaches the head to each value A.
//
// Example:
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	p := pair.MakePair("key", O.Some(42))
//	result := pair.MonadSequence(O.MonadMap[int, pair.Pair[string, int]], p)
//	// O.Some(Pair[string, int]{"key", 42})
//
//go:inline
func MonadSequence[L, A, HKTA, HKTPA any](
	mmap func(HKTA, Kleisli[L, A, A]) HKTPA,
	fas Pair[L, HKTA],
) HKTPA {
	return mmap(Tail(fas), FromHead[A](Head(fas)))
}

// MonadTraverse maps the tail of a [Pair] with an effectful function f and turns the result
// into a higher kinded type of a [Pair], keeping the head unchanged.
//
// It is equivalent to MonadSequence(mmap, MonadMapTail(fas, f)).
//
// Example:
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	p := pair.MakePair("key", 42)
//	result := pair.MonadTraverse(O.MonadMap[int, pair.Pair[string, int]], O.Of[int], p)
//	// O.Some(Pair[string, int]{"key", 42})
//
//go:inline
func MonadTraverse[L, A, HKTA, HKTPA any](
	mmap func(HKTA, Kleisli[L, A, A]) HKTPA,
	f func(A) HKTA,
	fas Pair[L, A],
) HKTPA {
	return MonadSequence(mmap, MonadMapTail(fas, f))
}

// Sequence is the curried, data-last version of [MonadSequence].
// It takes the curried map function of the higher kinded type and returns a function
// that turns a Pair[L, HKT[A]] into HKT[Pair[L, A]].
//
// Example:
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	seq := pair.Sequence[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
//	    O.Map[int, pair.Pair[string, int]],
//	)
//	result := seq(pair.MakePair("key", O.Some(42)))
//	// O.Some(Pair[string, int]{"key", 42})
//
//go:inline
func Sequence[L, A, HKTA, HKTPA any](
	mmap func(Kleisli[L, A, A]) func(HKTA) HKTPA,
) func(Pair[L, HKTA]) HKTPA {
	return Paired(F.Uncurry2(F.Flow2(
		FromHead[A, L],
		mmap,
	)))
}

// Traverse is the curried, data-last version of [MonadTraverse].
// It takes the curried map function of the higher kinded type and an effectful function f
// and returns a function that turns a Pair[L, A] into HKT[Pair[L, A]].
//
// It is equivalent to F.Flow2(MapTail[L](f), Sequence(mmap)).
//
// Example:
//
//	import O "github.com/IBM/fp-go/v2/option"
//
//	positive := O.FromPredicate(N.MoreThan(0))
//	trav := pair.Traverse[string, int, O.Option[int], O.Option[pair.Pair[string, int]]](
//	    O.Map[int, pair.Pair[string, int]],
//	)(positive)
//	trav(pair.MakePair("key", 42))  // O.Some(Pair[string, int]{"key", 42})
//	trav(pair.MakePair("key", -1))  // O.None[Pair[string, int]]()
//
//go:inline
func Traverse[L, A, HKTA, HKTPA any](
	mmap func(Kleisli[L, A, A]) func(HKTA) HKTPA,
) func(func(A) HKTA) func(Pair[L, A]) HKTPA {
	seq := Sequence(mmap)
	return func(f func(A) HKTA) func(Pair[L, A]) HKTPA {
		return F.Flow2(
			MapTail[L](f),
			seq,
		)
	}
}
