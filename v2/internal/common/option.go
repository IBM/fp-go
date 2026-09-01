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
package common

import (
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	P "github.com/IBM/fp-go/v2/predicate"
)

// OptionOptionize1 converts a function with 1 parameters returning a tuple of a return value R and a boolean into a function with 1 parameters returning an Option[R]
func OptionOptionize1[F ~func(T0) (R, bool), T0, R any](f F) OptionKleisli[T0, R] {
	none := OptionNone[R]()
	return func(t0 T0) Option[R] {
		if r, ok := f(t0); ok {
			return OptionSome(r)
		}
		return none
	}
}

// optionFromPredicate creates an Option based on a predicate function.
// If the predicate returns true for the value, it returns OptionSome(a), otherwise None.
func optionFromPredicate[A any](a A, pred Predicate[A]) Option[A] {
	if pred(a) {
		return OptionSome(a)
	}
	return OptionNone[A]()
}

// OptionFromPredicate returns a function that creates an Option based on a predicate.
// The returned function will wrap a value in Some if the predicate is satisfied, otherwise None.
//
// Example:
//
//	isPositive := OptionFromPredicate(N.MoreThan(0))
//	result := isPositive(5)  // OptionSome(5)
//	result := isPositive(-1) // None
func OptionFromPredicate[A any](pred Predicate[A]) OptionKleisli[A, A] {
	return F.Bind2nd(optionFromPredicate[A], pred)
}

//go:inline
func OptionFromZero[A comparable]() OptionKleisli[A, A] {
	return OptionFromPredicate(P.IsZero[A]())
}

//go:inline
func OptionFromNonZero[A comparable]() OptionKleisli[A, A] {
	return OptionFromPredicate(P.IsNonZero[A]())
}

//go:inline
func OptionFromEq[A any](pred eq.Eq[A]) func(A) OptionKleisli[A, A] {
	return F.Flow2(P.IsEqual(pred), OptionFromPredicate[A])
}

//go:inline
func OptionFromStrictEq[A comparable]() func(A) OptionKleisli[A, A] {
	return OptionFromEq(eq.FromStrictEquals[A]())
}

// OptionFromNillable converts a pointer to an Option wrapping the pointer itself.
// Returns OptionSome(*A) when the pointer is non-nil, None otherwise.
//
// Deprecated: Use FromNillable2 instead. FromNillable2 unwraps the pointed-to
// value into an Option[A], which is the more useful form: it avoids carrying a
// *A inside the Option and composes naturally with the rest of the Option API.
//
//go:inline
func OptionFromNillable[A any](a *A) Option[*A] {
	return optionFromPredicate(a, F.IsNonNil[A])
}

// OptionFromNillable2 converts a pointer to an Option of the pointed-to value.
// Returns OptionSome(value) when the pointer is non-nil (dereferencing it), None otherwise.
// This is the preferred alternative to FromNillable because the resulting Option[A]
// composes naturally with Map, Chain, and the rest of the Option API without
// having to deal with an intermediate pointer inside the Option.
//
// Type Parameters:
//   - A: the type that the pointer points to
//
// Parameters:
//   - a: a pointer to a value of type A; may be nil
//
// Returns:
//   - OptionSome(*a) when a is non-nil
//   - None when a is nil
//
// See Also:
//   - ToNillable2: the inverse — converts Option[A] back to *A
//   - FromNillable: the older variant that wraps the pointer itself as Option[*A]
func OptionFromNillable2[A any](a *A) Option[A] {
	if a != nil {
		return OptionOf(*a)
	}
	return OptionNone[A]()
}

// OptionToNillable2 converts an Option[A] back to a nullable pointer.
// Returns a pointer to the contained value when the Option is Some, nil otherwise.
// This is the inverse of FromNillable2: round-tripping through both functions
// preserves the value but allocates a new pointer on each call.
//
// Type Parameters:
//   - A: the type contained in the Option
//
// Parameters:
//   - fa: the Option value to convert
//
// Returns:
//   - a pointer to a copy of the contained value when fa is Some
//   - nil when fa is None
//
// See Also:
//   - FromNillable2: the inverse — converts *A to Option[A]
func OptionToNillable2[A any](fa Option[A]) *A {
	if fa.s {
		return &fa.a
	}
	return nil
}

// OptionFromValidation converts a validation function (returning value and bool) to an Option-returning function.
// This is an alias for Optionize1.
//
// Example:
//
//	parseNum := OptionFromValidation(func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	})
//	result := parseNum("42") // OptionSome(42)
//
//go:inline
func OptionFromValidation[A, B any](f func(A) (B, bool)) OptionKleisli[A, B] {
	return OptionOptionize1(f)
}

