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

package readeroption

import (
	G "github.com/IBM/fp-go/v2/readeroption/generic"
)

// TraverseArray transforms an array by applying a function that returns a ReaderOption to each element.
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
//	findUser := func(id int) readeroption.ReaderOption[DB, User] { ... }
//
//	userIDs := []int{1, 2, 3}
//	result := F.Pipe1(
//	    readeroption.Of[DB](userIDs),
//	    readeroption.Chain(readeroption.TraverseArray[DB](findUser)),
//	)
//	// result will be Some([]User) if all users are found, None otherwise
func TraverseArray[E, A, B any](f Kleisli[E, A, B]) Kleisli[E, []A, []B] {
	return G.TraverseArray[ReaderOption[E, B], ReaderOption[E, []B], []A](f)
}

// TraverseArrayWithIndex is like TraverseArray but the function also receives the index of each element.
//
// Example:
//
//	type DB struct { ... }
//
//	processWithIndex := func(idx int, value string) readeroption.ReaderOption[DB, Result] {
//	    // Use idx in processing
//	    return readeroption.Asks(func(db DB) option.Option[Result] { ... })
//	}
//
//	values := []string{"a", "b", "c"}
//	result := readeroption.TraverseArrayWithIndex[DB](processWithIndex)(values)
func TraverseArrayWithIndex[E, A, B any](f func(int, A) ReaderOption[E, B]) func([]A) ReaderOption[E, []B] {
	return G.TraverseArrayWithIndex[ReaderOption[E, B], ReaderOption[E, []B], []A](f)
}

// SequenceArray converts an array of ReaderOption values into a ReaderOption of an array.
// If any element is None, the entire result is None.
// Otherwise, returns Some containing an array of all the unwrapped values.
//
// This is useful when you have multiple independent ReaderOption computations and want to
// combine their results into a single array.
//
// Example:
//
//	type Config struct { ... }
//
//	user1 := readeroption.Of[Config](User{ID: 1, Name: "Alice"})
//	user2 := readeroption.Of[Config](User{ID: 2, Name: "Bob"})
//	user3 := readeroption.None[Config, User]()
//
//	result := readeroption.SequenceArray([]readeroption.ReaderOption[Config, User]{
//	    user1, user2, user3,
//	})
//	// result(config) will be option.None[[]User]() because user3 is None
//
//	result2 := readeroption.SequenceArray([]readeroption.ReaderOption[Config, User]{
//	    user1, user2,
//	})
//	// result2(config) will be option.Some([]User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}})
func SequenceArray[E, A any](ma []ReaderOption[E, A]) ReaderOption[E, []A] {
	return G.SequenceArray[ReaderOption[E, A], ReaderOption[E, []A]](ma)
}
