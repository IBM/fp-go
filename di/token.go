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
	"fmt"
	"strconv"
	"sync/atomic"

	DIE "github.com/IBM/fp-go/di/erasure"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	IOO "github.com/IBM/fp-go/iooption"
	O "github.com/IBM/fp-go/option"
)

type Token[T any] interface {
	DIE.Token
	ToType(any) O.Option[T]
}

type InjectionToken[T any] interface {
	Token[T]
	Option() Token[O.Option[T]]
	IOEither() Token[IOE.IOEither[error, T]]
	IOOption() Token[IOO.IOOption[T]]
}

func makeId() IO.IO[string] {
	var count int64
	return func() string {
		return strconv.FormatInt(atomic.AddInt64(&count, 1), 16)
	}
}

var genId = makeId()

type token[T any] struct {
	name string
	id   string
	typ  DIE.TokenType
}

func (t *token[T]) Id() string {
	return t.id
}

func (t *token[T]) Type() DIE.TokenType {
	return t.typ
}

func (t *token[T]) ToType(value any) O.Option[T] {
	return O.ToType[T](value)
}

func (t *token[T]) String() string {
	return t.name
}

func makeToken[T any](name string, id string, typ DIE.TokenType) Token[T] {
	return &token[T]{name, id, typ}
}

type injectionToken[T any] struct {
	token[T]
	option   Token[O.Option[T]]
	ioeither Token[IOE.IOEither[error, T]]
	iooption Token[IOO.IOOption[T]]
}

func (i *injectionToken[T]) Option() Token[O.Option[T]] {
	return i.option
}

func (i *injectionToken[T]) IOEither() Token[IOE.IOEither[error, T]] {
	return i.ioeither
}

func (i *injectionToken[T]) IOOption() Token[IOO.IOOption[T]] {
	return i.iooption
}

// MakeToken create a unique injection token for a specific type
func MakeToken[T any](name string) InjectionToken[T] {
	id := genId()
	return &injectionToken[T]{
		token[T]{name, id, DIE.Mandatory},
		makeToken[O.Option[T]](fmt.Sprintf("Option[%s]", name), id, DIE.Option),
		makeToken[IOE.IOEither[error, T]](fmt.Sprintf("IOEither[%s]", name), id, DIE.IOEither),
		makeToken[IOO.IOOption[T]](fmt.Sprintf("IOOption[%s]", name), id, DIE.IOOption),
	}
}
