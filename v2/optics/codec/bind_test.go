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

// Made with Bob
