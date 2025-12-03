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

package lens

import (
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type (
	Street struct {
		num  int
		name string
	}

	Address struct {
		city   string
		street *Street
	}

	Inner struct {
		Value int
		Foo   string
	}

	InnerOpt struct {
		Value *int
		Foo   *string
	}

	Outer struct {
		inner *Inner
	}

	OuterOpt struct {
		inner *InnerOpt
	}
)

func (outer Outer) GetInner() *Inner {
	return outer.inner
}

func (outer Outer) SetInner(inner *Inner) Outer {
	outer.inner = inner
	return outer
}

func (outer OuterOpt) GetInnerOpt() *InnerOpt {
	return outer.inner
}

func (outer OuterOpt) SetInnerOpt(inner *InnerOpt) OuterOpt {
	outer.inner = inner
	return outer
}

func (inner *Inner) GetValue() int {
	return inner.Value
}

func (inner *Inner) SetValue(value int) *Inner {
	inner.Value = value
	return inner
}

func (inner *InnerOpt) GetValue() *int {
	return inner.Value
}

func (inner *InnerOpt) SetValue(value *int) *InnerOpt {
	inner.Value = value
	return inner
}

func (street *Street) GetName() string {
	return street.name
}

func (street *Street) SetName(name string) *Street {
	street.name = name
	return street
}

func (addr *Address) GetStreet() *Street {
	return addr.street
}

func (addr *Address) SetStreet(s *Street) *Address {
	addr.street = s
	return addr
}

var (
	streetLens = MakeLensRef((*Street).GetName, (*Street).SetName)
	addrLens   = MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

	sampleStreet  = Street{num: 220, name: "Schönaicherstr"}
	sampleAddress = Address{city: "Böblingen", street: &sampleStreet}
)

func TestLens(t *testing.T) {
	// read the value
	assert.Equal(t, sampleStreet.name, streetLens.Get(&sampleStreet))
	// new street
	newName := "Böblingerstr"
	// update
	old := sampleStreet
	updated := streetLens.Set(newName)(&sampleStreet)
	assert.Equal(t, old, sampleStreet)
	// validate the new name
	assert.Equal(t, newName, streetLens.Get(updated))
}

func TestAddressCompose(t *testing.T) {
	// compose
	streetName := Compose[*Address](streetLens)(addrLens)
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))
	// new street
	newName := "Böblingerstr"
	updated := streetName.Set(newName)(&sampleAddress)
	// check that we have not modified the original
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))
	assert.Equal(t, newName, streetName.Get(updated))
}

func TestIMap(t *testing.T) {

	type S struct {
		a int
	}

	sa := F.Pipe1(
		Id[S](),
		IMap[S](
			func(s S) int { return s.a },
			func(a int) S { return S{a} },
		),
	)

	assert.Equal(t, 1, sa.Get(S{1}))
	assert.Equal(t, S{2}, sa.Set(2)(S{1}))
}

func TestPassByValue(t *testing.T) {

	testLens := MakeLens(func(s Street) string { return s.name }, func(s Street, value string) Street {
		s.name = value
		return s
	})

	s1 := Street{1, "value1"}
	s2 := testLens.Set("value2")(s1)

	assert.Equal(t, "value1", s1.name)
	assert.Equal(t, "value2", s2.name)
}

func TestIdRef(t *testing.T) {
	idLens := IdRef[Street]()
	street := &Street{num: 1, name: "Main"}

	assert.Equal(t, street, idLens.Get(street))

	newStreet := &Street{num: 2, name: "Oak"}
	result := idLens.Set(newStreet)(street)
	assert.Equal(t, newStreet, result)
	assert.Equal(t, 1, street.num) // Original unchanged
}

func TestComposeRef(t *testing.T) {
	composedLens := ComposeRef[Address](streetLens)(addrLens)

	assert.Equal(t, sampleStreet.name, composedLens.Get(&sampleAddress))

	newName := "NewStreet"
	updated := composedLens.Set(newName)(&sampleAddress)
	assert.Equal(t, newName, composedLens.Get(updated))
	assert.Equal(t, sampleStreet.name, sampleAddress.street.name) // Original unchanged
}

func TestMakeLensCurried(t *testing.T) {
	nameLens := MakeLensCurried(
		func(s Street) string { return s.name },
		func(name string) func(Street) Street {
			return func(s Street) Street {
				s.name = name
				return s
			}
		},
	)

	street := Street{num: 1, name: "Main"}
	assert.Equal(t, "Main", nameLens.Get(street))

	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name)
}