// OptionMonadAp applies a function wrapped in an Option to a value wrapped in an Option.
// If either the function or the value is None, returns None.
// This is the monadic form of the applicative functor.
//
// Example:
//
//	fab := OptionSome(N.Mul(2))
//	fa := OptionSome(5)
//	result := OptionMonadAp(fab, fa) // OptionSome(10)
func OptionMonadAp[B, A any](fab Option[func(A) B], fa Option[A]) Option[B] {
	if fab.s && fa.s {
		return OptionSome(fab.a(fa.a))
	}
	return OptionNone[B]()
}

// OptionAp is the curried applicative functor for Option.
// Returns a function that applies an Option-wrapped function to the given Option value.
//
// Example:
//
//	fa := OptionSome(5)
//	applyTo5 := OptionAp[int](fa)
//	fab := OptionSome(N.Mul(2))
//	result := applyTo5(fab) // OptionSome(10)
func OptionAp[B, A any](fa Option[A]) OptionOperator[func(A) B, B] {
	if fa.s {
		return func(fab Option[func(A) B]) Option[B] {
			if fab.s {
				return OptionSome(fab.a(fa.a))
			}
			return OptionNone[B]()
		}
	}
	// shortcut
	return F.Constant1[Option[func(A) B]](OptionNone[B]())
}

// OptionMonadMap applies a function to the value inside an Option.
// If the Option is None, returns None. This is the monadic form of Map.
//
// Example:
//
//	fa := OptionSome(5)
//	result := OptionMonadMap(fa, N.Mul(2)) // OptionSome(10)
func OptionMonadMap[A, B any](fa Option[A], f func(A) B) Option[B] {
	if fa.s {
		return OptionSome(f(fa.a))
	}
	return OptionNone[B]()
}

// OptionMap returns a function that applies a transformation to the value inside an Option.
// If the Option is None, returns None.
//
// Example:
//
//	double := OptionMap(N.Mul(2))
//	result := double(OptionSome(5)) // OptionSome(10)
//	result := double(OptionNone[int]()) // None
func OptionMap[A, B any](f func(a A) B) OptionOperator[A, B] {
	return func(fa Option[A]) Option[B] {
		if fa.s {
			return OptionSome(f(fa.a))
		}
		return OptionNone[B]()
	}
}

// OptionMonadMapTo replaces the value inside an Option with a constant value.
// If the Option is None, returns None. This is the monadic form of MapTo.
//
// Example:
//
//	fa := OptionSome(5)
//	result := OptionMonadMapTo(fa, "hello") // OptionSome("hello")
func OptionMonadMapTo[A, B any](fa Option[A], b B) Option[B] {
	if fa.s {
		return OptionSome(b)
	}
	return OptionNone[B]()
}

// OptionMapTo returns a function that replaces the value inside an Option with a constant.
//
// Example:
//
//	replaceWith42 := OptionMapTo[string, int](42)
//	result := replaceWith42(OptionSome("hello")) // OptionSome(42)
func OptionMapTo[A, B any](b B) OptionOperator[A, B] {
	return func(fa Option[A]) Option[B] {
		if fa.s {
			return OptionSome(b)
		}
		return OptionNone[B]()
	}
}

// OptionTryCatch executes a function that may return an error and converts the result to an Option.
// Returns OptionSome(value) if no error occurred, None if an error occurred.
//
// Example:
//
//	result := OptionTryCatch(func() (int, error) {
//	    return strconv.Atoi("42")
//	}) // OptionSome(42)
func OptionTryCatch[A any](f func() (A, error)) Option[A] {
	val, err := f()
	if err != nil {
		return OptionNone[A]()
	}
	return OptionSome(val)
}

