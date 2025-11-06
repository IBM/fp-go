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
	"fmt"
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
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

// SequenceArrayErrorTest tests if the sequence operation works in case the operation can error
func SequenceArrayErrorTest[
	HKTA,
	HKTB,
	HKTAA any, // HKT[[]A]
](
	eq EQ.Eq[HKTB],

	left func(error) HKTA,
	leftB func(error) HKTB,
	pa pointed.Pointed[string, HKTA],
	pb pointed.Pointed[bool, HKTB],
	faa functor.Functor[[]string, bool, HKTAA, HKTB],
	seq func([]HKTA) HKTAA,
) func(count int) func(t *testing.T) {

	return func(count int) func(t *testing.T) {

		expGood := make([]string, count)
		good := make([]HKTA, count)
		expBad := make([]error, count)
		bad := make([]HKTA, count)

		for i := 0; i < count; i++ {
			goodVal := fmt.Sprintf("TestData %d", i)
			badVal := fmt.Errorf("ErrorData %d", i)
			expGood[i] = goodVal
			good[i] = pa.Of(goodVal)
			expBad[i] = badVal
			bad[i] = left(badVal)
		}

		total := 1 << count

		return func(t *testing.T) {
			// test the good case
			res := F.Pipe2(
				good,
				seq,
				faa.Map(func(act []string) bool {
					return assert.Equal(t, expGood, act)
				}),
			)
			assert.True(t, eq.Equals(res, pb.Of(true)))
			// iterate and test the bad cases
			for i := 1; i < total; i++ {
				// run the test
				t.Run(fmt.Sprintf("Bitmask test %d", i), func(t1 *testing.T) {
					// the actual
					act := make([]HKTA, count)
					// the expected error
					var exp error
					// prepare the values bases on the bit mask
					mask := 1
					for j := 0; j < count; j++ {
						if (i & mask) == 0 {
							act[j] = good[j]
						} else {
							act[j] = bad[j]
							if exp == nil {
								exp = expBad[j]
							}
						}
						mask <<= 1
					}
					// test the good case
					res := F.Pipe2(
						act,
						seq,
						faa.Map(func(act []string) bool {
							return assert.Equal(t, expGood, act)
						}),
					)
					// validate the error
					assert.True(t, eq.Equals(res, leftB(exp)))
				})
			}
		}
	}
}

// SequenceRecordTest tests if the sequence operation works in case the operation cannot error
func SequenceRecordTest[
	HKTA,
	HKTB,
	HKTAA any, // HKT[map[string]string]
](
	eq EQ.Eq[HKTB],

	pa pointed.Pointed[string, HKTA],
	pb pointed.Pointed[bool, HKTB],
	faa functor.Functor[map[string]string, bool, HKTAA, HKTB],
	seq func(map[string]HKTA) HKTAA,
) func(count int) func(t *testing.T) {

	return func(count int) func(t *testing.T) {

		exp := make(map[string]string)
		good := make(map[string]HKTA)
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("KeyData %d", i)
			val := fmt.Sprintf("ValueData %d", i)
			exp[key] = val
			good[key] = pa.Of(val)
		}

		return func(t *testing.T) {
			res := F.Pipe2(
				good,
				seq,
				faa.Map(func(act map[string]string) bool {
					return assert.Equal(t, exp, act)
				}),
			)
			assert.True(t, eq.Equals(res, pb.Of(true)))
		}
	}
}
