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

package generic

import (
	"sync"

	L "github.com/IBM/fp-go/v2/internal/lazy"
)

// Memoize converts a unary function into a unary function that caches the value depending on the parameter
func Memoize[F ~func(K) T, K comparable, T any](f F) F {
	return ContramapMemoize[F](func(k K) K { return k })(f)
}

// ContramapMemoize converts a unary function into a unary function that caches the value depending on the parameter
func ContramapMemoize[F ~func(A) T, KF func(A) K, A any, K comparable, T any](kf KF) func(F) F {
	return CacheCallback[func(F) F](kf, getOrCreate[K, T]())
}

// getOrCreate is a naive implementation of a cache, without bounds
func getOrCreate[K comparable, T any]() func(K, func() func() T) func() T {
	cache := make(map[K]func() T)
	var l sync.Mutex

	return func(k K, cb func() func() T) func() T {
		// only lock to access a lazy accessor to the value
		l.Lock()
		existing, ok := cache[k]
		if !ok {
			existing = cb()
			cache[k] = existing
		}
		l.Unlock()
		// compute the value outside of the lock
		return existing
	}
}

// SingleElementCache is a cache with a capacity of a single element
func SingleElementCache[
	LLT ~func() LT, // generator of the generator
	K comparable, // key into the cache
	LT ~func() T, // generator of a value
	T any, // the cached data type
]() func(K, LLT) LT {
	var l sync.Mutex

	var key K
	var value LT
	hasKey := false

	return func(k K, gen LLT) LT {
		l.Lock()

		existing := value
		if !hasKey || key != k {
			existing = gen()
			// update state
			key = k
			value = existing
			hasKey = true
		}

		l.Unlock()

		return existing
	}
}

// CacheCallback converts a unary function into a unary function that caches the value depending on the parameter
func CacheCallback[
	EM ~func(F) F, // endomorphism of the function
	LLT ~func() LT, // generator of the generator
	LT ~func() T, // generator of a value
	F ~func(A) T, // function to actually cache
	KF func(A) K, // extracts the cache key from the input
	C ~func(K, LLT) LT, // the cache callback function
	A any, K comparable, T any](kf KF, getOrCreate C) EM {
	return func(f F) F {
		return func(a A) T {
			// cache entry
			return getOrCreate(kf(a), func() LT {
				return L.Memoize[LT](func() T {
					return f(a)
				})
			})()
		}
	}
}
