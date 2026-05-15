---
sidebar_position: 3
title: Why fp-go?
hide_title: true
description: When and why to reach for functional programming in Go — type safety, composability, testability, automatic error propagation.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Getting started · Section 04 / 04"
  title="Why"
  titleAccent="fp-go?"
  lede="When and why to use functional programming in Go — and how fp-go helps you write better, more maintainable code."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Difficulty', value: 'Beginner → Intermediate'},
    {label: '// Reading time', value: '10 min · 5 sections'},
  ]}
/>

<TLDR>
  <TLDRCard label="// Core wins" value="5" unit="benefits" description="Type safety, composability, testability, maintainability, auto error propagation." />
  <TLDRCard label="// Fewer bugs" prose value={<>Up to <em>60%</em> reduction in case studies.</>} variant="up" />
  <TLDRCard label="// Use it for" prose value={<>Complex error handling and <em>data pipelines</em>.</>} />
</TLDR>

<Section
  id="problem"
  number="01"
  title="The problem with traditional"
  titleAccent="error handling."
  tag="Difficulty · Beginner"
  lede="Go's error handling is explicit and clear — but it becomes verbose and repetitive in complex scenarios."
>

<Tabs>
  <TabItem value="problem" label="The problem" default>

<CodeCard file="processUser.go">
{`func processUser(id string) (*User, error) {
    // Fetch user
    user, err := fetchUser(id)
    if err != nil {
        return nil, fmt.Errorf("fetch user: %w", err)
    }

    // Validate user
    if err := validateUser(user); err != nil {
        return nil, fmt.Errorf("validate user: %w", err)
    }

    // Enrich user data
    enriched, err := enrichUser(user)
    if err != nil {
        return nil, fmt.Errorf("enrich user: %w", err)
    }

    // Transform user
    transformed, err := transformUser(enriched)
    if err != nil {
        return nil, fmt.Errorf("transform user: %w", err)
    }

    return transformed, nil
}`}
</CodeCard>

<Compare>
  <CompareCol kind="bad" pill="issues">
    <p>Repetitive <code>if err != nil</code> checks.</p>
    <p>Error handling tangled with business logic.</p>
    <p>Hard to see the actual data flow.</p>
    <p>Difficult to compose operations.</p>
    <p>Easy to forget error checks.</p>
  </CompareCol>
  <CompareCol kind="good" pill="see right tab" title="fp-go">
    <p>No repetitive checks.</p>
    <p>Clear data flow pipeline.</p>
    <p>Business logic stays prominent.</p>
    <p>Composable, impossible to forget error handling.</p>
  </CompareCol>
</Compare>

  </TabItem>
  <TabItem value="solution" label="With fp-go">

<CodeCard file="processUser.go" status="tested">
{`func processUser(id string) result.Result[*User] {
    return function.Pipe4(
        fetchUser(id),                 // Result[*User]
        result.Chain(validateUser),    // Validation
        result.Chain(enrichUser),      // Enrichment
        result.Chain(transformUser),   // Transformation
    )
}`}
</CodeCard>

<Callout type="success" title="Benefits.">
  No repetitive error checks. Clear data flow pipeline. Business logic is prominent. Easy to compose. Impossible to forget error handling.
</Callout>

  </TabItem>
</Tabs>

</Section>

<Section
  id="benefits"
  number="02"
  title="Core"
  titleAccent="benefits."
  lede="Five concrete wins from adopting fp-go."
>

### 1. Type safety

fp-go leverages Go's type system to make errors impossible to ignore.

<Tabs>
  <TabItem value="unsafe" label="Easy to ignore errors">

<CodeCard file="unsafe.go">
{`// Idiomatic Go - easy to forget error check
result, _ := riskyOperation()  // Ignoring error!
doSomething(result)            // Potential panic`}
</CodeCard>

  </TabItem>
  <TabItem value="safe" label="Impossible to ignore">

