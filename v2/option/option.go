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
	"github.com/IBM/fp-go/v2/internal/common"
)

// FromPredicate returns a function that creates an Option based on a predicate.
// The returned function will wrap a value in Some if the predicate is satisfied, otherwise None.
//
// Example:
//
//	isPositive := FromPredicate(N.MoreThan(0))
//	result := isPositive(5)  // Some(5)
//	result := isPositive(-1) // None
//
//go:inline
func FromPredicate[A any](pred Predicate[A]) Kleisli[A, A] {
	return common.OptionFromPredicate(pred)
}

//go:inline
func FromZero[A comparable]() Kleisli[A, A] {
	return common.OptionFromZero[A]()
}

//go:inline
func FromNonZero[A comparable]() Kleisli[A, A] {
	return common.OptionFromNonZero[A]()
}

//go:inline
func FromEq[A any](pred eq.Eq[A]) func(A) Kleisli[A, A] {
	return common.OptionFromEq(pred)
}

//go:inline
func FromStrictEq[A comparable]() func(A) Kleisli[A, A] {
	return common.OptionFromStrictEq[A]()
}

// FromNillable converts a pointer to an Option wrapping the pointer itself.
// Returns Some(*A) when the pointer is non-nil, None otherwise.
//
// Deprecated: Use FromNillable2 instead. FromNillable2 unwraps the pointed-to
// value into an Option[A], which is the more useful form: it avoids carrying a
// *A inside the Option and composes naturally with the rest of the Option API.
//
//go:inline
func FromNillable[A any](a *A) Option[*A] {
	return common.OptionFromNillable(a)
}

// FromNillable2 converts a pointer to an Option of the pointed-to value.
// Returns Some(value) when the pointer is non-nil (dereferencing it), None otherwise.
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
//   - Some(*a) when a is non-nil
//   - None when a is nil
//
// See Also:
//   - ToNillable2: the inverse — converts Option[A] back to *A
//   - FromNillable: the older variant that wraps the pointer itself as Option[*A]
//
//go:inline
func FromNillable2[A any](a *A) Option[A] {
	return common.OptionFromNillable2(a)
}

// ToNillable2 converts an Option[A] back to a nullable pointer.
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
//
//go:inline
func ToNillable2[A any](fa Option[A]) *A {
	return common.OptionToNillable2(fa)
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
	return common.OptionFromValidation(f)
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
//
//go:inline
func MonadAp[B, A any](fab Option[func(A) B], fa Option[A]) Option[B] {
	return common.OptionMonadAp[B, A](fab, fa)
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
//
//go:inline
func Ap[B, A any](fa Option[A]) Operator[func(A) B, B] {
	return common.OptionAp[B, A](fa)
}

// MonadMap applies a function to the value inside an Option.
// If the Option is None, returns None. This is the monadic form of Map.
//
// Example:
//
//	fa := Some(5)
//	result := MonadMap(fa, N.Mul(2)) // Some(10)
//
//go:inline
func MonadMap[A, B any](fa Option[A], f func(A) B) Option[B] {
	return common.OptionMonadMap(fa, f)
}

// Map returns a function that applies a transformation to the value inside an Option.
// If the Option is None, returns None.
//
// Example:
//
//	double := Map(N.Mul(2))
//	result := double(Some(5)) // Some(10)
//	result := double(None[int]()) // None
//
//go:inline
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return common.OptionMap(f)
}

// MonadMapTo replaces the value inside an Option with a constant value.
// If the Option is None, returns None. This is the monadic form of MapTo.
//
// Example:
//
//	fa := Some(5)
//	result := MonadMapTo(fa, "hello") // Some("hello")
//
//go:inline
func MonadMapTo[A, B any](fa Option[A], b B) Option[B] {
	return common.OptionMonadMapTo(fa, b)
}

// MapTo returns a function that replaces the value inside an Option with a constant.
//
// Example:
//
//	replaceWith42 := MapTo[string, int](42)
//	result := replaceWith42(Some("hello")) // Some(42)
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return common.OptionMapTo[A, B](b)
}

// TryCatch executes a function that may return an error and converts the result to an Option.
// Returns Some(value) if no error occurred, None if an error occurred.
//
// Example:
//
//	result := TryCatch(func() (int, error) {
//	    return strconv.Atoi("42")
//	}) // Some(42)
//
//go:inline
func TryCatch[A any](f func() (A, error)) Option[A] {
	return common.OptionTryCatch(f)
}

