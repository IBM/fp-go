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

package codec

import (
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Test types for ApSL
type Person struct {
	Name string
	Age  int
}

func TestApSL_EncodingCombination(t *testing.T) {
	t.Run("combines encodings using monoid", func(t *testing.T) {
		// Create a lens for Person.Name
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Create base codec that encodes to "Person:"
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected Person",
					},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "Person:" },
		)

		// Create field codec for Name
		nameCodec := MakeType(
			"Name",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSL to combine encodings
		operator := ApSL(S.Monoid, nameLens, nameCodec)
		enhancedCodec := operator(baseCodec)

		// Test encoding - should concatenate base encoding with field encoding
		person := Person{Name: "Alice", Age: 30}
		encoded := enhancedCodec.Encode(person)

		// The monoid concatenates: base encoding + field encoding
		// Note: The order depends on how the monoid is applied in ApSL
		assert.Contains(t, encoded, "Person:")
		assert.Contains(t, encoded, "Alice")
	})
}

func TestApSL_ValidationCombination(t *testing.T) {
	t.Run("validates field through lens", func(t *testing.T) {
		// Create a lens for Person.Age
		ageLens := lens.MakeLens(
			func(p Person) int { return p.Age },
			func(p Person, age int) Person {
				return Person{Name: p.Name, Age: age}
			},
		)

		// Create base codec that always succeeds
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected Person",
					},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Create field codec for Age that validates positive numbers
		ageCodec := MakeType(
			"Age",
			func(i any) validation.Result[int] {
				if n, ok := i.(int); ok {
					if n > 0 {
						return validation.ToResult(validation.Success(n))
					}
					return validation.ToResult(validation.Failures[int](validation.Errors{
						&validation.ValidationError{
							Value:    n,
							Messsage: "age must be positive",
						},
					}))
				}
				return validation.ToResult(validation.Failures[int](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected int",
					},
				}))
			},
			func(i any) Decode[Context, int] {
				return func(ctx Context) validation.Validation[int] {
					if n, ok := i.(int); ok {
						if n > 0 {
							return validation.Success(n)
						}
						return validation.FailureWithMessage[int](n, "age must be positive")(ctx)
					}
					return validation.FailureWithMessage[int](i, "expected int")(ctx)
				}
			},
			strconv.Itoa,
		)

		// Apply ApSL
		operator := ApSL(S.Monoid, ageLens, ageCodec)
		enhancedCodec := operator(baseCodec)

		// Test with invalid age (negative) - field validation should fail
		invalidPerson := Person{Name: "Charlie", Age: -5}
		invalidResult := enhancedCodec.Decode(invalidPerson)
		assert.True(t, either.IsLeft(invalidResult), "Should fail with negative age")

		// Extract and verify we have errors
		errors := either.MonadFold(invalidResult,
			F.Identity[validation.Errors],
			func(Person) validation.Errors { return nil },
		)
		assert.NotEmpty(t, errors, "Should have validation errors")
	})
}

func TestApSL_TypeChecking(t *testing.T) {
	t.Run("preserves base type checker", func(t *testing.T) {
		// Create a lens for Person.Name
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Create base codec with type checker
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected Person",
					},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Create field codec
		nameCodec := MakeType(
			"Name",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSL
		operator := ApSL(S.Monoid, nameLens, nameCodec)
		enhancedCodec := operator(baseCodec)

		// Test type checking with valid type
		person := Person{Name: "Eve", Age: 22}
		isResult := enhancedCodec.Is(person)
		assert.True(t, either.IsRight(isResult), "Should accept Person type")

		// Test type checking with invalid type
		invalidResult := enhancedCodec.Is("not a person")
		assert.True(t, either.IsLeft(invalidResult), "Should reject non-Person type")
	})
}

func TestApSL_Naming(t *testing.T) {
	t.Run("generates descriptive name", func(t *testing.T) {
		// Create a lens for Person.Name
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Create base codec
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected Person",
					},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Create field codec
		nameCodec := MakeType(
			"Name",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSL
		operator := ApSL(S.Monoid, nameLens, nameCodec)
		enhancedCodec := operator(baseCodec)

		// Check that the name includes ApS
		name := enhancedCodec.Name()
		assert.Contains(t, name, "ApS", "Name should contain 'ApS'")
	})
}