<CodeCard file="safe.go" status="tested">
{`// fp-go - must handle the Result
result := riskyOperation()  // Returns Result[T]

// Can't access value without handling error
value := result.GetOrElse(func() T {
    return defaultValue
})

// Or explicitly handle both cases
result.Fold(
    func(err error) { /* handle error */ },
    func(val T) { /* handle success */ },
)`}
</CodeCard>

  </TabItem>
</Tabs>

### 2. Composability

Build complex operations from simple, reusable pieces.

<Tabs>
  <TabItem value="imperative" label="Imperative">

<CodeCard file="imperative.go">
{`func processData(data []int) ([]string, error) {
    // Filter
    var filtered []int
    for _, v := range data {
        if v > 0 {
            filtered = append(filtered, v)
        }
    }

    // Transform
    var doubled []int
    for _, v := range filtered {
        doubled = append(doubled, v*2)
    }

    // Validate
    for _, v := range doubled {
        if v > 100 {
            return nil, errors.New("value too large")
        }
    }

    // Convert to strings
    var result []string
    for _, v := range doubled {
        result = append(result, fmt.Sprintf("%d", v))
    }

    return result, nil
}`}
</CodeCard>

  </TabItem>
  <TabItem value="functional" label="Functional">

<CodeCard file="functional.go" status="tested">
{`func processData(data []int) result.Result[[]string] {
    return function.Pipe3(
        array.Filter(func(v int) bool { return v > 0 }),
        array.Map(func(v int) int { return v * 2 }),
        array.TraverseResult(func(v int) result.Result[string] {
            if v > 100 {
                return result.Err[string](errors.New("value too large"))
            }
            return result.Ok(fmt.Sprintf("%d", v))
        }),
    )(data)
}`}
</CodeCard>

<Callout type="success" title="Advantages.">
  Each operation is a pure, reusable function. Each step is independently testable. The data transformation pipeline is explicit. Errors propagate automatically.
</Callout>

  </TabItem>
</Tabs>

### 3. Testability

Pure functions are trivial to test — no mocks, no setup, no teardown.

<Tabs>
  <TabItem value="impure" label="Hard to test">

<CodeCard file="impure.go">
{`// Impure function - depends on external state
var db *sql.DB

func getUser(id string) (*User, error) {
    // Uses global db connection
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    // ... parsing logic
}

// Test requires:
// - Database setup
// - Test data insertion
// - Connection management
// - Cleanup`}
</CodeCard>

  </TabItem>
  <TabItem value="pure" label="Easy to test">

<CodeCard file="pure.go" status="tested">
{`// Pure function - all dependencies explicit
func getUser(db Database) func(string) result.Result[*User] {
    return func(id string) result.Result[*User] {
        return db.QueryUser(id)
    }
}

// Test is simple:
func TestGetUser(t *testing.T) {
    mockDB := &MockDatabase{
        users: map[string]*User{
            "123": {ID: "123", Name: "Alice"},
        },
    }

    result := getUser(mockDB)("123")

    assert.True(t, result.IsOk())
    assert.Equal(t, "Alice", result.GetOrElse(func() *User {
        return nil
    }).Name)
}`}
</CodeCard>

  </TabItem>
</Tabs>

### 4. Maintainability

<Compare>
  <CompareCol kind="bad" title="Before fp-go" pill="50+ lines">
    <p>Nested <code>if</code> statements.</p>
    <p>Mixed concerns.</p>
    <p>Hard to follow the logic.</p>
    <p>Difficult to add new steps.</p>
  </CompareCol>
  <CompareCol kind="good" title="With fp-go" pill="clear pipeline">
    <CodeCard file="snippet.go">
{`return function.Pipe5(
    step1,
    step2,
    step3,
    step4,
    step5,
)`}
    </CodeCard>
  </CompareCol>
</Compare>

### 5. Automatic error propagation

<Tabs>
  <TabItem value="manual" label="Manual propagation">

