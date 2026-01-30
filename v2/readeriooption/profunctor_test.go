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

package readeriooption

import (
	"context"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestPromap_TransformBoth(t *testing.T) {
	// Test transforming both input environment and output value
	type GlobalConfig struct {
		Factor int
	}

	type LocalConfig struct {
		Multiplier int
	}

	// Original computation expects LocalConfig and returns int
	original := func(cfg LocalConfig) IOOption[int] {
		return func() Option[int] {
			return O.Of(10 * cfg.Multiplier)
		}
	}

	// Transform GlobalConfig to LocalConfig (contravariant)
	envTransform := func(g GlobalConfig) LocalConfig {
		return LocalConfig{Multiplier: g.Factor}
	}

	// Transform int to string (covariant)
	valueTransform := func(n int) string {
		return fmt.Sprintf("%d", n)
	}

	// Apply Promap
	adapted := F.Pipe1(
		original,
		Promap(envTransform, valueTransform),
	)

	globalCfg := GlobalConfig{Factor: 5}
	result := adapted(globalCfg)()

	expected := O.Of("50")
	assert.Equal(t, expected, result)
}

func TestPromap_WithNone(t *testing.T) {
	// Test that None is preserved through Promap
	type Config1 struct {
		Value int
	}

	type Config2 struct {
		Data int
	}

	original := None[Config1, int]()

	envTransform := func(c2 Config2) Config1 {
		return Config1{Value: c2.Data}
	}

	valueTransform := func(n int) string {
		return fmt.Sprintf("%d", n)
	}

	adapted := F.Pipe1(
		original,
		Promap(envTransform, valueTransform),
	)

	cfg := Config2{Data: 10}
	result := adapted(cfg)()

	expected := O.None[string]()
	assert.Equal(t, expected, result)
}

func TestPromap_Identity(t *testing.T) {
	// Test that Promap with identity functions is identity
	original := Of[context.Context](42)

	adapted := F.Pipe1(
		original,
		Promap(
			F.Identity[context.Context],
			F.Identity[int],
		),
	)

	result := adapted(context.Background())()
	expected := O.Of(42)

	assert.Equal(t, expected, result)
}

func TestPromap_Composition(t *testing.T) {
	// Test that Promap composes correctly
	type Config1 struct{ A int }
	type Config2 struct{ B int }
	type Config3 struct{ C int }

	original := func(c1 Config1) IOOption[int] {
		return func() Option[int] {
			return O.Of(c1.A * 2)
		}
	}

	// First transformation
	f1 := func(c2 Config2) Config1 { return Config1{A: c2.B + 1} }
	g1 := func(n int) int { return n * 3 }

	// Second transformation
	f2 := func(c3 Config3) Config2 { return Config2{B: c3.C + 2} }
	g2 := func(n int) string { return fmt.Sprintf("%d", n) }

	// Apply transformations separately
	step1 := F.Pipe1(original, Promap(f1, g1))
	step2 := F.Pipe1(step1, Promap(f2, g2))

	// Apply composed transformation
	composed := F.Pipe1(
		original,
		Promap(
			F.Flow2(f2, f1),
			F.Flow2(g1, g2),
		),
	)

	cfg := Config3{C: 5}

	result1 := step2(cfg)()
	result2 := composed(cfg)()

	// Both should give the same result: ((5+2+1)*2)*3 = 48
	expected := O.Of("48")
	assert.Equal(t, expected, result1)
	assert.Equal(t, expected, result2)
}

func TestContramap_TransformEnvironment(t *testing.T) {
	// Test transforming only the environment
	type GlobalConfig struct {
		DatabaseURL string
		Port        int
	}

	type DBConfig struct {
		URL string
	}

	// Original computation expects DBConfig
	original := func(cfg DBConfig) IOOption[string] {
		return func() Option[string] {
			return O.Of("Connected to: " + cfg.URL)
		}
	}

	// Transform GlobalConfig to DBConfig
	envTransform := func(g GlobalConfig) DBConfig {
		return DBConfig{URL: g.DatabaseURL}
	}

	// Apply Contramap
	adapted := F.Pipe1(
		original,
		Contramap[string](envTransform),
	)

	globalCfg := GlobalConfig{
		DatabaseURL: "localhost:5432",
		Port:        8080,
	}
	result := adapted(globalCfg)()

	expected := O.Of("Connected to: localhost:5432")
	assert.Equal(t, expected, result)
}

func TestContramap_WithNone(t *testing.T) {
	// Test that None is preserved through Contramap
	type Config1 struct {
		Value int
	}

	type Config2 struct {
		Data int
	}

	original := None[Config1, string]()

	envTransform := func(c2 Config2) Config1 {
		return Config1{Value: c2.Data}
	}

	adapted := F.Pipe1(
		original,
		Contramap[string](envTransform),
	)

	cfg := Config2{Data: 10}
	result := adapted(cfg)()

	expected := O.None[string]()
	assert.Equal(t, expected, result)
}

func TestContramap_Identity(t *testing.T) {
	// Test that Contramap with identity function is identity
	original := Of[context.Context](42)

	adapted := F.Pipe1(
		original,
		Contramap[int](F.Identity[context.Context]),
	)

	result := adapted(context.Background())()
	expected := O.Of(42)

	assert.Equal(t, expected, result)
}

func TestContramap_Composition(t *testing.T) {
	// Test that Contramap composes correctly
	type Config1 struct{ A int }
	type Config2 struct{ B int }
	type Config3 struct{ C int }

	original := func(c1 Config1) IOOption[int] {
		return func() Option[int] {
			return O.Of(c1.A * 10)
		}
	}

	f1 := func(c2 Config2) Config1 { return Config1{A: c2.B + 1} }
	f2 := func(c3 Config3) Config2 { return Config2{B: c3.C + 2} }

	// Apply transformations separately
	step1 := F.Pipe1(original, Contramap[int](f1))
	step2 := F.Pipe1(step1, Contramap[int](f2))

	// Apply composed transformation
	composed := F.Pipe1(
		original,
		Contramap[int](F.Flow2(f2, f1)),
	)

	cfg := Config3{C: 5}

	result1 := step2(cfg)()
	result2 := composed(cfg)()

	// Both should give the same result: (5+2+1)*10 = 80
	expected := O.Of(80)
	assert.Equal(t, expected, result1)
	assert.Equal(t, expected, result2)
}

func TestPromap_RealWorldExample(t *testing.T) {
	// Real-world example: adapting a database query function
	type AppConfig struct {
		DBHost     string
		DBPort     int
		DBUser     string
		DBPassword string
		LogLevel   string
	}

	type DBConnection struct {
		ConnectionString string
	}

	type User struct {
		ID   int
		Name string
	}

	type UserDTO struct {
		UserID      int
		DisplayName string
	}

	// Original function that queries database
	queryUser := func(conn DBConnection) IOOption[User] {
		return func() Option[User] {
			// Simulate database query
			if conn.ConnectionString != "" {
				return O.Of(User{ID: 1, Name: "Alice"})
			}
			return O.None[User]()
		}
	}

	// Adapt to work with AppConfig and return UserDTO
	adaptedQuery := F.Pipe1(
		queryUser,
		Promap(
			// Extract DB connection from app config
			func(cfg AppConfig) DBConnection {
				return DBConnection{
					ConnectionString: cfg.DBUser + "@" + cfg.DBHost,
				}
			},
			// Convert User to UserDTO
			func(u User) UserDTO {
				return UserDTO{
					UserID:      u.ID,
					DisplayName: "User: " + u.Name,
				}
			},
		),
	)

	appCfg := AppConfig{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "admin",
		DBPassword: "secret",
		LogLevel:   "info",
	}

	result := adaptedQuery(appCfg)()
	expected := O.Of(UserDTO{UserID: 1, DisplayName: "User: Alice"})

	assert.Equal(t, expected, result)
}

func TestContramap_RealWorldExample(t *testing.T) {
	// Real-world example: adapting a service that needs specific config
	type GlobalConfig struct {
		ServiceURL string
		APIKey     string
		Timeout    int
		RetryCount int
	}

	type ServiceConfig struct {
		Endpoint string
		Auth     string
	}

	// Service function that needs ServiceConfig
	callService := func(cfg ServiceConfig) IOOption[string] {
		return func() Option[string] {
			if cfg.Endpoint != "" && cfg.Auth != "" {
				return O.Of("Response from " + cfg.Endpoint)
			}
			return O.None[string]()
		}
	}

	// Adapt to work with GlobalConfig
	adaptedService := F.Pipe1(
		callService,
		Contramap[string](func(g GlobalConfig) ServiceConfig {
			return ServiceConfig{
				Endpoint: g.ServiceURL,
				Auth:     "Bearer " + g.APIKey,
			}
		}),
	)

	globalCfg := GlobalConfig{
		ServiceURL: "https://api.example.com",
		APIKey:     "secret-key",
		Timeout:    30,
		RetryCount: 3,
	}

	result := adaptedService(globalCfg)()
	expected := O.Of("Response from https://api.example.com")

	assert.Equal(t, expected, result)
}
