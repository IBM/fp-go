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

package stateless

import (
	"fmt"
	"testing"

	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestZipWith(t *testing.T) {
	left := From(1, 2, 3)
	right := From("a", "b", "c", "d")

	res := ZipWith(left, right, func(l int, r string) string {
		return fmt.Sprintf("%s%d", r, l)
	})

	assert.Equal(t, ToArray(From("a1", "b2", "c3")), ToArray(res))
}

func TestZip(t *testing.T) {
	left := From(1, 2, 3)
	right := From("a", "b", "c", "d")

	res := Zip[string](left)(right)

	assert.Equal(t, ToArray(From(P.MakePair("a", 1), P.MakePair("b", 2), P.MakePair("c", 3))), ToArray(res))
}
