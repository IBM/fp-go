package readerioresult

import (
	"github.com/IBM/fp-go/v2/reader"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

//go:inline
func Sequence[R1, R2, A any](ma ReaderIOResult[R2, ReaderIOResult[R1, A]]) Kleisli[R2, R1, A] {
	return RIOE.Sequence(ma)
}

//go:inline
func SequenceReader[R1, R2, A any](ma ReaderIOResult[R2, Reader[R1, A]]) Kleisli[R2, R1, A] {
	return RIOE.SequenceReader(ma)
}

//go:inline
func SequenceReaderIO[R1, R2, A any](ma ReaderIOResult[R2, ReaderIO[R1, A]]) Kleisli[R2, R1, A] {
	return RIOE.SequenceReaderIO(ma)
}

//go:inline
func SequenceReaderEither[R1, R2, A any](ma ReaderIOResult[R2, ReaderResult[R1, A]]) Kleisli[R2, R1, A] {
	return RIOE.SequenceReaderEither(ma)
}

//go:inline
func SequenceReaderResult[R1, R2, A any](ma ReaderIOResult[R2, ReaderResult[R1, A]]) Kleisli[R2, R1, A] {
	return RIOE.SequenceReaderEither(ma)
}

//go:inline
func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return RIOE.Traverse[R2](f)
}

func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderIOResult[R2, A]) Kleisli[R2, R1, B] {
	return RIOE.TraverseReader[R2, R1, error](f)
}
