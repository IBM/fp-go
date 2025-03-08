// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2025-03-07 23:16:11.0357415 +0100 CET m=+0.045725301
package generic

import (
	E "github.com/IBM/fp-go/v2/either"
	RD "github.com/IBM/fp-go/v2/reader/generic"
)

// From0 converts a function with 1 parameters returning a tuple into a function with 0 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From0[GRA ~func(C) GIOA, F ~func(C) func() (R, error), GIOA ~func() E.Either[error, R], C, R any](f F) func() GRA {
	return RD.From0[GRA](func(r C) GIOA {
		return E.Eitherize0(f(r))
	})
}

// Eitherize0 converts a function with 0 parameters returning a tuple into a function with 0 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize0[GRA ~func(C) GIOA, F ~func(C) (R, error), GIOA ~func() E.Either[error, R], C, R any](f F) func() GRA {
	return From0[GRA](func(r C) func() (R, error) {
		return func() (R, error) {
			return f(r)
		}
	})
}

// Uneitherize0 converts a function with 0 parameters returning a [GRA] into a function with 0 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize0[GRA ~func(C) GIOA, F ~func(C) (R, error), GIOA ~func() E.Either[error, R], C, R any](f func() GRA) F {
	return func(c C) (R, error) {
		return E.UnwrapError(f()(c)())
	}
}

// From1 converts a function with 2 parameters returning a tuple into a function with 1 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From1[GRA ~func(C) GIOA, F ~func(C, T0) func() (R, error), GIOA ~func() E.Either[error, R], T0, C, R any](f F) func(T0) GRA {
	return RD.From1[GRA](func(r C, t0 T0) GIOA {
		return E.Eitherize0(f(r, t0))
	})
}

// Eitherize1 converts a function with 1 parameters returning a tuple into a function with 1 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize1[GRA ~func(C) GIOA, F ~func(C, T0) (R, error), GIOA ~func() E.Either[error, R], T0, C, R any](f F) func(T0) GRA {
	return From1[GRA](func(r C, t0 T0) func() (R, error) {
		return func() (R, error) {
			return f(r, t0)
		}
	})
}

// Uneitherize1 converts a function with 1 parameters returning a [GRA] into a function with 1 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize1[GRA ~func(C) GIOA, F ~func(C, T0) (R, error), GIOA ~func() E.Either[error, R], T0, C, R any](f func(T0) GRA) F {
	return func(c C, t0 T0) (R, error) {
		return E.UnwrapError(f(t0)(c)())
	}
}

// From2 converts a function with 3 parameters returning a tuple into a function with 2 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From2[GRA ~func(C) GIOA, F ~func(C, T0, T1) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, C, R any](f F) func(T0, T1) GRA {
	return RD.From2[GRA](func(r C, t0 T0, t1 T1) GIOA {
		return E.Eitherize0(f(r, t0, t1))
	})
}

// Eitherize2 converts a function with 2 parameters returning a tuple into a function with 2 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize2[GRA ~func(C) GIOA, F ~func(C, T0, T1) (R, error), GIOA ~func() E.Either[error, R], T0, T1, C, R any](f F) func(T0, T1) GRA {
	return From2[GRA](func(r C, t0 T0, t1 T1) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1)
		}
	})
}

// Uneitherize2 converts a function with 2 parameters returning a [GRA] into a function with 2 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize2[GRA ~func(C) GIOA, F ~func(C, T0, T1) (R, error), GIOA ~func() E.Either[error, R], T0, T1, C, R any](f func(T0, T1) GRA) F {
	return func(c C, t0 T0, t1 T1) (R, error) {
		return E.UnwrapError(f(t0, t1)(c)())
	}
}

// From3 converts a function with 4 parameters returning a tuple into a function with 3 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From3[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, C, R any](f F) func(T0, T1, T2) GRA {
	return RD.From3[GRA](func(r C, t0 T0, t1 T1, t2 T2) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2))
	})
}

// Eitherize3 converts a function with 3 parameters returning a tuple into a function with 3 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize3[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, C, R any](f F) func(T0, T1, T2) GRA {
	return From3[GRA](func(r C, t0 T0, t1 T1, t2 T2) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2)
		}
	})
}

// Uneitherize3 converts a function with 3 parameters returning a [GRA] into a function with 3 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize3[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, C, R any](f func(T0, T1, T2) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2) (R, error) {
		return E.UnwrapError(f(t0, t1, t2)(c)())
	}
}

