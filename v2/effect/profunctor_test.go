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

package effect

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// Test types for profunctor tests
type AppConfig struct {
	DatabaseURL string
	APIKey      string
	Port        int
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Host string
	Port int
}

// TestPromapBasic tests basic Promap functionality
func TestPromapBasic(t *testing.T) {
	t.Run("transform both context and output", func(t *testing.T) {
		// Effect that uses DBConfig and returns an int
		getUserCount := Succeed[DBConfig](42)

		// Transform AppConfig to DBConfig
		extractDBConfig := func(app AppConfig) DBConfig {
			return DBConfig{URL: app.DatabaseURL}
		}

		// Transform int to string
		formatCount := func(count int) string {
			return fmt.Sprintf("Users: %d", count)
		}

		// Adapt the effect to work with AppConfig and return string
		adapted := Promap(extractDBConfig, formatCount)(getUserCount)
		result := adapted(AppConfig{
			DatabaseURL: "localhost:5432",
			APIKey:      "secret",
			Port:        8080,
		})(context.Background())()

		assert.Equal(t, R.Of("Users: 42"), result)
	})

	t.Run("identity transformations", func(t *testing.T) {
		// Effect that returns a value
		getValue := Succeed[DBConfig](100)

		// Identity transformations
		identity := func(x DBConfig) DBConfig { return x }
		identityInt := func(x int) int { return x }

		// Apply identity transformations
		adapted := Promap(identity, identityInt)(getValue)
		result := adapted(DBConfig{URL: "localhost"})(context.Background())()

		assert.Equal(t, R.Of(100), result)
	})
}

// TestPromapComposition tests that Promap composes correctly
func TestPromapComposition(t *testing.T) {
	t.Run("compose multiple transformations", func(t *testing.T) {
		// Effect that uses ServerConfig and returns the port
		getPort := Map[ServerConfig](func(cfg ServerConfig) int {
			return cfg.Port
		})(Ask[ServerConfig]())

		// First transformation: AppConfig -> ServerConfig
		extractServerConfig := func(app AppConfig) ServerConfig {
			return ServerConfig{Host: "localhost", Port: app.Port}
		}

		// Second transformation: int -> string
		formatPort := func(port int) string {
			return fmt.Sprintf(":%d", port)
		}

		// Apply transformations
		adapted := Promap(extractServerConfig, formatPort)(getPort)
		result := adapted(AppConfig{
			DatabaseURL: "db.example.com",
			APIKey:      "key123",
			Port:        9000,
		})(context.Background())()

		assert.Equal(t, R.Of(":9000"), result)
	})
}