func TestApSL_ErrorAccumulation(t *testing.T) {
	t.Run("accumulates validation errors", func(t *testing.T) {
		// Create a lens for Person.Age
		ageLens := lens.MakeLens(
			func(p Person) int { return p.Age },
			func(p Person, age int) Person {
				return Person{Name: p.Name, Age: age}
			},
		)

		// Create base codec that fails validation
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "base validation error",
					},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					return validation.FailureWithMessage[Person](i, "base validation error")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Create field codec that also fails
		ageCodec := MakeType(
			"Age",
			func(i any) validation.Result[int] {
				return validation.ToResult(validation.Failures[int](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "age validation error",
					},
				}))
			},
			func(i any) Decode[Context, int] {
				return func(ctx Context) validation.Validation[int] {
					return validation.FailureWithMessage[int](i, "age validation error")(ctx)
				}
			},
			strconv.Itoa,
		)

		// Apply ApSL
		operator := ApSL(S.Monoid, ageLens, ageCodec)
		enhancedCodec := operator(baseCodec)

		// Test validation - should accumulate errors
		person := Person{Name: "Dave", Age: 30}
		result := enhancedCodec.Decode(person)

		// Should fail
		assert.True(t, either.IsLeft(result), "Should fail validation")

		// Extract errors
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(Person) validation.Errors { return nil },
		)

		// Should have errors from both base and field validation
		assert.NotEmpty(t, errors, "Should have validation errors")
	})
}

// Test types for ApSO
type PersonWithNickname struct {
	Name     string
	Nickname *string
}

func TestApSO_EncodingWithPresentField(t *testing.T) {
	t.Run("encodes optional field when present", func(t *testing.T) {
		// Create an optional for PersonWithNickname.Nickname
		nicknameOpt := optional.MakeOptional(
			func(p PersonWithNickname) option.Option[string] {
				if p.Nickname != nil {
					return option.Some(*p.Nickname)
				}
				return option.None[string]()
			},
			func(p PersonWithNickname, nick string) PersonWithNickname {
				p.Nickname = &nick
				return p
			},
		)

		// Create base codec that encodes to "Person:"
		baseCodec := MakeType(
			"PersonWithNickname",
			func(i any) validation.Result[PersonWithNickname] {
				if p, ok := i.(PersonWithNickname); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[PersonWithNickname](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected PersonWithNickname",
					},
				}))
			},
			func(i any) Decode[Context, PersonWithNickname] {
				return func(ctx Context) validation.Validation[PersonWithNickname] {
					if p, ok := i.(PersonWithNickname); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[PersonWithNickname](i, "expected PersonWithNickname")(ctx)
				}
			},
			func(p PersonWithNickname) string { return "Person:" },
		)

		// Create field codec for Nickname
		nicknameCodec := MakeType(
			"Nickname",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSO to combine encodings
		operator := ApSO(S.Monoid, nicknameOpt, nicknameCodec)
		enhancedCodec := operator(baseCodec)

		// Test encoding with nickname present
		nickname := "Ali"
		person := PersonWithNickname{Name: "Alice", Nickname: &nickname}
		encoded := enhancedCodec.Encode(person)

		// Should include both base and nickname
		assert.Contains(t, encoded, "Person:")
		assert.Contains(t, encoded, "Ali")
	})
}

func TestApSO_EncodingWithAbsentField(t *testing.T) {
	t.Run("omits optional field when absent", func(t *testing.T) {
		// Create an optional for PersonWithNickname.Nickname
		nicknameOpt := optional.MakeOptional(
			func(p PersonWithNickname) option.Option[string] {
				if p.Nickname != nil {
					return option.Some(*p.Nickname)
				}
				return option.None[string]()
			},
			func(p PersonWithNickname, nick string) PersonWithNickname {
				p.Nickname = &nick
				return p
			},
		)

		// Create base codec
		baseCodec := MakeType(
			"PersonWithNickname",
			func(i any) validation.Result[PersonWithNickname] {
				if p, ok := i.(PersonWithNickname); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[PersonWithNickname](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected PersonWithNickname",
					},
				}))
			},
			func(i any) Decode[Context, PersonWithNickname] {
				return func(ctx Context) validation.Validation[PersonWithNickname] {
					if p, ok := i.(PersonWithNickname); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[PersonWithNickname](i, "expected PersonWithNickname")(ctx)
				}
			},
			func(p PersonWithNickname) string { return "Person:Bob" },
		)

		// Create field codec
		nicknameCodec := MakeType(
			"Nickname",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSO
		operator := ApSO(S.Monoid, nicknameOpt, nicknameCodec)
		enhancedCodec := operator(baseCodec)

		// Test encoding with nickname absent
		person := PersonWithNickname{Name: "Bob", Nickname: nil}
		encoded := enhancedCodec.Encode(person)

		// Should only have base encoding
		assert.Equal(t, "Person:Bob", encoded)
	})
}

