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

// Optional is an optic used to zoom inside a product. Unlike the `Lens`, the element that the `Optional` focuses
// on may not exist.
package optional

import (
	EM "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

// Optional is an optional reference to a subpart of a data type
type Optional[S, A any] struct {
	GetOption func(s S) O.Option[A]
	Set       func(a A) EM.Endomorphism[S]
}

// setCopy wraps a setter for a pointer into a setter that first creates a copy before
// modifying that copy
func setCopy[SET ~func(*S, A) *S, S, A any](setter SET) func(s *S, a A) *S {
	return func(s *S, a A) *S {
		cpy := *s
		return setter(&cpy, a)
	}
}

// MakeOptional creates an Optional based on a getter and a setter function. Make sure that the setter creates a (shallow) copy of the
// data. This happens automatically if the data is passed by value. For pointers consider to use `MakeOptionalRef`
// and for other kinds of data structures that are copied by reference make sure the setter creates the copy.
func MakeOptional[S, A any](get func(S) O.Option[A], set func(S, A) S) Optional[S, A] {
	return Optional[S, A]{GetOption: get, Set: EM.Curry2(F.Swap(set))}
}

// MakeOptionalRef creates an Optional based on a getter and a setter function. The setter passed in does not have to create a shallow
// copy, the implementation wraps the setter into one that copies the pointer before modifying it
func MakeOptionalRef[S, A any](get func(*S) O.Option[A], set func(*S, A) *S) Optional[*S, A] {
	return MakeOptional(get, setCopy(set))
}

// Id returns am optional implementing the identity operation
func id[S any](creator func(get func(S) O.Option[S], set func(S, S) S) Optional[S, S]) Optional[S, S] {
	return creator(O.Some[S], F.Second[S, S])
}

// Id returns am optional implementing the identity operation
func Id[S any]() Optional[S, S] {
	return id(MakeOptional[S, S])
}

// Id returns am optional implementing the identity operation
func IdRef[S any]() Optional[*S, *S] {
	return id(MakeOptionalRef[S, *S])
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
func compose[S, A, B any](creator func(get func(S) O.Option[B], set func(S, B) S) Optional[S, B], ab Optional[A, B]) func(Optional[S, A]) Optional[S, B] {
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
func Compose[S, A, B any](ab Optional[A, B]) func(Optional[S, A]) Optional[S, B] {
	return compose(MakeOptional[S, B], ab)
}

// ComposeRef combines two Optional and allows to narrow down the focus to a sub-Optional
func ComposeRef[S, A, B any](ab Optional[A, B]) func(Optional[*S, A]) Optional[*S, B] {
	return compose(MakeOptionalRef[S, B], ab)
}

// fromPredicate implements the function generically for both the ref and the direct case
func fromPredicate[S, A any](creator func(get func(S) O.Option[A], set func(S, A) S) Optional[S, A], pred func(A) bool) func(func(S) A, func(S, A) S) Optional[S, A] {
	fromPred := O.FromPredicate(pred)
	return func(get func(S) A, set func(S, A) S) Optional[S, A] {
		return creator(
			F.Flow2(get, fromPred),
			func(s S, a A) S {
				return F.Pipe3(
					s,
					get,
					fromPred,
					O.Fold(F.Constant(s), F.Bind1st(set, s)),
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
func IMap[S, A, B any](ab func(A) B, ba func(B) A) func(Optional[S, A]) Optional[S, B] {
	return func(sa Optional[S, A]) Optional[S, B] {
		return imap(sa, ab, ba)
	}
}

func ModifyOption[S, A any](f func(A) A) func(Optional[S, A]) func(S) O.Option[S] {
	return func(o Optional[S, A]) func(S) O.Option[S] {
		return func(s S) O.Option[S] {
			return optionalModifyOption(f, o, s)
		}
	}
}

func SetOption[S, A any](a A) func(Optional[S, A]) func(S) O.Option[S] {
	return ModifyOption[S](F.Constant1[A](a))
}

func ichain[S, A, B any](sa Optional[S, A], ab func(A) O.Option[B], ba func(B) O.Option[A]) Optional[S, B] {
	return MakeOptional(
		F.Flow2(sa.GetOption, O.Chain(ab)),
		func(s S, b B) S {
			return O.MonadFold(ba(b), EM.Identity[S], sa.Set)(s)
		},
	)
}

// IChain implements a bidirectional mapping of the transform if the transform can produce optionals (e.g. in case of type mappings)
func IChain[S, A, B any](ab func(A) O.Option[B], ba func(B) O.Option[A]) func(Optional[S, A]) Optional[S, B] {
	return func(sa Optional[S, A]) Optional[S, B] {
		return ichain(sa, ab, ba)
	}
}

// IChainAny implements a bidirectional mapping to and from any
func IChainAny[S, A any]() func(Optional[S, any]) Optional[S, A] {
	fromAny := O.ToType[A]
	toAny := O.ToAny[A]
	return func(sa Optional[S, any]) Optional[S, A] {
		return ichain(sa, fromAny, toAny)
	}
}
