// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioresult

import (
	"context"
	"errors"
	"testing"

	P "github.com/IBM/fp-go/v2/pair"
	RES "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// resourceState tracks the lifecycle of resources for testing
type resourceState struct {
	resourcesCreated  int
	resourcesReleased int
	lastError         error
}

// mockResource represents a test resource
type mockResource struct {
	id      int
	isValid bool
}

// TestWithResourceSuccess tests successful resource creation, usage, and release
func TestWithResourceSuccess(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx := t.Context()

	// Create a resource
	onCreate := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				newState := resourceState{
					resourcesCreated:  s.resourcesCreated + 1,
					resourcesReleased: s.resourcesReleased,
				}
				res := mockResource{id: newState.resourcesCreated, isValid: true}
				return RES.Of(P.MakePair(newState, res))
			}
		}
	}

	// Release a resource
	onRelease := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					newState := resourceState{
						resourcesCreated:  s.resourcesCreated,
						resourcesReleased: s.resourcesReleased + 1,
					}
					return RES.Of(P.MakePair(newState, 0))
				}
			}
		}
	}

	// Use the resource
	useResource := func(res mockResource) StateReaderIOResult[resourceState, string] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, string]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, string]] {
				return func() Result[Pair[resourceState, string]] {
					result := "used resource " + string(rune(res.id+'0'))
					return RES.Of(P.MakePair(s, result))
				}
			}
		}
	}

	withResource := WithResource[string](onCreate, onRelease)
	result := withResource(useResource)
	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsRight(outcome))
	RES.Map(func(p Pair[resourceState, string]) Pair[resourceState, string] {
		state := P.Head(p)
		value := P.Tail(p)

		// Verify state updates
		// Note: Final state comes from the use function, not the release function
		// onCreate: 0->1, use: sees 1 (doesn't modify), release: sees 1 and increments released
		// The final state is from use function which saw state=1 with resourcesReleased=0
		assert.Equal(t, 1, state.resourcesCreated, "Resource should be created once")
		assert.Equal(t, 0, state.resourcesReleased, "Final state is from use function, before release")

		// Verify result
		assert.Equal(t, "used resource 1", value)

		return p
	})(outcome)
}

// TestWithResourceErrorInCreate tests error handling when resource creation fails
func TestWithResourceErrorInCreate(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx := t.Context()

	createError := errors.New("failed to create resource")

	// onCreate that fails
	onCreate := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				return RES.Left[Pair[resourceState, mockResource]](createError)
			}
		}
	}

	// Release should not be called if onCreate fails
	onRelease := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					t.Error("onRelease should not be called when onCreate fails")
					return RES.Of(P.MakePair(s, 0))
				}
			}
		}
	}

	useResource := func(res mockResource) StateReaderIOResult[resourceState, string] {
		return Of[resourceState]("should not reach here")
	}

	withResource := WithResource[string](onCreate, onRelease)
	result := withResource(useResource)
	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsLeft(outcome))
	RES.Fold(
		func(err error) bool {
			assert.Equal(t, createError, err)
			return true
		},
		func(p Pair[resourceState, string]) bool {
			t.Error("Expected error but got success")
			return false
		},
	)(outcome)
}

// TestWithResourceErrorInUse tests that resources are released even when usage fails
func TestWithResourceErrorInUse(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx := t.Context()

	useError := errors.New("failed to use resource")

	// Create a resource
	onCreate := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				newState := resourceState{
					resourcesCreated:  s.resourcesCreated + 1,
					resourcesReleased: s.resourcesReleased,
				}
				res := mockResource{id: 1, isValid: true}
				return RES.Of(P.MakePair(newState, res))
			}
		}
	}

	releaseWasCalled := false

	// Release should still be called even if use fails
	onRelease := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					releaseWasCalled = true
					newState := resourceState{
						resourcesCreated:  s.resourcesCreated,
						resourcesReleased: s.resourcesReleased + 1,
					}
					return RES.Of(P.MakePair(newState, 0))
				}
			}
		}
	}

	// Use that fails
	useResource := func(res mockResource) StateReaderIOResult[resourceState, string] {
		return Left[resourceState, string](useError)
	}

	withResource := WithResource[string](onCreate, onRelease)
	result := withResource(useResource)
	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsLeft(outcome))
	assert.True(t, releaseWasCalled, "onRelease should be called even when use fails")

	RES.Fold(
		func(err error) bool {
			assert.Equal(t, useError, err)
			return true
		},
		func(p Pair[resourceState, string]) bool {
			t.Error("Expected error but got success")
			return false
		},
	)(outcome)
}

