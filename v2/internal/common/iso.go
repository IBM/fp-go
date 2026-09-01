package common

import (
	"fmt"

	EM "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
)

type (
	// isoTag is a zero-size non-generic type embedded in Iso[S, A] so that
	// String, Format, and LogValue are compiled once and shared across all
	// type-parameter instantiations, rather than being duplicated for every
	// distinct Iso[S, A].
	isoTag struct{}

	// Iso represents an isomorphism between types S and A.
	// An isomorphism is a bidirectional transformation that converts between two types
	// without any loss of information. It consists of two functions that are inverses
	// of each other.
	//
	// Type Parameters:
	//   - S: The source type
	//   - A: The target type
	//
	// Fields:
	//   - Get: Converts from S to A
	//   - ReverseGet: Converts from A back to S
	//
	// Laws:
	// An Iso must satisfy the round-trip laws:
	//  1. ReverseGet(Get(s)) == s for all s: S
	//  2. Get(ReverseGet(a)) == a for all a: A
	//
	// Example:
	//
	//	// Isomorphism between Celsius and Fahrenheit
	//	tempIso := Iso[float64, float64]{
	//	    Get: func(c float64) float64 { return c*9/5 + 32 },
	//	    ReverseGet: func(f float64) float64 { return (f - 32) * 5 / 9 },
	//	}
	//
	//	fahrenheit := tempIso.Get(20.0)        // 68.0
	//	celsius := tempIso.ReverseGet(68.0)    // 20.0
	Iso[S, A any] struct {
		isoTag

		// Get converts a value from the source type S to the target type A.
		Get func(s S) A

		// ReverseGet converts a value from the target type A back to the source type S.
		// This is the inverse of Get.
		ReverseGet func(a A) S
	}
)

// MakeIso constructs an isomorphism from two functions.
// The functions should be inverses of each other to satisfy the isomorphism laws.
//
// Type Parameters:
//   - S: The source type
//   - A: The target type
//
// Parameters:
//   - get: Function to convert from S to A
//   - reverse: Function to convert from A to S (inverse of get)
//
// Returns:
//   - An Iso[S, A] that uses the provided functions
//
// Example:
//
//	// Create an isomorphism between string and []byte
//	stringBytesIso := MakeIso(
//	    func(s string) []byte { return []byte(s) },
//	    func(b []byte) string { return string(b) },
//	)
//
//	bytes := stringBytesIso.Get("hello")           // []byte("hello")
//	str := stringBytesIso.ReverseGet([]byte("hi")) // "hi"
func MakeIso[S, A any](get func(S) A, reverse func(A) S) Iso[S, A] {
	return Iso[S, A]{Get: get, ReverseGet: reverse}
}

// IsoId returns an identity isomorphism that performs no transformation.
// Both Get and ReverseGet are the identity function.
//
// Type Parameters:
//   - S: The type for both source and target
//
// Returns:
//   - An Iso[S, S] where Get and ReverseGet are both identity functions
//
// Example:
//
//	idIso := IsoId[int]()
//	value := idIso.Get(42)        // 42
//	same := idIso.ReverseGet(42)  // 42
//
// Use cases:
//   - As a starting point for isomorphism composition
//   - When you need an isomorphism but don't want to transform the value
//   - In generic code that requires an isomorphism parameter
func IsoId[S any]() Iso[S, S] {
	return MakeIso(F.Identity[S], F.Identity[S])
}

// IsoComposeIso combines two isomorphisms to create a new isomorphism.
// Given Iso[S, A] and Iso[A, B], creates Iso[S, B].
// The resulting isomorphism first applies the outer iso (S → A),
// then the inner iso (A → B).
//
// Type Parameters:
//   - S: The outermost source type
//   - A: The intermediate type
//   - B: The innermost target type
//
// Parameters:
//   - ab: The inner isomorphism (A → B)
//
// Returns:
//   - A function that takes the outer isomorphism (S → A) and returns the composed isomorphism (S → B)
//
// Example:
//
//	metersToKm := MakeIso(
//	    func(m float64) float64 { return m / 1000 },
//	    func(km float64) float64 { return km * 1000 },
//	)
//
//	kmToMiles := MakeIso(
//	    func(km float64) float64 { return km * 0.621371 },
//	    func(mi float64) float64 { return mi / 0.621371 },
//	)
//
//	// IsoComposeIso: meters → kilometers → miles
//	metersToMiles := F.Pipe1(metersToKm, IsoComposeIso[float64](kmToMiles))
//
//	miles := metersToMiles.Get(5000)        // ~3.11 miles
//	meters := metersToMiles.ReverseGet(3.11) // ~5000 meters
func IsoComposeIso[S, A, B any](ab Iso[A, B]) func(Iso[S, A]) Iso[S, B] {
	return func(sa Iso[S, A]) Iso[S, B] {
		return MakeIso(
			F.Flow2(sa.Get, ab.Get),
			F.Flow2(ab.ReverseGet, sa.ReverseGet),
		)
	}
}

