---
title: Effect
hide_title: true
description: Complete solution for real-world applications - combines dependency injection, lazy evaluation, and error handling.
sidebar_position: 16
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Effect"
  lede="Complete solution for real-world applications. Effect[C, A] combines dependency injection (Reader), lazy evaluation (IO), and error handling (Result) into one powerful abstraction."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/effect' },
    { label: 'Type', value: 'ReaderIOResult[C, A]' },
    { label: 'Alias', value: 'func(C) func() Either[error, A]' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Effect is the most powerful type in fp-go, combining three essential abstractions:
- **Reader[C]**: Dependency injection and context passing
- **IO**: Lazy evaluation and side effect management  
- **Result[A]**: Type-safe error handling

<CodeCard file="type_definition.go">
{`package effect

// Effect is an alias for ReaderIOResult
type Effect[C, A any] = ReaderIOResult[C, A]
// Which expands to: func(C) IO[Result[A]]
// Or fully: func(C) func() Either[error, A]
`}
</CodeCard>

### Why Effect?

<ApiTable>
| Capability | Benefit |
|------------|---------|
| Dependency Injection | Type-safe, testable dependencies |
| Lazy Evaluation | Control when effects execute |
| Error Handling | Type-safe, composable errors |
| Testability | Easy to mock dependencies |
| Composability | Build complex operations from simple ones |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Succeed` | `func Succeed[C, A any](value A) Effect[C, A]` | Create success value |
| `Fail` | `func Fail[C, A any](err error) Effect[C, A]` | Create failure |
| `Of` | `func Of[C, A any](value A) Effect[C, A]` | Alias for Succeed |
| `FromIOResult` | `func FromIOResult[C, A any](ior IOResult[A]) Effect[C, A]` | Lift IOResult |
| `FromIO` | `func FromIO[C, A any](io IO[A]) Effect[C, A]` | Lift IO (always succeeds) |
| `Ask` | `func Ask[C, A any](f func(C) A) Effect[C, A]` | Access dependencies |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(Effect[C, A]) Effect[C, B]` | Transform success value |
| `Chain` | `func Chain[C, A, B any](f func(A) Effect[C, B]) func(Effect[C, A]) Effect[C, B]` | Sequence operations |
| `Tap` | `func Tap[C, A any](f func(A) Effect[C, any]) func(Effect[C, A]) Effect[C, A]` | Side effect, keep value |
| `BiMap` | `func BiMap[C, A, B any](fe func(error) error, fa func(A) B) func(Effect[C, A]) Effect[C, B]` | Transform both sides |
</ApiTable>

### Dependencies

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Asks` | `func Asks[C, A any](f func(C) A) Effect[C, A]` | Access and transform context |
| `Provide` | `func Provide[C, A any](ctx C) func(Effect[C, A]) IOResult[A]` | Supply dependencies |
| `Local` | `func Local[C1, C2, A any](f func(C1) C2) func(Effect[C2, A]) Effect[C1, A]` | Transform context |
</ApiTable>

### Error Handling

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `OrElse` | `func OrElse[C, A any](f func(error) Effect[C, A]) func(Effect[C, A]) Effect[C, A]` | Fallback on error |
| `MapError` | `func MapError[C, A any](f func(error) error) func(Effect[C, A]) Effect[C, A]` | Transform error |
| `Fold` | `func Fold[C, A, B any](onError func(error) B, onOk func(A) B) func(Effect[C, A]) Reader[C, IO[B]]` | Pattern match |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[C, A, B any](fa Effect[C, A]) func(Effect[C, func(A) B]) Effect[C, B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[C, A any]([]Effect[C, A]) Effect[C, []A]` | All-or-nothing for arrays |
| `TraverseArray` | `func TraverseArray[C, A, B any](f func(A) Effect[C, B]) func([]A) Effect[C, []B]` | Map and sequence |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "errors"
    E "github.com/IBM/fp-go/v2/effect"
)

type Dependencies struct {
    DB     *sql.DB
    Logger *log.Logger
}

