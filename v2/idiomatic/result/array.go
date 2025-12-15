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

package result

// TraverseArrayG transforms an array by applying a function that returns a Result (value, error) to each element.
// It processes elements from left to right, applying the function to each.
// If any element produces an error, the entire operation short-circuits and returns that error.
// Otherwise, it returns a successful result containing the array of all transformed values.
//
// The G suffix indicates support for generic slice types (e.g., custom slice types based on []T).
//
// Type Parameters:
//   - GA: Source slice type (must be based on []A)
//   - GB: Target slice type (must be based on []B)
//   - A: Source element type
//   - B: Target element type
//
// Parameters:
//   - f: A Kleisli arrow (A) -> (B, error) that transforms each element
//
// Returns:
//   - A Kleisli arrow (GA) -> (GB, error) that transforms the entire array
//
// Behavior:
//   - Short-circuits on the first error encountered
//   - Preserves the order of elements
//   - Returns an empty slice for empty input
//
// Example - Parse strings to integers:
//
//	parse := func(s string) (int, error) {
//	    return strconv.Atoi(s)
//	}
//	result := result.TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"})
//	// result is ([]int{1, 2, 3}, nil)
//
// Example - Short-circuit on error:
//
//	result := result.TraverseArrayG[[]string, []int](parse)([]string{"1", "bad", "3"})
//	// result is ([]int(nil), error) - stops at "bad"
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return func(ga GA) (GB, error) {
		bs := make(GB, len(ga))
		for i, a := range ga {
			b, err := f(a)
			if err != nil {
				return Left[GB](err)
			}
			bs[i] = b
		}
		return Of(bs)
	}
}

// TraverseArray transforms an array by applying a function that returns a Result (value, error) to each element.
// It processes elements from left to right, applying the function to each.
// If any element produces an error, the entire operation short-circuits and returns that error.
// Otherwise, it returns a successful result containing the array of all transformed values.
//
// This is a convenience wrapper around [TraverseArrayG] for standard slice types.
//
// Type Parameters:
//   - A: Source element type
//   - B: Target element type
//
// Parameters:
//   - f: A Kleisli arrow (A) -> (B, error) that transforms each element
//
// Returns:
//   - A Kleisli arrow ([]A) -> ([]B, error) that transforms the entire array
//
// Example - Validate and transform:
//
//	validate := func(s string) (int, error) {
//	    n, err := strconv.Atoi(s)
//	    if err != nil {
//	        return 0, err
//	    }
//	    if n < 0 {
//	        return 0, errors.New("negative number")
//	    }
//	    return n * 2, nil
//	}
//	result := result.TraverseArray(validate)([]string{"1", "2", "3"})
//	// result is ([]int{2, 4, 6}, nil)
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return TraverseArrayG[[]A, []B](f)
}

// TraverseArrayWithIndexG transforms an array by applying an indexed function that returns a Result (value, error).
// The function receives both the zero-based index and the element for each iteration.
// If any element produces an error, the entire operation short-circuits and returns that error.
// Otherwise, it returns a successful result containing the array of all transformed values.
//
// The G suffix indicates support for generic slice types (e.g., custom slice types based on []T).
//
// Type Parameters:
//   - GA: Source slice type (must be based on []A)
//   - GB: Target slice type (must be based on []B)
//   - A: Source element type
//   - B: Target element type
//
// Parameters:
//   - f: An indexed function (int, A) -> (B, error) that transforms each element
//
// Returns:
//   - A Kleisli arrow (GA) -> (GB, error) that transforms the entire array
//
// Behavior:
//   - Processes elements from left to right with their indices (0, 1, 2, ...)
//   - Short-circuits on the first error encountered
//   - Preserves the order of elements
//
// Example - Annotate with index:
//
//	annotate := func(i int, s string) (string, error) {
//	    if S.IsEmpty(s) {
//	        return "", fmt.Errorf("empty string at index %d", i)
//	    }
//	    return fmt.Sprintf("[%d]=%s", i, s), nil
//	}
//	result := result.TraverseArrayWithIndexG[[]string, []string](annotate)([]string{"a", "b", "c"})
//	// result is ([]string{"[0]=a", "[1]=b", "[2]=c"}, nil)
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) (B, error)) Kleisli[GA, GB] {
	return func(ga GA) (GB, error) {
		bs := make(GB, len(ga))
		for i, a := range ga {
			b, err := f(i, a)
			if err != nil {
				return Left[GB](err)
			}
			bs[i] = b
		}
		return Of(bs)
	}
}

// TraverseArrayWithIndex transforms an array by applying an indexed function that returns a Result (value, error).
// The function receives both the zero-based index and the element for each iteration.
// If any element produces an error, the entire operation short-circuits and returns that error.
// Otherwise, it returns a successful result containing the array of all transformed values.
//
// This is a convenience wrapper around [TraverseArrayWithIndexG] for standard slice types.
//
// Type Parameters:
//   - A: Source element type
//   - B: Target element type
//
// Parameters:
//   - f: An indexed function (int, A) -> (B, error) that transforms each element
//
// Returns:
//   - A Kleisli arrow ([]A) -> ([]B, error) that transforms the entire array
//
// Example - Validate with position info:
//
//	check := func(i int, s string) (string, error) {
//	    if S.IsEmpty(s) {
//	        return "", fmt.Errorf("empty value at position %d", i)
//	    }
//	    return strings.ToUpper(s), nil
//	}
//	result := result.TraverseArrayWithIndex(check)([]string{"a", "b", "c"})
//	// result is ([]string{"A", "B", "C"}, nil)
//
//go:inline
func TraverseArrayWithIndex[A, B any](f func(int, A) (B, error)) Kleisli[[]A, []B] {
	return TraverseArrayWithIndexG[[]A, []B](f)
}