// OptionFold provides a way to handle both Some and None cases of an Option.
// Returns a function that applies onNone if the Option is None, or onSome if it is Some.
//
// Example:
//
//	handler := OptionFold(
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	)
//	result := handler(OptionSome(42))    // "value: 42"
//	result := handler(OptionNone[int]()) // "no value"
//
// Relation to predicate.OptionFold:
//
// option.OptionFold and predicate.OptionFold are two specialisations of the same categorical
// pattern — eliminating a two-case sum type into a common result type B.
//
// option.Option[A] is a two-case sum type {None, OptionSome(A)}.  The Some constructor
// carries a payload of type A, so the onSome branch receives it; the None branch
// carries no payload and is therefore a thunk:
//
//	option.OptionFold :: (() → B) → (A → B) → Option[A] → B
//
// predicate.Predicate[A] is morally equivalent to A → bool, where bool is the
// smallest two-case sum type {false, true}.  Because bool carries no payload beyond
// the branch tag, both handlers in predicate.OptionFold must receive A to preserve context:
//
//	predicate.OptionFold :: (A → B) → (A → B) → (A → bool) → A → B
//
// The link between the two is FromPredicate, which converts a Predicate[A] into an
// Option[A]-producing function.  Using it, predicate.OptionFold can always be expressed in
// terms of option.OptionFold:
//
//	predicate.OptionFold(onFalse, onTrue)(p)(a)
//	  == option.OptionFold(func() B { return onFalse(a) }, onTrue)(FromPredicate(p)(a))
//
// Conversely, option.OptionFold cannot in general be expressed via predicate.OptionFold because
// None carries no A value for the false branch to inspect.
//
// See Also:
//   - MonadFold: The uncurried form of this function
//   - predicate.OptionFold: The analogous eliminator for the bool two-case sum type
//   - FromPredicate: Converts a Predicate[A] into a OptionKleisli[A, A] producing Option[A]
//
//go:inline
func OptionFold[A, B any](onNone func() B, onSome func(a A) B) func(Option[A]) B {
	return F.Bind23of3(OptionMonadFold[A, B])(onNone, onSome)
}

// OptionMonadGetOrElse extracts the value from an Option or returns a default value.
// This is the monadic form of GetOrElse.
//
// Example:
//
//	result := OptionMonadGetOrElse(OptionSome(42), lazy.Of(0)) // 42
//	result := OptionMonadGetOrElse(OptionNone[int](), lazy.Of(0)) // 0
//
//go:inline
func OptionMonadGetOrElse[A any](fa Option[A], onNone func() A) A {
	return OptionMonadFold(fa, onNone, F.Identity[A])
}

// OptionGetOrElse returns a function that extracts the value from an Option or returns a default.
//
// Example:
//
//	getOrZero := OptionGetOrElse(lazy.Of(0))
//	result := getOrZero(OptionSome(42)) // 42
//	result := getOrZero(OptionNone[int]()) // 0
//
//go:inline
func OptionGetOrElse[A any](onNone func() A) func(Option[A]) A {
	return OptionFold(onNone, F.Identity[A])
}

// OptionMonadChain applies a function that returns an Option to the value inside an Option.
// This is the monadic bind operation. If the input is None, returns None.
//
// Example:
//
//	fa := OptionSome(5)
//	result := OptionMonadChain(fa, F.Flow2(Predicate(N.MoreThan(0)), Map(N.Mul(2)))) // OptionSome(10)
//
//go:inline
func OptionMonadChain[A, B any](fa Option[A], f OptionKleisli[A, B]) Option[B] {
	return OptionMonadFold(fa, OptionNone[B], f)
}

// OptionChain returns a function that applies an Option-returning function to an Option value.
// This is the curried form of the monadic bind operation.
//
// Example:
//
//	validate := OptionChain(F.Flow2(Predicate(N.MoreThan(0)), Map(N.Mul(2))))
//	result := validate(OptionSome(5)) // OptionSome(10)
//
//go:inline
func OptionChain[A, B any](f OptionKleisli[A, B]) OptionOperator[A, B] {
	return OptionFold(OptionNone[B], f)
}

// OptionMonadChainTo ignores the first Option and returns the second Option.
// Useful for sequencing operations where the first result is not needed.
//
// Example:
//
//	result := OptionMonadChainTo(OptionSome(5), OptionSome("hello")) // OptionSome("hello")
func OptionMonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B] {
	if ma.s {
		return mb
	}
	return OptionNone[B]()
}

// OptionChainTo returns a function that ignores its input Option and returns a fixed Option.
//
// Example:
//
//	replaceWith := OptionChainTo(OptionSome("hello"))
//	result := replaceWith(OptionSome(42)) // OptionSome("hello")
func OptionChainTo[A, B any](mb Option[B]) OptionOperator[A, B] {
	if mb.s {
		return func(ma Option[A]) Option[B] {
			if ma.s {
				return mb
			}
			return OptionNone[B]()
		}
	}
	return F.Constant1[Option[A]](OptionNone[B]())
}

