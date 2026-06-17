package effect

import "github.com/IBM/fp-go/v2/context/readerreaderioresult"

// ContramapMemoize memoizes an Effect by deriving a comparable cache key from the
// effect environment.
//
// The returned operator transforms an Effect[C, A] into one that caches the thunk
// produced for each derived key. Calls whose environments map to the same key share
// the same memoized computation and result.
//
// This is useful when the environment is not itself comparable or when only part of
// the environment should determine cache identity.
//
// Correctness Contract:
//   - kf must faithfully partition the environment for the effect: if kf(c1) == kf(c2), then running the
//     effect with c1 and c2 must produce the same result.
//   - The thunk computed on the first run for a given key is cached and reused for all subsequent runs
//     whose environments map to the same key — regardless of the actual environment passed.
//   - Violating this contract causes silent incorrect behavior: no error is returned, but callers
//     receive stale results.
//
// Type Parameters:
//   - A: The result type of the effect
//   - C: The environment type required by the effect
//   - K: The comparable cache key type derived from C
//
// Parameters:
//   - kf: Extracts the cache key from the environment
//
// Returns:
//   - Operator[C, A, A]: An operator that memoizes by the derived key
func ContramapMemoize[A, C any, K comparable](kf Reader[C, K]) Operator[C, A, A] {
	return readerreaderioresult.ContramapMemoize[A](kf)
}

// Memoize memoizes an Effect using the environment value itself as the cache key.
//
// For each distinct environment value, the thunk produced by rdr is memoized and
// reused on subsequent runs with the same environment. The environment type must be
// comparable so it can be used directly as the cache key.
//
// Type Parameters:
//   - C: The comparable environment type used as the cache key
//   - A: The result type of the effect
//
// Parameters:
//   - rdr: The effect to memoize
//
// Returns:
//   - Effect[C, A]: A memoized effect keyed by the environment
func Memoize[C comparable, A any](rdr Effect[C, A]) Effect[C, A] {
	return readerreaderioresult.Memoize(rdr)
}
