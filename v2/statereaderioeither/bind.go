package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

//go:inline
func Do[ST, R, E, A any](
	empty A,
) StateReaderIOEither[ST, R, E, A] {
	return Of[ST, R, E](empty)
}

//go:inline
func Bind[ST, R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[ST, R, E, S1, T],
) Operator[ST, R, E, S1, S2] {
	return C.Bind(
		Chain[ST, R, E, S1, S2],
		Map[ST, R, E, T, S2],
		setter,
		f,
	)
}

//go:inline
func Let[ST, R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[ST, R, E, S1, S2] {
	return F.Let(
		Map[ST, R, E, S1, S2],
		key,
		f,
	)
}

//go:inline
func LetTo[ST, R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[ST, R, E, S1, S2] {
	return F.LetTo(
		Map[ST, R, E, S1, S2],
		key,
		b,
	)
}

//go:inline
func BindTo[ST, R, E, S1, T any](
	setter func(T) S1,
) Operator[ST, R, E, T, S1] {
	return C.BindTo(
		Map[ST, R, E, T, S1],
		setter,
	)
}

//go:inline
func ApS[ST, R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa StateReaderIOEither[ST, R, E, T],
) Operator[ST, R, E, S1, S2] {
	return A.ApS(
		Ap[S2, ST, R, E, T],
		Map[ST, R, E, S1, func(T) S2],
		setter,
		fa,
	)
}

//go:inline
func ApSL[ST, R, E, S, T any](
	lens Lens[S, T],
	fa StateReaderIOEither[ST, R, E, T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return ApS(lens.Set, fa)
}

//go:inline
func BindL[ST, R, E, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, R, E, T, T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

//go:inline
func LetL[ST, R, E, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return Let[ST, R, E](lens.Set, function.Flow2(lens.Get, f))
}

//go:inline
func LetToL[ST, R, E, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return LetTo[ST, R, E](lens.Set, b)
}
