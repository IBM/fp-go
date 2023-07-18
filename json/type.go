package json

import (
	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
)

func ToTypeE[A any](src any) E.Either[error, A] {
	return F.Pipe2(
		src,
		Marshal[any],
		E.Chain(Unmarshal[A]),
	)
}

func ToTypeO[A any](src any) O.Option[A] {
	return F.Pipe1(
		ToTypeE[A](src),
		E.ToOption[error, A](),
	)
}
