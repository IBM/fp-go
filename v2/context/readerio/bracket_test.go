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

package readerio

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/io"
	"github.com/stretchr/testify/assert"
)

// mockResource simulates a resource that tracks its lifecycle
type mockResource struct {
	id       int
	acquired bool
	released bool
	used     bool
}

// TestBracket_Success tests that Bracket properly manages resource lifecycle on success
func TestBracket_Success(t *testing.T) {
	resource := &mockResource{id: 1}

	// Acquire resource
	acquire := func(ctx context.Context) io.IO[*mockResource] {
		return func() *mockResource {
			resource.acquired = true
			return resource
		}
	}

	// Use resource
	use := func(r *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				r.used = true
				return "success"
			}
		}
	}

	// Release resource
	release := func(r *mockResource, result string) ReaderIO[any] {
		return func(ctx context.Context) io.IO[any] {
			return func() any {
				r.released = true
				return nil
			}
		}
	}

	// Execute bracket
	operation := Bracket(acquire, use, release)
	result := operation(context.Background())()

	// Verify lifecycle
	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.used, "Resource should be used")
	assert.True(t, resource.released, "Resource should be released")
	assert.Equal(t, "success", result)
}

// TestBracket_MultipleResources tests managing multiple resources
func TestBracket_MultipleResources(t *testing.T) {
	resource1 := &mockResource{id: 1}
	resource2 := &mockResource{id: 2}

	acquire1 := func(ctx context.Context) io.IO[*mockResource] {
		return func() *mockResource {
			resource1.acquired = true
			return resource1
		}
	}

	use1 := func(r1 *mockResource) ReaderIO[*mockResource] {
		return func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				r1.used = true
				resource2.acquired = true
				return resource2
			}
		}
	}

	release1 := func(r1 *mockResource, result string) ReaderIO[any] {
		return func(ctx context.Context) io.IO[any] {
			return func() any {
				r1.released = true
				return nil
			}
		}
	}

	// Nested bracket for second resource
	use2 := func(r2 *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				r2.used = true
				return "both used"
			}
		}
	}

	release2 := func(r2 *mockResource, result string) ReaderIO[any] {
		return func(ctx context.Context) io.IO[any] {
			return func() any {
				r2.released = true
				return nil
			}
		}
	}

	// Compose brackets
	operation := Bracket(acquire1, func(r1 *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			r2 := use1(r1)(ctx)()
			return Bracket(
				func(ctx context.Context) io.IO[*mockResource] {
					return func() *mockResource { return r2 }
				},
				use2,
				release2,
			)(ctx)
		}
	}, release1)

	result := operation(context.Background())()

	assert.True(t, resource1.acquired)
	assert.True(t, resource1.used)
	assert.True(t, resource1.released)
	assert.True(t, resource2.acquired)
	assert.True(t, resource2.used)
	assert.True(t, resource2.released)
	assert.Equal(t, "both used", result)
}

// TestWithResource_Success tests WithResource with successful operation
func TestWithResource_Success(t *testing.T) {
	resource := &mockResource{id: 1}

	// Define resource management
	withResource := WithResource[*mockResource, string, any](
		func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				resource.acquired = true
				return resource
			}
		},
		func(r *mockResource) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					r.released = true
					return nil
				}
			}
		},
	)

	// Use resource
	operation := withResource(func(r *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				r.used = true
				return "result"
			}
		}
	})

	result := operation(context.Background())()

	assert.True(t, resource.acquired)
	assert.True(t, resource.used)
	assert.True(t, resource.released)
	assert.Equal(t, "result", result)
}

// TestWithResource_Reusability tests that WithResource can be reused with different operations
func TestWithResource_Reusability(t *testing.T) {
	callCount := 0

	withResource := WithResource[*mockResource, int, any](
		func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				callCount++
				return &mockResource{id: callCount, acquired: true}
			}
		},
		func(r *mockResource) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					r.released = true
					return nil
				}
			}
		},
	)

	// First operation
	op1 := withResource(func(r *mockResource) ReaderIO[int] {
		return func(ctx context.Context) io.IO[int] {
			return func() int {
				r.used = true
				return r.id * 2
			}
		}
	})

	result1 := op1(context.Background())()
	assert.Equal(t, 2, result1)
	assert.Equal(t, 1, callCount)

	// Second operation (should create new resource)
	op2 := withResource(func(r *mockResource) ReaderIO[int] {
		return func(ctx context.Context) io.IO[int] {
			return func() int {
				r.used = true
				return r.id * 3
			}
		}
	})

	result2 := op2(context.Background())()
	assert.Equal(t, 6, result2)
	assert.Equal(t, 2, callCount)
}

