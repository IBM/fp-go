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

package stateio

import (
	"testing"

	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

type ResourceState struct {
	ResourceCreated bool
	Value           int
}

func TestWithResource(t *testing.T) {
	initial := ResourceState{ResourceCreated: false, Value: 0}

	// Create resource
	onCreate := func(s ResourceState) IO[Pair[ResourceState, string]] {
		return func() Pair[ResourceState, string] {
			newState := ResourceState{ResourceCreated: true, Value: s.Value}
			return pair.MakePair(newState, "resource-handle")
		}
	}

	// Release resource (cleanup function)
	onRelease := func(res string) StateIO[ResourceState, int] {
		return func(s ResourceState) IO[Pair[ResourceState, int]] {
			return func() Pair[ResourceState, int] {
				// Release doesn't modify state in this test, just returns
				return pair.MakePair(s, 0)
			}
		}
	}

	// Use resource
	useResource := func(res string) StateIO[ResourceState, int] {
		return func(s ResourceState) IO[Pair[ResourceState, int]] {
			return func() Pair[ResourceState, int] {
				// Verify we received the resource handle
				assert.Equal(t, "resource-handle", res)
				newState := ResourceState{ResourceCreated: s.ResourceCreated, Value: 42}
				return pair.MakePair(newState, 42)
			}
		}
	}

	// Create the resource management computation
	withRes := WithResource[int](onCreate, onRelease)
	computation := withRes(useResource)

	result := computation(initial)()

	// Verify the resource was created and used
	assert.Equal(t, 42, pair.Tail(result))
	finalState := pair.Head(result)
	assert.True(t, finalState.ResourceCreated)
	assert.Equal(t, 42, finalState.Value)
}

func TestWithResourceChained(t *testing.T) {
	initial := ResourceState{ResourceCreated: false, Value: 0}

	// Create resource
	onCreate := func(s ResourceState) IO[Pair[ResourceState, int]] {
		return func() Pair[ResourceState, int] {
			newState := ResourceState{ResourceCreated: true, Value: s.Value}
			return pair.MakePair(newState, 100)
		}
	}

	// Release resource
	onRelease := func(res int) StateIO[ResourceState, int] {
		return func(s ResourceState) IO[Pair[ResourceState, int]] {
			return func() Pair[ResourceState, int] {
				return pair.MakePair(s, 0)
			}
		}
	}

	// Use resource with chaining
	useResource := func(res int) StateIO[ResourceState, int] {
		return MonadChain(
			Of[ResourceState](res),
			func(r int) StateIO[ResourceState, int] {
				return func(s ResourceState) IO[Pair[ResourceState, int]] {
					return func() Pair[ResourceState, int] {
						newState := ResourceState{
							ResourceCreated: s.ResourceCreated,
							Value:           r * 2,
						}
						return pair.MakePair(newState, r*2)
					}
				}
			},
		)
	}

	// Create the resource management computation
	withRes := WithResource[int](onCreate, onRelease)
	computation := withRes(useResource)

	result := computation(initial)()

	// Verify the resource was created and used
	assert.Equal(t, 200, pair.Tail(result))
	finalState := pair.Head(result)
	assert.True(t, finalState.ResourceCreated)
	assert.Equal(t, 200, finalState.Value)
}
