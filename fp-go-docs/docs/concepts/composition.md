---
sidebar_position: 12
title: Composition
hide_title: true
description: Build complex functionality from simple, reusable functions — the essence of functional programming.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Concepts · 03 / 06"
  title="Composition."
  lede="Combine simple, reusable functions into pipelines. The essence of functional programming, applied to Go."
  meta={[
    {label: '// Difficulty', value: 'Beginner → Intermediate'},
    {label: '// Reading time', value: '10 min · 9 sections'},
    {label: '// Prereqs', value: 'Pure functions'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Tool of choice" prose value={<><em>Flow</em> — left-to-right pipelines.</>} variant="up" />
  <TLDRCard label="// Building block" prose value={<>Small, <em>pure</em> functions.</>} />
  <TLDRCard label="// Win" prose value={<>Reusable, <em>testable</em>, easy to modify.</>} />
</TLDR>

<Section id="what" number="01" title="What is" titleAccent="composition?">

**Composition** is combining simple functions to create more complex ones. Think of it like LEGO blocks — each block is simple, blocks connect in standard ways, complex structures emerge from simple pieces.

<CodeCard file="lego.go" status="tested">
{`// Simple functions (LEGO blocks)
func double(x int) int { return x * 2 }
func addOne(x int) int { return x + 1 }
func square(x int) int { return x * x }

// Compose them (build something)
result := square(addOne(double(5)))  // ((5*2)+1)^2 = 121`}
</CodeCard>

</Section>

<Section id="why" number="02" title="Why composition" titleAccent="matters.">

<Tabs groupId="composition-why">
<TabItem value="without" label="Without composition">

<CodeCard file="monolithic.go">
{`// Monolithic function - hard to reuse
func processUserData(data string) string {
    // Parse
    parsed := strings.TrimSpace(data)

    // Validate
    if len(parsed) == 0 {
        return ""
    }

    // Transform
    lower := strings.ToLower(parsed)

    // Format
    return fmt.Sprintf("user_%s", lower)
}

// Can't reuse individual steps
// Hard to test each step
// Hard to modify`}
</CodeCard>

</TabItem>
<TabItem value="with" label="With composition">

<CodeCard file="composed.go" status="tested">
{`// Small, reusable functions
func trim(s string) string { return strings.TrimSpace(s) }
func lower(s string) string { return strings.ToLower(s) }
func addPrefix(prefix string) func(string) string {
    return func(s string) string {
        return prefix + s
    }
}

// Compose them
import "github.com/IBM/fp-go/v2/function"

processUserData := function.Flow3(
    trim,
    lower,
    addPrefix("user_"),
)

result := processUserData("  JOHN  ")  // "user_john"

// Each step is reusable
// Each step is testable
// Easy to modify pipeline`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="patterns" number="03" title="Composition" titleAccent="patterns.">

### Pattern 1 — Manual composition

<Tabs groupId="composition">
<TabItem value="nested" label="Nested calls">

<CodeCard file="nested.go">
{`// Right-to-left (inside-out)
result := square(addOne(double(5)))
//         ^      ^      ^
//         3rd    2nd    1st

// Execution order: double → addOne → square
// (5*2) = 10
// (10+1) = 11
// (11*11) = 121`}
</CodeCard>

<Compare>
  <CompareCol kind="good" title="Pros" pill="simple">
    <p>No dependencies.</p>
    <p>Just function calls.</p>
  </CompareCol>
  <CompareCol kind="bad" title="Cons" pill="watch out">
    <p>Hard to read with many functions.</p>
  </CompareCol>
</Compare>

</TabItem>
<TabItem value="intermediate" label="Intermediate variables">

<CodeCard file="intermediate.go">
{`// Step by step
step1 := double(5)    // 10
step2 := addOne(step1) // 11
step3 := square(step2) // 121

// Execution order: clear`}
</CodeCard>

<Compare>
  <CompareCol kind="good" title="Pros" pill="readable">
    <p>Easy to debug.</p>
    <p>Clear order.</p>
  </CompareCol>
  <CompareCol kind="bad" title="Cons" pill="trade-offs">
    <p>Verbose.</p>
    <p>Temporary variables.</p>
  </CompareCol>
</Compare>

</TabItem>
</Tabs>

### Pattern 2 — Compose (right-to-left)

<CodeCard file="compose.go" status="tested">
{`import "github.com/IBM/fp-go/v2/function"

// Mathematical composition: (f ∘ g)(x) = f(g(x))
composed := function.Compose3(
    square,   // Applied LAST (3rd)
    addOne,   // Applied 2nd
    double,   // Applied FIRST (1st)
)

result := composed(5)  // 121

// Execution: double(5) → addOne(10) → square(11)`}
</CodeCard>

<Callout title="When to reach for Compose.">
  Mathematical, concise. But right-to-left can confuse Go readers — prefer <strong>Flow</strong> unless you specifically want the mathematical convention.
</Callout>

### Pattern 3 — Flow (left-to-right) ⭐

<CodeCard file="flow.go" status="tested">
{`import "github.com/IBM/fp-go/v2/function"

// Pipeline style: data flows left to right
pipeline := function.Flow3(
    double,   // Applied FIRST (1st)
    addOne,   // Applied 2nd
    square,   // Applied LAST (3rd)
)

result := pipeline(5)  // 121

// Execution: double(5) → addOne(10) → square(11)`}
</CodeCard>

<Callout type="success" title="Recommended.">
  Intuitive, reads like a pipeline. This is the default choice.
</Callout>

### Pattern 4 — Pipe (data-first) ⭐

<CodeCard file="pipe.go" status="tested">
{`import "github.com/IBM/fp-go/v2/function"

// Start with data, pipe through functions
result := function.Pipe3(
    5,        // Start with data
    double,   // 10
    addOne,   // 11
    square,   // 121
)`}
</CodeCard>

<Callout type="success" title="When to use Pipe.">
  Very clear, data-first. Best for one-off processing where you don't need to reuse the pipeline.
</Callout>

</Section>

<Section id="api" number="04" title="fp-go composition" titleAccent="functions.">

### Flow — recommended

Create reusable pipelines.

<CodeCard file="flow-signatures.go">
{`// Flow2 - 2 functions
func Flow2[A, B, C any](
    f func(A) B,
    g func(B) C,
) func(A) C

// Flow3 - 3 functions
func Flow3[A, B, C, D any](
    f func(A) B,
    g func(B) C,
    h func(C) D,
) func(A) D

// Up to Flow9`}
</CodeCard>

<CodeCard file="flow-example.go" status="tested">
{`// Create pipeline
processNumber := function.Flow4(
    func(x int) int { return x * 2 },      // double
    func(x int) int { return x + 1 },      // add one
    func(x int) int { return x * x },      // square
    func(x int) string { return fmt.Sprintf("Result: %d", x) },
)

// Use it
output := processNumber(5)  // "Result: 121"`}
</CodeCard>

### Pipe — data-first

Process data through a pipeline.

<CodeCard file="pipe-signatures.go">
{`// Pipe2 - data + 2 functions
func Pipe2[A, B, C any](
    a A,
    f func(A) B,
    g func(B) C,
) C

// Pipe3 - data + 3 functions
func Pipe3[A, B, C, D any](
    a A,
    f func(A) B,
    g func(B) C,
    h func(C) D,
) D

// Up to Pipe9`}
</CodeCard>

<CodeCard file="pipe-example.go" status="tested">
{`result := function.Pipe4(
    "  HELLO WORLD  ",
    strings.TrimSpace,
    strings.ToLower,
    func(s string) string { return strings.ReplaceAll(s, " ", "_") },
    func(s string) string { return "slug_" + s },
)
// "slug_hello_world"`}
</CodeCard>

### Compose — mathematical

Right-to-left composition.

<CodeCard file="compose-signatures.go">
{`// Compose2 - 2 functions (right-to-left)
func Compose2[A, B, C any](
    f func(B) C,  // Applied SECOND
    g func(A) B,  // Applied FIRST
) func(A) C

// Example
composed := function.Compose2(
    square,   // Applied second
    double,   // Applied first
)
result := composed(5)  // square(double(5)) = 100`}
</CodeCard>

<Callout title="Reminder.">
  Use Flow instead for better Go readability.
</Callout>

</Section>

<Section id="examples" number="05" title="Real-world" titleAccent="examples.">

### String processing

<Tabs groupId="example">
<TabItem value="standard" label="Without fp-go">

<CodeCard file="without.go">
{`func processUsername(input string) string {
    // Step 1: trim
    trimmed := strings.TrimSpace(input)

    // Step 2: lowercase
    lower := strings.ToLower(trimmed)

    // Step 3: remove special chars
    cleaned := regexp.MustCompile(\`[^a-z0-9]\`).ReplaceAllString(lower, "")

    // Step 4: limit length
    if len(cleaned) > 20 {
        cleaned = cleaned[:20]
    }

    return cleaned
}`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="With fp-go v2">

<CodeCard file="with.go" status="tested">
{`import "github.com/IBM/fp-go/v2/function"

var (
    trim = strings.TrimSpace
    lower = strings.ToLower
    removeSpecial = func(s string) string {
        return regexp.MustCompile(\`[^a-z0-9]\`).ReplaceAllString(s, "")
    }
    limitLength = func(max int) func(string) string {
        return func(s string) string {
            if len(s) > max {
                return s[:max]
            }
            return s
        }
    }
)

var processUsername = function.Flow4(
    trim,
    lower,
    removeSpecial,
    limitLength(20),
)

// Usage
username := processUsername("  John.Doe@123!  ")  // "johndoe123"`}
</CodeCard>

</TabItem>
</Tabs>

### Data transformation

<Tabs groupId="example">
<TabItem value="standard" label="Without fp-go">

<CodeCard file="without.go">
{`func processOrders(orders []Order) []OrderSummary {
    // Filter active orders
    var active []Order
    for _, order := range orders {
        if order.Status == "active" {
            active = append(active, order)
        }
    }

    // Calculate totals
    var withTotals []Order
    for _, order := range active {
        total := 0.0
        for _, item := range order.Items {
            total += item.Price
        }
        order.Total = total
        withTotals = append(withTotals, order)
    }

    // Convert to summaries
    var summaries []OrderSummary
    for _, order := range withTotals {
        summaries = append(summaries, OrderSummary{
            ID:    order.ID,
            Total: order.Total,
        })
    }

    return summaries
}`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="With fp-go v2">

<CodeCard file="with.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/array"
    "github.com/IBM/fp-go/v2/function"
)

var (
    isActive = func(o Order) bool {
        return o.Status == "active"
    }

    calculateTotal = func(o Order) Order {
        total := 0.0
        for _, item := range o.Items {
            total += item.Price
        }
        o.Total = total
        return o
    }

    toSummary = func(o Order) OrderSummary {
        return OrderSummary{
            ID:    o.ID,
            Total: o.Total,
        }
    }
)

var processOrders = function.Flow3(
    array.Filter(isActive),
    array.Map(calculateTotal),
    array.Map(toSummary),
)

// Usage
summaries := processOrders(orders)`}
</CodeCard>

</TabItem>
</Tabs>

### API response processing

<Tabs groupId="example">
<TabItem value="standard" label="Without fp-go">

<CodeCard file="without.go">
{`func processAPIResponse(data []byte) (Result, error) {
    // Parse JSON
    var raw RawResponse
    if err := json.Unmarshal(data, &raw); err != nil {
        return Result{}, err
    }

    // Validate
    if raw.Status != "success" {
        return Result{}, errors.New("invalid status")
    }

    // Transform
    result := Result{
        ID:   raw.Data.ID,
        Name: strings.ToUpper(raw.Data.Name),
    }

    return result, nil
}`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="With fp-go v2">

<CodeCard file="with.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/result"
)

var (
    parseJSON = func(data []byte) result.Result[RawResponse] {
        var raw RawResponse
        err := json.Unmarshal(data, &raw)
        return result.FromGoError(raw, err)
    }

    validateStatus = func(raw RawResponse) result.Result[RawResponse] {
        if raw.Status != "success" {
            return result.Err[RawResponse](errors.New("invalid status"))
        }
        return result.Ok(raw)
    }

    transform = func(raw RawResponse) Result {
        return Result{
            ID:   raw.Data.ID,
            Name: strings.ToUpper(raw.Data.Name),
        }
    }
)

var processAPIResponse = function.Flow3(
    parseJSON,
    result.Chain(validateStatus),
    result.Map(transform),
)

// Usage
res := processAPIResponse(data)
res.Fold(
    func(err error) { /* handle error */ },
    func(result Result) { /* use result */ },
)`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="monads" number="06" title="Composition with" titleAccent="monads." lede="Compose effectful functions by lifting them with Map / Chain.">

### Chaining operations

<CodeCard file="chain.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/function"
)

// Each function returns Result
func fetchUser(id string) result.Result[User] { /* ... */ }
func validateUser(user User) result.Result[User] { /* ... */ }
func enrichUser(user User) result.Result[User] { /* ... */ }
func saveUser(user User) result.Result[User] { /* ... */ }

// Compose with Chain
var processUser = function.Flow3(
    fetchUser,
    result.Chain(validateUser),
    result.Chain(enrichUser),
    result.Chain(saveUser),
)

// Usage
res := processUser("user-123")
// Stops at first error, or returns final success`}
</CodeCard>

### Mixing Map and Chain

<CodeCard file="mix.go">
{`// Map: transform the value
// Chain: transform and return new Result

var pipeline = function.Flow4(
    fetchUser,                           // Result[User]
    result.Map(normalizeUser),           // Result[User] - just transform
    result.Chain(validateUser),          // Result[User] - can fail
    result.Map(toDTO),                   // Result[UserDTO] - just transform
)`}
</CodeCard>

</Section>

<Section id="point-free" number="07" title="Point-free" titleAccent="style." lede="Define functions without naming their arguments.">

### With points (arguments)

<CodeCard file="with-points.go">
{`// Mentions 'x' explicitly
double := func(x int) int {
    return x * 2
}

// Mentions 'users' explicitly
activeUsers := func(users []User) []User {
    return array.Filter(func(u User) bool {
        return u.Active
    })(users)
}`}
</CodeCard>

### Point-free

<CodeCard file="point-free.go">
{`// No mention of arguments
var double = function.Flow1(func(x int) int { return x * 2 })

// No mention of 'users'
var activeUsers = array.Filter(func(u User) bool {
    return u.Active
})

// Use it
result := activeUsers(users)`}
</CodeCard>

<Compare>
  <CompareCol kind="good" title="Use point-free when" pill="cleaner">
    <p>The pipeline is clearer without naming intermediate values.</p>
    <CodeCard file="clear.go">
{`var processUsers = function.Flow3(
    array.Filter(isActive),
    array.Map(normalize),
    array.Map(toDTO),
)`}
    </CodeCard>
  </CompareCol>
  <CompareCol kind="bad" title="Don't force when" pill="harder to read">
    <p>Logic is complex — a normal function is clearer.</p>
    <CodeCard file="unclear.go">
{`var processUser = func(user User) User {
    // Complex logic here
    return user
}`}
    </CodeCard>
  </CompareCol>
</Compare>

</Section>

<Section id="best-practices" number="08" title="Best" titleAccent="practices.">

### 1 — Keep functions small

<Compare>
  <CompareCol kind="good" pill="recommended">
    <CodeCard file="good.go">
{`func trim(s string) string { return strings.TrimSpace(s) }
func lower(s string) string { return strings.ToLower(s) }
func addPrefix(p string) func(string) string {
    return func(s string) string { return p + s }
}`}
    </CodeCard>
  </CompareCol>
  <CompareCol kind="bad" pill="avoid">
    <CodeCard file="bad.go">
{`func processString(s string, prefix string, maxLen int) string {
    s = strings.TrimSpace(s)
    s = strings.ToLower(s)
    s = prefix + s
    if len(s) > maxLen {
        s = s[:maxLen]
    }
    return s
}`}
    </CodeCard>
  </CompareCol>
</Compare>

### 2 — Use descriptive names

<Compare>
  <CompareCol kind="good" pill="clear">
    <CodeCard file="good.go">
{`var processUsername = function.Flow3(
    trim,
    lower,
    removeSpecialChars,
)`}
    </CodeCard>
  </CompareCol>
  <CompareCol kind="bad" pill="unclear">
    <CodeCard file="bad.go">
{`var process = function.Flow3(
    f1,
    f2,
    f3,
)`}
    </CodeCard>
  </CompareCol>
</Compare>

### 3 — Limit pipeline length

<CodeCard file="limit.go">
{`// ✅ Good: 3-5 steps
var pipeline = function.Flow4(
    parse,
    validate,
    transform,
    format,
)

// ⚠️ Consider breaking up: 10+ steps
var hugePipeline = function.Flow10(
    step1, step2, step3, step4, step5,
    step6, step7, step8, step9, step10,
)

// Better: break into sub-pipelines
var preprocess = function.Flow3(step1, step2, step3)
var process = function.Flow3(step4, step5, step6)
var postprocess = function.Flow4(step7, step8, step9, step10)

var pipeline = function.Flow3(preprocess, process, postprocess)`}
</CodeCard>

### 4 — Test each function

<CodeCard file="testing.go" status="tested">
{`// Test individual functions
func TestTrim(t *testing.T) {
    assert.Equal(t, "hello", trim("  hello  "))
}

func TestLower(t *testing.T) {
    assert.Equal(t, "hello", lower("HELLO"))
}

// Test composition
func TestProcessUsername(t *testing.T) {
    assert.Equal(t, "user_john", processUsername("  JOHN  "))
}`}
</CodeCard>

### 5 — Use type aliases for clarity

<CodeCard file="aliases.go">
{`// Define types for clarity
type Username string
type Email string
type UserID string

// Functions with clear types
func normalizeUsername(s string) Username {
    return Username(strings.ToLower(strings.TrimSpace(s)))
}

func validateEmail(s string) result.Result[Email] {
    if !strings.Contains(s, "@") {
        return result.Err[Email](errors.New("invalid email"))
    }
    return result.Ok(Email(s))
}`}
</CodeCard>

</Section>

<Section id="performance" number="09" title="Performance" titleAccent="considerations.">

### Composition overhead

<CodeCard file="overhead.go">
{`// Minimal overhead
composed := function.Flow3(f1, f2, f3)
result := composed(input)

// Equivalent to:
result := f3(f2(f1(input)))

// Just function calls - very fast`}
</CodeCard>

### When to optimize

<CodeCard file="optimize.go">
{`// ✅ Fine for most cases
var process = function.Flow5(
    step1, step2, step3, step4, step5,
)

// ⚠️ Hot path with millions of calls?
// Consider inlining:
func process(x int) int {
    x = step1(x)
    x = step2(x)
    x = step3(x)
    x = step4(x)
    return step5(x)
}`}
</CodeCard>

<CodeCard file="bench.go">
{`func BenchmarkComposed(b *testing.B) {
    composed := function.Flow3(double, addOne, square)
    for i := 0; i < b.N; i++ {
        _ = composed(5)
    }
}

func BenchmarkInlined(b *testing.B) {
    for i := 0; i < b.N; i++ {
        x := 5
        x = double(x)
        x = addOne(x)
        _ = square(x)
    }
}

// Results: negligible difference for most cases`}
</CodeCard>

</Section>

<Section id="faq" number="10" title="Common" titleAccent="questions.">

<Callout title="Isn't this just function calls?">
  Yes. Composition is a structured way to call functions. The value is in <strong>reusability</strong>, <strong>testability</strong>, <strong>clarity</strong>, and <strong>maintainability</strong>.
</Callout>

<Callout title="When should I use composition?">
  Use it when you have a sequence of transformations, the steps are reusable, and you want clear/testable code. Don't force it when logic is complex and branching, steps are tightly coupled, or it makes the code less clear.
</Callout>

<Callout type="success" title="Flow vs. Pipe vs. Compose — quick rule.">
  <ul>
    <li><strong>Flow</strong> — create reusable pipelines.</li>
    <li><strong>Pipe</strong> — one-off data processing.</li>
    <li><strong>Compose</strong> — only if you prefer mathematical style.</li>
  </ul>
  Most people prefer Flow.
</Callout>

</Section>

<Section id="summary" number="11" title="Summary">

<Checklist
  title="Composition"
  items={[
    {label: 'Build complex from simple', done: true},
    {label: 'Reusable functions', done: true},
    {label: 'Clear data flow', done: true},
    {label: 'Easy to test', done: true},
    {label: 'Easy to modify', done: true},
  ]}
/>

<ApiTable
  columns={['Tool', 'Direction', 'Use for']}
  rows={[
    {symbol: 'Flow', signature: 'Left-to-right', description: 'Reusable pipelines (recommended).'},
    {symbol: 'Pipe', signature: 'Data-first', description: 'One-off processing of a known value.'},
    {symbol: 'Compose', signature: 'Right-to-left', description: 'Mathematical convention.'},
  ]}
/>

<Callout type="success" title="Key takeaway.">
  Composition is about building maintainable systems from simple, reusable pieces. Use it where it adds clarity, not complexity.
</Callout>

</Section>
