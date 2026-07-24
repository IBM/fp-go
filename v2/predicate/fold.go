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

package predicate

// Fold evaluates a predicate against a value and maps the boolean result to a value of type B.
//
// Given two mapping functions and a predicate, Fold returns a function that, when applied
// to a value of type A, tests it with the predicate and calls onTrue if the predicate
// returns true, or onFalse if it returns false. Both branches receive the original input
// value, so contextual information is preserved regardless of which branch is taken.
//
// Type Parameters:
//   - A: The input type tested by the predicate
//   - B: The output type produced by both branch functions
//
// Parameters:
//   - onFalse: Called with the input value when the predicate returns false
//   - onTrue:  Called with the input value when the predicate returns true
//
// Returns:
//   - A function that takes a Predicate[A] and returns a function from A to B
//
// Relation to option.Fold:
//
// predicate.Fold and option.Fold are two specialisations of the same categorical
// pattern — eliminating a two-case sum type into a common result type B.
//
// A Predicate[A] is morally equivalent to a function A → bool, where bool is the
// smallest two-case sum type {false, true}.  Because bool carries no payload beyond
// the branch tag, both handlers must receive A to preserve context:
//
//	predicate.Fold :: (A → B) → (A → B) → (A → bool) → A → B
//
// option.Option[A] is a richer two-case sum type {None, Some(A)}.  Here the Some
// constructor already carries the payload, so only the onSome branch needs A;
// onNone is a thunk:
//
//	option.Fold :: (() → B) → (A → B) → Option[A] → B
//
// The link between the two is option.FromPredicate, which converts a Predicate[A]
// into an Option[A]-producing function.  Using it, predicate.Fold can always be
// expressed in terms of option.Fold and FromPredicate:
//
//	predicate.Fold(onFalse, onTrue)(p)(a)
//	  == option.Fold(func() B { return onFalse(a) }, onTrue)(option.FromPredicate(p)(a))
//
// Conversely, option.Fold cannot in general be expressed via predicate.Fold because
// option.None carries no A value for the onFalse branch to inspect.
//
// See Also:
//   - Predicate: The boolean-valued function type used as the condition
//   - option.Fold: The richer analogue that eliminates an Option[A]
//   - option.FromPredicate: Converts a Predicate[A] into a Kleisli[A, A] over Option
func Fold[A, B any](onFalse, onTrue func(A) B) func(Predicate[A]) func(A) B {
	return func(p Predicate[A]) func(A) B {
		return func(a A) B {
			if p(a) {
				return onTrue(a)
			}
			return onFalse(a)
		}
	}
}
