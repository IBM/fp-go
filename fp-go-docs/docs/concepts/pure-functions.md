---
sidebar_position: 11
title: Pure Functions
hide_title: true
description: What makes a function pure, why it matters, and how to apply purity in Go with fp-go.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Concepts · 01 / 06"
  title="Pure"
  titleAccent="functions."
  lede="The foundation of functional programming — what makes a function pure, why it matters, and how to apply it in Go."
  meta={[
    {label: '// Difficulty', value: 'Beginner'},
    {label: '// Reading time', value: '8 min · 6 sections'},
    {label: '// Prereqs', value: 'Basic Go'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Rule 1" prose value={<>Same input → <em>same output.</em></>} />
  <TLDRCard label="// Rule 2" prose value={<><em>No</em> side effects.</>} />
  <TLDRCard label="// Benefit" prose value={<>Testable, composable, <em>parallelizable.</em></>} variant="up" />
</TLDR>

<Section id="what" number="01" title="What is a" titleAccent="pure function?">

A **pure function** is a function that:

1. **Always returns the same output for the same input** (deterministic)
2. **Has no side effects** (doesn't modify external state)

That's it. Simple concept, powerful implications.

</Section>

<Section id="rules" number="02" title="The two" titleAccent="rules.">

### Rule 1 — deterministic

<Tabs groupId="purity">
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`// Pure: always returns same result for same input
func add(a, b int) int {
    return a + b
}

// Always returns 5
result1 := add(2, 3) // 5
result2 := add(2, 3) // 5
result3 := add(2, 3) // 5`}
</CodeCard>

</TabItem>
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`// Impure: result depends on external state
var counter int

func addWithCounter(a, b int) int {
    counter++
    return a + b + counter
}

// Returns different results!
result1 := addWithCounter(2, 3) // 6  (counter=1)
result2 := addWithCounter(2, 3) // 7  (counter=2)
result3 := addWithCounter(2, 3) // 8  (counter=3)`}
</CodeCard>

</TabItem>
</Tabs>

### Rule 2 — no side effects

<Tabs groupId="purity">
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`// Pure: doesn't modify anything outside itself
func multiply(a, b int) int {
    return a * b
}

// No external state changed
result := multiply(3, 4) // 12`}
</CodeCard>

</TabItem>
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`// Impure: modifies external state
var total int

func addToTotal(value int) int {
    total += value  // Side effect!
    return total
}

result := addToTotal(5) // total is now 5`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="common-side-effects" number="03" title="Common" titleAccent="side effects.">

### Modifying variables

<CodeCard file="impure-cache.go">
{`var cache map[string]string

// Impure: modifies global cache
func getCached(key string) string {
    if val, ok := cache[key]; ok {
        return val
    }
    val := fetchFromDB(key)
    cache[key] = val  // Side effect!
    return val
}`}
</CodeCard>

### I/O operations

<CodeCard file="impure-io.go">
{`// Impure: reads from file system
func readConfig() Config {
    data, _ := os.ReadFile("config.json")  // Side effect!
    var config Config
    json.Unmarshal(data, &config)
    return config
}

// Impure: writes to console
func logMessage(msg string) {
    fmt.Println(msg)  // Side effect!
}`}
</CodeCard>

### Network calls

<CodeCard file="impure-net.go">
{`// Impure: makes HTTP request
func fetchUser(id string) User {
    resp, _ := http.Get("https://api.example.com/users/" + id)  // Side effect!
    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    return user
}`}
</CodeCard>

### Random and time

<CodeCard file="impure-nondet.go">
{`// Impure: uses random number generator
func generateID() string {
    return fmt.Sprintf("id-%d", rand.Int())  // Side effect!
}

// Impure: depends on current time
func isExpired(expiresAt time.Time) bool {
    return time.Now().After(expiresAt)  // Side effect!
}`}
</CodeCard>

</Section>

<Section id="benefits" number="04" title="Why pure functions" titleAccent="matter.">

### 1. Easy to test

<Tabs groupId="testing">
<TabItem value="pure" label="✅ Pure (easy)">

<CodeCard file="pure-test.go" status="tested">
{`// Pure function
func calculateDiscount(price float64, percentage float64) float64 {
    return price * (percentage / 100)
}

// Simple test - no setup needed
func TestCalculateDiscount(t *testing.T) {
    result := calculateDiscount(100, 10)
    assert.Equal(t, 10.0, result)
}`}
</CodeCard>

</TabItem>
<TabItem value="impure" label="❌ Impure (hard)">

<CodeCard file="impure-test.go">
{`// Impure function
var discountRate float64

func calculateDiscount(price float64) float64 {
    return price * (discountRate / 100)
}

// Complex test - needs setup
func TestCalculateDiscount(t *testing.T) {
    // Setup
    oldRate := discountRate
    discountRate = 10
    defer func() { discountRate = oldRate }()

    // Test
    result := calculateDiscount(100)
    assert.Equal(t, 10.0, result)
}`}
</CodeCard>

</TabItem>
</Tabs>

### 2. Easy to reason about

<CodeCard file="reasoning.go">
{`// Pure: you know exactly what it does
func fullName(first, last string) string {
    return first + " " + last
}

// No need to check:
// - What global variables it uses
// - What files it reads
// - What network calls it makes
// - What it logs
// Just look at the function!`}
</CodeCard>

### 3. Easy to compose

<CodeCard file="compose.go" status="tested">
{`// Pure functions compose naturally
func double(x int) int { return x * 2 }
func addOne(x int) int { return x + 1 }
func square(x int) int { return x * x }

// Compose them
result := square(addOne(double(5)))  // ((5*2)+1)^2 = 121

// Or with fp-go
import "github.com/IBM/fp-go/v2/function"

composed := function.Flow3(double, addOne, square)
result := composed(5)  // 121`}
</CodeCard>

### 4. Cacheable (memoization)

<CodeCard file="memo.go">
{`// Pure functions can be safely cached
var cache = make(map[int]int)

func expensiveCalculation(n int) int {
    if result, ok := cache[n]; ok {
        return result  // Return cached result
    }

    // Do expensive calculation
    result := /* ... */
    cache[n] = result
    return result
}

// Safe because function is pure!
// Same input always gives same output`}
</CodeCard>

### 5. Parallelizable

<CodeCard file="parallel.go">
{`// Pure functions are safe to run in parallel
func processItem(item Item) Result {
    // Pure processing
    return transform(item)
}

// Safe to parallelize
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(i Item) {
        defer wg.Done()
        result := processItem(i)  // No race conditions!
        results <- result
    }(item)
}
wg.Wait()`}
</CodeCard>

</Section>

<Section id="patterns" number="05" title="Making functions" titleAccent="pure.">

### Pattern 1 — pass dependencies as parameters

<Tabs groupId="refactor">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`var db *sql.DB

func getUser(id string) (User, error) {
    // Uses global db
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    var user User
    err := row.Scan(&user.ID, &user.Name)
    return user, err
}`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`// Pass db as parameter
func getUser(db *sql.DB, id string) (User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    var user User
    err := row.Scan(&user.ID, &user.Name)
    return user, err
}

// Now testable with mock db!`}
</CodeCard>

