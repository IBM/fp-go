---
title: IOOption
hide_title: true
description: Lazy side effects with optional values - combine IO and Option for effectful optional computations.
sidebar_position: 15
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="IOOption"
  lede="Combine lazy evaluation (IO) with optional values (Option). IOOption[A] represents a synchronous computation with side effects that may produce no value."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/iooption' },
    { label: 'Type', value: 'Monad (IO[Option[A]])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

<CodeCard file="type_definition.go">
{`package iooption

// IOOption is IO of Option
type IOOption[A any] = IO[Option[A]]
// Which expands to: func() Option[A]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use IOOption When | Use IOResult When |
|-------------------|-------------------|
| Optional results (cache miss, search) | Operations that can fail with errors |
| No error information needed | Need error messages |
| Absence is not an error | Failure needs explanation |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Some` | `func Some[A any](value A) IOOption[A]` | Create IOOption with value |
| `None` | `func None[A any]() IOOption[A]` | Create empty IOOption |
| `Of` | `func Of[A any](value A) IOOption[A]` | Alias for Some |
| `FromIO` | `func FromIO[A any](io IO[A]) IOOption[A]` | Lift IO to IOOption |
| `FromOption` | `func FromOption[A any](opt Option[A]) IOOption[A]` | Lift Option to IOOption |
| `FromNillable` | `func FromNillable[A any](io IO[*A]) IOOption[A]` | From IO of pointer |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(IOOption[A]) IOOption[B]` | Transform value if present |
| `Chain` | `func Chain[A, B any](f func(A) IOOption[B]) func(IOOption[A]) IOOption[B]` | Sequence operations |
| `Filter` | `func Filter[A any](pred func(A) bool) func(IOOption[A]) IOOption[A]` | Keep only if predicate holds |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Alt` | `func Alt[A any](second IOOption[A]) func(IOOption[A]) IOOption[A]` | Fallback on None |
| `GetOrElse` | `func GetOrElse[A any](f func() IO[A]) func(IOOption[A]) IO[A]` | Extract with default |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
)

func main() {
    // Create IOOption values
    some := IOO.Some(42)
    none := IOO.None[int]()
    
    // Execute
    result := some()  // Option[int] = Some(42)
    result = none()   // Option[int] = None
}
`}
</CodeCard>

### Cache Lookup

<CodeCard file="cache.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
)

func getFromCache(key string) IOO.IOOption[Data] {
    return func() option.Option[Data] {
        if val, ok := cache.Get(key); ok {
            return option.Some(val.(Data))
        }
        return option.None[Data]()
    }
}

func getWithFallback(key string) IO.IO[Data] {
    return IOO.GetOrElse(func() IO.IO[Data] {
        return fetchFromDB(key)
    })(getFromCache(key))
}

func main() {
    data := getWithFallback("user:123")()
}
`}
</CodeCard>

### Search Operations

<CodeCard file="search.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
    F "github.com/IBM/fp-go/v2/function"
)

func findUser(email string) IOO.IOOption[User] {
    return func() option.Option[User] {
        users := db.QueryUsers("email = ?", email)
        if len(users) > 0 {
            return option.Some(users[0])
        }
        return option.None[User]()
    }
}

func findActiveUser(email string) IOO.IOOption[User] {
    return F.Pipe1(
        findUser(email),
        IOO.Filter(func(u User) bool {
            return u.Active
        }),
    )
}

func main() {
    result := findActiveUser("alice@example.com")()
    // Option[User]
}
`}
</CodeCard>

### Chaining Optional Operations

<CodeCard file="chaining.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
    F "github.com/IBM/fp-go/v2/function"
)

func getUserProfile(id string) IOO.IOOption[Profile] {
    return F.Pipe2(
        findUser(id),
        IOO.Chain(func(u User) IOO.IOOption[Profile] {
            return loadProfile(u.ProfileID)
        }),
        IOO.Filter(func(p Profile) bool {
            return p.Visible
        }),
    )
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Fallback Chain

<CodeCard file="fallback.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
    F "github.com/IBM/fp-go/v2/function"
)

func getData(key string) IOO.IOOption[Data] {
    return F.Pipe2(
        getFromMemCache(key),
        IOO.Alt(getFromRedis(key)),
        IOO.Alt(getFromDatabase(key)),
    )
}
`}
</CodeCard>

### Pattern 2: Optional Configuration

<CodeCard file="config.go">
{`package main

import (
    IOO "github.com/IBM/fp-go/v2/iooption"
)

func loadOptionalConfig(path string) IOO.IOOption[Config] {
    return func() option.Option[Config] {
        data, err := os.ReadFile(path)
        if err != nil {
            return option.None[Config]()
        }
        
        var cfg Config
        if err := json.Unmarshal(data, &cfg); err != nil {
            return option.None[Config]()
        }
        
        return option.Some(cfg)
    }
}

func getConfig() IO.IO[Config] {
    return IOO.GetOrElse(func() IO.IO[Config] {
        return IO.Of(defaultConfig())
    })(loadOptionalConfig("config.json"))
}
`}
</CodeCard>

</Section>
