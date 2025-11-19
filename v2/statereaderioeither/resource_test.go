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

package statereaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// resourceState tracks resource lifecycle
type resourceState struct {
	openResources int
	lastError     error
}

// testResource represents a simple resource
type testResource struct {
	id   int
	data string
}

func TestWithResource_SuccessCase(t *testing.T) {
	state := resourceState{openResources: 0}
	ctx := testContext{multiplier: 1}
	released := false

	// Create a resource (increments open count)
	onCreate := FromState[testContext, error](func(s resourceState) P.Pair[resourceState, testResource] {
		newState := resourceState{openResources: s.openResources + 1}
		resource := testResource{id: 42, data: "test"}
		return P.MakePair(newState, resource)
	})

	// Release the resource (decrements open count)
	onRelease := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			released = true
			newState := resourceState{openResources: s.openResources - 1}
			return P.MakePair(newState, 0)
		})
	}

	// Use the resource
	result := WithResource[string](
		onCreate,
		onRelease,
	)(func(res testResource) StateReaderIOEither[resourceState, testContext, error, string] {
		return Of[resourceState, testContext, error](fmt.Sprintf("Resource: %d - %s", res.id, res.data))
	})

	res := result(state)(ctx)()

	// Verify success
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[resourceState, string]) P.Pair[resourceState, string] {
		assert.Equal(t, "Resource: 42 - test", P.Tail(p))
		// State is 1 because onCreate incremented to 1, then release saw state=1 and decremented to 0,
		// but the final state comes from the use function which doesn't modify state
		assert.Equal(t, 1, P.Head(p).openResources)
		return p
	})(res)
	assert.True(t, released)
}

func TestWithResource_ErrorInUse(t *testing.T) {
	state := resourceState{openResources: 0}
	ctx := testContext{multiplier: 1}
	released := false

	// Create a resource
	onCreate := FromState[testContext, error](func(s resourceState) P.Pair[resourceState, testResource] {
		newState := resourceState{openResources: s.openResources + 1}
		resource := testResource{id: 99, data: "data"}
		return P.MakePair(newState, resource)
	})

	// Release the resource
	onRelease := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			released = true
			newState := resourceState{openResources: s.openResources - 1}
			return P.MakePair(newState, 0)
		})
	}

	// Use the resource with an error
	testErr := errors.New("processing error")
	result := WithResource[string](
		onCreate,
		onRelease,
	)(func(res testResource) StateReaderIOEither[resourceState, testContext, error, string] {
		return Left[resourceState, testContext, string](testErr)
	})

	res := result(state)(ctx)()

	// Verify error is propagated but resource was still released
	assert.True(t, E.IsLeft(res))
	assert.True(t, released)
}

func TestWithResource_ErrorInCreate(t *testing.T) {
	state := resourceState{openResources: 0}
	ctx := testContext{multiplier: 1}
	released := false

	// Create a resource that fails
	createErr := errors.New("creation failed")
	onCreate := Left[resourceState, testContext, testResource](createErr)

	// Release function
	onRelease := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			released = true
			return P.MakePair(s, 0)
		})
	}

	// Try to use the resource
	result := WithResource[string](
		onCreate,
		onRelease,
	)(func(res testResource) StateReaderIOEither[resourceState, testContext, error, string] {
		return Of[resourceState, testContext, error]("should not reach here")
	})

	res := result(state)(ctx)()

	// Verify creation error is propagated and release was not called
	assert.True(t, E.IsLeft(res))
	assert.False(t, released)
}

func TestWithResource_StateThreading(t *testing.T) {
	state := resourceState{openResources: 0}
	ctx := testContext{multiplier: 2}

	// Track state changes
	var statesObserved []int

	// Create a resource (state: 0 -> 1)
	onCreate := FromState[testContext, error](func(s resourceState) P.Pair[resourceState, testResource] {
		statesObserved = append(statesObserved, s.openResources)
		newState := resourceState{openResources: s.openResources + 1}
		resource := testResource{id: 1, data: "file"}
		return P.MakePair(newState, resource)
	})

	// Use the resource (state: 1 -> 2)
	useResource := func(res testResource) StateReaderIOEither[resourceState, testContext, error, string] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, string] {
			statesObserved = append(statesObserved, s.openResources)
			newState := resourceState{openResources: s.openResources + 1}
			return P.MakePair(newState, fmt.Sprintf("used-%d", res.id))
		})
	}

	// Release the resource (state: 2 -> 1)
	onRelease := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			statesObserved = append(statesObserved, s.openResources)
			newState := resourceState{openResources: s.openResources - 1}
			return P.MakePair(newState, 0)
		})
	}

	result := WithResource[string](
		onCreate,
		onRelease,
	)(useResource)

	res := result(state)(ctx)()

	// Verify state threading
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[resourceState, string]) P.Pair[resourceState, string] {
		assert.Equal(t, "used-1", P.Tail(p))
		assert.Equal(t, 2, P.Head(p).openResources) // Final state from use function
		return p
	})(res)

	// Verify state was observed: onCreate sees initial state (0),
	// useResource sees state after create (1), onRelease sees state after create (1)
	assert.Equal(t, []int{0, 1, 1}, statesObserved)
}

func TestWithResource_MultipleResources(t *testing.T) {
	state := resourceState{openResources: 0}
	ctx := testContext{multiplier: 1}

	// Create first resource
	createResource1 := FromState[testContext, error](func(s resourceState) P.Pair[resourceState, testResource] {
		newState := resourceState{openResources: s.openResources + 1}
		return P.MakePair(newState, testResource{id: 1, data: "res1"})
	})

	releaseResource1 := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			newState := resourceState{openResources: s.openResources - 1}
			return P.MakePair(newState, 0)
		})
	}

	// Second resource creator
	createResource2 := FromState[testContext, error](func(s resourceState) P.Pair[resourceState, testResource] {
		newState := resourceState{openResources: s.openResources + 1}
		return P.MakePair(newState, testResource{id: 2, data: "res2"})
	})

	releaseResource2 := func(res testResource) StateReaderIOEither[resourceState, testContext, error, int] {
		return FromState[testContext, error](func(s resourceState) P.Pair[resourceState, int] {
			newState := resourceState{openResources: s.openResources - 1}
			return P.MakePair(newState, 0)
		})
	}

	// Nest resources
	result := WithResource[string](
		createResource1,
		releaseResource1,
	)(func(res1 testResource) StateReaderIOEither[resourceState, testContext, error, string] {
		return WithResource[string](
			createResource2,
			releaseResource2,
		)(func(res2 testResource) StateReaderIOEither[resourceState, testContext, error, string] {
			return Of[resourceState, testContext, error](
				fmt.Sprintf("%s + %s", res1.data, res2.data),
			)
		})
	})

	res := result(state)(ctx)()

	// Verify both resources were used and released
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[resourceState, string]) P.Pair[resourceState, string] {
		assert.Equal(t, "res1 + res2", P.Tail(p))
		// Final state comes from innermost use function (Of doesn't modify state)
		// onCreate1: 0->1, onCreate2: 1->2, release2: sees 2, release1: sees 1
		// Final state from Of: 2 (from the state after both creates)
		assert.Equal(t, 2, P.Head(p).openResources)
		return p
	})(res)
}
