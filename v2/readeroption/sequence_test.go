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

package readeroption

import (
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

func TestSequenceT1(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := None[MyContext, string]()

	res1 := SequenceT1(t1)
	assert.Equal(t, O.Of(T.MakeTuple1("s1")), res1(defaultContext))

	res2 := SequenceT1(e1)
	assert.Equal(t, O.None[T.Tuple1[string]](), res2(defaultContext))
}

func TestSequenceT2(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := None[MyContext, string]()
	t2 := Of[MyContext](2)
	e2 := None[MyContext, int]()

	res1 := SequenceT2(t1, t2)
	assert.Equal(t, O.Of(T.MakeTuple2("s1", 2)), res1(defaultContext))

	res2 := SequenceT2(e1, t2)
	assert.Equal(t, O.None[T.Tuple2[string, int]](), res2(defaultContext))

	res3 := SequenceT2(t1, e2)
	assert.Equal(t, O.None[T.Tuple2[string, int]](), res3(defaultContext))
}

func TestSequenceT3(t *testing.T) {

	t1 := Of[MyContext]("s1")
	e1 := None[MyContext, string]()
	t2 := Of[MyContext](2)
	e2 := None[MyContext, int]()
	t3 := Of[MyContext](true)
	e3 := None[MyContext, bool]()

	res1 := SequenceT3(t1, t2, t3)
	assert.Equal(t, O.Of(T.MakeTuple3("s1", 2, true)), res1(defaultContext))

	res2 := SequenceT3(e1, t2, t3)
	assert.Equal(t, O.None[T.Tuple3[string, int, bool]](), res2(defaultContext))

	res3 := SequenceT3(t1, e2, t3)
	assert.Equal(t, O.None[T.Tuple3[string, int, bool]](), res3(defaultContext))

	res4 := SequenceT3(t1, t2, e3)
	assert.Equal(t, O.None[T.Tuple3[string, int, bool]](), res4(defaultContext))
}

func TestSequenceT4(t *testing.T) {

	t1 := Of[MyContext]("s1")
	t2 := Of[MyContext](2)
	t3 := Of[MyContext](true)
	t4 := Of[MyContext](1.0)

	res := SequenceT4(t1, t2, t3, t4)

	assert.Equal(t, O.Of(T.MakeTuple4("s1", 2, true, 1.0)), res(defaultContext))
}
