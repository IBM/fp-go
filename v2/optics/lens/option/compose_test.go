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

package option

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// Test types for ComposeOption - using unique names to avoid conflicts
type (
	DatabaseCfg struct {
		Host string
		Port int
	}

	ServerConfig struct {
		Database *DatabaseCfg
	}

	AppSettings struct {
		MaxRetries int
		Timeout    int
	}

	ApplicationConfig struct {
		Settings *AppSettings
	}
)

// Helper methods for DatabaseCfg
func (db *DatabaseCfg) GetPort() int {
	return db.Port
}

func (db *DatabaseCfg) SetPort(port int) *DatabaseCfg {
	db.Port = port
	return db
}

// Helper methods for ServerConfig
func (c ServerConfig) GetDatabase() *DatabaseCfg {
	return c.Database
}

func (c ServerConfig) SetDatabase(db *DatabaseCfg) ServerConfig {
	c.Database = db
	return c
}

// Helper methods for AppSettings
func (s *AppSettings) GetMaxRetries() int {
	return s.MaxRetries
}

func (s *AppSettings) SetMaxRetries(retries int) *AppSettings {
	s.MaxRetries = retries
	return s
}

// Helper methods for ApplicationConfig
func (ac ApplicationConfig) GetSettings() *AppSettings {
	return ac.Settings
}

func (ac ApplicationConfig) SetSettings(s *AppSettings) ApplicationConfig {
	ac.Settings = s
	return ac
}

// TestComposeOptionBasicOperations tests basic get/set operations
func TestComposeOptionBasicOperations(t *testing.T) {
	// Create lenses
	dbLens := FromNillable(L.MakeLens(ServerConfig.GetDatabase, ServerConfig.SetDatabase))
	portLens := L.MakeLensRef((*DatabaseCfg).GetPort, (*DatabaseCfg).SetPort)

	defaultDB := &DatabaseCfg{Host: "localhost", Port: 5432}
	configPortLens := F.Pipe1(dbLens, ComposeOption[ServerConfig, int](defaultDB)(portLens))

	t.Run("Get from empty config returns None", func(t *testing.T) {
		config := ServerConfig{Database: nil}
		result := configPortLens.Get(config)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get from config with database returns Some", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
		result := configPortLens.Get(config)
		assert.Equal(t, O.Some(3306), result)
	})

	t.Run("Set Some on empty config creates database with default", func(t *testing.T) {
		config := ServerConfig{Database: nil}
		updated := configPortLens.Set(O.Some(3306))(config)
		assert.NotNil(t, updated.Database)
		assert.Equal(t, 3306, updated.Database.Port)
		assert.Equal(t, "localhost", updated.Database.Host) // From default
	})

	t.Run("Set Some on existing database updates port", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 5432}}
		updated := configPortLens.Set(O.Some(8080))(config)
		assert.NotNil(t, updated.Database)
		assert.Equal(t, 8080, updated.Database.Port)
		assert.Equal(t, "example.com", updated.Database.Host) // Preserved
	})

	t.Run("Set None removes database entirely", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
		updated := configPortLens.Set(O.None[int]())(config)
		assert.Nil(t, updated.Database)
	})

	t.Run("Set None on empty config is no-op", func(t *testing.T) {
		config := ServerConfig{Database: nil}
		updated := configPortLens.Set(O.None[int]())(config)
		assert.Nil(t, updated.Database)
	})
}

