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

package readeroption

import (
	G "github.com/IBM/fp-go/v2/readeroption/generic"
)

// Curry functions convert Go functions that take a context as the first parameter
// and return (value, bool) into curried ReaderOption functions.
//
// This follows the Go convention of passing context as the first parameter
// (see https://pkg.go.dev/context), while providing a functional programming interface.
//
// The bool return value indicates success (true) or failure (false), which maps to
// Some or None in the Option monad.

// Curry0 converts a function that takes only a context and returns (A, bool)
// into a ReaderOption[R, A].
//
// Example:
//
//	getConfig := func(ctx context.Context) (Config, bool) {
//	    cfg, ok := ctx.Value("config").(Config)
//	    return cfg, ok
//	}
//	ro := readeroption.Curry0(getConfig)
//	result := ro(ctx) // Returns option.Some(config) or option.None()
func Curry0[R, A any](f func(R) (A, bool)) ReaderOption[R, A] {
	return G.Curry0[ReaderOption[R, A]](f)
}

// Curry1 converts a function that takes a context and one argument, returning (A, bool),
// into a curried function that returns a ReaderOption.
//
// Example:
//
//	findUser := func(ctx context.Context, id int) (User, bool) {
//	    // Query database using context
//	    return user, found
//	}
//	ro := readeroption.Curry1(findUser)
//	result := ro(123)(ctx) // Returns option.Some(user) or option.None()
func Curry1[R, T1, A any](f func(R, T1) (A, bool)) Kleisli[R, T1, A] {
	return G.Curry1[ReaderOption[R, A]](f)
}

// Curry2 converts a function that takes a context and two arguments, returning (A, bool),
// into a curried function that returns a ReaderOption.
//
// Example:
//
//	query := func(ctx context.Context, table string, id int) (Record, bool) {
//	    // Query database using context
//	    return record, found
//	}
//	ro := readeroption.Curry2(query)
//	result := ro("users")(123)(ctx) // Returns option.Some(record) or option.None()
func Curry2[R, T1, T2, A any](f func(R, T1, T2) (A, bool)) func(T1) func(T2) ReaderOption[R, A] {
	return G.Curry2[ReaderOption[R, A]](f)
}

// Curry3 converts a function that takes a context and three arguments, returning (A, bool),
// into a curried function that returns a ReaderOption.
//
// Example:
//
//	complexQuery := func(ctx context.Context, db string, table string, id int) (Record, bool) {
//	    // Query database using context
//	    return record, found
//	}
//	ro := readeroption.Curry3(complexQuery)
//	result := ro("mydb")("users")(123)(ctx)
func Curry3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, bool)) func(T1) func(T2) func(T3) ReaderOption[R, A] {
	return G.Curry3[ReaderOption[R, A]](f)
}

// Uncurry1 converts a curried ReaderOption function back to a Go function
// that takes a context and one argument, returning (A, bool).
//
// Example:
//
//	ro := func(id int) readeroption.ReaderOption[context.Context, User] { ... }
//	findUser := readeroption.Uncurry1(ro)
//	user, found := findUser(ctx, 123)
func Uncurry1[R, T1, A any](f func(T1) ReaderOption[R, A]) func(R, T1) (A, bool) {
	return G.Uncurry1(f)
}

// Uncurry2 converts a curried ReaderOption function back to a Go function
// that takes a context and two arguments, returning (A, bool).
func Uncurry2[R, T1, T2, A any](f func(T1) func(T2) ReaderOption[R, A]) func(R, T1, T2) (A, bool) {
	return G.Uncurry2(f)
}

// Uncurry3 converts a curried ReaderOption function back to a Go function
// that takes a context and three arguments, returning (A, bool).
func Uncurry3[R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) ReaderOption[R, A]) func(R, T1, T2, T3) (A, bool) {
	return G.Uncurry3(f)
}
