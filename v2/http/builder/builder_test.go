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

package builder

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	FD "github.com/IBM/fp-go/v2/http/form"
	H "github.com/IBM/fp-go/v2/http/headers"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {

	name := H.ContentType
	withContentType := WithHeader(name)
	withoutContentType := WithoutHeader(name)

	b1 := F.Pipe1(
		Default,
		withContentType(C.JSON),
	)

	b2 := F.Pipe1(
		b1,
		withContentType(C.TextPlain),
	)

	b3 := F.Pipe1(
		b2,
		withoutContentType,
	)

	assert.Equal(t, O.None[string](), Default.GetHeader(name))
	assert.Equal(t, O.Of(C.JSON), b1.GetHeader(name))
	assert.Equal(t, O.Of(C.TextPlain), b2.GetHeader(name))
	assert.Equal(t, O.None[string](), b3.GetHeader(name))
}

func TestWithFormData(t *testing.T) {
	data := F.Pipe1(
		FD.Default,
		FD.WithValue("a")("b"),
	)

	res := F.Pipe1(
		Default,
		WithFormData(data),
	)

	assert.Equal(t, C.FormEncoded, Headers.Get(res).Get(H.ContentType))
}

func TestHash(t *testing.T) {

	b1 := F.Pipe4(
		Default,
		WithContentType(C.JSON),
		WithHeader(H.Accept)(C.JSON),
		WithURL("http://www.example.com"),
		WithJSON(map[string]string{"a": "b"}),
	)

	b2 := F.Pipe4(
		Default,
		WithURL("http://www.example.com"),
		WithHeader(H.Accept)(C.JSON),
		WithContentType(C.JSON),
		WithJSON(map[string]string{"a": "b"}),
	)

	assert.Equal(t, MakeHash(b1), MakeHash(b2))
	assert.NotEqual(t, MakeHash(Default), MakeHash(b2))

	fmt.Println(MakeHash(b1))
}
