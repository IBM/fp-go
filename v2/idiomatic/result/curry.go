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

// Curry0 converts a Go function that returns (R, error) into a curried version that returns (R, error).
//
// Example:
//
//	getConfig := func() (string, error) { return "config", nil }
//	curried := either.Curry0(getConfig)
//	result := curried() // Right("config")
func Curry0[R any](f func() (R, error)) func() (R, error) {
	return f
}

// Curry1 converts a Go function that returns (R, error) into a curried version that returns (R, error).
//
// Example:
//
//	parse := strconv.Atoi
//	curried := either.Curry1(parse)
//	result := curried("42") // Right(42)
func Curry1[T1, R any](f func(T1) (R, error)) func(T1) (R, error) {
	return f
}

// Curry2 converts a 2-argument Go function that returns (R, error) into a curried version.
//
// Example:
//
//	divide := func(a, b int) (int, error) {
//	    if b == 0 { return 0, errors.New("div by zero") }
//	    return a / b, nil
//	}
//	curried := either.Curry2(divide)
//	result := curried(10)(2) // Right(5)
func Curry2[T1, T2, R any](f func(T1, T2) (R, error)) func(T1) func(T2) (R, error) {
	return func(t1 T1) func(T2) (R, error) {
		return func(t2 T2) (R, error) {
			return f(t1, t2)
		}
	}
}

// Curry3 converts a 3-argument Go function that returns (R, error) into a curried version.
func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) (R, error)) func(T1) func(T2) func(T3) (R, error) {
	return func(t1 T1) func(T2) func(T3) (R, error) {
		return func(t2 T2) func(T3) (R, error) {
			return func(t3 T3) (R, error) {
				return f(t1, t2, t3)
			}
		}
	}
}

// Curry4 converts a 4-argument Go function that returns (R, error) into a curried version.
func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) (R, error)) func(T1) func(T2) func(T3) func(T4) (R, error) {
	return func(t1 T1) func(T2) func(T3) func(T4) (R, error) {
		return func(t2 T2) func(T3) func(T4) (R, error) {
			return func(t3 T3) func(T4) (R, error) {
				return func(t4 T4) (R, error) {
					return f(t1, t2, t3, t4)
				}
			}
		}
	}
}

// Uncurry0 converts a function returning (R, error) back to Go's (R, error) style.
//
// Example:
//
//	curried := func() either.Either[error, string] { return either.Right[error]("value") }
//	uncurried := either.Uncurry0(curried)
//	result, err := uncurried() // "value", nil
func Uncurry0[R any](f func() (R, error)) func() (R, error) {
	return func() (R, error) {
		return f()
	}
}

// Uncurry1 converts a function returning (R, error) back to Go's (R, error) style.
//
// Example:
//
//	curried := func(x int) either.Either[error, string] { return either.Right[error](strconv.Itoa(x)) }
//	uncurried := either.Uncurry1(curried)
//	result, err := uncurried(42) // "42", nil
func Uncurry1[T1, R any](f func(T1) (R, error)) func(T1) (R, error) {
	return func(t1 T1) (R, error) {
		return f(t1)
	}
}

// Uncurry2 converts a curried function returning (R, error) back to Go's (R, error) style.
func Uncurry2[T1, T2, R any](f func(T1) func(T2) (R, error)) func(T1, T2) (R, error) {
	return func(t1 T1, t2 T2) (R, error) {
		return f(t1)(t2)
	}
}

// Uncurry3 converts a curried function returning (R, error) back to Go's (R, error) style.
func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) (R, error)) func(T1, T2, T3) (R, error) {
	return func(t1 T1, t2 T2, t3 T3) (R, error) {
		return f(t1)(t2)(t3)
	}
}

// Uncurry4 converts a curried function returning (R, error) back to Go's (R, error) style.
func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) (R, error)) func(T1, T2, T3, T4) (R, error) {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
		return f(t1)(t2)(t3)(t4)
	}
}
