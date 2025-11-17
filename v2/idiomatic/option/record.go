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

// TraverseRecordG transforms a record (map) by applying a function that returns an Option to each value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
// This is the generic version that works with custom map types.
//
// Example:
//
//	validate := func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordG[map[string]int, map[string]int](validate)(input) // Some(map[a:2 b:4])
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return func(ga GA) (GB, bool) {
		gb := make(GB)
		for k, a := range ga {
			if b, ok := f(a); ok {
				gb[k] = b
			} else {
				return gb, false
			}

		}
		return gb, true
	}
}

// TraverseRecord transforms a record (map) by applying a function that returns an Option to each value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
//
// Example:
//
//	validate := func(x int) Option[string] {
//	    if x > 0 { return Some(fmt.Sprintf("%d", x)) }
//	    return None[string]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecord(validate)(input) // Some(map[a:"1" b:"2"])
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

// TraverseRecordWithIndexG transforms a record by applying a function that receives both key and value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
// This is the generic version that works with custom map types.
//
// Example:
//
//	f := func(k string, v int) Option[string] {
//	    return Some(fmt.Sprintf("%s:%d", k, v))
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordWithIndexG[map[string]int, map[string]string](f)(input) // Some(map[a:"a:1" b:"b:2"])
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) (B, bool)) Kleisli[GA, GB] {
	return func(ga GA) (GB, bool) {
		gb := make(GB)
		for k, a := range ga {
			if b, ok := f(k, a); ok {
				gb[k] = b
			} else {
				return gb, false
			}

		}
		return gb, true
	}
}

// TraverseRecordWithIndex transforms a record by applying a function that receives both key and value.
// Returns Some containing the map of results if all operations succeed, None if any fails.
//
// Example:
//
//	f := func(k string, v int) Option[int] {
//	    if v > 0 { return Some(v) }
//	    return None[int]()
//	}
//	input := map[string]int{"a": 1, "b": 2}
//	result := TraverseRecordWithIndex(f)(input) // Some(map[a:1 b:2])
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) (B, bool)) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordWithIndexG[map[K]A, map[K]B](f)
}