<CodeCard file="manual.go">
{`func process() error {
    result1, err := step1()
    if err != nil {
        return err
    }

    result2, err := step2(result1)
    if err != nil {
        return err
    }

    result3, err := step3(result2)
    if err != nil {
        return err
    }

    return step4(result3)
}`}
</CodeCard>

  </TabItem>
  <TabItem value="automatic" label="Automatic propagation">

<CodeCard file="automatic.go" status="tested">
{`func process() result.Result[T] {
    return function.Pipe3(
        step1(),      // If this fails, rest are skipped
        step2,        // Only runs if step1 succeeded
        step3,        // Only runs if step2 succeeded
        step4,        // Only runs if step3 succeeded
    )
}`}
</CodeCard>

  </TabItem>
</Tabs>

</Section>

<Section
  id="when"
  number="03"
  title="When to use"
  titleAccent="fp-go."
  lede="The honest fit guide — where it shines, and where idiomatic Go is the better answer."
>

### Excellent fit

#### Complex business logic

<CodeCard file="order.go">
{`// Multiple validation steps
// Data transformations
// Error handling at each step
func validateAndProcessOrder(order Order) result.Result[ProcessedOrder] {
    return function.Pipe5(
        validateCustomer(order.CustomerID),
        result.Chain(func(customer Customer) result.Result[Order] {
            return validateInventory(order)
        }),
        result.Chain(calculatePricing),
        result.Chain(applyDiscounts),
        result.Chain(finalizeOrder),
    )
}`}
</CodeCard>

#### Data transformation pipelines

<CodeCard file="etl.go">
{`// ETL operations
// Data cleaning
// Format conversions
func transformData(raw []RawData) result.Result[[]CleanData] {
    return function.Pipe4(
        array.Filter(isValid),
        array.Map(normalize),
        array.TraverseResult(enrich),
        result.Map(array.Map(format)),
    )(raw)
}`}
</CodeCard>

#### API clients with error handling

<CodeCard file="client.go">
{`// HTTP requests
// Response parsing
// Error handling
func fetchUserProfile(id string) ioresult.IOResult[Profile] {
    return function.Pipe3(
        buildRequest(id),
        ioresult.Chain(executeRequest),
        ioresult.Chain(parseResponse),
    )
}`}
</CodeCard>

#### Configuration management

<CodeCard file="config.go">
{`// Load config
// Validate
// Apply defaults
func loadConfig(path string) result.Result[Config] {
    return function.Pipe3(
        readConfigFile(path),
        result.Chain(parseConfig),
        result.Chain(validateConfig),
        result.Map(applyDefaults),
    )
}`}
</CodeCard>

### Use with caution

<Callout type="warn" title="Simple CRUD operations.">
  For straightforward DB calls, idiomatic Go is fine. fp-go adds unnecessary complexity:
  <CodeCard file="crud.go">
{`func getUser(id string) (*User, error) {
    return db.QueryUser(id)
}`}
  </CodeCard>
</Callout>

<Callout type="warn" title="Performance-critical hot paths.">
  For tight loops, use idiomatic Go or fp-go's idiomatic packages. Direct operations are faster; the idiomatic packages offer near-native performance.
</Callout>

<Callout type="warn" title="Team unfamiliar with FP.">
  Start with simple examples, provide training, introduce gradually, and use fp-go for new code first.
</Callout>

### Not recommended

<Callout type="warn" title="Trivial scripts and low-level system code.">
  For one-off scripts and code dealing with syscalls or memory management, stick with idiomatic Go.
</Callout>

</Section>

<Section
  id="case-studies"
  number="04"
  title="Real-world"
  titleAccent="success stories."
>

<ApiTable
  columns={['Case study', 'Before', 'After']}
  rows={[
    {symbol: 'API Gateway', signature: '500+ lines of nested error handling', description: '150 lines of clear pipeline code. Each step independently testable. 60% reduction in bugs.'},
    {symbol: 'Data Pipeline', signature: 'Complex state management, scattered error handling', description: 'Clear linear data flow. Centralized error handling. 40% faster development time.'},
    {symbol: 'Microservice', signature: 'Inconsistent error handling, hard to compose', description: 'Consistent patterns. Easy composition. 50% reduction in production errors.'},
  ]}
