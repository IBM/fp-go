// Copyright (c) 2025 IBM Corp.
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

package option

// TraverseArrayG transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	parse := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil { return None[int]() }
//	    return Some(n)
//	}
//	result := TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
//	result := TraverseArrayG[[]string, []int](parse)([]string{"1", "x", "3"}) // None
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return func(g GA) (GB, bool) {
		bs := make(GB, len(g))
		for i, a := range g {
			b, bok := f(a)
			if !bok {
				return bs, false
			}
			bs[i] = b
		}
		return bs, true
	}
}

// TraverseArray transforms an array by applying a function that returns an Option to each element.
// Returns Some containing the array of results if all operations succeed, None if any fails.
//
// Example:
//
//	validate := func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	}
//	result := TraverseArray(validate)([]int{1, 2, 3}) // Some([2, 4, 6])
//	result := TraverseArray(validate)([]int{1, -1, 3}) // None
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return TraverseArrayG[[]A, []B](f)
}

// TraverseArrayWithIndexG transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
// This is the generic version that works with custom slice types.
//
// Example:
//
//	f := func(i int, s string) Option[string] {
//	    return Some(fmt.Sprintf("%d:%s", i, s))
//	}
//	result := TraverseArrayWithIndexG[[]string, []string](f)([]string{"a", "b"}) // Some(["0:a", "1:b"])
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) (B, bool)) Kleisli[GA, GB] {
	return func(g GA) (GB, bool) {
		bs := make(GB, len(g))
		for i, a := range g {
			b, bok := f(i, a)
			if !bok {
				return bs, false
			}
			bs[i] = b
		}
		return bs, true
	}
}

// TraverseArrayWithIndex transforms an array by applying an indexed function that returns an Option.
// The function receives both the index and the element.
//
// Example:
//
//	f := func(i int, x int) Option[int] {
//	    if x > i { return Some(x) }
//	    return None[int]()
//	}
//	result := TraverseArrayWithIndex(f)([]int{1, 2, 3}) // Some([1, 2, 3])
func TraverseArrayWithIndex[A, B any](f func(int, A) (B, bool)) Kleisli[[]A, []B] {
	return TraverseArrayWithIndexG[[]A, []B](f)
}