func TestApSO_TypeChecking(t *testing.T) {
	t.Run("preserves base type checker", func(t *testing.T) {
		// Create an optional for PersonWithNickname.Nickname
		nicknameOpt := optional.MakeOptional(
			func(p PersonWithNickname) option.Option[string] {
				if p.Nickname != nil {
					return option.Some(*p.Nickname)
				}
				return option.None[string]()
			},
			func(p PersonWithNickname, nick string) PersonWithNickname {
				p.Nickname = &nick
				return p
			},
		)

		// Create base codec with type checker
		baseCodec := MakeType(
			"PersonWithNickname",
			func(i any) validation.Result[PersonWithNickname] {
				if p, ok := i.(PersonWithNickname); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[PersonWithNickname](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected PersonWithNickname",
					},
				}))
			},
			func(i any) Decode[Context, PersonWithNickname] {
				return func(ctx Context) validation.Validation[PersonWithNickname] {
					if p, ok := i.(PersonWithNickname); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[PersonWithNickname](i, "expected PersonWithNickname")(ctx)
				}
			},
			func(p PersonWithNickname) string { return "" },
		)

		// Create field codec
		nicknameCodec := MakeType(
			"Nickname",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSO
		operator := ApSO(S.Monoid, nicknameOpt, nicknameCodec)
		enhancedCodec := operator(baseCodec)

		// Test type checking with valid type
		nickname := "Eve"
		person := PersonWithNickname{Name: "Eve", Nickname: &nickname}
		isResult := enhancedCodec.Is(person)
		assert.True(t, either.IsRight(isResult), "Should accept PersonWithNickname type")

		// Test type checking with invalid type
		invalidResult := enhancedCodec.Is("not a person")
		assert.True(t, either.IsLeft(invalidResult), "Should reject non-PersonWithNickname type")
	})
}

func TestApSO_Naming(t *testing.T) {
	t.Run("generates descriptive name", func(t *testing.T) {
		// Create an optional for PersonWithNickname.Nickname
		nicknameOpt := optional.MakeOptional(
			func(p PersonWithNickname) option.Option[string] {
				if p.Nickname != nil {
					return option.Some(*p.Nickname)
				}
				return option.None[string]()
			},
			func(p PersonWithNickname, nick string) PersonWithNickname {
				p.Nickname = &nick
				return p
			},
		)

		// Create base codec
		baseCodec := MakeType(
			"PersonWithNickname",
			func(i any) validation.Result[PersonWithNickname] {
				if p, ok := i.(PersonWithNickname); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[PersonWithNickname](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected PersonWithNickname",
					},
				}))
			},
			func(i any) Decode[Context, PersonWithNickname] {
				return func(ctx Context) validation.Validation[PersonWithNickname] {
					if p, ok := i.(PersonWithNickname); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[PersonWithNickname](i, "expected PersonWithNickname")(ctx)
				}
			},
			func(p PersonWithNickname) string { return "" },
		)

		// Create field codec
		nicknameCodec := MakeType(
			"Nickname",
			func(i any) validation.Result[string] {
				if s, ok := i.(string); ok {
					return validation.ToResult(validation.Success(s))
				}
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "expected string",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					if s, ok := i.(string); ok {
						return validation.Success(s)
					}
					return validation.FailureWithMessage[string](i, "expected string")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSO
		operator := ApSO(S.Monoid, nicknameOpt, nicknameCodec)
		enhancedCodec := operator(baseCodec)

		// Check that the name includes ApS
		name := enhancedCodec.Name()
		assert.Contains(t, name, "ApS", "Name should contain 'ApS'")
	})
}

