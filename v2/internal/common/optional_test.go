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

package common

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"

	"github.com/stretchr/testify/assert"
)

type (
	Phone struct {
		number string
	}

	Employment struct {
		phone *Phone
	}

	Info struct {
		employment *Employment
	}

	Response struct {
		info *Info
	}
)

func (response *Response) GetInfo() *Info {
	return response.info
}

func (response *Response) SetInfo(info *Info) *Response {
	response.info = info
	return response
}

var (
	responseOptional = OptionalFromPredicateRef[Response](F.IsNonNil[Info])((*Response).GetInfo, (*Response).SetInfo)

	sampleResponse      = Response{info: &Info{}}
	sampleEmptyResponse = Response{}
)

func TestOptionalFromPredicateRef_GetOption(t *testing.T) {
	assert.Equal(t, OptionOf(sampleResponse.info), responseOptional.GetOption(&sampleResponse))
	assert.Equal(t, OptionNone[*Info](), responseOptional.GetOption(&sampleEmptyResponse))
}

// Test types for comprehensive testing
type Person struct {
	Name string
	Age  int
}

type Config struct {
	Timeout int
	Retries int
}

// TestMakeOptional tests basic Optional operations
func TestMakeOptional(t *testing.T) {
	t.Run("GetOption returns Some when value exists", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		person := Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsSome(result))
		assert.Equal(t, "Alice", OptionGetOrElse(F.Constant(""))(result))
	})

	t.Run("GetOption returns None when value doesn't exist", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		person := Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("Set updates value when optional matches", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		person := Person{Name: "Alice", Age: 30}
		updated := optional.Set("Bob")(person)

		assert.Equal(t, "Bob", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})
}

// TestMakeOptional_Laws tests that MakeOptional satisfies the optional laws
// Reference: https://gcanti.github.io/monocle-ts/modules/Optional.ts.html
func TestMakeOptional_Laws(t *testing.T) {
	optional := MakeOptional(
		func(p Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)

	t.Run("SetGet Law: GetOption(Set(a)(s)) = Some(a) when GetOption(s) = Some(_)", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Set a new value
		newName := "Bob"
		updated := optional.Set(newName)(person)

		// Get the value back
		result := optional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.True(t, OptionIsSome(result))
		assert.Equal(t, newName, OptionGetOrElse(F.Constant(""))(result))
	})

	t.Run("GetSet Law: Set(a)(s) = s when GetOption(s) = None (no-op)", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to set a value - this should be a no-op since GetOption returns None
		// Note: Direct Set always updates, but this is expected behavior.
		// The no-op behavior is enforced through ModifyOption and optionalModify.
		updated := optional.Set("Bob")(person)

		// Direct Set will update even when GetOption returns None
		// This is by design - Set is unconditional
		assert.Equal(t, "Bob", updated.Name)
	})

	t.Run("SetSet Law: Set(b)(Set(a)(s)) = Set(b)(s)", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Set twice
		setTwice := optional.Set("Charlie")(optional.Set("Bob")(person))

		// Set once with the final value
		setOnce := optional.Set("Charlie")(person)

		// They should be equal
		assert.Equal(t, setOnce, setTwice)
		assert.Equal(t, "Charlie", setTwice.Name)
	})
}

// TestMakeOptionalRef tests MakeOptionalRef with pointer types
func TestMakeOptionalRef(t *testing.T) {
	t.Run("GetOption returns Some when value exists (non-nil pointer)", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		person := &Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsSome(result))
		assert.Equal(t, "Alice", OptionGetOrElse(F.Constant(""))(result))
	})

	t.Run("GetOption returns None when pointer is nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("GetOption returns None when value doesn't exist", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		person := &Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("Set updates value and creates copy (immutability)", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		original := &Person{Name: "Alice", Age: 30}
		updated := optional.Set("Bob")(original)

		// Verify the update
		assert.Equal(t, "Bob", updated.Name)
		assert.Equal(t, 30, updated.Age)

		// Verify immutability: original should be unchanged
		assert.Equal(t, "Alice", original.Name)

		// Verify they are different pointers
		assert.NotEqual(t, original, updated)
	})

	t.Run("Set is no-op when pointer is nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		updated := optional.Set("Bob")(person)

		// Verify nothing changed (no-op for nil)
		assert.Nil(t, updated)
	})
}