// TestPromapWithErrors tests Promap with effects that can fail
func TestPromapWithErrors(t *testing.T) {
	t.Run("propagates errors correctly", func(t *testing.T) {
		// Effect that fails
		failingEffect := Fail[DBConfig, int](fmt.Errorf("database connection failed"))

		// Transformations
		extractDBConfig := func(app AppConfig) DBConfig {
			return DBConfig{URL: app.DatabaseURL}
		}
		formatCount := func(count int) string {
			return fmt.Sprintf("Count: %d", count)
		}

		// Apply transformations
		adapted := Promap(extractDBConfig, formatCount)(failingEffect)
		result := adapted(AppConfig{DatabaseURL: "localhost"})(context.Background())()

		assert.True(t, R.IsLeft(result))
		err := R.MonadFold(result,
			func(e error) error { return e },
			func(string) error { return nil },
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("output transformation not applied on error", func(t *testing.T) {
		callCount := 0

		// Effect that fails
		failingEffect := Fail[DBConfig, int](fmt.Errorf("error"))

		// Transformation that counts calls
		countingTransform := func(x int) string {
			callCount++
			return strconv.Itoa(x)
		}

		// Apply transformations
		adapted := Promap(
			func(app AppConfig) DBConfig { return DBConfig{URL: app.DatabaseURL} },
			countingTransform,
		)(failingEffect)
		result := adapted(AppConfig{DatabaseURL: "localhost"})(context.Background())()

		assert.True(t, R.IsLeft(result))
		assert.Equal(t, 0, callCount, "output transformation should not be called on error")
	})
}

// TestPromapWithComplexTypes tests Promap with more complex type transformations
func TestPromapWithComplexTypes(t *testing.T) {
	t.Run("transform struct to different struct", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		type UserDTO struct {
			UserID   int
			FullName string
		}

		// Effect that uses User and returns a string
		getUserInfo := Map[User](func(user User) string {
			return fmt.Sprintf("User %s (ID: %d)", user.Name, user.ID)
		})(Ask[User]())

		// Transform UserDTO to User
		dtoToUser := func(dto UserDTO) User {
			return User{ID: dto.UserID, Name: dto.FullName}
		}

		// Transform string to uppercase
		toUpper := func(s string) string {
			return fmt.Sprintf("INFO: %s", s)
		}

		// Apply transformations
		adapted := Promap(dtoToUser, toUpper)(getUserInfo)
		result := adapted(UserDTO{UserID: 123, FullName: "Alice"})(context.Background())()

		assert.Equal(t, R.Of("INFO: User Alice (ID: 123)"), result)
	})
}

// TestPromapChaining tests chaining multiple Promap operations
func TestPromapChaining(t *testing.T) {
	t.Run("chain multiple Promap operations", func(t *testing.T) {
		// Base effect that doubles the input
		baseEffect := Map[int](func(x int) int {
			return x * 2
		})(Ask[int]())

		// First Promap: string -> int, int -> string
		step1 := Promap(
			func(s string) int {
				n, _ := strconv.Atoi(s)
				return n
			},
			strconv.Itoa,
		)(baseEffect)

		// Second Promap: float64 -> string, string -> float64
		step2 := Promap(
			func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			},
			func(s string) float64 {
				f, _ := strconv.ParseFloat(s, 64)
				return f
			},
		)(step1)

		result := step2(21.0)(context.Background())()

		assert.Equal(t, R.Of(42.0), result)
	})
}

// TestPromapEdgeCases tests edge cases
func TestPromapEdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		effect := Map[int](func(x int) int {
			return x
		})(Ask[int]())

		adapted := Promap(
			func(s string) int { return 0 },
			func(x int) string { return "" },
		)(effect)

		result := adapted("anything")(context.Background())()

		assert.Equal(t, R.Of(""), result)
	})

	t.Run("nil context handling", func(t *testing.T) {
		effect := Succeed[int]("success")

		adapted := Promap(
			func(s string) int { return 42 },
			func(s string) string { return s + "!" },
		)(effect)

		// Using background context instead of nil
		result := adapted("test")(context.Background())()

		assert.Equal(t, R.Of("success!"), result)
	})
}

// TestPromapIntegration tests integration with other effect operations
func TestPromapIntegration(t *testing.T) {
	t.Run("Promap with Map", func(t *testing.T) {
		// Base effect that adds 10
		baseEffect := Map[int](func(x int) int {
			return x + 10
		})(Ask[int]())

		// Apply Promap
		promapped := Promap(
			func(s string) int {
				n, _ := strconv.Atoi(s)
				return n
			},
			func(x int) int { return x * 2 },
		)(baseEffect)

		// Apply Map on top
		mapped := Map[string](func(x int) string {
			return fmt.Sprintf("Result: %d", x)
		})(promapped)

		result := mapped("5")(context.Background())()

		assert.Equal(t, R.Of("Result: 30"), result)
	})

	t.Run("Promap with Chain", func(t *testing.T) {
		// Base effect
		baseEffect := Ask[int]()

		// Apply Promap
		promapped := Promap(
			func(s string) int {
				n, _ := strconv.Atoi(s)
				return n
			},
			func(x int) int { return x * 2 },
		)(baseEffect)

		// Chain with another effect
		chained := Chain(func(x int) Effect[string, string] {
			return Succeed[string](fmt.Sprintf("Value: %d", x))
		})(promapped)

		result := chained("10")(context.Background())()

		assert.Equal(t, R.Of("Value: 20"), result)
	})
}

// BenchmarkPromap benchmarks the Promap operation
func BenchmarkPromap(b *testing.B) {
	effect := Map[int](func(x int) int {
		return x * 2
	})(Ask[int]())

	adapted := Promap(
		func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		},
		strconv.Itoa,
	)(effect)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = adapted("42")(ctx)()
	}
}

// BenchmarkPromapChained benchmarks chained Promap operations
func BenchmarkPromapChained(b *testing.B) {
	baseEffect := Map[int](func(x int) int {
		return x * 2
	})(Ask[int]())

	step1 := Promap(
		func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		},
		strconv.Itoa,
	)(baseEffect)

	step2 := Promap(
		func(f float64) string {
			return fmt.Sprintf("%.0f", f)
		},
		func(s string) float64 {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		},
	)(step1)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = step2(21.0)(ctx)()
	}
}
