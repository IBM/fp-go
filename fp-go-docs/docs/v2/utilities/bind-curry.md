---
title: Bind & Curry
hide_title: true
description: Partial application and currying utilities for creating specialized functions in fp-go
sidebar_position: 21
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Bind & Curry"
  titleAccent="Partial Application"
  lede="Bind and curry functions allow you to fix some arguments of a function, creating specialized versions with fewer parameters."
  meta={[
    { label: 'Package', value: 'function' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Technique', value: 'Partial Application' }
  ]}
/>

---

<Section num="1" title="Overview">

**Partial application** is the technique of fixing some arguments of a function, producing a new function with fewer parameters.

**Key Concepts:**
- **Bind1st**: Fix the first argument
- **Bind2nd**: Fix the second argument
- **Currying**: Transform multi-argument function into chain of single-argument functions

<Callout type="info">
Partial application is useful for creating specialized versions of generic functions, reducing code duplication.
</Callout>

</Section>

---

<Section num="2" title="Bind1st - Fix First Argument">

Fix the first argument of a two-argument function:

<CodeCard file="bind_first.go" tag="example">
{`import F "github.com/IBM/fp-go/function"

// Original two-argument function
divide := func(a, b int) int {
    return a / b
}

// Fix first argument to 10
divideBy10 := F.Bind1st(divide, 10)

// Now it's a single-argument function
divideBy10(2)  // 5  (10 / 2)
divideBy10(5)  // 2  (10 / 5)
divideBy10(1)  // 10 (10 / 1)

// Equivalent to:
divide(10, 2)  // 5
divide(10, 5)  // 2
divide(10, 1)  // 10`}
</CodeCard>

</Section>

---

<Section num="3" title="Bind2nd - Fix Second Argument">

Fix the second argument of a two-argument function:

<CodeCard file="bind_second.go" tag="example">
{`import F "github.com/IBM/fp-go/function"

// Original two-argument function
divide := func(a, b int) int {
    return a / b
}

// Fix second argument to 10
divideTenBy := F.Bind2nd(divide, 10)

// Now it's a single-argument function
divideTenBy(100)  // 10  (100 / 10)
divideTenBy(50)   // 5   (50 / 10)
divideTenBy(20)   // 2   (20 / 10)

// Equivalent to:
divide(100, 10)  // 10
divide(50, 10)   // 5
divide(20, 10)   // 2`}
</CodeCard>

</Section>

---

<Section num="4" title="String Operations">

Create specialized string functions:

<CodeCard file="bind_strings.go" tag="example">
{`import (
    F "github.com/IBM/fp-go/function"
    "strings"
)

// Generic string operation
contains := func(haystack, needle string) bool {
    return strings.Contains(haystack, needle)
}

// Create specialized checkers
hasAt := F.Bind2nd(contains, "@")
hasDot := F.Bind2nd(contains, ".")

// Use specialized functions
hasAt("user@example.com")   // true
hasAt("username")           // false

hasDot("file.txt")          // true
hasDot("filename")          // false

// Combine for email validation
isEmail := func(s string) bool {
    return hasAt(s) && hasDot(s)
}

isEmail("user@example.com")  // true
isEmail("invalid")           // false`}
</CodeCard>

</Section>

---

<Section num="5" title="API Client Example">

Create specialized API request functions:

<CodeCard file="bind_api.go" tag="example">
{`type APIClient struct {
    BaseURL string
    Token   string
}

type Response struct {
    Status int
    Body   string
}

// Generic request function
func makeRequest(client *APIClient, endpoint string) Response {
    url := client.BaseURL + endpoint
    // Make HTTP request with client.Token
    return Response{Status: 200, Body: "..."}
}

// Create client
client := &APIClient{
    BaseURL: "https://api.example.com",
    Token:   "secret-token",
}

// Create specialized request function
request := F.Bind1st(makeRequest, client)

// Use with different endpoints
users := request("/users")
posts := request("/posts")
comments := request("/comments")

// Much cleaner than:
// makeRequest(client, "/users")
// makeRequest(client, "/posts")
// makeRequest(client, "/comments")`}
</CodeCard>

</Section>

---

<Section num="6" title="Configuration Pattern">

Fix configuration for processing functions:

<CodeCard file="bind_config.go" tag="example">
{`type Config struct {
    Env       string
    Debug     bool
    MaxRetries int
}

type Data struct {
    ID   int
    Name string
}

type Result struct {
    Success bool
    Output  string
}

func processWithConfig(config Config, data Data) Result {
    // Process data using config settings
    if config.Debug {
        fmt.Printf("Processing %s in %s\n", data.Name, config.Env)
    }
    return Result{Success: true, Output: data.Name}
}

// Create production processor
prodConfig := Config{
    Env:        "production",
    Debug:      false,
    MaxRetries: 3,
}
processInProd := F.Bind1st(processWithConfig, prodConfig)

// Create development processor
devConfig := Config{
    Env:        "development",
    Debug:      true,
    MaxRetries: 1,
}
processInDev := F.Bind1st(processWithConfig, devConfig)

// Process different data with same config
data1 := Data{ID: 1, Name: "Alice"}
data2 := Data{ID: 2, Name: "Bob"}

prodResult1 := processInProd(data1)
prodResult2 := processInProd(data2)

devResult1 := processInDev(data1)
devResult2 := processInDev(data2)`}
</CodeCard>

</Section>

---

<Section num="7" title="Array Operations">

Create specialized array filters and mappers:

<CodeCard file="bind_array.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
)

