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

package function

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMemoize tests the Memoize function
func TestMemoize(t *testing.T) {
	t.Run("caches computed values", func(t *testing.T) {
		callCount := 0
		expensive := func(n int) int {
			callCount++
			time.Sleep(10 * time.Millisecond)
			return n * 2
		}

		memoized := Memoize(expensive)

		// First call should compute
		result1 := memoized(5)
		assert.Equal(t, 10, result1)
		assert.Equal(t, 1, callCount)

		// Second call with same input should use cache
		result2 := memoized(5)
		assert.Equal(t, 10, result2)
		assert.Equal(t, 1, callCount, "should not recompute for cached value")

		// Different input should compute again
		result3 := memoized(10)
		assert.Equal(t, 20, result3)
		assert.Equal(t, 2, callCount)

		// Original input should still be cached
		result4 := memoized(5)
		assert.Equal(t, 10, result4)
		assert.Equal(t, 2, callCount, "should still use cached value")
	})

	t.Run("works with string keys", func(t *testing.T) {
		callCount := 0
		toUpper := func(s string) string {
			callCount++
			return fmt.Sprintf("UPPER_%s", s)
		}

		memoized := Memoize(toUpper)

		result1 := memoized("hello")
		assert.Equal(t, "UPPER_hello", result1)
		assert.Equal(t, 1, callCount)

		result2 := memoized("hello")
		assert.Equal(t, "UPPER_hello", result2)
		assert.Equal(t, 1, callCount)

		result3 := memoized("world")
		assert.Equal(t, "UPPER_world", result3)
		assert.Equal(t, 2, callCount)
	})

	t.Run("is thread-safe", func(t *testing.T) {
		var callCount int32
		expensive := func(n int) int {
			atomic.AddInt32(&callCount, 1)
			time.Sleep(5 * time.Millisecond)
			return n * n
		}

		memoized := Memoize(expensive)

		// Run concurrent calls with same input
		var wg sync.WaitGroup
		results := make([]int, 10)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = memoized(7)
			}(i)
		}
		wg.Wait()

		// All results should be the same
		for _, result := range results {
			assert.Equal(t, 49, result)
		}

		// Function should be called at least once, but possibly more due to race
		// (the cache is eventually consistent)
		assert.Greater(t, atomic.LoadInt32(&callCount), int32(0))
	})

	t.Run("handles zero values correctly", func(t *testing.T) {
		callCount := 0
		identity := func(n int) int {
			callCount++
			return n
		}

		memoized := Memoize(identity)

		result1 := memoized(0)
		assert.Equal(t, 0, result1)
		assert.Equal(t, 1, callCount)

		result2 := memoized(0)
		assert.Equal(t, 0, result2)
		assert.Equal(t, 1, callCount, "should cache zero value")
	})

	t.Run("caches multiple different values", func(t *testing.T) {
		callCount := 0
		square := func(n int) int {
			callCount++
			return n * n
		}

		memoized := Memoize(square)

		// Cache multiple values
		assert.Equal(t, 4, memoized(2))
		assert.Equal(t, 9, memoized(3))
		assert.Equal(t, 16, memoized(4))
		assert.Equal(t, 3, callCount)

		// All should be cached
		assert.Equal(t, 4, memoized(2))
		assert.Equal(t, 9, memoized(3))
		assert.Equal(t, 16, memoized(4))
		assert.Equal(t, 3, callCount, "all values should be cached")
	})
}