// TestComposeOptionLensLawsDetailed verifies that ComposeOption satisfies lens laws
func TestComposeOptionLensLawsDetailed(t *testing.T) {
	// Setup
	defaultDB := &DatabaseCfg{Host: "localhost", Port: 5432}
	dbLens := FromNillable(L.MakeLens(ServerConfig.GetDatabase, ServerConfig.SetDatabase))
	portLens := L.MakeLensRef((*DatabaseCfg).GetPort, (*DatabaseCfg).SetPort)
	configPortLens := F.Pipe1(dbLens, ComposeOption[ServerConfig, int](defaultDB)(portLens))

	// Equality predicates
	eqInt := EQT.Eq[int]()
	eqOptInt := O.Eq(eqInt)
	eqServerConfig := func(a, b ServerConfig) bool {
		if a.Database == nil && b.Database == nil {
			return true
		}
		if a.Database == nil || b.Database == nil {
			return false
		}
		return a.Database.Host == b.Database.Host && a.Database.Port == b.Database.Port
	}

	// Test structures
	configNil := ServerConfig{Database: nil}
	config3306 := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
	config5432 := ServerConfig{Database: &DatabaseCfg{Host: "test.com", Port: 5432}}

	// Law 1: GetSet - lens.Get(lens.Set(a)(s)) == a
	t.Run("Law1_GetSet_WithSome", func(t *testing.T) {
		// Setting Some(8080) and getting back should return Some(8080)
		result := configPortLens.Get(configPortLens.Set(O.Some(8080))(config3306))
		assert.True(t, eqOptInt.Equals(result, O.Some(8080)),
			"Get(Set(Some(8080))(s)) should equal Some(8080)")
	})

	t.Run("Law1_GetSet_WithNone", func(t *testing.T) {
		// Setting None and getting back should return None
		result := configPortLens.Get(configPortLens.Set(O.None[int]())(config3306))
		assert.True(t, eqOptInt.Equals(result, O.None[int]()),
			"Get(Set(None)(s)) should equal None")
	})

	t.Run("Law1_GetSet_OnEmptyWithSome", func(t *testing.T) {
		// Setting Some on empty config and getting back
		result := configPortLens.Get(configPortLens.Set(O.Some(9000))(configNil))
		assert.True(t, eqOptInt.Equals(result, O.Some(9000)),
			"Get(Set(Some(9000))(empty)) should equal Some(9000)")
	})

	// Law 2: SetGet - lens.Set(lens.Get(s))(s) == s
	t.Run("Law2_SetGet_WithDatabase", func(t *testing.T) {
		// Setting what we get should return the same structure
		result := configPortLens.Set(configPortLens.Get(config3306))(config3306)
		assert.True(t, eqServerConfig(result, config3306),
			"Set(Get(s))(s) should equal s")
	})

	t.Run("Law2_SetGet_WithoutDatabase", func(t *testing.T) {
		// Setting what we get from empty should return the same structure
		result := configPortLens.Set(configPortLens.Get(configNil))(configNil)
		assert.True(t, eqServerConfig(result, configNil),
			"Set(Get(empty))(empty) should equal empty")
	})

	t.Run("Law2_SetGet_DifferentConfigs", func(t *testing.T) {
		// Test with another config
		result := configPortLens.Set(configPortLens.Get(config5432))(config5432)
		assert.True(t, eqServerConfig(result, config5432),
			"Set(Get(s))(s) should equal s for any s")
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("Law3_SetSet_BothSome", func(t *testing.T) {
		// Setting twice with Some should be same as setting once
		setTwice := configPortLens.Set(O.Some(9000))(configPortLens.Set(O.Some(8080))(config3306))
		setOnce := configPortLens.Set(O.Some(9000))(config3306)
		assert.True(t, eqServerConfig(setTwice, setOnce),
			"Set(a2)(Set(a1)(s)) should equal Set(a2)(s)")
	})

	t.Run("Law3_SetSet_BothNone", func(t *testing.T) {
		// Setting None twice should be same as setting once
		setTwice := configPortLens.Set(O.None[int]())(configPortLens.Set(O.None[int]())(config3306))
		setOnce := configPortLens.Set(O.None[int]())(config3306)
		assert.True(t, eqServerConfig(setTwice, setOnce),
			"Set(None)(Set(None)(s)) should equal Set(None)(s)")
	})

	t.Run("Law3_SetSet_SomeThenNone", func(t *testing.T) {
		// Setting None after Some should be same as setting None directly
		setTwice := configPortLens.Set(O.None[int]())(configPortLens.Set(O.Some(8080))(config3306))
		setOnce := configPortLens.Set(O.None[int]())(config3306)
		assert.True(t, eqServerConfig(setTwice, setOnce),
			"Set(None)(Set(Some)(s)) should equal Set(None)(s)")
	})

	t.Run("Law3_SetSet_NoneThenSome", func(t *testing.T) {
		// Setting Some after None creates a new database with default values
		// This is different from setting Some directly which preserves existing fields
		setTwice := configPortLens.Set(O.Some(8080))(configPortLens.Set(O.None[int]())(config3306))
		// After setting None, the database is removed, so setting Some creates it with defaults
		assert.NotNil(t, setTwice.Database)
		assert.Equal(t, 8080, setTwice.Database.Port)
		assert.Equal(t, "localhost", setTwice.Database.Host) // From default, not "example.com"

		// This demonstrates that ComposeOption's behavior when setting None then Some
		// uses the default value for the intermediate structure
		setOnce := configPortLens.Set(O.Some(8080))(config3306)
		assert.Equal(t, 8080, setOnce.Database.Port)
		assert.Equal(t, "example.com", setOnce.Database.Host) // Preserved from original

		// They are NOT equal because the Host field differs
		assert.False(t, eqServerConfig(setTwice, setOnce),
			"Set(Some)(Set(None)(s)) uses default, Set(Some)(s) preserves fields")
	})

	t.Run("Law3_SetSet_OnEmpty", func(t *testing.T) {
		// Setting twice on empty config
		setTwice := configPortLens.Set(O.Some(9000))(configPortLens.Set(O.Some(8080))(configNil))
		setOnce := configPortLens.Set(O.Some(9000))(configNil)
		assert.True(t, eqServerConfig(setTwice, setOnce),
			"Set(a2)(Set(a1)(empty)) should equal Set(a2)(empty)")
	})
}

// TestComposeOptionWithModify tests the Modify operation
func TestComposeOptionWithModify(t *testing.T) {
	defaultDB := &DatabaseCfg{Host: "localhost", Port: 5432}
	dbLens := FromNillable(L.MakeLens(ServerConfig.GetDatabase, ServerConfig.SetDatabase))
	portLens := L.MakeLensRef((*DatabaseCfg).GetPort, (*DatabaseCfg).SetPort)
	configPortLens := F.Pipe1(dbLens, ComposeOption[ServerConfig, int](defaultDB)(portLens))

	t.Run("Modify with identity returns same structure", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
		result := L.Modify[ServerConfig](F.Identity[Option[int]])(configPortLens)(config)
		assert.Equal(t, config.Database.Port, result.Database.Port)
		assert.Equal(t, config.Database.Host, result.Database.Host)
	})

	t.Run("Modify with Some transformation", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
		// Double the port if it exists
		doublePort := O.Map(func(p int) int { return p * 2 })
		result := L.Modify[ServerConfig](doublePort)(configPortLens)(config)
		assert.Equal(t, 6612, result.Database.Port)
		assert.Equal(t, "example.com", result.Database.Host)
	})

	t.Run("Modify on empty config with Some transformation", func(t *testing.T) {
		config := ServerConfig{Database: nil}
		doublePort := O.Map(func(p int) int { return p * 2 })
		result := L.Modify[ServerConfig](doublePort)(configPortLens)(config)
		// Should remain empty since there's nothing to modify
		assert.Nil(t, result.Database)
	})
}

