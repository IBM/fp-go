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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

// FromReaderIOResult converts a ReaderIOResult[Reader] into a Reader.
//
// This function bridges the gap between context-aware, IO-based computations that may fail
// (ReaderIOResult) and the simpler Reader type used for test assertions. It executes the
// ReaderIOResult computation using the test's context, handles any potential errors by
// converting them to test failures via NoError, and returns the resulting Reader.
//
// The conversion process:
//  1. Executes the ReaderIOResult with the test context (t.Context())
//  2. Runs the resulting IO operation ()
//  3. Extracts the Result, converting errors to test failures using NoError
//  4. Returns a Reader that can be applied to *testing.T
//
// This is particularly useful when you have test assertions that need to:
//   - Access context for cancellation or deadlines
//   - Perform IO operations (file access, network calls, etc.)
//   - Handle potential errors gracefully in tests
//
// Parameters:
//   - ri: A ReaderIOResult that produces a Reader when given a context and executed
//
// Returns:
//   - A Reader that can be directly applied to *testing.T for assertion
//
// Example:
//
//	func TestWithContext(t *testing.T) {
//	    // Create a ReaderIOResult that performs an IO operation
//	    checkDatabase := func(ctx context.Context) func() result.Result[assert.Reader] {
//	        return func() result.Result[assert.Reader] {
//	            // Simulate database check
//	            if err := db.PingContext(ctx); err != nil {
//	                return result.Error[assert.Reader](err)
//	            }
//	            return result.Of[assert.Reader](assert.NoError(nil))
//	        }
//	    }
//
//	    // Convert to Reader and execute
//	    assertion := assert.FromReaderIOResult(checkDatabase)
//	    assertion(t)
//	}
func FromReaderIOResult(ri ReaderIOResult[Reader]) Reader {
	return func(t *testing.T) bool {
		return F.Pipe1(
			ri(t.Context())(),
			result.GetOrElse(NoError),
		)(t)
	}
}

// FromReaderIO converts a ReaderIO[Reader] into a Reader.
//
// This function bridges the gap between context-aware, IO-based computations (ReaderIO)
// and the simpler Reader type used for test assertions. It executes the ReaderIO
// computation using the test's context and returns the resulting Reader.
//
// Unlike FromReaderIOResult, this function does not handle errors explicitly - it assumes
// the IO operation will succeed or that any errors are handled within the ReaderIO itself.
//
// The conversion process:
//  1. Executes the ReaderIO with the test context (t.Context())
//  2. Runs the resulting IO operation ()
//  3. Returns a Reader that can be applied to *testing.T
//
// This is particularly useful when you have test assertions that need to:
//   - Access context for cancellation or deadlines
//   - Perform IO operations that don't fail (or handle failures internally)
//   - Integrate with context-aware testing utilities
//
// Parameters:
//   - ri: A ReaderIO that produces a Reader when given a context and executed
//
// Returns:
//   - A Reader that can be directly applied to *testing.T for assertion
//
// Example:
//
//	func TestWithIO(t *testing.T) {
//	    // Create a ReaderIO that performs an IO operation
//	    logAndCheck := func(ctx context.Context) func() assert.Reader {
//	        return func() assert.Reader {
//	            // Log something using context
//	            logger.InfoContext(ctx, "Running test")
//	            // Return an assertion
//	            return assert.Equal(42)(computeValue())
//	        }
//	    }
//
//	    // Convert to Reader and execute
//	    assertion := assert.FromReaderIO(logAndCheck)
//	    assertion(t)
//	}
func FromReaderIO(ri ReaderIO[Reader]) Reader {
	return func(t *testing.T) bool {
		return ri(t.Context())()(t)
	}
}
