package readerio

import (
	"context"

	"github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
)

//go:inline
func SequenceReader[R, A any](ma ReaderIO[Reader[R, A]]) Reader[R, ReaderIO[A]] {
	return RIO.SequenceReader(ma)
}

//go:inline
func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderIO[A]) Kleisli[R, B] {
	return RIO.TraverseReader[context.Context, R](f)
}
