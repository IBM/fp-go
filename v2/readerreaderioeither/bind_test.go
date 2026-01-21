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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

type OuterCtx struct {
	database string
}

type InnerCtx struct {
	apiKey string
}

func getLastName(s utils.Initial) ReaderReaderIOEither[OuterCtx, InnerCtx, error, string] {
	return Of[OuterCtx, InnerCtx, error]("Doe")
}

func getGivenName(s utils.WithLastName) ReaderReaderIOEither[OuterCtx, InnerCtx, error, string] {
	return Of[OuterCtx, InnerCtx, error]("John")
}

func TestDo(t *testing.T) {
	result := Do[OuterCtx, InnerCtx, error](utils.Empty)
	assert.Equal(t, E.Of[error](utils.Empty), result(OuterCtx{})(InnerCtx{})())
}

func TestBind(t *testing.T) {
	res := F.Pipe3(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[OuterCtx, InnerCtx, error](utils.GetFullName),
	)

	assert.Equal(t, E.Of[error]("John Doe"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindWithContext(t *testing.T) {
	outer := OuterCtx{database: "postgres"}
	inner := InnerCtx{apiKey: "secret"}

	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		Bind(utils.SetLastName, func(s utils.Initial) ReaderReaderIOEither[OuterCtx, InnerCtx, error, string] {
			return func(o OuterCtx) ReaderIOEither[InnerCtx, error, string] {
				return func(i InnerCtx) IOE.IOEither[error, string] {
					return IOE.Of[error](o.database + "-" + i.apiKey)
				}
			}
		}),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("postgres-secret"), res(outer)(inner)())
}

func TestLet(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.WithLastName{LastName: "Doe"}),
		Let[OuterCtx, InnerCtx, error](
			func(given string) func(utils.WithLastName) utils.WithGivenName {
				return func(s utils.WithLastName) utils.WithGivenName {
					return utils.WithGivenName{
						WithLastName: s,
						GivenName:    given,
					}
				}
			},
			func(s utils.WithLastName) string {
				return "John"
			},
		),
		Map[OuterCtx, InnerCtx, error](utils.GetFullName),
	)

	assert.Equal(t, E.Of[error]("John Doe"), res(OuterCtx{})(InnerCtx{})())
}

