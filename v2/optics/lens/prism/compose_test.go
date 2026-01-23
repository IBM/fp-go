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

package prism

import (
	"testing"

	"github.com/IBM/fp-go/v2/assert"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
)

// Test types for composition examples

// ConnectionType is a sum type representing different database connections
type ConnectionType interface {
	isConnection()
}

type PostgreSQL struct {
	Host string
	Port int
}

func (PostgreSQL) isConnection() {}

type MySQL struct {
	Host string
	Port int
}

func (MySQL) isConnection() {}

type MongoDB struct {
	Host string
	Port int
}

func (MongoDB) isConnection() {}

// Config is the top-level configuration
type Config struct {
	Connection ConnectionType
	AppName    string
}

// Helper functions to create prisms for each connection type

func postgresqlPrism() P.Prism[ConnectionType, PostgreSQL] {
	return P.MakePrism(
		func(ct ConnectionType) O.Option[PostgreSQL] {
			if pg, ok := ct.(PostgreSQL); ok {
				return O.Some(pg)
			}
			return O.None[PostgreSQL]()
		},
		func(pg PostgreSQL) ConnectionType { return pg },
	)
}

func mysqlPrism() P.Prism[ConnectionType, MySQL] {
	return P.MakePrism(
		func(ct ConnectionType) O.Option[MySQL] {
			if my, ok := ct.(MySQL); ok {
				return O.Some(my)
			}
			return O.None[MySQL]()
		},
		func(my MySQL) ConnectionType { return my },
	)
}

func mongodbPrism() P.Prism[ConnectionType, MongoDB] {
	return P.MakePrism(
		func(ct ConnectionType) O.Option[MongoDB] {
			if mg, ok := ct.(MongoDB); ok {
				return O.Some(mg)
			}
			return O.None[MongoDB]()
		},
		func(mg MongoDB) ConnectionType { return mg },
	)
}

// Helper function to create connection lens
func connectionLens() L.Lens[Config, ConnectionType] {
	return L.MakeLens(
		func(c Config) ConnectionType { return c.Connection },
		func(c Config, ct ConnectionType) Config {
			c.Connection = ct
			return c
		},
	)
}

// Helper function to create nil-safe connection lens for pointer types
func connectionLensRef() L.Lens[*Config, ConnectionType] {
	return L.MakeLensRef(
		func(c *Config) ConnectionType {
			if c == nil {
				return nil
			}
			return c.Connection
		},
		func(c *Config, ct ConnectionType) *Config {
			if c == nil {
				return &Config{Connection: ct}
			}
			c.Connection = ct
			return c
		},
	)
}

// TestComposeBasicFunctionality tests basic composition behavior
func TestComposeBasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when Prism matches", func(t *testing.T) {
		connLens := connectionLens()
		pgPrism := postgresqlPrism()

		// Compose connection lens with PostgreSQL prism
		configPgOptional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		result := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsSome(result))(t)

		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("localhost")(pg.Host)(t)
		assert.Equal(5432)(pg.Port)(t)
	})

	t.Run("GetOption returns None when Prism doesn't match", func(t *testing.T) {
		connLens := connectionLens()
		pgPrism := postgresqlPrism()

		configPgOptional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		result := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(result))(t)
	})

	t.Run("Set updates value when Prism matches", func(t *testing.T) {
		connLens := connectionLens()
		pgPrism := postgresqlPrism()

		configPgOptional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify the update
		result := configPgOptional.GetOption(updated)
		assert.Equal(true)(O.IsSome(result))(t)

		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("remote.example.com")(pg.Host)(t)
		assert.Equal(5433)(pg.Port)(t)

		// Verify other fields are unchanged
		assert.Equal("TestApp")(updated.AppName)(t)
	})

	t.Run("Set is no-op when Prism doesn't match", func(t *testing.T) {
		connLens := connectionLens()
		pgPrism := postgresqlPrism()

		configPgOptional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify nothing changed (no-op)
		assert.Equal(config)(updated)(t)

		// Verify the connection is still MySQL
		if my, ok := updated.Connection.(MySQL); ok {
			assert.Equal("localhost")(my.Host)(t)
			assert.Equal(3306)(my.Port)(t)
		} else {
			t.Fatal("Expected MySQL connection to remain unchanged")
		}
	})
}

