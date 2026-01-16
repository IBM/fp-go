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
	G "github.com/IBM/fp-go/v2/function/generic"
)

// Memoize converts a unary function into a memoized version that caches computed values.
//
// Behavior:
//   - On first call with a given input, the function executes and the result is cached
//   - Subsequent calls with the same input return the cached result without re-execution
//   - The cache uses the input parameter directly as the key (must be comparable)
//   - The cache is thread-safe using mutex locks
//   - The cache has no size limit and grows unbounded
//   - Each unique input creates a new cache entry that persists for the lifetime of the memoized function
//
// Implementation Details:
//   - Uses an internal map[K]func()T to store lazy values
//   - The cached value is wrapped in a lazy function to defer computation until needed
//   - Lock is held only to access the cache map, not during value computation
//   - This allows concurrent computations for different keys
//
// Type Parameters:
//   - K: The type of the function parameter, must be comparable (used as cache key)
//   - T: The return type of the function
//
// Parameters:
//   - f: The function to memoize
//
// Returns:
//   - A memoized version of the function that caches results by parameter value
//
// Example:
//
//	// Expensive computation
//	expensiveCalc := func(n int) int {
//	    time.Sleep(100 * time.Millisecond)
//	    return n * n
//	}
//
//	// Memoize to avoid redundant calculations
//	memoized := Memoize(expensiveCalc)
//	result1 := memoized(5) // Takes 100ms, computes and caches 25
//	result2 := memoized(5) // Instant, returns cached 25
//	result3 := memoized(10) // Takes 100ms, computes and caches 100
//
// Note: The cache grows unbounded. For bounded caches, use CacheCallback with a custom cache implementation.
func Memoize[K comparable, T any](f func(K) T) func(K) T {
	return G.Memoize(f)
}

// ContramapMemoize creates a higher-order function that memoizes functions using a custom key extraction strategy.
//
// Behavior:
//   - Allows caching based on a derived key rather than the full input parameter
//   - The key extraction function (kf) determines what constitutes a cache hit
//   - Two inputs that produce the same key will share the same cached result
//   - This enables caching for non-comparable types by extracting comparable keys
//   - The cache is thread-safe and unbounded
//
// Use Cases:
//   - Cache by a subset of struct fields (e.g., User.ID instead of entire User)
//   - Cache by a computed property (e.g., string length, hash value)
//   - Normalize inputs before caching (e.g., lowercase strings, rounded numbers)
//
// Implementation Details:
//   - Internally uses the same caching mechanism as Memoize
//   - The key function is applied to each input before cache lookup
//   - Returns a function transformer that can be applied to any function with matching signature
//
// Type Parameters:
//   - T: The return type of the function to be memoized
//   - A: The input type of the function to be memoized
//   - K: The type of the cache key, must be comparable
//
// Parameters:
//   - kf: A function that extracts a cache key from the input parameter
//
// Returns:
//   - A function that takes a function (A) -> T and returns its memoized version
//
// Example:
//
//	type User struct {
//	    ID   int
//	    Name string
//	    Email string
//	}
//
//	// Cache by user ID only, ignoring other fields
//	cacheByID := ContramapMemoize[string, User, int](func(u User) int {
//	    return u.ID
//	})
//
//	getUserData := func(u User) string {
//	    // Expensive database lookup
//	    return fmt.Sprintf("Data for user %d", u.ID)
//	}
//
//	memoized := cacheByID(getUserData)
//	result1 := memoized(User{ID: 1, Name: "Alice", Email: "a@example.com"}) // Computed
//	result2 := memoized(User{ID: 1, Name: "Bob", Email: "b@example.com"})   // Cached (same ID)
//	result3 := memoized(User{ID: 2, Name: "Alice", Email: "a@example.com"}) // Computed (different ID)
func ContramapMemoize[T, A any, K comparable](kf func(A) K) func(func(A) T) func(A) T {
	return G.ContramapMemoize[func(A) T](kf)
}

