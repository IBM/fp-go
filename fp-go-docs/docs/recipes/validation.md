---
sidebar_position: 3
title: Validation
description: Input validation with accumulating errors
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 03 / 17"
  title="Input"
  titleAccent="Validation"
  lede="Validate input data functionally with the ability to accumulate multiple validation errors instead of failing on the first error."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '4' },
    { label: 'Use Cases', value: 'Forms, APIs, Data Quality' }
  ]}
/>

<TLDR>
  <TLDRCard title="Accumulate Errors" icon="list">
    Don't stop at the first error—collect all validation issues to provide comprehensive feedback to users.
  </TLDRCard>
  <TLDRCard title="Compose Validators" icon="layers">
    Build complex validators from simple, reusable ones—create a library of validation building blocks.
  </TLDRCard>
  <TLDRCard title="Include Field Paths" icon="map-pin">
    For nested objects, include the full path to invalid fields—makes debugging and error display easier.
  </TLDRCard>
</TLDR>

<Section id="simple-field" number="01" title="Simple Field" titleAccent="Validation">

Validate individual fields and collect errors using Either for single-field validation.

<CodeCard file="simple-field-validation.go">
{`package main

import (
    "fmt"
    "regexp"
    "strings"
    
    E "github.com/IBM/fp-go/v2/either"
)

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Email validation
var emailRegex = regexp.MustCompile(\`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\`)

func validateEmail(email string) E.Either[ValidationError, string] {
    if email == "" {
        return E.Left[string](ValidationError{"email", "is required"})
    }
    if !emailRegex.MatchString(email) {
        return E.Left[string](ValidationError{"email", "invalid format"})
    }
    return E.Right[ValidationError](email)
}

// Password validation
func validatePassword(password string) E.Either[ValidationError, string] {
    if len(password) < 8 {
        return E.Left[string](ValidationError{"password", "must be at least 8 characters"})
    }
    if !containsDigit(password) {
        return E.Left[string](ValidationError{"password", "must contain at least one digit"})
    }
    return E.Right[ValidationError](password)
}

func containsDigit(s string) bool {
    for _, c := range s {
        if c >= '0' && c <= '9' {
            return true
        }
    }
    return false
}

func main() {
    // Valid email
    email1 := validateEmail("user@example.com")
    fmt.Println(E.IsRight(email1)) // true
    
    // Invalid email
    email2 := validateEmail("invalid")
    if E.IsLeft(email2) {
        fmt.Println(E.GetLeft(email2).Error())
        // email: invalid format
    }
    
    // Valid password
    pwd1 := validatePassword("secret123")
    fmt.Println(E.IsRight(pwd1)) // true
    
    // Invalid password
    pwd2 := validatePassword("short")
    if E.IsLeft(pwd2) {
        fmt.Println(E.GetLeft(pwd2).Error())
        // password: must be at least 8 characters
    }
}`}
</CodeCard>

</Section>

<Section id="accumulating-errors" number="02" title="Accumulating Multiple" titleAccent="Errors">

Collect all validation errors instead of stopping at the first one, providing comprehensive feedback.

<CodeCard file="accumulating-errors.go">
{`package main

import (
    "fmt"
    "strings"
)

type ValidationErrors []string

func (ve ValidationErrors) Error() string {
    return strings.Join(ve, "; ")
}

type SignupForm struct {
    Username string
    Email    string
    Password string
    Age      int
}

// Validate entire form and collect all errors
func validateSignupForm(form SignupForm) (SignupForm, ValidationErrors) {
    var errors ValidationErrors
    
    // Username validation
    if len(form.Username) < 3 {
        errors = append(errors, "username must be at least 3 characters")
    }
    if len(form.Username) > 20 {
        errors = append(errors, "username must be at most 20 characters")
    }
    
    // Email validation
    if form.Email == "" {
        errors = append(errors, "email is required")
    } else if !strings.Contains(form.Email, "@") {
        errors = append(errors, "email must contain @")
    }
    
    // Password validation
    if len(form.Password) < 8 {
        errors = append(errors, "password must be at least 8 characters")
    }
    if !containsDigit(form.Password) {
        errors = append(errors, "password must contain a digit")
    }
    if !containsUpper(form.Password) {
        errors = append(errors, "password must contain an uppercase letter")
    }
    
    // Age validation
    if form.Age < 13 {
        errors = append(errors, "must be at least 13 years old")
    }
    if form.Age > 120 {
        errors = append(errors, "age must be realistic")
    }
    
    return form, errors
}

func containsDigit(s string) bool {
    for _, c := range s {
        if c >= '0' && c <= '9' {
            return true
        }
    }
    return false
}

func containsUpper(s string) bool {
    for _, c := range s {
        if c >= 'A' && c <= 'Z' {
            return true
        }
    }
    return false
}

func main() {
    // Invalid form with multiple errors
    form := SignupForm{
        Username: "ab",           // too short
        Email:    "invalid",      // no @
        Password: "weak",         // too short, no digit, no uppercase
        Age:      10,             // too young
    }
    
    _, errors := validateSignupForm(form)
    if len(errors) > 0 {
        fmt.Printf("Validation failed with %d errors:\\n", len(errors))
        for i, err := range errors {
            fmt.Printf("%d. %s\\n", i+1, err)
        }
    }
    
    // Valid form
    validForm := SignupForm{
        Username: "alice",
        Email:    "alice@example.com",
        Password: "Secret123",
        Age:      25,
    }
    
    _, validErrors := validateSignupForm(validForm)
    if len(validErrors) == 0 {
        fmt.Println("Form is valid!")
    }
}`}
</CodeCard>

</Section>

<Section id="custom-types" number="03" title="Validation with Custom" titleAccent="Types">