// TestComposeOptionalLaws tests that the composition satisfies Optional laws
// Reference: https://gcanti.github.io/monocle-ts/modules/Optional.ts.html
func TestComposeOptionalLaws(t *testing.T) {
	connLens := connectionLens()
	pgPrism := postgresqlPrism()
	configPgOptional := Compose[Config](pgPrism)(connLens)

	t.Run("SetGet Law: GetOption(Set(b)(s)) = Some(b) when GetOption(s) = Some(_)", func(t *testing.T) {
		// Start with a config that has PostgreSQL
		config := Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		// Verify the prism matches
		initial := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsSome(initial))(t)

		// Set a new value
		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Get the value back
		result := configPgOptional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.Equal(true)(O.IsSome(result))(t)
		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal(newPg.Host)(pg.Host)(t)
		assert.Equal(newPg.Port)(pg.Port)(t)
	})

	t.Run("GetSet Law: Set(b)(s) = s when GetOption(s) = None (no-op)", func(t *testing.T) {
		// Start with a config that has MySQL (not PostgreSQL)
		config := Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		// Verify the prism doesn't match
		initial := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(initial))(t)

		// Try to set a PostgreSQL value
		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify GetSet law: structure should be unchanged (no-op)
		assert.Equal(config)(updated)(t)
	})

	t.Run("SetSet Law: Set(b2)(Set(b1)(s)) = Set(b2)(s)", func(t *testing.T) {
		config := Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		pg1 := PostgreSQL{Host: "server1.example.com", Port: 5433}
		pg2 := PostgreSQL{Host: "server2.example.com", Port: 5434}

		// Set twice
		setTwice := configPgOptional.Set(pg2)(configPgOptional.Set(pg1)(config))

		// Set once with the final value
		setOnce := configPgOptional.Set(pg2)(config)

		// They should be equal
		assert.Equal(setOnce)(setTwice)(t)

		// Verify the final value
		result := configPgOptional.GetOption(setTwice)
		assert.Equal(true)(O.IsSome(result))(t)
		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal(pg2.Host)(pg.Host)(t)
		assert.Equal(pg2.Port)(pg.Port)(t)
	})
}

// TestComposeMultipleVariants tests composition with different prism variants
func TestComposeMultipleVariants(t *testing.T) {
	connLens := connectionLens()

	t.Run("PostgreSQL variant", func(t *testing.T) {
		pgPrism := postgresqlPrism()
		optional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: PostgreSQL{Host: "pg.example.com", Port: 5432},
		}

		result := optional.GetOption(config)
		assert.Equal(true)(O.IsSome(result))(t)
	})

	t.Run("MySQL variant", func(t *testing.T) {
		myPrism := mysqlPrism()
		optional := Compose[Config](myPrism)(connLens)

		config := Config{
			Connection: MySQL{Host: "mysql.example.com", Port: 3306},
		}

		result := optional.GetOption(config)
		assert.Equal(true)(O.IsSome(result))(t)
	})

	t.Run("MongoDB variant", func(t *testing.T) {
		mgPrism := mongodbPrism()
		optional := Compose[Config](mgPrism)(connLens)

		config := Config{
			Connection: MongoDB{Host: "mongo.example.com", Port: 27017},
		}

		result := optional.GetOption(config)
		assert.Equal(true)(O.IsSome(result))(t)
	})

	t.Run("Cross-variant no-op", func(t *testing.T) {
		// Try to use PostgreSQL optional on MySQL config
		pgPrism := postgresqlPrism()
		optional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: MySQL{Host: "mysql.example.com", Port: 3306},
		}

		// GetOption should return None
		result := optional.GetOption(config)
		assert.Equal(true)(O.IsNone(result))(t)

		// Set should be no-op
		newPg := PostgreSQL{Host: "pg.example.com", Port: 5432}
		updated := optional.Set(newPg)(config)
		assert.Equal(config)(updated)(t)
	})
}

