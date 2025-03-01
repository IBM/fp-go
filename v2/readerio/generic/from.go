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

package generic

import (
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// these functions From a golang function with the context as the firsr parameter into a either reader with the context as the last parameter
// this goes back to the advice in https://pkg.go.dev/context to put the context as a first parameter as a convention

func From0[GEA ~func(R) GIOA, GIOA ~func() A, R, A any](f func(R) GIOA) func() GEA {
	return G.From0[GEA](f)
}

func From1[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, A any](f func(R, T1) GIOA) func(T1) GEA {
	return G.From1[GEA](f)
}

func From2[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, T2, A any](f func(R, T1, T2) GIOA) func(T1, T2) GEA {
	return G.From2[GEA](f)
}

func From3[GEA ~func(R) GIOA, GIOA ~func() A, R, T1, T2, T3, A any](f func(R, T1, T2, T3) GIOA) func(T1, T2, T3) GEA {
	return G.From3[GEA](f)
}
