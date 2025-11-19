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

package ioresult

import (
	"encoding/json"
	"log"

	"github.com/IBM/fp-go/v2/idiomatic/result"
)

// LogJSON converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst]
func LogJSON[A any](prefix string) Kleisli[A, any] {
	return func(a A) IOResult[any] {
		// convert to a string
		b, jsonerr := json.MarshalIndent(a, "", "  ")
		// log this
		return func() (any, error) {
			if jsonerr != nil {
				return result.Left[any](jsonerr)
			}
			log.Printf(prefix, string(b))
			return result.Of[any](b)
		}
	}
}