// TestComposeEdgeCases tests edge cases and boundary conditions
func TestComposeEdgeCases(t *testing.T) {
	t.Run("Identity lens with prism", func(t *testing.T) {
		// Identity lens that doesn't transform the value
		idLens := L.MakeLens(
			func(ct ConnectionType) ConnectionType { return ct },
			func(_ ConnectionType, ct ConnectionType) ConnectionType { return ct },
		)

		pgPrism := postgresqlPrism()
		optional := Compose[ConnectionType](pgPrism)(idLens)

		conn := ConnectionType(PostgreSQL{Host: "localhost", Port: 5432})
		result := optional.GetOption(conn)
		assert.Equal(true)(O.IsSome(result))(t)

		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("localhost")(pg.Host)(t)
	})

	t.Run("Multiple sets preserve structure", func(t *testing.T) {
		connLens := connectionLens()
		pgPrism := postgresqlPrism()
		optional := Compose[Config](pgPrism)(connLens)

		config := Config{
			Connection: PostgreSQL{Host: "host1", Port: 5432},
			AppName:    "TestApp",
		}

		// Apply multiple sets
		pg2 := PostgreSQL{Host: "host2", Port: 5433}
		pg3 := PostgreSQL{Host: "host3", Port: 5434}
		pg4 := PostgreSQL{Host: "host4", Port: 5435}

		updated := F.Pipe3(
			config,
			optional.Set(pg2),
			optional.Set(pg3),
			optional.Set(pg4),
		)

		// Verify final value
		result := optional.GetOption(updated)
		assert.Equal(true)(O.IsSome(result))(t)
		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("host4")(pg.Host)(t)
		assert.Equal(5435)(pg.Port)(t)

		// Verify structure is preserved
		assert.Equal("TestApp")(updated.AppName)(t)
	})
}

// TestComposeDocumentationExample tests the example from the documentation
func TestComposeDocumentationExample(t *testing.T) {
	// This test verifies the example code in the documentation works correctly

	// Lens to focus on Connection field
	connLens := L.MakeLens(
		func(c Config) ConnectionType { return c.Connection },
		func(c Config, ct ConnectionType) Config { c.Connection = ct; return c },
	)

	// Prism to extract PostgreSQL from ConnectionType
	pgPrism := P.MakePrism(
		func(ct ConnectionType) O.Option[PostgreSQL] {
			if pg, ok := ct.(PostgreSQL); ok {
				return O.Some(pg)
			}
			return O.None[PostgreSQL]()
		},
		func(pg PostgreSQL) ConnectionType { return pg },
	)

	// Compose to create Optional[Config, PostgreSQL]
	configPgOptional := Compose[Config](pgPrism)(connLens)

	config := Config{Connection: PostgreSQL{Host: "localhost"}}
	host := configPgOptional.GetOption(config) // Some(PostgreSQL{Host: "localhost"})
	assert.Equal(true)(O.IsSome(host))(t)

	updated := configPgOptional.Set(PostgreSQL{Host: "remote"})(config)
	// updated.Connection = PostgreSQL{Host: "remote"}
	result := configPgOptional.GetOption(updated)
	assert.Equal(true)(O.IsSome(result))(t)
	pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
	assert.Equal("remote")(pg.Host)(t)

	configMySQL := Config{Connection: MySQL{Host: "localhost"}}
	none := configPgOptional.GetOption(configMySQL) // None (Prism doesn't match)
	assert.Equal(true)(O.IsNone(none))(t)

	unchanged := configPgOptional.Set(PostgreSQL{Host: "remote"})(configMySQL)
	// unchanged == configMySQL (no-op because Prism doesn't match)
	assert.Equal(configMySQL)(unchanged)(t)
}

