package generic

import (
	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/array"
	O "github.com/ibm/fp-go/option"
	"github.com/ibm/fp-go/tuple"
)

// Of constructs a single element array
func Of[GA ~[]A, A any](value A) GA {
	return GA{value}
}

// From constructs an array from a set of variadic arguments
func From[GA ~[]A, A any](data ...A) GA {
	return data
}

func Lookup[GA ~[]A, A any](idx int) func(GA) O.Option[A] {
	none := O.None[A]()
	if idx < 0 {
		return F.Constant1[GA](none)
	}
	return func(as GA) O.Option[A] {
		if idx < len(as) {
			return O.Some(as[idx])
		}
		return none
	}
}

func Tail[GA ~[]A, A any](as GA) O.Option[GA] {
	if array.IsEmpty(as) {
		return O.None[GA]()
	}
	return O.Some(as[1:])
}

func Head[GA ~[]A, A any](as GA) O.Option[A] {
	if array.IsEmpty(as) {
		return O.None[A]()
	}
	return O.Some(as[0])
}

func First[GA ~[]A, A any](as GA) O.Option[A] {
	return Head(as)
}

func Last[GA ~[]A, A any](as GA) O.Option[A] {
	if array.IsEmpty(as) {
		return O.None[A]()
	}
	return O.Some(as[len(as)-1])
}

func Append[GA ~[]A, A any](as GA, a A) GA {
	return array.Append(as, a)
}

func Empty[GA ~[]A, A any]() GA {
	return array.Empty[GA]()
}

func UpsertAt[GA ~[]A, A any](a A) func(GA) GA {
	return array.UpsertAt[GA](a)
}

func MonadMap[GA ~[]A, GB ~[]B, A, B any](as GA, f func(a A) B) GB {
	return array.MonadMap[GA, GB](as, f)
}

func Size[GA ~[]A, A any](as GA) int {
	return len(as)
}

func filterMap[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(a A) O.Option[B]) GB {
	return array.Reduce(fa, func(bs GB, a A) GB {
		return O.MonadFold(f(a), F.Constant(bs), F.Bind1st(Append[GB, B], bs))
	}, Empty[GB]())
}

func MonadFilterMap[GA ~[]A, GB ~[]B, A, B any](fa GA, f func(a A) O.Option[B]) GB {
	return filterMap[GA, GB](fa, f)
}

func FilterMap[GA ~[]A, GB ~[]B, A, B any](f func(a A) O.Option[B]) func(GA) GB {
	return F.Bind2nd(MonadFilterMap[GA, GB, A, B], f)
}

func MonadPartition[GA ~[]A, A any](as GA, pred func(A) bool) tuple.Tuple2[GA, GA] {
	left := Empty[GA]()
	right := Empty[GA]()
	array.Reduce(as, func(c bool, a A) bool {
		if pred(a) {
			right = append(right, a)
		} else {
			left = append(left, a)
		}
		return c
	}, true)
	// returns the partition
	return tuple.MakeTuple2(left, right)
}

func Partition[GA ~[]A, A any](pred func(A) bool) func(GA) tuple.Tuple2[GA, GA] {
	return F.Bind2nd(MonadPartition[GA, A], pred)
}