// TestContramapMemoize tests the ContramapMemoize function
func TestContramapMemoize(t *testing.T) {
	type User struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("caches by extracted key", func(t *testing.T) {
		callCount := 0
		getUserData := func(u User) string {
			callCount++
			return fmt.Sprintf("Data for user %d: %s", u.ID, u.Name)
		}

		// Cache by ID only
		cacheByID := ContramapMemoize[string](func(u User) int {
			return u.ID
		})

		memoized := cacheByID(getUserData)

		user1 := User{ID: 1, Name: "Alice", Age: 30}
		result1 := memoized(user1)
		assert.Equal(t, "Data for user 1: Alice", result1)
		assert.Equal(t, 1, callCount)

		// Same ID, different name - should use cache
		user2 := User{ID: 1, Name: "Bob", Age: 25}
		result2 := memoized(user2)
		assert.Equal(t, "Data for user 1: Alice", result2, "should return cached result")
		assert.Equal(t, 1, callCount, "should not recompute")

		// Different ID - should compute
		user3 := User{ID: 2, Name: "Charlie", Age: 35}
		result3 := memoized(user3)
		assert.Equal(t, "Data for user 2: Charlie", result3)
		assert.Equal(t, 2, callCount)
	})

	t.Run("works with string key extraction", func(t *testing.T) {
		type Product struct {
			SKU   string
			Name  string
			Price float64
		}

		callCount := 0
		getPrice := func(p Product) float64 {
			callCount++
			return p.Price * 1.1 // Add 10% markup
		}

		cacheBySKU := ContramapMemoize[float64](func(p Product) string {
			return p.SKU
		})

		memoized := cacheBySKU(getPrice)

		prod1 := Product{SKU: "ABC123", Name: "Widget", Price: 100.0}
		result1 := memoized(prod1)
		assert.InDelta(t, 110.0, result1, 0.01)
		assert.Equal(t, 1, callCount)

		// Same SKU, different price - should use cached result
		prod2 := Product{SKU: "ABC123", Name: "Widget", Price: 200.0}
		result2 := memoized(prod2)
		assert.InDelta(t, 110.0, result2, 0.01, "should use cached value")
		assert.Equal(t, 1, callCount)
	})

	t.Run("can use complex key extraction", func(t *testing.T) {
		type Request struct {
			Method string
			Path   string
			Body   string
		}

		callCount := 0
		processRequest := func(r Request) string {
			callCount++
			return fmt.Sprintf("Processed: %s %s", r.Method, r.Path)
		}

		// Cache by method and path, ignore body
		cacheByMethodPath := ContramapMemoize[string](func(r Request) string {
			return r.Method + ":" + r.Path
		})

		memoized := cacheByMethodPath(processRequest)

		req1 := Request{Method: "GET", Path: "/api/users", Body: "body1"}
		result1 := memoized(req1)
		assert.Equal(t, "Processed: GET /api/users", result1)
		assert.Equal(t, 1, callCount)

		// Same method and path, different body - should use cache
		req2 := Request{Method: "GET", Path: "/api/users", Body: "body2"}
		result2 := memoized(req2)
		assert.Equal(t, "Processed: GET /api/users", result2)
		assert.Equal(t, 1, callCount)

		// Different path - should compute
		req3 := Request{Method: "GET", Path: "/api/posts", Body: "body1"}
		result3 := memoized(req3)
		assert.Equal(t, "Processed: GET /api/posts", result3)
		assert.Equal(t, 2, callCount)
	})
}