// IsoComposeLens combines an isomorphism with a lens to produce a new lens.
//
// Given an isomorphism Iso[S, A] and a lens Lens[A, B], IsoComposeLens creates a
// Lens[S, B]. Internally the isomorphism is first converted to a Lens[S, A] via
// IsoAsLens, then composed with the provided lens using LensComposeLens.
//
// This is useful when the outer view of a structure is captured by a total,
// invertible transformation (an iso) and you want to drill further into the result
// with a lens — for example, when S is a newtype wrapper around A and you need to
// focus on a field of A.
//
// The resulting lens satisfies all three lens laws whenever both the isomorphism
// and the inner lens individually satisfy them.
//
// The composition follows: (sa ∘ ab).Get(s) = ab.Get(sa.Get(s))
//
// Type Parameters:
//   - S: Outermost source type — the input to the resulting lens
//   - A: Intermediate type — the target of the iso and the source of the inner lens
//   - B: Inner focus type — the focus of the resulting lens
//
// Parameters:
//   - ab: Lens from A to B (inner lens)
//
// Returns:
//   - A function that takes an Iso[S, A] and returns a Lens[S, B]
//
// Example:
//
//	type Celsius float64
//	type Thermometer struct{ Reading Celsius }
//
//	// Isomorphism: float64 ↔ Celsius
//	floatToCelsius := MakeIso(
//	    func(f float64) Celsius { return Celsius(f) },
//	    func(c Celsius) float64 { return float64(c) },
//	)
//
//	// Lens: Thermometer → Celsius (focuses on the Reading field)
//	readingLens := MakeLens(
//	    func(t Thermometer) Celsius { return t.Reading },
//	    func(t Thermometer, c Celsius) Thermometer { t.Reading = c; return t },
//	)
//
//	// Compose: float64 → Celsius → Thermometer.Reading (Celsius)
//	// Result type: Lens[float64, Celsius]
//	tempLens := F.Pipe1(floatToCelsius, IsoComposeLens[float64](readingLens))
//
//	t := Thermometer{Reading: 100}
//	reading := tempLens.Get(100.0)          // Celsius(100)
//	updated := tempLens.Set(Celsius(0))(0.0) // Thermometer{Reading: 0}
//
// See Also:
//   - IsoComposeIso: the variant that composes with an Iso[A, B] instead of a Lens[A, B]
//   - LensComposeLens: the pure-lens equivalent
//   - IsoAsLens: converts an Iso to a Lens directly
//
//go:inline
func IsoComposeLens[S, A, B any](ab Lens[A, B]) func(Iso[S, A]) Lens[S, B] {
	return F.Flow2(
		IsoAsLens[S],
		LensComposeLens[S](ab),
	)
}

// IsoReverse swaps the direction of an isomorphism.
// Given Iso[S, A], creates Iso[A, S] where Get and ReverseGet are swapped.
//
// Type Parameters:
//   - S: The original source type (becomes target)
//   - A: The original target type (becomes source)
//
// Parameters:
//   - sa: The isomorphism to reverse
//
// Returns:
//   - An Iso[A, S] with Get and ReverseGet swapped
//
// Example:
//
//	celsiusToFahrenheit := MakeIso(
//	    func(c float64) float64 { return c*9/5 + 32 },
//	    func(f float64) float64 { return (f - 32) * 5 / 9 },
//	)
//
//	// IsoReverse to get Fahrenheit to Celsius
//	fahrenheitToCelsius := IsoReverse(celsiusToFahrenheit)
//
//	celsius := fahrenheitToCelsius.Get(68.0)        // 20.0
//	fahrenheit := fahrenheitToCelsius.ReverseGet(20.0) // 68.0
func IsoReverse[S, A any](sa Iso[S, A]) Iso[A, S] {
	return MakeIso(
		sa.ReverseGet,
		sa.Get,
	)
}

// isoModify is an internal helper that applies a transformation function through an isomorphism.
// It converts S to A, applies the function, then converts back to S.
func isoModify[FCT ~func(A) A, S, A any](f FCT, sa Iso[S, A], s S) S {
	return F.Pipe3(
		s,
		sa.Get,
		f,
		sa.ReverseGet,
	)
}

