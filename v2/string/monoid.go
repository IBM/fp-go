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
	M "github.com/IBM/fp-go/v2/monoid"
)

// Monoid is the monoid implementing string concatenation with empty string as identity
var Monoid = M.MakeMonoid(concat, "")

// IntersperseMonoid creates a monoid that concatenates strings with a middle string in between,
// with empty string as identity
func IntersperseMonoid(middle string) M.Monoid[string] {
	return M.MakeMonoid(Intersperse(middle), "")
}
