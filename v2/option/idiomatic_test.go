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

func TestMonadChainIdiomatic(t *testing.T) {
	assert.Equal(t, Some(42), MonadChainIdiomatic(Some("42"), parseInt))
	assert.Equal(t, None[int](), MonadChainIdiomatic(Some("abc"), parseInt))
	assert.Equal(t, None[int](), MonadChainIdiomatic(None[string](), parseInt))
}

func TestChainIdiomatic(t *testing.T) {
	parse := ChainIdiomatic(parseInt)
	assert.Equal(t, Some(7), parse(Some("7")))
	assert.Equal(t, None[int](), parse(Some("bad")))
	assert.Equal(t, None[int](), parse(None[string]()))
}

func TestMonadChainFirstIdiomatic(t *testing.T) {
	assert.Equal(t, Some(5), MonadChainFirstIdiomatic(Some(5), positiveToString))
	assert.Equal(t, None[int](), MonadChainFirstIdiomatic(Some(-1), positiveToString))
	assert.Equal(t, None[int](), MonadChainFirstIdiomatic(None[int](), positiveToString))
}

func TestChainFirstIdiomatic(t *testing.T) {
	keepIfPositive := ChainFirstIdiomatic(positiveToString)
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

func TestBindIdiomatic(t *testing.T) {
	res := F.Pipe3(
		Do(utils.Empty),
		BindIdiomatic(utils.SetLastName, getLastNameIdiomatic),
		BindIdiomatic(utils.SetGivenName, getGivenNameIdiomatic),
		Map(utils.GetFullName),
	)
	assert.Equal(t, Of("John Doe"), res)
}

func TestBindIdiomaticNone(t *testing.T) {
	res := F.Pipe3(
		Do(utils.Empty),
		BindIdiomatic(utils.SetLastName, func(_ utils.Initial) (string, bool) { return "", false }),
		BindIdiomatic(utils.SetGivenName, getGivenNameIdiomatic),
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

func TestBindLIdiomatic(t *testing.T) {
	inc := BindLIdiomatic(counterValueLens, incrementUnder100)
	assert.Equal(t, Some(counterState{value: 43}), inc(Some(counterState{value: 42})))
	assert.Equal(t, None[counterState](), inc(Some(counterState{value: 100})))
	assert.Equal(t, None[counterState](), inc(None[counterState]()))
}

func TestTraverseArrayIdiomatic(t *testing.T) {
	parse := TraverseArrayIdiomatic(parseInt)
	assert.Equal(t, Some([]int{1, 2, 3}), parse([]string{"1", "2", "3"}))
	assert.Equal(t, None[[]int](), parse([]string{"1", "x", "3"}))
	assert.Equal(t, Some([]int{}), parse([]string{}))
}

func TestTraverseIterIdiomatic_AllSome(t *testing.T) {
	parse := TraverseIterIdiomatic(parseInt)
	result := parse(slices.Values([]string{"1", "2", "3"}))
	assert.True(t, IsSome(result))
	collected := MonadFold(result, func() []int { return nil }, collectSeq[int])
	assert.Equal(t, []int{1, 2, 3}, collected)
}

func TestTraverseIterIdiomatic_ContainsNone(t *testing.T) {
	parse := TraverseIterIdiomatic(parseInt)
	result := parse(slices.Values([]string{"1", "x", "3"}))
	assert.True(t, IsNone(result))
}

func TestTraverseIterIdiomatic_Empty(t *testing.T) {
	parse := TraverseIterIdiomatic(parseInt)
	result := parse(slices.Values([]string{}))
	assert.True(t, IsSome(result))
	collected := MonadFold(result, func() []int { return nil }, collectSeq[int])
	assert.Empty(t, collected)
}