func TestMakeLensRefCurried(t *testing.T) {
	nameLens := MakeLensRefCurried(
		func(s *Street) string { return s.name },
		func(name string) func(*Street) *Street {
			return func(s *Street) *Street {
				s.name = name
				return s
			}
		},
	)

	street := &Street{num: 1, name: "Main"}
	assert.Equal(t, "Main", nameLens.Get(street))

	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name)
}

func TestMakeLensWithEq(t *testing.T) {
	// Create a lens with equality check for string
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}

	// Test getting value
	assert.Equal(t, "Main", nameLens.Get(street))

	// Test setting a different value - should create a new copy
	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name) // Original unchanged
	assert.NotSame(t, street, updated)   // Different pointers

	// Test setting the same value - should return original pointer (optimization)
	same := nameLens.Set("Main")(street)
	assert.Equal(t, "Main", same.name)
	assert.Same(t, street, same) // Same pointer (no copy made)
}

func TestMakeLensWithEq_IntField(t *testing.T) {
	// Create a lens with equality check for int
	numLens := MakeLensWithEq(
		EQ.FromStrictEquals[int](),
		func(s *Street) int { return s.num },
		func(s *Street, num int) *Street {
			s.num = num
			return s
		},
	)

	street := &Street{num: 42, name: "Main"}

	// Test getting value
	assert.Equal(t, 42, numLens.Get(street))

	// Test setting a different value
	updated := numLens.Set(100)(street)
	assert.Equal(t, 100, updated.num)
	assert.Equal(t, 42, street.num)
	assert.NotSame(t, street, updated)

	// Test setting the same value - should return original pointer
	same := numLens.Set(42)(street)
	assert.Equal(t, 42, same.num)
	assert.Same(t, street, same)
}

func TestMakeLensWithEq_CustomEq(t *testing.T) {
	// Create a custom equality that ignores case
	caseInsensitiveEq := EQ.FromEquals(func(a, b string) bool {
		return len(a) == len(b) && a == b // Simple equality for this test
	})

	nameLens := MakeLensWithEq(
		caseInsensitiveEq,
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}

	// Setting the exact same value should return original pointer
	same := nameLens.Set("Main")(street)
	assert.Same(t, street, same)

	// Setting a different value should create a copy
	updated := nameLens.Set("Oak")(street)
	assert.NotSame(t, street, updated)
	assert.Equal(t, "Main", street.name)
	assert.Equal(t, "Oak", updated.name)
}

func TestMakeLensWithEq_ComposedLens(t *testing.T) {
	// Create lenses with equality optimization
	streetLensEq := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		(*Street).GetName,
		(*Street).SetName,
	)

	addrLensEq := MakeLensWithEq(
		EQ.FromStrictEquals[*Street](),
		(*Address).GetStreet,
		(*Address).SetStreet,
	)

	// Compose the lenses
	streetName := Compose[*Address](streetLensEq)(addrLensEq)

	sampleStreet := Street{num: 220, name: "Schönaicherstr"}
	sampleAddress := Address{city: "Böblingen", street: &sampleStreet}

	// Test getting value
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))

	// Test setting a different value
	newName := "Böblingerstr"
	updated := streetName.Set(newName)(&sampleAddress)
	assert.Equal(t, newName, streetName.Get(updated))
	assert.Equal(t, sampleStreet.name, sampleAddress.street.name) // Original unchanged
}

func TestMakeLensWithEq_MultipleUpdates(t *testing.T) {
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}

	// First update - creates a copy
	updated1 := nameLens.Set("Oak")(street)
	assert.NotSame(t, street, updated1)
	assert.Equal(t, "Oak", updated1.name)

	// Second update with same value - returns same pointer
	updated2 := nameLens.Set("Oak")(updated1)
	assert.Same(t, updated1, updated2)

	// Third update with different value - creates new copy
	updated3 := nameLens.Set("Elm")(updated2)
	assert.NotSame(t, updated2, updated3)
	assert.Equal(t, "Elm", updated3.name)
	assert.Equal(t, "Oak", updated2.name)
}

func TestMakeLensStrict(t *testing.T) {
	// Create a lens with strict equality for string
	nameLens := MakeLensStrict(
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}

	// Test getting value
	assert.Equal(t, "Main", nameLens.Get(street))

	// Test setting a different value - should create a new copy
	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name) // Original unchanged
	assert.NotSame(t, street, updated)   // Different pointers

	// Test setting the same value - should return original pointer (optimization)
	same := nameLens.Set("Main")(street)
	assert.Equal(t, "Main", same.name)
	assert.Same(t, street, same) // Same pointer (no copy made)
}

