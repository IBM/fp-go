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
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockResource simulates a resource that needs cleanup
type mockResource struct {
	id       int
	closed   bool
	closeMu  sync.Mutex
	closeErr error
}

func (m *mockResource) Close() error {
	m.closeMu.Lock()
	defer m.closeMu.Unlock()
	m.closed = true
	return m.closeErr
}

func (m *mockResource) IsClosed() bool {
	m.closeMu.Lock()
	defer m.closeMu.Unlock()
	return m.closed
}

// mockCloser implements io.Closer for testing WithCloser
type mockCloser struct {
	*mockResource
}

func TestBracketExtended(t *testing.T) {
	t.Run("successful acquire, use, and release with real resource", func(t *testing.T) {
		resource := &mockResource{id: 1}
		released := false

		result := Bracket(
			// Acquire
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			// Use
			func(r *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					return r.id * 2, nil
				}
			},
			// Release
			func(r *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					released = true
					assert.Equal(t, 2, result)
					assert.NoError(t, err)
					return nil, r.Close()
				}
			},
		)

		value, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, value)
		assert.True(t, released)
		assert.True(t, resource.IsClosed())
	})

	t.Run("acquire fails - release not called", func(t *testing.T) {
		acquireErr := errors.New("acquire failed")
		released := false

		result := Bracket(
			// Acquire fails
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return nil, acquireErr
				}
			},
			// Use (should not be called)
			func(r *mockResource) ReaderResult[int] {
				t.Fatal("use should not be called when acquire fails")
				return func(ctx context.Context) (int, error) {
					return 0, nil
				}
			},
			// Release (should not be called)
			func(r *mockResource, result int, err error) ReaderResult[any] {
				released = true
				return func(ctx context.Context) (any, error) {
					return nil, nil
				}
			},
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, acquireErr, err)
		assert.False(t, released)
	})

	t.Run("use fails - release still called", func(t *testing.T) {
		resource := &mockResource{id: 1}
		useErr := errors.New("use failed")
		released := false

		result := Bracket(
			// Acquire
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			// Use fails
			func(r *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					return 0, useErr
				}
			},
			// Release (should still be called)
			func(r *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					released = true
					assert.Equal(t, 0, result)
					assert.Equal(t, useErr, err)
					return nil, r.Close()
				}
			},
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, useErr, err)
		assert.True(t, released)
		assert.True(t, resource.IsClosed())
	})

	t.Run("release fails - error propagated", func(t *testing.T) {
		resource := &mockResource{id: 1, closeErr: errors.New("close failed")}
		released := false

		result := Bracket(
			// Acquire
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			// Use succeeds
			func(r *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					return r.id * 2, nil
				}
			},
			// Release fails
			func(r *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					released = true
					return nil, r.Close()
				}
			},
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "close failed", err.Error())
		assert.True(t, released)
		assert.True(t, resource.IsClosed())
	})

	t.Run("both use and release fail - use error takes precedence", func(t *testing.T) {
		resource := &mockResource{id: 1, closeErr: errors.New("close failed")}
		useErr := errors.New("use failed")
		released := false

		result := Bracket(
			// Acquire
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			// Use fails
			func(r *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					return 0, useErr
				}
			},
			// Release also fails
			func(r *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					released = true
					assert.Equal(t, useErr, err)
					return nil, r.Close()
				}
			},
		)

		_, err := result(context.Background())
		assert.Error(t, err)
		// The use error should be returned
		assert.Equal(t, useErr, err)
		assert.True(t, released)
		assert.True(t, resource.IsClosed())
	})

	t.Run("context cancellation during use", func(t *testing.T) {
		resource := &mockResource{id: 1}
		released := false

		result := Bracket(
			// Acquire
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			// Use checks context
			func(r *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					select {
					case <-ctx.Done():
						return 0, ctx.Err()
					default:
						return r.id * 2, nil
					}
				}
			},
			// Release
			func(r *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					released = true
					return nil, r.Close()
				}
			},
		)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := result(ctx)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.True(t, released)
		assert.True(t, resource.IsClosed())
	})
}

