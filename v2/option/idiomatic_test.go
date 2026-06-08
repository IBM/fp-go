// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"slices"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

func parseInt(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	return n, err == nil
}

func positiveToString(n int) (string, bool) {
	if n > 0 {
		return strconv.Itoa(n), true
	}
	return "", false
}

func TestMonadChainI(t *testing.T) {
	assert.Equal(t, Some(42), MonadChainI(Some("42"), parseInt))
	assert.Equal(t, None[int](), MonadChainI(Some("abc"), parseInt))
	assert.Equal(t, None[int](), MonadChainI(None[string](), parseInt))
}

func TestChainI(t *testing.T) {
	parse := ChainI(parseInt)
	assert.Equal(t, Some(7), parse(Some("7")))
	assert.Equal(t, None[int](), parse(Some("bad")))
	assert.Equal(t, None[int](), parse(None[string]()))
}

func TestMonadChainFirstI(t *testing.T) {
	assert.Equal(t, Some(5), MonadChainFirstI(Some(5), positiveToString))
	assert.Equal(t, None[int](), MonadChainFirstI(Some(-1), positiveToString))
	assert.Equal(t, None[int](), MonadChainFirstI(None[int](), positiveToString))
}

func TestChainFirstI(t *testing.T) {
	keepIfPositive := ChainFirstI(positiveToString)
	assert.Equal(t, Some(5), keepIfPositive(Some(5)))
	assert.Equal(t, None[int](), keepIfPositive(Some(-1)))
	assert.Equal(t, None[int](), keepIfPositive(None[int]()))
}

func getLastNameIdiomatic(_ utils.Initial) (string, bool) {
	return "Doe", true
}

func getGivenNameIdiomatic(_ utils.WithLastName) (string, bool) {
	return "John", true
}

func TestBindI(t *testing.T) {
	res := F.Pipe3(
		Do(utils.Empty),
		BindI(utils.SetLastName, getLastNameIdiomatic),
		BindI(utils.SetGivenName, getGivenNameIdiomatic),
		Map(utils.GetFullName),
	)
	assert.Equal(t, Of("John Doe"), res)
}

func TestBindINone(t *testing.T) {
	res := F.Pipe3(
		Do(utils.Empty),
		BindI(utils.SetLastName, func(_ utils.Initial) (string, bool) { return "", false }),
		BindI(utils.SetGivenName, getGivenNameIdiomatic),
		Map(utils.GetFullName),
	)
	assert.Equal(t, None[string](), res)
}

type counterState struct {
	value int
}

var counterValueLens = L.MakeLens(
	func(c counterState) int { return c.value },
	func(c counterState, v int) counterState { c.value = v; return c },
)

func incrementUnder100(v int) (int, bool) {
	if v >= 100 {
		return 0, false
	}
	return v + 1, true
}

func TestBindIL(t *testing.T) {
	inc := BindIL(counterValueLens, incrementUnder100)
	assert.Equal(t, Some(counterState{value: 43}), inc(Some(counterState{value: 42})))
	assert.Equal(t, None[counterState](), inc(Some(counterState{value: 100})))
	assert.Equal(t, None[counterState](), inc(None[counterState]()))
}

func TestTraverseArrayI(t *testing.T) {
	parse := TraverseArrayI(parseInt)
	assert.Equal(t, Some([]int{1, 2, 3}), parse([]string{"1", "2", "3"}))
	assert.Equal(t, None[[]int](), parse([]string{"1", "x", "3"}))
	assert.Equal(t, Some([]int{}), parse([]string{}))
}

func TestTraverseIterI_AllSome(t *testing.T) {
	parse := TraverseIterI(parseInt)
	result := parse(slices.Values([]string{"1", "2", "3"}))
	assert.True(t, IsSome(result))
	collected := MonadFold(result, func() []int { return nil }, collectSeq[int])
	assert.Equal(t, []int{1, 2, 3}, collected)
}

func TestTraverseIterI_ContainsNone(t *testing.T) {
	parse := TraverseIterI(parseInt)
	result := parse(slices.Values([]string{"1", "x", "3"}))
	assert.True(t, IsNone(result))
}

func TestTraverseIterI_Empty(t *testing.T) {
	parse := TraverseIterI(parseInt)
	result := parse(slices.Values([]string{}))
	assert.True(t, IsSome(result))
	collected := MonadFold(result, func() []int { return nil }, collectSeq[int])
	assert.Empty(t, collected)
}
