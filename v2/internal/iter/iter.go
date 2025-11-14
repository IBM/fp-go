package iter

import (
	F "github.com/IBM/fp-go/v2/function"
)

func MonadReduceWithIndex[GA ~func(yield func(A) bool), A, B any](fa GA, f func(int, B, A) B, initial B) B {
	current := initial
	var i int
	for a := range fa {
		current = f(i, current, a)
		i += 1
	}
	return current
}

func MonadReduce[GA ~func(yield func(A) bool), A, B any](fa GA, f func(B, A) B, initial B) B {
	current := initial
	for a := range fa {
		current = f(current, a)
	}
	return current
}

// Concat concatenates two sequences, yielding all elements from left followed by all elements from right.
func Concat[GT ~func(yield func(T) bool), T any](left, right GT) GT {
	return func(yield func(T) bool) {
		for t := range left {
			if !yield(t) {
				return
			}
		}
		for t := range right {
			if !yield(t) {
				return
			}
		}
	}
}

func Of[GA ~func(yield func(A) bool), A any](a A) GA {
	return func(yield func(A) bool) {
		yield(a)
	}
}

func MonadAppend[GA ~func(yield func(A) bool), A any](f GA, tail A) GA {
	return Concat(f, Of[GA](tail))
}

func Append[GA ~func(yield func(A) bool), A any](tail A) func(GA) GA {
	return F.Bind2nd(Concat[GA], Of[GA](tail))
}

func Prepend[GA ~func(yield func(A) bool), A any](head A) func(GA) GA {
	return F.Bind1st(Concat[GA], Of[GA](head))
}

func Empty[GA ~func(yield func(A) bool), A any]() GA {
	return func(_ func(A) bool) {}
}
