// Copyright (c) 2023 - 2025 IBM Corp.
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

// Package optional provides an optic for focusing on values that may not exist.
//
// # Overview
//
// Optional is an optic used to zoom inside a product. Unlike the Lens, the element that the Optional focuses
// on may not exist. An Optional[S, A] represents a relationship between a source type S and a focus type A,
// where the focus may or may not be present.
//
// # Optional Laws
//
// An Optional must satisfy the following laws, which are consistent with other functional programming libraries
// such as monocle-ts (https://gcanti.github.io/monocle-ts/modules/Optional.ts.html) and the Haskell lens library
// (https://hackage.haskell.org/package/lens):
//
//  1. GetSet Law (No-op on None):
//     If GetOption(s) returns None, then Set(a)(s) must return s unchanged (no-op).
//     This ensures that attempting to update a value that doesn't exist has no effect.
//
//     Formally: GetOption(s) = None => Set(a)(s) = s
//
//  2. SetGet Law (Get what you Set):
//     If GetOption(s) returns Some(_), then GetOption(Set(a)(s)) must return Some(a).
//     This ensures that after setting a value, you can retrieve it.
//
//     Formally: GetOption(s) = Some(_) => GetOption(Set(a)(s)) = Some(a)
//
//  3. SetSet Law (Last Set Wins):
//     Setting twice is the same as setting once with the final value.
//
//     Formally: Set(b)(Set(a)(s)) = Set(b)(s)
//
// # No-op Behavior
//
// A key property of Optional is that updating a value for which GetOption returns None is a no-op.
// This behavior is implemented through the optionalModify function, which only applies the modification
// if the optional value exists. When GetOption returns None, the original structure is returned unchanged.
//
// This is consistent with the behavior in:
//   - monocle-ts: Optional.modify returns the original value when the optional doesn't match
//   - Haskell lens: over and set operations are no-ops when the traversal finds no targets
//
// # Example
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Create an optional that focuses on non-empty names
//	nameOptional := MakeOptional(
//	    func(p Person) option.Option[string] {
//	        if p.Name != "" {
//	            return option.Some(p.Name)
//	        }
//	        return option.None[string]()
//	    },
//	    func(p Person, name string) Person {
//	        p.Name = name
//	        return p
//	    },
//	)
//
//	// When the optional matches, Set updates the value
//	person1 := Person{Name: "Alice", Age: 30}
//	updated1 := nameOptional.Set("Bob")(person1)
//	// updated1.Name == "Bob"
//
//	// When the optional doesn't match (Name is empty), Set is a no-op
//	person2 := Person{Name: "", Age: 30}
//	updated2 := nameOptional.Set("Bob")(person2)
//	// updated2 == person2 (unchanged)
package optional

import (
	EM "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

type (
	// Optional is an optional reference to a subpart of a data type
	Optional[S, A any] struct {
		GetOption func(s S) O.Option[A]
		Set       func(a A) EM.Endomorphism[S]
		name      string
	}

	// Kleisli represents a function that takes a value of type A and returns an Optional[S, B].
	// This is commonly used for composing optionals in a monadic style.
	//
	// Type Parameters:
	//   - S: The source type of the resulting optional
	//   - A: The input type to the function
	//   - B: The focus type of the resulting optional
	Kleisli[S, A, B any] = func(A) Optional[S, B]

	// Operator represents a function that transforms one optional into another.
	// It takes an Optional[S, A] and returns an Optional[S, B], allowing for optional transformations.
	//
	// Type Parameters:
	//   - S: The source type (remains constant)
	//   - A: The original focus type
	//   - B: The new focus type
	Operator[S, A, B any] = func(Optional[S, A]) Optional[S, B]
)

// setCopyRef wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopyRef[SET ~func(A) func(*S) *S, S, A any](setter SET) func(a A) func(*S) *S {
	return func(a A) func(*S) *S {

		sa := setter(a)

		return func(s *S) *S {
			if s == nil {
				return s
			}
			cpy := *s
			return sa(&cpy)
		}
	}
}

