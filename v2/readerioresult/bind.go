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

package readerioresult

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/result"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//	result := readerioeither.Do[Env, error](State{})
//
//go:inline
func Do[R, S any](
	empty S,
) ReaderIOResult[R, S] {
	return RIOE.Do[R, error](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access the shared environment.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOResult[Env, error, User] {
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	                return env.UserRepo.FindUser()
//	            })
//	        },
//	    ),
//	    readerioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOResult[Env, error, []Post] {
//	            // This can access s.User from the previous step
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, []Post] {
//	                return env.PostRepo.FindPostsByUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
//
//go:inline
func Bind[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return RIOE.Bind(setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, S1, S2] {
	return RIOE.Let[R, error](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, S1, S2] {
	return RIOE.LetTo[R, error](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[R, S1, T any](
	setter func(T) S1,
) Operator[R, T, S1] {
	return RIOE.BindTo[R, error](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUser := readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	    return env.UserRepo.FindUser()
//	})
//	getPosts := readerioeither.Asks(func(env Env) ioeither.IOEither[error, []Post] {
//	    return env.PostRepo.FindPosts()
//	})
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.ApS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    readerioeither.ApS(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        getPosts,
//	    ),
//	)
//
//go:inline
func ApS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIOResult[R, T],
) Operator[R, S1, S2] {
	return RIOE.ApS(setter, fa)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// This eliminates the need to manually write setter functions.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	getUser := readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	    return env.UserRepo.FindUser()
//	})
//	result := F.Pipe2(
//	    readerioeither.Of[Env, error](State{}),
//	    readerioeither.ApSL(userLens, getUser),
//	)
//
//go:inline
func ApSL[R, S, T any](
	lens L.Lens[S, T],
	fa ReaderIOResult[R, T],
) Operator[R, S, S] {
	return RIOE.ApSL(lens, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a ReaderIOResult computation that produces an updated value.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.BindL(userLens, func(user User) readerioeither.ReaderIOResult[Env, error, User] {
//	        return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	            return env.UserRepo.FindUser()
//	        })
//	    }),
//	)
//
//go:inline
func BindL[R, S, T any](
	lens L.Lens[S, T],
	f Kleisli[R, T, T],
) Operator[R, S, S] {
	return RIOE.BindL(lens, f)
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a ReaderIOResult).
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do[any, error](State{User: User{Name: "Alice"}}),
//	    readerioeither.LetL(userLens, func(user User) User {
//	        user.Name = "Bob"
//	        return user
//	    }),
//	)
//
//go:inline
func LetL[R, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[R, S, S] {
	return RIOE.LetL[R, error](lens, f)
}

// LetToL is a variant of LetTo that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The value b is set directly to the focused field.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	newUser := User{Name: "Bob", ID: 123}
//	result := F.Pipe2(
//	    readerioeither.Do[any, error](State{}),
//	    readerioeither.LetToL(userLens, newUser),
//	)
//
//go:inline
func LetToL[R, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[R, S, S] {
	return RIOE.LetToL[R, error](lens, b)
}

// BindIOEitherK is a variant of Bind that works with IOEither computations.
// It lifts an IOEither Kleisli arrow into the ReaderIOResult context, allowing you to
// compose IOEither operations within a do-notation chain.
//
// This is useful when you have an existing IOEither computation that doesn't need
// access to the Reader environment, and you want to integrate it into a ReaderIOResult pipeline.
//
// Parameters:
//   - setter: A function that takes the result T and returns a function to update the state from S1 to S2
//   - f: An IOEither Kleisli arrow that takes S1 and returns IOEither[error, T]
//
// Returns:
//   - An Operator that can be used in a do-notation chain
//
// Example:
//
//	type State struct {
//	    UserID int
//	    Data   []byte
//	}
//
//	// An IOEither operation that reads a file
//	readFile := func(s State) ioeither.IOEither[error, []byte] {
//	    return ioeither.TryCatch(func() ([]byte, error) {
//	        return os.ReadFile(fmt.Sprintf("user_%d.json", s.UserID))
//	    })
//	}
//
//	result := F.Pipe2(
//	    readerioresult.Do[Env, error](State{UserID: 123}),
//	    readerioresult.BindIOEitherK(
//	        func(data []byte) func(State) State {
//	            return func(s State) State { s.Data = data; return s }
//	        },
//	        readFile,
//	    ),
//	)
func BindIOEitherK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOEither[R, T]))
}

func BindIOResultK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOResult[R, T]))
}

// BindIOK is a variant of Bind that works with IO computations.
// It lifts an IO Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IO Kleisli arrow (S1 -> IO[T])
//
// Example:
//
//	getCurrentTime := func(s State) io.IO[time.Time] {
//	    return func() time.Time { return time.Now() }
//	}
func BindIOK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIO[R, T]))
}

// BindReaderK is a variant of Bind that works with Reader computations.
// It lifts a Reader Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A Reader Kleisli arrow (S1 -> Reader[R, T])
//
// Example:
//
//	getConfig := func(s State) reader.Reader[Env, string] {
//	    return func(env Env) string { return env.ConfigValue }
//	}
func BindReaderK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReader[R, T]))
}

// BindReaderIOK is a variant of Bind that works with ReaderIO computations.
// It lifts a ReaderIO Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A ReaderIO Kleisli arrow (S1 -> ReaderIO[R, T])
//
// Example:
//
//	logState := func(s State) readerio.ReaderIO[Env, string] {
//	    return func(env Env) io.IO[string] {
//	        return func() string {
//	            env.Logger.Println(s)
//	            return "logged"
//	        }
//	    }
//	}
func BindReaderIOK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReaderIO[R, T]))
}

// BindEitherK is a variant of Bind that works with Either (Result) computations.
// It lifts an Either Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An Either Kleisli arrow (S1 -> Either[error, T])
//
// Example:
//
//	parseValue := func(s State) result.Result[int] {
//	    return result.TryCatch(func() (int, error) {
//	        return strconv.Atoi(s.StringValue)
//	    })
//	}
func BindEitherK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromEither[R, T]))
}

