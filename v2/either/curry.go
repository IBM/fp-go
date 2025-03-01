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

package either

import (
	F "github.com/IBM/fp-go/v2/function"
)

// these function curry a golang function that returns an error into its curried version that returns an either

func Curry0[R any](f func() (R, error)) func() Either[error, R] {
	return Eitherize0(f)
}

func Curry1[T1, R any](f func(T1) (R, error)) func(T1) Either[error, R] {
	return Eitherize1(f)
}

func Curry2[T1, T2, R any](f func(T1, T2) (R, error)) func(T1) func(T2) Either[error, R] {
	return F.Curry2(Eitherize2(f))
}

func Curry3[T1, T2, T3, R any](f func(T1, T2, T3) (R, error)) func(T1) func(T2) func(T3) Either[error, R] {
	return F.Curry3(Eitherize3(f))
}

func Curry4[T1, T2, T3, T4, R any](f func(T1, T2, T3, T4) (R, error)) func(T1) func(T2) func(T3) func(T4) Either[error, R] {
	return F.Curry4(Eitherize4(f))
}

func Uncurry0[R any](f func() Either[error, R]) func() (R, error) {
	return func() (R, error) {
		return UnwrapError(f())
	}
}

func Uncurry1[T1, R any](f func(T1) Either[error, R]) func(T1) (R, error) {
	uc := F.Uncurry1(f)
	return func(t1 T1) (R, error) {
		return UnwrapError(uc(t1))
	}
}

func Uncurry2[T1, T2, R any](f func(T1) func(T2) Either[error, R]) func(T1, T2) (R, error) {
	uc := F.Uncurry2(f)
	return func(t1 T1, t2 T2) (R, error) {
		return UnwrapError(uc(t1, t2))
	}
}

func Uncurry3[T1, T2, T3, R any](f func(T1) func(T2) func(T3) Either[error, R]) func(T1, T2, T3) (R, error) {
	uc := F.Uncurry3(f)
	return func(t1 T1, t2 T2, t3 T3) (R, error) {
		return UnwrapError(uc(t1, t2, t3))
	}
}

func Uncurry4[T1, T2, T3, T4, R any](f func(T1) func(T2) func(T3) func(T4) Either[error, R]) func(T1, T2, T3, T4) (R, error) {
	uc := F.Uncurry4(f)
	return func(t1 T1, t2 T2, t3 T3, t4 T4) (R, error) {
		return UnwrapError(uc(t1, t2, t3, t4))
	}
}
