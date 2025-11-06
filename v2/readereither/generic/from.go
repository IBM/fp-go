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

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[GEA ~func(R) ET.Either[error, A], R, A any](f func(R) (A, error)) func() GEA {
	return G.From0[GEA](ET.Eitherize1(f))
}

func From1[GEA ~func(R) ET.Either[error, A], R, T1, A any](f func(R, T1) (A, error)) func(T1) GEA {
	return G.From1[GEA](ET.Eitherize2(f))
}

func From2[GEA ~func(R) ET.Either[error, A], R, T1, T2, A any](f func(R, T1, T2) (A, error)) func(T1, T2) GEA {
	return G.From2[GEA](ET.Eitherize3(f))
}

func From3[GEA ~func(R) ET.Either[error, A], R, T1, T2, T3, A any](f func(R, T1, T2, T3) (A, error)) func(T1, T2, T3) GEA {
	return G.From3[GEA](ET.Eitherize4(f))
}
