---
sidebar_position: 13
title: Monads
hide_title: true
description: A practical pattern for sequencing computations with context — no category theory required.
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<PageHeader
  eyebrow="Concepts · 02 / 06"
  title="Monads."
  lede="A practical pattern for chaining operations with context. No category theory required."
  meta={[
    {label: '// Difficulty', value: 'Intermediate'},
    {label: '// Reading time', value: '12 min · 8 sections'},
    {label: '// Prereqs', value: 'Pure functions'}
  ]}
/>

<TLDR>
  <TLDRCard label="// Mental model" prose value={<>A box with rules for <em>chaining</em>.</>} />
  <TLDRCard label="// Common monads" prose value={<>Option · Result · IO · <em>Array</em>.</>} />
  <TLDRCard label="// Key ops" prose value={<><em>Map</em>, Chain, Fold.</>} variant="up" />
</TLDR>

<Section id="what" number="01" title="What is a" titleAccent="monad?">

A **monad** is a design pattern for chaining operations that have some "context" or "effect".

<Callout title="Think of it as.">
  A box that holds a value, plus rules for chaining operations on that value while preserving the context.
</Callout>

<ApiTable
  columns={['Monad', 'Context', 'When you reach for it']}
  rows={[
    {symbol: 'Option', signature: '"might not have a value"', description: 'Optional / nullable fields without nil pointers.'},
    {symbol: 'Result', signature: '"might be an error"', description: 'Error handling without manual if-err checks.'},
    {symbol: 'IO', signature: '"will perform side effects"', description: 'Effectful computations described lazily.'},
    {symbol: 'Array', signature: '"has multiple values"', description: 'Branching / one-to-many transformations.'},
  ]}
/>

</Section>

<Section id="why" number="02" title="Why do we" titleAccent="need monads?">

<Tabs groupId="problem">
<TabItem value="problem" label="The problem">

<CodeCard file="problem.go">
{`// Without monads: error handling is messy
func processUser(id string) (User, error) {
    user, err := fetchUser(id)
    if err != nil {
        return User{}, err
    }

    validated, err := validateUser(user)
    if err != nil {
        return User{}, err
    }

    enriched, err := enrichUser(validated)
    if err != nil {
        return User{}, err
    }

    saved, err := saveUser(enriched)
    if err != nil {
        return User{}, err
    }

    return saved, nil
}

// Repetitive error checking
// Hard to compose
// Lots of boilerplate`}
</CodeCard>

</TabItem>
<TabItem value="solution" label="With monads">

