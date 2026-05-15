---
title: Record - Traverse
hide_title: true
description: Working with maps of effectful computations using traverse operations.
sidebar_position: 14
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record Traverse"
  lede="Working with maps of effectful computations. Traverse operations allow you to map over maps with effects and sequence the results."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'Traverse, TraverseWithIndex, Sequence' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Traverse` | `func Traverse[K comparable, A, B, F any](Applicative[F]) func(func(A) HKT[F, B]) func(map[K]A) HKT[F, map[K]B]` | Map and sequence |
| `TraverseWithIndex` | `func TraverseWithIndex[K comparable, A, B, F any](Applicative[F]) func(func(K, A) HKT[F, B]) func(map[K]A) HKT[F, map[K]B]` | Traverse with keys |
| `Sequence` | `func Sequence[K comparable, A, F any](Applicative[F]) func(map[K]HKT[F, A]) HKT[F, map[K]A]` | Flip map and effect |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Traverse with Option

<CodeCard file="option.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

m := map[string]int{"a": 1, "b": 2, "c": 3}

// Validate all values
validated := R.Traverse(O.Applicative[int]())(
    func(v int) O.Option[int] {
        if v > 0 {
            return O.Some(v)
        }
        return O.None[int]()
    },
)(m)
// Some(map[string]int{"a": 1, "b": 2, "c": 3})

// With invalid value
m2 := map[string]int{"a": 1, "b": -1, "c": 3}
validated2 := R.Traverse(O.Applicative[int]())(
    func(v int) O.Option[int] {
        if v > 0 {
            return O.Some(v)
        }
        return O.None[int]()
    },
)(m2)
// None - one invalid value fails all
`}
</CodeCard>

### Traverse with Result

<CodeCard file="result.go">
{`import Res "github.com/IBM/fp-go/v2/result"

configs := map[string]string{
    "port":    "8080",
    "timeout": "30",
    "retries": "abc",  // Invalid
}

// Parse all values
parsed := R.Traverse(Res.Applicative[error, int]())(
    func(s string) Res.Result[int] {
        if n, err := strconv.Atoi(s); err == nil {
            return Res.Success(n)
        } else {
            return Res.Error[int](err)
        }
    },
)(configs)
// Error - "abc" is not a valid integer
`}
</CodeCard>

### TraverseWithIndex

<CodeCard file="with_index.go">
{`configs := map[string]string{
    "api_url": "https://api.example.com",
    "timeout": "30",
    "retries": "3",
}

// Parse with context from keys
parsed := R.TraverseWithIndex(O.Applicative[int]())(
    func(key string, value string) O.Option[int] {
        if key == "api_url" {
            // Skip non-numeric configs
            return O.Some(0)
        }
        if n, err := strconv.Atoi(value); err == nil {
            return O.Some(n)
        }
        return O.None[int]()
    },
)(configs)
// Option[map[string]int]
`}
</CodeCard>

### Sequence

<CodeCard file="sequence.go">
{`// Map of Options
options := map[string]O.Option[int]{
    "a": O.Some(1),
    "b": O.Some(2),
    "c": O.Some(3),
}

// Flip to Option of map
result := R.Sequence(O.Applicative[int]())(options)
// Some(map[string]int{"a": 1, "b": 2, "c": 3})

// With None
withNone := map[string]O.Option[int]{
    "a": O.Some(1),
    "b": O.None[int](),
    "c": O.Some(3),
}

result2 := R.Sequence(O.Applicative[int]())(withNone)
// None
`}
</CodeCard>

### Configuration Parsing

<CodeCard file="config_parse.go">
{`type Config struct {
    Port    int
    Timeout int
    Retries int
}

raw := map[string]string{
    "port":    "8080",
    "timeout": "30",
    "retries": "3",
}

// Parse all fields
parsed := R.Traverse(Res.Applicative[error, int]())(
    func(s string) Res.Result[int] {
        n, err := strconv.Atoi(s)
        if err != nil {
            return Res.Error[int](err)
        }
        return Res.Success(n)
    },
)(raw)

// Build config from result
config := F.Pipe2(
    parsed,
    Res.Map(func(m map[string]int) Config {
        return Config{
            Port:    m["port"],
            Timeout: m["timeout"],
            Retries: m["retries"],
        }
    }),
)
// Result[Config]
`}
</CodeCard>

