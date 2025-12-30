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

package assert_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/IBM/fp-go/v2/assert"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
)

// Example_basicAssertions demonstrates basic equality and inequality assertions
func Example_basicAssertions() {
	// This would be in a real test function
	var t *testing.T // placeholder for example

	// Basic equality
	value := 42
	assert.Equal(42)(value)(t)

	// String equality
	name := "Alice"
	assert.Equal("Alice")(name)(t)

	// Inequality
	assert.NotEqual(10)(value)(t)
}

// Example_arrayAssertions demonstrates array-related assertions
func Example_arrayAssertions() {
	var t *testing.T // placeholder for example

	numbers := []int{1, 2, 3, 4, 5}

	// Check array is not empty
	assert.ArrayNotEmpty(numbers)(t)

	// Check array length
	assert.ArrayLength[int](5)(numbers)(t)

	// Check array contains a value
	assert.ArrayContains(3)(numbers)(t)
}

// Example_mapAssertions demonstrates map-related assertions
func Example_mapAssertions() {
	var t *testing.T // placeholder for example

	config := map[string]int{
		"timeout": 30,
		"retries": 3,
		"maxSize": 1000,
	}

	// Check map is not empty
	assert.RecordNotEmpty(config)(t)

	// Check map length
	assert.RecordLength[string, int](3)(config)(t)

	// Check map contains key
	assert.ContainsKey[int]("timeout")(config)(t)

	// Check map does not contain key
	assert.NotContainsKey[int]("unknown")(config)(t)
}

// Example_errorAssertions demonstrates error-related assertions
func Example_errorAssertions() {
	var t *testing.T // placeholder for example

	// Assert no error
	err := doSomethingSuccessful()
	assert.NoError(err)(t)

	// Assert error exists
	err2 := doSomethingThatFails()
	assert.Error(err2)(t)
}

// Example_resultAssertions demonstrates Result type assertions
func Example_resultAssertions() {
	var t *testing.T // placeholder for example

	// Assert success
	successResult := result.Of(42)
	assert.Success(successResult)(t)

	// Assert failure
	failureResult := result.Left[int](errors.New("something went wrong"))
	assert.Failure(failureResult)(t)
}

// Example_predicateAssertions demonstrates custom predicate assertions
func Example_predicateAssertions() {
	var t *testing.T // placeholder for example

	// Test if a number is positive
	isPositive := N.MoreThan(0)
	assert.That(isPositive)(42)(t)

	// Test if a string is uppercase
	isUppercase := func(s string) bool { return s == strings.ToUpper(s) }
	assert.That(isUppercase)("HELLO")(t)

	// Test if a number is even
	isEven := func(n int) bool { return n%2 == 0 }
	assert.That(isEven)(10)(t)
}

// Example_allOf demonstrates combining multiple assertions
func Example_allOf() {
	var t *testing.T // placeholder for example

	type User struct {
		Name   string
		Age    int
		Active bool
	}

	user := User{Name: "Alice", Age: 30, Active: true}

	// Combine multiple assertions
	assertions := assert.AllOf([]assert.Reader{
		assert.Equal("Alice")(user.Name),
		assert.Equal(30)(user.Age),
		assert.Equal(true)(user.Active),
	})

	assertions(t)
}

// Example_runAll demonstrates running named test cases
func Example_runAll() {
	var t *testing.T // placeholder for example

	testcases := map[string]assert.Reader{
		"addition":       assert.Equal(4)(2 + 2),
		"multiplication": assert.Equal(6)(2 * 3),
		"subtraction":    assert.Equal(1)(3 - 2),
		"division":       assert.Equal(2)(10 / 5),
	}

	assert.RunAll(testcases)(t)
}

// Example_local demonstrates focusing assertions on specific properties
func Example_local() {
	var t *testing.T // placeholder for example

	type User struct {
		Name string
		Age  int
	}

	// Create an assertion that checks if age is positive
	ageIsPositive := assert.That(func(age int) bool { return age > 0 })

	// Focus this assertion on the Age field of User
	userAgeIsPositive := assert.Local(func(u User) int { return u.Age })(ageIsPositive)

	// Now we can test the whole User object
	user := User{Name: "Alice", Age: 30}
	userAgeIsPositive(user)(t)
}

// Example_composableAssertions demonstrates building complex assertions
func Example_composableAssertions() {
	var t *testing.T // placeholder for example

	type Config struct {
		Host    string
		Port    int
		Timeout int
		Retries int
	}

	config := Config{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
		Retries: 3,
	}

	// Create focused assertions for each field
	validHost := assert.Local(func(c Config) string { return c.Host })(
		assert.StringNotEmpty,
	)

	validPort := assert.Local(func(c Config) int { return c.Port })(
		assert.That(func(p int) bool { return p > 0 && p < 65536 }),
	)

	validTimeout := assert.Local(func(c Config) int { return c.Timeout })(
		assert.That(func(t int) bool { return t > 0 }),
	)

	validRetries := assert.Local(func(c Config) int { return c.Retries })(
		assert.That(func(r int) bool { return r >= 0 }),
	)

	// Combine all assertions
	validConfig := assert.AllOf([]assert.Reader{
		validHost(config),
		validPort(config),
		validTimeout(config),
		validRetries(config),
	})

	validConfig(t)
}

// Helper functions for examples
func doSomethingSuccessful() error {
	return nil
}

func doSomethingThatFails() error {
	return errors.New("operation failed")
}
