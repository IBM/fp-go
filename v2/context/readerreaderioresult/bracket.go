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

package readerreaderioresult

import (
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

// Bracket ensures that a resource is properly cleaned up regardless of whether the operation
// succeeds or fails. It follows the acquire-use-release pattern with access to both outer (R)
// and inner (C) reader contexts.
//
// The release action is always called after the use action completes, whether it succeeds or fails.
// This makes it ideal for managing resources like file handles, database connections, or locks.
//
// Parameters:
//   - acquire: Acquires the resource, returning a ReaderReaderIOEither[R, C, E, A]
//   - use: Uses the acquired resource to perform an operation, returning ReaderReaderIOEither[R, C, E, B]
//   - release: Releases the resource, receiving both the resource and the result of use
//
// Returns:
//   - A ReaderReaderIOEither[R, C, E, B] that safely manages the resource lifecycle
//
// The release function receives:
//   - The acquired resource (A)
//   - The result of the use function (Either[E, B])
//
// Example:
//
//	type OuterConfig struct {
//	    ConnectionPool string
//	}
//	type InnerConfig struct {
//	    Timeout time.Duration
//	}
//
//	// Acquire a database connection
//	acquire := func(outer OuterConfig) readerioeither.ReaderIOEither[InnerConfig, error, *sql.DB] {
//	    return func(inner InnerConfig) ioeither.IOEither[error, *sql.DB] {
//	        return ioeither.TryCatch(
//	            func() (*sql.DB, error) {
//	                return sql.Open("postgres", outer.ConnectionPool)
//	            },
//	            func(err error) error { return err },
//	        )
//	    }
//	}
//
//	// Use the connection
//	use := func(db *sql.DB) readerreaderioeither.ReaderReaderIOEither[OuterConfig, InnerConfig, error, []User] {
//	    return func(outer OuterConfig) readerioeither.ReaderIOEither[InnerConfig, error, []User] {
//	        return func(inner InnerConfig) ioeither.IOEither[error, []User] {
//	            return queryUsers(db, inner.Timeout)
//	        }
//	    }
//	}
//
//	// Release the connection
//	release := func(db *sql.DB, result either.Either[error, []User]) readerreaderioeither.ReaderReaderIOEither[OuterConfig, InnerConfig, error, any] {
//	    return func(outer OuterConfig) readerioeither.ReaderIOEither[InnerConfig, error, any] {
//	        return func(inner InnerConfig) ioeither.IOEither[error, any] {
//	            return ioeither.TryCatch(
//	                func() (any, error) {
//	                    return nil, db.Close()
//	                },
//	                func(err error) error { return err },
//	            )
//	        }
//	    }
//	}
//
//	result := readerreaderioeither.Bracket(acquire, use, release)
//
//go:inline
func Bracket[
	R, A, B, ANY any](
	acquire ReaderReaderIOResult[R, A],
	use Kleisli[R, A, B],
	release func(A, Result[B]) ReaderReaderIOResult[R, ANY],
) ReaderReaderIOResult[R, B] {
	return RRIOE.Bracket(acquire, use, release)
}
