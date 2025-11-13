package option

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	LI "github.com/IBM/fp-go/v2/optics/lens/iso"
	O "github.com/IBM/fp-go/v2/option"
)

// fromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func fromPredicate[GET ~func(S) Option[A], SET ~func(Option[A]) Endomorphism[S], S, A any](creator func(get GET, set SET) LensO[S, A], pred func(A) bool, nilValue A) func(sa Lens[S, A]) LensO[S, A] {
	fromPred := O.FromPredicate(pred)
	return func(sa Lens[S, A]) LensO[S, A] {
		return creator(F.Flow2(sa.Get, fromPred), O.Fold(F.Bind1of1(sa.Set)(nilValue), sa.Set))
	}
}

// FromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
//
//go:inline
func FromPredicate[S, A any](pred func(A) bool, nilValue A) func(sa Lens[S, A]) LensO[S, A] {
	return fromPredicate(lens.MakeLensCurried[func(S) Option[A], func(Option[A]) Endomorphism[S]], pred, nilValue)
}

// FromPredicateRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
//
//go:inline
func FromPredicateRef[S, A any](pred func(A) bool, nilValue A) func(sa Lens[*S, A]) LensO[*S, A] {
	return fromPredicate(lens.MakeLensRefCurried[S, Option[A]], pred, nilValue)
}

// FromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
//
//go:inline
func FromNillable[S, A any](sa Lens[S, *A]) LensO[S, *A] {
	return FromPredicate[S](F.IsNonNil[A], nil)(sa)
}

// FromNillableRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
//
//go:inline
func FromNillableRef[S, A any](sa Lens[*S, *A]) LensO[*S, *A] {
	return FromPredicateRef[S](F.IsNonNil[A], nil)(sa)
}

// fromNullableProp returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func fromNullableProp[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](creator func(get GET, set SET) Lens[S, A], isNullable O.Kleisli[A, A], defaultValue A) func(sa Lens[S, A]) Lens[S, A] {
	orElse := O.GetOrElse(F.Constant(defaultValue))
	return func(sa Lens[S, A]) Lens[S, A] {
		return creator(F.Flow3(
			sa.Get,
			isNullable,
			orElse,
		), sa.Set)
	}
}

// FromNullableProp returns a `Lens` from a property that may be optional. The getter returns a default value for these items
//
//go:inline
func FromNullableProp[S, A any](isNullable O.Kleisli[A, A], defaultValue A) lens.Operator[S, A, A] {
	return fromNullableProp(lens.MakeLensCurried[func(S) A, func(A) Endomorphism[S]], isNullable, defaultValue)
}

// FromNullablePropRef returns a `Lens` from a property that may be optional. The getter returns a default value for these items
//
//go:inline
func FromNullablePropRef[S, A any](isNullable O.Kleisli[A, A], defaultValue A) lens.Operator[*S, A, A] {
	return fromNullableProp(lens.MakeLensRefCurried[S, A], isNullable, defaultValue)
}

// fromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func fromOption[GET ~func(S) A, SET ~func(A) Endomorphism[S], S, A any](creator func(get GET, set SET) Lens[S, A], defaultValue A) func(LensO[S, A]) Lens[S, A] {
	orElse := O.GetOrElse(F.Constant(defaultValue))
	return func(sa LensO[S, A]) Lens[S, A] {
		return creator(F.Flow2(
			sa.Get,
			orElse,
		), F.Flow2(O.Of[A], sa.Set))
	}
}

// FromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
//
//go:inline
func FromOption[S, A any](defaultValue A) func(LensO[S, A]) Lens[S, A] {
	return fromOption(lens.MakeLensCurried[func(S) A, func(A) Endomorphism[S]], defaultValue)
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
//
//go:inline
func FromOptionRef[S, A any](defaultValue A) func(LensO[*S, A]) Lens[*S, A] {
	return fromOption(lens.MakeLensRefCurried[S, A], defaultValue)
}

// FromIso converts a Lens[S, A] to a LensO[S, A] using an isomorphism.
//
// This function takes an isomorphism between A and Option[A] and uses it to
// transform a regular lens into an optional lens. It's particularly useful when
// you have a custom isomorphism that defines how to convert between a value
// and its optional representation.
//
// The isomorphism must satisfy the round-trip laws:
//  1. iso.ReverseGet(iso.Get(a)) == a for all a: A
//  2. iso.Get(iso.ReverseGet(opt)) == opt for all opt: Option[A]
//
// Type Parameters:
//   - S: The structure type containing the field
//   - A: The type of the field being focused on
//
// Parameters:
//   - iso: An isomorphism between A and Option[A] that defines the conversion
//
// Returns:
//   - A function that takes a Lens[S, A] and returns a LensO[S, A]
//
// Example:
//
//	type Config struct {
//	    timeout int
//	}
//
//	// Create a lens to the timeout field
//	timeoutLens := lens.MakeLens(
//	    func(c Config) int { return c.timeout },
//	    func(c Config, t int) Config { c.timeout = t; return c },
//	)
//
//	// Create an isomorphism that treats 0 as None
//	zeroAsNone := iso.MakeIso(
//	    func(t int) option.Option[int] {
//	        if t == 0 {
//	            return option.None[int]()
//	        }
//	        return option.Some(t)
//	    },
//	    func(opt option.Option[int]) int {
//	        return option.GetOrElse(func() int { return 0 })(opt)
//	    },
//	)
//
//	// Convert to optional lens
//	optTimeoutLens := FromIso[Config, int](zeroAsNone)(timeoutLens)
//
//	config := Config{timeout: 0}
//	opt := optTimeoutLens.Get(config)        // None[int]()
//	updated := optTimeoutLens.Set(option.Some(30))(config) // Config{timeout: 30}
//
// Common Use Cases:
//   - Converting between sentinel values (like 0, -1, "") and Option
//   - Applying custom validation logic when converting to/from Option
//   - Integrating with existing isomorphisms like FromNillable
//
// See also:
//   - FromPredicate: For predicate-based optional conversion
//   - FromNillable: For pointer-based optional conversion
//   - FromOption: For converting from optional to non-optional with defaults
//
//go:inline
func FromIso[S, A any](iso Iso[A, Option[A]]) func(Lens[S, A]) LensO[S, A] {
	return LI.Compose[S](iso)
}
