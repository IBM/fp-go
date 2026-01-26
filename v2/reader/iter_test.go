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

package reader

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	INTI "github.com/IBM/fp-go/v2/internal/iter"
	"github.com/stretchr/testify/assert"
)

// Helper function to collect iterator values into a slice
func collectIter[A any](seq Seq[A]) []A {
	return INTI.ToArray[Seq[A], []A](seq)
}

// Helper function to create an iterator from a slice
func fromSlice[A any](items []A) Seq[A] {
	return slices.Values(items)
}

func TestTraverseIter(t *testing.T) {
	type Config struct {
		Multiplier int
		Prefix     string
	}

	t.Run("traverses empty iterator", func(t *testing.T) {
		empty := INTI.Empty[Seq[int]]()

		multiplyByConfig := func(x int) Reader[Config, int] {
			return func(c Config) int { return x * c.Multiplier }
		}

		traversed := TraverseIter(multiplyByConfig)(empty)
		result := traversed(Config{Multiplier: 10})

		collected := collectIter(result)
		assert.Empty(t, collected)
	})

	t.Run("traverses single element iterator", func(t *testing.T) {
		single := INTI.Of[Seq[int]](5)

		multiplyByConfig := func(x int) Reader[Config, int] {
			return func(c Config) int { return x * c.Multiplier }
		}

		traversed := TraverseIter(multiplyByConfig)(single)
		result := traversed(Config{Multiplier: 3})

		collected := collectIter(result)
		assert.Equal(t, []int{15}, collected)
	})

	t.Run("traverses multiple elements", func(t *testing.T) {
		numbers := INTI.From(1, 2, 3, 4)

		multiplyByConfig := func(x int) Reader[Config, int] {
			return func(c Config) int { return x * c.Multiplier }
		}

		traversed := TraverseIter(multiplyByConfig)(numbers)
		result := traversed(Config{Multiplier: 10})

		collected := collectIter(result)
		assert.Equal(t, []int{10, 20, 30, 40}, collected)
	})

	t.Run("transforms types during traversal", func(t *testing.T) {
		numbers := INTI.From(1, 2, 3)

		intToString := func(x int) Reader[Config, string] {
			return func(c Config) string {
				return fmt.Sprintf("%s%d", c.Prefix, x)
			}
		}

		traversed := TraverseIter(intToString)(numbers)
		result := traversed(Config{Prefix: "num-"})

		collected := collectIter(result)
		assert.Equal(t, []string{"num-1", "num-2", "num-3"}, collected)
	})

	t.Run("all readers share same environment", func(t *testing.T) {
		numbers := INTI.From(1, 2, 3)

		// Each reader accesses the same config
		addBase := func(x int) Reader[Config, int] {
			return func(c Config) int {
				return x + c.Multiplier
			}
		}

		traversed := TraverseIter(addBase)(numbers)
		result := traversed(Config{Multiplier: 100})

		collected := collectIter(result)
		assert.Equal(t, []int{101, 102, 103}, collected)
	})

	t.Run("works with complex transformations", func(t *testing.T) {
		words := INTI.From("hello", "world")

		wordLength := func(s string) Reader[Config, int] {
			return func(c Config) int {
				return len(s) * c.Multiplier
			}
		}

		traversed := TraverseIter(wordLength)(words)
		result := traversed(Config{Multiplier: 2})

		collected := collectIter(result)
		assert.Equal(t, []int{10, 10}, collected) // "hello" = 5*2, "world" = 5*2
	})

	t.Run("preserves order of elements", func(t *testing.T) {
		numbers := fromSlice([]int{5, 3, 8, 1, 9})

		identity := func(x int) Reader[Config, int] {
			return Of[Config](x)
		}

		traversed := TraverseIter(identity)(numbers)
		result := traversed(Config{})

		collected := collectIter(result)
		assert.Equal(t, []int{5, 3, 8, 1, 9}, collected)
	})

	t.Run("can be used with different config types", func(t *testing.T) {
		type StringConfig struct {
			Suffix string
		}

		words := INTI.From("test", "data")

		addSuffix := func(s string) Reader[StringConfig, string] {
			return func(c StringConfig) string {
				return s + c.Suffix
			}
		}

		traversed := TraverseIter(addSuffix)(words)
		result := traversed(StringConfig{Suffix: ".txt"})

		collected := collectIter(result)
		assert.Equal(t, []string{"test.txt", "data.txt"}, collected)
	})
}