func TestWithResource(t *testing.T) {
	t.Run("reusable resource manager - successful operations", func(t *testing.T) {
		resource := &mockResource{id: 42}
		createCount := 0
		releaseCount := 0

		withResource := WithResource[int](
			// onCreate
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					createCount++
					return resource, nil
				}
			},
			// onRelease
			func(r *mockResource) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					releaseCount++
					return nil, r.Close()
				}
			},
		)

		// First operation
		operation1 := withResource(func(r *mockResource) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return r.id * 2, nil
			}
		})

		result1, err1 := operation1(context.Background())
		assert.NoError(t, err1)
		assert.Equal(t, 84, result1)
		assert.Equal(t, 1, createCount)
		assert.Equal(t, 1, releaseCount)

		// Reset for second operation
		resource.closed = false

		// Second operation with same resource manager
		operation2 := withResource(func(r *mockResource) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return r.id + 10, nil
			}
		})

		result2, err2 := operation2(context.Background())
		assert.NoError(t, err2)
		assert.Equal(t, 52, result2)
		assert.Equal(t, 2, createCount)
		assert.Equal(t, 2, releaseCount)
	})

	t.Run("resource manager with failing operation", func(t *testing.T) {
		resource := &mockResource{id: 42}
		releaseCount := 0
		opErr := errors.New("operation failed")

		withResource := WithResource[int](
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource, nil
				}
			},
			func(r *mockResource) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					releaseCount++
					return nil, r.Close()
				}
			},
		)

		operation := withResource(func(r *mockResource) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return 0, opErr
			}
		})

		_, err := operation(context.Background())
		assert.Error(t, err)
		assert.Equal(t, opErr, err)
		assert.Equal(t, 1, releaseCount)
		assert.True(t, resource.IsClosed())
	})

	t.Run("nested resource managers", func(t *testing.T) {
		resource1 := &mockResource{id: 1}
		resource2 := &mockResource{id: 2}

		withResource1 := WithResource[int](
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource1, nil
				}
			},
			func(r *mockResource) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					return nil, r.Close()
				}
			},
		)

		withResource2 := WithResource[int](
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return resource2, nil
				}
			},
			func(r *mockResource) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					return nil, r.Close()
				}
			},
		)

		// Nest the resource managers
		operation := withResource1(func(r1 *mockResource) ReaderResult[int] {
			return withResource2(func(r2 *mockResource) ReaderResult[int] {
				return func(ctx context.Context) (int, error) {
					return r1.id + r2.id, nil
				}
			})
		})

		result, err := operation(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 3, result)
		assert.True(t, resource1.IsClosed())
		assert.True(t, resource2.IsClosed())
	})
}

