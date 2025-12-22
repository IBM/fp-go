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

package iter

import "github.com/IBM/fp-go/v2/option"

// First returns the first element from an [Iterator] wrapped in an [Option].
//
// This function attempts to retrieve the first element from the iterator. If the iterator
// contains at least one element, it returns Some(element). If the iterator is empty,
// it returns None. The function consumes only the first element of the iterator.
//
// RxJS Equivalent: [first] - https://rxjs.dev/api/operators/first
//
// Type Parameters:
//   - U: The type of elements in the iterator
//
// Parameters:
//   - mu: The input iterator to get the first element from
//
// Returns:
//   - Option[U]: Some(first element) if the iterator is non-empty, None otherwise
//
// Example with non-empty sequence:
//
//	seq := iter.From(1, 2, 3, 4, 5)
//	first := iter.First(seq)
//	// Returns: Some(1)
//
// Example with empty sequence:
//
//	seq := iter.Empty[int]()
//	first := iter.First(seq)
//	// Returns: None
//
// Example with filtered sequence:
//
//	seq := iter.From(1, 2, 3, 4, 5)
//	filtered := iter.Filter(func(x int) bool { return x > 3 })(seq)
//	first := iter.First(filtered)
//	// Returns: Some(4)
func First[U any](mu Seq[U]) Option[U] {
	for u := range mu {
		return option.Some(u)
	}
	return option.None[U]()
}
