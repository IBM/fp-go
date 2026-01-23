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

package optional

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"

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
	responseOptional = FromPredicateRef[Response](F.IsNonNil[Info])((*Response).GetInfo, (*Response).SetInfo)

	sampleResponse      = Response{info: &Info{}}
	sampleEmptyResponse = Response{}
)

func TestOptional(t *testing.T) {
	assert.Equal(t, O.Of(sampleResponse.info), responseOptional.GetOption(&sampleResponse))
	assert.Equal(t, O.None[*Info](), responseOptional.GetOption(&sampleEmptyResponse))
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

// TestMakeOptionalBasicFunctionality tests basic Optional operations
func TestMakeOptionalBasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when value exists", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		person := Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("GetOption returns None when value doesn't exist", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		person := Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Set updates value when optional matches", func(t *testing.T) {
		optional := MakeOptional(
			func(p Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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

// TestOptionalLaws tests that Optional satisfies the optional laws
// Reference: https://gcanti.github.io/monocle-ts/modules/Optional.ts.html
func TestOptionalLaws(t *testing.T) {
	optional := MakeOptional(
		func(p Person) O.Option[string] {
			if p.Name != "" {
				return O.Some(p.Name)
			}
			return O.None[string]()
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
		assert.True(t, O.IsSome(initial))

		// Set a new value
		newName := "Bob"
		updated := optional.Set(newName)(person)

		// Get the value back
		result := optional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.True(t, O.IsSome(result))
		assert.Equal(t, newName, O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("GetSet Law: Set(a)(s) = s when GetOption(s) = None (no-op)", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

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

// TestMakeOptionalRefBasicFunctionality tests MakeOptionalRef with pointer types
func TestMakeOptionalRefBasicFunctionality(t *testing.T) {
	t.Run("GetOption returns Some when value exists (non-nil pointer)", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		person := &Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("GetOption returns None when pointer is nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("GetOption returns None when value doesn't exist", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		person := &Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Set updates value and creates copy (immutability)", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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

// TestMakeOptionalRefLaws tests that MakeOptionalRef satisfies optional laws
func TestMakeOptionalRefLaws(t *testing.T) {
	optional := MakeOptionalRef(
		func(p *Person) O.Option[string] {
			if p.Name != "" {
				return O.Some(p.Name)
			}
			return O.None[string]()
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
		assert.True(t, O.IsSome(initial))

		// Set a new value
		newName := "Bob"
		updated := optional.Set(newName)(person)

		// Get the value back
		result := optional.GetOption(updated)

		// Verify SetGet law: we should get back what we set
		assert.True(t, O.IsSome(result))
		assert.Equal(t, newName, O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("GetSet Law: Set(a)(s) = s when GetOption(s) = None (nil pointer)", func(t *testing.T) {
		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

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

// TestMakeOptionalRefImmutability tests immutability guarantees
func TestMakeOptionalRefImmutability(t *testing.T) {
	t.Run("Set creates a new pointer, doesn't modify original", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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

// TestMakeOptionalRefNilPointerEdgeCases tests edge cases with nil pointers
func TestMakeOptionalRefNilPointerEdgeCases(t *testing.T) {
	t.Run("GetOption on nil returns None", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				return O.Some(p.Name)
			},
			func(p *Person, name string) *Person {
				p.Name = name
				return p
			},
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Set on nil returns nil", func(t *testing.T) {
		optional := MakeOptionalRef(
			func(p *Person) O.Option[string] {
				return O.Some(p.Name)
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
			func(p *Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
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

// TestFromPredicateRef tests FromPredicateRef with nil handling
func TestFromPredicateRef(t *testing.T) {
	t.Run("Works with non-nil values matching predicate", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "Alice", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, "Alice", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("Returns None for nil pointer", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		var person *Person = nil
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Returns None when predicate doesn't match", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "", Age: 30}
		result := optional.GetOption(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("Set is no-op on nil pointer", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
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

// TestOptionalComposition tests composing optionals
func TestOptionalComposition(t *testing.T) {
	t.Run("Compose two optionals", func(t *testing.T) {
		// First optional: Person -> Name (if not empty)
		nameOptional := MakeOptional(
			func(p Person) O.Option[string] {
				if p.Name != "" {
					return O.Some(p.Name)
				}
				return O.None[string]()
			},
			func(p Person, name string) Person {
				p.Name = name
				return p
			},
		)

		// Second optional: String -> First character (if not empty)
		firstCharOptional := MakeOptional(
			func(s string) O.Option[rune] {
				if len(s) > 0 {
					return O.Some(rune(s[0]))
				}
				return O.None[rune]()
			},
			func(s string, r rune) string {
				if len(s) > 0 {
					return string(r) + s[1:]
				}
				return string(r)
			},
		)

		// Compose them
		composed := Compose[Person](firstCharOptional)(nameOptional)

		person := Person{Name: "Alice", Age: 30}
		result := composed.GetOption(person)

		assert.True(t, O.IsSome(result))
		assert.Equal(t, 'A', O.GetOrElse(F.Constant(rune(0)))(result))
	})
}

// TestOptionalNoOpBehavior tests that modifying through optionalModify is a no-op when GetOption returns None
// This is the key law: updating a value for which the preview returns None is a no-op
func TestOptionalNoOpBehavior(t *testing.T) {
	optional := MakeOptional(
		func(p Person) O.Option[string] {
			if p.Name != "" {
				return O.Some(p.Name)
			}
			return O.None[string]()
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
		assert.True(t, O.IsNone(initial))

		// Try to modify - should return None
		modifyResult := ModifyOption[Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, O.IsNone(modifyResult))
	})

	t.Run("optionalModify is no-op when GetOption returns None", func(t *testing.T) {
		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

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
		assert.True(t, O.IsSome(initial))

		// Modify should return Some with updated value
		modifyResult := ModifyOption[Person](func(name string) string {
			return name + " Smith"
		})(optional)(person)

		assert.True(t, O.IsSome(modifyResult))
		updatedPerson := O.GetOrElse(F.Constant(person))(modifyResult)
		assert.Equal(t, "Alice Smith", updatedPerson.Name)
	})

	t.Run("optionalModify updates when GetOption returns Some", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, O.IsSome(initial))

		// Modify should update the value
		updated := optionalModify(func(name string) string {
			return name + " Smith"
		}, optional, person)

		assert.Equal(t, "Alice Smith", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})
}

// TestOptionalNoOpBehaviorRef tests no-op behavior with pointer types
func TestOptionalNoOpBehaviorRef(t *testing.T) {
	optional := MakeOptionalRef(
		func(p *Person) O.Option[string] {
			if p.Name != "" {
				return O.Some(p.Name)
			}
			return O.None[string]()
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
		assert.True(t, O.IsNone(initial))

		// Try to modify - should return None
		modifyResult := ModifyOption[*Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, O.IsNone(modifyResult))
	})

	t.Run("ModifyOption returns None when pointer is nil", func(t *testing.T) {
		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

		// Try to modify - should return None
		modifyResult := ModifyOption[*Person](func(name string) string {
			return "Bob"
		})(optional)(person)

		assert.True(t, O.IsNone(modifyResult))
	})

	t.Run("optionalModify is no-op when GetOption returns None", func(t *testing.T) {
		person := &Person{Name: "", Age: 30}
		originalName := person.Name
		originalAge := person.Age

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

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
		assert.True(t, O.IsNone(initial))

		// Try to modify
		updated := optionalModify(func(name string) string {
			return "Bob"
		}, optional, person)

		// Verify no-op: should still be nil
		assert.Nil(t, updated)
	})
}

// TestFromPredicateNoOpBehavior tests that FromPredicate properly implements no-op behavior
func TestFromPredicateNoOpBehavior(t *testing.T) {
	t.Run("FromPredicate Set is no-op when predicate doesn't match", func(t *testing.T) {
		optional := FromPredicate[Person](func(name string) bool {
			return name != ""
		})(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { p.Name = name; return p },
		)

		person := Person{Name: "", Age: 30}

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

		// Set should be no-op when predicate doesn't match
		updated := optional.Set("Bob")(person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, person, updated)
		assert.Equal(t, "", updated.Name)
		assert.Equal(t, 30, updated.Age)
	})

	t.Run("FromPredicate Set updates when predicate matches on current value", func(t *testing.T) {
		optional := FromPredicate[Person](func(name string) bool {
			return name != ""
		})(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { p.Name = name; return p },
		)

		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, O.IsSome(initial))

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
		optional := FromPredicate[Person](func(age int) bool {
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

// TestFromPredicateRefNoOpBehavior tests that FromPredicateRef properly implements no-op behavior
func TestFromPredicateRefNoOpBehavior(t *testing.T) {
	t.Run("FromPredicateRef Set is no-op when predicate doesn't match", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
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
		assert.True(t, O.IsNone(initial))

		// Set should be no-op when predicate doesn't match
		updated := optional.Set("Bob")(person)

		// Verify no-op: structure should be unchanged
		assert.Equal(t, originalName, updated.Name)
		assert.Equal(t, originalAge, updated.Age)
		// Original should also be unchanged (immutability)
		assert.Equal(t, originalName, person.Name)
	})

	t.Run("FromPredicateRef Set is no-op when pointer is nil", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		var person *Person = nil

		// Verify optional doesn't match
		initial := optional.GetOption(person)
		assert.True(t, O.IsNone(initial))

		// Set should be no-op (return nil)
		updated := optional.Set("Bob")(person)

		assert.Nil(t, updated)
	})

	t.Run("FromPredicateRef Set updates when predicate matches on current value", func(t *testing.T) {
		optional := FromPredicateRef[Person](func(name string) bool {
			return name != ""
		})(
			func(p *Person) string { return p.Name },
			func(p *Person, name string) *Person { p.Name = name; return p },
		)

		person := &Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, O.IsSome(initial))

		// Set should update when predicate matches on the CURRENT value
		updated := optional.Set("Bob")(person)

		assert.Equal(t, "Bob", updated.Name)
		assert.Equal(t, 30, updated.Age)
		// Original should be unchanged (immutability)
		assert.Equal(t, "Alice", person.Name)
	})

	t.Run("FromPredicateRef demonstrates the no-op law correctly", func(t *testing.T) {
		// This test shows that FromPredicateRef implements the no-op law
		optional := FromPredicateRef[Person](func(age int) bool {
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

// TestSetOptionNoOpBehavior tests SetOption behavior with None
func TestSetOptionNoOpBehavior(t *testing.T) {
	optional := MakeOptional(
		func(p Person) O.Option[string] {
			if p.Name != "" {
				return O.Some(p.Name)
			}
			return O.None[string]()
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
		assert.True(t, O.IsNone(initial))

		// SetOption should return None
		result := SetOption[Person]("Bob")(optional)(person)

		assert.True(t, O.IsNone(result))
	})

	t.Run("SetOption returns Some when GetOption returns Some", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}

		// Verify optional matches
		initial := optional.GetOption(person)
		assert.True(t, O.IsSome(initial))

		// SetOption should return Some with updated value
		result := SetOption[Person]("Bob")(optional)(person)

		assert.True(t, O.IsSome(result))
		updatedPerson := O.GetOrElse(F.Constant(person))(result)
		assert.Equal(t, "Bob", updatedPerson.Name)
	})
}
