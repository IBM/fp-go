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

package assert

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

var (
	errTest = fmt.Errorf("test failure")

	// Eq is the equal predicate checking if objects are equal
	Eq = eq.FromEquals(assert.ObjectsAreEqual)
)

func wrap1[T any](wrapped func(t assert.TestingT, expected, actual any, msgAndArgs ...any) bool, t *testing.T, expected T) result.Kleisli[T, T] {
	return func(actual T) Result[T] {
		ok := wrapped(t, expected, actual)
		if ok {
			return result.Of(actual)
		}
		return result.Left[T](errTest)
	}
}

// NotEqual tests if the expected and the actual values are not equal
func NotEqual[T any](t *testing.T, expected T) result.Kleisli[T, T] {
	return wrap1(assert.NotEqual, t, expected)
}

// Equal tests if the expected and the actual values are equal
func Equal[T any](t *testing.T, expected T) result.Kleisli[T, T] {
	return wrap1(assert.Equal, t, expected)
}

// Length tests if an array has the expected length
func Length[T any](t *testing.T, expected int) result.Kleisli[[]T, []T] {
	return func(actual []T) Result[[]T] {
		ok := assert.Len(t, actual, expected)
		if ok {
			return result.Of(actual)
		}
		return result.Left[[]T](errTest)
	}
}

// NoError validates that there is no error
func NoError[T any](t *testing.T) result.Operator[T, T] {
	return func(actual Result[T]) Result[T] {
		return result.MonadFold(actual, func(e error) Result[T] {
			assert.NoError(t, e)
			return result.Left[T](e)
		}, func(value T) Result[T] {
			assert.NoError(t, nil)
			return result.Of(value)
		})
	}
}

// ArrayContains tests if a value is contained in an array
func ArrayContains[T any](t *testing.T, expected T) result.Kleisli[[]T, []T] {
	return func(actual []T) Result[[]T] {
		ok := assert.Contains(t, actual, expected)
		if ok {
			return result.Of(actual)
		}
		return result.Left[[]T](errTest)
	}
}

// ContainsKey tests if a key is contained in a map
func ContainsKey[T any, K comparable](t *testing.T, expected K) result.Kleisli[map[K]T, map[K]T] {
	return func(actual map[K]T) Result[map[K]T] {
		ok := assert.Contains(t, actual, expected)
		if ok {
			return result.Of(actual)
		}
		return result.Left[map[K]T](errTest)
	}
}

// NotContainsKey tests if a key is not contained in a map
func NotContainsKey[T any, K comparable](t *testing.T, expected K) result.Kleisli[map[K]T, map[K]T] {
	return func(actual map[K]T) Result[map[K]T] {
		ok := assert.NotContains(t, actual, expected)
		if ok {
			return result.Of(actual)
		}
		return result.Left[map[K]T](errTest)
	}
}
