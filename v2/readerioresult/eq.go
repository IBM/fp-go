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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/eq"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

// Eq implements the equals predicate for values contained in the IOEither monad
//
//go:inline
func Eq[R, A any](eq eq.Eq[Result[A]]) func(R) eq.Eq[ReaderIOResult[R, A]] {
	return RIOE.Eq[R](eq)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
//
//go:inline
func FromStrictEquals[R any, A comparable]() func(R) eq.Eq[ReaderIOResult[R, A]] {
	return RIOE.FromStrictEquals[R, error, A]()
}
