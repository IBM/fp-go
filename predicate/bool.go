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

package predicate

func Not[A any](predicate func(A) bool) func(A) bool {
	return func(a A) bool {
		return !predicate(a)
	}
}

// And creates a predicate that combines other predicates via &&
func And[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) && second(a)
		}
	}
}

// Or creates a predicate that combines other predicates via ||
func Or[A any](second func(A) bool) func(func(A) bool) func(A) bool {
	return func(first func(A) bool) func(A) bool {
		return func(a A) bool {
			return first(a) || second(a)
		}
	}
}