// Fold provides a way to handle both Some and None cases of an Option.
// Returns a function that applies onNone if the Option is None, or onSome if it is Some.
//
// Example:
//
//	handler := Fold(
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	)
//	result := handler(Some(42))    // "value: 42"
//	result := handler(None[int]()) // "no value"
//
// Relation to predicate.Fold:
//
// option.Fold and predicate.Fold are two specialisations of the same categorical
// pattern — eliminating a two-case sum type into a common result type B.
//
// option.Option[A] is a two-case sum type {None, Some(A)}.  The Some constructor
// carries a payload of type A, so the onSome branch receives it; the None branch
// carries no payload and is therefore a thunk:
//
//	option.Fold :: (() → B) → (A → B) → Option[A] → B
//
// predicate.Predicate[A] is morally equivalent to A → bool, where bool is the
// smallest two-case sum type {false, true}.  Because bool carries no payload beyond
// the branch tag, both handlers in predicate.Fold must receive A to preserve context:
//
//	predicate.Fold :: (A → B) → (A → B) → (A → bool) → A → B
//
// The link between the two is FromPredicate, which converts a Predicate[A] into an
// Option[A]-producing function.  Using it, predicate.Fold can always be expressed in
// terms of option.Fold:
//
//	predicate.Fold(onFalse, onTrue)(p)(a)
//	  == option.Fold(func() B { return onFalse(a) }, onTrue)(FromPredicate(p)(a))
//
// Conversely, option.Fold cannot in general be expressed via predicate.Fold because
// None carries no A value for the false branch to inspect.
//
// See Also:
//   - MonadFold: The uncurried form of this function
//   - predicate.Fold: The analogous eliminator for the bool two-case sum type
//   - FromPredicate: Converts a Predicate[A] into a Kleisli[A, A] producing Option[A]
//
//go:inline
func Fold[A, B any](onNone func() B, onSome func(a A) B) func(Option[A]) B {
	return common.OptionFold(onNone, onSome)
}

// MonadGetOrElse extracts the value from an Option or returns a default value.
// This is the monadic form of GetOrElse.
//
// Example:
//
//	result := MonadGetOrElse(Some(42), lazy.Of(0)) // 42
//	result := MonadGetOrElse(None[int](), lazy.Of(0)) // 0
//
//go:inline
func MonadGetOrElse[A any](fa Option[A], onNone func() A) A {
	return common.OptionMonadGetOrElse(fa, onNone)
}

// GetOrElse returns a function that extracts the value from an Option or returns a default.
//
// Example:
//
//	getOrZero := GetOrElse(lazy.Of(0))
//	result := getOrZero(Some(42)) // 42
//	result := getOrZero(None[int]()) // 0
//
//go:inline
func GetOrElse[A any](onNone func() A) func(Option[A]) A {
	return common.OptionGetOrElse(onNone)
}

// MonadChain applies a function that returns an Option to the value inside an Option.
// This is the monadic bind operation. If the input is None, returns None.
//
// Example:
//
//	fa := Some(5)
//	result := MonadChain(fa, F.Flow2(Predicate(N.MoreThan(0)), Map(N.Mul(2)))) // Some(10)
//
//go:inline
func MonadChain[A, B any](fa Option[A], f Kleisli[A, B]) Option[B] {
	return common.OptionMonadChain(fa, f)
}

// Chain returns a function that applies an Option-returning function to an Option value.
// This is the curried form of the monadic bind operation.
//
// Example:
//
//	validate := Chain(F.Flow2(Predicate(N.MoreThan(0)), Map(N.Mul(2))))
//	result := validate(Some(5)) // Some(10)
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return common.OptionChain(f)
}

// MonadChainTo ignores the first Option and returns the second Option.
// Useful for sequencing operations where the first result is not needed.
//
// Example:
//
//	result := MonadChainTo(Some(5), Some("hello")) // Some("hello")
//
//go:inline
func MonadChainTo[A, B any](ma Option[A], mb Option[B]) Option[B] {
	return common.OptionMonadChainTo(ma, mb)
}

// ChainTo returns a function that ignores its input Option and returns a fixed Option.
//
// Example:
//
//	replaceWith := ChainTo(Some("hello"))
//	result := replaceWith(Some(42)) // Some("hello")
//
//go:inline
func ChainTo[A, B any](mb Option[B]) Operator[A, B] {
	return common.OptionChainTo[A, B](mb)
}

