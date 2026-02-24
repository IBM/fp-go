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

package effect

import "github.com/IBM/fp-go/v2/context/readerreaderioresult"

// TraverseArray applies an effectful function to each element of an array,
// collecting the results into a new array. If any effect fails, the entire
// traversal fails and returns the first error encountered.
//
// This is useful for performing effectful operations on collections while
// maintaining the sequential order of results.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - A: The input element type
//   - B: The output element type
//
// # Parameters
//
//   - f: An effectful function to apply to each element
//
// # Returns
//
//   - Kleisli[C, []A, []B]: A function that transforms an array of A to an effect producing an array of B
//
// # Example
//
//	parseIntEff := func(s string) Effect[MyContext, int] {
//		val, err := strconv.Atoi(s)
//		if err != nil {
//			return effect.Fail[MyContext, int](err)
//		}
//		return effect.Of[MyContext](val)
//	}
//	input := []string{"1", "2", "3"}
//	eff := effect.TraverseArray[MyContext](parseIntEff)(input)
//	// eff produces []int{1, 2, 3}
func TraverseArray[C, A, B any](f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return readerreaderioresult.TraverseArray(f)
}
