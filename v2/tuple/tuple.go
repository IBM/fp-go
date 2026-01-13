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

// Package tuple contains type definitions and functions for data structures for tuples of heterogenous types. For homogeneous types
// consider to use arrays for simplicity
package tuple

import (
	"encoding/json"
	"fmt"
	"strings"

	N "github.com/IBM/fp-go/v2/number"
)

// Of creates a [Tuple1] from a single value.
// This is a convenience function equivalent to [MakeTuple1].
//
// Example:
//
//	t := tuple.Of(42) // Creates Tuple1[int]{F1: 42}
func Of[T1 any](t T1) Tuple1[T1] {
	return MakeTuple1(t)
}

// First returns the first element of a [Tuple2].
// This is a convenience accessor for the F1 field.
//
// Example:
//
//	t := tuple.MakeTuple2("hello", 42)
//	s := tuple.First(t) // Returns "hello"
func First[T1, T2 any](t Tuple2[T1, T2]) T1 {
	return t.F1
}

// Second returns the second element of a [Tuple2].
// This is a convenience accessor for the F2 field.
//
// Example:
//
//	t := tuple.MakeTuple2("hello", 42)
//	n := tuple.Second(t) // Returns 42
func Second[T1, T2 any](t Tuple2[T1, T2]) T2 {
	return t.F2
}

// Swap exchanges the positions of the two elements in a [Tuple2].
// The first element becomes the second, and the second becomes the first.
//
// Example:
//
//	t := tuple.MakeTuple2("hello", 42)
//	swapped := tuple.Swap(t) // Returns Tuple2[int, string]{F1: 42, F2: "hello"}
func Swap[T1, T2 any](t Tuple2[T1, T2]) Tuple2[T2, T1] {
	return MakeTuple2(t.F2, t.F1)
}

// Of2 creates a curried function that pairs a value with a constant second element.
// It returns a function that takes a value of type T1 and creates a [Tuple2] with
// the provided constant value e as the second element.
//
// This is useful for partial application and functional composition.
//
// Example:
//
//	pairWith42 := tuple.Of2[string](42)
//	t := pairWith42("hello") // Returns Tuple2[string, int]{F1: "hello", F2: 42}
func Of2[T1, T2 any](e T2) func(T1) Tuple2[T1, T2] {
	return func(t T1) Tuple2[T1, T2] {
		return MakeTuple2(t, e)
	}
}

// BiMap applies two mapping functions to both elements of a [Tuple2].
// The first function (mapSnd) is applied to the second element,
// and the second function (mapFst) is applied to the first element.
//
// This is a bifunctor map operation that allows independent transformation
// of both tuple elements.
//
// Example:
//
//	t := tuple.MakeTuple2(5, "hello")
//	mapper := tuple.BiMap(
//	    S.Size,
//	    func(n int) string { return fmt.Sprintf("%d", n*2) },
//	)
//	result := mapper(t) // Returns Tuple2[string, int]{F1: "10", F2: 5}
func BiMap[E, G, A, B any](mapSnd func(E) G, mapFst func(A) B) func(Tuple2[A, E]) Tuple2[B, G] {
	return func(t Tuple2[A, E]) Tuple2[B, G] {
		return MakeTuple2(mapFst(First(t)), mapSnd(Second(t)))
	}
}

// marshalJSON marshals the tuple into a JSON array
func tupleMarshalJSON(src ...any) ([]byte, error) {
	return json.Marshal(src)
}

// tupleUnmarshalJSON unmarshals a JSON array into a tuple
func tupleUnmarshalJSON(data []byte, dst ...any) error {
	var src []json.RawMessage
	if err := json.Unmarshal(data, &src); err != nil {
		return err
	}
	l := N.Min(len(src), len(dst))
	// unmarshal
	for i := 0; i < l; i++ {
		if err := json.Unmarshal(src[i], dst[i]); err != nil {
			return err
		}
	}
	// successfully decoded the tuple
	return nil
}

// tupleString converts a tuple to a string
func tupleString(src ...any) string {
	l := len(src)
	return fmt.Sprintf("Tuple%d[%s](%s)", l, fmt.Sprintf(strings.Repeat(", %T", l)[2:], src...), fmt.Sprintf(strings.Repeat(", %v", l)[2:], src...))
}
