// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2023-09-12 13:27:21.0181793 +0200 CEST m=+0.024632901

package readerioeither


import (
	G "github.com/IBM/fp-go/readerioeither/generic"	
)

// From0 converts a function with 1 parameters returning a tuple into a function with 0 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From0[F ~func(C) func() (R, error), C, R any](f F) func() ReaderIOEither[C, error, R] {
  return G.From0[ReaderIOEither[C, error, R]](f)
}

// Eitherize0 converts a function with 1 parameters returning a tuple into a function with 0 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize0[F ~func(C) (R, error), C, R any](f F) func() ReaderIOEither[C, error, R] {
  return G.Eitherize0[ReaderIOEither[C, error, R]](f)
}

// From1 converts a function with 2 parameters returning a tuple into a function with 1 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From1[F ~func(C, T0) func() (R, error), T0, C, R any](f F) func(T0) ReaderIOEither[C, error, R] {
  return G.From1[ReaderIOEither[C, error, R]](f)
}

// Eitherize1 converts a function with 2 parameters returning a tuple into a function with 1 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize1[F ~func(C, T0) (R, error), T0, C, R any](f F) func(T0) ReaderIOEither[C, error, R] {
  return G.Eitherize1[ReaderIOEither[C, error, R]](f)
}

// From2 converts a function with 3 parameters returning a tuple into a function with 2 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From2[F ~func(C, T0, T1) func() (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOEither[C, error, R] {
  return G.From2[ReaderIOEither[C, error, R]](f)
}

// Eitherize2 converts a function with 3 parameters returning a tuple into a function with 2 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize2[F ~func(C, T0, T1) (R, error), T0, T1, C, R any](f F) func(T0, T1) ReaderIOEither[C, error, R] {
  return G.Eitherize2[ReaderIOEither[C, error, R]](f)
}

// From3 converts a function with 4 parameters returning a tuple into a function with 3 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From3[F ~func(C, T0, T1, T2) func() (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOEither[C, error, R] {
  return G.From3[ReaderIOEither[C, error, R]](f)
}

// Eitherize3 converts a function with 4 parameters returning a tuple into a function with 3 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize3[F ~func(C, T0, T1, T2) (R, error), T0, T1, T2, C, R any](f F) func(T0, T1, T2) ReaderIOEither[C, error, R] {
  return G.Eitherize3[ReaderIOEither[C, error, R]](f)
}

// From4 converts a function with 5 parameters returning a tuple into a function with 4 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From4[F ~func(C, T0, T1, T2, T3) func() (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOEither[C, error, R] {
  return G.From4[ReaderIOEither[C, error, R]](f)
}

// Eitherize4 converts a function with 5 parameters returning a tuple into a function with 4 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize4[F ~func(C, T0, T1, T2, T3) (R, error), T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) ReaderIOEither[C, error, R] {
  return G.Eitherize4[ReaderIOEither[C, error, R]](f)
}

// From5 converts a function with 6 parameters returning a tuple into a function with 5 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From5[F ~func(C, T0, T1, T2, T3, T4) func() (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOEither[C, error, R] {
  return G.From5[ReaderIOEither[C, error, R]](f)
}

// Eitherize5 converts a function with 6 parameters returning a tuple into a function with 5 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize5[F ~func(C, T0, T1, T2, T3, T4) (R, error), T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) ReaderIOEither[C, error, R] {
  return G.Eitherize5[ReaderIOEither[C, error, R]](f)
}

// From6 converts a function with 7 parameters returning a tuple into a function with 6 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From6[F ~func(C, T0, T1, T2, T3, T4, T5) func() (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOEither[C, error, R] {
  return G.From6[ReaderIOEither[C, error, R]](f)
}

// Eitherize6 converts a function with 7 parameters returning a tuple into a function with 6 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize6[F ~func(C, T0, T1, T2, T3, T4, T5) (R, error), T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) ReaderIOEither[C, error, R] {
  return G.Eitherize6[ReaderIOEither[C, error, R]](f)
}

// From7 converts a function with 8 parameters returning a tuple into a function with 7 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) func() (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOEither[C, error, R] {
  return G.From7[ReaderIOEither[C, error, R]](f)
}

// Eitherize7 converts a function with 8 parameters returning a tuple into a function with 7 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize7[F ~func(C, T0, T1, T2, T3, T4, T5, T6) (R, error), T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) ReaderIOEither[C, error, R] {
  return G.Eitherize7[ReaderIOEither[C, error, R]](f)
}

// From8 converts a function with 9 parameters returning a tuple into a function with 8 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOEither[C, error, R] {
  return G.From8[ReaderIOEither[C, error, R]](f)
}

// Eitherize8 converts a function with 9 parameters returning a tuple into a function with 8 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize8[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) ReaderIOEither[C, error, R] {
  return G.Eitherize8[ReaderIOEither[C, error, R]](f)
}

// From9 converts a function with 10 parameters returning a tuple into a function with 9 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOEither[C, error, R] {
  return G.From9[ReaderIOEither[C, error, R]](f)
}

// Eitherize9 converts a function with 10 parameters returning a tuple into a function with 9 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize9[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) ReaderIOEither[C, error, R] {
  return G.Eitherize9[ReaderIOEither[C, error, R]](f)
}

// From10 converts a function with 11 parameters returning a tuple into a function with 10 parameters returning a [ReaderIOEither[R]]
// The first parameter is considered to be the context [C].
func From10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) func() (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOEither[C, error, R] {
  return G.From10[ReaderIOEither[C, error, R]](f)
}

// Eitherize10 converts a function with 11 parameters returning a tuple into a function with 10 parameters returning a [ReaderIOEither[C, error, R]]
// The first parameter is considered to be the context [C].
func Eitherize10[F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) ReaderIOEither[C, error, R] {
  return G.Eitherize10[ReaderIOEither[C, error, R]](f)
}
