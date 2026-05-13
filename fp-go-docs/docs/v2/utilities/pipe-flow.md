---
title: Pipe & Flow
hide_title: true
description: Function composition for readable data transformations with left-to-right execution.
sidebar_position: 20
---

import { PageHeader, Section, CodeCard, ApiTable, Compare, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Utilities"
  title="Pipe & Flow"
  lede="Function composition for readable data transformations. Pipe and Flow enable left-to-right composition, making code more readable than nested function calls."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/function' },
    { label: 'Operations', value: 'Pipe2-9, Flow2-9' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Pipe2` | `func Pipe2[A, B, C any](A, func(A) B, func(B) C) C` | Apply value through 2 functions |
| `Pipe3` | `func Pipe3[A, B, C, D any](A, func(A) B, func(B) C, func(C) D) D` | Apply through 3 functions |
| `Pipe4-9` | Similar pattern | Up to 9 functions |
| `Flow2` | `func Flow2[A, B, C any](func(A) B, func(B) C) func(A) C` | Compose 2 functions |
| `Flow3` | `func Flow3[A, B, C, D any](func(A) B, func(B) C, func(C) D) func(A) D` | Compose 3 functions |
| `Flow4-9` | Similar pattern | Up to 9 functions |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Pipe - Immediate Execution

<CodeCard file="pipe.go">
{`import F "github.com/IBM/fp-go/v2/function"

// Instead of: g(f(x))
result := F.Pipe2(
    x,
    f,  // First transformation
    g,  // Second transformation
)

// Instead of: h(g(f(x)))
result := F.Pipe3(
    x,
    f,
    g,
    h,
)

// Up to Pipe9 available
result := F.Pipe5(x, f1, f2, f3, f4, f5)
`}
</CodeCard>

### Flow - Reusable Pipeline

<CodeCard file="flow.go">
{`// Create composed function
transform := F.Flow2(
    f,  // A -> B
    g,  // B -> C
)   // Returns: A -> C

// Use later
result := transform(x)

// Compose multiple functions
pipeline := F.Flow4(
    parse,      // string -> int
    validate,   // int -> Result[int]
    transform,  // Result[int] -> Result[string]
    format,     // Result[string] -> string
)

output := pipeline(input)
`}
</CodeCard>

### Data Processing Pipeline

<CodeCard file="data_pipeline.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
)

type User struct {
    Name  string
    Age   int
    Email string
}

users := []User{
    {Name: "Alice", Age: 30, Email: "alice@example.com"},
    {Name: "Bob", Age: 17, Email: ""},
    {Name: "Charlie", Age: 25, Email: "charlie@example.com"},
}

// Process users: filter adults, extract emails, remove empty
emails := F.Pipe3(
    users,
    A.Filter(func(u User) bool { return u.Age >= 18 }),
    A.Map(func(u User) string { return u.Email }),
    A.Filter(func(e string) bool { return e != "" }),
)
// []string{"alice@example.com", "charlie@example.com"}
`}
</CodeCard>

### String Transformation

<CodeCard file="string.go">
{`import S "github.com/IBM/fp-go/v2/string"

// Clean and format string
cleanString := F.Flow3(
    strings.TrimSpace,
    strings.ToLower,
    func(s string) string {
        return strings.ReplaceAll(s, " ", "-")
    },
)

slug := cleanString("  Hello World  ")
// "hello-world"
`}
</CodeCard>

### API Request Processing

<CodeCard file="api.go">
{`type Request struct {
    Body string
}

type Response struct {
    Status int
    Data   string
}

// Create processing pipeline
processRequest := F.Flow4(
    parseBody,      // Request -> Result[Data]
    validateData,   // Result[Data] -> Result[Data]
    transformData,  // Result[Data] -> Result[string]
    formatResponse, // Result[string] -> Response
)

// Use pipeline
response := processRequest(request)
`}
</CodeCard>

### Validation Chain

<CodeCard file="validation.go">
{`import R "github.com/IBM/fp-go/v2/result"

type ValidationError struct {
    Field   string
    Message string
}

// Compose validators
validateUser := F.Flow3(
    validateEmail,              // User -> Result[User]
    R.Chain(validateAge),       // Result[User] -> Result[User]
    R.Chain(validateName),      // Result[User] -> Result[User]
)

result := validateUser(user)
// Result[User] - Success or first error
`}
</CodeCard>

### Complex Product Processing

<CodeCard file="products.go">
{`type Product struct {
    ID    int
    Name  string
    Price float64
    Tags  []string
}

// Process products
result := F.Pipe5(
    products,
    // Filter by price range
    A.Filter(func(p Product) bool {
        return p.Price >= 10 && p.Price <= 100
    }),
    // Sort by price
    A.Sort(productPriceOrd),
    // Add discount
    A.Map(func(p Product) Product {
        p.Price = p.Price * 0.9
        return p
    }),
    // Extract names
    A.Map(func(p Product) string { return p.Name }),
    // Take first 10
    A.Slice(0, 10),
)
`}
</CodeCard>

### Error Handling Pipeline

<CodeCard file="error_handling.go">
{`// Safe division with validation
safeDivide := F.Flow3(
    // Validate denominator
    func(args [2]int) R.Result[[2]int] {
        if args[1] == 0 {
            return R.Error[[2]int](errors.New("division by zero"))
        }
        return R.Success(args)
    },
    // Perform division
    R.Map(func(args [2]int) float64 {
        return float64(args[0]) / float64(args[1])
    }),
    // Round result
    R.Map(func(f float64) float64 {
        return math.Round(f*100) / 100
    }),
)

result := safeDivide([2]int{10, 3})
// Success(3.33)
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Pipe vs Flow

<CodeCard file="comparison.go">
{`// Pipe - immediate execution with value
result := F.Pipe3(
    initialValue,
    step1,
    step2,
    step3,
)

// Flow - create reusable pipeline
pipeline := F.Flow3(
    step1,
    step2,
    step3,
)

// Use multiple times
result1 := pipeline(value1)
result2 := pipeline(value2)
`}
</CodeCard>

### Optional Transformation

<CodeCard file="optional.go">
{`// Transform if present
result := F.Pipe3(
    maybeValue,
    O.Map(transform),
    O.GetOrElse(F.Constant(defaultValue)),
)
`}
</CodeCard>

### Array Processing Pattern

<CodeCard file="array_pattern.go">
{`// Filter, map, reduce pattern
total := F.Pipe3(
    items,
    A.Filter(predicate),
    A.Map(extract),
    A.Reduce(sum, 0),
)
`}
</CodeCard>

### Conditional Pipeline

<CodeCard file="conditional.go">
{`// Different paths based on condition
process := func(useAdvanced bool) func(Data) Result {
    if useAdvanced {
        return F.Flow3(validate, advancedTransform, format)
    }
    return F.Flow2(validate, simpleTransform)
}

result := process(true)(data)
`}
</CodeCard>

</Section>

---

<Callout type="info">

**When to Use**:
- **Pipe**: When you have a value and want to transform it immediately
- **Flow**: When you want to create a reusable transformation pipeline

</Callout>

<Callout type="info">

**Readability**: Pipe and Flow make code more readable by showing the transformation flow from left to right, matching how we naturally read code.

</Callout>

<Callout type="warn">

**Performance**: Each Pipe/Flow call has minimal overhead. For performance-critical code with simple transformations, direct function calls may be faster.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/utilities/ord', title: 'Ord (Ordering)' }}
  next={{ to: '/docs/v2/collections/nonempty-array', title: 'NonEmpty Array' }}
/>

