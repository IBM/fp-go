package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/reader"
	RR "github.com/IBM/fp-go/v2/readerresult"
)

//go:inline
func SequenceReader[R, A any](ma ReaderResult[Reader[R, A]]) reader.Kleisli[context.Context, R, Result[A]] {
	return RR.SequenceReader(ma)
}

func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderResult[A]) Kleisli[R, B] {
	return RR.TraverseReader[context.Context](f)
}
