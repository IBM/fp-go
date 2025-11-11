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

// From functions convert Go functions that take a context as the first parameter
// and return (value, bool) into ReaderOption functions with uncurried parameters.
//
// Unlike Curry functions which return fully curried functions, From functions
// return functions that take multiple parameters at once (uncurried style).
//
// This follows the Go convention of passing context as the first parameter
// (see https://pkg.go.dev/context), while providing a functional programming interface.
//
// The bool return value indicates success (true) or failure (false), which maps to
// Some or None in the Option monad.

// From0 converts a function that takes only a context and returns (A, bool)
// into a function that returns a ReaderOption[R, A].
//
// Example:
//
//	getConfig := func(ctx context.Context) (Config, bool) {
//	    cfg, ok := ctx.Value("config").(Config)
//	    return cfg, ok
//	}
//	roFunc := readeroption.From0(getConfig)
//	ro := roFunc() // Returns a ReaderOption[context.Context, Config]
//	result := ro(ctx) // Returns option.Some(config) or option.None()
func From0[R, A any](f func(R) (A, bool)) func() ReaderOption[R, A] {
	return G.From0[ReaderOption[R, A]](f)
}

// From1 converts a function that takes a context and one argument, returning (A, bool),
// into a function that takes one argument and returns a ReaderOption.
//
// This is equivalent to Curry1 but provided for consistency with the From naming convention.
//
// Example:
//
//	findUser := func(ctx context.Context, id int) (User, bool) {
//	    // Query database using context
//	    return user, found
//	}
//	roFunc := readeroption.From1(findUser)
//	ro := roFunc(123) // Returns a ReaderOption[context.Context, User]
//	result := ro(ctx) // Returns option.Some(user) or option.None()
func From1[R, T1, A any](f func(R, T1) (A, bool)) Kleisli[R, T1, A] {
	return G.From1[ReaderOption[R, A]](f)
}

// From2 converts a function that takes a context and two arguments, returning (A, bool),
// into a function that takes two arguments (uncurried) and returns a ReaderOption.
//
// Example:
//
//	query := func(ctx context.Context, table string, id int) (Record, bool) {
//	    // Query database using context
//	    return record, found
//	}
//	roFunc := readeroption.From2(query)
//	ro := roFunc("users", 123) // Returns a ReaderOption[context.Context, Record]
//	result := ro(ctx) // Returns option.Some(record) or option.None()
func From2[R, T1, T2, A any](f func(R, T1, T2) (A, bool)) func(T1, T2) ReaderOption[R, A] {
	return G.From2[ReaderOption[R, A]](f)
}

// From3 converts a function that takes a context and three arguments, returning (A, bool),
// into a function that takes three arguments (uncurried) and returns a ReaderOption.
//
// Example:
//
//	complexQuery := func(ctx context.Context, db string, table string, id int) (Record, bool) {
//	    // Query database using context
//	    return record, found
//	}
//	roFunc := readeroption.From3(complexQuery)
//	ro := roFunc("mydb", "users", 123) // Returns a ReaderOption[context.Context, Record]
//	result := ro(ctx) // Returns option.Some(record) or option.None()
func From3[R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, bool)) func(T1, T2, T3) ReaderOption[R, A] {
	return G.From3[ReaderOption[R, A]](f)
}