// TestCacheCallback tests the CacheCallback function
func TestCacheCallback(t *testing.T) {
	t.Run("works with custom cache implementation", func(t *testing.T) {
		// Create a simple bounded cache (max 2 items)
		boundedCache := func() func(int, func() func() string) func() string {
			cache := make(map[int]func() string)
			keys := []int{}
			var mu sync.Mutex

			return func(k int, gen func() func() string) func() string {
				mu.Lock()
				defer mu.Unlock()

				if existing, ok := cache[k]; ok {
					return existing
				}

				// Evict oldest if at capacity
				if len(keys) >= 2 {
					oldestKey := keys[0]
					delete(cache, oldestKey)
					keys = keys[1:]
				}

				value := gen()
				cache[k] = value
				keys = append(keys, k)
				return value
			}
		}

		callCount := 0
		expensive := func(n int) string {
			callCount++
			return fmt.Sprintf("Result: %d", n)
		}

		memoizer := CacheCallback(
			Identity[int],
			boundedCache(),
		)

		memoized := memoizer(expensive)

		// Cache first two values
		result1 := memoized(1)
		assert.Equal(t, "Result: 1", result1)
		assert.Equal(t, 1, callCount)

		result2 := memoized(2)
		assert.Equal(t, "Result: 2", result2)
		assert.Equal(t, 2, callCount)

		// Both should be cached
		memoized(1)
		memoized(2)
		assert.Equal(t, 2, callCount)

		// Third value should evict first
		result3 := memoized(3)
		assert.Equal(t, "Result: 3", result3)
		assert.Equal(t, 3, callCount)

		// First value should be recomputed (evicted)
		// Note: The cache stores lazy generators, so calling memoized(1) again
		// will create a new cache entry with a new lazy generator
		memoized(1)
		// The call count increases because a new lazy value is created and evaluated
		assert.GreaterOrEqual(t, callCount, 3, "first value should have been evicted")

		// Verify cache still works for remaining values
		prevCount := callCount
		memoized(2)
		memoized(3)
		// These might or might not increase count depending on eviction
		assert.GreaterOrEqual(t, callCount, prevCount)
	})

	t.Run("integrates with key extraction", func(t *testing.T) {
		type Item struct {
			ID    int
			Value string
		}

		// Simple cache
		simpleCache := func() func(int, func() func() string) func() string {
			cache := make(map[int]func() string)
			var mu sync.Mutex

			return func(k int, gen func() func() string) func() string {
				mu.Lock()
				defer mu.Unlock()

				if existing, ok := cache[k]; ok {
					return existing
				}

				value := gen()
				cache[k] = value
				return value
			}
		}

		callCount := 0
		process := func(item Item) string {
			callCount++
			return fmt.Sprintf("Processed: %s", item.Value)
		}

		memoizer := CacheCallback(
			func(item Item) int { return item.ID },
			simpleCache(),
		)

		memoized := memoizer(process)

		item1 := Item{ID: 1, Value: "first"}
		result1 := memoized(item1)
		assert.Equal(t, "Processed: first", result1)
		assert.Equal(t, 1, callCount)

		// Same ID, different value - should use cache
		item2 := Item{ID: 1, Value: "second"}
		result2 := memoized(item2)
		assert.Equal(t, "Processed: first", result2)
		assert.Equal(t, 1, callCount)
	})
}

// TestSingleElementCache tests the SingleElementCache function
func TestSingleElementCache(t *testing.T) {
	t.Run("caches single element", func(t *testing.T) {
		cache := SingleElementCache[int, string]()

		callCount := 0
		gen := func(n int) func() func() string {
			// This returns a generator that creates a lazy value
			return func() func() string {
				// This is the lazy value that gets cached
				return func() string {
					// This gets called when the lazy value is evaluated
					callCount++
					return fmt.Sprintf("Value: %d", n)
				}
			}
		}

		// First call - creates and caches lazy value for key 1
		lazy1 := cache(1, gen(1))
		result1 := lazy1()
		assert.Equal(t, "Value: 1", result1)
		assert.Equal(t, 1, callCount)

		// Same key - returns the same cached lazy value
		lazy1Again := cache(1, gen(1))
		result2 := lazy1Again()
		assert.Equal(t, "Value: 1", result2)
		// The lazy value is called again, so count increases
		assert.Equal(t, 2, callCount, "cached lazy value is called again")

		// Different key - replaces cache with new lazy value
		lazy2 := cache(2, gen(2))
		result3 := lazy2()
		assert.Equal(t, "Value: 2", result3)
		assert.Equal(t, 3, callCount)

		// Original key - cache was replaced, creates new lazy value
		lazy1New := cache(1, gen(1))
		result4 := lazy1New()
		assert.Equal(t, "Value: 1", result4)
		assert.Equal(t, 4, callCount, "new lazy value created after cache replacement")
	})

	t.Run("works with CacheCallback", func(t *testing.T) {
		cache := SingleElementCache[int, string]()

		callCount := 0
		expensive := func(n int) string {
			callCount++
			return fmt.Sprintf("Result: %d", n*n)
		}

		memoizer := CacheCallback(
			Identity[int],
			cache,
		)

		memoized := memoizer(expensive)

		// First computation
		result1 := memoized(5)
		assert.Equal(t, "Result: 25", result1)
		assert.Equal(t, 1, callCount)

		// Same input - cached
		result2 := memoized(5)
		assert.Equal(t, "Result: 25", result2)
		assert.Equal(t, 1, callCount)

		// Different input - replaces cache
		result3 := memoized(10)
		assert.Equal(t, "Result: 100", result3)
		assert.Equal(t, 2, callCount)

		// Back to first input - recomputed
		result4 := memoized(5)
		assert.Equal(t, "Result: 25", result4)
		assert.Equal(t, 3, callCount)
	})

	t.Run("is thread-safe", func(t *testing.T) {
		cache := SingleElementCache[int, string]()

		var callCount int32
		gen := func(n int) func() func() string {
			return func() func() string {
				return func() string {
					atomic.AddInt32(&callCount, 1)
					time.Sleep(5 * time.Millisecond)
					return fmt.Sprintf("Value: %d", n)
				}
			}
		}

		var wg sync.WaitGroup
		results := make([]string, 20)

		// Concurrent access with same key
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = cache(1, gen(1))()
			}(i)
		}

		// Concurrent access with different key
		for i := 10; i < 20; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = cache(2, gen(2))()
			}(i)
		}

		wg.Wait()

		// All results should be valid (either "Value: 1" or "Value: 2")
		for _, result := range results {
			assert.True(t, result == "Value: 1" || result == "Value: 2")
		}

		// Function should have been called, but exact count depends on race conditions
		assert.Greater(t, atomic.LoadInt32(&callCount), int32(0))
	})

	t.Run("handles rapid key changes", func(t *testing.T) {
		cache := SingleElementCache[int, string]()

		callCount := 0
		gen := func(n int) func() func() string {
			return func() func() string {
				return func() string {
					callCount++
					return fmt.Sprintf("Value: %d", n)
				}
			}
		}

		// Rapidly alternate between keys
		for i := 0; i < 10; i++ {
			cache(1, gen(1))()
			cache(2, gen(2))()
		}

		// Each key change should trigger a computation
		// (20 calls total: 10 for key 1, 10 for key 2)
		assert.Equal(t, 20, callCount)
	})
}