// TestComposeOptionComposition tests composing multiple ComposeOption lenses
func TestComposeOptionComposition(t *testing.T) {
	type Level3 struct {
		Value int
	}

	type Level2 struct {
		Level3 *Level3
	}

	type Level1 struct {
		Level2 *Level2
	}

	// Create lenses
	level2Lens := FromNillable(L.MakeLens(
		func(l1 Level1) *Level2 { return l1.Level2 },
		func(l1 Level1, l2 *Level2) Level1 { l1.Level2 = l2; return l1 },
	))

	level3Lens := L.MakeLensRef(
		func(l2 *Level2) *Level3 { return l2.Level3 },
		func(l2 *Level2, l3 *Level3) *Level2 { l2.Level3 = l3; return l2 },
	)

	valueLens := L.MakeLensRef(
		func(l3 *Level3) int { return l3.Value },
		func(l3 *Level3, v int) *Level3 { l3.Value = v; return l3 },
	)

	// Compose: Level1 -> Option[Level2] -> Option[Level3] -> Option[int]
	defaultLevel2 := &Level2{Level3: &Level3{Value: 0}}
	defaultLevel3 := &Level3{Value: 0}

	// First composition: Level1 -> Option[Level3]
	level1ToLevel3 := F.Pipe1(level2Lens, ComposeOption[Level1, *Level3](defaultLevel2)(level3Lens))

	// Second composition: Level1 -> Option[int]
	level1ToValue := F.Pipe1(level1ToLevel3, ComposeOption[Level1, int](defaultLevel3)(valueLens))

	t.Run("Get from fully populated structure", func(t *testing.T) {
		l1 := Level1{Level2: &Level2{Level3: &Level3{Value: 42}}}
		result := level1ToValue.Get(l1)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("Get from empty structure", func(t *testing.T) {
		l1 := Level1{Level2: nil}
		result := level1ToValue.Get(l1)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Set on empty structure creates all levels", func(t *testing.T) {
		l1 := Level1{Level2: nil}
		updated := level1ToValue.Set(O.Some(100))(l1)
		assert.NotNil(t, updated.Level2)
		assert.NotNil(t, updated.Level2.Level3)
		assert.Equal(t, 100, updated.Level2.Level3.Value)
	})

	t.Run("Set None removes top level", func(t *testing.T) {
		l1 := Level1{Level2: &Level2{Level3: &Level3{Value: 42}}}
		updated := level1ToValue.Set(O.None[int]())(l1)
		assert.Nil(t, updated.Level2)
	})
}

// TestComposeOptionEdgeCasesExtended tests additional edge cases
func TestComposeOptionEdgeCasesExtended(t *testing.T) {
	defaultSettings := &AppSettings{MaxRetries: 3, Timeout: 30}
	settingsLens := FromNillable(L.MakeLens(ApplicationConfig.GetSettings, ApplicationConfig.SetSettings))
	retriesLens := L.MakeLensRef((*AppSettings).GetMaxRetries, (*AppSettings).SetMaxRetries)
	configRetriesLens := F.Pipe1(settingsLens, ComposeOption[ApplicationConfig, int](defaultSettings)(retriesLens))

	t.Run("Multiple sets with different values", func(t *testing.T) {
		config := ApplicationConfig{Settings: nil}
		// Set multiple times
		config = configRetriesLens.Set(O.Some(5))(config)
		assert.Equal(t, 5, config.Settings.MaxRetries)

		config = configRetriesLens.Set(O.Some(10))(config)
		assert.Equal(t, 10, config.Settings.MaxRetries)

		config = configRetriesLens.Set(O.None[int]())(config)
		assert.Nil(t, config.Settings)
	})

	t.Run("Get after Set maintains consistency", func(t *testing.T) {
		config := ApplicationConfig{Settings: nil}
		updated := configRetriesLens.Set(O.Some(7))(config)
		retrieved := configRetriesLens.Get(updated)
		assert.Equal(t, O.Some(7), retrieved)
	})

	t.Run("Default values are used correctly", func(t *testing.T) {
		config := ApplicationConfig{Settings: nil}
		updated := configRetriesLens.Set(O.Some(15))(config)
		// Check that default timeout is used
		assert.Equal(t, 30, updated.Settings.Timeout)
		assert.Equal(t, 15, updated.Settings.MaxRetries)
	})

	t.Run("Preserves other fields when updating", func(t *testing.T) {
		config := ApplicationConfig{Settings: &AppSettings{MaxRetries: 5, Timeout: 60}}
		updated := configRetriesLens.Set(O.Some(10))(config)
		assert.Equal(t, 10, updated.Settings.MaxRetries)
		assert.Equal(t, 60, updated.Settings.Timeout) // Preserved
	})
}

// TestComposeOptionWithZeroValues tests behavior with zero values
func TestComposeOptionWithZeroValues(t *testing.T) {
	defaultDB := &DatabaseCfg{Host: "", Port: 0}
	dbLens := FromNillable(L.MakeLens(ServerConfig.GetDatabase, ServerConfig.SetDatabase))
	portLens := L.MakeLensRef((*DatabaseCfg).GetPort, (*DatabaseCfg).SetPort)
	configPortLens := F.Pipe1(dbLens, ComposeOption[ServerConfig, int](defaultDB)(portLens))

	t.Run("Set zero value", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 3306}}
		updated := configPortLens.Set(O.Some(0))(config)
		assert.Equal(t, 0, updated.Database.Port)
		assert.Equal(t, "example.com", updated.Database.Host)
	})

	t.Run("Get zero value returns Some(0)", func(t *testing.T) {
		config := ServerConfig{Database: &DatabaseCfg{Host: "example.com", Port: 0}}
		result := configPortLens.Get(config)
		assert.Equal(t, O.Some(0), result)
	})

	t.Run("Default with zero values", func(t *testing.T) {
		config := ServerConfig{Database: nil}
		updated := configPortLens.Set(O.Some(8080))(config)
		assert.Equal(t, "", updated.Database.Host) // From default
		assert.Equal(t, 8080, updated.Database.Port)
	})
}

