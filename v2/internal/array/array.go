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

package array

import "slices"

// Of wraps a single value in a one-element array. The result is never nil.
func Of[GA ~[]A, A any](a A) GA {
	return GA{a}
}

// Slice returns a sub-slice of as in the half-open range [low, high).
// Negative indices count backward from the end of the array.
// Returns nil (representing an empty array) when the resulting range is empty.
func Slice[GA ~[]A, A any](low, high int) func(as GA) GA {
	return func(as GA) GA {
		length := len(as)

		// Handle negative indices - count backward from the end
		if low < 0 {
			low = max(length+low, 0)
		}
		if high < 0 {
			high = max(length+high, 0)
		}

		if low > length {
			return nil
		}

		// End index > array length: slice to the end
		if high > length {
			high = length
		}

		// Start >= end: return empty array
		if low >= high {
			return nil
		}

		return as[low:high]
	}
}

// SliceRight returns the suffix of as starting at start.
// Negative indices count backward from the end of the array.
// Returns nil (representing an empty array) when start is at or beyond the end.
func SliceRight[GA ~[]A, A any](start int) func(as GA) GA {
	return func(as GA) GA {
		length := len(as)

		// Handle negative indices - count backward from the end
		if start < 0 {
			start = max(length+start, 0)
		}

		// Start index >= array length: return empty array
		if start >= length {
			return nil
		}

		return as[start:]
	}
}

func IsEmpty[GA ~[]A, A any](as GA) bool {
	return len(as) == 0
}

// EmptyToNil returns nil if as is empty and as otherwise, normalising every
// empty array to nil, the canonical representation of an empty array.
func EmptyToNil[GA ~[]A, A any](as GA) GA {
	if len(as) == 0 {
		return nil
	}
	return as
}

func IsNil[GA ~[]A, A any](as GA) bool {
	return as == nil
}

func IsNonNil[GA ~[]A, A any](as GA) bool {
	return as != nil
}

func Reduce[GA ~[]A, A, B any](fa GA, f func(B, A) B, initial B) B {
	current := initial
	for i := range len(fa) {
		current = f(current, fa[i])
	}
	return current
}

func ReduceWithIndex[GA ~[]A, A, B any](fa GA, f func(int, B, A) B, initial B) B {
	current := initial
	for i := range len(fa) {
		current = f(i, current, fa[i])
	}
	return current
}

func ReduceRight[GA ~[]A, A, B any](fa GA, f func(A, B) B, initial B) B {
	current := initial
	count := len(fa)
	for i := count - 1; i >= 0; i-- {
		current = f(fa[i], current)
	}
	return current
}

func ReduceRightWithIndex[GA ~[]A, A, B any](fa GA, f func(int, A, B) B, initial B) B {
	current := initial
	count := len(fa)
	for i := count - 1; i >= 0; i-- {
		current = f(i, fa[i], current)
	}
	return current
}

// Append appends a single element to as and returns the result.
// The result is never nil.
func Append[GA ~[]A, A any](as GA, a A) GA {
	return append(as, a)
}

// Push appends a single element to a copy of as and returns the copy.
// The result is never nil.
func Push[GA ~[]A, A any](as GA, a A) GA {
	l := len(as)
	cpy := make(GA, l+1)
	copy(cpy, as)
	cpy[l] = a
	return cpy
}

// Empty returns nil, the canonical representation of an empty array.
func Empty[GA ~[]A, A any]() GA {
	return nil
}

func upsertAt[GA ~[]A, A any](fa GA, a A) GA {
	buf := make(GA, len(fa)+1)
	buf[copy(buf, fa)] = a
	return buf
}

// UpsertAt appends a at the end of a copy of ma. The result is never nil.
func UpsertAt[GA ~[]A, A any](a A) func(GA) GA {
	return func(ma GA) GA {
		return upsertAt(ma, a)
	}
}

// MonadMap applies f to every element of as and returns the results.
// Returns nil (representing an empty array) when as is empty.
func MonadMap[GA ~[]A, GB ~[]B, A, B any](as GA, f func(a A) B) GB {
	count := len(as)
	// shortcut for empty list
	if count == 0 {
		return nil
	}
	// traverse the map
	bs := make(GB, count)
	for i := range count {
		bs[i] = f(as[i])
	}
	return bs
}

// Map returns a curried form of MonadMap.
// The returned function returns nil (representing an empty array) when its input is empty.
func Map[GA ~[]A, GB ~[]B, A, B any](f func(a A) B) func(GA) GB {
	return func(as GA) GB {
		return MonadMap[GA, GB](as, f)
	}
}

// MonadMapWithIndex applies f (with element index) to every element of as and returns the results.
// Returns nil (representing an empty array) when as is empty.
func MonadMapWithIndex[GA ~[]A, GB ~[]B, A, B any](as GA, f func(idx int, a A) B) GB {
	count := len(as)
	// shortcut for empty lists
	if count == 0 {
		return nil
	}
	// proceed with the mapping
	bs := make(GB, count)
	for i := range count {
		bs[i] = f(i, as[i])
	}
	return bs
}

// ConstNil always returns nil, the canonical representation of an empty array.
func ConstNil[GA ~[]A, A any]() GA {
	return GA(nil)
}

// Concat concatenates left and right into a new array.
// When either side is empty it is returned directly, so the result may be nil
// if that side was nil.
func Concat[GT ~[]T, T any](left, right GT) GT {
	// some performance checks
	ll := len(left)
	if ll == 0 {
		return right
	}
	lr := len(right)
	if lr == 0 {
		return left
	}
	// need to copy
	buf := make(GT, ll+lr)
	copy(buf[copy(buf, left):], right)
	return buf
}

// Reverse returns a new array with elements in reversed order.
// When as has zero or one elements it is returned as-is, so the result may be nil
// if as was nil.
func Reverse[GT ~[]T, T any](as GT) GT {
	l := len(as)
	if l <= 1 {
		return as
	}
	ras := make(GT, l)
	l1 := l - 1
	for i := range l {
		ras[i] = as[l1-i]
	}
	return ras
}

// UnsafeUpdateAt returns a copy of as with the element at index i replaced by v.
// The result is never nil.
func UnsafeUpdateAt[GT ~[]T, T any](as GT, i int, v T) GT {
	c := slices.Clone(as)
	c[i] = v
	return c
}

// MakeBy returns an array of length n with element i initialised by f(i).
// Returns nil (representing an empty array) when n <= 0.
func MakeBy[AS ~[]A, F ~func(int) A, A any](n int, f F) AS {
	// sanity check
	if n <= 0 {
		return nil
	}
	// run the generator function across the input
	as := make(AS, n)
	for i := range n {
		as[i] = f(i)
	}
	return as
}