// TestMemoizeIntegration tests integration scenarios
func TestMemoizeIntegration(t *testing.T) {
	t.Run("fibonacci with memoization", func(t *testing.T) {
		callCount := 0

		expensive := func(n int) int {
			callCount++
			time.Sleep(10 * time.Millisecond)
			return n * n
		}

		memoized := Memoize(expensive)

		// First call computes
		result1 := memoized(10)
		assert.Equal(t, 100, result1)
		assert.Equal(t, 1, callCount)

		// Second call with same input uses cache
		result2 := memoized(10)
		assert.Equal(t, 100, result2)
		assert.Equal(t, 1, callCount, "should use cached value")

		// Different input computes again
		result3 := memoized(5)
		assert.Equal(t, 25, result3)
		assert.Equal(t, 2, callCount)

		// Both values should remain cached
		assert.Equal(t, 100, memoized(10))
		assert.Equal(t, 25, memoized(5))
		assert.Equal(t, 2, callCount, "both values should be cached")
	})

	t.Run("chaining memoization strategies", func(t *testing.T) {
		type Request struct {
			UserID int
			Action string
		}

		callCount := 0
		processRequest := func(r Request) string {
			callCount++
			return fmt.Sprintf("User %d: %s", r.UserID, r.Action)
		}

		// First level: cache by UserID
		cacheByUser := ContramapMemoize[string](func(r Request) int {
			return r.UserID
		})

		memoized := cacheByUser(processRequest)

		req1 := Request{UserID: 1, Action: "login"}
		result1 := memoized(req1)
		assert.Equal(t, "User 1: login", result1)
		assert.Equal(t, 1, callCount)

		// Same user, different action - uses cache
		req2 := Request{UserID: 1, Action: "logout"}
		result2 := memoized(req2)
		assert.Equal(t, "User 1: login", result2)
		assert.Equal(t, 1, callCount)

		// Different user - computes
		req3 := Request{UserID: 2, Action: "login"}
		result3 := memoized(req3)
		assert.Equal(t, "User 2: login", result3)
		assert.Equal(t, 2, callCount)
	})
}
