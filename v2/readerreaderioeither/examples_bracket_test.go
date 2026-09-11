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

package readerreaderioeither

import (
	"errors"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

// exampleConfig is the outer environment of the examples, e.g. application wide settings
type exampleConfig struct {
	pool string
}

// exampleRequest is the inner environment of the examples, e.g. per request data
type exampleRequest struct {
	user string
}

// exampleConn is the resource managed by the examples
type exampleConn struct {
	pool string
}

// openConn acquires a connection to the pool named by the outer environment
func openConn(cfg exampleConfig) ReaderIOEither[exampleRequest, error, *exampleConn] {
	return func(req exampleRequest) IOEither[error, *exampleConn] {
		return func() Either[error, *exampleConn] {
			fmt.Printf("open %s for %s\n", cfg.pool, req.user)
			return E.Right[error](&exampleConn{pool: cfg.pool})
		}
	}
}

// closeConn releases a connection, independent of the outcome of its use
func closeConn[B any](conn *exampleConn, _ Either[error, B]) ReaderReaderIOEither[exampleConfig, exampleRequest, error, F.Void] {
	return FromIO[exampleConfig, exampleRequest, error](func() F.Void {
		fmt.Printf("close %s\n", conn.pool)
		return F.VOID
	})
}

// ExampleBracket acquires a connection using the outer environment, uses it with the
// inner environment and closes it again.
func ExampleBracket() {
	query := func(conn *exampleConn) ReaderReaderIOEither[exampleConfig, exampleRequest, error, string] {
		return func(_ exampleConfig) ReaderIOEither[exampleRequest, error, string] {
			return func(req exampleRequest) IOEither[error, string] {
				return IOE.Of[error](fmt.Sprintf("%s queried %s", req.user, conn.pool))
			}
		}
	}

	program := Bracket(openConn, query, closeConn[string])

	// nothing has happened so far, the program only runs once both environments
	// are provided and the resulting IO is executed
	fmt.Println(program(exampleConfig{pool: "users-db"})(exampleRequest{user: "alice"})())

	// Output:
	// open users-db for alice
	// close users-db
	// Right[string](alice queried users-db)
}

// ExampleBracket_useFails shows that the resource is released even if using it fails,
// and that the error of use is returned.
func ExampleBracket_useFails() {
	query := func(_ *exampleConn) ReaderReaderIOEither[exampleConfig, exampleRequest, error, string] {
		return Left[exampleConfig, exampleRequest, string](errors.New("query failed"))
	}

	program := Bracket(openConn, query, closeConn[string])

	fmt.Println(program(exampleConfig{pool: "users-db"})(exampleRequest{user: "bob"})())

	// Output:
	// open users-db for bob
	// close users-db
	// Left[*errors.errorString](query failed)
}

// ExampleBracket_acquireFails shows that neither use nor release run if the resource
// cannot be acquired.
func ExampleBracket_acquireFails() {
	acquire := Left[exampleConfig, exampleRequest, *exampleConn](errors.New("pool exhausted"))

	query := func(conn *exampleConn) ReaderReaderIOEither[exampleConfig, exampleRequest, error, string] {
		fmt.Println("query")
		return Of[exampleConfig, exampleRequest, error](conn.pool)
	}

	program := Bracket(acquire, query, closeConn[string])

	fmt.Println(program(exampleConfig{pool: "users-db"})(exampleRequest{user: "carol"})())

	// Output:
	// Left[*errors.errorString](pool exhausted)
}

// ExampleBracket_transaction uses the outcome passed to release to either commit or
// roll back a transaction. The Bracket is wrapped in a Kleisli arrow so that the same
// resource handling applies to every input.
func ExampleBracket_transaction() {
	type tx struct{ id int }

	begin := FromIO[exampleConfig, exampleRequest, error](func() *tx {
		fmt.Println("begin")
		return &tx{id: 1}
	})

	commitOrRollback := func(t *tx, outcome Either[error, int]) ReaderReaderIOEither[exampleConfig, exampleRequest, error, F.Void] {
		return FromIO[exampleConfig, exampleRequest, error](func() F.Void {
			if E.IsRight(outcome) {
				fmt.Printf("commit %d\n", t.id)
			} else {
				fmt.Printf("rollback %d\n", t.id)
			}
			return F.VOID
		})
	}

	transfer := func(amount int) ReaderReaderIOEither[exampleConfig, exampleRequest, error, int] {
		return Bracket(
			begin,
			func(_ *tx) ReaderReaderIOEither[exampleConfig, exampleRequest, error, int] {
				if amount <= 0 {
					return Left[exampleConfig, exampleRequest, int](fmt.Errorf("invalid amount %d", amount))
				}
				return Of[exampleConfig, exampleRequest, error](amount)
			},
			commitOrRollback,
		)
	}

	cfg := exampleConfig{pool: "ledger"}
	req := exampleRequest{user: "dave"}

	fmt.Println(transfer(100)(cfg)(req)())
	fmt.Println(transfer(-5)(cfg)(req)())

	// Output:
	// begin
	// commit 1
	// Right[int](100)
	// begin
	// rollback 1
	// Left[*errors.errorString](invalid amount -5)
}