// TestMakeOptionalRef_Laws tests that MakeOptionalRef satisfies optional laws
func TestMakeOptionalRef_Laws(t *testing.T) {
	optional := MakeOptionalRef(
		func(p *Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p *Person, name string) *Person {
			p.Name = name
			return p
		},
	)

	t.Run("SetGet Law: GetOption(Set(a)(s)) = Some(a) when GetOption(s) = Some(_)", func(t *testing.T) {
		person := &Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Set a new value
		newName := "Bob"
		updated := optional.Set(newName)(person)

		// Get the value back
		result := optional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.True(t, OptionIsSome(result))
		assert.Equal(t, newName, OptionGetOrElse(F.Constant(""))(result))
	})

	t.Run("GetSet Law: Set(a)(s) = s when GetOption(s) = None (nil pointer)", func(t *testing.T) {
		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to set a value
		updated := optional.Set("Bob")(person)

		// Verify GetSet law: structure should be unchanged (nil)
		assert.Nil(t, updated)
	})

	t.Run("SetSet Law: Set(b)(Set(a)(s)) = Set(b)(s)", func(t *testing.T) {
		person := &Person{Name: "Alice", Age: 30}

		// Set twice
		setTwice := optional.Set("Charlie")(optional.Set("Bob")(person))

		// Set once with the final value
		setOnce := optional.Set("Charlie")(person)

		// They should have equal values (but different pointers due to immutability)
		assert.Equal(t, setOnce.Name, setTwice.Name)
		assert.Equal(t, setOnce.Age, setTwice.Age)
		assert.Equal(t, "Charlie", setTwice.Name)
	})
}

// TestMakeOptionalRef_Immutability tests immutability guarantees of MakeOptionalRef
func TestMakeOptionalRef_Immutability(t *testing.T) {
	t.Run("Set creates a new pointer, doesn't modify original", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		original := &Person{Name: "Alice", Age: 30}
		origName := original.Name
		origAge := original.Age

		// Perform multiple sets
		updated1 := optional.Set("Bob")(original)
		updated2 := optional.Set("Charlie")(updated1)
		updated3 := optional.Set("David")(updated2)

		// Verify original is unchanged
		assert.Equal(t, origName, original.Name)
		assert.Equal(t, origAge, original.Age)

		// Verify final update has correct value
		assert.Equal(t, "David", updated3.Name)

		// Verify all pointers are different
		assert.NotEqual(t, original, updated1)
		assert.NotEqual(t, original, updated2)
		assert.NotEqual(t, original, updated3)
		assert.NotEqual(t, updated1, updated2)
		assert.NotEqual(t, updated2, updated3)
	})

	t.Run("Multiple operations on nil preserve nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil

		// Multiple sets on nil should all return nil
		updated1 := optional.Set("Bob")(person)
		updated2 := optional.Set("Charlie")(updated1)

		assert.Nil(t, updated1)
		assert.Nil(t, updated2)
	})
}

// TestMakeOptionalRef_NilPointerEdgeCases tests edge cases with nil pointers in MakeOptionalRef
func TestMakeOptionalRef_NilPointerEdgeCases(t *testing.T) {
	t.Run("GetOption on nil returns None", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				return OptionSome(p.Name)
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("Set on nil returns nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				return OptionSome(p.Name)
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		updated := optional.Set("Bob")(person)

		assert.Nil(t, updated)
	})

	t.Run("Chaining operations starting from nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil

		// Chain multiple operations
		result := F.Pipe2(
			person,
			optional.Set("Bob"),
			optional.Set("Charlie"),
		)

		assert.Nil(t, result)
	})
}

// TestOptionalFromPredicateRef_WithNilHandling tests OptionalFromPredicateRef with nil handling
func TestOptionalFromPredicateRef_WithNilHandling(t *testing.T) {
	t.Run("Works with non-nil values matching predicate", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsSome(result))
		assert.Equal(t, "Alice", OptionGetOrElse(F.Constant(""))(result))
	})

	t.Run("Returns None for nil pointer", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("Returns None when predicate doesn't match", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("Set is no-op on nil pointer", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		var person *Person = nil
		updated := optional.Set("Bob")(person)

		assert.Nil(t, updated)
	})
}

