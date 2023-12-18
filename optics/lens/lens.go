// Copyright (c) 2023 IBM Corp.
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

// Lens is an optic used to zoom inside a product.
package lens

import (
	EM "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

type (
	// Lens is a reference to a subpart of a data type
	Lens[S, A any] struct {
		Get func(s S) A
		Set func(a A) EM.Endomorphism[S]
	}
)

// setCopy wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopy[SET ~func(*S, A) *S, S, A any](setter SET) func(s *S, a A) *S {
	return func(s *S, a A) *S {
		copy := *s
		return setter(&copy, a)
	}
}

// setCopyCurried wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopyCurried[SET ~func(A) EM.Endomorphism[*S], S, A any](setter SET) func(a A) EM.Endomorphism[*S] {
	return func(a A) EM.Endomorphism[*S] {
		seta := setter(a)
		return func(s *S) *S {
			copy := *s
			return seta(&copy)
		}
	}
}

// MakeLens creates a [Lens] based on a getter and a setter function. Make sure that the setter creates a (shallow) copy of the
// data. This happens automatically if the data is passed by value. For pointers consider to use `MakeLensRef`
// and for other kinds of data structures that are copied by reference make sure the setter creates the copy.
func MakeLens[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET) Lens[S, A] {
	return MakeLensCurried(get, EM.Curry2(F.Swap(set)))
}

// MakeLensCurried creates a [Lens] based on a getter and a setter function. Make sure that the setter creates a (shallow) copy of the
// data. This happens automatically if the data is passed by value. For pointers consider to use `MakeLensRef`
// and for other kinds of data structures that are copied by reference make sure the setter creates the copy.
func MakeLensCurried[GET ~func(S) A, SET ~func(A) EM.Endomorphism[S], S, A any](get GET, set SET) Lens[S, A] {
	return Lens[S, A]{Get: get, Set: set}
}

// MakeLensRef creates a [Lens] based on a getter and a setter function. The setter passed in does not have to create a shallow
// copy, the implementation wraps the setter into one that copies the pointer before modifying it
//
// Such a [Lens] assumes that property A of S always exists
func MakeLensRef[GET ~func(*S) A, SET func(*S, A) *S, S, A any](get GET, set SET) Lens[*S, A] {
	return MakeLens(get, setCopy(set))
}

// MakeLensRefCurried creates a [Lens] based on a getter and a setter function. The setter passed in does not have to create a shallow
// copy, the implementation wraps the setter into one that copies the pointer before modifying it
//
// Such a [Lens] assumes that property A of S always exists
func MakeLensRefCurried[S, A any](get func(*S) A, set func(A) EM.Endomorphism[*S]) Lens[*S, A] {
	return MakeLensCurried(get, setCopyCurried(set))
}

// id returns a [Lens] implementing the identity operation
func id[GET ~func(S) S, SET ~func(S, S) S, S any](creator func(get GET, set SET) Lens[S, S]) Lens[S, S] {
	return creator(F.Identity[S], F.Second[S, S])
}

// Id returns a [Lens] implementing the identity operation
func Id[S any]() Lens[S, S] {
	return id(MakeLens[EM.Endomorphism[S], func(S, S) S])
}

// IdRef returns a [Lens] implementing the identity operation
func IdRef[S any]() Lens[*S, *S] {
	return id(MakeLensRef[EM.Endomorphism[*S], func(*S, *S) *S])
}

// Compose combines two lenses and allows to narrow down the focus to a sub-lens
func compose[GET ~func(S) B, SET ~func(S, B) S, S, A, B any](creator func(get GET, set SET) Lens[S, B], ab Lens[A, B]) func(Lens[S, A]) Lens[S, B] {
	abget := ab.Get
	abset := ab.Set
	return func(sa Lens[S, A]) Lens[S, B] {
		saget := sa.Get
		saset := sa.Set
		return creator(
			F.Flow2(saget, abget),
			func(s S, b B) S {
				return saset(abset(b)(saget(s)))(s)
			},
		)
	}
}

// Compose combines two lenses and allows to narrow down the focus to a sub-lens
func Compose[S, A, B any](ab Lens[A, B]) func(Lens[S, A]) Lens[S, B] {
	return compose(MakeLens[func(S) B, func(S, B) S], ab)
}

