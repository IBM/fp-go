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

package readerresult

import (
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/tuple"
)

// SequenceT functions convert multiple ReaderResult values into a single ReaderResult containing a tuple.
// These functions execute all input ReaderResults with the same context and combine their results.
//
// IMPORTANT: All ReaderResults are executed, even if one fails. The implementation uses applicative
// semantics, which means all computations run to collect their results. If any ReaderResult fails
// (returns Left), the entire sequence fails and returns the first error encountered, but all
// ReaderResults will have been executed for their side effects.
//
// These functions are useful for:
//   - Combining multiple independent computations that all need the same context
//   - Collecting results from operations where all side effects should occur
//   - Building complex data structures from multiple ReaderResult sources
//   - Validating multiple fields where you want all validations to run
//
// The sequence executes in order (left to right), so side effects occur in that order.

// SequenceT1 converts a single ReaderResult into a ReaderResult containing a 1-tuple.
// This is primarily useful for consistency in generic code or when you need to wrap
// a single value in a tuple structure.
//
// Type Parameters:
//   - A: The type of the value in the ReaderResult
//
// Parameters:
//   - a: The ReaderResult to wrap in a tuple
//
// Returns:
//   - A ReaderResult containing a Tuple1 with the value from the input
//
// Example:
//
//	rr := readerresult.Of(42)
//	tupled := readerresult.SequenceT1(rr)
//	result := tupled(t.Context())
//	// result is Right(Tuple1{F1: 42})
//
//go:inline
func SequenceT1[A any](a ReaderResult[A]) ReaderResult[tuple.Tuple1[A]] {
	return readereither.SequenceT1(a)
}

// SequenceT2 combines two ReaderResults into a single ReaderResult containing a 2-tuple.
// Both ReaderResults are executed with the same context. If either fails, the entire
// sequence fails.
//
// Type Parameters:
//   - A: The type of the first value
//   - B: The type of the second value
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//
// Returns:
//   - A ReaderResult containing a Tuple2 with both values
//
// Example:
//
//	getName := readerresult.Of("Alice")
//	getAge := readerresult.Of(30)
//	combined := readerresult.SequenceT2(getName, getAge)
//	result := combined(t.Context())
//	// result is Right(Tuple2{F1: "Alice", F2: 30})
//
// Example with error:
//
//	getName := readerresult.Of("Alice")
//	getAge := readerresult.Left[int](errors.New("age not found"))
//	combined := readerresult.SequenceT2(getName, getAge)
//	result := combined(t.Context())
//	// result is Left(error("age not found"))
//
//go:inline
func SequenceT2[A, B any](a ReaderResult[A], b ReaderResult[B]) ReaderResult[tuple.Tuple2[A, B]] {
	return readereither.SequenceT2(a, b)
}

// SequenceT3 combines three ReaderResults into a single ReaderResult containing a 3-tuple.
// All ReaderResults are executed sequentially with the same context. If any fails,
// the entire sequence fails immediately.
//
// Type Parameters:
//   - A: The type of the first value
//   - B: The type of the second value
//   - C: The type of the third value
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//   - c: The third ReaderResult
//
// Returns:
//   - A ReaderResult containing a Tuple3 with all three values
//
// Example:
//
//	getUserID := readerresult.Of(123)
//	getUserName := readerresult.Of("Alice")
//	getUserEmail := readerresult.Of("alice@example.com")
//	combined := readerresult.SequenceT3(getUserID, getUserName, getUserEmail)
//	result := combined(t.Context())
//	// result is Right(Tuple3{F1: 123, F2: "Alice", F3: "alice@example.com"})
//
// Example with context-aware operations:
//
//	fetchUser := func(ctx context.Context) result.Result[string] {
//	    if ctx.Err() != nil {
//	        return result.Error[string](ctx.Err())
//	    }
//	    return result.Of("Alice")
//	}
//	fetchAge := func(ctx context.Context) result.Result[int] {
//	    return result.Of(30)
//	}
//	fetchCity := func(ctx context.Context) result.Result[string] {
//	    return result.Of("NYC")
//	}
//	combined := readerresult.SequenceT3(fetchUser, fetchAge, fetchCity)
//	result := combined(t.Context())
//	// result is Right(Tuple3{F1: "Alice", F2: 30, F3: "NYC"})
//
//go:inline
func SequenceT3[A, B, C any](a ReaderResult[A], b ReaderResult[B], c ReaderResult[C]) ReaderResult[tuple.Tuple3[A, B, C]] {
	return readereither.SequenceT3(a, b, c)
}

// SequenceT4 combines four ReaderResults into a single ReaderResult containing a 4-tuple.
// All ReaderResults are executed sequentially with the same context. If any fails,
// the entire sequence fails immediately without executing the remaining ones.
//
// Type Parameters:
//   - A: The type of the first value
//   - B: The type of the second value
//   - C: The type of the third value
//   - D: The type of the fourth value
//
// Parameters:
//   - a: The first ReaderResult
//   - b: The second ReaderResult
//   - c: The third ReaderResult
//   - d: The fourth ReaderResult
//
// Returns:
//   - A ReaderResult containing a Tuple4 with all four values
//
// Example:
//
//	getID := readerresult.Of(123)
//	getName := readerresult.Of("Alice")
//	getEmail := readerresult.Of("alice@example.com")
//	getAge := readerresult.Of(30)
//	combined := readerresult.SequenceT4(getID, getName, getEmail, getAge)
//	result := combined(t.Context())
//	// result is Right(Tuple4{F1: 123, F2: "Alice", F3: "alice@example.com", F4: 30})
//
// Example with early failure:
//
//	getID := readerresult.Of(123)
//	getName := readerresult.Left[string](errors.New("name not found"))
//	getEmail := readerresult.Of("alice@example.com")  // Not executed
//	getAge := readerresult.Of(30)                      // Not executed
//	combined := readerresult.SequenceT4(getID, getName, getEmail, getAge)
//	result := combined(t.Context())
//	// result is Left(error("name not found"))
//	// getEmail and getAge are never executed due to early failure
//
// Example building a complex structure:
//
//	type UserProfile struct {
//	    ID    int
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	fetchUserData := readerresult.SequenceT4(
//	    fetchUserID(userID),
//	    fetchUserName(userID),
//	    fetchUserEmail(userID),
//	    fetchUserAge(userID),
//	)
//
//	buildProfile := readerresult.Map(func(t tuple.Tuple4[int, string, string, int]) UserProfile {
//	    return UserProfile{
//	        ID:    t.F1,
//	        Name:  t.F2,
//	        Email: t.F3,
//	        Age:   t.F4,
//	    }
//	})
//
//	userProfile := F.Pipe1(fetchUserData, buildProfile)
//	result := userProfile(t.Context())
//
//go:inline
func SequenceT4[A, B, C, D any](a ReaderResult[A], b ReaderResult[B], c ReaderResult[C], d ReaderResult[D]) ReaderResult[tuple.Tuple4[A, B, C, D]] {
	return readereither.SequenceT4(a, b, c, d)
}