// IsoModify creates a function that applies a transformation in the target space.
// It converts the source value to the target type, applies the transformation,
// then converts back to the source type.
//
// Type Parameters:
//   - S: The source type
//   - FCT: The transformation function type (A → A)
//   - A: The target type
//
// Parameters:
//   - f: The transformation function to apply in the target space
//
// Returns:
//   - A function that takes an Iso[S, A] and returns an endomorphism (S → S)
//
// Example:
//
//	type Meters float64
//	type Kilometers float64
//
//	mToKm := MakeIso(
//	    func(m Meters) Kilometers { return Kilometers(m / 1000) },
//	    func(km Kilometers) Meters { return Meters(km * 1000) },
//	)
//
//	// Double the distance in kilometers, result in meters
//	doubled := IsoModify[Meters](func(km Kilometers) Kilometers {
//	    return km * 2
//	})(mToKm)(Meters(5000))
//	// Result: Meters(10000)
func IsoModify[S any, FCT ~func(A) A, A any](f FCT) func(Iso[S, A]) Endomorphism[S] {
	return F.Curry3(isoModify[FCT, S, A])(f)
}

// IsoUnwrap extracts the target value from a source value using an isomorphism.
// This is a convenience function that applies the Get function of the isomorphism.
//
// Type Parameters:
//   - A: The target type to extract
//   - S: The source type
//
// Parameters:
//   - s: The source value to unwrap
//
// Returns:
//   - A function that takes an Iso[S, A] and returns the unwrapped value of type A
//
// Example:
//
//	type UserId int
//
//	userIdIso := MakeIso(
//	    func(id UserId) int { return int(id) },
//	    func(i int) UserId { return UserId(i) },
//	)
//
//	rawId := IsoUnwrap[int](UserId(42))(userIdIso) // 42
//
// Note: This function is also available as To for semantic clarity.
func IsoUnwrap[A, S any](s S) func(Iso[S, A]) A {
	return func(sa Iso[S, A]) A {
		return sa.Get(s)
	}
}

// IsoWrap wraps a target value into a source value using an isomorphism.
// This is a convenience function that applies the ReverseGet function of the isomorphism.
//
// Type Parameters:
//   - S: The source type to wrap into
//   - A: The target type
//
// Parameters:
//   - a: The target value to wrap
//
// Returns:
//   - A function that takes an Iso[S, A] and returns the wrapped value of type S
//
// Example:
//
//	type UserId int
//
//	userIdIso := MakeIso(
//	    func(id UserId) int { return int(id) },
//	    func(i int) UserId { return UserId(i) },
//	)
//
//	userId := IsoWrap[UserId](42)(userIdIso) // UserId(42)
//
// Note: This function is also available as From for semantic clarity.
func IsoWrap[S, A any](a A) func(Iso[S, A]) S {
	return func(sa Iso[S, A]) S {
		return sa.ReverseGet(a)
	}
}

// IsoTo extracts the target value from a source value using an isomorphism.
// This is an alias for Unwrap, provided for semantic clarity when the
// direction of conversion is important.
//
// Type Parameters:
//   - A: The target type to convert to
//   - S: The source type
//
// Parameters:
//   - s: The source value to convert
//
// Returns:
//   - A function that takes an Iso[S, A] and returns the converted value of type A
//
// Example:
//
//	type Email string
//	type ValidatedEmail struct{ value Email }
//
//	emailIso := MakeIso(
//	    func(ve ValidatedEmail) Email { return ve.value },
//	    func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
//	)
//
//	// Convert to Email
//	email := IsoTo[Email](ValidatedEmail{value: "user@example.com"})(emailIso)
//	// "user@example.com"
func IsoTo[A, S any](s S) func(Iso[S, A]) A {
	return IsoUnwrap[A](s)
}

// IsoFrom wraps a target value into a source value using an isomorphism.
// This is an alias for Wrap, provided for semantic clarity when the
// direction of conversion is important.
//
// Type Parameters:
//   - S: The source type to convert from
//   - A: The target type
//
// Parameters:
//   - a: The target value to convert
//
// Returns:
//   - A function that takes an Iso[S, A] and returns the converted value of type S
//
// Example:
//
//	type Email string
//	type ValidatedEmail struct{ value Email }
//
//	emailIso := MakeIso(
//	    func(ve ValidatedEmail) Email { return ve.value },
//	    func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
//	)
//
//	// Convert from Email
//	validated := IsoFrom[ValidatedEmail](Email("admin@example.com"))(emailIso)
//	// ValidatedEmail{value: "admin@example.com"}
func IsoFrom[S, A any](a A) func(Iso[S, A]) S {
	return IsoWrap[S](a)
}

