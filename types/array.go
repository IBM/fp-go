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
	"reflect"

	AR "github.com/IBM/fp-go/array/generic"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	T "github.com/IBM/fp-go/tuple"
)

func toUnknownArray[A any](item func(reflect.Value, Context) E.Either[Errors, A], c Context, val reflect.Value) []E.Either[Errors, A] {
	l := val.Len()
	res := make([]E.Either[Errors, A], l)
	for i := l - 1; i >= 0; i-- {
		v := val.Index(i)
		res[i] = item(v, AR.Push[Context](&ContextEntry{Key: fmt.Sprintf("[%d]", i), Value: v})(c))
	}
	return res
}

func flattenUnknownArray[GA ~[]A, A any](as []E.Either[Errors, A]) E.Either[Errors, GA] {
	return F.Pipe1(
		AR.Reduce(as, func(t T.Tuple2[GA, Errors], item E.Either[Errors, A]) T.Tuple2[GA, Errors] {
			return E.MonadFold(item, func(e Errors) T.Tuple2[GA, Errors] {
				return T.MakeTuple2(t.F1, append(t.F2, e...))
			}, func(a A) T.Tuple2[GA, Errors] {
				return T.MakeTuple2(append(t.F1, a), t.F2)
			})
		}, T.MakeTuple2(make(GA, len(as)), make(Errors, 0))),
		func(t T.Tuple2[GA, Errors]) E.Either[Errors, GA] {
			if AR.IsEmpty(t.F2) {
				return E.Of[Errors](t.F1)
			}
			return E.Left[GA](t.F2)
		},
	)
}

func toValidatedArray[GA ~[]A, A any](item func(reflect.Value, Context) E.Either[Errors, A], c Context, val reflect.Value) E.Either[Errors, GA] {
	return F.Pipe1(
		toUnknownArray(item, c, val),
		flattenUnknownArray[GA, A],
	)
}

func validateArray[GA ~[]A, A any](item Validate[reflect.Value, A]) func(i reflect.Value, c Context) E.Either[Errors, GA] {
	var r func(i reflect.Value, c Context) E.Either[Errors, GA]

	r = func(i reflect.Value, c Context) E.Either[Errors, GA] {
		// check for unknow array
		switch i.Kind() {
		case reflect.Slice:
			return toValidatedArray[GA](item, c, i)
		case reflect.Array:
			return toValidatedArray[GA](item, c, i)
		case reflect.Pointer:
			return r(i.Elem(), c)
		default:
			return Failure[GA](c, fmt.Sprintf("Type %T is neither an array nor a slice nor a pointer to these values", i))
		}
	}

	return r
}

// ArrayG returns the type validator for an array
func ArrayG[GA ~[]A, A any](item Validate[reflect.Value, A]) *Type[GA, GA, reflect.Value] {
	return FromValidate(validateArray[GA, A](item))
}

// Array returns the type validator for an array
func Array[A any](item Validate[reflect.Value, A]) *Type[[]A, []A, reflect.Value] {
	return ArrayG[[]A, A](item)
}
