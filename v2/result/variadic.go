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

package result

import "github.com/IBM/fp-go/v2/either"

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
//
//go:inline
func Variadic0[V, R any](f func([]V) (R, error)) func(...V) Result[R] {
	return either.Variadic0(f)
}

// Variadic1 converts a function with 1 fixed parameter and a slice into a variadic function returning Either.
//
//go:inline
func Variadic1[T1, V, R any](f func(T1, []V) (R, error)) func(T1, ...V) Result[R] {
	return either.Variadic1(f)
}

// Variadic2 converts a function with 2 fixed parameters and a slice into a variadic function returning Either.
//
//go:inline
func Variadic2[T1, T2, V, R any](f func(T1, T2, []V) (R, error)) func(T1, T2, ...V) Result[R] {
	return either.Variadic2(f)
}

// Variadic3 converts a function with 3 fixed parameters and a slice into a variadic function returning Either.
//
//go:inline
func Variadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, []V) (R, error)) func(T1, T2, T3, ...V) Result[R] {
	return either.Variadic3(f)
}

// Variadic4 converts a function with 4 fixed parameters and a slice into a variadic function returning Either.
//
//go:inline
func Variadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, []V) (R, error)) func(T1, T2, T3, T4, ...V) Result[R] {
	return either.Variadic4(f)
}

// Unvariadic0 converts a variadic function returning (R, error) into a function taking a slice and returning Either.
//
//go:inline
func Unvariadic0[V, R any](f func(...V) (R, error)) func([]V) Result[R] {
	return either.Unvariadic0(f)
}

// Unvariadic1 converts a variadic function with 1 fixed parameter into a function taking a slice and returning Either.
//
//go:inline
func Unvariadic1[T1, V, R any](f func(T1, ...V) (R, error)) func(T1, []V) Result[R] {
	return either.Unvariadic1(f)
}

// Unvariadic2 converts a variadic function with 2 fixed parameters into a function taking a slice and returning Either.
//
//go:inline
func Unvariadic2[T1, T2, V, R any](f func(T1, T2, ...V) (R, error)) func(T1, T2, []V) Result[R] {
	return either.Unvariadic2(f)
}

// Unvariadic3 converts a variadic function with 3 fixed parameters into a function taking a slice and returning Either.
//
//go:inline
func Unvariadic3[T1, T2, T3, V, R any](f func(T1, T2, T3, ...V) (R, error)) func(T1, T2, T3, []V) Result[R] {
	return either.Unvariadic3(f)
}

// Unvariadic4 converts a variadic function with 4 fixed parameters into a function taking a slice and returning Either.
//
//go:inline
func Unvariadic4[T1, T2, T3, T4, V, R any](f func(T1, T2, T3, T4, ...V) (R, error)) func(T1, T2, T3, T4, []V) Result[R] {
	return either.Unvariadic4(f)
}