// ============================================================================
// Tests for Compose function (both lenses return Option values)
// ============================================================================

// TestComposeBasicOperations tests basic get/set operations for Compose
func TestComposeBasicOperations(t *testing.T) {
	type Value struct {
		Data *string
	}

	type Container struct {
		Value *Value
	}

	// Create lenses
	valueLens := FromNillable(L.MakeLens(
		func(c Container) *Value { return c.Value },
		func(c Container, v *Value) Container { c.Value = v; return c },
	))

	dataLens := L.MakeLensRef(
		func(v *Value) *string { return v.Data },
		func(v *Value, d *string) *Value { v.Data = d; return v },
	)

	defaultValue := &Value{Data: nil}
	composedLens := F.Pipe1(valueLens, Compose[Container, *string](defaultValue)(
		FromNillable(dataLens),
	))

	t.Run("Get from empty container returns None", func(t *testing.T) {
		container := Container{Value: nil}
		result := composedLens.Get(container)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get from container with nil data returns None", func(t *testing.T) {
		container := Container{Value: &Value{Data: nil}}
		result := composedLens.Get(container)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get from container with data returns Some", func(t *testing.T) {
		data := "test"
		container := Container{Value: &Value{Data: &data}}
		result := composedLens.Get(container)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, &data, O.GetOrElse(func() *string { return nil })(result))
	})

	t.Run("Set Some on empty container creates structure with default", func(t *testing.T) {
		container := Container{Value: nil}
		data := "new"
		updated := composedLens.Set(O.Some(&data))(container)
		assert.NotNil(t, updated.Value)
		assert.NotNil(t, updated.Value.Data)
		assert.Equal(t, "new", *updated.Value.Data)
	})

	t.Run("Set Some on existing container updates data", func(t *testing.T) {
		oldData := "old"
		container := Container{Value: &Value{Data: &oldData}}
		newData := "new"
		updated := composedLens.Set(O.Some(&newData))(container)
		assert.NotNil(t, updated.Value)
		assert.NotNil(t, updated.Value.Data)
		assert.Equal(t, "new", *updated.Value.Data)
	})

	t.Run("Set None when container is empty is no-op", func(t *testing.T) {
		container := Container{Value: nil}
		updated := composedLens.Set(O.None[*string]())(container)
		assert.Nil(t, updated.Value)
	})

	t.Run("Set None when container exists unsets data", func(t *testing.T) {
		data := "test"
		container := Container{Value: &Value{Data: &data}}
		updated := composedLens.Set(O.None[*string]())(container)
		assert.NotNil(t, updated.Value)
		assert.Nil(t, updated.Value.Data)
	})
}

