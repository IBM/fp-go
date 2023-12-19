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

package http

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/http/content"
	"github.com/stretchr/testify/assert"
)

func NoError[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(func(err error) bool {
		return assert.NoError(t, err)
	}, F.Constant1[A](true))
}

func Error[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(F.Constant1[error](true), func(A) bool {
		return assert.Error(t, nil)
	})
}

func TestValidateJsonContentTypeString(t *testing.T) {

	res := F.Pipe1(
		validateJsonContentTypeString(C.Json),
		NoError[ParsedMediaType](t),
	)

	assert.True(t, res)
}

func TestValidateInvalidJsonContentTypeString(t *testing.T) {

	res := F.Pipe1(
		validateJsonContentTypeString("application/xml"),
		Error[ParsedMediaType](t),
	)

	assert.True(t, res)
}
