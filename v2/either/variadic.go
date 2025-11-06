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

package either

// Variadic0 converts a function taking a slice and returning (R, error) into a variadic function returning Either.
//
// Example:
//
//	sum := func(nums []int) (int, error) {
//	    total := 0
//	    for _, n := range nums { total += n }
//	    return total, nil
//	}
//	variadicSum := either.Variadic0(sum)
//	result := variadicSum(1, 2, 3) // Right(6)
func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Either[error, R] {
	return func(v ...V) Either[error, R] {
		return TryCatchError(f(v))
	}
}

// Variadic1 converts a function with 1 fixed parameter and a slice into a variadic function returning Either.
func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Either[error, R] {
	return func(t1 T1, v ...V) Either[error, R] {
		return TryCatchError(f(t1, v))
	}
}

// Variadic2 converts a function with 2 fixed parameters and a slice into a variadic function returning Either.
func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, v))
	}
}

// Variadic3 converts a function with 3 fixed parameters and a slice into a variadic function returning Either.
func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, v))
	}
}

// Variadic4 converts a function with 4 fixed parameters and a slice into a variadic function returning Either.
func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v ...V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, t4, v))
	}
}

// Unvariadic0 converts a variadic function returning (R, error) into a function taking a slice and returning Either.
func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Either[error, R] {
	return func(v []V) Either[error, R] {
		return TryCatchError(f(v...))
	}
}

// Unvariadic1 converts a variadic function with 1 fixed parameter into a function taking a slice and returning Either.
func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Either[error, R] {
	return func(t1 T1, v []V) Either[error, R] {
		return TryCatchError(f(t1, v...))
	}
}

// Unvariadic2 converts a variadic function with 2 fixed parameters into a function taking a slice and returning Either.
func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Either[error, R] {
	return func(t1 T1, t2 T2, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, v...))
	}
}

// Unvariadic3 converts a variadic function with 3 fixed parameters into a function taking a slice and returning Either.
func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, v...))
	}
}

// Unvariadic4 converts a variadic function with 4 fixed parameters into a function taking a slice and returning Either.
func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Either[error, R] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, v []V) Either[error, R] {
		return TryCatchError(f(t1, t2, t3, t4, v...))
	}
}