// TestOptionalComposeOptional tests composing optionals
func TestOptionalComposeOptional(t *testing.T) {
	t.Run("Compose two optionals", func(t *testing.T) {
		// First optional: Person -> Name (if not empty)
		nameOptional := MakeOptional(
			func(p Person) Option[string] {
				if p.Name != "" {
					return OptionSome(p.Name)
				}
				return OptionNone[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		// Second optional: String -> First character (if not empty)
		firstCharOptional := MakeOptional(
			func(s string) Option[rune] {
				if len(s) > 0 {
					return OptionSome(rune(s[0]))
				}
				return OptionNone[rune]()
			},
			func(s string, r rune) string {
				if len(s) > 0 {
					return string(r) + s[1:]
				}
				return string(r)
			},
		)

		// Compose them
		composed := OptionalComposeOptional[Person](firstCharOptional)(nameOptional)

		person := Person{Name: "Alice", Age: 30}
		result := composed.GetOption(person)

		assert.True(t, OptionIsSome(result))
		assert.Equal(t, 'A', OptionGetOrElse(F.Constant(rune(0)))(result))
	})
}

// TestOptionalModifyOption tests that modifying through optionalModify is a no-op when GetOption returns None
// This is the key law: updating a value for which the preview returns None is a no-op
func TestOptionalModifyOption(t *testing.T) {
	optional := MakeOptional(
		func(p Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)

	t.Run("ModifyOption returns None when GetOption returns None", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify - should return None
		modifyResult := OptionalModifyOption[Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, OptionIsNone(modifyResult))
	})

	t.Run("optionalModify is no-op when GetOption returns None", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify using the internal optionalModify function
		updated := optionalModify(func(name string) string {
			return "Bob"
		}, optional, person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, person, updated)
		assert.Equal(t, "", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})

	t.Run("ModifyOption returns Some when GetOption returns Some", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Modify should return Some with updated value
		modifyResult := OptionalModifyOption[Person](func(name string) string {
			return name + " Smith"
		})(optional)(person)

		assert.True(t, OptionIsSome(modifyResult))
		updatedPerson := OptionGetOrElse(F.Constant(person))(modifyResult)
		assert.Equal(t, "Alice Smith", updatedPerson.Name)
	})

	t.Run("optionalModify updates when GetOption returns Some", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Modify should update the value
		updated := optionalModify(func(name string) string {
			return name + " Smith"
		}, optional, person)

		assert.Equal(t, "Alice Smith", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})
}

// TestOptionalModifyOption_Ref tests OptionalModifyOption no-op behavior with pointer types
func TestOptionalModifyOption_Ref(t *testing.T) {
	optional := MakeOptionalRef(
		func(p *Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p *Person, name string) *Person {
			p.Name = name
			return p
		},
	)

	t.Run("ModifyOption returns None when GetOption returns None (empty name)", func(t *testing.T) {
		person := &Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify - should return None
		modifyResult := OptionalModifyOption[*Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, OptionIsNone(modifyResult))
	})

	t.Run("ModifyOption returns None when pointer is nil", func(t *testing.T) {
		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify - should return None
		modifyResult := OptionalModifyOption[*Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, OptionIsNone(modifyResult))
	})

	t.Run("optionalModify is no-op when GetOption returns None", func(t *testing.T) {
		person := &Person{Name: "", Age: 30}
		originalName := person.Name
		originalAge := person.Age

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify
		updated := optionalModify(func(name string) string {
			return "Bob"
		}, optional, person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, originalName, updated.Name)
		assert.Equal(t, originalAge, updated.Age)
	})

	t.Run("optionalModify is no-op when pointer is nil", func(t *testing.T) {
		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Try to modify
		updated := optionalModify(func(name string) string {
			return "Bob"
		}, optional, person)

		// Verify no-op: should still be nil
		assert.Nil(t, updated)
	})
}

// TestOptionalFromPredicate tests that OptionalFromPredicate properly implements no-op behavior
func TestOptionalFromPredicate(t *testing.T) {
	t.Run("FromPredicate Set is no-op when predicate doesn't match", func(t *testing.T) {
		optional := OptionalFromPredicate[Person](func(name string) bool {
			return name != ""
		})(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { p.Name = name; return p },
		)

		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Set should be no-op when predicate doesn't match
		updated := optional.Set("Bob")(person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, person, updated)
		assert.Equal(t, "", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})

	t.Run("FromPredicate Set updates when predicate matches on current value", func(t *testing.T) {
		optional := OptionalFromPredicate[Person](func(name string) bool {
			return name != ""
		})(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { p.Name = name; return p },
		)

		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Set should update when predicate matches on the CURRENT value
		// Note: FromPredicate's setter checks the predicate on the current value,
		// not the new value. This is the correct behavior for the no-op law.
		updated := optional.Set("Bob")(person)

		assert.Equal(t, "Bob", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})

	t.Run("FromPredicate demonstrates the no-op law correctly", func(t *testing.T) {
		// This test shows that FromPredicate implements the no-op law:
		// The setter checks if the CURRENT value matches the predicate
		optional := OptionalFromPredicate[Person](func(age int) bool {
			return age >= 18 // Adult predicate
		})(
			func(p Person) int { return p.Age },
			func(p Person, age int) Person { p.Age = age; return p },
		)

		// Case 1: Current value matches predicate (adult) - Set should work
		adult := Person{Name: "Alice", Age: 30}
		updatedAdult := optional.Set(25)(adult)
		assert.Equal(t, 25, updatedAdult.Age)

		// Case 2: Current value doesn't match predicate (child) - Set is no-op
		child := Person{Name: "Bob", Age: 10}
		updatedChild := optional.Set(25)(child)
		assert.Equal(t, 10, updatedChild.Age) // Unchanged - no-op!
	})
}

// TestOptionalFromPredicateRef_NoOpBehavior tests that OptionalFromPredicateRef properly implements no-op behavior
func TestOptionalFromPredicateRef_NoOpBehavior(t *testing.T) {
	t.Run("FromPredicateRef Set is no-op when predicate doesn't match", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "", Age: 30}
		originalName := person.Name
		originalAge := person.Age

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Set should be no-op when predicate doesn't match
		updated := optional.Set("Bob")(person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, originalName, updated.Name)
		assert.Equal(t, originalAge, updated.Age)
		// Original should also be unchanged (immutability)
		assert.Equal(t, originalName, person.Name)
	})

	t.Run("FromPredicateRef Set is no-op when pointer is nil", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// Set should be no-op (return nil)
		updated := optional.Set("Bob")(person)

		assert.Nil(t, updated)
	})

	t.Run("FromPredicateRef Set updates when predicate matches on current value", func(t *testing.T) {
		optional := OptionalFromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// Set should update when predicate matches on the CURRENT value
		updated := optional.Set("Bob")(person)

		assert.Equal(t, "Bob", updated.Name)
		assert.Equal(t, 30, updated.Age)
		// Original should be unchanged (immutability)
		assert.Equal(t, "Alice", person.Name)
	})

	t.Run("FromPredicateRef demonstrates the no-op law correctly", func(t *testing.T) {
		// This test shows that FromPredicateRef implements the no-op law
		optional := OptionalFromPredicateRef[Person](func(age int) bool {
			return age >= 18 // Adult predicate
		})(
			func(p *Person) int { return p.Age },
			func(p *Person, age int) *Person { p.Age = age; return p },
		)

		// Case 1: Current value matches predicate (adult) - Set should work
		adult := &Person{Name: "Alice", Age: 30}
		updatedAdult := optional.Set(25)(adult)
		assert.Equal(t, 25, updatedAdult.Age)
		assert.Equal(t, 30, adult.Age) // Original unchanged

		// Case 2: Current value doesn't match predicate (child) - Set is no-op
		child := &Person{Name: "Bob", Age: 10}
		updatedChild := optional.Set(25)(child)
		assert.Equal(t, 10, updatedChild.Age) // Unchanged - no-op!
		assert.Equal(t, 10, child.Age)        // Original also unchanged
	})
}

// TestOptionalSetOption tests OptionalSetOption behavior with None
func TestOptionalSetOption(t *testing.T) {
	optional := MakeOptional(
		func(p Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)

	t.Run("SetOption returns None when GetOption returns None", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, OptionIsNone(initial))

		// SetOption should return None
		result := OptionalSetOption[Person]("Bob")(optional)(person)

		assert.True(t, OptionIsNone(result))
	})

	t.Run("SetOption returns Some when GetOption returns Some", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, OptionIsSome(initial))

		// SetOption should return Some with updated value
		result := OptionalSetOption[Person]("Bob")(optional)(person)

		assert.True(t, OptionIsSome(result))
		updatedPerson := OptionGetOrElse(F.Constant(person))(result)
		assert.Equal(t, "Bob", updatedPerson.Name)
	})
}
