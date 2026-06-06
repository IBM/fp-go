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

// MonadChainIdiomatic applies a Go-idiomatic function (returning (B, error)) to the Right
// value inside a Result. Returns Left(err) if the input is Left or the function returns a
// non-nil error.
//
// [KleisliIdiomatic][A, B] is func(A) (B, error), the standard Go convention for functions
// that may fail. This function adapts that convention for use in Result pipelines.
//
// Example:
//
//	parse := strconv.Atoi // func(string) (int, error)
//	MonadChainIdiomatic(Right("42"), parse) // Right(42)
//	MonadChainIdiomatic(Right("x"), parse)  // Left(error)
//	MonadChainIdiomatic(Left[string](errors.New("earlier failure")), parse) // Left(error)
func MonadChainIdiomatic[A, B any](fa Result[A], f KleisliIdiomatic[A, B]) Result[B] {
	return MonadChain(fa, Eitherize1(f))
}

// ChainIdiomatic returns a function that applies a Go-idiomatic function (returning (B, error))
// to the Right value inside a Result. Returns Left(err) if the input is Left or the function
// returns a non-nil error.
//
// This is the curried form of [MonadChainIdiomatic].
//
// Example:
//
//	parse := ChainIdiomatic(strconv.Atoi)
//	parse(Right("42")) // Right(42)
//	parse(Right("x"))  // Left(error)
//	parse(Left[string](errors.New("fail"))) // Left(error)
func ChainIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Operator[A, B] {
	return Chain(Eitherize1(f))
}

// MonadChainLeftIdiomatic applies a Go-idiomatic function (returning (A, error)) to the Left
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
//	MonadChainLeftIdiomatic(Left[int](errors.New("not found")), recover) // Right(0)
//	MonadChainLeftIdiomatic(Left[int](errors.New("other")), recover)     // Left(other)
//	MonadChainLeftIdiomatic(Right(42), recover)                          // Right(42)
func MonadChainLeftIdiomatic[A any](fa Result[A], f KleisliIdiomatic[error, A]) Result[A] {
	return MonadChainLeft(fa, Eitherize1(f))
}

// ChainLeftIdiomatic returns a function that applies a Go-idiomatic function (returning (A, error))
// to the Left (error) value inside a Result. Right values pass through unchanged.
//
// This is the curried form of [MonadChainLeftIdiomatic].
//
// Example:
//
//	recover := func(err error) (int, error) {
//	    if err.Error() == "not found" { return 0, nil }
//	    return 0, err
//	}
//	ChainLeftIdiomatic(recover)(Left[int](errors.New("not found"))) // Right(0)
//	ChainLeftIdiomatic(recover)(Right(42))                          // Right(42)
func ChainLeftIdiomatic[A any](f KleisliIdiomatic[error, A]) Kleisli[Result[A], A] {
	return ChainLeft(Eitherize1(f))
}

// MonadChainFirstIdiomatic applies a Go-idiomatic function (returning (B, error)) to the Right
// value inside a Result, but discards the function's result and keeps the original value on success.
// Returns Left(err) if the input is Left or the function returns a non-nil error.
//
// Example:
//
//	validate := func(n int) (string, error) {
//	    if n > 0 { return strconv.Itoa(n), nil }
//	    return "", errors.New("non-positive")
//	}
//	MonadChainFirstIdiomatic(Right(5), validate)  // Right(5) — original value kept
//	MonadChainFirstIdiomatic(Right(-1), validate) // Left(error)
//	MonadChainFirstIdiomatic(Left[int](errors.New("prior")), validate) // Left(prior)
func MonadChainFirstIdiomatic[A, B any](ma Result[A], f KleisliIdiomatic[A, B]) Result[A] {
	return MonadChainFirst(ma, Eitherize1(f))
}

// ChainFirstIdiomatic returns a function that applies a Go-idiomatic function (returning (B, error))
// but keeps the original Right value on success.
// Returns Left(err) if the input is Left or the function returns a non-nil error.
//
// This is the curried form of [MonadChainFirstIdiomatic].
//
// Example:
//
//	validate := func(n int) (string, error) {
//	    if n > 0 { return strconv.Itoa(n), nil }
//	    return "", errors.New("non-positive")
//	}
//	ChainFirstIdiomatic(validate)(Right(5))  // Right(5)
//	ChainFirstIdiomatic(validate)(Right(-1)) // Left(error)
//	ChainFirstIdiomatic(validate)(Left[int](errors.New("prior"))) // Left(prior)
func ChainFirstIdiomatic[A, B any](f KleisliIdiomatic[A, B]) Operator[A, A] {
	return ChainFirst(Eitherize1(f))
}
