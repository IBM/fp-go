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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// these functions curry a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func Curry0[GEA ~func(R) ET.Either[error, A], R, A any](f func(R) (A, error)) GEA {
	return G.Curry0[GEA](ET.Eitherize1(f))
}

func Curry1[GEA ~func(R) ET.Either[error, A], R, T1, A any](f func(R, T1) (A, error)) func(T1) GEA {
	return G.Curry1[GEA](ET.Eitherize2(f))
}

func Curry2[GEA ~func(R) ET.Either[error, A], R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1) func(T2) GEA {
	return G.Curry2[GEA](ET.Eitherize3(f))
}

func Curry3[GEA ~func(R) ET.Either[error, A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1) func(T2) func(T3) GEA {
	return G.Curry3[GEA](ET.Eitherize4(f))
}

func Uncurry1[GEA ~func(R) ET.Either[error, A], R, T1, A any](f func(T1) GEA) func(R, T1) (A, error) {
	return ET.Uneitherize2(G.Uncurry1(f))
}

func Uncurry2[GEA ~func(R) ET.Either[error, A], R, T1, T2, A any](f func(T1) func(T2) GEA) func(R, T1, T2) (A, error) {
	return ET.Uneitherize3(G.Uncurry2(f))
}

func Uncurry3[GEA ~func(R) ET.Either[error, A], R, T1, T2, T3, A any](f func(T1) func(T2) func(T3) GEA) func(R, T1, T2, T3) (A, error) {
	return ET.Uneitherize4(G.Uncurry3(f))
}
