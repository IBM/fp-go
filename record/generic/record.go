// Copyright (c) 2023 IBM Corp.
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

package generic

import (
	F "github.com/IBM/fp-go/function"
	G "github.com/IBM/fp-go/internal/record"
	Mg "github.com/IBM/fp-go/magma"
	Mo "github.com/IBM/fp-go/monoid"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

func IsEmpty[M ~map[K]V, K comparable, V any](r M) bool {
	return len(r) == 0
}

func IsNonEmpty[M ~map[K]V, K comparable, V any](r M) bool {
	return len(r) > 0
}

func Keys[M ~map[K]V, GK ~[]K, K comparable, V any](r M) GK {
	return collect[M, GK](r, F.First[K, V])
}

func Values[M ~map[K]V, GV ~[]V, K comparable, V any](r M) GV {
	return collect[M, GV](r, F.Second[K, V])
}

func collect[M ~map[K]V, GR ~[]R, K comparable, V, R any](r M, f func(K, V) R) GR {
	count := len(r)
	result := make(GR, count)
	idx := 0
	for k, v := range r {
		result[idx] = f(k, v)
		idx++
	}
	return result
}

func Collect[M ~map[K]V, GR ~[]R, K comparable, V, R any](f func(K, V) R) func(M) GR {
	return F.Bind2nd(collect[M, GR, K, V, R], f)
}

func Reduce[M ~map[K]V, K comparable, V, R any](f func(R, V) R, initial R) func(M) R {
	return func(r M) R {
		return G.Reduce(r, f, initial)
	}
}

func ReduceWithIndex[M ~map[K]V, K comparable, V, R any](f func(K, R, V) R, initial R) func(M) R {
	return func(r M) R {
		return G.ReduceWithIndex(r, f, initial)
	}
}

func ReduceRef[M ~map[K]V, K comparable, V, R any](f func(R, *V) R, initial R) func(M) R {
	return func(r M) R {
		return G.ReduceRef(r, f, initial)
	}
}

func ReduceRefWithIndex[M ~map[K]V, K comparable, V, R any](f func(K, R, *V) R, initial R) func(M) R {
	return func(r M) R {
		return G.ReduceRefWithIndex(r, f, initial)
	}
}

func MonadMap[M ~map[K]V, N ~map[K]R, K comparable, V, R any](r M, f func(V) R) N {
	return MonadMapWithIndex[M, N](r, F.Ignore1of2[K](f))
}

func MonadChainWithIndex[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N], r M, f func(K, V1) N) N {
	return G.ReduceWithIndex(r, func(k K, dst N, b V1) N {
		return m.Concat(dst, f(k, b))
	}, m.Empty())
}

func MonadChain[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N], r M, f func(V1) N) N {
	return G.Reduce(r, func(dst N, b V1) N {
		return m.Concat(dst, f(b))
	}, m.Empty())
}

func ChainWithIndex[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N]) func(func(K, V1) N) func(M) N {
	return func(f func(K, V1) N) func(M) N {
		return func(ma M) N {
			return MonadChainWithIndex(m, ma, f)
		}
	}
}

func Chain[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N]) func(func(V1) N) func(M) N {
	return func(f func(V1) N) func(M) N {
		return func(ma M) N {
			return MonadChain(m, ma, f)
		}
	}
}

func MonadMapWithIndex[M ~map[K]V, N ~map[K]R, K comparable, V, R any](r M, f func(K, V) R) N {
	return G.ReduceWithIndex(r, func(k K, dst N, v V) N {
		return upsertAtReadWrite(dst, k, f(k, v))
	}, make(N, len(r)))
}

func MonadMapRefWithIndex[M ~map[K]V, N ~map[K]R, K comparable, V, R any](r M, f func(K, *V) R) N {
	return G.ReduceRefWithIndex(r, func(k K, dst N, v *V) N {
		return upsertAtReadWrite(dst, k, f(k, v))
	}, make(N, len(r)))
}

func MonadMapRef[M ~map[K]V, N ~map[K]R, K comparable, V, R any](r M, f func(*V) R) N {
	return MonadMapRefWithIndex[M, N](r, F.Ignore1of2[K](f))
}

