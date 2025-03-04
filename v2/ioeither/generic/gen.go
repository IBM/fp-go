// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2025-03-04 23:59:33.8343102 +0100 CET m=+0.003180601
package generic


import (
	ET "github.com/IBM/fp-go/v2/either"
)

// Eitherize0 converts a function with 0 parameters returning a tuple into a function with 0 parameters returning a [GIOA]
func Eitherize0[GIOA ~func() ET.Either[error, R], F ~func() (R, error), R any](f F) func() GIOA {
  e := ET.Eitherize0(f)
  return func() GIOA {
    return func() ET.Either[error, R] {
      return e()
    }}
}

// Uneitherize0 converts a function with 0 parameters returning a tuple into a function with 0 parameters returning a [GIOA]
func Uneitherize0[GIOA ~func() ET.Either[error, R], GTA ~func() GIOA, R any](f GTA) func() (R, error) {
  return func() (R, error) {
    return ET.Unwrap(f()())
  }
}

// Eitherize1 converts a function with 1 parameters returning a tuple into a function with 1 parameters returning a [GIOA]
func Eitherize1[GIOA ~func() ET.Either[error, R], F ~func(T1) (R, error), T1, R any](f F) func(T1) GIOA {
  e := ET.Eitherize1(f)
  return func(t1 T1) GIOA {
    return func() ET.Either[error, R] {
      return e(t1)
    }}
}

// Uneitherize1 converts a function with 1 parameters returning a tuple into a function with 1 parameters returning a [GIOA]
func Uneitherize1[GIOA ~func() ET.Either[error, R], GTA ~func(T1) GIOA, T1, R any](f GTA) func(T1) (R, error) {
  return func(t1 T1) (R, error) {
    return ET.Unwrap(f(t1)())
  }
}

// Eitherize2 converts a function with 2 parameters returning a tuple into a function with 2 parameters returning a [GIOA]
func Eitherize2[GIOA ~func() ET.Either[error, R], F ~func(T1, T2) (R, error), T1, T2, R any](f F) func(T1, T2) GIOA {
  e := ET.Eitherize2(f)
  return func(t1 T1, t2 T2) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2)
    }}
}

// Uneitherize2 converts a function with 2 parameters returning a tuple into a function with 2 parameters returning a [GIOA]
func Uneitherize2[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2) GIOA, T1, T2, R any](f GTA) func(T1, T2) (R, error) {
  return func(t1 T1, t2 T2) (R, error) {
    return ET.Unwrap(f(t1, t2)())
  }
}

// Eitherize3 converts a function with 3 parameters returning a tuple into a function with 3 parameters returning a [GIOA]
func Eitherize3[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3) (R, error), T1, T2, T3, R any](f F) func(T1, T2, T3) GIOA {
  e := ET.Eitherize3(f)
  return func(t1 T1, t2 T2, t3 T3) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3)
    }}
}

// Uneitherize3 converts a function with 3 parameters returning a tuple into a function with 3 parameters returning a [GIOA]
func Uneitherize3[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3) GIOA, T1, T2, T3, R any](f GTA) func(T1, T2, T3) (R, error) {
  return func(t1 T1, t2 T2, t3 T3) (R, error) {
    return ET.Unwrap(f(t1, t2, t3)())
  }
}

// Eitherize4 converts a function with 4 parameters returning a tuple into a function with 4 parameters returning a [GIOA]
func Eitherize4[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4) (R, error), T1, T2, T3, T4, R any](f F) func(T1, T2, T3, T4) GIOA {
  e := ET.Eitherize4(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4)
    }}
}

// Uneitherize4 converts a function with 4 parameters returning a tuple into a function with 4 parameters returning a [GIOA]
func Uneitherize4[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4) GIOA, T1, T2, T3, T4, R any](f GTA) func(T1, T2, T3, T4) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4)())
  }
}

