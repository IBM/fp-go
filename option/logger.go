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

package option

import (
	"log"

	F "github.com/IBM/fp-go/function"
	L "github.com/IBM/fp-go/logging"
)

func _log[A any](left func(string, ...any), right func(string, ...any), prefix string) func(Option[A]) Option[A] {
	return Fold(
		func() Option[A] {
			left("%s", prefix)
			return None[A]()
		},
		func(a A) Option[A] {
			right("%s: %v", prefix, a)
			return Some(a)
		})
}

func Logger[A any](loggers ...*log.Logger) func(string) func(Option[A]) Option[A] {
	left, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) func(Option[A]) Option[A] {
		delegate := _log[A](left, right, prefix)
		return func(ma Option[A]) Option[A] {
			return F.Pipe1(
				delegate(ma),
				ChainTo[A](ma),
			)
		}
	}
}