</TabItem>
</Tabs>

### Pattern 2 — return new values, don't modify

<Tabs groupId="refactor">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`// Modifies the slice
func addItem(items []Item, item Item) {
    items = append(items, item)  // Modifies input!
}

original := []Item{{ID: 1}}
addItem(original, Item{ID: 2})
// original is now modified`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`// Returns new slice
func addItem(items []Item, item Item) []Item {
    result := make([]Item, len(items)+1)
    copy(result, items)
    result[len(items)] = item
    return result
}

original := []Item{{ID: 1}}
updated := addItem(original, Item{ID: 2})
// original unchanged, updated has new item`}
</CodeCard>

</TabItem>
</Tabs>

### Pattern 3 — separate pure logic from effects

<Tabs groupId="refactor">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`// Mixed pure logic and effects
func processOrder(orderID string) error {
    // Effect: fetch from DB
    order, err := db.GetOrder(orderID)
    if err != nil {
        return err
    }

    // Pure: calculate total
    total := 0.0
    for _, item := range order.Items {
        total += item.Price
    }

    // Effect: save to DB
    order.Total = total
    return db.SaveOrder(order)
}`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`// Pure: just calculates
func calculateTotal(items []Item) float64 {
    total := 0.0
    for _, item := range items {
        total += item.Price
    }
    return total
}

// Impure: handles effects
func processOrder(orderID string) error {
    order, err := db.GetOrder(orderID)
    if err != nil {
        return err
    }

    // Use pure function
    order.Total = calculateTotal(order.Items)

    return db.SaveOrder(order)
}

// Now calculateTotal is easily testable!`}
</CodeCard>