func BindResultK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, F.Flow2(f, FromResult[R, T]))
}

// BindIOEitherKL is a lens-based variant of BindIOEitherK.
// It combines a lens with an IOEither Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IOEither Kleisli arrow (T -> IOEither[error, T])
//
// Example:
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//	updateUser := func(u User) ioeither.IOEither[error, User] {
//	    return ioeither.TryCatch(func() (User, error) {
//	        return saveUser(u)
//	    })
//	}
//	result := F.Pipe2(
//	    readerioresult.Do[Env](State{}),
//	    readerioresult.BindIOEitherKL(userLens, updateUser),
//	)
func BindIOEitherKL[R, S, T any](
	lens L.Lens[S, T],
	f ioresult.Kleisli[T, T],
) Operator[R, S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[R, T]))
}

// BindIOResultKL is a lens-based variant of BindIOResultK.
// It combines a lens with an IOResult Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IOResult Kleisli arrow (T -> IOResult[T])
func BindIOResultKL[R, S, T any](
	lens L.Lens[S, T],
	f ioresult.Kleisli[T, T],
) Operator[R, S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[R, T]))
}

// BindIOKL is a lens-based variant of BindIOK.
// It combines a lens with an IO Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IO Kleisli arrow (T -> IO[T])
//
// Example:
//
//	timestampLens := lens.MakeLens(
//	    func(s State) time.Time { return s.Timestamp },
//	    func(s State, t time.Time) State { s.Timestamp = t; return s },
//	)
//	updateTimestamp := func(t time.Time) io.IO[time.Time] {
//	    return func() time.Time { return time.Now() }
//	}
func BindIOKL[R, S, T any](
	lens L.Lens[S, T],
	f io.Kleisli[T, T],
) Operator[R, S, S] {
	return BindL(lens, F.Flow2(f, FromIO[R, T]))
}

// BindReaderKL is a lens-based variant of BindReaderK.
// It combines a lens with a Reader Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A Reader Kleisli arrow (T -> Reader[R, T])
//
// Example:
//
//	configLens := lens.MakeLens(
//	    func(s State) string { return s.Config },
//	    func(s State, c string) State { s.Config = c; return s },
//	)
//	getConfigFromEnv := func(c string) reader.Reader[Env, string] {
//	    return func(env Env) string { return env.ConfigValue }
//	}
func BindReaderKL[R, S, T any](
	lens L.Lens[S, T],
	f reader.Kleisli[R, T, T],
) Operator[R, S, S] {
	return BindL(lens, F.Flow2(f, FromReader[R, T]))
}

// BindReaderIOKL is a lens-based variant of BindReaderIOK.
// It combines a lens with a ReaderIO Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A ReaderIO Kleisli arrow (T -> ReaderIO[R, T])
//
// Example:
//
//	logLens := lens.MakeLens(
//	    func(s State) string { return s.LogMessage },
//	    func(s State, l string) State { s.LogMessage = l; return s },
//	)
//	logMessage := func(msg string) readerio.ReaderIO[Env, string] {
//	    return func(env Env) io.IO[string] {
//	        return func() string {
//	            env.Logger.Println(msg)
//	            return "logged: " + msg
//	        }
//	    }
//	}
func BindReaderIOKL[R, S, T any](
	lens L.Lens[S, T],
	f readerio.Kleisli[R, T, T],
) Operator[R, S, S] {
	return BindL(lens, F.Flow2(f, FromReaderIO[R, T]))
}

// ApIOEitherS is an applicative variant that works with IOEither values.
// Unlike BindIOEitherK, this uses applicative composition (ApS) instead of monadic
// composition (Bind), allowing independent computations to be combined.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IOEither value (not a Kleisli arrow)
//
// Example:
//
//	readConfig := ioeither.TryCatch(func() (Config, error) {
//	    return loadConfig()
//	})
func ApIOEitherS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[R, S1, S2] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[R, S1], ioresult.Operator[S1, S2]], ioeither.ApS(setter, fa))
}

// ApIOResultS is an applicative variant that works with IOResult values.
// This is an alias for ApIOEitherS for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IOResult value
func ApIOResultS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[R, S1, S2] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[R, S1], ioresult.Operator[S1, S2]], ioeither.ApS(setter, fa))
}