// TestComposeRefBasicFunctionality tests basic ComposeRef behavior with pointer types
func TestComposeRefBasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when Prism matches (non-nil pointer)", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		config := &Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		result := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsSome(result))(t)

		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("localhost")(pg.Host)(t)
		assert.Equal(5432)(pg.Port)(t)
	})

	t.Run("GetOption returns None when pointer is nil", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil

		result := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(result))(t)
	})

	t.Run("GetOption returns None when Prism doesn't match", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		config := &Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		result := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(result))(t)
	})

	t.Run("Set updates value when Prism matches (creates copy)", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		original := &Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(original)

		// Verify the update
		result := configPgOptional.GetOption(updated)
		assert.Equal(true)(O.IsSome(result))(t)

		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("remote.example.com")(pg.Host)(t)
		assert.Equal(5433)(pg.Port)(t)

		// Verify immutability: original should be unchanged
		if origPg, ok := original.Connection.(PostgreSQL); ok {
			assert.Equal("localhost")(origPg.Host)(t)
			assert.Equal(5432)(origPg.Port)(t)
		} else {
			t.Fatal("Original config should still have PostgreSQL connection")
		}

		// Verify they are different pointers
		if original == updated {
			t.Fatal("Set should create a new pointer, not modify in place")
		}
	})

	t.Run("Set is no-op when pointer is nil", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil

		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify nothing changed (no-op for nil)
		if updated != nil {
			t.Fatalf("Expected nil, got %v", updated)
		}
	})

	t.Run("Set is no-op when Prism doesn't match", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()

		configPgOptional := ComposeRef[Config](pgPrism)(connLens)

		original := &Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(original)

		// Verify nothing changed (no-op)
		assert.Equal(original)(updated)(t)

		// Verify the connection is still MySQL
		if my, ok := updated.Connection.(MySQL); ok {
			assert.Equal("localhost")(my.Host)(t)
			assert.Equal(3306)(my.Port)(t)
		} else {
			t.Fatal("Expected MySQL connection to remain unchanged")
		}
	})
}

// TestComposeRefOptionalLaws tests that ComposeRef satisfies Optional laws
func TestComposeRefOptionalLaws(t *testing.T) {
	connLens := connectionLensRef()
	pgPrism := postgresqlPrism()
	configPgOptional := ComposeRef[Config](pgPrism)(connLens)

	t.Run("SetGet Law: GetOption(Set(b)(s)) = Some(b) when GetOption(s) = Some(_)", func(t *testing.T) {
		// Start with a config that has PostgreSQL
		config := &Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		// Verify the prism matches
		initial := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsSome(initial))(t)

		// Set a new value
		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Get the value back
		result := configPgOptional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.Equal(true)(O.IsSome(result))(t)
		pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal(newPg.Host)(pg.Host)(t)
		assert.Equal(newPg.Port)(pg.Port)(t)
	})

	t.Run("GetSet Law: Set(b)(s) = s when GetOption(s) = None (no-op for nil)", func(t *testing.T) {
		// Start with nil config
		var config *Config = nil

		// Verify the prism doesn't match
		initial := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(initial))(t)

		// Try to set a PostgreSQL value
		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify GetSet law: structure should be unchanged (nil)
		if updated != nil {
			t.Fatalf("Expected nil, got %v", updated)
		}
	})

	t.Run("GetSet Law: Set(b)(s) = s when GetOption(s) = None (no-op for mismatched prism)", func(t *testing.T) {
		// Start with a config that has MySQL (not PostgreSQL)
		config := &Config{
			Connection: MySQL{Host: "localhost", Port: 3306},
			AppName:    "TestApp",
		}

		// Verify the prism doesn't match
		initial := configPgOptional.GetOption(config)
		assert.Equal(true)(O.IsNone(initial))(t)

		// Try to set a PostgreSQL value
		newPg := PostgreSQL{Host: "remote.example.com", Port: 5433}
		updated := configPgOptional.Set(newPg)(config)

		// Verify GetSet law: structure should be unchanged
		assert.Equal(config)(updated)(t)
	})

	t.Run("SetSet Law: Set(b2)(Set(b1)(s)) = Set(b2)(s)", func(t *testing.T) {
		config := &Config{
			Connection: PostgreSQL{Host: "localhost", Port: 5432},
			AppName:    "TestApp",
		}

		pg1 := PostgreSQL{Host: "server1.example.com", Port: 5433}
		pg2 := PostgreSQL{Host: "server2.example.com", Port: 5434}

		// Set twice
		setTwice := configPgOptional.Set(pg2)(configPgOptional.Set(pg1)(config))

		// Set once with the final value
		setOnce := configPgOptional.Set(pg2)(config)

		// They should be equal in value (but different pointers due to immutability)
		result1 := configPgOptional.GetOption(setTwice)
		result2 := configPgOptional.GetOption(setOnce)

		assert.Equal(true)(O.IsSome(result1))(t)
		assert.Equal(true)(O.IsSome(result2))(t)

		pg1Result := O.GetOrElse(F.Constant(PostgreSQL{}))(result1)
		pg2Result := O.GetOrElse(F.Constant(PostgreSQL{}))(result2)

		assert.Equal(pg2.Host)(pg1Result.Host)(t)
		assert.Equal(pg2.Port)(pg1Result.Port)(t)
		assert.Equal(pg2.Host)(pg2Result.Host)(t)
		assert.Equal(pg2.Port)(pg2Result.Port)(t)
	})
}

