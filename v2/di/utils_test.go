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

package di

import (
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

var (
	toInt    = toType[int]()
	toString = toType[string]()
)

func TestToType(t *testing.T) {
	// good cases
	assert.Equal(t, result.Of(10), toInt(any(10)))
	assert.Equal(t, result.Of("Carsten"), toString(any("Carsten")))
	assert.Equal(t, result.Of(O.Of("Carsten")), toType[Option[string]]()(any(O.Of("Carsten"))))
	assert.Equal(t, result.Of(O.Of(any("Carsten"))), toType[Option[any]]()(any(O.Of(any("Carsten")))))
	// failure
	assert.False(t, E.IsRight(toInt(any("Carsten"))))
	assert.False(t, E.IsRight(toType[Option[string]]()(O.Of(any("Carsten")))))
}

func TestToOptionType(t *testing.T) {
	// shortcuts
	toOptInt := toOptionType(toInt)
	toOptString := toOptionType(toString)
	// good cases
	assert.Equal(t, result.Of(O.Of(10)), toOptInt(any(O.Of(any(10)))))
	assert.Equal(t, result.Of(O.Of("Carsten")), toOptString(any(O.Of(any("Carsten")))))
	// bad cases
	assert.False(t, E.IsRight(toOptInt(any(10))))
	assert.False(t, E.IsRight(toOptInt(any(O.Of(10)))))
}

func invokeIOEither[T any](e Result[IOResult[T]]) Result[T] {
	return F.Pipe1(
		e,
		E.Chain(func(ioe IOResult[T]) Result[T] {
			return ioe()
		}),
	)
}

func TestToIOEitherType(t *testing.T) {
	// shortcuts
	toIOEitherInt := toIOEitherType(toInt)
	toIOEitherString := toIOEitherType(toString)
	// good cases
	assert.Equal(t, result.Of(10), invokeIOEither(toIOEitherInt(any(ioresult.Of(any(10))))))
	assert.Equal(t, result.Of("Carsten"), invokeIOEither(toIOEitherString(any(ioresult.Of(any("Carsten"))))))
	// bad cases
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherString(any(ioresult.Of(any(10)))))))
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherString(any(ioresult.Of("Carsten"))))))
	assert.False(t, E.IsRight(invokeIOEither(toIOEitherString(any("Carsten")))))
}

func TestToArrayType(t *testing.T) {
	// shortcuts
	toArrayString := toArrayType(toString)
	// good cases
	assert.Equal(t, result.Of(A.From("a", "b")), toArrayString(any(A.From(any("a"), any("b")))))
}