// TestWithResourceStateThreading tests that state is properly threaded through all operations
func TestWithResourceStateThreading(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx := t.Context()

	// Create increments counter
	onCreate := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				newState := resourceState{
					resourcesCreated:  s.resourcesCreated + 1,
					resourcesReleased: s.resourcesReleased,
				}
				res := mockResource{id: newState.resourcesCreated, isValid: true}
				return RES.Of(P.MakePair(newState, res))
			}
		}
	}

	// Use observes the state after creation
	useResource := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					// Verify state was updated by onCreate
					assert.Equal(t, 1, s.resourcesCreated)
					assert.Equal(t, 0, s.resourcesReleased)
					return RES.Of(P.MakePair(s, s.resourcesCreated))
				}
			}
		}
	}

	// Release increments released counter
	onRelease := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					// Verify state was updated by onCreate and use
					assert.Equal(t, 1, s.resourcesCreated)
					assert.Equal(t, 0, s.resourcesReleased)

					newState := resourceState{
						resourcesCreated:  s.resourcesCreated,
						resourcesReleased: s.resourcesReleased + 1,
					}
					return RES.Of(P.MakePair(newState, 0))
				}
			}
		}
	}

	withResource := WithResource[int](onCreate, onRelease)
	result := withResource(useResource)
	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsRight(outcome))
	RES.Map(func(p Pair[resourceState, int]) Pair[resourceState, int] {
		finalState := P.Head(p)
		value := P.Tail(p)

		// Verify final state
		// Note: Final state is from the use function, which preserves the state it received
		// onCreate: 0->1, use: sees 1, release: sees 1 and increments released to 1
		// But final state is from use function where resourcesReleased=0
		assert.Equal(t, 1, finalState.resourcesCreated)
		assert.Equal(t, 0, finalState.resourcesReleased, "Final state is from use function, before release")
		assert.Equal(t, 1, value)

		return p
	})(outcome)
}

// TestWithResourceMultipleResources tests using WithResource multiple times (nesting)
func TestWithResourceMultipleResources(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx := t.Context()

	createResource := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				newState := resourceState{
					resourcesCreated:  s.resourcesCreated + 1,
					resourcesReleased: s.resourcesReleased,
				}
				res := mockResource{id: newState.resourcesCreated, isValid: true}
				return RES.Of(P.MakePair(newState, res))
			}
		}
	}

	releaseResource := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					newState := resourceState{
						resourcesCreated:  s.resourcesCreated,
						resourcesReleased: s.resourcesReleased + 1,
					}
					return RES.Of(P.MakePair(newState, 0))
				}
			}
		}
	}

	// Create two nested resources
	withResource1 := WithResource[int](createResource, releaseResource)
	withResource2 := WithResource[int](createResource, releaseResource)

	result := withResource1(func(res1 mockResource) StateReaderIOResult[resourceState, int] {
		return withResource2(func(res2 mockResource) StateReaderIOResult[resourceState, int] {
			// Both resources should be available
			return Of[resourceState](res1.id + res2.id)
		})
	})

	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsRight(outcome))
	RES.Map(func(p Pair[resourceState, int]) Pair[resourceState, int] {
		finalState := P.Head(p)
		value := P.Tail(p)

		// Both resources created, but final state is from innermost use function
		// onCreate1: 0->1, onCreate2: 1->2, use (Of): sees 2
		// Release functions execute but their state changes aren't in the final result
		assert.Equal(t, 2, finalState.resourcesCreated)
		assert.Equal(t, 0, finalState.resourcesReleased, "Final state is from use function, before releases")
		// res1.id = 1, res2.id = 2, sum = 3
		assert.Equal(t, 3, value)

		return p
	})(outcome)
}

// TestWithResourceContextCancellation tests behavior with context cancellation
func TestWithResourceContextCancellation(t *testing.T) {
	initialState := resourceState{resourcesCreated: 0, resourcesReleased: 0}
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // Cancel immediately

	cancelError := errors.New("context cancelled")

	// Create should respect context cancellation
	onCreate := func(s resourceState) ReaderIOResult[Pair[resourceState, mockResource]] {
		return func(ctx context.Context) IOResult[Pair[resourceState, mockResource]] {
			return func() Result[Pair[resourceState, mockResource]] {
				if ctx.Err() != nil {
					return RES.Left[Pair[resourceState, mockResource]](cancelError)
				}
				newState := resourceState{
					resourcesCreated:  s.resourcesCreated + 1,
					resourcesReleased: s.resourcesReleased,
				}
				res := mockResource{id: 1, isValid: true}
				return RES.Of(P.MakePair(newState, res))
			}
		}
	}

	onRelease := func(res mockResource) StateReaderIOResult[resourceState, int] {
		return func(s resourceState) ReaderIOResult[Pair[resourceState, int]] {
			return func(ctx context.Context) IOResult[Pair[resourceState, int]] {
				return func() Result[Pair[resourceState, int]] {
					newState := resourceState{
						resourcesCreated:  s.resourcesCreated,
						resourcesReleased: s.resourcesReleased + 1,
					}
					return RES.Of(P.MakePair(newState, 0))
				}
			}
		}
	}

	useResource := func(res mockResource) StateReaderIOResult[resourceState, string] {
		return Of[resourceState]("should not reach here")
	}

	withResource := WithResource[string](onCreate, onRelease)
	result := withResource(useResource)
	outcome := result(initialState)(ctx)()

	assert.True(t, RES.IsLeft(outcome))
	RES.Fold(
		func(err error) bool {
			assert.Equal(t, cancelError, err)
			return true
		},
		func(p Pair[resourceState, string]) bool {
			t.Error("Expected error but got success")
			return false
		},
	)(outcome)
}