func TestApSO_ErrorAccumulation(t *testing.T) {
	t.Run("accumulates validation errors", func(t *testing.T) {
		// Create an optional for PersonWithNickname.Nickname
		nicknameOpt := optional.MakeOptional(
			func(p PersonWithNickname) option.Option[string] {
				if p.Nickname != nil {
					return option.Some(*p.Nickname)
				}
				return option.None[string]()
			},
			func(p PersonWithNickname, nick string) PersonWithNickname {
				p.Nickname = &nick
				return p
			},
		)

		// Create base codec that fails validation
		baseCodec := MakeType(
			"PersonWithNickname",
			func(i any) validation.Result[PersonWithNickname] {
				return validation.ToResult(validation.Failures[PersonWithNickname](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "base validation error",
					},
				}))
			},
			func(i any) Decode[Context, PersonWithNickname] {
				return func(ctx Context) validation.Validation[PersonWithNickname] {
					return validation.FailureWithMessage[PersonWithNickname](i, "base validation error")(ctx)
				}
			},
			func(p PersonWithNickname) string { return "" },
		)

		// Create field codec that also fails
		nicknameCodec := MakeType(
			"Nickname",
			func(i any) validation.Result[string] {
				return validation.ToResult(validation.Failures[string](validation.Errors{
					&validation.ValidationError{
						Value:    i,
						Messsage: "nickname validation error",
					},
				}))
			},
			func(i any) Decode[Context, string] {
				return func(ctx Context) validation.Validation[string] {
					return validation.FailureWithMessage[string](i, "nickname validation error")(ctx)
				}
			},
			F.Identity[string],
		)

		// Apply ApSO
		operator := ApSO(S.Monoid, nicknameOpt, nicknameCodec)
		enhancedCodec := operator(baseCodec)

		// Test validation with present nickname - should accumulate errors
		nickname := "Dave"
		person := PersonWithNickname{Name: "Dave", Nickname: &nickname}
		result := enhancedCodec.Decode(person)

		// Should fail
		assert.True(t, either.IsLeft(result), "Should fail validation")

		// Extract errors
		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(PersonWithNickname) validation.Errors { return nil },
		)

		// Should have errors from both base and field validation
		assert.NotEmpty(t, errors, "Should have validation errors")
	})
}

// TestBind_EncodingCombination verifies that Bind combines the base encoding with
// the field encoding produced by the Kleisli arrow using the monoid.
func TestBind_EncodingCombination(t *testing.T) {
	t.Run("combines base and field encodings using monoid", func(t *testing.T) {
		// Lens for Person.Name
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Base codec encodes to "Person:"
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "Person:" },
		)

		// Kleisli arrow: always returns a string identity codec regardless of struct value
		kleisli := func(p Person) Type[string, string, any] {
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if s, ok := i.(string); ok {
						return validation.ToResult(validation.Success(s))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected string"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if s, ok := i.(string); ok {
							return validation.Success(s)
						}
						return validation.FailureWithMessage[string](i, "expected string")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		person := Person{Name: "Alice", Age: 30}
		encoded := enhancedCodec.Encode(person)

		// Encoding should include both the base prefix and the field value
		assert.Contains(t, encoded, "Person:")
		assert.Contains(t, encoded, "Alice")
	})
}

// TestBind_KleisliArrowReceivesCurrentValue verifies that the Kleisli arrow f
// receives the current struct value when producing the field codec.
func TestBind_KleisliArrowReceivesCurrentValue(t *testing.T) {
	t.Run("kleisli arrow receives current struct value during encoding", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Kleisli arrow that uses the struct value to produce a prefix in the encoding
		var capturedPerson Person
		kleisli := func(p Person) Type[string, string, any] {
			capturedPerson = p
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if s, ok := i.(string); ok {
						return validation.ToResult(validation.Success(s))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected string"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if s, ok := i.(string); ok {
							return validation.Success(s)
						}
						return validation.FailureWithMessage[string](i, "expected string")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		person := Person{Name: "Bob", Age: 25}
		enhancedCodec.Encode(person)

		// The Kleisli arrow should have been called with the actual struct value
		assert.Equal(t, person, capturedPerson)
	})
}

// TestBind_ValidationSuccess verifies that Bind correctly validates and decodes
// a struct when both the base and field validations succeed.
func TestBind_ValidationSuccess(t *testing.T) {
	t.Run("succeeds when base and field validations pass", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// The field codec receives the same input I (any = Person struct).
		// It must extract the Name field from the Person input.
		kleisli := func(p Person) Type[string, string, any] {
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if person, ok := i.(Person); ok {
						return validation.ToResult(validation.Success(person.Name))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected Person"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if person, ok := i.(Person); ok {
							return validation.Success(person.Name)
						}
						return validation.FailureWithMessage[string](i, "expected Person")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		person := Person{Name: "Carol", Age: 28}
		result := enhancedCodec.Decode(person)

		assert.True(t, either.IsRight(result), "Should succeed when both validations pass")
	})
}

