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

package readerresult

import (
	"fmt"
	"testing"

	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

var (
	errFoo = fmt.Errorf("error")
)

func TestSequenceT1(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := Left[MyContext, string](errFoo)

	res1 := SequenceT1(t1)
	v1, err1 := res1(defaultContext)
	assert.NoError(t, err1)
	assert.Equal(t, T.MakeTuple1("s1"), v1)

	res2 := SequenceT1(e1)
	_, err2 := res2(defaultContext)
	assert.Equal(t, errFoo, err2)
}

func TestSequenceT2(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := Left[MyContext, string](errFoo)
	t2 := Of[MyContext](2)
	e2 := Left[MyContext, int](errFoo)

	res1 := SequenceT2(t1, t2)
	v1, err1 := res1(defaultContext)
	assert.NoError(t, err1)
	assert.Equal(t, T.MakeTuple2("s1", 2), v1)

	res2 := SequenceT2(e1, t2)
	_, err2 := res2(defaultContext)
	assert.Equal(t, errFoo, err2)

	res3 := SequenceT2(t1, e2)
	_, err3 := res3(defaultContext)
	assert.Equal(t, errFoo, err3)
}

func TestSequenceT3(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := Left[MyContext, string](errFoo)
	t2 := Of[MyContext](2)
	e2 := Left[MyContext, int](errFoo)
	t3 := Of[MyContext](true)
	e3 := Left[MyContext, bool](errFoo)

	res1 := SequenceT3(t1, t2, t3)
	v1, err1 := res1(defaultContext)
	assert.NoError(t, err1)
	assert.Equal(t, T.MakeTuple3("s1", 2, true), v1)

	res2 := SequenceT3(e1, t2, t3)
	_, err2 := res2(defaultContext)
	assert.Equal(t, errFoo, err2)

	res3 := SequenceT3(t1, e2, t3)
	_, err3 := res3(defaultContext)
	assert.Equal(t, errFoo, err3)

	res4 := SequenceT3(t1, t2, e3)
	_, err4 := res4(defaultContext)
	assert.Equal(t, errFoo, err4)
}

func TestSequenceT4(t *testing.T) {

	t1 := Of[MyContext]("s1")
	t2 := Of[MyContext](2)
	t3 := Of[MyContext](true)
	t4 := Of[MyContext](1.0)

	res := SequenceT4(t1, t2, t3, t4)

	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, T.MakeTuple4("s1", 2, true, 1.0), v)
}