<CodeCard file="solution.go" status="tested">
{`// With Result monad: clean and composable
import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/function"
)

func processUser(id string) result.Result[User] {
    return function.Pipe4(
        fetchUser(id),
        result.Chain(validateUser),
        result.Chain(enrichUser),
        result.Chain(saveUser),
    )
}

// No repetitive error checking
// Composable
// Clear data flow
// Stops at first error automatically`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="pattern" number="03" title="The monad" titleAccent="pattern." lede="Every monad has three things.">

### 1. A type constructor

<CodeCard file="constructors.go">
{`// Option monad
type Option[A any] struct { /* ... */ }

// Result monad
type Result[A any] struct { /* ... */ }

// IO monad
type IO[A any] func() A`}
</CodeCard>

### 2. A way to put values in (`Return`/`Of`)

<CodeCard file="of.go">
{`// Option
opt := option.Some(42)           // Put 42 in Option context

// Result
res := result.Ok(42)             // Put 42 in Result context

// IO
io := io.Of(func() int { return 42 })  // Put 42 in IO context`}
</CodeCard>

### 3. A way to chain operations (`Chain`/`FlatMap`)

<CodeCard file="chain.go">
{`// Option
result := option.Chain(func(x int) option.Option[int] {
    return option.Some(x * 2)
})(opt)

// Result
result := result.Chain(func(x int) result.Result[int] {
    return result.Ok(x * 2)
})(res)

// IO
result := io.Chain(func(x int) io.IO[int] {
    return io.Of(func() int { return x * 2 })
})(myIO)`}
</CodeCard>

</Section>

<Section id="common" number="04" title="Common monads in" titleAccent="fp-go.">

### Option monad — "might not have a value"

<Tabs groupId="monad">
<TabItem value="without" label="Without Option">

<CodeCard file="without.go">
{`func findUser(id string) *User {
    user := db.FindByID(id)
    return user
}

func getEmail(user *User) *string {
    if user == nil {
        return nil
    }
    return &user.Email
}

// Lots of nil checks
user := findUser("123")
if user != nil {
    email := getEmail(user)
    if email != nil {
        // Use email
    }
}`}
</CodeCard>

</TabItem>
<TabItem value="with" label="With Option">

<CodeCard file="with.go" status="tested">
{`import "github.com/IBM/fp-go/v2/option"

func findUser(id string) option.Option[User] {
    user := db.FindByID(id)
    if user == nil {
        return option.None[User]()
    }
    return option.Some(*user)
}

func getEmail(user User) option.Option[string] {
    if user.Email == "" {
        return option.None[string]()
    }
    return option.Some(user.Email)
}

// Chain operations
email := option.Chain(getEmail)(findUser("123"))

// Or with Pipe
email := function.Pipe2(
    findUser("123"),
    option.Chain(getEmail),
)

// Handle result
email.Fold(
    func() { fmt.Println("No email") },
    func(e string) { fmt.Println("Email:", e) },
)`}
</CodeCard>

</TabItem>
</Tabs>

### Result monad — "might be an error"

<Tabs groupId="monad">
<TabItem value="without" label="Without Result">

<CodeCard file="without.go">
{`func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func processNumbers(a, b, c int) (int, error) {
    result1, err := divide(a, b)
    if err != nil {
        return 0, err
    }

    result2, err := divide(result1, c)
    if err != nil {
        return 0, err
    }

    return result2, nil
}`}
</CodeCard>

</TabItem>
<TabItem value="with" label="With Result">

<CodeCard file="with.go" status="tested">
{`import "github.com/IBM/fp-go/v2/result"

func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

func processNumbers(a, b, c int) result.Result[int] {
    return function.Pipe2(
        divide(a, b),
        result.Chain(func(x int) result.Result[int] {
            return divide(x, c)
        }),
    )
}

// Stops at first error automatically`}
</CodeCard>

</TabItem>
</Tabs>

### IO monad — "will perform side effects"

<Tabs groupId="monad">
<TabItem value="without" label="Without IO">

<CodeCard file="without.go">
{`func readFile(path string) ([]byte, error) {
    return os.ReadFile(path) // Executes immediately
}

func parseJSON(data []byte) (Config, error) {
    var config Config
    err := json.Unmarshal(data, &config)
    return config, err
}

// Hard to test, executes immediately
config, err := readFile("config.json")
if err != nil {
    return err
}
parsed, err := parseJSON(config)`}
</CodeCard>

</TabItem>
<TabItem value="with" label="With IO">

<CodeCard file="with.go" status="tested">
{`import "github.com/IBM/fp-go/v2/ioresult"

func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

func parseJSON(data []byte) result.Result[Config] {
    var config Config
    err := json.Unmarshal(data, &config)
    return result.FromGoError(config, err)
}

// Build pipeline (doesn't execute yet)
loadConfig := function.Pipe2(
    readFile("config.json"),
    ioresult.Chain(func(data []byte) ioresult.IOResult[Config] {
        return ioresult.FromResult(parseJSON(data))
    }),
)

// Execute when ready
config := loadConfig()  // Now it runs`}
</CodeCard>

</TabItem>
</Tabs>

### Array monad — "has multiple values"

<CodeCard file="array.go" status="tested">
{`import "github.com/IBM/fp-go/v2/array"

// Array is a monad!
numbers := []int{1, 2, 3}

// Chain (FlatMap) - each element produces an array
result := array.Chain(func(x int) []int {
    return []int{x, x * 2}
})(numbers)
// [1, 2, 2, 4, 3, 6]

// Map - each element produces a single value
doubled := array.Map(func(x int) int {
    return x * 2
})(numbers)
// [2, 4, 6]`}
</CodeCard>

</Section>

<Section id="operations" number="05" title="Monad" titleAccent="operations.">

### Map vs. Chain

<Compare>
  <CompareCol kind="good" title="Map" pill="A → B">
    <p>Transform the value inside.</p>
    <CodeCard file="map.go">
{`opt := option.Some(5)
doubled := option.Map(func(x int) int {
    return x * 2
})(opt)
// Some(10)`}
    </CodeCard>
  </CompareCol>
  <CompareCol kind="good" title="Chain" pill="A → M[B]">
    <p>Transform and return a new monad.</p>
    <CodeCard file="chain.go">
{`opt := option.Some(5)
result := option.Chain(func(x int) option.Option[int] {
    if x > 0 {
        return option.Some(x * 2)
    }
    return option.None[int]()
})(opt)
// Some(10)`}
    </CodeCard>
  </CompareCol>
</Compare>

<Callout type="success" title="Rule of thumb.">
  Use <strong>Map</strong> when your function returns a plain value. Use <strong>Chain</strong> when your function returns another monad.
</Callout>

### Common operations across monads

<CodeCard file="ops-of.go">
{`// Of/Return - Put value in monad
option.Some(42)
result.Ok(42)
io.Of(func() int { return 42 })`}
</CodeCard>

<CodeCard file="ops-map.go">
{`// Map - Transform value
option.Map(func(x int) int { return x * 2 })(opt)
result.Map(func(x int) int { return x * 2 })(res)
io.Map(func(x int) int { return x * 2 })(myIO)`}
</CodeCard>

<CodeCard file="ops-chain.go">
{`// Chain/FlatMap - Transform and flatten
option.Chain(func(x int) option.Option[int] {
    return option.Some(x * 2)
})(opt)

result.Chain(func(x int) result.Result[int] {
    return result.Ok(x * 2)
})(res)`}
</CodeCard>

<CodeCard file="ops-fold.go">
{`// Fold - Extract value
option.Fold(
    func() int { return 0 },           // None case
    func(x int) int { return x },      // Some case
)(opt)

result.Fold(
    func(err error) int { return 0 },  // Error case
    func(x int) int { return x },      // Success case
)(res)`}
</CodeCard>

</Section>

<Section id="laws" number="06" title="The monad" titleAccent="laws." lede="Three laws guarantee predictable composition. You don't have to memorize them.">

### Law 1 — Left identity

`of(a).chain(f) === f(a)`

<CodeCard file="law1.go">
{`a := 5
f := func(x int) option.Option[int] { return option.Some(x * 2) }

// These are equivalent:
result1 := option.Chain(f)(option.Some(a))
result2 := f(a)
// Both give Some(10)`}
</CodeCard>

### Law 2 — Right identity

`m.chain(of) === m`

<CodeCard file="law2.go">
{`m := option.Some(5)

// These are equivalent:
result1 := option.Chain(option.Some[int])(m)
result2 := m
// Both give Some(5)`}
</CodeCard>

### Law 3 — Associativity

`m.chain(f).chain(g) === m.chain(x => f(x).chain(g))`

<CodeCard file="law3.go">
{`m := option.Some(5)
f := func(x int) option.Option[int] { return option.Some(x * 2) }
g := func(x int) option.Option[int] { return option.Some(x + 1) }

// These are equivalent:
result1 := option.Chain(g)(option.Chain(f)(m))
result2 := option.Chain(func(x int) option.Option[int] {
    return option.Chain(g)(f(x))
})(m)
// Both give Some(11)`}
</CodeCard>

<Callout type="success" title="Why these matter.">
  They ensure monads compose predictably.
</Callout>

</Section>

<Section id="examples" number="07" title="Real-world" titleAccent="examples.">

### User registration

<Tabs groupId="example">
<TabItem value="standard" label="Without monads">

<CodeCard file="without.go">
{`func registerUser(email, password string) (User, error) {
    if !isValidEmail(email) {
        return User{}, errors.New("invalid email")
    }

    existing, err := db.FindByEmail(email)
    if err != nil {
        return User{}, err
    }
    if existing != nil {
        return User{}, errors.New("email already exists")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return User{}, err
    }

    user := User{
        Email:    email,
        Password: string(hash),
    }

    if err := db.Save(&user); err != nil {
        return User{}, err
    }

    return user, nil
}`}
</CodeCard>

</TabItem>
<TabItem value="monads" label="With monads">

<CodeCard file="with.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/function"
)

