// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package option implements the Option monad, a data type that can have a defined value or none
package option

import (
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	P "github.com/IBM/fp-go/v2/predicate"
)

// fromPredicate creates an Option based on a predicate function.
// If the predicate returns true for the value, it returns Some(a), otherwise None.
func fromPredicate[A any](a A, pred func(A) bool) Option[A] {
	if pred(a) {
		return Some(a)
	}
	return None[A]()
}

// FromPredicate returns a function that creates an Option based on a predicate.
// The returned function will wrap a value in Some if the predicate is satisfied, otherwise None.
//
// Example:
//
//	isPositive := FromPredicate(N.MoreThan(0))
//	result := isPositive(5)  // Some(5)
//	result := isPositive(-1) // None
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A] {
	return F.Bind2nd(fromPredicate[A], pred)
}

//go:inline
func FromZero[A comparable]() Kleisli[A, A] {
	return FromPredicate(P.IsZero[A]())
}

//go:inline
func FromNonZero[A comparable]() Kleisli[A, A] {
	return FromPredicate(P.IsNonZero[A]())
}

//go:inline
func FromEq[A any](pred eq.Eq[A]) func(A) Kleisli[A, A] {
	return F.Flow2(P.IsEqual(pred), FromPredicate[A])
}

//go:inline
func FromStrictEq[A comparable]() func(A) Kleisli[A, A] {
	return FromEq(eq.FromStrictEquals[A]())
}

// FromNillable converts a pointer to an Option.
// Returns Some if the pointer is non-nil, None otherwise.
//
// Example:
//
//	var ptr *int = nil
//	result := FromNillable(ptr) // None
//	val := 42
//	result := FromNillable(&val) // Some(&val)
//
//go:inline
func FromNillable[A any](a *A) Option[*A] {
	return fromPredicate(a, F.IsNonNil[A])
}

// FromValidation converts a validation function (returning value and bool) to an Option-returning function.
// This is an alias for Optionize1.
//
// Example:
//
//	parseNum := FromValidation(func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	})
//	result := parseNum("42") // Some(42)
//
//go:inline
func FromValidation[A, B any](f func(A) (B, bool)) Kleisli[A, B] {
	return Optionize1(f)
}

// MonadAp applies a function wrapped in an Option to a value wrapped in an Option.
// If either the function or the value is None, returns None.
// This is the monadic form of the applicative functor.
//
// Example:
//
//	fab := Some(N.Mul(2))
//	fa := Some(5)
//	result := MonadAp(fab, fa) // Some(10)
func MonadAp[B, A any](fab Option[func(A) B], fa Option[A]) Option[B] {
	if fab.isSome && fa.isSome {
		return Some(fab.value(fa.value))
	}
	return None[B]()
}

// Ap is the curried applicative functor for Option.
// Returns a function that applies an Option-wrapped function to the given Option value.
//
// Example:
//
//	fa := Some(5)
//	applyTo5 := Ap[int](fa)
//	fab := Some(N.Mul(2))
//	result := applyTo5(fab) // Some(10)
func Ap[B, A any](fa Option[A]) Operator[func(A) B, B] {
	if fa.isSome {
		return func(fab Option[func(A) B]) Option[B] {
			if fab.isSome {
				return Some(fab.value(fa.value))
			}
			return None[B]()
		}
	}
	// shortcut
	return F.Constant1[Option[func(A) B]](None[B]())
}

// MonadMap applies a function to the value inside an Option.
// If the Option is None, returns None. This is the monadic form of Map.
//
// Example:
//
//	fa := Some(5)
//	result := MonadMap(fa, N.Mul(2)) // Some(10)
func MonadMap[A, B any](fa Option[A], f func(A) B) Option[B] {
	if fa.isSome {
		return Some(f(fa.value))
	}
	return None[B]()
}

// Map returns a function that applies a transformation to the value inside an Option.
// If the Option is None, returns None.
//
// Example:
//
//	double := Map(N.Mul(2))
//	result := double(Some(5)) // Some(10)
//	result := double(None[int]()) // None
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return func(fa Option[A]) Option[B] {
		if fa.isSome {
			return Some(f(fa.value))
		}
		return None[B]()
	}
}

// MonadMapTo replaces the value inside an Option with a constant value.
// If the Option is None, returns None. This is the monadic form of MapTo.
//
// Example:
//
//	fa := Some(5)
//	result := MonadMapTo(fa, "hello") // Some("hello")
func MonadMapTo[A, B any](fa Option[A], b B) Option[B] {
	if fa.isSome {
		return Some(b)
	}
	return None[B]()
}

