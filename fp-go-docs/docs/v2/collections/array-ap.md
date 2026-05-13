---
title: Array - Applicative
hide_title: true
description: Applicative operations for arrays - apply arrays of functions to arrays of values.
sidebar_position: 9
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Applicative"
  lede="Applicative operations for arrays. Apply arrays of functions to arrays of values, producing all combinations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Operations', value: 'Ap, Flap' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[A, B any](fns []func(A) B) func([]A) []B` | Apply functions to values |
| `Flap` | `func Flap[A, B any](value A) func([]func(A) B) []B` | Apply value to functions |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Ap - Apply Functions

<CodeCard file="ap.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

// Array of functions
fns := []func(int) int{
    func(n int) int { return n * 2 },
    func(n int) int { return n + 10 },
}

// Array of values
values := []int{1, 2, 3}

// Apply all functions to all values
result := F.Pipe2(
    values,
    A.Ap(fns),
)
// []int{2, 4, 6, 11, 12, 13}
// (1*2, 2*2, 3*2, 1+10, 2+10, 3+10)
`}
</CodeCard>

### Flap - Apply Value

<CodeCard file="flap.go">
{`fns := []func(int) string{
    func(n int) string { return fmt.Sprintf("Double: %d", n*2) },
    func(n int) string { return fmt.Sprintf("Square: %d", n*n) },
}

result := F.Pipe2(
    fns,
    A.Flap(5),
)
// []string{"Double: 10", "Square: 25"}
`}
</CodeCard>

### Validation Example

<CodeCard file="validation.go">
{`type Validator func(string) bool

validators := []Validator{
    func(s string) bool { return len(s) > 5 },
    func(s string) bool { return strings.Contains(s, "@") },
}

inputs := []string{"test", "user@example.com", "admin"}

// Apply all validators to all inputs
results := F.Pipe2(
    inputs,
    A.Ap(validators),
)
// []bool{false, true, false, false, true, false}
`}
</CodeCard>

</Section>