func validateEmail(email string) result.Result[string] {
    if !isValidEmail(email) {
        return result.Err[string](errors.New("invalid email"))
    }
    return result.Ok(email)
}

func checkNotExists(email string) result.Result[string] {
    existing, err := db.FindByEmail(email)
    if err != nil {
        return result.Err[string](err)
    }
    if existing != nil {
        return result.Err[string](errors.New("email already exists"))
    }
    return result.Ok(email)
}

func hashPassword(password string) result.Result[string] {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return result.FromGoError(string(hash), err)
}

func createUser(email, hash string) result.Result[User] {
    user := User{Email: email, Password: hash}
    err := db.Save(&user)
    return result.FromGoError(user, err)
}

func registerUser(email, password string) result.Result[User] {
    return function.Pipe4(
        validateEmail(email),
        result.Chain(checkNotExists),
        result.Chain(func(e string) result.Result[string] {
            return hashPassword(password)
        }),
        result.Chain(func(hash string) result.Result[User] {
            return createUser(email, hash)
        }),
    )
}`}
</CodeCard>

</TabItem>
</Tabs>

### Configuration loading

<Tabs groupId="example">
<TabItem value="standard" label="Without monads">

<CodeCard file="without.go">
{`func loadConfig() (Config, error) {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return Config{}, err
    }

    var raw RawConfig
    if err := json.Unmarshal(data, &raw); err != nil {
        return Config{}, err
    }

    if raw.Port == 0 {
        return Config{}, errors.New("port required")
    }

    config := Config{
        Port:    raw.Port,
        Host:    raw.Host,
        Timeout: time.Duration(raw.TimeoutSec) * time.Second,
    }

    return config, nil
}`}
</CodeCard>

</TabItem>
<TabItem value="monads" label="With monads">

<CodeCard file="with.go" status="tested">
{`import (
    "github.com/IBM/fp-go/v2/ioresult"
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/function"
)