// TestComposeRefImmutability tests that ComposeRef preserves immutability
func TestComposeRefImmutability(t *testing.T) {
	t.Run("Set creates a new pointer, doesn't modify original", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()
		optional := ComposeRef[Config](pgPrism)(connLens)

		original := &Config{
			Connection: PostgreSQL{Host: "original", Port: 5432},
			AppName:    "OriginalApp",
		}

		// Store original values
		origPg := original.Connection.(PostgreSQL)
		origAppName := original.AppName

		// Perform multiple sets
		pg1 := PostgreSQL{Host: "host1", Port: 5433}
		pg2 := PostgreSQL{Host: "host2", Port: 5434}
		pg3 := PostgreSQL{Host: "host3", Port: 5435}

		updated1 := optional.Set(pg1)(original)
		updated2 := optional.Set(pg2)(updated1)
		updated3 := optional.Set(pg3)(updated2)

		// Verify original is unchanged
		currentPg := original.Connection.(PostgreSQL)
		assert.Equal(origPg.Host)(currentPg.Host)(t)
		assert.Equal(origPg.Port)(currentPg.Port)(t)
		assert.Equal(origAppName)(original.AppName)(t)

		// Verify final update has correct value
		result := optional.GetOption(updated3)
		assert.Equal(true)(O.IsSome(result))(t)
		finalPg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
		assert.Equal("host3")(finalPg.Host)(t)
		assert.Equal(5435)(finalPg.Port)(t)

		// Verify all pointers are different
		if original == updated1 || original == updated2 || original == updated3 {
			t.Fatal("Set should create new pointers, not modify in place")
		}
		if updated1 == updated2 || updated2 == updated3 || updated1 == updated3 {
			t.Fatal("Each Set should create a new pointer")
		}
	})

	t.Run("Multiple operations on nil preserve nil", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()
		optional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil

		// Multiple sets on nil should all return nil
		pg1 := PostgreSQL{Host: "host1", Port: 5433}
		pg2 := PostgreSQL{Host: "host2", Port: 5434}

		updated1 := optional.Set(pg1)(config)
		updated2 := optional.Set(pg2)(updated1)

		if updated1 != nil {
			t.Fatalf("Expected nil after first set, got %v", updated1)
		}
		if updated2 != nil {
			t.Fatalf("Expected nil after second set, got %v", updated2)
		}
	})
}

