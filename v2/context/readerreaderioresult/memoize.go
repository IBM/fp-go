package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
)

// ContramapMemoize memoizes a ReaderReaderIOResult by deriving a comparable cache key
// from the outer environment.
//
// The returned operator transforms a ReaderReaderIOResult[C, A] into one that caches the
// ReaderIOResult[context.Context, A] produced for each derived key. Calls whose outer
// environment values map to the same key share the same memoized inner computation.
//
// This is useful when the outer environment is not itself comparable or when only part
// of the environment should determine cache identity.
//
// Type Parameters:
//   - A: The result type of the computation
//   - C: The outer environment type
//   - K: The comparable cache key type derived from C
//
// Parameters:
//   - kf: Extracts the cache key from the outer environment
//
// Returns:
//   - Operator[C, A, A]: An operator that memoizes by the derived key
func ContramapMemoize[A, C any, K comparable](kf Reader[C, K]) Operator[C, A, A] {
	memoize := F.ContramapMemoize[ReaderIOResult[context.Context, A]](kf)

	return func(rdr ReaderReaderIOResult[C, A]) ReaderReaderIOResult[C, A] {
		return memoize(F.Flow2(
			rdr,
			readerioresult.Memoize,
		))
	}
}

// Memoize memoizes a ReaderReaderIOResult using the outer environment value itself as
// the cache key.
//
// For each distinct outer environment value, the computation produced by rdr is
// memoized and reused on subsequent calls with the same environment. The outer
// environment type must be comparable so it can be used directly as the cache key.
//
// Type Parameters:
//   - C: The comparable outer environment type used as the cache key
//   - A: The result type of the computation
//
// Parameters:
//   - rdr: The computation to memoize
//
// Returns:
//   - ReaderReaderIOResult[C, A]: A memoized computation keyed by the outer environment
func Memoize[C comparable, A any](rdr ReaderReaderIOResult[C, A]) ReaderReaderIOResult[C, A] {
	return ContramapMemoize[A](F.Identity[C])(rdr)
}