// TestBind_ValidationFailsOnBaseFailure verifies that Bind uses fail-fast (monadic)
// semantics: if the base codec fails, the Kleisli arrow is never evaluated.
func TestBind_ValidationFailsOnBaseFailure(t *testing.T) {
	t.Run("fails fast when base validation fails", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Base codec always fails
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "base always fails"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					return validation.FailureWithMessage[Person](i, "base always fails")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		kleisliCalled := false
		kleisli := func(p Person) Type[string, string, any] {
			kleisliCalled = true
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if s, ok := i.(string); ok {
						return validation.ToResult(validation.Success(s))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected string"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if s, ok := i.(string); ok {
							return validation.Success(s)
						}
						return validation.FailureWithMessage[string](i, "expected string")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		person := Person{Name: "Dave", Age: 40}
		result := enhancedCodec.Decode(person)

		assert.True(t, either.IsLeft(result), "Should fail when base validation fails")
		assert.False(t, kleisliCalled, "Kleisli arrow should NOT be called when base fails")
	})
}

// TestBind_ValidationFailsOnFieldFailure verifies that Bind propagates field
// validation errors when the Kleisli arrow's codec fails.
func TestBind_ValidationFailsOnFieldFailure(t *testing.T) {
	t.Run("fails when field validation from kleisli codec fails", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Base codec succeeds
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Kleisli arrow returns a codec that always fails regardless of input
		kleisli := func(p Person) Type[string, string, any] {
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "field always fails"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						return validation.FailureWithMessage[string](i, "field always fails")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		// The field codec receives the same input (Person) and always fails
		person := Person{Name: "Eve", Age: 22}
		result := enhancedCodec.Decode(person)

		assert.True(t, either.IsLeft(result), "Should fail when field validation fails")

		errors := either.MonadFold(result,
			F.Identity[validation.Errors],
			func(Person) validation.Errors { return nil },
		)
		assert.NotEmpty(t, errors, "Should have validation errors from field codec")
	})
}

// TestBind_TypeCheckingPreserved verifies that Bind preserves the base type checker.
func TestBind_TypeCheckingPreserved(t *testing.T) {
	t.Run("preserves base type checker", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		kleisli := func(p Person) Type[string, string, any] {
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if s, ok := i.(string); ok {
						return validation.ToResult(validation.Success(s))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected string"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if s, ok := i.(string); ok {
							return validation.Success(s)
						}
						return validation.FailureWithMessage[string](i, "expected string")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		// Valid type
		person := Person{Name: "Frank", Age: 35}
		isResult := enhancedCodec.Is(person)
		assert.True(t, either.IsRight(isResult), "Should accept Person type")

		// Invalid type
		invalidResult := enhancedCodec.Is("not a person")
		assert.True(t, either.IsLeft(invalidResult), "Should reject non-Person type")
	})
}