// OptionMonadChainFirst applies a function that returns an Option but keeps the original value.
// If either operation results in None, returns None.
//
// Example:
//
//	result := OptionMonadChainFirst(OptionSome(5), func(x int) Option[string] {
//	    return OptionSome(fmt.Sprintf("%d", x))
//	}) // OptionSome(5) - original value is kept
func OptionMonadChainFirst[A, B any](ma Option[A], f OptionKleisli[A, B]) Option[A] {
	return C.MonadChainFirst(
		OptionMonadChain[A, A],
		OptionMonadMap[B, A],
		ma,
		f,
	)
}

// OptionChainFirst returns a function that applies an Option-returning function but keeps the original value.
//
// Example:
//
//	logAndKeep := OptionChainFirst(func(x int) Option[string] {
//	    fmt.Println(x)
//	    return OptionSome("logged")
//	})
//	result := logAndKeep(OptionSome(5)) // OptionSome(5)
func OptionChainFirst[A, B any](f OptionKleisli[A, B]) OptionOperator[A, A] {
	return C.ChainFirst(
		OptionChain[A, A],
		OptionMap[B, A],
		f,
	)
}

// OptionFlatten removes one level of nesting from a nested Option.
//
// Example:
//
//	nested := OptionSome(OptionSome(42))
//	result := OptionFlatten(nested) // OptionSome(42)
//	nested := OptionSome(OptionNone[int]())
//	result := OptionFlatten(nested) // None
func OptionFlatten[A any](mma Option[Option[A]]) Option[A] {
	return OptionMonadChain(mma, F.Identity[Option[A]])
}

// MonadAlt returns the first Option if it's Some, otherwise returns the alternative.
// This is the monadic form of the Alt operation.
//
// Example:
//
//	result := MonadAlt(OptionSome(5), func() Option[int] { return OptionSome(10) }) // OptionSome(5)
//	result := MonadAlt(OptionNone[int](), func() Option[int] { return OptionSome(10) }) // OptionSome(10)
func OptionMonadAlt[A any](fa Option[A], that func() Option[A]) Option[A] {
	return OptionMonadFold(fa, that, OptionOf[A])
}

// OptionAlt returns a function that provides an alternative Option if the input is None.
//
// Example:
//
//	withDefault := OptionAlt(func() Option[int] { return OptionSome(0) })
//	result := withDefault(OptionSome(5)) // OptionSome(5)
//	result := withDefault(OptionNone[int]()) // OptionSome(0)
func OptionAlt[A any](that func() Option[A]) OptionOperator[A, A] {
	return OptionFold(that, OptionOf[A])
}

// OptionMonadSequence2 sequences two Options and applies a function to their values.
// Returns None if either Option is None.
//
// Example:
//
//	result := OptionMonadSequence2(OptionSome(2), OptionSome(3), func(a, b int) Option[int] {
//	    return OptionSome(a + b)
//	}) // OptionSome(5)
func OptionMonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R] {
	if o1.s && o2.s {
		return f(o1.a, o2.a)
	}
	return OptionNone[R]()
}

// OptionSequence2 returns a function that sequences two Options with a combining function.
//
// Example:
//
//	add := OptionSequence2(func(a, b int) Option[int] { return OptionSome(a + b) })
//	result := add(OptionSome(2), OptionSome(3)) // OptionSome(5)
func OptionSequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R] {
	return func(o1 Option[T1], o2 Option[T2]) Option[R] {
		return OptionMonadSequence2(o1, o2, f)
	}
}

// OptionReduce folds an Option into a single value using a reducer function.
// If the Option is None, returns the initial value.
//
// Example:
//
//	sum := OptionReduce(func(acc, val int) int { return acc + val }, 0)
//	result := sum(OptionSome(5)) // 5
//	result := sum(OptionNone[int]()) // 0
func OptionReduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B {
	return func(ma Option[A]) B {
		if ma.s {
			return f(initial, ma.a)
		}
		return initial
	}
}

// OptionFilter keeps the Option if it's Some and the predicate is satisfied, otherwise returns None.
//
// Example:
//
//	isPositive := OptionFilter(N.MoreThan(0))
//	result := isPositive(OptionSome(5)) // OptionSome(5)
//	result := isPositive(OptionSome(-1)) // None
//	result := isPositive(OptionNone[int]()) // None
func OptionFilter[A any](pred Predicate[A]) OptionOperator[A, A] {
	return func(ma Option[A]) Option[A] {
		if ma.s && pred(ma.a) {
			return ma
		}
		return OptionNone[A]()
	}
}