func TestMakeLensStrict_IntField(t *testing.T) {
	// Create a lens with strict equality for int
	numLens := MakeLensStrict(
		func(s *Street) int { return s.num },
		func(s *Street, num int) *Street {
			s.num = num
			return s
		},
	)

	street := &Street{num: 42, name: "Main"}

	// Test getting value
	assert.Equal(t, 42, numLens.Get(street))

	// Test setting a different value
	updated := numLens.Set(100)(street)
	assert.Equal(t, 100, updated.num)
	assert.Equal(t, 42, street.num)
	assert.NotSame(t, street, updated)

	// Test setting the same value - should return original pointer
	same := numLens.Set(42)(street)
	assert.Equal(t, 42, same.num)
	assert.Same(t, street, same)
}

func TestMakeLensStrict_PointerField(t *testing.T) {
	// Test with pointer field (comparable type)
	type Container struct {
		ptr *int
	}

	ptrLens := MakeLensStrict(
		func(c *Container) *int { return c.ptr },
		func(c *Container, ptr *int) *Container {
			c.ptr = ptr
			return c
		},
	)

	value1 := 42
	value2 := 100
	container := &Container{ptr: &value1}

	// Test getting value
	assert.Equal(t, &value1, ptrLens.Get(container))

	// Test setting a different pointer
	updated := ptrLens.Set(&value2)(container)
	assert.Equal(t, &value2, updated.ptr)
	assert.Equal(t, &value1, container.ptr)
	assert.NotSame(t, container, updated)

	// Test setting the same pointer - should return original
	same := ptrLens.Set(&value1)(container)
	assert.Same(t, container, same)
}

func TestMakeLensStrict_ComposedLens(t *testing.T) {
	// Create lenses with strict equality optimization
	streetLensStrict := MakeLensStrict(
		(*Street).GetName,
		(*Street).SetName,
	)

	addrLensStrict := MakeLensStrict(
		(*Address).GetStreet,
		(*Address).SetStreet,
	)

	// Compose the lenses
	streetName := Compose[*Address](streetLensStrict)(addrLensStrict)

	sampleStreet := Street{num: 220, name: "Schönaicherstr"}
	sampleAddress := Address{city: "Böblingen", street: &sampleStreet}

	// Test getting value
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))

	// Test setting a different value
	newName := "Böblingerstr"
	updated := streetName.Set(newName)(&sampleAddress)
	assert.Equal(t, newName, streetName.Get(updated))
	assert.Equal(t, sampleStreet.name, sampleAddress.street.name) // Original unchanged
}

func TestMakeLensStrict_MultipleUpdates(t *testing.T) {
	nameLens := MakeLensStrict(
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}

	// First update - creates a copy
	updated1 := nameLens.Set("Oak")(street)
	assert.NotSame(t, street, updated1)
	assert.Equal(t, "Oak", updated1.name)

	// Second update with same value - returns same pointer
	updated2 := nameLens.Set("Oak")(updated1)
	assert.Same(t, updated1, updated2)

	// Third update with different value - creates new copy
	updated3 := nameLens.Set("Elm")(updated2)
	assert.NotSame(t, updated2, updated3)
	assert.Equal(t, "Elm", updated3.name)
	assert.Equal(t, "Oak", updated2.name)
}

func TestMakeLensStrict_BoolField(t *testing.T) {
	type Config struct {
		enabled bool
	}

	enabledLens := MakeLensStrict(
		func(c *Config) bool { return c.enabled },
		func(c *Config, enabled bool) *Config {
			c.enabled = enabled
			return c
		},
	)

	config := &Config{enabled: true}

	// Test getting value
	assert.True(t, enabledLens.Get(config))

	// Test setting a different value
	updated := enabledLens.Set(false)(config)
	assert.False(t, updated.enabled)
	assert.True(t, config.enabled)
	assert.NotSame(t, config, updated)

	// Test setting the same value - should return original pointer
	same := enabledLens.Set(true)(config)
	assert.Same(t, config, same)
}