// ApIOS is an applicative variant that works with IO values.
// It lifts an IO value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IO value
//
// Example:
//
//	getCurrentTime := func() time.Time { return time.Now() }
func ApIOS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[R, S1, S2] {
	return ApS(setter, FromIO[R](fa))
}

// ApReaderS is an applicative variant that works with Reader values.
// It lifts a Reader value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A Reader value
//
// Example:
//
//	getEnvConfig := func(env Env) string { return env.ConfigValue }
func ApReaderS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, S1, S2] {
	return ApS(setter, FromReader(fa))
}

// ApReaderIOS is an applicative variant that works with ReaderIO values.
// It lifts a ReaderIO value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A ReaderIO value
func ApReaderIOS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[R, T],
) Operator[R, S1, S2] {
	return ApS(setter, FromReaderIO(fa))
}

// ApEitherS is an applicative variant that works with Either (Result) values.
// It lifts an Either value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An Either value
//
// Example:
//
//	parseResult := result.TryCatch(func() (int, error) {
//	    return strconv.Atoi("123")
//	})
func ApEitherS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[R, S1, S2] {
	return ApS(setter, FromEither[R](fa))
}

// ApResultS is an applicative variant that works with Result values.
// This is an alias for ApEitherS for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A Result value
func ApResultS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[R, S1, S2] {
	return ApS(setter, FromResult[R](fa))
}

// ApIOEitherSL is a lens-based variant of ApIOEitherS.
// It combines a lens with an IOEither value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IOEither value
//
// Example:
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//	loadUser := ioeither.TryCatch(func() (User, error) {
//	    return fetchUser()
//	})
func ApIOEitherSL[R, S, T any](
	lens L.Lens[S, T],
	fa IOResult[T],
) Operator[R, S, S] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[R, S], ioresult.Operator[S, S]], ioresult.ApSL(lens, fa))
}

// ApIOResultSL is a lens-based variant of ApIOResultS.
// This is an alias for ApIOEitherSL for consistency with the Result naming convention.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IOResult value
func ApIOResultSL[R, S, T any](
	lens L.Lens[S, T],
	fa IOResult[T],
) Operator[R, S, S] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[R, S], ioresult.Operator[S, S]], ioresult.ApSL(lens, fa))
}

// ApIOSL is a lens-based variant of ApIOS.
// It combines a lens with an IO value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IO value
//
// Example:
//
//	timestampLens := lens.MakeLens(
//	    func(s State) time.Time { return s.Timestamp },
//	    func(s State, t time.Time) State { s.Timestamp = t; return s },
//	)
//	getCurrentTime := func() time.Time { return time.Now() }
func ApIOSL[R, S, T any](
	lens L.Lens[S, T],
	fa IO[T],
) Operator[R, S, S] {
	return ApSL(lens, FromIO[R](fa))
}

// ApReaderSL is a lens-based variant of ApReaderS.
// It combines a lens with a Reader value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A Reader value
//
// Example:
//
//	configLens := lens.MakeLens(
//	    func(s State) string { return s.Config },
//	    func(s State, c string) State { s.Config = c; return s },
//	)
//	getConfig := func(env Env) string { return env.ConfigValue }
func ApReaderSL[R, S, T any](
	lens L.Lens[S, T],
	fa Reader[R, T],
) Operator[R, S, S] {
	return ApSL(lens, FromReader(fa))
}

// ApReaderIOSL is a lens-based variant of ApReaderIOS.
// It combines a lens with a ReaderIO value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A ReaderIO value
//
// Example:
//
//	logLens := lens.MakeLens(
//	    func(s State) string { return s.LogMessage },
//	    func(s State, l string) State { s.LogMessage = l; return s },
//	)
//	logWithEnv := func(env Env) io.IO[string] {
//	    return func() string {
//	        env.Logger.Println("Processing")
//	        return "logged"
//	    }
//	}
func ApReaderIOSL[R, S, T any](
	lens L.Lens[S, T],
	fa ReaderIO[R, T],
) Operator[R, S, S] {
	return ApSL(lens, FromReaderIO(fa))
}

// ApEitherSL is a lens-based variant of ApEitherS.
// It combines a lens with an Either value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An Either value
//
// Example:
//
//	valueLens := lens.MakeLens(
//	    func(s State) int { return s.Value },
//	    func(s State, v int) State { s.Value = v; return s },
//	)
//	parseValue := result.TryCatch(func() (int, error) {
//	    return strconv.Atoi("123")
//	})
func ApEitherSL[R, S, T any](
	lens L.Lens[S, T],
	fa Result[T],
) Operator[R, S, S] {
	return ApSL(lens, FromEither[R](fa))
}

// ApResultSL is a lens-based variant of ApResultS.
// This is an alias for ApEitherSL for consistency with the Result naming convention.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A Result value
func ApResultSL[R, S, T any](
	lens L.Lens[S, T],
	fa Result[T],
) Operator[R, S, S] {
	return ApSL(lens, FromResult[R](fa))
}
