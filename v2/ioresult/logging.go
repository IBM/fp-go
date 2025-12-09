// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/pair"
)

// LogJSON converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst]
//
//go:inline
func LogJSON[A any](prefix string) Kleisli[A, string] {
	return ioeither.LogJSON[A](prefix)
}

//go:inline
func LogEntryExitF[A, STARTTOKEN, ANY any](
	onEntry IO[STARTTOKEN],
	onExit io.Kleisli[pair.Pair[STARTTOKEN, Result[A]], ANY],
) Operator[A, A] {
	return ioeither.LogEntryExitF(onEntry, onExit)
}

//go:inline
func LogEntryExit[A any](name string) Operator[A, A] {
	return ioeither.LogEntryExit[error, A](name)
}