func main() {
    // Create effects
    success := E.Succeed[Dependencies, string]("Hello")
    failure := E.Fail[Dependencies, string](errors.New("error"))
    
    // Access dependencies
    getDB := E.Ask[Dependencies, *sql.DB](func(deps Dependencies) *sql.DB {
        return deps.DB
    })
    
    // Execute with dependencies
    deps := Dependencies{DB: connectDB(), Logger: log.Default()}
    result := success(deps)()  // Result[string] = Ok("Hello")
}
`}
</CodeCard>

### Service Layer

<CodeCard file="service.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

type Dependencies struct {
    DB     *sql.DB
    Cache  *Cache
    Logger *log.Logger
}

type UserService struct{}

func (s *UserService) GetUser(id string) E.Effect[Dependencies, User] {
    return F.Pipe3(
        // Try cache first
        s.getUserFromCache(id),
        E.OrElse(func(err error) E.Effect[Dependencies, User] {
            // Fallback to database
            return s.getUserFromDB(id)
        }),
        E.Tap(s.logUser),
        E.Chain(s.validateUser),
    )
}

func (s *UserService) getUserFromCache(id string) E.Effect[Dependencies, User] {
    return E.Ask[Dependencies, User](func(deps Dependencies) User {
        if user, ok := deps.Cache.Get(id); ok {
            return user
        }
        return User{} // Will trigger OrElse
    })
}

func (s *UserService) getUserFromDB(id string) E.Effect[Dependencies, User] {
    return E.Ask[Dependencies, User](func(deps Dependencies) User {
        var user User
        err := deps.DB.QueryRow(
            "SELECT id, name, email FROM users WHERE id = ?", id,
        ).Scan(&user.ID, &user.Name, &user.Email)
        
        if err != nil {
            return User{}
        }
        return user
    })
}

func (s *UserService) logUser(user User) E.Effect[Dependencies, any] {
    return E.Ask[Dependencies, any](func(deps Dependencies) any {
        deps.Logger.Printf("Processing user: %s", user.ID)
        return nil
    })
}

func (s *UserService) validateUser(user User) E.Effect[Dependencies, User] {
    if user.Name == "" {
        return E.Fail[Dependencies, User](errors.New("invalid user"))
    }
    return E.Succeed[Dependencies, User](user)
}
`}
</CodeCard>

### HTTP Handler

<CodeCard file="http.go">
{`package main

import (
    "encoding/json"
    "net/http"
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

type HTTPDeps struct {
    UserService *UserService
    Logger      *log.Logger
}

func GetUserHandler(deps HTTPDeps) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        
        // Create effect
        eff := F.Pipe2(
            deps.UserService.GetUser(id),
            E.Map[HTTPDeps](func(u User) []byte {
                data, _ := json.Marshal(u)
                return data
            }),
        )
        
        // Execute effect
        result := eff(deps)()
        
        // Handle result
        if result.IsError(result) {
            http.Error(w, "Error", http.StatusInternalServerError)
            return
        }
        
        data, _ := result.Unwrap(result)
        w.Header().Set("Content-Type", "application/json")
        w.Write(data)
    }
}
`}
</CodeCard>

### Transaction Management

<CodeCard file="transaction.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

type TxDeps struct {
    DB *sql.DB
}

func WithTransaction[A any](
    eff E.Effect[TxDeps, A],
) E.Effect[TxDeps, A] {
    return E.Ask[TxDeps, A](func(deps TxDeps) A {
        tx, err := deps.DB.Begin()
        if err != nil {
            return *new(A)
        }
        
        // Execute with transaction
        txDeps := TxDeps{DB: tx}
        result := eff(txDeps)()
        
        // Commit or rollback
        if result.IsError(result) {
            tx.Rollback()
            return *new(A)
        }
        
        if err := tx.Commit(); err != nil {
            return *new(A)
        }
        
        val, _ := result.Unwrap(result)
        return val
    })
}

