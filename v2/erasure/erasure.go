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

package erasure

import (
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
)

// Erase converts a variable of type T to an any by returning a pointer to that variable
func Erase[T any](t T) any {
	return &t
}

// Unerase converts an erased variable back to its original value
func Unerase[T any](t any) T {
	return *t.(*T)
}

// SafeUnerase converts an erased variable back to its original value
func SafeUnerase[T any](t any) E.Either[error, T] {
	return F.Pipe2(
		t,
		E.ToType[*T](errors.OnSome[any]("Value of type [%T] is not erased")),
		E.Map[error](F.Deref[T]),
	)
}

// Erase0 converts a type safe function into an erased function
func Erase0[T1 any](f func() T1) func() any {
	return F.Nullary2(f, Erase[T1])
}

// Erase1 converts a type safe function into an erased function
func Erase1[T1, T2 any](f func(T1) T2) func(any) any {
	return F.Flow3(
		Unerase[T1],
		f,
		Erase[T2],
	)
}

// Erase2 converts a type safe function into an erased function
func Erase2[T1, T2, T3 any](f func(T1, T2) T3) func(any, any) any {
	return func(t1, t2 any) any {
		return Erase(f(Unerase[T1](t1), Unerase[T2](t2)))
	}
}