func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

func parseJSON(data []byte) result.Result[RawConfig] {
    var raw RawConfig
    err := json.Unmarshal(data, &raw)
    return result.FromGoError(raw, err)
}

func validateConfig(raw RawConfig) result.Result[RawConfig] {
    if raw.Port == 0 {
        return result.Err[RawConfig](errors.New("port required"))
    }
    return result.Ok(raw)
}

func transformConfig(raw RawConfig) Config {
    return Config{
        Port:    raw.Port,
        Host:    raw.Host,
        Timeout: time.Duration(raw.TimeoutSec) * time.Second,
    }
}

func loadConfig() ioresult.IOResult[Config] {
    return function.Pipe4(
        readFile("config.json"),
        ioresult.Chain(func(data []byte) ioresult.IOResult[RawConfig] {
            return ioresult.FromResult(parseJSON(data))
        }),
        ioresult.Chain(func(raw RawConfig) ioresult.IOResult[RawConfig] {
            return ioresult.FromResult(validateConfig(raw))
        }),
        ioresult.Map(transformConfig),
    )
}`}
</CodeCard>

</TabItem>
</Tabs>

</Section>

<Section id="patterns" number="08" title="Common" titleAccent="patterns.">

### Sequential operations

<CodeCard file="sequential.go">
{`result := function.Pipe4(
    fetchUser(id),
    result.Chain(validateUser),
    result.Chain(enrichUser),
    result.Chain(saveUser),
)`}
</CodeCard>

### Conditional logic

<CodeCard file="conditional.go">
{`result := result.Chain(func(user User) result.Result[User] {
    if user.Age < 18 {
        return result.Err[User](errors.New("too young"))
    }
    return result.Ok(user)
})(userResult)`}
</CodeCard>

### Combining results

<CodeCard file="combine.go">
{`func createOrder(userID, productID string) result.Result[Order] {
    user := fetchUser(userID)
    product := fetchProduct(productID)

    return result.Chain(func(u User) result.Result[Order] {
        return result.Map(func(p Product) Order {
            return Order{User: u, Product: p}
        })(product)
    })(user)
}`}
</CodeCard>

### Error recovery

<CodeCard file="recover.go">
{`result := result.OrElse(func(err error) result.Result[User] {
    log.Printf("Error: %v, using default", err)
    return result.Ok(defaultUser)
})(userResult)`}
</CodeCard>

</Section>

<Section id="when" number="09" title="When to use" titleAccent="monads.">

<Compare>
  <CompareCol kind="good" title="Use monads when" pill="fits">
    <p>Chaining operations that can fail.</p>
    <p>Handling optional values.</p>
    <p>Managing side effects.</p>
    <p>Building composable pipelines.</p>
    <p>Need consistent error handling.</p>
  </CompareCol>
  <CompareCol kind="bad" title="Don't force when" pill="trade-offs">
    <p>Simple, one-off operations.</p>
    <p>Performance is critical (hot path).</p>
    <p>Team unfamiliar with the pattern.</p>
    <p>Standard Go is clearer.</p>
  </CompareCol>
</Compare>

</Section>

<Section id="faq" number="10" title="Common" titleAccent="questions.">

<Callout title="Do I need to understand category theory?">
  No. Think of monads as a design pattern for chaining operations with context. The theory is interesting but not required.
</Callout>

<Callout title="Aren't monads just error handling?">
  No. Error handling is one use case. Monads handle any "context": Option (optional values), Result (errors), IO (side effects), Array (multiple values), Reader (dependency injection).
</Callout>

<Callout title="Is this overengineering?">
  Depends. For simple cases, standard Go is fine. For complex error handling and composition, monads shine.
</Callout>

<Callout type="success" title="Map vs. Chain — quick rule.">
  Map: function returns a plain value. Chain: function returns a monad.
  <CodeCard file="rule.go">
{`// Map: int → string
result.Map(func(x int) string { return fmt.Sprint(x) })

// Chain: int → Result[string]
result.Chain(func(x int) result.Result[string] {
    return result.Ok(fmt.Sprint(x))
})`}
  </CodeCard>
</Callout>

</Section>

<Section id="summary" number="11" title="Summary">

<Checklist
  title="Monads"
  items={[
    {label: 'Pattern for chaining with context', done: true},
    {label: 'Three parts: type, return, chain', done: true},
    {label: 'Common types: Option, Result, IO, Array', done: true},
    {label: 'Operations: Map, Chain, Fold', done: true},
    {label: 'Laws ensure predictable behavior', done: true},
  ]}
/>

<Callout type="success" title="Key takeaway.">
  Monads are a practical pattern for managing context in a composable way. You don't need to understand the theory to use them effectively.
</Callout>

</Section>