// MonadChainFirst applies a function that returns an Option but keeps the original value.
// If either operation results in None, returns None.
//
// Example:
//
//	result := MonadChainFirst(Some(5), func(x int) Option[string] {
//	    return Some(fmt.Sprintf("%d", x))
//	}) // Some(5) - original value is kept
//
//go:inline
func MonadChainFirst[A, B any](ma Option[A], f Kleisli[A, B]) Option[A] {
	return common.OptionMonadChainFirst(ma, f)
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
//
//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return common.OptionChainFirst(f)
}

// Flatten removes one level of nesting from a nested Option.
//
// Example:
//
//	nested := Some(Some(42))
//	result := Flatten(nested) // Some(42)
//	nested := Some(None[int]())
//	result := Flatten(nested) // None
//
//go:inline
func Flatten[A any](mma Option[Option[A]]) Option[A] {
	return common.OptionFlatten(mma)
}

// MonadAlt returns the first Option if it's Some, otherwise returns the alternative.
// This is the monadic form of the Alt operation.
//
// Example:
//
//	result := MonadAlt(Some(5), func() Option[int] { return Some(10) }) // Some(5)
//	result := MonadAlt(None[int](), func() Option[int] { return Some(10) }) // Some(10)
//
//go:inline
func MonadAlt[A any](fa Option[A], that func() Option[A]) Option[A] {
	return common.OptionMonadAlt(fa, that)
}

// Alt returns a function that provides an alternative Option if the input is None.
//
// Example:
//
//	withDefault := Alt(func() Option[int] { return Some(0) })
//	result := withDefault(Some(5)) // Some(5)
//	result := withDefault(None[int]()) // Some(0)
//
//go:inline
func Alt[A any](that func() Option[A]) Operator[A, A] {
	return common.OptionAlt(that)
}

// MonadSequence2 sequences two Options and applies a function to their values.
// Returns None if either Option is None.
//
// Example:
//
//	result := MonadSequence2(Some(2), Some(3), func(a, b int) Option[int] {
//	    return Some(a + b)
//	}) // Some(5)
//
//go:inline
func MonadSequence2[T1, T2, R any](o1 Option[T1], o2 Option[T2], f func(T1, T2) Option[R]) Option[R] {
	return common.OptionMonadSequence2(o1, o2, f)
}

// Sequence2 returns a function that sequences two Options with a combining function.
//
// Example:
//
//	add := Sequence2(func(a, b int) Option[int] { return Some(a + b) })
//	result := add(Some(2), Some(3)) // Some(5)
//
//go:inline
func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R] {
	return common.OptionSequence2(f)
}

// Reduce folds an Option into a single value using a reducer function.
// If the Option is None, returns the initial value.
//
// Example:
//
//	sum := Reduce(func(acc, val int) int { return acc + val }, 0)
//	result := sum(Some(5)) // 5
//	result := sum(None[int]()) // 0
//
//go:inline
func Reduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B {
	return common.OptionReduce(f, initial)
}

// Filter keeps the Option if it's Some and the predicate is satisfied, otherwise returns None.
//
// Example:
//
//	isPositive := Filter(N.MoreThan(0))
//	result := isPositive(Some(5)) // Some(5)
//	result := isPositive(Some(-1)) // None
//	result := isPositive(None[int]()) // None
//
//go:inline
func Filter[A any](pred Predicate[A]) Operator[A, A] {
	return common.OptionFilter(pred)
}

// MonadFlap applies a value to a function wrapped in an Option.
// This is the monadic form of Flap.
//
// Example:
//
//	fab := Some(N.Mul(2))
//	result := MonadFlap(fab, 5) // Some(10)
//
//go:inline
func MonadFlap[B, A any](fab Option[func(A) B], a A) Option[B] {
	return common.OptionMonadFlap[B, A](fab, a)
}

// Flap returns a function that applies a value to an Option-wrapped function.
//
// Example:
//
//	applyFive := Flap[int](5)
//	fab := Some(N.Mul(2))
//	result := applyFive(fab) // Some(10)
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return common.OptionFlap[B, A](a)
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
//
//go:inline
func Zero[A any]() Option[A] {
	return common.OptionZero[A]()
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
//	    return option.Some(42)
//	})
//
//	// Use in a pipeline
//	result := F.Pipe1(
//	    option.None[int](),
//	    provideDefault,
//	) // Some(42)
//
//	// Some values pass through unchanged
//	result := F.Pipe1(
//	    option.Some(10),
//	    provideDefault,
//	) // Some(10)
//
//go:inline
func ChainNone[A any](onNone func() Option[A]) Operator[A, A] {
	return common.OptionChainOptionNone(onNone)
}
