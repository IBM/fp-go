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

package record

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// isEven is a reusable predicate used across partition tests.
func isEven(v int) bool { return v%2 == 0 }

// TestMonadPartition_SplitsOnPredicate verifies that matching entries go to
// Tail and non-matching entries go to Head.
func TestMonadPartition_SplitsOnPredicate(t *testing.T) {
	src := Record[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := MonadPartition(src, isEven)

	assert.Equal(t, Record[string, int]{"a": 1, "c": 3}, P.Head(result))
	assert.Equal(t, Record[string, int]{"b": 2, "d": 4}, P.Tail(result))
}

// TestMonadPartition_AllMatch verifies that Head is empty when all entries match.
func TestMonadPartition_AllMatch(t *testing.T) {
	src := Record[string, int]{"a": 2, "b": 4}

	result := MonadPartition(src, isEven)

	assert.Empty(t, P.Head(result))
	assert.Equal(t, src, P.Tail(result))
}

// TestMonadPartition_NoneMatch verifies that Tail is empty when no entries match.
func TestMonadPartition_NoneMatch(t *testing.T) {
	src := Record[string, int]{"a": 1, "b": 3}

	result := MonadPartition(src, isEven)

	assert.Equal(t, src, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartition_EmptyRecord verifies that an empty record produces two
// empty records.
func TestMonadPartition_EmptyRecord(t *testing.T) {
	src := Record[string, int]{}

	result := MonadPartition(src, isEven)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartition_NilRecord verifies that a nil map is treated identically
// to an empty record.
func TestMonadPartition_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := MonadPartition(src, isEven)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionWithIndex_SplitsOnKeyAndValue verifies that the predicate
// receives both key and value, and entries are routed correctly.
func TestMonadPartitionWithIndex_SplitsOnKeyAndValue(t *testing.T) {
	src := Record[string, int]{
		"keep_1": 1,
		"keep_2": 2,
		"drop_3": 3,
		"drop_4": 4,
	}

	// Keep only entries whose key starts with "keep"
	result := MonadPartitionWithIndex(src, func(k string, _ int) bool {
		return strings.HasPrefix(k, "keep")
	})

	assert.Equal(t, Record[string, int]{"drop_3": 3, "drop_4": 4}, P.Head(result))
	assert.Equal(t, Record[string, int]{"keep_1": 1, "keep_2": 2}, P.Tail(result))
}

// TestMonadPartitionWithIndex_ValueAndKeyTogether verifies that the predicate
// can use both key and value together.
func TestMonadPartitionWithIndex_ValueAndKeyTogether(t *testing.T) {
	src := Record[string, int]{
		"a":   1,
		"bb":  2,
		"ccc": 3,
	}

	// Match when key length equals value
	result := MonadPartitionWithIndex(src, func(k string, v int) bool {
		return len(k) == v
	})

	assert.Equal(t, Record[string, int]{}, P.Head(result))
	assert.Equal(t, Record[string, int]{"a": 1, "bb": 2, "ccc": 3}, P.Tail(result))
}

// TestMonadPartitionWithIndex_EmptyRecord verifies that an empty record
// produces two empty records.
func TestMonadPartitionWithIndex_EmptyRecord(t *testing.T) {
	src := Record[string, int]{}

	result := MonadPartitionWithIndex(src, func(k string, v int) bool { return true })

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestPartition_CurriedReturnsFunction verifies that Partition returns a
// reusable function applied to different records.
func TestPartition_CurriedReturnsFunction(t *testing.T) {
	partitionEvens := Partition[string](isEven)

	src1 := Record[string, int]{"a": 1, "b": 2}
	src2 := Record[string, int]{"c": 3, "d": 4}

	r1 := partitionEvens(src1)
	r2 := partitionEvens(src2)

	assert.Equal(t, Record[string, int]{"a": 1}, P.Head(r1))
	assert.Equal(t, Record[string, int]{"b": 2}, P.Tail(r1))

	assert.Equal(t, Record[string, int]{"c": 3}, P.Head(r2))
	assert.Equal(t, Record[string, int]{"d": 4}, P.Tail(r2))
}

// TestPartition_NilRecord verifies that the curried form handles a nil map.
func TestPartition_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := Partition[string](isEven)(src)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestPartitionWithIndex_CurriedReturnsFunction verifies that PartitionWithIndex
// returns a reusable function applied to different records.
func TestPartitionWithIndex_CurriedReturnsFunction(t *testing.T) {
	partitionByPrefix := PartitionWithIndex(func(k string, _ int) bool {
		return strings.HasPrefix(k, "yes")
	})

	src := Record[string, int]{
		"yes_a": 1,
		"no_b":  2,
		"yes_c": 3,
	}

	result := partitionByPrefix(src)

	assert.Equal(t, Record[string, int]{"no_b": 2}, P.Head(result))
	assert.Equal(t, Record[string, int]{"yes_a": 1, "yes_c": 3}, P.Tail(result))
}

// TestPartitionWithIndex_NilRecord verifies that the curried form handles a
// nil map.
func TestPartitionWithIndex_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := PartitionWithIndex(func(_ string, _ int) bool { return true })(src)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// ExampleMonadPartition demonstrates splitting a record of scores into
// passing (>= 60) and failing (< 60) entries.
func ExampleMonadPartition() {
	scores := Record[string, int]{
		"alice": 80,
		"bob":   55,
		"carol": 70,
		"dave":  45,
	}

	result := MonadPartition(scores, N.MoreThan(59))

	failing := P.Head(result)
	passing := P.Tail(result)

	fmt.Println("alice passing:", passing["alice"])
	fmt.Println("carol passing:", passing["carol"])
	fmt.Println("bob failing:", failing["bob"])
	fmt.Println("dave failing:", failing["dave"])

	// Output:
	// alice passing: 80
	// carol passing: 70
	// bob failing: 55
	// dave failing: 45
}

// ExampleMonadPartitionWithIndex demonstrates splitting a record using both
// key prefix and value parity.
func ExampleMonadPartitionWithIndex() {
	data := Record[string, int]{
		"keep_2": 2,
		"keep_4": 4,
		"drop_1": 1,
		"drop_3": 3,
	}

	result := MonadPartitionWithIndex(data, func(k string, _ int) bool {
		return strings.HasPrefix(k, "keep")
	})

	fmt.Println("kept keep_2:", P.Tail(result)["keep_2"])
	fmt.Println("kept keep_4:", P.Tail(result)["keep_4"])
	fmt.Println("dropped drop_1:", P.Head(result)["drop_1"])
	fmt.Println("dropped drop_3:", P.Head(result)["drop_3"])

	// Output:
	// kept keep_2: 2
	// kept keep_4: 4
	// dropped drop_1: 1
	// dropped drop_3: 3
}

// ExamplePartition demonstrates applying the same curried partition function
// to two separate monthly budget records to split expenses into over-budget
// and within-budget categories.
func ExamplePartition() {
	// Budget limit per category is 100.
	splitBudget := Partition[string](N.MoreThan(100))

	january := Record[string, int]{
		"rent":     900,
		"groceries": 80,
	}
	february := Record[string, int]{
		"rent":     900,
		"groceries": 110,
	}

	janResult := splitBudget(january)
	febResult := splitBudget(february)

	// January: only rent exceeds the budget.
	fmt.Println("jan over-budget:    rent =", P.Tail(janResult)["rent"])
	fmt.Println("jan within-budget:  groceries =", P.Head(janResult)["groceries"])

	// February: rent and groceries both exceeded the budget.
	fmt.Println("feb over-budget:    rent =", P.Tail(febResult)["rent"])
	fmt.Println("feb over-budget:    groceries =", P.Tail(febResult)["groceries"])

	// Output:
	// jan over-budget:    rent = 900
	// jan within-budget:  groceries = 80
	// feb over-budget:    rent = 900
	// feb over-budget:    groceries = 110
}

// ExamplePartitionWithIndex demonstrates applying the same curried
// index-aware partition function to two service health records, routing
// entries to degraded or healthy based on both the service name and its
// error-rate value.
func ExamplePartitionWithIndex() {
	// A service is considered degraded when its error rate exceeds 5,
	// unless it is a "batch" service which tolerates up to 10.
	isDegraded := PartitionWithIndex(func(name string, errorRate int) bool {
		if strings.HasPrefix(name, "batch") {
			return errorRate > 10
		}
		return errorRate > 5
	})

	morning := Record[string, int]{
		"api-gateway": 3,
		"batch-jobs":  8,
	}
	evening := Record[string, int]{
		"api-gateway": 7,
		"batch-jobs":  12,
	}

	morningResult := isDegraded(morning)
	eveningResult := isDegraded(evening)

	// Morning: all services within tolerance.
	fmt.Println("morning healthy:   api-gateway =", P.Head(morningResult)["api-gateway"])
	fmt.Println("morning healthy:   batch-jobs =", P.Head(morningResult)["batch-jobs"])

	// Evening: both services are degraded.
	fmt.Println("evening degraded:  api-gateway =", P.Tail(eveningResult)["api-gateway"])
	fmt.Println("evening degraded:  batch-jobs =", P.Tail(eveningResult)["batch-jobs"])

	// Output:
	// morning healthy:   api-gateway = 3
	// morning healthy:   batch-jobs = 8
	// evening degraded:  api-gateway = 7
	// evening degraded:  batch-jobs = 12
}

// --- MonadPartitionMap ---

// classifyInt routes even numbers to Right and odd numbers to Left (as strings).
func classifyInt(v int) E.Either[string, int] {
	if v%2 == 0 {
		return E.Right[string](v)
	}
	return E.Left[int](strconv.Itoa(v))
}

// TestMonadPartitionMap_SplitsOnEither verifies that Right results go to Tail
// and Left results go to Head, with values transformed.
func TestMonadPartitionMap_SplitsOnEither(t *testing.T) {
	src := Record[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := MonadPartitionMap[string, int, string, int](src, classifyInt)

	assert.Equal(t, Record[string, string]{"a": "1", "c": "3"}, P.Head(result))
	assert.Equal(t, Record[string, int]{"b": 2, "d": 4}, P.Tail(result))
}

// TestMonadPartitionMap_AllRight verifies that Head is empty when every entry
// produces a Right.
func TestMonadPartitionMap_AllRight(t *testing.T) {
	src := Record[string, int]{"a": 2, "b": 4}

	result := MonadPartitionMap[string, int, string, int](src, classifyInt)

	assert.Empty(t, P.Head(result))
	assert.Equal(t, Record[string, int]{"a": 2, "b": 4}, P.Tail(result))
}

// TestMonadPartitionMap_AllLeft verifies that Tail is empty when every entry
// produces a Left.
func TestMonadPartitionMap_AllLeft(t *testing.T) {
	src := Record[string, int]{"a": 1, "b": 3}

	result := MonadPartitionMap[string, int, string, int](src, classifyInt)

	assert.Equal(t, Record[string, string]{"a": "1", "b": "3"}, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMap_EmptyRecord verifies that an empty source produces two
// empty records.
func TestMonadPartitionMap_EmptyRecord(t *testing.T) {
	src := Record[string, int]{}

	result := MonadPartitionMap[string, int, string, int](src, classifyInt)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMap_NilRecord verifies that a nil map is handled safely.
func TestMonadPartitionMap_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := MonadPartitionMap[string, int, string, int](src, classifyInt)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMap_TypeTransformation verifies that Left and Right can
// carry different types (heterogeneous output).
func TestMonadPartitionMap_TypeTransformation(t *testing.T) {
	// Classify strings: numeric → Right(parsed int), non-numeric → Left(original string)
	classify := func(s string) E.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return E.Left[int](s)
		}
		return E.Right[string](n)
	}

	src := Record[string, string]{
		"x": "42",
		"y": "hello",
		"z": "7",
	}

	result := MonadPartitionMap[string, string, string, int](src, classify)

	assert.Equal(t, Record[string, string]{"y": "hello"}, P.Head(result))
	assert.Equal(t, Record[string, int]{"x": 42, "z": 7}, P.Tail(result))
}

// --- PartitionMap ---

// TestPartitionMap_CurriedReturnsFunction verifies that PartitionMap produces a
// reusable function that can be applied to multiple records.
func TestPartitionMap_CurriedReturnsFunction(t *testing.T) {
	splitEvens := PartitionMap[string, int, string, int](classifyInt)

	src1 := Record[string, int]{"a": 1, "b": 2}
	src2 := Record[string, int]{"c": 3, "d": 4}

	r1 := splitEvens(src1)
	r2 := splitEvens(src2)

	assert.Equal(t, Record[string, string]{"a": "1"}, P.Head(r1))
	assert.Equal(t, Record[string, int]{"b": 2}, P.Tail(r1))

	assert.Equal(t, Record[string, string]{"c": "3"}, P.Head(r2))
	assert.Equal(t, Record[string, int]{"d": 4}, P.Tail(r2))
}

// TestPartitionMap_NilRecord verifies that the curried form handles a nil map.
func TestPartitionMap_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := PartitionMap[string, int, string, int](classifyInt)(src)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// --- Examples ---

// ExampleMonadPartitionMap demonstrates splitting a record of raw strings into
// parse errors (Left/Head) and successfully parsed integers (Right/Tail).
func ExampleMonadPartitionMap() {
	// Parse each value as an integer; Left = parse error, Right = parsed int.
	parseOrFail := func(s string) E.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return E.Left[int]("bad: " + s)
		}
		return E.Right[string](n)
	}

	data := Record[string, string]{
		"port":    "8080",
		"timeout": "thirty",
		"retries": "3",
	}

	result := MonadPartitionMap[string, string, string, int](data, parseOrFail)

	fmt.Println("parsed port:", P.Tail(result)["port"])
	fmt.Println("parsed retries:", P.Tail(result)["retries"])
	fmt.Println("invalid timeout:", P.Head(result)["timeout"])

	// Output:
	// parsed port: 8080
	// parsed retries: 3
	// invalid timeout: bad: thirty
}

// ExamplePartitionMap demonstrates applying the same curried PartitionMap
// function to two separate configuration records to split valid numeric values
// from invalid ones.
func ExamplePartitionMap() {
	// Reusable classifier: numeric string → Right(int), non-numeric → Left(string).
	splitNumeric := PartitionMap[string, string, string, int](func(s string) E.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return E.Left[int](s)
		}
		return E.Right[string](n)
	})

	serviceA := Record[string, string]{
		"workers": "4",
		"mode":    "fast",
	}
	serviceB := Record[string, string]{
		"workers": "eight",
		"timeout": "30",
	}

	rA := splitNumeric(serviceA)
	rB := splitNumeric(serviceB)

	fmt.Println("A valid workers:", P.Tail(rA)["workers"])
	fmt.Println("A invalid mode:", P.Head(rA)["mode"])
	fmt.Println("B invalid workers:", P.Head(rB)["workers"])
	fmt.Println("B valid timeout:", P.Tail(rB)["timeout"])

	// Output:
	// A valid workers: 4
	// A invalid mode: fast
	// B invalid workers: eight
	// B valid timeout: 30
}

// --- MonadPartitionMapWithIndex ---

// classifyWithKey routes even values to Right and odd values to Left, tagging
// the Left message with the key for traceability.
func classifyWithKey(k string, v int) E.Either[string, int] {
	if v%2 == 0 {
		return E.Right[string](v)
	}
	return E.Left[int](k + "=" + strconv.Itoa(v))
}

// TestMonadPartitionMapWithIndex_SplitsOnKeyAndValue verifies that the
// predicate receives both key and value, Right values go to Tail and Left
// values (tagged with the key) go to Head.
func TestMonadPartitionMapWithIndex_SplitsOnKeyAndValue(t *testing.T) {
	src := Record[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	result := MonadPartitionMapWithIndex[string, int, string, int](src, classifyWithKey)

	assert.Equal(t, Record[string, string]{"a": "a=1", "c": "c=3"}, P.Head(result))
	assert.Equal(t, Record[string, int]{"b": 2, "d": 4}, P.Tail(result))
}

// TestMonadPartitionMapWithIndex_AllRight verifies that Head is empty when
// every entry produces a Right.
func TestMonadPartitionMapWithIndex_AllRight(t *testing.T) {
	src := Record[string, int]{"a": 2, "b": 4}

	result := MonadPartitionMapWithIndex[string, int, string, int](src, classifyWithKey)

	assert.Empty(t, P.Head(result))
	assert.Equal(t, Record[string, int]{"a": 2, "b": 4}, P.Tail(result))
}

// TestMonadPartitionMapWithIndex_AllLeft verifies that Tail is empty when
// every entry produces a Left.
func TestMonadPartitionMapWithIndex_AllLeft(t *testing.T) {
	src := Record[string, int]{"a": 1, "b": 3}

	result := MonadPartitionMapWithIndex[string, int, string, int](src, classifyWithKey)

	assert.Equal(t, Record[string, string]{"a": "a=1", "b": "b=3"}, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMapWithIndex_EmptyRecord verifies that an empty source
// produces two empty records.
func TestMonadPartitionMapWithIndex_EmptyRecord(t *testing.T) {
	src := Record[string, int]{}

	result := MonadPartitionMapWithIndex[string, int, string, int](src, classifyWithKey)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMapWithIndex_NilRecord verifies that a nil map is handled
// safely.
func TestMonadPartitionMapWithIndex_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := MonadPartitionMapWithIndex[string, int, string, int](src, classifyWithKey)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// TestMonadPartitionMapWithIndex_KeyUsedInClassification verifies that the key
// itself drives the classification decision.
func TestMonadPartitionMapWithIndex_KeyUsedInClassification(t *testing.T) {
	// Accept when key starts with "ok", reject otherwise; Left carries the
	// original value, Right carries the uppercased value.
	classify := func(k string, v string) E.Either[string, string] {
		if strings.HasPrefix(k, "ok") {
			return E.Right[string](strings.ToUpper(v))
		}
		return E.Left[string](v)
	}

	src := Record[string, string]{
		"ok_a":  "hello",
		"bad_b": "world",
		"ok_c":  "foo",
	}

	result := MonadPartitionMapWithIndex[string, string, string, string](src, classify)

	assert.Equal(t, Record[string, string]{"bad_b": "world"}, P.Head(result))
	assert.Equal(t, Record[string, string]{"ok_a": "HELLO", "ok_c": "FOO"}, P.Tail(result))
}

// --- PartitionMapWithIndex ---

// TestPartitionMapWithIndex_CurriedReturnsFunction verifies that
// PartitionMapWithIndex produces a reusable function that can be applied to
// multiple records.
func TestPartitionMapWithIndex_CurriedReturnsFunction(t *testing.T) {
	splitEvens := PartitionMapWithIndex[string, int, string, int](classifyWithKey)

	src1 := Record[string, int]{"a": 1, "b": 2}
	src2 := Record[string, int]{"c": 3, "d": 4}

	r1 := splitEvens(src1)
	r2 := splitEvens(src2)

	assert.Equal(t, Record[string, string]{"a": "a=1"}, P.Head(r1))
	assert.Equal(t, Record[string, int]{"b": 2}, P.Tail(r1))

	assert.Equal(t, Record[string, string]{"c": "c=3"}, P.Head(r2))
	assert.Equal(t, Record[string, int]{"d": 4}, P.Tail(r2))
}

// TestPartitionMapWithIndex_NilRecord verifies that the curried form handles a
// nil map safely.
func TestPartitionMapWithIndex_NilRecord(t *testing.T) {
	var src Record[string, int]

	result := PartitionMapWithIndex[string, int, string, int](classifyWithKey)(src)

	assert.Empty(t, P.Head(result))
	assert.Empty(t, P.Tail(result))
}

// --- Examples ---

// ExampleMonadPartitionMapWithIndex demonstrates splitting a record of
// configuration values into validation errors (Left/Head) and accepted values
// (Right/Tail), where the key name determines the acceptable range.
func ExampleMonadPartitionMapWithIndex() {
	// Ports must be in 1024–65535; the key is included in the error message.
	classifyPort := func(name string, port int) E.Either[string, int] {
		if port < 1024 || port > 65535 {
			return E.Left[int](fmt.Sprintf("%s: port %d out of range", name, port))
		}
		return E.Right[string](port)
	}

	services := Record[string, int]{
		"api":      8080,
		"admin":    80,
		"metrics":  9090,
		"internal": 70000,
	}

	result := MonadPartitionMapWithIndex[string, int, string, int](services, classifyPort)

	fmt.Println("valid api:", P.Tail(result)["api"])
	fmt.Println("valid metrics:", P.Tail(result)["metrics"])
	fmt.Println("invalid admin:", P.Head(result)["admin"])
	fmt.Println("invalid internal:", P.Head(result)["internal"])

	// Output:
	// valid api: 8080
	// valid metrics: 9090
	// invalid admin: admin: port 80 out of range
	// invalid internal: internal: port 70000 out of range
}

// ExamplePartitionMapWithIndex demonstrates applying the same curried
// PartitionMapWithIndex function to two configuration records to split entries
// into validation errors and valid values, using the key to build the error.
func ExamplePartitionMapWithIndex() {
	// Reusable classifier: non-empty string → Right(string), empty → Left with key.
	requireNonEmpty := PartitionMapWithIndex[string, string, string, string](
		func(k string, v string) E.Either[string, string] {
			if strings.TrimSpace(v) == "" {
				return E.Left[string](k + " is required")
			}
			return E.Right[string](v)
		},
	)

	configA := Record[string, string]{
		"host": "localhost",
		"port": "",
	}
	configB := Record[string, string]{
		"host": "",
		"port": "8080",
	}

	rA := requireNonEmpty(configA)
	rB := requireNonEmpty(configB)

	fmt.Println("A valid host:", P.Tail(rA)["host"])
	fmt.Println("A invalid port:", P.Head(rA)["port"])
	fmt.Println("B invalid host:", P.Head(rB)["host"])
	fmt.Println("B valid port:", P.Tail(rB)["port"])

	// Output:
	// A valid host: localhost
	// A invalid port: port is required
	// B invalid host: host is required
	// B valid port: 8080
}
