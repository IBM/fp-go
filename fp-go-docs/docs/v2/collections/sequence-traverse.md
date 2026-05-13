---
title: Sequence & Traverse
hide_title: true
description: Working with arrays of effectful computations - converting between arrays of effects and effects of arrays.
sidebar_position: 7
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Sequence & Traverse"
  lede="Working with arrays of effectful computations. Convert between []Effect[A] and Effect[[]A] for powerful composition patterns."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Operations', value: 'Sequence, Traverse, TraverseWithIndex' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Sequence` | `func Sequence[A, F any](Applicative[F]) func([]HKT[F, A]) HKT[F, []A]` | Flip array and effect |
| `Traverse` | `func Traverse[A, B, F any](Applicative[F]) func(func(A) HKT[F, B]) func([]A) HKT[F, []B]` | Map and sequence |
| `TraverseWithIndex` | `func TraverseWithIndex[A, B, F any](Applicative[F]) func(func(int, A) HKT[F, B]) func([]A) HKT[F, []B]` | Traverse with index |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Sequence with Option

<CodeCard file="sequence_option.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

// Array of Options -> Option of Array
options := []O.Option[int]{
    O.Some(1),
    O.Some(2),
    O.Some(3),
}

result := A.Sequence(O.Applicative[int]())(options)
// Some([]int{1, 2, 3})

// With None - short circuits
withNone := []O.Option[int]{
    O.Some(1),
    O.None[int](),
    O.Some(3),
}

result2 := A.Sequence(O.Applicative[int]())(withNone)
// None - one None fails all
`}
</CodeCard>

### Sequence with Result

<CodeCard file="sequence_result.go">
{`import Res "github.com/IBM/fp-go/v2/result"

// Array of Results -> Result of Array
results := []Res.Result[int]{
    Res.Success(1),
    Res.Success(2),
    Res.Success(3),
}

combined := A.Sequence(Res.Applicative[error, int]())(results)
// Success([]int{1, 2, 3})

// With error
withError := []Res.Result[int]{
    Res.Success(1),
    Res.Error[int](errors.New("failed")),
    Res.Success(3),
}

combined2 := A.Sequence(Res.Applicative[error, int]())(withError)
// Error("failed")
`}
</CodeCard>

### Traverse - Map and Sequence

<CodeCard file="traverse.go">
{`type User struct {
    ID   int
    Name string
}

ids := []int{1, 2, 3}

// Fetch users (returns Option)
users := A.Traverse(O.Applicative[User]())(
    func(id int) O.Option[User] {
        return fetchUser(id)  // Returns Option[User]
    },
)(ids)
// Option[[]User] - Some if all found, None if any missing
`}
</CodeCard>

### TraverseWithIndex

<CodeCard file="traverse_index.go">
{`values := []string{"10", "20", "30"}

// Parse with index for error messages
parsed := A.TraverseWithIndex(Res.Applicative[error, int]())(
    func(i int, s string) Res.Result[int] {
        n, err := strconv.Atoi(s)
        if err != nil {
            return Res.Error[int](
                fmt.Errorf("index %d: %w", i, err),
            )
        }
        return Res.Success(n)
    },
)(values)
// Result[[]int]
`}
</CodeCard>

### Validating All Items

<CodeCard file="validate.go">
{`type Item struct {
    ID    int
    Value string
}

items := []Item{
    {ID: 1, Value: "valid"},
    {ID: 2, Value: "also-valid"},
    {ID: 3, Value: ""},  // Invalid
}

// Validate all items
validated := A.Traverse(Res.Applicative[error, Item]())(
    func(item Item) Res.Result[Item] {
        if item.Value == "" {
            return Res.Error[Item](
                fmt.Errorf("item %d: empty value", item.ID),
            )
        }
        return Res.Success(item)
    },
)(items)
// Error("item 3: empty value")
`}
</CodeCard>

### Parsing Configuration

<CodeCard file="parse_config.go">
{`type Config struct {
    Port    int
    Timeout int
    Retries int
}

raw := []string{"8080", "30", "3"}

// Parse all values
parsed := A.Traverse(Res.Applicative[error, int]())(
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
    Res.Map(func(values []int) Config {
        return Config{
            Port:    values[0],
            Timeout: values[1],
            Retries: values[2],
        }
    }),
)
// Result[Config]
`}
</CodeCard>

### Parallel API Calls

<CodeCard file="api.go">
{`import IOE "github.com/IBM/fp-go/v2/ioeither"

