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

package record

import (
	G "github.com/IBM/fp-go/v2/internal/record"
)

// TraverseWithIndex transforms a map of values into a value of a map by applying an effectful function
// to each key-value pair. The function has access to both the key and value.
//
// This is useful when you need to perform an operation that may fail or have side effects on each
// element of a map, and you want to collect the results in the same applicative context.
//
// Type parameters:
//   - K: The key type (must be comparable)
//   - A: The input value type
//   - B: The output value type
//   - HKTB: Higher-kinded type representing the effect containing B (e.g., Option[B], Either[E, B])
//   - HKTAB: Higher-kinded type representing a function from B to map[K]B in the effect
//   - HKTRB: Higher-kinded type representing the effect containing map[K]B
//
// Parameters:
//   - fof: Lifts a pure map[K]B into the effect (the "of" or "pure" function)
//   - fmap: Maps a function over the effect (the "map" or "fmap" function)
//   - fap: Applies an effectful function to an effectful value (the "ap" function)
//   - f: The transformation function that takes a key and value and returns an effect
//
// Example with Option:
//
//	f := func(k string, n int) O.Option[int] {
//	    if n > 0 {
//	        return O.Some(n * 2)
//	    }
//	    return O.None[int]()
//	}
//	traverse := TraverseWithIndex(O.Of[map[string]int], O.Map[...], O.Ap[...], f)
//	result := traverse(map[string]int{"a": 1, "b": 2}) // O.Some(map[string]int{"a": 2, "b": 4})
func TraverseWithIndex[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(K, A) HKTB) func(map[K]A) HKTRB {
	return G.TraverseWithIndex[map[K]A](fof, fmap, fap, f)
}

// Traverse transforms a map of values into a value of a map by applying an effectful function
// to each value. Unlike TraverseWithIndex, this function does not provide access to the keys.
//
// This is useful when you need to perform an operation that may fail or have side effects on each
// element of a map, and you want to collect the results in the same applicative context.
//
// Type parameters:
//   - K: The key type (must be comparable)
//   - A: The input value type
//   - B: The output value type
//   - HKTB: Higher-kinded type representing the effect containing B (e.g., Option[B], Either[E, B])
//   - HKTAB: Higher-kinded type representing a function from B to map[K]B in the effect
//   - HKTRB: Higher-kinded type representing the effect containing map[K]B
//
// Parameters:
//   - fof: Lifts a pure map[K]B into the effect (the "of" or "pure" function)
//   - fmap: Maps a function over the effect (the "map" or "fmap" function)
//   - fap: Applies an effectful function to an effectful value (the "ap" function)
//   - f: The transformation function that takes a value and returns an effect
//
// Example with Option:
//
//	f := func(s string) O.Option[string] {
//	    if s != "" {
//	        return O.Some(strings.ToUpper(s))
//	    }
//	    return O.None[string]()
//	}
//	traverse := Traverse(O.Of[map[string]string], O.Map[...], O.Ap[...], f)
//	result := traverse(map[string]string{"a": "hello"}) // O.Some(map[string]string{"a": "HELLO"})
func Traverse[K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(map[K]B) HKTRB,
	fmap func(func(map[K]B) func(B) map[K]B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
	f func(A) HKTB) func(map[K]A) HKTRB {
	return G.Traverse[map[K]A](fof, fmap, fap, f)
}

// Sequence transforms a map of effects into an effect of a map.
// This is the dual of Traverse where the transformation function is the identity.
//
// This is useful when you have a map where each value is already in an effect context
// (like Option, Either, etc.) and you want to "flip" the nesting to get a single effect
// containing a map of plain values.
//
// If any value in the map is a "failure" (e.g., None, Left), the entire result will be
// a failure. If all values are "successes", the result will be a success containing a map
// of all the unwrapped values.
//
// Type parameters:
//   - K: The key type (must be comparable)
//   - A: The value type inside the effect
//   - HKTA: Higher-kinded type representing the effect containing A (e.g., Option[A])
//   - HKTAA: Higher-kinded type representing a function from A to map[K]A in the effect
//   - HKTRA: Higher-kinded type representing the effect containing map[K]A
//
// Parameters:
//   - fof: Lifts a pure map[K]A into the effect (the "of" or "pure" function)
//   - fmap: Maps a function over the effect (the "map" or "fmap" function)
//   - fap: Applies an effectful function to an effectful value (the "ap" function)
//   - ma: The input map where each value is in an effect context
//
// Example with Option:
//
//	input := map[string]O.Option[int]{"a": O.Some(1), "b": O.Some(2)}
//	result := Sequence(O.Of[map[string]int], O.Map[...], O.Ap[...], input)
//	// result: O.Some(map[string]int{"a": 1, "b": 2})
//
//	input2 := map[string]O.Option[int]{"a": O.Some(1), "b": O.None[int]()}
//	result2 := Sequence(O.Of[map[string]int], O.Map[...], O.Ap[...], input2)
//	// result2: O.None[map[string]int]()
func Sequence[K comparable, A, HKTA, HKTAA, HKTRA any](
	fof func(map[K]A) HKTRA,
	fmap func(func(map[K]A) func(A) map[K]A) func(HKTRA) HKTAA,
	fap func(HKTA) func(HKTAA) HKTRA,
	ma map[K]HKTA) HKTRA {
	return G.Sequence(fof, fmap, fap, ma)

}
