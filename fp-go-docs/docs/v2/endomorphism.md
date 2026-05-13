---
title: Endomorphism
hide_title: true
description: Functions from a type to itself - special properties for composition and transformation pipelines.
sidebar_position: 28
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Endomorphism"
  lede="Functions from a type to itself (A → A). Endomorphisms have special properties that make them ideal for composition, transformation pipelines, and middleware patterns."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/endomorphism' },
    { label: 'Type', value: 'func(A) A' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

An endomorphism transforms a value while preserving its type:
- **Input and output**: Same type `A`
- **Composable**: Naturally chains together
- **Monoid**: Forms a monoid under composition

<CodeCard file="type_definition.go">
{`package endomorphism

// Endomorphism is a function from A to A
type Endomorphism[A any] = func(A) A
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Transformations | Modify values of the same type |
| Pipelines | Chain transformations |
| Middleware | Request/response processing |
| State updates | Functional state modifications |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Identity` | `func Identity[A any]() Endomorphism[A]` | Identity function (returns input) |
| `Monoid` | `func Monoid[A any]() monoid.Monoid[Endomorphism[A]]` | Monoid instance for composition |
</ApiTable>

### Composition

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Compose` | `func Compose[A any](f, g Endomorphism[A]) Endomorphism[A]` | Compose two endomorphisms (g ∘ f) |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "strings"
    E "github.com/IBM/fp-go/v2/endomorphism"
)

func main() {
    // Simple endomorphisms
    increment := func(n int) int {
        return n + 1
    }
    
    double := func(n int) int {
        return n * 2
    }
    
    // Compose: first double, then increment
    transform := E.Compose(increment, double)
    
    result := transform(5)  // (5 * 2) + 1 = 11
    fmt.Println(result)
    
    // String transformation
    toUpper := func(s string) string {
        return strings.ToUpper(s)
    }
    
    trim := func(s string) string {
        return strings.TrimSpace(s)
    }
    
    clean := E.Compose(toUpper, trim)
    fmt.Println(clean("  hello  "))  // "HELLO"
}
`}
</CodeCard>

### Struct Updates

<CodeCard file="struct_updates.go">
{`package main

import (
    "strings"
    E "github.com/IBM/fp-go/v2/endomorphism"
)

type User struct {
    Name  string
    Email string
    Age   int
}

// Endomorphisms for user updates
func updateName(name string) E.Endomorphism[User] {
    return func(u User) User {
        u.Name = name
        return u
    }
}

func normalizeEmail(u User) User {
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    return u
}

func incrementAge(u User) User {
    u.Age++
    return u
}

func main() {
    user := User{Name: "Alice", Email: " BOB@EXAMPLE.COM ", Age: 30}
    
    // Compose updates
    update := E.Compose(
        incrementAge,
        E.Compose(normalizeEmail, updateName("Bob")),
    )
    
    updated := update(user)
    fmt.Printf("%+v\n", updated)
    // {Name:Bob Email:bob@example.com Age:31}
}
`}
</CodeCard>

### Monoid Composition

<CodeCard file="monoid.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/endomorphism"
    M "github.com/IBM/fp-go/v2/monoid"
)

func main() {
    // Get monoid instance
    m := E.Monoid[int]()
    
    // Identity (no-op)
    identity := m.Empty()
    fmt.Println(identity(42))  // 42
    
    // Combine endomorphisms
    double := func(n int) int { return n * 2 }
    addTen := func(n int) int { return n + 10 }
    
    combined := m.Concat(double, addTen)
    result := combined(5)  // (5 * 2) + 10 = 20
    fmt.Println(result)
}
`}
</CodeCard>

### Middleware Pattern

<CodeCard file="middleware.go">
{`package main

import (
    "log"
    "net/http"
    E "github.com/IBM/fp-go/v2/endomorphism"
    F "github.com/IBM/fp-go/v2/function"
)

type Handler = E.Endomorphism[http.Handler]

func LoggingMiddleware() Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Printf("%s %s", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }
}

func AuthMiddleware() Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !isAuthenticated(r) {
                http.Error(w, "Unauthorized", 401)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func CORSMiddleware() Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            next.ServeHTTP(w, r)
        })
    }
}