type UserData struct {
    ID   int
    Name string
}

userIDs := []int{1, 2, 3, 4, 5}

// Fetch all users
fetchAll := A.Traverse(IOE.Applicative[error, UserData]())(
    func(id int) IOE.IOEither[error, UserData] {
        return fetchUserAPI(id)
    },
)(userIDs)
// IOEither[error, []UserData]

// Execute
users := fetchAll()
// Either[error, []UserData]
`}
</CodeCard>

### File Operations

<CodeCard file="files.go">
{`import IO "github.com/IBM/fp-go/v2/io"

filenames := []string{"file1.txt", "file2.txt", "file3.txt"}

// Read all files
readAll := A.Traverse(IO.Applicative[[]byte]())(
    func(name string) IO.IO[[]byte] {
        return func() []byte {
            data, _ := os.ReadFile(name)
            return data
        }
    },
)(filenames)
// IO[[][]byte]

// Execute
contents := readAll()
// [][]byte - all file contents
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### All or Nothing Processing

<CodeCard file="all_or_nothing.go">
{`// Process all items - fail if any fails
func ProcessAll(items []string) Res.Result[[]int] {
    return A.Traverse(Res.Applicative[error, int]())(
        func(s string) Res.Result[int] {
            return processItem(s)
        },
    )(items)
}
`}
</CodeCard>

### Batch Operations

<CodeCard file="batch.go">
{`// Batch database inserts
func InsertAll(users []User) IOE.IOEither[error, []int] {
    return A.Traverse(IOE.Applicative[error, int]())(
        func(u User) IOE.IOEither[error, int] {
            return insertUser(u)
        },
    )(users)
}
`}
</CodeCard>

### Conditional Processing

<CodeCard file="conditional.go">
{`// Process only valid items
func ProcessValid(items []Item) O.Option[[]Result] {
    return A.Traverse(O.Applicative[Result]())(
        func(item Item) O.Option[Result] {
            if item.IsValid() {
                return O.Some(process(item))
            }
            return O.None[Result]()
        },
    )(items)
}
`}
</CodeCard>

### Collecting Results

<CodeCard file="collect.go">
{`// Collect successful results, skip failures
func CollectSuccesses(items []string) []int {
    results := A.Map(func(s string) O.Option[int] {
        if n, err := strconv.Atoi(s); err == nil {
            return O.Some(n)
        }
        return O.None[int]()
    })(items)
    
    // Filter out Nones
    return A.FilterMap(F.Identity[O.Option[int]])(results)
}
`}
</CodeCard>

### Parallel with Limit

<CodeCard file="parallel_limit.go">
{`// Process in batches to limit parallelism
func ProcessInBatches(items []Item, batchSize int) IOE.IOEither[error, []Result] {
    batches := chunkArray(items, batchSize)
    
    return A.Traverse(IOE.Applicative[error, []Result]())(
        func(batch []Item) IOE.IOEither[error, []Result] {
            return A.Traverse(IOE.Applicative[error, Result]())(
                processItem,
            )(batch)
        },
    )(batches)
}
`}
</CodeCard>

### Error Accumulation

<CodeCard file="errors.go">
{`// Collect all errors instead of short-circuiting
type ValidationResult struct {
    Valid  []Item
    Errors []error
}

func ValidateAll(items []Item) ValidationResult {
    var valid []Item
    var errors []error
    
    for _, item := range items {
        if err := validate(item); err != nil {
            errors = append(errors, err)
        } else {
            valid = append(valid, item)
        }
    }
    
    return ValidationResult{Valid: valid, Errors: errors}
}
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Short-Circuit Behavior**: Both Sequence and Traverse short-circuit on the first failure:
- With Option: first None returns None
- With Result: first Error returns Error
- With IOEither: first Left returns Left

</Callout>

<Callout type="info">

**Use Cases**: Traverse is ideal for:
- Validating all array elements
- Batch API requests
- Parallel file operations
- All-or-nothing transformations
- Converting imperative loops to functional pipelines

</Callout>

<Callout type="warn">

**Performance**: Traverse operations process elements sequentially by default. For true parallelism, use specialized parallel execution utilities or batch processing patterns.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/nonempty-array', title: 'NonEmpty Array' }}
  next={{ to: '/docs/v2/collections/record', title: 'Record (Map)' }}
/>