</TabItem>
</Tabs>

### Pattern 4 — use fp-go for effects

<Tabs groupId="refactor">
<TabItem value="standard" label="Without fp-go">

<CodeCard file="standard.go">
{`// Impure: executes immediately
func fetchUser(id string) (User, error) {
    resp, err := http.Get("https://api.example.com/users/" + id)
    if err != nil {
        return User{}, err
    }
    defer resp.Body.Close()

    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    return user, err
}`}
</CodeCard>

</TabItem>
<TabItem value="v2" label="With fp-go v2">

<CodeCard file="fp-go.go" status="tested">
{`// Pure: returns a description of the effect
func fetchUser(id string) ioresult.IOResult[User] {
    return func() result.Result[User] {
        resp, err := http.Get("https://api.example.com/users/" + id)
        if err != nil {
            return result.Err[User](err)
        }
        defer resp.Body.Close()

        var user User
        err = json.NewDecoder(resp.Body).Decode(&user)
        return result.FromGoError(user, err)
    }
}

// Function is pure - it just returns a function
// Effect happens when you execute it:
io := fetchUser("123")  // Pure! No HTTP call yet
user := io()            // Now the HTTP call happens`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="practical" number="06" title="Pure functions" titleAccent="in Go." lede="Go is not a pure functional language — and that's fine. Pure functions are a tool, not a religion.">

<Compare>
  <CompareCol kind="good" title="Use pure functions for" pill="fits">
    <p>Business logic.</p>
    <p>Calculations.</p>
    <p>Transformations.</p>
    <p>Validations.</p>
    <p>Formatting.</p>
    <p>Parsing (when possible).</p>
  </CompareCol>
  <CompareCol kind="bad" title="Don't force purity for" pill="trade-offs">
    <p>I/O operations (use fp-go types instead).</p>
    <p>Logging (use structured logging).</p>
    <p>Metrics (use dedicated libraries).</p>
    <p>Performance-critical code (if purity hurts performance).</p>
  </CompareCol>
</Compare>

<CodeCard file="balance.go">
{`// Pure: business logic
func calculateShipping(weight float64, distance float64) float64 {
    baseRate := 5.0
    weightRate := weight * 0.5
    distanceRate := distance * 0.1
    return baseRate + weightRate + distanceRate
}

// Pure: validation
func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return errors.New("invalid email")
    }
    return nil
}

// Impure but necessary: I/O
func saveOrder(order Order) error {
    // Use fp-go to make it more manageable
    return ioresult.Of(func() result.Result[Order] {
        // Database operation
        return result.FromGoError(order, db.Save(order))
    })()
}`}
</CodeCard>

</Section>

<Section id="examples" number="07" title="Real-world" titleAccent="examples.">

### E-commerce pricing

<Tabs groupId="example">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`var taxRate float64
var discountRate float64

func calculatePrice(basePrice float64) float64 {
    price := basePrice
    price -= price * (discountRate / 100)
    price += price * (taxRate / 100)
    return price
}

// Hard to test - depends on global state`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`type PricingConfig struct {
    TaxRate      float64
    DiscountRate float64
}

func calculatePrice(basePrice float64, config PricingConfig) float64 {
    price := basePrice
    price -= price * (config.DiscountRate / 100)
    price += price * (config.TaxRate / 100)
    return price
}

// Easy to test with different configs
func TestCalculatePrice(t *testing.T) {
    config := PricingConfig{TaxRate: 10, DiscountRate: 20}
    result := calculatePrice(100, config)
    assert.Equal(t, 88.0, result) // (100 - 20) + 8 = 88
}`}
</CodeCard>

</TabItem>
</Tabs>

### Data transformation

