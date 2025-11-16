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

import (
	I "github.com/IBM/fp-go/v2/iterator/iter"
)

// TraverseIter transforms a sequence by applying a function that returns an Option to each element.
// Returns Some containing a sequence of results if all operations succeed, None if any fails.
// This function is useful for processing sequences where each element may fail validation or transformation.
//
// The traversal short-circuits on the first None encountered, making it efficient for validation pipelines.
// The resulting sequence is lazy and will only be evaluated when iterated.
//
// Example:
//
//	// Parse a sequence of strings to integers
//	parse := func(s string) Option[int] {
//	    n, err := strconv.Atoi(s)
//	    if err != nil { return None[int]() }
//	    return Some(n)
//	}
//
//	// Create a sequence of strings
//	strings := func(yield func(string) bool) {
//	    for _, s := range []string{"1", "2", "3"} {
//	        if !yield(s) { return }
//	    }
//	}
//
//	result := TraverseIter(parse)(strings)
//	// result is Some(sequence of [1, 2, 3])
//
//	// With invalid input
//	invalidStrings := func(yield func(string) bool) {
//	    for _, s := range []string{"1", "invalid", "3"} {
//	        if !yield(s) { return }
//	    }
//	}
//
//	result := TraverseIter(parse)(invalidStrings)
//	// result is None because "invalid" cannot be parsed
func TraverseIter[A, B any](f Kleisli[A, B]) Kleisli[Seq[A], Seq[B]] {
	return func(s Seq[A]) (Seq[B], bool) {
		var bs []B
		for a := range s {
			b, bok := f(a)
			if !bok {
				return nil, false
			}
			bs = append(bs, b)
		}
		return I.From(bs...), true
	}
}
