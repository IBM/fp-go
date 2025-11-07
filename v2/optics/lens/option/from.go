package option

import (
	EM "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
)

// fromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func fromPredicate[GET ~func(S) Option[A], SET ~func(S, Option[A]) S, S, A any](creator func(get GET, set SET) LensO[S, A], pred func(A) bool, nilValue A) func(sa Lens[S, A]) LensO[S, A] {
	fromPred := O.FromPredicate(pred)
	return func(sa Lens[S, A]) LensO[S, A] {
		fold := O.Fold(F.Bind1of1(sa.Set)(nilValue), sa.Set)
		return creator(F.Flow2(sa.Get, fromPred), func(s S, a Option[A]) S {
			return F.Pipe2(
				a,
				fold,
				EM.Ap(s),
			)
		})
	}
}

// FromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func FromPredicate[S, A any](pred func(A) bool, nilValue A) func(sa Lens[S, A]) LensO[S, A] {
	return fromPredicate(lens.MakeLens[func(S) Option[A], func(S, Option[A]) S], pred, nilValue)
}

// FromPredicateRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func FromPredicateRef[S, A any](pred func(A) bool, nilValue A) func(sa Lens[*S, A]) Lens[*S, Option[A]] {
	return fromPredicate(lens.MakeLensRef[func(*S) Option[A], func(*S, Option[A]) *S], pred, nilValue)
}

// FromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
func FromNillable[S, A any](sa Lens[S, *A]) Lens[S, Option[*A]] {
	return FromPredicate[S](F.IsNonNil[A], nil)(sa)
}

// FromNillableRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
func FromNillableRef[S, A any](sa Lens[*S, *A]) Lens[*S, Option[*A]] {
	return FromPredicateRef[S](F.IsNonNil[A], nil)(sa)
}

// fromNullableProp returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func fromNullableProp[GET ~func(S) A, SET ~func(S, A) S, S, A any](creator func(get GET, set SET) Lens[S, A], isNullable func(A) Option[A], defaultValue A) func(sa Lens[S, A]) Lens[S, A] {
	return func(sa Lens[S, A]) Lens[S, A] {
		return creator(F.Flow3(
			sa.Get,
			isNullable,
			O.GetOrElse(F.Constant(defaultValue)),
		), func(s S, a A) S {
			return sa.Set(a)(s)
		},
		)
	}
}

// FromNullableProp returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func FromNullableProp[S, A any](isNullable func(A) Option[A], defaultValue A) func(sa Lens[S, A]) Lens[S, A] {
	return fromNullableProp(lens.MakeLens[func(S) A, func(S, A) S], isNullable, defaultValue)
}

// FromNullablePropRef returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func FromNullablePropRef[S, A any](isNullable func(A) Option[A], defaultValue A) func(sa Lens[*S, A]) Lens[*S, A] {
	return fromNullableProp(lens.MakeLensRef[func(*S) A, func(*S, A) *S], isNullable, defaultValue)
}

// fromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func fromOption[GET ~func(S) A, SET ~func(S, A) S, S, A any](creator func(get GET, set SET) Lens[S, A], defaultValue A) func(sa LensO[S, A]) Lens[S, A] {
	return func(sa LensO[S, A]) Lens[S, A] {
		return creator(F.Flow2(
			sa.Get,
			O.GetOrElse(F.Constant(defaultValue)),
		), func(s S, a A) S {
			return sa.Set(O.Some(a))(s)
		},
		)
	}
}

// FromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func FromOption[S, A any](defaultValue A) func(sa LensO[S, A]) Lens[S, A] {
	return fromOption(lens.MakeLens[func(S) A, func(S, A) S], defaultValue)
}

// FromOptionRef creates a lens from an Option property with a default value for pointer structures.
//
// This is the pointer version of [FromOption], with automatic copying to ensure immutability.
// The getter returns the value inside Some[A], or the defaultValue if it's None[A].
// The setter always wraps the value in Some[A].
//
// Type Parameters:
//   - S: Structure type (will be used as *S)
//   - A: Focus type
//
// Parameters:
//   - defaultValue: Value to return when the Option is None
//
// Returns:
//   - A function that takes a Lens[*S, Option[A]] and returns a Lens[*S, A]
func FromOptionRef[S, A any](defaultValue A) func(sa Lens[*S, Option[A]]) Lens[*S, A] {
	return fromOption(lens.MakeLensRef[func(*S) A, func(*S, A) *S], defaultValue)
}
