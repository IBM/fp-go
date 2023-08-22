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

package types

import (
	"fmt"

	AR "github.com/IBM/fp-go/array/generic"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/option"
)

type (
	ValidationError struct {
		Context Context
		Message string
	}

	ContextEntry struct {
		Key string
	}

	Errors []*ValidationError

	Context []*ContextEntry

	Encoder[A, O any] interface {
		Encode(A) O
	}

	Decoder[I, A any] interface {
		Validate(I, Context) E.Either[Errors, A]
		Decode(I) E.Either[Errors, A]
	}

	Type[A, O, I any] struct {
		validate func(I, Context) E.Either[Errors, A]
		encode   func(A) O
		is       func(any) option.Option[A]
	}
)

func (t *Type[A, O, I]) Validate(i I, c Context) E.Either[Errors, A] {
	return t.validate(i, c)
}

func (t *Type[A, O, I]) Decode(i I) E.Either[Errors, A] {
	return t.validate(i, AR.Of[Context](&ContextEntry{}))
}

func (t *Type[A, O, I]) Encode(a A) O {
	return t.encode(a)
}

func (t *Type[A, O, I]) Is(a any) option.Option[A] {
	return t.is(a)
}

func (t *Type[A, O, I]) AsEncoder() Encoder[A, O] {
	return t
}

func (t *Type[A, O, I]) AsDecoder() Decoder[I, A] {
	return t
}

func (val *ValidationError) Error() string {
	return val.Message
}

func Pipe[O, I, A, B any](ab Type[B, A, A]) func(a Type[A, O, I]) Type[B, O, I] {
	return func(a Type[A, O, I]) Type[B, O, I] {
		return Type[B, O, I]{
			is: ab.Is,
			validate: func(i I, c Context) E.Either[Errors, B] {
				return F.Pipe1(
					a.Validate(i, c),
					E.Chain(F.Bind2nd(ab.Validate, c)),
				)
			},
			encode: F.Flow2(
				ab.Encode,
				a.Encode,
			),
		}
	}
}

func Success[A any](value A) E.Either[Errors, A] {
	return E.Of[Errors](value)
}

func Failures[A any](err Errors) E.Either[Errors, A] {
	return E.Left[A](err)
}

func Failure[A any](c Context, message string) E.Either[Errors, A] {
	return Failures[A](AR.Of[Errors](&ValidationError{Context: c, Message: message}))
}

func makeCanonicalType[A any]() Type[A, A, any] {
	is := option.ToType[A]

	return Type[A, A, any]{
		is: is,
		validate: func(u any, c Context) E.Either[Errors, A] {
			return F.Pipe2(
				u,
				is,
				option.Fold(
					func() E.Either[Errors, A] {
						return Failure[A](c, fmt.Sprintf("source is of type %T", u))
					},
					Success[A],
				),
			)
		},
		encode: F.Identity[A],
	}
}

var String = makeCanonicalType[string]()
var Int = makeCanonicalType[int]()
var Bool = makeCanonicalType[bool]()
