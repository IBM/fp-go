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
	"log"

	Logging "github.com/IBM/fp-go/logging"
)

func Logger[GA ~func() any, A any](loggers ...*log.Logger) func(string) func(A) GA {
	_, right := Logging.LoggingCallbacks(loggers...)
	return func(prefix string) func(A) GA {
		return func(a A) GA {
			return FromImpure[GA](func() {
				right("%s: %v", prefix, a)
			})
		}
	}
}

func Logf[GA ~func() any, A any](prefix string) func(A) GA {
	return func(a A) GA {
		return FromImpure[GA](func() {
			log.Printf(prefix, a)
		})
	}
}