// ComposeOption combines a `Lens` that returns an optional value with a `Lens` that returns a definite value
// the getter returns an `Option[B]` because the container `A` could already be an option
// if the setter is invoked with `Some[B]` then the value of `B` will be set, potentially on a default value of `A` if `A` did not exist
// if the setter is invoked with `None[B]` then the container `A` is reset to `None[A]` because this is the only way to remove `B`
func ComposeOption[S, B, A any](defaultA A) func(ab Lens[A, B]) func(Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
	defa := F.Constant(defaultA)
	return func(ab Lens[A, B]) func(Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
		foldab := O.Fold(O.None[B], F.Flow2(ab.Get, O.Some[B]))
		return func(sa Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
			// set A on S
			seta := F.Flow2(
				O.Some[A],
				sa.Set,
			)
			// remove A from S
			unseta := F.Nullary2(
				O.None[A],
				sa.Set,
			)
			return MakeLens(
				F.Flow2(sa.Get, foldab),
				func(s S, ob O.Option[B]) S {
					return F.Pipe2(
						ob,
						O.Fold(unseta, func(b B) EM.Endomorphism[S] {
							setbona := F.Flow2(
								ab.Set(b),
								seta,
							)
							return F.Pipe2(
								s,
								sa.Get,
								O.Fold(
									F.Nullary2(
										defa,
										setbona,
									),
									setbona,
								),
							)
						}),
						EM.Ap(s),
					)
				},
			)
		}
	}
}

// ComposeOptions combines a `Lens` that returns an optional value with a `Lens` that returns another optional value
// the getter returns `None[B]` if either `A` or `B` is `None`
// if the setter is called with `Some[B]` and `A` exists, 'A' is updated with `B`
// if the setter is called with `Some[B]` and `A` does not exist, the default of 'A' is updated with `B`
// if the setter is called with `None[B]` and `A` does not exist this is the identity operation on 'S'
// if the setter is called with `None[B]` and `A` does exist, 'B' is removed from 'A'
func ComposeOptions[S, B, A any](defaultA A) func(ab Lens[A, O.Option[B]]) func(Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
	defa := F.Constant(defaultA)
	noops := EM.Identity[S]
	noneb := O.None[B]()
	return func(ab Lens[A, O.Option[B]]) func(Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
		unsetb := ab.Set(noneb)
		return func(sa Lens[S, O.Option[A]]) Lens[S, O.Option[B]] {
			// sets an A onto S
			seta := F.Flow2(
				O.Some[A],
				sa.Set,
			)
			return MakeLensCurried(
				F.Flow2(
					sa.Get,
					O.Chain(ab.Get),
				),
				func(b O.Option[B]) EM.Endomorphism[S] {
					return func(s S) S {
						return O.MonadFold(b, func() EM.Endomorphism[S] {
							return F.Pipe2(
								s,
								sa.Get,
								O.Fold(noops, F.Flow2(unsetb, seta)),
							)
						}, func(b B) EM.Endomorphism[S] {
							// sets a B onto an A
							setb := F.Flow2(
								ab.Set(O.Some(b)),
								seta,
							)
							return F.Pipe2(
								s,
								sa.Get,
								O.Fold(F.Nullary2(defa, setb), setb),
							)
						})(s)
					}
				},
			)
		}
	}
}

// Compose combines two lenses and allows to narrow down the focus to a sub-lens
func ComposeRef[S, A, B any](ab Lens[A, B]) func(Lens[*S, A]) Lens[*S, B] {
	return compose(MakeLensRef[func(*S) B, func(*S, B) *S], ab)
}

func modify[FCT ~func(A) A, S, A any](f FCT, sa Lens[S, A], s S) S {
	return sa.Set(f(sa.Get(s)))(s)
}

// Modify changes a property of a [Lens] by invoking a transformation function
// if the transformed property has not changes, the method returns the original state
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Lens[S, A]) EM.Endomorphism[S] {
	return EM.Curry3(modify[FCT, S, A])(f)
}