// CacheCallback creates a higher-order function that memoizes functions using a custom cache implementation.
//
// Behavior:
//   - Provides complete control over caching strategy through the getOrCreate callback
//   - Separates cache key extraction (kf) from cache storage (getOrCreate)
//   - The getOrCreate function receives a key and a lazy value generator
//   - The cache implementation decides when to store, evict, or retrieve values
//   - Enables advanced caching strategies: LRU, LFU, TTL, bounded size, etc.
//
// How It Works:
//  1. When the memoized function is called with input A:
//  2. The key function (kf) extracts a cache key K from A
//  3. A lazy value generator is created that will compute f(A) when called
//  4. The getOrCreate callback is invoked with the key and lazy generator
//  5. The cache implementation returns a lazy value (either cached or newly created)
//  6. The lazy value is evaluated to produce the final result T
//
// Cache Implementation Contract:
//   - getOrCreate receives: (key K, generator func() func() T)
//   - getOrCreate returns: func() T (a lazy value)
//   - The generator creates a new lazy value when called
//   - The cache should store and return lazy values, not final results
//   - This allows deferred computation and proper lazy evaluation
//
// Type Parameters:
//   - T: The return type of the function to be memoized
//   - A: The input type of the function to be memoized
//   - K: The type of the cache key, must be comparable
//
// Parameters:
//   - kf: A function that extracts a cache key from the input parameter
//   - getOrCreate: A cache implementation that stores and retrieves lazy values
//
// Returns:
//   - A function that takes a function (A) -> T and returns its memoized version
//
// Example:
//
//	// Create a bounded LRU cache (max 100 items)
//	lruCache := func() func(int, func() func() string) func() string {
//	    cache := make(map[int]func() string)
//	    keys := []int{}
//	    var mu sync.Mutex
//	    maxSize := 100
//
//	    return func(k int, gen func() func() string) func() string {
//	        mu.Lock()
//	        defer mu.Unlock()
//
//	        if existing, ok := cache[k]; ok {
//	            return existing // Cache hit
//	        }
//
//	        // Evict oldest if at capacity
//	        if len(keys) >= maxSize {
//	            delete(cache, keys[0])
//	            keys = keys[1:]
//	        }
//
//	        // Create and store new lazy value
//	        value := gen()
//	        cache[k] = value
//	        keys = append(keys, k)
//	        return value
//	    }
//	}
//
//	// Use custom cache with memoization
//	memoizer := CacheCallback[string, int, int](
//	    Identity[int], // Use input as key
//	    lruCache(),
//	)
//
//	expensiveFunc := func(n int) string {
//	    time.Sleep(100 * time.Millisecond)
//	    return fmt.Sprintf("Result: %d", n)
//	}
//
//	memoized := memoizer(expensiveFunc)
//	result := memoized(42) // Computed and cached
//	result = memoized(42)  // Retrieved from cache
//
// See also: SingleElementCache for a simple bounded cache implementation.
func CacheCallback[
	T, A any, K comparable](kf func(A) K, getOrCreate func(K, func() func() T) func() T) func(func(A) T) func(A) T {
	return G.CacheCallback[func(func(A) T) func(A) T](kf, getOrCreate)
}

// SingleElementCache creates a thread-safe cache implementation that stores at most one element.
//
// Behavior:
//   - Stores only the most recently accessed key-value pair
//   - When a new key is accessed, it replaces the previous cached entry
//   - If the same key is accessed again, the cached value is returned
//   - Thread-safe: uses mutex to protect concurrent access
//   - Memory-efficient: constant O(1) space regardless of usage
//
// How It Works:
//  1. Initially, the cache is empty (hasKey = false)
//  2. On first access with key K1:
//     - Calls the generator to create a lazy value
//     - Stores K1 and the lazy value
//     - Returns the lazy value
//  3. On subsequent access with same key K1:
//     - Returns the stored lazy value without calling generator
//  4. On access with different key K2:
//     - Calls the generator to create a new lazy value
//     - Replaces K1 with K2 and updates the stored lazy value
//     - Returns the new lazy value
//  5. If K1 is accessed again, it's treated as a new key (cache miss)
//
// Use Cases:
//   - Sequential processing where the same key is accessed multiple times in a row
//   - Memory-constrained environments where unbounded caches are not feasible
//   - Scenarios where only the most recent computation needs caching
//   - Testing or debugging with controlled cache behavior
//
// Important Notes:
//   - The cache stores the lazy value (func() T), not the computed result
//   - Each time the returned lazy value is called, it may recompute (depends on lazy implementation)
//   - For true result caching, combine with lazy memoization (as done in CacheCallback)
//   - Alternating between two keys will cause constant cache misses
//
// Type Parameters:
//   - K: The type of the cache key, must be comparable
//   - T: The type of the cached value
//
// Returns:
//   - A cache function suitable for use with CacheCallback
//
// Example:
//
//	// Create a single-element cache
//	cache := SingleElementCache[int, string]()
//
//	// Use with CacheCallback
//	memoizer := CacheCallback[string, int, int](
//	    Identity[int],  // Use input as key
//	    cache,
//	)
//
//	expensiveFunc := func(n int) string {
//	    time.Sleep(100 * time.Millisecond)
//	    return fmt.Sprintf("Result: %d", n)
//	}
//
//	memoized := memoizer(expensiveFunc)
//	result1 := memoized(42) // Computed (100ms) and cached
//	result2 := memoized(42) // Instant - returns cached value
//	result3 := memoized(99) // Computed (100ms) - replaces cache entry for 42
//	result4 := memoized(99) // Instant - returns cached value
//	result5 := memoized(42) // Computed (100ms) - cache was replaced, must recompute
//
// Performance Characteristics:
//   - Space: O(1) - stores exactly one key-value pair
//   - Time: O(1) - cache lookup and update are constant time
//   - Best case: Same key accessed repeatedly (100% hit rate)
//   - Worst case: Alternating keys (0% hit rate)
func SingleElementCache[K comparable, T any]() func(K, func() func() T) func() T {
	return G.SingleElementCache[func() func() T, K]()
}