func getRef[GET ~func(*S) O.Option[A], S, A any](getter GET) func(*S) O.Option[A] {
	return func(s *S) O.Option[A] {
		if s == nil {
			return O.None[A]()
		}
		return getter(s)
	}
}

// MakeOptional creates an Optional based on a getter and a setter function. Make sure that the setter creates a (shallow) copy of the
// data. This happens automatically if the data is passed by value. For pointers consider to use `MakeOptionalRef`
// and for other kinds of data structures that are copied by reference make sure the setter creates the copy.
//
//go:inline
func MakeOptional[S, A any](get O.Kleisli[S, A], set func(S, A) S) Optional[S, A] {
	return MakeOptionalWithName(get, set, "GenericOptional")
}

//go:inline
func MakeOptionalCurried[S, A any](get O.Kleisli[S, A], set func(A) func(S) S) Optional[S, A] {
	return MakeOptionalCurriedWithName(get, set, "GenericOptional")
}

//go:inline
func MakeOptionalWithName[S, A any](get O.Kleisli[S, A], set func(S, A) S, name string) Optional[S, A] {
	return MakeOptionalCurriedWithName(get, F.Bind2of2(set), name)
}

func MakeOptionalCurriedWithName[S, A any](get O.Kleisli[S, A], set func(A) func(S) S, name string) Optional[S, A] {
	return Optional[S, A]{GetOption: get, Set: set, name: name}
}

// MakeOptionalRef creates an Optional based on a getter and a setter function. The setter passed in does not have to create a shallow
// copy, the implementation wraps the setter into one that copies the pointer before modifying it
//
//go:inline
func MakeOptionalRef[S, A any](get O.Kleisli[*S, A], set func(*S, A) *S) Optional[*S, A] {
	return MakeOptionalCurried(getRef(get), setCopyRef(F.Bind2of2(set)))
}

//go:inline
func MakeOptionalRefWithName[S, A any](get O.Kleisli[*S, A], set func(*S, A) *S, name string) Optional[*S, A] {
	return MakeOptionalCurriedWithName(getRef(get), setCopyRef(F.Bind2of2(set)), name)
}

//go:inline
func MakeOptionalRefCurriedWithName[S, A any](get O.Kleisli[*S, A], set func(A) func(*S) *S, name string) Optional[*S, A] {
	return MakeOptionalCurriedWithName(getRef(get), setCopyRef(set), name)
}

// Id returns am optional implementing the identity operation
func idWithName[S any](creator func(get O.Kleisli[S, S], set func(S, S) S, name string) Optional[S, S], name string) Optional[S, S] {
	return creator(O.Some[S], F.Second[S, S], name)
}

// Id returns am optional implementing the identity operation
func Id[S any]() Optional[S, S] {
	return idWithName(MakeOptionalWithName[S, S], "Identity")
}

// Id returns am optional implementing the identity operation
func IdRef[S any]() Optional[*S, *S] {
	return idWithName(MakeOptionalRefWithName[S, *S], "Identity")
}

func optionalModifyOption[S, A any](f func(A) A, optional Optional[S, A], s S) O.Option[S] {
	return F.Pipe1(
		optional.GetOption(s),
		O.Map(func(a A) S {
			return optional.Set(f(a))(s)
		}),
	)
}

func optionalModify[S, A any](f func(A) A, optional Optional[S, A], s S) S {
	return F.Pipe1(
		optionalModifyOption(f, optional, s),
		O.GetOrElse(F.Constant(s)),
	)
}

// Compose combines two Optional and allows to narrow down the focus to a sub-Optional
func compose[S, A, B any](creator func(get O.Kleisli[S, B], set func(S, B) S) Optional[S, B], ab Optional[A, B]) Operator[S, A, B] {
	abget := ab.GetOption
	abset := ab.Set
	return func(sa Optional[S, A]) Optional[S, B] {
		saget := sa.GetOption
		return creator(
			F.Flow2(saget, O.Chain(abget)),
			func(s S, b B) S {
				return optionalModify(abset(b), sa, s)
			},
		)
	}
}