// isoImap is an internal helper that bidirectionally maps an isomorphism.
// It transforms both directions of the isomorphism using the provided functions.
func isoImap[S, A, B any](sa Iso[S, A], ab func(A) B, ba func(B) A) Iso[S, B] {
	return MakeIso(
		F.Flow2(sa.Get, ab),
		F.Flow2(ba, sa.ReverseGet),
	)
}

// IsoIMap bidirectionally maps the target type of an isomorphism.
// Given Iso[S, A] and functions A → B and B → A, creates Iso[S, B].
// This allows you to transform both directions of an isomorphism.
//
// Type Parameters:
//   - S: The source type (unchanged)
//   - A: The original target type
//   - B: The new target type
//
// Parameters:
//   - ab: Function to map from A to B
//   - ba: Function to map from B to A (inverse of ab)
//
// Returns:
//   - A function that transforms Iso[S, A] to Iso[S, B]
//
// Example:
//
//	type Celsius float64
//	type Kelvin float64
//
//	celsiusIso := Id[Celsius]()
//
//	// Create isomorphism to Kelvin
//	celsiusToKelvin := F.Pipe1(
//	    celsiusIso,
//	    IsoIMap(
//	        func(c Celsius) Kelvin { return Kelvin(c + 273.15) },
//	        func(k Kelvin) Celsius { return Celsius(k - 273.15) },
//	    ),
//	)
//
//	kelvin := celsiusToKelvin.Get(Celsius(20))      // 293.15 K
//	celsius := celsiusToKelvin.ReverseGet(Kelvin(293.15)) // 20°C
//
// Note: The functions ab and ba must be inverses of each other to maintain
// the isomorphism laws.
func IsoIMap[S, A, B any](ab func(A) B, ba func(B) A) func(Iso[S, A]) Iso[S, B] {
	return func(sa Iso[S, A]) Iso[S, B] {
		return isoImap(sa, ab, ba)
	}
}

// IsoAsLens converts an Iso[S, A] into a Lens[S, A].
//
// The resulting lens uses the iso's Get as its getter and its ReverseGet as its setter.
// Because an isomorphism is total and invertible, Set always replaces the focused value
// by applying ReverseGet — the original S is discarded.
//
// Type Parameters:
//   - S: The source type of the isomorphism and of the resulting lens
//   - A: The target type of the isomorphism and the focus type of the resulting lens
//
// Parameters:
//   - sa: The isomorphism to convert
//
// Returns:
//   - Lens[S, A]: a lens whose Get is sa.Get and whose Set is sa.ReverseGet
//
// See Also:
//   - IsoAsLensRef: the pointer-receiver variant for Iso[*S, A]
func IsoAsLens[S, A any](sa Iso[S, A]) Lens[S, A] {
	return MakeLensCurriedWithName(sa.Get, F.Flow2(sa.ReverseGet, F.Flow2(F.Constant1[S, S], EM.Of[func(S) S])), fmt.Sprintf("FromIso[%s]", sa))
}

// IsoAsLensRef converts an Iso[*S, A] into a Lens[*S, A] for pointer-based structures.
//
// This is the pointer-receiver variant of IsoAsLens. It accepts an isomorphism whose
// source type is already a pointer (*S) and wraps the resulting lens with copy-on-write
// semantics: before the setter is applied, a shallow copy of *S is made so that the
// original value is never mutated. A nil pointer is handled safely by substituting a
// zero-value S.
//
// In practice the copy-on-write wrapper is a no-op for an iso-derived setter, because
// the setter is simply ReverseGet(a) and ignores the incoming *S entirely. The wrapper
// exists for consistency with other pointer lenses and to guard against nil inputs.
//
// Deprecated: IsoAsLens[*S, A] is sufficient because the copy-on-write wrapper added
// by this function has no observable effect for iso-derived setters. Prefer IsoAsLens
// unless you specifically need the nil-pointer safety on the incoming *S, which the
// iso's ReverseGet already cannot depend on.
//
// Type Parameters:
//   - S: The pointee type; the isomorphism source is *S
//   - A: The target type of the isomorphism and the focus type of the resulting lens
//
// Parameters:
//   - sa: The isomorphism from *S to A to convert
//
// Returns:
//   - Lens[*S, A]: a lens whose Get is sa.Get and whose Set applies sa.ReverseGet
//     with copy-on-write protection on the *S receiver
//
// See Also:
//   - IsoAsLens: the general form that accepts any Iso[S, A], including Iso[*S, A]
func IsoAsLensRef[S, A any](sa Iso[*S, A]) Lens[*S, A] {
	return MakeLensRefCurriedWithName(sa.Get, F.Flow2(sa.ReverseGet, F.Flow2(F.Constant1[*S, *S], EM.Of[func(*S) *S])), fmt.Sprintf("FromIso[%s]", sa))
}