// From4 converts a function with 5 parameters returning a tuple into a function with 4 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From4[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) GRA {
	return RD.From4[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3))
	})
}

// Eitherize4 converts a function with 4 parameters returning a tuple into a function with 4 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize4[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, C, R any](f F) func(T0, T1, T2, T3) GRA {
	return From4[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3)
		}
	})
}

// Uneitherize4 converts a function with 4 parameters returning a [GRA] into a function with 4 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize4[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, C, R any](f func(T0, T1, T2, T3) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3)(c)())
	}
}

// From5 converts a function with 6 parameters returning a tuple into a function with 5 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From5[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) GRA {
	return RD.From5[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4))
	})
}

// Eitherize5 converts a function with 5 parameters returning a tuple into a function with 5 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize5[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, C, R any](f F) func(T0, T1, T2, T3, T4) GRA {
	return From5[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4)
		}
	})
}

// Uneitherize5 converts a function with 5 parameters returning a [GRA] into a function with 5 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize5[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, C, R any](f func(T0, T1, T2, T3, T4) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4)(c)())
	}
}

// From6 converts a function with 7 parameters returning a tuple into a function with 6 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From6[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) GRA {
	return RD.From6[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4, t5))
	})
}

// Eitherize6 converts a function with 6 parameters returning a tuple into a function with 6 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize6[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, C, R any](f F) func(T0, T1, T2, T3, T4, T5) GRA {
	return From6[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4, t5)
		}
	})
}

// Uneitherize6 converts a function with 6 parameters returning a [GRA] into a function with 6 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize6[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, C, R any](f func(T0, T1, T2, T3, T4, T5) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5)(c)())
	}
}

// From7 converts a function with 8 parameters returning a tuple into a function with 7 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From7[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) GRA {
	return RD.From7[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4, t5, t6))
	})
}

// Eitherize7 converts a function with 7 parameters returning a tuple into a function with 7 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize7[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6) GRA {
	return From7[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4, t5, t6)
		}
	})
}

// Uneitherize7 converts a function with 7 parameters returning a [GRA] into a function with 7 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize7[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, C, R any](f func(T0, T1, T2, T3, T4, T5, T6) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6)(c)())
	}
}

// From8 converts a function with 9 parameters returning a tuple into a function with 8 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From8[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) GRA {
	return RD.From8[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4, t5, t6, t7))
	})
}

// Eitherize8 converts a function with 8 parameters returning a tuple into a function with 8 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize8[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7) GRA {
	return From8[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4, t5, t6, t7)
		}
	})
}

// Uneitherize8 converts a function with 8 parameters returning a [GRA] into a function with 8 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize8[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, C, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7)(c)())
	}
}

// From9 converts a function with 10 parameters returning a tuple into a function with 9 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From9[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) GRA {
	return RD.From9[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4, t5, t6, t7, t8))
	})
}

// Eitherize9 converts a function with 9 parameters returning a tuple into a function with 9 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize9[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8) GRA {
	return From9[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4, t5, t6, t7, t8)
		}
	})
}

// Uneitherize9 converts a function with 9 parameters returning a [GRA] into a function with 9 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize9[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, C, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7, T8) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7, t8)(c)())
	}
}

// From10 converts a function with 11 parameters returning a tuple into a function with 10 parameters returning a [GRA]
// The first parameter is considerd to be the context [C].
func From10[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) func() (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) GRA {
	return RD.From10[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) GIOA {
		return E.Eitherize0(f(r, t0, t1, t2, t3, t4, t5, t6, t7, t8, t9))
	})
}

// Eitherize10 converts a function with 10 parameters returning a tuple into a function with 10 parameters returning a [GRA]
// The first parameter is considered to be the context [C].
func Eitherize10[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f F) func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) GRA {
	return From10[GRA](func(r C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) func() (R, error) {
		return func() (R, error) {
			return f(r, t0, t1, t2, t3, t4, t5, t6, t7, t8, t9)
		}
	})
}

// Uneitherize10 converts a function with 10 parameters returning a [GRA] into a function with 10 parameters returning a tuple.
// The first parameter is considered to be the context [C].
func Uneitherize10[GRA ~func(C) GIOA, F ~func(C, T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) (R, error), GIOA ~func() E.Either[error, R], T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, C, R any](f func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) GRA) F {
	return func(c C, t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) (R, error) {
		return E.UnwrapError(f(t0, t1, t2, t3, t4, t5, t6, t7, t8, t9)(c)())
	}
}
