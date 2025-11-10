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
	EQ "github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/result"
)

// Eq implements the equals predicate for values contained in the IOResult monad
//
//go:inline
func Eq[A any](eq EQ.Eq[Result[A]]) EQ.Eq[IOResult[A]] {
	return io.Eq(eq)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
//
//go:inline
func FromStrictEquals[A comparable]() EQ.Eq[IOResult[A]] {
	return Eq(result.FromStrictEquals[A]())
}