func TestLetTo(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.WithLastName{LastName: "Doe"}),
		LetTo[OuterCtx, InnerCtx, error](
			func(given string) func(utils.WithLastName) utils.WithGivenName {
				return func(s utils.WithLastName) utils.WithGivenName {
					return utils.WithGivenName{
						WithLastName: s,
						GivenName:    given,
					}
				}
			},
			"Jane",
		),
		Map[OuterCtx, InnerCtx, error](utils.GetFullName),
	)

	assert.Equal(t, E.Of[error]("Jane Doe"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindTo(t *testing.T) {
	res := F.Pipe1(
		Of[OuterCtx, InnerCtx, error]("Doe"),
		BindTo[OuterCtx, InnerCtx, error](func(lastName string) utils.WithLastName {
			return utils.WithLastName{LastName: lastName}
		}),
	)

	expected := utils.WithLastName{LastName: "Doe"}
	assert.Equal(t, E.Of[error](expected), res(OuterCtx{})(InnerCtx{})())
}

func TestApS(t *testing.T) {
	res := F.Pipe3(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApS(utils.SetLastName, Of[OuterCtx, InnerCtx, error]("Doe")),
		ApS(utils.SetGivenName, Of[OuterCtx, InnerCtx, error]("John")),
		Map[OuterCtx, InnerCtx, error](utils.GetFullName),
	)

	assert.Equal(t, E.Of[error]("John Doe"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindIOEitherK(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindIOEitherK[OuterCtx, InnerCtx](
			utils.SetLastName,
			func(s utils.Initial) IOE.IOEither[error, string] {
				return IOE.Of[error]("Smith")
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Smith"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindIOEitherKError(t *testing.T) {
	err := errors.New("io error")
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindIOEitherK[OuterCtx, InnerCtx](
			utils.SetLastName,
			func(s utils.Initial) IOE.IOEither[error, string] {
				return IOE.Left[string](err)
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Left[string](err), res(OuterCtx{})(InnerCtx{})())
}

func TestBindIOK(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindIOK[OuterCtx, InnerCtx, error](
			utils.SetLastName,
			func(s utils.Initial) io.IO[string] {
				return io.Of("Johnson")
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Johnson"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindReaderK(t *testing.T) {
	outer := OuterCtx{database: "mysql"}

	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindReaderK[InnerCtx, error](
			utils.SetLastName,
			func(s utils.Initial) R.Reader[OuterCtx, string] {
				return R.Asks(func(ctx OuterCtx) string {
					return ctx.database
				})
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("mysql"), res(outer)(InnerCtx{})())
}

func TestBindReaderIOK(t *testing.T) {
	outer := OuterCtx{database: "postgres"}

	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindReaderIOK[InnerCtx, error](
			utils.SetLastName,
			func(s utils.Initial) readerio.ReaderIO[OuterCtx, string] {
				return func(ctx OuterCtx) io.IO[string] {
					return io.Of(ctx.database + "-io")
				}
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("postgres-io"), res(outer)(InnerCtx{})())
}

func TestBindEitherK(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindEitherK[OuterCtx, InnerCtx](
			utils.SetLastName,
			func(s utils.Initial) E.Either[error, string] {
				return E.Of[error]("Brown")
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Brown"), res(OuterCtx{})(InnerCtx{})())
}

func TestBindEitherKError(t *testing.T) {
	err := errors.New("either error")
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		BindEitherK[OuterCtx, InnerCtx](
			utils.SetLastName,
			func(s utils.Initial) E.Either[error, string] {
				return E.Left[string](err)
			},
		),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Left[string](err), res(OuterCtx{})(InnerCtx{})())
}

func TestApIOEitherS(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApIOEitherS[OuterCtx, InnerCtx](utils.SetLastName, IOE.Of[error]("Williams")),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Williams"), res(OuterCtx{})(InnerCtx{})())
}

func TestApIOS(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApIOS[OuterCtx, InnerCtx, error](utils.SetLastName, io.Of("Davis")),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Davis"), res(OuterCtx{})(InnerCtx{})())
}

func TestApReaderS(t *testing.T) {
	outer := OuterCtx{database: "cassandra"}

	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApReaderS[InnerCtx, error](utils.SetLastName, R.Asks(func(ctx OuterCtx) string {
			return ctx.database
		})),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("cassandra"), res(outer)(InnerCtx{})())
}

func TestApReaderIOS(t *testing.T) {
	outer := OuterCtx{database: "neo4j"}

	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApReaderIOS[InnerCtx, error](utils.SetLastName, func(ctx OuterCtx) io.IO[string] {
			return io.Of(ctx.database + "-graph")
		}),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("neo4j-graph"), res(outer)(InnerCtx{})())
}

func TestApEitherS(t *testing.T) {
	res := F.Pipe2(
		Do[OuterCtx, InnerCtx, error](utils.Empty),
		ApEitherS[OuterCtx, InnerCtx](utils.SetLastName, E.Of[error]("Miller")),
		Map[OuterCtx, InnerCtx, error](func(s utils.WithLastName) string {
			return s.LastName
		}),
	)

	assert.Equal(t, E.Of[error]("Miller"), res(OuterCtx{})(InnerCtx{})())
}

func TestComplexBindChain(t *testing.T) {
	outer := OuterCtx{database: "postgres"}
	inner := InnerCtx{apiKey: "secret123"}

	type ComplexState struct {
		Database string
		APIKey   string
		Count    int
		Status   string
	}

	res := F.Pipe4(
		Do[OuterCtx, InnerCtx, error](ComplexState{}),
		Bind(
			func(db string) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState { s.Database = db; return s }
			},
			func(s ComplexState) ReaderReaderIOEither[OuterCtx, InnerCtx, error, string] {
				return func(o OuterCtx) ReaderIOEither[InnerCtx, error, string] {
					return func(i InnerCtx) IOE.IOEither[error, string] {
						return IOE.Of[error](o.database)
					}
				}
			},
		),
		Bind(
			func(key string) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState { s.APIKey = key; return s }
			},
			func(s ComplexState) ReaderReaderIOEither[OuterCtx, InnerCtx, error, string] {
				return func(o OuterCtx) ReaderIOEither[InnerCtx, error, string] {
					return func(i InnerCtx) IOE.IOEither[error, string] {
						return IOE.Of[error](i.apiKey)
					}
				}
			},
		),
		Let[OuterCtx, InnerCtx, error](
			func(count int) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState { s.Count = count; return s }
			},
			func(s ComplexState) int {
				return len(s.Database) + len(s.APIKey)
			},
		),
		LetTo[OuterCtx, InnerCtx, error](
			func(status string) func(ComplexState) ComplexState {
				return func(s ComplexState) ComplexState { s.Status = status; return s }
			},
			"ready",
		),
	)

	expected := ComplexState{
		Database: "postgres",
		APIKey:   "secret123",
		Count:    17, // len("postgres") + len("secret123")
		Status:   "ready",
	}
	assert.Equal(t, E.Of[error](expected), res(outer)(inner)())
}
