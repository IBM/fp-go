package result

import "github.com/IBM/fp-go/v2/either"

func TraversableIter[A, B any]() Traversable[A, B, Seq[A], Seq[B]] {
	return either.TraversableIter[error, A, B]()
}
