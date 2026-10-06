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

package readerreaderioeither

import (
	E "github.com/IBM/fp-go/v2/either"
)

// exCfg is the outer environment of the examples, e.g. application wide settings.
type exCfg struct {
	Factor int
}

// exReq is the inner environment of the examples, e.g. per request data.
type exReq struct {
	User string
}

// exRRIOE fixes the environments and the error type of the examples.
type exRRIOE[A any] = ReaderReaderIOEither[exCfg, exReq, string, A]

var (
	exConfig  = exCfg{Factor: 10}
	exRequest = exReq{User: "alice"}
)

// exRun provides both environments and executes the IO.
func exRun[A any](m exRRIOE[A]) Either[string, A] {
	return m(exConfig)(exRequest)()
}

// exFactor reads the factor from the outer environment.
func exFactor(c exCfg) int { return c.Factor }

// exUser reads the user from the inner environment.
func exUser(r exReq) string { return r.User }

// exPositive fails for numbers that are not positive.
func exPositive(n int) Either[string, int] {
	if n > 0 {
		return E.Right[string](n)
	}
	return E.Left[int]("not positive")
}
