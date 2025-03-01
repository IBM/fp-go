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

package ioeither

import (
	G "github.com/IBM/fp-go/v2/ioeither/generic"
)

// LogJson converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst]
//
// Deprecated: use [LogJSON] instead
func LogJson[A any](prefix string) func(A) IOEither[error, any] {
	return G.LogJson[IOEither[error, any], A](prefix)
}

// LogJSON converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst]
func LogJSON[A any](prefix string) func(A) IOEither[error, any] {
	return G.LogJSON[IOEither[error, any], A](prefix)
}
