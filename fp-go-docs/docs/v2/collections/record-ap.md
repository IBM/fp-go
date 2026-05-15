---
title: Record - Applicative
hide_title: true
description: Applying maps of functions to maps of values with applicative operations.
sidebar_position: 17
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Record"
  titleAccent="Applicative"
  lede="Applying maps of functions to maps of values. Applicative operations enable function application across map structures."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'Ap, Flap' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[K comparable, A, B any](Monoid[K, B]) func(map[K]func(A) B) func(map[K]A) map[K]B` | Apply functions to values |
| `Flap` | `func Flap[K comparable, A, B any](A) func(map[K]func(A) B) map[K]B` | Apply value to functions |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Ap - Apply Functions

<CodeCard file="ap.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    M "github.com/IBM/fp-go/v2/monoid"
    F "github.com/IBM/fp-go/v2/function"
)

fns := map[string]func(int) int{
    "double": func(n int) int { return n * 2 },
    "square": func(n int) int { return n * n },
}

values := map[string]int{
    "a": 5,
    "b": 3,
}

// Apply all functions to all values
result := F.Pipe2(
    values,
    R.Ap(M.MergeMonoid[string, int]())(fns),
)
// Result depends on merge strategy
// Combines function and value maps
`}
</CodeCard>

### Flap - Apply Value

<CodeCard file="flap.go">
{`fns := map[string]func(int) string{
    "double": func(n int) string {
        return fmt.Sprintf("Double: %d", n*2)
    },
    "square": func(n int) string {
        return fmt.Sprintf("Square: %d", n*n)
    },
    "triple": func(n int) string {
        return fmt.Sprintf("Triple: %d", n*3)
    },
}

// Apply single value to all functions
result := F.Pipe2(
    fns,
    R.Flap(5),
)
// map[string]string{
//   "double": "Double: 10",
//   "square": "Square: 25",
//   "triple": "Triple: 15",
// }
`}
</CodeCard>

### Multiple Transformations

<CodeCard file="multiple.go">
{`type Validator func(string) bool

validators := map[string]Validator{
    "length":    func(s string) bool { return len(s) >= 5 },
    "uppercase": func(s string) bool { return s == strings.ToUpper(s) },
    "numeric":   func(s string) bool { return regexp.MustCompile(\`^[0-9]+$\`).MatchString(s) },
}

// Test a value against all validators
results := F.Pipe2(
    validators,
    R.Flap("HELLO"),
)
// map[string]bool{
//   "length": true,
//   "uppercase": true,
//   "numeric": false,
// }
`}
</CodeCard>

### Formatting Pipeline

<CodeCard file="format.go">
{`type Formatter func(float64) string

formatters := map[string]Formatter{
    "currency": func(n float64) string {
        return fmt.Sprintf("$%.2f", n)
    },
    "percent": func(n float64) string {
        return fmt.Sprintf("%.1f%%", n*100)
    },
    "scientific": func(n float64) string {
        return fmt.Sprintf("%.2e", n)
    },
}

value := 0.12345

formatted := F.Pipe2(
    formatters,
    R.Flap(value),
)
// map[string]string{
//   "currency": "$0.12",
//   "percent": "12.3%",
//   "scientific": "1.23e-01",
// }
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Validation Suite

<CodeCard file="validation.go">
{`type ValidationRule func(User) error

rules := map[string]ValidationRule{
    "email": func(u User) error {
        if !strings.Contains(u.Email, "@") {
            return errors.New("invalid email")
        }
        return nil
    },
    "age": func(u User) error {
        if u.Age < 18 {
            return errors.New("must be 18+")
        }
        return nil
    },
}

user := User{Email: "test@example.com", Age: 25}

// Run all validations
results := F.Pipe2(
    rules,
    R.Flap(user),
)
// map[string]error{
//   "email": nil,
//   "age": nil,
// }
`}
</CodeCard>

### Data Transformers

<CodeCard file="transformers.go">
{`transformers := map[string]func(string) string{
    "lowercase": strings.ToLower,
    "uppercase": strings.ToUpper,
    "title":     strings.Title,
}

input := "hello world"

transformed := F.Pipe2(
    transformers,
    R.Flap(input),
)
// map[string]string{
//   "lowercase": "hello world",
//   "uppercase": "HELLO WORLD",
//   "title": "Hello World",
// }
`}
</CodeCard>

</Section>

<Callout type="info">

**Applicative Pattern**: `Flap` is particularly useful when you have a collection of functions and want to apply a single value to all of them, producing a map of results.

</Callout>
