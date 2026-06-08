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

package result

// MonadChainI applies a Go-idiomatic function (returning (B, error)) to the Right
// value inside a Result. Returns Left(err) if the input is Left or the function returns a
// non-nil error.
//
// [KleisliI][A, B] is func(A) (B, error), the standard Go convention for functions
// that may fail. This function adapts that convention for use in Result pipelines.
//
// Example:
//
//	parse := strconv.Atoi // func(string) (int, error)
//	MonadChainI(Right("42"), parse) // Right(42)
//	MonadChainI(Right("x"), parse)  // Left(error)
//	MonadChainI(Left[string](errors.New("earlier failure")), parse) // Left(error)
func MonadChainI[A, B any](fa Result[A], f KleisliI[A, B]) Result[B] {
	return MonadChain(fa, Eitherize1(f))
}

// ChainI returns a function that applies a Go-idiomatic function (returning (B, error))
// to the Right value inside a Result. Returns Left(err) if the input is Left or the function
// returns a non-nil error.
//
// This is the curried form of [MonadChainI].
//
// Example:
//
//	parse := ChainI(strconv.Atoi)
//	parse(Right("42")) // Right(42)
//	parse(Right("x"))  // Left(error)
//	parse(Left[string](errors.New("fail"))) // Left(error)
func ChainI[A, B any](f KleisliI[A, B]) Operator[A, B] {
	return Chain(Eitherize1(f))
}

// MonadChainLeftI applies a Go-idiomatic function (returning (A, error)) to the Left
// (error) value inside a Result. If the Result is Right, it is returned unchanged.
//
// This enables Go-idiomatic error-recovery functions to participate in Result pipelines.
// The function f receives the error and may:
//   - Return (value, nil) to recover and produce a Right
//   - Return (zero, err) to transform or propagate the error as a new Left
//
// Example:
//
//	recover := func(err error) (int, error) {
//	    if err.Error() == "not found" { return 0, nil } // recover with default
//	    return 0, err                                   // propagate other errors
//	}
//	MonadChainLeftI(Left[int](errors.New("not found")), recover) // Right(0)
//	MonadChainLeftI(Left[int](errors.New("other")), recover)     // Left(other)
//	MonadChainLeftI(Right(42), recover)                          // Right(42)
func MonadChainLeftI[A any](fa Result[A], f KleisliI[error, A]) Result[A] {
	return MonadChainLeft(fa, Eitherize1(f))
}

// ChainLeftI returns a function that applies a Go-idiomatic function (returning (A, error))
// to the Left (error) value inside a Result. Right values pass through unchanged.
//
// This is the curried form of [MonadChainLeftI].
//
// Example:
//
//	recover := func(err error) (int, error) {
//	    if err.Error() == "not found" { return 0, nil }
//	    return 0, err
//	}
//	ChainLeftI(recover)(Left[int](errors.New("not found"))) // Right(0)
//	ChainLeftI(recover)(Right(42))                          // Right(42)
func ChainLeftI[A any](f KleisliI[error, A]) Kleisli[Result[A], A] {
	return ChainLeft(Eitherize1(f))
}

// MonadChainFirstI applies a Go-idiomatic function (returning (B, error)) to the Right
// value inside a Result, but discards the function's result and keeps the original value on success.
// Returns Left(err) if the input is Left or the function returns a non-nil error.
//
// Example:
//
//	validate := func(n int) (string, error) {
//	    if n > 0 { return strconv.Itoa(n), nil }
//	    return "", errors.New("non-positive")
//	}
//	MonadChainFirstI(Right(5), validate)  // Right(5) — original value kept
//	MonadChainFirstI(Right(-1), validate) // Left(error)
//	MonadChainFirstI(Left[int](errors.New("prior")), validate) // Left(prior)
func MonadChainFirstI[A, B any](ma Result[A], f KleisliI[A, B]) Result[A] {
	return MonadChainFirst(ma, Eitherize1(f))
}

// ChainFirstI returns a function that applies a Go-idiomatic function (returning (B, error))
// but keeps the original Right value on success.
// Returns Left(err) if the input is Left or the function returns a non-nil error.
//
// This is the curried form of [MonadChainFirstI].
//
// Example:
//
//	validate := func(n int) (string, error) {
//	    if n > 0 { return strconv.Itoa(n), nil }
//	    return "", errors.New("non-positive")
//	}
//	ChainFirstI(validate)(Right(5))  // Right(5)
//	ChainFirstI(validate)(Right(-1)) // Left(error)
//	ChainFirstI(validate)(Left[int](errors.New("prior"))) // Left(prior)
func ChainFirstI[A, B any](f KleisliI[A, B]) Operator[A, A] {
	return ChainFirst(Eitherize1(f))
}