// MapTo returns a function that replaces the value inside an Option with a constant.
//
// Example:
//
//	replaceWith42 := MapTo[string, int](42)
//	result := replaceWith42(Some("hello")) // Some(42)
func MapTo[A, B any](b B) Operator[A, B] {
	return func(fa Option[A]) Option[B] {
		if fa.isSome {
			return Some(b)
		}
		return None[B]()
	}
}

// TryCatch executes a function that may return an error and converts the result to an Option.
// Returns Some(value) if no error occurred, None if an error occurred.
//
// Example:
//
//	result := TryCatch(func() (int, error) {
//	    return strconv.Atoi("42")
//	}) // Some(42)
func TryCatch[A any](f func() (A, error)) Option[A] {
	val, err := f()
	if err != nil {
		return None[A]()
	}
	return Some(val)
}

// Fold provides a way to handle both Some and None cases of an Option.
// Returns a function that applies onNone if the Option is None, or onSome if it's Some.
//
// Example:
//
//	handler := Fold(
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	)
//	result := handler(Some(42)) // "value: 42"
//	result := handler(None[int]()) // "no value"
//
//go:inline
func Fold[A, B any](onNone func() B, onSome func(a A) B) func(ma Option[A]) B {
	return func(fa Option[A]) B {
		return MonadFold(fa, onNone, onSome)
	}
}

// MonadGetOrElse extracts the value from an Option or returns a default value.
// This is the monadic form of GetOrElse.
//
// Example:
//
//	result := MonadGetOrElse(Some(42), func() int { return 0 }) // 42
//	result := MonadGetOrElse(None[int](), func() int { return 0 }) // 0
//
//go:inline
func MonadGetOrElse[A any](fa Option[A], onNone func() A) A {
	return MonadFold(fa, onNone, F.Identity[A])
}

// GetOrElse returns a function that extracts the value from an Option or returns a default.
//
// Example:
//
//	getOrZero := GetOrElse(func() int { return 0 })
//	result := getOrZero(Some(42)) // 42
//	result := getOrZero(None[int]()) // 0
//
//go:inline
func GetOrElse[A any](onNone func() A) func(Option[A]) A {
	return Fold(onNone, F.Identity[A])
}

// MonadChain applies a function that returns an Option to the value inside an Option.
// This is the monadic bind operation. If the input is None, returns None.
//
// Example:
//
//	fa := Some(5)
//	result := MonadChain(fa, func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	}) // Some(10)
//
//go:inline
func MonadChain[A, B any](fa Option[A], f Kleisli[A, B]) Option[B] {
	return MonadFold(fa, None[B], f)
}

// Chain returns a function that applies an Option-returning function to an Option value.
// This is the curried form of the monadic bind operation.
//
// Example:
//
//	validate := Chain(func(x int) Option[int] {
//	    if x > 0 { return Some(x * 2) }
//	    return None[int]()
//	})
//	result := validate(Some(5)) // Some(10)
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return Fold(None[B], f)
}

// MonadChainTo ignores the first Option and returns the second Option.
// Useful for sequencing operations where the first result is not needed.
//
// Example:
//
//	result := MonadChainTo(Some(5), Some("hello")) // Some("hello")
func MonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B] {
	if ma.isSome {
		return mb
	}
	return None[B]()
}

// ChainTo returns a function that ignores its input Option and returns a fixed Option.
//
// Example:
//
//	replaceWith := ChainTo(Some("hello"))
//	result := replaceWith(Some(42)) // Some("hello")
func ChainTo[A, B any](mb Option[B]) Operator[A, B] {
	if mb.isSome {
		return func(ma Option[A]) Option[B] {
			if ma.isSome {
				return mb
			}
			return None[B]()
		}
	}
	return F.Constant1[Option[A]](None[B]())
}

// MonadChainFirst applies a function that returns an Option but keeps the original value.
// If either operation results in None, returns None.
//
// Example:
//
//	result := MonadChainFirst(Some(5), func(x int) Option[string] {
//	    return Some(fmt.Sprintf("%d", x))
//	}) // Some(5) - original value is kept
func MonadChainFirst[A, B any](ma Option[A], f Kleisli[A, B]) Option[A] {
	return C.MonadChainFirst(
		MonadChain[A, A],
		MonadMap[B, A],
		ma,
		f,
	)
}

// ChainFirst returns a function that applies an Option-returning function but keeps the original value.
//
// Example:
//
//	logAndKeep := ChainFirst(func(x int) Option[string] {
//	    fmt.Println(x)
//	    return Some("logged")
//	})
//	result := logAndKeep(Some(5)) // Some(5)
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return C.ChainFirst(
		Chain[A, A],
		Map[B, A],
		f,
	)
}

// Flatten removes one level of nesting from a nested Option.
//
// Example:
//
//	nested := Some(Some(42))
//	result := Flatten(nested) // Some(42)
//	nested := Some(None[int]())
//	result := Flatten(nested) // None
func Flatten[A any](mma Option[Option[A]]) Option[A] {
	return MonadChain(mma, F.Identity[Option[A]])
}