// TestComposeRefNilPointerEdgeCases tests edge cases with nil pointers
func TestComposeRefNilPointerEdgeCases(t *testing.T) {
	t.Run("GetOption on nil returns None", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()
		optional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil
		result := optional.GetOption(config)

		assert.Equal(true)(O.IsNone(result))(t)
	})

	t.Run("Set on nil with matching prism returns nil", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()
		optional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil
		newPg := PostgreSQL{Host: "remote", Port: 5432}
		updated := optional.Set(newPg)(config)

		if updated != nil {
			t.Fatalf("Expected nil, got %v", updated)
		}
	})

	t.Run("Chaining operations starting from nil", func(t *testing.T) {
		connLens := connectionLensRef()
		pgPrism := postgresqlPrism()
		optional := ComposeRef[Config](pgPrism)(connLens)

		var config *Config = nil

		// Chain multiple operations
		pg1 := PostgreSQL{Host: "host1", Port: 5433}
		pg2 := PostgreSQL{Host: "host2", Port: 5434}

		result := F.Pipe2(
			config,
			optional.Set(pg1),
			optional.Set(pg2),
		)

		if result != nil {
			t.Fatalf("Expected nil after chained operations, got %v", result)
		}
	})
}

// TestComposeRefDocumentationExample tests the example from the ComposeRef documentation
func TestComposeRefDocumentationExample(t *testing.T) {
	// Lens to focus on Connection field (pointer-based)
	connLens := connectionLensRef()

	// Prism to extract PostgreSQL from ConnectionType
	pgPrism := P.MakePrism(
		func(ct ConnectionType) O.Option[PostgreSQL] {
			if pg, ok := ct.(PostgreSQL); ok {
				return O.Some(pg)
			}
			return O.None[PostgreSQL]()
		},
		func(pg PostgreSQL) ConnectionType { return pg },
	)

	// Compose to create Optional[*Config, PostgreSQL]
	configPgOptional := ComposeRef[Config](pgPrism)(connLens)

	// Works with non-nil pointers
	config := &Config{Connection: PostgreSQL{Host: "localhost"}}
	host := configPgOptional.GetOption(config) // Some(PostgreSQL{Host: "localhost"})
	assert.Equal(true)(O.IsSome(host))(t)

	updated := configPgOptional.Set(PostgreSQL{Host: "remote"})(config)
	// updated is a new *Config with Connection = PostgreSQL{Host: "remote"}
	result := configPgOptional.GetOption(updated)
	assert.Equal(true)(O.IsSome(result))(t)
	pg := O.GetOrElse(F.Constant(PostgreSQL{}))(result)
	assert.Equal("remote")(pg.Host)(t)

	// original config is unchanged (immutability preserved)
	origPg := config.Connection.(PostgreSQL)
	assert.Equal("localhost")(origPg.Host)(t)

	// Handles nil pointers safely
	var nilConfig *Config = nil
	none := configPgOptional.GetOption(nilConfig) // None (nil pointer)
	assert.Equal(true)(O.IsNone(none))(t)

	unchanged := configPgOptional.Set(PostgreSQL{Host: "remote"})(nilConfig)
	// unchanged == nil (no-op because source is nil)
	if unchanged != nil {
		t.Fatalf("Expected nil, got %v", unchanged)
	}

	// Works with mismatched prisms
	configMySQL := &Config{Connection: MySQL{Host: "localhost"}}
	none = configPgOptional.GetOption(configMySQL) // None (Prism doesn't match)
	assert.Equal(true)(O.IsNone(none))(t)

	unchanged = configPgOptional.Set(PostgreSQL{Host: "remote"})(configMySQL)
	// unchanged == configMySQL (no-op because Prism doesn't match)
	assert.Equal(configMySQL)(unchanged)(t)
}