// Eitherize5 converts a function with 5 parameters returning a tuple into a function with 5 parameters returning a [GIOA]
func Eitherize5[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5) (R, error), T1, T2, T3, T4, T5, R any](f F) func(T1, T2, T3, T4, T5) GIOA {
  e := ET.Eitherize5(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5)
    }}
}

// Uneitherize5 converts a function with 5 parameters returning a tuple into a function with 5 parameters returning a [GIOA]
func Uneitherize5[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5) GIOA, T1, T2, T3, T4, T5, R any](f GTA) func(T1, T2, T3, T4, T5) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5)())
  }
}

// Eitherize6 converts a function with 6 parameters returning a tuple into a function with 6 parameters returning a [GIOA]
func Eitherize6[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5, T6) (R, error), T1, T2, T3, T4, T5, T6, R any](f F) func(T1, T2, T3, T4, T5, T6) GIOA {
  e := ET.Eitherize6(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5, t6)
    }}
}

// Uneitherize6 converts a function with 6 parameters returning a tuple into a function with 6 parameters returning a [GIOA]
func Uneitherize6[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5, T6) GIOA, T1, T2, T3, T4, T5, T6, R any](f GTA) func(T1, T2, T3, T4, T5, T6) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5, t6)())
  }
}

// Eitherize7 converts a function with 7 parameters returning a tuple into a function with 7 parameters returning a [GIOA]
func Eitherize7[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5, T6, T7) (R, error), T1, T2, T3, T4, T5, T6, T7, R any](f F) func(T1, T2, T3, T4, T5, T6, T7) GIOA {
  e := ET.Eitherize7(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5, t6, t7)
    }}
}

// Uneitherize7 converts a function with 7 parameters returning a tuple into a function with 7 parameters returning a [GIOA]
func Uneitherize7[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5, T6, T7) GIOA, T1, T2, T3, T4, T5, T6, T7, R any](f GTA) func(T1, T2, T3, T4, T5, T6, T7) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5, t6, t7)())
  }
}

// Eitherize8 converts a function with 8 parameters returning a tuple into a function with 8 parameters returning a [GIOA]
func Eitherize8[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8) GIOA {
  e := ET.Eitherize8(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5, t6, t7, t8)
    }}
}

// Uneitherize8 converts a function with 8 parameters returning a tuple into a function with 8 parameters returning a [GIOA]
func Uneitherize8[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5, T6, T7, T8) GIOA, T1, T2, T3, T4, T5, T6, T7, T8, R any](f GTA) func(T1, T2, T3, T4, T5, T6, T7, T8) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5, t6, t7, t8)())
  }
}

// Eitherize9 converts a function with 9 parameters returning a tuple into a function with 9 parameters returning a [GIOA]
func Eitherize9[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) GIOA {
  e := ET.Eitherize9(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5, t6, t7, t8, t9)
    }}
}

// Uneitherize9 converts a function with 9 parameters returning a tuple into a function with 9 parameters returning a [GIOA]
func Uneitherize9[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9) GIOA, T1, T2, T3, T4, T5, T6, T7, T8, T9, R any](f GTA) func(T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5, t6, t7, t8, t9)())
  }
}

// Eitherize10 converts a function with 10 parameters returning a tuple into a function with 10 parameters returning a [GIOA]
func Eitherize10[GIOA ~func() ET.Either[error, R], F ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error), T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f F) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) GIOA {
  e := ET.Eitherize10(f)
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10) GIOA {
    return func() ET.Either[error, R] {
      return e(t1, t2, t3, t4, t5, t6, t7, t8, t9, t10)
    }}
}

// Uneitherize10 converts a function with 10 parameters returning a tuple into a function with 10 parameters returning a [GIOA]
func Uneitherize10[GIOA ~func() ET.Either[error, R], GTA ~func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) GIOA, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R any](f GTA) func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) (R, error) {
  return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10) (R, error) {
    return ET.Unwrap(f(t1, t2, t3, t4, t5, t6, t7, t8, t9, t10)())
  }
}