Create reusable validators with custom error types for maximum flexibility and composability.

<CodeCard file="custom-validators.go">
{`package main

import (
    "fmt"
    "strings"
)

// Validator function type
type Validator[T any] func(T) []string

// Combine multiple validators
func combineValidators[T any](validators ...Validator[T]) Validator[T] {
    return func(value T) []string {
        var errors []string
        for _, validator := range validators {
            errors = append(errors, validator(value)...)
        }
        return errors
    }
}

// String validators
func minLength(min int) Validator[string] {
    return func(s string) []string {
        if len(s) < min {
            return []string{fmt.Sprintf("must be at least %d characters", min)}
        }
        return nil
    }
}

func maxLength(max int) Validator[string] {
    return func(s string) []string {
        if len(s) > max {
            return []string{fmt.Sprintf("must be at most %d characters", max)}
        }
        return nil
    }
}

func required() Validator[string] {
    return func(s string) []string {
        if strings.TrimSpace(s) == "" {
            return []string{"is required"}
        }
        return nil
    }
}

func pattern(regex string, message string) Validator[string] {
    return func(s string) []string {
        // Simple pattern check (in real code, use regexp package)
        if !strings.Contains(s, "@") && message == "invalid email format" {
            return []string{message}
        }
        return nil
    }
}

// Number validators
func min(minVal int) Validator[int] {
    return func(n int) []string {
        if n < minVal {
            return []string{fmt.Sprintf("must be at least %d", minVal)}
        }
        return nil
    }
}

func max(maxVal int) Validator[int] {
    return func(n int) []string {
        if n > maxVal {
            return []string{fmt.Sprintf("must be at most %d", maxVal)}
        }
        return nil
    }
}

// Field validation result
type FieldValidation struct {
    Field  string
    Errors []string
}

func validateField[T any](field string, value T, validator Validator[T]) *FieldValidation {
    errors := validator(value)
    if len(errors) > 0 {
        return &FieldValidation{Field: field, Errors: errors}
    }
    return nil
}

type User struct {
    Username string
    Email    string
    Age      int
}

func validateUser(user User) []FieldValidation {
    var validations []FieldValidation
    
    // Validate username
    usernameValidator := combineValidators(
        required(),
        minLength(3),
        maxLength(20),
    )
    if v := validateField("username", user.Username, usernameValidator); v != nil {
        validations = append(validations, *v)
    }
    
    // Validate email
    emailValidator := combineValidators(
        required(),
        pattern("@", "invalid email format"),
    )
    if v := validateField("email", user.Email, emailValidator); v != nil {
        validations = append(validations, *v)
    }
    
    // Validate age
    ageValidator := combineValidators(
        min(13),
        max(120),
    )
    if v := validateField("age", user.Age, ageValidator); v != nil {
        validations = append(validations, *v)
    }
    
    return validations
}

func main() {
    user := User{
        Username: "ab",
        Email:    "invalid",
        Age:      10,
    }
    
    validations := validateUser(user)
    if len(validations) > 0 {
        fmt.Println("Validation errors:")
        for _, v := range validations {
            fmt.Printf("%s:\\n", v.Field)
            for _, err := range v.Errors {
                fmt.Printf("  - %s\\n", err)
            }
        }
    }
}`}
</CodeCard>

</Section>

<Section id="nested-objects" number="04" title="Nested Object" titleAccent="Validation">

Validate complex nested structures with full path tracking for precise error reporting.

<CodeCard file="nested-validation.go">
{`package main

import (
    "fmt"
)

type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Person struct {
    Name    string
    Age     int
    Address Address
}

type ValidationResult struct {
    Path   string
    Errors []string
}

func validateAddress(addr Address, path string) []ValidationResult {
    var results []ValidationResult
    
    if addr.Street == "" {
        results = append(results, ValidationResult{
            Path:   path + ".street",
            Errors: []string{"is required"},
        })
    }
    
    if addr.City == "" {
        results = append(results, ValidationResult{
            Path:   path + ".city",
            Errors: []string{"is required"},
        })
    }
    
    if len(addr.ZipCode) != 5 {
        results = append(results, ValidationResult{
            Path:   path + ".zipCode",
            Errors: []string{"must be 5 digits"},
        })
    }
    
    return results
}

func validatePerson(person Person) []ValidationResult {
    var results []ValidationResult
    
    if person.Name == "" {
        results = append(results, ValidationResult{
            Path:   "name",
            Errors: []string{"is required"},
        })
    }
    
    if person.Age < 0 {
        results = append(results, ValidationResult{
            Path:   "age",
            Errors: []string{"must be non-negative"},
        })
    }
    
    // Validate nested address
    addressResults := validateAddress(person.Address, "address")
    results = append(results, addressResults...)
    
    return results
}

func main() {
    person := Person{
        Name: "",
        Age:  -5,
        Address: Address{
            Street:  "",
            City:    "New York",
            ZipCode: "123", // invalid
        },
    }
    
    results := validatePerson(person)
    if len(results) > 0 {
        fmt.Println("Validation errors:")
        for _, r := range results {
            fmt.Printf("%s: %v\\n", r.Path, r.Errors)
        }
    }
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="05" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Accumulate errors** — Don't stop at the first error; collect all validation issues
  </ChecklistItem>
  <ChecklistItem status="required">
    **Use descriptive messages** — Make error messages clear and actionable
  </ChecklistItem>
  <ChecklistItem status="required">
    **Validate at boundaries** — Validate input at system boundaries (API, forms, etc.)
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Compose validators** — Build complex validators from simple, reusable ones
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Include field paths** — For nested objects, include the full path to the invalid field
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Separate validation logic** — Keep validation separate from business logic
  </ChecklistItem>
</Checklist>

</Section>
