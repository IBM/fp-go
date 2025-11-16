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

package reader

import (
	"github.com/IBM/fp-go/v2/function"
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// MonadTraverseArray transforms each element of an array using a function that returns a Reader,
// then collects the results into a single Reader containing an array.
// This is the monadic version that takes the array as the first parameter.
//
// All Readers share the same environment R and are evaluated with it.
//
// Example:
//
//	type Config struct { Prefix string }
//	numbers := []int{1, 2, 3}
//	addPrefix := func(n int) reader.Reader[Config, string] {
//	    return reader.Asks(func(c Config) string {
//	        return fmt.Sprintf("%s%d", c.Prefix, n)
//	    })
//	}
//	r := reader.MonadTraverseArray(numbers, addPrefix)
//	result := r(Config{Prefix: "num"}) // ["num1", "num2", "num3"]
func MonadTraverseArray[R, A, B any](ma []A, f Kleisli[R, A, B]) Reader[R, []B] {
	return G.MonadTraverseArray[Reader[R, B], Reader[R, []B]](ma, f)
}

// TraverseArray transforms each element of an array using a function that returns a Reader,
// then collects the results into a single Reader containing an array.
//
// This is useful for performing a Reader computation on each element of an array
// where all computations share the same environment.
//
// Example:
//
//	type Config struct { Multiplier int }
//	multiply := func(n int) reader.Reader[Config, int] {
//	    return reader.Asks(func(c Config) int { return n * c.Multiplier })
//	}
//	transform := reader.TraverseArray(multiply)
//	r := transform([]int{1, 2, 3})
//	result := r(Config{Multiplier: 10}) // [10, 20, 30]
func TraverseArray[R, A, B any](f Kleisli[R, A, B]) func([]A) Reader[R, []B] {
	return G.TraverseArray[Reader[R, B], Reader[R, []B], []A](f)
}

// TraverseArrayWithIndex transforms each element of an array using a function that takes
// both the index and the element, returning a Reader. The results are collected into
// a single Reader containing an array.
//
// This is useful when the transformation needs to know the position of each element.
//
// Example:
//
//	type Config struct { Prefix string }
//	addIndexPrefix := func(i int, s string) reader.Reader[Config, string] {
//	    return reader.Asks(func(c Config) string {
//	        return fmt.Sprintf("%s[%d]:%s", c.Prefix, i, s)
//	    })
//	}
//	transform := reader.TraverseArrayWithIndex(addIndexPrefix)
//	r := transform([]string{"a", "b", "c"})
//	result := r(Config{Prefix: "item"}) // ["item[0]:a", "item[1]:b", "item[2]:c"]
func TraverseArrayWithIndex[R, A, B any](f func(int, A) Reader[R, B]) func([]A) Reader[R, []B] {
	return G.TraverseArrayWithIndex[Reader[R, B], Reader[R, []B], []A](f)
}

// SequenceArray converts an array of Readers into a single Reader containing an array.
// All Readers in the input array share the same environment and are evaluated with it.
//
// This is useful when you have multiple independent Reader computations and want to
// collect all their results.
//
// Example:
//
//	type Config struct { X, Y, Z int }
//	readers := []reader.Reader[Config, int]{
//	    reader.Asks(func(c Config) int { return c.X }),
//	    reader.Asks(func(c Config) int { return c.Y }),
//	    reader.Asks(func(c Config) int { return c.Z }),
//	}
//	r := reader.SequenceArray(readers)
//	result := r(Config{X: 1, Y: 2, Z: 3}) // [1, 2, 3]
func SequenceArray[R, A any](ma []Reader[R, A]) Reader[R, []A] {
	return MonadTraverseArray(ma, function.Identity[Reader[R, A]])
}