// TestWithResource_DifferentResultTypes tests WithResource with different result types
func TestWithResource_DifferentResultTypes(t *testing.T) {
	resource := &mockResource{id: 42}

	withResourceInt := WithResource[*mockResource, int, any](
		func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				resource.acquired = true
				return resource
			}
		},
		func(r *mockResource) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					r.released = true
					return nil
				}
			}
		},
	)

	// Operation returning int
	opInt := withResourceInt(func(r *mockResource) ReaderIO[int] {
		return func(ctx context.Context) io.IO[int] {
			return func() int {
				return r.id
			}
		}
	})

	resultInt := opInt(context.Background())()
	assert.Equal(t, 42, resultInt)

	// Reset resource state
	resource.acquired = false
	resource.released = false

	// Create new WithResource for string type
	withResourceString := WithResource[*mockResource, string, any](
		func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				resource.acquired = true
				return resource
			}
		},
		func(r *mockResource) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					r.released = true
					return nil
				}
			}
		},
	)

	// Operation returning string
	opString := withResourceString(func(r *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				return "value"
			}
		}
	})

	resultString := opString(context.Background())()
	assert.Equal(t, "value", resultString)
	assert.True(t, resource.released)
}

// TestWithResource_ContextPropagation tests that context is properly propagated
func TestWithResource_ContextPropagation(t *testing.T) {
	type contextKey string
	const key contextKey = "test-key"

	withResource := WithResource[string, string, any](
		func(ctx context.Context) io.IO[string] {
			return func() string {
				value := ctx.Value(key)
				if value != nil {
					return value.(string)
				}
				return "no-value"
			}
		},
		func(r string) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					return nil
				}
			}
		},
	)

	operation := withResource(func(r string) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				return r + "-processed"
			}
		}
	})

	ctx := context.WithValue(context.Background(), key, "test-value")
	result := operation(ctx)()

	assert.Equal(t, "test-value-processed", result)
}

// TestWithResource_ErrorInRelease tests behavior when release function encounters an error
func TestWithResource_ErrorInRelease(t *testing.T) {
	resource := &mockResource{id: 1}
	releaseError := errors.New("release failed")

	withResource := WithResource[*mockResource, string, error](
		func(ctx context.Context) io.IO[*mockResource] {
			return func() *mockResource {
				resource.acquired = true
				return resource
			}
		},
		func(r *mockResource) ReaderIO[error] {
			return func(ctx context.Context) io.IO[error] {
				return func() error {
					r.released = true
					return releaseError
				}
			}
		},
	)

	operation := withResource(func(r *mockResource) ReaderIO[string] {
		return func(ctx context.Context) io.IO[string] {
			return func() string {
				r.used = true
				return "success"
			}
		}
	})

	result := operation(context.Background())()

	// Operation should succeed even if release returns error
	assert.Equal(t, "success", result)
	assert.True(t, resource.acquired)
	assert.True(t, resource.used)
	assert.True(t, resource.released)
}

// BenchmarkBracket benchmarks the Bracket function
func BenchmarkBracket(b *testing.B) {
	acquire := func(ctx context.Context) io.IO[int] {
		return func() int {
			return 42
		}
	}

	use := func(n int) ReaderIO[int] {
		return func(ctx context.Context) io.IO[int] {
			return func() int {
				return n * 2
			}
		}
	}

	release := func(n int, result int) ReaderIO[any] {
		return func(ctx context.Context) io.IO[any] {
			return func() any {
				return nil
			}
		}
	}

	operation := Bracket(acquire, use, release)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		operation(ctx)()
	}
}

// BenchmarkWithResource benchmarks the WithResource function
func BenchmarkWithResource(b *testing.B) {
	withResource := WithResource[int, int, any](
		func(ctx context.Context) io.IO[int] {
			return func() int {
				return 42
			}
		},
		func(n int) ReaderIO[any] {
			return func(ctx context.Context) io.IO[any] {
				return func() any {
					return nil
				}
			}
		},
	)

	operation := withResource(func(n int) ReaderIO[int] {
		return func(ctx context.Context) io.IO[int] {
			return func() int {
				return n * 2
			}
		}
	})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		operation(ctx)()
	}
}

// Made with Bob