// MonadAlt returns the first Option if it's Some, otherwise returns the alternative.
// This is the monadic form of the Alt operation.
//
// Example:
//
//	result := MonadAlt(Some(5), func() Option[int] { return Some(10) }) // Some(5)
//	result := MonadAlt(None[int](), func() Option[int] { return Some(10) }) // Some(10)
func MonadAlt[A any](fa Option[A], that func() Option[A]) Option[A] {
	return MonadFold(fa, that, Of[A])
}

// Alt returns a function that provides an alternative Option if the input is None.
//
// Example:
//
//	withDefault := Alt(func() Option[int] { return Some(0) })
//	result := withDefault(Some(5)) // Some(5)
//	result := withDefault(None[int]()) // Some(0)
func Alt[A any](that func() Option[A]) Operator[A, A] {
	return Fold(that, Of[A])
}

// MonadSequence2 sequences two Options and applies a function to their values.
// Returns None if either Option is None.
//
// Example:
//
//	result := MonadSequence2(Some(2), Some(3), func(a, b int) Option[int] {
//	    return Some(a + b)
//	}) // Some(5)
func MonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R] {
	if o1.isSome && o2.isSome {
		return f(o1.value, o2.value)
	}
	return None[R]()
}

// Sequence2 returns a function that sequences two Options with a combining function.
//
// Example:
//
//	add := Sequence2(func(a, b int) Option[int] { return Some(a + b) })
//	result := add(Some(2), Some(3)) // Some(5)
func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R] {
	return func(o1 Option[T1], o2 Option[T2]) Option[R] {
		return MonadSequence2(o1, o2, f)
	}
}

// Reduce folds an Option into a single value using a reducer function.
// If the Option is None, returns the initial value.
//
// Example:
//
//	sum := Reduce(func(acc, val int) int { return acc + val }, 0)
//	result := sum(Some(5)) // 5
//	result := sum(None[int]()) // 0
func Reduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B {
	return func(ma Option[A]) B {
		if ma.isSome {
			return f(initial, ma.value)
		}
		return initial
	}
}

// Filter keeps the Option if it's Some and the predicate is satisfied, otherwise returns None.
//
// Example:
//
//	isPositive := Filter(N.MoreThan(0))
//	result := isPositive(Some(5)) // Some(5)
//	result := isPositive(Some(-1)) // None
//	result := isPositive(None[int]()) // None
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return func(ma Option[A]) Option[A] {
		if ma.isSome && pred(ma.value) {
			return ma
		}
		return None[A]()
	}
}

// MonadFlap applies a value to a function wrapped in an Option.
// This is the monadic form of Flap.
//
// Example:
//
//	fab := Some(N.Mul(2))
//	result := MonadFlap(fab, 5) // Some(10)
func MonadFlap[B, A any](fab Option[func(A) B], a A) Option[B] {
	if fab.isSome {
		return Some(fab.value(a))
	}
	return None[B]()
}

// Flap returns a function that applies a value to an Option-wrapped function.
//
// Example:
//
//	applyFive := Flap[int](5)
//	fab := Some(N.Mul(2))
//	result := applyFive(fab) // Some(10)
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return func(fab Option[func(A) B]) Option[B] {
		if fab.isSome {
			return Some(fab.value(a))
		}
		return None[B]()
	}
}

// Zero returns the zero value of an [Option], which is None.
// This function is useful as an identity element in monoid operations or for creating an empty Option.
//
// The zero value for Option[A] is always None, representing the absence of a value.
// This is consistent with the Option monad's semantics where None represents "no value"
// and Some represents "a value".
//
// Important: Zero() returns the same value as the default initialization of Option[A].
// When you declare `var o Option[A]` without initialization, it has the same value as Zero[A]().
//
// Note: Unlike other types where zero might be a default value, Option's zero is explicitly
// the absence of any value (None), not Some with a zero value.
//
// Example:
//
//	// Zero Option of any type is always None
//	o1 := option.Zero[int]()     // None
//	o2 := option.Zero[string]()  // None
//	o3 := option.Zero[*int]()    // None
//
//	// Zero equals default initialization
//	var defaultInit Option[int]
//	zero := option.Zero[int]()
//	assert.Equal(t, defaultInit, zero) // true
//
//	// Verify it's None
//	o := option.Zero[int]()
//	assert.True(t, option.IsNone(o))   // true
//	assert.False(t, option.IsSome(o))  // false
//
//	// Different from Some with zero value
//	someZero := option.Some(0)         // Some(0)
//	zero := option.Zero[int]()         // None
//	assert.NotEqual(t, someZero, zero) // they are different
func Zero[A any]() Option[A] {
	return None[A]()
}