// TestComposeLensLawsDetailed verifies that Compose satisfies lens laws
func TestComposeLensLawsDetailed(t *testing.T) {
	type Inner struct {
		Value *int
		Extra string
	}

	type Outer struct {
		Inner *Inner
	}

	// Setup
	defaultInner := &Inner{Value: nil, Extra: "default"}
	innerLens := FromNillable(L.MakeLens(
		func(o Outer) *Inner { return o.Inner },
		func(o Outer, i *Inner) Outer { o.Inner = i; return o },
	))
	valueLens := L.MakeLensRef(
		func(i *Inner) *int { return i.Value },
		func(i *Inner, v *int) *Inner { i.Value = v; return i },
	)
	composedLens := F.Pipe1(innerLens, Compose[Outer, *int](defaultInner)(
		FromNillable(valueLens),
	))

	// Equality predicates
	eqIntPtr := EQT.Eq[*int]()
	eqOptIntPtr := O.Eq(eqIntPtr)
	eqOuter := func(a, b Outer) bool {
		if a.Inner == nil && b.Inner == nil {
			return true
		}
		if a.Inner == nil || b.Inner == nil {
			return false
		}
		aVal := a.Inner.Value
		bVal := b.Inner.Value
		if aVal == nil && bVal == nil {
			return a.Inner.Extra == b.Inner.Extra
		}
		if aVal == nil || bVal == nil {
			return false
		}
		return *aVal == *bVal && a.Inner.Extra == b.Inner.Extra
	}

	// Test structures
	val42 := 42
	val100 := 100
	outerNil := Outer{Inner: nil}
	outerWithNilValue := Outer{Inner: &Inner{Value: nil, Extra: "test"}}
	outer42 := Outer{Inner: &Inner{Value: &val42, Extra: "test"}}

	// Law 1: GetSet - lens.Get(lens.Set(a)(s)) == a
	t.Run("Law1_GetSet_WithSome", func(t *testing.T) {
		result := composedLens.Get(composedLens.Set(O.Some(&val100))(outer42))
		assert.True(t, eqOptIntPtr.Equals(result, O.Some(&val100)),
			"Get(Set(Some(100))(s)) should equal Some(100)")
	})

	t.Run("Law1_GetSet_WithNone", func(t *testing.T) {
		result := composedLens.Get(composedLens.Set(O.None[*int]())(outer42))
		assert.True(t, eqOptIntPtr.Equals(result, O.None[*int]()),
			"Get(Set(None)(s)) should equal None")
	})

	t.Run("Law1_GetSet_OnEmpty", func(t *testing.T) {
		result := composedLens.Get(composedLens.Set(O.Some(&val100))(outerNil))
		assert.True(t, eqOptIntPtr.Equals(result, O.Some(&val100)),
			"Get(Set(Some(100))(empty)) should equal Some(100)")
	})

	// Law 2: SetGet - lens.Set(lens.Get(s))(s) == s
	t.Run("Law2_SetGet_WithValue", func(t *testing.T) {
		result := composedLens.Set(composedLens.Get(outer42))(outer42)
		assert.True(t, eqOuter(result, outer42),
			"Set(Get(s))(s) should equal s")
	})

	t.Run("Law2_SetGet_WithNilValue", func(t *testing.T) {
		result := composedLens.Set(composedLens.Get(outerWithNilValue))(outerWithNilValue)
		assert.True(t, eqOuter(result, outerWithNilValue),
			"Set(Get(s))(s) should equal s when value is nil")
	})

	t.Run("Law2_SetGet_WithNilInner", func(t *testing.T) {
		result := composedLens.Set(composedLens.Get(outerNil))(outerNil)
		assert.True(t, eqOuter(result, outerNil),
			"Set(Get(empty))(empty) should equal empty")
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("Law3_SetSet_BothSome", func(t *testing.T) {
		val200 := 200
		setTwice := composedLens.Set(O.Some(&val200))(composedLens.Set(O.Some(&val100))(outer42))
		setOnce := composedLens.Set(O.Some(&val200))(outer42)
		assert.True(t, eqOuter(setTwice, setOnce),
			"Set(a2)(Set(a1)(s)) should equal Set(a2)(s)")
	})

	t.Run("Law3_SetSet_BothNone", func(t *testing.T) {
		setTwice := composedLens.Set(O.None[*int]())(composedLens.Set(O.None[*int]())(outer42))
		setOnce := composedLens.Set(O.None[*int]())(outer42)
		assert.True(t, eqOuter(setTwice, setOnce),
			"Set(None)(Set(None)(s)) should equal Set(None)(s)")
	})

	t.Run("Law3_SetSet_SomeThenNone", func(t *testing.T) {
		setTwice := composedLens.Set(O.None[*int]())(composedLens.Set(O.Some(&val100))(outer42))
		setOnce := composedLens.Set(O.None[*int]())(outer42)
		assert.True(t, eqOuter(setTwice, setOnce),
			"Set(None)(Set(Some)(s)) should equal Set(None)(s)")
	})

	t.Run("Law3_SetSet_NoneThenSome", func(t *testing.T) {
		// This case is interesting: setting None then Some uses default
		setTwice := composedLens.Set(O.Some(&val100))(composedLens.Set(O.None[*int]())(outer42))
		// After None, inner still exists but value is nil
		// Then setting Some updates the value
		assert.NotNil(t, setTwice.Inner)
		assert.NotNil(t, setTwice.Inner.Value)
		assert.Equal(t, 100, *setTwice.Inner.Value)
		assert.Equal(t, "test", setTwice.Inner.Extra) // Preserved from original
	})
}

// TestComposeWithModify tests the Modify operation for Compose
func TestComposeWithModify(t *testing.T) {
	type Data struct {
		Count *int
	}

	type Store struct {
		Data *Data
	}

	defaultData := &Data{Count: nil}
	dataLens := FromNillable(L.MakeLens(
		func(s Store) *Data { return s.Data },
		func(s Store, d *Data) Store { s.Data = d; return s },
	))
	countLens := L.MakeLensRef(
		func(d *Data) *int { return d.Count },
		func(d *Data, c *int) *Data { d.Count = c; return d },
	)
	composedLens := F.Pipe1(dataLens, Compose[Store, *int](defaultData)(
		FromNillable(countLens),
	))

	t.Run("Modify with identity returns same structure", func(t *testing.T) {
		count := 5
		store := Store{Data: &Data{Count: &count}}
		result := L.Modify[Store](F.Identity[Option[*int]])(composedLens)(store)
		assert.Equal(t, 5, *result.Data.Count)
	})

	t.Run("Modify with Some transformation", func(t *testing.T) {
		count := 5
		store := Store{Data: &Data{Count: &count}}
		// Double the count if it exists
		doubleCount := O.Map(func(c *int) *int {
			doubled := *c * 2
			return &doubled
		})
		result := L.Modify[Store](doubleCount)(composedLens)(store)
		assert.Equal(t, 10, *result.Data.Count)
	})

	t.Run("Modify on empty store", func(t *testing.T) {
		store := Store{Data: nil}
		doubleCount := O.Map(func(c *int) *int {
			doubled := *c * 2
			return &doubled
		})
		result := L.Modify[Store](doubleCount)(composedLens)(store)
		// Should remain empty since there's nothing to modify
		assert.Nil(t, result.Data)
	})
}

// TestComposeMultiLevel tests composing multiple Compose operations
func TestComposeMultiLevel(t *testing.T) {
	type Level3 struct {
		Value *string
	}

	type Level2 struct {
		Level3 *Level3
	}

	type Level1 struct {
		Level2 *Level2
	}

	// Create lenses
	level2Lens := FromNillable(L.MakeLens(
		func(l1 Level1) *Level2 { return l1.Level2 },
		func(l1 Level1, l2 *Level2) Level1 { l1.Level2 = l2; return l1 },
	))

	level3Lens := L.MakeLensRef(
		func(l2 *Level2) *Level3 { return l2.Level3 },
		func(l2 *Level2, l3 *Level3) *Level2 { l2.Level3 = l3; return l2 },
	)

	valueLens := L.MakeLensRef(
		func(l3 *Level3) *string { return l3.Value },
		func(l3 *Level3, v *string) *Level3 { l3.Value = v; return l3 },
	)

	// Compose: Level1 -> Option[Level2] -> Option[Level3] -> Option[string]
	defaultLevel2 := &Level2{Level3: nil}
	defaultLevel3 := &Level3{Value: nil}

	// First composition: Level1 -> Option[Level3]
	level1ToLevel3 := F.Pipe1(level2Lens, Compose[Level1, *Level3](defaultLevel2)(
		FromNillable(level3Lens),
	))

	// Second composition: Level1 -> Option[string]
	level1ToValue := F.Pipe1(level1ToLevel3, Compose[Level1, *string](defaultLevel3)(
		FromNillable(valueLens),
	))

	t.Run("Get from fully populated structure", func(t *testing.T) {
		value := "test"
		l1 := Level1{Level2: &Level2{Level3: &Level3{Value: &value}}}
		result := level1ToValue.Get(l1)
		assert.True(t, O.IsSome(result))
	})

	t.Run("Get from partially populated structure", func(t *testing.T) {
		l1 := Level1{Level2: &Level2{Level3: &Level3{Value: nil}}}
		result := level1ToValue.Get(l1)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get from empty structure", func(t *testing.T) {
		l1 := Level1{Level2: nil}
		result := level1ToValue.Get(l1)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Set on empty structure creates all levels", func(t *testing.T) {
		l1 := Level1{Level2: nil}
		value := "new"
		updated := level1ToValue.Set(O.Some(&value))(l1)
		assert.NotNil(t, updated.Level2)
		assert.NotNil(t, updated.Level2.Level3)
		assert.NotNil(t, updated.Level2.Level3.Value)
		assert.Equal(t, "new", *updated.Level2.Level3.Value)
	})

	t.Run("Set None when structure exists unsets value", func(t *testing.T) {
		value := "test"
		l1 := Level1{Level2: &Level2{Level3: &Level3{Value: &value}}}
		updated := level1ToValue.Set(O.None[*string]())(l1)
		assert.NotNil(t, updated.Level2)
		assert.NotNil(t, updated.Level2.Level3)
		assert.Nil(t, updated.Level2.Level3.Value)
	})
}

// TestComposeEdgeCasesExtended tests additional edge cases for Compose
func TestComposeEdgeCasesExtended(t *testing.T) {
	type Metadata struct {
		Tags *[]string
	}

	type Document struct {
		Metadata *Metadata
	}

	defaultMetadata := &Metadata{Tags: nil}
	metadataLens := FromNillable(L.MakeLens(
		func(d Document) *Metadata { return d.Metadata },
		func(d Document, m *Metadata) Document { d.Metadata = m; return d },
	))
	tagsLens := L.MakeLensRef(
		func(m *Metadata) *[]string { return m.Tags },
		func(m *Metadata, t *[]string) *Metadata { m.Tags = t; return m },
	)
	composedLens := F.Pipe1(metadataLens, Compose[Document, *[]string](defaultMetadata)(
		FromNillable(tagsLens),
	))

	t.Run("Multiple sets with different values", func(t *testing.T) {
		doc := Document{Metadata: nil}
		tags1 := []string{"tag1"}
		tags2 := []string{"tag2", "tag3"}

		// Set first value
		doc = composedLens.Set(O.Some(&tags1))(doc)
		assert.NotNil(t, doc.Metadata)
		assert.NotNil(t, doc.Metadata.Tags)
		assert.Equal(t, 1, len(*doc.Metadata.Tags))

		// Set second value
		doc = composedLens.Set(O.Some(&tags2))(doc)
		assert.Equal(t, 2, len(*doc.Metadata.Tags))

		// Set None
		doc = composedLens.Set(O.None[*[]string]())(doc)
		assert.NotNil(t, doc.Metadata)
		assert.Nil(t, doc.Metadata.Tags)
	})

	t.Run("Get after Set maintains consistency", func(t *testing.T) {
		doc := Document{Metadata: nil}
		tags := []string{"test"}
		updated := composedLens.Set(O.Some(&tags))(doc)
		retrieved := composedLens.Get(updated)
		assert.True(t, O.IsSome(retrieved))
	})

	t.Run("Default values are used when creating structure", func(t *testing.T) {
		doc := Document{Metadata: nil}
		tags := []string{"new"}
		updated := composedLens.Set(O.Some(&tags))(doc)
		// Metadata should be created with default (Tags: nil initially, then set)
		assert.NotNil(t, updated.Metadata)
		assert.NotNil(t, updated.Metadata.Tags)
		assert.Equal(t, []string{"new"}, *updated.Metadata.Tags)
	})
}
