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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse(t *testing.T) {
	// Test with positive predicate
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	// Test value that passes predicate
	result, err := F.Pipe2(5, Of, FilterOrElse(isPositive, onNegative))()
	assert.NoError(t, err)
	assert.Equal(t, 5, result)

	// Test value that fails predicate
	_, err = F.Pipe2(-3, Of, FilterOrElse(isPositive, onNegative))()
	assert.Error(t, err)
	assert.Equal(t, "-3 is not positive", err.Error())

	// Test error value (should pass through unchanged)
	originalError := errors.New("original error")
	_, err = F.Pipe2(originalError, Left[int], FilterOrElse(isPositive, onNegative))()
	assert.Error(t, err)
	assert.Equal(t, originalError, err)
}

func TestFilterOrElse_WithChain(t *testing.T) {
	// Test FilterOrElse in a chain with other IO operations
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	double := func(x int) (int, error) { return x * 2, nil }

	// Test successful chain
	result, err := F.Pipe3(5, Of, FilterOrElse(isPositive, onNegative), ChainResultK(double))()
	assert.NoError(t, err)
	assert.Equal(t, 10, result)

	// Test chain with filter failure
	_, err = F.Pipe3(-5, Of, FilterOrElse(isPositive, onNegative), ChainResultK(double))()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not positive")
}

func TestFilterOrElse_MultipleFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	isEven := func(x int) bool { return x%2 == 0 }
	onOdd := func(x int) error { return fmt.Errorf("%d is not even", x) }

	// Test value that passes both filters
	result, err := F.Pipe3(4, Of, FilterOrElse(isPositive, onNegative), FilterOrElse(isEven, onOdd))()
	assert.NoError(t, err)
	assert.Equal(t, 4, result)

	// Test value that fails second filter
	_, err = F.Pipe3(3, Of, FilterOrElse(isPositive, onNegative), FilterOrElse(isEven, onOdd))()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not even")
}
