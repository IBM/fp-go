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

	"github.com/IBM/fp-go/v2/result"
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
	assert.Equal(t, result.Of(T.MakeTuple1("s1")), res1(defaultContext))

	res2 := SequenceT1(e1)
	assert.Equal(t, result.Left[T.Tuple1[string]](errFoo), res2(defaultContext))
}

func TestSequenceT2(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := Left[MyContext, string](errFoo)
	t2 := Of[MyContext](2)
	e2 := Left[MyContext, int](errFoo)

	res1 := SequenceT2(t1, t2)
	assert.Equal(t, result.Of(T.MakeTuple2("s1", 2)), res1(defaultContext))

	res2 := SequenceT2(e1, t2)
	assert.Equal(t, result.Left[T.Tuple2[string, int]](errFoo), res2(defaultContext))

	res3 := SequenceT2(t1, e2)
	assert.Equal(t, result.Left[T.Tuple2[string, int]](errFoo), res3(defaultContext))
}

func TestSequenceT3(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := Left[MyContext, string](errFoo)
	t2 := Of[MyContext](2)
	e2 := Left[MyContext, int](errFoo)
	t3 := Of[MyContext](true)
	e3 := Left[MyContext, bool](errFoo)

	res1 := SequenceT3(t1, t2, t3)
	assert.Equal(t, result.Of(T.MakeTuple3("s1", 2, true)), res1(defaultContext))

	res2 := SequenceT3(e1, t2, t3)
	assert.Equal(t, result.Left[T.Tuple3[string, int, bool]](errFoo), res2(defaultContext))

	res3 := SequenceT3(t1, e2, t3)
	assert.Equal(t, result.Left[T.Tuple3[string, int, bool]](errFoo), res3(defaultContext))

	res4 := SequenceT3(t1, t2, e3)
	assert.Equal(t, result.Left[T.Tuple3[string, int, bool]](errFoo), res4(defaultContext))
}

func TestSequenceT4(t *testing.T) {

	t1 := Of[MyContext]("s1")
	t2 := Of[MyContext](2)
	t3 := Of[MyContext](true)
	t4 := Of[MyContext](1.0)

	res := SequenceT4(t1, t2, t3, t4)

	assert.Equal(t, result.Of(T.MakeTuple4("s1", 2, true, 1.0)), res(defaultContext))
}
