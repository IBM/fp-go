---
title: ReaderIOEither
hide_title: true
description: The ultimate monad stack - dependency injection + lazy effects + custom error handling.
sidebar_position: 20
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="ReaderIOEither"
  lede="The most powerful combination: dependency injection (Reader) + lazy effects (IO) + custom error handling (Either). ReaderIOEither[C, E, A] is the ultimate monad stack."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/readerioeither' },
    { label: 'Type', value: 'Monad (func(C) IO[Either[E, A]])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

<CodeCard file="type_definition.go">
{`package readerioeither

// ReaderIOEither combines Reader, IO, and Either
type ReaderIOEither[C, E, A any] = Reader[C, IOEither[E, A]]
// Which expands to: func(C) func() Either[E, A]
`}
</CodeCard>

### The Ultimate Stack

<ApiTable>
| Layer | Provides |
|-------|----------|
| **Reader[C, ...]** | Dependency injection |
| **IO[...]** | Lazy evaluation of side effects |
| **Either[E, A]** | Custom error handling |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Right` | `func Right[C, E, A any](value A) ReaderIOEither[C, E, A]` | Create successful value |
| `Left` | `func Left[C, E, A any](err E) ReaderIOEither[C, E, A]` | Create error value |
| `Of` | `func Of[C, E, A any](value A) ReaderIOEither[C, E, A]` | Alias for Right |
| `Ask` | `func Ask[C, E any]() ReaderIOEither[C, E, C]` | Access context |
| `Asks` | `func Asks[C, E, A any](f func(C) A) ReaderIOEither[C, E, A]` | Access and transform context |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, E, A, B any](f func(A) B) func(ReaderIOEither[C, E, A]) ReaderIOEither[C, E, B]` | Transform success value |
| `MapLeft` | `func MapLeft[C, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[C, E1, A]) ReaderIOEither[C, E2, A]` | Transform error |
| `Chain` | `func Chain[C, E, A, B any](f func(A) ReaderIOEither[C, E, B]) func(ReaderIOEither[C, E, A]) ReaderIOEither[C, E, B]` | Sequence operations |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[C, E, A, B any](fa ReaderIOEither[C, E, A]) func(ReaderIOEither[C, E, func(A) B]) ReaderIOEither[C, E, B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[C, E, A any]([]ReaderIOEither[C, E, A]) ReaderIOEither[C, E, []A]` | All-or-nothing |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Application

<CodeCard file="basic.go">
{`package main

import (
    RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

type Dependencies struct {
    DB     *sql.DB
    Logger *log.Logger
}

type AppError struct {
    Code    int
    Message string
}

func fetchUser(id string) RIOE.ReaderIOEither[Dependencies, AppError, User] {
    return RIOE.Ask[Dependencies, AppError, *sql.DB](func(deps Dependencies) *sql.DB {
        return deps.DB
    }).Chain(func(db *sql.DB) RIOE.ReaderIOEither[Dependencies, AppError, User] {
        return RIOE.TryCatch(func() (User, AppError) {
            user, err := queryUser(db, id)
            if err != nil {
                return User{}, AppError{Code: 500, Message: err.Error()}
            }
            return user, AppError{}
        })
    })
}

func main() {
    deps := Dependencies{
        DB:     connectDB(),
        Logger: log.New(os.Stdout, "", 0),
    }
    
    result := fetchUser("user-123")(deps)()
    // Either[AppError, User]
}
`}
</CodeCard>

### Complete Application

<CodeCard file="complete_app.go">
{`package main

import (
    RIOE "github.com/IBM/fp-go/v2/readerioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type Config struct {
    DBUrl string
    Port  int
}

type Dependencies struct {
    DB     *sql.DB
    Logger *log.Logger
    Config Config
}

type AppError struct {
    Code    int
    Message string
}

func getDB() RIOE.ReaderIOEither[Dependencies, AppError, *sql.DB] {
    return RIOE.Asks(func(deps Dependencies) *sql.DB {
        return deps.DB
    })
}

func logInfo(msg string) RIOE.ReaderIOEither[Dependencies, AppError, unit.Unit] {
    return RIOE.Ask[Dependencies, AppError, *log.Logger](func(deps Dependencies) *log.Logger {
        return deps.Logger
    }).Chain(func(logger *log.Logger) RIOE.ReaderIOEither[Dependencies, AppError, unit.Unit] {
        return RIOE.FromIO(IO.FromImpure(func() {
            logger.Println(msg)
        }))
    })
}

func fetchUser(id string) RIOE.ReaderIOEither[Dependencies, AppError, User] {
    return F.Pipe3(
        logInfo("Fetching user: " + id),
        RIOE.Chain(func(_ unit.Unit) RIOE.ReaderIOEither[Dependencies, AppError, *sql.DB] {
            return getDB()
        }),
        RIOE.Chain(func(db *sql.DB) RIOE.ReaderIOEither[Dependencies, AppError, User] {
            return RIOE.TryCatch(func() (User, AppError) {
                user, err := db.QueryUser(id)
                if err != nil {
                    return User{}, AppError{Code: 500, Message: err.Error()}
                }
                return user, AppError{}
            })
        }),
        RIOE.ChainFirst(func(u User) RIOE.ReaderIOEither[Dependencies, AppError, unit.Unit] {
            return logInfo("Found user: " + u.Name)
        }),
    )
}

func main() {
    deps := Dependencies{
        DB:     connectDB(),
        Logger: log.New(os.Stdout, "", 0),
        Config: loadConfig(),
    }
    
    result := fetchUser("user-123")(deps)()
    // Either[AppError, User]
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern: Full-Stack Application

<CodeCard file="fullstack.go">
{`package main

import (
    RIOE "github.com/IBM/fp-go/v2/readerioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func processOrder(orderID string) RIOE.ReaderIOEither[Dependencies, AppError, Receipt] {
    return F.Pipe4(
        fetchOrder(orderID),
        RIOE.Chain(validateOrder),
        RIOE.Chain(chargePayment),
        RIOE.Chain(sendConfirmation),
        RIOE.MapLeft(func(err AppError) AppError {
            return AppError{
                Code:    err.Code,
                Message: fmt.Sprintf("Order processing failed: %s", err.Message),
            }
        }),
    )
}
`}
</CodeCard>

### When to Use

<ApiTable>
| Use ReaderIOEither When | Consider Simpler Alternative |
|-------------------------|------------------------------|
| Need all three: DI + effects + custom errors | Use Reader if no effects needed |
| Building full applications | Use IOEither if no DI needed |
| Complex business logic | Use ReaderIO if errors are simple |
| Maximum composability required | Simpler types if features not needed |
</ApiTable>

</Section>
