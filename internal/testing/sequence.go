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

package testing

import (
	"fmt"
	"testing"

	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/pointed"
	"github.com/stretchr/testify/assert"
)

// SequenceArrayTest tests if the sequence operation works in case the operation cannot error
func SequenceArrayTest[
	HKTA,
	HKTB,
	HKTAA any, // HKT[[]A]
](
	eq EQ.Eq[HKTB],

	pa pointed.Pointed[string, HKTA],
	pb pointed.Pointed[bool, HKTB],
	faa functor.Functor[[]string, bool, HKTAA, HKTB],
	seq func([]HKTA) HKTAA,
) func(count int) func(t *testing.T) {

	return func(count int) func(t *testing.T) {

		exp := make([]string, count)
		good := make([]HKTA, count)
		for i := 0; i < count; i++ {
			val := fmt.Sprintf("TestData %d", i)
			exp[i] = val
			good[i] = pa.Of(val)
		}

		return func(t *testing.T) {
			res := F.Pipe2(
				good,
				seq,
				faa.Map(func(act []string) bool {
					return assert.Equal(t, exp, act)
				}),
			)
			assert.True(t, eq.Equals(res, pb.Of(true)))
		}
	}
}
