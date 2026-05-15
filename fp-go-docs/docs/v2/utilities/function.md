---
title: Function Utilities
hide_title: true
description: Core function manipulation utilities including identity, constant, and combinators.
sidebar_position: 19
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Utilities"
  title="Function Utilities"
  lede="Core function manipulation utilities. Essential tools for working with functions in a functional style."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/function' },
    { label: 'Utilities', value: 'Identity, Constant, Swap, First, Second' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Identity` | `func Identity[A any](A) A` | Return input unchanged |
| `Constant` | `func Constant[A any](A) func() A` | Return constant value |
| `Constant1` | `func Constant1[A, B any](B) func(A) B` | Ignore input, return constant |
| `Constant2` | `func Constant2[A, B, C any](C) func(A, B) C` | Ignore inputs, return constant |
| `Swap` | `func Swap[A, B, C any](func(A, B) C) func(B, A) C` | Reverse parameters |
| `First` | `func First[A, B any](A, B) A` | Return first argument |
| `Second` | `func Second[A, B any](A, B) B` | Return second argument |
| `IsNil` | `func IsNil[A any](*A) bool` | Check if pointer is nil |
| `IsNonNil` | `func IsNonNil[A any](*A) bool` | Check if pointer is not nil |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Identity

<CodeCard file="identity.go">
{`import F "github.com/IBM/fp-go/v2/function"

value := F.Identity(42)
// 42

// Useful in higher-order functions
import O "github.com/IBM/fp-go/v2/option"

result := F.Pipe2(
    someOption,
    O.Map(F.Identity[int]),  // No transformation
)

// Conditional transformation
transform := func(shouldTransform bool) func(int) int {
    if shouldTransform {
        return func(n int) int { return n * 2 }
    }
    return F.Identity[int]  // No transformation
}
`}
</CodeCard>

### Constant Functions

<CodeCard file="constant.go">
{`// Constant - returns value (0 args)
getDefault := F.Constant(42)
value := getDefault()
// 42

// Constant1 - ignores 1 input
alwaysTrue := F.Constant1[string, bool](true)
result := alwaysTrue("anything")
// true

// Useful for default handlers
handleError := F.Constant1[error, int](0)

// Constant2 - ignores 2 inputs
defaultValue := F.Constant2[string, int, bool](false)
result := defaultValue("ignored", 42)
// false
`}
</CodeCard>

### Swap Parameters

<CodeCard file="swap.go">
{`// Original function
divide := func(a, b int) int { return a / b }

// Swapped version
divideSwapped := F.Swap(divide)

divide(10, 2)         // 5
divideSwapped(2, 10)  // 5 (parameters swapped)

// API expects (context, id)
fetchUser := func(ctx context.Context, id int) User {
    // ...
}

// But we have (id, context) - swap to match
fetchUserSwapped := F.Swap(fetchUser)
user := fetchUserSwapped(123, ctx)
`}
</CodeCard>

### First & Second

<CodeCard file="first_second.go">
{`// First returns first argument
first := F.First(42, "hello")
// 42

// Second returns second argument
second := F.Second(42, "hello")
// "hello"

// Useful in callbacks
type Handler func(Result, error) Result

onSuccess := func(r Result, _ error) Result {
    return r
}

onError := func(_ Result, err error) Result {
    return Result{Error: err}
}
`}
</CodeCard>

### Nullability Checks

<CodeCard file="nil_checks.go">
{`var ptr *int = nil
var val *int = new(int)

F.IsNil(ptr)     // true
F.IsNil(val)     // false

F.IsNonNil(ptr)  // false
F.IsNonNil(val)  // true

// Only process non-nil values
processIfNotNil := func(ptr *Data) Result {
    if F.IsNonNil(ptr) {
        return process(*ptr)
    }
    return defaultResult
}
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Default Values

<CodeCard file="defaults.go">
{`type Config struct {
    Timeout int
    Retries int
}

// Provide defaults
getTimeout := F.Pipe3(
    config.Timeout,
    O.FromPredicate(func(t int) bool { return t > 0 }),
    O.GetOrElse(F.Constant(30)),
)

// Always return default on error
import R "github.com/IBM/fp-go/v2/result"

safeGet := F.Pipe2(
    dangerousOperation(),
    R.GetOrElse(F.Constant(defaultValue)),
)
`}
</CodeCard>

### Callback Handlers

<CodeCard file="handlers.go">
{`type Handler func(data string, err error) string

// Success handler (ignore error)
onSuccess := func(data string, _ error) string {
    return data
}

// Error handler (ignore data)
onError := func(_ string, err error) string {
    return err.Error()
}

// Using Constant for fixed responses
notFound := F.Constant2[string, error, string]("Not Found")

// Default handler
defaultHandler := F.Constant1[Request, Response](
    Response{Status: 404, Body: "Not Found"},
)
`}
</CodeCard>

### Parameter Reordering

<CodeCard file="reorder.go">
{`// Create specialized version
divideBy10 := func(n int) int {
    return F.Swap(divide)(10, n)
}

divideBy10(2)  // 5
divideBy10(5)  // 2
`}
</CodeCard>

### Optional Transformation

<CodeCard file="optional.go">
{`// Transform only if condition met
maybeTransform := func(condition bool, f func(A) A) func(A) A {
    if condition {
        return f
    }
    return F.Identity[A]
}

result := F.Pipe2(
    value,
    maybeTransform(shouldTransform, transform),
)
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Identity**: Use `Identity` for no-op transformations in pipelines. It's clearer than `func(x T) T { return x }`.

</Callout>

<Callout type="info">

**Constant**: `Constant` functions are useful for providing default values, especially with `GetOrElse` and error handlers.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/utilities/eq', title: 'Eq (Equality)' }}
  next={{ to: '/docs/v2/utilities/magma', title: 'Magma' }}
/>

---
