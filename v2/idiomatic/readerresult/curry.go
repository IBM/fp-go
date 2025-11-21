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

package readerresult

// These functions curry/uncurry Go functions with context as the first parameter into/from ReaderResult form.
// This follows the Go convention of putting context as the first parameter as advised in https://pkg.go.dev/context.
//
// Unlike the From* functions which return partially applied functions, Curry* functions return fully curried
// functions where each parameter is applied one at a time.

// Curry0 converts a context-only function into a ReaderResult (same as From0 but emphasizes immediate application).
//
// Example:
//
//	getConfig := func(ctx context.Context) (Config, error) { ... }
//	rr := readerresult.Curry0(getConfig)
//	// rr is a ReaderResult[context.Context, Config]
func Curry0[R, A any](f func(R) (A, error)) ReaderResult[R, A] {
	return f
}

// Curry1 converts a function with one parameter into a curried function returning a ReaderResult.
//
// Example:
//
//	getUser := func(ctx context.Context, id int) (User, error) { ... }
//	curried := readerresult.Curry1(getUser)
//	// curried(42) returns ReaderResult[context.Context, User]
func Curry1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A] {
	return func(t T1) ReaderResult[R, A] {
		return func(r R) (A, error) {
			return f(r, t)
		}
	}
}

// Curry2 converts a function with two parameters into a fully curried function.
// Each parameter is applied one at a time.
//
// Example:
//
//	queryDB := func(ctx context.Context, table string, id int) (Record, error) { ... }
//	curried := readerresult.Curry2(queryDB)
//	// curried("users")(42) returns ReaderResult[context.Context, Record]
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) ReaderResult[R, A] {
	return func(t1 T1) func(T2) ReaderResult[R, A] {
		return func(t2 T2) ReaderResult[R, A] {
			return func(r R) (A, error) {
				return f(r, t1, t2)
			}
		}
	}
}

// Curry3 converts a function with three parameters into a fully curried function.
//
// Example:
//
//	updateRecord := func(ctx context.Context, table string, id int, data string) (Result, error) { ... }
//	curried := readerresult.Curry3(updateRecord)
//	// curried("users")(42)("data") returns ReaderResult[context.Context, Result]
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) ReaderResult[R, A] {
	return func(t1 T1) func(T2) func(T3) ReaderResult[R, A] {
		return func(t2 T2) func(T3) ReaderResult[R, A] {
			return func(t3 T3) ReaderResult[R, A] {
				return func(r R) (A, error) {
					return f(r, t1, t2, t3)
				}
			}
		}
	}
}

// Uncurry1 converts a ReaderResult-returning function back into an idiomatic Go function.
// This is useful for adapting functional code to work with traditional Go APIs.
//
// Example:
//
//	rrf := func(id int) readerresult.ReaderResult[context.Context, User] { ... }
//	gofunc := readerresult.Uncurry1(rrf)
//	// gofunc(ctx, 42) returns (User, error)
func Uncurry1[R, T1, A any](f func(T1) ReaderResult[R, A]) func(R, T1) (A, error) {
	return func(r R, t T1) (A, error) {
		return f(t)(r)
	}
}

// Uncurry2 converts a curried two-parameter ReaderResult function into an idiomatic Go function.
//
// Example:
//
//	rrf := func(table string) func(int) readerresult.ReaderResult[context.Context, Record] { ... }
//	gofunc := readerresult.Uncurry2(rrf)
//	// gofunc(ctx, "users", 42) returns (Record, error)
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderResult[R, A]) func(R, T1, T2) (A, error) {
	return func(r R, t1 T1, t2 T2) (A, error) {
		return f(t1)(t2)(r)
	}
}

// Uncurry3 converts a curried three-parameter ReaderResult function into an idiomatic Go function.
//
// Example:
//
//	rrf := func(table string) func(int) func(string) readerresult.ReaderResult[context.Context, Result] { ... }
//	gofunc := readerresult.Uncurry3(rrf)
//	// gofunc(ctx, "users", 42, "data") returns (Result, error)
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderResult[R, A]) func(R, T1, T2, T3) (A, error) {
	return func(r R, t1 T1, t2 T2, t3 T3) (A, error) {
		return f(t1)(t2)(t3)(r)
	}
}