func TestWithCloser(t *testing.T) {
	t.Run("successful operation with io.Closer", func(t *testing.T) {
		resource := &mockCloser{mockResource: &mockResource{id: 100}}

		withCloser := WithCloser[string](
			func() ReaderResult[*mockCloser] {
				return func(ctx context.Context) (*mockCloser, error) {
					return resource, nil
				}
			},
		)

		operation := withCloser(func(r *mockCloser) ReaderResult[string] {
			return func(ctx context.Context) (string, error) {
				return "success", nil
			}
		})

		result, err := operation(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.True(t, resource.IsClosed())
	})

	t.Run("operation fails but closer still called", func(t *testing.T) {
		resource := &mockCloser{mockResource: &mockResource{id: 100}}
		opErr := errors.New("operation failed")

		withCloser := WithCloser[string](
			func() ReaderResult[*mockCloser] {
				return func(ctx context.Context) (*mockCloser, error) {
					return resource, nil
				}
			},
		)

		operation := withCloser(func(r *mockCloser) ReaderResult[string] {
			return func(ctx context.Context) (string, error) {
				return "", opErr
			}
		})

		_, err := operation(context.Background())
		assert.Error(t, err)
		assert.Equal(t, opErr, err)
		assert.True(t, resource.IsClosed())
	})

	t.Run("closer fails", func(t *testing.T) {
		closeErr := errors.New("close failed")
		resource := &mockCloser{mockResource: &mockResource{id: 100, closeErr: closeErr}}

		withCloser := WithCloser[string](
			func() ReaderResult[*mockCloser] {
				return func(ctx context.Context) (*mockCloser, error) {
					return resource, nil
				}
			},
		)

		operation := withCloser(func(r *mockCloser) ReaderResult[string] {
			return func(ctx context.Context) (string, error) {
				return "success", nil
			}
		})

		_, err := operation(context.Background())
		assert.Error(t, err)
		assert.Equal(t, closeErr, err)
		assert.True(t, resource.IsClosed())
	})

	t.Run("with strings.Reader (real io.Closer)", func(t *testing.T) {
		content := "Hello, World!"

		withReader := WithCloser[string](
			func() ReaderResult[io.ReadCloser] {
				return func(ctx context.Context) (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader(content)), nil
				}
			},
		)

		operation := withReader(func(r io.ReadCloser) ReaderResult[string] {
			return func(ctx context.Context) (string, error) {
				data, err := io.ReadAll(r)
				return string(data), err
			}
		})

		result, err := operation(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("multiple operations with same closer", func(t *testing.T) {
		createCount := 0

		withCloser := WithCloser[int](
			func() ReaderResult[*mockCloser] {
				return func(ctx context.Context) (*mockCloser, error) {
					createCount++
					return &mockCloser{mockResource: &mockResource{id: createCount}}, nil
				}
			},
		)

		// First operation
		op1 := withCloser(func(r *mockCloser) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return r.id * 10, nil
			}
		})

		result1, err1 := op1(context.Background())
		assert.NoError(t, err1)
		assert.Equal(t, 10, result1)

		// Second operation
		op2 := withCloser(func(r *mockCloser) ReaderResult[int] {
			return func(ctx context.Context) (int, error) {
				return r.id * 20, nil
			}
		})

		result2, err2 := op2(context.Background())
		assert.NoError(t, err2)
		assert.Equal(t, 40, result2)

		assert.Equal(t, 2, createCount)
	})
}

func TestOnClose(t *testing.T) {
	t.Run("onClose helper function", func(t *testing.T) {
		resource := &mockCloser{mockResource: &mockResource{id: 1}}

		closeFunc := onClose(resource)
		_, err := closeFunc(context.Background())

		assert.NoError(t, err)
		assert.True(t, resource.IsClosed())
	})

	t.Run("onClose with error", func(t *testing.T) {
		closeErr := errors.New("close error")
		resource := &mockCloser{mockResource: &mockResource{id: 1, closeErr: closeErr}}

		closeFunc := onClose(resource)
		_, err := closeFunc(context.Background())

		assert.Error(t, err)
		assert.Equal(t, closeErr, err)
		assert.True(t, resource.IsClosed())
	})
}

// Integration test combining multiple bracket patterns
func TestBracketIntegration(t *testing.T) {
	t.Run("complex resource management scenario", func(t *testing.T) {
		// Simulate a scenario with multiple resources
		db := &mockResource{id: 1}
		cache := &mockResource{id: 2}
		logger := &mockResource{id: 3}

		result := Bracket(
			// Acquire DB
			func() ReaderResult[*mockResource] {
				return func(ctx context.Context) (*mockResource, error) {
					return db, nil
				}
			},
			// Use DB to get cache and logger
			func(dbRes *mockResource) ReaderResult[int] {
				return Bracket(
					// Acquire cache
					func() ReaderResult[*mockResource] {
						return func(ctx context.Context) (*mockResource, error) {
							return cache, nil
						}
					},
					// Use cache to get logger
					func(cacheRes *mockResource) ReaderResult[int] {
						return Bracket(
							// Acquire logger
							func() ReaderResult[*mockResource] {
								return func(ctx context.Context) (*mockResource, error) {
									return logger, nil
								}
							},
							// Use all resources
							func(logRes *mockResource) ReaderResult[int] {
								return func(ctx context.Context) (int, error) {
									return dbRes.id + cacheRes.id + logRes.id, nil
								}
							},
							// Release logger
							func(logRes *mockResource, result int, err error) ReaderResult[any] {
								return func(ctx context.Context) (any, error) {
									return nil, logRes.Close()
								}
							},
						)
					},
					// Release cache
					func(cacheRes *mockResource, result int, err error) ReaderResult[any] {
						return func(ctx context.Context) (any, error) {
							return nil, cacheRes.Close()
						}
					},
				)
			},
			// Release DB
			func(dbRes *mockResource, result int, err error) ReaderResult[any] {
				return func(ctx context.Context) (any, error) {
					return nil, dbRes.Close()
				}
			},
		)

		value, err := result(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 6, value) // 1 + 2 + 3
		assert.True(t, db.IsClosed())
		assert.True(t, cache.IsClosed())
		assert.True(t, logger.IsClosed())
	})
}