func TransferFunds(from, to string, amount float64) E.Effect[TxDeps, unit.Unit] {
    return WithTransaction(
        F.Pipe2(
            debit(from, amount),
            E.Chain[TxDeps](func(_ unit.Unit) E.Effect[TxDeps, unit.Unit] {
                return credit(to, amount)
            }),
        ),
    )
}
`}
</CodeCard>

### Do Notation

<CodeCard file="do_notation.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

type State struct {
    User   User
    Orders []Order
    Total  float64
}

func GenerateReport(userID string) E.Effect[Dependencies, Report] {
    return F.Pipe4(
        E.Do[Dependencies](State{}),
        E.Bind("User", func(s State) E.Effect[Dependencies, User] {
            return fetchUser(userID)
        }),
        E.Bind("Orders", func(s State) E.Effect[Dependencies, []Order] {
            return fetchOrders(s.User.ID)
        }),
        E.Let("Total", func(s State) float64 {
            var total float64
            for _, order := range s.Orders {
                total += order.Amount
            }
            return total
        }),
        E.Map[Dependencies](func(s State) Report {
            return Report{
                User:   s.User,
                Orders: s.Orders,
                Total:  s.Total,
            }
        }),
    )
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Middleware

<CodeCard file="middleware.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/effect"
)

func LoggingMiddleware[A any](
    next E.Effect[HTTPDeps, A],
) E.Effect[HTTPDeps, A] {
    return E.Tap[HTTPDeps](func(resp A) E.Effect[HTTPDeps, any] {
        return E.Ask[HTTPDeps, any](func(deps HTTPDeps) any {
            deps.Logger.Printf("Response: %+v", resp)
            return nil
        })
    })(next)
}

func AuthMiddleware[A any](
    next E.Effect[HTTPDeps, A],
) E.Effect[HTTPDeps, A] {
    return E.Chain[HTTPDeps](func(_ any) E.Effect[HTTPDeps, A] {
        return next
    })(validateAuth())
}
`}
</CodeCard>

### Pattern 2: Fallback Chain

<CodeCard file="fallback.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

func GetConfig(key string) E.Effect[Dependencies, string] {
    return F.Pipe3(
        GetFromEnv(key),
        E.OrElse(func(err error) E.Effect[Dependencies, string] {
            return GetFromFile(key)
        }),
        E.OrElse(func(err error) E.Effect[Dependencies, string] {
            return GetFromDefaults(key)
        }),
    )
}
`}
</CodeCard>

### Pattern 3: Testing

<CodeCard file="testing.go">
{`package main

import (
    "testing"
    E "github.com/IBM/fp-go/v2/effect"
)

func TestProcessOrder(t *testing.T) {
    // Mock dependencies
    mockDeps := Dependencies{
        DB: &MockDB{
            Orders: map[string]Order{
                "123": {ID: "123", Total: 100.0},
            },
        },
        Logger: log.New(io.Discard, "", 0),
        Config: Config{},
    }
    
    // Execute effect with mocks
    result := ProcessOrder("123")(mockDeps)()
    
    // Assert
    assert.True(t, result.IsOk(result))
    order, _ := result.Unwrap(result)
    assert.Equal(t, "123", order.ID)
}
`}
</CodeCard>

</Section>

<Section id="comparison" number="05" title="Type" titleAccent="Comparison">

<ApiTable>
| Feature | Effect[C, A] | IOResult[A] | ReaderIOResult[C, A] |
|---------|--------------|-------------|----------------------|
| Dependencies | ✅ Yes | ❌ No | ✅ Yes |
| Lazy | ✅ Yes | ✅ Yes | ✅ Yes |
| Can fail | ✅ Yes | ✅ Yes | ✅ Yes |
| Type params | 2 | 1 | 2 |
| Use case | Full applications | Simple IO | Same as Effect |
</ApiTable>

<Callout type="info">

**Note**: Effect is an alias for ReaderIOResult - they are the same type with different naming conventions inspired by effect-ts.

</Callout>

</Section>

<Section id="best-practices" number="06" title="Best" titleAccent="Practices">

### ✅ Do

<CodeCard file="do.go">
{`// Use Effect for application logic
func ProcessUser(id string) effect.Effect[Dependencies, User]

// Compose effects
result := F.Pipe3(
    fetchUser(id),
    effect.Chain[Dependencies](validateUser),
    effect.Chain[Dependencies](saveUser),
)

// Provide dependencies at the edge
func main() {
    deps := setupDependencies()
    result := ProcessUser("123")(deps)()
}

// Test with mock dependencies
func TestProcessUser(t *testing.T) {
    mockDeps := createMockDeps()
    result := ProcessUser("123")(mockDeps)()
}
`}
</CodeCard>

### ❌ Don't

<CodeCard file="dont.go">
{`// Don't execute effects in constructors
func NewService() *Service {
    data := fetchData()(deps)()  // ❌ Executes immediately
    return &Service{data: data}
}

// Don't pass dependencies as parameters
func ProcessUser(deps Dependencies, id string) effect.Effect[Dependencies, User]  // ❌
func ProcessUser(id string) effect.Effect[Dependencies, User]  // ✅

// Don't ignore errors
effect.Map[Dependencies](func(u User) User {
    saveUser(u)(deps)()  // ❌ Error is lost
    return u
})
`}
</CodeCard>

</Section>