func Map[M ~map[K]V, N ~map[K]R, K comparable, V, R any](f func(V) R) func(M) N {
	return F.Bind2nd(MonadMap[M, N, K, V, R], f)
}

func MapRef[M ~map[K]V, N ~map[K]R, K comparable, V, R any](f func(*V) R) func(M) N {
	return F.Bind2nd(MonadMapRef[M, N, K, V, R], f)
}

func MapWithIndex[M ~map[K]V, N ~map[K]R, K comparable, V, R any](f func(K, V) R) func(M) N {
	return F.Bind2nd(MonadMapWithIndex[M, N, K, V, R], f)
}

func MapRefWithIndex[M ~map[K]V, N ~map[K]R, K comparable, V, R any](f func(K, *V) R) func(M) N {
	return F.Bind2nd(MonadMapRefWithIndex[M, N, K, V, R], f)
}

func Lookup[M ~map[K]V, K comparable, V any](k K) func(M) O.Option[V] {
	n := O.None[V]()
	return func(m M) O.Option[V] {
		if val, ok := m[k]; ok {
			return O.Some(val)
		}
		return n
	}
}

func Has[M ~map[K]V, K comparable, V any](k K, r M) bool {
	_, ok := r[k]
	return ok
}

func union[M ~map[K]V, K comparable, V any](m Mg.Magma[V], left M, right M) M {
	lenLeft := len(left)

	if lenLeft == 0 {
		return right
	}

	lenRight := len(right)
	if lenRight == 0 {
		return left
	}

	result := make(M, lenLeft+lenRight)

	for k, v := range left {
		if val, ok := right[k]; ok {
			result[k] = m.Concat(v, val)
		} else {
			result[k] = v
		}
	}

	for k, v := range right {
		if _, ok := left[k]; !ok {
			result[k] = v
		}
	}

	return result
}

func unionLast[M ~map[K]V, K comparable, V any](left M, right M) M {
	lenLeft := len(left)

	if lenLeft == 0 {
		return right
	}

	lenRight := len(right)
	if lenRight == 0 {
		return left
	}

	result := make(M, lenLeft+lenRight)

	for k, v := range left {
		result[k] = v
	}

	for k, v := range right {
		result[k] = v
	}

	return result
}

func Union[M ~map[K]V, K comparable, V any](m Mg.Magma[V]) func(M) func(M) M {
	return func(right M) func(M) M {
		return func(left M) M {
			return union(m, left, right)
		}
	}
}

func UnionLast[M ~map[K]V, K comparable, V any](right M) func(M) M {
	return func(left M) M {
		return unionLast(left, right)
	}
}

func Merge[M ~map[K]V, K comparable, V any](right M) func(M) M {
	return UnionLast(right)
}

func UnionFirst[M ~map[K]V, K comparable, V any](right M) func(M) M {
	return func(left M) M {
		return unionLast(right, left)
	}
}

func Empty[M ~map[K]V, K comparable, V any]() M {
	return make(M)
}

func Size[M ~map[K]V, K comparable, V any](r M) int {
	return len(r)
}

func ToArray[M ~map[K]V, GT ~[]T.Tuple2[K, V], K comparable, V any](r M) GT {
	return collect[M, GT](r, T.MakeTuple2[K, V])
}

func ToEntries[M ~map[K]V, GT ~[]T.Tuple2[K, V], K comparable, V any](r M) GT {
	return ToArray[M, GT](r)
}

func FromEntries[M ~map[K]V, GT ~[]T.Tuple2[K, V], K comparable, V any](fa GT) M {
	m := make(M)
	for _, t := range fa {
		upsertAtReadWrite(m, t.F1, t.F2)
	}
	return m
}

func duplicate[M ~map[K]V, K comparable, V any](r M) M {
	return MonadMap[M, M](r, F.Identity[V])
}

func upsertAt[M ~map[K]V, K comparable, V any](r M, k K, v V) M {
	dup := duplicate(r)
	dup[k] = v
	return dup
}

