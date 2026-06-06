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

package option

import L "github.com/IBM/fp-go/v2/optics/lens"

// MonadChainIdiomatic applies a Go-idiomatic function (returning (B, bool)) to the value
// inside an Option. Returns None if the input is None or the function returns false.
//
// [KleisliIdiomatic][A, B] is func(A) (B, bool), the conventional Go pattern for
// functions that may fail without an error value (e.g., map lookups, type assertions).
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	MonadChainIdiomatic(Some("42"), parse) // Some(42)
//	MonadChainIdiomatic(Some("x"), parse)  // None
//	MonadChainIdiomatic(None[string](), parse) // None
func MonadChainIdiomatic[A, B any](fa Option[A], f KleisliIdiomatic[A, B]) Option[B] {
	return MonadChain(fa, FromValidation(f))
}

// ChainIdiomatic returns a function that applies a Go-idiomatic function (returning (B, bool))
// to the value inside an Option. Returns None if the input is None or the function returns false.
//
// This is the curried form of [MonadChainIdiomatic].
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	ChainIdiomatic(parse)(Some("42")) // Some(42)
//	ChainIdiomatic(parse)(None[string]()) // None
func ChainIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Operator[A, B] {
	return Fold(None[B], FromValidation(f))
}

// MonadChainFirstIdiomatic applies a Go-idiomatic function (returning (B, bool)) to the value
// inside an Option, but discards the function's result and keeps the original value on success.
// Returns None if the input is None or the function returns false.
//
// Example:
//
//	validate := func(n int) (string, bool) {
//	    if n > 0 { return strconv.Itoa(n), true }
//	    return "", false
//	}
//	MonadChainFirstIdiomatic(Some(5), validate)  // Some(5) — original value kept
//	MonadChainFirstIdiomatic(Some(-1), validate) // None
//	MonadChainFirstIdiomatic(None[int](), validate) // None
func MonadChainFirstIdiomatic[A, B any](ma Option[A], f KleisliIdiomatic[A, B]) Option[A] {
	return MonadChainFirst(ma, FromValidation(f))
}

// ChainFirstIdiomatic returns a function that applies a Go-idiomatic function (returning (B, bool))
// but keeps the original Option value if the function succeeds.
// Returns None if the input is None or the function returns false.
//
// This is the curried form of [MonadChainFirstIdiomatic].
//
// Example:
//
//	validate := func(n int) (string, bool) {
//	    if n > 0 { return strconv.Itoa(n), true }
//	    return "", false
//	}
//	ChainFirstIdiomatic(validate)(Some(5))  // Some(5)
//	ChainFirstIdiomatic(validate)(None[int]()) // None
func ChainFirstIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Operator[A, A] {
	return ChainFirst(FromValidation(f))
}

// BindIdiomatic attaches the result of a Go-idiomatic computation (returning (A, bool))
// to a context S1 to produce a context S2, using the given setter.
// Returns None if the input is None or the computation returns false.
//
// This is the idiomatic-function variant of [Bind], suitable for use in do-notation
// style pipelines when the computation follows the Go (value, bool) convention.
//
// Example:
//
//	res := F.Pipe3(
//	    Do(utils.Empty),
//	    BindIdiomatic(utils.SetLastName, func(utils.Initial) (string, bool) {
//	        return "Doe", true
//	    }),
//	    BindIdiomatic(utils.SetGivenName, func(utils.WithLastName) (string, bool) {
//	        return "John", true
//	    }),
//	    Map(utils.GetFullName),
//	) // Some("John Doe")
func BindIdiomatic[S1, S2, A any](
	setter func(A) func(S1) S2,
	f KleisliIdiomatic[S1, A],
) Operator[S1, S2] {
	return Bind(setter, FromValidation(f))
}

// BindLIdiomatic attaches the result of a Go-idiomatic computation (returning (T, bool))
// to a context using a lens-based setter. The computation receives the current value at
// the lens focus and returns a new value of the same type.
// Returns None if the input is None or the computation returns false.
//
// This is the idiomatic-function variant of [BindL].
//
// Example:
//
//	type Counter struct { Value int }
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//	increment := func(v int) (int, bool) {
//	    if v >= 100 { return 0, false }
//	    return v + 1, true
//	}
//	BindLIdiomatic(valueLens, increment)(Some(Counter{Value: 42})) // Some(Counter{Value: 43})
//	BindLIdiomatic(valueLens, increment)(Some(Counter{Value: 100})) // None
func BindLIdiomatic[S, T any](
	lens L.Lens[S, T],
	f KleisliIdiomatic[T, T],
) Operator[S, S] {
	return BindL(lens, FromValidation(f))
}

// TraverseIterIdiomatic transforms a sequence by applying a Go-idiomatic function
// (returning (B, bool)) to each element.
// Returns Some containing a sequence of results if all calls return true, None if any returns false.
//
// This is the idiomatic-function variant of [TraverseIter].
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	TraverseIterIdiomatic(parse)(slices.Values([]string{"1", "2", "3"})) // Some(iter{1,2,3})
//	TraverseIterIdiomatic(parse)(slices.Values([]string{"1", "x"}))      // None
func TraverseIterIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Kleisli[Seq[A], Seq[B]] {
	return TraverseIter(FromValidation(f))
}

// TraverseArrayIdiomatic transforms a slice by applying a Go-idiomatic function
// (returning (B, bool)) to each element.
// Returns Some containing the slice of results if all calls return true, None if any returns false.
//
// This is the idiomatic-function variant of [TraverseArray].
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	TraverseArrayIdiomatic(parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
//	TraverseArrayIdiomatic(parse)([]string{"1", "x", "3"}) // None
func TraverseArrayIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Kleisli[[]A, []B] {
	return TraverseArray(FromValidation(f))
}