// Generic comparison function
greaterThan := func(threshold, value int) bool {
    return value > threshold
}

// Create specialized filters
isGreaterThan10 := F.Bind1st(greaterThan, 10)
isGreaterThan100 := F.Bind1st(greaterThan, 100)

numbers := []int{5, 15, 50, 150, 200}

// Filter with specialized predicates
above10 := F.Pipe2(
    numbers,
    A.Filter(isGreaterThan10),
)
// []int{15, 50, 150, 200}

above100 := F.Pipe2(
    numbers,
    A.Filter(isGreaterThan100),
)
// []int{150, 200}

// Generic multiply function
multiply := func(factor, value int) int {
    return factor * value
}

// Create specialized mappers
double := F.Bind1st(multiply, 2)
triple := F.Bind1st(multiply, 3)

doubled := F.Pipe2(numbers, A.Map(double))
// []int{10, 30, 100, 300, 400}

tripled := F.Pipe2(numbers, A.Map(triple))
// []int{15, 45, 150, 450, 600}`}
</CodeCard>

</Section>

---

<Section num="8" title="Logging and Middleware">

Create specialized logging functions:

<CodeCard file="bind_logging.go" tag="example">
{`type Logger struct {
    Prefix string
    Level  string
}

func log(logger *Logger, message string) {
    fmt.Printf("[%s] %s: %s\n", logger.Level, logger.Prefix, message)
}

// Create specialized loggers
errorLogger := &Logger{Prefix: "API", Level: "ERROR"}
infoLogger := &Logger{Prefix: "API", Level: "INFO"}

logError := F.Bind1st(log, errorLogger)
logInfo := F.Bind1st(log, infoLogger)

// Use specialized loggers
logInfo("Server started")
// [INFO] API: Server started

logError("Connection failed")
// [ERROR] API: Connection failed

logInfo("Request processed")
// [INFO] API: Request processed`}
</CodeCard>

</Section>

---

<Section num="9" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `Bind1st[A, B, C]` | `(func(A, B) C, A) -> func(B) C` | Fix first argument |
| `Bind2nd[A, B, C]` | `(func(A, B) C, B) -> func(A) C` | Fix second argument |
</ApiTable>

**Usage Pattern:**
```go
// Original function
f := func(a A, b B) C { ... }

// Fix first argument
g := Bind1st(f, valueA)  // g(b B) C

// Fix second argument
h := Bind2nd(f, valueB)  // h(a A) C
```

</Section>

---

<Section num="10" title="Related Concepts">

**Common Use Cases:**
- Creating specialized versions of generic functions
- Configuration injection
- API client builders
- Logging and middleware
- Reducing code duplication

**See Also:**
- [Function](./function.md) - Core function utilities
- [Pipe & Flow](./pipe-flow.md) - Function composition
- [Compose](./compose.md) - Mathematical composition

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/compose', title: 'Compose' }}
  next={{ to: '/docs/v2/utilities/eq', title: 'Eq' }}
/>