// Compose combines two Optional and allows to narrow down the focus to a sub-Optional
func Compose[S, A, B any](ab Optional[A, B]) Operator[S, A, B] {
	return compose(MakeOptional[S, B], ab)
}

// ComposeRef combines two Optional and allows to narrow down the focus to a sub-Optional
func ComposeRef[S, A, B any](ab Optional[A, B]) Operator[*S, A, B] {
	return compose(MakeOptionalRef[S, B], ab)
}

// fromPredicate implements the function generically for both the ref and the direct case
func fromPredicate[S, A any](creator func(get O.Kleisli[S, A], set func(S, A) S) Optional[S, A], pred func(A) bool) func(func(S) A, func(S, A) S) Optional[S, A] {
	fromPred := O.FromPredicate(pred)
	return func(get func(S) A, set func(S, A) S) Optional[S, A] {
		return creator(
			F.Flow2(get, fromPred),
			func(s S, a A) S {
				return F.Pipe3(
					s,
					get,
					fromPred,
					O.Fold(F.Constant(s), func(_ A) S {
						return set(s, a)
					}),
				)
			},
		)
	}
}

// FromPredicate creates an optional from getter and setter functions. It checks
// for optional values and the correct update procedure
func FromPredicate[S, A any](pred func(A) bool) func(func(S) A, func(S, A) S) Optional[S, A] {
	return fromPredicate(MakeOptional[S, A], pred)
}

// FromPredicate creates an optional from getter and setter functions. It checks
// for optional values and the correct update procedure
func FromPredicateRef[S, A any](pred func(A) bool) func(func(*S) A, func(*S, A) *S) Optional[*S, A] {
	return fromPredicate(MakeOptionalRef[S, A], pred)
}

func imap[S, A, B any](sa Optional[S, A], ab func(A) B, ba func(B) A) Optional[S, B] {
	return MakeOptional(
		F.Flow2(sa.GetOption, O.Map(ab)),
		func(s S, b B) S {
			return sa.Set(ba(b))(s)
		},
	)
}

// IMap implements a bidirectional mapping of the transform
func IMap[S, A, B any](ab func(A) B, ba func(B) A) Operator[S, A, B] {
	return func(sa Optional[S, A]) Optional[S, B] {
		return imap(sa, ab, ba)
	}
}

func ModifyOption[S, A any](f func(A) A) func(Optional[S, A]) O.Kleisli[S, S] {
	return func(o Optional[S, A]) O.Kleisli[S, S] {
		return func(s S) O.Option[S] {
			return optionalModifyOption(f, o, s)
		}
	}
}

func SetOption[S, A any](a A) func(Optional[S, A]) O.Kleisli[S, S] {
	return ModifyOption[S](F.Constant1[A](a))
}

func ichain[S, A, B any](sa Optional[S, A], ab O.Kleisli[A, B], ba O.Kleisli[B, A]) Optional[S, B] {
	return MakeOptional(
		F.Flow2(sa.GetOption, O.Chain(ab)),
		func(s S, b B) S {
			return O.MonadFold(ba(b), EM.Identity[S], sa.Set)(s)
		},
	)
}

// IChain implements a bidirectional mapping of the transform if the transform can produce optionals (e.g. in case of type mappings)
func IChain[S, A, B any](ab O.Kleisli[A, B], ba O.Kleisli[B, A]) Operator[S, A, B] {
	return func(sa Optional[S, A]) Optional[S, B] {
		return ichain(sa, ab, ba)
	}
}

// IChainAny implements a bidirectional mapping to and from any
func IChainAny[S, A any]() Operator[S, any, A] {
	fromAny := O.InstanceOf[A]
	toAny := O.ToAny[A]
	return func(sa Optional[S, any]) Optional[S, A] {
		return ichain(sa, fromAny, toAny)
	}
}
