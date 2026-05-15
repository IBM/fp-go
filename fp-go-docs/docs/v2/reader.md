---
title: Reader
hide_title: true
description: Functional dependency injection - computations that depend on a shared context.
sidebar_position: 17
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Reader"
  lede="Functional dependency injection. Reader[C, A] represents a computation that depends on a context C and produces a value A."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/reader' },
    { label: 'Type', value: 'Monad (func(C) A)' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Reader is simply a function from context to value:

<CodeCard file="type_definition.go">
{`package reader

// Reader is a function from context to value
type Reader[C, A any] = func(C) A
`}
</CodeCard>

### Why Reader?

<Compare>
<CompareCol kind="bad">
<CodeCard file="traditional_di.go">
{`// ❌ Traditional DI with structs
type Service struct {
    db     *sql.DB
    logger *log.Logger
    config Config
}

func (s *Service) GetUser(id string) User {
    // Hard to compose
    // Difficult to test
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="reader_di.go">
{`// ✅ Reader-based DI
type Dependencies struct {
    DB     *sql.DB
    Logger *log.Logger
    Config Config
}

func GetUser(id string) reader.Reader[Dependencies, User] {
    return func(deps Dependencies) User {
        // Easy to compose
        // Simple to test
    }
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[C, A any](value A) Reader[C, A]` | Wrap pure value (ignores context) |
| `Ask` | `func Ask[C, A any](f func(C) A) Reader[C, A]` | Access and transform context |
| `Asks` | `func Asks[C, A any](f func(C) A) Reader[C, A]` | Alias for Ask |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(Reader[C, A]) Reader[C, B]` | Transform result |
| `Chain` | `func Chain[C, A, B any](f func(A) Reader[C, B]) func(Reader[C, A]) Reader[C, B]` | Sequence operations |
| `Flatten` | `func Flatten[C, A any](Reader[C, Reader[C, A]]) Reader[C, A]` | Unwrap nested Reader |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[C, A, B any](fa Reader[C, A]) func(Reader[C, func(A) B]) Reader[C, B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[C, A any]([]Reader[C, A]) Reader[C, []A]` | All-or-nothing |
| `TraverseArray` | `func TraverseArray[C, A, B any](f func(A) Reader[C, B]) func([]A) Reader[C, []B]` | Map and sequence |
</ApiTable>

### Context Manipulation

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Local` | `func Local[C1, C2, A any](f func(C1) C2) func(Reader[C2, A]) Reader[C1, A]` | Transform context before use |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    R "github.com/IBM/fp-go/v2/reader"
)

type Config struct {
    DBUrl string
    Port  int
}

func main() {
    // Wrap pure value
    answer := R.Of[Config, int](42)
    result := answer(Config{})  // 42
    
    // Access context
    getDBUrl := R.Ask(func(c Config) string {
        return c.DBUrl
    })
    
    config := Config{DBUrl: "localhost:5432", Port: 8080}
    url := getDBUrl(config)  // "localhost:5432"
    
    // Access and transform
    getPort := R.Asks(func(c Config) int {
        return c.Port
    })
    port := getPort(config)  // 8080
}
`}
</CodeCard>

### Dependency Injection

<CodeCard file="dependency_injection.go">
{`package main

import (
    "database/sql"
    "log"
    R "github.com/IBM/fp-go/v2/reader"
    F "github.com/IBM/fp-go/v2/function"
)

type Dependencies struct {
    DB     *sql.DB
    Logger *log.Logger
    Config Config
}

func getDB() R.Reader[Dependencies, *sql.DB] {
    return R.Ask(func(deps Dependencies) *sql.DB {
        return deps.DB
    })
}

func getLogger() R.Reader[Dependencies, *log.Logger] {
    return R.Ask(func(deps Dependencies) *log.Logger {
        return deps.Logger
    })
}

func fetchUser(id string) R.Reader[Dependencies, User] {
    return F.Pipe1(
        getDB(),
        R.Chain(func(db *sql.DB) R.Reader[Dependencies, User] {
            return R.Of[Dependencies](queryUser(db, id))
        }),
    )
}

func main() {
    deps := Dependencies{
        DB:     connectDB(),
        Logger: log.New(os.Stdout, "", 0),
        Config: loadConfig(),
    }
    
    user := fetchUser("user-123")(deps)
}
`}
</CodeCard>

### Composing Operations

<CodeCard file="composing.go">
{`package main

import (
    R "github.com/IBM/fp-go/v2/reader"
    F "github.com/IBM/fp-go/v2/function"
)

func getUserEmail(id string) R.Reader[Dependencies, string] {
    return F.Pipe1(
        fetchUser(id),
        R.Map(func(u User) string {
            return u.Email
        }),
    )
}

func sendEmail(to, subject, body string) R.Reader[Dependencies, bool] {
    return R.Ask(func(deps Dependencies) bool {
        // Use deps.Logger, deps.Config, etc.
        return sendEmailViaService(to, subject, body)
    })
}

func notifyUser(id string) R.Reader[Dependencies, bool] {
    return F.Pipe1(
        getUserEmail(id),
        R.Chain(func(email string) R.Reader[Dependencies, bool] {
            return sendEmail(email, "Notification", "Hello!")
        }),
    )
}

func main() {
    deps := Dependencies{...}
    success := notifyUser("user-123")(deps)
}
`}
</CodeCard>

### Local Context Transformation

<CodeCard file="local_context.go">
{`package main

import (
    R "github.com/IBM/fp-go/v2/reader"
)

type AppConfig struct {
    Database DatabaseConfig
    Server   ServerConfig
}

type DatabaseConfig struct {
    Host string
    Port int
}

// Reader that needs DatabaseConfig
func connectDB() R.Reader[DatabaseConfig, *sql.DB] {
    return R.Ask(func(cfg DatabaseConfig) *sql.DB {
        return sql.Open("postgres", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
    })
}

// Transform AppConfig to DatabaseConfig
func withDBConfig() R.Reader[AppConfig, *sql.DB] {
    return R.Local(func(app AppConfig) DatabaseConfig {
        return app.Database
    })(connectDB())
}

func main() {
    appConfig := AppConfig{
        Database: DatabaseConfig{Host: "localhost", Port: 5432},
        Server:   ServerConfig{Port: 8080},
    }
    
    db := withDBConfig()(appConfig)
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Testing with Mocks

<CodeCard file="testing.go">
{`package main

import (
    "testing"
    R "github.com/IBM/fp-go/v2/reader"
)

func TestFetchUser(t *testing.T) {
    // Mock dependencies
    mockDeps := Dependencies{
        DB:     mockDB(),
        Logger: mockLogger(),
        Config: testConfig(),
    }
    
    // Test the Reader
    user := fetchUser("test-id")(mockDeps)
    
    assert.Equal(t, "test-id", user.ID)
}
`}
</CodeCard>

### Pattern 2: Configuration Management

<CodeCard file="config_management.go">
{`package main

import (
    R "github.com/IBM/fp-go/v2/reader"
)

type Config struct {
    APIKey    string
    Timeout   time.Duration
    MaxRetries int
}

func getAPIKey() R.Reader[Config, string] {
    return R.Asks(func(c Config) string {
        return c.APIKey
    })
}

func getTimeout() R.Reader[Config, time.Duration] {
    return R.Asks(func(c Config) time.Duration {
        return c.Timeout
    })
}

func makeAPICall(endpoint string) R.Reader[Config, Response] {
    return F.Pipe2(
        getAPIKey(),
        R.Chain(func(key string) R.Reader[Config, time.Duration] {
            return getTimeout()
        }),
        R.Chain(func(timeout time.Duration) R.Reader[Config, Response] {
            return R.Of[Config](callAPI(endpoint, key, timeout))
        }),
    )
}
`}
</CodeCard>

### When to Use Reader

<ApiTable>
| Use Reader When | Consider Alternative |
|-----------------|---------------------|
| Need dependency injection | Simple function parameters sufficient |
| Multiple operations share context | Single operation doesn't need DI |
| Want testable, composable code | Struct-based DI is adequate |
| Context is read-only | Need mutable state (use State) |
</ApiTable>

</Section>