func main() {
    baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Compose middleware
    handler := F.Pipe3(
        baseHandler,
        LoggingMiddleware(),
        AuthMiddleware(),
        CORSMiddleware(),
    )
    
    http.ListenAndServe(":8080", handler)
}
`}
</CodeCard>

### State Updates

<CodeCard file="state.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/endomorphism"
    F "github.com/IBM/fp-go/v2/function"
)

type AppState struct {
    Counter int
    Message string
    Active  bool
}

// State update endomorphisms
func incrementCounter(s AppState) AppState {
    s.Counter++
    return s
}

func setMessage(msg string) E.Endomorphism[AppState] {
    return func(s AppState) AppState {
        s.Message = msg
        return s
    }
}

func toggleActive(s AppState) AppState {
    s.Active = !s.Active
    return s
}

func main() {
    initialState := AppState{Counter: 0, Message: "", Active: false}
    
    // Compose state updates
    update := F.Pipe3(
        initialState,
        incrementCounter,
        setMessage("Updated"),
        toggleActive,
    )
    
    fmt.Printf("%+v\n", update)
    // {Counter:1 Message:Updated Active:true}
}
`}
</CodeCard>

### Data Transformation Pipeline

<CodeCard file="pipeline.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/endomorphism"
    F "github.com/IBM/fp-go/v2/function"
)

type Data struct {
    Value  int
    Status string
    Tags   []string
}

func main() {
    data := Data{Value: 10, Status: "pending", Tags: []string{}}
    
    // Build transformation pipeline
    pipeline := F.Pipe4(
        data,
        func(d Data) Data {
            d.Value *= 2
            return d
        },
        func(d Data) Data {
            d.Status = "processed"
            return d
        },
        func(d Data) Data {
            d.Tags = append(d.Tags, "validated")
            return d
        },
        func(d Data) Data {
            d.Value += 10
            return d
        },
    )
    
    fmt.Printf("%+v\n", pipeline)
    // {Value:30 Status:processed Tags:[validated]}
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Configuration Builder

<CodeCard file="config_builder.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/endomorphism"
)

type Config struct {
    Host    string
    Port    int
    Timeout int
    Debug   bool
}

func WithHost(host string) E.Endomorphism[Config] {
    return func(c Config) Config {
        c.Host = host
        return c
    }
}

func WithPort(port int) E.Endomorphism[Config] {
    return func(c Config) Config {
        c.Port = port
        return c
    }
}

func WithDebug(debug bool) E.Endomorphism[Config] {
    return func(c Config) Config {
        c.Debug = debug
        return c
    }
}

func main() {
    config := F.Pipe3(
        Config{},
        WithHost("localhost"),
        WithPort(8080),
        WithDebug(true),
    )
    
    fmt.Printf("%+v\n", config)
}
`}
</CodeCard>

### Pattern 2: Validation Chain

<CodeCard file="validation.go">
{`package main

import (
    "strings"
    E "github.com/IBM/fp-go/v2/endomorphism"
)

type FormData struct {
    Email    string
    Password string
    Valid    bool
    Errors   []string
}

func validateEmail(f FormData) FormData {
    if !strings.Contains(f.Email, "@") {
        f.Errors = append(f.Errors, "invalid email")
        f.Valid = false
    }
    return f
}

func validatePassword(f FormData) FormData {
    if len(f.Password) < 8 {
        f.Errors = append(f.Errors, "password too short")
        f.Valid = false
    }
    return f
}

func normalizeEmail(f FormData) FormData {
    f.Email = strings.ToLower(strings.TrimSpace(f.Email))
    return f
}

func main() {
    form := FormData{
        Email:    " USER@EXAMPLE.COM ",
        Password: "pass",
        Valid:    true,
    }
    
    validated := F.Pipe3(
        form,
        normalizeEmail,
        validateEmail,
        validatePassword,
    )
    
    fmt.Printf("%+v\n", validated)
}
`}
</CodeCard>

</Section>

<Callout type="info">

**Mathematical Property**: Endomorphisms form a monoid under composition with the identity function as the empty element. This makes them perfect for building composable transformation pipelines.

</Callout>