func IMap[E any, AB ~func(A) B, BA ~func(B) A, A, B any](ab AB, ba BA) func(Lens[E, A]) Lens[E, B] {
	return func(ea Lens[E, A]) Lens[E, B] {
		return Lens[E, B]{Get: F.Flow2(ea.Get, ab), Set: F.Flow2(ba, ea.Set)}
	}
}

// fromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func fromPredicate[GET ~func(S) O.Option[A], SET ~func(S, O.Option[A]) S, S, A any](creator func(get GET, set SET) Lens[S, O.Option[A]], pred func(A) bool, nilValue A) func(sa Lens[S, A]) Lens[S, O.Option[A]] {
	fromPred := O.FromPredicate(pred)
	return func(sa Lens[S, A]) Lens[S, O.Option[A]] {
		fold := O.Fold(F.Bind1of1(sa.Set)(nilValue), sa.Set)
		return creator(F.Flow2(sa.Get, fromPred), func(s S, a O.Option[A]) S {
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
func FromPredicate[S, A any](pred func(A) bool, nilValue A) func(sa Lens[S, A]) Lens[S, O.Option[A]] {
	return fromPredicate(MakeLens[func(S) O.Option[A], func(S, O.Option[A]) S], pred, nilValue)
}

// FromPredicateRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the nil value will be set instead
func FromPredicateRef[S, A any](pred func(A) bool, nilValue A) func(sa Lens[*S, A]) Lens[*S, O.Option[A]] {
	return fromPredicate(MakeLensRef[func(*S) O.Option[A], func(*S, O.Option[A]) *S], pred, nilValue)
}

// FromPredicate returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
func FromNillable[S, A any](sa Lens[S, *A]) Lens[S, O.Option[*A]] {
	return FromPredicate[S](F.IsNonNil[A], nil)(sa)
}

// FromNillableRef returns a `Lens` for a property accessibly as a getter and setter that can be optional
// if the optional value is set then the `nil` value will be set instead
func FromNillableRef[S, A any](sa Lens[*S, *A]) Lens[*S, O.Option[*A]] {
	return FromPredicateRef[S](F.IsNonNil[A], nil)(sa)
}

// fromNullableProp returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func fromNullableProp[GET ~func(S) A, SET ~func(S, A) S, S, A any](creator func(get GET, set SET) Lens[S, A], isNullable func(A) O.Option[A], defaultValue A) func(sa Lens[S, A]) Lens[S, A] {
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
func FromNullableProp[S, A any](isNullable func(A) O.Option[A], defaultValue A) func(sa Lens[S, A]) Lens[S, A] {
	return fromNullableProp(MakeLens[func(S) A, func(S, A) S], isNullable, defaultValue)
}

// FromNullablePropRef returns a `Lens` from a property that may be optional. The getter returns a default value for these items
func FromNullablePropRef[S, A any](isNullable func(A) O.Option[A], defaultValue A) func(sa Lens[*S, A]) Lens[*S, A] {
	return fromNullableProp(MakeLensRef[func(*S) A, func(*S, A) *S], isNullable, defaultValue)
}

// fromFromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func fromOption[GET ~func(S) A, SET ~func(S, A) S, S, A any](creator func(get GET, set SET) Lens[S, A], defaultValue A) func(sa Lens[S, O.Option[A]]) Lens[S, A] {
	return func(sa Lens[S, O.Option[A]]) Lens[S, A] {
		return creator(F.Flow2(
			sa.Get,
			O.GetOrElse(F.Constant(defaultValue)),
		), func(s S, a A) S {
			return sa.Set(O.Some(a))(s)
		},
		)
	}
}

// FromFromOption returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func FromOption[S, A any](defaultValue A) func(sa Lens[S, O.Option[A]]) Lens[S, A] {
	return fromOption(MakeLens[func(S) A, func(S, A) S], defaultValue)
}

// FromFromOptionRef returns a `Lens` from an option property. The getter returns a default value the setter will always set the some option
func FromOptionRef[S, A any](defaultValue A) func(sa Lens[*S, O.Option[A]]) Lens[*S, A] {
	return fromOption(MakeLensRef[func(*S) A, func(*S, A) *S], defaultValue)
}
