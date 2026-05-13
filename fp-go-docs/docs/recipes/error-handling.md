---
sidebar_position: 2
title: Error Handling Patterns
hide_title: true
description: Common patterns for handling errors functionally
---

<PageHeader
  eyebrow="Recipes · Error Handling"
  title="Error Handling"
  titleAccent="Patterns."
  lede="Common patterns for handling errors functionally using fp-go's Result, Either, and IOResult types. Learn to chain, recover, and transform errors."
  meta={[
    {label: '// Patterns', value: '6 core patterns'},
    {label: '// Types', value: 'Result, Either, IOResult'},
    {label: '// Difficulty', value: 'Beginner → Intermediate'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Use Result" prose value={<>For simple <em>success/failure</em> with error.</>} variant="up" />
  <TLDRCard label="// Use Either" prose value={<>For <em>custom error types</em> with context.</>} />
  <TLDRCard label="// Chain operations" prose value={<>Compose with <em>Chain/Bind</em> for pipelines.</>} />
</TLDR>

<Section id="pattern-1" number="01" title="Pattern 1: Basic error" titleAccent="handling.">

Use `Result[A]` when you need simple success/failure handling with Go's `error` type.

<CodeCard file="basic-result.go">
{`package main

import (
    "errors"
    "fmt"
    
    R "github.com/IBM/fp-go/v2/result"
)

// Divide safely returns a Result
func safeDivide(a, b float64) R.Result[float64] {
    if b == 0 {
        return R.Left[float64](errors.New("division by zero"))
    }
    return R.Right[float64](a / b)
}

func main() {
    // Success case
    result1 := safeDivide(10, 2)
    fmt.Println(R.IsRight(result1)) // true
    
    // Error case
    result2 := safeDivide(10, 0)
    fmt.Println(R.IsLeft(result2)) // true
    
    // Extract value with default
    value := R.GetOrElse(func() float64 { return 0 })(result2)
    fmt.Println(value) // 0
}`}
</CodeCard>

</Section>

<Section id="pattern-2" number="02" title="Pattern 2: Chaining" titleAccent="operations.">

Chain multiple operations that can fail, short-circuiting on the first error.

<CodeCard file="chaining.go">
{`package main

import (
    "errors"
    "fmt"
    "strconv"
    
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

// Parse string to int
func parseInt(s string) R.Result[int] {
    val, err := strconv.Atoi(s)
    if err != nil {
        return R.Left[int](err)
    }
    return R.Right[int](val)
}

// Validate positive number
func validatePositive(n int) R.Result[int] {
    if n <= 0 {
        return R.Left[int](errors.New("number must be positive"))
    }
    return R.Right[int](n)
}

// Double the number
func double(n int) R.Result[int] {
    return R.Right[int](n * 2)
}

func processNumber(s string) R.Result[int] {
    return F.Pipe3(
        parseInt(s),
        R.Chain(validatePositive),
        R.Chain(double),
    )
}

func main() {
    // Success: "5" -> 5 -> validate -> 10
    result1 := processNumber("5")
    fmt.Println(R.GetOrElse(func() int { return -1 })(result1)) // 10
    
    // Fails at parsing
    result2 := processNumber("abc")
    fmt.Println(R.IsLeft(result2)) // true
    
    // Fails at validation
    result3 := processNumber("-5")
    fmt.Println(R.IsLeft(result3)) // true
}`}
</CodeCard>

</Section>

<Section id="pattern-3" number="03" title="Pattern 3: Custom error" titleAccent="types.">

Use `Either[E, A]` when you need custom error types with more context.

<CodeCard file="custom-errors.go">
{`package main

import (
    "fmt"
    
    E "github.com/IBM/fp-go/v2/either"
)

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type User struct {
    Name  string
    Email string
    Age   int
}

// Validation functions
func validateName(name string) E.Either[ValidationError, string] {
    if len(name) < 2 {
        return E.Left[string](ValidationError{
            Field:   "name",
            Message: "must be at least 2 characters",
        })
    }
    return E.Right[ValidationError](name)
}

func validateEmail(email string) E.Either[ValidationError, string] {
    if len(email) < 5 || !contains(email, "@") {
        return E.Left[string](ValidationError{
            Field:   "email",
            Message: "invalid email format",
        })
    }
    return E.Right[ValidationError](email)
}

func validateAge(age int) E.Either[ValidationError, int] {
    if age < 18 {
        return E.Left[int](ValidationError{
            Field:   "age",
            Message: "must be at least 18",
        })
    }
    return E.Right[ValidationError](age)
}

// Create user with validation
func createUser(name, email string, age int) E.Either[ValidationError, User] {
    validName := validateName(name)
    if E.IsLeft(validName) {
        return E.Left[User](E.GetLeft(validName))
    }
    
    validEmail := validateEmail(email)
    if E.IsLeft(validEmail) {
        return E.Left[User](E.GetLeft(validEmail))
    }
    
    validAge := validateAge(age)
    if E.IsLeft(validAge) {
        return E.Left[User](E.GetLeft(validAge))
    }
    
    return E.Right[ValidationError](User{
        Name:  E.GetRight(validName),
        Email: E.GetRight(validEmail),
        Age:   E.GetRight(validAge),
    })
}`}
</CodeCard>

</Section>

<Section id="pattern-4" number="04" title="Pattern 4: Error" titleAccent="recovery.">

Provide fallback values or alternative computations when errors occur.

<CodeCard file="recovery.go">
{`package main

import (
    "errors"
    "fmt"
    
    R "github.com/IBM/fp-go/v2/result"
)

// Try to get config from file
func getConfigFromFile() R.Result[string] {
    return R.Left[string](errors.New("file not found"))
}

// Fallback to environment variable
func getConfigFromEnv() R.Result[string] {
    return R.Left[string](errors.New("env var not set"))
}

// Final fallback to default
func getDefaultConfig() R.Result[string] {
    return R.Right[error]("default-config")
}

func getConfig() R.Result[string] {
    return R.OrElse(
        func() R.Result[string] {
            return R.OrElse(
                getConfigFromEnv,
            )(getConfigFromFile())
        },
    )(getDefaultConfig())
}

func main() {
    config := getConfig()
    value := R.GetOrElse(func() string { return "" })(config)
    fmt.Println(value) // "default-config"
}`}
</CodeCard>

</Section>

<Section id="pattern-5" number="05" title="Pattern 5: Mapping" titleAccent="errors.">

Transform error types while preserving the error state.

<CodeCard file="map-errors.go">
{`package main

import (
    "fmt"
    
    E "github.com/IBM/fp-go/v2/either"
)

type DatabaseError struct {
    Code    int
    Message string
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("DB Error %d: %s", e.Code, e.Message)
}

type APIError struct {
    Status  int
    Message string
}

func (e APIError) Error() string {
    return fmt.Sprintf("API Error %d: %s", e.Status, e.Message)
}

// Database operation that returns DatabaseError
func queryDatabase() E.Either[DatabaseError, string] {
    return E.Left[string](DatabaseError{
        Code:    1001,
        Message: "connection timeout",
    })
}

// Convert DatabaseError to APIError
func toAPIError(dbErr DatabaseError) APIError {
    return APIError{
        Status:  500,
        Message: fmt.Sprintf("Database error: %s", dbErr.Message),
    }
}

func handleRequest() E.Either[APIError, string] {
    result := queryDatabase()
    return E.MapLeft(toAPIError)(result)
}

func main() {
    result := handleRequest()
    if E.IsLeft(result) {
        apiErr := E.GetLeft(result)
        fmt.Printf("%s\\n", apiErr.Error())
        // API Error 500: Database error: connection timeout
    }
}`}
</CodeCard>

</Section>

<Section id="pattern-6" number="06" title="Pattern 6: Collecting" titleAccent="errors.">

When you need to collect all errors instead of short-circuiting.

<CodeCard file="collect-errors.go">
{`package main

import (
    "fmt"
    "strings"
    
    E "github.com/IBM/fp-go/v2/either"
)

type ValidationErrors []string

func (ve ValidationErrors) Error() string {
    return strings.Join(ve, "; ")
}

func validateUserData(name, email string, age int) E.Either[ValidationErrors, bool] {
    var errors ValidationErrors
    
    if len(name) < 2 {
        errors = append(errors, "name too short")
    }
    
    if !contains(email, "@") {
        errors = append(errors, "invalid email")
    }
    
    if age < 18 {
        errors = append(errors, "age must be 18+")
    }
    
    if len(errors) > 0 {
        return E.Left[bool](errors)
    }
    
    return E.Right[ValidationErrors](true)
}

func main() {
    // Multiple validation errors
    result := validateUserData("A", "invalid", 15)
    
    if E.IsLeft(result) {
        errors := E.GetLeft(result)
        fmt.Printf("Validation failed: %s\\n", errors.Error())
        // Validation failed: name too short; invalid email; age must be 18+
    }
}`}
</CodeCard>

</Section>
