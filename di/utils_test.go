// Copyright (c) 2023 IBM Corp.
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
package di

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestToType(t *testing.T) {
	// good cases
	assert.Equal(t, E.Of[error](10), toType[int]()(any(10)))
	assert.Equal(t, E.Of[error]("Carsten"), toType[string]()(any("Carsten")))
	assert.Equal(t, E.Of[error](O.Of("Carsten")), toType[O.Option[string]]()(any(O.Of("Carsten"))))
	assert.Equal(t, E.Of[error](O.Of(any("Carsten"))), toType[O.Option[any]]()(any(O.Of(any("Carsten")))))
	// failure
	assert.False(t, E.IsRight(toType[int]()(any("Carsten"))))
	assert.False(t, E.IsRight(toType[O.Option[string]]()(O.Of(any("Carsten")))))
}

func TestToOptionType(t *testing.T) {
	// good cases
	assert.Equal(t, E.Of[error](O.Of(10)), toOptionType[int]()(any(O.Of(any(10)))))
	assert.Equal(t, E.Of[error](O.Of("Carsten")), toOptionType[string]()(any(O.Of(any("Carsten")))))
	// bad cases
	assert.False(t, E.IsRight(toOptionType[int]()(any(10))))
	assert.False(t, E.IsRight(toOptionType[int]()(any(O.Of(10)))))
}

func invokeIOEither[T any](e E.Either[error, IOE.IOEither[error, T]]) E.Either[error, T] {
	return F.Pipe1(
		e,
		E.Chain(func(ioe IOE.IOEither[error, T]) E.Either[error, T] {
			return ioe()
		}),
	)
}

func TestToIOEitherType(t *testing.T) {
	// good cases
	assert.Equal(t, E.Of[error](10), invokeIOEither(toIOEitherType[int]()(any(IOE.Of[error](any(10))))))
	assert.Equal(t, E.Of[error]("Carsten"), invokeIOEither(toIOEitherType[string]()(any(IOE.Of[error](any("Carsten"))))))
	// bad cases
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherType[string]()(any(IOE.Of[error](any(10)))))))
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherType[string]()(any(IOE.Of[error]("Carsten"))))))
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherType[string]()(any("Carsten")))))
}
