package state

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

//go:inline
func Do[ST, A any](
	empty A,
) State[ST, A] {
	return Of[ST](empty)
}

//go:inline
func Bind[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[ST, S1, T],
) Operator[ST, S1, S2] {
	return C.Bind(
		Chain[ST, Kleisli[ST, S1, S2], S1, S2],
		Map[ST, func(T) S2, T, S2],
		setter,
		f,
	)
}

//go:inline
func Let[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[ST, S1, S2] {
	return F.Let(
		Map[ST, func(S1) S2, S1, S2],
		key,
		f,
	)
}

//go:inline
func LetTo[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[ST, S1, S2] {
	return F.LetTo(
		Map[ST, func(S1) S2, S1, S2],
		key,
		b,
	)
}

//go:inline
func BindTo[ST, S1, T any](
	setter func(T) S1,
) Operator[ST, T, S1] {
	return C.BindTo(
		Map[ST, func(T) S1, T, S1],
		setter,
	)
}

//go:inline
func ApS[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa State[ST, T],
) Operator[ST, S1, S2] {
	return A.ApS(
		Ap[S2, ST, T],
		Map[ST, func(S1) func(T) S2, S1, func(T) S2],
		setter,
		fa,
	)
}

//go:inline
func ApSL[ST, S, T any](
	lens Lens[S, T],
	fa State[ST, T],
) Endomorphism[State[ST, S]] {
	return ApS(lens.Set, fa)
}

//go:inline
func BindL[ST, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, T, T],
) Endomorphism[State[ST, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

//go:inline
func LetL[ST, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[State[ST, S]] {
	return Let[ST](lens.Set, function.Flow2(lens.Get, f))
}

//go:inline
func LetToL[ST, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[State[ST, S]] {
	return LetTo[ST](lens.Set, b)
}
