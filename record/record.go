package record

import (
	Mg "github.com/ibm/fp-go/magma"
	O "github.com/ibm/fp-go/option"
	G "github.com/ibm/fp-go/record/generic"
	T "github.com/ibm/fp-go/tuple"
)

func IsEmpty[K comparable, V any](r map[K]V) bool {
	return G.IsEmpty(r)
}

func IsNonEmpty[K comparable, V any](r map[K]V) bool {
	return G.IsNonEmpty(r)
}

func Keys[K comparable, V any](r map[K]V) []K {
	return G.Keys[map[K]V, []K](r)
}

func Values[K comparable, V any](r map[K]V) []V {
	return G.Values[map[K]V, []V](r)
}

func Collect[K comparable, V, R any](f func(K, V) R) func(map[K]V) []R {
	return G.Collect[map[K]V, []R](f)
}

func Reduce[K comparable, V, R any](f func(R, V) R, initial R) func(map[K]V) R {
	return G.Reduce[map[K]V](f, initial)
}

func ReduceWithIndex[K comparable, V, R any](f func(K, R, V) R, initial R) func(map[K]V) R {
	return G.ReduceWithIndex[map[K]V](f, initial)
}

func ReduceRef[K comparable, V, R any](f func(R, *V) R, initial R) func(map[K]V) R {
	return G.ReduceRef[map[K]V](f, initial)
}

func ReduceRefWithIndex[K comparable, V, R any](f func(K, R, *V) R, initial R) func(map[K]V) R {
	return G.ReduceRefWithIndex[map[K]V](f, initial)
}

func MonadMap[K comparable, V, R any](r map[K]V, f func(V) R) map[K]R {
	return G.MonadMap[map[K]V, map[K]R](r, f)
}

func MonadMapWithIndex[K comparable, V, R any](r map[K]V, f func(K, V) R) map[K]R {
	return G.MonadMapWithIndex[map[K]V, map[K]R](r, f)
}

func MonadMapRefWithIndex[K comparable, V, R any](r map[K]V, f func(K, *V) R) map[K]R {
	return G.MonadMapRefWithIndex[map[K]V, map[K]R](r, f)
}

func MonadMapRef[K comparable, V, R any](r map[K]V, f func(*V) R) map[K]R {
	return G.MonadMapRef[map[K]V, map[K]R](r, f)
}

func Map[K comparable, V, R any](f func(V) R) func(map[K]V) map[K]R {
	return G.Map[map[K]V, map[K]R](f)
}

func MapRef[K comparable, V, R any](f func(*V) R) func(map[K]V) map[K]R {
	return G.MapRef[map[K]V, map[K]R](f)
}

func MapWithIndex[K comparable, V, R any](f func(K, V) R) func(map[K]V) map[K]R {
	return G.MapWithIndex[map[K]V, map[K]R](f)
}

func MapRefWithIndex[K comparable, V, R any](f func(K, *V) R) func(map[K]V) map[K]R {
	return G.MapRefWithIndex[map[K]V, map[K]R](f)
}

func Lookup[K comparable, V any](k K) func(map[K]V) O.Option[V] {
	return G.Lookup[map[K]V](k)
}

func Has[K comparable, V any](k K, r map[K]V) bool {
	return G.Has(k, r)
}

func Union[K comparable, V any](m Mg.Magma[V]) func(map[K]V) func(map[K]V) map[K]V {
	return G.Union[map[K]V](m)
}

func Empty[K comparable, V any]() map[K]V {
	return G.Empty[map[K]V]()
}

func Size[K comparable, V any](r map[K]V) int {
	return G.Size(r)
}

func ToArray[K comparable, V any](r map[K]V) []T.Tuple2[K, V] {
	return G.ToArray[map[K]V, []T.Tuple2[K, V]](r)
}

func ToEntries[K comparable, V any](r map[K]V) []T.Tuple2[K, V] {
	return G.ToEntries[map[K]V, []T.Tuple2[K, V]](r)
}

func FromEntries[K comparable, V any](fa []T.Tuple2[K, V]) map[K]V {
	return G.FromEntries[map[K]V](fa)
}

func UpsertAt[K comparable, V any](k K, v V) func(map[K]V) map[K]V {
	return G.UpsertAt[map[K]V](k, v)
}

func DeleteAt[K comparable, V any](k K) func(map[K]V) map[K]V {
	return G.DeleteAt[map[K]V](k)
}

func Singleton[K comparable, V any](k K, v V) map[K]V {
	return G.Singleton[map[K]V](k, v)
}

// FilterMapWithIndex creates a new map with only the elements for which the transformation function creates a Some
func FilterMapWithIndex[K comparable, V1, V2 any](f func(K, V1) O.Option[V2]) func(map[K]V1) map[K]V2 {
	return G.FilterMapWithIndex[map[K]V1, map[K]V2](f)
}

// FilterMap creates a new map with only the elements for which the transformation function creates a Some
func FilterMap[K comparable, V1, V2 any](f func(V1) O.Option[V2]) func(map[K]V1) map[K]V2 {
	return G.FilterMap[map[K]V1, map[K]V2](f)
}

// Filter creates a new map with only the elements that match the predicate
func Filter[K comparable, V any](f func(K) bool) func(map[K]V) map[K]V {
	return G.Filter[map[K]V](f)
}

// FilterWithIndex creates a new map with only the elements that match the predicate
func FilterWithIndex[K comparable, V any](f func(K, V) bool) func(map[K]V) map[K]V {
	return G.FilterWithIndex[map[K]V](f)
}

// IsNil checks if the map is set to nil
func IsNil[K comparable, V any](m map[K]V) bool {
	return G.IsNil(m)
}

// IsNonNil checks if the map is set to nil
func IsNonNil[K comparable, V any](m map[K]V) bool {
	return G.IsNonNil(m)
}

// ConstNil return a nil map
func ConstNil[K comparable, V any]() map[K]V {
	return (map[K]V)(nil)
}
