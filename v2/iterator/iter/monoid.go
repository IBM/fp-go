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

package iter

import (
	G "github.com/IBM/fp-go/v2/internal/iter"
	M "github.com/IBM/fp-go/v2/monoid"
)

// Monoid returns a Monoid instance for Seq[T].
// The monoid's concat operation concatenates sequences, and the empty value is an empty sequence.
//
// Example:
//
//	m := Monoid[int]()
//	seq1 := From(1, 2)
//	seq2 := From(3, 4)
//	result := m.Concat(seq1, seq2)
//	// yields: 1, 2, 3, 4
//
//go:inline
func Monoid[T any]() M.Monoid[Seq[T]] {
	return G.Monoid[Seq[T]]()
}
