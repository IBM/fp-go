package array

import (
	"github.com/IBM/fp-go/internal/array"
	M "github.com/IBM/fp-go/monoid"
)

func concat[T any](left, right []T) []T {
	// some performance checks
	ll := len(left)
	lr := len(right)
	if ll == 0 {
		return right
	}
	if lr == 0 {
		return left
	}
	// need to copy
	buf := make([]T, ll+lr)
	copy(buf[copy(buf, left):], right)
	return buf
}

func Monoid[T any]() M.Monoid[[]T] {
	return M.MakeMonoid(concat[T], Empty[T]())
}

func addLen[A any](count int, data []A) int {
	return count + len(data)
}

// ConcatAll efficiently concatenates the input arrays into a final array
func ArrayConcatAll[A any](data ...[]A) []A {
	// get the full size
	count := array.Reduce(data, addLen[A], 0)
	buf := make([]A, count)
	// copy
	array.Reduce(data, func(idx int, seg []A) int {
		return idx + copy(buf[idx:], seg)
	}, 0)
	// returns the final array
	return buf
}
