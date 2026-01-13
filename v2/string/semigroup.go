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

package string

import (
	S "github.com/IBM/fp-go/v2/semigroup"
)

// concat concatenates two strings using simple string concatenation.
// This is an internal helper function used by the Semigroup and Monoid implementations.
func concat(left, right string) string {
	return left + right
}

// Semigroup is the semigroup implementing string concatenation
var Semigroup = S.MakeSemigroup(concat)

// IntersperseSemigroup creates a semigroup that concatenates strings with a middle string in between
func IntersperseSemigroup(middle string) S.Semigroup[string] {
	return S.MakeSemigroup(Intersperse(middle))
}