// OptionMonadFlap applies a value to a function wrapped in an Option.
// This is the monadic form of Flap.
//
// Example:
//
//	fab := OptionSome(N.Mul(2))
//	result := OptionMonadFlap(fab, 5) // OptionSome(10)
func OptionMonadFlap[B, A any](fab Option[func(A) B], a A) Option[B] {
	if fab.s {
		return OptionSome(fab.a(a))
	}
	return OptionNone[B]()
}

// OptionFlap returns a function that applies a value to an Option-wrapped function.
//
// Example:
//
//	applyFive := OptionFlap[int](5)
//	fab := OptionSome(N.Mul(2))
//	result := applyFive(fab) // OptionSome(10)
func OptionFlap[B, A any](a A) OptionOperator[func(A) B, B] {
	return func(fab Option[func(A) B]) Option[B] {
		if fab.s {
			return OptionSome(fab.a(a))
		}
		return OptionNone[B]()
	}
}

// OptionZero returns the zero value of an [Option], which is None.
// This function is useful as an identity element in monoid operations or for creating an empty Option.
//
// The zero value for Option[A] is always None, representing the absence of a value.
// This is consistent with the Option monad's semantics where None represents "no value"
// and Some represents "a value".
//
// Important: OptionZero() returns the same value as the default initialization of Option[A].
// When you declare `var o Option[A]` without initialization, it has the same value as OptionZero[A]().
//
// Note: Unlike other types where zero might be a default value, Option's zero is explicitly
// the absence of any value (None), not Some with a zero value.
//
// Example:
//
//	// OptionZero Option of any type is always None
//	o1 := option.OptionZero[int]()     // None
//	o2 := option.OptionZero[string]()  // None
//	o3 := option.OptionZero[*int]()    // None
//
//	// OptionZero equals default initialization
//	var defaultInit Option[int]
//	zero := option.OptionZero[int]()
//	assert.Equal(t, defaultInit, zero) // true
//
//	// Verify it's None
//	o := option.OptionZero[int]()
//	assert.True(t, option.IsNone(o))   // true
//	assert.False(t, option.IsOptionSome(o))  // false
//
//	// Different from Some with zero value
//	someZero := option.OptionSome(0)         // OptionSome(0)
//	zero := option.OptionZero[int]()         // None
//	assert.NotEqual(t, someZero, zero) // they are different
func OptionZero[A any]() Option[A] {
	return OptionNone[A]()
}

// ChainNone is the curried version that sequences a computation on the None (empty) value.
// Returns a function that applies the provided function when the Option is None,
// or returns None when the Option is Some.
//
// Note: ChainNone is identical to Alt - both provide the same functionality for providing
// alternative values when an Option is None.
//
// The naming convention follows the pattern established by ChainLeft in the Either type,
// where operations on the "error" or "empty" case use the Left/None suffix to distinguish
// them from operations on the "success" or "present" case (Chain operates on Some values,
// while ChainNone operates on None values).
//
// This is useful for creating reusable default value providers or transformers that can be
// composed with other Option operations using pipes or function composition.
//
// Example:
//
//	// Create a reusable default provider
//	provideDefault := option.ChainNone(func() option.Option[int] {
//	    return option.OptionSome(42)
//	})
//
//	// Use in a pipeline
//	result := F.Pipe1(
//	    option.OptionNone[int](),
//	    provideDefault,
//	) // OptionSome(42)
//
//	// Some values pass through unchanged
//	result := F.Pipe1(
//	    option.OptionSome(10),
//	    provideDefault,
//	) // OptionSome(10)
//
//go:inline
func OptionChainOptionNone[A any](onNone func() Option[A]) OptionOperator[A, A] {
	return OptionAlt(onNone)
}

func optionToType[T any](a any) (T, bool) {
	b, ok := a.(T)
	return b, ok
}

// OptionInstanceOf attempts to convert a value of type any to a specific type T using type assertion.
// Returns Some(value) if the type assertion succeeds, None if it fails.
//
// Example:
//
//	var x any = 42
//	result := OptionInstanceOf[int](x) // Some(42)
//
//	var y any = "hello"
//	result := OptionInstanceOf[int](y) // None (wrong type)
func OptionInstanceOf[T any](src any) Option[T] {
	return F.Pipe1(
		src,
		OptionOptionize1(optionToType[T]),
	)
}

// OptionToAny converts a value of any type to Option[any].
// This always succeeds and returns Some containing the value as any.
//
// Example:
//
//	result := OptionToAny(42) // Some(any(42))
//	result := OptionToAny("hello") // Some(any("hello"))
func OptionToAny[T any](src T) Option[any] {
	return OptionOf(any(src))
}
