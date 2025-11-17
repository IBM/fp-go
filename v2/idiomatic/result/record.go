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

// TraverseRecordG transforms a map by applying a function that returns an Either to each value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the map of all Right values.
// The G suffix indicates support for generic map types.
//
// Example:
//
//	parse := func(s string) either.Either[error, int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseRecordG[map[string]string, map[string]int](parse)(map[string]string{"a": "1", "b": "2"})
//	// result is Right(map[string]int{"a": 1, "b": 2})
//
//go:inline
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f Kleisli[A, B]) Kleisli[GA, GB] {
	return func(ga GA) (GB, error) {
		bs := make(GB)
		for k, a := range ga {
			b, err := f(a)
			if err != nil {
				return Left[GB](err)
			}
			bs[k] = b
		}
		return Of(bs)
	}
}

// TraverseRecord transforms a map by applying a function that returns an Either to each value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// Otherwise, returns Right containing the map of all Right values.
//
// Example:
//
//	parse := func(s string) either.Either[error, int] {
//	    v, err := strconv.Atoi(s)
//	    return either.FromError(v, err)
//	}
//	result := either.TraverseRecord[string](parse)(map[string]string{"a": "1", "b": "2"})
//	// result is Right(map[string]int{"a": 1, "b": 2})
//
//go:inline
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

// TraverseRecordWithIndexG transforms a map by applying an indexed function that returns an Either.
// The function receives both the key and the value.
// If any value produces a Left, the entire result is that Left (short-circuits).
// The G suffix indicates support for generic map types.
//
// Example:
//
//	validate := func(k string, v string) either.Either[error, string] {
//	    if len(v) > 0 {
//	        return either.Right[error](k + ":" + v)
//	    }
//	    return either.Left[string](fmt.Errorf("empty value for key %s", k))
//	}
//	result := either.TraverseRecordWithIndexG[map[string]string, map[string]string](validate)(map[string]string{"a": "1"})
//	// result is Right(map[string]string{"a": "a:1"})
//
//go:inline
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) (B, error)) Kleisli[GA, GB] {
	return func(ga GA) (GB, error) {
		bs := make(GB)
		for k, a := range ga {
			b, err := f(k, a)
			if err != nil {
				return Left[GB](err)
			}
			bs[k] = b
		}
		return Of(bs)
	}
}

// TraverseRecordWithIndex transforms a map by applying an indexed function that returns an Either.
// The function receives both the key and the value.
// If any value produces a Left, the entire result is that Left (short-circuits).
//
// Example:
//
//	validate := func(k string, v string) either.Either[error, string] {
//	    if len(v) > 0 {
//	        return either.Right[error](k + ":" + v)
//	    }
//	    return either.Left[string](fmt.Errorf("empty value for key %s", k))
//	}
//	result := either.TraverseRecordWithIndex[string](validate)(map[string]string{"a": "1"})
//	// result is Right(map[string]string{"a": "a:1"})
//
//go:inline
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) (B, error)) Kleisli[map[K]A, map[K]B] {
	return TraverseRecordWithIndexG[map[K]A, map[K]B](f)
}
