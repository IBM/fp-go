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

// MonadChainI applies a Go-idiomatic function (returning (B, bool)) to the value
// inside an Option. Returns None if the input is None or the function returns false.
//
// [KleisliI][A, B] is func(A) (B, bool), the conventional Go pattern for
// functions that may fail without an error value (e.g., map lookups, type assertions).
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	MonadChainI(Some("42"), parse) // Some(42)
//	MonadChainI(Some("x"), parse)  // None
//	MonadChainI(None[string](), parse) // None
func MonadChainI[A, B any](fa Option[A], f KleisliI[A, B]) Option[B] {
	return MonadChain(fa, FromValidation(f))
}

// ChainI returns a function that applies a Go-idiomatic function (returning (B, bool))
// to the value inside an Option. Returns None if the input is None or the function returns false.
//
// This is the curried form of [MonadChainI].
//
// Example:
//
//	parse := func(s string) (int, bool) {
//	    n, err := strconv.Atoi(s)
//	    return n, err == nil
//	}
//	ChainI(parse)(Some("42")) // Some(42)
//	ChainI(parse)(None[string]()) // None
func ChainI[A, B any](f KleisliI[A, B]) Operator[A, B] {
	return Fold(None[B], FromValidation(f))
}

// MonadChainFirstI applies a Go-idiomatic function (returning (B, bool)) to the value
// inside an Option, but discards the function's result and keeps the original value on success.
// Returns None if the input is None or the function returns false.
//
// Example:
//
//	validate := func(n int) (string, bool) {
//	    if n > 0 { return strconv.Itoa(n), true }
//	    return "", false
//	}
//	MonadChainFirstI(Some(5), validate)  // Some(5) — original value kept
//	MonadChainFirstI(Some(-1), validate) // None
//	MonadChainFirstI(None[int](), validate) // None
func MonadChainFirstI[A, B any](ma Option[A], f KleisliI[A, B]) Option[A] {
	return MonadChainFirst(ma, FromValidation(f))
}

// ChainFirstI returns a function that applies a Go-idiomatic function (returning (B, bool))
// but keeps the original Option value if the function succeeds.
// Returns None if the input is None or the function returns false.
//
// This is the curried form of [MonadChainFirstI].
//
// Example:
//
//	validate := func(n int) (string, bool) {
//	    if n > 0 { return strconv.Itoa(n), true }
//	    return "", false
//	}
//	ChainFirstI(validate)(Some(5))  // Some(5)
//	ChainFirstI(validate)(None[int]()) // None
func ChainFirstI[A, B any](f KleisliI[A, B]) Operator[A, A] {
	return ChainFirst(FromValidation(f))
}

// BindI attaches the result of a Go-idiomatic computation (returning (A, bool))
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
//	    BindI(utils.SetLastName, func(utils.Initial) (string, bool) {
//	        return "Doe", true
//	    }),
//	    BindI(utils.SetGivenName, func(utils.WithLastName) (string, bool) {
//	        return "John", true
//	    }),
//	    Map(utils.GetFullName),
//	) // Some("John Doe")
func BindI[S1, S2, A any](
	setter func(A) func(S1) S2,
	f KleisliI[S1, A],
) Operator[S1, S2] {
	return Bind(setter, FromValidation(f))
}

// BindIL attaches the result of a Go-idiomatic computation (returning (T, bool))
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
//	BindIL(valueLens, increment)(Some(Counter{Value: 42})) // Some(Counter{Value: 43})
//	BindIL(valueLens, increment)(Some(Counter{Value: 100})) // None
func BindIL[S, T any](
	lens L.Lens[S, T],
	f KleisliI[T, T],
) Operator[S, S] {
	return BindL(lens, FromValidation(f))
}

// TraverseIterI transforms a sequence by applying a Go-idiomatic function
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
//	TraverseIterI(parse)(slices.Values([]string{"1", "2", "3"})) // Some(iter{1,2,3})
//	TraverseIterI(parse)(slices.Values([]string{"1", "x"}))      // None
func TraverseIterI[A, B any](f KleisliI[A, B]) Kleisli[Seq[A], Seq[B]] {
	return TraverseIter(FromValidation(f))
}

// TraverseArrayI transforms a slice by applying a Go-idiomatic function
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
//	TraverseArrayI(parse)([]string{"1", "2", "3"}) // Some([1, 2, 3])
//	TraverseArrayI(parse)([]string{"1", "x", "3"}) // None
func TraverseArrayI[A, B any](f KleisliI[A, B]) Kleisli[[]A, []B] {
	return TraverseArray(FromValidation(f))
}