### Validation Suite

<CodeCard file="validation.go">
{`type User struct {
    Email string
    Age   int
}

user := User{Email: "test@example.com", Age: 25}

validators := map[string]func(User) Res.Result[bool]{
    "email": func(u User) Res.Result[bool] {
        if strings.Contains(u.Email, "@") {
            return Res.Success(true)
        }
        return Res.Error[bool](errors.New("invalid email"))
    },
    "age": func(u User) Res.Result[bool] {
        if u.Age >= 18 {
            return Res.Success(true)
        }
        return Res.Error[bool](errors.New("must be 18+"))
    },
}

// Run all validations
results := R.Traverse(Res.Applicative[error, bool]())(
    func(validate func(User) Res.Result[bool]) Res.Result[bool] {
        return validate(user)
    },
)(validators)
// Result[map[string]bool] - Success if all pass
`}
</CodeCard>

### API Batch Requests

<CodeCard file="api.go">
{`import IOE "github.com/IBM/fp-go/v2/ioeither"

type UserData struct {
    ID   int
    Name string
}

userIDs := map[string]int{
    "alice": 1,
    "bob":   2,
    "charlie": 3,
}

// Fetch all users
fetchAll := R.Traverse(IOE.Applicative[error, UserData]())(
    func(id int) IOE.IOEither[error, UserData] {
        return fetchUserAPI(id)
    },
)(userIDs)
// IOEither[error, map[string]UserData]

// Execute
users := fetchAll()
// Either[error, map[string]UserData]
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### All or Nothing Validation

<CodeCard file="all_or_nothing.go">
{`// Validate all entries - fail if any invalid
func ValidateAll(data map[string]string) Res.Result[map[string]int] {
    return R.Traverse(Res.Applicative[error, int]())(
        func(s string) Res.Result[int] {
            n, err := strconv.Atoi(s)
            if err != nil {
                return Res.Error[int](
                    fmt.Errorf("invalid value: %s", s),
                )
            }
            return Res.Success(n)
        },
    )(data)
}
`}
</CodeCard>

### Collecting Errors

<CodeCard file="collect_errors.go">
{`// Validate with error collection
type ValidationErrors []error

func ValidateWithErrors(
    data map[string]string,
) (map[string]int, ValidationErrors) {
    result := make(map[string]int)
    var errors ValidationErrors
    
    for k, v := range data {
        if n, err := strconv.Atoi(v); err == nil {
            result[k] = n
        } else {
            errors = append(errors, 
                fmt.Errorf("%s: %w", k, err))
        }
    }
    
    return result, errors
}
`}
</CodeCard>

### Conditional Processing

<CodeCard file="conditional.go">
{`// Process only certain keys
func ProcessSelected(
    data map[string]string,
    keys []string,
) O.Option[map[string]int] {
    selected := make(map[string]string)
    for _, k := range keys {
        if v, ok := data[k]; ok {
            selected[k] = v
        }
    }
    
    return R.Traverse(O.Applicative[int]())(
        func(s string) O.Option[int] {
            if n, err := strconv.Atoi(s); err == nil {
                return O.Some(n)
            }
            return O.None[int]()
        },
    )(selected)
}
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Short-Circuit Behavior**: Traverse operations short-circuit on the first failure. With Option, the first None returns None. With Result, the first Error returns Error.

</Callout>

<Callout type="info">

**Use Cases**: Traverse is ideal for:
- Validating all map values
- Parsing configuration maps
- Batch API requests
- All-or-nothing transformations

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/record-monoid', title: 'Record Monoid' }}
  next={{ to: '/docs/v2/collections/record-eq', title: 'Record Equality' }}
/>

---
