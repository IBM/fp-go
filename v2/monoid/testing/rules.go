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

package testing

import (
	"testing"

	AR "github.com/IBM/fp-go/v2/internal/array"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

func assertLaws[A any](t *testing.T, m M.Monoid[A]) func(a A) bool {
	e := m.Empty()
	return func(a A) bool {
		return assert.Equal(t, a, m.Concat(a, e), "Monoid right identity") &&
			assert.Equal(t, a, m.Concat(e, a), "Monoid left identity")
	}
}

// AssertLaws asserts the monoid laws for a dataset
func AssertLaws[A any](t *testing.T, m M.Monoid[A]) func(data []A) bool {
	law := assertLaws(t, m)

	return func(data []A) bool {
		return AR.Reduce(data, func(result bool, value A) bool {
			return result && law(value)
		}, true)
	}
}
