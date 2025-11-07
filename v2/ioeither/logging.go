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

package ioeither

import (
	"encoding/json"
	"log"

	"github.com/IBM/fp-go/v2/bytes"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
)

// LogJSON converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst]
func LogJSON[A any](prefix string) Kleisli[error, A, any] {
	return func(a A) IOEither[error, any] {
		// log this
		return function.Pipe3(
			either.TryCatchError(json.MarshalIndent(a, "", "  ")),
			either.Map[error](bytes.ToString),
			FromEither[error, string],
			Chain(func(data string) IOEither[error, any] {
				return FromImpure[error](func() {
					log.Printf(prefix, data)
				})
			}),
		)
	}
}
