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
	"context"
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/result"
)

func TestFromReaderIOResult(t *testing.T) {
	t.Run("should pass when ReaderIOResult returns success with passing assertion", func(t *testing.T) {
		// Create a ReaderIOResult that returns a successful Reader
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				// Return a Reader that always passes
				return result.Of[Reader](func(t *testing.T) bool {
					return true
				})
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIOResult to pass when ReaderIOResult returns success")
		}
	})

	t.Run("should pass when ReaderIOResult returns success with Equal assertion", func(t *testing.T) {
		// Create a ReaderIOResult that returns a successful Equal assertion
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				return result.Of[Reader](Equal(42)(42))
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIOResult to pass with Equal assertion")
		}
	})

	t.Run("should fail when ReaderIOResult returns error", func(t *testing.T) {
		mockT := &testing.T{}

		// Create a ReaderIOResult that returns an error
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				return result.Left[Reader](errors.New("test error"))
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(mockT)
		if res {
			t.Error("Expected FromReaderIOResult to fail when ReaderIOResult returns error")
		}
	})

	t.Run("should fail when ReaderIOResult returns success but assertion fails", func(t *testing.T) {
		mockT := &testing.T{}

		// Create a ReaderIOResult that returns a failing assertion
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				return result.Of[Reader](Equal(42)(43))
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(mockT)
		if res {
			t.Error("Expected FromReaderIOResult to fail when assertion fails")
		}
	})

	t.Run("should use test context", func(t *testing.T) {
		contextUsed := false

		// Create a ReaderIOResult that checks if context is provided
		ri := func(ctx context.Context) func() result.Result[Reader] {
			if ctx != nil {
				contextUsed = true
			}
			return func() result.Result[Reader] {
				return result.Of[Reader](func(t *testing.T) bool {
					return true
				})
			}
		}

		reader := FromReaderIOResult(ri)
		reader(t)

		if !contextUsed {
			t.Error("Expected FromReaderIOResult to use test context")
		}
	})

	t.Run("should work with NoError assertion", func(t *testing.T) {
		// Create a ReaderIOResult that returns NoError assertion
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				return result.Of[Reader](NoError(nil))
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIOResult to pass with NoError assertion")
		}
	})

	t.Run("should work with complex assertions", func(t *testing.T) {
		// Create a ReaderIOResult with multiple composed assertions
		ri := func(ctx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				arr := []int{1, 2, 3}
				assertions := AllOf([]Reader{
					ArrayNotEmpty(arr),
					ArrayLength[int](3)(arr),
					ArrayContains(2)(arr),
				})
				return result.Of[Reader](assertions)
			}
		}

		reader := FromReaderIOResult(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIOResult to pass with complex assertions")
		}
	})
}

func TestFromReaderIO(t *testing.T) {
	t.Run("should pass when ReaderIO returns passing assertion", func(t *testing.T) {
		// Create a ReaderIO that returns a Reader that always passes
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				return func(t *testing.T) bool {
					return true
				}
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass when ReaderIO returns passing assertion")
		}
	})

	t.Run("should pass when ReaderIO returns Equal assertion", func(t *testing.T) {
		// Create a ReaderIO that returns an Equal assertion
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				return Equal(42)(42)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with Equal assertion")
		}
	})

	t.Run("should fail when ReaderIO returns failing assertion", func(t *testing.T) {
		mockT := &testing.T{}

		// Create a ReaderIO that returns a failing assertion
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				return Equal(42)(43)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(mockT)
		if res {
			t.Error("Expected FromReaderIO to fail when assertion fails")
		}
	})

	t.Run("should use test context", func(t *testing.T) {
		contextUsed := false

		// Create a ReaderIO that checks if context is provided
		ri := func(ctx context.Context) func() Reader {
			if ctx != nil {
				contextUsed = true
			}
			return func() Reader {
				return func(t *testing.T) bool {
					return true
				}
			}
		}

		reader := FromReaderIO(ri)
		reader(t)

		if !contextUsed {
			t.Error("Expected FromReaderIO to use test context")
		}
	})

	t.Run("should work with NoError assertion", func(t *testing.T) {
		// Create a ReaderIO that returns NoError assertion
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				return NoError(nil)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with NoError assertion")
		}
	})

	t.Run("should work with Error assertion", func(t *testing.T) {
		// Create a ReaderIO that returns Error assertion
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				return Error(errors.New("expected error"))
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with Error assertion")
		}
	})

	t.Run("should work with complex assertions", func(t *testing.T) {
		// Create a ReaderIO with multiple composed assertions
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				mp := map[string]int{"a": 1, "b": 2}
				return AllOf([]Reader{
					RecordNotEmpty(mp),
					RecordLength[string, int](2)(mp),
					ContainsKey[int]("a")(mp),
				})
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with complex assertions")
		}
	})

	t.Run("should work with string assertions", func(t *testing.T) {
		// Create a ReaderIO with string assertions
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				str := "hello world"
				return AllOf([]Reader{
					StringNotEmpty(str),
					StringLength[any, any](11)(str),
				})
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with string assertions")
		}
	})

	t.Run("should work with Result assertions", func(t *testing.T) {
		// Create a ReaderIO with Result assertions
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				successResult := result.Of[int](42)
				return Success(successResult)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with Success assertion")
		}
	})

	t.Run("should work with Failure assertion", func(t *testing.T) {
		// Create a ReaderIO with Failure assertion
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				failureResult := result.Left[int](errors.New("test error"))
				return Failure(failureResult)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)
		if !res {
			t.Error("Expected FromReaderIO to pass with Failure assertion")
		}
	})
}

// TestFromReaderIOResultIntegration tests integration scenarios
func TestFromReaderIOResultIntegration(t *testing.T) {
	t.Run("should work in a realistic scenario with context cancellation", func(t *testing.T) {
		// Create a ReaderIOResult that uses the context
		ri := func(testCtx context.Context) func() result.Result[Reader] {
			return func() result.Result[Reader] {
				// Check if context is valid
				if testCtx == nil {
					return result.Left[Reader](errors.New("context is nil"))
				}

				// Return a successful assertion
				return result.Of[Reader](Equal("test")("test"))
			}
		}

		// Use the actual testing.T from the subtest
		reader := FromReaderIOResult(ri)
		res := reader(t)
		if !res {
			t.Error("Expected integration test to pass")
		}
	})
}

// TestFromReaderIOIntegration tests integration scenarios
func TestFromReaderIOIntegration(t *testing.T) {
	t.Run("should work in a realistic scenario with logging", func(t *testing.T) {
		logCalled := false

		// Create a ReaderIO that simulates logging
		ri := func(ctx context.Context) func() Reader {
			return func() Reader {
				// Simulate logging with context
				if ctx != nil {
					logCalled = true
				}

				// Return an assertion
				return Equal(100)(100)
			}
		}

		reader := FromReaderIO(ri)
		res := reader(t)

		if !res {
			t.Error("Expected integration test to pass")
		}

		if !logCalled {
			t.Error("Expected logging to be called")
		}
	})
}