func TestMakeLensRef_WithNilState(t *testing.T) {
	// Test that MakeLensRef creates a total lens that works with nil pointers
	nameLens := MakeLensRef(
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	// Test Get with nil - should handle gracefully
	var nilStreet *Street = nil
	name := nameLens.Get(nilStreet)
	assert.Equal(t, "", name)

	// Test Set with nil - should create a new object with zero values except the set field
	updated := nameLens.Set("NewStreet")(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, "NewStreet", updated.name)
	assert.Equal(t, 0, updated.num) // Zero value for int

	// Verify original nil pointer is unchanged
	assert.Nil(t, nilStreet)
}

func TestMakeLensRef_WithNilState_IntField(t *testing.T) {
	// Test with an int field lens
	numLens := MakeLensRef(
		func(s *Street) int {
			if s == nil {
				return 0
			}
			return s.num
		},
		func(s *Street, num int) *Street {
			s.num = num
			return s
		},
	)

	var nilStreet *Street = nil

	// Get from nil should return zero value
	num := numLens.Get(nilStreet)
	assert.Equal(t, 0, num)

	// Set on nil should create new object
	updated := numLens.Set(42)(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, 42, updated.num)
	assert.Equal(t, "", updated.name) // Zero value for string
}

func TestMakeLensRef_WithNilState_Composed(t *testing.T) {
	// Test composed lenses with nil state
	streetLens := MakeLensRef(
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		(*Street).SetName,
	)

	addrLens := MakeLensRef(
		func(a *Address) *Street {
			if a == nil {
				return nil
			}
			return a.street
		},
		(*Address).SetStreet,
	)

	// Compose the lenses
	streetName := ComposeRef[Address](streetLens)(addrLens)

	var nilAddress *Address = nil

	// Get from nil should handle gracefully
	name := streetName.Get(nilAddress)
	assert.Equal(t, "", name)

	// Set on nil should create new nested structure
	updated := streetName.Set("TestStreet")(nilAddress)
	assert.NotNil(t, updated)
	assert.NotNil(t, updated.street)
	assert.Equal(t, "TestStreet", updated.street.name)
	assert.Equal(t, "", updated.city) // Zero value for city
}

func TestMakeLensRefCurried_WithNilState(t *testing.T) {
	// Test that MakeLensRefCurried creates a total lens that works with nil pointers
	nameLens := MakeLensRefCurried(
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(name string) func(*Street) *Street {
			return func(s *Street) *Street {
				s.name = name
				return s
			}
		},
	)

	// Test Get with nil
	var nilStreet *Street = nil
	name := nameLens.Get(nilStreet)
	assert.Equal(t, "", name)

	// Test Set with nil - should create a new object
	updated := nameLens.Set("CurriedStreet")(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, "CurriedStreet", updated.name)
	assert.Equal(t, 0, updated.num) // Zero value for int

	// Verify original nil pointer is unchanged
	assert.Nil(t, nilStreet)
}

func TestMakeLensRefCurried_WithNilState_IntField(t *testing.T) {
	// Test with an int field lens using curried setter
	numLens := MakeLensRefCurried(
		func(s *Street) int {
			if s == nil {
				return 0
			}
			return s.num
		},
		func(num int) func(*Street) *Street {
			return func(s *Street) *Street {
				s.num = num
				return s
			}
		},
	)

	var nilStreet *Street = nil

	// Get from nil should return zero value
	num := numLens.Get(nilStreet)
	assert.Equal(t, 0, num)

	// Set on nil should create new object
	updated := numLens.Set(99)(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, 99, updated.num)
	assert.Equal(t, "", updated.name) // Zero value for string
}

func TestMakeLensRefCurried_WithNilState_MultipleOperations(t *testing.T) {
	// Test multiple operations on nil and non-nil states
	nameLens := MakeLensRefCurried(
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(name string) func(*Street) *Street {
			return func(s *Street) *Street {
				s.name = name
				return s
			}
		},
	)

	var nilStreet *Street = nil

	// First operation on nil
	street1 := nameLens.Set("First")(nilStreet)
	assert.NotNil(t, street1)
	assert.Equal(t, "First", street1.name)

	// Second operation on non-nil result
	street2 := nameLens.Set("Second")(street1)
	assert.NotNil(t, street2)
	assert.Equal(t, "Second", street2.name)
	assert.Equal(t, "First", street1.name) // Original unchanged

	// Third operation back to nil (edge case)
	street3 := nameLens.Set("Third")(nilStreet)
	assert.NotNil(t, street3)
	assert.Equal(t, "Third", street3.name)
}

func TestMakeLensRef_WithNilState_NestedStructure(t *testing.T) {
	// Test with nested structure where inner can be nil
	innerLens := MakeLensRef(
		func(o *Outer) *Inner {
			if o == nil {
				return nil
			}
			return o.inner
		},
		func(o *Outer, i *Inner) *Outer {
			o.inner = i
			return o
		},
	)

	var nilOuter *Outer = nil

	// Get from nil outer
	inner := innerLens.Get(nilOuter)
	assert.Nil(t, inner)

	// Set on nil outer
	newInner := &Inner{Value: 42, Foo: "test"}
	updated := innerLens.Set(newInner)(nilOuter)
	assert.NotNil(t, updated)
	assert.Equal(t, newInner, updated.inner)
}

func TestMakeLensWithEq_WithNilState(t *testing.T) {
	// Test that MakeLensWithEq creates a total lens that works with nil pointers
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	// Test Get with nil - should handle gracefully
	var nilStreet *Street = nil
	name := nameLens.Get(nilStreet)
	assert.Equal(t, "", name)

	// Test Set with nil - should create a new object with zero values except the set field
	updated := nameLens.Set("NewStreet")(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, "NewStreet", updated.name)
	assert.Equal(t, 0, updated.num) // Zero value for int

	// Verify original nil pointer is unchanged
	assert.Nil(t, nilStreet)
}

func TestMakeLensWithEq_WithNilState_IntField(t *testing.T) {
	// Test with an int field lens with equality optimization
	numLens := MakeLensWithEq(
		EQ.FromStrictEquals[int](),
		func(s *Street) int {
			if s == nil {
				return 0
			}
			return s.num
		},
		func(s *Street, num int) *Street {
			s.num = num
			return s
		},
	)

	var nilStreet *Street = nil

	// Get from nil should return zero value
	num := numLens.Get(nilStreet)
	assert.Equal(t, 0, num)

	// Set on nil should create new object
	updated := numLens.Set(42)(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, 42, updated.num)
	assert.Equal(t, "", updated.name) // Zero value for string
}

func TestMakeLensWithEq_WithNilState_EqualityOptimization(t *testing.T) {
	// Test that equality optimization works with nil state
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	var nilStreet *Street = nil

	// Setting empty string on nil should return a new object with empty string
	// (since the zero value equals the set value)
	updated1 := nameLens.Set("")(nilStreet)
	assert.NotNil(t, updated1)
	assert.Equal(t, "", updated1.name)

	// Setting the same empty string again should return the same pointer (optimization)
	updated2 := nameLens.Set("")(updated1)
	assert.Same(t, updated1, updated2)

	// Setting a different value should create a new copy
	updated3 := nameLens.Set("Different")(updated1)
	assert.NotSame(t, updated1, updated3)
	assert.Equal(t, "Different", updated3.name)
	assert.Equal(t, "", updated1.name)
}

func TestMakeLensWithEq_WithNilState_CustomEq(t *testing.T) {
	// Test with custom equality predicate on nil state
	customEq := EQ.FromEquals(func(a, b string) bool {
		return len(a) == len(b) && a == b
	})

	nameLens := MakeLensWithEq(
		customEq,
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	var nilStreet *Street = nil

	// Get from nil
	name := nameLens.Get(nilStreet)
	assert.Equal(t, "", name)

	// Set on nil with non-empty string
	updated := nameLens.Set("Test")(nilStreet)
	assert.NotNil(t, updated)
	assert.Equal(t, "Test", updated.name)

	// Set same value should return same pointer
	same := nameLens.Set("Test")(updated)
	assert.Same(t, updated, same)
}

func TestMakeLensWithEq_WithNilState_MultipleOperations(t *testing.T) {
	// Test multiple operations on nil and non-nil states with equality optimization
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string {
			if s == nil {
				return ""
			}
			return s.name
		},
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	var nilStreet *Street = nil

	// First operation on nil
	street1 := nameLens.Set("First")(nilStreet)
	assert.NotNil(t, street1)
	assert.Equal(t, "First", street1.name)

	// Second operation with same value - should return same pointer
	street2 := nameLens.Set("First")(street1)
	assert.Same(t, street1, street2)

	// Third operation with different value - should create new copy
	street3 := nameLens.Set("Second")(street2)
	assert.NotSame(t, street2, street3)
	assert.Equal(t, "Second", street3.name)
	assert.Equal(t, "First", street2.name)

	// Fourth operation back to nil with zero value
	street4 := nameLens.Set("")(nilStreet)
	assert.NotNil(t, street4)
	assert.Equal(t, "", street4.name)
}
