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

package consumer

import (
	"strconv"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestLocal(t *testing.T) {
	t.Run("basic type transformation", func(t *testing.T) {
		var captured int
		consumeInt := func(x int) {
			captured = x
		}

		// Transform string to int before consuming
		stringToInt := func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		}

		consumeString := Local(stringToInt)(consumeInt)
		consumeString("42")

		assert.Equal(t, 42, captured)
	})

	t.Run("field extraction from struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		var capturedName string
		consumeName := func(name string) {
			capturedName = name
		}

		extractName := func(u User) string {
			return u.Name
		}

		consumeUser := Local(extractName)(consumeName)
		consumeUser(User{Name: "Alice", Age: 30})

		assert.Equal(t, "Alice", capturedName)
	})

	t.Run("simplifying complex types", func(t *testing.T) {
		type DetailedConfig struct {
			Host     string
			Port     int
			Timeout  time.Duration
			MaxRetry int
		}

		type SimpleConfig struct {
			Host string
			Port int
		}

		var captured SimpleConfig
		consumeSimple := func(c SimpleConfig) {
			captured = c
		}

		simplify := func(d DetailedConfig) SimpleConfig {
			return SimpleConfig{Host: d.Host, Port: d.Port}
		}

		consumeDetailed := Local(simplify)(consumeSimple)
		consumeDetailed(DetailedConfig{
			Host:     "localhost",
			Port:     8080,
			Timeout:  time.Second,
			MaxRetry: 3,
		})

		assert.Equal(t, SimpleConfig{Host: "localhost", Port: 8080}, captured)
	})

	t.Run("multiple transformations", func(t *testing.T) {
		type Response struct {
			StatusCode int
			Body       string
		}

		var capturedStatus int
		consumeStatus := func(code int) {
			capturedStatus = code
		}

		getStatus := func(r Response) int {
			return r.StatusCode
		}

		consumeResponse := Local(getStatus)(consumeStatus)
		consumeResponse(Response{StatusCode: 200, Body: "OK"})

		assert.Equal(t, 200, capturedStatus)
	})

	t.Run("chaining Local transformations", func(t *testing.T) {
		type Level3 struct{ Value int }
		type Level2 struct{ L3 Level3 }
		type Level1 struct{ L2 Level2 }

		var captured int
		consumeInt := func(x int) {
			captured = x
		}

		// Chain multiple Local transformations
		extract3 := func(l3 Level3) int { return l3.Value }
		extract2 := func(l2 Level2) Level3 { return l2.L3 }
		extract1 := func(l1 Level1) Level2 { return l1.L2 }

		// Compose the transformations
		consumeLevel3 := Local(extract3)(consumeInt)
		consumeLevel2 := Local(extract2)(consumeLevel3)
		consumeLevel1 := Local(extract1)(consumeLevel2)

		consumeLevel1(Level1{L2: Level2{L3: Level3{Value: 42}}})

		assert.Equal(t, 42, captured)
	})

	t.Run("identity transformation", func(t *testing.T) {
		var captured string
		consumeString := func(s string) {
			captured = s
		}

		identity := function.Identity[string]

		consumeIdentity := Local(identity)(consumeString)
		consumeIdentity("test")

		assert.Equal(t, "test", captured)
	})

	t.Run("transformation with calculation", func(t *testing.T) {
		type Rectangle struct {
			Width  int
			Height int
		}

		var capturedArea int
		consumeArea := func(area int) {
			capturedArea = area
		}

		calculateArea := func(r Rectangle) int {
			return r.Width * r.Height
		}

		consumeRectangle := Local(calculateArea)(consumeArea)
		consumeRectangle(Rectangle{Width: 5, Height: 10})

		assert.Equal(t, 50, capturedArea)
	})

	t.Run("multiple consumers with same transformation", func(t *testing.T) {
		type Event struct {
			Type      string
			Timestamp time.Time
		}

		var capturedType string
		var capturedTime time.Time

		consumeType := func(t string) {
			capturedType = t
		}

		consumeTime := func(t time.Time) {
			capturedTime = t
		}

		extractType := func(e Event) string { return e.Type }
		extractTime := func(e Event) time.Time { return e.Timestamp }

		consumeEventType := Local(extractType)(consumeType)
		consumeEventTime := Local(extractTime)(consumeTime)

		now := time.Now()
		event := Event{Type: "UserLogin", Timestamp: now}

		consumeEventType(event)
		consumeEventTime(event)

		assert.Equal(t, "UserLogin", capturedType)
		assert.Equal(t, now, capturedTime)
	})

	t.Run("transformation with slice", func(t *testing.T) {
		var captured int
		consumeLength := func(n int) {
			captured = n
		}

		getLength := func(s []string) int {
			return len(s)
		}

		consumeSlice := Local(getLength)(consumeLength)
		consumeSlice([]string{"a", "b", "c"})

		assert.Equal(t, 3, captured)
	})

	t.Run("transformation with map", func(t *testing.T) {
		var captured int
		consumeCount := func(n int) {
			captured = n
		}

		getCount := func(m map[string]int) int {
			return len(m)
		}

		consumeMap := Local(getCount)(consumeCount)
		consumeMap(map[string]int{"a": 1, "b": 2, "c": 3})

		assert.Equal(t, 3, captured)
	})

	t.Run("transformation with pointer", func(t *testing.T) {
		var captured int
		consumeInt := func(x int) {
			captured = x
		}

		dereference := func(p *int) int {
			if p == nil {
				return 0
			}
			return *p
		}

		consumePointer := Local(dereference)(consumeInt)

		value := 42
		consumePointer(&value)
		assert.Equal(t, 42, captured)

		consumePointer(nil)
		assert.Equal(t, 0, captured)
	})

	t.Run("transformation with custom type", func(t *testing.T) {
		type MyType struct {
			Value string
		}

		var captured string
		consumeString := func(s string) {
			captured = s
		}

		extractValue := func(m MyType) string {
			return m.Value
		}

		consumeMyType := Local(extractValue)(consumeString)
		consumeMyType(MyType{Value: "test"})

		assert.Equal(t, "test", captured)
	})

	t.Run("accumulation through multiple calls", func(t *testing.T) {
		var sum int
		accumulate := func(x int) {
			sum += x
		}

		double := func(x int) int {
			return x * 2
		}

		accumulateDoubled := Local(double)(accumulate)

		accumulateDoubled(1)
		accumulateDoubled(2)
		accumulateDoubled(3)

		assert.Equal(t, 12, sum) // (1*2) + (2*2) + (3*2) = 2 + 4 + 6 = 12
	})

	t.Run("transformation with error handling", func(t *testing.T) {
		type Result struct {
			Value int
			Error error
		}

		var captured int
		consumeInt := func(x int) {
			captured = x
		}

		extractValue := func(r Result) int {
			if r.Error != nil {
				return -1
			}
			return r.Value
		}

		consumeResult := Local(extractValue)(consumeInt)

		consumeResult(Result{Value: 42, Error: nil})
		assert.Equal(t, 42, captured)

		consumeResult(Result{Value: 100, Error: assert.AnError})
		assert.Equal(t, -1, captured)
	})

	t.Run("transformation preserves consumer behavior", func(t *testing.T) {
		callCount := 0
		consumer := func(x int) {
			callCount++
		}

		transform := func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		}

		transformedConsumer := Local(transform)(consumer)

		transformedConsumer("1")
		transformedConsumer("2")
		transformedConsumer("3")

		assert.Equal(t, 3, callCount)
	})

	t.Run("comparison with reader.Local behavior", func(t *testing.T) {
		// This test demonstrates the dual nature of Consumer and Reader
		// Consumer: transforms input before consumption (contravariant)
		// Reader: transforms environment before reading (also contravariant on input)

		type DetailedEnv struct {
			Value int
			Extra string
		}

		type SimpleEnv struct {
			Value int
		}

		var captured int
		consumeSimple := func(e SimpleEnv) {
			captured = e.Value
		}

		simplify := func(d DetailedEnv) SimpleEnv {
			return SimpleEnv{Value: d.Value}
		}

		consumeDetailed := Local(simplify)(consumeSimple)
		consumeDetailed(DetailedEnv{Value: 42, Extra: "ignored"})

		assert.Equal(t, 42, captured)
	})
}
