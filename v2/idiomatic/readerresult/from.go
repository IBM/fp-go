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

import (
	"github.com/IBM/fp-go/v2/function"
)

// These functions convert idiomatic Go functions (with context as first parameter and (value, error) return)
// into ReaderResult computations. This follows the Go convention of putting context as the first parameter
// as advised in https://pkg.go.dev/context.

// From0 converts a function that takes only a context and returns (A, error) into a ReaderResult.
//
// Example:
//
//	getConfig := func(ctx context.Context) (Config, error) { ... }
//	rr := readerresult.From0(getConfig)()
//	// rr is a ReaderResult[context.Context, Config]
func From0[R, A any](f func(R) (A, error)) func() ReaderResult[R, A] {
	return function.Constant(f)
}

// From1 converts a function with one parameter into a ReaderResult-returning function.
// The context parameter is moved to the end (ReaderResult style).
//
// Example:
//
//	getUser := func(ctx context.Context, id int) (User, error) { ... }
//	rr := readerresult.From1(getUser)
//	// rr(42) returns ReaderResult[context.Context, User]
func From1[R, T1, A any](f func(R, T1) (A, error)) func(T1) ReaderResult[R, A] {
	return func(t1 T1) ReaderResult[R, A] {
		return func(r R) (A, error) {
			return f(r, t1)
		}
	}
}

// From2 converts a function with two parameters into a ReaderResult-returning function.
// The context parameter is moved to the end (ReaderResult style).
//
// Example:
//
//	queryDB := func(ctx context.Context, table string, id int) (Record, error) { ... }
//	rr := readerresult.From2(queryDB)
//	// rr("users", 42) returns ReaderResult[context.Context, Record]
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) ReaderResult[R, A] {
	return func(t1 T1, t2 T2) ReaderResult[R, A] {
		return func(r R) (A, error) {
			return f(r, t1, t2)
		}
	}
}

// From3 converts a function with three parameters into a ReaderResult-returning function.
// The context parameter is moved to the end (ReaderResult style).
//
// Example:
//
//	updateRecord := func(ctx context.Context, table string, id int, data string) (Result, error) { ... }
//	rr := readerresult.From3(updateRecord)
//	// rr("users", 42, "data") returns ReaderResult[context.Context, Result]
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) ReaderResult[R, A] {
	return func(t1 T1, t2 T2, t3 T3) ReaderResult[R, A] {
		return func(r R) (A, error) {
			return f(r, t1, t2, t3)
		}
	}
}
