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

package array

func Slice[GA ~[]A, A any](low, high int) func(as GA) GA {
	return func(as GA) GA {
		length := len(as)

		// Handle negative indices - count backward from the end
		if low < 0 {
			low = length + low
			if low < 0 {
				low = 0
			}
		}
		if high < 0 {
			high = length + high
			if high < 0 {
				high = 0
			}
		}

		// Start index > array length: return empty array
		if low > length {
			return Empty[GA, A]()
		}

		// End index > array length: slice to the end
		if high > length {
			high = length
		}

		// Start >= end: return empty array
		if low >= high {
			return Empty[GA, A]()
		}

		return as[low:high]
	}
}

func IsEmpty[GA ~[]A, A any](as GA) bool {
	return len(as) == 0
}

func IsNil[GA ~[]A, A any](as GA) bool {
	return as == nil
}

func IsNonNil[GA ~[]A, A any](as GA) bool {
	return as != nil
}

func Reduce[GA ~[]A, A, B any](fa GA, f func(B, A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
		current = f(current, fa[i])
	}
	return current
}

func ReduceWithIndex[GA ~[]A, A, B any](fa GA, f func(int, B, A) B, initial B) B {
	current := initial
	count := len(fa)
	for i := 0; i < count; i++ {
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

func Append[GA ~[]A, A any](as GA, a A) GA {
	return append(as, a)
}

func Push[GA ~[]A, A any](as GA, a A) GA {
	l := len(as)
	cpy := make(GA, l+1)
	copy(cpy, as)
	cpy[l] = a
	return cpy
}

func Empty[GA ~[]A, A any]() GA {
	return make(GA, 0)
}

func upsertAt[GA ~[]A, A any](fa GA, a A) GA {
	buf := make(GA, len(fa)+1)
	buf[copy(buf, fa)] = a
	return buf
}

func UpsertAt[GA ~[]A, A any](a A) func(GA) GA {
	return func(ma GA) GA {
		return upsertAt(ma, a)
	}
}

func MonadMap[GA ~[]A, GB ~[]B, A, B any](as GA, f func(a A) B) GB {
	count := len(as)
	bs := make(GB, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(as[i])
	}
	return bs
}

func Map[GA ~[]A, GB ~[]B, A, B any](f func(a A) B) func(GA) GB {
	return func(as GA) GB {
		return MonadMap[GA, GB](as, f)
	}
}

func MonadMapWithIndex[GA ~[]A, GB ~[]B, A, B any](as GA, f func(idx int, a A) B) GB {
	count := len(as)
	bs := make(GB, count)
	for i := count - 1; i >= 0; i-- {
		bs[i] = f(i, as[i])
	}
	return bs
}

func ConstNil[GA ~[]A, A any]() GA {
	return (GA)(nil)
}
