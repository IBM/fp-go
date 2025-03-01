// Copyright (c) 2023 IBM Corp.
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

package array

import (
	E "github.com/IBM/fp-go/v2/eq"
)

func equals[T any](left []T, right []T, eq func(T, T) bool) bool {
	if len(left) != len(right) {
		return false
	}
	for i, v1 := range left {
		v2 := right[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

func Eq[T any](e E.Eq[T]) E.Eq[[]T] {
	eq := e.Equals
	return E.FromEquals(func(left, right []T) bool {
		return equals(left, right, eq)
	})
}
