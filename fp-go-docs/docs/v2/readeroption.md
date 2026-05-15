---
title: ReaderOption
hide_title: true
description: Dependency injection with optional values - combines Reader context with Option for computations that may not return a result.
sidebar_position: 21
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="ReaderOption"
  lede="Dependency injection with optional values. ReaderOption[C, A] combines Reader context with Option for computations that require dependencies and may not return a result."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/readeroption' },
    { label: 'Type', value: 'Reader[C, Option[A]]' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

ReaderOption represents a computation that requires context and may not return a value:
- **Reader**: Dependency injection pattern
- **Option**: Optional result (Some/None)

<CodeCard file="type_definition.go">
{`package readeroption

// ReaderOption is Reader specialized for optional results
type ReaderOption[C, A any] = Reader[C, Option[A]]
// Which expands to: func(C) Option[A]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Optional results with dependencies | Cache lookups, database queries |
| Configuration | Optional settings from context |
| Not an error | Absence is a valid outcome |
| Composable lookups | Chain optional operations with context |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Some` | `func Some[C, A any](value A) ReaderOption[C, A]` | Create with value |
| `None` | `func None[C, A any]() ReaderOption[C, A]` | Create empty |
| `Of` | `func Of[C, A any](value A) ReaderOption[C, A]` | Alias for Some |
| `FromReader` | `func FromReader[C, A any](r Reader[C, A]) ReaderOption[C, A]` | Lift Reader to ReaderOption |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(ReaderOption[C, A]) ReaderOption[C, B]` | Transform value if present |
| `Chain` | `func Chain[C, A, B any](f func(A) ReaderOption[C, B]) func(ReaderOption[C, A]) ReaderOption[C, B]` | FlatMap - sequence operations |
| `Ap` | `func Ap[C, A, B any](fa ReaderOption[C, A]) func(ReaderOption[C, func(A) B]) ReaderOption[C, B]` | Apply wrapped function |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `GetOrElse` | `func GetOrElse[C, A any](f func() A) func(ReaderOption[C, A]) Reader[C, A]` | Extract with default |
| `Fold` | `func Fold[C, A, B any](onNone func() B, onSome func(A) B) func(ReaderOption[C, A]) Reader[C, B]` | Pattern match both cases |
</ApiTable>

### Testing

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `IsSome` | `func IsSome[C, A any](ReaderOption[C, A]) Reader[C, bool]` | Test for value |
| `IsNone` | `func IsNone[C, A any](ReaderOption[C, A]) Reader[C, bool]` | Test for absence |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "github.com/IBM/fp-go/v2/readeroption"
    "github.com/IBM/fp-go/v2/option"
)

type Dependencies struct {
    Cache *Cache
    DB    *sql.DB
}

func FindInCache(key string) readeroption.ReaderOption[Dependencies, Data] {
    return func(deps Dependencies) option.Option[Data] {
        if val, ok := deps.Cache.Get(key); ok {
            return option.Some(val)
        }
        return option.None[Data]()
    }
}

func main() {
    deps := Dependencies{Cache: NewCache()}
    
    // Execute with dependencies
    result := FindInCache("key")(deps)
    
    // Check result
    if option.IsSome(result) {
        fmt.Println("Found in cache")
    }
}
`}
</CodeCard>

### Optional Configuration

<CodeCard file="config.go">
{`package main

import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    O "github.com/IBM/fp-go/v2/option"
)

type Config struct {
    Settings map[string]string
}

func GetSetting(key string) RO.ReaderOption[Config, string] {
    return func(cfg Config) O.Option[string] {
        if val, ok := cfg.Settings[key]; ok {
            return O.Some(val)
        }
        return O.None[string]()
    }
}

func main() {
    cfg := Config{Settings: map[string]string{"theme": "dark"}}
    
    // With fallback
    theme := RO.GetOrElse(func() string {
        return "light"
    })(GetSetting("theme"))(cfg)
    
    fmt.Println(theme) // "dark"
}
`}
</CodeCard>

### Chaining Lookups

<CodeCard file="chaining.go">
{`package main

import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    F "github.com/IBM/fp-go/v2/function"
)

type Dependencies struct {
    Cache *Cache
    DB    *sql.DB
}

func FindInCache(id string) RO.ReaderOption[Dependencies, User] {
    return func(deps Dependencies) option.Option[User] {
        if user, ok := deps.Cache.Get(id); ok {
            return option.Some(user)
        }
        return option.None[User]()
    }
}

func FindInDB(id string) RO.ReaderOption[Dependencies, User] {
    return func(deps Dependencies) option.Option[User] {
        user, err := deps.DB.Query(id)
        if err != nil {
            return option.None[User]()
        }
        return option.Some(user)
    }
}

func FindUser(id string) RO.ReaderOption[Dependencies, User] {
    return F.Pipe2(
        FindInCache(id),
        RO.Alt(func() RO.ReaderOption[Dependencies, User] {
            return FindInDB(id)
        }),
    )
}
`}
</CodeCard>

### Validation with Context

<CodeCard file="validation.go">
{`package main

import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    O "github.com/IBM/fp-go/v2/option"
)

type ValidationContext struct {
    AllowedDomains []string
}

func ValidateEmail(email string) RO.ReaderOption[ValidationContext, string] {
    return func(ctx ValidationContext) O.Option[string] {
        domain := extractDomain(email)
        for _, allowed := range ctx.AllowedDomains {
            if domain == allowed {
                return O.Some(email)
            }
        }
        return O.None[string]()
    }
}

func main() {
    ctx := ValidationContext{
        AllowedDomains: []string{"example.com", "test.com"},
    }
    
    result := ValidateEmail("user@example.com")(ctx)
    // Some("user@example.com")
    
    result = ValidateEmail("user@invalid.com")(ctx)
    // None
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Fallback Chain

<CodeCard file="fallback.go">
{`package main

import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    F "github.com/IBM/fp-go/v2/function"
)

// Try multiple sources in order
func GetConfig(key string) RO.ReaderOption[Dependencies, string] {
    return F.Pipe3(
        GetFromEnv(key),
        RO.Alt(func() RO.ReaderOption[Dependencies, string] {
            return GetFromFile(key)
        }),
        RO.Alt(func() RO.ReaderOption[Dependencies, string] {
            return GetFromDefaults(key)
        }),
    )
}
`}
</CodeCard>

### Pattern 2: Conditional Execution

<CodeCard file="conditional.go">
{`package main

import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    O "github.com/IBM/fp-go/v2/option"
)

func GetFeature(name string) RO.ReaderOption[AppContext, Feature] {
    return func(ctx AppContext) O.Option[Feature] {
        if !ctx.FeatureFlags.IsEnabled(name) {
            return O.None[Feature]()
        }
        return O.Some(ctx.Features[name])
    }
}
`}
</CodeCard>

</Section>
