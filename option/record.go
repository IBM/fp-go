package option

import (
	F "github.com/IBM/fp-go/function"
	RR "github.com/IBM/fp-go/internal/record"
)

// TraverseRecord transforms a record of options into an option of a record
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(A) Option[B]) func(GA) Option[GB] {
	return RR.Traverse[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseRecord transforms a record of options into an option of a record
func TraverseRecord[K comparable, A, B any](f func(A) Option[B]) func(map[K]A) Option[map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Option[A], K comparable, A any](ma GOA) Option[GA] {
	return TraverseRecordG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]Option[A]) Option[map[K]A] {
	return SequenceRecordG[map[K]A](ma)
}