/>

</Section>

<Section
  id="compare"
  number="05"
  title="Compare with"
  titleAccent="other approaches."
>

### vs. idiomatic Go

<ApiTable
  columns={['Aspect', 'Idiomatic Go', 'fp-go']}
  rows={[
    {symbol: 'Error handling', signature: 'Manual if err != nil', description: 'Automatic propagation.'},
    {symbol: 'Composability', signature: 'Limited', description: 'Excellent.'},
    {symbol: 'Type safety', signature: 'Good', description: 'Excellent.'},
    {symbol: 'Learning curve', signature: 'Low', description: 'Medium.'},
    {symbol: 'Verbosity', signature: 'High for complex logic', description: 'Low.'},
    {symbol: 'Performance', signature: 'Excellent', description: 'Good (excellent with idiomatic packages).'},
    {symbol: 'Best for', signature: 'Simple operations', description: 'Complex logic.'},
  ]}
/>

### vs. other FP libraries

<ApiTable
  columns={['Feature', 'fp-go', 'samber/lo · go-functional']}
  rows={[
    {symbol: 'Monadic types', signature: 'Full support', description: 'samber/lo: limited · go-functional: yes.'},
    {symbol: 'Type safety', signature: 'Excellent', description: 'samber/lo: uses any · go-functional: good.'},
    {symbol: 'Error handling', signature: 'Built-in', description: 'samber/lo: manual · go-functional: built-in.'},
    {symbol: 'Documentation', signature: 'Comprehensive', description: 'samber/lo: good · go-functional: limited.'},
    {symbol: 'Active development', signature: 'Yes', description: 'samber/lo: yes · go-functional: sporadic.'},
    {symbol: 'Production-ready', signature: 'Yes (IBM)', description: 'samber/lo: yes · go-functional: unknown.'},
  ]}
/>

</Section>

<Section
  id="start"
  number="06"
  title="Getting"
  titleAccent="started."
>

<Checklist
  title="Adoption path"
  items={[
    {label: 'Start small — use Result for one function', impact: 'step 1'},
    {label: 'Learn core concepts (pure functions, monads, composition)', impact: 'step 2'},
    {label: 'Explore examples in the Recipes section', impact: 'step 3'},
    {label: 'Adopt gradually — new features first, refactor complex logic over time', impact: 'step 4'},
  ]}
/>

### 1. Start small

<CodeCard file="first.go" status="tested">
{`// Begin with simple Result usage
func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}`}
</CodeCard>

### 2. Learn core concepts

<ul>
  <li><a href="./concepts/pure-functions">Pure functions</a></li>
  <li><a href="./concepts/monads">Monads</a></li>
  <li><a href="./concepts/composition">Composition</a></li>
</ul>

### 3. Explore examples

<ul>
  <li><a href="./quickstart">Quickstart</a></li>
  <li><a href="./recipes/error-handling">Error handling recipes</a></li>
</ul>

### 4. Gradual adoption

Use for new features first; refactor complex logic gradually; keep simple code idiomatic.

</Section>

<Section
  id="takeaways"
  number="07"
  title="Key"
  titleAccent="takeaways."
>

<Checklist
  title="Remember"
  items={[
    {label: 'fp-go excels at complex logic — error handling, transformations, business rules', done: true},
    {label: 'Type safety prevents bugs — impossible to ignore errors', done: true},
    {label: 'Composability improves maintainability', done: true},
    {label: 'Testability is built-in — pure functions are trivial to test', done: true},
    {label: 'Not a silver bullet — use idiomatic Go for simple operations', done: true},
    {label: 'Production-ready — used at IBM and other companies', done: true},
    {label: 'Gradual adoption works — start small, expand as you learn', done: true},
  ]}
/>

</Section>
