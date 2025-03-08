package generic

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2025-03-07 23:15:45.4688315 +0100 CET m=+0.002075201

import (
	"context"

	E "github.com/IBM/fp-go/v2/either"
	RE "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// Eitherize0 converts a function with 0 parameters returning a tuple into a function with 0 parameters returning a [GRA]
// The inverse function is [Uneitherize0]
func Eitherize0[GRA ~func(context.Context) GIOA, F ~func(context.Context) (R, error), GIOA ~func() E.Either[error, R], R any](f F) func() GRA {
	return RE.Eitherize0[GRA](f)
}

// Uneitherize0 converts a function with 0 parameters returning a [GRA] into a function with 0 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize0[GRA ~func(context.Context) GIOA, F ~func(context.Context) (R, error), GIOA ~func() E.Either[error, R], R any](f func() GRA) F {
	return func(c context.Context) (R, error) {
		return E.UnwrapError(f()(c)())
	}
}

// Eitherize1 converts a function with 1 parameters returning a tuple into a function with 1 parameters returning a [GRA]
// The inverse function is [Uneitherize1]
func Eitherize1[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0) (R, error), GIOA ~func() E.Either[error, R], T0, R any](f F) func(T0) GRA {
	return RE.Eitherize1[GRA](f)
}

// Uneitherize1 converts a function with 1 parameters returning a [GRA] into a function with 1 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize1[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0) (R, error), GIOA ~func() E.Either[error, R], T0, R any](f func(T0) GRA) F {
	return func(c context.Context, t0 T0) (R, error) {
		return E.UnwrapError(f(t0)(c)())
	}
}

// Eitherize2 converts a function with 2 parameters returning a tuple into a function with 2 parameters returning a [GRA]
// The inverse function is [Uneitherize2]
func Eitherize2[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1) (R, error), GIOA ~func() E.Either[error, R], T0, T1, R any](f F) func(T0, T1) GRA {
	return RE.Eitherize2[GRA](f)
}

// Uneitherize2 converts a function with 2 parameters returning a [GRA] into a function with 2 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize2[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1) (R, error), GIOA ~func() E.Either[error, R], T0, T1, R any](f func(T0, T1) GRA) F {
	return func(c context.Context, t0 T0, t1 T1) (R, error) {
		return E.UnwrapError(f(t0, t1)(c)())
	}
}

// Eitherize3 converts a function with 3 parameters returning a tuple into a function with 3 parameters returning a [GRA]
// The inverse function is [Uneitherize3]
func Eitherize3[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, R any](f F) func(T0, T1, T2) GRA {
	return RE.Eitherize3[GRA](f)
}

// Uneitherize3 converts a function with 3 parameters returning a [GRA] into a function with 3 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize3[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, R any](f func(T0, T1, T2) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2) (R, error) {
		return E.UnwrapError(f(t0, t1, t2)(c)())
	}
}

// Eitherize4 converts a function with 4 parameters returning a tuple into a function with 4 parameters returning a [GRA]
// The inverse function is [Uneitherize4]
func Eitherize4[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, R any](f F) func(T0, T1, T2, T3) GRA {
	return RE.Eitherize4[GRA](f)
}

// Uneitherize4 converts a function with 4 parameters returning a [GRA] into a function with 4 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize4[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, R any](f func(T0, T1, T2, T3) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3)(c)())
	}
}

// Eitherize5 converts a function with 5 parameters returning a tuple into a function with 5 parameters returning a [GRA]
// The inverse function is [Uneitherize5]
func Eitherize5[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, R any](f F) func(T0, T1, T2, T3, T4) GRA {
	return RE.Eitherize5[GRA](f)
}

// Uneitherize5 converts a function with 5 parameters returning a [GRA] into a function with 5 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize5[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, R any](f func(T0, T1, T2, T3, T4) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4)(c)())
	}
}

// Eitherize6 converts a function with 6 parameters returning a tuple into a function with 6 parameters returning a [GRA]
// The inverse function is [Uneitherize6]
func Eitherize6[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, R any](f F) func(T0, T1, T2, T3, T4, T5) GRA {
	return RE.Eitherize6[GRA](f)
}

// Uneitherize6 converts a function with 6 parameters returning a [GRA] into a function with 6 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize6[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, R any](f func(T0, T1, T2, T3, T4, T5) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5)(c)())
	}
}

// Eitherize7 converts a function with 7 parameters returning a tuple into a function with 7 parameters returning a [GRA]
// The inverse function is [Uneitherize7]
func Eitherize7[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) GRA {
	return RE.Eitherize7[GRA](f)
}

// Uneitherize7 converts a function with 7 parameters returning a [GRA] into a function with 7 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize7[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, R any](f func(T0, T1, T2, T3, T4, T5, T6) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6)(c)())
	}
}

// Eitherize8 converts a function with 8 parameters returning a tuple into a function with 8 parameters returning a [GRA]
// The inverse function is [Uneitherize8]
func Eitherize8[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) GRA {
	return RE.Eitherize8[GRA](f)
}

// Uneitherize8 converts a function with 8 parameters returning a [GRA] into a function with 8 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize8[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7)(c)())
	}
}

// Eitherize9 converts a function with 9 parameters returning a tuple into a function with 9 parameters returning a [GRA]
// The inverse function is [Uneitherize9]
func Eitherize9[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) GRA {
	return RE.Eitherize9[GRA](f)
}

// Uneitherize9 converts a function with 9 parameters returning a [GRA] into a function with 9 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize9[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7, T8) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7, t8)(c)())
	}
}

// Eitherize10 converts a function with 10 parameters returning a tuple into a function with 10 parameters returning a [GRA]
// The inverse function is [Uneitherize10]
func Eitherize10[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) GRA {
	return RE.Eitherize10[GRA](f)
}

// Uneitherize10 converts a function with 10 parameters returning a [GRA] into a function with 10 parameters returning a tuple.
// The first parameter is considered to be the [context.Context].
func Uneitherize10[GRA ~func(context.Context) GIOA, F ~func(context.Context, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) GRA) F {
	return func(c context.Context, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7, t8, t9)(c)())
	}
}