<Tabs groupId="example">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`func processUsers(users []User) {
    for i := range users {
        users[i].Name = strings.ToUpper(users[i].Name)
        users[i].Email = strings.ToLower(users[i].Email)
        users[i].Active = true
    }
}

// Modifies input - surprising behavior`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`func normalizeUser(user User) User {
    return User{
        ID:     user.ID,
        Name:   strings.ToUpper(user.Name),
        Email:  strings.ToLower(user.Email),
        Active: true,
    }
}

func processUsers(users []User) []User {
    result := make([]User, len(users))
    for i, user := range users {
        result[i] = normalizeUser(user)
    }
    return result
}

// Or with fp-go
import "github.com/IBM/fp-go/v2/array"

func processUsers(users []User) []User {
    return array.Map(normalizeUser)(users)
}

// Original unchanged, clear behavior`}
</CodeCard>

</TabItem>
</Tabs>

### Configuration

<Tabs groupId="example">
<TabItem value="impure" label="❌ Impure">

<CodeCard file="impure.go">
{`var config Config

func init() {
    data, _ := os.ReadFile("config.json")
    json.Unmarshal(data, &config)
}

func getTimeout() time.Duration {
    return config.Timeout
}

// Global state, hard to test`}
</CodeCard>

</TabItem>
<TabItem value="pure" label="✅ Pure">

<CodeCard file="pure.go" status="tested">
{`type Config struct {
    Timeout time.Duration
}

// Pure: just parses
func parseConfig(data []byte) (Config, error) {
    var config Config
    err := json.Unmarshal(data, &config)
    return config, err
}

// Pure: just extracts
func getTimeout(config Config) time.Duration {
    return config.Timeout
}

// Impure: I/O isolated
func loadConfig() (Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return Config{}, err
    }
    return parseConfig(data)
}`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="testing" number="08" title="Testing pure" titleAccent="functions.">

### Simple tests

<CodeCard file="test_basic.go" status="tested">
{`func TestPureFunctions(t *testing.T) {
    // No setup needed!

    t.Run("add", func(t *testing.T) {
        assert.Equal(t, 5, add(2, 3))
        assert.Equal(t, 0, add(-1, 1))
    })

    t.Run("multiply", func(t *testing.T) {
        assert.Equal(t, 12, multiply(3, 4))
        assert.Equal(t, 0, multiply(0, 100))
    })
}`}
</CodeCard>

### Property-based testing

<CodeCard file="test_properties.go" status="tested">
{`func TestAddCommutative(t *testing.T) {
    // Pure functions have mathematical properties
    for i := 0; i < 100; i++ {
        a := rand.Int()
        b := rand.Int()

        // Commutative: a + b = b + a
        assert.Equal(t, add(a, b), add(b, a))
    }
}

func TestAddAssociative(t *testing.T) {
    for i := 0; i < 100; i++ {
        a := rand.Int()
        b := rand.Int()
        c := rand.Int()

        // Associative: (a + b) + c = a + (b + c)
        assert.Equal(t, add(add(a, b), c), add(a, add(b, c)))
    }
}`}
</CodeCard>

</Section>

<Section id="faq" number="09" title="Common" titleAccent="questions.">

<Callout title="Aren't all Go functions impure?">
  No. Many Go functions are pure: <code>strings.ToUpper</code>, <code>math.Max</code>, <code>strconv.Itoa</code>, most of <code>encoding/json</code> parsing.
</Callout>

<Callout title="Should I never use global variables?">
  Use them wisely. Global constants are fine. Global configuration loaded once at startup is often okay. Global <em>mutable</em> state is problematic.
</Callout>

<Callout title="What about logging?">
  Logging is a side effect — but often acceptable. Prefer structured logging, log at boundaries (not in pure functions), and use context for request-scoped logging.
</Callout>

<Callout type="success" title="Is this practical in Go?">
  Yes. Many successful Go projects use pure functions extensively. It's about balance, not dogma.
</Callout>

</Section>

<Section id="summary" number="10" title="Summary">

<Checklist
  title="Pure functions"
  items={[
    {label: 'Same input → same output', done: true},
    {label: 'No side effects', done: true},
    {label: 'Easy to test', done: true},
    {label: 'Easy to reason about', done: true},
    {label: 'Easy to compose', done: true},
    {label: 'Cacheable', done: true},
    {label: 'Parallelizable', done: true},
  ]}
/>

<Callout type="success" title="Key takeaway.">
  Pure functions are a tool for writing better code. Use them where they help, don't force them where they don't.
</Callout>

</Section>