func TestSequenceIter(t *testing.T) {
	type Config struct {
		Base       int
		Multiplier int
	}

	t.Run("sequences empty iterator", func(t *testing.T) {
		empty := func(yield func(Reader[Config, int]) bool) {}

		sequenced := SequenceIter(empty)
		result := sequenced(Config{Base: 10})

		collected := collectIter(result)
		assert.Empty(t, collected)
	})

	t.Run("sequences single reader", func(t *testing.T) {
		single := func(yield func(Reader[Config, int]) bool) {
			yield(func(c Config) int { return c.Base + 5 })
		}

		sequenced := SequenceIter(single)
		result := sequenced(Config{Base: 10})

		collected := collectIter(result)
		assert.Equal(t, []int{15}, collected)
	})

	t.Run("sequences multiple readers", func(t *testing.T) {
		readers := INTI.From(
			func(c Config) int { return c.Base + 1 },
			func(c Config) int { return c.Base + 2 },
			func(c Config) int { return c.Base + 3 },
		)

		sequenced := SequenceIter(readers)
		result := sequenced(Config{Base: 10})

		collected := collectIter(result)
		assert.Equal(t, []int{11, 12, 13}, collected)
	})

	t.Run("all readers receive same environment", func(t *testing.T) {
		readers := INTI.From(
			func(c Config) int { return c.Base * c.Multiplier },
			func(c Config) int { return c.Base + c.Multiplier },
			func(c Config) int { return c.Base - c.Multiplier },
		)

		sequenced := SequenceIter(readers)
		result := sequenced(Config{Base: 10, Multiplier: 3})

		collected := collectIter(result)
		assert.Equal(t, []int{30, 13, 7}, collected)
	})

	t.Run("works with string readers", func(t *testing.T) {
		type StringConfig struct {
			Prefix string
			Suffix string
		}

		readers := INTI.From(
			func(c StringConfig) string { return c.Prefix + "first" },
			func(c StringConfig) string { return c.Prefix + "second" },
			func(c StringConfig) string { return "third" + c.Suffix },
		)

		sequenced := SequenceIter(readers)
		result := sequenced(StringConfig{Prefix: "pre-", Suffix: "-post"})

		collected := collectIter(result)
		assert.Equal(t, []string{"pre-first", "pre-second", "third-post"}, collected)
	})

	t.Run("preserves order of readers", func(t *testing.T) {
		readers := INTI.From(
			Of[Config](5),
			Of[Config](3),
			Of[Config](8),
			Of[Config](1),
		)

		sequenced := SequenceIter(readers)
		result := sequenced(Config{})

		collected := collectIter(result)
		assert.Equal(t, []int{5, 3, 8, 1}, collected)
	})

	t.Run("works with complex reader logic", func(t *testing.T) {
		readers := INTI.From(
			func(c Config) string {
				return strconv.Itoa(c.Base * 2)
			},
			func(c Config) string {
				return fmt.Sprintf("mult-%d", c.Multiplier)
			},
			func(c Config) string {
				return fmt.Sprintf("sum-%d", c.Base+c.Multiplier)
			},
		)

		sequenced := SequenceIter(readers)
		result := sequenced(Config{Base: 5, Multiplier: 3})

		collected := collectIter(result)
		assert.Equal(t, []string{"10", "mult-3", "sum-8"}, collected)
	})

	t.Run("can handle large number of readers", func(t *testing.T) {

		readers := func(yield func(Reader[Config, int]) bool) {
			for i := 0; i < 100; i++ {
				i := i // capture loop variable
				yield(func(c Config) int { return c.Base + i })
			}
		}

		sequenced := SequenceIter(readers)
		result := sequenced(Config{Base: 1000})

		collected := collectIter(result)
		assert.Len(t, collected, 100)
		assert.Equal(t, 1000, collected[0])
		assert.Equal(t, 1099, collected[99])
	})
}

func TestTraverseIterAndSequenceIterRelationship(t *testing.T) {
	type Config struct {
		Value int
	}

	t.Run("SequenceIter is TraverseIter with identity", func(t *testing.T) {
		// Create an iterator of readers
		readers := INTI.From(
			func(c Config) int { return c.Value + 1 },
			func(c Config) int { return c.Value + 2 },
			func(c Config) int { return c.Value + 3 },
		)

		// Using SequenceIter
		sequenced := SequenceIter(readers)
		sequencedResult := sequenced(Config{Value: 10})

		// Using TraverseIter with identity function
		identity := Asks[Config, int]
		traversed := TraverseIter(identity)(readers)
		traversedResult := traversed(Config{Value: 10})

		// Both should produce the same results
		sequencedCollected := collectIter(sequencedResult)
		traversedCollected := collectIter(traversedResult)

		assert.Equal(t, sequencedCollected, traversedCollected)
		assert.Equal(t, []int{11, 12, 13}, sequencedCollected)
	})
}

func TestIteratorIntegration(t *testing.T) {
	type AppConfig struct {
		DatabaseURL string
		APIKey      string
		Port        int
	}

	t.Run("real-world example: processing configuration values", func(t *testing.T) {
		// Iterator of field names
		fields := INTI.From(
			"database",
			"api",
			"port",
		)

		// Function that creates a reader for each field
		getConfigValue := func(field string) Reader[AppConfig, string] {
			return func(c AppConfig) string {
				switch field {
				case "database":
					return c.DatabaseURL
				case "api":
					return c.APIKey
				case "port":
					return strconv.Itoa(c.Port)
				default:
					return "unknown"
				}
			}
		}

		// Traverse to get all config values
		traversed := TraverseIter(getConfigValue)(fields)
		result := traversed(AppConfig{
			DatabaseURL: "postgres://localhost",
			APIKey:      "secret-key",
			Port:        8080,
		})

		collected := collectIter(result)
		assert.Equal(t, []string{
			"postgres://localhost",
			"secret-key",
			"8080",
		}, collected)
	})

	t.Run("real-world example: batch processing with shared config", func(t *testing.T) {
		type ProcessConfig struct {
			Prefix string
			Suffix string
		}

		// Iterator of items to process
		items := fromSlice([]string{"item1", "item2", "item3"})

		// Processing function that uses config
		processItem := func(item string) Reader[ProcessConfig, string] {
			return func(c ProcessConfig) string {
				return c.Prefix + item + c.Suffix
			}
		}

		// Process all items with shared config
		traversed := TraverseIter(processItem)(items)
		result := traversed(ProcessConfig{
			Prefix: "[",
			Suffix: "]",
		})

		collected := collectIter(result)
		assert.Equal(t, []string{"[item1]", "[item2]", "[item3]"}, collected)
	})
}
