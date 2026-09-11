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
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/stretchr/testify/assert"
)

func readTestComputation(calls *int) ReaderReaderIOEither[OuterConfig, InnerConfig, error, string] {
	return func(o OuterConfig) RIOE.ReaderIOEither[InnerConfig, error, string] {
		return func(i InnerConfig) IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				*calls++
				return E.Right[error](o.database + ":" + i.apiKey)
			}
		}
	}
}

func TestReadIOEitherRight(t *testing.T) {
	calls := 0
	rio := IOE.Right[error](OuterConfig{database: "db"})

	res := ReadIOEither[string, OuterConfig, InnerConfig](rio)(readTestComputation(&calls))(InnerConfig{apiKey: "key"})

	assert.Equal(t, 0, calls, "computation must not run before the IO is executed")
	assert.Equal(t, E.Right[error]("db:key"), res())
	assert.Equal(t, E.Right[error]("db:key"), res())
	assert.Equal(t, 2, calls)
}

func TestReadIOEitherLeftEnvironment(t *testing.T) {
	calls := 0
	err := errors.New("no env")
	rio := IOE.Left[OuterConfig](err)

	res := ReadIOEither[string, OuterConfig, InnerConfig](rio)(readTestComputation(&calls))(InnerConfig{apiKey: "key"})()

	assert.Equal(t, E.Left[string](err), res)
	assert.Equal(t, 0, calls, "computation must not run when the environment fails")
}

func TestReadIOEitherLeftComputation(t *testing.T) {
	err := errors.New("failed")
	rio := IOE.Right[error](OuterConfig{database: "db"})
	rri := Left[OuterConfig, InnerConfig, string](err)

	res := ReadIOEither[string, OuterConfig, InnerConfig](rio)(rri)(InnerConfig{apiKey: "key"})()

	assert.Equal(t, E.Left[string](err), res)
}

func TestReadIO(t *testing.T) {
	calls := 0
	envCalls := 0
	rio := func() OuterConfig {
		envCalls++
		return OuterConfig{database: "db"}
	}

	res := ReadIO[InnerConfig, error, string](io.IO[OuterConfig](rio))(readTestComputation(&calls))(InnerConfig{apiKey: "key"})

	assert.Equal(t, 0, envCalls, "environment must not be read before the IO is executed")
	assert.Equal(t, E.Right[error]("db:key"), res())
	assert.Equal(t, 1, envCalls)
	assert.Equal(t, 1, calls)
}

func TestReadIOLeftComputation(t *testing.T) {
	err := errors.New("failed")
	rri := Left[OuterConfig, InnerConfig, string](err)

	res := ReadIO[InnerConfig, error, string](io.Of(OuterConfig{database: "db"}))(rri)(InnerConfig{apiKey: "key"})()

	assert.Equal(t, E.Left[string](err), res)
}
