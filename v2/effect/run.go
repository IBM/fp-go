package effect

import (
	"context"

	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
	"github.com/IBM/fp-go/v2/result"
)

func Provide[C, A any](c C) func(Effect[C, A]) ReaderIOResult[A] {
	return readerreaderioresult.Read[A](c)
}

func RunSync[A any](fa ReaderIOResult[A]) readerresult.ReaderResult[A] {
	return func(ctx context.Context) (A, error) {
		return result.Unwrap(fa(ctx)())
	}
}