func deleteAt[M ~map[K]V, K comparable, V any](r M, k K) M {
	dup := duplicate(r)
	delete(dup, k)
	return dup
}

func upsertAtReadWrite[M ~map[K]V, K comparable, V any](r M, k K, v V) M {
	r[k] = v
	return r
}

func UpsertAt[M ~map[K]V, K comparable, V any](k K, v V) func(M) M {
	return func(ma M) M {
		return upsertAt(ma, k, v)
	}
}

func DeleteAt[M ~map[K]V, K comparable, V any](k K) func(M) M {
	return F.Bind2nd(deleteAt[M, K, V], k)
}

func Singleton[M ~map[K]V, K comparable, V any](k K, v V) M {
	return M{k: v}
}

func filterMapWithIndex[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](fa M, f func(K, V1) O.Option[V2]) N {
	return G.ReduceWithIndex(fa, func(key K, n N, value V1) N {
		return O.MonadFold(f(key, value), F.Constant(n), func(v V2) N {
			return upsertAtReadWrite(n, key, v)
		})
	}, make(N))
}

func filterWithIndex[M ~map[K]V, K comparable, V any](fa M, f func(K, V) bool) M {
	return filterMapWithIndex[M, M](fa, func(k K, v V) O.Option[V] {
		if f(k, v) {
			return O.Of(v)
		}
		return O.None[V]()
	})
}

func filter[M ~map[K]V, K comparable, V any](fa M, f func(K) bool) M {
	return filterWithIndex(fa, F.Ignore2of2[V](f))
}

// Filter creates a new map with only the elements that match the predicate
func Filter[M ~map[K]V, K comparable, V any](f func(K) bool) func(M) M {
	return F.Bind2nd(filter[M, K, V], f)
}

// FilterWithIndex creates a new map with only the elements that match the predicate
func FilterWithIndex[M ~map[K]V, K comparable, V any](f func(K, V) bool) func(M) M {
	return F.Bind2nd(filterWithIndex[M, K, V], f)
}

// FilterMapWithIndex creates a new map with only the elements for which the transformation function creates a Some
func FilterMapWithIndex[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](f func(K, V1) O.Option[V2]) func(M) N {
	return F.Bind2nd(filterMapWithIndex[M, N, K, V1, V2], f)
}

// FilterMap creates a new map with only the elements for which the transformation function creates a Some
func FilterMap[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](f func(V1) O.Option[V2]) func(M) N {
	return F.Bind2nd(filterMapWithIndex[M, N, K, V1, V2], F.Ignore1of2[K](f))
}

// Flatten converts a nested map into a regular map
func Flatten[M ~map[K]N, N ~map[K]V, K comparable, V any](m Mo.Monoid[N]) func(M) N {
	return Chain[M, N](m)(F.Identity[N])
}

// FilterChainWithIndex creates a new map with only the elements for which the transformation function creates a Some
func FilterChainWithIndex[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N]) func(func(K, V1) O.Option[N]) func(M) N {
	flatten := Flatten[map[K]N, N](m)
	return func(f func(K, V1) O.Option[N]) func(M) N {
		return F.Flow2(
			FilterMapWithIndex[M, map[K]N](f),
			flatten,
		)
	}
}

// FilterChain creates a new map with only the elements for which the transformation function creates a Some
func FilterChain[M ~map[K]V1, N ~map[K]V2, K comparable, V1, V2 any](m Mo.Monoid[N]) func(func(V1) O.Option[N]) func(M) N {
	flatten := Flatten[map[K]N, N](m)
	return func(f func(V1) O.Option[N]) func(M) N {
		return F.Flow2(
			FilterMap[M, map[K]N](f),
			flatten,
		)
	}
}

// IsNil checks if the map is set to nil
func IsNil[M ~map[K]V, K comparable, V any](m M) bool {
	return m == nil
}

// IsNonNil checks if the map is set to nil
func IsNonNil[M ~map[K]V, K comparable, V any](m M) bool {
	return m != nil
}

// ConstNil return a nil map
func ConstNil[M ~map[K]V, K comparable, V any]() M {
	return (M)(nil)
}
