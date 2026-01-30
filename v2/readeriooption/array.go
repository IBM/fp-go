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

package readeriooption

import (
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// TraverseArray transforms an array by applying a function that returns a ReaderIOOption to each element.
// If any element results in None, the entire result is None.
// Otherwise, returns Some containing an array of all the unwrapped values.
//
// This is useful for performing a sequence of operations that may fail on each element of an array,
// where you want all operations to succeed or the entire computation to fail.
//
// Example:
//
//	type DB struct { ... }
//
//	findUser := func(id int) readeroption.ReaderIOOption[DB, User] { ... }
//
//	userIDs := []int{1, 2, 3}
//	result := F.Pipe1(
//	    readeroption.Of[DB](userIDs),
//	    readeroption.Chain(readeroption.TraverseArray[DB](findUser)),
//	)
//	// result will be Some([]User) if all users are found, None otherwise
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return RA.Traverse[[]A, []B](
		Of,
		Map,
		Ap,

		f,
	)
}

// TraverseArrayWithIndex is like TraverseArray but the function also receives the index of each element.
//
// Example:
//
//	type DB struct { ... }
//
//	processWithIndex := func(idx int, value string) readeroption.ReaderIOOption[DB, Result] {
//	    // Use idx in processing
//	    return readeroption.Asks(func(db DB) option.Option[Result] { ... })
//	}
//
//	values := []string{"a", "b", "c"}
//	result := readeroption.TraverseArrayWithIndex[DB](processWithIndex)(values)
func TraverseArrayWithIndex[E, A, B any](f func(int, A) ReaderIOOption[E, B]) func([]A) ReaderIOOption[E, []B] {
	return RA.TraverseWithIndex[[]A, []B](
		Of,
		Map,
		Ap,

		f,
	)
}