// TestBind_Naming verifies that Bind generates a descriptive name for the codec.
func TestBind_Naming(t *testing.T) {
	t.Run("generates descriptive name containing Bind and lens info", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		kleisli := func(p Person) Type[string, string, any] {
			return MakeType(
				"Name",
				func(i any) validation.Result[string] {
					if s, ok := i.(string); ok {
						return validation.ToResult(validation.Success(s))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected string"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if s, ok := i.(string); ok {
							return validation.Success(s)
						}
						return validation.FailureWithMessage[string](i, "expected string")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		name := enhancedCodec.Name()
		assert.Contains(t, name, "Bind", "Name should contain 'Bind'")
	})
}

// TestBind_DependentFieldCodec verifies that the Kleisli arrow can produce
// different codecs based on the current struct value (the key differentiator
// from ApSL).
//
// The field codec Type[T, O, I] receives the same input I as the base codec.
// It must extract the field value from that input. The Kleisli arrow f(s)
// produces a different codec depending on the already-decoded struct value s.
func TestBind_DependentFieldCodec(t *testing.T) {
	t.Run("kleisli arrow produces different codecs based on struct value", func(t *testing.T) {
		// Lens for Person.Name
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person {
				return Person{Name: name, Age: p.Age}
			},
		)

		// Base codec succeeds for any Person
		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Kleisli arrow: the field codec receives the same input I (any = Person).
		// It extracts the Name from the Person input.
		// If the decoded struct's Age > 18, accept any name (including empty).
		// If Age <= 18, reject empty names.
		kleisli := func(p Person) Type[string, string, any] {
			if p.Age > 18 {
				// Adult: accept any name extracted from the Person input
				return MakeType(
					"AnyName",
					func(i any) validation.Result[string] {
						if person, ok := i.(Person); ok {
							return validation.ToResult(validation.Success(person.Name))
						}
						return validation.ToResult(validation.Failures[string](validation.Errors{
							&validation.ValidationError{Value: i, Messsage: "expected Person"},
						}))
					},
					func(i any) Decode[Context, string] {
						return func(ctx Context) validation.Validation[string] {
							if person, ok := i.(Person); ok {
								return validation.Success(person.Name)
							}
							return validation.FailureWithMessage[string](i, "expected Person")(ctx)
						}
					},
					F.Identity[string],
				)
			}
			// Minor: reject empty names
			return MakeType(
				"NonEmptyName",
				func(i any) validation.Result[string] {
					if person, ok := i.(Person); ok {
						if person.Name != "" {
							return validation.ToResult(validation.Success(person.Name))
						}
						return validation.ToResult(validation.Failures[string](validation.Errors{
							&validation.ValidationError{Value: person.Name, Messsage: "name must not be empty for minors"},
						}))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected Person"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if person, ok := i.(Person); ok {
							if person.Name != "" {
								return validation.Success(person.Name)
							}
							return validation.FailureWithMessage[string](person.Name, "name must not be empty for minors")(ctx)
						}
						return validation.FailureWithMessage[string](i, "expected Person")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		// Adult (Age=30) with empty name: should succeed (adult codec accepts any name)
		adultPerson := Person{Name: "", Age: 30}
		adultResult := enhancedCodec.Decode(adultPerson)
		assert.True(t, either.IsRight(adultResult), "Adult should accept empty name")

		// Minor (Age=15) with empty name: should fail (minor codec rejects empty names)
		minorPerson := Person{Name: "", Age: 15}
		minorResult := enhancedCodec.Decode(minorPerson)
		assert.True(t, either.IsLeft(minorResult), "Minor with empty name should fail")

		// Minor (Age=15) with non-empty name: should succeed
		minorWithName := Person{Name: "Junior", Age: 15}
		minorWithNameResult := enhancedCodec.Decode(minorWithName)
		assert.True(t, either.IsRight(minorWithNameResult), "Minor with non-empty name should succeed")
	})
}

// TestBind_EncodingOrder verifies that the base encoding comes before the field
// encoding in the monoid concatenation: m.Concat(base, field).
// String concatenation is non-commutative, so order is observable.
func TestBind_EncodingOrder(t *testing.T) {
	t.Run("base encoding precedes field encoding in Concat", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { return Person{Name: name, Age: p.Age} },
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "BASE-" },
		)

		kleisli := func(p Person) Type[string, string, any] {
			return MakeType("Name",
				func(i any) validation.Result[string] {
					return validation.ToResult(validation.Success(""))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] { return validation.Success("") }
				},
				func(name string) string { return "FIELD" },
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		encoded := enhancedCodec.Encode(Person{Name: "Alice", Age: 30})

		// S.Monoid.Concat("BASE-", "FIELD") = "BASE-FIELD", not "FIELDBASE-"
		assert.Equal(t, "BASE-FIELD", encoded)
	})
}

// TestBind_DecodedStructHasFieldValue verifies that after a successful decode the
// lens setter has been applied: the returned struct contains the validated field value.
func TestBind_DecodedStructHasFieldValue(t *testing.T) {
	t.Run("decoded struct contains the field value written by the lens setter", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { return Person{Name: name, Age: p.Age} },
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// Field codec extracts Name from the Person input.
		kleisli := func(p Person) Type[string, string, any] {
			return MakeType("Name",
				func(i any) validation.Result[string] {
					if person, ok := i.(Person); ok {
						return validation.ToResult(validation.Success(person.Name))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected Person"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if person, ok := i.(Person); ok {
							return validation.Success(person.Name)
						}
						return validation.FailureWithMessage[string](i, "expected Person")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		operator := Bind(S.Monoid, nameLens, kleisli)
		enhancedCodec := operator(baseCodec)

		result := enhancedCodec.Decode(Person{Name: "Grace", Age: 50})

		decoded := either.MonadFold(result,
			func(validation.Errors) Person { t.Fatal("unexpected Left"); return Person{} },
			F.Identity[Person],
		)
		assert.Equal(t, "Grace", decoded.Name, "lens setter must write the validated name into the struct")
		assert.Equal(t, 50, decoded.Age, "unrelated field must be preserved")
	})
}

// TestBind_ChainedBindOperators verifies that two consecutive Bind operators
// compose correctly: the second Kleisli arrow receives the struct with the first
// field already set, and both fields appear in the final decoded struct.
func TestBind_ChainedBindOperators(t *testing.T) {
	t.Run("second kleisli receives struct with first field already set", func(t *testing.T) {
		nameLens := lens.MakeLens(
			func(p Person) string { return p.Name },
			func(p Person, name string) Person { return Person{Name: name, Age: p.Age} },
		)
		ageLens := lens.MakeLens(
			func(p Person) int { return p.Age },
			func(p Person, age int) Person { return Person{Name: p.Name, Age: age} },
		)

		baseCodec := MakeType(
			"Person",
			func(i any) validation.Result[Person] {
				if p, ok := i.(Person); ok {
					return validation.ToResult(validation.Success(p))
				}
				return validation.ToResult(validation.Failures[Person](validation.Errors{
					&validation.ValidationError{Value: i, Messsage: "expected Person"},
				}))
			},
			func(i any) Decode[Context, Person] {
				return func(ctx Context) validation.Validation[Person] {
					if p, ok := i.(Person); ok {
						return validation.Success(p)
					}
					return validation.FailureWithMessage[Person](i, "expected Person")(ctx)
				}
			},
			func(p Person) string { return "" },
		)

		// First bind: decodes Name.
		nameKleisli := func(p Person) Type[string, string, any] {
			return MakeType("Name",
				func(i any) validation.Result[string] {
					if person, ok := i.(Person); ok {
						return validation.ToResult(validation.Success(person.Name))
					}
					return validation.ToResult(validation.Failures[string](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected Person"},
					}))
				},
				func(i any) Decode[Context, string] {
					return func(ctx Context) validation.Validation[string] {
						if person, ok := i.(Person); ok {
							return validation.Success(person.Name)
						}
						return validation.FailureWithMessage[string](i, "expected Person")(ctx)
					}
				},
				F.Identity[string],
			)
		}

		// Second bind: decodes Age. Captures the struct it receives so we can assert
		// that Name was already set by the first bind.
		var capturedStructInSecondBind Person
		ageKleisli := func(p Person) Type[int, string, any] {
			capturedStructInSecondBind = p
			return MakeType("Age",
				func(i any) validation.Result[int] {
					if person, ok := i.(Person); ok {
						return validation.ToResult(validation.Success(person.Age))
					}
					return validation.ToResult(validation.Failures[int](validation.Errors{
						&validation.ValidationError{Value: i, Messsage: "expected Person"},
					}))
				},
				func(i any) Decode[Context, int] {
					return func(ctx Context) validation.Validation[int] {
						if person, ok := i.(Person); ok {
							return validation.Success(person.Age)
						}
						return validation.FailureWithMessage[int](i, "expected Person")(ctx)
					}
				},
				func(age int) string { return strconv.Itoa(age) },
			)
		}

		codec := F.Pipe2(
			baseCodec,
			Bind(S.Monoid, nameLens, nameKleisli),
			Bind(S.Monoid, ageLens, ageKleisli),
		)

		input := Person{Name: "Henry", Age: 42}
		result := codec.Decode(input)

		decoded := either.MonadFold(result,
			func(validation.Errors) Person { t.Fatal("unexpected Left"); return Person{} },
			F.Identity[Person],
		)

		assert.Equal(t, "Henry", decoded.Name, "Name must be set by first Bind")
		assert.Equal(t, 42, decoded.Age, "Age must be set by second Bind")
		// The second Kleisli received the struct *after* the first bind wrote Name into it.
		assert.Equal(t, "Henry", capturedStructInSecondBind.Name,
			"second Kleisli arrow must see the Name set by the first Bind")
	})
}

// Wrapper is a simple sum-type analogue used by the BindTo tests.
// Its only variant holds a string value.
type Wrapper struct{ Value string }

// wrapperPrism focuses on Wrapper.Value.
var wrapperPrism = prism.MakePrism(
	func(w Wrapper) option.Option[string] { return option.Some(w.Value) },
	func(s string) Wrapper { return Wrapper{Value: s} },
)

// TestBindTo_DecodesInnerType verifies that BindTo decodes the input using the
// inner codec and constructs S from T via the prism's ReverseGet.
func TestBindTo_DecodesInnerType(t *testing.T) {
	t.Run("decodes input to inner type and wraps with ReverseGet", func(t *testing.T) {
		wrapperCodec := BindTo(S.Monoid, wrapperPrism, String())

		result := wrapperCodec.Decode("hello")

		assert.True(t, either.IsRight(result), "Should succeed for valid string input")
		wrapper := either.MonadFold(result,
			func(validation.Errors) Wrapper { t.Fatal("unexpected Left"); return Wrapper{} },
			F.Identity[Wrapper],
		)
		assert.Equal(t, Wrapper{Value: "hello"}, wrapper)
	})
}

// TestBindTo_DecodesFailurePropagate verifies that a decoding failure in the inner
// codec is propagated as a Left by the BindTo codec.
func TestBindTo_DecodesFailurePropagate(t *testing.T) {
	t.Run("propagates inner codec validation failure", func(t *testing.T) {
		// IntFromString expects a string that parses as an int; passing a non-numeric
		// string should cause the inner codec to fail, and BindTo must forward that.
		type Box struct{ N int }
		boxPrism := prism.MakePrism(
			func(b Box) option.Option[int] { return option.Some(b.N) },
			func(n int) Box { return Box{N: n} },
		)

		boxCodec := BindTo(S.Monoid, boxPrism, IntFromString())

		result := boxCodec.Decode("not-a-number")

		assert.True(t, either.IsLeft(result), "Should fail when inner codec fails")
	})
}

// TestBindTo_EncodesWithGetOption verifies that BindTo encodes using the inner
// codec when the prism's GetOption succeeds.
func TestBindTo_EncodesWithGetOption(t *testing.T) {
	t.Run("encodes inner value when prism matches", func(t *testing.T) {
		wrapperCodec := BindTo(S.Monoid, wrapperPrism, String())

		encoded := wrapperCodec.Encode(Wrapper{Value: "world"})

		assert.Equal(t, "world", encoded)
	})
}

// TestBindTo_EncodesEmptyWhenPrismMisses verifies that BindTo returns the
// monoid's empty value when the prism's GetOption returns None.
func TestBindTo_EncodesEmptyWhenPrismMisses(t *testing.T) {
	t.Run("returns monoid empty when prism GetOption returns None", func(t *testing.T) {
		// Build a prism that never matches (always returns None for GetOption).
		neverMatchPrism := prism.MakePrism(
			func(w Wrapper) option.Option[string] { return option.None[string]() },
			func(s string) Wrapper { return Wrapper{Value: s} },
		)

		wrapperCodec := BindTo(S.Monoid, neverMatchPrism, String())

		encoded := wrapperCodec.Encode(Wrapper{Value: "ignored"})

		// S.Monoid.Empty() returns ""
		assert.Equal(t, S.Monoid.Empty(), encoded)
	})
}

// TestBindTo_NameContainsBindTo verifies the generated codec name.
func TestBindTo_NameContainsBindTo(t *testing.T) {
	t.Run("codec name contains BindTo", func(t *testing.T) {
		wrapperCodec := BindTo(S.Monoid, wrapperPrism, String())

		assert.Contains(t, wrapperCodec.Name(), "BindTo")
	})
}

// TestBindTo_RoundTrip verifies that encoding a decoded value is a no-op for a
// well-formed prism (GetOption ∘ ReverseGet = Some).
func TestBindTo_RoundTrip(t *testing.T) {
	t.Run("decode then encode returns original string", func(t *testing.T) {
		wrapperCodec := BindTo(S.Monoid, wrapperPrism, String())

		decoded := either.MonadFold(wrapperCodec.Decode("fp-go"),
			func(validation.Errors) Wrapper { t.Fatal("unexpected Left"); return Wrapper{} },
			F.Identity[Wrapper],
		)
		encoded := wrapperCodec.Encode(decoded)

		assert.Equal(t, "fp-go", encoded)
	})
}

// TestBindTo_IntFromString verifies BindTo with a non-string inner codec to
// confirm the generic type parameters work correctly.
func TestBindTo_IntFromString(t *testing.T) {
	t.Run("decodes integer via IntFromString inner codec", func(t *testing.T) {
		type Box struct{ N int }
		boxPrism := prism.MakePrism(
			func(b Box) option.Option[int] { return option.Some(b.N) },
			func(n int) Box { return Box{N: n} },
		)

		// IntFromString has I = string, so BindTo is instantiated with I = string.
		boxCodec := BindTo(S.Monoid, boxPrism, IntFromString())

		result := boxCodec.Decode("42")

		box := either.MonadFold(result,
			func(validation.Errors) Box { t.Fatal("unexpected Left"); return Box{} },
			F.Identity[Box],
		)
		assert.Equal(t, Box{N: 42}, box)

		// Also verify encoding roundtrips.
		encoded := boxCodec.Encode(box)
		assert.Equal(t, strconv.Itoa(42), encoded)
	})
}
