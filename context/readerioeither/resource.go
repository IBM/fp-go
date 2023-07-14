package readerioeither

import (
	F "github.com/ibm/fp-go/function"
	RIE "github.com/ibm/fp-go/readerioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[R, A any](onCreate ReaderIOEither[R], onRelease func(R) ReaderIOEither[any]) func(func(R) ReaderIOEither[A]) ReaderIOEither[A] {
	// wraps the callback functions with a context check
	return F.Flow2(
		F.Bind2nd(F.Flow2[func(R) ReaderIOEither[A], func(ReaderIOEither[A]) ReaderIOEither[A], R, ReaderIOEither[A], ReaderIOEither[A]], WithContext[A]),
		RIE.WithResource[ReaderIOEither[A]](WithContext(onCreate), onRelease),
	)
}
